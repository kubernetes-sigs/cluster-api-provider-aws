/*
Copyright 2018 The Kubernetes Authors.

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

package asg

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/golang/mock/gomock"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha4"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/autoscaling/mock_autoscalingiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
)

func TestService_GetASGByName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name            string
		machinePoolName string
		expect          func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
		check           func(autoscalinggroup *expinfrav1.AutoScalingGroup, err error)
	}{
		{
			name:            "ASG is not found",
			machinePoolName: "test-asg-is-not-present",
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DescribeAutoScalingGroups(gomock.Eq(&autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: []*string{
						aws.String("test-asg-is-not-present"),
					},
				})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
			check: func(autoscalinggroup *expinfrav1.AutoScalingGroup, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if autoscalinggroup != nil {
					t.Fatalf("did not expect anything but got something: %+v", autoscalinggroup)
				}
			},
		},
		{
			name:            "ASG should be found",
			machinePoolName: "test-group-is-present",
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DescribeAutoScalingGroups(gomock.Eq(&autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: []*string{
						aws.String("test-group-is-present"),
					},
				})).
					Return(&autoscaling.DescribeAutoScalingGroupsOutput{
						AutoScalingGroups: []*autoscaling.Group{
							{
								AutoScalingGroupName: aws.String("test-group-is-present"),
								MixedInstancesPolicy: &autoscaling.MixedInstancesPolicy{
									InstancesDistribution: &autoscaling.InstancesDistribution{
										OnDemandAllocationStrategy: aws.String("prioritized"),
									},
									LaunchTemplate: &autoscaling.LaunchTemplate{},
								},
							},
						}}, nil)
			},
			check: func(autoscalinggroup *expinfrav1.AutoScalingGroup, err error) {
				if err != nil {
					t.Fatalf("did not expect error: %v", err)
				}
				if autoscalinggroup == nil {
					t.Fatalf("Expected autoscaling group, but didn't get any: %+v", autoscalinggroup)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			cs, err := scope.NewClusterScope(scope.ClusterScopeParams{
				Client:     client,
				Cluster:    &clusterv1.Cluster{},
				AWSCluster: &infrav1.AWSCluster{},
			})
			if err != nil {
				t.Fatalf("Failed to create test context: %v", err)
			}

			tt.expect(asgMock.EXPECT())
			s := NewService(cs)
			s.ASGClient = asgMock

			mps := &scope.MachinePoolScope{
				AWSMachinePool: &expinfrav1.AWSMachinePool{},
			}
			mps.AWSMachinePool.Name = tt.machinePoolName

			asg, err := s.GetASGByName(mps)
			tt.check(asg, err)
		})
	}
}

func TestService_SDKToAutoScalingGroup(t *testing.T) {
	tests := []struct {
		name    string
		input   *autoscaling.Group
		want    *expinfrav1.AutoScalingGroup
		wantErr bool
	}{
		{
			name: "valid input - all required fields filled",
			input: &autoscaling.Group{
				AutoScalingGroupARN:  aws.String("test-id"),
				AutoScalingGroupName: aws.String("test-name"),
				DesiredCapacity:      aws.Int64(1234),
				MaxSize:              aws.Int64(1234),
				MinSize:              aws.Int64(1234),
				CapacityRebalance:    aws.Bool(true),
				AvailabilityZones:    aws.StringSlice([]string{"test-az"}),
				MixedInstancesPolicy: &autoscaling.MixedInstancesPolicy{
					InstancesDistribution: &autoscaling.InstancesDistribution{
						OnDemandAllocationStrategy:          aws.String("prioritized"),
						OnDemandBaseCapacity:                aws.Int64(1234),
						OnDemandPercentageAboveBaseCapacity: aws.Int64(1234),
						SpotAllocationStrategy:              aws.String("lowest-price"),
					},
					LaunchTemplate: &autoscaling.LaunchTemplate{
						Overrides: []*autoscaling.LaunchTemplateOverrides{
							{
								InstanceType:     aws.String("t2.medium"),
								WeightedCapacity: aws.String("test-weighted-cap"),
							},
						},
					},
				},
			},
			want: &expinfrav1.AutoScalingGroup{
				ID:                "test-id",
				Name:              "test-name",
				DesiredCapacity:   aws.Int32(1234),
				MaxSize:           int32(1234),
				MinSize:           int32(1234),
				CapacityRebalance: true,
				AvailabilityZones: []string{"test-az"},
				MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
					InstancesDistribution: &expinfrav1.InstancesDistribution{
						OnDemandAllocationStrategy:          expinfrav1.OnDemandAllocationStrategyPrioritized,
						OnDemandBaseCapacity:                aws.Int64(1234),
						OnDemandPercentageAboveBaseCapacity: aws.Int64(1234),
						SpotAllocationStrategy:              expinfrav1.SpotAllocationStrategyLowestPrice,
					},
					Overrides: []expinfrav1.Overrides{
						{
							InstanceType: "t2.medium",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid input - without mixedInstancesPolicy",
			input: &autoscaling.Group{
				AutoScalingGroupARN:  aws.String("test-id"),
				AutoScalingGroupName: aws.String("test-name"),
				DesiredCapacity:      aws.Int64(1234),
				MaxSize:              aws.Int64(1234),
				MinSize:              aws.Int64(1234),
				CapacityRebalance:    aws.Bool(true),
				AvailabilityZones:    aws.StringSlice([]string{"test-az"}),
				MixedInstancesPolicy: nil,
			},
			want: &expinfrav1.AutoScalingGroup{
				ID:                   "test-id",
				Name:                 "test-name",
				DesiredCapacity:      aws.Int32(1234),
				MaxSize:              int32(1234),
				MinSize:              int32(1234),
				CapacityRebalance:    true,
				AvailabilityZones:    []string{"test-az"},
				MixedInstancesPolicy: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			got, err := s.SDKToAutoScalingGroup(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SDKToAutoScalingGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.SDKToAutoScalingGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
