// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudformation

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

func (s *Service) createStack(stackName string, yaml string) error {

	input := &cfn.CreateStackInput{
		Capabilities: aws.StringSlice([]string{cfn.CapabilityCapabilityIam, cfn.CapabilityCapabilityNamedIam}),
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stackName),
	}
	glog.V(2).Infof("creating AWS CloudFormation stack %q", stackName)
	if _, err := s.CFN.CreateStack(input); err != nil {
		return errors.Wrap(err, "failed to create AWS CloudFormation stack")
	}

	desInput := &cfn.DescribeStacksInput{StackName: aws.String(stackName)}
	glog.V(2).Infof("waiting for stack %q to create", stackName)
	if err := s.CFN.WaitUntilStackCreateComplete(desInput); err != nil {
		return errors.Wrap(err, "failed to create AWS CloudFormation stack")
	}

	glog.V(2).Infof("stack %q created", stackName)
	return nil
}

func (s *Service) updateStack(stackName string, yaml string) error {

	input := &cfn.UpdateStackInput{
		Capabilities: aws.StringSlice([]string{cfn.CapabilityCapabilityIam, cfn.CapabilityCapabilityNamedIam}),
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stackName),
	}
	glog.V(2).Infof("updating AWS CloudFormation stack %q", stackName)
	if _, err := s.CFN.UpdateStack(input); err != nil {
		return errors.Wrap(err, "failed to update AWS CloudFormation stack")
	}
	desInput := &cfn.DescribeStacksInput{StackName: aws.String(stackName)}
	glog.V(2).Infof("waiting for stack %q to update", stackName)
	if err := s.CFN.WaitUntilStackUpdateComplete(desInput); err != nil {
		return errors.Wrap(err, "failed to update AWS CloudFormation stack")
	}

	glog.V(2).Infof("stack %q updated", stackName)
	return nil
}

// ShowStackResources prints out in tabular format the resources in the
// stack
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
