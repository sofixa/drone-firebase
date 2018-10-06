// Copyright 2016 Google Inc. All Rights Reserved.
// Copyright 2018 Adrian Todorov
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

var (
	buildCommit string
)

func main() {

	app := cli.NewApp()
	app.Name = "hugo plugin"
	app.Usage = "hugo plugin"
	app.Action = run
	app.Version = fmt.Sprintf("%s", buildCommit)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token",
			Usage:  "The token to use for login",
			EnvVar: "PLUGIN_TOKEN",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "project",
			Usage:  "(optional) The project alias to deploy to. If not set, uses the default specified in the .firebaserc file. Note: This does currently not work for Firebase project IDs, so you must run firebase use --add to add aliases for your different environments. Then use this field to specify the alias name such as default, staging or production",
			EnvVar: "PLUGIN_PROJECT",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "message",
			Usage:  "(optional) The message to use for your commit. You can use variables available from Drone.io, such as $$COMMIT as part of the message. The message does not have to be quoted",
			EnvVar: "PLUGIN_MESSAGE",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "targets",
			Usage:  "(optional) The type of deployment to be done. Must be a comma separated list of the following types: hosting, database, and storage",
			EnvVar: "PLUGIN_TARGETS",
			Value:  "",
		},
		cli.BoolFlag{
			Name:   "dryrun",
			Usage:  "Whether the deploy commands should be executed, or just printed to stdout",
			EnvVar: "PLUGIN_DRYRUN",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "(optional) A bool indicating whether commands should run in debug mode. WARNING: This will print all information about requests such as authentication information!",
			EnvVar: "PLUGIN_DEBUG",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
func run(c *cli.Context) error {
	plugin := Plugin{
		Token:   c.String("token"),
		Project: c.String("project"),
		Message: c.String("message"),
		Targets: c.String("targets"),
		DryRun:  c.Bool("dryrun"),
		Debug:   c.Bool("debug"),
	}
	return plugin.Exec()

}
