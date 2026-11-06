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
	"errors"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	ocmerrors "github.com/openshift-online/ocm-sdk-go/errors"
)

func TestNewOCMAPIError(t *testing.T) {
	t.Run("returns the original error when there is no OCM response", func(t *testing.T) {
		g := NewWithT(t)
		netErr := errors.New("connection refused")
		g.Expect(newOCMAPIError(netErr)).To(BeIdenticalTo(netErr))
		g.Expect(OCMAPIErrorStatus(netErr)).To(Equal(0))
	})

	t.Run("keeps status, reason and operation ID", func(t *testing.T) {
		g := NewWithT(t)
		apiErr, err := ocmerrors.NewError().
			Status(400).
			Code("CLUSTERS-MGMT-400").
			Reason("Invalid max_node_provision_time '5m': value should be between 15m and 60m").
			OperationID("op-123").
			Build()
		g.Expect(err).ToNot(HaveOccurred())

		got := newOCMAPIError(apiErr)
		g.Expect(got.Error()).To(Equal("Invalid max_node_provision_time '5m': value should be between 15m and 60m (operation ID: op-123)"))
		g.Expect(OCMAPIErrorStatus(got)).To(Equal(400))
		g.Expect(OCMAPIErrorStatus(fmt.Errorf("failed to update cluster autoscaler: %w", got))).To(Equal(400))
	})

	t.Run("falls back to the HTTP status text without a reason", func(t *testing.T) {
		g := NewWithT(t)
		apiErr, err := ocmerrors.NewError().Status(503).Build()
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(newOCMAPIError(apiErr).Error()).To(Equal("Service Unavailable"))
	})
}
