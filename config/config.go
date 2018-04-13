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
	LogLevel        int    `json:"LogLevel"`
	MaxLogSize      int64  `json:"MaxLogSize"`
	Port            string `json:"Port"`
	ONTWSAddress    string `json:"ONTWSAddress"`
	ONTRPCAdress    string `json:"ONTRPCAdress"`
	ScannerInterval int    `json:"ScannerInterval"`
	CodeHash        string `json:"CodeHash"`
}

var Configuration *Config

// NewConfig returns the config with the environment variables set to their
// respective fields, or defaults if not present.
func init() {
	file, e := ioutil.ReadFile(DEFAULT_CONFIG_FILE_NAME)
	if e != nil {
		fmt.Errorf("File error: %v\n", e)
		os.Exit(1)
	}
	// Remove the UTF-8 Byte Order Mark
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

	config := Config{}
	e = json.Unmarshal(file, &config)
	if e != nil {
		fmt.Errorf("Unmarshal json file error %v\n", e)
		os.Exit(1)
	}

	Configuration = &config
}
