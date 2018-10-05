// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iam

import (
	"encoding/json"
	"fmt"
)

const (

	// Any is an IAM policy grammar wildcard
	Any = "*"

	// CurrentVersion is the latest version of the IAM policy grammar
	CurrentVersion = "2012-10-17"

	// EffectAllow is the Allow effect in an IAM policy statement entry
	EffectAllow = "Allow"

	// EffectDeny is the Deny effect in an IAM policy statement entry
	EffectDeny = "Deny"

	// IAMSuffix is a standard prefix for resources for Cluster API Provider AWS
	IAMSuffix = "cluster-api-provider-aws.sigs.k8s.io"

	// PrincipalAWS is the principal covering AWS arns.
	PrincipalAWS = "AWS"

	// PrincipalFederated is the principal covering federated identities.
	PrincipalFederated = "Federated"

	// PrincipalService is the principal covering AWS services.
	PrincipalService = "Service"
)

// PolicyDocument represents an AWS IAM policy document
type PolicyDocument struct {
	Version   string
	Statement Statements
	ID        string `json:"id,omitempty"`
}

// StatementEntry represents each "statement" block in an IAM policy document
type StatementEntry struct {
	Sid          string     `json:",omitempty"`
	Principal    Principals `json:",omitempty"`
	NotPrincipal Principals `json:",omitempty"`
	Effect       string     `json:"Effect"`
	Action       Actions    `json:"Action"`
	Resource     Resources  `json:",omitempty"`
	Condition    Conditions `json:",omitempty"`
}

// Statements is the list of StatementEntries
type Statements []StatementEntry

// Principals is the map of all principals a statement entry refers to
type Principals map[string]PrincipalID

// Actions is the list of actions
type Actions []string

// Resources is the list of resources
type Resources []string

// PrincipalID represents the list of all principals, such as ARNs
type PrincipalID []string

// Conditions is the map of all conditions in the statement entry.
type Conditions map[string]ConditionValues

// ConditionValues are a list of condition values in a condition statement
type ConditionValues []string

// JSON is the JSON output of the policy document
func (p *PolicyDocument) JSON() (string, error) {
	b, err := json.Marshal(&p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// NewManagedName creates an IAM acceptable name prefixed with this Cluster API
// implementation's prefix.
func NewManagedName(prefix string) string {
	return fmt.Sprintf("%s.%s", prefix, IAMSuffix)
}
