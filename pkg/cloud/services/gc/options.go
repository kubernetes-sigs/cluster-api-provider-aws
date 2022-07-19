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

package gc

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"
)

// ServiceOption is an option for creating the service.
type ServiceOption func(*Service)

// WithELBClient is an option for specifying a AWS ELB Client.
func WithELBClient(client elbiface.ELBAPI) ServiceOption {
	return func(s *Service) {
		s.elbClient = client
	}
}

// WithELBv2Client is an option for specifying a AWS ELBv2 Client.
func WithELBv2Client(client elbv2iface.ELBV2API) ServiceOption {
	return func(s *Service) {
		s.elbv2Client = client
	}
}

// WithResourceTaggingClient is an option for specifying a AWS Resource Tagging Client.
func WithResourceTaggingClient(client resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI) ServiceOption {
	return func(s *Service) {
		s.resourceTaggingClient = client
	}
}

// WithEC2Client is an option for specifying a AWS EC2 Client.
func WithEC2Client(client ec2iface.EC2API) ServiceOption {
	return func(s *Service) {
		s.ec2Client = client
	}
}

// WithCleanupFuns is an option for specifying clean-up functions.
func WithCleanupFuns(funcs map[string]ResourceCleanupFunc) ServiceOption {
	return func(s *Service) {
		for serviceName, fn := range funcs {
			s.cleanupFuncs[serviceName] = fn
		}
	}
}
