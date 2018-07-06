package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/log"
	"github.com/ontio/ontology-oracle/models"
	"github.com/ontio/ontology-oracle/utils"
)

var Version = byte(0)

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

	value, err := app.RPC.GetStorage(contractAddress, []byte("UndoTxHash"))
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
	operation := "setOracleOutcome"
	txHash := jr.JobID
	dataString := jr.Result.Data.Get("value").String()

	args := []interface{}{operation, txHash, dataString}
	contractAddress, err := utils.GetContractAddress()
	if err != nil {
		return fmt.Errorf("utils.GetContractAddress error:%s", err)
	}
	_, err = app.RPC.InvokeNeoVMContract(config.Configuration.GasPrice, config.Configuration.GasLimit, app.Account,
		contractAddress, []interface{}{args})
	if err != nil {
		return fmt.Errorf("InvokeNeoVMContract error:%s", err)
	}
	return nil
}
