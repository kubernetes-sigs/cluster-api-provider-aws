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

package ec2_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	providerconfigv1 "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2/mock_ec2iface"
)

func TestInstanceIfExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name       string
		instanceID string
		expect     func(m *mock_ec2iface.MockEC2API)
		check      func(instance *ec2svc.Instance, err error)
	}{
		{
			name:       "does not exist",
			instanceID: "hello",
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
						InstanceIds: []*string{aws.String("hello")},
					})).
					Return(nil, ec2svc.NewNotFound(errors.New("not found")))
			},
			check: func(instance *ec2svc.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if instance != nil {
					t.Fatalf("Did not expect anything but got something: %+v", instance)
				}
			},
		},
		{
			name:       "instance exists",
			instanceID: "id-1",
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeInstances(gomock.Eq(&ec2.DescribeInstancesInput{
						InstanceIds: []*string{aws.String("id-1")},
					})).
					Return(&ec2.DescribeInstancesOutput{
						Reservations: []*ec2.Reservation{
							{
								Instances: []*ec2.Instance{
									{
										State: &ec2.InstanceState{
											Code: aws.Int64(16),
											Name: aws.String(ec2.StateAvailable),
										},
										InstanceId: aws.String("id-1"),
									},
								},
							},
						},
					}, nil)
			},
			check: func(instance *ec2svc.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}

				if instance == nil {
					t.Fatalf("expected instance but got nothing")
				}

				if instance.ID != "id-1" {
					t.Fatalf("expected id-1 but got: %v", instance.ID)
				}
			},
		},
		{
			name:       "error describing instances",
			instanceID: "one",
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: []*string{aws.String("one")}}).
					Return(nil, errors.New("some unknown error"))
			},
			check: func(i *ec2svc.Instance, err error) {
				if err == nil {
					t.Fatalf("expected an error but got none.")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)
			s := ec2svc.NewService(ec2Mock)
			instance, err := s.InstanceIfExists(&tc.instanceID)
			tc.check(instance, err)
		})
	}
}

func TestTerminateInstance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	instanceNotFoundError := errors.New("instance not found")

	testCases := []struct {
		name       string
		instanceID string
		expect     func(m *mock_ec2iface.MockEC2API)
		check      func(err error)
	}{
		{
			name:       "instance exists",
			instanceID: "i-exist",
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					TerminateInstances(gomock.Eq(&ec2.TerminateInstancesInput{
						InstanceIds: []*string{aws.String("i-exist")},
					})).
					Return(&ec2.TerminateInstancesOutput{}, nil)
			},
			check: func(err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
		{
			name:       "instance does not exist",
			instanceID: "i-donotexist",
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					TerminateInstances(gomock.Eq(&ec2.TerminateInstancesInput{
						InstanceIds: []*string{aws.String("i-donotexist")},
					})).
					Return(&ec2.TerminateInstancesOutput{}, instanceNotFoundError)
			},
			check: func(err error) {
				if err == nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)
			s := ec2svc.NewService(ec2Mock)
			err := s.TerminateInstance(&tc.instanceID)
			tc.check(err)
		})
	}
}

func TestCreateInstance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testcases := []struct {
		name                  string
		machine               clusterv1.Machine
		providerConfig        providerconfigv1.AWSMachineProviderConfig
		cluster               clusterv1.Cluster
		clusterProviderconfig providerconfigv1.AWSClusterProviderConfig
		clusterStatus         providerconfigv1.AWSClusterProviderStatus
		expect                func(m *mock_ec2iface.MockEC2API)
		check                 func(instance *ec2svc.Instance, err error)
	}{
		{
			name: "simple",
			machine: clusterv1.Machine{
				Spec: clusterv1.MachineSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-1",
					},
					Versions: clusterv1.MachineVersionInfo{
						Kubelet:      "v1.11.3",
						ControlPlane: "v1.11.3",
					},
				},
			},
			providerConfig: providerconfigv1.AWSMachineProviderConfig{},
			cluster: clusterv1.Cluster{
				Spec: clusterv1.ClusterSpec{
					ClusterNetwork: clusterv1.ClusterNetworkingConfig{
						ServiceDomain: "cluster.local",
						Pods: clusterv1.NetworkRanges{
							CIDRBlocks: []string{"10.96.0.0/12"},
						},
						Services: clusterv1.NetworkRanges{
							CIDRBlocks: []string{"192.168.0.0/16"},
						},
					},
				},
			},
			clusterStatus: providerconfigv1.AWSClusterProviderStatus{
				Network: providerconfigv1.Network{
					Subnets: providerconfigv1.Subnets{
						&providerconfigv1.Subnet{
							ID:       "subnet-id",
							IsPublic: false,
						},
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					RunInstances(gomock.Any()).
					Return(&ec2.Reservation{
						Instances: []*ec2.Instance{
							&ec2.Instance{
								State: &ec2.InstanceState{
									Name: aws.String(ec2.InstanceStateNamePending),
								},
								InstanceId: aws.String("two"),
							},
						},
					}, nil)
			},
			check: func(instance *ec2svc.Instance, err error) {
				if err != nil {
					t.Errorf("did not expect error: %v", err)
				}
			},
		},
		{
			name: "simple controlplane instance",
			machine: clusterv1.Machine{
				Spec: clusterv1.MachineSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-cp-1",
					},
					Versions: clusterv1.MachineVersionInfo{
						Kubelet:      "v1.11.3",
						ControlPlane: "v1.11.3",
					},
				},
			},
			providerConfig: providerconfigv1.AWSMachineProviderConfig{
				Subnet: &providerconfigv1.AWSResourceReference{
					ID: aws.String("subnet-id"),
				},
				AMI: providerconfigv1.AWSResourceReference{
					ID: aws.String("ami-id"),
				},
				InstanceType: "t3.medium",
				KeyName:      "default",
				NodeRole:     "controlplane",
			},
			cluster: clusterv1.Cluster{
				Spec: clusterv1.ClusterSpec{
					ClusterNetwork: clusterv1.ClusterNetworkingConfig{
						ServiceDomain: "cluster.local",
						Pods: clusterv1.NetworkRanges{
							CIDRBlocks: []string{"10.96.0.0/12"},
						},
						Services: clusterv1.NetworkRanges{
							CIDRBlocks: []string{"192.168.0.0/16"},
						},
					},
				},
			},
			clusterStatus: providerconfigv1.AWSClusterProviderStatus{
				Network: providerconfigv1.Network{
					Subnets: providerconfigv1.Subnets{
						&providerconfigv1.Subnet{
							ID:       "subnet-id",
							IsPublic: false,
						},
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					RunInstances(gomock.Any()).
					Return(&ec2.Reservation{
						Instances: []*ec2.Instance{
							&ec2.Instance{
								State: &ec2.InstanceState{
									Name: aws.String(ec2.InstanceStateNamePending),
								},
								InstanceId: aws.String("two"),
							},
						},
					}, nil)
			},
			check: func(instance *ec2svc.Instance, err error) {
				if err != nil {
					t.Errorf("did not expect error: %v", err)
				}
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)
			s := ec2svc.NewService(ec2Mock)
			instance, err := s.CreateInstance(&tc.machine, &tc.providerConfig, &tc.cluster, &tc.clusterProviderconfig, &tc.clusterStatus)
			tc.check(instance, err)
		})
	}
}
