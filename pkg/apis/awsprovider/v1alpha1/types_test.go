/*
Copyright 2018 The Kubernetes Authors.

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
package v1alpha1_test

import (
	"reflect"
	"testing"

	providerv1 "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

func TestEquals(t *testing.T) {
	testcases := []struct {
		name                  string
		inputIngressRule      *providerv1.IngressRule
		inputIngressRuleParam *providerv1.IngressRule
		expected              bool
	}{
		{
			name: "Equals should return false for different Description",
			inputIngressRule: &providerv1.IngressRule{
				Description:            "allow access from vpc production",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               80,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			inputIngressRuleParam: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               80,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			expected: false,
		},
		{
			name: "Equals should return false for different Protocols",
			inputIngressRule: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			inputIngressRuleParam: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolUDP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			expected: false,
		},
		{
			name: "Equals should return false for different FromPort",
			inputIngressRule: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               80,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			inputIngressRuleParam: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               443,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			expected: false,
		},
		{
			name: "Equals should return false for different ToPort",
			inputIngressRule: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			inputIngressRuleParam: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 443,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			expected: false,
		},
		{
			name: "Equals should return false for different CidrBlocks",
			inputIngressRule: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "192.168.1.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			inputIngressRuleParam: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"172.16.12.0/24", "10.20.30.0/24", "10.30.40.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			expected: false,
		},
		{
			name: "Equals should return false for CidrBlocks of different lengths",
			inputIngressRule: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			inputIngressRuleParam: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"172.16.12.0/24", "10.20.30.0/24", "10.30.40.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			expected: false,
		},
		{
			name: "Equals should return false for different SourceSecurityGroupIDs",
			inputIngressRule: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-3031748"},
			},
			inputIngressRuleParam: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"172.16.12.0/24", "10.20.30.0/24", "10.30.40.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1236745", "sg-1234567"},
			},
			expected: false,
		},
		{
			name: "Equals should return false for SourceSecurityGroupIDs of different lengths",
			inputIngressRule: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567"},
			},
			inputIngressRuleParam: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"172.16.12.0/24", "10.20.30.0/24", "10.30.40.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1236745", "sg-1234567"},
			},
			expected: false,
		},
		{
			name: "Equals should return true for same SourceSecurityGroupIDs and CidrBlocks",
			inputIngressRule: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"10.20.30.0/24", "10.30.40.0/24", "172.16.12.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1234567", "sg-1236745"},
			},
			inputIngressRuleParam: &providerv1.IngressRule{
				Description:            "allow access from vpc",
				Protocol:               providerv1.SecurityGroupProtocolTCP,
				FromPort:               0,
				ToPort:                 80,
				CidrBlocks:             []string{"172.16.12.0/24", "10.20.30.0/24", "10.30.40.0/24"},
				SourceSecurityGroupIDs: []string{"sg-1236745", "sg-1234567"},
			},
			expected: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.inputIngressRule.Equals(tc.inputIngressRuleParam)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("For [%s], equality check failed, Got: %t, want: %t", tc.name, actual, tc.expected)
			}
		})
	}
}
