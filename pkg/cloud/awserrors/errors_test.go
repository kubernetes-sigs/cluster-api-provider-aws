/*
Copyright 2024 The Kubernetes Authors.

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

package awserrors

import (
	"errors"
	"testing"

	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

func TestIsInsufficientResourcesError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "InsufficientFreeAddressesInSubnet error",
			err:      &smithy.GenericAPIError{Code: InsufficientFreeAddressesInSubnet, Message: "not enough free addresses in subnet"},
			expected: true,
		},
		{
			name:     "InsufficientInstanceCapacity error",
			err:      &smithy.GenericAPIError{Code: InsufficientInstanceCapacity, Message: "insufficient instance capacity"},
			expected: true,
		},
		{
			name:     "InstanceLimitExceeded error",
			err:      &smithy.GenericAPIError{Code: InstanceLimitExceeded, Message: "instance limit exceeded"},
			expected: true,
		},
		{
			name:     "InsufficientVolumeCapacity error",
			err:      &smithy.GenericAPIError{Code: InsufficientVolumeCapacity, Message: "insufficient volume capacity"},
			expected: true,
		},
		{
			name:     "VcpuLimitExceeded error",
			err:      &smithy.GenericAPIError{Code: VcpuLimitExceeded, Message: "vCPU limit exceeded"},
			expected: true,
		},
		{
			name:     "non-resource-exhaustion AWS error",
			err:      &smithy.GenericAPIError{Code: InvalidInstanceID, Message: "instance not found"},
			expected: false,
		},
		{
			name:     "AuthFailure error",
			err:      &smithy.GenericAPIError{Code: AuthFailure, Message: "auth failure"},
			expected: false,
		},
		{
			name:     "plain error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "wrapped InsufficientFreeAddressesInSubnet error",
			err:      pkgerrors.Wrap(&smithy.GenericAPIError{Code: InsufficientFreeAddressesInSubnet, Message: "not enough free addresses"}, "failed to run instance"),
			expected: true,
		},
		{
			name:     "double-wrapped InsufficientInstanceCapacity error",
			err:      pkgerrors.Wrap(pkgerrors.Wrap(&smithy.GenericAPIError{Code: InsufficientInstanceCapacity, Message: "insufficient capacity"}, "failed to run instance"), "unable to create instance"),
			expected: true,
		},
		{
			name:     "wrapped non-matching error",
			err:      pkgerrors.Wrap(&smithy.GenericAPIError{Code: AuthFailure, Message: "auth failure"}, "failed to run instance"),
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsInsufficientResourcesError(tc.err)
			if result != tc.expected {
				t.Errorf("IsInsufficientResourcesError(%v) = %v, want %v", tc.err, result, tc.expected)
			}
		})
	}
}
