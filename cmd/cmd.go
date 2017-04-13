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

package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jennal/goplay/log"
	"github.com/jennal/goplay/transfer"
)

func Start(s transfer.IServer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGABRT, syscall.SIGBUS, syscall.SIGFPE, syscall.SIGSEGV, syscall.SIGTERM)

	err := s.Start()
	if err != nil {
		log.Error(err)
		return
	}

	<-c
	err = s.Stop()
	if err != nil {
		log.Error(err)
	}
}
