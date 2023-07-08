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

package stream

import (
	"io"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// Producer is an interface for anything that can generate an io.Reader from
// which we can read from (typically JSON output).
type Producer interface {
	// The first two io.Readers are expected to be the stdout and stderr streams
	// from the process, respectively.
	Produce() (io.Reader, io.Reader, error)
	Close() error
}

// An ExternalRequest is anything that can create and then consume any stream.
// The request comes bundled with something that can produce a stream
// (io.Reader), and something that can read from that stream to populate some
// arbitrary data structure.
type ExternalRequest struct {
	RequestParams  interface{}
	StreamProducer Producer
}

// BackoffDefault is the default Backoff behavior for network call retries.
//
// Previous values from k8s.io/apimachinery/pkg/util/wait.`Backoff`:
// - Duration: time.Second
// - Factor:   2
// - Jitter:   0.1
// - Steps:    45
// - Cap:      time.Second * 60
func BackoffDefault() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = time.Second
	b.Multiplier = 2
	b.RandomizationFactor = 0.1
	b.MaxElapsedTime = time.Second * 60

	return b
}
