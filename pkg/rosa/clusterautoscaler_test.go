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
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/gomega"
	sdk "github.com/openshift-online/ocm-sdk-go"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
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

// TestOCMClientClusterAutoscalerErrorMapping makes sure the OCM errors of the cluster autoscaler
// requests are still reported as OCMAPIError. The controller decides whether to retry from that
// status, so an error that loses it turns a rejected request into an endless reconcile loop.
func TestOCMClientClusterAutoscalerErrorMapping(t *testing.T) {
	// A JWT with token_use=access, so the SDK uses it as an access token instead of trying to
	// refresh it against the test server.
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"token_use":"access","exp":4102444800}`))
	token := "eyJhbGciOiJub25lIn0." + payload + ".signature"

	newClient := func(t *testing.T, status int, body string) *ocmclient {
		t.Helper()
		g := NewWithT(t)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, err := w.Write([]byte(body))
			g.Expect(err).ToNot(HaveOccurred())
		}))
		t.Cleanup(server.Close)

		connection, err := sdk.NewConnectionBuilder().URL(server.URL).Tokens(token).Build()
		g.Expect(err).ToNot(HaveOccurred())
		t.Cleanup(func() { _ = connection.Close() })

		return &ocmclient{connection: connection}
	}

	t.Run("Get returns the HTTP status of the OCM error", func(t *testing.T) {
		g := NewWithT(t)
		client := newClient(t, http.StatusNotFound,
			`{"kind":"Error","code":"CLUSTERS-MGMT-404","reason":"Autoscaler for cluster ID 'cluster-123' is not found"}`)

		_, err := client.GetClusterAutoscaler(context.Background(), "cluster-123")
		g.Expect(err).To(HaveOccurred())
		g.Expect(OCMAPIErrorStatus(err)).To(Equal(http.StatusNotFound))
		g.Expect(err.Error()).To(ContainSubstring("is not found"))
	})

	t.Run("Update returns the HTTP status of the OCM error", func(t *testing.T) {
		g := NewWithT(t)
		client := newClient(t, http.StatusBadRequest,
			`{"kind":"Error","code":"CLUSTERS-MGMT-400","reason":"cluster size exceeds the limit"}`)

		_, err := client.UpdateClusterAutoscaler(context.Background(), "cluster-123", &cmv1.ClusterAutoscaler{})
		g.Expect(err).To(HaveOccurred())
		g.Expect(OCMAPIErrorStatus(err)).To(Equal(http.StatusBadRequest))
		g.Expect(err.Error()).To(ContainSubstring("cluster size exceeds the limit"))
	})

	t.Run("a client without a connection reports a clear error", func(t *testing.T) {
		g := NewWithT(t)
		client := &ocmclient{}

		_, err := client.GetClusterAutoscaler(context.Background(), "cluster-123")
		g.Expect(err).To(MatchError(errNoOCMConnection))
	})
}
