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
	"math/big"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/ontio/ontology-oracle/models"
	"github.com/ontio/ontology/smartcontract/service/neovm"
	"github.com/ontio/ontology/vm/neovm/types"
)

type JSONParse struct {
	Data []*OracleParamAbi `json:"data"`
}

type OracleParamAbi struct {
	Type    string            `json:"type"`
	Path    []string          `json:"path"`
	Decimal uint64            `json:"decimal"`
	SubType []*OracleParamAbi `json:"sub_type"`
}

func (jsonParse *JSONParse) Perform(input models.RunResult) models.RunResult {
	jsa, err := simplejson.NewJson(input.Data)
	if err != nil {
		return input.WithError(err)
	}

	results, err := parseStruct(jsa, jsonParse.Data)
	if err != nil {
		return input.WithError(err)
	}
	stakeItem := types.NewStruct(results)
	b, err := neovm.SerializeStackItem(stakeItem)
	if err != nil {
		return input.WithError(err)
	}
	return input.WithValue(b)
}

func parseStruct(jsa *simplejson.Json, dataList []*OracleParamAbi) ([]types.StackItems, error) {
	temp1 := []types.StackItems{}
	for _, data := range dataList {
		js, err := getByPath(jsa, data.Path)
		if err != nil {
			return nil, err
		}
		switch strings.ToLower(data.Type) {
		case "array":
			temp2 := types.NewArray(nil)
			tempArray, err := js.Array()
			if err != nil {
				return nil, err
			}
			for _, item := range tempArray {
				b, err := json.Marshal(item)
				if err != nil {
					return nil, err
				}
				jst, err := simplejson.NewJson(b)
				if err != nil {
					return nil, err
				}
				results, err := parseStruct(jst, data.SubType)
				if err != nil {
					return nil, err
				}
				for _, result := range results {
					temp2.Add(result)
				}
			}
			temp1 = append(temp1, temp2)
		case "map":
			temp3 := types.NewMap()
			tempMap, err := js.Map()
			if err != nil {
				return nil, err
			}
			for k, v := range tempMap {
				b, err := json.Marshal(v)
				if err != nil {
					return nil, err
				}
				jst, err := simplejson.NewJson(b)
				if err != nil {
					return nil, err
				}
				results, err := parseStruct(jst, data.SubType)
				if err != nil {
					return nil, err
				}
				for _, result := range results {
					temp3.Add(types.NewByteArray([]byte(k)), result)
				}
			}
			temp1 = append(temp1, temp3)
		case "struct":
			temp4 := types.NewStruct(nil)
			results, err := parseStruct(js, data.SubType)
			if err != nil {
				return nil, err
			}
			for _, result := range results {
				temp4.Add(result)
			}
			temp1 = append(temp1, temp4)
		default:
			switch strings.ToLower(data.Type) {
			case "string":
				r, err := getStringValue(js)
				if err != nil {
					return nil, err
				}
				ba := types.NewByteArray([]byte(r))
				temp1 = append(temp1, ba)
			case "int":
				r, err := getIntValue(js)
				if err != nil {
					return nil, err
				}
				int := types.NewInteger(new(big.Int).SetInt64(r))
				temp1 = append(temp1, int)
			case "float":
				r, err := getFloatValue(js)
				if err != nil {
					return nil, err
				}
				if data.Decimal == 0 {
					data.Decimal = 1
				}
				float := types.NewInteger(new(big.Int).SetInt64(int64(r * float64(data.Decimal))))
				temp1 = append(temp1, float)
			default:
				return nil, fmt.Errorf("data.Type is not supported")
			}
		}
	}
	return temp1, nil
}

func getStringValue(js *simplejson.Json) (string, error) {
	str, err := js.String()
	if err != nil {
		return str, err
	}
	return str, nil
}

func getIntValue(js *simplejson.Json) (int64, error) {
	int64, err := js.Int64()
	if err != nil {
		return int64, err
	}
	return int64, nil
}

func getFloatValue(js *simplejson.Json) (float64, error) {
	float64, err := js.Float64()
	if err != nil {
		return float64, err
	}
	return float64, nil
}

func getByPath(js *simplejson.Json, path []string) (*simplejson.Json, error) {
	var ok bool
	for _, k := range path {
		if isArray(js) {
			js, ok = arrayGet(js, k)
		} else {
			js, ok = js.CheckGet(k)
		}
		if !ok {
			return js, fmt.Errorf("No value could be found for the key '" + k + "'")
		}
	}
	return js, nil
}

func isArray(js *simplejson.Json) bool {
	if _, err := js.Array(); err != nil {
		return false
	}
	return true
}

func arrayGet(js *simplejson.Json, key string) (*simplejson.Json, bool) {
	i, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return js, false
	}
	a, err := js.Array()
	if err != nil || len(a) < int(i-1) {
		return js, false
	}
	return js.GetIndex(int(i)), true
}
