/*
Copyright 2021 The Kubernetes Authors.

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

package scope

import (
	"testing"

	"github.com/go-logr/logr"
	. "github.com/onsi/gomega"

	"k8s.io/klog/klogr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

func TestSubnetPlacement(t *testing.T) {
	testCases := []struct {
		name                string
		specSubnetIDs       []string
		specAZs             []string
		parentAZs           []string
		controlPlaneSubnets infrav1.Subnets
		logger              logr.Logger
		expectedSubnetIDs   []string
		expectError         bool
	}{
		{
			name:          "spec subnets expected",
			specSubnetIDs: []string{"az1"},
			specAZs:       []string{"eu-west-1b"},
			parentAZs:     []string{"eu-west-1c"},
			controlPlaneSubnets: infrav1.Subnets{
				&infrav1.SubnetSpec{
					ID:               "az1",
					AvailabilityZone: "eu-west-1a",
				},
				&infrav1.SubnetSpec{
					ID:               "az2",
					AvailabilityZone: "eu-west-1b",
				},
				&infrav1.SubnetSpec{
					ID:               "az3",
					AvailabilityZone: "eu-west-1c",
				},
			},
			logger:            klogr.New(),
			expectedSubnetIDs: []string{"az1"},
			expectError:       false,
		},
		{
			name:          "spec azs expected",
			specSubnetIDs: []string{},
			specAZs:       []string{"eu-west-1b"},
			parentAZs:     []string{"eu-west-1c"},
			controlPlaneSubnets: infrav1.Subnets{
				&infrav1.SubnetSpec{
					ID:               "az1",
					AvailabilityZone: "eu-west-1a",
				},
				&infrav1.SubnetSpec{
					ID:               "az2",
					AvailabilityZone: "eu-west-1b",
				},
				&infrav1.SubnetSpec{
					ID:               "az3",
					AvailabilityZone: "eu-west-1c",
				},
			},
			logger:            klogr.New(),
			expectedSubnetIDs: []string{"az2"},
			expectError:       false,
		},
		{
			name:          "parent azs expected",
			specSubnetIDs: []string{},
			specAZs:       []string{},
			parentAZs:     []string{"eu-west-1c"},
			controlPlaneSubnets: infrav1.Subnets{
				&infrav1.SubnetSpec{
					ID:               "az1",
					AvailabilityZone: "eu-west-1a",
				},
				&infrav1.SubnetSpec{
					ID:               "az2",
					AvailabilityZone: "eu-west-1b",
				},
				&infrav1.SubnetSpec{
					ID:               "az3",
					AvailabilityZone: "eu-west-1c",
				},
			},
			logger:            klogr.New(),
			expectedSubnetIDs: []string{"az3"},
			expectError:       false,
		},
		{
			name:          "use control plane subnets",
			specSubnetIDs: []string{},
			specAZs:       []string{},
			parentAZs:     []string{},
			controlPlaneSubnets: infrav1.Subnets{
				&infrav1.SubnetSpec{
					ID:               "az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				&infrav1.SubnetSpec{
					ID:               "az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				&infrav1.SubnetSpec{
					ID:               "az3",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
			},
			logger:            klogr.New(),
			expectedSubnetIDs: []string{"az1", "az2"},
			expectError:       false,
		},
		{
			name:                "no placement",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{},
			controlPlaneSubnets: infrav1.Subnets{},
			logger:              klogr.New(),
			expectedSubnetIDs:   []string{},
			expectError:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			strategy, err := newDefaultSubnetPlacementStrategy(tc.logger)
			g.Expect(err).NotTo(HaveOccurred())

			actualSubnetIDs, err := strategy.Place(&placementInput{
				SpecSubnetIDs:           tc.specSubnetIDs,
				SpecAvailabilityZones:   tc.specAZs,
				ParentAvailabilityZones: tc.parentAZs,
				ControlplaneSubnets:     tc.controlPlaneSubnets,
			})

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
				return
			} else {
				g.Expect(err).To(BeNil())
			}

			g.Expect(actualSubnetIDs).To(Equal(tc.expectedSubnetIDs))

		})
	}
}
