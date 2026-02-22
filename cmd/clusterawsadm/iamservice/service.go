/*
Copyright 2024 The Kubernetes Authors.

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

// Package iamservice is a redefined way to create and delete IAM instance profiles, policies and roles.
package iamservice

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	configservice "github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"
	go_cfn "github.com/awslabs/goformation/v4/cloudformation"
	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
)

const (
	// EntityAlreadyExists is the AWS IAM error code for when an entity already exists.
	EntityAlreadyExists string = "EntityAlreadyExists"
	// NoSuchEntity is the AWS IAM error code for when an entity does not exist.
	NoSuchEntity string = "NoSuchEntity"
)

// Service defines methods for managing the AWS IAM resources using the AWS SDK.
type Service interface {
	CreateResources(ctx context.Context, t go_cfn.Template, tags map[string]string) error
	DeleteResources(ctx context.Context, t go_cfn.Template, tags map[string]string) error
}

type serviceImpl struct {
	IAM *iam.Client
}

// New creates a new IAM service object to interact with the AWS SDK.
func New(iamSvc *iam.Client) Service {
	return &serviceImpl{
		IAM: iamSvc,
	}
}

func prioritySet(t go_cfn.Template) (rmap map[string][]go_cfn.Resource, err error) {
	rmap = map[string][]go_cfn.Resource{}
	for _, resource := range t.Resources {
		if resource.AWSCloudFormationType() == string(configservice.ResourceTypeRole) {
			rmap["roles"] = append(rmap["roles"], resource)
		} else if resource.AWSCloudFormationType() == string(configservice.ResourceTypeIAMInstanceProfile) {
			rmap["instanceProfiles"] = append(rmap["instanceProfiles"], resource)
		} else if resource.AWSCloudFormationType() == "AWS::IAM::ManagedPolicy" {
			rmap["policies"] = append(rmap["policies"], resource)
		} else {
			return nil, errors.Wrapf(err, "unknown resource type %v", resource)
		}
	}
	return rmap, nil
}

// CreateServices manages the order in which the IAM resources should be created.
func (s *serviceImpl) CreateResources(ctx context.Context, t go_cfn.Template, tags map[string]string) error {
	rmap, err := prioritySet(t)
	if err != nil {
		return err
	}
	for _, resource := range rmap["roles"] {
		err := s.CreateRole(ctx, resource, tags)
		if err != nil {
			return err
		}
	}
	for _, resource := range rmap["instanceProfiles"] {
		err := s.CreateInstanceProfile(ctx, resource, tags)
		if err != nil {
			return err
		}
	}
	for _, resource := range rmap["policies"] {
		err := s.CreatePolicy(ctx, resource, tags)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Created CAPA managed entities.")
	return nil
}

// CreateRole creates a new CAPA Managed IAM Role and attaches it to the policies as defined in the bootstrap configuration file.
func (s *serviceImpl) CreateRole(ctx context.Context, resource go_cfn.Resource, tags map[string]string) error {
	res := resource.(*cfn_iam.Role)
	tgs := []iamtypes.Tag{}
	for k, v := range tags {
		tag := iamtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		tgs = append(tgs, tag)
	}
	rawdata := res.AssumeRolePolicyDocument.(*iamv1.PolicyDocument)
	data, err := json.Marshal(rawdata)
	if err != nil {
		return errors.Wrapf(err, "corrupt policy document format for IAM role \"%s\"", res.RoleName)
	}
	_, err = s.IAM.CreateRole(ctx, &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(string(data)),
		Description:              &res.Description,
		RoleName:                 &res.RoleName,
		Tags:                     tgs,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case EntityAlreadyExists:
				klog.Warningf("IAM role \"%s\" already exists", res.RoleName)
			default:
				return errors.Wrapf(err, "failed to create IAM role \"%s\"", res.RoleName)
			}
		} else {
			return errors.Wrapf(err, "unexpected error occurred")
		}
	}
	err = attachPoliciesToRole(ctx, &res.RoleName, res.ManagedPolicyArns, s.IAM)
	if err != nil {
		return err
	}
	klog.V(2).Infof("created \"%s\" CAPA managed IAM role", res.RoleName)
	return nil
}

// CreateInstanceProfile creates a new CAPA Managed Instance Profile and attaches it to the role as defined in the bootstrap configuration file.
func (s *serviceImpl) CreateInstanceProfile(ctx context.Context, resource go_cfn.Resource, tags map[string]string) error {
	res := resource.(*cfn_iam.InstanceProfile)
	tgs := []iamtypes.Tag{}
	for k, v := range tags {
		tag := iamtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		tgs = append(tgs, tag)
	}
	_, err := s.IAM.CreateInstanceProfile(ctx, &iam.CreateInstanceProfileInput{
		InstanceProfileName: &res.InstanceProfileName,
		Tags:                tgs,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case EntityAlreadyExists:
				klog.Warningf("instance profile \"%s\" already exists", res.InstanceProfileName)
			default:
				return errors.Wrapf(err, "failed to create instance profile \"%s\"", res.InstanceProfileName)
			}
		} else {
			return errors.Wrapf(err, "unexpected error occurred")
		}
	}
	err = attachRoleToInstanceProf(ctx, resource, s.IAM)
	if err != nil {
		return err
	}
	klog.V(2).Infof("created \"%s\" CAPA managed instance profile", res.InstanceProfileName)
	return nil
}

// CreatePolicy creates a new CAPA Managed IAM Policy and attaches it to the roles as defined in the bootstrap configuration file.
func (s *serviceImpl) CreatePolicy(ctx context.Context, resource go_cfn.Resource, tags map[string]string) error {
	res := resource.(*cfn_iam.ManagedPolicy)
	tgs := []iamtypes.Tag{}
	for k, v := range tags {
		tag := iamtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		tgs = append(tgs, tag)
	}
	rawdata := res.PolicyDocument.(*iamv1.PolicyDocument)
	data, err := json.Marshal(rawdata)
	if err != nil {
		return errors.Wrapf(err, "corrupt policy document format for policy \"%s\"", res.ManagedPolicyName)
	}
	create, err := s.IAM.CreatePolicy(ctx, &iam.CreatePolicyInput{
		Description:    &res.Description,
		PolicyDocument: aws.String(string(data)),
		PolicyName:     &res.ManagedPolicyName,
		Tags:           tgs,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case EntityAlreadyExists:
				klog.Warningf("policy \"%s\" already exists", res.ManagedPolicyName)
				policies, err := listpolicies(ctx, s.IAM)
				if err != nil {
					return err
				}
				for _, policy := range policies {
					if *policy.PolicyName == res.ManagedPolicyName {
						return attachRolesToPolicy(ctx, &policy, res.Roles, s.IAM)
					}
				}
			default:
				return errors.Wrapf(err, "failed to create CAPA managed IAM policy \"%s\"", res.ManagedPolicyName)
			}
		} else {
			return errors.Wrapf(err, "unexpected error occurred")
		}
	}
	err = attachRolesToPolicy(ctx, create.Policy, res.Roles, s.IAM)
	if err != nil {
		return err
	}
	klog.V(2).Infof("created \"%s\" CAPA managed IAM policy", res.ManagedPolicyName)
	return nil
}

func attachRoleToInstanceProf(ctx context.Context, resource go_cfn.Resource, client *iam.Client) error {
	res := resource.(*cfn_iam.InstanceProfile)
	roleName, err := getRoleName(res.Roles[0])
	if err != nil {
		return err
	}
	_, err = client.AddRoleToInstanceProfile(ctx, &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: &res.InstanceProfileName,
		RoleName:            &roleName,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case EntityAlreadyExists:
				klog.Warningf("instance profile \"%s\" is already attached to its IAM role", res.InstanceProfileName)
			default:
				return errors.Wrapf(err, "failed to attach instance profile \"%s\" to IAM role \"%s\"", res.InstanceProfileName, roleName)
			}
		} else {
			return errors.Wrapf(err, "unexpected error occurred")
		}
	}
	return nil
}

func attachPoliciesToRole(ctx context.Context, rolename *string, awsManagedPolicies []string, client *iam.Client) error {
	if awsManagedPolicies == nil {
		// klog.Warningf("no policies defined to attach to the IAM role \"%s\"", *rolename) // TODO
		return nil
	}
	for _, policy := range awsManagedPolicies {
		// making a copy of policyArn to avoid implicit memory aliasing
		policyArn := policy
		_, err := client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
			RoleName:  rolename,
			PolicyArn: &policyArn,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				switch apiErr.ErrorCode() {
				case EntityAlreadyExists:
					klog.Warningf("IAM role \"%s\" is already attached to policy", *rolename) // TODO should we output the policy arn?
				default:
					return errors.Wrapf(err, "failed to attach IAM role \"%s\" to policy \"%v\"", *rolename, policyArn) // TODO should we output the policy arn?
				}
			} else {
				return errors.Wrapf(err, "unexpected error occurred")
			}
		}
	}
	return nil
}

func attachRolesToPolicy(ctx context.Context, policy *iamtypes.Policy, roles []string, client *iam.Client) error {
	if roles == nil {
		// klog.Warningf("no IAM roles defined to attach to the policy \"%s\"", *policy.PolicyName) //TODO
		return nil
	}
	policyarn := policy.Arn
	for _, encodedRole := range roles {
		roleName, err := getRoleName(encodedRole)
		if err != nil {
			return err
		}
		_, err = client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
			PolicyArn: policyarn,
			RoleName:  &roleName,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				switch apiErr.ErrorCode() {
				case EntityAlreadyExists:
					klog.Warningf("policy \"%s\" is already attached to IAM role \"%s\"", *policy.PolicyName, roleName)
				default:
					return errors.Wrapf(err, "failed to attach policy \"%s\" to IAM role \"%s\"", *policy.PolicyName, roleName)
				}
			} else {
				return errors.Wrapf(err, "unexpected error occurred")
			}
		}
	}
	return nil
}

func getRoleName(encodedRole string) (string, error) {
	var roleName string
	bytes, err := base64.StdEncoding.DecodeString(encodedRole)
	if err != nil {
		return "", err
	}
	roleRef := string(regexp.MustCompile(`(AWSIAMRole[a-zA-Z]+)`).Find(bytes))
	switch roleRef {
	case "AWSIAMRoleControllers":
		roleName = fmt.Sprintf("controllers%s", iamv1.DefaultNameSuffix)
	case "AWSIAMRoleNodes":
		roleName = fmt.Sprintf("nodes%s", iamv1.DefaultNameSuffix)
	case "AWSIAMRoleEKSControlPlane":
		roleName = fmt.Sprintf("eks-controlplane%s", iamv1.DefaultNameSuffix)
	case "AWSIAMRoleControlPlane":
		roleName = fmt.Sprintf("control-plane%s", iamv1.DefaultNameSuffix)
	default:
		return "", fmt.Errorf("unrecognised or no role found: \"%s\"", roleName)
	}
	return roleName, nil
}

func listpolicies(ctx context.Context, client *iam.Client) ([]iamtypes.Policy, error) {
	list, err := client.ListPolicies(ctx, &iam.ListPoliciesInput{
		OnlyAttached: false,
		Scope:        iamtypes.PolicyScopeType("Local"),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list CAPA managed IAM policies")
	}
	if list.Policies == nil {
		klog.Warningf("no CAPA managed IAM policies detected on the AWS console")
		return nil, nil
	}
	return list.Policies, nil
}

// DeleteServices manages the order in which the IAM resources should be deleted.
func (s *serviceImpl) DeleteResources(ctx context.Context, t go_cfn.Template, tags map[string]string) error {
	rmap, err := prioritySet(t)
	if err != nil {
		return err
	}
	for _, resource := range rmap["instanceProfiles"] {
		err := DeleteInstanceProfile(ctx, resource, s.IAM)
		if err != nil {
			return err
		}
	}
	for _, resource := range rmap["roles"] {
		err := DeleteRole(ctx, resource, s.IAM)
		if err != nil {
			return err
		}
	}
	policies, err := listpolicies(ctx, s.IAM)
	if err != nil {
		return err
	}
	for _, resource := range rmap["policies"] {
		templatePolicy := resource.(*cfn_iam.ManagedPolicy)
		for _, policy := range policies {
			if templatePolicy.ManagedPolicyName == *policy.PolicyName {
				err := DeletePolicy(ctx, &policy, s.IAM)
				if err != nil {
					return err
				}
			}
		}
	}
	fmt.Printf("Deleted CAPA managed entities.")
	return nil
}

// DeleteRole securely deletes CAPA Managed IAM Role.
func DeleteRole(ctx context.Context, resource go_cfn.Resource, client *iam.Client) error {
	res := resource.(*cfn_iam.Role)
	_, err := client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: &res.RoleName,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case NoSuchEntity:
				klog.Warningf("IAM role \"%s\" does not exist", res.RoleName)
				return nil
			default:
				return errors.Wrapf(err, "failed to get \"%s\" IAM role", res.RoleName)
			}
		} else {
			return errors.Wrapf(err, "unexpected error occurred")
		}
	}
	rolePolicies, err := client.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: &res.RoleName,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to list policies attached to the \"%s\" IAM role", res.RoleName)
	}
	if rolePolicies.AttachedPolicies != nil {
		for _, policy := range rolePolicies.AttachedPolicies {
			_, err := client.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
				RoleName:  &res.RoleName,
				PolicyArn: policy.PolicyArn,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to detach \"%s\" IAM role from \"%s\" policy", res.RoleName, *policy.PolicyArn)
			}
		}
	}
	_, err = client.DeleteRole(ctx, &iam.DeleteRoleInput{
		RoleName: &res.RoleName,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case NoSuchEntity:
				klog.Warningf("IAM role \"%s\" does not exist", res.RoleName)
			default:
				return errors.Wrapf(err, "failed to delete \"%s\" IAM role", res.RoleName)
			}
		} else {
			return errors.Wrapf(err, "unexpected error occurred")
		}
	}
	klog.V(2).Infof("deleted \"%s\" CAPA managed IAM role", res.RoleName)
	return nil
}

// DeleteInstanceProfile securely deletes CAPA Managed IAM Instance Profile.
func DeleteInstanceProfile(ctx context.Context, resource go_cfn.Resource, client *iam.Client) error {
	res := resource.(*cfn_iam.InstanceProfile)
	instanceProfileExists, err := client.GetInstanceProfile(ctx, &iam.GetInstanceProfileInput{
		InstanceProfileName: &res.InstanceProfileName,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case NoSuchEntity:
				klog.Warningf("instance profile \"%s\" does not exist", res.InstanceProfileName)
				return nil
			default:
				return errors.Wrapf(err, "failed to get \"%s\" instance profile", res.InstanceProfileName)
			}
		} else {
			return errors.Wrapf(err, "unexpected error occurred")
		}
	}
	_, err = client.RemoveRoleFromInstanceProfile(ctx, &iam.RemoveRoleFromInstanceProfileInput{
		InstanceProfileName: instanceProfileExists.InstanceProfile.InstanceProfileName,
		RoleName:            &res.InstanceProfileName,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to detach \"%s\" IAM role from \"%s\" instance profile", res.InstanceProfileName, res.InstanceProfileName)
	}
	_, err = client.DeleteInstanceProfile(ctx, &iam.DeleteInstanceProfileInput{
		InstanceProfileName: &res.InstanceProfileName,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case NoSuchEntity:
				klog.Warningf("instance profile \"%s\" does not exist", res.InstanceProfileName)
				return nil
			default:
				return errors.Wrapf(err, "failed to delete \"%s\" instance profile", res.InstanceProfileName)
			}
		} else {
			return errors.Wrapf(err, "unexpected error occurred")
		}
	}
	klog.V(2).Infof("deleted \"%s\" CAPA managed instance profile", res.InstanceProfileName)
	return nil
}

// DeletePolicy securely deletes CAPA Managed IAM Policy.
func DeletePolicy(ctx context.Context, policy *iamtypes.Policy, client *iam.Client) error {
	_, err := client.DeletePolicy(ctx, &iam.DeletePolicyInput{
		PolicyArn: policy.Arn,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case NoSuchEntity:
				klog.Warningf("IAM policy \"%s\" does not exist", *policy.Arn)
				return nil
			default:
				return errors.Wrapf(err, "failed to delete \"%s\" IAM policy", *policy.PolicyName)
			}
		} else {
			return errors.Wrapf(err, "unexpected error occurred")
		}
	}
	klog.V(2).Infof("deleted \"%s\" CAPA managed IAM policy", *policy.PolicyName)
	return nil
}
