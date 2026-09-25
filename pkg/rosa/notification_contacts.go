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
	"reflect"
	"sort"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// ReconcileNotificationContacts syncs spec.notificationContacts to the cluster
// subscription's notification contacts sub-resource.
// When the field is unset (nil), contacts are not managed.
// When set (including empty), the subscription is reconciled to exactly that set.
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

	desired := append([]string(nil), rosaScope.ControlPlane.Spec.NotificationContacts...)
	sort.Strings(desired)

	current, err := ocmClient.GetSubscriptionNotificationContacts(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get notification contacts for cluster '%s': %w", cluster.ID(), err)
	}
	if current == nil {
		current = []string{}
	}
	sortedCurrent := append([]string(nil), current...)
	sort.Strings(sortedCurrent)

	if reflect.DeepEqual(desired, sortedCurrent) {
		return nil
	}

	if err := ocmClient.UpdateSubscriptionNotificationContacts(ctx, subscriptionID, desired); err != nil {
		return fmt.Errorf("failed to update notification contacts for cluster '%s': %w", cluster.ID(), err)
	}

	return nil
}
