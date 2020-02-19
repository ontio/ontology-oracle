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

package core

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/log"
	"github.com/ontio/ontology-oracle/models"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const DB_PATH = "./Store"

// Application implements the common functions used in the core node.
type Application interface {
	Start()
	Stop()
}

type OracleApplication struct {
	Account *sdk.Account
	JobList chan *models.JobSpec
	Store   *leveldb.DB
	Ont     *sdk.OntologySdk
	Exiter  func(int)
}

func NewApplication(acct *sdk.Account) Application {
	jobList := make(chan *models.JobSpec, 10)
	ontSdk := sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(config.Configuration.ONTRPCAddress)

	//init leveldb store
	// default Options
	o := opt.Options{
		NoSync: false,
		Filter: filter.NewBloomFilter(10),
	}
	db, err := leveldb.OpenFile(DB_PATH, &o)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(DB_PATH, nil)
	}
	if err != nil {
		log.Fatalf("Can't init leveldb: %s", err)
		os.Exit(1)
	}

	return &OracleApplication{
		Account: acct,
		JobList: jobList,
		Store:   db,
		Ont:     ontSdk,
		Exiter:  os.Exit,
	}
}

func (app *OracleApplication) Start() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		app.Stop()
		app.Exiter(1)
	}()

	go app.JobRunner()
	go app.OntScanner()
}

func (app *OracleApplication) Stop() {
	defer log.ClosePrintLog()
	log.Info("Gracefully Stop Oracle Node...")
}

func (app *OracleApplication) JobRunner() {
	for {
		job := <-app.JobList
		_, err := app.Store.Get([]byte(job.ID), nil)
		if err != leveldb.ErrNotFound {
			log.Debugf("job %v is already on process", job.ID)
			continue
		}

		switch strings.ToLower(job.Scheduler.Type) {
		default:
			go app.RunJob(job)
		}
	}
}

func (app *OracleApplication) OntScanner() {
	log.Info("Start getting undo request in oracle contract.")

	timer := time.NewTimer(time.Duration(config.Configuration.ScannerInterval) * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			err := app.AddUndoRequests()
			if err != nil {
				log.Errorf("OntScanner error: %v", err)
			}
			timer.Reset(time.Duration(config.Configuration.ScannerInterval) * time.Second)
		}
	}
}

func (app *OracleApplication) AddJob(job *models.JobSpec) {
	app.JobList <- job
}
