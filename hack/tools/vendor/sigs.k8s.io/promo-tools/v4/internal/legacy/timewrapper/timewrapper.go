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

package timewrapper

import (
	"time"
)

// Time groups the functions mocked by FakeTime.
type Time interface {
	Now() time.Time
	Sleep(d time.Duration)
}

// RealTime is a wrapper for actual time functions.
type RealTime struct{}

// Now simply calls time.Now()
func (RealTime) Now() time.Time { return time.Now() }

// Sleep simply calls time.Sleep(d), using the given duration.
func (RealTime) Sleep(d time.Duration) { time.Sleep(d) }

// FakeTime holds the global fake time.
type FakeTime struct {
	Time time.Time
}

// Now returns the global fake time.
func (ft *FakeTime) Now() time.Time { return ft.Time }

// Sleep adds the given duration to the global fake time.
func (ft *FakeTime) Sleep(d time.Duration) {
	ft.Time = ft.Time.Add(d)
}
