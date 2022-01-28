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

var (
	payload = new(Payload)
)

// payload is the values from request
type Payload map[string]interface{}

type Error interface {
	Error() string
	// GetResponse returns the Response struct
	GetResponse() Response
	// GetStatusCode returns http status code.
	GetStatusCode() int
}

type ErrorProcess struct {
	Cause   error
	Code    int
	Message string
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   error  `json:"error,omitempty"`
}
