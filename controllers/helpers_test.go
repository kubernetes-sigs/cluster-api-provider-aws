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
	"context"
	"sort"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	elbService "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/helpers"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

const DNSName = "www.google.com"

var (
	lbName          = aws.String("test-cluster-apiserver")
	lbArn           = aws.String("loadbalancer::arn")
	tgArn           = aws.String("arn::target-group")
	describeLBInput = &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []string{"test-cluster-apiserver"},
	}
	describeLBAttributesInput = &elb.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: lbName,
	}
	describeLBOutput = &elb.DescribeLoadBalancersOutput{
		LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{
			{
				Scheme:            aws.String(string(infrav1.ELBSchemeInternetFacing)),
				Subnets:           []string{"subnet-1"},
				AvailabilityZones: []string{"us-east-1a"},
				VPCId:             aws.String("vpc-exists"),
			},
		},
	}
	describeLBOutputV2 = &elbv2.DescribeLoadBalancersOutput{
		LoadBalancers: []elbv2types.LoadBalancer{
			{
				Scheme: elbService.SchemeToSDKScheme(infrav1.ELBSchemeInternetFacing),
				AvailabilityZones: []elbv2types.AvailabilityZone{
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
		Attributes: []elbv2types.LoadBalancerAttribute{
			{
				Key:   aws.String("cross-zone"),
				Value: aws.String("true"),
			},
		},
	}
	describeLBAttributesOutput = &elb.DescribeLoadBalancerAttributesOutput{
		LoadBalancerAttributes: &elbtypes.LoadBalancerAttributes{
			CrossZoneLoadBalancing: &elbtypes.CrossZoneLoadBalancing{
				Enabled: false,
			},
		},
	}
	expectedTags = []elbtypes.Tag{
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
	expectedV2Tags = []elbv2types.Tag{
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
	maxWaitActiveUpdateDelete = time.Minute
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
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSCluster",
			APIVersion: infrav1.GroupVersion.Identifier(),
		},
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
			AWSCluster:                   &awsCluster,
			TagUnmanagedNetworkResources: true,
			MaxWaitActiveUpdateDelete:    maxWaitActiveUpdateDelete,
		},
	)
}

func mockedCreateLBCalls(t *testing.T, m *mocks.MockELBAPIMockRecorder, withHealthCheckUpdate bool) {
	t.Helper()
	m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(describeLBInput)).
		Return(describeLBOutput, nil).MinTimes(1)
	m.DescribeLoadBalancerAttributes(gomock.Any(), gomock.Eq(describeLBAttributesInput)).
		Return(describeLBAttributesOutput, nil)
	m.DescribeTags(gomock.Any(), &elb.DescribeTagsInput{LoadBalancerNames: []string{*lbName}}).Return(
		&elb.DescribeTagsOutput{
			TagDescriptions: []elbtypes.TagDescription{
				{
					LoadBalancerName: lbName,
					Tags: []elbtypes.Tag{{
						Key:   aws.String(infrav1.ClusterTagKey("test-cluster-apiserver")),
						Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
					}},
				},
			},
		}, nil)
	m.ModifyLoadBalancerAttributes(gomock.Any(), gomock.Eq(&elb.ModifyLoadBalancerAttributesInput{
		LoadBalancerAttributes: &elbtypes.LoadBalancerAttributes{
			ConnectionSettings:     &elbtypes.ConnectionSettings{IdleTimeout: aws.Int32(600)},
			CrossZoneLoadBalancing: &elbtypes.CrossZoneLoadBalancing{Enabled: false},
		},
		LoadBalancerName: aws.String(""),
	})).MaxTimes(1)

	if withHealthCheckUpdate {
		m.ConfigureHealthCheck(gomock.Any(), gomock.Any()).MaxTimes(1)
	}

	m.AddTags(gomock.Any(), gomock.AssignableToTypeOf(&elb.AddTagsInput{})).Return(&elb.AddTagsOutput{}, nil).Do(
		func(ctx context.Context, actual *elb.AddTagsInput, optFns ...func(*elb.Options)) {
			sortTagsByKey := func(tags []elbtypes.Tag) {
				sort.Slice(tags, func(i, j int) bool {
					return *(tags[i].Key) < *(tags[j].Key)
				})
			}

			sortTagsByKey(actual.Tags)
			if !cmp.Equal(expectedTags, actual.Tags, cmp.AllowUnexported(elbtypes.Tag{})) {
				t.Fatalf("Actual AddTagsInput did not match expected, Actual : %v, Expected: %v", actual.Tags, expectedTags)
			}
		}).AnyTimes()
	m.RemoveTags(gomock.Any(), gomock.Eq(&elb.RemoveTagsInput{
		LoadBalancerNames: []string{""},
		Tags: []elbtypes.TagKeyOnly{
			{
				Key: aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster-apiserver"),
			},
		},
	})).MaxTimes(1)
	m.ApplySecurityGroupsToLoadBalancer(gomock.Any(), gomock.Eq(&elb.ApplySecurityGroupsToLoadBalancerInput{
		LoadBalancerName: aws.String(""),
		SecurityGroups:   []string{"sg-apiserver-lb"},
	})).MaxTimes(1)
	m.RegisterInstancesWithLoadBalancer(gomock.Any(), gomock.Eq(&elb.RegisterInstancesWithLoadBalancerInput{Instances: []elbtypes.Instance{{InstanceId: aws.String("two")}}, LoadBalancerName: lbName})).MaxTimes(1)
}

func mockedCreateLBV2Calls(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
		Names: []string{aws.ToString(lbName)},
	})).
		Return(describeLBOutputV2, nil).MinTimes(1)
	m.DescribeLoadBalancerAttributes(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: lbArn,
	})).Return(describeLBAttributesOutputV2, nil)
	m.DescribeTags(gomock.Any(), &elbv2.DescribeTagsInput{ResourceArns: []string{aws.ToString(lbArn)}}).Return(
		&elbv2.DescribeTagsOutput{
			TagDescriptions: []elbv2types.TagDescription{
				{
					ResourceArn: lbArn,
					Tags: []elbv2types.Tag{{
						Key:   aws.String(infrav1.ClusterTagKey("test-cluster-apiserver")),
						Value: aws.String(string(infrav1.ResourceLifecycleOwned)),
					}},
				},
			},
		}, nil)
	m.ModifyLoadBalancerAttributes(gomock.Any(), gomock.Eq(&elbv2.ModifyLoadBalancerAttributesInput{
		Attributes: []elbv2types.LoadBalancerAttribute{
			{
				Key:   aws.String(infrav1.LoadBalancerAttributeEnableLoadBalancingCrossZone),
				Value: aws.String("false"),
			},
		},
		LoadBalancerArn: lbArn,
	})).MaxTimes(1)
	m.AddTags(gomock.Any(), gomock.AssignableToTypeOf(&elbv2.AddTagsInput{})).Return(&elbv2.AddTagsOutput{}, nil).Do(
		func(ctx context.Context, actual *elbv2.AddTagsInput, optFns ...func(*elbv2.Options)) {
			sortTagsByKey := func(tags []elbv2types.Tag) {
				sort.Slice(tags, func(i, j int) bool {
					return *(tags[i].Key) < *(tags[j].Key)
				})
			}

			sortTagsByKey(actual.Tags)
			if !cmp.Equal(expectedV2Tags, actual.Tags, cmp.AllowUnexported(elbv2types.Tag{})) {
				t.Fatalf("Actual AddTagsInput did not match expected, Actual : %v, Expected: %v", actual.Tags, expectedV2Tags)
			}
		}).AnyTimes()
	m.RemoveTags(gomock.Any(), gomock.Eq(&elbv2.RemoveTagsInput{
		ResourceArns: []string{aws.ToString(lbArn)},
		TagKeys:      []string{"sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster-apiserver"},
	})).MaxTimes(1)
	m.SetSecurityGroups(gomock.Any(), gomock.Eq(&elbv2.SetSecurityGroupsInput{
		LoadBalancerArn: lbArn,
		SecurityGroups:  []string{"sg-apiserver-lb"},
	})).MaxTimes(1)
	m.WaitUntilLoadBalancerAvailable(gomock.Any(), gomock.Eq(&elbv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []string{aws.ToString(lbArn)},
	}), maxWaitActiveUpdateDelete).MaxTimes(1)
}

func mockedDescribeTargetGroupsCall(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.DescribeTargetGroups(gomock.Any(), gomock.Eq(&elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: lbArn,
	})).
		Return(&elbv2.DescribeTargetGroupsOutput{
			NextMarker: new(string),
			TargetGroups: []elbv2types.TargetGroup{
				{
					HealthCheckEnabled:         aws.Bool(true),
					HealthCheckIntervalSeconds: new(int32),
					HealthCheckPath:            new(string),
					HealthCheckPort:            new(string),
					HealthCheckProtocol:        elbv2types.ProtocolEnumTcp,
					HealthCheckTimeoutSeconds:  new(int32),
					HealthyThresholdCount:      new(int32),
					IpAddressType:              elbv2types.TargetGroupIpAddressTypeEnumIpv4,
					LoadBalancerArns:           []string{aws.ToString(lbArn)},
					Matcher:                    &elbv2types.Matcher{},
					Port:                       new(int32),
					Protocol:                   elbv2types.ProtocolEnumTcp,
					ProtocolVersion:            new(string),
					TargetGroupArn:             tgArn,
					TargetGroupName:            new(string),
					TargetType:                 elbv2types.TargetTypeEnumIp,
					UnhealthyThresholdCount:    new(int32),
					VpcId:                      new(string),
				},
			},
		}, nil)
}

func mockedCreateTargetGroupCall(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.CreateTargetGroup(gomock.Any(), helpers.PartialMatchCreateTargetGroupInput(t, &elbv2.CreateTargetGroupInput{
		HealthCheckEnabled:         aws.Bool(true),
		HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
		HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
		HealthCheckProtocol:        elbv2types.ProtocolEnumTcp,
		HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
		HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
		// Note: this is treated as a prefix with the partial matcher.
		Name:     aws.String("apiserver-target"),
		Port:     aws.Int32(infrav1.DefaultAPIServerPort),
		Protocol: elbv2types.ProtocolEnumTcp,
		Tags: []elbv2types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("bar-apiserver"),
			},
			{
				Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
				Value: aws.String("owned"),
			},
			{
				Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
				Value: aws.String("apiserver"),
			},
		},
		UnhealthyThresholdCount: aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
		VpcId:                   aws.String("vpc-exists"),
	})).Return(&elbv2.CreateTargetGroupOutput{
		TargetGroups: []elbv2types.TargetGroup{{
			HealthCheckEnabled:         aws.Bool(true),
			HealthCheckIntervalSeconds: aws.Int32(infrav1.DefaultAPIServerHealthCheckIntervalSec),
			HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
			HealthCheckProtocol:        elbv2types.ProtocolEnumTcp,
			HealthCheckTimeoutSeconds:  aws.Int32(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
			HealthyThresholdCount:      aws.Int32(infrav1.DefaultAPIServerHealthThresholdCount),
			LoadBalancerArns:           []string{aws.ToString(lbArn)},
			Matcher:                    &elbv2types.Matcher{},
			Port:                       aws.Int32(infrav1.DefaultAPIServerPort),
			Protocol:                   elbv2types.ProtocolEnumTcp,
			TargetGroupArn:             tgArn,
			TargetGroupName:            aws.String("apiserver-target"),
			UnhealthyThresholdCount:    aws.Int32(infrav1.DefaultAPIServerUnhealthThresholdCount),
			VpcId:                      aws.String("vpc-exists"),
		}},
	}, nil)
}

func mockedModifyTargetGroupAttributes(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.ModifyTargetGroupAttributes(gomock.Any(), gomock.Eq(&elbv2.ModifyTargetGroupAttributesInput{
		TargetGroupArn: tgArn,
		Attributes: []elbv2types.TargetGroupAttribute{
			{
				Key:   aws.String(infrav1.TargetGroupAttributeEnableConnectionTermination),
				Value: aws.String("false"),
			},
			{
				Key:   aws.String(infrav1.TargetGroupAttributeUnhealthyDrainingIntervalSeconds),
				Value: aws.String("300"),
			},
			{
				Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
				Value: aws.String("false"),
			},
		},
	})).Return(nil, nil)
}

func mockedDescribeListenersCall(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.DescribeListeners(gomock.Any(), gomock.Eq(&elbv2.DescribeListenersInput{
		LoadBalancerArn: lbArn,
	})).
		Return(&elbv2.DescribeListenersOutput{
			Listeners: []elbv2types.Listener{{
				DefaultActions: []elbv2types.Action{{
					TargetGroupArn: aws.String("arn::targetgroup-not-found"),
				}},
				ListenerArn:     aws.String("arn::listener"),
				LoadBalancerArn: lbArn,
			}},
		}, nil)
}

func mockedCreateListenerCall(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.CreateListener(gomock.Any(), gomock.Eq(&elbv2.CreateListenerInput{
		DefaultActions: []elbv2types.Action{
			{
				TargetGroupArn: tgArn,
				Type:           elbv2types.ActionTypeEnumForward,
			},
		},
		LoadBalancerArn: lbArn,
		Port:            aws.Int32(infrav1.DefaultAPIServerPort),
		Protocol:        elbv2types.ProtocolEnumTcp,
		Tags: []elbv2types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("test-cluster-apiserver"),
			},
			{
				Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/cluster/test-cluster"),
				Value: aws.String("owned"),
			},
			{
				Key:   aws.String("sigs.k8s.io/cluster-api-provider-aws/role"),
				Value: aws.String("apiserver"),
			},
		},
	})).Return(&elbv2.CreateListenerOutput{
		Listeners: []elbv2types.Listener{
			{
				DefaultActions: []elbv2types.Action{
					{
						TargetGroupArn: tgArn,
						Type:           elbv2types.ActionTypeEnumForward,
					},
				},
				ListenerArn: aws.String("listener::arn"),
				Port:        aws.Int32(infrav1.DefaultAPIServerPort),
				Protocol:    elbv2types.ProtocolEnumTcp,
			},
		},
	}, nil)
}

func mockedDeleteLBCalls(expectV2Call bool, mv2 *mocks.MockELBV2APIMockRecorder, m *mocks.MockELBAPIMockRecorder) {
	if expectV2Call {
		mv2.DescribeLoadBalancers(gomock.Any(), gomock.Any()).Return(describeLBOutputV2, nil)
		mv2.DescribeLoadBalancerAttributes(gomock.Any(), gomock.Any()).
			Return(describeLBAttributesOutputV2, nil).MaxTimes(1)
		mv2.DescribeTags(gomock.Any(), gomock.Any()).Return(
			&elbv2.DescribeTagsOutput{
				TagDescriptions: []elbv2types.TagDescription{
					{
						Tags: []elbv2types.Tag{
							{
								Key:   aws.String("name"),
								Value: lbName,
							},
						},
					},
				},
			}, nil).MaxTimes(1)
		mv2.DescribeTargetGroups(gomock.Any(), gomock.Any()).Return(&elbv2.DescribeTargetGroupsOutput{}, nil)
		mv2.DescribeListeners(gomock.Any(), gomock.Any()).Return(&elbv2.DescribeListenersOutput{}, nil)
		mv2.DeleteLoadBalancer(gomock.Any(), gomock.Eq(&elbv2.DeleteLoadBalancerInput{LoadBalancerArn: lbArn})).
			Return(&elbv2.DeleteLoadBalancerOutput{}, nil).MaxTimes(1)
		mv2.DescribeLoadBalancers(gomock.Any(), gomock.Any()).Return(&elbv2.DescribeLoadBalancersOutput{}, nil)
	}
	m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(describeLBInput)).
		Return(describeLBOutput, nil)
	m.DescribeLoadBalancers(gomock.Any(), gomock.Eq(describeLBInput)).
		Return(&elb.DescribeLoadBalancersOutput{}, nil).AnyTimes()
	m.DescribeTags(gomock.Any(), &elb.DescribeTagsInput{LoadBalancerNames: []string{*lbName}}).Return(
		&elb.DescribeTagsOutput{
			TagDescriptions: []elbtypes.TagDescription{
				{
					LoadBalancerName: lbName,
				},
			},
		}, nil).MaxTimes(1)
	m.DescribeLoadBalancerAttributes(gomock.Any(), gomock.Eq(describeLBAttributesInput)).
		Return(describeLBAttributesOutput, nil).MaxTimes(1)
	m.DeleteLoadBalancer(gomock.Any(), gomock.Eq(&elb.DeleteLoadBalancerInput{LoadBalancerName: lbName})).
		Return(&elb.DeleteLoadBalancerOutput{}, nil).MaxTimes(1)
	m.DescribeLoadBalancersPages(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
}
