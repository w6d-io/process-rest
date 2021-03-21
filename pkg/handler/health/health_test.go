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

package health_test

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/w6d-io/process-rest/pkg/handler/health"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("health", func() {
	Context("", func() {
		It("", func() {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			health.Health(c)
			//c.Writer.Status()
			Expect(c.Writer.Status()).To(Equal(200))
		})
	})
})
