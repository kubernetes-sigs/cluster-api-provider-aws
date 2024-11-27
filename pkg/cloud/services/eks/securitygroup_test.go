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

package eks

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/aws/aws-sdk-go/service/ec2"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/converters"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/eks/mock_eksiface"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/iamauth/mock_iamauth"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestDesiredTags(t *testing.T) {
	g := NewWithT(t)

	existingTags := infrav1.Tags{
		"aws:eks:cluster-name":               "cluster-name",
		"kubernetes.io/cluster/cluster-name": "owned",
		"Name":                               "eks-cluster-sg-cluster-name-1488456827",
		"tagKeyNotToBeIncluded":              "tagValueNotToBeIncluded",
		"sampleTag":                          "sampleOldValue",
	}

	additionalTags := infrav1.Tags{
		"hello":     "there",
		"new":       "tag",
		"tagKey":    "tagValue",
		"sampleTag": "sampleNewValue",
	}

	expectedDesiredTags := infrav1.Tags{
		"aws:eks:cluster-name":               "cluster-name",
		"kubernetes.io/cluster/cluster-name": "owned",
		"Name":                               "eks-cluster-sg-cluster-name-1488456827",
		"hello":                              "there",
		"new":                                "tag",
		"tagKey":                             "tagValue",
		"sampleTag":                          "sampleNewValue",
	}

	desiredTags := desiredTags(existingTags, additionalTags)
	g.Expect(desiredTags).To(Equal(expectedDesiredTags))
}

func TestUpdateTagsForEKSManagedSecurityGroup(t *testing.T) {
	g := NewWithT(t)

	mockControl := gomock.NewController(t)
	defer mockControl.Finish()

	eksMock := mock_eksiface.NewMockEKSAPI(mockControl)
	iamMock := mock_iamauth.NewMockIAMAPI(mockControl)
	ec2Mock := mocks.NewMockEC2API(mockControl)

	scheme := runtime.NewScheme()
	_ = infrav1.AddToScheme(scheme)
	_ = ekscontrolplanev1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	vpcSpec := infrav1.VPCSpec{
		IPv6: &infrav1.IPv6{
			CidrBlock: "2001:db8:85a3::/56",
		},
	}
	scope, err := scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
		Client: client,
		Cluster: &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "ns",
				Name:      "cluster-name",
			},
		},
		ControlPlane: &ekscontrolplanev1.AWSManagedControlPlane{
			Spec: ekscontrolplanev1.AWSManagedControlPlaneSpec{
				RoleName: ptr.To[string]("arn-role"),
				Version:  aws.String("1.29"),
				Region:   "us-east-1",
				NetworkSpec: infrav1.NetworkSpec{
					Subnets: []infrav1.SubnetSpec{
						{
							ID:               "sub-1",
							CidrBlock:        "10.0.10.0/24",
							AvailabilityZone: "us-west-2a",
							IsPublic:         true,
							IsIPv6:           true,
							IPv6CidrBlock:    "2001:db8:85a3:1::/64",
						},
						{
							ID:               "sub-2",
							CidrBlock:        "10.0.10.0/24",
							AvailabilityZone: "us-west-2b",
							IsPublic:         false,
							IsIPv6:           true,
							IPv6CidrBlock:    "2001:db8:85a3:2::/64",
						},
					},
					VPC: vpcSpec,
				},
			},
		},
	})
	g.Expect(err).To(BeNil())

	s := NewService(scope)
	s.EKSClient = eksMock
	s.IAMClient = iamMock
	s.EC2Client = ec2Mock

	existingTags := infrav1.Tags{
		"aws:eks:cluster-name":               "cluster-name",
		"kubernetes.io/cluster/cluster-name": "owned",
		"Name":                               "eks-cluster-sg-cluster-name-1488456827",
	}

	desiredTags := infrav1.Tags{
		"hello":                              "there",
		"new":                                "tag",
		"aws:eks:cluster-name":               "cluster-name",
		"kubernetes.io/cluster/cluster-name": "owned",
		"Name":                               "eks-cluster-sg-cluster-name-1488456827",
		"tagKey":                             "tagValue",
		"sampleTag":                          "sampleNewValue",
	}

	sampleEKSecurityGroupID := aws.String("sg-025f2495c64d5")

	createSGinput := &ec2.CreateSecurityGroupInput{
		GroupName: sampleEKSecurityGroupID,
		TagSpecifications: []*ec2.TagSpecification{
			{
				Tags: converters.MapToTags(desiredTags),
			},
		},
	}

	ec2Mock.EXPECT().CreateSecurityGroup(createSGinput).Return(&ec2.CreateSecurityGroupOutput{
		GroupId: sampleEKSecurityGroupID,
	}, nil)

	_, err = ec2Mock.CreateSecurityGroup(createSGinput)
	g.Expect(err).To(BeNil())

	ec2Mock.EXPECT().CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{sampleEKSecurityGroupID},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("hello"),
				Value: aws.String("there"),
			},
			{
				Key:   aws.String("new"),
				Value: aws.String("tag"),
			},
			{
				Key:   aws.String("sampleTag"),
				Value: aws.String("sampleNewValue"),
			},
			{
				Key:   aws.String("tagKey"),
				Value: aws.String("tagValue"),
			},
		},
	}).Return(nil, nil)

	err = s.updateTagsForEKSManagedSecurityGroup(sampleEKSecurityGroupID, existingTags, desiredTags)
	g.Expect(err).To(BeNil())

	describeSGinput := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{sampleEKSecurityGroupID},
	}

	ec2Mock.EXPECT().DescribeSecurityGroups(describeSGinput).Return(&ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: []*ec2.SecurityGroup{
			{
				GroupId:     sampleEKSecurityGroupID,
				Tags:        converters.MapToTags(desiredTags),
				Description: aws.String("EKS created security group applied to ENI that is attached to EKS Control Plane master nodes, as well as any managed workloads."),
			},
		},
	}, nil)

	describeOutput, err := ec2Mock.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{sampleEKSecurityGroupID},
	})
	g.Expect(err).To(BeNil())

	// Denotes that this EKS security group is created by AWS.
	g.Expect(*describeOutput.SecurityGroups[0].Description).To(Equal("EKS created security group applied to ENI that is attached to EKS Control Plane master nodes, as well as any managed workloads."))

	g.Expect(len(describeOutput.SecurityGroups[0].Tags)).To(Equal(len(desiredTags)))

	for key, value := range desiredTags {
		g.Expect(converters.TagsToMap(describeOutput.SecurityGroups[0].Tags)).To(HaveKeyWithValue(key, value))
	}

}
