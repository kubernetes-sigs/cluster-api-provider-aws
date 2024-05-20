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

	. "github.com/onsi/gomega"
	"k8s.io/klog/v2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

func TestSubnetPlacement(t *testing.T) {
	testCases := []struct {
		name                string
		specSubnetIDs       []string
		specAZs             []string
		parentAZs           []string
		subnetPlacementType *expinfrav1.AZSubnetType
		controlPlaneSubnets infrav1.Subnets
		logger              *logger.Logger
		expectedSubnetIDs   []string
		expectError         bool
	}{
		{
			name:                "spec subnets expected",
			specSubnetIDs:       []string{"subnet-az1"},
			specAZs:             []string{"eu-west-1b"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: nil,
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az1"},
			expectError:       false,
		},
		{
			name:                "spec azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{"eu-west-1b"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: nil,
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az2", "subnet-az3"},
			expectError:       false,
		},
		{
			name:                "parent azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: nil,
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az4", "subnet-az5"},
			expectError:       false,
		},
		{
			name:                "spec private azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{"eu-west-1b"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypePrivate),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az6",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
					Tags: infrav1.Tags{
						infrav1.NameAWSSubnetAssociation: infrav1.SecondarySubnetTagValue,
					},
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az3"},
			expectError:       false,
		},
		{
			name:                "spec public azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{"eu-west-1b"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypePublic),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az2"},
			expectError:       false,
		},
		{
			name:                "spec all azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{"eu-west-1b"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypeAll),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az2", "subnet-az3"},
			expectError:       false,
		},
		{
			name:                "spec public no azs found",
			specSubnetIDs:       []string{},
			specAZs:             []string{"eu-west-1a"},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypePublic),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{},
			expectError:       true,
		},
		{
			name:                "parent private azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypePrivate),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az5"},
			expectError:       false,
		},
		{
			name:                "parent public azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypePublic),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az4"},
			expectError:       false,
		},
		{
			name:                "parent all azs expected",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{"eu-west-1c"},
			subnetPlacementType: expinfrav1.NewAZSubnetType(expinfrav1.AZSubnetTypeAll),
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az4", "subnet-az5"},
			expectError:       false,
		},
		{
			name:          "use control plane subnets",
			specSubnetIDs: []string{},
			specAZs:       []string{},
			parentAZs:     []string{},
			controlPlaneSubnets: infrav1.Subnets{
				infrav1.SubnetSpec{
					ID:               "subnet-az1",
					AvailabilityZone: "eu-west-1a",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az2",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az3",
					AvailabilityZone: "eu-west-1b",
					IsPublic:         false,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az4",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         true,
				},
				infrav1.SubnetSpec{
					ID:               "subnet-az5",
					AvailabilityZone: "eu-west-1c",
					IsPublic:         false,
				},
			},
			logger:            logger.NewLogger(klog.Background()),
			expectedSubnetIDs: []string{"subnet-az1", "subnet-az3", "subnet-az5"},
			expectError:       false,
		},
		{
			name:                "no placement",
			specSubnetIDs:       []string{},
			specAZs:             []string{},
			parentAZs:           []string{},
			controlPlaneSubnets: infrav1.Subnets{},
			logger:              logger.NewLogger(klog.Background()),
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
				SubnetPlacementType:     tc.subnetPlacementType,
			})

			if tc.expectError {
				g.Expect(err).ToNot(BeNil())
				return
			}

			g.Expect(err).To(BeNil())
			g.Expect(actualSubnetIDs).To(Equal(tc.expectedSubnetIDs))
		})
	}
}
