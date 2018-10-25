// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	providerv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/cluster"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/cluster/mock_clusteriface" // nolint
	service "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clientv1 "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
)

type clusterGetter struct {
	ci *mock_clusteriface.MockClusterInterface
}

func (c *clusterGetter) Clusters(ns string) clientv1.ClusterInterface {
	return c.ci
}

type testServicesGetter struct {
	ec2 *mocks.MockEC2Interface
	elb *mocks.MockELBInterface
}

func (d *testServicesGetter) Session(clusterConfig *providerv1.AWSClusterProviderConfig) *session.Session {
	return nil
}

func (d *testServicesGetter) EC2(session *session.Session) service.EC2Interface {
	return d.ec2
}

func (d *testServicesGetter) ELB(session *session.Session) service.ELBInterface {
	return d.elb
}

func TestGetIP(t *testing.T) {
	testcases := []struct {
		name       string
		cluster    *clusterv1.Cluster
		expectedIP string
		elbExpects func(*mocks.MockELBInterface)
	}{
		{
			name: "sunny day test",
			cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test"},
				Spec: clusterv1.ClusterSpec{
					ProviderConfig: clusterv1.ProviderConfig{
						Value: RuntimeRawExtension(&providerv1.AWSClusterProviderConfig{}, t),
					},
				},
				Status: clusterv1.ClusterStatus{
					ProviderStatus: RuntimeRawExtension(&providerv1.AWSClusterProviderStatus{}, t),
				},
			},
			expectedIP: "something",
			elbExpects: func(m *mocks.MockELBInterface) {
				m.EXPECT().
					GetAPIServerDNSName("test").
					Return("something", nil)
			},
		},
		{
			name: "lookup IP if the status is empty",
			cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test"},
				Spec: clusterv1.ClusterSpec{
					ProviderConfig: clusterv1.ProviderConfig{
						Value: RuntimeRawExtension(&providerv1.AWSClusterProviderConfig{}, t),
					},
				},
			},
			expectedIP: "dunno",
			elbExpects: func(m *mocks.MockELBInterface) {
				m.EXPECT().
					GetAPIServerDNSName("test").
					Return("dunno", nil)
			},
		},
		{
			name: "return the IP if it is stored in the status",
			cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test"},
				Spec: clusterv1.ClusterSpec{
					ProviderConfig: clusterv1.ProviderConfig{
						Value: RuntimeRawExtension(&providerv1.AWSClusterProviderConfig{}, t),
					},
				},
				Status: clusterv1.ClusterStatus{
					ProviderStatus: RuntimeRawExtension(&providerv1.AWSClusterProviderStatus{
						Network: providerv1.Network{
							APIServerELB: providerv1.ClassicELB{
								DNSName: "banana",
							},
						},
					}, t),
				},
			},
			expectedIP: "banana",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			clusters := &clusterGetter{
				ci: mock_clusteriface.NewMockClusterInterface(mockCtrl),
			}

			services := &testServicesGetter{
				ec2: mocks.NewMockEC2Interface(mockCtrl),
				elb: mocks.NewMockELBInterface(mockCtrl),
			}

			ap := cluster.ActuatorParams{
				ClustersGetter: clusters,
				ServicesGetter: services,
			}

			actuator, err := cluster.NewActuator(ap)
			if err != nil {
				t.Fatalf("could not create an actuator: %v", err)
			}

			if tc.elbExpects != nil {
				tc.elbExpects(services.elb)
			}

			ip, err := actuator.GetIP(tc.cluster, nil)
			if err != nil {
				t.Fatalf("failed to get API server address: %v", err)
			}

			if ip != tc.expectedIP {
				t.Fatalf("got the wrong IP. Found %v, wanted %v", ip, tc.expectedIP)
			}
		})
	}
}

func TestReconcile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	clusters := &clusterGetter{
		ci: mock_clusteriface.NewMockClusterInterface(mockCtrl),
	}

	services := &testServicesGetter{
		ec2: mocks.NewMockEC2Interface(mockCtrl),
		elb: mocks.NewMockELBInterface(mockCtrl),
	}

	ap := cluster.ActuatorParams{
		ClustersGetter: clusters,
		ServicesGetter: services,
	}

	actuator, err := cluster.NewActuator(ap)
	if err != nil {
		t.Fatalf("could not create an actuator: %v", err)
	}

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test", ClusterName: "test"},
		Spec: clusterv1.ClusterSpec{
			ProviderConfig: clusterv1.ProviderConfig{
				Value: RuntimeRawExtension(&providerv1.AWSClusterProviderConfig{
					Region: "us-east-1",
				}, t),
			},
		},
	}

	clusters.ci.EXPECT().
		UpdateStatus(gomock.AssignableToTypeOf(cluster)).
		Return(cluster, nil)

	services.ec2.EXPECT().
		ReconcileNetwork("test", gomock.AssignableToTypeOf(&providerv1.Network{})).
		Return(nil)

	services.ec2.EXPECT().
		ReconcileBastion("test", "", gomock.AssignableToTypeOf(&providerv1.AWSClusterProviderStatus{})).
		Return(nil)

	services.elb.EXPECT().
		ReconcileLoadbalancers("test", gomock.AssignableToTypeOf(&providerv1.Network{})).
		Return(nil)

	if err := actuator.Reconcile(cluster); err != nil {
		t.Fatalf("failed to reconcile cluster: %v", err)
	}
}

func RuntimeRawExtension(p interface{}, t *testing.T) *runtime.RawExtension {
	t.Helper()
	out, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	return &runtime.RawExtension{
		Raw: out,
	}
}
