package models

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"net/url"
)

type RunStatus string

const (

	// RunStatusInProgress is used for when a run is actively being executed.
	RunStatusInProgress = RunStatus("in progress")
	// RunStatusErrored is used for when a run has errored and will not complete.
	RunStatusErrored = RunStatus("errored")
	// RunStatusCompleted is used for when a run has successfully completed execution.
	RunStatusCompleted = RunStatus("completed")
)

// Completed returns true if the status is RunStatusCompleted.
func (s RunStatus) Completed() bool {
	return s == RunStatusCompleted
}

// Errored returns true if the status is RunStatusErrored.
func (s RunStatus) Errored() bool {
	return s == RunStatusErrored
}

// Runnable returns true if the status is ready to be run.
func (s RunStatus) Runnable() bool {
	return !s.Errored()
}

// JSON stores the json types string, number, bool, and null.
// Arrays and Objects are returned as their raw json types.
type JSON struct {
	gjson.Result
}

// ParseJSON attempts to coerce the input byte array into valid JSON
// and parse it into a JSON object.
func ParseJSON(b []byte) (JSON, error) {
	var j JSON
	str := string(b)
	if len(str) == 0 {
		str = `{}`
	}
	return j, json.Unmarshal([]byte(str), &j)
}

// UnmarshalJSON parses the JSON bytes and stores in the *JSON pointer.
func (j *JSON) UnmarshalJSON(b []byte) error {
	str := string(b)
	if !gjson.Valid(str) {
		return fmt.Errorf("invalid JSON: %v", str)
	}
	*j = JSON{gjson.Parse(str)}
	return nil
}

// MarshalJSON returns the JSON data if it already exists, returns
// an empty JSON object as bytes if not.
func (j JSON) MarshalJSON() ([]byte, error) {
	if j.Exists() {
		return j.Bytes(), nil
	}
	return []byte("{}"), nil
}

// Merge combines the given JSON with the existing JSON.
func (j JSON) Merge(j2 JSON) (JSON, error) {
	body := j.Map()
	for key, value := range j2.Map() {
		body[key] = value
	}

	cleaned := map[string]interface{}{}
	for k, v := range body {
		cleaned[k] = v.Value()
	}

	b, err := json.Marshal(cleaned)
	if err != nil {
		return JSON{}, err
	}

	var rval JSON
	return rval, gjson.Unmarshal(b, &rval)
}

// Empty returns true if the JSON does not exist.
func (j JSON) Empty() bool {
	return !j.Exists()
}

// Bytes returns the raw JSON.
func (j JSON) Bytes() []byte {
	return []byte(j.String())
}

// Add returns a new instance of JSON with the new value added.
func (j JSON) Add(key string, val interface{}) (JSON, error) {
	var j2 JSON
	b, err := json.Marshal(val)
	if err != nil {
		return j2, err
	}
	str := fmt.Sprintf(`{"%v":%v}`, key, string(b))
	if err = json.Unmarshal([]byte(str), &j2); err != nil {
		return j2, err
	}
	return j.Merge(j2)
}

// WebURL contains the URL of the endpoint.
type WebURL struct {
	*url.URL
}

// UnmarshalJSON parses the raw URL stored in JSON-encoded
// data to a URL structure and sets it to the URL field.
func (w *WebURL) UnmarshalJSON(j []byte) error {
	var v string
	err := json.Unmarshal(j, &v)
	if err != nil {
		return err
	}
	u, err := url.ParseRequestURI(v)
	if err != nil {
		return err
	}
	w.URL = u
	return nil
}

// MarshalJSON returns the JSON-encoded string of the given data.
func (w *WebURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

// String delegates to the wrapped URL struct or an empty string when it is nil
func (w *WebURL) String() string {
	if w.URL == nil {
		return ""
	}
	return w.URL.String()
}
