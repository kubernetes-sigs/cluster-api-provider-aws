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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/iam"
	go_cfn "github.com/awslabs/goformation/v4/cloudformation"
	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
)

// Service defines methods for managing the AWS IAM resources using the AWS SDK.
type Service interface {
	CreateResources(t go_cfn.Template, tags map[string]string) error
	DeleteResources(t go_cfn.Template, tags map[string]string) error
}

type serviceImpl struct {
	IAM *iam.IAM
}

// New creates a new IAM service object to interact with the AWS SDK.
func New(iamSvc *iam.IAM) Service {
	return &serviceImpl{
		IAM: iamSvc,
	}
}

func createClient() (*iam.IAM, error) {
	s, err := session.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "internal server error")
	}
	return iam.New(s), nil
}

func prioritySet(t go_cfn.Template) (rmap map[string][]go_cfn.Resource, err error) {
	rmap = map[string][]go_cfn.Resource{}
	for _, resource := range t.Resources {
		if resource.AWSCloudFormationType() == configservice.ResourceTypeAwsIamRole {
			rmap["roles"] = append(rmap["roles"], resource)
		} else if resource.AWSCloudFormationType() == "AWS::IAM::InstanceProfile" {
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
func (s *serviceImpl) CreateResources(t go_cfn.Template, tags map[string]string) error {
	client, err := createClient()
	if err != nil {
		return err
	}
	rmap, err := prioritySet(t)
	if err != nil {
		return err
	}
	for _, resource := range rmap["roles"] {
		err := CreateRole(resource, tags, client)
		if err != nil {
			return err
		}
	}
	for _, resource := range rmap["instanceProfiles"] {
		err := CreateInstanceProfile(resource, tags, client)
		if err != nil {
			return err
		}
	}
	for _, resource := range rmap["policies"] {
		err := CreatePolicy(resource, tags, client)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateRole creates a new CAPA Managed IAM Role and attaches it to the policies as defined in the bootstrap configuration file.
func CreateRole(resource go_cfn.Resource, tags map[string]string, client *iam.IAM) error {
	res := resource.(*cfn_iam.Role)
	tgs := []*iam.Tag{}
	for k, v := range tags {
		tag := iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		tgs = append(tgs, &tag)
	}
	rawdata := res.AssumeRolePolicyDocument.(*iamv1.PolicyDocument)
	data, err := json.Marshal(rawdata)
	if err != nil {
		return errors.Wrapf(err, "corrupt policy document format for IAM role \"%s\"", res.RoleName)
	}
	_, err = client.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(string(data)),
		Description:              &res.Description,
		RoleName:                 &res.RoleName,
		Tags:                     tgs,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				klog.Warningf("IAM role \"%s\" already exists", res.RoleName)
			default:
				return errors.Wrapf(err, "failed to create IAM role \"%s\"", res.RoleName)
			}
		}
	}
	err = attachPoliciesToRole(&res.RoleName, res.ManagedPolicyArns, client)
	if err != nil {
		return err
	}
	klog.V(2).Infof("created \"%s\" CAPA managed IAM role", res.RoleName)
	return nil
}

// CreateInstanceProfile creates a new CAPA Managed Instance Profile and attaches it to the role as defined in the bootstrap configuration file.
func CreateInstanceProfile(resource go_cfn.Resource, tags map[string]string, client *iam.IAM) error {
	res := resource.(*cfn_iam.InstanceProfile)
	tgs := []*iam.Tag{}
	for k, v := range tags {
		tag := iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		tgs = append(tgs, &tag)
	}
	_, err := client.CreateInstanceProfile(&iam.CreateInstanceProfileInput{
		InstanceProfileName: &res.InstanceProfileName,
		Tags:                tgs,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				klog.Warningf("instance profile \"%s\" already exists", res.InstanceProfileName)
			default:
				return errors.Wrapf(err, "failed to create instance profile \"%s\"", res.InstanceProfileName)
			}
		}
	}
	err = attachRoleToInstanceProf(resource, client)
	if err != nil {
		return err
	}
	klog.V(2).Infof("created \"%s\" CAPA managed instance profile", res.InstanceProfileName)
	return nil
}

// CreatePolicy creates a new CAPA Managed IAM Policy and attaches it to the roles as defined in the bootstrap configuration file.
func CreatePolicy(resource go_cfn.Resource, tags map[string]string, client *iam.IAM) error {
	res := resource.(*cfn_iam.ManagedPolicy)
	tgs := []*iam.Tag{}
	for k, v := range tags {
		tag := iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}
		tgs = append(tgs, &tag)
	}
	rawdata := res.PolicyDocument.(*iamv1.PolicyDocument)
	data, err := json.Marshal(rawdata)
	if err != nil {
		return errors.Wrapf(err, "corrupt policy document format for policy \"%s\"", res.ManagedPolicyName)
	}
	create, err := client.CreatePolicy(&iam.CreatePolicyInput{
		Description:    &res.Description,
		PolicyDocument: aws.String(string(data)),
		PolicyName:     &res.ManagedPolicyName,
		Tags:           tgs,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				klog.Warningf("policy \"%s\" already exists", res.ManagedPolicyName)
				policies, err := listpolicies(client)
				if err != nil {
					return err
				}
				for _, policy := range policies {
					if *policy.PolicyName == res.ManagedPolicyName {
						return attachRolesToPolicy(policy, res.Roles, client)
					}
				}
			default:
				return errors.Wrapf(err, "failed to create CAPA managed IAM policy \"%s\"", res.ManagedPolicyName)
			}
		}
	}
	err = attachRolesToPolicy(create.Policy, res.Roles, client)
	if err != nil {
		return err
	}
	klog.V(2).Infof("created \"%s\" CAPA managed IAM policy", res.ManagedPolicyName)
	return nil
}

func attachRoleToInstanceProf(resource go_cfn.Resource, client *iam.IAM) error {
	res := resource.(*cfn_iam.InstanceProfile)
	roleName, err := getRoleName(res.Roles[0])
	if err != nil {
		return err
	}
	_, err = client.AddRoleToInstanceProfile(&iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: &res.InstanceProfileName,
		RoleName:            &roleName,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException, iam.ErrCodeLimitExceededException:
				klog.Warningf("instance profile \"%s\" is already attached to its IAM role", res.InstanceProfileName)
			default:
				return errors.Wrapf(err, "failed to attach instance profile \"%s\" to IAM role \"%s\"", res.InstanceProfileName, roleName)
			}
		}
	}
	return nil
}

func attachPoliciesToRole(rolename *string, awsManagedPolicies []string, client *iam.IAM) error {
	if awsManagedPolicies == nil {
		// klog.Warningf("no policies defined to attach to the IAM role \"%s\"", *rolename) // TODO
		return nil
	}
	for _, policyArn := range awsManagedPolicies {
		// making a copy of policyArn to avoid implicit memory aliasing
		policy := policyArn
		_, err := client.AttachRolePolicy(&iam.AttachRolePolicyInput{
			RoleName:  rolename,
			PolicyArn: &policy,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeEntityAlreadyExistsException:
					klog.Warningf("IAM role \"%s\" is already attached to policy", *rolename) // TODO should we output the policy arn? how safe is it
					continue
				default:
					return errors.Wrapf(err, "failed to attach IAM role \"%s\" to policy", *rolename) // TODO should we output the policy arn? how safe is it
				}
			}
		}
	}
	return nil
}

func attachRolesToPolicy(policy *iam.Policy, roles []string, client *iam.IAM) error {
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
		_, err = client.AttachRolePolicy(&iam.AttachRolePolicyInput{
			PolicyArn: policyarn,
			RoleName:  &roleName,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeEntityAlreadyExistsException:
					klog.Warningf("policy \"%s\" is already attached to IAM role \"%s\"", *policy.PolicyName, roleName)
					continue
				default:
					return errors.Wrapf(err, "failed to attach policy \"%s\" to IAM role \"%s\"", *policy.PolicyName, roleName)
				}
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
	if roleRef == "AWSIAMRoleControllers" {
		roleName = fmt.Sprintf("controllers%s", iamv1.DefaultNameSuffix)
	} else if roleRef == "AWSIAMRoleNodes" {
		roleName = fmt.Sprintf("nodes%s", iamv1.DefaultNameSuffix)
	} else if roleRef == "AWSIAMRoleEKSControlPlane" {
		roleName = fmt.Sprintf("eks-controlplane%s", iamv1.DefaultNameSuffix)
	} else if roleRef == "AWSIAMRoleControlPlane" {
		roleName = fmt.Sprintf("control-plane%s", iamv1.DefaultNameSuffix)
	} else {
		return "", fmt.Errorf("unrecognised or no role found: \"%s\"", roleName)
	}
	return roleName, nil
}

func listpolicies(client *iam.IAM) ([]*iam.Policy, error) {
	list, err := client.ListPolicies(&iam.ListPoliciesInput{
		OnlyAttached: aws.Bool(false),
		Scope:        aws.String("Local"),
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
func (s *serviceImpl) DeleteResources(t go_cfn.Template, tags map[string]string) error {
	client, err := createClient()
	if err != nil {
		return err
	}
	rmap, err := prioritySet(t)
	if err != nil {
		return err
	}
	for _, resource := range rmap["instanceProfiles"] {
		err := DeleteInstanceProfile(resource, client)
		if err != nil {
			return err
		}
	}
	for _, resource := range rmap["roles"] {
		err := DeleteRole(resource, client)
		if err != nil {
			return err
		}
	}
	policies, err := listpolicies(client)
	if err != nil {
		return err
	}
	for _, resource := range rmap["policies"] {
		templatePolicy := resource.(*cfn_iam.ManagedPolicy)
		for _, policy := range policies {
			if templatePolicy.ManagedPolicyName == *policy.PolicyName {
				err := DeletePolicy(policy, client)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// DeleteRole securely deletes CAPA Managed IAM Role.
func DeleteRole(resource go_cfn.Resource, client *iam.IAM) error {
	res := resource.(*cfn_iam.Role)
	_, err := client.GetRole(&iam.GetRoleInput{
		RoleName: &res.RoleName,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				klog.Warningf("\"%s\" IAM role does not exist", res.RoleName)
				return nil
			default:
				return errors.Wrapf(err, "failed to get \"%s\" IAM role", res.RoleName)
			}
		}
	}
	rolePolicies, err := client.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{
		RoleName: &res.RoleName,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to list policies attached to the \"%s\" IAM role", res.RoleName)
	}
	if rolePolicies.AttachedPolicies != nil {
		for _, policy := range rolePolicies.AttachedPolicies {
			_, err := client.DetachRolePolicy(&iam.DetachRolePolicyInput{
				RoleName:  &res.RoleName,
				PolicyArn: policy.PolicyArn,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to detach \"%s\" IAM role from \"%s\" policy", res.RoleName, *policy.PolicyArn)
			}
		}
	}
	_, err = client.DeleteRole(&iam.DeleteRoleInput{
		RoleName: &res.RoleName,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
			default:
				klog.Warningf("\"%s\" IAM role does not exist", res.RoleName)
				return errors.Wrapf(err, "failed to delete \"%s\" IAM role", res.RoleName)
			}
		}
	}
	klog.V(2).Infof("deleted \"%s\" CAPA managed IAM role", res.RoleName)
	return nil
}

// DeleteInstanceProfile securely deletes CAPA Managed IAM Instance Profile.
func DeleteInstanceProfile(resource go_cfn.Resource, client *iam.IAM) error {
	res := resource.(*cfn_iam.InstanceProfile)
	instanceProfileExists, err := client.GetInstanceProfile(&iam.GetInstanceProfileInput{
		InstanceProfileName: &res.InstanceProfileName,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				klog.Warningf("\"%s\" instance profile does not exist", res.InstanceProfileName)
				return nil
			default:
				return errors.Wrapf(err, "failed to get \"%s\" instance profile", res.InstanceProfileName)
			}
		}
	}
	_, err = client.RemoveRoleFromInstanceProfile(&iam.RemoveRoleFromInstanceProfileInput{
		InstanceProfileName: instanceProfileExists.InstanceProfile.InstanceProfileName,
		RoleName:            &res.InstanceProfileName,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to detach \"%s\" IAM role from \"%s\" instance profile", res.InstanceProfileName, res.InstanceProfileName)
	}
	_, err = client.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: &res.InstanceProfileName,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				klog.Warningf("\"%s\" instance profile does not exist", res.InstanceProfileName)
				return nil
			default:
				return errors.Wrapf(err, "failed to delete \"%s\" instance profile", res.InstanceProfileName)
			}
		}
	}
	klog.V(2).Infof("deleted \"%s\" CAPA managed instance profile", res.InstanceProfileName)
	return nil
}

// DeletePolicy securely deletes CAPA Managed IAM Policy.
func DeletePolicy(policy *iam.Policy, client *iam.IAM) error {
	_, err := client.DeletePolicy(&iam.DeletePolicyInput{
		PolicyArn: policy.Arn,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				klog.Warningf("\"%s\" IAM policy does not exist", *policy.Arn)
				return nil
			default:
				return errors.Wrapf(err, "failed to delete \"%s\" IAM policy", *policy.PolicyName)
			}
		}
	}
	klog.V(2).Infof("deleted \"%s\" CAPA managed IAM policy", *policy.PolicyName)
	return nil
}
