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

package wait_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"k8s.io/apimachinery/pkg/util/wait"

	. "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/wait"
)

var (
	errRetryable        = awserr.New("retryable error", "", nil)
	errNonRetryable     = errors.New("non retryable error")
	retryableErrorCodes = []string{"retryable error"}
)

func TestWaitForWithRetryable(t *testing.T) {
	backoff := wait.Backoff{
		Duration: 1 * time.Millisecond,
		Factor:   0,
		Jitter:   0,
		Steps:    2,
		Cap:      3 * time.Millisecond,
	}
	var firstIteration bool

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
			retryableErrors: retryableErrorCodes,
			expectedError:   nil,
		},
		{
			name:    "returns timeout error",
			backoff: backoff,
			condition: func() (done bool, err error) {
				return false, nil
			},
			retryableErrors: retryableErrorCodes,
			expectedError:   wait.ErrWaitTimeout,
		},
		{
			name:    "error occurred in conditionFunc, returns actual error",
			backoff: backoff,
			condition: func() (done bool, err error) {
				return false, errNonRetryable
			},
			retryableErrors: retryableErrorCodes,
			expectedError:   errNonRetryable,
		},
		{
			name:    "timed out in retryable error, returns the retryable error",
			backoff: backoff,
			condition: func() (done bool, err error) {
				return false, errRetryable
			},
			retryableErrors: retryableErrorCodes,
			expectedError:   errRetryable,
		},
		{
			name:    "retryable err at first, non-retryable err after that, returns latest error",
			backoff: backoff,
			condition: func() (done bool, err error) {
				if firstIteration {
					firstIteration = false
					return false, errRetryable
				}
				firstIteration = false
				return false, errNonRetryable
			},
			retryableErrors: retryableErrorCodes,
			expectedError:   errNonRetryable,
		},
		{
			name:    "retryable err at first, failure but no error after that, returns timeout error",
			backoff: backoff,
			condition: func() (done bool, err error) {
				if firstIteration {
					firstIteration = false
					return false, errRetryable
				}
				return false, nil
			},
			retryableErrors: retryableErrorCodes,
			expectedError:   wait.ErrWaitTimeout,
		},
		{
			name:    "retryable error at first, success after that, returns nil",
			backoff: backoff,
			condition: func() (done bool, err error) {
				if firstIteration {
					firstIteration = false
					return false, errRetryable
				}
				return true, nil
			},
			retryableErrors: retryableErrorCodes,
			expectedError:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstIteration = true

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
