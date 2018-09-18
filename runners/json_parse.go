package runners

import (
	"fmt"
	"strconv"

	"github.com/bitly/go-simplejson"
	"github.com/ontio/ontology-oracle/models"
	"github.com/ontio/ontology/smartcontract/service/neovm"
	"github.com/ontio/ontology/vm/neovm/types"
	"math/big"
	"strings"
)

type JSONParse struct {
	Data []Data `json:"data"`
}

type Data struct {
	Type string   `json:"type"`
	Path []string `json:"path"`
}

func (jsonParse *JSONParse) Perform(input models.RunResult) models.RunResult {
	js, err := simplejson.NewJson(input.Data)
	if err != nil {
		return input.WithError(err)
	}

	stackArray := []types.StackItems{}
	for _, data := range jsonParse.Data {
		js, err = getByPath(js, data.Path)
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

func getByPath(js *simplejson.Json, path []string) (*simplejson.Json, error) {
	var ok bool
	for _, k := range path[:len(path)-1] {
		if isArray(js, k) {
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

func isArray(js *simplejson.Json, key string) bool {
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
