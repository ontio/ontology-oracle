package models

import (
	"fmt"
	"github.com/tidwall/gjson"
)

type JobRun struct {
	ID        string        `json:"id"`
	JobID     string        `json:"jobId"`
	Scheduler SchedulerSpec `json:"scheduler"`
	Result    RunResult     `json:"result"`
	Status    RunStatus     `json:"status"`
	TaskRuns  []TaskRun     `json:"taskRuns"`
}

func (jr JobRun) ApplyResult(result RunResult) JobRun {
	jr.Result = result
	jr.Status = result.Status

	return jr
}

type TaskRun struct {
	ID     string    `json:"id"`
	Result RunResult `json:"result"`
	Status RunStatus `json:"status"`
	Task   TaskSpec  `json:"task"`
}

func (tr TaskRun) MarkCompleted() TaskRun {
	tr.Status = RunStatusCompleted
	tr.Result.Status = RunStatusCompleted
	return tr
}

func (tr TaskRun) ApplyResult(result RunResult) TaskRun {
	tr.Result = result
	tr.Status = result.Status
	return tr
}

type RunResult struct {
	JobRunID     string    `json:"jobRunId"`
	Data         JSON      `json:"data"`
	Status       RunStatus `json:"status"`
	ErrorMessage string    `json:"error"`
}

func (rr RunResult) Value() (string, error) {
	val, err := rr.value()
	if err != nil {
		return "", err
	}
	if val.Type != gjson.String {
		return "", fmt.Errorf("non string value")
	}
	return val.String(), nil
}

func (rr RunResult) Get(path string) (gjson.Result, error) {
	return rr.Data.Get(path), nil
}

func (rr RunResult) value() (gjson.Result, error) {
	return rr.Get("value")
}

func (rr RunResult) WithError(err error) RunResult {
	rr.ErrorMessage = err.Error()
	rr.Status = RunStatusErrored
	return rr
}

func (rr RunResult) WithValue(val string) RunResult {
	data, err := rr.Data.Add("value", val)
	if err != nil {
		return rr.WithError(err)
	}
	rr.Status = RunStatusCompleted
	rr.Data = data
	return rr
}
