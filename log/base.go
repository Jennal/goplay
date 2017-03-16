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
)

type Logger interface {
	Log(line string)
	Logf(format string, args ...interface{})
	Error(err error)
	Errorf(format string, args ...interface{})
	NewErrorf(format string, args ...interface{}) error
	NewError(msg string) error
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

func (logger _logger) Log(line string) {
	setStdout()
	l.Output(logger.depth, logger.prefix+line)
}

func (logger _logger) Logf(format string, args ...interface{}) {
	setStdout()
	l.Output(logger.depth, logger.prefix+fmt.Sprintf(format, args...))
}

func (logger _logger) Error(err error) {
	setStderr()
	l.Output(logger.depth, logger.prefix+err.Error())
}

func (logger _logger) Errorf(format string, args ...interface{}) {
	setStderr()
	l.Output(logger.depth, logger.prefix+fmt.Sprintf(format, args...))
}

func (logger _logger) NewErrorf(format string, args ...interface{}) error {
	setStderr()
	err := fmt.Errorf(format, args...)
	l.Output(logger.depth, logger.prefix+err.Error())
	return err
}

func (logger _logger) NewError(msg string) error {
	setStderr()
	err := errors.New(msg)
	l.Output(logger.depth, logger.prefix+err.Error())
	return err
}
