// Package iam provides a service for managing AWS IAM Identity Providers.
package iam

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/s3"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
// One alternative is to have a large list of functions from the ec2 client.
type Service struct {
	scope     scope.EC2Scope
	IAMClient iamiface.IAMAPI
}

// NewService returns a new service given the api clients.
func NewService(clusterScope scope.EC2Scope) *Service {
	iamClient := scope.NewIAMClient(clusterScope, clusterScope, clusterScope, clusterScope.InfraCluster())

	return &Service{
		scope:     clusterScope,
		IAMClient: iamClient,
	}
}

const (
	// TrustPolicyJSON is the data key for the configmap containing the oidc trust policy.
	TrustPolicyJSON = "trust-policy.json"

	// S3HostFormat format for the host format for s3 for the oidc provider.
	S3HostFormat = "s3.%s.amazonaws.com"

	// STSAWSAudience security token service url.
	STSAWSAudience = "sts.amazonaws.com"

	// TrustPolicyConfigMapName name of the configmap that contains the trust policy template.
	TrustPolicyConfigMapName = "boilerplate-oidc-trust-policy"

	// TrustPolicyConfigMapNamespace namespace the trust policy template is put into.
	TrustPolicyConfigMapNamespace = metav1.NamespaceDefault
)

var (
	// WhitespaceRegex defines a regex pattern for detecting whitespace.
	WhitespaceRegex = regexp.MustCompile(`(?m)[\t\n]`)
)

// ReconcileOIDCProvider replicates functionality already built into managed clusters by auto-deploying an oidc provider
// that trusts the clusters Service Account issuer.
// For more details see: https://github.com/aws/amazon-eks-pod-identity-webhook/blob/master/SELF_HOSTED_SETUP.md
// 1. create an oidc provider in aws which points to the s3 bucket
// 2. build openid discovery config and upload to S3 bucket
// 3. build JWKs document based on SA signer public key and upload to the s3 bucket.
// 4. add trust Policy config to awscluster status.
func (s *Service) ReconcileOIDCProvider(ctx context.Context) error {
	log := s.scope.GetLogger()
	if !s.scope.AssociateOIDCProvider() {
		log.V(4).Info("OIDC provider association is disabled, skipping reconciliation")
		return nil
	}

	log.Info("Reconciling OIDC Provider")

	if s.scope.Bucket() == nil {
		return errors.New("s3 bucket configuration required to associate OIDC provider")
	}

	if err := s.reconcileBucketContents(ctx); err != nil {
		return err
	}

	if err := s.reconcileIdentityProvider(ctx); err != nil {
		return err
	}

	conditions.MarkTrue(s.scope.InfraCluster(), infrav1.OIDCProviderReadyCondition)

	if err := s.reconcileTrustPolicyConfigMap(ctx); err != nil {
		return fmt.Errorf("failed to reconcile trust policy config map: %w", err)
	}

	return nil
}

// DeleteOIDCProvider will delete the IAM OIDC provider. Note: that the bucket is cleaned up in the s3 service.
func (s *Service) DeleteOIDCProvider(ctx context.Context) error {
	if !s.scope.AssociateOIDCProvider() {
		return nil
	}

	log := s.scope.GetLogger()
	log.Info("Deleting OIDC Provider")

	if s.scope.Bucket() != nil {
		if err := s.deleteBucketContents(ctx, s3.NewService(s.scope)); err != nil {
			return err
		}
	}

	if s.scope.OIDCProviderStatus() != nil {
		return deleteOIDCProvider(s.scope.OIDCProviderStatus().ARN, s.IAMClient)
	}

	return nil
}
