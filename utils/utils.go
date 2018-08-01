package utils

import (
	"fmt"
	"strings"

	"github.com/ontio/ontology/common"
	"github.com/satori/go.uuid"
	"github.com/ontio/ontology-oracle/config"
)

// NewBytes32ID returns a randomly generated UUID
func NewBytes32ID() string {
	return strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
}

func GetContractAddress() (common.Address, error) {
	contractAddress, err := common.AddressFromHexString(config.Configuration.ContractAddress)
	if err != nil {
		return common.Address{}, fmt.Errorf("common.AddressFromHexString error:%s", err)
	}
	return contractAddress, nil
}
