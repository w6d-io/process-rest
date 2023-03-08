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
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/w6d-io/process-rest/pkg/router"
	"time"
)

var _ = Describe("Router", func() {
	Context("route engine", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("add a post handler", func() {
			router.AddPost("/test/unit", func(c *gin.Context) {})
		})
		It("add a get handler", func() {
			router.AddGet("/test/unit", func(c *gin.Context) {})
		})
		It("set listen", func() {
			router.SetListen(":8080")
		})
	})
	Context("run", func() {
		It("fail", func() {
			done := make(chan bool)
			var err error
			go func() {
				err = router.Run()
			}()
			go func() {
				for {
					select {
					case <-done:
						if err = router.Stop(); err != nil {
							return
						}
					}
				}
			}()
			time.Sleep(500 * time.Millisecond)
			done <- true
			Expect(err).To(Succeed())
		})
	})
})
