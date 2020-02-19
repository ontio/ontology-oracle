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

package utils

import (
	"fmt"
	"strings"

	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology/common"
	"github.com/satori/go.uuid"
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
