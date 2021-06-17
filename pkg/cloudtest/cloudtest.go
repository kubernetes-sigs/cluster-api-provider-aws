/*
Copyright 2018 The Kubernetes Authors.

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

package cloudtest

import (
	"encoding/json"
	"testing"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
)

// RuntimeRawExtension takes anything and turns it into a *runtime.RawExtension.
// This is helpful for creating clusterv1.Cluster/Machine objects that need
// a specific AWSClusterProviderSpec or Status.
func RuntimeRawExtension(t *testing.T, p interface{}) *runtime.RawExtension {
	t.Helper()
	out, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	return &runtime.RawExtension{
		Raw: out,
	}
}

// Log implements logr.Logger for testing. Do not use if you actually want to
// test log messages.
type Log struct{}

// Error implements Log errors.
func (l *Log) Error(err error, msg string, keysAndValues ...interface{}) {}

// V returns the Logger's log level.
func (l *Log) V(level int) logr.Logger { return l }

// WithValues returns logs with specific values.
func (l *Log) WithValues(keysAndValues ...interface{}) logr.Logger { return l }

// WithName returns the logger with a specific name.
func (l *Log) WithName(name string) logr.Logger { return l }

// Info implements info messages for the logger.
func (l *Log) Info(msg string, keysAndValues ...interface{}) {}

// Enabled returns the state of the logger.
func (l *Log) Enabled() bool { return false }
