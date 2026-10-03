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
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// TestReconcileNodegroupVersionLaunchTemplate covers launch template handling in reconcileNodegroupVersion,
// including the nil pointer dereference from
// https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/6287.
func TestReconcileNodegroupVersionLaunchTemplate(t *testing.T) {
	const ltID = "lt-0123456789abcdef0"

	tests := []struct {
		name        string
		ngLT        *ekstypes.LaunchTemplateSpecification
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		errContains string
	}{
		{
			name:        "nodegroup without launch template returns error",
			ngLT:        nil,
			expect:      func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			errContains: "no launch template version",
		},
		{
			name:        "nodegroup launch template without version returns error",
			ngLT:        &ekstypes.LaunchTemplateSpecification{Id: aws.String(ltID)},
			expect:      func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			errContains: "no launch template version",
		},
		{
			name:   "same launch template version is a no-op",
			ngLT:   &ekstypes.LaunchTemplateSpecification{Id: aws.String(ltID), Version: aws.String("2")},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
		},
		{
			name: "different launch template version updates the nodegroup",
			ngLT: &ekstypes.LaunchTemplateSpecification{Id: aws.String(ltID), Version: aws.String("1")},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.UpdateNodegroupVersion(gomock.Any(), gomock.Eq(&eks.UpdateNodegroupVersionInput{
					ClusterName:   aws.String("test-cluster"),
					NodegroupName: aws.String("test-nodegroup"),
					LaunchTemplate: &ekstypes.LaunchTemplateSpecification{
						Id:      aws.String(ltID),
						Version: aws.String("2"),
					},
				})).Return(&eks.UpdateNodegroupVersionOutput{}, nil).Times(1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)
			tc.expect(eksMock.EXPECT())

			s := &NodegroupService{
				scope: &scope.ManagedMachinePoolScope{
					ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
						Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{EKSClusterName: "test-cluster"},
					},
					ManagedMachinePool: &expinfrav1.AWSManagedMachinePool{
						Spec: expinfrav1.AWSManagedMachinePoolSpec{EKSNodegroupName: "test-nodegroup"},
						Status: expinfrav1.AWSManagedMachinePoolStatus{
							LaunchTemplateID:      aws.String(ltID),
							LaunchTemplateVersion: aws.String("2"),
						},
					},
					MachinePool: &clusterv1.MachinePool{},
				},
				EKSClient: eksMock,
			}
			ng := &ekstypes.Nodegroup{
				Version:        aws.String("1.30"),
				ReleaseVersion: aws.String("1.30.0-20240101"),
				LaunchTemplate: tc.ngLT,
			}

			var err error
			g.Expect(func() {
				err = s.reconcileNodegroupVersion(context.TODO(), ng)
			}).NotTo(Panic())

			if tc.errContains != "" {
				g.Expect(err).To(MatchError(ContainSubstring(tc.errContains)))
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
		})
	}
}
