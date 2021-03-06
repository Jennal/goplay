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

// +build windows

package color

func echo(str string) string {
	return str
}

var Info = echo
var Trace = echo
var Error = echo

var Red = echo
var Green = echo
var Cyan = echo
var Blue = echo
var Yellow = echo
var Magenta = echo
var Black = echo
var White = echo

var HiRed = echo
var HiGreen = echo
var HiCyan = echo
var HiBlue = echo
var HiYellow = echo
var HiMagenta = echo
var HiBlack = echo
var HiWhite = echo
