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
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/helpers"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/mocks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

const DNSName = "www.google.com"

var (
	lbName          = aws.String("test-cluster-apiserver")
	lbArn           = aws.String("loadbalancer::arn")
	tgArn           = aws.String("arn::target-group")
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
	fakeKubeconfig = `{
  "apiVersion": "v1",
  "clusters": [
    {
      "cluster": {
        "certificate-authority-data": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvakNDQWVhZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJek1ESXlNVEU0TURZeE5sb1hEVE16TURJeE9ERTRNRFl4Tmxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBSnNqCk0yWHJJSXFUaE9ibDA0UG04TmZhRzBmd3NpOHZVUEZhaTJydGNFUCtpVVZPWDVKSXBxZzBjSklCQU4rUU9RZjAKbWR4bjhNNStzWXVuN3BDTThSWnp4Rzd1OWk4U3FtemhlSHZVU2t0QjNWSDNyY0lqUGlpdW5XQWZWR2VJaC9aZQpvck10RE9oVndwbXFiU2k5NlJaSFhHanV4aVBreXd4UUowTFlNNG02d0ZhNTFFMzJvSU5ZUXJ3ZnBJMWgvZXp4ClZhVkhPTU1aeng0Z1RyNWFWYVFrdk0raFRWYVpUT0p4QWdoek8wNmc1anRhTjlVYy94a3dBRUI3dFczWWtURHUKdjJKUVdiWHhucDRDZ2luWWlTbU51aDBwQlVLQTFhbDA2ejhkbnhnUEFRRkhVRUlqTy9CZ3IyTFB2SmpPNHlDYwpYbjd6ZUF4U1BCREo3ZGlQaTFzQ0F3RUFBYU5aTUZjd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZCaVVMditWa0NtRjd0b2lJSE00Z2tOWWQyZWdNQlVHQTFVZEVRUU8KTUF5Q0NtdDFZbVZ5Ym1WMFpYTXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBQ0J3RXNpR3l1OW96VHJXVUp4bAphMlJmenArQU5MWG1YcWJzb09YLzIvVGVoUlB1Y1VnQml4aXQvdDNud3E2MFNsOWJVaHd4S2VSd29FRjlFY2xGCmNuQUZySytjM1FDSVR6ek15aW9MZEdVYUNGdEVlWnVld3p4T2dpTm42SDR0b1lVR2pDTTlhdDY2OG53RHBQc3oKcWRFSjdndnhNcjB3cUhNT2Q5SE55eUlJUFdjZzJXRThjdzQrQzllTUNVWWI2Y2lHMHl6VmViSXExc2tNT0hFRQpFTGtCNytRd25Gcjd1Y1huZ29BUHpVMzg2VW9pWkpNbHFrN1djUEUzdk1SZmhOOW5hTmExZjk4T3NMRFgvdG1OCjBzaWpFZSt4RCtvR0s4ZzhsNXlLbUlad1BkNElCUWlqSjBGcFZaOFJmNE5hSEd1QTYxb2hzZkU5VTJkVVM1aEMKQ0F3PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
        "server": "https://127.0.0.1:12345"
      },
      "name": "capa"
    }
  ],
  "contexts": [
    {
      "context": {
        "cluster": "capa",
        "user": "capa"
      },
      "name": "capa"
    }
  ],
  "current-context": "capa",
  "kind": "Config",
  "preferences": {},
  "users": [
    {
      "name": "capa",
      "user": {
        "client-certificate-data": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURJVENDQWdtZ0F3SUJBZ0lJWTFTa3dwcTFwQm93RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TXpBeU1qRXhPREEyTVRaYUZ3MHlOREF5TWpFeE9EQTJNVGhhTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXJBT24zTEdFSVE2aFB2b3YKbFlMNTlEWmFVTkhXT1hYdkNJWlIzVmRSSVFVSjJEQTFjaEpyazNEbVNZN016dmdYRjBRNktETFNZOE5veEVlSQpoZ0NOZlZtN1pnaVlUNk9zUXpTNzVIdDdONWpReHp5ZkdsTFZOdjFRenM1M01DTUUrMjI4czRPdzdWUkpySVJNClErRzl0L1dveG8rNjhlRTZOMmlZSzdJUE9peU85VzFVQ3FPc285dVM1dHBzeWtZYkRQYlI4aG1NZGdRWFpLZUQKVzAzRHRuekRKOFVjN3crQXRTT0pPSHRMbmlONG5WelJ6cUdqWHlDZXQ5RGhFWHQyNkFKejJMRzZXbDEwZElYdApBTEh5Y1BkdWtGUXBxT0czQjU0VXVXbWJDcTlGZ0pKQ0cyT1dsL2NpdlhoUWw2UEpRR1prK1RRVkl0eS95OWhNCkgxai9QUUlEQVFBQm8xWXdWREFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0RBWURWUjBUQVFIL0JBSXdBREFmQmdOVkhTTUVHREFXZ0JRWWxDNy9sWkFwaGU3YUlpQnpPSUpEV0hkbgpvREFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBTlNPbHRpT2JFb2lKTFk5c1hvOHZZT05HYUYvdFNiV1BOU1RnCkEvTmJzUFAzSjJnVmpHa1pUL3d6TCtXd3FmL2FXa2xnQjZIL2pUTjBpbnp3UXZFQ2RhQ054U0gyUXZoTUx1Rm4KeVBHdHdOKzFFQ1VoSnMvcUJlM0F2elRkVlorbWw3SDEraW9oQ0k0T2JNakZJNkhiN0tEL1RTbXE5cEhhdGxFbApmeDRPamlSU0krU09uR2QwU01zZUNVS0loaUNzcVdHWU5wcnRuZDNBM0E5bmxiZVBVZ3BPN0tnMDR2SUZnenVuCmRDTzFpaytOcjhJR0JaM1RKTThGc09LZmlVMjNacGRMbU1TWmZXaUU5U1I1WmdUMVcySEpPZS9UWW04TzJIdU8KcTI5a2VrK3VBL3JYMml3MTh6d0JkbGlUZEtTVnN6V2pHYUUxY09laktJZTdlaWtKU1E9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
        "client-key-data": "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBckFPbjNMR0VJUTZoUHZvdmxZTDU5RFphVU5IV09YWHZDSVpSM1ZkUklRVUoyREExCmNoSnJrM0RtU1k3TXp2Z1hGMFE2S0RMU1k4Tm94RWVJaGdDTmZWbTdaZ2lZVDZPc1F6Uzc1SHQ3TjVqUXh6eWYKR2xMVk52MVF6czUzTUNNRSsyMjhzNE93N1ZSSnJJUk1RK0c5dC9Xb3hvKzY4ZUU2TjJpWUs3SVBPaXlPOVcxVQpDcU9zbzl1UzV0cHN5a1liRFBiUjhobU1kZ1FYWktlRFcwM0R0bnpESjhVYzd3K0F0U09KT0h0TG5pTjRuVnpSCnpxR2pYeUNldDlEaEVYdDI2QUp6MkxHNldsMTBkSVh0QUxIeWNQZHVrRlFwcU9HM0I1NFV1V21iQ3E5RmdKSkMKRzJPV2wvY2l2WGhRbDZQSlFHWmsrVFFWSXR5L3k5aE1IMWovUFFJREFRQUJBb0lCQUZDY3Nma3FNK3F1Q1lFVApER083NnRtNDh5QjNCamNOZnNUMjlieU9KQlllWDF1eVZBbXVlajJ4cGlxSXdwQ29FV0kwOWpCS2lQQjU3a28xCjM5UTB4Qm9maUVRcitQMHhqaFE4eldVcFBUaHo1RHZLdkNGQSttejN6L3ZySkU0cHl6YnRXWWFqUzdJZFV4MkgKTHBxTld6U3Y2cldMUENnSXpoaVRhdVRuWUYxNUtaZFhNY2FHRzg3Sk5OTW5hQnNWY2pMOWhvS3JhZjltK3UwYgpPL2VBZmJXY3F1b1daUzVuQll1THdjWVBxcDVaQkhWSFB0OXFjczZnSGMwYitTRlNrL1ZmdGkzemw2OTVNMWNvCkR6NmxnamlXRFEzOUNPVjVDRTZYSGxoMnlIU3UxeFp2ZmZBV2dDSGxMZmtEaGJvYm44VVNwVlZPcWxiZDg5M1kKQVF1VWJjRUNnWUVBelljUm9YUUd3enpjZkxPTjNLbk9LVU1xNFpWbFk0aElTdjdma3lBaHZxZVRPdkNsdENJZwpXOGhjRGFWSjZOZWhhdzRBUFIramFEQ2VXVnI2QmpSNHpiVWdoS1Iyc0hDZ2ZTWGxiWHRxRG1JRmRBb3d5Q0NOCmkvMDM3R2laRlZ2TDFnSklBa3hOS1ZjRUNKQkJBeFdHWFJ4bWgwQVpGTS9qZlFRZ29zU1JHTkVDZ1lFQTFrR3oKSm1BWitxZzZWQllqOTk4Uzg5T21zNFdKNXl4c1RYV29oYldGWVVsbjM0ZG5zdGdWd3dYYmdGUGJyNGF3TEdYYwo3T295RU80RkVudUY3Umt0TVhydlg5L0p6SjE5Qnl5b1ZIUjJYb1RSMzlGTDExYW13V1VLaVlJS2xncVU4L0FZCkwwckErTjJWTEhJaXpOMklaL204ZzVnMjg4WndjTE5kSm5IbUdxMENnWUJRUFRWQzVUdG1xYklpOVM2alFaLzkKTi8zYnlDbW5MQ09kTTlneFpsQUdVUUlIOXIrYWplQTROUWJMUlFhUDR1OWdEajFGbjc1NkJORXRiWGxEUnVVSwoyblh4a3d2TFlvMGxqcy85YUR0Rmxqc3V6SE96RGhKMDNzMGdmTTJYS3hsdldjQ25OUWJDNXZmcUovZFpydU9nCkltQmMyVWR4ZXFBRHhNTEJDU0Rrb1FLQmdGNnVQOTNJZ0JKOWZ2RWpxNWRnMDc0K0hKK2VkbmRhOFMwMXZsZ0EKQkVZZXF6RmpZOWJybUlwTEwxbkJOUWFYRFlsQkptVG5oV0puM0lQelpCYUhscW9UK200eXRibWZLdDRkeFBFMQpXZzJnd2lJWEdsMjVwQTA0ZW5TVHE5dnNKekM1TytiQ01RNkkxT0FFUEE2dUl4WlhqUS9XRndxWStaMUVGZmprClFsd3RBb0dBU3FnNGNEd2Q4d2QwOFd2aEtaYkNoeXh4cVYvczdkT2YrZWhFaWtqdGlMc0pMc3IxdUEzZ0oxQS8KY1p5QWZVaURVZlE1cHVUQWVzUXo3b0VsYU9YQjlTRTdCQlhQZmswdnJBOVRyRmJ1UytwVjlkUjc3WkhMMGlnUgoxQXg1UE1YLzdGY1ZPRWZselR5V2U2TFkvYVU5MFJkVEtGQ1NKa0VzaTZPNzl4eUVnbDQ9Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg=="
      }
    }
  ]
}`
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
	m.SetSecurityGroups(gomock.Eq(&elbv2.SetSecurityGroupsInput{
		LoadBalancerArn: lbArn,
		SecurityGroups:  aws.StringSlice([]string{"sg-apiserver-lb"}),
	})).MaxTimes(1)
}

func mockedDescribeTargetGroupsCall(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.DescribeTargetGroups(gomock.Eq(&elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: lbArn,
	})).
		Return(&elbv2.DescribeTargetGroupsOutput{
			NextMarker: new(string),
			TargetGroups: []*elbv2.TargetGroup{
				{
					HealthCheckEnabled:         aws.Bool(true),
					HealthCheckIntervalSeconds: new(int64),
					HealthCheckPath:            new(string),
					HealthCheckPort:            new(string),
					HealthCheckProtocol:        new(string),
					HealthCheckTimeoutSeconds:  new(int64),
					HealthyThresholdCount:      new(int64),
					IpAddressType:              new(string),
					LoadBalancerArns:           []*string{lbArn},
					Matcher:                    &elbv2.Matcher{},
					Port:                       new(int64),
					Protocol:                   new(string),
					ProtocolVersion:            new(string),
					TargetGroupArn:             tgArn,
					TargetGroupName:            new(string),
					TargetType:                 new(string),
					UnhealthyThresholdCount:    new(int64),
					VpcId:                      new(string),
				}},
		}, nil)
}

func mockedCreateTargetGroupCall(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.CreateTargetGroup(helpers.PartialMatchCreateTargetGroupInput(t, &elbv2.CreateTargetGroupInput{
		HealthCheckEnabled:         aws.Bool(true),
		HealthCheckIntervalSeconds: aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
		HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
		HealthCheckProtocol:        aws.String("TCP"),
		HealthCheckTimeoutSeconds:  aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
		HealthyThresholdCount:      aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
		// Note: this is treated as a prefix with the partial matcher.
		Name:     aws.String("apiserver-target"),
		Port:     aws.Int64(infrav1.DefaultAPIServerPort),
		Protocol: aws.String("TCP"),
		Tags: []*elbv2.Tag{
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
		UnhealthyThresholdCount: aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
		VpcId:                   aws.String("vpc-exists"),
	})).Return(&elbv2.CreateTargetGroupOutput{
		TargetGroups: []*elbv2.TargetGroup{{
			HealthCheckEnabled:         aws.Bool(true),
			HealthCheckIntervalSeconds: aws.Int64(infrav1.DefaultAPIServerHealthCheckIntervalSec),
			HealthCheckPort:            aws.String(infrav1.DefaultAPIServerPortString),
			HealthCheckProtocol:        aws.String("TCP"),
			HealthCheckTimeoutSeconds:  aws.Int64(infrav1.DefaultAPIServerHealthCheckTimeoutSec),
			HealthyThresholdCount:      aws.Int64(infrav1.DefaultAPIServerHealthThresholdCount),
			LoadBalancerArns:           []*string{lbArn},
			Matcher:                    &elbv2.Matcher{},
			Port:                       aws.Int64(infrav1.DefaultAPIServerPort),
			Protocol:                   aws.String("TCP"),
			TargetGroupArn:             tgArn,
			TargetGroupName:            aws.String("apiserver-target"),
			UnhealthyThresholdCount:    aws.Int64(infrav1.DefaultAPIServerUnhealthThresholdCount),
			VpcId:                      aws.String("vpc-exists"),
		}},
	}, nil)
}

func mockedModifyTargetGroupAttributes(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.ModifyTargetGroupAttributes(gomock.Eq(&elbv2.ModifyTargetGroupAttributesInput{
		TargetGroupArn: tgArn,
		Attributes: []*elbv2.TargetGroupAttribute{
			{
				Key:   aws.String(infrav1.TargetGroupAttributeEnablePreserveClientIP),
				Value: aws.String("false"),
			},
		},
	})).Return(nil, nil)
}

func mockedDescribeListenersCall(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.DescribeListeners(gomock.Eq(&elbv2.DescribeListenersInput{
		LoadBalancerArn: lbArn,
	})).
		Return(&elbv2.DescribeListenersOutput{
			Listeners: []*elbv2.Listener{{
				DefaultActions: []*elbv2.Action{{
					TargetGroupArn: aws.String("arn::targetgroup-not-found"),
				}},
				ListenerArn:     aws.String("arn::listener"),
				LoadBalancerArn: lbArn,
			}},
		}, nil)
}

func mockedCreateListenerCall(t *testing.T, m *mocks.MockELBV2APIMockRecorder) {
	t.Helper()
	m.CreateListener(gomock.Eq(&elbv2.CreateListenerInput{
		DefaultActions: []*elbv2.Action{
			{
				TargetGroupArn: tgArn,
				Type:           aws.String(elbv2.ActionTypeEnumForward),
			},
		},
		LoadBalancerArn: lbArn,
		Port:            aws.Int64(infrav1.DefaultAPIServerPort),
		Protocol:        aws.String("TCP"),
		Tags: []*elbv2.Tag{
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
		Listeners: []*elbv2.Listener{
			{
				DefaultActions: []*elbv2.Action{
					{
						TargetGroupArn: tgArn,
						Type:           aws.String(elbv2.ActionTypeEnumForward),
					},
				},
				ListenerArn: aws.String("listener::arn"),
				Port:        aws.Int64(infrav1.DefaultAPIServerPort),
				Protocol:    aws.String("TCP"),
			},
		}}, nil)
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
