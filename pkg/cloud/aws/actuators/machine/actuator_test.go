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

package machine_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators/machine/mock_machineiface"
	service "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloudtest"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
)

type ec2svc struct {
	instance *v1alpha1.Instance
}

func (m *ec2svc) InstanceIfExists(instanceID *string) (*v1alpha1.Instance, error) { return nil, nil }
func (m *ec2svc) CreateInstance(machine *clusterv1.Machine, config *v1alpha1.AWSMachineProviderConfig, clusterStatus *v1alpha1.AWSClusterProviderStatus, cluster *clusterv1.Cluster) (*v1alpha1.Instance, error) {
	return nil, nil
}
func (m *ec2svc) TerminateInstance(instanceID string) error { return nil }
func (m *ec2svc) DeleteBastion(instanceID string, status *v1alpha1.AWSClusterProviderStatus) error {
	return nil
}
func (m *ec2svc) CreateOrGetMachine(machine *clusterv1.Machine, status *v1alpha1.AWSMachineProviderStatus, config *v1alpha1.AWSMachineProviderConfig, clusterStatus *v1alpha1.AWSClusterProviderStatus, cluster *clusterv1.Cluster) (*v1alpha1.Instance, error) {
	return m.instance, nil
}
func (m *ec2svc) UpdateInstanceSecurityGroups(instanceID string, securityGroups []string) error {
	return nil
}
func (m *ec2svc) UpdateResourceTags(resourceID *string, create map[string]string, remove map[string]string) error {
	return nil
}
func (m *ec2svc) ReconcileNetwork(clusterName string, network *v1alpha1.Network) error { return nil }
func (m *ec2svc) ReconcileBastion(clusterName, keyName string, status *v1alpha1.AWSClusterProviderStatus) error {
	return nil
}
func (m *ec2svc) DeleteNetwork(clusterName string, network *v1alpha1.Network) error { return nil }

type machinesvc struct {
	mock *mock_machineiface.MockMachineInterface
}

func (m *machinesvc) Machines(namespace string) client.MachineInterface {
	return m.mock
}

type getter struct {
	ec2 service.EC2Interface
}

func (g *getter) Session(*v1alpha1.AWSClusterProviderConfig) *session.Session {
	return nil
}
func (g *getter) EC2(*session.Session) service.EC2Interface {
	return g.ec2
}
func (g *getter) ELB(*session.Session) service.ELBInterface {
	return nil
}

func TestCreate(t *testing.T) {
	testcases := []struct {
		name           string
		machine        *clusterv1.Machine
		service        *ec2svc
		machineExpects func(*mock_machineiface.MockMachineInterface)
	}{
		{
			name: "ensure machine and status are both updated",
			machine: &clusterv1.Machine{
				Spec: clusterv1.MachineSpec{
					ProviderConfig: clusterv1.ProviderConfig{
						Value: cloudtest.RuntimeRawExtension(t, v1alpha1.AWSMachineProviderConfig{}),
					},
				},
				Status: clusterv1.MachineStatus{},
			},
			service: &ec2svc{
				instance: &v1alpha1.Instance{
					ID: "hello world",
				},
			},
			machineExpects: func(mock *mock_machineiface.MockMachineInterface) {
				mock.EXPECT().Update(&clusterv1.Machine{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"cluster-api-provider-aws": "true",
						},
					},
					Spec: clusterv1.MachineSpec{
						ProviderConfig: clusterv1.ProviderConfig{
							Value: cloudtest.RuntimeRawExtension(t, v1alpha1.AWSMachineProviderConfig{}),
						},
					},
				}).Return(nil, nil)

				mock.EXPECT().UpdateStatus(&clusterv1.Machine{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"cluster-api-provider-aws": "true",
						},
					},
					Spec: clusterv1.MachineSpec{
						ProviderConfig: clusterv1.ProviderConfig{
							Value: cloudtest.RuntimeRawExtension(t, v1alpha1.AWSMachineProviderConfig{}),
						},
					},
					Status: clusterv1.MachineStatus{
						ProviderStatus: cloudtest.RuntimeRawExtension(t, v1alpha1.AWSMachineProviderStatus{
							InstanceID:    aws.String("hello world"),
							InstanceState: aws.String(""),
						}),
					},
				}).Return(nil, nil)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockMachine := mock_machineiface.NewMockMachineInterface(mockCtrl)
			tc.machineExpects(mockMachine)

			a := machine.NewActuator(machine.ActuatorParams{
				ServicesGetter: &getter{
					ec2: tc.service,
				},
				MachinesGetter: &machinesvc{
					mock: mockMachine,
				},
			})
			err := a.Create(&clusterv1.Cluster{}, tc.machine)
			if err != nil {
				t.Fatalf("did not expect an error: %v", err)
			}
		})
	}
}
