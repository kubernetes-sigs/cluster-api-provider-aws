/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha1

type (
	Effect            string
	ConditionOperator string
	PrincipalType     string
)

const (

	// Any is the AWS IAM policy grammar wildcard.
	Any = "*"

	// CurrentVersion is the latest version of the AWS IAM policy grammar.
	CurrentVersion = "2012-10-17"

	// EffectAllow is the Allow effect in an AWS IAM policy statement entry.
	EffectAllow Effect = "Allow"

	// EffectDeny is the Deny effect in an AWS IAM policy statement entry.
	EffectDeny Effect = "Deny"

	// PrincipalAWS is the principal type covering AWS ARNs.
	PrincipalAWS PrincipalType = "AWS"

	// PrincipalFederated is the principal type covering federated identities.
	PrincipalFederated PrincipalType = "Federated"

	// PrincipalService is the principal type covering AWS services.
	PrincipalService PrincipalType = "Service"

	// StringEquals is an AWS IAM policy condition operator.
	StringEquals ConditionOperator = "StringEquals"

	// StringNotEquals is an AWS IAM policy condition operator.
	StringNotEquals ConditionOperator = "StringNotEquals"

	// StringEqualsIgnoreCase is an AWS IAM policy condition operator.
	StringEqualsIgnoreCase ConditionOperator = "StringEqualsIgnoreCase"

	// StringLike is an AWS IAM policy condition operator.
	StringLike ConditionOperator = "StringLike"

	// StringNotLike is an AWS IAM policy condition operator.
	StringNotLike ConditionOperator = "StringNotLike"
)

// PolicyDocument represents an AWS IAM policy document, and can be
// converted into JSON using "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/converters"
type PolicyDocument struct {
	Version   string
	Statement Statements
	ID        string `json:"id,omitempty"`
}

// StatementEntry represents each "statement" block in an AWS IAM policy document.
type StatementEntry struct {
	Sid          string     `json:",omitempty"`
	Principal    Principals `json:",omitempty"`
	NotPrincipal Principals `json:",omitempty"`
	Effect       Effect     `json:"Effect"`
	Action       Actions    `json:"Action"`
	Resource     Resources  `json:",omitempty"`
	Condition    Conditions `json:"Condition,omitempty"`
}

// Statements is the list of StatementEntries
type Statements []StatementEntry

// Principals is the map of all principals a statement entry refers to
type Principals map[PrincipalType]PrincipalID

// Actions is the list of actions
type Actions []string

// Resources is the list of resources
type Resources []string

// PrincipalID represents the list of all principals, such as ARNs
type PrincipalID []string

// Conditions is the map of all conditions in the statement entry.
type Conditions map[ConditionOperator]interface{}

// DeepCopyInto copies the receiver, writing into out. in must be non-nil.
func (in Conditions) DeepCopyInto(out *Conditions) {
	{
		in := &in
		*out = make(Conditions, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy copies the receiver, creating a new Conditions.
func (in Conditions) DeepCopy() Conditions {
	if in == nil {
		return nil
	}
	out := new(Conditions)
	in.DeepCopyInto(out)
	return *out
}
