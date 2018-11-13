package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	DEFAULT_CONFIG_FILE_NAME = "./config.json"
)

// Config holds parameters used by the application which can be overridden
// by setting environment variables.
type Config struct {
	WalletFile      string `json:"WalletFile"`
	LogLevel        int    `json:"LogLevel"`
	MaxLogSize      int64  `json:"MaxLogSize"`
	ONTRPCAddress   string `json:"ONTRPCAddress"`
	ScannerInterval int    `json:"ScannerInterval"`
	GasPrice        uint64 `json:"GasPrice"`
	GasLimit        uint64 `json:"GasLimit"`
	ContractAddress string `json:"ContractAddress"`
	RandomOrgKey    string `json:"RandomOrgKey"`
}

var Configuration *Config

// NewConfig returns the config with the environment variables set to their
// respective fields, or defaults if not present.
func init() {
	file, e := ioutil.ReadFile(DEFAULT_CONFIG_FILE_NAME)
	if e != nil {
		fmt.Println("File error: ", e)
		os.Exit(1)
	}
	// Remove the UTF-8 Byte Order Mark
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

	config := Config{}
	e = json.Unmarshal(file, &config)
	if e != nil {
		fmt.Println("Unmarshal json file error: ", e)
		os.Exit(1)
	}

	Configuration = &config
}
