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

package client

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/version"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	machineapiapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"
	apimachineryerrors "k8s.io/apimachinery/pkg/api/errors"
)

//go:generate go run ../../vendor/github.com/golang/mock/mockgen -source=./client.go -destination=./mock/client_generated.go -package=mock

const (
	// AwsCredsSecretIDKey is secret key containing AWS KeyId
	AwsCredsSecretIDKey = "aws_access_key_id"
	// AwsCredsSecretAccessKey is secret key containing AWS Secret Key
	AwsCredsSecretAccessKey = "aws_secret_access_key"
)

// AwsClientBuilderFuncType is function type for building aws client
type AwsClientBuilderFuncType func(client client.Client, secretName, namespace, region string) (Client, error)

// Client is a wrapper object for actual AWS SDK clients to allow for easier testing.
type Client interface {
	DescribeImages(*ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error)
	DescribeDHCPOptions(input *ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error)
	DescribeVpcs(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
	DescribeSubnets(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	DescribeAvailabilityZones(*ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error)
	DescribeSecurityGroups(*ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error)
	RunInstances(*ec2.RunInstancesInput) (*ec2.Reservation, error)
	DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
	TerminateInstances(*ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error)
	DescribeVolumes(*ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error)

	RegisterInstancesWithLoadBalancer(*elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error)
	ELBv2DescribeLoadBalancers(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error)
	ELBv2DescribeTargetGroups(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error)
	ELBv2RegisterTargets(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error)
}

type awsClient struct {
	ec2Client   ec2iface.EC2API
	elbClient   elbiface.ELBAPI
	elbv2Client elbv2iface.ELBV2API
}

func (c *awsClient) DescribeDHCPOptions(input *ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
	return c.ec2Client.DescribeDhcpOptions(input)
}

func (c *awsClient) DescribeImages(input *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	return c.ec2Client.DescribeImages(input)
}

func (c *awsClient) DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return c.ec2Client.DescribeVpcs(input)
}

func (c *awsClient) DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return c.ec2Client.DescribeSubnets(input)
}

func (c *awsClient) DescribeAvailabilityZones(input *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
	return c.ec2Client.DescribeAvailabilityZones(input)
}

func (c *awsClient) DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return c.ec2Client.DescribeSecurityGroups(input)
}

func (c *awsClient) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	return c.ec2Client.RunInstances(input)
}

func (c *awsClient) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return c.ec2Client.DescribeInstances(input)
}

func (c *awsClient) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return c.ec2Client.TerminateInstances(input)
}

func (c *awsClient) DescribeVolumes(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	return c.ec2Client.DescribeVolumes(input)
}

func (c *awsClient) RegisterInstancesWithLoadBalancer(input *elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {
	return c.elbClient.RegisterInstancesWithLoadBalancer(input)
}

func (c *awsClient) ELBv2DescribeLoadBalancers(input *elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
	return c.elbv2Client.DescribeLoadBalancers(input)
}

func (c *awsClient) ELBv2DescribeTargetGroups(input *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
	return c.elbv2Client.DescribeTargetGroups(input)
}

func (c *awsClient) ELBv2RegisterTargets(input *elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
	return c.elbv2Client.RegisterTargets(input)
}

// NewClient creates our client wrapper object for the actual AWS clients we use.
// For authentication the underlying clients will use either the cluster AWS credentials
// secret if defined (i.e. in the root cluster),
// otherwise the IAM profile of the master where the actuator will run. (target clusters)
func NewClient(ctrlRuntimeClient client.Client, secretName, namespace, region string) (Client, error) {
	awsConfig := &aws.Config{Region: aws.String(region)}

	if secretName != "" {
		var secret corev1.Secret
		if err := ctrlRuntimeClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: secretName}, &secret); err != nil {
			if apimachineryerrors.IsNotFound(err) {
				return nil, machineapiapierrors.InvalidMachineConfiguration("aws credentials secret %s/%s: %v not found", namespace, secretName, err)
			}
			return nil, err
		}
		accessKeyID, ok := secret.Data[AwsCredsSecretIDKey]
		if !ok {
			return nil, machineapiapierrors.InvalidMachineConfiguration("AWS credentials secret %v did not contain key %v",
				secretName, AwsCredsSecretIDKey)
		}
		secretAccessKey, ok := secret.Data[AwsCredsSecretAccessKey]
		if !ok {
			return nil, machineapiapierrors.InvalidMachineConfiguration("AWS credentials secret %v did not contain key %v",
				secretName, AwsCredsSecretAccessKey)
		}

		awsConfig.Credentials = credentials.NewStaticCredentials(
			string(accessKeyID), string(secretAccessKey), "")
	}

	// Otherwise default to relying on the IAM role of the masters where the actuator is running:
	s, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	s.Handlers.Build.PushBackNamed(addProviderVersionToUserAgent)

	return &awsClient{
		ec2Client:   ec2.New(s),
		elbClient:   elb.New(s),
		elbv2Client: elbv2.New(s),
	}, nil
}

// NewClientFromKeys creates our client wrapper object for the actual AWS clients we use.
// For authentication the underlying clients will use AWS credentials.
func NewClientFromKeys(accessKey, secretAccessKey, region string) (Client, error) {
	awsConfig := &aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretAccessKey,
			"",
		),
	}

	s, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	s.Handlers.Build.PushBackNamed(addProviderVersionToUserAgent)

	return &awsClient{
		ec2Client:   ec2.New(s),
		elbClient:   elb.New(s),
		elbv2Client: elbv2.New(s),
	}, nil
}

// addProviderVersionToUserAgent is a named handler that will add cluster-api-provider-aws
// version information to requests made by the AWS SDK.
var addProviderVersionToUserAgent = request.NamedHandler{
	Name: "openshift.io/cluster-api-provider-aws",
	Fn:   request.MakeAddToUserAgentHandler("openshift.io cluster-api-provider-aws", version.Version.String()),
}
