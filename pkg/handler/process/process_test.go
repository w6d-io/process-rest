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
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/util/framer"

	"github.com/w6d-io/process-rest/pkg/handler/process"
)

var _ = Describe("Process", func() {
	Context("", func() {
		BeforeEach(func() {
			process.YamlMarshal = yaml.Marshal
			process.IoTempFile = os.CreateTemp
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
			r := io.NopCloser(strings.NewReader(payload))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			URL, err := url.Parse("http://localhost:8888/process?id=a9bac696-f21e-4149-9018-cf882e5bf8e7")
			Expect(err).To(Succeed())
			c.Request = &http.Request{
				Body: framer.NewJSONFramedReader(r),
				URL:  URL,
			}
			process.Process(c)
			Expect(c.Writer.Status()).To(Equal(200))
		})
		It("payload well consisted, force yaml error", func() {
			payload := `
{
  "global": { "label": "test-integration" },
  "redis": { "enabled": true }
}
`
			r := io.NopCloser(strings.NewReader(payload))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			URL, err := url.Parse("http://localhost:8888/process?id=a9bac696-f21e-4149-9018-cf882e5bf8e7")
			Expect(err).To(Succeed())
			c.Request = &http.Request{
				Body: framer.NewJSONFramedReader(r),
				URL:  URL,
			}
			process.YamlMarshal = func(in interface{}) (out []byte, err error) {
				return nil, errors.New("yaml marshal error")
			}
			process.Process(c)
			Expect(c.Writer.Status()).To(Equal(500))
		})
		It("payload well consisted, force iotemp creation error", func() {
			payload := `
{
  "global": { "label": "test-integration" },
  "redis": { "enabled": true }
}
`
			r := io.NopCloser(strings.NewReader(payload))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			URL, err := url.Parse("http://localhost:8888/process?id=a9bac696-f21e-4149-9018-cf882e5bf8e7")
			Expect(err).To(Succeed())
			c.Request = &http.Request{
				Body: framer.NewJSONFramedReader(r),
				URL:  URL,
			}
			process.IoTempFile = func(dir, pattern string) (f *os.File, err error) {
				return nil, errors.New("io temp error")
			}
			process.Process(c)
			Expect(c.Writer.Status()).To(Equal(500))
		})
		It("return 500 due to malformed payload", func() {
			payload := `
{
  "global": { "label": "test-integration",
}
`
			r := io.NopCloser(strings.NewReader(payload))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			URL, err := url.Parse("http://localhost:8888/process?id=a9bac696-f21e-4149-9018-cf882e5bf8e7")
			Expect(err).To(Succeed())
			c.Request = &http.Request{
				Body: framer.NewJSONFramedReader(r),
				URL:  URL,
			}
			process.Process(c)
			Expect(c.Writer.Status()).To(Equal(400))
		})
		It("", func() {
			payload := `
{
  "global": { "label": "test-integration" },
  "redis": { "enabled": true }
}
`
			r := io.NopCloser(strings.NewReader(payload))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			URL, err := url.Parse("http://localhost:8888/process?id=a9bac696-f21e-4149-9018-cf882e5bf8e7")
			Expect(err).To(Succeed())
			c.Request = &http.Request{
				Body: framer.NewJSONFramedReader(r),
				URL:  URL,
			}
			process.Process(c)
			Expect(c.Writer.Status()).To(Equal(200))
		})
		It("get error Message", func() {
			e := process.ErrorProcess{
				Cause:   errors.New("test"),
				Code:    500,
				Message: "Test",
			}
			Expect(e.Error()).To(Equal("Test : test"))
			e.Cause = nil
			Expect(e.Error()).To(Equal("Test"))
		})
		It("", func() {
			cid := process.GetCorrelationID(nil)
			Expect(cid).ToNot(BeEmpty())
		})
	})
})
