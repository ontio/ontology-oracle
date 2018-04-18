package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/payload"
	"github.com/ontio/ontology/core/signature"
	"github.com/ontio/ontology/core/types"
	cstates "github.com/ontio/ontology/smartcontract/states"
	vmtypes "github.com/ontio/ontology/smartcontract/types"
	"github.com/ontio/ontology/vm/neovm"
	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/log"
	"github.com/ontio/ontology-oracle/models"
	"github.com/ontio/ontology-oracle/utils"
)

func (app *OracleApplication) InvokeNeoVM(
	account *account.Account,
	address common.Address,
	params []interface{}) error {

	code, err := app.BuildNVMInvokeCode(address, params)
	if err != nil {
		return fmt.Errorf("BuildNVMInvokeCode error:%s", err)
	}
	tx := app.NewInvokeTransaction(0, vmtypes.NEOVM, code)

	return app.SendTransaction(account, tx)
}

func (app *OracleApplication) BuildNVMInvokeCode(addr common.Address, params []interface{}) ([]byte, error) {
	builder := neovm.NewParamsBuilder(new(bytes.Buffer))
	err := app.buildNVMParamInter(builder, params)
	if err != nil {
		return nil, err
	}
	args := builder.ToArray()

	crt := &cstates.Contract{
		Address: addr,
		Args:    args,
	}
	crtBuf := bytes.NewBuffer(nil)
	err = crt.Serialize(crtBuf)
	if err != nil {
		return nil, fmt.Errorf("Serialize contract error:%s", err)
	}

	buf := bytes.NewBuffer(nil)
	buf.Write(append([]byte{0x67}, crtBuf.Bytes()[:]...))
	return buf.Bytes(), nil
}

func (app *OracleApplication) SendTransaction(signer *account.Account, tx *types.Transaction) error {
	err := app.SignTransaction(tx, signer)
	if err != nil {
		return fmt.Errorf("SignTransaction error:%s", err)
	}

	return app.RPC.SendRawTransaction(tx)
}

func (app *OracleApplication) NewInvokeTransaction(gasLimit common.Fixed64, vmType vmtypes.VmType, code []byte) *types.Transaction {
	invokePayload := &payload.InvokeCode{
		GasLimit: gasLimit,
		Code: vmtypes.VmCode{
			VmType: vmType,
			Code:   code,
		},
	}
	tx := &types.Transaction{
		Version:    0,
		TxType:     types.Invoke,
		Nonce:      uint32(time.Now().Unix()),
		Payload:    invokePayload,
		Attributes: make([]*types.TxAttribute, 0, 0),
		Fee:        make([]*types.Fee, 0, 0),
		NetWorkFee: 0,
		Sigs:       make([]*types.Sig, 0, 0),
	}
	return tx
}

func (app *OracleApplication) buildNVMParamInter(builder *neovm.ParamsBuilder, smartContractParams []interface{}) error {
	//虚拟机参数入栈时会反序
	for i := len(smartContractParams) - 1; i >= 0; i-- {
		switch v := smartContractParams[i].(type) {
		case bool:
			builder.EmitPushBool(v)
		case int:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case uint:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case int32:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case uint32:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case int64:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case common.Fixed64:
			builder.EmitPushInteger(big.NewInt(int64(v.GetData())))
		case uint64:
			val := big.NewInt(0)
			builder.EmitPushInteger(val.SetUint64(uint64(v)))
		case string:
			builder.EmitPushByteArray([]byte(v))
		case *big.Int:
			builder.EmitPushInteger(v)
		case []byte:
			builder.EmitPushByteArray(v)
		case []interface{}:
			err := app.buildNVMParamInter(builder, v)
			if err != nil {
				return err
			}
			builder.EmitPushInteger(big.NewInt(int64(len(v))))
			builder.Emit(neovm.PACK)
		default:
			return fmt.Errorf("unsupported param:%s", v)
		}
	}
	return nil
}

func (app *OracleApplication) SignTransaction(tx *types.Transaction, signer *account.Account) error {
	txHash := tx.Hash()
	data, err := signature.Sign(signer, txHash.ToArray())
	if err != nil {
		return fmt.Errorf("crypto sign error:%s", err)
	}

	sig := &types.Sig{
		PubKeys: []keypair.PublicKey{signer.PublicKey},
		M:       1,
		SigData: [][]byte{data},
	}
	tx.Sigs = []*types.Sig{sig}
	return nil
}

func (app *OracleApplication) ParseResp(resp map[string]interface{}) error {
	action := resp["Action"]
	errs := resp["Error"]
	if action.(string) != "InvokeTransaction" {
		return nil
	}
	if errs.(float64) != 0 {
		return errors.New("Recieve failed InvokeTransaction")
	}

	result := resp["Result"].([]interface{})
	codeHash := result[0].(map[string]interface{})["CodeHash"]

	if codeHash != config.Configuration.CodeHash {
		return nil
	}
	states := result[0].(map[string]interface{})["States"]
	name, _ := utils.ConvertToString(states.([]interface{})[0].([]interface{})[0])
	if name == "setOracleOutcome" {
		jobId := states.([]interface{})[0].([]interface{})[2]
		status, _ := utils.ConvertToString(states.([]interface{})[0].([]interface{})[1])
		if status == "success" {
			log.Infof("SetOutcome of Job ID %v is successfully committed!", jobId)
		} else {
			log.Errorf("SetOutcome of Job ID %v is failed!", jobId)
		}
		return nil
	}
	if name != "createOracleRequest" {
		return nil
	}

	request, _ := utils.ConvertToString(states.([]interface{})[0].([]interface{})[1])
	j := models.JobSpec{}
	err := json.Unmarshal([]byte(request), &j)
	if err != nil {
		return err
	}
	j.ID = result[0].(map[string]interface{})["TxHash"].(string)
	app.AddJob(&j)

	log.Infof("Ontology listener get request txHash: %v", j.ID)

	return nil
}

type UndoRequests struct {
	Requests map[string]interface{} `json:"requests"`
}

func (app *OracleApplication) AddUndoRequests() error {
	address, err := utils.GetContractAddress()
	if err != nil {
		return fmt.Errorf("GetContractAddress error: %v", err)
	}

	value, err := app.RPC.GetStorage(address, []byte("UndoTxHash"))
	if err != nil {
		return fmt.Errorf("GetStorage UndoTxHash error:%s", err)
	}

	undoRequests := &UndoRequests{
		Requests: make(map[string]interface{}),
	}
	err = json.Unmarshal(value, &undoRequests)
	if err != nil {
		return fmt.Errorf("Unmarshal UndoRequests: %s", err)
	}

	for txHash := range undoRequests.Requests {
		tx, err := utils.ParseUint256FromHexString(txHash)
		events, err := app.RPC.GetSmartContractEvent(tx)
		if err != nil {
			return fmt.Errorf("GetSmartContractEvent error:%s", err)
		}

		name := (events[0].States[0]).(string)
		if name != "createOracleRequest" {
			return nil
		}

		request := (events[0].States[1]).(string)

		j := models.JobSpec{}
		err = json.Unmarshal([]byte(request), &j)
		if err != nil {
			return fmt.Errorf("json.Unmarshal error:%s", err)
		}
		j.ID = txHash
		app.AddJob(&j)

		log.Debugf("Ontology Scanner get request txHash: %v", j.ID)
	}
	return nil
}

func (app *OracleApplication) sendDataToContract(jr models.JobRun) error {
	address, err := utils.GetContractAddress()
	if err != nil {
		return fmt.Errorf("GetContractAddress error: %v", err)
	}

	operation := "setOracleOutcome"
	txHash, _ := hex.DecodeString(jr.JobID)
	dataString := jr.Result.Data.Get("value").String()
	args := []interface{}{txHash, dataString}
	err = app.InvokeNeoVM(
		app.Account,
		address,
		[]interface{}{operation, args})
	return err
}

func (app *OracleApplication) sendCronDataToContract(jr models.JobRun) error {
	address, err := utils.GetContractAddress()
	if err != nil {
		return fmt.Errorf("GetContractAddress error: %v", err)
	}

	operation := "setOracleCronOutcome"
	txHash, _ := hex.DecodeString(jr.JobID)
	dataString := jr.Result.Data.Get("value").String()
	args := []interface{}{txHash, dataString}
	err = app.InvokeNeoVM(
		app.Account,
		address,
		[]interface{}{operation, args})
	return err
}
