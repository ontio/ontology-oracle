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

package runners

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ontio/ontology-oracle/models"
)

// The Runner interface applies to all core runners.
// Each implementation must return a RunResult.
type Runner interface {
	Perform(models.RunResult) models.RunResult
}

// For determines the runner type to use for a given task
func For(task models.TaskSpec) (r Runner, err error) {
	switch strings.ToLower(task.Type) {
	case "randomorg":
		r = &RandomOrg{}
		err = unmarshalParams(task.Params, r)
	case "httpget":
		r = &HTTPGet{}
		err = unmarshalParams(task.Params, r)
	case "httppost":
		r = &HTTPPost{}
		err = unmarshalParams(task.Params, r)
	case "jsonparse":
		r = &JSONParse{}
		err = unmarshalParams(task.Params, r)
	default:
		return nil, fmt.Errorf("%s is not a supported runner type", task.Type)
	}
	return r, err
}

func unmarshalParams(params models.JSON, dst interface{}) error {
	bytes, err := params.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
