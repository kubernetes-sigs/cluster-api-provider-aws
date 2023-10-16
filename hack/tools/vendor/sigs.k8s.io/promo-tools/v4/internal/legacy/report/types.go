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

package report

import (
	"io"

	"cloud.google.com/go/errorreporting"
)

// Reporter requires a single method, called Report(), which corresponds to
// calling Stackdriver Error Reporting for the real implementation.
type Reporter interface {
	Report(errorreporting.Entry)
}

// ReportingFacility has a Reporter and Closer. Unlike LoggingFacility, there is
// no need to have a struct type because the same thing is used to call Report()
// and Close() on. This is because the real implementation calls both Report()
// and Close() methods on the same errorreporting.Client type (and so, the fake
// implementation follows suit and does the same). As such, there is no need to
// have a wrapper struct around the two (separate) interfaces, and
// ReportingFacility is just a plain interface.
type ReportingFacility interface {
	Reporter
	io.Closer
}
