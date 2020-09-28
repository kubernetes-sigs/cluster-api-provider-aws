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

import "errors"

var (
	ErrInvalidBackendType = errors.New("invalid backend type")
	ErrClientRequired     = errors.New("k8s client required")
	ErrRoleARNRequired    = errors.New("rolearn is required")
	ErrUserARNRequired    = errors.New("userarn is required")
	ErrUserNameRequired   = errors.New("username is required")
	ErrGroupsRequired     = errors.New("groups are required")
	ErrIsNotARN           = errors.New("supplied value is not a ARN")
	ErrIsNotRoleARN       = errors.New("supplied ARN is not a role ARN")
	ErrIsNotUserARN       = errors.New("supplied ARN is not a user ARN")
)
