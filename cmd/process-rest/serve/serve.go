/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 20/03/2021
*/

package serve

import (
	"github.com/spf13/cobra"

	"github.com/w6d-io/x/logx"
	"github.com/w6d-io/x/pflagx"

	config "github.com/w6d-io/process-rest/internal/config"
	"github.com/w6d-io/process-rest/pkg/handler"
	router "github.com/w6d-io/process-rest/pkg/router"
)

var (
	Cmd = &cobra.Command{
		Use:   "serve",
		Short: "Run the project server",
		RunE:  serve,
	}

	_ = handler.Handler{}
)

func init() {
	cobra.OnInitialize(config.Init)

	pflagx.CallerSkip = -2
	pflagx.Init(Cmd, &config.CfgFile)
}

func serve(_ *cobra.Command, _ []string) error {
	log := logx.WithName(nil, "Serve.Command")

	if err := router.Run(); err != nil {
		log.Error(err, "run server")
		return err
	}

	return nil
}
