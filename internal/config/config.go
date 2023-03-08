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
package config

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/w6d-io/hook"

	"github.com/w6d-io/x/cmdx"
	"github.com/w6d-io/x/logx"
)

var (
	// Version microservice version
	Version = ""

	// Revision git commit
	Revision = ""

	// Built Date built
	Built = ""

	// CfgFile contain the path of the config file
	CfgFile string

	// OsExit is hack for unit-test
	OsExit = os.Exit
)

// Init load the config file
func Init() {
	log := logx.WithName(nil, "Config.Init")
	config = new(Config)
	data, err := os.ReadFile(CfgFile)
	if err != nil {
		log.Error(err, "error reading the configuration")
		OsExit(2)
		return
	}
	err = yaml.Unmarshal(data, config)
	cmdx.Must(err, "error unmarshal the configuration")

	err = yaml.Unmarshal(data, config)
	cmdx.Must(err, "Error unmarshal the configuration")

	err = config.AddPreScript()
	cmdx.Must(err, "Error checking AddPreScript")

	err = config.AddProcessScript()
	cmdx.Must(err, "Error checking AddProcessScript")

	err = config.AddPostScript()
	cmdx.Must(err, "Error checking AddPostScript")

	if !Validate() {
		log.Error(errors.New("a process script should be set"), "")
		OsExit(2)
		return
	}
	for _, wh := range config.Hooks {
		if err := hook.Subscribe(context.Background(), wh.URL, wh.Scope); err != nil {
			log.Error(err, "hook subscription failed")
			OsExit(2)
			return
		}
	}
}

func (c *Config) AddPostScript() error {
	log := logx.WithName(nil, "Config.AddPostScript")
	if c.PostScriptFolder == "" {
		return nil
	}
	files, err := ioutil.ReadDir(c.PostScriptFolder)
	if err != nil {
		log.Error(err, "get file in folder failed", "folder", c.PostScriptFolder)
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := fmt.Sprintf("%s%c%s", c.PostScriptFolder, os.PathSeparator, file.Name())
		AddPostScript(path)
	}
	return nil
}

func (c *Config) AddPreScript() error {
	log := logx.WithName(nil, "Config.AddPreScript")
	if c.PreScriptFolder == "" {
		return nil
	}
	files, err := ioutil.ReadDir(c.PreScriptFolder)
	if err != nil {
		log.Error(err, "get file in folder failed", "folder", c.PreScriptFolder)
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := fmt.Sprintf("%s%c%s", c.PreScriptFolder, os.PathSeparator, file.Name())
		AddPreScript(path)
	}
	return nil
}

func (c *Config) AddProcessScript() error {
	log := logx.WithName(nil, "Config.AddMainScript")
	if c.MainScriptFolder == "" {
		return nil
	}
	files, err := ioutil.ReadDir(c.MainScriptFolder)
	if err != nil {
		log.Error(err, "get file in folder failed", "folder", c.MainScriptFolder)
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := fmt.Sprintf("%s%c%s", c.MainScriptFolder, os.PathSeparator, file.Name())
		AddMainScript(path)
	}
	return nil
}

// AddPostScript appends the path to post script
func AddPostScript(path string) {
	if path == "" {
		return
	}
	postScript = append(postScript, path)
}

// AddPreScript appends the path to pre script
func AddPreScript(path string) {
	if path == "" {
		return
	}
	preScript = append(preScript, path)
}

// AddMainScript appends the path to pre script
func AddMainScript(path string) {
	if path == "" {
		return
	}
	mainScript = append(mainScript, path)
}

func Reset() {
	preScript = []string{}
	mainScript = []string{}
	postScript = []string{}
}

func Validate() bool {
	log := logx.WithName(nil, "Config.Validate")
	log.V(1).Info("contain", "pre_script", preScript,
		"main_script", mainScript,
		"post_script", postScript)
	return len(mainScript) != 0
}

func GetPreScript() []string {
	return preScript
}

func GetMainScript() []string {
	return mainScript
}

func GetPostScript() []string {
	return postScript
}
