package runners

import (
	"encoding/json"
	"fmt"

	"github.com/AkshatM/caprice"
	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/models"
)

type RandomOrg struct {
	Method      string `json:"method"`
	N           int    `json:"n"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Replacement bool   `json:"replacement"`
}

type SignedIntegerData struct {
	Raw          json.RawMessage `json:"raw"`
	HashedApiKey string          `json:"hashedApiKey"`
	SerialNumber int             `json:"serialNumber"`
	Data         []int           `json:"data"`
	Signature    string          `json:"signature"`
}

func (randomOrg *RandomOrg) Perform(input models.RunResult) models.RunResult {
	rng := caprice.TrueRNG(config.Configuration.RandomOrgKey)

	var result caprice.Result
	var signedResult caprice.SignedIntegerData
	var err1 caprice.Error
	var err error
	var bytes []byte
	switch randomOrg.Method {
	case "GenerateIntegers":
		result, err1 = rng.GenerateIntegersRaw(randomOrg.N, randomOrg.Min, randomOrg.Max, randomOrg.Replacement)
		if err1.Message != "" {
			return input.WithError(fmt.Errorf(err1.Message))
		}
		bytes, err = json.Marshal(result)
		if err != nil {
			return input.WithError(err)
		}
	case "GenerateSignedIntegers":
		signedResult, err1 = rng.GenerateSignedIntegers(randomOrg.N, randomOrg.Min, randomOrg.Max, randomOrg.Replacement)
		if err1.Message != "" {
			return input.WithError(fmt.Errorf(err1.Message))
		}
		data := SignedIntegerData{
			Raw:          signedResult.Raw,
			HashedApiKey: signedResult.HashedApiKey,
			SerialNumber: signedResult.SerialNumber,
			Data:         signedResult.Data,
			Signature:    signedResult.Signature,
		}
		bytes, err = json.Marshal(data)
		if err != nil {
			return input.WithError(err)
		}
	default:
		return input.WithError(fmt.Errorf("randomOrg method is not supported"))
	}

	return input.WithValue(bytes)
}
