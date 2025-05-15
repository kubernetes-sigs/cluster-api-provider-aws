package iam

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	stderr "errors"
	"fmt"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/go-jose/go-jose/v4"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/s3"
)

const (
	jwksKey                 = "/openid/v1/jwks"
	opendIDConfigurationKey = "/.well-known/openid-configuration"
)

type oidcDiscovery struct {
	Issuer                string   `json:"issuer"`
	JWKSURI               string   `json:"jwks_uri"`
	AuthorizationEndpoint string   `json:"authorization_endpoint"`
	ResponseTypes         []string `json:"response_types_supported"`
	SubjectTypes          []string `json:"subject_types_supported"`
	SigningAlgs           []string `json:"id_token_signing_alg_values_supported"`
	ClaimsSupported       []string `json:"claims_supported"`
}

type jwksDocument struct {
	Keys []jose.JSONWebKey `json:"keys"`
}

func buildDiscoveryJSON(issuerURL string) ([]byte, error) {
	d := oidcDiscovery{
		Issuer:                issuerURL,
		JWKSURI:               fmt.Sprintf("%v/openid/v1/jwks", issuerURL),
		AuthorizationEndpoint: "urn:kubernetes:programmatic_authorization",
		ResponseTypes:         []string{"id_token"},
		SubjectTypes:          []string{"public"},
		SigningAlgs:           []string{"RS256"},
		ClaimsSupported:       []string{"sub", "iss"},
	}
	return json.MarshalIndent(d, "", "")
}

func (s *Service) deleteBucketContents(ctx context.Context, s3 *s3.Service) error {
	if err := s3.DeleteKey(ctx, "/"+path.Join(s.scope.Name(), jwksKey)); err != nil {
		return err
	}

	return s3.DeleteKey(ctx, "/"+path.Join(s.scope.Name(), opendIDConfigurationKey))
}

func (s *Service) buildIssuerURL() string {
	// e.g. s3-us-west-2.amazonaws.com/<bucketname>/<clustername>
	return fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", s.scope.Region(), s.scope.Bucket().Name, s.scope.Name())
}

func (s *Service) reconcileBucketContents(ctx context.Context) error {
	// create the OpenID Connect discovery document
	openIDConfig, err := buildDiscoveryJSON(s.buildIssuerURL())
	if err != nil {
		return err
	}

	s3scope := s3.NewService(s.scope)

	if _, err := s3scope.CreatePublicKey(ctx, path.Join(s.scope.Name(), opendIDConfigurationKey), openIDConfig); err != nil {
		return err
	}

	// Read the <cluster>-sa secret that contains Service Account signing key
	secret := &corev1.Secret{}
	if err := s.scope.ManagementClient().Get(ctx, client.ObjectKey{
		Name:      s.scope.Name() + "-sa",
		Namespace: s.scope.Namespace(),
	}, secret); err != nil {
		return fmt.Errorf("failed to get service account signing secret: %w", err)
	}

	// Create jwks document that will be published to S3
	key, err := createJwksKey(secret.Data["tls.key"])
	if err != nil {
		return fmt.Errorf("failed to create jwks key: %w", err)
	}
	jwksDocument := jwksDocument{Keys: []jose.JSONWebKey{*key}}
	jwksBytes, err := json.MarshalIndent(jwksDocument, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal jwks payload to json: %w", err)
	}

	if _, err := s3scope.CreatePublicKey(ctx, path.Join(s.scope.Name(), jwksKey), jwksBytes); err != nil {
		return err
	}

	return nil
}

func buildOIDCTrustPolicy(arn string) iamv1.PolicyDocument {
	conditionValue := arn[strings.Index(arn, "/")+1:] + ":sub"

	return iamv1.PolicyDocument{
		Version: "2012-10-17",
		Statement: iamv1.Statements{
			iamv1.StatementEntry{
				Sid:    "",
				Effect: "Allow",
				Principal: iamv1.Principals{
					iamv1.PrincipalFederated: iamv1.PrincipalID{arn},
				},
				Action: iamv1.Actions{"sts:AssumeRoleWithWebIdentity"},
				Condition: iamv1.Conditions{
					"ForAnyValue:StringLike": map[string][]string{
						conditionValue: {"system:serviceaccount:${SERVICE_ACCOUNT_NAMESPACE}:${SERVICE_ACCOUNT_NAME}"},
					},
				},
			},
		},
	}
}

// FindAndVerifyOIDCProvider will try to find an OIDC provider. It will return an error if the found provider does not
// match the cluster spec.
func findAndVerifyOIDCProvider(issuerURL, thumbprint string, iamClient iamiface.IAMAPI) (string, error) {
	output, err := iamClient.ListOpenIDConnectProviders(&iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return "", errors.Wrap(err, "error listing providers")
	}
	for _, r := range output.OpenIDConnectProviderList {
		provider, err := iamClient.GetOpenIDConnectProvider(&iam.GetOpenIDConnectProviderInput{OpenIDConnectProviderArn: r.Arn})
		if err != nil {
			return "", errors.Wrap(err, "error getting provider")
		}
		// URL should always contain `https`.
		if "https://"+aws.StringValue(provider.Url) != issuerURL {
			continue
		}
		if len(provider.ThumbprintList) != 1 || aws.StringValue(provider.ThumbprintList[0]) != thumbprint {
			return "", errors.Wrap(err, "found provider with matching issuerURL but with non-matching thumbprint")
		}
		if len(provider.ClientIDList) != 1 || aws.StringValue(provider.ClientIDList[0]) != STSAWSAudience {
			return "", errors.Wrap(err, "found provider with matching issuerURL but with non-matching clientID")
		}
		return aws.StringValue(r.Arn), nil
	}
	return "", nil
}

func fetchRootCAThumbprint(url string, port int) (_ string, err error) {
	// Parse cmdline arguments using flag package
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", url, port), &tls.Config{
		MinVersion: tls.VersionTLS12,
	})
	if err != nil {
		return "", err
	}

	defer func() {
		if cerr := conn.Close(); cerr != nil {
			err = stderr.Join(err, cerr)
		}
	}()
	// Get the ConnectionState struct as that's the one which gives us x509.Certificate struct
	cert := conn.ConnectionState().PeerCertificates[0]
	fingerprint := sha1.Sum(cert.Raw) //nolint:gosec // this is not used for real security
	var buf bytes.Buffer
	for _, f := range fingerprint {
		if _, err := fmt.Fprintf(&buf, "%02X", f); err != nil {
			return "", err
		}
	}

	return strings.ToLower(buf.String()), nil
}

// DeleteOIDCProvider will delete an OIDC provider.
func deleteOIDCProvider(arn string, iamClient iamiface.IAMAPI) error {
	if arn == "" {
		return nil
	}

	input := iam.DeleteOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: aws.String(arn),
	}

	_, err := iamClient.DeleteOpenIDConnectProvider(&input)
	if err != nil {
		var aerr awserr.Error
		ok := errors.As(err, &aerr)
		if !ok {
			return errors.Wrap(err, "deleting OIDC provider")
		}

		switch aerr.Code() {
		case iam.ErrCodeNoSuchEntityException:
			return nil
		default:
			return errors.Wrap(err, "deleting OIDC provider")
		}
	}
	return nil
}

// createJwksKey generates a JSON Web Key (JWK) from the given private key bytes (in PEM format).
// It returns a pointer to the JSONWebKey or an error if the operation fails.
func createJwksKey(privKeyBytes []byte) (*jose.JSONWebKey, error) {
	keyBlock, _ := pem.Decode(privKeyBytes)
	if keyBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM block for private key")
	}

	cert, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	keyID, err := keyIDFromPublicKey(&cert.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key ID from public key: %w", err)
	}

	return &jose.JSONWebKey{
		Key:       &cert.PublicKey,
		KeyID:     keyID,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}, nil
}

// keyIDFromPublicKey derives a key ID non-reversibly from a public key.
// taken from: https://github.com/kubernetes/kubernetes/blob/master/pkg/serviceaccount/jwt.go
func keyIDFromPublicKey(publicKey interface{}) (string, error) {
	publicKeyDERBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to serialize public key to DER format: %v", err)
	}

	hasher := crypto.SHA256.New()
	hasher.Write(publicKeyDERBytes)
	publicKeyDERHash := hasher.Sum(nil)

	keyID := base64.RawURLEncoding.EncodeToString(publicKeyDERHash)

	return keyID, nil
}
