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
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
)

// TrustRelationshipPolicyDocument represesnts an IAM policy docyment
type TrustRelationshipPolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

// ToJSONString converts the document to a JSON string
func (d *TrustRelationshipPolicyDocument) ToJSONString() (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// StatementEntry represents a statement within an IAM policy document
type StatementEntry struct {
	Effect    string
	Action    []string
	Principal map[string][]string
}

func (s *Service) reconcileControlPlaneIAMRole() error {
	s.scope.V(2).Info("Reconciling EKS Control Plane IAM Role")

	if s.scope.ControlPlane.Spec.RoleName == nil {
		s.scope.ControlPlane.Spec.RoleName = aws.String(fmt.Sprintf("%s-iam-service-role", s.scope.Name()))
	}

	role, err := s.getIAMRole(*s.scope.ControlPlane.Spec.RoleName)
	if err != nil {
		if !isNotFound(err) {
			return err
		}

		role, err = s.createRole(*s.scope.ControlPlane.Spec.RoleName)
		if err != nil {
			record.Warnf(s.scope.ControlPlane, "FailedIAMRoleCreation", "Failed to create control plane IAM role %q: %v", *s.scope.ControlPlane.Spec.RoleName, err)
			return err
		}
		record.Eventf(s.scope.ControlPlane, "SucessfulIAMRoleCreation", "Created control plane IAM role %q", *s.scope.ControlPlane.Spec.RoleName)
	}

	if s.isUnmanaged(role) {
		s.scope.V(2).Info("Skipping, EKS control plane role policy assignment as role is unamanged")
		return nil
	}

	//TODO: check tags and trust relationship to see if they need updating

	policies := []*string{
		aws.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
	}
	if s.scope.ControlPlane.Spec.RoleAdditionalPolicies != nil {
		for _, policy := range *s.scope.ControlPlane.Spec.RoleAdditionalPolicies {
			additionalPolicy := policy
			policies = append(policies, &additionalPolicy)
		}
	}
	err = s.ensurePoliciesAttached(role, policies)
	if err != nil {
		return errors.Wrapf(err, "error ensuring policies are attached: %v", policies)
	}

	return nil
}

func (s *Service) getIAMRole(name string) (*iam.Role, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}

	out, err := s.IAMClient.GetRole(input)
	if err != nil {
		return nil, err
	}

	return out.Role, nil
}

func (s *Service) getIAMPolicy(policyArn string) (*iam.Policy, error) {
	input := &iam.GetPolicyInput{
		PolicyArn: &policyArn,
	}

	out, err := s.IAMClient.GetPolicy(input)
	if err != nil {
		return nil, err
	}

	return out.Policy, nil
}

func (s *Service) getIAMRolePolicies(roleName string) ([]*string, error) {
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	}

	out, err := s.IAMClient.ListAttachedRolePolicies(input)
	if err != nil {
		return nil, errors.Wrapf(err, "error listing role polices for %s", roleName)
	}

	policies := []*string{}
	for _, policy := range out.AttachedPolicies {
		policies = append(policies, policy.PolicyArn)
	}

	return policies, nil
}

func (s *Service) detachIAMRolePolicy(roleName string, policyARN string) error {
	input := &iam.DetachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(policyARN),
	}

	_, err := s.IAMClient.DetachRolePolicy(input)
	if err != nil {
		return errors.Wrapf(err, "error detaching policy %s from role %s", policyARN, roleName)
	}

	return nil
}

func (s *Service) attachIAMRolePolicy(roleName string, policyARN string) error {
	input := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(policyARN),
	}

	_, err := s.IAMClient.AttachRolePolicy(input)
	if err != nil {
		return errors.Wrapf(err, "error attaching policy %s to role %s", policyARN, roleName)
	}

	return nil
}

func (s *Service) ensurePoliciesAttached(role *iam.Role, policies []*string) error {
	s.scope.V(2).Info("Ensuring Polices are attached to EKS Control Plane IAM Role")
	existingPolices, err := s.getIAMRolePolicies(*role.RoleName)
	if err != nil {
		return err
	}

	// Remove polices that aren't in the list
	for _, existingPolicy := range existingPolices {
		found := findStringInSlice(policies, *existingPolicy)
		if !found {
			err = s.detachIAMRolePolicy(*role.RoleName, *existingPolicy)
			if err != nil {
				return err
			}
			s.scope.V(2).Info("Detached policy from role", "role", role.RoleName, "policy", existingPolicy)
		}
	}

	// Add any policies that aren't currently attached
	for _, policy := range policies {
		found := findStringInSlice(existingPolices, *policy)
		if !found {
			// Make sure policy exists before attaching
			_, err := s.getIAMPolicy(*policy)
			if err != nil {
				return errors.Wrapf(err, "error getting policy %s", *policy)
			}

			err = s.attachIAMRolePolicy(*role.RoleName, *policy)
			if err != nil {
				return err
			}
			s.scope.V(2).Info("Attached policy to role", "role", role.RoleName, "policy", *policy)
		}
	}

	return nil
}

func (s *Service) createRole(name string) (*iam.Role, error) {
	//TODO: tags also needs a separate sync
	additionalTags := s.scope.AdditionalTags()
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())] = string(infrav1.ResourceLifecycleOwned)
	tags := []*iam.Tag{}
	for k, v := range additionalTags {
		tags = append(tags, &iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	trustRelationship := s.controlPlaneTrustRelationship(false)
	trustRelationShipJSON, err := trustRelationship.ToJSONString()
	if err != nil {
		return nil, errors.Wrap(err, "error converting trust relationship to json")
	}

	input := &iam.CreateRoleInput{
		RoleName:                 aws.String(name),
		Tags:                     tags,
		AssumeRolePolicyDocument: aws.String(trustRelationShipJSON),
	}

	out, err := s.IAMClient.CreateRole(input)
	if err != nil {
		return nil, err
	}

	return out.Role, nil
}

func (s *Service) detachAllPoliciesForRole(name string) error {
	s.scope.V(3).Info("Detaching all policies for role", "role", name)
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: &name,
	}
	policies, err := s.IAMClient.ListAttachedRolePolicies(input)
	if err != nil {
		return errors.Wrapf(err, "error fetching policies for role %s", name)
	}
	for _, p := range policies.AttachedPolicies {
		s.scope.V(2).Info("Detaching policy", "policy", *p)
		if err := s.detachIAMRolePolicy(name, *p.PolicyArn); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) deleteRole(name string) error {
	if err := s.detachAllPoliciesForRole(name); err != nil {
		return errors.Wrapf(err, "error detaching policies for role %s", name)
	}

	input := &iam.DeleteRoleInput{
		RoleName: aws.String(name),
	}

	_, err := s.IAMClient.DeleteRole(input)
	if err != nil {
		return errors.Wrapf(err, "error deleting role %s", name)
	}

	return nil
}

func (s *Service) deleteControlPlaneIAMRole() error {
	s.scope.V(2).Info("Deleting EKS Control Plane IAM Role")

	role, err := s.getIAMRole(*s.scope.ControlPlane.Spec.RoleName)
	if err != nil {
		if isNotFound(err) {
			s.scope.V(2).Info("EKS Control Plane IAM Role already deleted")
			return nil
		}

		return errors.Wrap(err, "getting eks control plane iam role")
	}

	if s.isUnmanaged(role) {
		s.scope.V(2).Info("Skipping, EKS control plane iam role deletion as role is unamanged")
		return nil
	}

	err = s.deleteRole(*s.scope.ControlPlane.Spec.RoleName)
	if err != nil {
		record.Eventf(s.scope.ControlPlane, "FailedIAMRoleDeletion", "Failed to delete control Plane IAM role %q: %v", *s.scope.ControlPlane.Spec.RoleName, err)
		return err
	}

	record.Eventf(s.scope.ControlPlane, "SucessfulIAMRoleDeletion", "Deleted Control Plane IAM role %q", *s.scope.ControlPlane.Spec.RoleName)
	return nil
}

func (s *Service) isUnmanaged(role *iam.Role) bool {
	keyToFind := infrav1.ClusterAWSCloudProviderTagKey(s.scope.Name())
	for _, tag := range role.Tags {
		if *tag.Key == keyToFind && *tag.Value == string(infrav1.ResourceLifecycleOwned) {
			return false
		}
	}

	return true
}

func (s *Service) controlPlaneTrustRelationship(enableFargate bool) *TrustRelationshipPolicyDocument {
	principal := make(map[string][]string)
	principal["Service"] = []string{"eks.amazonaws.com"}
	if enableFargate {
		principal["Service"] = append(principal["Service"], "eks-fargate-pods.amazonaws.com")
	}

	policy := &TrustRelationshipPolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			{
				Effect: "Allow",
				Action: []string{
					"sts:AssumeRole",
				},
				Principal: principal,
			},
		},
	}

	return policy
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

func findStringInSlice(slice []*string, toFind string) bool {
	for _, item := range slice {
		if *item == toFind {
			return true
		}
	}

	return false
}
