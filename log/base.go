// Copyright (C) 2017 Jennal(jennalcn@gmail.com). All rights reserved.
//
// Licensed under the MIT License (the "License"); you may not use this file except
// in compliance with the License. You may obtain a copy of the License at
//
// http://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

//Package log is a wrapper to system log package
package log

import (
	"errors"
	"fmt"
	l "log"
	"os"
	"runtime"
	"strings"

	"github.com/jennal/goplay/log/color"
)

const (
	SKIP = 5
)

var (
	INFO  = color.Info("[Info] ")
	TRACE = color.Trace("[Trace] ")
	ERROR = color.Error("[Error] ")
)

type Logger interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Error(err error)
	Errorf(format string, args ...interface{})
	NewErrorf(format string, args ...interface{}) error
	NewError(args ...interface{}) error

	RecoverErrorf(format string, args ...interface{})
}

func setStdout() {
	l.SetFlags(l.LstdFlags | l.Lshortfile)
	l.SetOutput(os.Stdout)
}

func setStderr() {
	l.SetFlags(l.LstdFlags | l.Lshortfile)
	l.SetOutput(os.Stderr)
}

type _logger struct {
	prefix string
	depth  int
}

func NewLogger(prefix string, depth int) Logger {
	if len(prefix) > 0 {
		prefix = prefix + " "
	}

	return _logger{
		prefix: prefix,
		depth:  depth,
	}
}

func (logger _logger) Log(args ...interface{}) {
	setStdout()
	line := fmt.Sprint(args...)
	l.Output(logger.depth, INFO+logger.prefix+line)
}

func (logger _logger) Logf(format string, args ...interface{}) {
	setStdout()
	l.Output(logger.depth, INFO+logger.prefix+fmt.Sprintf(format, args...))
}

func (logger _logger) Trace(args ...interface{}) {
	setStdout()
	line := fmt.Sprint(args...)
	l.Output(logger.depth, TRACE+logger.prefix+line+"\n"+GetStack(5))
}

func (logger _logger) Tracef(format string, args ...interface{}) {
	setStdout()
	l.Output(logger.depth, TRACE+logger.prefix+fmt.Sprintf(format, args...)+"\n"+GetStack(5))
}

func (logger _logger) Error(err error) {
	setStderr()
	l.Output(logger.depth, ERROR+logger.prefix+color.HiRed(err.Error())+"\n"+GetStack(5))
}

func (logger _logger) Errorf(format string, args ...interface{}) {
	setStderr()
	l.Output(logger.depth, ERROR+logger.prefix+color.HiRed(fmt.Sprintf(format, args...))+"\n"+GetStack(5))
}

func (logger _logger) NewErrorf(format string, args ...interface{}) error {
	setStderr()
	err := fmt.Errorf(format, args...)
	l.Output(logger.depth, ERROR+logger.prefix+color.HiRed(err.Error())+"\n"+GetStack(5))

	return err
}

func (logger _logger) NewError(args ...interface{}) error {
	setStderr()
	msg := fmt.Sprint(args...)
	err := errors.New(msg)
	l.Output(logger.depth, ERROR+logger.prefix+color.HiRed(err.Error())+"\n"+GetStack(5))

	return err
}

func (logger _logger) RecoverErrorf(format string, args ...interface{}) {
	setStderr()
	l.Output(logger.depth, ERROR+logger.prefix+"recover from panic: "+color.HiRed(fmt.Sprintf(format, args...))+"\n"+GetStack(9))
}

func GetStack(skip int) string {
	gopath := os.Getenv("GOPATH") + "/src/"
	gopath = strings.Replace(gopath, "\\", "/", -1) // fix windows slash
	result := ""
	pc := make([]uintptr, 10)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return result
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	// Loop to get frames.
	// A fixed number of pcs can expand to an indefinite number of Frames.
	for {
		frame, more := frames.Next()
		if strings.HasPrefix(frame.Function, "runtime.") &&
			!strings.HasPrefix(frame.Function, "runtime.call") &&
			!strings.HasPrefix(frame.Function, "runtime.gopanic") &&
			strings.Contains(frame.File, "runtime/") {
			break
		}

		filename := color.HiMagenta(strings.TrimPrefix(frame.File, gopath))
		line := color.HiGreen(fmt.Sprint(frame.Line))
		function := color.HiBlue(frame.Function)
		result += fmt.Sprintf("\t=> (%v:%v) %v\n", filename, line, function)
		if !more {
			break
		}
	}

	return result
}
