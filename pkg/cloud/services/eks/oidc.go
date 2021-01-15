/*
Copyright 2020 The Kubernetes Authors.

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

package eks

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	apiiam "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/iam/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/converters"
	"sigs.k8s.io/cluster-api/controllers/remote"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	trustPolicyConfigMapName      = "boilerplate-oidc-trust-policy"
	trustPolicyConfigMapNamespace = metav1.NamespaceDefault

	whitespaceRe = regexp.MustCompile(`(?m)[\t\n]`)
)

func (s *Service) reconcileOIDCProvider(cluster *eks.Cluster) error {
	if !s.scope.ControlPlane.Spec.AssociateOIDCProvider || s.scope.ControlPlane.Status.OIDCProvider.ARN != "" {
		return nil
	}

	if !s.scope.EnableIAM() {
		return errors.New("'AssociateOIDCProvider' provided without enabling the 'EKSEnableIAM' feature flag")
	}

	s.scope.Info("Reconciling EKS OIDC Provider", "cluster-name", cluster.Name)
	oidcProvider, err := s.CreateOIDCProvider(cluster)
	if err != nil {
		return errors.Wrap(err, "failed to create OIDC provider")
	}
	s.scope.ControlPlane.Status.OIDCProvider.ARN = oidcProvider

	policy, err := converters.IAMPolicyDocumentToJSON(s.buildOIDCTrustPolicy())
	if err != nil {
		return errors.Wrap(err, "failed to parse IAM policy")
	}
	s.scope.ControlPlane.Status.OIDCProvider.TrustPolicy = whitespaceRe.ReplaceAllString(policy, "")
	if err := s.scope.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to update control plane with OIDC provider ARN")
	}

	if err := s.reconcileTrustPolicy(); err != nil {
		return errors.Wrap(err, "failed to reconcile trust policy in workload cluster")
	}

	return nil
}

func (s *Service) reconcileTrustPolicy() error {
	ctx := context.Background()

	clusterKey := client.ObjectKey{
		Name:      s.scope.Name(),
		Namespace: s.scope.Namespace(),
	}

	restConfig, err := remote.RESTConfig(ctx, s.scope.Client, clusterKey)
	if err != nil {
		return fmt.Errorf("getting remote client for %s/%s: %w", s.scope.Namespace(), s.scope.Name(), err)
	}

	remoteClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		return fmt.Errorf("getting client for remote cluster: %w", err)
	}

	configMapRef := types.NamespacedName{
		Name:      trustPolicyConfigMapName,
		Namespace: trustPolicyConfigMapNamespace,
	}

	trustPolicyConfigMap := &corev1.ConfigMap{}

	err = remoteClient.Get(ctx, configMapRef, trustPolicyConfigMap)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("getting %s/%s config map: %w", trustPolicyConfigMapNamespace, trustPolicyConfigMapName, err)
	}

	policy, err := converters.IAMPolicyDocumentToJSON(s.buildOIDCTrustPolicy())
	if err != nil {
		return errors.Wrap(err, "failed to parse IAM policy")
	}

	trustPolicyConfigMap.Data = map[string]string{
		"trust-policy.json": policy,
	}

	if trustPolicyConfigMap.UID == "" {
		trustPolicyConfigMap.Name = trustPolicyConfigMapName
		trustPolicyConfigMap.Namespace = trustPolicyConfigMapNamespace
		s.V(2).Info("Creating new Trust Policy ConfigMap", "cluster", s.scope.Name(), "configmap", trustPolicyConfigMapName)
		return remoteClient.Create(ctx, trustPolicyConfigMap)
	}

	s.V(2).Info("Updating existing Trust Policy ConfigMap", "cluster", s.scope.Name(), "configmap", trustPolicyConfigMapName)
	return remoteClient.Update(ctx, trustPolicyConfigMap)
}

func (s *Service) deleteOIDCProvider() error {
	if !s.scope.ControlPlane.Spec.AssociateOIDCProvider || s.scope.ControlPlane.Status.OIDCProvider.ARN == "" {
		return nil
	}

	providerARN := s.scope.ControlPlane.Status.OIDCProvider.ARN
	if err := s.DeleteOIDCProvider(&providerARN); err != nil {
		return errors.Wrap(err, "failed to delete OIDC provider")
	}

	s.scope.ControlPlane.Status.OIDCProvider.ARN = ""
	if err := s.scope.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to update control plane with OIDC provider ARN")
	}

	return nil
}

func (s *Service) buildOIDCTrustPolicy() apiiam.PolicyDocument {
	providerARN := s.scope.ControlPlane.Status.OIDCProvider.ARN
	conditionValue := providerARN[strings.Index(providerARN, "/")+1:] + ":sub"

	return apiiam.PolicyDocument{
		Version: "2012-10-17",
		Statement: apiiam.Statements{
			apiiam.StatementEntry{
				Sid:    "",
				Effect: "Allow",
				Principal: apiiam.Principals{
					apiiam.PrincipalFederated: apiiam.PrincipalID{providerARN},
				},
				Action: apiiam.Actions{"sts:AssumeRoleWithWebIdentity"},
				Condition: apiiam.Conditions{
					"ForAnyValue:StringLike": map[string][]string{
						conditionValue: {"system:serviceaccount:${SERVICE_ACCOUNT_NAMESPACE}:${SERVICE_ACCOUNT_NAME}"},
					},
				},
			},
		},
	}
}
