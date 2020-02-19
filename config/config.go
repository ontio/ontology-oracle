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
