/*
Copyright The Kubernetes Authors.
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

package termination

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/klog/klogr"
)

var _ = Describe("Handler Suite", func() {
	var terminationServer *httptest.Server
	var httpHandler http.Handler
	var nodeName string
	var stop chan struct{}
	var errs chan error
	var h *handler

	BeforeEach(func() {
		// Reset test vars
		terminationServer = nil
		httpHandler = nil
		nodeName = "testNode"

		h = &handler{
			client:       k8sClient,
			pollInterval: 100 * time.Millisecond,
			nodeName:     nodeName,
			log:          klogr.New(),
		}
	})

	JustBeforeEach(func() {
		Expect(httpHandler).ToNot(BeNil())
		terminationServer = httptest.NewServer(httpHandler)

		pollURL, err := url.Parse(terminationServer.URL)
		Expect(err).ToNot(HaveOccurred())
		h.pollURL = pollURL

		stop, errs = StartTestHandler(h)
	})

	AfterEach(func() {
		if !isClosed(stop) {
			close(stop)
		}
		terminationServer.Close()
	})

	Context("when the handler is stopped", func() {
		BeforeEach(func() {
			httpHandler = newMockHTTPHandler(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(404)
			})
		})

		JustBeforeEach(func() {
			close(stop)
		})

		It("should not return an error", func() {
			Eventually(errs).Should(Receive(BeNil()))
		})
	})
})

// mockHTTPHandler is used to mock the pollURL responses during tests
type mockHTTPHandler struct {
	handleFunc func(rw http.ResponseWriter, req *http.Request)
}

// ServeHTTP implements the http.Handler interface
func (m *mockHTTPHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	m.handleFunc(rw, req)
}

// newMockHTTPHandler constructs a mockHTTPHandler with the given handleFunc
func newMockHTTPHandler(handleFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	return &mockHTTPHandler{handleFunc: handleFunc}
}

// isClosed checks if a channel is closed already
func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}
