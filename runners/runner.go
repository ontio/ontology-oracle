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
