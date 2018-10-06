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
									{
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
			err := s.TerminateInstance(tc.instanceID)
			tc.check(err)
		})
	}
}

func TestCreateInstance(t *testing.T) {
	testcases := []struct {
		name          string
		machine       clusterv1.Machine
		machineConfig *v1alpha1.AWSMachineProviderConfig
		clusterStatus *v1alpha1.AWSClusterProviderStatus
		expect        func(m *mock_ec2iface.MockEC2API)
		check         func(instance *v1alpha1.Instance, err error)
	}{
		{
			name: "simple",
			machine: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"set": "node"},
				},
			},
			machineConfig: &v1alpha1.AWSMachineProviderConfig{
				AMI: v1alpha1.AWSResourceReference{
					ID: aws.String("abc"),
				},
				InstanceType: "m5.large",
			},
			clusterStatus: &v1alpha1.AWSClusterProviderStatus{
				Network: v1alpha1.Network{
					Subnets: v1alpha1.Subnets{
						&v1alpha1.Subnet{
							ID:       "subnet-1",
							IsPublic: false,
						},
						&v1alpha1.Subnet{
							IsPublic: false,
						},
					},
					SecurityGroups: map[v1alpha1.SecurityGroupRole]*v1alpha1.SecurityGroup{
						v1alpha1.SecurityGroupControlPlane: {
							ID: "1",
						},
						v1alpha1.SecurityGroupNode: {
							ID: "2",
						},
					},
				},
			},
			expect: func(m *mock_ec2iface.MockEC2API) {
				m.EXPECT().
					RunInstances(&ec2.RunInstancesInput{
						ImageId:          aws.String("abc"),
						InstanceType:     aws.String("m5.large"),
						KeyName:          aws.String("default"),
						MaxCount:         aws.Int64(1),
						MinCount:         aws.Int64(1),
						SubnetId:         aws.String("subnet-1"),
						SecurityGroupIds: aws.StringSlice([]string{"2"}),
					}).
					Return(&ec2.Reservation{
						Instances: []*ec2.Instance{
							{
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
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)
			s := ec2svc.NewService(ec2Mock)
			instance, err := s.CreateInstance(&tc.machine, tc.machineConfig, tc.clusterStatus)
			tc.check(instance, err)
		})
	}
}
