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

package identity

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	token "sigs.k8s.io/cluster-api-provider-aws/pkg/internal/token"
)

// AWSPrincipalTypeProvider defines the interface for AWS Principal Type Provider.
type AWSPrincipalTypeProvider interface {
	credentials.Provider
	// Hash returns a unique hash of the data forming the credentials
	// for this Principal
	Hash() (string, error)
	Name() string
}

// NewAWSStaticPrincipalTypeProvider will create a new AWSStaticPrincipalTypeProvider from a given AWSClusterStaticIdentity.
func NewAWSStaticPrincipalTypeProvider(identity *infrav1.AWSClusterStaticIdentity, secret *corev1.Secret) *AWSStaticPrincipalTypeProvider {
	accessKeyID := string(secret.Data["AccessKeyID"])
	secretAccessKey := string(secret.Data["SecretAccessKey"])
	sessionToken := string(secret.Data["SessionToken"])

	return &AWSStaticPrincipalTypeProvider{
		Principal:       identity,
		credentials:     credentials.NewStaticCredentials(accessKeyID, secretAccessKey, sessionToken),
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		SessionToken:    sessionToken,
	}
}

// GetAssumeRoleCredentials will return the Credentials of a given AWSRolePrincipalTypeProvider.
func GetAssumeRoleCredentials(roleIdentityProvider *AWSRolePrincipalTypeProvider, awsConfig *aws.Config) *credentials.Credentials {
	sess := session.Must(session.NewSession(awsConfig))

	creds := stscreds.NewCredentials(sess, roleIdentityProvider.Principal.Spec.RoleArn, func(p *stscreds.AssumeRoleProvider) {
		if roleIdentityProvider.Principal.Spec.ExternalID != "" {
			p.ExternalID = aws.String(roleIdentityProvider.Principal.Spec.ExternalID)
		}
		p.RoleSessionName = roleIdentityProvider.Principal.Spec.SessionName
		if roleIdentityProvider.Principal.Spec.InlinePolicy != "" {
			p.Policy = aws.String(roleIdentityProvider.Principal.Spec.InlinePolicy)
		}
		p.Duration = time.Duration(roleIdentityProvider.Principal.Spec.DurationSeconds) * time.Second
		// For testing
		if roleIdentityProvider.stsClient != nil {
			p.Client = roleIdentityProvider.stsClient
		}
	})
	return creds
}

// NewAWSRolePrincipalTypeProvider will create a new AWSRolePrincipalTypeProvider from an AWSClusterRoleIdentity.
func NewAWSRolePrincipalTypeProvider(identity *infrav1.AWSClusterRoleIdentity, sourceProvider *AWSPrincipalTypeProvider, log logr.Logger) *AWSRolePrincipalTypeProvider {
	return &AWSRolePrincipalTypeProvider{
		credentials:    nil,
		stsClient:      nil,
		Principal:      identity,
		sourceProvider: sourceProvider,
		log:            log.WithName("AWSRolePrincipalTypeProvider"),
	}
}

// AWSStaticPrincipalTypeProvider defines the specs for a static AWSPrincipalTypeProvider.
type AWSStaticPrincipalTypeProvider struct {
	Principal   *infrav1.AWSClusterStaticIdentity
	credentials *credentials.Credentials
	// these are for tests :/
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// Hash returns the byte encoded AWSStaticPrincipalTypeProvider.
func (p *AWSStaticPrincipalTypeProvider) Hash() (string, error) {
	var roleIdentityValue bytes.Buffer
	err := gob.NewEncoder(&roleIdentityValue).Encode(p)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	return string(hash.Sum(roleIdentityValue.Bytes())), nil
}

// Retrieve returns the credential values for the AWSStaticPrincipalTypeProvider.
func (p *AWSStaticPrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	return p.credentials.Get()
}

// Name returns the name of the AWSStaticPrincipalTypeProvider.
func (p *AWSStaticPrincipalTypeProvider) Name() string {
	return p.Principal.Name
}

// IsExpired checks the expiration state of the AWSStaticPrincipalTypeProvider.
func (p *AWSStaticPrincipalTypeProvider) IsExpired() bool {
	return p.credentials.IsExpired()
}

// AWSRolePrincipalTypeProvider defines the specs for a AWSPrincipalTypeProvider with a role.
type AWSRolePrincipalTypeProvider struct {
	Principal      *infrav1.AWSClusterRoleIdentity
	credentials    *credentials.Credentials
	sourceProvider *AWSPrincipalTypeProvider
	log            logr.Logger
	stsClient      stsiface.STSAPI
}

// Hash returns the byte encoded AWSRolePrincipalTypeProvider.
func (p *AWSRolePrincipalTypeProvider) Hash() (string, error) {
	var roleIdentityValue bytes.Buffer
	err := gob.NewEncoder(&roleIdentityValue).Encode(p)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	return string(hash.Sum(roleIdentityValue.Bytes())), nil
}

// Name returns the name of the AWSRolePrincipalTypeProvider.
func (p *AWSRolePrincipalTypeProvider) Name() string {
	return p.Principal.Name
}

// Retrieve returns the credential values for the AWSRolePrincipalTypeProvider.
func (p *AWSRolePrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	if p.credentials == nil || p.IsExpired() {
		awsConfig := aws.NewConfig()
		if p.sourceProvider != nil {
			sourceCreds, err := (*p.sourceProvider).Retrieve()
			if err != nil {
				return credentials.Value{}, err
			}
			awsConfig = awsConfig.WithCredentials(credentials.NewStaticCredentialsFromCreds(sourceCreds))
		}

		creds := GetAssumeRoleCredentials(p, awsConfig)
		// Update credentials
		p.credentials = creds
	}
	return p.credentials.Get()
}

// IsExpired checks the expiration state of the AWSRolePrincipalTypeProvider.
func (p *AWSRolePrincipalTypeProvider) IsExpired() bool {
	return p.credentials.IsExpired()
}

type AWSServiceAccountPrincipalTypeProvider struct {
	Principal    *infrav1.AWSServiceAccountIdentity
	credentials  *credentials.Credentials
	log          logr.Logger
	stsClient    stsiface.STSAPI
	k8sClient    client.Client
	tokenFetcher *token.ServiceAcountTokenFetcher
}

// NewAWSServiceAccountPrincipalTypeProvider will create a new AWSServiceAccountPrincipalTypeProvider from an AWSClusterRoleIdentity.
func NewAWSServiceAccountPrincipalTypeProvider(identity *infrav1.AWSServiceAccountIdentity, log logr.Logger, k8sClient client.Client) (*AWSServiceAccountPrincipalTypeProvider, error) {
	awsConfig := aws.NewConfig()
	sess := session.Must(session.NewSession(awsConfig))
	tokenFetcher, err := token.NewServiceAccountTokenFetcher()
	if err != nil {
		return nil, err
	}
	return &AWSServiceAccountPrincipalTypeProvider{
		Principal:    identity,
		credentials:  nil,
		stsClient:    sts.New(sess),
		log:          log.WithName("AWSServiceAccountPrincipalTypeProvider"),
		k8sClient:    k8sClient,
		tokenFetcher: tokenFetcher,
	}, nil
}

// Hash returns the byte encoded AWSServiceAccountPrincipalTypeProvider.
func (p *AWSServiceAccountPrincipalTypeProvider) Hash() (string, error) {
	var serviceAccountIdentityValue bytes.Buffer
	err := gob.NewEncoder(&serviceAccountIdentityValue).Encode(p)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	return string(hash.Sum(serviceAccountIdentityValue.Bytes())), nil
}

// Name returns the name of the AWSServiceAccountPrincipalTypeProvider.
func (p *AWSServiceAccountPrincipalTypeProvider) Name() string {
	return p.Principal.Name
}

// Retrieve returns the credential values for the AWSServiceAccountPrincipalTypeProvider.
func (p *AWSServiceAccountPrincipalTypeProvider) Retrieve() (credentials.Value, error) {
	if p.credentials == nil || p.IsExpired() {
		if err := p.checkOrCreateServiceAccount(); err != nil {
			return credentials.Value{}, err
		}
		params := &token.ServiceAccountTokenFetcherParams{
			ServiceAccount:    p.Principal.Name,
			Namespace:         p.Principal.Namespace,
			ExpirationSeconds: int64(p.Principal.Spec.ExpirationSeconds),
			Audiences:         p.Principal.Spec.Audience,
		}
		token, err := p.tokenFetcher.FetchToken(params)
		if err != nil {
			return credentials.Value{}, err
		}
		// Update credentials
		p.credentials = GetAssumeRoleWithWebIdentityCredentials(p, token)
	}
	return p.credentials.Get()
}

func (p *AWSServiceAccountPrincipalTypeProvider) checkOrCreateServiceAccount() error {
	serviceAccount := &corev1.ServiceAccount{}
	serviceAccountKey := client.ObjectKey{Name: p.Principal.Name, Namespace: p.Principal.Namespace}
	err := p.k8sClient.Get(context.TODO(), serviceAccountKey, serviceAccount)
	if err != nil && apierrors.IsNotFound(err) {
		serviceAccount := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      p.Principal.Name,
				Namespace: p.Principal.Namespace,
			},
		}
		if err = p.k8sClient.Create(context.TODO(), serviceAccount, &client.CreateOptions{}); err != nil {
			return err
		}
	}
	return err
}

// IsExpired checks the expiration state of the AWSServiceAccountPrincipalTypeProvider.
func (p *AWSServiceAccountPrincipalTypeProvider) IsExpired() bool {
	return p.credentials.IsExpired()
}

func GetAssumeRoleWithWebIdentityCredentials(serviceAccountIdentityProvider *AWSServiceAccountPrincipalTypeProvider, token []byte) *credentials.Credentials {
	creds := credentials.NewCredentials(stscreds.NewWebIdentityRoleProviderWithToken(serviceAccountIdentityProvider.stsClient,
		serviceAccountIdentityProvider.Principal.Spec.RoleArn, "session-name", Token(token)))
	return creds
}

// Token type implements TokenFetcher interface as expected by stscreds.NewWebIdentityRoleProviderWithToken().
type Token []byte

func (t Token) FetchToken(ctx credentials.Context) ([]byte, error) {
	return t, nil
}
