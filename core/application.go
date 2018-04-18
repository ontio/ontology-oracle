package core

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/log"
	"github.com/ontio/ontology-oracle/models"
	"github.com/ontio/ontology-oracle/utils"
	"strings"
	"time"
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
	WS        *utils.WebSocketClient
	RPC       *utils.RpcClient
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

	//go app.OntListener()
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
		case "cron":
			go app.ExecuteCron(job)
		default:
			jobRun := job.NewRun()
			go app.ExecuteRun(jobRun)
		}
	}
}

func (app *OracleApplication) OntScanner() {
	log.Info("Start getting undo request in oracle contract.")
	app.RPC = utils.NewRpcClient()
	app.RPC.SetAddress(config.Configuration.ONTRPCAdress)

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

func (app *OracleApplication) OntListener() {
	address := config.Configuration.ONTWSAddress
	wsClient := utils.NewWebSocketClient(address)
	app.WS = wsClient
	recvCh, existCh, err := wsClient.Connet()
	if err != nil {
		log.Fatalf("NewWebSocketClient error: %v", err.Error())
	} else {
		log.Infof("NewWebSocketClient connected: %v", address)
	}
	go func() {
		for {
			select {
			case <-existCh:
				return
			case data := <-recvCh:
				resp := make(map[string]interface{}, 0)
				err := json.Unmarshal(data, &resp)
				if err != nil {
					log.Errorf("WS json.Unmarshal error: %v", err.Error())
					continue
				}
				err = app.ParseResp(resp)
				if err != nil {
					log.Errorf("Parse ontology transaction error: %v", err.Error())
				}
			}
		}
	}()
}
