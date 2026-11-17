/*
Copyright 2026 The Kubernetes Authors.

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

package eks

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// TestCanPruneLaunchTemplateVersions: prune only when the nodegroup is ACTIVE on the status LT ID and version; always report the in-use version.
func TestCanPruneLaunchTemplateVersions(t *testing.T) {
	lt := func(id, version string) *ekstypes.LaunchTemplateSpecification {
		return &ekstypes.LaunchTemplateSpecification{Id: aws.String(id), Version: aws.String(version)}
	}
	tests := []struct {
		name          string
		statusVersion *string
		ng            *ekstypes.Nodegroup
		describeErr   error
		want          bool
		wantInUse     []string
		wantErr       bool
	}{
		{
			name:          "ACTIVE on status version: prune allowed, status version in use",
			statusVersion: aws.String("7"),
			ng:            &ekstypes.Nodegroup{Status: ekstypes.NodegroupStatusActive, LaunchTemplate: lt("lt-1", "7")},
			want:          true,
			wantInUse:     []string{"7"},
		},
		{
			name:          "UPDATING rollout in flight: prune blocked",
			statusVersion: aws.String("7"),
			ng:            &ekstypes.Nodegroup{Status: ekstypes.NodegroupStatusUpdating, LaunchTemplate: lt("lt-1", "6")},
			want:          false,
			wantInUse:     []string{"6"},
		},
		{
			name:          "ACTIVE but still on older version (failed/stuck rollout): prune blocked",
			statusVersion: aws.String("7"),
			ng:            &ekstypes.Nodegroup{Status: ekstypes.NodegroupStatusActive, LaunchTemplate: lt("lt-1", "6")},
			want:          false,
			wantInUse:     []string{"6"},
		},
		{
			name:          "stale status (v7 created, status still v6), ACTIVE on v6: prune allowed but v6 reported in use",
			statusVersion: aws.String("6"),
			ng:            &ekstypes.Nodegroup{Status: ekstypes.NodegroupStatusActive, LaunchTemplate: lt("lt-1", "6")},
			want:          true,
			wantInUse:     []string{"6"},
		},
		{
			name:          "DEGRADED on status version: prune blocked",
			statusVersion: aws.String("7"),
			ng:            &ekstypes.Nodegroup{Status: ekstypes.NodegroupStatusDegraded, LaunchTemplate: lt("lt-1", "7")},
			want:          false,
			wantInUse:     []string{"7"},
		},
		{
			name:          "ACTIVE on a different launch template ID: prune blocked, nothing of ours in use",
			statusVersion: aws.String("7"),
			ng:            &ekstypes.Nodegroup{Status: ekstypes.NodegroupStatusActive, LaunchTemplate: lt("lt-other", "7")},
			want:          false,
		},
		{
			name:          "ACTIVE without launch template: prune blocked",
			statusVersion: aws.String("7"),
			ng:            &ekstypes.Nodegroup{Status: ekstypes.NodegroupStatusActive},
			want:          false,
		},
		{
			name:          "nodegroup not found (first create): prune allowed",
			statusVersion: aws.String("7"),
			describeErr:   &ekstypes.ResourceNotFoundException{Message: aws.String("not found")},
			want:          true,
		},
		{
			name:          "describe error is propagated",
			statusVersion: aws.String("7"),
			describeErr:   &ekstypes.ServerException{Message: aws.String("boom")},
			wantErr:       true,
		},
		{
			name:      "no status version yet: prune allowed, in-use version still reported",
			ng:        &ekstypes.Nodegroup{Status: ekstypes.NodegroupStatusActive, LaunchTemplate: lt("lt-1", "6")},
			want:      true,
			wantInUse: []string{"6"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()
			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)
			eksMock.EXPECT().
				DescribeNodegroup(gomock.Any(), gomock.Eq(&eks.DescribeNodegroupInput{ClusterName: aws.String("test-cluster"), NodegroupName: aws.String("test-ng")})).
				Return(&eks.DescribeNodegroupOutput{Nodegroup: tc.ng}, tc.describeErr)

			log := logger.NewLogger(logr.Discard())
			s := &NodegroupService{
				scope: &scope.ManagedMachinePoolScope{
					Logger: *log,
					ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
						Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{EKSClusterName: "test-cluster"},
					},
					ManagedMachinePool: &expinfrav1.AWSManagedMachinePool{
						Spec:   expinfrav1.AWSManagedMachinePoolSpec{EKSNodegroupName: "test-ng"},
						Status: expinfrav1.AWSManagedMachinePoolStatus{LaunchTemplateID: aws.String("lt-1"), LaunchTemplateVersion: tc.statusVersion},
					},
					MachinePool: &clusterv1.MachinePool{},
				},
				EKSClient: eksMock,
			}

			got, inUse, err := s.CanPruneLaunchTemplateVersions(context.TODO())
			if tc.wantErr {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(got).To(Equal(tc.want))
			g.Expect(inUse).To(Equal(tc.wantInUse))
		})
	}
}
