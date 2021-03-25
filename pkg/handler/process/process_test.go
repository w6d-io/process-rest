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
package process_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/w6d-io/process-rest/pkg/handler/process"
	"k8s.io/apimachinery/pkg/util/framer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Process", func() {
	Context("", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("payload well consisted", func() {
			payload := `
{
  "global": { "label": "test-integration" },
  "redis": { "enabled": true }
}
`
			r := ioutil.NopCloser(strings.NewReader(payload))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = &http.Request{
				Body: framer.NewJSONFramedReader(r),
			}
			process.Process(c)
			Expect(c.Writer.Status()).To(Equal(200))
		})
		It("return 500 due to malformed payload", func() {
			payload := `
{
  "global": { "label": "test-integration",
}
`
			r := ioutil.NopCloser(strings.NewReader(payload))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = &http.Request{
				Body: framer.NewJSONFramedReader(r),
			}
			process.Process(c)
			Expect(c.Writer.Status()).To(Equal(500))
		})
		It("", func() {
			payload := `
{
  "global": { "label": "test-integration" },
  "redis": { "enabled": true }
}
`
			r := ioutil.NopCloser(strings.NewReader(payload))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = &http.Request{
				Body: framer.NewJSONFramedReader(r),
			}
			process.Process(c)
			Expect(c.Writer.Status()).To(Equal(200))
		})
	})
})