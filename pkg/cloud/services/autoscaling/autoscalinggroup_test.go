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
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/autoscaling/mock_autoscalingiface"
	"sigs.k8s.io/cluster-api-provider-aws/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

func TestService_GetASGByName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name            string
		machinePoolName string
		wantErr         bool
		wantASG         bool
		expect          func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
	}{
		{
			name:            "should return nil if ASG is not found",
			machinePoolName: "test-asg-is-not-present",
			wantErr:         false,
			wantASG:         false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DescribeAutoScalingGroups(gomock.Eq(&autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: []*string{
						aws.String("test-asg-is-not-present"),
					},
				})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
		},
		{
			name:            "should return error if describe asg failed",
			machinePoolName: "dependency-failure-occurred",
			wantErr:         true,
			wantASG:         false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DescribeAutoScalingGroups(gomock.Eq(&autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: []*string{
						aws.String("dependency-failure-occurred"),
					},
				})).
					Return(nil, awserrors.NewFailedDependency("unknown error occurred"))
			},
		},
		{
			name:            "should return ASG, if found",
			machinePoolName: "test-group-is-present",
			wantErr:         false,
			wantASG:         true,
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fakeClient := getFakeClient()

			clusterScope, err := getClusterScope(fakeClient)
			g.Expect(err).ToNot(HaveOccurred())
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)
			tt.expect(asgMock.EXPECT())
			s := NewService(clusterScope)
			s.ASGClient = asgMock

			mps, err := getMachinePoolScope(fakeClient, clusterScope)
			g.Expect(err).ToNot(HaveOccurred())
			mps.AWSMachinePool.Name = tt.machinePoolName

			asg, err := s.GetASGByName(mps)
			checkErr(tt.wantErr, err, g)
			checkASG(tt.wantASG, asg, g)
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
			name: "valid input - all fields filled",
			input: &autoscaling.Group{
				AutoScalingGroupARN:  aws.String("test-id"),
				AutoScalingGroupName: aws.String("test-name"),
				DesiredCapacity:      aws.Int64(1234),
				MaxSize:              aws.Int64(1234),
				MinSize:              aws.Int64(1234),
				CapacityRebalance:    aws.Bool(true),
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
				Status: aws.String("status"),
				Tags: []*autoscaling.TagDescription{
					{
						Key:   aws.String("key"),
						Value: aws.String("value"),
					},
				},
				Instances: []*autoscaling.Instance{
					{
						InstanceId:     aws.String("instanceId"),
						LifecycleState: aws.String("lifecycleState"),
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
				Status: "status",
				Tags: map[string]string{
					"key": "value",
				},
				Instances: []infrav1.Instance{
					{
						ID:    "instanceId",
						State: "lifecycleState",
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
				MixedInstancesPolicy: nil,
			},
			want: &expinfrav1.AutoScalingGroup{
				ID:                   "test-id",
				Name:                 "test-name",
				DesiredCapacity:      aws.Int32(1234),
				MaxSize:              int32(1234),
				MinSize:              int32(1234),
				CapacityRebalance:    true,
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
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Service.SDKToAutoScalingGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ASGIfExists(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name    string
		asgName *string
		wantErr bool
		wantASG bool
		expect  func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
	}{
		{
			name:    "should return nil if ASG name is not given",
			asgName: nil,
			wantErr: false,
			wantASG: false,
			expect:  func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {},
		},
		{
			name:    "should return without error if ASG is not found",
			asgName: aws.String("asgName"),
			wantErr: false,
			wantASG: false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DescribeAutoScalingGroups(gomock.Eq(&autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: []*string{
						aws.String("asgName"),
					},
				})).
					Return(nil, awserrors.NewNotFound("resource not found"))
			},
		},
		{
			name:    "should return error if describe ASG fails",
			asgName: aws.String("asgName"),
			wantErr: true,
			wantASG: false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DescribeAutoScalingGroups(gomock.Eq(&autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: []*string{
						aws.String("asgName"),
					},
				})).
					Return(nil, awserrors.NewFailedDependency("unknown error occurred"))
			},
		},
		{
			name:    "should return ASG, if found",
			asgName: aws.String("asgName"),
			wantErr: false,
			wantASG: true,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DescribeAutoScalingGroups(gomock.Eq(&autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: []*string{
						aws.String("asgName"),
					},
				})).
					Return(&autoscaling.DescribeAutoScalingGroupsOutput{
						AutoScalingGroups: []*autoscaling.Group{
							{
								AutoScalingGroupName: aws.String("asgName"),
								MixedInstancesPolicy: &autoscaling.MixedInstancesPolicy{
									InstancesDistribution: &autoscaling.InstancesDistribution{
										OnDemandAllocationStrategy: aws.String("prioritized"),
									},
									LaunchTemplate: &autoscaling.LaunchTemplate{},
								},
							},
						}}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fakeClient := getFakeClient()

			clusterScope, err := getClusterScope(fakeClient)
			g.Expect(err).ToNot(HaveOccurred())
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)
			tt.expect(asgMock.EXPECT())
			s := NewService(clusterScope)
			s.ASGClient = asgMock

			asg, err := s.ASGIfExists(tt.asgName)
			checkErr(tt.wantErr, err, g)
			checkASG(tt.wantASG, asg, g)
		})
	}
}

func TestService_CreateASG(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tests := []struct {
		name                  string
		machinePoolName       string
		setupMachinePoolScope func(*scope.MachinePoolScope)
		wantErr               bool
		wantASG               bool
		expect                func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
	}{
		{
			name:                  "should return without error if create ASG is successful",
			machinePoolName:       "create-asg-success",
			setupMachinePoolScope: func(mps *scope.MachinePoolScope) {},
			wantErr:               false,
			wantASG:               false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				expected := &autoscaling.CreateAutoScalingGroupInput{
					AutoScalingGroupName: aws.String("create-asg-success"),
					CapacityRebalance:    aws.Bool(false),
					DefaultCooldown:      aws.Int64(0),
					MixedInstancesPolicy: &autoscaling.MixedInstancesPolicy{
						InstancesDistribution: &autoscaling.InstancesDistribution{
							OnDemandAllocationStrategy:          aws.String("prioritized"),
							OnDemandBaseCapacity:                aws.Int64(0),
							OnDemandPercentageAboveBaseCapacity: aws.Int64(100),
							SpotAllocationStrategy:              aws.String(""),
						},
						LaunchTemplate: &autoscaling.LaunchTemplate{
							LaunchTemplateSpecification: &autoscaling.LaunchTemplateSpecification{
								LaunchTemplateName: aws.String("create-asg-success"),
								Version:            aws.String("$Latest"),
							},
							Overrides: []*autoscaling.LaunchTemplateOverrides{
								{
									InstanceType: aws.String("t1.large"),
								},
							},
						},
					},
					MaxSize: aws.Int64(0),
					MinSize: aws.Int64(0),
					Tags: []*autoscaling.Tag{
						{
							Key:               aws.String("kubernetes.io/cluster/test"),
							PropagateAtLaunch: aws.Bool(false),
							ResourceId:        aws.String("create-asg-success"),
							ResourceType:      aws.String("auto-scaling-group"),
							Value:             aws.String("owned"),
						},
						{
							Key:               aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test"),
							PropagateAtLaunch: aws.Bool(false),
							ResourceId:        aws.String("create-asg-success"),
							ResourceType:      aws.String("auto-scaling-group"),
							Value:             aws.String("owned"),
						},
						{
							Key:               aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
							PropagateAtLaunch: aws.Bool(false),
							ResourceId:        aws.String("create-asg-success"),
							ResourceType:      aws.String("auto-scaling-group"),
							Value:             aws.String("node"),
						},
						{
							Key:               aws.String("Name"),
							PropagateAtLaunch: aws.Bool(false),
							ResourceId:        aws.String("create-asg-success"),
							ResourceType:      aws.String("auto-scaling-group"),
							Value:             aws.String("create-asg-success"),
						},
					},
					VPCZoneIdentifier: aws.String("subnet1"),
				}

				m.CreateAutoScalingGroup(gomock.AssignableToTypeOf(&autoscaling.CreateAutoScalingGroupInput{})).Do(
					func(actual *autoscaling.CreateAutoScalingGroupInput) (*autoscaling.CreateAutoScalingGroupOutput, error) {
						sortTagsByKey := func(tags []*autoscaling.Tag) {
							sort.Slice(tags, func(i, j int) bool {
								return *(tags[i].Key) < *(tags[j].Key)
							})
						}
						// sorting tags to avoid failure due to different ordering of tags
						sortTagsByKey(actual.Tags)
						sortTagsByKey(expected.Tags)
						if !cmp.Equal(expected, actual) {
							t.Fatalf("Actual CreateAutoScalingGroupInput did not match expected, Actual : %v, Expected: %v", actual, expected)
						}
						return &autoscaling.CreateAutoScalingGroupOutput{}, nil
					})
			},
		},
		{
			name:            "should return error if subnet not found for asg",
			machinePoolName: "create-asg-fail",
			setupMachinePoolScope: func(mps *scope.MachinePoolScope) {
				mps.AWSMachinePool.Spec.Subnets = nil
				mps.InfraCluster.(*scope.ClusterScope).AWSCluster.Spec.NetworkSpec.Subnets = nil
			},
			wantErr: true,
			wantASG: false,
			expect:  func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {},
		},
		{
			name:            "should return error if create ASG fails",
			machinePoolName: "create-asg-fail",
			setupMachinePoolScope: func(mps *scope.MachinePoolScope) {
				mps.AWSMachinePool.Spec.MixedInstancesPolicy = nil
			},
			wantErr: true,
			wantASG: false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.CreateAutoScalingGroup(gomock.AssignableToTypeOf(&autoscaling.CreateAutoScalingGroupInput{})).Return(nil, awserrors.NewFailedDependency("dependency failure"))
			},
		},
		{
			name:            "should return error if launch template is missing",
			machinePoolName: "create-asg-fail",
			setupMachinePoolScope: func(mps *scope.MachinePoolScope) {
				mps.AWSMachinePool.Spec.MixedInstancesPolicy = nil
				mps.AWSMachinePool.Status.LaunchTemplateID = ""
			},
			wantErr: true,
			wantASG: false,
			expect:  func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fakeClient := getFakeClient()

			clusterScope, err := getClusterScope(fakeClient)
			g.Expect(err).ToNot(HaveOccurred())
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)
			tt.expect(asgMock.EXPECT())
			s := NewService(clusterScope)
			s.ASGClient = asgMock

			mps, err := getMachinePoolScope(fakeClient, clusterScope)
			g.Expect(err).ToNot(HaveOccurred())
			mps.AWSMachinePool.Name = tt.machinePoolName
			tt.setupMachinePoolScope(mps)
			asg, err := s.CreateASG(mps)
			checkErr(tt.wantErr, err, g)
			checkASG(tt.wantASG, asg, g)
		})
	}
}

func TestService_UpdateASG(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name                  string
		machinePoolName       string
		setupMachinePoolScope func(*scope.MachinePoolScope)
		wantErr               bool
		expect                func(e *mocks.MockEC2APIMockRecorder, m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
	}{
		{
			name:            "should return without error if update ASG is successful",
			machinePoolName: "update-asg-success",
			wantErr:         false,
			setupMachinePoolScope: func(mps *scope.MachinePoolScope) {
				mps.AWSMachinePool.Spec.Subnets = nil
			},
			expect: func(e *mocks.MockEC2APIMockRecorder, m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.UpdateAutoScalingGroup(gomock.AssignableToTypeOf(&autoscaling.UpdateAutoScalingGroupInput{})).Return(&autoscaling.UpdateAutoScalingGroupOutput{}, nil)
			},
		},
		{
			name:            "should return error if update ASG fails",
			machinePoolName: "update-asg-fail",
			wantErr:         true,
			setupMachinePoolScope: func(mps *scope.MachinePoolScope) {
				mps.AWSMachinePool.Spec.MixedInstancesPolicy = nil
			},
			expect: func(e *mocks.MockEC2APIMockRecorder, m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.UpdateAutoScalingGroup(gomock.AssignableToTypeOf(&autoscaling.UpdateAutoScalingGroupInput{})).Return(nil, awserrors.NewFailedDependency("dependency failure"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fakeClient := getFakeClient()

			clusterScope, err := getClusterScope(fakeClient)
			g.Expect(err).ToNot(HaveOccurred())
			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)
			tt.expect(ec2Mock.EXPECT(), asgMock.EXPECT())
			s := NewService(clusterScope)
			s.ASGClient = asgMock

			mps, err := getMachinePoolScope(fakeClient, clusterScope)
			g.Expect(err).ToNot(HaveOccurred())
			mps.AWSMachinePool.Name = tt.machinePoolName

			err = s.UpdateASG(mps)
			checkErr(tt.wantErr, err, g)
		})
	}
}

func TestService_UpdateASGWithSubnetFilters(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name                 string
		machinePoolName      string
		awsResourceReference []infrav1.AWSResourceReference
		wantErr              bool
		expect               func(e *mocks.MockEC2APIMockRecorder, m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
	}{
		{
			name:            "should return without error if update ASG is successful",
			machinePoolName: "update-asg-success",
			wantErr:         false,
			awsResourceReference: []infrav1.AWSResourceReference{
				{
					Filters: []infrav1.Filter{{Name: "availability-zone", Values: []string{"us-east-1a"}}},
				},
			},
			expect: func(e *mocks.MockEC2APIMockRecorder, m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				e.DescribeSubnets(gomock.AssignableToTypeOf(&ec2.DescribeSubnetsInput{})).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []*ec2.Subnet{{SubnetId: aws.String("subnet-02")}},
				}, nil)
				m.UpdateAutoScalingGroup(gomock.AssignableToTypeOf(&autoscaling.UpdateAutoScalingGroupInput{})).Return(&autoscaling.UpdateAutoScalingGroupOutput{}, nil)
			},
		},
		{
			name:            "should return error if update ASG fails",
			machinePoolName: "update-asg-fail",
			wantErr:         true,
			awsResourceReference: []infrav1.AWSResourceReference{
				{
					ID: aws.String("subnet-01"),
				},
			},
			expect: func(e *mocks.MockEC2APIMockRecorder, m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.UpdateAutoScalingGroup(gomock.AssignableToTypeOf(&autoscaling.UpdateAutoScalingGroupInput{})).Return(nil, awserrors.NewFailedDependency("dependency failure"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fakeClient := getFakeClient()

			clusterScope, err := getClusterScope(fakeClient)
			g.Expect(err).ToNot(HaveOccurred())

			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)
			if tt.expect != nil {
				tt.expect(ec2Mock.EXPECT(), asgMock.EXPECT())
			}
			s := NewService(clusterScope)
			s.ASGClient = asgMock
			s.EC2Client = ec2Mock

			mps, err := getMachinePoolScope(fakeClient, clusterScope)
			g.Expect(err).ToNot(HaveOccurred())
			mps.AWSMachinePool.Name = tt.machinePoolName
			mps.AWSMachinePool.Spec.Subnets = tt.awsResourceReference

			err = s.UpdateASG(mps)
			checkErr(tt.wantErr, err, g)
		})
	}
}

func TestService_UpdateResourceTags(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		resourceID *string
		create     map[string]string
		remove     map[string]string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		expect  func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
	}{
		{
			name: "should return nil if nothing to update",
			args: args{
				resourceID: aws.String("mock-resource-id"),
			},
			wantErr: false,
			expect:  func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {},
		},
		{
			name: "should create tags if new tags are passed",
			args: args{
				resourceID: aws.String("mock-resource-id"),
				create: map[string]string{
					"key1": "value1",
				},
			},
			wantErr: false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.CreateOrUpdateTags(gomock.Eq(&autoscaling.CreateOrUpdateTagsInput{
					Tags: mapToTags(map[string]string{
						"key1": "value1",
					}, aws.String("mock-resource-id")),
				})).
					Return(nil, nil)
			},
		},
		{
			name: "should return error if new tags creation failed",
			args: args{
				resourceID: aws.String("mock-resource-id"),
				create: map[string]string{
					"key1": "value1",
				},
			},
			wantErr: true,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.CreateOrUpdateTags(gomock.Eq(&autoscaling.CreateOrUpdateTagsInput{
					Tags: mapToTags(map[string]string{
						"key1": "value1",
					}, aws.String("mock-resource-id")),
				})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
		},
		{
			name: "should remove tags successfully if tags to be deleted",
			args: args{
				resourceID: aws.String("mock-resource-id"),
				remove: map[string]string{
					"key1": "value1",
				},
			},
			wantErr: false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DeleteTags(gomock.Eq(&autoscaling.DeleteTagsInput{
					Tags: mapToTags(map[string]string{
						"key1": "value1",
					}, aws.String("mock-resource-id")),
				})).
					Return(nil, nil)
			},
		},
		{
			name: "should return error if removing existing tags failed",
			args: args{
				resourceID: aws.String("mock-resource-id"),
				remove: map[string]string{
					"key1": "value1",
				},
			},
			wantErr: true,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DeleteTags(gomock.Eq(&autoscaling.DeleteTagsInput{
					Tags: mapToTags(map[string]string{
						"key1": "value1",
					}, aws.String("mock-resource-id")),
				})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fakeClient := getFakeClient()

			clusterScope, err := getClusterScope(fakeClient)
			g.Expect(err).ToNot(HaveOccurred())
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)
			tt.expect(asgMock.EXPECT())
			s := NewService(clusterScope)
			s.ASGClient = asgMock

			err = s.UpdateResourceTags(tt.args.resourceID, tt.args.create, tt.args.remove)
			checkErr(tt.wantErr, err, g)
		})
	}
}

func TestService_DeleteASG(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name    string
		wantErr bool
		expect  func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
	}{
		{
			name:    "Delete ASG successful",
			wantErr: false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DeleteAutoScalingGroup(gomock.Eq(&autoscaling.DeleteAutoScalingGroupInput{
					AutoScalingGroupName: aws.String("asgName"),
					ForceDelete:          aws.Bool(true),
				})).
					Return(nil, nil)
			},
		},
		{
			name:    "Delete ASG should fail when ASG is not found",
			wantErr: true,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DeleteAutoScalingGroup(gomock.Eq(&autoscaling.DeleteAutoScalingGroupInput{
					AutoScalingGroupName: aws.String("asgName"),
					ForceDelete:          aws.Bool(true),
				})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fakeClient := getFakeClient()

			clusterScope, err := getClusterScope(fakeClient)
			g.Expect(err).ToNot(HaveOccurred())
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)
			tt.expect(asgMock.EXPECT())
			s := NewService(clusterScope)
			s.ASGClient = asgMock

			err = s.DeleteASG("asgName")
			checkErr(tt.wantErr, err, g)
		})
	}
}

func TestService_DeleteASGAndWait(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name    string
		wantErr bool
		expect  func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
	}{
		{
			name:    "Delete ASG with wait passed",
			wantErr: false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DeleteAutoScalingGroup(gomock.Eq(&autoscaling.DeleteAutoScalingGroupInput{
					AutoScalingGroupName: aws.String("asgName"),
					ForceDelete:          aws.Bool(true),
				})).
					Return(nil, nil)
				m.WaitUntilGroupNotExists(gomock.Eq(&autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: aws.StringSlice([]string{"asgName"}),
				})).
					Return(nil)
			},
		},
		{
			name:    "should return error if delete ASG failed while waiting",
			wantErr: true,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DeleteAutoScalingGroup(gomock.Eq(&autoscaling.DeleteAutoScalingGroupInput{
					AutoScalingGroupName: aws.String("asgName"),
					ForceDelete:          aws.Bool(true),
				})).
					Return(nil, nil)
				m.WaitUntilGroupNotExists(gomock.Eq(&autoscaling.DescribeAutoScalingGroupsInput{
					AutoScalingGroupNames: aws.StringSlice([]string{"asgName"}),
				})).
					Return(awserrors.NewFailedDependency("dependency error"))
			},
		},
		{
			name:    "should return error if delete ASG failed during ASG deletion",
			wantErr: true,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DeleteAutoScalingGroup(gomock.Eq(&autoscaling.DeleteAutoScalingGroupInput{
					AutoScalingGroupName: aws.String("asgName"),
					ForceDelete:          aws.Bool(true),
				})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fakeClient := getFakeClient()

			clusterScope, err := getClusterScope(fakeClient)
			g.Expect(err).ToNot(HaveOccurred())
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)
			tt.expect(asgMock.EXPECT())
			s := NewService(clusterScope)
			s.ASGClient = asgMock

			err = s.DeleteASGAndWait("asgName")
			checkErr(tt.wantErr, err, g)
		})
	}
}

func TestService_CanStartASGInstanceRefresh(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name     string
		wantErr  bool
		canStart bool
		expect   func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
	}{
		{
			name:     "should return error if describe instance refresh failed",
			wantErr:  true,
			canStart: false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DescribeInstanceRefreshes(gomock.Eq(&autoscaling.DescribeInstanceRefreshesInput{
					AutoScalingGroupName: aws.String("machinePoolName"),
				})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
		},
		{
			name:     "should return true if no instance available for refresh",
			wantErr:  false,
			canStart: true,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DescribeInstanceRefreshes(gomock.Eq(&autoscaling.DescribeInstanceRefreshesInput{
					AutoScalingGroupName: aws.String("machinePoolName"),
				})).
					Return(&autoscaling.DescribeInstanceRefreshesOutput{}, nil)
			},
		},
		{
			name:     "should return false if some instances have unfinished refresh",
			wantErr:  false,
			canStart: false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.DescribeInstanceRefreshes(gomock.Eq(&autoscaling.DescribeInstanceRefreshesInput{
					AutoScalingGroupName: aws.String("machinePoolName"),
				})).
					Return(&autoscaling.DescribeInstanceRefreshesOutput{
						InstanceRefreshes: []*autoscaling.InstanceRefresh{
							{
								Status: aws.String(autoscaling.InstanceRefreshStatusInProgress),
							},
						},
					}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fakeClient := getFakeClient()

			clusterScope, err := getClusterScope(fakeClient)
			g.Expect(err).ToNot(HaveOccurred())
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)
			tt.expect(asgMock.EXPECT())
			s := NewService(clusterScope)
			s.ASGClient = asgMock

			mps, err := getMachinePoolScope(fakeClient, clusterScope)
			g.Expect(err).ToNot(HaveOccurred())
			mps.AWSMachinePool.Name = "machinePoolName"

			out, err := s.CanStartASGInstanceRefresh(mps)
			checkErr(tt.wantErr, err, g)
			if tt.canStart {
				g.Expect(out).To(BeTrue())
				return
			}
			g.Expect(out).To(BeFalse())
		})
	}
}

func TestService_StartASGInstanceRefresh(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name    string
		wantErr bool
		expect  func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder)
	}{
		{
			name:    "should return error if start instance refresh failed",
			wantErr: true,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.StartInstanceRefresh(gomock.Eq(&autoscaling.StartInstanceRefreshInput{
					AutoScalingGroupName: aws.String("mpn"),
					Strategy:             aws.String("Rolling"),
					Preferences: &autoscaling.RefreshPreferences{
						InstanceWarmup:       aws.Int64(100),
						MinHealthyPercentage: aws.Int64(80),
					},
				})).
					Return(nil, awserrors.NewNotFound("not found"))
			},
		},
		{
			name:    "should return nil if start instance refresh is success",
			wantErr: false,
			expect: func(m *mock_autoscalingiface.MockAutoScalingAPIMockRecorder) {
				m.StartInstanceRefresh(gomock.Eq(&autoscaling.StartInstanceRefreshInput{
					AutoScalingGroupName: aws.String("mpn"),
					Strategy:             aws.String("Rolling"),
					Preferences: &autoscaling.RefreshPreferences{
						InstanceWarmup:       aws.Int64(100),
						MinHealthyPercentage: aws.Int64(80),
					},
				})).
					Return(&autoscaling.StartInstanceRefreshOutput{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fakeClient := getFakeClient()

			clusterScope, err := getClusterScope(fakeClient)
			g.Expect(err).ToNot(HaveOccurred())
			asgMock := mock_autoscalingiface.NewMockAutoScalingAPI(mockCtrl)
			tt.expect(asgMock.EXPECT())
			s := NewService(clusterScope)
			s.ASGClient = asgMock

			mps, err := getMachinePoolScope(fakeClient, clusterScope)
			g.Expect(err).ToNot(HaveOccurred())
			mps.AWSMachinePool.Name = "mpn"

			err = s.StartASGInstanceRefresh(mps)
			checkErr(tt.wantErr, err, g)
		})
	}
}

func getFakeClient() client.Client {
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	_ = expinfrav1.AddToScheme(scheme)
	return fake.NewClientBuilder().WithScheme(scheme).Build()
}

func checkErr(wantErr bool, err error, g *WithT) {
	if wantErr {
		g.Expect(err).To(HaveOccurred())
		return
	}
	g.Expect(err).To(BeNil())
}

func checkASG(wantASG bool, asg *expinfrav1.AutoScalingGroup, g *WithT) {
	if wantASG {
		g.Expect(asg).To(Not(BeNil()))
		return
	}
	g.Expect(asg).To(BeNil())
}

func getClusterScope(client client.Client) (*scope.ClusterScope, error) {
	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	cs, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:  client,
		Cluster: cluster,
		AWSCluster: &infrav1.AWSCluster{
			Spec: infrav1.AWSClusterSpec{
				NetworkSpec: infrav1.NetworkSpec{
					Subnets: []infrav1.SubnetSpec{
						{
							ID: "subnetId",
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func getMachinePoolScope(client client.Client, clusterScope *scope.ClusterScope) (*scope.MachinePoolScope, error) {
	awsMachinePool := &expinfrav1.AWSMachinePool{
		Spec: expinfrav1.AWSMachinePoolSpec{
			Subnets: []infrav1.AWSResourceReference{
				{
					ID: aws.String("subnet1"),
				},
			},
			RefreshPreferences: &expinfrav1.RefreshPreferences{
				Strategy:             aws.String("Rolling"),
				InstanceWarmup:       aws.Int64(100),
				MinHealthyPercentage: aws.Int64(80),
			},
			MixedInstancesPolicy: &expinfrav1.MixedInstancesPolicy{
				InstancesDistribution: &expinfrav1.InstancesDistribution{
					OnDemandAllocationStrategy:          "prioritized",
					OnDemandBaseCapacity:                aws.Int64(0),
					OnDemandPercentageAboveBaseCapacity: aws.Int64(100),
				},
				Overrides: []expinfrav1.Overrides{
					{
						InstanceType: "t1.large",
					},
				},
			},
		},
		Status: expinfrav1.AWSMachinePoolStatus{
			LaunchTemplateID: "launchTemplateID",
		},
	}
	mps, err := scope.NewMachinePoolScope(scope.MachinePoolScopeParams{
		Client:         client,
		Cluster:        clusterScope.Cluster,
		MachinePool:    &expclusterv1.MachinePool{},
		InfraCluster:   clusterScope,
		AWSMachinePool: awsMachinePool,
	})
	if err != nil {
		return nil, err
	}
	return mps, nil
}
