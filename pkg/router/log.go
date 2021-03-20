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
package router

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	ctrl "sigs.k8s.io/controller-runtime"
)

// LogMiddleware logs a gin HTTP request in JSON format, with some additional custom key/values
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := ctrl.Log.WithName("gin")
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := GetDurationInMilliseconds(start)
		var (
			status = http.StatusNotImplemented
			corID  string
		)
		if c.Writer != nil {
			status = c.Writer.Status()
			corID = c.Writer.Header().Get(CorrelationId)
		}

		entry := log.WithValues(
			"client_ip", GetClientIP(c),
			"duration", duration,
			"method", c.Request.Method,
			"uri", c.Request.RequestURI,
			"status", status,
			"referrer", c.Request.Referer(),
			"correlation_id", corID,
		)

		if status >= http.StatusInternalServerError {
			entry.Error(errors.New(c.Errors.String()), "logging")
		} else {
			entry.Info(c.Request.Method + " " + c.Request.RequestURI)
		}
	}
}

// GetClientIP gets the correct IP for the end client instead of the proxy
func GetClientIP(c *gin.Context) string {
	// first check the X-Forwarded-For header
	requester := c.Request.Header.Get("X-Forwarded-For")
	// if empty, check the Real-IP header
	if len(requester) == 0 {
		requester = c.Request.Header.Get("X-Real-IP")
	}
	// if the requester is still empty, use the hard-coded address from the socket
	if len(requester) == 0 {
		requester = c.Request.RemoteAddr
	}

	// if requester is a comma delimited list, take the first one
	// (this happens when proxied via elastic load balancer then again through nginx)
	if strings.Contains(requester, ",") {
		requester = strings.Split(requester, ",")[0]
	}

	return requester
}

// GetDurationInMilliseconds takes a start time and returns a duration in milliseconds
func GetDurationInMilliseconds(start time.Time) time.Duration {
	end := time.Now()
	duration := end.Sub(start)
	//milliseconds := float64(duration) / float64(time.Millisecond)
	//rounded := float64(int(milliseconds*100+.5)) / 100
	return duration
}

// CorrelationID return a gin handler and add a field in the http header
func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var correlationID string
		if c.Writer != nil {
			correlationID = c.Writer.Header().Get(CorrelationId)
			if correlationID == "" {
				c.Header(CorrelationId, uuid.New().String())
			}
		}
	}
}
