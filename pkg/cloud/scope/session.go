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
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	sessionCache sync.Map
)

func sessionForClusterWithRegion(k8sClient client.Client, awsCluster *infrav1.AWSCluster, region string) (*session.Session, error) {
	s, ok := sessionCache.Load(region)
	if ok {
		return s.(*session.Session), nil
	}

	awsConfig := aws.NewConfig()
	provider, err := getProviderForCluster(context.Background(), k8sClient, awsCluster)
	if err != nil {
		return nil, fmt.Errorf("Failed to get provider for cluster: %v", err)
	}
	if provider != nil {
		awsConfig.Credentials = credentials.NewCredentials(provider)
	}

	ns, err := session.NewSession(awsConfig.WithRegion(region))
	if err != nil {
		return nil, err
	}

	sessionCache.Store(region, ns)
	return ns, nil
}

func getProviderForCluster(ctx context.Context, k8sClient client.Client, awsCluster *infrav1.AWSCluster) (AWSPrincipalTypeProvider, error) {
	var provider AWSPrincipalTypeProvider
	if awsCluster.Spec.PrincipalRef != nil {
		principalObjectKey := client.ObjectKey{Namespace: awsCluster.Spec.PrincipalRef.Namespace, Name: awsCluster.Spec.PrincipalRef.Name}
		switch awsCluster.Spec.PrincipalRef.Kind {
		case "AWSClusterStaticPrincipal":
			principal := &infrav1.AWSClusterStaticPrincipal{}
			err := k8sClient.Get(ctx, principalObjectKey, principal)
			if err != nil {
				return nil, err
			}
			secret := &corev1.Secret{}
			err = k8sClient.Get(ctx, client.ObjectKey{Name: principal.Spec.SecretRef.Name, Namespace: principal.Spec.SecretRef.Namespace}, secret)
			if err != nil {
				return nil, err
			}
			provider = NewAWSStaticPrincipalTypeProvider(principal, secret)
		case "AWSClusterRolePrincipal":
			principal := &infrav1.AWSClusterRolePrincipal{}
			err := k8sClient.Get(ctx, principalObjectKey, principal)
			if err != nil {
				return nil, err
			}
			provider = NewAWSRolePrincipalTypeProvider(principal)
		case "AWSServiceAccountPrincipal":
			principal := &infrav1.AWSServiceAccountPrincipal{}
			err := k8sClient.Get(ctx, principalObjectKey, principal)
			if err != nil {
				return nil, err
			}
			provider = NewAWSServiceAccountPrincipalTypeProvider(principal)
		default:
			return nil, fmt.Errorf("No such provider known: '%s'",awsCluster.Spec.PrincipalRef.Kind)
		}
	}

	return provider, nil
}

type AWSPrincipalTypeProvider interface {
	credentials.Provider
	// Hash returns a unique hash of the data forming the credentials
	// for this principal
	Hash() (string, error)
}

func NewAWSStaticPrincipalTypeProvider(principal *infrav1.AWSClusterStaticPrincipal, secret *corev1.Secret) (*AWSStaticPrincipalTypeProvider) {
	return &AWSStaticPrincipalTypeProvider{
		principal: principal,
		accessKeyId: string(secret.Data["AccessKeyID"]),
		secretAccessKey: string(secret.Data["SecretAccessKey"]),
		sessionToken: string(secret.Data["SessionToken"]),
	}
}

func NewAWSRolePrincipalTypeProvider(principal *infrav1.AWSClusterRolePrincipal) (*AWSRolePrincipalTypeProvider) {
	return &AWSRolePrincipalTypeProvider{
		principal: principal,
	}
}

func NewAWSServiceAccountPrincipalTypeProvider(principal *infrav1.AWSServiceAccountPrincipal) (*AWSServiceAccountPrincipalTypeProvider) {
	return &AWSServiceAccountPrincipalTypeProvider{
		principal: principal,
	}
}

type AWSStaticPrincipalTypeProvider struct {
	principal *infrav1.AWSClusterStaticPrincipal
	accessKeyId string
	secretAccessKey string
	sessionToken string
}
func (p *AWSStaticPrincipalTypeProvider) Hash() (string,error) {
	return "", nil
}
func (p *AWSStaticPrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{}, nil
}
func (p *AWSStaticPrincipalTypeProvider) IsExpired() bool {
	return false
}

type AWSRolePrincipalTypeProvider struct {
	principal * infrav1.AWSClusterRolePrincipal
}

func (p *AWSRolePrincipalTypeProvider) Hash() (string,error) {
	return "", nil
}
func (p *AWSRolePrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{}, nil
}
func (p *AWSRolePrincipalTypeProvider) IsExpired() bool {
	return false
}

type AWSServiceAccountPrincipalTypeProvider struct {
	principal * infrav1.AWSServiceAccountPrincipal
}

func (p *AWSServiceAccountPrincipalTypeProvider) Hash() (string,error) {
	return "", nil
}
func (p *AWSServiceAccountPrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{}, nil
}
func (p *AWSServiceAccountPrincipalTypeProvider) IsExpired() bool {
	return false
}