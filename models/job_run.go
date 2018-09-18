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
