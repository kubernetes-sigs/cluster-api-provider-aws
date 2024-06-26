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

package podidentities

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
)

// errPodIdentityAssociationNotFound defines an error for when an eks pod identity is not found.
var errPodIdentityAssociationNotFound = errors.New("eks pod identity association not found")

// DeletePodIdentityAssociationProcedure is a procedure that will delete an EKS eks pod identity.
type DeletePodIdentityAssociationProcedure struct {
	eksClient             eksiface.EKSAPI
	clusterName           string
	existingAssociationID string
}

// Do implements the logic for the procedure.
func (p *DeletePodIdentityAssociationProcedure) Do(_ context.Context) error {
	input := &eks.DeletePodIdentityAssociationInput{
		AssociationId: aws.String(p.existingAssociationID),
		ClusterName:   aws.String(p.clusterName),
	}

	if _, err := p.eksClient.DeletePodIdentityAssociation(input); err != nil {
		return fmt.Errorf("deleting eks pod identity %s: %w", p.existingAssociationID, err)
	}

	return nil
}

// Name is the name of the procedure.
func (p *DeletePodIdentityAssociationProcedure) Name() string {
	return "eks_pod_identity_delete"
}

// CreatePodIdentityAssociationProcedure is a procedure that will create an EKS eks pod identity for a cluster.
type CreatePodIdentityAssociationProcedure struct {
	eksClient      eksiface.EKSAPI
	clusterName    string
	newAssociation *EKSPodIdentityAssociation
}

// Do implements the logic for the procedure.
func (p *CreatePodIdentityAssociationProcedure) Do(_ context.Context) error {
	if p.newAssociation == nil {
		return fmt.Errorf("getting desired eks pod identity for cluster %s: %w", p.clusterName, errPodIdentityAssociationNotFound)
	}

	input := &eks.CreatePodIdentityAssociationInput{
		ClusterName:    aws.String(p.clusterName),
		Namespace:      &p.newAssociation.ServiceAccountNamespace,
		RoleArn:        &p.newAssociation.RoleARN,
		ServiceAccount: &p.newAssociation.ServiceAccountName,
	}

	_, err := p.eksClient.CreatePodIdentityAssociation(input)
	if err != nil {
		return fmt.Errorf("creating desired eks pod identity for cluster %s: %w", p.clusterName, err)
	}

	return nil
}

// Name is the name of the procedure.
func (p *CreatePodIdentityAssociationProcedure) Name() string {
	return "eks_pod_identity_create"
}
