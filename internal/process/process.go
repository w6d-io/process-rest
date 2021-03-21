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
	"os/exec"

	ctrl "sigs.k8s.io/controller-runtime"
)

// AddPostScript appends the path to post script
func AddPostScript(path string) {
	if path == "" {
		return
	}
	postScript = append(postScript, path)
}

// AddPreScript appends the path to pre script
func AddPreScript(path string) {
	if path == "" {
		return
	}
	preScript = append(preScript, path)
}

// AddMainScript appends the path to pre script
func AddMainScript(path string) {
	if path == "" {
		return
	}
	mainScript = append(mainScript, path)
}

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
	return LoopProcess(preScript, outputs, arg...)
}

func PostProcess(outputs map[string]Output, arg ...string) error {
	ctrl.Log.WithName("PostProcess").V(1).Info("loop process")
	return LoopProcess(postScript, outputs, arg...)
}

func MainProcess(outputs map[string]Output, arg ...string) error {
	ctrl.Log.WithName("MainProcess").V(1).Info("loop process")
	return LoopProcess(mainScript, outputs, arg...)
}

func Execute(arg ...string) error {
	log := ctrl.Log.WithName("Execute")
	outputs := make(map[string]Output)

	// do pre-process
	if err := PreProcess(outputs, arg...); err != nil {
		log.Error(err, "pre process failed")
		return NewError(err, 550, "pre process failed")
	}

	// do main process
	if err := MainProcess(outputs, arg...); err != nil {
		log.Error(err, "main process failed")
		return NewError(err, 551, "main process failed")
	}

	// do post-process
	if err := PostProcess(outputs, arg...); err != nil {
		log.Error(err, "post process failed")
		return NewError(err, 552, "post process failed")
	}

	return nil
}

func Reset() {
	preScript = []string{}
	mainScript = []string{}
	postScript = []string{}
}

func Validate() bool {
	log := ctrl.Log.WithName("Validate")
	log.V(1).Info("contain", "pre_script", preScript,
		"main_script", mainScript,
		"post_script", postScript)
	return len(mainScript) != 0
}
