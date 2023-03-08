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
package router_test

import (
	"net/http"

	"github.com/w6d-io/process-rest/pkg/router"

	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Log", func() {
	var (
		c *gin.Context
	)
	Describe("For gin gonic", func() {
		Context("log output", func() {
			outputFunc := router.LogMiddleware()
			It("gin output handlerFunc", func() {
				Expect(outputFunc).ToNot(BeNil())
			})
			Context("Correlation ID", func() {
				correlationID := router.CorrelationID()
				It("gin id header handlerFunc", func() {
					Expect(correlationID).ToNot(BeNil())
				})
				When("is not empty", func() {
					It("set the correlation id in the http header", func() {
						c2 := &gin.Context{Request: &http.Request{Header: http.Header{}}}
						cr := router.CorrelationID()
						cr(c2)
					})
				})
			})
			Context("Get client ip address", func() {
				c = &gin.Context{Request: &http.Request{Header: http.Header{}}}
				c.Request.RemoteAddr = "10.0.0.2,10.0.0.3"
				When("X-Real-IP is not set", func() {
					It("return the remote address", func() {
						Expect(router.GetClientIP(c)).Should(Equal("10.0.0.2"))
					})
				})
				When("X-Real-IP is set", func() {
					It("returns the ip address", func() {
						c.Request.Header.Set("X-Real-IP", "10.0.0.1")
						Expect(router.GetClientIP(c)).Should(Equal("10.0.0.1"))
					})
				})
			})
			Context("Gin handler function", func() {
				It("Json log", func() {
					jsonLog := router.LogMiddleware()
					c.Request.Method = "POST"
					jsonLog(c)
				})
			})
		})
	})
})
