/*
Copyright 2022 The Kubernetes Authors.

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
package controllers

import (
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

const DNSName = "www.google.com"

var (
	lbName          = aws.String("test-cluster-apiserver")
	lbArn           = aws.String("loadbalancer::arn")
	describeLBInput = &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice([]string{"test-cluster-apiserver"}),
	}
	describeLBAttributesInput = &elb.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: lbName,
	}
	describeLBOutput = &elb.DescribeLoadBalancersOutput{
		LoadBalancerDescriptions: []*elb.LoadBalancerDescription{
			{
				Scheme:            aws.String(string(infrav1.ELBSchemeInternetFacing)),
				Subnets:           []*string{aws.String("subnet-1")},
				AvailabilityZones: []*string{aws.String("us-east-1a")},
				VPCId:             aws.String("vpc-exists"),
			},
		},
	}
	describeLBOutputV2 = &elbv2.DescribeLoadBalancersOutput{
		LoadBalancers: []*elbv2.LoadBalancer{
			{
				Scheme: aws.String(string(infrav1.ELBSchemeInternetFacing)),
				AvailabilityZones: []*elbv2.AvailabilityZone{
					{
						SubnetId: aws.String("subnet-1"),
						ZoneName: aws.String("us-east-1a"),
					},
				},
				LoadBalancerArn: aws.String(*lbArn),
				VpcId:           aws.String("vpc-exists"),
				DNSName:         aws.String("dns"),
			},
		},
	}
	describeLBAttributesOutputV2 = &elbv2.DescribeLoadBalancerAttributesOutput{
		Attributes: []*elbv2.LoadBalancerAttribute{
			{
				Key:   aws.String("cross-zone"),
				Value: aws.String("true"),
			},
		},
	}
	describeLBAttributesOutput = &elb.DescribeLoadBalancerAttributesOutput{
		LoadBalancerAttributes: &elb.LoadBalancerAttributes{
			CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{
				Enabled: aws.Bool(false),
			},
		},
	}
	expectedTags = []*elb.Tag{
		{
			Key:   aws.String("Name"),
			Value: lbName,
		},
		{
			Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
			Value: aws.String("owned"),
		},
		{
			Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
			Value: aws.String("apiserver"),
		},
	}
	expectedV2Tags = []*elbv2.Tag{
		{
			Key:   aws.String("Name"),
			Value: lbName,
		},
		{
			Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
			Value: aws.String("owned"),
		},
		{
			Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
			Value: aws.String("apiserver"),
		},
	}
)

func expectAWSClusterConditions(g *WithT, m *infrav1.AWSCluster, expected []conditionAssertion) {
	g.Expect(len(m.Status.Conditions)).To(BeNumerically(">=", len(expected)), "number of conditions")
	for _, c := range expected {
		actual := conditions.Get(m, c.conditionType)
		g.Expect(actual).To(Not(BeNil()))
		g.Expect(actual.Type).To(Equal(c.conditionType))
		g.Expect(actual.Status).To(Equal(c.status))
		g.Expect(actual.Severity).To(Equal(c.severity))
		g.Expect(actual.Reason).To(Equal(c.reason))
	}
}

func getAWSCluster(name, namespace string) infrav1.AWSCluster {
	return infrav1.AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: infrav1.AWSClusterSpec{
			Region: "us-east-1",
			ControlPlaneLoadBalancer: &infrav1.AWSLoadBalancerSpec{
				LoadBalancerType: infrav1.LoadBalancerTypeClassic,
			},
			NetworkSpec: infrav1.NetworkSpec{
				VPC: infrav1.VPCSpec{
					ID:        "vpc-exists",
					CidrBlock: "10.0.0.0/8",
				},
				Subnets: infrav1.Subnets{
					{
						ID:               "subnet-1",
						AvailabilityZone: "us-east-1a",
						CidrBlock:        "10.0.10.0/24",
						IsPublic:         false,
					},
					{
						ID:               "subnet-2",
						AvailabilityZone: "us-east-1c",
						CidrBlock:        "10.0.11.0/24",
						IsPublic:         true,
					},
				},
				SecurityGroupOverrides: map[infrav1.SecurityGroupRole]string{},
			},
			Bastion: infrav1.Bastion{Enabled: true},
		},
	}
}

func getClusterScope(awsCluster infrav1.AWSCluster) (*scope.ClusterScope, error) {
	return scope.NewClusterScope(
		scope.ClusterScopeParams{
			Client: testEnv.Client,
			Cluster: &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
			},
			AWSCluster: &awsCluster,
		},
	)
}

func mockedCreateLBCalls(t *testing.T, m *mocks.MockELBAPIMockRecorder) {
	t.Helper()
	m.DescribeLoadBalancers(gomock.Eq(describeLBInput)).
		Return(describeLBOutput, nil).MinTimes(1)
	m.DescribeLoadBalancerAttributes(gomock.Eq(describeLBAttributesInput)).
		Return(describeLBAttributesOutput, nil)
	m.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: aws.StringSlice([]string{*lbName})}).Return(
		&elb.DescribeTagsOutput{
			TagDescriptions: []*elb.TagDescription{
				{
					LoadBalancerName: lbName,
					Tags: []*elb.Tag{{
						Key:   aws.String(infrav1.ClusterTagKey("test-cluster-apiserver")),
						Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
					}},
				},
			},
		}, nil)
	m.ModifyLoadBalancerAttributes(gomock.Eq(&elb.ModifyLoadBalancerAttributesInput{
		LoadBalancerAttributes: &elb.LoadBalancerAttributes{
			ConnectionSettings:     &elb.ConnectionSettings{IdleTimeout: aws.Int64(600)},
			CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{Enabled: aws.Bool(false)},
		},
		LoadBalancerName: aws.String(""),
	})).MaxTimes(1)

	m.AddTags(gomock.AssignableToTypeOf(&elb.AddTagsInput{})).Return(&elb.AddTagsOutput{}, nil).Do(
		func(actual *elb.AddTagsInput) {
			sortTagsByKey := func(tags []*elb.Tag) {
				sort.Slice(tags, func(i, j int) bool {
					return *(tags[i].Key) < *(tags[j].Key)
				})
			}

			sortTagsByKey(actual.Tags)
			if !cmp.Equal(expectedTags, actual.Tags) {
				t.Fatalf("Actual AddTagsInput did not match expected, Actual : %v, Expected: %v", actual.Tags, expectedTags)
			}
		}).AnyTimes()
	m.RemoveTags(gomock.Eq(&elb.RemoveTagsInput{
		LoadBalancerNames: aws.StringSlice([]string{""}),
		Tags: []*elb.TagKeyOnly{
			{
				Key: aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster-apiserver"),
			},
		},
	})).MaxTimes(1)
	m.ApplySecurityGroupsToLoadBalancer(gomock.Eq(&elb.ApplySecurityGroupsToLoadBalancerInput{
		LoadBalancerName: aws.String(""),
		SecurityGroups:   aws.StringSlice([]string{"sg-apiserver-lb"}),
	})).MaxTimes(1)
	m.RegisterInstancesWithLoadBalancer(gomock.Eq(&elb.RegisterInstancesWithLoadBalancerInput{Instances: []*elb.Instance{{InstanceId: aws.String("two")}}, LoadBalancerName: lbName})).MaxTimes(1)
}

func mockedCreateLBV2Calls(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.DescribeLoadBalancers(gomock.Eq(&elbv2.DescribeLoadBalancersInput{
		Names: []*string{lbName},
	})).
		Return(describeLBOutputV2, nil).MinTimes(1)
	m.DescribeLoadBalancerAttributes(gomock.Eq(&elbv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: lbArn,
	})).Return(describeLBAttributesOutputV2, nil)
	m.DescribeTags(&elbv2.DescribeTagsInput{ResourceArns: []*string{lbArn}}).Return(
		&elbv2.DescribeTagsOutput{
			TagDescriptions: []*elbv2.TagDescription{
				{
					ResourceArn: lbArn,
					Tags: []*elbv2.Tag{{
						Key:   aws.String(infrav1.ClusterTagKey("test-cluster-apiserver")),
						Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
					}},
				},
			},
		}, nil)
	m.ModifyLoadBalancerAttributes(gomock.Eq(&elbv2.ModifyLoadBalancerAttributesInput{
		Attributes: []*elbv2.LoadBalancerAttribute{
			{
				Key:   aws.String(infrav1.LoadBalancerAttributeEnableLoadBalancingCrossZone),
				Value: aws.String("false"),
			},
		},
		LoadBalancerArn: lbArn,
	})).MaxTimes(1)
	m.AddTags(gomock.AssignableToTypeOf(&elbv2.AddTagsInput{})).Return(&elbv2.AddTagsOutput{}, nil).Do(
		func(actual *elbv2.AddTagsInput) {
			sortTagsByKey := func(tags []*elbv2.Tag) {
				sort.Slice(tags, func(i, j int) bool {
					return *(tags[i].Key) < *(tags[j].Key)
				})
			}

			sortTagsByKey(actual.Tags)
			if !cmp.Equal(expectedV2Tags, actual.Tags) {
				t.Fatalf("Actual AddTagsInput did not match expected, Actual : %v, Expected: %v", actual.Tags, expectedV2Tags)
			}
		}).AnyTimes()
	m.RemoveTags(gomock.Eq(&elbv2.RemoveTagsInput{
		ResourceArns: []*string{lbArn},
		TagKeys:      []*string{aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster-apiserver")},
	})).MaxTimes(1)
}

func mockedDeleteLBCalls(expectV2Call bool, mv2 *mocks.MockELBV2APIMockRecorder, m *mocks.MockELBAPIMockRecorder) {
	if expectV2Call {
		mv2.DescribeLoadBalancers(gomock.Any()).Return(describeLBOutputV2, nil)
		mv2.DescribeLoadBalancerAttributes(gomock.Any()).
			Return(describeLBAttributesOutputV2, nil).MaxTimes(1)
		mv2.DescribeTags(gomock.Any()).Return(
			&elbv2.DescribeTagsOutput{
				TagDescriptions: []*elbv2.TagDescription{
					{
						Tags: []*elbv2.Tag{
							{
								Key:   aws.String("name"),
								Value: lbName,
							},
						},
					},
				},
			}, nil).MaxTimes(1)
		mv2.DescribeTargetGroups(gomock.Any()).Return(&elbv2.DescribeTargetGroupsOutput{}, nil)
		mv2.DescribeListeners(gomock.Any()).Return(&elbv2.DescribeListenersOutput{}, nil)
		mv2.DeleteLoadBalancer(gomock.Eq(&elbv2.DeleteLoadBalancerInput{LoadBalancerArn: lbArn})).
			Return(&elbv2.DeleteLoadBalancerOutput{}, nil).MaxTimes(1)
		mv2.DescribeLoadBalancers(gomock.Any()).Return(&elbv2.DescribeLoadBalancersOutput{}, nil)
	}
	m.DescribeLoadBalancers(gomock.Eq(describeLBInput)).
		Return(describeLBOutput, nil)
	m.DescribeLoadBalancers(gomock.Eq(describeLBInput)).
		Return(&elb.DescribeLoadBalancersOutput{}, nil).AnyTimes()
	m.DescribeTags(&elb.DescribeTagsInput{LoadBalancerNames: aws.StringSlice([]string{*lbName})}).Return(
		&elb.DescribeTagsOutput{
			TagDescriptions: []*elb.TagDescription{
				{
					LoadBalancerName: lbName,
				},
			},
		}, nil).MaxTimes(1)
	m.DescribeLoadBalancerAttributes(gomock.Eq(describeLBAttributesInput)).
		Return(describeLBAttributesOutput, nil).MaxTimes(1)
	m.DeleteLoadBalancer(gomock.Eq(&elb.DeleteLoadBalancerInput{LoadBalancerName: lbName})).
		Return(&elb.DeleteLoadBalancerOutput{}, nil).MaxTimes(1)
	m.DescribeLoadBalancersPages(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
}
