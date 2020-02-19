/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package core

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/log"
	"github.com/ontio/ontology-oracle/models"
	"github.com/ontio/ontology-oracle/utils"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/neovm"
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

	value, err := app.Ont.GetStorage(contractAddress.ToHexString(), []byte("UndoRequest"))
	if err != nil {
		return fmt.Errorf("GetStorage UndoTxHash error:%s", err)
	}
	if len(value) <= 2 {
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
		txHashBytes, err := k.GetByteArray()
		if err != nil {
			log.Errorf("k.GetByteArray error:%s", err)
		}
		requestBytes, err := v.GetByteArray()
		if err != nil {
			log.Errorf("v.GetByteArray error:%s", err)
		}
		request := string(requestBytes)

		tx, err := common.Uint256ParseFromBytes(txHashBytes)
		if err != nil {
			log.Errorf("common.Uint256ParseFromBytes error:%s", err)
		}

		j := models.JobSpec{}
		j.ID = tx.ToHexString()
		j.Request = requestBytes
		err = json.Unmarshal([]byte(request), &j)
		if err != nil {
			log.Errorf("json.Unmarshal error:%s", err)
			rr := models.RunResult{
				JobRunID:     j.ID,
				Status:       models.RunStatusErrored,
				ErrorMessage: err.Error(),
			}
			jr := models.JobRun{
				JobID: j.ID,
			}
			jr = jr.ApplyResult(rr)
			err := app.sendDataToContract(jr)
			if err != nil {
				log.Errorf("send error data to contract error: %v", err.Error())
			} else {
				log.Infof("send error data to contract success, Job ID is: %v", jr.JobID)
			}
			err = app.Store.Put([]byte(jr.JobID), requestBytes, nil)
			if err != nil {
				log.Errorf("put request into db error : %v", err)
			}
			continue
		}
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

	args := []interface{}{operation, []interface{}{txHash[:], jr.Result.Data, string(jr.Result.Status), jr.Result.ErrorMessage}}
	contractAddress, err := utils.GetContractAddress()
	if err != nil {
		return fmt.Errorf("utils.GetContractAddress error:%s", err)
	}
	_, err = app.Ont.NeoVM.InvokeNeoVMContract(config.Configuration.GasPrice, config.Configuration.GasLimit, app.Account,
		contractAddress, args)
	if err != nil {
		return fmt.Errorf("InvokeNeoVMContract error:%s", err)
	}
	return nil
}
