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
