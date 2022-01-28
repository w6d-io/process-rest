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

package config_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/w6d-io/process-rest/internal/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/w6d-io/x/cmdx"
)

var fileTest = `#!/bin/bash
echo "test"
`

var configTestFile = `
%s: %s
`
var configTestFileWithHook = `
%s: %s
hooks:
  - url: %s
    scope: "test"
`

var _ = Describe("Config", func() {
	When("Parse yaml config", func() {
		var configExitCode int
		var cmdExitCode int
		BeforeEach(func() {
			config.OsExit = func(code int) {
				configExitCode = code
			}
			cmdx.OsExit = func(code int) {
				cmdExitCode = code
			}
		})
		Context("New", func() {
			BeforeEach(func() {
				cmdExitCode = 0
			})
			AfterEach(func() {
				config.CfgFile = ""
			})
			It("fail due to file does not exist", func() {
				config.CfgFile = "testdata/no-file"
				config.Init()
				Expect(configExitCode).To(Equal(2))
			})
			It("fail unmarshal", func() {
				config.CfgFile = "testdata/fail-marshal.yaml"
				config.Init()
				Expect(cmdExitCode).To(Equal(1))
			})
			It("pre script folder does not exist", func() {
				config.CfgFile = "testdata/pre_script_does_not_exists.yaml"
				config.Init()
				Expect(cmdExitCode).To(Equal(1))
			})
			It("post script folder does not exist", func() {
				config.CfgFile = "testdata/post_script_does_not_exists.yaml"
				config.Init()
				Expect(cmdExitCode).To(Equal(1))
			})
			It("process script folder does not exist", func() {
				config.CfgFile = "testdata/process_script_does_not_exists.yaml"
				config.Init()
				Expect(cmdExitCode).To(Equal(1))
			})
			It("failed because of validation", func() {
				dir, err := ioutil.TempDir("", "post_dir")
				Expect(err).To(Succeed())
				_, err = ioutil.TempDir(dir, "directory")
				Expect(err).To(Succeed())
				filename := dir + string(os.PathSeparator) + "script1.sh"
				err = ioutil.WriteFile(filename, []byte(fileTest), 0644)
				Expect(err).To(Succeed())
				configFile := dir + string(os.PathSeparator) + "config.yaml"
				data := fmt.Sprintf(configTestFile, "post_script_folder", dir)
				err = ioutil.WriteFile(configFile, []byte(data), 0444)
				Expect(err).To(Succeed())
				config.CfgFile = configFile
				config.Init()
				Expect(configExitCode).To(Equal(2))
				config.Reset()
				err = os.RemoveAll(dir)
				Expect(err).To(Succeed())

			})
			It("success", func() {
				dir, err := ioutil.TempDir("", "process_dir")
				Expect(err).To(Succeed())
				_, err = ioutil.TempDir(dir, "directory")
				Expect(err).To(Succeed())
				filename := dir + string(os.PathSeparator) + "script1.sh"
				err = ioutil.WriteFile(filename, []byte(fileTest), 0644)
				Expect(err).To(Succeed())
				configFile := dir + string(os.PathSeparator) + "config.yaml"
				data := fmt.Sprintf(configTestFile, "main_script_folder", dir)
				err = ioutil.WriteFile(configFile, []byte(data), 0444)
				Expect(err).To(Succeed())
				config.CfgFile = configFile
				config.Init()
				config.Reset()
				err = os.RemoveAll(dir)
				Expect(err).To(Succeed())
			})
			It("success with hook", func() {
				dir, err := ioutil.TempDir("", "process_dir")
				Expect(err).To(Succeed())
				_, err = ioutil.TempDir(dir, "directory")
				Expect(err).To(Succeed())
				filename := dir + string(os.PathSeparator) + "script1.sh"
				err = ioutil.WriteFile(filename, []byte(fileTest), 0644)
				Expect(err).To(Succeed())
				configFile := dir + string(os.PathSeparator) + "config.yaml"
				data := fmt.Sprintf(configTestFileWithHook, "main_script_folder", dir, "http://localhost")
				err = ioutil.WriteFile(configFile, []byte(data), 0444)
				Expect(err).To(Succeed())
				config.CfgFile = configFile
				config.Init()
				config.Reset()
				err = os.RemoveAll(dir)
				Expect(err).To(Succeed())
			})
			It("failed on hook", func() {
				dir, err := ioutil.TempDir("", "process_dir")
				Expect(err).To(Succeed())
				_, err = ioutil.TempDir(dir, "directory")
				Expect(err).To(Succeed())
				filename := dir + string(os.PathSeparator) + "script1.sh"
				err = ioutil.WriteFile(filename, []byte(fileTest), 0644)
				Expect(err).To(Succeed())
				configFile := dir + string(os.PathSeparator) + "config.yaml"
				data := fmt.Sprintf(configTestFileWithHook, "main_script_folder", dir, "http://{}")
				err = ioutil.WriteFile(configFile, []byte(data), 0444)
				Expect(err).To(Succeed())
				config.Init()
				Expect(configExitCode).To(Equal(2))
				config.Reset()
				err = os.RemoveAll(dir)
				Expect(err).To(Succeed())
			})
		})
		Context("add script", func() {
			BeforeEach(func() {
			})
			AfterEach(func() {
			})
			It("succeed for postscript", func() {
				dir, err := ioutil.TempDir("", "post_dir")
				Expect(err).To(Succeed())
				_, err = ioutil.TempDir(dir, "directory")
				Expect(err).To(Succeed())
				filename := dir + string(os.PathSeparator) + "script1.sh"
				err = ioutil.WriteFile(filename, []byte(fileTest), 0644)
				Expect(err).To(Succeed())
				c := &config.Config{
					PostScriptFolder: dir,
				}
				err = c.AddPostScript()
				Expect(err).To(Succeed())
				config.Reset()
				err = os.RemoveAll(dir)
				Expect(err).To(Succeed())
			})
			It("succeed for postscript", func() {
				c := &config.Config{
					PostScriptFolder: "/no_such_folder",
				}
				err := c.AddPostScript()
				Expect(err).ToNot(Succeed())
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
			It("got nothing for postscript", func() {
				c := &config.Config{}
				err := c.AddPostScript()
				Expect(err).To(Succeed())
			})
			It("succeed for prescript", func() {
				dir, err := ioutil.TempDir("", "pre_dir")
				Expect(err).To(Succeed())
				_, err = ioutil.TempDir(dir, "directory")
				Expect(err).To(Succeed())
				filename := dir + string(os.PathSeparator) + "script1.sh"
				err = ioutil.WriteFile(filename, []byte(fileTest), 0644)
				Expect(err).To(Succeed())
				c := &config.Config{
					PreScriptFolder: dir,
				}
				err = c.AddPreScript()
				Expect(err).To(Succeed())
				config.Reset()
				err = os.RemoveAll(dir)
				Expect(err).To(Succeed())
			})
			It("succeed for prescript", func() {
				c := &config.Config{
					PreScriptFolder: "/no_such_folder",
				}
				err := c.AddPreScript()
				Expect(err).ToNot(Succeed())
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
			It("got nothing for postscript", func() {
				c := &config.Config{}
				err := c.AddPreScript()
				Expect(err).To(Succeed())
			})
			It("succeed for process script", func() {
				dir, err := ioutil.TempDir("", "process_dir")
				Expect(err).To(Succeed())
				_, err = ioutil.TempDir(dir, "directory")
				Expect(err).To(Succeed())
				filename := dir + string(os.PathSeparator) + "script1.sh"
				err = ioutil.WriteFile(filename, []byte(fileTest), 0644)
				Expect(err).To(Succeed())
				c := &config.Config{
					MainScriptFolder: dir,
				}
				err = c.AddProcessScript()
				Expect(err).To(Succeed())
				config.Reset()
				err = os.RemoveAll(dir)
				Expect(err).To(Succeed())
			})
			It("succeed for prescript", func() {
				c := &config.Config{
					MainScriptFolder: "/no_such_folder",
				}
				err := c.AddProcessScript()
				Expect(err).ToNot(Succeed())
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
			It("got nothing for postscript", func() {
				c := &config.Config{}
				err := c.AddProcessScript()
				Expect(err).To(Succeed())
			})
			It("check add function with empty path", func() {
				config.AddPreScript("")
				config.AddMainScript("")
				config.AddPostScript("")
			})
			It("", func() {
				scripts := config.GetPreScript()
				Expect(len(scripts)).To(Equal(0))
				scripts = config.GetMainScript()
				Expect(len(scripts)).To(Equal(0))
				scripts = config.GetPostScript()
				Expect(len(scripts)).To(Equal(0))
			})
		})
	})
})
