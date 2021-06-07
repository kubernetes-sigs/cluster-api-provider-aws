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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/version"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	configv1 "github.com/openshift/api/config/v1"
	machineapiapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"
	apimachineryerrors "k8s.io/apimachinery/pkg/api/errors"
)

//go:generate go run ../../vendor/github.com/golang/mock/mockgen -source=./client.go -destination=./mock/client_generated.go -package=mock

const (
	// AwsCredsSecretIDKey is secret key containing AWS KeyId
	AwsCredsSecretIDKey = "aws_access_key_id"
	// AwsCredsSecretAccessKey is secret key containing AWS Secret Key
	AwsCredsSecretAccessKey = "aws_secret_access_key"

	// GlobalInfrastuctureName default name for infrastructure object
	GlobalInfrastuctureName = "cluster"

	// KubeCloudConfigNamespace is the namespace where the kube cloud config ConfigMap is located
	KubeCloudConfigNamespace = "openshift-config-managed"
	// kubeCloudConfigName is the name of the kube cloud config ConfigMap
	kubeCloudConfigName = "kube-cloud-config"
	// cloudCABundleKey is the key in the kube cloud config ConfigMap where the custom CA bundle is located
	cloudCABundleKey = "ca-bundle.pem"
)

// AwsClientBuilderFuncType is function type for building aws client
type AwsClientBuilderFuncType func(client client.Client, secretName, namespace, region string, configManagedClient client.Client) (Client, error)

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
	CreateTags(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error)

	RegisterInstancesWithLoadBalancer(*elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error)
	ELBv2DescribeLoadBalancers(*elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error)
	ELBv2DescribeTargetGroups(*elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error)
	ELBv2DescribeTargetHealth(*elbv2.DescribeTargetHealthInput) (*elbv2.DescribeTargetHealthOutput, error)
	ELBv2RegisterTargets(*elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error)
	ELBv2DeregisterTargets(*elbv2.DeregisterTargetsInput) (*elbv2.DeregisterTargetsOutput, error)
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

func (c *awsClient) CreateTags(input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return c.ec2Client.CreateTags(input)
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

func (c *awsClient) ELBv2DescribeTargetHealth(input *elbv2.DescribeTargetHealthInput) (*elbv2.DescribeTargetHealthOutput, error) {
	return c.elbv2Client.DescribeTargetHealth(input)
}

func (c *awsClient) ELBv2RegisterTargets(input *elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
	return c.elbv2Client.RegisterTargets(input)
}

func (c *awsClient) ELBv2DeregisterTargets(input *elbv2.DeregisterTargetsInput) (*elbv2.DeregisterTargetsOutput, error) {
	return c.elbv2Client.DeregisterTargets(input)
}

// NewClient creates our client wrapper object for the actual AWS clients we use.
// For authentication the underlying clients will use either the cluster AWS credentials
// secret if defined (i.e. in the root cluster),
// otherwise the IAM profile of the master where the actuator will run. (target clusters)
func NewClient(ctrlRuntimeClient client.Client, secretName, namespace, region string, configManagedClient client.Client) (Client, error) {
	s, err := newAWSSession(ctrlRuntimeClient, secretName, namespace, region, configManagedClient)
	if err != nil {
		return nil, err
	}

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

// NewValidatedClient creates our client wrapper object for the actual AWS clients we use.
// This should behave the same as NewClient except it will validate the client configuration
// (eg the region) before returning the client.
func NewValidatedClient(ctrlRuntimeClient client.Client, secretName, namespace, region string, configManagedClient client.Client) (Client, error) {
	s, err := newAWSSession(ctrlRuntimeClient, secretName, namespace, region, configManagedClient)
	if err != nil {
		return nil, err
	}

	// Check that the endpoint can be resolved by the endpoint resolver.
	// If the endpoint is not known, it is not a standard or configured custom region.
	// If this is the case, the client will likely not be able to connect
	_, err = s.Config.EndpointResolver.EndpointFor("ec2", region, func(opts *endpoints.Options) {
		opts.StrictMatching = true
	})
	if err != nil {
		return nil, fmt.Errorf("region %q not resolved: %w", region, err)
	}

	return &awsClient{
		ec2Client:   ec2.New(s),
		elbClient:   elb.New(s),
		elbv2Client: elbv2.New(s),
	}, nil
}

func newAWSSession(ctrlRuntimeClient client.Client, secretName, namespace, region string, configManagedClient client.Client) (*session.Session, error) {
	sessionOptions := session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
	}

	if secretName != "" {
		var secret corev1.Secret
		if err := ctrlRuntimeClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: secretName}, &secret); err != nil {
			if apimachineryerrors.IsNotFound(err) {
				return nil, machineapiapierrors.InvalidMachineConfiguration("aws credentials secret %s/%s: %v not found", namespace, secretName, err)
			}
			return nil, err
		}
		sharedCredsFile, err := sharedCredentialsFileFromSecret(&secret)
		if err != nil {
			return nil, fmt.Errorf("failed to create shared credentials file from Secret: %v", err)
		}
		sessionOptions.SharedConfigState = session.SharedConfigEnable
		sessionOptions.SharedConfigFiles = []string{sharedCredsFile}
	}

	// Resolve custom endpoints
	if err := resolveEndpoints(&sessionOptions.Config, ctrlRuntimeClient, region); err != nil {
		return nil, err
	}

	if err := useCustomCABundle(&sessionOptions, configManagedClient); err != nil {
		return nil, fmt.Errorf("failed to set the custom CA bundle: %w", err)
	}

	// Otherwise default to relying on the IAM role of the masters where the actuator is running:
	s, err := session.NewSessionWithOptions(sessionOptions)
	if err != nil {
		return nil, err
	}

	// Remove any temporary shared credentials files after session creation so they don't accumulate
	if len(sessionOptions.SharedConfigFiles) > 0 {
		os.Remove(sessionOptions.SharedConfigFiles[0])
	}

	s.Handlers.Build.PushBackNamed(addProviderVersionToUserAgent)

	return s, nil
}

// addProviderVersionToUserAgent is a named handler that will add cluster-api-provider-aws
// version information to requests made by the AWS SDK.
var addProviderVersionToUserAgent = request.NamedHandler{
	Name: "openshift.io/cluster-api-provider-aws",
	Fn:   request.MakeAddToUserAgentHandler("openshift.io cluster-api-provider-aws", version.Version.String()),
}

func resolveEndpoints(awsConfig *aws.Config, ctrlRuntimeClient client.Client, region string) error {
	infra := &configv1.Infrastructure{}
	infraName := client.ObjectKey{Name: GlobalInfrastuctureName}

	if err := ctrlRuntimeClient.Get(context.Background(), infraName, infra); err != nil {
		return err
	}

	// Do nothing when custom endpoints are missing
	if infra.Status.PlatformStatus == nil || infra.Status.PlatformStatus.AWS == nil {
		return nil
	}

	customEndpointsMap := buildCustomEndpointsMap(infra.Status.PlatformStatus.AWS.ServiceEndpoints)

	if len(customEndpointsMap) == 0 {
		return nil
	}

	customResolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if url, ok := customEndpointsMap[service]; ok {
			return endpoints.ResolvedEndpoint{
				URL:           url,
				SigningRegion: region,
			}, nil

		}
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	awsConfig.EndpointResolver = endpoints.ResolverFunc(customResolver)

	return nil
}

// buildCustomEndpointsMap constructs a map that links endpoint name and it's url
func buildCustomEndpointsMap(customEndpoints []configv1.AWSServiceEndpoint) map[string]string {
	customEndpointsMap := make(map[string]string)

	for _, customEndpoint := range customEndpoints {
		customEndpointsMap[customEndpoint.Name] = customEndpoint.URL
	}

	return customEndpointsMap
}

// sharedCredentialsFileFromSecret returns a location (path) to the shared credentials
// file that was created using the provided secret
func sharedCredentialsFileFromSecret(secret *corev1.Secret) (string, error) {
	var data []byte
	switch {
	case len(secret.Data["credentials"]) > 0:
		data = secret.Data["credentials"]
	case len(secret.Data["aws_access_key_id"]) > 0 && len(secret.Data["aws_secret_access_key"]) > 0:
		data = newConfigForStaticCreds(
			string(secret.Data["aws_access_key_id"]),
			string(secret.Data["aws_secret_access_key"]),
		)
	default:
		return "", fmt.Errorf("invalid secret for aws credentials")
	}

	f, err := ioutil.TempFile("", "aws-shared-credentials")
	if err != nil {
		return "", fmt.Errorf("failed to create file for shared credentials: %v", err)
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return "", fmt.Errorf("failed to write credentials to %s: %v", f.Name(), err)
	}
	return f.Name(), nil
}

func newConfigForStaticCreds(accessKey string, accessSecret string) []byte {
	buf := &bytes.Buffer{}
	fmt.Fprint(buf, "[default]\n")
	fmt.Fprintf(buf, "aws_access_key_id = %s\n", accessKey)
	fmt.Fprintf(buf, "aws_secret_access_key = %s\n", accessSecret)
	return buf.Bytes()
}

// useCustomCABundle will set up a custom CA bundle in the AWS options if a CA bundle is configured in the
// kube cloud config.
func useCustomCABundle(awsOptions *session.Options, configManagedClient client.Client) error {
	cm := &corev1.ConfigMap{}
	switch err := configManagedClient.Get(
		context.Background(),
		client.ObjectKey{Namespace: KubeCloudConfigNamespace, Name: kubeCloudConfigName},
		cm,
	); {
	case apimachineryerrors.IsNotFound(err):
		// no cloud config ConfigMap, so no custom CA bundle
		return nil
	case err != nil:
		return fmt.Errorf("failed to get kube-cloud-config ConfigMap: %w", err)
	}
	caBundle, ok := cm.Data[cloudCABundleKey]
	if !ok {
		// no "ca-bundle.pem" key in the ConfigMap, so no custom CA bundle
		return nil
	}
	klog.Info("using a custom CA bundle")
	awsOptions.CustomCABundle = strings.NewReader(caBundle)
	return nil
}
