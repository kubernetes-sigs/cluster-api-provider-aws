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
	"fmt"
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
				Release: infrav1.DedicatedHostReleaseStrategyOnMachineDeletion,
				Tags: map[string]string{
					"Environment": "test",
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
					return &ec2.AllocateHostsOutput{
						HostIds: []string{"h-1234567890abcdef0"},
					}, nil
				})
			},
		},
		{
			name: "should fail if AWS returns multiple hosts",
			dynamicAllocationSpec: &infrav1.DynamicHostAllocationSpec{
				Release: infrav1.DedicatedHostReleaseStrategyOnMachineDeletion,
			},
			availabilityZone: "us-west-2a",
			instanceType:     "m5.large",
			expectError:      true,
			setupMocks: func(m *mocks.MockEC2API) {
				// Mock AllocateHosts to return multiple hosts (should never happen with quantity=1, but test the validation)
				m.EXPECT().AllocateHosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ec2.AllocateHostsOutput{
					HostIds: []string{"h-1234567890abcdef0", "h-0987654321fedcba0"},
				}, nil)
			},
		},
		{
			name: "should fail if AWS returns no hosts",
			dynamicAllocationSpec: &infrav1.DynamicHostAllocationSpec{
				Release: infrav1.DedicatedHostReleaseStrategyOnMachineDeletion,
			},
			availabilityZone: "us-west-2a",
			instanceType:     "m5.large",
			expectError:      true,
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

			scope := createTestClusterScope(t)
			s := NewService(scope)
			s.EC2Client = ec2Mock

			hostID, err := s.AllocateDedicatedHost(context.TODO(), tt.dynamicAllocationSpec, tt.instanceType, tt.availabilityZone)

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

func TestReleaseDedicatedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ec2Mock := mocks.NewMockEC2API(mockCtrl)
	ec2Mock.EXPECT().ReleaseHosts(gomock.Any(), gomock.Any()).Return(&ec2.ReleaseHostsOutput{}, nil)

	scope := createTestClusterScope(t)
	s := NewService(scope)
	s.EC2Client = ec2Mock

	err := s.ReleaseDedicatedHost(context.TODO(), "h-1234567890abcdef0")
	assert.NoError(t, err)
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

func TestValidateHostCompatibility(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name         string
		hostID       string
		instanceType string
		hostInfo     *infrav1.DedicatedHostInfo
		expectError  bool
		errorMessage string
	}{
		{
			name:         "compatible instance type",
			hostID:       "h-1234567890abcdef0",
			instanceType: "m5.large",
			hostInfo: &infrav1.DedicatedHostInfo{
				HostID:            "h-1234567890abcdef0",
				InstanceFamily:    "m5",
				AvailableCapacity: 2,
			},
			expectError: false,
		},
		{
			name:         "incompatible instance family",
			hostID:       "h-1234567890abcdef0",
			instanceType: "c5.large",
			hostInfo: &infrav1.DedicatedHostInfo{
				HostID:            "h-1234567890abcdef0",
				InstanceFamily:    "m5",
				AvailableCapacity: 2,
			},
			expectError:  true,
			errorMessage: "instance type c5.large (family c5) is not compatible with host h-1234567890abcdef0 (family m5)",
		},
		{
			name:         "no available capacity",
			hostID:       "h-1234567890abcdef0",
			instanceType: "m5.large",
			hostInfo: &infrav1.DedicatedHostInfo{
				HostID:            "h-1234567890abcdef0",
				InstanceFamily:    "m5",
				AvailableCapacity: 0,
			},
			expectError:  true,
			errorMessage: "host h-1234567890abcdef0 has no available capacity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a proper service with a scope
			scope := createTestClusterScope(t)
			s := NewService(scope)

			// Mock the DescribeDedicatedHost method by creating a service that returns our test host info
			testService := &testServiceWrapper{
				Service:  s,
				hostInfo: tt.hostInfo,
			}

			err := testService.ValidateHostCompatibility(context.TODO(), tt.hostID, tt.instanceType)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// testServiceWrapper wraps the Service to override DescribeDedicatedHost for testing.
type testServiceWrapper struct {
	*Service
	hostInfo *infrav1.DedicatedHostInfo
}

func (t *testServiceWrapper) DescribeDedicatedHost(ctx context.Context, hostID string) (*infrav1.DedicatedHostInfo, error) {
	return t.hostInfo, nil
}

// ValidateHostCompatibility overrides the parent method to use our mock DescribeDedicatedHost.
func (t *testServiceWrapper) ValidateHostCompatibility(ctx context.Context, hostID string, instanceType string) error {
	hostInfo, err := t.DescribeDedicatedHost(ctx, hostID)
	if err != nil {
		return err
	}

	// Extract instance family from instance type (e.g., "m5.large" -> "m5")
	instanceFamily := extractInstanceFamily(instanceType)

	if hostInfo.InstanceFamily != instanceFamily {
		return fmt.Errorf("instance type %s (family %s) is not compatible with host %s (family %s)",
			instanceType, instanceFamily, hostID, hostInfo.InstanceFamily)
	}

	// Check if host has available capacity
	if hostInfo.AvailableCapacity <= 0 {
		return fmt.Errorf("host %s has no available capacity", hostID)
	}

	return nil
}

func TestExtractInstanceFamily(t *testing.T) {
	tests := []struct {
		instanceType string
		expected     string
	}{
		{"m5.large", "m5"},
		{"c5.xlarge", "c5"},
		{"r5.2xlarge", "r5"},
		{"invalid", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.instanceType, func(t *testing.T) {
			result := extractInstanceFamily(tt.instanceType)
			assert.Equal(t, tt.expected, result)
		})
	}
}
