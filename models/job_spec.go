package models

import (
	"github.com/ontio/ontology-oracle/utils"
)

type JobSpec struct {
	ID        string        `json:"id"`
	Scheduler SchedulerSpec `json:"scheduler"`
	Tasks     []TaskSpec    `json:"tasks"`
	Request   []byte        `json:"request"`
}

func (j JobSpec) NewRun() JobRun {
	jrid := utils.NewBytes32ID()
	taskRuns := make([]TaskRun, len(j.Tasks))
	for i, task := range j.Tasks {
		taskRuns[i] = TaskRun{
			ID:     utils.NewBytes32ID(),
			Task:   task,
			Result: RunResult{JobRunID: jrid},
		}
	}

	return JobRun{
		ID:        jrid,
		JobID:     j.ID,
		Scheduler: j.Scheduler,
		TaskRuns:  taskRuns,
	}
}

type TaskSpec struct {
	Type   string `json:"type"`
	Params JSON
}

type SchedulerSpec struct {
	Type   string `json:"type"`
	Params string `json:"params"`
}
