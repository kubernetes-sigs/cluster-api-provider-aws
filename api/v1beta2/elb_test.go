/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	healthz     = ELBHealthCheckTarget("HTTPS:6443/healthz")
	httpHealthz = ELBHealthCheckTarget("HTTP:6443/healthz")
	dcp         = ELBHealthCheckTarget("DCCP:6443")
	tcpPath     = ELBHealthCheckTarget("TCP:6443/healthz")
)

func TestELBValidate(t *testing.T) {
	tests := []struct {
		name     string
		self     *ELBHealthCheckTarget
		expected []*field.Error
	}{
		{
			name:     "no errors",
			self:     &healthz,
			expected: nil,
		},
		{
			name:     "no errors http with path",
			self:     &httpHealthz,
			expected: nil,
		},
		{
			name: "bad protocol",
			self: &dcp,
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "value cannot have characters other than alphabets, numbers, and ':' .",
					Field:    "spec.controlPlaneLoadBalancer.healthCheck",
					BadValue: "DCCP:6443",
				},
			},
		},
		{
			name: "wrong protocol with path",
			self: &tcpPath,
			expected: []*field.Error{
				{
					Type:     field.ErrorTypeInvalid,
					Detail:   "cannot specify paths with protocol other than HTTP and HTTPS",
					Field:    "spec.controlPlaneLoadBalancer.healthCheck",
					BadValue: "TCP:6443/healthz",
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.self.ValidateELBHealthCheck()
			if !cmp.Equal(out, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, out)
			}
		})
	}
}
