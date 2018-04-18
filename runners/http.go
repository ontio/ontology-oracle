package runners

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ontio/ontology-oracle/models"
)

type HTTPGet struct {
	URL models.WebURL `json:"url"`
}

func (httpGet *HTTPGet) Perform(input models.RunResult) models.RunResult {
	response, err := http.Get(httpGet.URL.String())
	if err != nil {
		return input.WithError(err)
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil {
		return input.WithError(err)
	}

	if response.StatusCode >= 400 {
		return input.WithError(fmt.Errorf(body))
	}

	return input.WithValue(body)
}

type HTTPPost struct {
	URL models.WebURL `json:"url"`
}

func (httpPost *HTTPPost) Perform(input models.RunResult) models.RunResult {
	reqBody := bytes.NewBufferString(input.Data.String())
	response, err := http.Post(httpPost.URL.String(), "application/json", reqBody)
	if err != nil {
		return input.WithError(err)
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil {
		return input.WithError(err)
	}

	if response.StatusCode >= 400 {
		return input.WithError(fmt.Errorf(body))
	}

	return input.WithValue(body)
}
