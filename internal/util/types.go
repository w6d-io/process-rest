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
package util

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/w6d-io/appdeploy/internal/config"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// LEVEL

// levelStrings contains level string supported
var levelStrings = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"error": zapcore.ErrorLevel,
}

// LevelFlag contains structure for managing zap level
type LevelFlag struct {
	ZapOptions *zap.Options
	value      string
}

func (l LevelFlag) String() string {
	return l.value
}

func (l LevelFlag) Set(flagValue string) error {
	if flagValue == "" {
		return errors.New("log-level cannot be empty")
	}
	level, validLevel := levelStrings[strings.ToLower(flagValue)]
	if !validLevel {
		logLevel, err := strconv.Atoi(flagValue)
		if err != nil {
			return fmt.Errorf("invalid log level \"%s\"", flagValue)
		}
		if logLevel > 0 {
			intLevel := -1 * logLevel
			l.ZapOptions.Level = zapcore.Level(int8(intLevel))
		} else {
			return fmt.Errorf("invalid log level \"%s\"", flagValue)
		}
	} else {
		l.ZapOptions.Level = level
	}
	l.value = flagValue
	return nil
}

// OUTPUT FORMAT

// JsonEncoderConfig returns an opinionated EncoderConfig
func JsonEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// TextEncoderConfig returns an opinionated EncoderConfig
func TextEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// OutputFormatFlag contains structure for managing zap encoding
type OutputFormatFlag struct {
	ZapOptions *zap.Options
	value      string
}

func (o *OutputFormatFlag) String() string {
	return o.value
}

func (o *OutputFormatFlag) Set(flagValue string) error {
	if flagValue == "" {
		return errors.New("log-format cannot be empty")
	}
	val := strings.ToLower(flagValue)
	switch val {
	case "json":
		o.ZapOptions.Encoder = zapcore.NewJSONEncoder(JsonEncoderConfig())
	case "text":
		o.ZapOptions.Encoder = zapcore.NewConsoleEncoder(TextEncoderConfig())
	default:
		return fmt.Errorf(`invalid "%s"`, flagValue)
	}
	o.value = flagValue
	return nil
}

// CONFIG PART

type ConfigFlag struct {
	value string
}

func (f ConfigFlag) String() string {
	return f.value
}

func (f ConfigFlag) Set(flagValue string) error {
	if flagValue == "" {
		return errors.New("config cannot be empty")
	}
	isFileExists := func(filename string) bool {
		info, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return false
		}
		return !info.IsDir()
	}
	if !isFileExists(flagValue) {
		return fmt.Errorf("file %s does not exist", flagValue)
	}
	if err := config.New(flagValue); err != nil {
		return fmt.Errorf("instantiate config returns %s", err)
	}
	f.value = flagValue
	return nil
}

// INIT PART

// BindFlags custom flags
func BindFlags(o *zap.Options, fs *flag.FlagSet) {

	var outputFormat OutputFormatFlag
	outputFormat.ZapOptions = o
	fs.Var(&outputFormat, "log-format", "log encoding ( 'json' or 'text')")

	var level LevelFlag
	level.ZapOptions = o
	fs.Var(&level, "log-level", "log level verbosity. Can be 'debug', 'info', 'error', "+
		"or any integer value > 0 which corresponds to custom debug levels of increasing verbosity")

	// TODO add auth file/config ??

	var c ConfigFlag
	fs.Var(&c, "config", "config file")
}
