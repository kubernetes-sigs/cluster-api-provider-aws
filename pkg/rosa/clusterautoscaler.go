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
	"errors"
	"fmt"
	"net/http"

	ocmcommonclient "github.com/openshift-online/ocm-common/pkg/ocm/client"
	sdk "github.com/openshift-online/ocm-sdk-go"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	ocmerrors "github.com/openshift-online/ocm-sdk-go/errors"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// ClusterAutoscalerClient reads and updates the cluster autoscaler of a ROSA HCP cluster through the OCM API.
// It uses the client shared through ocm-common (also used by other OCM clients such as Terraform) on top of
// the OCM SDK, and does not depend on the rosa CLI library.
type ClusterAutoscalerClient interface {
	// GetClusterAutoscaler returns the cluster autoscaler of the given cluster.
	GetClusterAutoscaler(ctx context.Context, clusterID string) (*cmv1.ClusterAutoscaler, error)
	// UpdateClusterAutoscaler sends a PATCH with the given (possibly partial) cluster autoscaler.
	UpdateClusterAutoscaler(ctx context.Context, clusterID string, autoscaler *cmv1.ClusterAutoscaler) (*cmv1.ClusterAutoscaler, error)
	// Close releases the underlying OCM connection.
	Close() error
}

type clusterAutoscalerClient struct {
	connection *sdk.Connection
	autoscaler ocmcommonclient.ClusterAutoscalerClient
}

// NewClusterAutoscalerClient creates a ClusterAutoscalerClient using the OCM credentials of the given control plane.
func NewClusterAutoscalerClient(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (ClusterAutoscalerClient, error) {
	connection, err := newOCMRawConnection(ctx, rosaScope)
	if err != nil {
		return nil, err
	}
	return &clusterAutoscalerClient{
		connection: connection,
		autoscaler: ocmcommonclient.NewClusterAutoscalerClient(connection.ClustersMgmt().V1().Clusters()),
	}, nil
}

func (c *clusterAutoscalerClient) GetClusterAutoscaler(ctx context.Context, clusterID string) (*cmv1.ClusterAutoscaler, error) {
	autoscaler, err := c.autoscaler.Get(ctx, clusterID)
	if err != nil {
		return nil, newOCMAPIError(err)
	}
	return autoscaler, nil
}

func (c *clusterAutoscalerClient) UpdateClusterAutoscaler(ctx context.Context, clusterID string, autoscaler *cmv1.ClusterAutoscaler) (*cmv1.ClusterAutoscaler, error) {
	updated, err := c.autoscaler.Update(ctx, clusterID, autoscaler)
	if err != nil {
		return nil, newOCMAPIError(err)
	}
	return updated, nil
}

func (c *clusterAutoscalerClient) Close() error {
	return c.connection.Close()
}

// OCMAPIError is an error returned by the OCM API. Unlike handleErr, it keeps the HTTP status
// and the operation ID so callers can classify the error and support can trace the request.
type OCMAPIError struct {
	// Status is the HTTP status code returned by OCM.
	Status int
	// Code is the OCM error code, for example "CLUSTERS-MGMT-400".
	Code string
	// Reason is the human readable error message returned by OCM.
	Reason string
	// OperationID identifies the request in the OCM logs.
	OperationID string
}

func (e *OCMAPIError) Error() string {
	msg := e.Reason
	if msg == "" {
		msg = http.StatusText(e.Status)
	}
	if e.OperationID != "" {
		return fmt.Sprintf("%s (operation ID: %s)", msg, e.OperationID)
	}
	return msg
}

// newOCMAPIError converts an error returned by the OCM SDK into an OCMAPIError. For HTTP error
// responses the SDK returns an *ocmerrors.Error; any other error, such as a network failure, is
// returned unchanged.
func newOCMAPIError(err error) error {
	var apiErr *ocmerrors.Error
	if !errors.As(err, &apiErr) || apiErr.Status() == 0 {
		return err
	}
	return &OCMAPIError{
		Status:      apiErr.Status(),
		Code:        apiErr.Code(),
		Reason:      apiErr.Reason(),
		OperationID: apiErr.OperationID(),
	}
}

// OCMAPIErrorStatus returns the HTTP status of an OCMAPIError, or 0 if err is not an OCMAPIError.
func OCMAPIErrorStatus(err error) int {
	var apiErr *OCMAPIError
	if errors.As(err, &apiErr) {
		return apiErr.Status
	}
	return 0
}
