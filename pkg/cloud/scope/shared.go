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
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
)

var (
	// ErrAZSubnetsNotFound is an error when a availability zone is specified but there are
	// no matching subnets for that availability zone (a.k.a. fault domain).
	ErrAZSubnetsNotFound = errors.New("no subnets found for supplied availability zone")
	// ErrLoggerRequired is an error if a logger isn't specified.
	ErrLoggerRequired = errors.New("logger is required")
	// ErrNotPlaced is an error if there is no placement determined.
	ErrNotPlaced = errors.New("placement not determined")
)

type placementInput struct {
	SpecSubnetIDs           []string
	SpecAvailabilityZones   []string
	ParentAvailabilityZones []string
	ControlplaneSubnets     infrav1.Subnets
}

type subnetsPlacementStratgey interface {
	Place(input *placementInput) ([]string, error)
}

func newDefaultSubnetPlacementStrategy(logger *logr.Logger) (subnetsPlacementStratgey, error) {
	if logger == nil {
		return nil, ErrLoggerRequired
	}

	return &defaultSubnetPlacementStrategy{
		logger: *logger,
	}, nil
}

// defaultSubnetPlacementStrategy is the default strategy for subnet placement.
type defaultSubnetPlacementStrategy struct {
	logger logr.Logger
}

// Place works out the subnet placement based on the following precedence:
// 1. Explicit definition of subnet IDs in the spec
// 2. If the spec has Availability Zones then get the subnets for these AZs
// 3. If the parent resource has Availability Zones then get the subnets for these AZs
// 4. All the private subnets from the control plane are used
// In Cluster API Availability Zone can also be referred to by the name `Failure Domain`.
func (p *defaultSubnetPlacementStrategy) Place(input *placementInput) ([]string, error) {
	if len(input.SpecSubnetIDs) > 0 {
		p.logger.V(2).Info("using subnets from the spec")
		return input.SpecSubnetIDs, nil
	}

	if len(input.SpecAvailabilityZones) > 0 {
		p.logger.V(2).Info("determining subnets to use from the spec availability zones")
		subnetIDs, err := p.getSubnetsForAZs(input.SpecAvailabilityZones, input.ControlplaneSubnets)
		if err != nil {
			return nil, fmt.Errorf("getting subnets for spec azs: %w", err)
		}

		return subnetIDs, nil
	}

	if len(input.ParentAvailabilityZones) > 0 {
		p.logger.V(2).Info("determining subnets to use from the parents availability zones")
		subnetIDs, err := p.getSubnetsForAZs(input.ParentAvailabilityZones, input.ControlplaneSubnets)
		if err != nil {
			return nil, fmt.Errorf("getting subnets for parent azs: %w", err)
		}

		return subnetIDs, nil
	}

	controlPlaneSubnetIDs := input.ControlplaneSubnets.FilterPrivate().IDs()
	if len(controlPlaneSubnetIDs) > 0 {
		p.logger.V(2).Info("using all the private subnets from the control plane")
		return controlPlaneSubnetIDs, nil
	}

	return nil, ErrNotPlaced
}

func (p *defaultSubnetPlacementStrategy) getSubnetsForAZs(azs []string, controlPlaneSubnets infrav1.Subnets) ([]string, error) {
	subnetIDs := []string{}

	for _, zone := range azs {
		subnets := controlPlaneSubnets.FilterByZone(zone)
		if len(subnets) == 0 {
			return nil, fmt.Errorf("getting subnets for availability zone %s: %w", zone, ErrAZSubnetsNotFound)
		}
		subnetIDs = append(subnetIDs, subnets.IDs()...)
	}

	return subnetIDs, nil
}
