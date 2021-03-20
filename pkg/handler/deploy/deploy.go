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
package deploy

import (
	"github.com/gin-gonic/gin"
	"github.com/w6d-io/appdeploy/internal/process"
	"github.com/w6d-io/appdeploy/pkg/handler"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// @Param {payload}
// @Success 200 {object}
// @Failure 500 {object} httputil.HTTPError
// Deploy handle POST on /deploy
func Deploy(c *gin.Context) {

	log := logger
	if err := c.BindJSON(payload); err != nil {
		c.JSON(500, handler.Response{Message: "unmarshal failed", Error: err, Status: "error"})
		return
	}

	values, err := yaml.Marshal(payload)
	if err != nil {
		c.JSON(500, handler.Response{Message: "marshal values failed", Error: err, Status: "error"})
		return
	}

	file, err := ioutil.TempFile(os.TempDir(), "values-*.yaml")
	if err != nil {
		c.JSON(500, handler.Response{Message: "create values failed", Error: err, Status: "error"})
		return
	}
	defer func() {
		err := file.Close()
		log.Error(err, "Close file failed")
	}()

	if _, err := file.Write(values); err != nil {
		c.JSON(500, handler.Response{Message: "write values failed", Error: err, Status: "error"})
		return
	}

	err = process.Execute(file.Name())
	if err == nil {
		c.JSON(200, handler.Response{Message: "processing...", Status: "succeed"})
		return
	}
	deployError, ok := err.(Error)
	if !ok {
		c.JSON(500, handler.Response{Message: "process failed", Error: err})
		return
	}
	c.JSON(deployError.GetStatusCode(), deployError.GetResponse())
}
