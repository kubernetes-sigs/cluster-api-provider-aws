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
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceEndpoint defines a tuple containing AWS Service resolution information
type ServiceEndpoint struct {
	ServiceID     string
	URL           string
	SigningRegion string
}

var sessionCache sync.Map

func sessionForClusterWithRegion(k8sClient client.Client, awsCluster *infrav1.AWSCluster, region string, endpoint []ServiceEndpoint, logger logr.Logger) (*session.Session, error) {
	log := logger.WithName("identity")
	log.V(4).Info("Creating an AWS Session")
	s, ok := sessionCache.Load(region)
	if ok {
		return s.(*session.Session), nil
	}

	// start
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

	awsConfig := &aws.Config{
		Region:           aws.String(region),
		EndpointResolver: endpoints.ResolverFunc(resolver),
	}
	providers, err := getProvidersForCluster(context.Background(), k8sClient, awsCluster, awsConfig, log)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get providers for cluster")
	}
	// end

	awsProviders := make([]credentials.Provider, len(providers))
	for i, provider := range providers {
		// load an existing matching providers from the cache if such a providers exists
		providerHash, err := provider.Hash()
		cachedProvider, ok := sessionCache.Load(providerHash)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to retrieve providers from cache")
		}
		if ok {
			provider = cachedProvider.(AWSPrincipalTypeProvider)
		} else {
			// add this providers to the cache
			sessionCache.Store(providerHash, provider)
		}
		awsProviders[i] = provider.(credentials.Provider)
	}

	if len(awsProviders) > 0 {
		awsConfig.Credentials = credentials.NewChainCredentials(awsProviders)
	}

	ns, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create a new AWS session")
	}

	sessionCache.Store(region, ns)
	return ns, nil
}

func buildProvidersForRef(ctx context.Context, providers []AWSPrincipalTypeProvider, k8sClient client.Client, awsCluster *infrav1.AWSCluster, ref *corev1.ObjectReference, awsConfig *aws.Config, log logr.Logger) ([]AWSPrincipalTypeProvider, error) {
	if ref == nil {
		log.V(4).Info("AWSCluster does not have a PrincipalRef specified")
		return providers, nil
	}

	// if the namespace isn't specified then assume it's in the same namespace as the AWSCluster
	namespace := ref.Namespace
	if namespace == "" {
		namespace = awsCluster.Namespace
	}

	var provider AWSPrincipalTypeProvider
	principalObjectKey := client.ObjectKey{Name: ref.Name, Namespace: namespace}
	log.V(4).Info("Get Principal", "Key", principalObjectKey)
	switch ref.Kind {
	case "AWSClusterStaticPrincipal":
		principal := &infrav1.AWSClusterStaticPrincipal{}
		err := k8sClient.Get(ctx, principalObjectKey, principal)
		if err != nil {
			return providers, err
		}
		secret := &corev1.Secret{}
		err = k8sClient.Get(ctx, client.ObjectKey{Name: principal.Spec.SecretRef.Name, Namespace: principal.Spec.SecretRef.Namespace}, secret)
		if err != nil {
			return providers, err
		}
		log.V(4).Info("Found an AWSClusterStaticPrincipal", "principal", principal.GetName())
		if !clusterIsPermittedToUsePrincipal(awsCluster, principal.Spec.AWSClusterPrincipalSpec) {
			return providers, errors.Errorf("cluster %s%s is not permitted to use principal %s/%s", awsCluster.Namespace, awsCluster.Name, principal.Namespace, principal.Name)
		}
		provider = NewAWSStaticPrincipalTypeProvider(principal, secret)
		providers = append(providers, provider)
		if principal.Spec.SourcePrincipalRef != nil {
			providers, err = buildProvidersForRef(ctx, providers, k8sClient, awsCluster, principal.Spec.SourcePrincipalRef, awsConfig, log)
			if err != nil {
				return providers, client.IgnoreNotFound(err)
			}
		}
	case "AWSClusterRolePrincipal":
		principal := &infrav1.AWSClusterRolePrincipal{}
		err := k8sClient.Get(ctx, principalObjectKey, principal)
		if err != nil {
			return providers, err
		}
		log.V(4).Info("Found an AWSClusterRolePrincipal", "principal", principal.GetName())
		provider = NewAWSRolePrincipalTypeProvider(principal, awsConfig, log)
		providers = append(providers, provider)
		if principal.Spec.SourcePrincipalRef != nil {
			providers, err = buildProvidersForRef(ctx, providers, k8sClient, awsCluster, principal.Spec.SourcePrincipalRef, awsConfig, log)
			if err != nil {
				return providers, client.IgnoreNotFound(err)
			}
		}
	default:
		return providers, errors.Errorf("No such provider known: '%s'", ref.Kind)
	}

	return providers, nil
}

func getProvidersForCluster(ctx context.Context, k8sClient client.Client, awsCluster *infrav1.AWSCluster, awsConfig *aws.Config, log logr.Logger) ([]AWSPrincipalTypeProvider, error) {
	providers := make([]AWSPrincipalTypeProvider, 0)
	providers, err := buildProvidersForRef(ctx, providers, k8sClient, awsCluster, awsCluster.Spec.PrincipalRef, awsConfig, log)
	if err != nil {
		return nil, err
	}

	return providers, nil
}

func clusterIsPermittedToUsePrincipal(awsCluster *infrav1.AWSCluster, principalSpec infrav1.AWSClusterPrincipalSpec) bool {
	// TODO (andrewmy):
	// https://github.com/randomvariable/cluster-api-provider-aws/blob/2f7b382b70ccbf7c2b4b56f9a14227c5b422b698/docs/proposal/20200506-single-controller-multitenancy.md#implementation-detailsnotesconstraints
	// AllowedNamespaces is a selector of namespaces that AWSClusters can
	// use this ClusterPrincipal from. This is a standard Kubernetes LabelSelector,
	// a label query over a set of resources. The result of matchLabels and
	// matchExpressions are ANDed. Controllers must not support AWSClusters in
	// namespaces outside this selector.
	//
	// An empty selector (default) indicates that AWSClusters can use this
	// AWSClusterPrincipal from any namespace. This field is intentionally not a
	// pointer because the nil behavior (no namespaces) is undesirable here.
	return true
}
