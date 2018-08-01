package core

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ontio/ontology-go-sdk/rpc"
	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/log"
	"github.com/ontio/ontology-oracle/models"
	"github.com/ontio/ontology/account"
)

// Application implements the common functions used in the core node.
type Application interface {
	Start()
	Stop()
}

type OracleApplication struct {
	Account   *account.Account
	JobList   chan *models.JobSpec
	DoingJobs map[string]interface{}
	RPC       *rpc.RpcClient
	Exiter    func(int)
}

func NewApplication(acct *account.Account) Application {
	jobList := make(chan *models.JobSpec, 10)
	return &OracleApplication{
		Account:   acct,
		JobList:   jobList,
		DoingJobs: make(map[string]interface{}),
		Exiter:    os.Exit,
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
		_, ok := app.DoingJobs[job.ID]
		if ok {
			log.Debugf("job %v is already on process", job.ID)
			continue
		}
		switch strings.ToLower(job.Scheduler.Type) {
		default:
			jobRun := job.NewRun()
			go app.ExecuteRun(jobRun)
		}
	}
}

func (app *OracleApplication) OntScanner() {
	log.Info("Start getting undo request in oracle contract.")
	app.RPC = rpc.NewRpcClient()
	app.RPC.SetAddress(config.Configuration.ONTRPCAdress)
	err := app.AddUndoRequests()
	if err != nil {
		log.Errorf("OntScanner error: %v", err)
	}

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
