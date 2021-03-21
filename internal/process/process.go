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
	postScript = append(postScript, path)
}

// AddPreScript appends the path to pre script
func AddPreScript(path string) {
	preScript = append(preScript, path)
}

// AddDeployScript appends the path to pre script
func AddDeployScript(path string) {
	deployScript = append(deployScript, path)
}

func Process(name string, arg ...string) (string, error) {
	log := ctrl.Log.WithName("Process")
	log.V(1).Info("build command")
	cmd := exec.Command(name, arg...)

	log.V(1).Info("exec command and get output", "script", name)
	output, err := cmd.Output()
	if err != nil {
		log.Error(err, "script failed", "script", name)
		return "", err
	}
	log.V(1).Info("script succeeded", "output", string(output))
	return string(output), nil
}

func LoopProcess(scripts []string, outputs map[string]Output, arg ...string) error {
	log := ctrl.Log.WithName("LoopProcess")
	for _, script := range scripts {
		output, err := Process(script, arg...)
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

func PreDeploy(outputs map[string]Output, arg ...string) error {
	return LoopProcess(preScript, outputs, arg...)
}

func PostDeploy(outputs map[string]Output, arg ...string) error {
	return LoopProcess(postScript, outputs, arg...)
}

func Deploy(outputs map[string]Output, arg ...string) error {
	return LoopProcess(deployScript, outputs, arg...)
}

func Execute(arg ...string) error {
	log := ctrl.Log.WithName("Execute")
	outputs := make(map[string]Output)

	// do pre-deploy
	if err := PreDeploy(outputs, arg...); err != nil {
		log.Error(err, "pre deploy failed")
		return NewError(err, 550, "pre deploy failed")
	}

	// do deployment
	if err := Deploy(outputs, arg...); err != nil {
		log.Error(err, "deploy failed")
		return NewError(err, 551, "Deploy failed")
	}

	// do post-deploy
	if err := PostDeploy(outputs, arg...); err != nil {
		log.Error(err, "post deploy failed")
		return NewError(err, 552, "post deploy failed")
	}

	return nil
}

func Reset() {
	preScript = []string{}
	deployScript = []string{}
	postScript = []string{}
}

func Validate() bool {
	log := ctrl.Log.WithName("Validate")
	log.V(1).Info("contain", "pre_script", preScript,
		"deploy_script", deployScript,
		"post_script", postScript)
	return len(deployScript) != 0
}
