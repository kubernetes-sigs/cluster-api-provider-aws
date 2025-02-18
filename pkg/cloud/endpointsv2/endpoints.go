/*
Copyright 2020 The Kubernetes Authors.

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

// Package endpoints contains aws endpoint related utilities.
package endpointsv2

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

var (
	errServiceEndpointFormat             = errors.New("must be formatted as ${ServiceID}=${URL}")
	errServiceEndpointSigningRegion      = errors.New("must be formatted as ${SigningRegion}:${ServiceID1}=${URL1},${ServiceID2}=${URL2...}")
	errServiceEndpointURL                = errors.New("must use a valid URL as a service-endpoint")
	errServiceEndpointDuplicateServiceID = errors.New("same serviceID defined twice for signing region")
	// ServiceEndpointsMap Can be made private after Go SDK V2 migration
	ServiceEndpointsMap map[string]ServiceEndpoint
)

// ServiceEndpoint contains AWS Service resolution information for SDK V2.
type ServiceEndpoint struct {
	ServiceID     string
	URL           string
	SigningRegion string
}

// ParseFlag parses the command line flag of service endponts in the format ${SigningRegion1}:${ServiceID1}=${URL1},${ServiceID2}=${URL2}...;${SigningRegion2}...
// returning a set of ServiceEndpoints.
func ParseFlag(serviceEndpoints string) error {
	if serviceEndpoints == "" {
		return nil
	}

	signingRegionConfigs := strings.Split(serviceEndpoints, ";")
	for _, regionConfig := range signingRegionConfigs {
		components := strings.SplitN(regionConfig, ":", 2)
		if len(components) != 2 {
			return errServiceEndpointSigningRegion
		}
		signingRegion := components[0]
		servicePairs := strings.Split(components[1], ",")
		for _, servicePair := range servicePairs {
			kv := strings.Split(servicePair, "=")
			if len(kv) != 2 {
				return errServiceEndpointFormat
			}
			var serviceID = kv[0]
			if _, ok := ServiceEndpointsMap[serviceID]; ok {
				return errServiceEndpointDuplicateServiceID
			}

			URL, err := url.ParseRequestURI(kv[1])
			if err != nil {
				return errServiceEndpointURL
			}
			ServiceEndpointsMap[serviceID] = ServiceEndpoint{
				ServiceID:     serviceID,
				URL:           URL.String(),
				SigningRegion: signingRegion,
			}
		}
	}
	return nil
}

// Custom EndpointResolverV2 ResolveEndpoint handlers

// MultiServiceEndpointResolver implements EndpointResolverV2 interface for services
type MultiServiceEndpointResolver struct {
	endpoints map[string]ServiceEndpoint
}

// NewMultiServiceEndpointResolver returns new MultiServiceEndpointResolver
func NewMultiServiceEndpointResolver() *MultiServiceEndpointResolver {
	return &MultiServiceEndpointResolver{
		endpoints: ServiceEndpointsMap,
	}
}

// S3EndpointResolver implements EndpointResolverV2 interface for S3
type S3EndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for S3
func (s *S3EndpointResolver) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	if _, ok := s.endpoints[s3.ServiceID]; !ok {
		return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	endpoint := ServiceEndpointsMap[s3.ServiceID]
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}
