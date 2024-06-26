/*
Copyright 2021 The Kubernetes Authors.

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

package podidentities

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
)

func TestEKSPodIdentityAssociationsPlan(t *testing.T) {
	clusterName := "default.cluster"
	namespace := "my-namespace"
	roleArn := "aws://rolearn"
	responseAssociationArn := "aws://association-arn"
	associationID := "aws://association-id"
	serviceAccount := "my-service-account"
	created := time.Now()

	testCases := []struct {
		name              string
		desired           []EKSPodIdentityAssociation
		current           []EKSPodIdentityAssociation
		expect            func(m *mock_eksiface.MockEKSAPIMockRecorder)
		expectCreateError bool
		expectDoError     bool
	}{
		{
			name: "no desired and no current",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				// Do nothing
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "no current and 1 desired",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					CreatePodIdentityAssociation(gomock.Eq(&eks.CreatePodIdentityAssociationInput{
						Namespace:      aws.String(namespace),
						RoleArn:        aws.String(roleArn),
						ServiceAccount: aws.String(serviceAccount),
						ClusterName:    aws.String(clusterName),
					})).
					Return(&eks.CreatePodIdentityAssociationOutput{
						Association: &eks.PodIdentityAssociation{
							AssociationArn: aws.String(responseAssociationArn),
							AssociationId:  aws.String(associationID),
							Namespace:      aws.String(namespace),
							RoleArn:        aws.String(roleArn),
							ServiceAccount: aws.String(serviceAccount),
							ClusterName:    aws.String(clusterName),
							CreatedAt:      &created,
							ModifiedAt:     &created,
						},
					}, nil)
			},
			desired: []EKSPodIdentityAssociation{
				{
					ServiceAccountName:      serviceAccount,
					ServiceAccountNamespace: namespace,
					RoleARN:                 roleArn,
				},
			},
			current:           []EKSPodIdentityAssociation{},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name:   "1 current and 1 desired",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {},
			desired: []EKSPodIdentityAssociation{
				{
					ServiceAccountName:      serviceAccount,
					ServiceAccountNamespace: namespace,
					RoleARN:                 roleArn,
				},
			},
			current: []EKSPodIdentityAssociation{
				{
					ServiceAccountName:      serviceAccount,
					ServiceAccountNamespace: namespace,
					RoleARN:                 roleArn,
				},
			},
			expectCreateError: false,
			expectDoError:     false,
		},
		{
			name: "1 current and 0 desired",
			expect: func(m *mock_eksiface.MockEKSAPIMockRecorder) {
				m.
					DeletePodIdentityAssociation(gomock.Eq(&eks.DeletePodIdentityAssociationInput{
						AssociationId: aws.String(associationID),
						ClusterName:   aws.String(clusterName),
					})).
					Return(&eks.DeletePodIdentityAssociationOutput{
						Association: &eks.PodIdentityAssociation{
							AssociationArn: aws.String(responseAssociationArn),
							AssociationId:  aws.String(associationID),
							Namespace:      aws.String(namespace),
							RoleArn:        aws.String(roleArn),
							ServiceAccount: aws.String(serviceAccount),
							ClusterName:    aws.String(clusterName),
							CreatedAt:      &created,
							ModifiedAt:     &created,
						},
					}, nil)
			},
			desired: []EKSPodIdentityAssociation{},
			current: []EKSPodIdentityAssociation{
				{
					ServiceAccountName:      serviceAccount,
					ServiceAccountNamespace: namespace,
					RoleARN:                 roleArn,
					AssociationID:           associationID,
				},
			},
			expectCreateError: false,
			expectDoError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			mockControl := gomock.NewController(t)
			defer mockControl.Finish()

			eksMock := mock_eksiface.NewMockEKSAPI(mockControl)
			tc.expect(eksMock.EXPECT())

			ctx := context.TODO()

			planner := NewPlan(clusterName, tc.desired, tc.current, eksMock)
			procedures, err := planner.Create(ctx)
			if tc.expectCreateError {
				g.Expect(err).To(HaveOccurred())
				return
			}
			g.Expect(err).To(BeNil())
			g.Expect(procedures).NotTo(BeNil())

			for _, proc := range procedures {
				procErr := proc.Do(ctx)
				if tc.expectDoError {
					g.Expect(procErr).To(HaveOccurred())
					return
				}
				g.Expect(procErr).To(BeNil())
			}
		})
	}
}
