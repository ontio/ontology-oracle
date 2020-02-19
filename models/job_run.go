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

package models

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
	Data         []byte    `json:"data"`
	Status       RunStatus `json:"status"`
	ErrorMessage string    `json:"error"`
}

func (rr RunResult) WithError(err error) RunResult {
	rr.Data = nil
	rr.ErrorMessage = err.Error()
	rr.Status = RunStatusErrored
	return rr
}

func (rr RunResult) WithValue(data []byte) RunResult {
	rr.Status = RunStatusCompleted
	rr.Data = data
	return rr
}
