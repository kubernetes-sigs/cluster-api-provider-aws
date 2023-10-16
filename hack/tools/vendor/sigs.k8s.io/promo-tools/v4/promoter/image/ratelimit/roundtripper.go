/*
Copyright 2023 The Kubernetes Authors.

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

package ratelimit

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"
)

const (
	burst = 1

	// Number of events to pass to the crane rate limiter to match the
	// AR api limits:
	// https://github.com/kubernetes/registry.k8s.io/issues/153#issuecomment-1460913153
	// bentheelder: (83*60=4980)
	// Temp: We are temporarily dropping this from the max of 83 to
	// two thirds while the rate limiter is instrumented everywhere.
	MaxEvents = 50
)

// RoundTripper wraps an http.RoundTripper with rate limiting
type RoundTripper struct {
	rateLimiter  *rate.Limiter
	roundTripper http.RoundTripper
}

var _ http.RoundTripper = &RoundTripper{}

var Limiter *RoundTripper

func init() {
	if Limiter == nil {
		Limiter = NewRoundTripper(MaxEvents)
	}
}

func NewRoundTripper(limit rate.Limit) *RoundTripper {
	return &RoundTripper{
		rateLimiter:  rate.NewLimiter(limit, burst),
		roundTripper: http.DefaultTransport,
	}
}

func (rt *RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// only rate limit read type calls, not writes
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		err := rt.rateLimiter.Wait(context.Background())
		if err != nil {
			return nil, err
		}
	}

	return rt.roundTripper.RoundTrip(r)
}
