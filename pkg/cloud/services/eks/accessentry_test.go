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

package eks

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	clusterName        = "test-cluster"
	principalARN       = "arn:aws:iam::123456789012:role/my-role"
	secondPrincipalARN = "arn:aws:iam::123456789012:role/second-role"
	policyARN          = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
)

func TestReconcileAccessEntries(t *testing.T) {
	tests := []struct {
		name          string
		accessEntries []ekscontrolplanev1.AccessEntry
		expect        func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError   bool
	}{
		{
			name:          "no access entries",
			accessEntries: []ekscontrolplanev1.AccessEntry{},
			expect:        func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError:   false,
		},
		{
			name: "create new access entry",
			accessEntries: []ekscontrolplanev1.AccessEntry{
				{
					PrincipalARN:     principalARN,
					Type:             ekscontrolplanev1.AccessEntryTypeStandard,
					Username:         "admin",
					KubernetesGroups: []string{"system:masters"},
					AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{
						{
							PolicyARN: policyARN,
							AccessScope: ekscontrolplanev1.AccessScope{
								Type: ekscontrolplanev1.AccessScopeTypeCluster,
							},
						},
					},
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListAccessEntries(gomock.Any(), gomock.Any()).Return(&eks.ListAccessEntriesOutput{
					AccessEntries: []string{},
				}, nil)

				m.CreateAccessEntry(gomock.Any(), &eks.CreateAccessEntryInput{
					ClusterName:      aws.String(clusterName),
					PrincipalArn:     aws.String(principalARN),
					Type:             ekscontrolplanev1.AccessEntryTypeStandard.APIValue(),
					Username:         aws.String("admin"),
					KubernetesGroups: []string{"system:masters"},
				}).Return(&eks.CreateAccessEntryOutput{}, nil)

				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{},
				}, nil)

				m.AssociateAccessPolicy(gomock.Any(), &eks.AssociateAccessPolicyInput{
					ClusterName:  aws.String(clusterName),
					PrincipalArn: aws.String(principalARN),
					PolicyArn:    aws.String(policyARN),
					AccessScope: &ekstypes.AccessScope{
						Type: ekstypes.AccessScopeTypeCluster,
					},
				}).Return(&eks.AssociateAccessPolicyOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "update existing access entry",
			accessEntries: []ekscontrolplanev1.AccessEntry{
				{
					PrincipalARN:     principalARN,
					Type:             ekscontrolplanev1.AccessEntryTypeStandard,
					Username:         "admin-updated",
					KubernetesGroups: []string{"system:masters", "developers"},
					AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{
						{
							PolicyARN: policyARN,
							AccessScope: ekscontrolplanev1.AccessScope{
								Type: ekscontrolplanev1.AccessScopeTypeCluster,
							},
						},
					},
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListAccessEntries(gomock.Any(), gomock.Any()).Return(&eks.ListAccessEntriesOutput{
					AccessEntries: []string{principalARN},
				}, nil)

				m.DescribeAccessEntry(gomock.Any(), gomock.Any()).Return(&eks.DescribeAccessEntryOutput{
					AccessEntry: &ekstypes.AccessEntry{
						PrincipalArn:     aws.String(principalARN),
						Username:         aws.String("admin"),
						Type:             ekscontrolplanev1.AccessEntryTypeStandard.APIValue(),
						KubernetesGroups: []string{"system:masters"},
					},
				}, nil)

				m.UpdateAccessEntry(gomock.Any(), gomock.Any()).Return(&eks.UpdateAccessEntryOutput{}, nil)

				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{
						{
							PolicyArn: aws.String(policyARN),
							AccessScope: &ekstypes.AccessScope{
								Type: ekstypes.AccessScopeTypeCluster,
							},
						},
					},
				}, nil)

				m.AssociateAccessPolicy(gomock.Any(), &eks.AssociateAccessPolicyInput{
					ClusterName:  aws.String(clusterName),
					PrincipalArn: aws.String(principalARN),
					PolicyArn:    aws.String(policyARN),
					AccessScope: &ekstypes.AccessScope{
						Type: ekstypes.AccessScopeTypeCluster,
					},
				}).Return(&eks.AssociateAccessPolicyOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "delete access entry",
			accessEntries: []ekscontrolplanev1.AccessEntry{
				{
					PrincipalARN:     principalARN,
					Type:             "STANDARD",
					Username:         "admin",
					KubernetesGroups: []string{"system:masters"},
					AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{
						{
							PolicyARN: policyARN,
							AccessScope: ekscontrolplanev1.AccessScope{
								Type: "cluster",
							},
						},
					},
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListAccessEntries(gomock.Any(), gomock.Any()).Return(&eks.ListAccessEntriesOutput{
					AccessEntries: []string{principalARN, secondPrincipalARN},
				}, nil)

				m.DescribeAccessEntry(gomock.Any(), gomock.Any()).Return(&eks.DescribeAccessEntryOutput{
					AccessEntry: &ekstypes.AccessEntry{
						PrincipalArn:     aws.String(principalARN),
						Username:         aws.String("admin"),
						Type:             ekscontrolplanev1.AccessEntryTypeStandard.APIValue(),
						KubernetesGroups: []string{"system:masters"},
					},
				}, nil)

				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{
						{
							PolicyArn: aws.String(policyARN),
							AccessScope: &ekstypes.AccessScope{
								Type: ekstypes.AccessScopeTypeCluster,
							},
						},
					},
				}, nil)

				m.AssociateAccessPolicy(gomock.Any(), &eks.AssociateAccessPolicyInput{
					ClusterName:  aws.String(clusterName),
					PrincipalArn: aws.String(principalARN),
					PolicyArn:    aws.String(policyARN),
					AccessScope: &ekstypes.AccessScope{
						Type: ekstypes.AccessScopeTypeCluster,
					},
				}).Return(&eks.AssociateAccessPolicyOutput{}, nil)

				m.DeleteAccessEntry(gomock.Any(), &eks.DeleteAccessEntryInput{
					ClusterName:  aws.String(clusterName),
					PrincipalArn: aws.String(secondPrincipalARN),
				}).Return(&eks.DeleteAccessEntryOutput{}, nil)
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      clusterName,
				},
				Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
					EKSClusterName: clusterName,
					AccessConfig: &ekscontrolplanev1.AccessConfig{
						AuthenticationMode: ekscontrolplanev1.EKSAuthenticationModeAPIAndConfigMap,
					},
					AccessEntries: tc.accessEntries,
				},
			}

			scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
				},
				ControlPlane: controlPlane,
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())
			s := NewService(scope)
			s.EKSClient = eksMock

			err = s.reconcileAccessEntries(context.TODO())
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

func TestReconcileAccessPolicies(t *testing.T) {
	tests := []struct {
		name        string
		accessEntry ekscontrolplanev1.AccessEntry
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
	}{
		{
			name: "ec2_linux Type skips policy reconciliation",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN: principalARN,
				Type:         ekscontrolplanev1.AccessEntryTypeEC2Linux,
			},
			expect:      func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError: false,
		},
		{
			name: "ec2_windows Type skips policy reconciliation",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN: principalARN,
				Type:         ekscontrolplanev1.AccessEntryTypeEC2Windows,
			},
			expect:      func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			expectError: false,
		},
		{
			name: "associate new policy",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN: principalARN,
				Type:         ekscontrolplanev1.AccessEntryTypeStandard,
				AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{
					{
						PolicyARN: policyARN,
						AccessScope: ekscontrolplanev1.AccessScope{
							Type: ekscontrolplanev1.AccessScopeTypeCluster,
						},
					},
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{},
				}, nil)

				m.AssociateAccessPolicy(gomock.Any(), &eks.AssociateAccessPolicyInput{
					ClusterName:  aws.String(clusterName),
					PrincipalArn: aws.String(principalARN),
					PolicyArn:    aws.String(policyARN),
					AccessScope: &ekstypes.AccessScope{
						Type: ekstypes.AccessScopeTypeCluster,
					},
				}).Return(&eks.AssociateAccessPolicyOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "disassociate policy",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN:   principalARN,
				Type:           ekscontrolplanev1.AccessEntryTypeStandard,
				AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{
						{
							PolicyArn: aws.String(policyARN),
							AccessScope: &ekstypes.AccessScope{
								Type: ekstypes.AccessScopeTypeCluster,
							},
						},
					},
				}, nil)

				m.DisassociateAccessPolicy(gomock.Any(), &eks.DisassociateAccessPolicyInput{
					ClusterName:  aws.String(clusterName),
					PrincipalArn: aws.String(principalARN),
					PolicyArn:    aws.String(policyARN),
				}).Return(&eks.DisassociateAccessPolicyOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "namespace scoped policy",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN: principalARN,
				Type:         ekscontrolplanev1.AccessEntryTypeStandard,
				AccessPolicies: []ekscontrolplanev1.AccessPolicyReference{
					{
						PolicyARN: policyARN,
						AccessScope: ekscontrolplanev1.AccessScope{
							Type:       ekscontrolplanev1.AccessScopeTypeNamespace,
							Namespaces: []string{"kube-system", "default"},
						},
					},
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{},
				}, nil)

				m.AssociateAccessPolicy(gomock.Any(), &eks.AssociateAccessPolicyInput{
					ClusterName:  aws.String(clusterName),
					PrincipalArn: aws.String(principalARN),
					PolicyArn:    aws.String(policyARN),
					AccessScope: &ekstypes.AccessScope{
						Type:       ekstypes.AccessScopeTypeNamespace,
						Namespaces: []string{"kube-system", "default"},
					},
				}).Return(&eks.AssociateAccessPolicyOutput{}, nil)
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName: clusterName,
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())
			s := NewService(scope)
			s.EKSClient = eksMock

			err = s.reconcileAccessPolicies(context.TODO(), tc.accessEntry)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

func TestCreateAccessEntry(t *testing.T) {
	tests := []struct {
		name        string
		accessEntry ekscontrolplanev1.AccessEntry
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
	}{
		{
			name: "basic access entry",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN: principalARN,
				Type:         ekscontrolplanev1.AccessEntryTypeStandard,
				Username:     "admin",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.CreateAccessEntry(gomock.Any(), &eks.CreateAccessEntryInput{
					ClusterName:  aws.String(clusterName),
					PrincipalArn: aws.String(principalARN),
					Type:         ekscontrolplanev1.AccessEntryTypeStandard.APIValue(),
					Username:     aws.String("admin"),
				}).Return(&eks.CreateAccessEntryOutput{}, nil)

				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{},
				}, nil)
			},
			expectError: false,
		},
		{
			name: "access entry with groups",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN:     principalARN,
				Type:             ekscontrolplanev1.AccessEntryTypeStandard,
				Username:         "admin",
				KubernetesGroups: []string{"system:masters", "developers"},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.CreateAccessEntry(gomock.Any(), &eks.CreateAccessEntryInput{
					ClusterName:      aws.String(clusterName),
					PrincipalArn:     aws.String(principalARN),
					Type:             ekscontrolplanev1.AccessEntryTypeStandard.APIValue(),
					Username:         aws.String("admin"),
					KubernetesGroups: []string{"system:masters", "developers"},
				}).Return(&eks.CreateAccessEntryOutput{}, nil)

				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{},
				}, nil)
			},
			expectError: false,
		},
		{
			name: "api error",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN:     principalARN,
				Type:             ekscontrolplanev1.AccessEntryTypeStandard,
				Username:         "admin",
				KubernetesGroups: []string{"system:masters"},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.CreateAccessEntry(gomock.Any(), gomock.Any()).Return(nil, &ekstypes.InvalidParameterException{
					Message: aws.String("error creating access entry"),
				})
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName: clusterName,
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())
			s := NewService(scope)
			s.EKSClient = eksMock

			err = s.createAccessEntry(context.TODO(), tc.accessEntry)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

func TestUpdateAccessEntry(t *testing.T) {
	tests := []struct {
		name        string
		accessEntry ekscontrolplanev1.AccessEntry
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
	}{
		{
			name: "no updates needed",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN:     principalARN,
				Type:             ekscontrolplanev1.AccessEntryTypeStandard,
				Username:         "admin",
				KubernetesGroups: []string{"system:masters"},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribeAccessEntry(gomock.Any(), gomock.Any()).Return(&eks.DescribeAccessEntryOutput{
					AccessEntry: &ekstypes.AccessEntry{
						PrincipalArn:     aws.String(principalARN),
						Type:             ekscontrolplanev1.AccessEntryTypeStandard.APIValue(),
						Username:         aws.String("admin"),
						KubernetesGroups: []string{"system:masters"},
					},
				}, nil)

				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{},
				}, nil)
			},
			expectError: false,
		},
		{
			name: "type change requires recreate",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN:     principalARN,
				Type:             ekscontrolplanev1.AccessEntryTypeFargateLinux,
				Username:         "admin",
				KubernetesGroups: []string{"system:masters"},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribeAccessEntry(gomock.Any(), gomock.Any()).Return(&eks.DescribeAccessEntryOutput{
					AccessEntry: &ekstypes.AccessEntry{
						PrincipalArn:     aws.String(principalARN),
						Type:             ekscontrolplanev1.AccessEntryTypeStandard.APIValue(),
						Username:         aws.String("admin"),
						KubernetesGroups: []string{"system:masters"},
					},
				}, nil)

				m.DeleteAccessEntry(gomock.Any(), gomock.Any()).Return(&eks.DeleteAccessEntryOutput{}, nil)

				m.CreateAccessEntry(gomock.Any(), gomock.Any()).Return(&eks.CreateAccessEntryOutput{}, nil)

				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{},
				}, nil)
			},
			expectError: false,
		},
		{
			name: "username update",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN:     principalARN,
				Type:             ekscontrolplanev1.AccessEntryTypeStandard,
				Username:         "new-admin",
				KubernetesGroups: []string{"system:masters"},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribeAccessEntry(gomock.Any(), gomock.Any()).Return(&eks.DescribeAccessEntryOutput{
					AccessEntry: &ekstypes.AccessEntry{
						PrincipalArn:     aws.String(principalARN),
						Type:             ekscontrolplanev1.AccessEntryTypeStandard.APIValue(),
						Username:         aws.String("admin"),
						KubernetesGroups: []string{"system:masters"},
					},
				}, nil)

				m.UpdateAccessEntry(gomock.Any(), &eks.UpdateAccessEntryInput{
					ClusterName:  aws.String(clusterName),
					PrincipalArn: aws.String(principalARN),
					Username:     aws.String("new-admin"),
				}).Return(&eks.UpdateAccessEntryOutput{}, nil)

				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{},
				}, nil)
			},
			expectError: false,
		},
		{
			name: "kubernetes groups update",
			accessEntry: ekscontrolplanev1.AccessEntry{
				PrincipalARN:     principalARN,
				Type:             ekscontrolplanev1.AccessEntryTypeStandard,
				Username:         "admin",
				KubernetesGroups: []string{"developers"},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribeAccessEntry(gomock.Any(), gomock.Any()).Return(&eks.DescribeAccessEntryOutput{
					AccessEntry: &ekstypes.AccessEntry{
						PrincipalArn:     aws.String(principalARN),
						Type:             ekscontrolplanev1.AccessEntryTypeStandard.APIValue(),
						Username:         aws.String("admin"),
						KubernetesGroups: []string{"system:masters"},
					},
				}, nil)

				m.UpdateAccessEntry(gomock.Any(), &eks.UpdateAccessEntryInput{
					ClusterName:      aws.String(clusterName),
					PrincipalArn:     aws.String(principalARN),
					KubernetesGroups: []string{"developers"},
				}).Return(&eks.UpdateAccessEntryOutput{}, nil)

				m.ListAssociatedAccessPolicies(gomock.Any(), gomock.Any()).Return(&eks.ListAssociatedAccessPoliciesOutput{
					AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{},
				}, nil)
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName: clusterName,
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())
			s := NewService(scope)
			s.EKSClient = eksMock

			err = s.updateAccessEntry(context.TODO(), tc.accessEntry)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

func TestDeleteAccessEntry(t *testing.T) {
	tests := []struct {
		name        string
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
	}{
		{
			name: "successful delete",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DeleteAccessEntry(gomock.Any(), &eks.DeleteAccessEntryInput{
					ClusterName:  aws.String(clusterName),
					PrincipalArn: aws.String(principalARN),
				}).Return(&eks.DeleteAccessEntryOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "api error",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DeleteAccessEntry(gomock.Any(), gomock.Any()).Return(nil, &ekstypes.ResourceNotFoundException{
					Message: aws.String("access entry not found"),
				})
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)

			scheme := runtime.NewScheme()
			_ = infrav1.AddToScheme(scheme)
			_ = ekscontrolplanev1.AddToScheme(scheme)
			client := fake.NewClientBuilder().WithScheme(scheme).Build()

			scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
				Client: client,
				Cluster: &clusterv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "ns",
						Name:      clusterName,
					},
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName: clusterName,
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())
			s := NewService(scope)
			s.EKSClient = eksMock

			err = s.deleteAccessEntry(context.TODO(), principalARN)
			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
		})
	}
}

