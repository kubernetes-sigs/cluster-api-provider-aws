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
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	go_cfn "github.com/awslabs/goformation/v4/cloudformation"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
)

const (
	// MaxWaitCreateUpdateDelete is the default maximum amount of time to wait for a cfn stack to complete.
	MaxWaitCreateUpdateDelete = 30 * time.Minute
)

// CFNAPI defines the CFN API interface.
type CFNAPI interface {
	CreateStack(ctx context.Context, params *cfn.CreateStackInput, optFns ...func(*cfn.Options)) (*cfn.CreateStackOutput, error)
	DeleteStack(ctx context.Context, params *cfn.DeleteStackInput, optFns ...func(*cfn.Options)) (*cfn.DeleteStackOutput, error)
	DescribeStacks(ctx context.Context, params *cfn.DescribeStacksInput, optFns ...func(*cfn.Options)) (*cfn.DescribeStacksOutput, error)
	DescribeStackResources(ctx context.Context, params *cfn.DescribeStackResourcesInput, optFns ...func(*cfn.Options)) (*cfn.DescribeStackResourcesOutput, error)
	UpdateStack(ctx context.Context, params *cfn.UpdateStackInput, optFns ...func(*cfn.Options)) (*cfn.UpdateStackOutput, error)

	// Waiters for CFN stacks
	WaitUntilStackCreateComplete(ctx context.Context, input *cfn.DescribeStacksInput, maxWait time.Duration) error
	WaitUntilStackUpdateComplete(ctx context.Context, input *cfn.DescribeStacksInput, maxWait time.Duration) error
	WaitUntilStackDeleteComplete(ctx context.Context, input *cfn.DescribeStacksInput, maxWait time.Duration) error
}

// CFNClient is a wrapper over cfn.Client for implementing custom methods of CFNAPI.
type CFNClient struct {
	*cfn.Client
}

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	CFN CFNAPI
}

// NewService returns a new service given the CloudFormation api client.
func NewService(i CFNAPI) *Service {
	return &Service{
		CFN: i,
	}
}

// ReconcileBootstrapStack creates or updates bootstrap CloudFormation.
func (s *Service) ReconcileBootstrapStack(ctx context.Context, stackName string, t go_cfn.Template, tags map[string]string) error {
	yaml, err := t.YAML()
	processedYaml := string(yaml)
	if err != nil {
		return errors.Wrap(err, "failed to generate AWS CloudFormation YAML")
	}

	stackTags := []cfntypes.Tag{}
	for k, v := range tags {
		stackTags = append(stackTags, cfntypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	//nolint:nestif
	if err := s.createStack(ctx, stackName, processedYaml, stackTags); err != nil {
		if code, _ := awserrors.Code(errors.Cause(err)); code == (&cfntypes.AlreadyExistsException{}).ErrorCode() {
			klog.Infof("AWS Cloudformation stack %q already exists, updating", klog.KRef("", stackName))
			updateErr := s.updateStack(ctx, stackName, processedYaml, stackTags)
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

// ReconcileBootstrapNoUpdate creates or updates bootstrap CloudFormation without updating the stack.
func (s *Service) ReconcileBootstrapNoUpdate(ctx context.Context, stackName string, t go_cfn.Template, tags map[string]string) error {
	yaml, err := t.YAML()
	processedYaml := string(yaml)
	if err != nil {
		return errors.Wrap(err, "failed to generate AWS CloudFormation YAML")
	}

	stackTags := []cfntypes.Tag{}
	for k, v := range tags {
		stackTags = append(stackTags, cfntypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	//nolint:nestif
	if err := s.createStack(ctx, stackName, processedYaml, stackTags); err != nil {
		if code, _ := awserrors.Code(errors.Cause(err)); code == (&cfntypes.AlreadyExistsException{}).ErrorCode() {
			desInput := &cfn.DescribeStacksInput{StackName: aws.String(stackName)}
			if err := s.CFN.WaitUntilStackCreateComplete(ctx, desInput, MaxWaitCreateUpdateDelete); err != nil {
				return errors.Wrap(err, "failed to wait for AWS CloudFormation stack to be CreateComplete")
			}
			return nil
		}
		return fmt.Errorf("failed to create CF stack: %w", err)
	}
	return nil
}

func (s *Service) createStack(ctx context.Context, stackName, yaml string, tags []cfntypes.Tag) error {
	input := &cfn.CreateStackInput{
		Capabilities: []cfntypes.Capability{cfntypes.CapabilityCapabilityIam, cfntypes.CapabilityCapabilityNamedIam},
		TemplateBody: aws.String(yaml),
		StackName:    aws.String(stackName),
		Tags:         tags,
	}
	klog.V(2).Infof("creating AWS CloudFormation stack %q", stackName)
	if _, err := s.CFN.CreateStack(ctx, input); err != nil {
		return errors.Wrap(err, "failed to create AWS CloudFormation stack")
	}

	desInput := &cfn.DescribeStacksInput{StackName: aws.String(stackName)}
	klog.V(2).Infof("waiting for stack %q to create", stackName)
	if err := s.CFN.WaitUntilStackCreateComplete(ctx, desInput, MaxWaitCreateUpdateDelete); err != nil {
		return errors.Wrap(err, "failed to wait for AWS CloudFormation stack to be CreateComplete")
	}

	klog.V(2).Infof("stack %q created", stackName)
	return nil
}

func (s *Service) updateStack(ctx context.Context, stackName, yaml string, tags []cfntypes.Tag) error {
	input := &cfn.UpdateStackInput{
		Capabilities: []cfntypes.Capability{cfntypes.CapabilityCapabilityIam, cfntypes.CapabilityCapabilityNamedIam},
		TemplateBody: aws.String(yaml),
		StackName:    aws.String(stackName),
		Tags:         tags,
	}
	klog.V(2).Infof("updating AWS CloudFormation stack %q", stackName)
	if _, err := s.CFN.UpdateStack(ctx, input); err != nil {
		return errors.Wrap(err, "failed to update AWS CloudFormation stack")
	}
	desInput := &cfn.DescribeStacksInput{StackName: aws.String(stackName)}
	klog.V(2).Infof("waiting for stack %q to update", stackName)
	if err := s.CFN.WaitUntilStackUpdateComplete(ctx, desInput, MaxWaitCreateUpdateDelete); err != nil {
		return errors.Wrap(err, "failed to update AWS CloudFormation stack")
	}

	klog.V(2).Infof("stack %q updated", stackName)
	return nil
}

// DeleteStack deletes a cloudformation stack.
func (s *Service) DeleteStack(ctx context.Context, stackName string, retainResources []string) error {
	klog.V(2).Infof("deleting AWS CloudFormation stack %q", stackName)
	var err error
	if retainResources == nil {
		_, err = s.CFN.DeleteStack(ctx, &cfn.DeleteStackInput{StackName: aws.String(stackName)})
	} else {
		_, err = s.CFN.DeleteStack(ctx, &cfn.DeleteStackInput{StackName: aws.String(stackName), RetainResources: retainResources})
	}
	if err != nil {
		return errors.Wrap(err, "failed to delete AWS CloudFormation stack")
	}

	klog.V(2).Infof("waiting for stack %q to delete", stackName)
	if err := s.CFN.WaitUntilStackDeleteComplete(ctx, &cfn.DescribeStacksInput{StackName: aws.String(stackName)}, MaxWaitCreateUpdateDelete); err != nil {
		return errors.Wrap(err, "failed to delete AWS CloudFormation stack")
	}

	klog.V(2).Infof("stack %q deleted", stackName)
	return nil
}

// ShowStackResources prints out in tabular format the resources in the
// stack.
func (s *Service) ShowStackResources(ctx context.Context, stackName string) error {
	input := &cfn.DescribeStackResourcesInput{
		StackName: aws.String(stackName),
	}
	out, err := s.CFN.DescribeStackResources(ctx, input)
	if err != nil {
		return errors.Wrap(err, "unable to describe stack resources")
	}

	fmt.Print("\nFollowing resources are in the stack: \n\n")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "Resource\tType\tStatus")

	for _, r := range out.StackResources {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			aws.ToString(r.ResourceType),
			aws.ToString(r.PhysicalResourceId),
			r.ResourceStatus)

		switch r.ResourceStatus {
		case cfntypes.ResourceStatusCreateComplete, cfntypes.ResourceStatusUpdateComplete:
			continue
		default:
			fmt.Println(aws.ToString(r.ResourceStatusReason))
		}
	}

	w.Flush()

	fmt.Print("\n\n")

	return nil
}

// WaitUntilStackCreateComplete is blocking function to wait until CFN Stack is successfully created.
func (c *CFNClient) WaitUntilStackCreateComplete(ctx context.Context, input *cfn.DescribeStacksInput, maxWait time.Duration) error {
	waiter := cfn.NewStackCreateCompleteWaiter(c, func(o *cfn.StackCreateCompleteWaiterOptions) {
		o.LogWaitAttempts = true
	})

	return waiter.Wait(ctx, input, maxWait)
}

// WaitUntilStackUpdateComplete is blocking function to wait until CFN Stack is successfully updated.
func (c *CFNClient) WaitUntilStackUpdateComplete(ctx context.Context, input *cfn.DescribeStacksInput, maxWait time.Duration) error {
	waiter := cfn.NewStackUpdateCompleteWaiter(c, func(o *cfn.StackUpdateCompleteWaiterOptions) {
		o.LogWaitAttempts = true
	})

	return waiter.Wait(ctx, input, maxWait)
}

// WaitUntilStackDeleteComplete is blocking function to wait until CFN Stack is successfully deleted.
func (c *CFNClient) WaitUntilStackDeleteComplete(ctx context.Context, input *cfn.DescribeStacksInput, maxWait time.Duration) error {
	waiter := cfn.NewStackDeleteCompleteWaiter(c, func(o *cfn.StackDeleteCompleteWaiterOptions) {
		o.LogWaitAttempts = true
	})

	return waiter.Wait(ctx, input, maxWait)
}
