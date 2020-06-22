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
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var notFoundFunc = func(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(404)
}

var _ = Describe("Handler Suite", func() {
	var terminationServer *httptest.Server
	var httpHandler http.Handler
	var nodeName string
	var testNode *corev1.Node
	var stop chan struct{}
	var errs chan error
	var h *handler

	nodeMarkedForDeletion := func(nodeName string) func() (bool, error) {
		key := client.ObjectKey{Name: nodeName}
		return func() (bool, error) {
			n := &corev1.Node{}
			err := k8sClient.Get(ctx, key, n)
			if err != nil {
				return false, err
			}
			for _, condition := range n.Status.Conditions {
				if condition.Type == terminatingConditionType {
					if condition.Status == corev1.ConditionTrue {
						return true, nil
					}
					// Found the condition with the right type, but wrong status
					return false, nil
				}
			}
			return false, nil
		}
	}

	BeforeEach(func() {
		// Reset test vars
		terminationServer = nil
		httpHandler = nil
		nodeName = "test-node"
		httpHandler = newMockHTTPHandler(notFoundFunc)
		stop = nil
		errs = nil

		testNode = newTestNode(nodeName)
		createNode(testNode)

		// use NewHandler() instead of manual construction in order to test NewHandler() logic
		// like checking that machine api is added to scheme
		handlerInterface, err := NewHandler(klogr.New(), cfg, 100*time.Millisecond, "", nodeName)
		Expect(err).ToNot(HaveOccurred())

		h = handlerInterface.(*handler)

		// set pollURL so we can override initial value later
		h.pollURL = nil
	})

	JustBeforeEach(func() {
		Expect(httpHandler).ToNot(BeNil())
		terminationServer = httptest.NewServer(httpHandler)

		if h.pollURL == nil {
			pollURL, err := url.Parse(terminationServer.URL)
			Expect(err).ToNot(HaveOccurred())
			h.pollURL = pollURL
		}
	})

	AfterEach(func() {
		if stop != nil && !isClosed(stop) {
			close(stop)
		}
		terminationServer.Close()

		Expect(deleteAllNodes(k8sClient)).To(Succeed())
	})

	Context("when running the handler", func() {
		JustBeforeEach(func() {
			stop, errs = StartTestHandler(h)
		})

		Context("when the handler is stopped", func() {
			JustBeforeEach(func() {
				close(stop)
			})

			It("should not return an error", func() {
				Eventually(errs).Should(Receive(BeNil()))
			})
		})

		Context("when polling the termination endpoint", func() {
			var counter int32

			BeforeEach(func() {
				// Ensure the polling logic is excercised in tests
				httpHandler = newMockHTTPHandler(func(rw http.ResponseWriter, req *http.Request) {
					if atomic.LoadInt32(&counter) == 4 {
						rw.WriteHeader(200)
					} else {
						atomic.AddInt32(&counter, 1)
						rw.WriteHeader(404)
					}
				})
			})

			JustBeforeEach(func() {
				// Ensure the polling logic is excercised in tests
				for atomic.LoadInt32(&counter) < 4 {
					continue
				}
			})

			Context("and the handler is stopped", func() {
				JustBeforeEach(func() {
					close(stop)
				})

				It("should not return an error", func() {
					Eventually(errs).Should(Receive(BeNil()))
				})

				It("should not mark the node for deletion", func() {
					Consistently(nodeMarkedForDeletion(testNode.Name)).Should(BeFalse())
				})
			})

			Context("and the instance termination notice is fulfilled", func() {
				It("should mark the node for deletion", func() {
					Eventually(nodeMarkedForDeletion(testNode.Name)).Should(BeTrue())
				})
			})

			Context("and the instance termination notice is not fulfilled", func() {
				BeforeEach(func() {
					httpHandler = newMockHTTPHandler(notFoundFunc)
				})

				It("should not mark the node for deletion", func() {
					Consistently(nodeMarkedForDeletion(testNode.Name)).Should(BeFalse())
				})
			})

			Context("and the instance termination endpoint returns an unknown status", func() {
				BeforeEach(func() {
					httpHandler = newMockHTTPHandler(func(rw http.ResponseWriter, req *http.Request) {
						if atomic.LoadInt32(&counter) == 4 {
							rw.WriteHeader(500)
						} else {
							atomic.AddInt32(&counter, 1)
							rw.WriteHeader(404)
						}
					})
				})

				It("should return an error", func() {
					Eventually(errs).Should(Receive(MatchError("error polling termination endpoint: unexpected status: 500")))
				})

				It("should not delete the machine", func() {
					Consistently(nodeMarkedForDeletion(testNode.Name)).Should(BeFalse())
				})
			})

			Context("and the poll URL cannot be reached", func() {
				BeforeEach(func() {
					h.pollURL = &url.URL{Opaque: "abc#1://localhost"}
				})

				It("should return an error", func() {
					Eventually(errs).Should(Receive(MatchError(ContainSubstring("error polling termination endpoint: could not get URL \"abc#1://localhost\":"))))
				})

				It("should not delete the machine", func() {
					Consistently(nodeMarkedForDeletion(testNode.Name)).Should(BeFalse())
				})
			})
		})

		Context("addNodeTerminationCondition", func() {
			JustBeforeEach(func() {
				addNodeTerminationCondition(testNode)
			})

			Context("with no existing conditions", func() {
				BeforeEach(func() {
					Expect(testNode.Status.Conditions).To(HaveLen(0))
				})

				It("should add the condition to the node", func() {
					Expect(testNode.Status.Conditions).To(HaveLen(1))
					condition := testNode.Status.Conditions[0]
					Expect(condition.Type).To(Equal(terminatingConditionType))
					Expect(condition.Status).To(Equal(corev1.ConditionTrue))
					Expect(condition.Reason).To(Equal(terminationRequestedReason))
				})
			})

			Context("with the terminating condition with the correct status", func() {
				var updated *metav1.Time

				BeforeEach(func() {
					now := metav1.Now()
					updated = &now
					testNode.Status.Conditions = []corev1.NodeCondition{
						{
							Type:               terminatingConditionType,
							Status:             corev1.ConditionTrue,
							Reason:             terminationRequestedReason,
							LastTransitionTime: now,
							LastHeartbeatTime:  now,
						},
					}
				})

				It("should not update the condition on the node", func() {
					Expect(testNode.Status.Conditions).To(HaveLen(1))
					condition := testNode.Status.Conditions[0]
					Expect(condition.Type).To(Equal(terminatingConditionType))
					Expect(condition.Status).To(Equal(corev1.ConditionTrue))
					Expect(condition.Reason).To(Equal(terminationRequestedReason))
					Expect(condition.LastTransitionTime).To(Equal(*updated))
					Expect(condition.LastHeartbeatTime).To(Equal(*updated))
				})
			})

			Context("with the terminating condition with the incorrect status", func() {
				var updated *metav1.Time

				BeforeEach(func() {
					now := metav1.Now()
					updated = &now
					testNode.Status.Conditions = []corev1.NodeCondition{
						{
							Type:               terminatingConditionType,
							Status:             corev1.ConditionFalse,
							Reason:             terminationRequestedReason,
							LastTransitionTime: now,
							LastHeartbeatTime:  now,
						},
					}
				})

				It("should update the condition on the node", func() {
					Expect(testNode.Status.Conditions).To(HaveLen(1))
					condition := testNode.Status.Conditions[0]
					Expect(condition.Type).To(Equal(terminatingConditionType))
					Expect(condition.Status).To(Equal(corev1.ConditionTrue))
					Expect(condition.Reason).To(Equal(terminationRequestedReason))
					Expect(condition.LastTransitionTime).ToNot(Equal(*updated))
					Expect(condition.LastHeartbeatTime).ToNot(Equal(*updated))
				})
			})

			Context("with existing conditions", func() {
				var existingCondition *corev1.NodeCondition

				BeforeEach(func() {
					now := metav1.Now()
					existingCondition = &corev1.NodeCondition{
						Type:               corev1.NodeReady,
						Status:             corev1.ConditionFalse,
						Reason:             "Unknown",
						LastTransitionTime: now,
						LastHeartbeatTime:  now,
					}

					testNode.Status.Conditions = []corev1.NodeCondition{*existingCondition}
				})

				It("should not modify the existing conditions", func() {
					Expect(testNode.Status.Conditions).To(HaveLen(2))
					Expect(testNode.Status.Conditions).To(ContainElement(Equal(*existingCondition)))
				})
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

func deleteAllNodes(c client.Client) error {
	nodeList := &corev1.NodeList{}
	err := c.List(ctx, nodeList)
	if err != nil {
		return fmt.Errorf("error listing machines: %v", err)
	}

	// Delete all nodes found
	for _, node := range nodeList.Items {
		n := node
		err := c.Delete(ctx, &n)
		if err != nil {
			return err
		}
	}
	return nil
}

func newTestNode(name string) *corev1.Node {
	return &corev1.Node{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Node",
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func createNode(n *corev1.Node) {
	typeMeta := n.TypeMeta
	status := n.Status
	Expect(k8sClient.Create(ctx, n)).To(Succeed())
	n.Status = status
	Expect(k8sClient.Status().Update(ctx, n)).To(Succeed())

	// Fetch object to sync back to latest changes
	key := client.ObjectKey{Namespace: n.Namespace, Name: n.Name}
	Expect(k8sClient.Get(ctx, key, n)).To(Succeed())
	// Restore TypeMeta as not restored by Get
	n.TypeMeta = typeMeta
}
