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

package runners

import (
	"encoding/json"
	"fmt"

	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/models"
	"github.com/siovanus/caprice"
)

type RandomOrg struct {
	Method      string `json:"method"`
	N           int    `json:"n"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Replacement bool   `json:"replacement"`
}

type IntegerData struct {
	Data           []interface{} `json:"data"`
	CompletionTime string        `json:"completionTime"`
}

type SignedIntegerData struct {
	CompletionTime string `json:"completionTime"`
	HashedApiKey   string `json:"hashedApiKey"`
	SerialNumber   int    `json:"serialNumber"`
	Data           []int  `json:"data"`
	Signature      string `json:"signature"`
}

func (randomOrg *RandomOrg) Perform(input models.RunResult) models.RunResult {
	rng := caprice.TrueRNG(config.Configuration.RandomOrgKey)

	var result caprice.Result
	var signedResult caprice.SignedIntegerData
	var err1 caprice.Error
	var err error
	var bytes []byte
	switch randomOrg.Method {
	case "GenerateIntegers":
		result, err1 = rng.GenerateIntegersRaw(randomOrg.N, randomOrg.Min, randomOrg.Max, randomOrg.Replacement)
		if err1.Message != "" {
			return input.WithError(fmt.Errorf(err1.Message))
		}
		data := IntegerData{
			CompletionTime: result.Random.CompletionTime,
			Data:           result.Random.Data,
		}
		bytes, err = json.Marshal(data)
		if err != nil {
			return input.WithError(err)
		}
	case "GenerateSignedIntegers":
		signedResult, err1 = rng.GenerateSignedIntegers(randomOrg.N, randomOrg.Min, randomOrg.Max, randomOrg.Replacement)
		if err1.Message != "" {
			return input.WithError(fmt.Errorf(err1.Message))
		}
		data := SignedIntegerData{
			CompletionTime: signedResult.CompletionTime,
			HashedApiKey:   signedResult.HashedApiKey,
			SerialNumber:   signedResult.SerialNumber,
			Data:           signedResult.Data,
			Signature:      signedResult.Signature,
		}
		bytes, err = json.Marshal(data)
		if err != nil {
			return input.WithError(err)
		}
	default:
		return input.WithError(fmt.Errorf("randomOrg method is not supported"))
	}

	return input.WithValue(bytes)
}
