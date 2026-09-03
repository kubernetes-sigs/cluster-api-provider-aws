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

package webhooks

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestROSAMachinePoolValidateCreate(t *testing.T) {
	tests := []struct {
		name             string
		pool             *expinfrav1.ROSAMachinePool
		wantErrToContain *string
	}{
		{
			name: "nil spotMarketOptions is accepted",
			pool: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{},
			},
			wantErrToContain: nil,
		},
		{
			name: "empty spotMarketOptions is accepted",
			pool: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{},
				},
			},
			wantErrToContain: nil,
		},
		{
			name: "spotMarketOptions with maxPrice is accepted",
			pool: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{MaxPrice: ptr.To("0.05")},
				},
			},
			wantErrToContain: nil,
		},
		{
			name: "spotMarketOptions with capacityReservationID is rejected",
			pool: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions:     &expinfrav1.SpotMarketOptions{},
					CapacityReservationID: "cr-123",
				},
			},
			wantErrToContain: ptr.To[string]("spec.spotMarketOptions"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			warn, err := (&ROSAMachinePool{}).ValidateCreate(context.Background(), tt.pool)
			if tt.wantErrToContain != nil {
				g.Expect(err).ToNot(BeNil())
				if err != nil {
					g.Expect(err.Error()).To(ContainSubstring(*tt.wantErrToContain))
				}
			} else {
				g.Expect(err).To(Succeed())
			}
			// Nothing emits warnings yet
			g.Expect(warn).To(BeEmpty())
		})
	}
}

func TestROSAMachinePoolValidateUpdate(t *testing.T) {
	tests := []struct {
		name             string
		old              *expinfrav1.ROSAMachinePool
		new              *expinfrav1.ROSAMachinePool
		wantErrToContain *string
	}{
		{
			name: "unchanged nil spotMarketOptions is accepted",
			old: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{},
			},
			new: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{},
			},
			wantErrToContain: nil,
		},
		{
			name: "unchanged spotMarketOptions with maxPrice is accepted",
			old: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{MaxPrice: ptr.To("0.05")},
				},
			},
			new: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{MaxPrice: ptr.To("0.05")},
				},
			},
			wantErrToContain: nil,
		},
		{
			name: "unchanged empty spotMarketOptions is accepted",
			old: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{},
				},
			},
			new: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{},
				},
			},
			wantErrToContain: nil,
		},
		{
			name: "adding spotMarketOptions (nil -> set) is rejected as immutable",
			old: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{},
			},
			new: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{},
				},
			},
			wantErrToContain: ptr.To[string]("spec.spotMarketOptions"),
		},
		{
			name: "removing spotMarketOptions (set -> nil) is rejected as immutable",
			old: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{},
				},
			},
			new: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{},
			},
			wantErrToContain: ptr.To[string]("spec.spotMarketOptions"),
		},
		{
			name: "changing spotMarketOptions maxPrice is rejected as immutable",
			old: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{MaxPrice: ptr.To("0.05")},
				},
			},
			new: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{MaxPrice: ptr.To("0.10")},
				},
			},
			wantErrToContain: ptr.To[string]("spec.spotMarketOptions"),
		},
		{
			name: "typed-nil vs empty struct is treated as a change (immutable)",
			old: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: nil,
				},
			},
			new: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions: &expinfrav1.SpotMarketOptions{},
				},
			},
			wantErrToContain: ptr.To[string]("spec.spotMarketOptions"),
		},
		{
			name: "adding spotMarketOptions with capacityReservationID is rejected",
			old: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					CapacityReservationID: "cr-123",
				},
			},
			new: &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					SpotMarketOptions:     &expinfrav1.SpotMarketOptions{},
					CapacityReservationID: "cr-123",
				},
			},
			wantErrToContain: ptr.To[string]("spec.spotMarketOptions"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			warn, err := (&ROSAMachinePool{}).ValidateUpdate(context.Background(), tt.old.DeepCopy(), tt.new)
			if tt.wantErrToContain != nil {
				g.Expect(err).ToNot(BeNil())
				if err != nil {
					g.Expect(err.Error()).To(ContainSubstring(*tt.wantErrToContain))
				}
			} else {
				g.Expect(err).To(Succeed())
			}
			// Nothing emits warnings yet
			g.Expect(warn).To(BeEmpty())
		})
	}
}
