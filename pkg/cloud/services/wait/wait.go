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

package wait

import (
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
)

/*
 Ideally, this entire file would be replaced with returning a retryable
 error and letting the actuator requeue deletion. Unfortunately, since
 the retry behaviour is not tunable, with a max retry limit of 10, we
 implement waits manually here.
*/

// NewBackoff creates a new API Machinery backoff parameter set suitable
// for use with AWS services.
func NewBackoff() wait.Backoff {
	// Return a exponential backoff configuration which
	// returns durations for a total time of ~10m.
	// Example: 1s, 2s, 4s, 8s, 16s, 20s, ... 20s — for a total of N steps.
	return wait.Backoff{
		Duration: time.Second,
		Factor:   2,
		Steps:    32,
		Jitter:   4,
		Cap:      20 * time.Second,
	}
}

// WaitForWithRetryable repeats a condition check with exponential backoff.
func WaitForWithRetryable(backoff wait.Backoff, condition wait.ConditionFunc, retryableErrors ...string) error { //nolint
	return wait.ExponentialBackoff(backoff, func() (bool, error) {
		ok, err := condition()
		if ok {
			// All done!
			return true, nil
		}
		if err == nil {
			// Not done, but no error, so keep waiting.
			return false, nil
		}

		// If the returned error isn't empty, check if the error is a retryable one,
		// or return immediately.
		code, ok := awserrors.Code(errors.Cause(err))
		if !ok {
			return false, err
		}

		for _, r := range retryableErrors {
			if code == r {
				// We should retry.
				return false, nil
			}
		}

		// Got an error that we can't retry, so return it.
		return false, err
	})
}
