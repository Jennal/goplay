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

// +build !windows

//Package color makes logs colorful
package color

import "github.com/fatih/color"

var Info = color.New(color.Cyan, color.Bold).SprintFunc()
var Trace = color.New(color.Yellow, color.Bold).SprintFunc()
var Error = color.New(color.Red, color.Bold).SprintFunc()

var Red = color.New(color.FgRed).SprintFunc()
var Green = color.New(color.FgGreen).SprintFunc()
var Cyan = color.New(color.FgCyan).SprintFunc()
var Blue = color.New(color.FgBlue).SprintFunc()
var Yellow = color.New(color.FgYellow).SprintFunc()
var Magenta = color.New(color.FgMagenta).SprintFunc()
var Black = color.New(color.FgBlack).SprintFunc()
var White = color.New(color.FgWhite).SprintFunc()

var HiRed = color.New(color.FgHiRed).SprintFunc()
var HiGreen = color.New(color.FgHiGreen).SprintFunc()
var HiCyan = color.New(color.FgHiCyan).SprintFunc()
var HiBlue = color.New(color.FgHiBlue).SprintFunc()
var HiYellow = color.New(color.FgHiYellow).SprintFunc()
var HiMagenta = color.New(color.FgHiMagenta).SprintFunc()
var HiBlack = color.New(color.FgHiBlack).SprintFunc()
var HiWhite = color.New(color.FgHiWhite).SprintFunc()
