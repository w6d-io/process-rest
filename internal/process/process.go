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
	"context"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/w6d-io/hook"
	"github.com/w6d-io/process-rest/internal/config"
	"github.com/w6d-io/x/logx"
)

func Run(name string, arg ...string) (string, error) {
	log := logx.WithName(nil, "Process.Run")
	log.V(1).Info("build command")
	cmd := exec.Command(name, arg...)

	log.V(1).Info("exec command and get output", "script", cmd.String())
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Error(err, "script failed", "script", cmd.String(), "stdout", string(output), "exit_code", exitErr.ExitCode(), "stderr", string(exitErr.Stderr))
			return string(exitErr.Stderr), exitErr
		}
		log.Error(err, "script failed", "script", cmd.String(), "output", string(output))
		return string(output), err
	}
	log.V(1).Info("script succeeded", "output", string(output))
	return string(output), nil
}

func (p *Process) LoopProcess(scripts []string, arg ...string) error {
	log := logx.WithName(nil, "Process.LoopProcess")
	for _, script := range scripts {
		log.Info("run", "script", script)
		arg = append([]string{script}, arg...)
		args := strings.Join(arg, " ")
		output, err := Run("bash", "-c", args)
		o := Output{
			Name:   path.Base(script),
			Status: "succeeded",
			Log:    output,
			Error:  "",
		}
		if err != nil {
			log.Error(err, "process failed", "script", script)
			o.Status = "failed"
			o.Error = err.Error()
			p.Outputs = append(p.Outputs, o)
			return err
		}
		p.Outputs = append(p.Outputs, o)
	}
	return nil
}

func (p *Process) PreProcess(arg ...string) error {
	log := logx.WithName(nil, "Process.PreProcess")
	log.V(1).Info("loop process")
	return p.LoopProcess(config.GetPreScript(), arg...)
}

func (p *Process) PostProcess(arg ...string) error {
	log := logx.WithName(nil, "Process.PostProcess")
	log.V(1).Info("loop process")
	return p.LoopProcess(config.GetPostScript(), arg...)
}

func (p *Process) MainProcess(arg ...string) error {
	log := logx.WithName(nil, "Process.MainProcess")
	log.V(1).Info("loop process")
	return p.LoopProcess(config.GetMainScript(), arg...)
}

func Execute(id string, arg ...string) {
	log := logx.WithName(nil, "Process.Execute")
	log.V(1).Info("loop process")
	p := new(Process)
	errc := make(chan error)
	go func() {
		// do pre-process
		if err := p.PreProcess(arg...); err != nil {
			log.Error(err, "pre process failed")
			p.Notify(id, "pre-process-failed", err)
			errc <- NewError(err, 550, "pre process failed")
			return
		}

		// do main process
		if err := p.MainProcess(arg...); err != nil {
			log.Error(err, "main process failed")
			p.Notify(id, "main-process-failed", err)
			errc <- NewError(err, 551, "main process failed")
			return
		}

		// do post-process
		if err := p.PostProcess(arg...); err != nil {
			log.Error(err, "post process failed")
			p.Notify(id, "post-process-failed", err)
			errc <- NewError(err, 552, "post process failed")
			return
		}
		p.Notify(id, "process-succeeded", nil)
		errc <- nil
	}()

	//for range []string{"1", "2"} {
	if err := <-errc; err != nil {
		log.Error(err, "process failed")
	}
	//}
}

func (p *Process) Notify(id string, scope string, err error) {
	log := logx.WithName(nil, "Process.Notify")

	status := &Status{
		Success: err == nil,
		Log:     p.GetLogMessage(err),
		ID:      id,
	}
	log.V(1).Info("send", "scope", scope)
	_ = hook.Send(context.Background(), status, scope)
}

func (p *Process) GetLogMessage(err error) string {
	var messages []string
	if err != nil {
		messages = append(messages, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
	for _, o := range p.Outputs {
		message := fmt.Sprintf(`{"script":"%s", "error":"%s", "status":"%s", "log":"%s"}`,
			o.Name, o.Error, o.Status, o.Log)
		messages = append(messages, message)
	}
	return fmt.Sprintf("{%s}", strings.Join(messages, ","))
}
