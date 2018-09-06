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

	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2/mock_ec2iface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
)

type instances struct{}

func (i *instances) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{}, nil
}
func (i *instances) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	return &ec2.Reservation{}, nil
}

func TestInstanceIfExists_DoesNotExist(t *testing.T) {
	s := ec2svc.Service{
		Instances: &instances{},
	}
	id := "hello"
	out, err := s.InstanceIfExists(&id)
	if err != nil {
		t.Fatalf("did not expect error: %v", err)
	}
	if out != nil {
		t.Fatal("Did not expect anything but got something")
	}
}

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
					Return(&ec2.DescribeInstancesOutput{}, nil)
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)
			tc.expect(ec2Mock)
			s := ec2svc.Service{Instances: ec2Mock}
			instance, err := s.InstanceIfExists(&tc.instanceID)
			tc.check(instance, err)
		})
	}
}
