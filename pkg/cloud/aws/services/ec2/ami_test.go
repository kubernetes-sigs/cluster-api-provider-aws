/*
Copyright 2019 The Kubernetes Authors.

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

package ec2

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/actuators"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2/mock_ec2iface"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func TestAMIs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testCases := []struct {
		name   string
		expect func(m *mock_ec2iface.MockEC2APIMockRecorder)
	}{
		{
			name: "simple test",
			expect: func(m *mock_ec2iface.MockEC2APIMockRecorder) {
				m.DescribeImages(gomock.AssignableToTypeOf(&ec2.DescribeImagesInput{})).
					Return(&ec2.DescribeImagesOutput{
						Images: []*ec2.Image{
							{
								ImageId:      aws.String("ancient"),
								CreationDate: aws.String("2011-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("pretty new"),
								CreationDate: aws.String("2019-02-08T17:02:31.000Z"),
							},
							{
								ImageId:      aws.String("pretty old"),
								CreationDate: aws.String("2014-02-08T17:02:31.000Z"),
							},
						},
					}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec2Mock := mock_ec2iface.NewMockEC2API(mockCtrl)

			scope, err := actuators.NewScope(actuators.ScopeParams{
				Cluster: &clusterv1.Cluster{},
				AWSClients: actuators.AWSClients{
					EC2: ec2Mock,
				},
			})
			if err != nil {
				t.Fatalf("did not expect err: %v", err)
			}

			tc.expect(ec2Mock.EXPECT())

			s := NewService(scope)
			id, err := s.defaultAMILookup("", "base os", "baseos version", "1.11.1")
			if err != nil {
				t.Fatalf("did not expect error calling a mock: %v", err)
			}
			if id != "pretty new" {
				t.Fatalf("returned %q expected 'pretty new'", id)
			}
		})
	}
}
