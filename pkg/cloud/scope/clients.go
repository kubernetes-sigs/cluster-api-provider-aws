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

package scope

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	ec2v2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	elasticloadbalancing "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elasticloadbalancingv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	resourcegroupstaggingapiv2 "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/endpointsv2"
	awslogs "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/logs"
	awsmetrics "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/metrics"
	awsmetricsv2 "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/metricsv2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/v2/version"
)

// ResourceGroupsTaggingAPIAPI is a compatibility layer for the v1 resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI interface.
// It is used to provide a consistent interface for the GetResources method.
type ResourceGroupsTaggingAPIAPI interface {
	resourcegroupstaggingapiv2.GetResourcesAPIClient
}

// ELBV2API is a compatibility layer for the v2 elasticloadbalancingv2.Client interface.
type ELBV2API interface {
	DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancingv2.DescribeLoadBalancersInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeLoadBalancersOutput, error)
	DescribeTargetGroups(ctx context.Context, params *elasticloadbalancingv2.DescribeTargetGroupsInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeTargetGroupsOutput, error)
	DescribeListeners(ctx context.Context, params *elasticloadbalancingv2.DescribeListenersInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeListenersOutput, error)
	DescribeTags(ctx context.Context, params *elasticloadbalancingv2.DescribeTagsInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeTagsOutput, error)
	DeleteLoadBalancer(ctx context.Context, params *elasticloadbalancingv2.DeleteLoadBalancerInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DeleteLoadBalancerOutput, error)
	DeleteTargetGroup(ctx context.Context, params *elasticloadbalancingv2.DeleteTargetGroupInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DeleteTargetGroupOutput, error)
	DeleteListener(ctx context.Context, params *elasticloadbalancingv2.DeleteListenerInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DeleteListenerOutput, error)
}

// EC2API is a compatibility layer for the v1 ec2iface.EC2API interface.
// It is used to provide a consistent interface for the DeleteSecurityGroup method.
type EC2API interface {
	DeleteSecurityGroup(context.Context, *ec2v2.DeleteSecurityGroupInput, ...func(*ec2v2.Options)) (*ec2v2.DeleteSecurityGroupOutput, error)
	ec2v2.DescribeSecurityGroupsAPIClient
}

// ELBAPI is a compatibility layer for the v2 elasticloadbalancing.Client interface.
type ELBAPI interface {
	DeleteLoadBalancer(ctx context.Context, params *elasticloadbalancing.DeleteLoadBalancerInput, optFns ...func(*elasticloadbalancing.Options)) (*elasticloadbalancing.DeleteLoadBalancerOutput, error)
	DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancing.DescribeLoadBalancersInput, optFns ...func(*elasticloadbalancing.Options)) (*elasticloadbalancing.DescribeLoadBalancersOutput, error)
	DescribeTags(ctx context.Context, params *elasticloadbalancing.DescribeTagsInput, optFns ...func(*elasticloadbalancing.Options)) (*elasticloadbalancing.DescribeTagsOutput, error)
}

// NewASGClient creates a new ASG API client for a given session.
func NewASGClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *autoscaling.Client {
	cfg := session.SessionV2()

	autoscalingOpts := []func(*autoscaling.Options){
		func(o *autoscaling.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevelV2(logger.GetLogger())
		},
		autoscaling.WithAPIOptions(
			awsmetricsv2.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetricsv2.WithCAPAUserAgentMiddleware(),
		),
	}

	return autoscaling.NewFromConfig(cfg, autoscalingOpts...)
}

// NewEC2Client creates a new EC2 API client for a given session.
func NewEC2Client(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) ec2iface.EC2API {
	ec2Client := ec2.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	ec2Client.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	if session.ServiceLimiter(ec2.ServiceID) != nil {
		ec2Client.Handlers.Sign.PushFront(session.ServiceLimiter(ec2.ServiceID).LimitRequest)
	}
	ec2Client.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	if session.ServiceLimiter(ec2.ServiceID) != nil {
		ec2Client.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(ec2.ServiceID).ReviewResponse)
	}
	ec2Client.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return ec2Client
}

// NewELBClient creates a new ELB API client for a given session.
func NewELBClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) elbiface.ELBAPI {
	elbClient := elb.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	elbClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	elbClient.Handlers.Sign.PushFront(session.ServiceLimiter(elb.ServiceID).LimitRequest)
	elbClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	elbClient.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(elb.ServiceID).ReviewResponse)
	elbClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return elbClient
}

// NewELBv2Client creates a new ELB v2 API client for a given session.
func NewELBv2Client(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) elbv2iface.ELBV2API {
	elbClient := elbv2.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	elbClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	elbClient.Handlers.Sign.PushFront(session.ServiceLimiter(elbv2.ServiceID).LimitRequest)
	elbClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	elbClient.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(elbv2.ServiceID).ReviewResponse)
	elbClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return elbClient
}

// NewEventBridgeClient creates a new EventBridge API client for a given session.
func NewEventBridgeClient(scopeUser cloud.ScopeUsage, session cloud.Session, target runtime.Object) *eventbridge.Client {
	cfg := session.SessionV2()

	opts := []func(*eventbridge.Options){
		eventbridge.WithAPIOptions(
			awsmetricsv2.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetricsv2.WithCAPAUserAgentMiddleware(),
		),
	}

	return eventbridge.NewFromConfig(cfg, opts...)
}

// NewSQSClient creates a new SQS API client for a given session.
func NewSQSClient(scopeUser cloud.ScopeUsage, session cloud.Session, target runtime.Object) *sqs.Client {
	cfg := session.SessionV2()

	opts := []func(*sqs.Options){
		sqs.WithAPIOptions(
			awsmetricsv2.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetricsv2.WithCAPAUserAgentMiddleware(),
		),
	}

	return sqs.NewFromConfig(cfg, opts...)
}

// NewGlobalSQSClient for creating a new SQS API client that isn't tied to a cluster.
func NewGlobalSQSClient(scopeUser cloud.ScopeUsage, session cloud.Session) *sqs.Client {
	cfg := session.SessionV2()

	opts := []func(*sqs.Options){
		sqs.WithAPIOptions(
			awsmetricsv2.WithRequestMetricContextMiddleware(),
			awsmetricsv2.WithCAPAUserAgentMiddleware(),
		),
	}

	return sqs.NewFromConfig(cfg, opts...)
}

// NewResourgeTaggingClient creates a new Resource Tagging API client for a given session.
func NewResourgeTaggingClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI {
	resourceTagging := resourcegroupstaggingapi.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	resourceTagging.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	resourceTagging.Handlers.Sign.PushFront(session.ServiceLimiter(resourceTagging.ServiceID).LimitRequest)
	resourceTagging.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	resourceTagging.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(resourceTagging.ServiceID).ReviewResponse)
	resourceTagging.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return resourceTagging
}

// NewSecretsManagerClient creates a new Secrets API client for a given session..
func NewSecretsManagerClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) secretsmanageriface.SecretsManagerAPI {
	secretsClient := secretsmanager.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	secretsClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	secretsClient.Handlers.Sign.PushFront(session.ServiceLimiter(secretsClient.ServiceID).LimitRequest)
	secretsClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	secretsClient.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(secretsClient.ServiceID).ReviewResponse)
	secretsClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return secretsClient
}

// NewEKSClient creates a new EKS API client for a given session.
func NewEKSClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *eks.Client {
	cfg := session.SessionV2()
	multiSvcEndpointResolver := endpointsv2.NewMultiServiceEndpointResolver()
	eksEndpointResolver := &endpointsv2.EKSEndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}
	s3Opts := []func(*eks.Options){
		func(o *eks.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevelV2(logger.GetLogger())
			o.EndpointResolverV2 = eksEndpointResolver
		},
		eks.WithAPIOptions(awsmetricsv2.WithMiddlewares(scopeUser.ControllerName(), target), awsmetricsv2.WithCAPAUserAgentMiddleware()),
	}
	return eks.NewFromConfig(cfg, s3Opts...)
}

// NewIAMClient creates a new IAM API client for a given session.
func NewIAMClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *iam.Client {
	cfg := session.SessionV2()

	iamOpts := []func(*iam.Options){
		func(o *iam.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevelV2(logger.GetLogger())
		},
		iam.WithAPIOptions(
			awsmetricsv2.WithMiddlewares(scopeUser.ControllerName(), target),
			awsmetricsv2.WithCAPAUserAgentMiddleware(),
		),
	}

	return iam.NewFromConfig(cfg, iamOpts...)
}

// NewSTSClient creates a new STS API client for a given session.
func NewSTSClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) stsiface.STSAPI {
	stsClient := sts.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	stsClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	stsClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	stsClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return stsClient
}

// NewSSMClient creates a new Secrets API client for a given session.
func NewSSMClient(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) ssmiface.SSMAPI {
	ssmClient := ssm.New(session.Session(), aws.NewConfig().WithLogLevel(awslogs.GetAWSLogLevel(logger.GetLogger())).WithLogger(awslogs.NewWrapLogr(logger.GetLogger())))
	ssmClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	ssmClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	ssmClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return ssmClient
}

// NewResourceTaggingClientV2 creates a new Resource Tagging API client for a given session using AWS SDK v2.
func NewResourceTaggingClientV2(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) ResourceGroupsTaggingAPIAPI {
	cfg := session.SessionV2()
	resourceTaggingOpts := []func(*resourcegroupstaggingapiv2.Options){
		func(o *resourcegroupstaggingapiv2.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevelV2(logger.GetLogger())
		},
		resourcegroupstaggingapiv2.WithAPIOptions(awsmetricsv2.WithMiddlewares(scopeUser.ControllerName(), target), awsmetricsv2.WithCAPAUserAgentMiddleware()),
	}
	return resourcegroupstaggingapiv2.NewFromConfig(cfg, resourceTaggingOpts...)
}

// NewS3Client creates a new S3 API client for a given session.
func NewS3Client(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) *s3.Client {
	cfg := session.SessionV2()
	multiSvcEndpointResolver := endpointsv2.NewMultiServiceEndpointResolver()
	s3EndpointResolver := &endpointsv2.S3EndpointResolver{
		MultiServiceEndpointResolver: multiSvcEndpointResolver,
	}
	s3Opts := []func(*s3.Options){
		func(o *s3.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevelV2(logger.GetLogger())
			o.EndpointResolverV2 = s3EndpointResolver
		},
		s3.WithAPIOptions(awsmetricsv2.WithMiddlewares(scopeUser.ControllerName(), target), awsmetricsv2.WithCAPAUserAgentMiddleware()),
	}
	return s3.NewFromConfig(cfg, s3Opts...)
}

// NewEC2ClientV2 creates a new EC2 API client for a given session using AWS SDK v2.
func NewEC2ClientV2(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) EC2API {
	cfg := session.SessionV2()
	ec2Opts := []func(*ec2v2.Options){
		func(o *ec2v2.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevelV2(logger.GetLogger())
		},
		ec2v2.WithAPIOptions(awsmetricsv2.WithMiddlewares(scopeUser.ControllerName(), target), awsmetricsv2.WithCAPAUserAgentMiddleware()),
	}
	return ec2v2.NewFromConfig(cfg, ec2Opts...)
}

// NewELBV2ClientV2 creates a new ELBV2 API client for a given session using AWS SDK v2.
func NewELBV2ClientV2(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) ELBV2API {
	cfg := session.SessionV2()
	elbOpts := []func(*elasticloadbalancingv2.Options){
		func(o *elasticloadbalancingv2.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevelV2(logger.GetLogger())
		},
		elasticloadbalancingv2.WithAPIOptions(awsmetricsv2.WithMiddlewares(scopeUser.ControllerName(), target), awsmetricsv2.WithCAPAUserAgentMiddleware()),
	}
	return elasticloadbalancingv2.NewFromConfig(cfg, elbOpts...)
}

// NewELBClientV2 creates a new ELB API client for a given session using AWS SDK v2.
func NewELBClientV2(scopeUser cloud.ScopeUsage, session cloud.Session, logger logger.Wrapper, target runtime.Object) ELBAPI {
	cfg := session.SessionV2()
	elbOpts := []func(*elasticloadbalancing.Options){
		func(o *elasticloadbalancing.Options) {
			o.Logger = logger.GetAWSLogger()
			o.ClientLogMode = awslogs.GetAWSLogLevelV2(logger.GetLogger())
		},
		elasticloadbalancing.WithAPIOptions(awsmetricsv2.WithMiddlewares(scopeUser.ControllerName(), target), awsmetricsv2.WithCAPAUserAgentMiddleware()),
	}
	return elasticloadbalancing.NewFromConfig(cfg, elbOpts...)
}

func recordAWSPermissionsIssue(target runtime.Object) func(r *request.Request) {
	return func(r *request.Request) {
		if awsErr, ok := r.Error.(awserr.Error); ok {
			switch awsErr.Code() {
			case "AuthFailure", "UnauthorizedOperation", "NoCredentialProviders":
				record.Warnf(target, awsErr.Code(), "Operation %s failed with a credentials or permission issue", r.Operation.Name)
			}
		}
	}
}

func getUserAgentHandler() request.NamedHandler {
	return request.NamedHandler{
		Name: "capa/user-agent",
		Fn:   request.MakeAddToUserAgentHandler("aws.cluster.x-k8s.io", version.Get().String()),
	}
}

// AWSClients contains all the aws clients used by the scopes.
type AWSClients struct {
	ASG               autoscaling.Client
	EC2               ec2iface.EC2API
	EC2V2             *ec2v2.Client
	ELB               elbiface.ELBAPI
	ELBV2             *elbv2.ELBV2
	SecretsManager    secretsmanageriface.SecretsManagerAPI
	ResourceTagging   resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
	ResourceTaggingV2 *resourcegroupstaggingapiv2.Client
}
