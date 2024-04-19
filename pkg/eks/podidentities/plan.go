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

// Package podidentities provides a plan to manage EKS podidentities associations.
package podidentities

import (
	"context"

	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/planner"
)

// NewPlan creates a new Plan to manage EKS pod identities.
func NewPlan(clusterName string, desiredAssociations, currentAssociations []EKSPodIdentityAssociation, client eksiface.EKSAPI) planner.Plan {
	return &plan{
		currentAssociations: currentAssociations,
		desiredAssociations: desiredAssociations,
		eksClient:           client,
		clusterName:         clusterName,
	}
}

// Plan is a plan that will manage EKS addons.
type plan struct {
	currentAssociations []EKSPodIdentityAssociation
	desiredAssociations []EKSPodIdentityAssociation
	eksClient           eksiface.EKSAPI
	clusterName         string
}

func (a *plan) getCurrentAssociation(association EKSPodIdentityAssociation) bool {
	for _, current := range a.currentAssociations {
		if current.ServiceAccountName == association.ServiceAccountName && current.ServiceAccountNamespace == association.ServiceAccountNamespace {
			return true
		}
	}
	return false
}

func (a *plan) getDesiredAssociation(association EKSPodIdentityAssociation) bool {
	for _, desired := range a.desiredAssociations {
		if desired.ServiceAccountName == association.ServiceAccountName && desired.ServiceAccountNamespace == association.ServiceAccountNamespace {
			return true
		}
	}
	return false
}

// Create will create the plan (i.e. list of procedures) for managing EKS addons.
func (a *plan) Create(_ context.Context) ([]planner.Procedure, error) {
	procedures := []planner.Procedure{}

	for _, d := range a.desiredAssociations {
		desired := d
		existsInCurrent := a.getCurrentAssociation(desired)
		existsInDesired := a.getDesiredAssociation(desired)

		// Create pod association if is doesnt already exist
		if existsInDesired && !existsInCurrent {
			procedures = append(procedures,
				&CreatePodIdentityAssociationProcedure{
					eksClient:      a.eksClient,
					clusterName:    a.clusterName,
					newAssociation: &desired,
				},
			)
		}
	}

	for _, current := range a.currentAssociations {
		existsInCurrent := a.getCurrentAssociation(current)
		existsInDesired := a.getDesiredAssociation(current)

		if !existsInDesired && existsInCurrent {
			// Delete pod association if it exists
			procedures = append(procedures,
				&DeletePodIdentityAssociationProcedure{
					eksClient:             a.eksClient,
					clusterName:           a.clusterName,
					existingAssociationID: current.AssociationID,
				},
			)
		}
	}

	return procedures, nil
}
