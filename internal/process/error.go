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

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

func (e *Error) GetStatusCode() int {
	return e.Code
}

//func (e *Error) GetResponse() Response {
//	return Response{
//		Status:  "error",
//		Message: e.Message,
//		Error:   e.Cause,
//	}
//}

func NewError(cause error, code int, message string) error {
	return &Error{
		Code:    code,
		Cause:   cause,
		Message: message,
	}
}
