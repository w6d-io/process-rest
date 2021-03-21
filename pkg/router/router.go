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
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func init() {
	engine.Use(LogMiddleware())
	engine.Use(gin.Recovery())
	engine.Use(CorrelationID())
}

// AddPost binds a function/method to a relative path in POST http method
func AddPost(relativePath string, handlers ...gin.HandlerFunc) {
	engine.POST(relativePath, handlers...)
}

// AddGet binds a function/method to a relative path in POST http method
func AddGet(relativePath string, handlers ...gin.HandlerFunc) {
	engine.GET(relativePath, handlers...)
}

func SetListen(address string) {
	server.Addr = address
}

func Run() error {
	server.Handler = engine

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)
	go func() {
		<-quit
		logger.Info("receive interrupt or kill signal")
		if err := server.Close(); err != nil {
			logger.Error(err, "Server closed")
			os.Exit(1)
		}
	}()
	logger.WithValues("address", server.Addr).Info("Listening and serving HTTP")
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logger.Info("Server closed under request")
			return nil
		}
		logger.Error(err, "Server closed unexpect")
		return err
	}
	return nil
}

// Stop the http server
func Stop() error {
	if server != nil {
		return server.Close()
	}
	return nil
}
