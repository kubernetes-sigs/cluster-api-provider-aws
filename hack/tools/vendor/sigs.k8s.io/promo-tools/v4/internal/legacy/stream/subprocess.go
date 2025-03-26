/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package stream

import (
	"io"
	"os/exec"
)

// Subprocess can spawn a subprocess and read from it. It can be used to read
// from an io.Reader that produces JSON, or whatever else.
type Subprocess struct {
	CmdInvocation []string
	cmd           *exec.Cmd
}

// Produce runs the external process and returns two io.Readers (to stdout and
// stderr).
func (sp *Subprocess) Produce() (stdOut, stdErr io.Reader, err error) {
	invocation := sp.CmdInvocation
	cmd := exec.Command(invocation[0], invocation[1:]...)
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	sp.cmd = cmd
	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	return stdoutReader, stderrReader, nil
}

// Close the subprocess by waiting for it.
func (sp *Subprocess) Close() error {
	// The call to Wait() _cannot_ come before all reads from the pipe have been
	// completed. Otherwise, we may end up closing the pipe before the JSON
	// decoder has had time to read from it.
	// See https://golang.org/pkg/os/exec/#Cmd.StdoutPipe.
	return sp.cmd.Wait()
}
