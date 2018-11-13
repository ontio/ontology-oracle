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
	if err != nil {
		return input.WithError(err)
	}

	if response.StatusCode >= 400 {
		return input.WithError(fmt.Errorf(string(bytes)))
	}

	return input.WithValue(bytes)
}

type HTTPPost struct {
	URL         models.WebURL `json:"url"`
	ContentType string        `json:"contentType"`
	Body        string        `json:"body"`
}

func (httpPost *HTTPPost) Perform(input models.RunResult) models.RunResult {
	response, err := http.Post(httpPost.URL.String(), httpPost.ContentType, bytes.NewReader([]byte(httpPost.Body)))
	if err != nil {
		return input.WithError(err)
	}

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return input.WithError(err)
	}

	if response.StatusCode >= 400 {
		return input.WithError(fmt.Errorf(string(bytes)))
	}

	return input.WithValue(bytes)
}
