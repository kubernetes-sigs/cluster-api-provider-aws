// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"os"
	"os/exec"
)

func runShellWithWait(cmd string) {
	runCommandWithWait("sh", "-c", cmd)
}

func runCommandWithWait(cmd string, args ...string) bool {
	command := runCommand(cmd, args...)
	if err := command.Wait(); err != nil {
		log.Println(err)
	}
	return command.ProcessState.Success()
}

func runShell(cmd string) *exec.Cmd {
	return runCommand("sh", "-c", cmd)
}

func runCommand(cmd string, args ...string) *exec.Cmd {
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Start(); err != nil {
		log.Println(err)
	}
	return command
}
