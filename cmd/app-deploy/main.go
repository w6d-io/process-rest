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
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/w6d-io/app-deploy/internal/util"
	"github.com/w6d-io/app-deploy/pkg/router"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	zapraw "go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	setupLog = ctrl.Log.WithName("setup")

	// Version microservice version
	Version = ""

	// Revision git commit
	Revision = ""

	// GoVersion ...
	GoVersion = ""

	// Built Date built
	Built = ""

	// OsArch ...
	OsArch = ""
)

func main() {
	setupLog.Info("managed flag")
	opts := zap.Options{
		Encoder: zapcore.NewConsoleEncoder(util.TextEncoderConfig()),
	}
	util.BindFlags(&opts, flag.CommandLine)
	flag.Parse()
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	if !seen["config"] {
		fmt.Print("config file is missing\n")
		setupLog.Error(errors.New("flag error"), "config file is missing")
		os.Exit(1)
	}

	opts.Development = os.Getenv("RELEASE") != "prod"
	opts.StacktraceLevel = zapcore.PanicLevel
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts), zap.RawZapOpts(zapraw.AddCaller())))

	setupLog.Info("starting app-deploy", "Version", Version, "Built",
		Built, "Revision", Revision, "Arch", OsArch, "GoVersion", GoVersion)
	if err := router.Run(); err != nil {
		setupLog.Error(err, "run server")
	}
}
