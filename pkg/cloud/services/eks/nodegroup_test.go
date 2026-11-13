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

package eks

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/gomega"
)

func TestNodeRepairConfigNeedsUpdate(t *testing.T) {
	disabled := &ekstypes.NodeRepairConfig{Enabled: aws.Bool(false)}
	enabled := &ekstypes.NodeRepairConfig{Enabled: aws.Bool(true)}

	tests := []struct {
		name    string
		current *ekstypes.NodeRepairConfig
		desired *ekstypes.NodeRepairConfig
		want    bool
	}{
		{name: "region without node auto repair, repair disabled", current: nil, desired: disabled, want: false},
		{name: "region without node auto repair, repair enabled", current: nil, desired: enabled, want: true},
		{name: "repair config without enabled, repair disabled", current: &ekstypes.NodeRepairConfig{}, desired: disabled, want: false},
		{name: "desired without enabled, repair disabled", current: nil, desired: &ekstypes.NodeRepairConfig{}, want: false},
		{name: "desired without enabled, repair enabled", current: enabled, desired: &ekstypes.NodeRepairConfig{}, want: true},
		{name: "disabled and disabled", current: disabled, desired: disabled, want: false},
		{name: "enabled and enabled", current: enabled, desired: enabled, want: false},
		{name: "disabled to enabled", current: disabled, desired: enabled, want: true},
		{name: "enabled to disabled", current: enabled, desired: disabled, want: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(nodeRepairConfigNeedsUpdate(tc.current, tc.desired)).To(Equal(tc.want))
		})
	}
}
