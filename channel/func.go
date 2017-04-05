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

package channel

import "github.com/jennal/goplay/pkg"
import "strings"

func IsChannel(header *pkg.Header) bool {
	if header.Type&pkg.PKG_PUSH != pkg.PKG_PUSH {
		return false
	}

	if !strings.HasPrefix(header.Route, CHANNEL_PREFIX) {
		return false
	}

	return true
}

func GetChannelName(route string) string {
	return route[len(CHANNEL_PREFIX):]
}
