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
	"fmt"
	"github.com/w6d-io/hook"
	"github.com/w6d-io/process-rest/internal/config"
	"os/exec"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"
)

func Run(name string, arg ...string) (string, error) {
	log := ctrl.Log.WithName("Process")
	log.V(1).Info("build command")
	cmd := exec.Command(name, arg...)

	log.V(1).Info("exec command and get output", "script", cmd.String())
	output, err := cmd.Output()
	if err != nil {
		log.Error(err, "script failed", "script", cmd.String())
		return "", err
	}
	log.V(1).Info("script succeeded", "output", string(output))
	return string(output), nil
}

func LoopProcess(scripts []string, outputs map[string]Output, arg ...string) error {
	log := ctrl.Log.WithName("LoopProcess")
	for _, script := range scripts {
		log.Info("run", "script", script)
		arg = append([]string{"-c", script}, arg...)
		output, err := Run("bash", arg...)
		if err != nil {
			log.Error(err, "process failed", "script", script)
			outputs[script] = Output{
				Status: "failed",
				Log:    output,
				Error:  err.Error(),
			}
			return err
		}
	}
	return nil
}

func PreProcess(outputs map[string]Output, arg ...string) error {
	ctrl.Log.WithName("PreProcess").V(1).Info("loop process")
	return LoopProcess(config.GetPreScript(), outputs, arg...)
}

func PostProcess(outputs map[string]Output, arg ...string) error {
	ctrl.Log.WithName("PostProcess").V(1).Info("loop process")
	return LoopProcess(config.GetPostScript(), outputs, arg...)
}

func MainProcess(outputs map[string]Output, arg ...string) error {
	ctrl.Log.WithName("MainProcess").V(1).Info("loop process")
	return LoopProcess(config.GetMainScript(), outputs, arg...)
}

func Execute(id string, arg ...string) {
	log := ctrl.Log.WithName("Execute")
	outputs := make(map[string]Output)
	errc := make(chan error, 2)
	go func() {
		// do pre-process
		if err := PreProcess(outputs, arg...); err != nil {
			log.Error(err, "pre process failed")
			Notify(id, outputs, "pre-process-failed", err)
			errc <- NewError(err, 550, "pre process failed")
			return
		}

		// do main process
		if err := MainProcess(outputs, arg...); err != nil {
			log.Error(err, "main process failed")
			Notify(id, outputs, "main-process-failed", err)
			errc <- NewError(err, 551, "main process failed")
			return
		}

		// do post-process
		if err := PostProcess(outputs, arg...); err != nil {
			log.Error(err, "post process failed")
			Notify(id, outputs, "post-process-failed", err)
			errc <- NewError(err, 552, "post process failed")
			return
		}
		Notify(id, outputs, "process-succeeded", nil)

	}()

	for range []string{"1", "2"} {
		if err := <-errc; err != nil {
			log.Error(err, "process failed")
		}
	}
}

func Notify(id string, outputs map[string]Output, scope string, err error) {
	log := ctrl.Log.WithName("Notify")

	status := &Status{
		Succes: err != nil,
		Log:    GetLogMessage(err, outputs),
		ID:     id,
	}
	if err := hook.Send(status, ctrl.Log, scope); err != nil {
		log.Error(err, "notification failed")
	}

}

func GetLogMessage(err error, outputs map[string]Output) string {
	var messages []string
	if err != nil {
		messages = append(messages, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
	for key := range outputs {
		message := fmt.Sprintf(`{"script":"%s", "error":"%s", "status":"%s", "log":"%s"}`,
			key, outputs[key].Error, outputs[key].Status, outputs[key].Log)
		messages = append(messages, message)
	}
	return fmt.Sprintf("{%s}", strings.Join(messages, ","))
}
