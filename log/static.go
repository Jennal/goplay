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

var defaultLogger Logger = NewLogger("", 4)

func Log(line string) {
	defaultLogger.Log(line)
}

func Logf(format string, args ...interface{}) {
	defaultLogger.Logf(format, args...)
}

func Error(err error) {
	defaultLogger.Error(err)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func NewErrorf(format string, args ...interface{}) error {
	return defaultLogger.NewErrorf(format, args...)
}

func NewError(msg string) error {
	return defaultLogger.NewError(msg)
}
