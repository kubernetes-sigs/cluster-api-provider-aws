package iamservice

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/iam"
	go_cfn "github.com/awslabs/goformation/v4/cloudformation"
	cfn_iam "github.com/awslabs/goformation/v4/cloudformation/iam"
	"k8s.io/klog/v2"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
)

type Service interface {
	CreateServices(t go_cfn.Template, tags map[string]string) error
}

type serviceImpl struct {
	IAM *iam.IAM
}

func New(iamSvc *iam.IAM) Service {
	return &serviceImpl{
		IAM: iamSvc,
	}
}

func CreateClient() *iam.IAM {
	s, err := session.NewSession()
	if err != nil {
		fmt.Print(err)
	}
	return iam.New(s)
}

func prioritySet(t go_cfn.Template, client *iam.IAM) (rmap map[string][]go_cfn.Resource, err error) {
	rmap = map[string][]go_cfn.Resource{}
	for _, resource := range t.Resources {
		if resource.AWSCloudFormationType() == configservice.ResourceTypeAwsIamRole {
			rmap["roles"] = append(rmap["roles"], resource)
		} else if resource.AWSCloudFormationType() == "AWS::IAM::InstanceProfile" {
			rmap["instanceProfiles"] = append(rmap["instanceProfiles"], resource)
		} else if resource.AWSCloudFormationType() == "AWS::IAM::ManagedPolicy" {
			rmap["policies"] = append(rmap["policies"], resource)
		} else {
			return nil, fmt.Errorf("unknown resource type: %v", resource)
		}
	}
	return rmap, nil
}

func (s *serviceImpl) CreateServices(t go_cfn.Template, tags map[string]string) error {
	client := CreateClient()
	rmap, err := prioritySet(t, client)
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
		return fmt.Errorf("marshalling \"%s\" metadata: %w", res.RoleName, err)
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
				klog.Warningf("role: \"%s\" already exists, checking if the role is attached to desired policies", res.RoleName)
			// TODO
			// case iam.ErrCodeConcurrentModificationException:
			// case iam.ErrCodeUnmodifiableEntityException:
			// case iam.ErrCodeMalformedPolicyDocumentException:
			default:
				return fmt.Errorf(aerr.Error())
			}
		}
	}
	return attachPoliciesToRole(&res.RoleName, res.ManagedPolicyArns, client)
}

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
				klog.Warningf("instance profile: \"%s\" already exists, checking if the instance profile is attached to the desired role", res.InstanceProfileName)
			// TODO
			// case iam.ErrCodeConcurrentModificationException:
			default:
				return fmt.Errorf(aerr.Error())
			}
		}
	}
	return attachRoleToInstanceProf(resource, client)
}

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
		return fmt.Errorf("marshalling \"%s\" metadata: %w", res.ManagedPolicyName, err)
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
				klog.Warningf("policy: \"%s\" already exists, checking if the policy is attached to desired roles", res.ManagedPolicyName)
				policies, _ := listpolicies(client)
				for _, policy := range policies {
					if *policy.PolicyName == res.ManagedPolicyName {
						return attachRolesToPolicy(policy, res.Roles, client)
					}
				}
			// TODO
			// case iam.ErrCodeConcurrentModificationException:
			// case iam.ErrCodeUnmodifiableEntityException:
			// case iam.ErrCodeMalformedPolicyDocumentException:
			default:
				return fmt.Errorf(aerr.Error())
			}
		}
	}
	return attachRolesToPolicy(create.Policy, res.Roles, client)
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
				klog.Warningf("the instance profile \"%s\" is already attached to the IAM role", res.InstanceProfileName)
			// TODO
			// case iam.ErrCodeNoSuchEntityException:
			default:
				return fmt.Errorf(aerr.Error())
			}
		}
	}
	fmt.Printf("Successfully linked the instance profile \"%s\" to IAM role\n", res.InstanceProfileName)
	return nil
}

func attachPoliciesToRole(rolename *string, awsManagedPolicies []string, client *iam.IAM) error {
	if awsManagedPolicies == nil {
		fmt.Printf("no AWS managed policies set for the IAM role: \"%s\"\n", *rolename)
		return nil
	}
	for _, policy := range awsManagedPolicies {
		_, err := client.AttachRolePolicy(&iam.AttachRolePolicyInput{
			RoleName:  rolename,
			PolicyArn: &policy,
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				// TODO
				// case iam.ErrCodeNoSuchEntityException:
				// case iam.ErrCodePolicyNotAttachableException:
				default:
					return fmt.Errorf(aerr.Error())
				}
			}
		}
		//TODO: this line is printed even if the policy is already attached to the role, change this behaviour
		fmt.Printf("Successfully attached policy \"%s\" to the IAM role \"%s\"\n", policy, *rolename)
	}
	return nil
}

func attachRolesToPolicy(policy *iam.Policy, roles []string, client *iam.IAM) error {
	if roles == nil {
		fmt.Printf("no IAM roles set for the policy: \"%s\"\n", *policy.PolicyName)
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
				// TODO
				// cannot add a case for entities already attached
				// case iam.ErrCodeNoSuchEntityException:
				// case iam.ErrCodePolicyNotAttachableException:
				default:
					return fmt.Errorf(aerr.Error())
				}
			}
		}
		//TODO: this line is printed even if the role is already attached to the policy, change this behaviour
		fmt.Printf("Successfully attached IAM role \"%s\" to the policy \"%s\"\n", roleName, *policy.PolicyName)
	}
	return nil
}

func getRoleName(encodedRole string) (string, error) {
	var roleName string
	bytes, err := base64.RawStdEncoding.DecodeString(encodedRole)
	if err != nil {
		klog.Warningf("decoding \"%s\": %w", encodedRole, err)
	}
	roleRef := strings.Trim(strings.TrimLeft(string(bytes), "{Ref:\\ \""), "\\ \"")
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
		return nil, err
	}
	if list.Policies == nil {
		klog.Warningf("cannot find any CAPA managed policies on the AWS console")
		return nil, nil
	}
	return list.Policies, nil
}
