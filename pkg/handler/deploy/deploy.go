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
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/w6d-io/app-deploy/internal/process"
	"github.com/w6d-io/app-deploy/pkg/handler"
	"gopkg.in/yaml.v3"
)

// @Param {payload}
// @Success 200 {object}
// @Failure 500 {object} httputil.HTTPError
// Deploy handle POST on /deploy
func Deploy(c *gin.Context) {
	file := new(os.File)
	err := InitDeploy(c, file)
	if err != nil {
		deployError := err.(Error)
		c.JSON(deployError.GetStatusCode(), deployError.GetResponse())
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

func InitDeploy(c *gin.Context, file *os.File) error {
	if err := c.BindJSON(payload); err != nil {
		return &ErrorDeploy{Code: 500, Cause: err, Message: "unmarshal failed"}
	}
	values, err := yaml.Marshal(payload)
	if err != nil {
		return &ErrorDeploy{Code: 500, Cause: err, Message: "marshal values failed"}
	}
	file, err = ioutil.TempFile("", "values-*.yaml")
	if err != nil {
		return &ErrorDeploy{Code: 500, Cause: err, Message: "create values failed"}
	}
	defer func() {
		err := file.Close()
		logger.Error(err, "Close file failed")
	}()
	if _, err := file.Write(values); err != nil {
		return &ErrorDeploy{Code: 500, Cause: err, Message: "write values failed"}
	}
	return nil
}

func (e *ErrorDeploy) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

func (e *ErrorDeploy) GetStatusCode() int {
	return e.Code
}

func (e *ErrorDeploy) GetResponse() handler.Response {
	return handler.Response{
		Status: "error",
		Message: e.Message,
		Error: e.Cause,
	}
}