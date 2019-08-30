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
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/awserrors"
)

// PollWithRetryable repeats a condition check until a timeout.
func PollWithRetryable(condition wait.ConditionFunc, retryableErrors ...string) error { //nolint
	return wait.PollImmediate(10*time.Second, 10*time.Minute, func() (bool, error) {
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
