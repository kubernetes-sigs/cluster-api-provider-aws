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
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

func TestReconcilePodIdentityAssociations(t *testing.T) {
	clusterName := "test-cluster"

	tests := []struct {
		name         string
		associations []ekscontrolplanev1.PodIdentityAssociation
		expect       func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError  bool
	}{
		{
			name: "creates new association",
			associations: []ekscontrolplanev1.PodIdentityAssociation{
				{
					ServiceAccountName:      "my-sa",
					ServiceAccountNamespace: "default",
					RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(&eks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{},
					}, nil)
				m.CreatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.CreatePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "updates existing association when role ARN changes",
			associations: []ekscontrolplanev1.PodIdentityAssociation{
				{
					ServiceAccountName:      "my-sa",
					ServiceAccountNamespace: "default",
					RoleARN:                 "arn:aws:iam::123456789012:role/new-role",
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(&eks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{
							{
								AssociationId:  aws.String("assoc-123"),
								ClusterName:    aws.String(clusterName),
								AssociationArn: aws.String("arn:aws:eks:us-west-2:123456789012:podidentityassociation/assoc-123"),
							},
						},
					}, nil)
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/old-role"),
							Tags: map[string]string{
								"kubernetes.io/cluster/test-cluster": "owned",
							},
						},
					}, nil).Times(2)
				m.UpdatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.UpdatePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "recreates association when namespace changes",
			associations: []ekscontrolplanev1.PodIdentityAssociation{
				{
					ServiceAccountName:      "my-sa",
					ServiceAccountNamespace: "new-namespace",
					RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(&eks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{
							{
								AssociationId:  aws.String("assoc-123"),
								ClusterName:    aws.String(clusterName),
								AssociationArn: aws.String("arn:aws:eks:us-west-2:123456789012:podidentityassociation/assoc-123"),
							},
						},
					}, nil)
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("old-namespace"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/my-role"),
							Tags: map[string]string{
								"kubernetes.io/cluster/test-cluster": "owned",
							},
						},
					}, nil)
				m.DeletePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DeletePodIdentityAssociationOutput{}, nil)
				m.CreatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.CreatePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name:         "deletes orphaned association",
			associations: []ekscontrolplanev1.PodIdentityAssociation{},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(&eks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{
							{
								AssociationId:  aws.String("assoc-123"),
								ClusterName:    aws.String(clusterName),
								AssociationArn: aws.String("arn:aws:eks:us-west-2:123456789012:podidentityassociation/assoc-123"),
							},
						},
					}, nil)
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("orphaned-sa"),
							Tags: map[string]string{
								"kubernetes.io/cluster/test-cluster": "owned",
							},
						},
					}, nil)
				m.DeletePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DeletePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "no changes needed",
			associations: []ekscontrolplanev1.PodIdentityAssociation{
				{
					ServiceAccountName:      "my-sa",
					ServiceAccountNamespace: "default",
					RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(&eks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{
							{
								AssociationId:  aws.String("assoc-123"),
								ClusterName:    aws.String(clusterName),
								AssociationArn: aws.String("arn:aws:eks:us-west-2:123456789012:podidentityassociation/assoc-123"),
							},
						},
					}, nil)
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/my-role"),
							Tags: map[string]string{
								"kubernetes.io/cluster/test-cluster": "owned",
							},
						},
					}, nil).Times(2)
			},
			expectError: false,
		},
		{
			name:         "empty spec and no existing associations",
			associations: []ekscontrolplanev1.PodIdentityAssociation{},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(&eks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{},
					}, nil)
			},
			expectError: false,
		},
		{
			name: "error handling list fails",
			associations: []ekscontrolplanev1.PodIdentityAssociation{
				{
					ServiceAccountName:      "my-sa",
					ServiceAccountNamespace: "default",
					RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("list failed"))
			},
			expectError: true,
		},
		{
			name: "error handling create fails",
			associations: []ekscontrolplanev1.PodIdentityAssociation{
				{
					ServiceAccountName:      "my-sa",
					ServiceAccountNamespace: "default",
					RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
				},
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(&eks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{},
					}, nil)
				m.CreatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("create failed"))
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
					ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: "default"},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: "default"},
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName:          clusterName,
						PodIdentityAssociations: tc.associations,
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())

			s := NewService(scope)
			s.EKSClient = eksMock

			err = s.reconcilePodIdentityAssociations(context.TODO())

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestGetManagedPodIdentityAssociations(t *testing.T) {
	clusterName := "test-cluster"

	tests := []struct {
		name        string
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
		expectCount int
	}{
		{
			name: "returns only associations with managed tag",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(&eks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{
							{
								AssociationId:  aws.String("assoc-1"),
								ClusterName:    aws.String(clusterName),
								AssociationArn: aws.String("arn:aws:eks:us-west-2:123456789012:podidentityassociation/assoc-1"),
							},
							{
								AssociationId:  aws.String("assoc-2"),
								ClusterName:    aws.String(clusterName),
								AssociationArn: aws.String("arn:aws:eks:us-west-2:123456789012:podidentityassociation/assoc-2"),
							},
						},
					}, nil)
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Eq(&eks.DescribePodIdentityAssociationInput{
					AssociationId: aws.String("assoc-1"),
					ClusterName:   aws.String(clusterName),
				})).Return(&eks.DescribePodIdentityAssociationOutput{
					Association: &ekstypes.PodIdentityAssociation{
						Namespace:      aws.String("default"),
						ServiceAccount: aws.String("sa1"),
						Tags: map[string]string{
							"kubernetes.io/cluster/test-cluster": "owned",
						},
					},
				}, nil)
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Eq(&eks.DescribePodIdentityAssociationInput{
					AssociationId: aws.String("assoc-2"),
					ClusterName:   aws.String(clusterName),
				})).Return(&eks.DescribePodIdentityAssociationOutput{
					Association: &ekstypes.PodIdentityAssociation{
						Namespace:      aws.String("default"),
						ServiceAccount: aws.String("sa2"),
						Tags:           map[string]string{},
					},
				}, nil)
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "handles pagination",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Eq(&eks.ListPodIdentityAssociationsInput{
					ClusterName: aws.String(clusterName),
					NextToken:   nil,
				})).Return(&eks.ListPodIdentityAssociationsOutput{
					Associations: []ekstypes.PodIdentityAssociationSummary{
						{
							AssociationId:  aws.String("assoc-1"),
							ClusterName:    aws.String(clusterName),
							AssociationArn: aws.String("arn:aws:eks:us-west-2:123456789012:podidentityassociation/assoc-1"),
						},
					},
					NextToken: aws.String("token-1"),
				}, nil)
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("sa1"),
							Tags: map[string]string{
								"kubernetes.io/cluster/test-cluster": "owned",
							},
						},
					}, nil)
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Eq(&eks.ListPodIdentityAssociationsInput{
					ClusterName: aws.String(clusterName),
					NextToken:   aws.String("token-1"),
				})).Return(&eks.ListPodIdentityAssociationsOutput{
					Associations: []ekstypes.PodIdentityAssociationSummary{
						{
							AssociationId:  aws.String("assoc-2"),
							ClusterName:    aws.String(clusterName),
							AssociationArn: aws.String("arn:aws:eks:us-west-2:123456789012:podidentityassociation/assoc-2"),
						},
					},
				}, nil)
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("sa2"),
							Tags: map[string]string{
								"kubernetes.io/cluster/test-cluster": "owned",
							},
						},
					}, nil)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "handles empty results",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(&eks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{},
					}, nil)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "returns error on list failure",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("list failed"))
			},
			expectError: true,
			expectCount: 0,
		},
		{
			name: "continues on describe error",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.ListPodIdentityAssociations(gomock.Any(), gomock.Any()).
					Return(&eks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{
							{
								AssociationId:  aws.String("assoc-1"),
								ClusterName:    aws.String(clusterName),
								AssociationArn: aws.String("arn:aws:eks:us-west-2:123456789012:podidentityassociation/assoc-1"),
							},
							{
								AssociationId:  aws.String("assoc-2"),
								ClusterName:    aws.String(clusterName),
								AssociationArn: aws.String("arn:aws:eks:us-west-2:123456789012:podidentityassociation/assoc-2"),
							},
						},
					}, nil)
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("describe failed"))
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("sa2"),
							Tags: map[string]string{
								"kubernetes.io/cluster/test-cluster": "owned",
							},
						},
					}, nil)
			},
			expectError: false,
			expectCount: 1,
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
					ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: "default"},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: "default"},
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName: clusterName,
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())

			s := NewService(scope)
			s.EKSClient = eksMock

			result, err := s.getManagedPodIdentityAssociations(context.TODO(), clusterName)

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(result).To(HaveLen(tc.expectCount))
			}
		})
	}
}

func TestCreatePodIdentityAssociation(t *testing.T) {
	clusterName := "test-cluster"

	tests := []struct {
		name        string
		association ekscontrolplanev1.PodIdentityAssociation
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
	}{
		{
			name: "creates association with required fields only",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.CreatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.CreatePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "creates association with optional target role ARN",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
				TargetRoleARN:           "arn:aws:iam::123456789012:role/target-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.CreatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.CreatePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "returns error on API failure",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.CreatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("API error"))
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
					ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: "default"},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: "default"},
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName: clusterName,
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())

			s := NewService(scope)
			s.EKSClient = eksMock

			err = s.createPodIdentityAssociation(context.TODO(), tc.association)

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestUpdatePodIdentityAssociation(t *testing.T) {
	clusterName := "test-cluster"
	assocID := "assoc-123"

	tests := []struct {
		name        string
		association ekscontrolplanev1.PodIdentityAssociation
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
	}{
		{
			name: "updates role ARN when changed",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/new-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/old-role"),
						},
					}, nil)
				m.UpdatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.UpdatePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "updates target role ARN when changed",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
				TargetRoleARN:           "arn:aws:iam::123456789012:role/new-target",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/my-role"),
							TargetRoleArn:  aws.String("arn:aws:iam::123456789012:role/old-target"),
						},
					}, nil)
				m.UpdatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.UpdatePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "removes target role ARN when cleared from spec",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/my-role"),
							TargetRoleArn:  aws.String("arn:aws:iam::123456789012:role/old-target"),
						},
					}, nil)
				m.UpdatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.UpdatePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "recreates when namespace changes",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "new-namespace",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("old-namespace"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/my-role"),
						},
					}, nil)
				m.DeletePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DeletePodIdentityAssociationOutput{}, nil)
				m.CreatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.CreatePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "recreates when service account name changes",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "new-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("old-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/my-role"),
						},
					}, nil)
				m.DeletePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DeletePodIdentityAssociationOutput{}, nil)
				m.CreatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.CreatePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "no-op when nothing changed",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/my-role"),
						},
					}, nil)
			},
			expectError: false,
		},
		{
			name: "handles target role ARN being nil in existing association",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/my-role"),
							TargetRoleArn:  nil,
						},
					}, nil)
			},
			expectError: false,
		},
		{
			name: "returns error on describe failure",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("describe failed"))
			},
			expectError: true,
		},
		{
			name: "returns error on update failure",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "default",
				RoleARN:                 "arn:aws:iam::123456789012:role/new-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("default"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/old-role"),
						},
					}, nil)
				m.UpdatePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("update failed"))
			},
			expectError: true,
		},
		{
			name: "returns error on delete recreate failure",
			association: ekscontrolplanev1.PodIdentityAssociation{
				ServiceAccountName:      "my-sa",
				ServiceAccountNamespace: "new-namespace",
				RoleARN:                 "arn:aws:iam::123456789012:role/my-role",
			},
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DescribePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(&eks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							Namespace:      aws.String("old-namespace"),
							ServiceAccount: aws.String("my-sa"),
							RoleArn:        aws.String("arn:aws:iam::123456789012:role/my-role"),
						},
					}, nil)
				m.DeletePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("delete failed"))
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
					ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: "default"},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: "default"},
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName: clusterName,
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())

			s := NewService(scope)
			s.EKSClient = eksMock

			err = s.updatePodIdentityAssociation(context.TODO(), assocID, tc.association)

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}

func TestDeletePodIdentityAssociation(t *testing.T) {
	clusterName := "test-cluster"
	assocID := "assoc-123"

	tests := []struct {
		name        string
		expect      func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectError bool
	}{
		{
			name: "deletes association successfully",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DeletePodIdentityAssociation(gomock.Any(), gomock.Eq(&eks.DeletePodIdentityAssociationInput{
					AssociationId: aws.String(assocID),
					ClusterName:   aws.String(clusterName),
				})).Return(&eks.DeletePodIdentityAssociationOutput{}, nil)
			},
			expectError: false,
		},
		{
			name: "returns error on API failure",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.DeletePodIdentityAssociation(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("API error"))
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
					ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: "default"},
				},
				ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
					ObjectMeta: metav1.ObjectMeta{Name: clusterName, Namespace: "default"},
					Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
						EKSClusterName: clusterName,
					},
				},
			})
			g.Expect(err).To(BeNil())

			tc.expect(eksMock.EXPECT())

			s := NewService(scope)
			s.EKSClient = eksMock

			err = s.deletePodIdentityAssociation(context.TODO(), assocID)

			if tc.expectError {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
