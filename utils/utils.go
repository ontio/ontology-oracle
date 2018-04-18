package utils

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology-oracle/config"
	"github.com/satori/go.uuid"
)

// NewBytes32ID returns a randomly generated UUID
func NewBytes32ID() string {
	return strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
}

// BasicPost sends a POST request to the HTTP client with contentType and returns a response.
func BasicPost(url string, contentType string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", contentType)
	resp, err := client.Do(request)
	return resp, err
}

func ConvertToString(v interface{}) (string, error) {
	value, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%v ConvertToString failed", v)
	}
	data, _ := hex.DecodeString(value)
	return string(data), nil
}

//ParseUint256FromHexString return Uint256 parse from hex string
func ParseUint256FromHexString(value string) (common.Uint256, error) {
	data, err := hex.DecodeString(value)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("hex.DecodeString error:%s", err)
	}
	res, err := common.Uint256ParseFromBytes(data)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("Uint160ParseFromBytes error:%s", err)
	}
	return res, nil
}

func GetContractAddress() (common.Address, error) {
	addressBytes, err := hex.DecodeString(config.Configuration.CodeHash)
	if err != nil {
		return common.Address{}, fmt.Errorf("Decode codeHash config error: %v", err)
	}
	address, err := common.AddressParseFromBytes(addressBytes)
	if err != nil {
		return common.Address{}, fmt.Errorf("Decode codeHash config error: %v", err)
	}
	return address, nil
}
