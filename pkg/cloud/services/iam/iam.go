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
// 2. pause until kubeconfig and cluster access is ready
// 3. build openid discovery config and upload to S3 bucket
// 4. copy Service Account public signing JWKs to the s3 bucket.
func (s *Service) ReconcileOIDCProvider(ctx context.Context) error {
	if !s.scope.AssociateOIDCProvider() {
		return nil
	}

	log := s.scope.GetLogger()
	log.Info("Associating OIDC Provider")

	if s.scope.Bucket() == nil {
		return errors.New("s3 bucket configuration required to associate oidc provider")
	}

	if err := s.reconcileIdentityProvider(ctx); err != nil {
		return err
	}

	// the following can only run with a working workload cluster, return nil until then
	_, ok := s.scope.InfraCluster().GetAnnotations()[scope.KubeconfigReadyAnnotation]
	if !ok {
		log.Info("Associating OIDC Provider paused, kubeconfig and workload cluster API access is not ready")
		return nil
	}

	log.Info("Associating OIDC Provider continuing, kubeconfig for the workload cluster is available")
	if err := s.reconcileBucketContents(ctx); err != nil {
		return err
	}

	if err := s.reconcileTrustPolicyConfigMap(ctx); err != nil {
		return fmt.Errorf("failed to reconcile trust policy config map: %w", err)
	}

	conditions.MarkTrue(s.scope.InfraCluster(), infrav1.OIDCProviderReadyCondition)

	return nil
}

// DeleteOIDCProvider will delete the IAM OIDC provider. Note: that the bucket is cleaned up in the s3 service.
func (s *Service) DeleteOIDCProvider(_ context.Context) error {
	if !s.scope.AssociateOIDCProvider() {
		return nil
	}

	log := s.scope.GetLogger()
	log.Info("Deleting OIDC Provider")

	if s.scope.Bucket() != nil {
		if err := s.deleteBucketContents(s3.NewService(s.scope)); err != nil {
			return err
		}
	}

	return deleteOIDCProvider(s.scope.OIDCProviderStatus().ARN, s.IAMClient)
}
