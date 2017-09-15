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

package color

import (
	"testing"
)

func TestColor(t *testing.T) {
	t.Log(Red("Red"))
	t.Log(Green("Green"))
	t.Log(Cyan("Cyan"))
	t.Log(Blue("Blue"))
	t.Log(Yellow("Yellow"))
	t.Log(Magenta("Magenta"))
	t.Log(Black("Black"))
	t.Log(White("White"))
	t.Log(HiRed("HiRed"))
	t.Log(HiGreen("HiGreen"))
	t.Log(HiCyan("HiCyan"))
	t.Log(HiBlue("HiBlue"))
	t.Log(HiYellow("HiYellow"))
	t.Log(HiMagenta("HiMagenta"))
	t.Log(HiBlack("HiBlack"))
	t.Log(HiWhite("HiWhite"))
}
