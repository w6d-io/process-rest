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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/w6d-io/process-rest/internal/process"
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
	if !process.Validate() {
		return errors.New("a process script should be set")
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
		process.AddPostScript(path)
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
		process.AddPreScript(path)
	}
	return nil
}

func (c *Config) AddProcessScript() error {
	log := ctrl.Log.WithName("Config").WithName("AddMainScript")
	if c.ProcessScriptFolder == "" {
		return nil
	}
	files, err := ioutil.ReadDir(c.ProcessScriptFolder)
	if err != nil {
		log.Error(err, "get file in folder failed", "folder", c.ProcessScriptFolder)
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := fmt.Sprintf("%s%c%s", c.ProcessScriptFolder, os.PathSeparator, file.Name())
		process.AddMainScript(path)
	}
	return nil
}
