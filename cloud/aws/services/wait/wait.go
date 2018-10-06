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

	"github.com/aws/aws-sdk-go/aws/awserr"
	aMW "k8s.io/apimachinery/pkg/util/wait"
)

/*
 Ideally, this entire file would be replaced with returning a retryable
 error and letting the actuator requeue deletion. Unfortunately, since
 the retry behaviour is not tunable, with a max retry limit of 10, we
 implement waits manually here.
*/

// NewBackoff creates a new API Machinery backoff parameter set suitable
// for use with AWS services, with values based loosely on
// https://github.com/Netflix/edda/blob/master/src/main/scala/com/netflix/edda/Crawler.scala#L159
func NewBackoff() aMW.Backoff {

	duration, _ := time.ParseDuration("2s")
	return aMW.Backoff{
		Duration: duration,
		Factor:   1.5,
		Jitter:   1.0,
		Steps:    100,
	}
}

// WaitForWithRetryable repeats a condition check with exponential backoff.
//
// It takes a list of string slice of AWS API errors that can be considered as retriable.
//
// It checks the condition up to Steps times.
//
// If Jitter is greater than zero, a random amount of each duration is added
// (between duration and duration*(1+jitter)).
//
// If it receives a RequestLimitExceeded or Throttling error, then we
// multiply the delay by the specified factor.
//
// If the condition never returns true, ErrWaitTimeout is returned. All other
// errors terminate immediately.
func WaitForWithRetryable(backoff aMW.Backoff, condition aMW.ConditionFunc, retryableErrors []string) error {
	duration := backoff.Duration
	for i := 0; i < backoff.Steps; i++ {
		if i != 0 {
			adjusted := duration
			if backoff.Jitter > 0.0 {
				adjusted = aMW.Jitter(duration, backoff.Jitter)
			}
			time.Sleep(adjusted)
			duration = time.Duration(float64(duration) * backoff.Factor)
		}
		ok, err := condition()
		if ok {
			return nil
		}
		if err != nil {
			awserr, ok := err.(awserr.Error)
			if !ok {
				return err
			}
			isRetryable := false
			for _, r := range retryableErrors {
				if awserr.Code() == r {
					isRetryable = true
				}
			}
			if !isRetryable {
				return err
			}
		}
	}
	return aMW.ErrWaitTimeout
}
