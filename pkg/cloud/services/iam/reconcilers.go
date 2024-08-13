package iam

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/converters"
)

// CreateOIDCProvider will create an OIDC provider in IAM and store the arn/trustpolicy on the cluster status.
func (s *Service) reconcileIdentityProvider(_ context.Context) error {
	s3Host := fmt.Sprintf(S3HostFormat, s.scope.Region())
	thumbprint, err := fetchRootCAThumbprint(s3Host, 443)
	if err != nil {
		return err
	}

	// e.g. s3-us-west-2.amazonaws.com/<bucketname>/<clustername>
	oidcURL := "https://" + path.Join(s3Host, s.scope.Bucket().Name, s.scope.Name())
	arn, err := findAndVerifyOIDCProvider(oidcURL, thumbprint, s.IAMClient)
	if err != nil {
		return err
	}

	// find and verify confirms it's in IAM but if the status is not set we still want update
	providerStatus := s.scope.OIDCProviderStatus()
	if providerStatus.ARN != "" && providerStatus.ARN == arn {
		return nil
	}

	if arn == "" {
		var tags []*iam.Tag
		tags = append(tags, &iam.Tag{
			Key:   aws.String(v1beta2.ClusterAWSCloudProviderTagKey(s.scope.Name())),
			Value: aws.String(string(v1beta2.ResourceLifecycleOwned)),
		})

		input := iam.CreateOpenIDConnectProviderInput{
			ClientIDList:   aws.StringSlice([]string{STSAWSAudience}),
			ThumbprintList: aws.StringSlice([]string{thumbprint}),
			Url:            aws.String(oidcURL),
			Tags:           tags,
		}
		provider, err := s.IAMClient.CreateOpenIDConnectProvider(&input)
		if err != nil {
			return errors.Wrap(err, "error creating provider")
		}
		arn = aws.StringValue(provider.OpenIDConnectProviderArn)
	}

	providerStatus.ARN = arn
	oidcTrustPolicy := buildOIDCTrustPolicy(providerStatus.ARN)
	policy, err := converters.IAMPolicyDocumentToJSON(oidcTrustPolicy)
	if err != nil {
		return errors.Wrap(err, "failed to parse IAM policy")
	}
	providerStatus.TrustPolicy = WhitespaceRegex.ReplaceAllString(policy, "")
	return s.scope.PatchObject()
}

// reconcileTrustPolicyConfigMap make sure the remote cluster has the config map of the trust policy, this enables
// the remote cluster to have everything it needs to create roles for services accounts.
func (s *Service) reconcileTrustPolicyConfigMap(ctx context.Context) error {
	remoteClient, err := s.scope.RemoteClient()
	if err != nil {
		return err
	}

	configMapRef := types.NamespacedName{
		Name:      TrustPolicyConfigMapName,
		Namespace: TrustPolicyConfigMapNamespace,
	}

	trustPolicyConfigMap := &corev1.ConfigMap{}
	err = remoteClient.Get(ctx, configMapRef, trustPolicyConfigMap)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("getting %s/%s config map: %w", TrustPolicyConfigMapNamespace, TrustPolicyConfigMapName, err)
	}

	policy, err := converters.IAMPolicyDocumentToJSON(buildOIDCTrustPolicy(s.scope.OIDCProviderStatus().ARN))
	if err != nil {
		return errors.Wrap(err, "failed to parse IAM policy")
	}

	if tp, ok := trustPolicyConfigMap.Data[TrustPolicyJSON]; ok && tp == policy {
		return nil // trust policy in the kube is the same as generated, don't update
	}

	trustPolicyConfigMap.Data = map[string]string{
		TrustPolicyJSON: policy,
	}

	if trustPolicyConfigMap.UID == "" {
		trustPolicyConfigMap.Name = TrustPolicyConfigMapName
		trustPolicyConfigMap.Namespace = TrustPolicyConfigMapNamespace
		return remoteClient.Create(ctx, trustPolicyConfigMap)
	}

	return remoteClient.Update(ctx, trustPolicyConfigMap)
}
