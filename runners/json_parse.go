package runners

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"encoding/json"

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

	result, err := parseStruct(jsa, jsonParse.Data)
	if err != nil {
		return input.WithError(err)
	}
	b, err := neovm.SerializeStackItem(result)
	if err != nil {
		return input.WithError(err)
	}
	return input.WithValue(b)
}

func parseStruct(jsa *simplejson.Json, dataList []*OracleParamAbi) (types.StackItems, error) {
	temp1 := types.NewStruct(nil)
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
				result, err := parseStruct(jst, data.SubType)
				if err != nil {
					return nil, err
				}
				temp2.Add(result)
			}
			temp1.Add(temp2)
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
				result, err := parseStruct(jst, data.SubType)
				if err != nil {
					return nil, err
				}
				temp3.Add(types.NewByteArray([]byte(k)), result)
			}
			temp1.Add(temp3)
		default:
			switch strings.ToLower(data.Type) {
			case "string":
				r, err := getStringValue(js)
				if err != nil {
					return nil, err
				}
				ba := types.NewByteArray([]byte(r))
				temp1.Add(ba)
			case "int":
				r, err := getIntValue(js)
				if err != nil {
					return nil, err
				}
				int := types.NewInteger(new(big.Int).SetInt64(r))
				temp1.Add(int)
			case "float":
				r, err := getFloatValue(js)
				if err != nil {
					return nil, err
				}
				if data.Decimal == 0 {
					data.Decimal = 1
				}
				float := types.NewInteger(new(big.Int).SetInt64(int64(r * float64(data.Decimal))))
				temp1.Add(float)
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
