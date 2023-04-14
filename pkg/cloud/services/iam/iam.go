package iam

import (
	"context"
	"errors"
	"regexp"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/s3"
)

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

	// PodIdentityWebhookCertificateFormat the format for the cert name used for the pod-identity-webhook.
	PodIdentityWebhookCertificateFormat = "%s-pod-id-wh"

	// SelfsignedIssuerFormat format for the self signed issuer used for the cluster to make the pod-identity-webhook cert.
	SelfsignedIssuerFormat = "%s-selfsigned-issuer"

	// S3HostFormat format for the host format for s3 for the oidc provider.
	S3HostFormat = "s3-%s.amazonaws.com"

	// STSAWSAudience security token service url.
	STSAWSAudience = "sts.amazonaws.com"

	// TrustPolicyConfigMapName name of the configmap that contains the trust policy template.
	TrustPolicyConfigMapName = "boilerplate-oidc-trust-policy"

	// TrustPolicyConfigMapNamespace namespace the trust policy template is put into.
	TrustPolicyConfigMapNamespace = metav1.NamespaceDefault
)

var (
	WhitespaceRegex = regexp.MustCompile(`(?m)[\t\n]`)
)

// ReconcileOIDCProvider replicates functionality already built into managed clusters by auto-deploying the
// modifying kube-apiserver args, deploying the pod identity webhook and setting/configuring an oidc provider
// for more details see: https://github.com/aws/amazon-eks-pod-identity-webhook/blob/master/SELF_HOSTED_SETUP.md
// 1. create a self signed issuer for the mutating webhook
// 2. add create a json patch for kube-apiserver and use capi config to add to the kubeadm.yml
// 3. create an oidc provider in aws which points to the s3 bucket
// 4. pause until kubeconfig and cluster acccess is ready
// 5. move openid config and jwks to the s3 bucket
// 6. add the pod identity webhook to the workload cluster
// 7. add the configmap to the workload cluster.
func (s *Service) ReconcileOIDCProvider(ctx context.Context) error {
	if !s.scope.AssociateOIDCProvider() {
		return nil
	}

	log := ctrl.LoggerFrom(ctx)
	log.Info("Associating OIDC Provider")

	if s.scope.Bucket() == nil {
		return errors.New("s3 bucket configuration required to associate oidc provider")
	}

	if err := s.reconcileSelfsignedIssuer(ctx); err != nil {
		return err
	}

	if err := s.reconcileKubeAPIParameters(ctx); err != nil {
		return err
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

	if err := s.reconcilePodIdentityWebhook(ctx); err != nil {
		return err
	}

	return s.reconcileTrustPolicyConfigMap(ctx)
}

// DeleteOIDCProvider will delete the iam resources note that the bucket is cleaned up in the s3 service
// 1. delete oidc provider
// 2. delete mwh certificate
// 3. delete cert-manager issuer.
func (s *Service) DeleteOIDCProvider(ctx context.Context) error {
	if !s.scope.AssociateOIDCProvider() {
		return nil
	}

	log := ctrl.LoggerFrom(ctx)
	log.Info("Deleting OIDC Provider")

	if s.scope.Bucket() != nil {
		if err := deleteBucketContents(s3.NewService(s.scope)); err != nil {
			return err
		}
	}

	if err := deleteCertificatesAndIssuer(ctx, s.scope.Name(), s.scope.Namespace(), s.scope.ManagementClient()); err != nil {
		return err
	}

	return deleteOIDCProvider(s.scope.OIDCProviderStatus().ARN, s.IAMClient)
}
