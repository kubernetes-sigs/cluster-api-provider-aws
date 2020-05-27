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
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

var (
	sessionCache sync.Map
)

func sessionForClusterWithRegion(k8sClient client.Client, awsCluster *infrav1.AWSCluster, region string, logger logr.Logger) (*session.Session, error) {
	log := logger.WithName("AWSSession")
	log.Info("Creating a new AWS Session")
	s, ok := sessionCache.Load(region)
	if ok {
		return s.(*session.Session), nil
	}

	awsConfig := aws.NewConfig()
	provider, err := getProviderForCluster(context.Background(), k8sClient, awsCluster, log)
	if err != nil {
		return nil, fmt.Errorf("Failed to get provider for cluster: %v", err)
	}
	if provider != nil {

		// load an existing matching provider from the cache if such a provider exists
		providerHash, err := provider.Hash()
		cachedProvider, ok := sessionCache.Load(providerHash)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve provider from cache: %v", err)
		}
		if ok {
			provider = cachedProvider.(AWSPrincipalTypeProvider)
		} else {
			// add this provider to the cache
			sessionCache.Store(providerHash, provider)
		}

		// set the awsconfig to use the credential provider we located
		awsConfig.Credentials = credentials.NewCredentials(provider)
	}


	ns, err := session.NewSession(awsConfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("Failed to create a new AWS session: %v", err)
	}

	sessionCache.Store(region, ns)
	return ns, nil
}

func getProviderForCluster(ctx context.Context, k8sClient client.Client, awsCluster *infrav1.AWSCluster, log logr.Logger) (AWSPrincipalTypeProvider, error) {
	log.Info("Checking for a provider for the cluster")
	var provider AWSPrincipalTypeProvider
	if awsCluster.Spec.PrincipalRef != nil {
		principalObjectKey := client.ObjectKey{Name: awsCluster.Spec.PrincipalRef.Name}
		log.Info("Looking for a provider", "key", principalObjectKey)
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
	} else {
		log.Info("AWSCluster does not have a PrincipalRef specified")
	}

	return provider, nil
}

type AWSPrincipalTypeProvider interface {
	credentials.Provider
	// Hash returns a unique hash of the data forming the credentials
	// for this Principal
	Hash() (string, error)
}

func NewAWSStaticPrincipalTypeProvider(principal *infrav1.AWSClusterStaticPrincipal, secret *corev1.Secret) (*AWSStaticPrincipalTypeProvider) {
	accessKeyId := string(secret.Data["AccessKeyID"])
	secretAccessKey := string(secret.Data["SecretAccessKey"])
	sessionToken := string(secret.Data["SessionToken"])

	return &AWSStaticPrincipalTypeProvider{
		Principal:       principal,
		credentials:     credentials.NewStaticCredentials(accessKeyId,secretAccessKey,sessionToken),
		accessKeyId:     accessKeyId,
		secretAccessKey: secretAccessKey,
		sessionToken:    sessionToken,
	}
}

func NewAWSRolePrincipalTypeProvider(principal *infrav1.AWSClusterRolePrincipal) (*AWSRolePrincipalTypeProvider) {
	roleProvider := &stscreds.AssumeRoleProvider{
		RoleARN: principal.Spec.RoleArn,
		ExternalID: aws.String(principal.Spec.ExternalID),
		// Duration: time.Second * Principal.Spec.DurationSeconds,// TODO: fixme
		RoleSessionName: principal.Spec.SessionName,
		Policy: aws.String(principal.Spec.InlinePolicy),
	}
	return &AWSRolePrincipalTypeProvider{
		credentials: credentials.NewCredentials(roleProvider),
		Principal:   principal,
	}
}

func NewAWSServiceAccountPrincipalTypeProvider(principal *infrav1.AWSServiceAccountPrincipal) (*AWSServiceAccountPrincipalTypeProvider) {
	return &AWSServiceAccountPrincipalTypeProvider{
		Principal: principal,
	}
}

type AWSStaticPrincipalTypeProvider struct {
	Principal   *infrav1.AWSClusterStaticPrincipal
	credentials *credentials.Credentials
	// these are for tests :/
	accessKeyId string
	secretAccessKey string
	sessionToken string
}
func (p *AWSStaticPrincipalTypeProvider) Hash() (string,error) {
	var roleIdentityValue bytes.Buffer
	err := gob.NewEncoder(&roleIdentityValue).Encode(p)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	return string(hash.Sum(roleIdentityValue.Bytes())), nil
}
func (p *AWSStaticPrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	return p.credentials.Get()
}
func (p *AWSStaticPrincipalTypeProvider) IsExpired() bool {
	return p.credentials.IsExpired()
}

type AWSRolePrincipalTypeProvider struct {
	Principal   *infrav1.AWSClusterRolePrincipal
	credentials *credentials.Credentials
}

func (p *AWSRolePrincipalTypeProvider) Hash() (string,error) {
	var roleIdentityValue bytes.Buffer
	err := gob.NewEncoder(&roleIdentityValue).Encode(p)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	return string(hash.Sum(roleIdentityValue.Bytes())), nil
}

func (p *AWSRolePrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	return p.credentials.Get()
}
func (p *AWSRolePrincipalTypeProvider) IsExpired() bool {
	return p.credentials.IsExpired()
}

type AWSServiceAccountPrincipalTypeProvider struct {
	Principal * infrav1.AWSServiceAccountPrincipal
}

func (p *AWSServiceAccountPrincipalTypeProvider) Hash() (string,error) {
	var roleIdentityValue bytes.Buffer
	err := gob.NewEncoder(&roleIdentityValue).Encode(p)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	return string(hash.Sum(roleIdentityValue.Bytes())), nil
}
func (p *AWSServiceAccountPrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{}, nil
}
func (p *AWSServiceAccountPrincipalTypeProvider) IsExpired() bool {
	return false
}