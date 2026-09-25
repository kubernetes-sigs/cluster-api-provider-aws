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
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"k8s.io/apimachinery/pkg/util/sets"

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

	t.Run("adds new contacts before removing stale ones", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		ocmMock.EXPECT().
			GetSubscriptionNotificationContacts(ctx, "sub-1").
			Return([]string{"old-user"}, nil).
			Times(1)
		gomock.InOrder(
			ocmMock.EXPECT().
				UpdateSubscriptionNotificationContacts(ctx, "sub-1", []string{"new-user", "old-user", "user@example.com"}).
				Return(nil).
				Times(1),
			ocmMock.EXPECT().
				UpdateSubscriptionNotificationContacts(ctx, "sub-1", []string{"new-user", "user@example.com"}).
				Return(nil).
				Times(1),
		)

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

	t.Run("only adds when no contacts need removing", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		ocmMock.EXPECT().
			GetSubscriptionNotificationContacts(ctx, "sub-1").
			Return([]string{"user-a"}, nil).
			Times(1)
		ocmMock.EXPECT().
			UpdateSubscriptionNotificationContacts(ctx, "sub-1", []string{"user-a", "user-b"}).
			Return(nil).
			Times(1)

		rosaScope := newNotificationContactsScope([]string{"user-a", "user-b"})
		g.Expect(rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, cluster)).To(Succeed())
	})

	t.Run("does not remove contacts when additions fail", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		ocmMock.EXPECT().
			GetSubscriptionNotificationContacts(ctx, "sub-1").
			Return([]string{"old-user"}, nil).
			Times(1)
		// A single update that keeps old-user; no follow-up call that would remove it.
		ocmMock.EXPECT().
			UpdateSubscriptionNotificationContacts(ctx, "sub-1", []string{"invalid-user", "old-user"}).
			Return(errors.New("failed to add 1 of 1 notification contact(s)")).
			Times(1)

		rosaScope := newNotificationContactsScope([]string{"invalid-user"})
		err := rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, cluster)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("existing contacts were preserved"))
		g.Expect(err.Error()).To(ContainSubstring("invalid-user"))
	})

	t.Run("returns error when removal fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		ocmMock.EXPECT().
			GetSubscriptionNotificationContacts(ctx, "sub-1").
			Return([]string{"user-a", "user-b"}, nil).
			Times(1)
		ocmMock.EXPECT().
			UpdateSubscriptionNotificationContacts(ctx, "sub-1", []string{"user-a"}).
			Return(errors.New("can't remove notification contact 'user-b'")).
			Times(1)

		rosaScope := newNotificationContactsScope([]string{"user-a"})
		err := rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, cluster)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("failed to remove notification contacts"))
	})

	t.Run("no-op when set to empty and no contacts exist", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		ocmMock.EXPECT().
			GetSubscriptionNotificationContacts(ctx, "sub-1").
			Return(nil, nil).
			Times(1)

		rosaScope := newNotificationContactsScope([]string{})
		g.Expect(rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, cluster)).To(Succeed())
	})

	t.Run("no-op when desired only differs by surrounding whitespace", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		ocmMock.EXPECT().
			GetSubscriptionNotificationContacts(ctx, "sub-1").
			Return([]string{"user-a"}, nil).
			Times(1)

		rosaScope := newNotificationContactsScope([]string{" user-a "})
		g.Expect(rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, cluster)).To(Succeed())
	})

	t.Run("sends trimmed usernames when updating", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ocmMock := mocks.NewMockOCMClient(mockCtrl)
		ocmMock.EXPECT().
			GetSubscriptionNotificationContacts(ctx, "sub-1").
			Return([]string{"user-a"}, nil).
			Times(1)
		gomock.InOrder(
			ocmMock.EXPECT().
				UpdateSubscriptionNotificationContacts(ctx, "sub-1", []string{"user-a", "user-b"}).
				Return(nil).
				Times(1),
			ocmMock.EXPECT().
				UpdateSubscriptionNotificationContacts(ctx, "sub-1", []string{"user-b"}).
				Return(nil).
				Times(1),
		)

		rosaScope := newNotificationContactsScope([]string{"\tuser-b "})
		g.Expect(rosa.ReconcileNotificationContacts(ctx, rosaScope, ocmMock, cluster)).To(Succeed())
	})
}

// TestReconcileNotificationContactsPreservesRecipients runs the reconciler against a
// fake that mirrors the github.com/openshift/rosa UpdateSubscriptionNotificationContacts
// semantics: additions are attempted individually and failures are accumulated, but
// deletions still proceed before the error is returned.
func TestReconcileNotificationContactsPreservesRecipients(t *testing.T) {
	ctx := context.Background()

	cluster, err := cmv1.NewCluster().
		ID("cluster-1").
		Subscription(cmv1.NewSubscription().ID("sub-1")).
		Build()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		current     []string
		desired     []string
		failingAdds []string
		wantErr     bool
		want        []string
	}{
		{
			name:        "complete addition failure keeps existing contact",
			current:     []string{"old-user"},
			desired:     []string{"invalid-user"},
			failingAdds: []string{"invalid-user"},
			wantErr:     true,
			want:        []string{"old-user"},
		},
		{
			name:        "partial addition failure keeps existing contacts and successful additions",
			current:     []string{"old-user-1", "old-user-2"},
			desired:     []string{"new-user", "invalid-user"},
			failingAdds: []string{"invalid-user"},
			wantErr:     true,
			want:        []string{"new-user", "old-user-1", "old-user-2"},
		},
		{
			name:    "successful replacement removes stale contacts",
			current: []string{"old-user"},
			desired: []string{"new-user"},
			want:    []string{"new-user"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			fake := &fakeNotificationContactsClient{
				contacts:    sets.New(tt.current...),
				failingAdds: sets.New(tt.failingAdds...),
			}

			err := rosa.ReconcileNotificationContacts(ctx, newNotificationContactsScope(tt.desired), fake, cluster)
			if tt.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
			g.Expect(sets.List(fake.contacts)).To(Equal(tt.want))
		})
	}

	t.Run("retry after transient failure converges", func(t *testing.T) {
		g := NewWithT(t)
		fake := &fakeNotificationContactsClient{
			contacts:    sets.New("old-user"),
			failingAdds: sets.New("new-user"),
		}
		rosaScope := newNotificationContactsScope([]string{"new-user"})

		g.Expect(rosa.ReconcileNotificationContacts(ctx, rosaScope, fake, cluster)).NotTo(Succeed())
		g.Expect(sets.List(fake.contacts)).To(Equal([]string{"old-user"}))

		fake.failingAdds = sets.New[string]()
		g.Expect(rosa.ReconcileNotificationContacts(ctx, rosaScope, fake, cluster)).To(Succeed())
		g.Expect(sets.List(fake.contacts)).To(Equal([]string{"new-user"}))
	})
}

func newNotificationContactsScope(contacts []string) *scope.ROSAControlPlaneScope {
	return &scope.ROSAControlPlaneScope{
		ControlPlane: &rosacontrolplanev1.ROSAControlPlane{
			Spec: rosacontrolplanev1.RosaControlPlaneSpec{
				NotificationContacts: contacts,
			},
		},
	}
}

// fakeNotificationContactsClient implements the notification contact methods of
// rosa.OCMClient; calling any other method panics.
type fakeNotificationContactsClient struct {
	rosa.OCMClient
	contacts    sets.Set[string]
	failingAdds sets.Set[string]
}

func (f *fakeNotificationContactsClient) GetSubscriptionNotificationContacts(_ context.Context, _ string) ([]string, error) {
	if f.contacts.Len() == 0 {
		return nil, nil
	}
	return sets.List(f.contacts), nil
}

func (f *fakeNotificationContactsClient) UpdateSubscriptionNotificationContacts(_ context.Context, _ string, contacts []string) error {
	desired := sets.New[string]()
	for _, c := range contacts {
		if c = strings.TrimSpace(c); c != "" {
			desired.Insert(c)
		}
	}

	var addErrors []string
	for _, c := range sets.List(desired.Difference(f.contacts)) {
		if f.failingAdds.Has(c) {
			addErrors = append(addErrors, c)
			continue
		}
		f.contacts.Insert(c)
	}

	f.contacts = f.contacts.Intersection(desired)

	if len(addErrors) > 0 {
		return fmt.Errorf("failed to add notification contact(s): %s", strings.Join(addErrors, ", "))
	}
	return nil
}
