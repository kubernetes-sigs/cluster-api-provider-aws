/*
Copyright 2021 The Kubernetes Authors.

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

package signals

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	// ExitSignals are used to determine if an incoming os.Signal should cause termination.
	ExitSignals = map[os.Signal]bool{
		syscall.SIGHUP:  true,
		syscall.SIGINT:  true,
		syscall.SIGABRT: true,
		syscall.SIGILL:  true,
		syscall.SIGQUIT: true,
		syscall.SIGTERM: true,
		syscall.SIGSEGV: true,
		syscall.SIGTSTP: true,
	}
	// ExitChannel is for gracefully terminating the LogSignals() function.
	ExitChannel = make(chan bool, 1)
	// SignalChannel is for listening to OS signals.
	SignalChannel = make(chan os.Signal, 1)
	// Debug is defined globally for mocking logrus.Debug statements.
	Debug func(args ...interface{}) = logrus.Debug
)

// Watch concurrently logs debug statements when encountering interrupt signals from the OS.
// This function relies on the os/signal package which may not capture SIGKILL and SIGSTOP.
func Watch() {
	Debug("Watching for OS Signals...")
	// Observe all signals, excluding SIGKILL and SIGSTOP.
	signal.Notify(SignalChannel)
	// Continuously log signals.
	go LogSignals()
}

// Stop gracefully terminates the concurrent signal logging.
// TODO: @tylerferrara Currently we don't gracefully exit, since the Auditor is designed
// to run indefinitely. In the future, we should enable graceful termination to allow
// this to be called from a Shutdown() function.
func Stop() {
	ExitChannel <- true
}

// LogSignals continuously prints logging statements for each signal it observes,
// exiting only when observing an exit signal.
func LogSignals() {
	for {
		select {
		case sig := <-SignalChannel:
			LogSignal(sig)
		case <-ExitChannel:
			// Gracefully terminate.
			return
		}
	}
}

// LogSignal prints a logging statements for the given signal,
// exiting only when observing an exit signal.
func LogSignal(sig os.Signal) {
	Debug("Encoutered signal: ", sig.String())
	// Handle exit signals.
	if _, found := ExitSignals[sig]; found {
		Debug("Exiting from signal: ", sig.String())
		// If we get here, an exit signal was seen. We must handle this by forcing
		// the program to exit. Without this, the program would ignore all exit signals
		// except SIGKILL.
		Stop()
	}
}
