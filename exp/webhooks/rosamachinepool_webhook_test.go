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

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

func TestROSAMachinePoolEc2MetadataHTTPTokensImmutability(t *testing.T) {
	w := &ROSAMachinePool{}

	t.Run("changing ec2MetadataHttpTokens is rejected", func(t *testing.T) {
		g := NewWithT(t)
		oldPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}
		newPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensOptional,
			},
		}

		_, err := w.ValidateUpdate(context.Background(), oldPool, newPool)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("ec2MetadataHttpTokens"))
	})

	t.Run("adding ec2MetadataHttpTokens after creation is rejected", func(t *testing.T) {
		g := NewWithT(t)
		oldPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName: "test-pool",
				InstanceType: "m5.large",
			},
		}
		newPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}

		_, err := w.ValidateUpdate(context.Background(), oldPool, newPool)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("ec2MetadataHttpTokens"))
	})

	t.Run("same ec2MetadataHttpTokens value is allowed", func(t *testing.T) {
		g := NewWithT(t)
		oldPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}
		newPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				Ec2MetadataHTTPTokens: rosacontrolplanev1.Ec2MetadataHTTPTokensRequired,
			},
		}

		_, err := w.ValidateUpdate(context.Background(), oldPool, newPool)
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("both empty ec2MetadataHttpTokens is allowed", func(t *testing.T) {
		g := NewWithT(t)
		oldPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName: "test-pool",
				InstanceType: "m5.large",
			},
		}
		newPool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName: "test-pool",
				InstanceType: "m5.large",
			},
		}

		_, err := w.ValidateUpdate(context.Background(), oldPool, newPool)
		g.Expect(err).ToNot(HaveOccurred())
	})
}

func TestROSAMachinePoolValidateSpotMarketOptions(t *testing.T) {
	w := &ROSAMachinePool{}

	tests := []struct {
		name                  string
		spotMarketOptions     *expinfrav1.RosaSpotMarketOptions
		capacityReservationID string
		expectError           bool
	}{
		{
			name:              "no spotMarketOptions and no capacityReservationID is allowed",
			spotMarketOptions: nil,
			expectError:       false,
		},
		{
			name:              "spotMarketOptions without maxPrice is allowed",
			spotMarketOptions: &expinfrav1.RosaSpotMarketOptions{},
			expectError:       false,
		},
		{
			name:              "spotMarketOptions with maxPrice is allowed",
			spotMarketOptions: &expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("0.5")},
			expectError:       false,
		},
		{
			name:                  "capacityReservationID without spotMarketOptions is allowed",
			spotMarketOptions:     nil,
			capacityReservationID: "cr-123456789012345678",
			expectError:           false,
		},
		{
			name:                  "spotMarketOptions with capacityReservationID is rejected",
			spotMarketOptions:     &expinfrav1.RosaSpotMarketOptions{},
			capacityReservationID: "cr-123456789012345678",
			expectError:           true,
		},
		{
			name:                  "spotMarketOptions with maxPrice and capacityReservationID is rejected",
			spotMarketOptions:     &expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("0.5")},
			capacityReservationID: "cr-123456789012345678",
			expectError:           true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			pool := &expinfrav1.ROSAMachinePool{
				Spec: expinfrav1.RosaMachinePoolSpec{
					NodePoolName:          "test-pool",
					InstanceType:          "m5.large",
					SpotMarketOptions:     tc.spotMarketOptions,
					CapacityReservationID: tc.capacityReservationID,
				},
			}

			allErrs := w.validateSpotMarketOptions(pool)
			if !tc.expectError {
				g.Expect(allErrs).To(BeEmpty())
				return
			}

			g.Expect(allErrs).To(HaveLen(1))
			g.Expect(allErrs[0].Field).To(Equal("spec.spotMarketOptions"))
			g.Expect(allErrs[0].Detail).To(ContainSubstring("setting both spotMarketOptions and capacityReservationID is not allowed"))
		})
	}
}

func TestROSAMachinePoolSpotMarketOptionsCreate(t *testing.T) {
	w := &ROSAMachinePool{}

	t.Run("spotMarketOptions alone is allowed", func(t *testing.T) {
		g := NewWithT(t)
		pool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:      "test-pool",
				InstanceType:      "m5.large",
				SpotMarketOptions: &expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("0.5")},
			},
		}

		_, err := w.ValidateCreate(context.Background(), pool)
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("spotMarketOptions with capacityReservationID is rejected", func(t *testing.T) {
		g := NewWithT(t)
		pool := &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				SpotMarketOptions:     &expinfrav1.RosaSpotMarketOptions{},
				CapacityReservationID: "cr-123456789012345678",
			},
		}

		_, err := w.ValidateCreate(context.Background(), pool)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("setting both spotMarketOptions and capacityReservationID is not allowed"))
	})
}

func TestROSAMachinePoolSpotMarketOptionsUpdate(t *testing.T) {
	w := &ROSAMachinePool{}

	newPoolWith := func(spot *expinfrav1.RosaSpotMarketOptions, capacityReservationID string) *expinfrav1.ROSAMachinePool {
		return &expinfrav1.ROSAMachinePool{
			Spec: expinfrav1.RosaMachinePoolSpec{
				NodePoolName:          "test-pool",
				InstanceType:          "m5.large",
				SpotMarketOptions:     spot,
				CapacityReservationID: capacityReservationID,
			},
		}
	}

	t.Run("adding spotMarketOptions after creation is not allowed", func(t *testing.T) {
		g := NewWithT(t)

		_, err := w.ValidateUpdate(context.Background(),
			newPoolWith(nil, ""),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("0.5")}, ""))
		g.Expect(err).To(HaveOccurred())
	})

	t.Run("changing maxPrice is allowed", func(t *testing.T) {
		g := NewWithT(t)

		_, err := w.ValidateUpdate(context.Background(),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("0.5")}, ""),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("0.7")}, ""))
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("deleting spotMarketOptions is not allowed", func(t *testing.T) {
		g := NewWithT(t)

		_, err := w.ValidateUpdate(context.Background(),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("0.5")}, ""),
			newPoolWith(nil, ""))
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("spotMarketOptions cannot be added/removed after the MachinePool is created"))
	})

	// The three day-2 transitions ROSA supports, verified against the live OCM API.
	t.Run("day-2: on-demand pool can be given a maxPrice", func(t *testing.T) {
		g := NewWithT(t)

		_, err := w.ValidateUpdate(context.Background(),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To(expinfrav1.RosaSpotMarketOptionsOnDemandPrice)}, ""),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("9.0")}, ""))
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("day-2: an existing maxPrice can be raised or lowered", func(t *testing.T) {
		g := NewWithT(t)

		_, err := w.ValidateUpdate(context.Background(),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("10.0")}, ""),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("9.0")}, ""))
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("day-2: maxPrice cannot be set back to on-demand", func(t *testing.T) {
		g := NewWithT(t)

		// OCM rejects an update that clears max_price, so the node pool can never reach
		// this spec. Without the guard the controller retries the same doomed update on
		// every reconcile. The sentinel makes this look like an ordinary value change,
		// which is exactly why it has to be checked explicitly.
		_, err := w.ValidateUpdate(context.Background(),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("9.0")}, ""),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To(expinfrav1.RosaSpotMarketOptionsOnDemandPrice)}, ""))
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring(`maxPrice cannot be changed back to "on-demand" after a price has been set`))
	})

	t.Run("day-2: dropping the maxPrice field is rejected the same way", func(t *testing.T) {
		g := NewWithT(t)

		// Equivalent to the sentinel for objects written before the default existed.
		_, err := w.ValidateUpdate(context.Background(),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{MaxPrice: ptr.To("9.0")}, ""),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{}, ""))
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring(`maxPrice cannot be changed back to "on-demand" after a price has been set`))
	})

	t.Run("adding capacityReservationID to a spot pool is rejected", func(t *testing.T) {
		g := NewWithT(t)

		_, err := w.ValidateUpdate(context.Background(),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{}, ""),
			newPoolWith(&expinfrav1.RosaSpotMarketOptions{}, "cr-123456789012345678"))
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("setting both spotMarketOptions and capacityReservationID is not allowed"))
	})
}
