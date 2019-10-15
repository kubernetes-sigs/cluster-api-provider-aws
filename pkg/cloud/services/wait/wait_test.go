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

package wait

import (
	"errors"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

func TestWaitForWithRetryable(t *testing.T) {
	backoff := wait.Backoff{
		Duration: 1 * time.Millisecond,
		Factor:   0,
		Jitter:   0,
		Steps:    1,
		Cap:      2 * time.Millisecond,
	}

	tests := []struct {
		name            string
		backoff         wait.Backoff
		condition       wait.ConditionFunc
		retryableErrors []string
		expectedError   error
	}{
		{
			name:    "return nil",
			backoff: backoff,
			condition: func() (done bool, err error) {
				return true, nil
			},
			retryableErrors: nil,
			expectedError:   nil,
		},
		{
			name:    "returns timeout error",
			backoff: backoff,
			condition: func() (done bool, err error) {
				return false, nil
			},
			retryableErrors: nil,
			expectedError:   wait.ErrWaitTimeout,
		},
		{
			name:    "error occurred in conditionFunc, returns actual error",
			backoff: backoff,
			condition: func() (done bool, err error) {
				return false, errors.New("new error")
			},
			retryableErrors: nil,
			expectedError:   errors.New("new error"),
		},
		{
			name:    "timed out in retryable error, returns actual error",
			backoff: backoff,
			condition: func() (done bool, err error) {
				return false, errors.New("retryable error")
			},
			retryableErrors: []string{"retryable error"},
			expectedError:   errors.New("retryable error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WaitForWithRetryable(tt.backoff, tt.condition, tt.retryableErrors...)

			var expected, actual string
			if err != nil {
				actual = err.Error()
			}
			if tt.expectedError != nil {
				expected = tt.expectedError.Error()
			}

			if expected != actual {
				t.Errorf("expected error: %v, got error: %v", expected, actual)
			}
		})
	}
}
