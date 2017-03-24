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

package log

var defaultLogger = NewLogger("", 4)

//Log logs normal info with a string
func Log(args ...interface{}) {
	defaultLogger.Log(args)
}

//Logf logs normal info with a string format
func Logf(format string, args ...interface{}) {
	defaultLogger.Logf(format, args...)
}

//Error logs error info with an error
func Error(err error) {
	defaultLogger.Error(err)
}

//Errorf logs error info with a string format
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

//NewErrorf logs error info with a string format and return an error with that message
func NewErrorf(format string, args ...interface{}) error {
	return defaultLogger.NewErrorf(format, args...)
}

//NewError logs error info with a string and return an error with that message
func NewError(args ...interface{}) error {
	return defaultLogger.NewError(args)
}
