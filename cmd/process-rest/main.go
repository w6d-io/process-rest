/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 20/03/2021
*/

package main

import (
	"os"

	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"

	"github.com/w6d-io/process-rest/cmd/process-rest/serve"
	"github.com/w6d-io/process-rest/internal/config"
	"github.com/w6d-io/x/logx"
)

var rootCmd = &cobra.Command{
	Use: "project",
	Run: func(cmd *cobra.Command, args []string) {
		log := logx.WithName(nil, "Main.Command")
		err := cmd.Help()
		if err != nil {
			log.Error(err, "cannot show help")
		}
	},
}

func main() {
	log := logx.WithName(nil, "Main.Command")

	rootCmd.AddCommand(cmdx.Version(&config.Version, &config.Revision, &config.Built))
	rootCmd.AddCommand(serve.Cmd)
	if err := rootCmd.Execute(); err != nil {
		log.Error(err, "exec command failed")
		os.Exit(1)
	}
}
