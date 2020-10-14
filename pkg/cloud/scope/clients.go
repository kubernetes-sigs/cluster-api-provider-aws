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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"

	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	awsmetrics "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/metrics"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/version"
)

// NewEC2Client creates a new EC2 API client for a given session
func NewEC2Client(scopeUser cloud.ScopeUsage, session cloud.Session, target runtime.Object) ec2iface.EC2API {
	ec2Client := ec2.New(session.Session())
	ec2Client.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	ec2Client.Handlers.Sign.PushFront(session.ServiceLimiter(ec2.ServiceID).LimitRequest)
	ec2Client.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	ec2Client.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(ec2.ServiceID).ReviewResponse)
	ec2Client.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return ec2Client
}

// NewELBClient creates a new ELB API client for a given session
func NewELBClient(scopeUser cloud.ScopeUsage, session cloud.Session, target runtime.Object) elbiface.ELBAPI {
	elbClient := elb.New(session.Session())
	elbClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	elbClient.Handlers.Sign.PushFront(session.ServiceLimiter(elb.ServiceID).LimitRequest)
	elbClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	elbClient.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(elb.ServiceID).ReviewResponse)
	elbClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return elbClient
}

// NewResourgeTaggingClient creates a new Resource Tagging API client for a given session
func NewResourgeTaggingClient(scopeUser cloud.ScopeUsage, session cloud.Session, target runtime.Object) resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI {
	resourceTagging := resourcegroupstaggingapi.New(session.Session())
	resourceTagging.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	resourceTagging.Handlers.Sign.PushFront(session.ServiceLimiter(resourceTagging.ServiceID).LimitRequest)
	resourceTagging.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	resourceTagging.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(resourceTagging.ServiceID).ReviewResponse)
	resourceTagging.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return resourceTagging
}

// NewSecretsManagerClient creates a new Secrets API client for a given session
func NewSecretsManagerClient(scopeUser cloud.ScopeUsage, session cloud.Session, target runtime.Object) secretsmanageriface.SecretsManagerAPI {
	secretsClient := secretsmanager.New(session.Session())
	secretsClient.Handlers.Build.PushFrontNamed(getUserAgentHandler())
	secretsClient.Handlers.Sign.PushFront(session.ServiceLimiter(secretsClient.ServiceID).LimitRequest)
	secretsClient.Handlers.CompleteAttempt.PushFront(awsmetrics.CaptureRequestMetrics(scopeUser.ControllerName()))
	secretsClient.Handlers.CompleteAttempt.PushFront(session.ServiceLimiter(secretsClient.ServiceID).ReviewResponse)
	secretsClient.Handlers.Complete.PushBack(recordAWSPermissionsIssue(target))

	return secretsClient
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
