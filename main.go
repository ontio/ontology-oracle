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
			Name:   "node",
			Usage:  "Run the oracle node",
			Action: client.RunNode,
		},
	}
	app.Run(args)
}

func NewProductionClient() *node.Client {
	return &node.Client{
		AppFactory: node.OracleAppFactory{},
	}
}
