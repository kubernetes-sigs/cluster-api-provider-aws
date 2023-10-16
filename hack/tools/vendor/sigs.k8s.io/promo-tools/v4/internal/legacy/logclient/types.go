/*
Copyright 2020 The Kubernetes Authors.

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

package logclient

import (
	"io"
	"log"
)

// These constants refer to the logging levels.
const (
	IndexLogInfo = iota
	IndexLogError
	IndexLogAlert
)

// GetLoggers extracts 3 loggers, corresponding to the logging levels defined
// above.
type GetLoggers interface {
	GetInfoLogger() *log.Logger
	GetErrorLogger() *log.Logger
	GetAlertLogger() *log.Logger
}

// LoggingFacility bundles 3 loggers together.
type LoggingFacility interface {
	GetLoggers
	io.Closer
}
