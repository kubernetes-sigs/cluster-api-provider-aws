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
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/providerconfig/v1alpha1"
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
		check      func(instance *v1alpha1.Instance, err error)
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
			check: func(instance *v1alpha1.Instance, err error) {
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
									&ec2.Instance{
										InstanceId:   aws.String("id-1"),
										InstanceType: aws.String("m5.large"),
										SubnetId:     aws.String("subnet-1"),
										ImageId:      aws.String("ami-1"),
										State: &ec2.InstanceState{
											Code: aws.Int64(16),
											Name: aws.String(ec2.StateAvailable),
										},
									},
								},
							},
						},
					}, nil)
			},
			check: func(instance *v1alpha1.Instance, err error) {
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
			check: func(i *v1alpha1.Instance, err error) {
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
		name    string
		machine clusterv1.Machine
		expect  func(m *mock_ec2iface.MockEC2API)
		check   func(instance *v1alpha1.Instance, err error)
	}{
		{
			name: "simple",
			machine: clusterv1.Machine{
				Spec: clusterv1.MachineSpec{
					ProviderConfig: clusterv1.ProviderConfig{
						Value: &runtime.RawExtension{
							Raw: []byte(`apiVersion: "cluster.k8s.io/v1alpha1"
kind: Machine
metadata:
  generateName: aws-controlplane-
  labels:
    set: controlplane
spec:
  versions:
    kubelet: v1.11.2
    controlPlane: v1.11.2`),
						},
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					RunInstances(&ec2.RunInstancesInput{
						ImageId:      aws.String("abc"),
						InstanceType: aws.String("something"),
						MaxCount:     aws.Int64(1),
						MinCount:     aws.Int64(1),
						SubnetId:     aws.String(""),
					}).
					Return(&ec2.Reservation{
						Instances: []*ec2.Instance{
							&ec2.Instance{
								State: &ec2.InstanceState{
									Name: aws.String(ec2.InstanceStateNamePending),
								},
								InstanceId:   aws.String("two"),
								InstanceType: aws.String("m5.large"),
								SubnetId:     aws.String("subnet-1"),
								ImageId:      aws.String("ami-1"),
							},
						},
					}, nil)
			},
			check: func(instance *v1alpha1.Instance, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)
			s := ec2svc.NewService(ec2Mock)
			instance, err := s.CreateInstance(&tc.machine, &v1alpha1.AWSMachineProviderConfig{
				AMI: v1alpha1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType: "something",
			}, &v1alpha1.AWSClusterProviderStatus{
				Network: v1alpha1.Network{
					Subnets: v1alpha1.Subnets{
						&v1alpha1.Subnet{
							IsPublic: true,
						},
					},
				},
			})
			tc.check(instance, err)
		})
	}
}
