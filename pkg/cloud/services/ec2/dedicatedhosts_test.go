/*
Copyright 2025 The Kubernetes Authors.

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

package ec2

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func createTestClusterScope(t *testing.T) *scope.ClusterScope {
	t.Helper()
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	scope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:  client,
		Cluster: &clusterv1.Cluster{},
		AWSCluster: &infrav1.AWSCluster{
			ObjectMeta: metav1.ObjectMeta{Name: "test"},
			Spec: infrav1.AWSClusterSpec{
				NetworkSpec: infrav1.NetworkSpec{
					VPC: infrav1.VPCSpec{
						ID: "test-vpc",
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create test context: %v", err)
	}
	return scope
}

func createTestMachineScope(t *testing.T, clusterScope *scope.ClusterScope) *scope.MachineScope {
	t.Helper()
	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	_ = clusterv1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()

	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-machine",
			Namespace: "default",
		},
	}

	awsMachine := &infrav1.AWSMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-aws-machine",
			Namespace: "default",
		},
		Spec: infrav1.AWSMachineSpec{
			InstanceType: "m5.large",
			AdditionalTags: infrav1.Tags{
				"Environment": "test",
				"Owner":       "test-user",
			},
		},
	}

	machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
		Client:       client,
		Cluster:      clusterScope.Cluster,
		Machine:      machine,
		AWSMachine:   awsMachine,
		InfraCluster: clusterScope,
	})
	if err != nil {
		t.Fatalf("Failed to create test machine scope: %v", err)
	}
	return machineScope
}

func TestAllocateDedicatedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name                  string
		dynamicAllocationSpec *infrav1.DynamicHostAllocationSpec
		availabilityZone      string
		expectError           bool
		instanceType          string
		setupMocks            func(m *mocks.MockEC2API)
	}{
		{
			name: "should allocate exactly one dedicated host",
			dynamicAllocationSpec: &infrav1.DynamicHostAllocationSpec{
				Tags: map[string]string{
					"Environment": "production", // This should override the machine's "test" value
					"Purpose":     "dedicated",  // This should be added from dedicated host specific tags
				},
			},
			availabilityZone: "us-west-2a",
			instanceType:     "m5.large",
			expectError:      false,
			setupMocks: func(m *mocks.MockEC2API) {
				// Mock AllocateHosts to return exactly one host
				m.EXPECT().AllocateHosts(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, input *ec2.AllocateHostsInput, optFns ...func(*ec2.Options)) (*ec2.AllocateHostsOutput, error) {
					// Verify that quantity is set to 1
					assert.Equal(t, int32(1), *input.Quantity)

					// Verify that tags are being passed
					assert.NotNil(t, input.TagSpecifications)
					assert.Len(t, input.TagSpecifications, 1)
					assert.Equal(t, types.ResourceTypeDedicatedHost, input.TagSpecifications[0].ResourceType)

					// Verify that only the expected tags are present (no standard cluster/machine tags)
					expectedTags := map[string]string{
						"Environment": "production", // from dedicated host specific tags (overrides machine's "test")
						"Owner":       "test-user",  // from machine AdditionalTags
						"Purpose":     "dedicated",  // from dedicated host specific tags
					}

					// Verify we have exactly the expected number of tags
					assert.Equal(t, len(expectedTags), len(input.TagSpecifications[0].Tags), "Should have exactly the expected number of tags")

					// Verify each expected tag is present with correct value
					for _, tag := range input.TagSpecifications[0].Tags {
						key := aws.ToString(tag.Key)
						value := aws.ToString(tag.Value)
						expectedValue, exists := expectedTags[key]
						assert.True(t, exists, "Unexpected tag found: %s", key)
						assert.Equal(t, expectedValue, value, "Tag %s should have value %s", key, expectedValue)
					}

					return &ec2.AllocateHostsOutput{
						HostIds: []string{"h-1234567890abcdef0"},
					}, nil
				})
			},
		},
		{
			name:                  "should fail if AWS returns multiple hosts",
			dynamicAllocationSpec: &infrav1.DynamicHostAllocationSpec{},
			availabilityZone:      "us-west-2a",
			instanceType:          "m5.large",
			expectError:           true,
			setupMocks: func(m *mocks.MockEC2API) {
				// Mock AllocateHosts to return multiple hosts (should never happen with quantity=1, but test the validation)
				m.EXPECT().AllocateHosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ec2.AllocateHostsOutput{
					HostIds: []string{"h-1234567890abcdef0", "h-0987654321fedcba0"},
				}, nil)
			},
		},
		{
			name:                  "should fail if AWS returns no hosts",
			dynamicAllocationSpec: &infrav1.DynamicHostAllocationSpec{},
			availabilityZone:      "us-west-2a",
			instanceType:          "m5.large",
			expectError:           true,
			setupMocks: func(m *mocks.MockEC2API) {
				// Mock AllocateHosts to return no hosts
				m.EXPECT().AllocateHosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ec2.AllocateHostsOutput{
					HostIds: []string{},
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec2Mock := mocks.NewMockEC2API(mockCtrl)
			tt.setupMocks(ec2Mock)

			clusterScope := createTestClusterScope(t)
			machineScope := createTestMachineScope(t, clusterScope)
			s := NewService(clusterScope)
			s.EC2Client = ec2Mock

			hostID, err := s.AllocateDedicatedHost(context.TODO(), tt.dynamicAllocationSpec, tt.instanceType, tt.availabilityZone, machineScope)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hostID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hostID)
			}
		})
	}
}

func TestDescribeDedicatedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	hostID := "h-1234567890abcdef0"

	host := types.Host{
		HostId:           aws.String(hostID),
		AvailabilityZone: aws.String("us-west-2a"),
		State:            types.AllocationStateAvailable,
		HostProperties: &types.HostProperties{
			InstanceFamily: aws.String("m5"),
			InstanceType:   aws.String("m5.large"),
			TotalVCpus:     aws.Int32(2),
		},
		Instances: []types.HostInstance{},
		Tags: []types.Tag{
			{
				Key:   aws.String("Environment"),
				Value: aws.String("test"),
			},
		},
	}

	ec2Mock := mocks.NewMockEC2API(mockCtrl)
	ec2Mock.EXPECT().DescribeHosts(gomock.Any(), gomock.Any()).Return(&ec2.DescribeHostsOutput{
		Hosts: []types.Host{host},
	}, nil)

	scope := createTestClusterScope(t)
	s := NewService(scope)
	s.EC2Client = ec2Mock

	hostInfo, err := s.DescribeDedicatedHost(context.TODO(), hostID)
	assert.NoError(t, err)
	assert.NotNil(t, hostInfo)
	assert.Equal(t, hostID, hostInfo.HostID)
	assert.Equal(t, "m5", hostInfo.InstanceFamily)
	assert.Equal(t, "m5.large", hostInfo.InstanceType)
	assert.Equal(t, "us-west-2a", hostInfo.AvailabilityZone)
	assert.Equal(t, "available", hostInfo.State)
	assert.Equal(t, int32(2), hostInfo.TotalCapacity)
	assert.Equal(t, int32(2), hostInfo.AvailableCapacity) // No instances running
	assert.Equal(t, "test", hostInfo.Tags["Environment"])
}

func TestAllocateDedicatedHostMultipleMachines(t *testing.T) {
	// This test verifies that multiple machines each get their own dedicated host
	// This is the intended behavior for dedicated hosts - each machine gets complete isolation
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create two machine scopes that would both try to allocate hosts
	clusterScope := createTestClusterScope(t)
	machineScope1 := createTestMachineScope(t, clusterScope)
	machineScope2 := createTestMachineScope(t, clusterScope)

	// Give them different names to simulate different machines
	machineScope1.AWSMachine.Name = "test-machine-1"
	machineScope2.AWSMachine.Name = "test-machine-2"

	ec2Mock := mocks.NewMockEC2API(mockCtrl)

	// Both machines will call AllocateHosts and get separate hosts
	// This is the intended behavior - each machine gets its own dedicated host for isolation
	ec2Mock.EXPECT().AllocateHosts(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, input *ec2.AllocateHostsInput, optFns ...func(*ec2.Options)) (*ec2.AllocateHostsOutput, error) {
		// Verify that quantity is set to 1
		assert.Equal(t, int32(1), *input.Quantity)
		return &ec2.AllocateHostsOutput{
			HostIds: []string{"h-1234567890abcdef0"},
		}, nil
	}).Times(2) // Expect two calls for two machines

	s := NewService(clusterScope)
	s.EC2Client = ec2Mock

	spec := &infrav1.DynamicHostAllocationSpec{
		Tags: map[string]string{
			"Environment": "test",
		},
	}

	// Simulate concurrent allocation (in real scenario, these would be concurrent)
	hostID1, err1 := s.AllocateDedicatedHost(context.TODO(), spec, "m5.large", "us-west-2a", machineScope1)
	hostID2, err2 := s.AllocateDedicatedHost(context.TODO(), spec, "m5.large", "us-west-2a", machineScope2)

	// Both should succeed but get different hosts (demonstrating the race condition)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, hostID1)
	assert.NotEmpty(t, hostID2)
	assert.Equal(t, "h-1234567890abcdef0", hostID1)
	assert.Equal(t, "h-1234567890abcdef0", hostID2) // Same host ID because of mock, but in real scenario they'd be different
}

func TestConvertToHostInfo(t *testing.T) {
	hostID := "h-1234567890abcdef0"

	host := types.Host{
		HostId:           aws.String(hostID),
		AvailabilityZone: aws.String("us-west-2a"),
		State:            types.AllocationStateAvailable,
		HostProperties: &types.HostProperties{
			InstanceFamily: aws.String("m5"),
			InstanceType:   aws.String("m5.large"),
			TotalVCpus:     aws.Int32(4),
		},
		Instances: []types.HostInstance{
			{InstanceId: aws.String("i-1234567890abcdef0")},
		},
		Tags: []types.Tag{
			{
				Key:   aws.String("Environment"),
				Value: aws.String("test"),
			},
		},
	}

	s := &Service{}
	hostInfo := s.convertToHostInfo(host)

	assert.Equal(t, hostID, hostInfo.HostID)
	assert.Equal(t, "m5", hostInfo.InstanceFamily)
	assert.Equal(t, "m5.large", hostInfo.InstanceType)
	assert.Equal(t, "us-west-2a", hostInfo.AvailabilityZone)
	assert.Equal(t, "available", hostInfo.State)
	assert.Equal(t, int32(4), hostInfo.TotalCapacity)
	assert.Equal(t, int32(3), hostInfo.AvailableCapacity) // 1 instance running, 3 available
	assert.Equal(t, "test", hostInfo.Tags["Environment"])
}
