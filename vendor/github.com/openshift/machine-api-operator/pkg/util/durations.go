/*
Copyright 2021 Red Hat.

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

package util

import (
	"time"
)

// The default durations for the leader election operations.
const (
	// LeaseDuration is the default duration for the leader election lease.
	LeaseDuration = 137 * time.Second
	// RenewDeadline is the default duration for the leader renewal.
	RenewDeadline = 107 * time.Second
	// RetryPeriod is the default duration for the leader election retrial.
	RetryPeriod = 26 * time.Second
)

// TimeDuration returns a pointer to the time.Duration.
func TimeDuration(i time.Duration) *time.Duration {
	return &i
}
