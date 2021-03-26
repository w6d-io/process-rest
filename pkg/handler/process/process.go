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
package process

import (
	"github.com/gin-gonic/gin"
	"github.com/w6d-io/process-rest/internal/process"
	"github.com/w6d-io/process-rest/pkg/router"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func init() {
	router.AddPost("/process", Process)
}

// Process handle POST on /process
func Process(c *gin.Context) {
	filename, err := InitProcess(c)
	if err != nil {
		processError := err.(Error)
		c.JSON(processError.GetStatusCode(), processError.GetResponse())
		return
	}
	ID, _ := c.GetQuery("id")
	go process.Execute(ID, filename)
	c.JSON(200, Response{Message: "processing...", Status: "succeed"})
}

func InitProcess(c *gin.Context) (string, error) {
	if err := c.BindJSON(payload); err != nil {
		return "", &ErrorProcess{Code: 500, Cause: err, Message: "unmarshal failed"}
	}
	values, err := yaml.Marshal(payload)
	if err != nil {
		return "", &ErrorProcess{Code: 500, Cause: err, Message: "marshal values failed"}
	}
	file, err := ioutil.TempFile("", "values-*.yaml")
	if err != nil {
		return "", &ErrorProcess{Code: 500, Cause: err, Message: "create values failed"}
	}
	defer func() {
		err := file.Close()
		logger.Error(err, "Close file failed")
	}()
	if _, err := file.Write(values); err != nil {
		return "", &ErrorProcess{Code: 500, Cause: err, Message: "write values failed"}
	}
	return file.Name(), nil
}

func (e *ErrorProcess) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

func (e *ErrorProcess) GetStatusCode() int {
	return e.Code
}

func (e *ErrorProcess) GetResponse() Response {
	return Response{
		Status:  "error",
		Message: e.Message,
		Error:   e.Cause,
	}
}
