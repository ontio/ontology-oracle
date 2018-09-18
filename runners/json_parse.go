package runners

import (
	"fmt"
	"strconv"
	"math/big"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/ontio/ontology-oracle/models"
	"github.com/ontio/ontology/smartcontract/service/neovm"
	"github.com/ontio/ontology/vm/neovm/types"
)

type JSONParse struct {
	Data []Data `json:"data"`
}

type Data struct {
	Type    string   `json:"type"`
	SubType string   `json:"sub_type"`
	Decimal uint64   `json:"decimal"`
	Path    []string `json:"path"`
}

func (jsonParse *JSONParse) Perform(input models.RunResult) models.RunResult {
	jsa, err := simplejson.NewJson(input.Data)
	if err != nil {
		return input.WithError(err)
	}

	stackArray := []types.StackItems{}
	for _, data := range jsonParse.Data {
		js, err := getByPath(jsa, data.Path)
		if err != nil {
			return input.WithError(err)
		}
		switch strings.ToLower(data.Type) {
		case "string":
			result, err := getStringValue(js)
			if err != nil {
				return input.WithError(err)
			}
			ba := types.NewByteArray([]byte(result))
			stackArray = append(stackArray, ba)
		case "int":
			result, err := getIntValue(js)
			if err != nil {
				return input.WithError(err)
			}
			int := types.NewInteger(new(big.Int).SetInt64(result))
			stackArray = append(stackArray, int)
		case "float":
			result, err := getFloatValue(js)
			if err != nil {
				return input.WithError(err)
			}
			float := types.NewInteger(new(big.Int).SetInt64(int64(result * float64(data.Decimal))))
			stackArray = append(stackArray, float)
		case "array":
			tempArray, err := js.Array()
			if err != nil {
				return input.WithError(err)
			}
			stackArrayTemp := []types.StackItems{}

			switch strings.ToLower(data.SubType) {
			case "string":
				for _, temp := range tempArray {
					result, ok := temp.(string)
					if !ok {
						return input.WithError(fmt.Errorf("array field is not string"))
					}
					ba := types.NewByteArray([]byte(result))
					stackArrayTemp = append(stackArray, ba)
				}
			case "int":
				for _, temp := range tempArray {
					result, ok := temp.(int64)
					if !ok {
						return input.WithError(fmt.Errorf("array field is not int"))
					}
					int := types.NewInteger(new(big.Int).SetInt64(result))
					stackArrayTemp = append(stackArray, int)
				}
			case "float":
				for _, temp := range tempArray {
					result, ok := temp.(float64)
					if !ok {
						return input.WithError(fmt.Errorf("array field is not float"))
					}
					float := types.NewInteger(new(big.Int).SetInt64(int64(result * float64(data.Decimal))))
					stackArrayTemp = append(stackArray, float)
				}
			default:
				return input.WithError(fmt.Errorf("data.SubType is not supported"))
			}
			array := types.NewArray(stackArrayTemp)
			stackArray = append(stackArray, array)
		case "map":
			tempMap, err := js.Map()
			if err != nil {
				return input.WithError(err)
			}
			mp := types.NewMap()
			switch strings.ToLower(data.SubType) {
			case "string":
				for k, v := range tempMap {
					result, ok := v.(string)
					if !ok {
						return input.WithError(fmt.Errorf("map field is not string"))
					}
					vStackItems := types.NewByteArray([]byte(result))
					mp.Add(types.NewByteArray([]byte(k)), vStackItems)
				}
			case "int":
				for k, v := range tempMap {
					result, ok := v.(int64)
					if !ok {
						return input.WithError(fmt.Errorf("map field is not int"))
					}
					vStackItems := types.NewInteger(new(big.Int).SetInt64(result))
					mp.Add(types.NewByteArray([]byte(k)), vStackItems)
				}
			case "float":
				for k, v := range tempMap {
					result, ok := v.(float64)
					if !ok {
						return input.WithError(fmt.Errorf("map field is not float"))
					}
					vStackItems := types.NewInteger(new(big.Int).SetInt64(int64(result * float64(data.Decimal))))
					mp.Add(types.NewByteArray([]byte(k)), vStackItems)
				}
			}
			stackArray = append(stackArray, mp)
		default:
			return input.WithError(fmt.Errorf("data.Type is not supported"))
		}
	}
	stru := types.NewStruct(stackArray)
	result, err := neovm.SerializeStackItem(stru)
	if err != nil {
		return input.WithError(err)
	}
	return input.WithValue(result)
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
