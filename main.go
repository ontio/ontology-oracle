package main

import (
	"os"

	"github.com/ontio/ontology-oracle/log"
	"github.com/ontio/ontology-oracle/node"
	"github.com/urfave/cli"
)

func main() {
	log.Init(log.PATH, log.Stdout)
	Run(NewProductionClient(), os.Args...)
}

func Run(client *node.Client, args ...string) {
	app := cli.NewApp()
	app.Usage = "cli for ontology-oracle"
	app.Commands = []cli.Command{
		{
			Name:    "node",
			Aliases: []string{"n"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "password, p",
					Usage: "password for the node's account",
				},
			},
			Usage:  "Run the oracle node",
			Action: client.RunNode,
		},
	}
	app.Run(args)
}

func NewProductionClient() *node.Client {
	return &node.Client{
		AppFactory: node.OracleAppFactory{},
		Runner:     node.OracleRunner{},
	}
}
