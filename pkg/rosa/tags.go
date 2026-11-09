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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/blang/semver"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	ocmerrors "github.com/openshift-online/ocm-sdk-go/errors"
)

const (
	clustersPath = "/api/clusters_mgmt/v1/clusters/"

	clusterAWSField  = "aws"
	nodePoolAWSField = "aws_node_pool"
)

// DayTwoTagsMinVersion is the first control plane version that reconciles AWS tag edits in place.
// The prerelease suffix makes 4.20 prereleases eligible too.
var DayTwoTagsMinVersion = semver.MustParse("4.20.0-0")

var errNoRawConnection = errors.New("OCM client has no raw connection; AWS tag edits are unavailable")

// SupportsDayTwoTags reports whether a control plane version supports day-2 AWS tag edits.
// It accepts plain versions and the prefixed raw IDs returned by OCM.
func SupportsDayTwoTags(controlPlaneVersion string) bool {
	version := strings.TrimPrefix(controlPlaneVersion, "openshift-v")
	version = strings.TrimPrefix(version, "openshift-")
	parsed, err := semver.Parse(version)
	if err != nil {
		return false
	}
	return parsed.GTE(DayTwoTagsMinVersion)
}

// TagChanges returns the OCM PATCH tag changes needed to replace current with desired.
// A nil result means the sets are equal; nil values in the result delete keys.
func TagChanges(current, desired map[string]string) map[string]*string {
	var changes map[string]*string
	for key, value := range desired {
		if currentValue, ok := current[key]; ok && currentValue == value {
			continue
		}
		if changes == nil {
			changes = make(map[string]*string)
		}
		changes[key] = &value
	}
	for key := range current {
		if _, ok := desired[key]; ok {
			continue
		}
		if changes == nil {
			changes = make(map[string]*string)
		}
		changes[key] = nil
	}
	return changes
}

// NodePoolTagChanges returns the OCM PATCH tag changes needed to make the node pool's own tags
// match desired. effective is the node pool's customer tag view, which includes tags inherited
// from the cluster, so tags whose values match the inherited cluster tags aren't treated as
// node pool overrides. A nil result means no update is needed.
func NodePoolTagChanges(clusterTags, effective, desired map[string]string) map[string]*string {
	desiredOverrides := poolOverrides(clusterTags, desired)
	changes := TagChanges(poolOverrides(clusterTags, effective), desiredOverrides)
	if changes == nil {
		return nil
	}
	// The effective view can't distinguish an inherited tag from a stored node pool override
	// with the same value. OCM ignores null for an inherited-only key but removes a stored
	// override, so clear those keys while an update is being sent anyway.
	for key, inherited := range clusterTags {
		if value, ok := effective[key]; !ok || value != inherited {
			continue
		}
		if _, keep := desiredOverrides[key]; keep {
			continue
		}
		changes[key] = nil
	}
	return changes
}

// TagsDrifted reports whether any desired tag is missing from, or differs in, actual.
// It can't detect removals, but works with tag views that include OCM system tags.
func TagsDrifted(actual, desired map[string]string) bool {
	for key, value := range desired {
		if actualValue, ok := actual[key]; !ok || actualValue != value {
			return true
		}
	}
	return false
}

// poolOverrides removes tags matching inherited cluster values from a node pool tag set.
func poolOverrides(clusterTags, poolTags map[string]string) map[string]string {
	overrides := make(map[string]string, len(poolTags))
	for key, value := range poolTags {
		if inherited, ok := clusterTags[key]; !ok || inherited != value {
			overrides[key] = value
		}
	}
	return overrides
}

// GetClusterUserTags returns the customer cluster AWS tags, omitting OCM system tags.
func (c *ocmclient) GetClusterUserTags(clusterID string) (map[string]string, error) {
	return c.getUserTags(clustersPath+url.PathEscape(clusterID), clusterAWSField)
}

// GetNodePoolUserTags returns the customer node pool AWS tags, including tags inherited from the cluster.
func (c *ocmclient) GetNodePoolUserTags(clusterID string, nodePoolID string) (map[string]string, error) {
	return c.getUserTags(nodePoolPath(clusterID, nodePoolID), nodePoolAWSField)
}

// UpdateClusterTags patches the cluster AWS tags. Nil values delete keys.
func (c *ocmclient) UpdateClusterTags(clusterID string, changes map[string]*string) error {
	cluster, err := cmv1.NewCluster().Build()
	if err != nil {
		return fmt.Errorf("failed to build cluster update: %w", err)
	}
	var typed bytes.Buffer
	if err := cmv1.MarshalCluster(cluster, &typed); err != nil {
		return fmt.Errorf("failed to encode cluster update: %w", err)
	}
	return c.patchTags(clustersPath+url.PathEscape(clusterID), typed.Bytes(), clusterAWSField, changes)
}

// UpdateNodePoolTags patches the node pool AWS tags. Nil values delete keys.
func (c *ocmclient) UpdateNodePoolTags(clusterID string, nodePoolID string, changes map[string]*string) error {
	nodePool, err := cmv1.NewNodePool().ID(nodePoolID).Build()
	if err != nil {
		return fmt.Errorf("failed to build node pool update: %w", err)
	}
	var typed bytes.Buffer
	if err := cmv1.MarshalNodePool(nodePool, &typed); err != nil {
		return fmt.Errorf("failed to encode node pool update: %w", err)
	}
	return c.patchTags(nodePoolPath(clusterID, nodePoolID), typed.Bytes(), nodePoolAWSField, changes)
}

func nodePoolPath(clusterID, nodePoolID string) string {
	return clustersPath + url.PathEscape(clusterID) + "/node_pools/" + url.PathEscape(nodePoolID)
}

// getUserTags requests customer-only AWS tags and omits null deletion markers from the response.
func (c *ocmclient) getUserTags(path, field string) (map[string]string, error) {
	if c.connection == nil {
		return nil, errNoRawConnection
	}
	response, err := c.connection.Get().Path(path).Parameter("fetchUserTagsOnly", true).Send()
	if err != nil {
		return nil, fmt.Errorf("failed to read AWS tags: %w", err)
	}
	if err := responseError(response.Status(), response.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to read AWS tags: %w", err)
	}
	return decodeUserTags(response.Bytes(), field)
}

// decodeUserTags extracts the tags map under field from an OCM resource body.
func decodeUserTags(body []byte, field string) (map[string]string, error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(body, &object); err != nil {
		return nil, fmt.Errorf("failed to decode AWS tags response: %w", err)
	}
	var awsFields struct {
		Tags map[string]*string `json:"tags"`
	}
	if len(object[field]) > 0 {
		if err := json.Unmarshal(object[field], &awsFields); err != nil {
			return nil, fmt.Errorf("failed to decode AWS tags: %w", err)
		}
	}
	tags := make(map[string]string, len(awsFields.Tags))
	for key, value := range awsFields.Tags {
		if value != nil {
			tags[key] = *value
		}
	}
	return tags, nil
}

// patchTags inserts nullable tag changes into a typed JSON update and sends it as a PATCH.
// The typed SDK can't encode null map values, which OCM uses to delete tags.
func (c *ocmclient) patchTags(path string, typed []byte, field string, changes map[string]*string) error {
	if c.connection == nil {
		return errNoRawConnection
	}
	body, err := tagPatchBody(typed, field, changes)
	if err != nil {
		return err
	}
	response, err := c.connection.Patch().Path(path).Header("Content-Type", "application/json").Bytes(body).Send()
	if err != nil {
		return fmt.Errorf("failed to update AWS tags: %w", err)
	}
	if err := responseError(response.Status(), response.Bytes()); err != nil {
		return fmt.Errorf("failed to update AWS tags: %w", err)
	}
	return nil
}

// tagPatchBody sets field.tags in the typed JSON object to changes.
func tagPatchBody(typed []byte, field string, changes map[string]*string) ([]byte, error) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(typed, &object); err != nil {
		return nil, fmt.Errorf("failed to decode typed update: %w", err)
	}
	awsFields := make(map[string]json.RawMessage)
	if len(object[field]) > 0 {
		if err := json.Unmarshal(object[field], &awsFields); err != nil {
			return nil, fmt.Errorf("failed to decode typed AWS update: %w", err)
		}
	}
	tagsJSON, err := json.Marshal(changes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode AWS tag changes: %w", err)
	}
	awsFields["tags"] = tagsJSON
	if object[field], err = json.Marshal(awsFields); err != nil {
		return nil, fmt.Errorf("failed to encode AWS update: %w", err)
	}
	body, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("failed to encode update: %w", err)
	}
	return body, nil
}

// responseError preserves the OCM error reason when a raw request fails.
func responseError(status int, body []byte) error {
	if status >= http.StatusOK && status < http.StatusMultipleChoices {
		return nil
	}
	err := fmt.Errorf("OCM returned HTTP %d", status)
	if serviceErr, unmarshalErr := ocmerrors.UnmarshalErrorStatus(body, status); unmarshalErr == nil && serviceErr != nil {
		return handleErr(serviceErr, err)
	}
	return err
}
