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
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	clientv1 "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/machine"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/actuators/machine/mock_machineiface"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2/mock_ec2iface"
)

type machinesGetter struct {
	mi *mock_machineiface.MockMachineInterface
}

func (m *machinesGetter) Machines(ns string) clientv1.MachineInterface {
	return m.mi
}

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mg := &machinesGetter{
		mi: mock_machineiface.NewMockMachineInterface(mockCtrl),
	}
	me := mock_ec2iface.NewMockEC2API(mockCtrl)
	defer mockCtrl.Finish()

	// clusterapi calls
	mg.mi.EXPECT().
		UpdateStatus(&clusterv1.Machine{
			ObjectMeta: v1.ObjectMeta{
				Labels:      map[string]string{"set": "node"},
				Annotations: map[string]string{"cluster-api-provider-aws": "true"},
			},
			Status: clusterv1.MachineStatus{
				ProviderStatus: &runtime.RawExtension{
					Raw: []byte(`{"kind":"AWSMachineProviderStatus","apiVersion":"awsproviderconfig/v1alpha1","instanceID":"1234","instanceState":"running"}
`),
				},
			},
			Spec: clusterv1.MachineSpec{
				ProviderConfig: clusterv1.ProviderConfig{
					Value: &runtime.RawExtension{
						Raw: []byte(`{"kind":"AWSMachineProviderConfig","apiVersion":"awsproviderconfig/v1alpha1","ami":{"id":"aws-instance-id"}}
`),
					},
				},
			},
		}).
		Return(&clusterv1.Machine{}, nil)

	// ec2 calls
	me.EXPECT().
		RunInstances(&ec2.RunInstancesInput{
			ImageId:          aws.String("aws-instance-id"),
			InstanceType:     aws.String(""),
			MaxCount:         aws.Int64(1),
			MinCount:         aws.Int64(1),
			SubnetId:         aws.String("subnet-1"),
			SecurityGroupIds: aws.StringSlice([]string{"2"}),
		}).
		Return(&ec2.Reservation{
			Instances: []*ec2.Instance{
				&ec2.Instance{
					State: &ec2.InstanceState{
						Name: aws.String(ec2.InstanceStateNameRunning),
					},
					InstanceId:   aws.String("1234"),
					InstanceType: aws.String("m5.large"),
					SubnetId:     aws.String("subnet-1"),
					ImageId:      aws.String("ami-1"),
				},
			},
		}, nil)

	codec, err := v1alpha1.NewCodec()
	if err != nil {
		t.Fatalf("failed to create a codec: %v", err)
	}
	ap := machine.ActuatorParams{
		Codec:          codec,
		MachinesGetter: mg,
		EC2Service: &ec2svc.Service{
			EC2: me,
		},
	}
	actuator, err := machine.NewActuator(ap)
	if err != nil {
		t.Fatalf("failed to create an actuator: %v", err)
	}

	if err := actuator.Create(&clusterv1.Cluster{
		Status: clusterv1.ClusterStatus{
			ProviderStatus: &runtime.RawExtension{
				Raw: []byte(`{"kind":"AWSClusterProviderStatus","apiVersion":"awsproviderconfig/v1alpha1","network":{"subnets":[{"id": "subnet-1", "public": false}],"securityGroups":{"node":{"id":"2"}}}}
`),
			},
		},
	}, &clusterv1.Machine{
		ObjectMeta: v1.ObjectMeta{
			Labels: map[string]string{
				"set": "node",
			},
		},
		Spec: clusterv1.MachineSpec{
			ProviderConfig: clusterv1.ProviderConfig{
				Value: &runtime.RawExtension{
					Raw: []byte(`{"kind":"AWSMachineProviderConfig","apiVersion":"awsproviderconfig/v1alpha1","ami":{"id":"aws-instance-id"}}
`),
				},
			},
		},
	}); err != nil {
		t.Fatalf("failed to create machine: %v", err)
	}
}

func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mg := &machinesGetter{
		mi: mock_machineiface.NewMockMachineInterface(mockCtrl),
	}
	me := mock_ec2iface.NewMockEC2API(mockCtrl)
	defer mockCtrl.Finish()

	gomock.InOrder(
		// ec2 calls
		me.EXPECT().
			DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String("2345")},
			}).
			Return(&ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					&ec2.Reservation{
						Instances: []*ec2.Instance{
							&ec2.Instance{
								State: &ec2.InstanceState{
									Name: aws.String(ec2.InstanceStateNameRunning),
								},
								InstanceId:   aws.String("2345"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-1"),
								ImageId:      aws.String("ami-1"),
							},
						},
					},
				},
			}, nil),
	)
	me.EXPECT().
		TerminateInstances(&ec2.TerminateInstancesInput{
			InstanceIds: []*string{aws.String("2345")},
		}).
		Return(nil, nil)

	codec, err := v1alpha1.NewCodec()
	if err != nil {
		t.Fatalf("failed to create a codec: %v", err)
	}

	ap := machine.ActuatorParams{
		Codec:          codec,
		MachinesGetter: mg,
		EC2Service: &ec2svc.Service{
			EC2: me,
		},
	}

	actuator, err := machine.NewActuator(ap)
	if err != nil {
		t.Fatalf("failed to create an actuator: %v", err)
	}

	testMachine := &clusterv1.Machine{
		Status: clusterv1.MachineStatus{
			ProviderStatus: &runtime.RawExtension{
				Raw: []byte(`{"kind":"AWSMachineProviderStatus","apiVersion":"awsproviderconfig/v1alpha1","instanceID":"2345","instanceState":"running"}
`),
			},
		},
	}

	// Delete the machine.
	if err := actuator.Delete(&clusterv1.Cluster{}, testMachine); err != nil {
		t.Fatalf("failed to delete machine: %v", err)
	}
}

func TestDeleteNotExisting(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mg := &machinesGetter{
		mi: mock_machineiface.NewMockMachineInterface(mockCtrl),
	}
	me := mock_ec2iface.NewMockEC2API(mockCtrl)
	defer mockCtrl.Finish()

	// ec2 calls
	me.EXPECT().
		DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{nil},
		}).
		Return(nil, ec2svc.NewNotFound(errors.New("")))

	codec, err := v1alpha1.NewCodec()
	if err != nil {
		t.Fatalf("failed to create a codec: %v", err)
	}

	ap := machine.ActuatorParams{
		Codec:          codec,
		MachinesGetter: mg,
		EC2Service: &ec2svc.Service{
			EC2: me,
		},
	}

	actuator, err := machine.NewActuator(ap)
	if err != nil {
		t.Fatalf("failed to create an actuator: %v", err)
	}

	// Get some empty cluster and machine structs.
	testCluster := &clusterv1.Cluster{}
	testMachine := &clusterv1.Machine{}

	// Delete the machine.
	if err := actuator.Delete(testCluster, testMachine); err != nil {
		t.Fatalf("failed to delete machine: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mg := &machinesGetter{
		mi: mock_machineiface.NewMockMachineInterface(mockCtrl),
	}
	me := mock_ec2iface.NewMockEC2API(mockCtrl)
	defer mockCtrl.Finish()

	mg.mi.EXPECT().
		UpdateStatus(&clusterv1.Machine{
			Status: clusterv1.MachineStatus{
				ProviderStatus: &runtime.RawExtension{
					Raw: []byte(`{"kind":"AWSMachineProviderStatus","apiVersion":"awsproviderconfig/v1alpha1"}
`),
				},
			},
		}).
		Return(&clusterv1.Machine{}, nil)

	codec, err := v1alpha1.NewCodec()
	if err != nil {
		t.Fatalf("failed to create a codec: %v", err)
	}

	ap := machine.ActuatorParams{
		Codec:          codec,
		MachinesGetter: mg,
		EC2Service: &ec2svc.Service{
			EC2: me,
		},
	}

	actuator, err := machine.NewActuator(ap)
	if err != nil {
		t.Fatalf("failed to create an actuator: %v", err)
	}

	testCluster := &clusterv1.Cluster{}
	testMachine := &clusterv1.Machine{}

	// Update the machine.
	if err := actuator.Update(testCluster, testMachine); err != nil {
		t.Fatalf("failed to delete machine: %v", err)
	}
}
