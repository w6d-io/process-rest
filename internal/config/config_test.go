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
	"github.com/w6d-io/process-rest/internal/config"
	"github.com/w6d-io/process-rest/internal/process"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var fileTest = `#!/bin/bash
echo "test"
`

var configTestFile = `
%s: %s
`

var _ = Describe("Config", func() {
	Context("New", func() {
		It("fail due to file does not exist", func() {
			err := config.New("testdata/no-file")
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("no such file or directory"))
		})
		It("fail unmarshal", func() {
			err := config.New("testdata/fail-marshal.yaml")
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("did not find expected node content"))
		})
		It("pre script folder does not exist", func() {
			err := config.New("testdata/pre_script_does_not_exists.yaml")
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("pre_not_exists"))
		})
		It("post script folder does not exist", func() {
			err := config.New("testdata/post_script_does_not_exists.yaml")
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("post_not_exists"))
		})
		It("process script folder does not exist", func() {
			err := config.New("testdata/process_script_does_not_exists.yaml")
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("process_not_exists"))
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
			err = config.New(configFile)
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(Equal("a process script should be set"))
			process.Reset()
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
			data := fmt.Sprintf(configTestFile, "process_script_folder", dir)
			err = ioutil.WriteFile(configFile, []byte(data), 0444)
			Expect(err).To(Succeed())
			err = config.New(configFile)
			Expect(err).To(Succeed())
			process.Reset()
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
			process.Reset()
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
			process.Reset()
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
				ProcessScriptFolder: dir,
			}
			err = c.AddProcessScript()
			Expect(err).To(Succeed())
			process.Reset()
			err = os.RemoveAll(dir)
			Expect(err).To(Succeed())
		})
		It("succeed for prescript", func() {
			c := &config.Config{
				ProcessScriptFolder: "/no_such_folder",
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

	})
})
