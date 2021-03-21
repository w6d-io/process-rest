/*
Copyright 2020 WILDCARD SA.

Licensed under the WILDCARD SA License, Version 1.0 (the "License");
WILDCARD SA is register in french corporation.
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.w6d.io/licenses/LICENSE-1.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is prohibited.
Created on 21/03/2021
*/

package process_test

import (
	"io/ioutil"
	"os"

	"github.com/w6d-io/process-rest/internal/process"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	configTestFile = `
%s: %s
`
	successTest = `#!/bin/bash
echo "test"
exit 0
`
	failTest = `#!/bin/bash
echo "failing test"
exit 1
`
)

var _ = Describe("Process", func() {
	Context("Execute", func() {
		var (
			dir       string
			filename  string
			filename2 string
			err       error
		)
		BeforeEach(func() {
			dir, err = ioutil.TempDir("", "test_dir")
			Expect(err).To(Succeed())
			_, err = ioutil.TempDir(dir, "directory")
			Expect(err).To(Succeed())
			filename = dir + string(os.PathSeparator) + "script1.sh"
			err = ioutil.WriteFile(filename, []byte(successTest), 0755)
			Expect(err).To(Succeed())
			filename2 = dir + string(os.PathSeparator) + "script2.sh"
			err = ioutil.WriteFile(filename2, []byte(failTest), 0755)
			Expect(err).To(Succeed())
		})
		AfterEach(func() {
			process.Reset()
			err = os.RemoveAll(dir)
			Expect(err).To(Succeed())
		})
		It("Do nothing", func() {
			process.AddPreScript("")
			process.AddMainScript("")
			process.AddPostScript("")
			Expect(process.Validate()).To(Equal(false))
			err = process.Execute()
			Expect(err).To(Succeed())
		})
		It("runs pre script with success", func() {
			process.AddPreScript(filename)
			err = process.Execute()
			Expect(err).To(Succeed())
		})
		It("runs main script with success", func() {
			process.AddMainScript(filename)
			err = process.Execute()
			Expect(err).To(Succeed())
		})
		It("runs post script with success", func() {
			process.AddPostScript(filename)
			err = process.Execute()
			Expect(err).To(Succeed())
		})
		It("runs pre script with failure", func() {
			process.AddPreScript(filename2)
			err = process.Execute()
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("pre process failed"))
		})
		It("runs main script with failure", func() {
			process.AddMainScript(filename2)
			err = process.Execute()
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("main process failed "))
		})
		It("runs post script with failure", func() {
			process.AddPostScript(filename2)
			err = process.Execute()
			Expect(err).ToNot(Succeed())
			Expect(err.Error()).To(ContainSubstring("post process failed"))
		})

	})
})
