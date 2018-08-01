package core

import (
	"encoding/json"
	"fmt"

	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/log"
	"github.com/ontio/ontology-oracle/models"
	"github.com/ontio/ontology-oracle/utils"
	"github.com/ontio/ontology/smartcontract/service/neovm"
	"bytes"
	"github.com/ontio/ontology/common"
)

type UndoRequests struct {
	Requests map[string]interface{} `json:"requests"`
}

type CreateOracleRequestParam struct {
	Request    string `json:"request"`
	OracleNode string `json:"oracleNode"`
	Address    string `json:"address"`
}

func (app *OracleApplication) AddUndoRequests() error {
	contractAddress, err := utils.GetContractAddress()
	if err != nil {
		return fmt.Errorf("utils.GetContractAddress error:%s", err)
	}

	value, err := app.RPC.GetStorage(contractAddress, []byte("UndoRequest"))
	if err != nil {
		return fmt.Errorf("GetStorage UndoTxHash error:%s", err)
	}
	if len(value) == 0 {
		return nil
	}

	bf := bytes.NewBuffer(value)
	items, err := neovm.DeserializeStackItem(bf)
	if err != nil {
		return fmt.Errorf("neovm.DeserializeStackItem error:%s", err)
	}
	requestMap, err := items.GetMap()
	if err != nil {
		return fmt.Errorf("items.GetMap error:%s", err)
	}
	for k, v := range requestMap {
		txHashBytes , err := k.GetByteArray()
		if err != nil {
			return fmt.Errorf("k.GetByteArray error:%s", err)
		}
		requestBytes , err := v.GetByteArray()
		if err != nil {
			return fmt.Errorf("v.GetByteArray error:%s", err)
		}
		request := string(requestBytes)

		tx, err := common.Uint256ParseFromBytes(txHashBytes)
		if err != nil {
			return fmt.Errorf("common.Uint256ParseFromBytes error:%s", err)
		}

		j := models.JobSpec{}
		err = json.Unmarshal([]byte(request), &j)
		if err != nil {
			return fmt.Errorf("json.Unmarshal error:%s", err)
		}
		j.ID = tx.ToHexString()
		app.AddJob(&j)

		log.Debugf("Ontology Scanner get request txHash: %v", j.ID)
	}

	return nil
}

func (app *OracleApplication) sendDataToContract(jr models.JobRun) error {
	operation := "SetOracleOutcome"
	txHash, err := common.Uint256FromHexString(jr.JobID)
	if err != nil {
		return fmt.Errorf("common.AddressFromHexString error:%s", err)
	}
	dataString := jr.Result.Data.Get("value").String()

	args := []interface{}{operation, []interface{}{txHash[:], []byte(dataString)}}
	contractAddress, err := utils.GetContractAddress()
	if err != nil {
		return fmt.Errorf("utils.GetContractAddress error:%s", err)
	}
	_, err = app.RPC.InvokeNeoVMContract(config.Configuration.GasPrice, config.Configuration.GasLimit, app.Account,
		contractAddress, args)
	if err != nil {
		return fmt.Errorf("InvokeNeoVMContract error:%s", err)
	}
	return nil
}
