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
