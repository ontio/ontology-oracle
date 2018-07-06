package utils

import (
	"encoding/hex"
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
	contractAddress, err := common.AddressFromBase58(config.Configuration.ContractAddress)
	if err != nil {
		return common.Address{}, fmt.Errorf("common.AddressFromBase58 error:%s", err)
	}
	return contractAddress, nil
}
