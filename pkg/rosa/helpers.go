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

package rosa

import (
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	ocmerrors "github.com/openshift-online/ocm-sdk-go/errors"
	errors "github.com/zgalor/weberr"
)

// IsNodePoolReady checkes whether the nodepool is provisoned and all replicas are available.
// If autosacling is enabled, NodePool must have replicas >= autosacling.MinReplica to be considered ready.
func IsNodePoolReady(nodePool *cmv1.NodePool) bool {
	if nodePool.Status().Message() != "" {
		return false
	}

	if nodePool.Replicas() != 0 {
		return nodePool.Replicas() == nodePool.Status().CurrentReplicas()
	}

	if nodePool.Autoscaling() != nil {
		return nodePool.Status().CurrentReplicas() >= nodePool.Autoscaling().MinReplica()
	}

	return false
}

func handleErr(res *ocmerrors.Error, err error) error {
	msg := res.Reason()
	if msg == "" {
		msg = err.Error()
	}
	// Hack to always display the correct terms and conditions message
	if res.Code() == "CLUSTERS-MGMT-451" {
		msg = "You must accept the Terms and Conditions in order to continue.\n" +
			"Go to https://www.redhat.com/wapps/tnc/ackrequired?site=ocm&event=register\n" +
			"Once you accept the terms, you will need to retry the action that was blocked."
	}
	errType := errors.ErrorType(res.Status()) //#nosec G115
	return errType.Set(errors.Errorf("%s", msg))
}
