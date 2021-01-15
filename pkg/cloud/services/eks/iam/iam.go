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

package iam

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	apiiam "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/iam/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/converters"
)

type IAMService struct {
	logr.Logger
	IAMClient iamiface.IAMAPI
}

func (s *IAMService) GetIAMRole(name string) (*iam.Role, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}

	out, err := s.IAMClient.GetRole(input)
	if err != nil {
		return nil, err
	}

	return out.Role, nil
}

func (s *IAMService) getIAMPolicy(policyArn string) (*iam.Policy, error) {
	input := &iam.GetPolicyInput{
		PolicyArn: &policyArn,
	}

	out, err := s.IAMClient.GetPolicy(input)
	if err != nil {
		return nil, err
	}

	return out.Policy, nil
}

func (s *IAMService) getIAMRolePolicies(roleName string) ([]*string, error) {
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

func (s *IAMService) detachIAMRolePolicy(roleName string, policyARN string) error {
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

func (s *IAMService) attachIAMRolePolicy(roleName string, policyARN string) error {
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

func (s *IAMService) EnsurePoliciesAttached(role *iam.Role, policies []*string) error {
	s.V(2).Info("Ensuring Polices are attached to role")
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
			s.V(2).Info("Detached policy from role", "role", role.RoleName, "policy", existingPolicy)
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
			s.V(2).Info("Attached policy to role", "role", role.RoleName, "policy", *policy)
		}
	}

	return nil
}

func RoleTags(key string, additionalTags infrav1.Tags) []*iam.Tag {
	additionalTags[infrav1.ClusterAWSCloudProviderTagKey(key)] = string(infrav1.ResourceLifecycleOwned)
	tags := []*iam.Tag{}
	for k, v := range additionalTags {
		tags = append(tags, &iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return tags
}

func (s *IAMService) CreateRole(
	roleName string,
	key string,
	trustRelationship *apiiam.PolicyDocument,
	additionalTags infrav1.Tags,
) (*iam.Role, error) {
	tags := RoleTags(key, additionalTags)

	trustRelationshipJSON, err := converters.IAMPolicyDocumentToJSON(*trustRelationship)
	if err != nil {
		return nil, errors.Wrap(err, "error converting trust relationship to json")
	}

	input := &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		Tags:                     tags,
		AssumeRolePolicyDocument: aws.String(trustRelationshipJSON),
	}

	out, err := s.IAMClient.CreateRole(input)
	if err != nil {
		return nil, err
	}

	return out.Role, nil
}

func (s *IAMService) EnsureTagsAndPolicy(
	role *iam.Role,
	key string,
	trustRelationship *apiiam.PolicyDocument,
	additionalTags infrav1.Tags,
) error {
	s.V(2).Info("Ensuring tags and AssumeRolePolicyDocument are set on role")
	trustRelationshipJSON, err := converters.IAMPolicyDocumentToJSON(*trustRelationship)
	if err != nil {
		return errors.Wrap(err, "error converting trust relationship to json")
	}

	if trustRelationshipJSON != *role.AssumeRolePolicyDocument {
		policyInput := &iam.UpdateAssumeRolePolicyInput{
			RoleName:       role.RoleName,
			PolicyDocument: aws.String(trustRelationshipJSON),
		}
		_, err := s.IAMClient.UpdateAssumeRolePolicy(policyInput)
		if err != nil {
			return err
		}

	}

	tagInput := &iam.TagRoleInput{
		RoleName: role.RoleName,
	}
	untagInput := &iam.UntagRoleInput{
		RoleName: role.RoleName,
	}
	currentTags := make(map[string]string)
	for _, tag := range role.Tags {
		currentTags[*tag.Key] = *tag.Value
		if *tag.Key == infrav1.ClusterAWSCloudProviderTagKey(key) {
			continue
		}
		if _, ok := additionalTags[*tag.Key]; !ok {
			untagInput.TagKeys = append(untagInput.TagKeys, tag.Key)
		}
	}
	for key, value := range additionalTags {
		if currentV, ok := currentTags[key]; !ok || value != currentV {
			tagInput.Tags = append(tagInput.Tags, &iam.Tag{
				Key:   aws.String(key),
				Value: aws.String(value),
			})
		}
	}

	if len(tagInput.Tags) > 0 {
		_, err = s.IAMClient.TagRole(tagInput)
		if err != nil {
			return err
		}
	}

	if len(untagInput.TagKeys) > 0 {
		_, err = s.IAMClient.UntagRole(untagInput)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *IAMService) detachAllPoliciesForRole(name string) error {
	s.V(3).Info("Detaching all policies for role", "role", name)
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: &name,
	}
	policies, err := s.IAMClient.ListAttachedRolePolicies(input)
	if err != nil {
		return errors.Wrapf(err, "error fetching policies for role %s", name)
	}
	for _, p := range policies.AttachedPolicies {
		s.V(2).Info("Detaching policy", "policy", *p)
		if err := s.detachIAMRolePolicy(name, *p.PolicyArn); err != nil {
			return err
		}
	}
	return nil
}

func (s *IAMService) DeleteRole(name string) error {
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

func (s *IAMService) IsUnmanaged(role *iam.Role, key string) bool {
	keyToFind := infrav1.ClusterAWSCloudProviderTagKey(key)
	for _, tag := range role.Tags {
		if *tag.Key == keyToFind && *tag.Value == string(infrav1.ResourceLifecycleOwned) {
			return false
		}
	}

	return true
}

func ControlPlaneTrustRelationship(enableFargate bool) *apiiam.PolicyDocument {
	principal := make(apiiam.Principals)
	principal["Service"] = []string{"eks.amazonaws.com"}
	if enableFargate {
		principal["Service"] = append(principal["Service"], "eks-fargate-pods.amazonaws.com")
	}

	policy := &apiiam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []apiiam.StatementEntry{
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

func NodegroupTrustRelationship() *apiiam.PolicyDocument {
	principal := make(apiiam.Principals)
	principal["Service"] = []string{"ec2.amazonaws.com"}

	policy := &apiiam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []apiiam.StatementEntry{
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

func findStringInSlice(slice []*string, toFind string) bool {
	for _, item := range slice {
		if *item == toFind {
			return true
		}
	}

	return false
}

const stsAWSAudience = "sts.amazonaws.com"

func (s *IAMService) CreateOIDCProvider(cluster *eks.Cluster) (string, error) {
	issuerURL, err := url.Parse(*cluster.Identity.Oidc.Issuer)
	if err != nil {
		return "", err
	}
	if issuerURL.Scheme != "https" {
		return "", errors.Errorf("invalid scheme for issuer URL %s", issuerURL.String())
	}

	thumbprint, err := fetchRootCAThumbprint(issuerURL.String())
	if err != nil {
		return "", err
	}
	input := iam.CreateOpenIDConnectProviderInput{
		ClientIDList:   aws.StringSlice([]string{stsAWSAudience}),
		ThumbprintList: aws.StringSlice([]string{thumbprint}),
		Url:            aws.String(issuerURL.String()),
	}
	provider, err := s.IAMClient.CreateOpenIDConnectProvider(&input)
	if err != nil {
		return "", errors.Wrap(err, "error creating provider")
	}
	return *provider.OpenIDConnectProviderArn, nil
}

func fetchRootCAThumbprint(issuerURL string) (string, error) {
	response, err := http.Get(issuerURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	rootCA := response.TLS.PeerCertificates[len(response.TLS.PeerCertificates)-1]
	sha1Sum := sha1.Sum(rootCA.Raw) //nolint:gosec
	return hex.EncodeToString(sha1Sum[:]), nil
}

func (s *IAMService) DeleteOIDCProvider(arn *string) error {
	input := iam.DeleteOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: arn,
	}

	_, err := s.IAMClient.DeleteOpenIDConnectProvider(&input)
	if err != nil {
		return errors.Wrap(err, "error deleting provider")
	}
	return nil
}
