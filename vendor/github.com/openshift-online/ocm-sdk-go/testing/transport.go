/*
Copyright (c) 2021 Red Hat, Inc.

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

package testing

import (
	"io"
	"net/http"
	"strings"

	. "github.com/onsi/gomega" // nolint
)

// TransportFunc is a function that implements the http.RoundTripper interface. This is intended to
// siplify writing tests that require custom round trippers because it doesn't require a new type.
// For example, to create a round tripper that always returns an error:
//
//	transport := testing.TransportFunc(func (*http.Request) (*http.Response, error) {
//		return nil, errors.New("my error")
//	})
type TransportFunc func(*http.Request) (*http.Response, error)

// RoundTrip is the implementation of the http.RoundTripper interface.
func (f TransportFunc) RoundTrip(request *http.Request) (response *http.Response, err error) {
	response, err = f(request)
	return
}

// ErrorTransport creates a transport that always returns the given error.
func ErrorTransport(err error) http.RoundTripper {
	return TransportFunc(func(*http.Request) (*http.Response, error) {
		return nil, err
	})
}

// TextTransport creates a transport that always returns the given status code and plan text body.
func TextTransport(code int, body string) http.RoundTripper {
	return TransportFunc(func(request *http.Request) (response *http.Response, err error) {
		response = &http.Response{
			StatusCode: code,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header: http.Header{
				"Content-Type": []string{"plain/text"},
			},
		}
		return
	})
}

// JSONTransport creates a transport that always returns the given status code and JSON body.
func JSONTransport(code int, body string) http.RoundTripper {
	return TransportFunc(func(request *http.Request) (response *http.Response, err error) {
		response = &http.Response{
			StatusCode: code,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
		}
		return
	})
}

// CombineTransports returns a transport that delegates to the given array of transports, in the
// given order. First request will go to the first transport, sceond request to the second transport,
// so on.
func CombineTransports(transports ...http.RoundTripper) http.RoundTripper {
	i := 0
	return TransportFunc(func(request *http.Request) (response *http.Response, err error) {
		Expect(i).To(BeNumerically("<", len(transports)))
		response, err = transports[i].RoundTrip(request)
		i++
		return
	})
}
