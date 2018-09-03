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

package ec2_test

import (
	"testing"

	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type instances struct{}

func (i *instances) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{}, nil
}
func (i *instances) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	return &ec2.Reservation{}, nil
}

func TestInstanceIfExists_DoesNotExist(t *testing.T) {
	s := ec2svc.Service{
		Instances: &instances{},
	}
	id := "hello"
	out, err := s.InstanceIfExists(&id)
	if err != nil {
		t.Fatalf("did not expect error: %v", err)
	}
	if out != nil {
		t.Fatal("Did not expect anything but got something")
	}
}
