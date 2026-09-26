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

package rosa_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
)

func TestReconcileNotificationContacts(t *testing.T) {
	g := NewWithT(t)
	ctx := context.Background()

	cluster, err := cmv1.NewCluster().
		ID("cluster-1").
		Subscription(cmv1.NewSubscription().ID("sub-1")).
		Build()
	g.Expect(err).NotTo(HaveOccurred())

	t.Run("no-op when notificationContacts is unset", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)

		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{},
			},
		}

		g.Expect(rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, cluster)).To(Succeed())
	})

	t.Run("no-op when desired matches current", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		ocmMock.EXPECT().
			GetSubscriptionNotificationContacts(ctx, "sub-1").
			Return([]string{"user-b", "user-a"}, nil).
			Times(1)

		contacts := []string{"user-a", "user-b"}
		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{
					NotificationContacts: contacts,
				},
			},
		}

		g.Expect(rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, cluster)).To(Succeed())
	})

	t.Run("updates when desired differs from current", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		ocmMock.EXPECT().
			GetSubscriptionNotificationContacts(ctx, "sub-1").
			Return([]string{"old-user"}, nil).
			Times(1)
		ocmMock.EXPECT().
			UpdateSubscriptionNotificationContacts(ctx, "sub-1", []string{"new-user", "user@example.com"}).
			Return(nil).
			Times(1)

		contacts := []string{"new-user", "user@example.com"}
		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{
					NotificationContacts: contacts,
				},
			},
		}

		g.Expect(rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, cluster)).To(Succeed())
	})

	t.Run("clears contacts when set to empty list", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		ocmMock.EXPECT().
			GetSubscriptionNotificationContacts(ctx, "sub-1").
			Return([]string{"old-user"}, nil).
			Times(1)
		ocmMock.EXPECT().
			UpdateSubscriptionNotificationContacts(ctx, "sub-1", gomock.Len(0)).
			Return(nil).
			Times(1)

		contacts := []string{}
		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{
					NotificationContacts: contacts,
				},
			},
		}

		g.Expect(rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, cluster)).To(Succeed())
	})

	t.Run("errors when subscription ID is missing", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)

		noSubCluster, err := cmv1.NewCluster().ID("cluster-1").Build()
		g.Expect(err).NotTo(HaveOccurred())

		contacts := []string{"user-a"}
		rosaScope := &scope.ROSAControlPlaneScope{
			ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
				Spec: rosacontrolplanev1.RosaControlPlaneSpec{
					NotificationContacts: contacts,
				},
			},
		}

		err = rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, noSubCluster)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("subscription ID"))
	})
}
