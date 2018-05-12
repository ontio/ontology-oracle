package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-oracle/log"
	"github.com/ontio/ontology-oracle/models"
	"github.com/ontio/ontology-oracle/utils"
	"github.com/ontio/ontology/common"
	cstates "github.com/ontio/ontology/smartcontract/states"
	vmtypes "github.com/ontio/ontology/smartcontract/types"
	"github.com/ontio/ontology/smartcontract/service/native/oracle"
)

func (app *OracleApplication) InvokeOracleContract(
	address common.Address,
	operation string,
	args []byte) error {

	crt := &cstates.Contract{
		Address: address,
		Method:  operation,
		Args:    args,
	}
	buf := bytes.NewBuffer(nil)
	err := crt.Serialize(buf)
	if err != nil {
		return fmt.Errorf("Serialize contract error:%s", err)
	}
	tx := sdkcom.NewInvokeTransaction(0, 0, vmtypes.Native, buf.Bytes())

	err = sdkcom.SignTransaction(sdkcom.CRYPTO_SCHEME_DEFAULT, tx, app.Account)
	if err != nil {
		return fmt.Errorf("SignTransaction error:%s", err)
	}

	_, err = app.RPC.SendRawTransaction(tx)
	if err != nil {
		return fmt.Errorf("SendTransaction error:%s", err)
	}

	return err
}

type UndoRequests struct {
	Requests map[string]interface{} `json:"requests"`
}

type CreateOracleRequestParam struct {
	Request    string `json:"request"`
	OracleNode string `json:"oracleNode"`
	Address    string `json:"address"`
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
	if len(value) == 0 {
		return nil
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

		name := (events[1].States[0]).(string)
		if name != "createOracleRequest" {
			return nil
		}

		request := (events[1].States[1]).(map[string]interface{})

		address := hex.EncodeToString(app.Account.Address[:])
		if request["oracleNode"].(string) != address {
			fmt.Println("a")
			return nil
		}

		j := models.JobSpec{}
		err = json.Unmarshal([]byte(request["request"].(string)), &j)
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
	txHash := jr.JobID
	dataString := jr.Result.Data.Get("value").String()
	params := &oracle.SetOracleOutcomeParam{
		TxHash:  txHash,
		Address: hex.EncodeToString(app.Account.Address[:]),
		Outcome: dataString,
	}

	bf := new(bytes.Buffer)
	if err := params.Serialize(bf); err != nil {
		return fmt.Errorf("Serialize params error: %v", err)
	}
	err = app.InvokeOracleContract(
		address,
		operation,
		bf.Bytes())
	return err
}

//func (app *OracleApplication) sendCronDataToContract(jr models.JobRun) error {
//	address, err := utils.GetContractAddress()
//	if err != nil {
//		return fmt.Errorf("GetContractAddress error: %v", err)
//	}
//
//	operation := "setOracleCronOutcome"
//	txHash := jr.JobID
//	dataString := jr.Result.Data.Get("value").String()
//	params := &oracle.SetOracleCronOutcomeParam{
//		TxHash:  txHash,
//		Address: hex.EncodeToString(app.Account.Address[:]),
//		Outcome: dataString,
//	}
//
//	bf := new(bytes.Buffer)
//	if err := params.Serialize(bf); err != nil {
//		return fmt.Errorf("Serialize params error: %v", err)
//	}
//	err = app.InvokeOracleContract(
//		address,
//		operation,
//		bf.Bytes())
//	return err
//}
