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

package rosa

import (
	"context"
	"fmt"
	"strings"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// ReconcileNotificationContacts syncs spec.notificationContacts to the cluster
// subscription's notification contacts sub-resource.
// When the field is unset (nil), contacts are not managed.
// When set (including empty), the subscription is reconciled to exactly that set.
//
// Contacts are reconciled in two phases so that a failed addition never leaves the
// subscription without its existing recipients: missing contacts are added first
// while keeping all current contacts, and stale contacts are only removed once every
// addition has succeeded.
func ReconcileNotificationContacts(
	ctx context.Context,
	rosaScope *scope.ROSAControlPlaneScope,
	ocmClient OCMClient,
	cluster *cmv1.Cluster,
) error {
	if rosaScope.ControlPlane.Spec.NotificationContacts == nil {
		return nil
	}

	if cluster == nil {
		return fmt.Errorf("cluster is nil")
	}

	sub := cluster.Subscription()
	if sub == nil || sub.ID() == "" {
		return fmt.Errorf("cluster subscription ID is not available for cluster '%s'", cluster.ID())
	}
	subscriptionID := sub.ID()

	desired := normalizeNotificationContacts(rosaScope.ControlPlane.Spec.NotificationContacts)

	currentList, err := ocmClient.GetSubscriptionNotificationContacts(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get notification contacts for cluster '%s': %w", cluster.ID(), err)
	}
	current := normalizeNotificationContacts(currentList)

	toAdd := desired.Difference(current)
	toRemove := current.Difference(desired)

	if toAdd.Len() > 0 {
		// Only additions: every current contact is part of the requested set, so nothing is deleted.
		if err := ocmClient.UpdateSubscriptionNotificationContacts(ctx, subscriptionID, sets.List(current.Union(desired))); err != nil {
			return fmt.Errorf("failed to add notification contacts %v for cluster '%s', existing contacts were preserved: %w",
				sets.List(toAdd), cluster.ID(), err)
		}
	}

	if toRemove.Len() > 0 {
		// Only deletions: all desired contacts are now present.
		if err := ocmClient.UpdateSubscriptionNotificationContacts(ctx, subscriptionID, sets.List(desired)); err != nil {
			return fmt.Errorf("failed to remove notification contacts %v for cluster '%s': %w",
				sets.List(toRemove), cluster.ID(), err)
		}
	}

	return nil
}

// normalizeNotificationContacts trims whitespace and drops blank and duplicate
// entries, matching how OCM stores notification contact usernames.
func normalizeNotificationContacts(contacts []string) sets.Set[string] {
	normalized := sets.New[string]()
	for _, contact := range contacts {
		if contact = strings.TrimSpace(contact); contact != "" {
			normalized.Insert(contact)
		}
	}
	return normalized
}
