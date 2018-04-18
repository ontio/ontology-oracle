package node

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/core"
	"github.com/ontio/ontology-oracle/http"
	"github.com/ontio/ontology-oracle/log"
	"github.com/urfave/cli"
	"os"
)

// Client is the shell for the node. It has fields for the Renderer,
// Config, AppFactory (the services application), Authenticator, and Runner.
type Client struct {
	AppFactory AppFactory
	Runner     Runner
}

// RunNode starts the oracle node.
func (client *Client) RunNode(c *cli.Context) error {

	log.Info("Starting Oracle Node... ")
	log.Info("Open the account")
	if !common.FileExisted(account.WALLET_FILENAME) {
		log.Fatal(fmt.Sprintf("No %s detected, please create a wallet.", account.WALLET_FILENAME))
		os.Exit(1)
	}
	accountClient := account.Open(account.WALLET_FILENAME, []byte(c.String("password")))
	if client == nil {
		log.Fatal("Can't get local account.")
		os.Exit(1)
	}
	acct := accountClient.GetDefaultAccount()
	app := client.AppFactory.NewApplication(acct)

	app.Start()
	defer app.Stop()

	return client.errorOut(client.Runner.Run(app))
}

func (client *Client) errorOut(err error) error {
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// AppFactory implements the NewApplication method.
type AppFactory interface {
	NewApplication(*account.Account) core.Application
}

// AppFactory is used to create a new Application.
type OracleAppFactory struct{}

// NewApplication returns a new instance of the node with the given config.
func (n OracleAppFactory) NewApplication(account *account.Account) core.Application {
	return core.NewApplication(account)
}

// Runner implements the Run method.
type Runner interface {
	Run(core.Application) error
}

// OracleRunner is used to run the node application.
type OracleRunner struct{}

// Run sets the log level based on config and starts the web router to listen
// for input and return data.
func (n OracleRunner) Run(app core.Application) error {
	port := config.Configuration.Port
	gin.SetMode(gin.DebugMode)
	return http.Router(app.(*core.OracleApplication)).Run(":" + port)
}
