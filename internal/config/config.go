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
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/w6d-io/hook"
	"gopkg.in/yaml.v3"

	ctrl "sigs.k8s.io/controller-runtime"
)

// get the config file and process it
func New(filename string) error {
	log := ctrl.Log.WithName("Config")
	log.V(1).Info("read config file")
	config = new(Config)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err, "error reading the configuration")
		return err
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		log.Error(err, "Error unmarshal the configuration")
		return err
	}
	if err := config.AddPreScript(); err != nil {
		return err
	}
	if err := config.AddProcessScript(); err != nil {
		return err
	}
	if err := config.AddPostScript(); err != nil {
		return err
	}
	if !Validate() {
		return errors.New("a process script should be set")
	}
	for _, wh := range config.Hooks {
		if err := hook.Subscribe(wh.URL, wh.Scope); err != nil {
			log.Error(err, "hook subscription failed")
			return err
		}
	}
	return nil
}

func (c *Config) AddPostScript() error {
	log := ctrl.Log.WithName("Config").WithName("AddPostScript")
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
	log := ctrl.Log.WithName("Config").WithName("AddPreScript")
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
	log := ctrl.Log.WithName("Config").WithName("AddMainScript")
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
	log := ctrl.Log.WithName("Validate")
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
