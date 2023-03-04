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

// Package cloudformation provides the API operation methods for making requests to
// AWS CloudFormation.
package cloudformation

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	go_cfn "github.com/awslabs/goformation/v4/cloudformation"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	CFN cloudformationiface.CloudFormationAPI
}

// NewService returns a new service given the CloudFormation api client.
func NewService(i cloudformationiface.CloudFormationAPI) *Service {
	return &Service{
		CFN: i,
	}
}

// ReconcileBootstrapStack creates or updates bootstrap CloudFormation.
func (s *Service) ReconcileBootstrapStack(stackName string, t go_cfn.Template, tags map[string]string) error {
	yaml, err := t.YAML()
	processedYaml := string(yaml)
	if err != nil {
		return errors.Wrap(err, "failed to generate AWS CloudFormation YAML")
	}

	stackTags := []*cfn.Tag{}
	for k, v := range tags {
		stackTags = append(stackTags, &cfn.Tag{
			Key:   pointer.String(k),
			Value: pointer.String(v),
		})
	}
	//nolint:nestif
	if err := s.createStack(stackName, processedYaml, stackTags); err != nil {
		if code, _ := awserrors.Code(errors.Cause(err)); code == "AlreadyExistsException" {
			klog.Infof("AWS Cloudformation stack %q already exists, updating", klog.KRef("", stackName))
			updateErr := s.updateStack(stackName, processedYaml, stackTags)
			if updateErr != nil {
				code, ok := awserrors.Code(errors.Cause(updateErr))
				message := awserrors.Message(errors.Cause(updateErr))
				if !ok || code != "ValidationError" || message != "No updates are to be performed." {
					return updateErr
				}
			}
			return nil
		}
		return err
	}
	return nil
}

func (s *Service) ReconcileBootstrapNoUpdate(stackName string, t go_cfn.Template, tags map[string]string) error {
	yaml, err := t.YAML()
	processedYaml := string(yaml)
	if err != nil {
		return errors.Wrap(err, "failed to generate AWS CloudFormation YAML")
	}

	stackTags := []*cfn.Tag{}
	for k, v := range tags {
		stackTags = append(stackTags, &cfn.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	//nolint:nestif
	if err := s.createStack(stackName, processedYaml, stackTags); err != nil {
		if code, _ := awserrors.Code(errors.Cause(err)); code == "AlreadyExistsException" {
			desInput := &cfn.DescribeStacksInput{StackName: aws.String(stackName)}
			if err := s.CFN.WaitUntilStackCreateComplete(desInput); err != nil {
				return errors.Wrap(err, "failed to wait for AWS CloudFormation stack to be CreateComplete")
			}
			return nil
		}
		return fmt.Errorf("failed to create CF stack: %w", err)
	}
	return nil
}

func (s *Service) createStack(stackName, yaml string, tags []*cfn.Tag) error {
	input := &cfn.CreateStackInput{
		Capabilities: aws.StringSlice([]string{cfn.CapabilityCapabilityIam, cfn.CapabilityCapabilityNamedIam}),
		TemplateBody: aws.String(yaml),
		StackName:    aws.String(stackName),
		Tags:         tags,
	}
	klog.V(2).Infof("creating AWS CloudFormation stack %q", stackName)
	if _, err := s.CFN.CreateStack(input); err != nil {
		return errors.Wrap(err, "failed to create AWS CloudFormation stack")
	}

	desInput := &cfn.DescribeStacksInput{StackName: aws.String(stackName)}
	klog.V(2).Infof("waiting for stack %q to create", stackName)
	if err := s.CFN.WaitUntilStackCreateComplete(desInput); err != nil {
		return errors.Wrap(err, "failed to wait for AWS CloudFormation stack to be CreateComplete")
	}

	klog.V(2).Infof("stack %q created", stackName)
	return nil
}

func (s *Service) updateStack(stackName, yaml string, tags []*cfn.Tag) error {
	input := &cfn.UpdateStackInput{
		Capabilities: aws.StringSlice([]string{cfn.CapabilityCapabilityIam, cfn.CapabilityCapabilityNamedIam}),
		TemplateBody: aws.String(yaml),
		StackName:    aws.String(stackName),
		Tags:         tags,
	}
	klog.V(2).Infof("updating AWS CloudFormation stack %q", stackName)
	if _, err := s.CFN.UpdateStack(input); err != nil {
		return errors.Wrap(err, "failed to update AWS CloudFormation stack")
	}
	desInput := &cfn.DescribeStacksInput{StackName: aws.String(stackName)}
	klog.V(2).Infof("waiting for stack %q to update", stackName)
	if err := s.CFN.WaitUntilStackUpdateComplete(desInput); err != nil {
		return errors.Wrap(err, "failed to update AWS CloudFormation stack")
	}

	klog.V(2).Infof("stack %q updated", stackName)
	return nil
}

// DeleteStack deletes a cloudformation stack.
func (s *Service) DeleteStack(stackName string, retainResources []*string) error {
	klog.V(2).Infof("deleting AWS CloudFormation stack %q", stackName)
	var err error
	if retainResources == nil {
		_, err = s.CFN.DeleteStack(&cfn.DeleteStackInput{StackName: aws.String(stackName)})
	} else {
		_, err = s.CFN.DeleteStack(&cfn.DeleteStackInput{StackName: aws.String(stackName), RetainResources: retainResources})
	}
	if err != nil {
		return errors.Wrap(err, "failed to delete AWS CloudFormation stack")
	}

	klog.V(2).Infof("waiting for stack %q to delete", stackName)
	if err := s.CFN.WaitUntilStackDeleteComplete(&cfn.DescribeStacksInput{StackName: aws.String(stackName)}); err != nil {
		return errors.Wrap(err, "failed to delete AWS CloudFormation stack")
	}

	klog.V(2).Infof("stack %q deleted", stackName)
	return nil
}

// ShowStackResources prints out in tabular format the resources in the
// stack.
func (s *Service) ShowStackResources(stackName string) error {
	input := &cfn.DescribeStackResourcesInput{
		StackName: aws.String(stackName),
	}
	out, err := s.CFN.DescribeStackResources(input)
	if err != nil {
		return errors.Wrap(err, "unable to describe stack resources")
	}

	fmt.Print("\nFollowing resources are in the stack: \n\n")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "Resource\tType\tStatus")

	for _, r := range out.StackResources {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			aws.StringValue(r.ResourceType),
			aws.StringValue(r.PhysicalResourceId),
			aws.StringValue(r.ResourceStatus))

		switch aws.StringValue(r.ResourceStatus) {
		case cfn.ResourceStatusCreateComplete, cfn.ResourceStatusUpdateComplete:
			continue
		default:
			fmt.Println(aws.StringValue(r.ResourceStatusReason))
		}
	}

	w.Flush()

	fmt.Print("\n\n")

	return nil
}
