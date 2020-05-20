package scope

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
)

type AWSPrincipalTypeProvider interface {
	credentials.Provider
	// Hash returns a unique hash of the data forming the credentials
	// for this Principal
	Hash() (string, error)
}

func NewAWSStaticPrincipalTypeProvider(principal *infrav1.AWSClusterStaticPrincipal, secret *corev1.Secret) *AWSStaticPrincipalTypeProvider {
	accessKeyID := string(secret.Data["AccessKeyID"])
	secretAccessKey := string(secret.Data["SecretAccessKey"])
	sessionToken := string(secret.Data["SessionToken"])

	return &AWSStaticPrincipalTypeProvider{
		Principal:       principal,
		credentials:     credentials.NewStaticCredentials(accessKeyID, secretAccessKey, sessionToken),
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		sessionToken:    sessionToken,
	}
}

func NewAWSRolePrincipalTypeProvider(principal *infrav1.AWSClusterRolePrincipal, awsConfig *aws.Config, log logr.Logger) *AWSRolePrincipalTypeProvider {
	sess := session.Must(session.NewSession(awsConfig))

	creds := stscreds.NewCredentials(sess, principal.Spec.RoleArn, func(p *stscreds.AssumeRoleProvider) {
		if principal.Spec.ExternalID != "" {
			p.ExternalID = aws.String(principal.Spec.ExternalID)
		}
		p.RoleSessionName = principal.Spec.SessionName
		if principal.Spec.InlinePolicy != "" {
			p.Policy = aws.String(principal.Spec.InlinePolicy)
		}
	})

	return &AWSRolePrincipalTypeProvider{
		credentials: creds,
		Principal:   principal,
		log:         log.WithName("AWSRolePrincipalTypeProvider"),
	}
}

type AWSStaticPrincipalTypeProvider struct {
	Principal   *infrav1.AWSClusterStaticPrincipal
	credentials *credentials.Credentials
	// these are for tests :/
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
}

func (p *AWSStaticPrincipalTypeProvider) Hash() (string, error) {
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
	log         logr.Logger
}

func (p *AWSRolePrincipalTypeProvider) Hash() (string, error) {
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
