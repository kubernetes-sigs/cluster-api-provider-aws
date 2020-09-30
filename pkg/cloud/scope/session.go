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
)

// ServiceEndpoint defines a tuple containing AWS Service resolution information
type ServiceEndpoint struct {
	ServiceID     string
	URL           string
	SigningRegion string
}

var sessionCache sync.Map

func sessionForRegion(region string, endpoint []ServiceEndpoint) (*session.Session, error) {
	s, ok := sessionCache.Load(region)
	if ok {
		return s.(*session.Session), nil
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
		return nil, err
	}

	sessionCache.Store(region, ns)
	return ns, nil
}
