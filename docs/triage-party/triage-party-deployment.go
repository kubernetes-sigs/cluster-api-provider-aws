/*
Copyright 2021 The Kubernetes Authors.

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

package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk"

	"os"

	"github.com/aws/aws-cdk-go/awscdk/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/awsecspatterns"
	"github.com/aws/aws-cdk-go/awscdk/awselasticloadbalancingv2"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type TriagePartyDeploymentStackProps struct {
	awscdk.StackProps
}

func NewTriagePartyDeploymentStack(scope constructs.Construct, id string, props *TriagePartyDeploymentStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	githubToken := os.Getenv("TRIAGE_PARTY_GITHUB_TOKEN")
	if githubToken == "" {
		fmt.Println("Require TRIAGE_PARTY_GITHUB_TOKEN")
		os.Exit(10)
	}

	service := awsecspatterns.NewApplicationLoadBalancedFargateService(stack, jsii.String("TriageParty"), &awsecspatterns.ApplicationLoadBalancedFargateServiceProps{
		HealthCheckGracePeriod: awscdk.Duration_Seconds(jsii.Number(600.0)),
		TaskImageOptions: &awsecspatterns.ApplicationLoadBalancedTaskImageOptions{
			ContainerPort: jsii.Number(8080),
			Image:         awsecs.ContainerImage_FromAsset(jsii.String("."), nil),
			Environment: &map[string]*string{
				"GITHUB_TOKEN": jsii.String(githubToken),
			},
		},
		PublicLoadBalancer: jsii.Bool(true),
	})

	service.TargetGroup().ConfigureHealthCheck(
		&awselasticloadbalancingv2.HealthCheck{
			Path: jsii.String("/healthz"),
		},
	)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewTriagePartyDeploymentStack(app, "TriagePartyDeploymentStack", &TriagePartyDeploymentStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// TODO @randomvariable: Currently dev account, change to heptio OSS / CNCF Publishing
	return &awscdk.Environment{
		Account: jsii.String("483887913171"),
		Region:  jsii.String("us-east-1"),
	}
}
