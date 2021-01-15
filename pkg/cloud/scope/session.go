/*
Copyright 2019 The Kubernetes Authors.

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
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/throttle"
)

// ServiceEndpoint defines a tuple containing AWS Service resolution information
type ServiceEndpoint struct {
	ServiceID     string
	URL           string
	SigningRegion string
}

var sessionCache sync.Map

type sessionCacheEntry struct {
	session         *session.Session
	serviceLimiters throttle.ServiceLimiters
}

func sessionForRegion(region string, endpoint []ServiceEndpoint) (*session.Session, throttle.ServiceLimiters, error) {
	if s, ok := sessionCache.Load(region); ok {
		entry := s.(*sessionCacheEntry)
		return entry.session, entry.serviceLimiters, nil
	}

	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		for _, s := range endpoint {
			if service == s.ServiceID {
				return endpoints.ResolvedEndpoint{
					URL:           s.URL,
					SigningRegion: s.SigningRegion,
				}, nil
			}
		}
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}
	ns, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		EndpointResolver: endpoints.ResolverFunc(resolver),
	})
	if err != nil {
		return nil, nil, err
	}

	sl := newServiceLimiters()
	sessionCache.Store(region, &sessionCacheEntry{
		session:         ns,
		serviceLimiters: sl,
	})
	return ns, sl, nil
}

func newServiceLimiters() throttle.ServiceLimiters {
	return throttle.ServiceLimiters{
		ec2.ServiceID:                      newEC2ServiceLimiter(),
		elb.ServiceID:                      newGenericServiceLimiter(),
		resourcegroupstaggingapi.ServiceID: newGenericServiceLimiter(),
		secretsmanager.ServiceID:           newGenericServiceLimiter(),
	}
}

func newGenericServiceLimiter() *throttle.ServiceLimiter {
	return &throttle.ServiceLimiter{
		{
			Operation:  throttle.NewMultiOperationMatch("Describe", "Get", "List"),
			RefillRate: 20.0,
			Burst:      100,
		},
		{
			Operation:  ".*",
			RefillRate: 5.0,
			Burst:      200,
		},
	}
}

func newEC2ServiceLimiter() *throttle.ServiceLimiter {
	return &throttle.ServiceLimiter{
		{
			Operation:  throttle.NewMultiOperationMatch("Describe", "Get"),
			RefillRate: 20.0,
			Burst:      100,
		},
		{
			Operation: throttle.NewMultiOperationMatch(
				"AuthorizeSecurityGroupIngress",
				"CancelSpotInstanceRequests",
				"CreateKeyPair",
				"RequestSpotInstances",
			),
			RefillRate: 20.0,
			Burst:      100,
		},
		{
			Operation:  "RunInstances",
			RefillRate: 2.0,
			Burst:      5,
		},
		{
			Operation:  "StartInstances",
			RefillRate: 2.0,
			Burst:      5,
		},
		{
			Operation:  ".*",
			RefillRate: 5.0,
			Burst:      200,
		},
	}
}
