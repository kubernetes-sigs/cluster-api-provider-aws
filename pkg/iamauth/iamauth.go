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

package iamauth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// EC2NodeUserName is the username required for EC2 nodes
	EC2NodeUserName = "system:node:{{EC2PrivateDNSName}}"
)

var (
	// NodeGroups is the groups that are required for a node
	NodeGroups = []string{"system:bootstrappers", "system:nodes"}
)

// AuthenticatorBackend is the interface that represents an aws-iam-authenticator backend
type AuthenticatorBackend interface {
	// MapRole is used to map a role ARN to a user and set of groups
	MapRole(mapping RoleMapping) error
	// MapUser is used to map a user ARN to a user and set of groups
	MapUser(mapping UserMapping) error
}

// BackendType is a type that represents the different aws-iam-authenticator backends
type BackendType string

var (
	// BackendTypeConfigMap is the Kubernetes config map backend
	BackendTypeConfigMap = BackendType("config-map")
	// BackendTypeCRD is the CRD based backend
	BackendTypeCRD = BackendType("crd")
)

// New will create a new authenticate backend for a given type
func New(backendType BackendType, client crclient.Client) (AuthenticatorBackend, error) {
	if client == nil {
		return nil, ErrClientRequired
	}

	switch backendType {
	case BackendTypeConfigMap:
		return &configMapBackend{client: client}, nil
	case BackendTypeCRD:
		return &crdBackend{client: client}, nil
	default:
		return nil, ErrInvalidBackendType
	}
}

type IAMAuthenticatorConfig struct {
	RoleMappings []RoleMapping `json:"mapRoles,omitempty"`
	UserMappings []UserMapping `json:"mapUsers,omitempty"`
}

type KubernetesMapping struct {
	UserName string   `json:"username,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

type RoleMapping struct {
	KubernetesMapping
	RoleARN string `json:"roleARN,omitempty"`
}

type UserMapping struct {
	KubernetesMapping
	UserARN string `json:"userARN,omitempty"`
}

func (r *RoleMapping) Validate() error {
	errs := []error{}

	if strings.TrimSpace(r.RoleARN) == "" {
		errs = append(errs, ErrRoleARNRequired)
	}
	if strings.TrimSpace(r.UserName) == "" {
		errs = append(errs, ErrUserNameRequired)
	}
	if len(r.Groups) == 0 {
		errs = append(errs, ErrGroupsRequired)
	}

	if !arn.IsARN(r.RoleARN) {
		errs = append(errs, ErrIsNotARN)
	} else {
		parsedARN, err := arn.Parse(r.RoleARN)
		if err != nil {
			errs = append(errs, err)
		} else if !strings.Contains(parsedARN.Resource, "role/") {
			errs = append(errs, ErrIsNotRoleARN)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	err := "Validation errors:\n"
	for i, e := range errs {
		err += fmt.Sprintf("\t%d: %s\n", i, e.Error())
	}

	return errors.New(err) //nolint: err113
}

func (u *UserMapping) Validate() error {
	errs := []error{}

	if strings.TrimSpace(u.UserARN) == "" {
		errs = append(errs, ErrUserARNRequired)
	}
	if strings.TrimSpace(u.UserName) == "" {
		errs = append(errs, ErrUserNameRequired)
	}
	if len(u.Groups) == 0 {
		errs = append(errs, ErrGroupsRequired)
	}

	if !arn.IsARN(u.UserARN) {
		errs = append(errs, ErrIsNotARN)
	} else {
		parsedARN, err := arn.Parse(u.UserARN)
		if err != nil {
			errs = append(errs, err)
		} else if !strings.Contains(parsedARN.Resource, "user/") {
			errs = append(errs, ErrIsNotUserARN)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	err := "Validation errors:\n"
	for i, e := range errs {
		err += fmt.Sprintf("\t%d: %s\n", i, e.Error())
	}

	return errors.New(err) //nolint: err113
}
