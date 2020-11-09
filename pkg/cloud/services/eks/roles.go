/*
Copyright 2020 The Kubernetes Authors.

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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1alpha3"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	eksiam "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/eks/iam"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

// NodegroupRolePolicies gives the policies required for a nodegroup role
func NodegroupRolePolicies() []string {
	return []string{
		"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
		"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy", //TODO: Can remove when CAPA supports provisioning of OIDC web identity federation with service account token volume projection
		"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
	}
}

func (s *Service) reconcileControlPlaneIAMRole() error {
	s.scope.V(2).Info("Reconciling EKS Control Plane IAM Role")

	if s.scope.ControlPlane.Spec.RoleName == nil {
		if !s.scope.EnableIAM() {
			s.scope.Info("no eks control plane role specified, using default eks control plane role")
			s.scope.ControlPlane.Spec.RoleName = &ekscontrolplanev1.DefaultEKSControlPlaneRole
		} else {
			s.scope.Info("no eks control plane role specified, using role based on cluster name")
			s.scope.ControlPlane.Spec.RoleName = aws.String(fmt.Sprintf("%s-iam-service-role", s.scope.Name()))
		}
	}
	s.scope.Info("using eks control plane role", "role-name", *s.scope.ControlPlane.Spec.RoleName)

	role, err := s.GetIAMRole(*s.scope.ControlPlane.Spec.RoleName)
	if err != nil {
		if !isNotFound(err) {
			return err
		}

		// If the disable IAM flag is used then the role must exist
		if !s.scope.EnableIAM() {
			return fmt.Errorf("getting role %s: %w", *s.scope.ControlPlane.Spec.RoleName, ErrClusterRoleNotFound)
		}

		role, err = s.CreateRole(*s.scope.ControlPlane.Spec.RoleName, s.scope.Name(), eksiam.ControlPlaneTrustRelationship(false), s.scope.AdditionalTags())
		if err != nil {
			record.Warnf(s.scope.ControlPlane, "FailedIAMRoleCreation", "Failed to create control plane IAM role %q: %v", *s.scope.ControlPlane.Spec.RoleName, err)

			return fmt.Errorf("creating role %s: %w", *s.scope.ControlPlane.Spec.RoleName, err)
		}
		record.Eventf(s.scope.ControlPlane, "SucessfulIAMRoleCreation", "Created control plane IAM role %q", *s.scope.ControlPlane.Spec.RoleName)
	}

	if s.IsUnmanaged(role, s.scope.Name()) {
		s.scope.V(2).Info("Skipping, EKS control plane role policy assignment as role is unamanged")
		return nil
	}

	//TODO: check tags and trust relationship to see if they need updating

	policies := []*string{
		aws.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
	}
	if s.scope.ControlPlane.Spec.RoleAdditionalPolicies != nil {
		if !s.scope.AllowAdditionalRoles() && len(*s.scope.ControlPlane.Spec.RoleAdditionalPolicies) > 0 {
			return ErrCannotUseAdditionalRoles
		}

		for _, policy := range *s.scope.ControlPlane.Spec.RoleAdditionalPolicies {
			additionalPolicy := policy
			policies = append(policies, &additionalPolicy)
		}
	}
	err = s.EnsurePoliciesAttached(role, policies)
	if err != nil {
		return errors.Wrapf(err, "error ensuring policies are attached: %v", policies)
	}

	return nil
}

func (s *Service) deleteControlPlaneIAMRole() error {
	if s.scope.ControlPlane.Spec.RoleName == nil {
		return nil
	}
	roleName := *s.scope.ControlPlane.Spec.RoleName
	if !s.scope.EnableIAM() {
		s.scope.V(2).Info("EKS IAM disabled, skipping deleting EKS Control Plane IAM Role")
		return nil
	}

	s.scope.V(2).Info("Deleting EKS Control Plane IAM Role")

	role, err := s.GetIAMRole(roleName)
	if err != nil {
		if isNotFound(err) {
			s.V(2).Info("EKS Control Plane IAM Role already deleted")
			return nil
		}

		return errors.Wrap(err, "getting eks control plane iam role")
	}

	if s.IsUnmanaged(role, s.scope.Name()) {
		s.V(2).Info("Skipping, EKS control plane iam role deletion as role is unamanged")
		return nil
	}

	err = s.DeleteRole(*s.scope.ControlPlane.Spec.RoleName)
	if err != nil {
		record.Eventf(s.scope.ControlPlane, "FailedIAMRoleDeletion", "Failed to delete control Plane IAM role %q: %v", *s.scope.ControlPlane.Spec.RoleName, err)
		return err
	}

	record.Eventf(s.scope.ControlPlane, "SucessfulIAMRoleDeletion", "Deleted Control Plane IAM role %q", *s.scope.ControlPlane.Spec.RoleName)
	return nil
}

func (s *NodegroupService) reconcileNodegroupIAMRole() error {
	s.scope.V(2).Info("Reconciling EKS Nodegroup IAM Role")

	if s.scope.RoleName() == "" {
		var roleName string
		if !s.scope.EnableIAM() {
			s.scope.Info("no EKS nodegroup role specified, using default EKS nodegroup role")
			roleName = infrav1exp.DefaultEKSNodegroupRole
		} else {
			s.scope.Info("no EKS nodegroup role specified, using role based on nodegroup name")
			roleName = fmt.Sprintf("%s-%s-nodegroup-iam-service-role", s.scope.KubernetesClusterName(), s.scope.NodegroupName())
		}
		s.scope.ManagedMachinePool.Spec.RoleName = roleName
	}

	role, err := s.GetIAMRole(s.scope.RoleName())
	if err != nil {
		if !isNotFound(err) {
			return err
		}

		// If the disable IAM flag is used then the role must exist
		if !s.scope.EnableIAM() {
			return ErrNodegroupRoleNotFound
		}

		role, err = s.CreateRole(s.scope.ManagedMachinePool.Spec.RoleName, s.scope.ClusterName(), eksiam.NodegroupTrustRelationship(), s.scope.AdditionalTags())
		if err != nil {
			record.Warnf(s.scope.ManagedMachinePool, "FailedIAMRoleCreation", "Failed to create nodegroup IAM role %q: %v", s.scope.RoleName(), err)
			return err
		}
		record.Eventf(s.scope.ManagedMachinePool, "SucessfulIAMRoleCreation", "Created nodegroup IAM role %q", s.scope.RoleName())
	}

	if s.IsUnmanaged(role, s.scope.ClusterName()) {
		s.scope.V(2).Info("Skipping, EKS nodegroup role policy assignment as role is unamanged")
		return nil
	}

	err = s.EnsureTagsAndPolicy(role, s.scope.ClusterName(), eksiam.NodegroupTrustRelationship(), s.scope.AdditionalTags())
	if err != nil {
		return errors.Wrapf(err, "error ensuring tags and policy document are set on node role")
	}

	policies := NodegroupRolePolicies()
	err = s.EnsurePoliciesAttached(role, aws.StringSlice(policies))
	if err != nil {
		return errors.Wrapf(err, "error ensuring policies are attached: %v", policies)
	}

	return nil
}

func (s *NodegroupService) deleteNodegroupIAMRole() (reterr error) {
	if err := s.scope.IAMReadyFalse(clusterv1.DeletingReason, ""); err != nil {
		return err
	}
	defer func() {
		if reterr != nil {
			record.Warnf(
				s.scope.ManagedMachinePool, "FailedDeleteIAMNodegroupRole", "Failed to delete EKS nodegroup role %s: %v", s.scope.ManagedMachinePool.Spec.RoleName, reterr,
			)
			if err := s.scope.IAMReadyFalse("DeletingFailed", reterr.Error()); err != nil {
				reterr = err
			}
		} else if err := s.scope.IAMReadyFalse(clusterv1.DeletedReason, ""); err != nil {
			reterr = err
		}
	}()
	roleName := s.scope.RoleName()
	if !s.scope.EnableIAM() {
		s.scope.V(2).Info("EKS IAM disabled, skipping deleting EKS Nodegroup IAM Role")
		return nil
	}

	s.scope.V(2).Info("Deleting EKS Nodegroup IAM Role")

	role, err := s.GetIAMRole(roleName)
	if err != nil {
		if isNotFound(err) {
			s.V(2).Info("EKS Nodegroup IAM Role already deleted")
			return nil
		}

		return errors.Wrap(err, "getting EKS nodegroup iam role")
	}

	if s.IsUnmanaged(role, s.scope.ClusterName()) {
		s.V(2).Info("Skipping, EKS Nodegroup iam role deletion as role is unamanged")
		return nil
	}

	err = s.DeleteRole(s.scope.RoleName())
	if err != nil {
		record.Eventf(s.scope.ManagedMachinePool, "FailedIAMRoleDeletion", "Failed to delete Nodegroup IAM role %q: %v", s.scope.ManagedMachinePool.Spec.RoleName, err)
		return err
	}

	record.Eventf(s.scope.ManagedMachinePool, "SucessfulIAMRoleDeletion", "Deleted Nodegroup IAM role %q", s.scope.ManagedMachinePool.Spec.RoleName)
	return nil
}

func isNotFound(err error) bool {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case iam.ErrCodeNoSuchEntityException:
			return true
		default:
			return false
		}
	}

	return false
}
