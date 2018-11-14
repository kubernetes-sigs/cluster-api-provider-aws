// Copyright © 2018 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	ec2svc "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/ec2"
	elbsvc "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/elb"
)

// SDKGetter is the default getter based on the AWS SDK, its main purpose
// is to provide cluster-aware clients to be used in machine and cluster actuators.
type SDKGetter struct{}

// Session returns a new AWS session based on the cluster configuration.
func (*SDKGetter) Session(clusterConfig *v1alpha1.AWSClusterProviderConfig) *session.Session {
	return session.Must(session.NewSession(aws.NewConfig().WithRegion(clusterConfig.Region)))
}

// EC2 returns a new EC2 service client.
func (*SDKGetter) EC2(session *session.Session) EC2Interface {
	return ec2svc.NewService(ec2.New(session))
}

// ELB returns a new ELB service client.
func (*SDKGetter) ELB(session *session.Session) ELBInterface {
	return elbsvc.NewService(elb.New(session))
}

// NewSDKGetter returns a new SDKGetter.
func NewSDKGetter() *SDKGetter {
	return new(SDKGetter)
}
