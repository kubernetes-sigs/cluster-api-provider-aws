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
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/klogr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var notFoundFunc = func(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(404)
}

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
		nodeName = "test-node"
		httpHandler = newMockHTTPHandler(notFoundFunc)

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

		Expect(deleteAllMachines(k8sClient)).To(Succeed())
	})

	Context("when the handler is stopped", func() {
		JustBeforeEach(func() {
			close(stop)
		})

		It("should not return an error", func() {
			Eventually(errs).Should(Receive(BeNil()))
		})
	})

	Context("when no machine exists for the node", func() {
		It("should return an error upon starting", func() {
			Eventually(errs).Should(Receive(MatchError("error fetching machine for node (\"test-node\"): machine not found for node \"test-node\"")))
		})
	})

	Context("getMachineForNode", func() {
		var machine *machinev1.Machine
		var err error

		JustBeforeEach(func() {
			machine, err = h.getMachineForNode(ctx)
		})

		Context("with a broken client", func() {
			BeforeEach(func() {
				brokenClient, err := client.New(cfg, client.Options{Scheme: runtime.NewScheme()})
				Expect(err).ToNot(HaveOccurred())
				h.client = brokenClient
			})

			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix("error listing machines: no kind is registered for the type v1beta1.MachineList in scheme"))
			})

			It("should not return a machine", func() {
				Expect(machine).To(BeNil())
			})
		})

		Context("with no machine for the node name", func() {
			It("should return an error", func() {
				Expect(err).To(MatchError("machine not found for node \"test-node\""))
			})

			It("should not return a machine", func() {
				Expect(machine).To(BeNil())
			})
		})

		Context("with a machine matching the node name", func() {
			var testMachine *machinev1.Machine

			BeforeEach(func() {
				testMachine = newTestMachine("test-machine", "test-namespace", nodeName)
				createMachine(testMachine)
			})

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return a machine", func() {
				Expect(machine).To(Equal(testMachine))
			})
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

func deleteAllMachines(c client.Client) error {
	machineList := &machinev1.MachineList{}
	err := c.List(ctx, machineList)
	if err != nil {
		return fmt.Errorf("error listing machines: %v", err)
	}

	// Delete all machines found
	for _, machine := range machineList.Items {
		m := machine
		err := c.Delete(ctx, &m)
		if err != nil {
			return err
		}
	}
	return nil
}

func newTestMachine(name, namespace, nodeName string) *machinev1.Machine {
	return &machinev1.Machine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Machine",
			APIVersion: machinev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Status: machinev1.MachineStatus{
			NodeRef: &corev1.ObjectReference{
				Name: nodeName,
			},
		},
	}
}

func createMachine(m *machinev1.Machine) {
	typeMeta := m.TypeMeta
	status := m.Status
	Expect(k8sClient.Create(ctx, m)).To(Succeed())
	m.Status = status
	Expect(k8sClient.Status().Update(ctx, m)).To(Succeed())

	// Fetch object to sync back to latest changes
	key := client.ObjectKey{Namespace: m.Namespace, Name: m.Name}
	Expect(k8sClient.Get(ctx, key, m)).To(Succeed())
	// Restore TypeMeta as not restored by Get
	m.TypeMeta = typeMeta
}
