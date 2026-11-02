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
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	ocmerrors "github.com/openshift-online/ocm-sdk-go/errors"
)

// errNoOCMConnection is returned when an operation needs the OCM SDK connection but the client was
// built without one, for example from an existing rosa-cli client.
var errNoOCMConnection = errors.New("OCM client has no OCM connection")

// GetClusterAutoscaler returns the cluster autoscaler of the given cluster.
func (c *ocmclient) GetClusterAutoscaler(ctx context.Context, clusterID string) (*cmv1.ClusterAutoscaler, error) {
	if c.connection == nil {
		return nil, errNoOCMConnection
	}
	autoscaler, err := ocmcommonclient.NewClusterAutoscalerClient(c.connection.ClustersMgmt().V1().Clusters()).Get(ctx, clusterID)
	if err != nil {
		return nil, newOCMAPIError(err)
	}
	return autoscaler, nil
}

// UpdateClusterAutoscaler sends a PATCH with the given (possibly partial) cluster autoscaler.
func (c *ocmclient) UpdateClusterAutoscaler(ctx context.Context, clusterID string, autoscaler *cmv1.ClusterAutoscaler) (*cmv1.ClusterAutoscaler, error) {
	if c.connection == nil {
		return nil, errNoOCMConnection
	}
	updated, err := ocmcommonclient.NewClusterAutoscalerClient(c.connection.ClustersMgmt().V1().Clusters()).Update(ctx, clusterID, autoscaler)
	if err != nil {
		return nil, newOCMAPIError(err)
	}
	return updated, nil
}

// OCMAPIError is an error returned by the OCM API, including its HTTP status and operation ID.
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
