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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/cloud/aws/services/ec2"
)

type vpcs struct{}

func (v *vpcs) DescribeVpcs(i *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return &ec2.DescribeVpcsOutput{
		Vpcs: []*ec2.Vpc{
			&ec2.Vpc{
				VpcId: aws.String("world"),
			},
		},
	}, nil
}

func TestReconcileVPC(t *testing.T) {
	s := ec2svc.Service{
		Instances: nil,
		VPCs:      &vpcs{},
	}

	id := "world"
	vpc, err := s.ReconcileVPC(id)
	if err != nil {
		t.Fatalf("got an unexpected error: %v", err)
	}
	if vpc.ID != "world" {
		t.Fatalf("Expected an id of %v but found %v", "world", vpc.ID)
	}
}
