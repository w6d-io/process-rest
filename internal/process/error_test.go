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
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/w6d-io/process-rest/internal/process"
	"github.com/w6d-io/process-rest/pkg/handler"
)

var _ = Describe("Error", func() {
	Context("check method", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("returns message because there is no cause", func() {
			err := process.NewError(nil, 500, "test with no cause")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("test with no cause"))
		})
		It("returns message and its cause", func() {
			err := process.NewError(errors.New("all goes wrong"), 500, "test with cause")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("test with cause : all goes wrong"))

		})
		It("get status code", func() {
			err := &process.Error{
				Code:    500,
				Cause:   errors.New("all goes wrong"),
				Message: "test with cause",
			}
			Expect(err.GetStatusCode()).To(Equal(500))
			Expect(err.GetResponse()).To(Equal(handler.Response{Status: "error", Message: err.Message, Error: err.Cause}))
		})
	})
})
