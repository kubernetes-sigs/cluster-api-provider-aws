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

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/cluster-api/util/secret"

	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
)

func (s *Service) reconcileKubeconfig(ctx context.Context, cluster *eks.Cluster) error {
	s.scope.V(2).Info("Reconciling EKS kubeconfig for cluster: %s", *cluster.Name)

	clusterRef := types.NamespacedName{
		Name:      s.scope.Cluster.Name,
		Namespace: s.scope.Cluster.Namespace,
	}

	_, err := secret.GetFromNamespacedName(ctx, s.scope.Client, clusterRef, secret.Kubeconfig)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return errors.Wrap(err, "failed to get kubeconfig secret")
		}

		createErr := s.createKubeconfigSecret(
			ctx,
			cluster,
			&clusterRef,
		)
		if createErr != nil {
			return err
		}

	}

	//TODO: does this need cert rotation? I don't think so

	return nil
}

func (s *Service) createKubeconfigSecret(ctx context.Context, cluster *eks.Cluster, clusterRef *types.NamespacedName) error {
	controllerOwnerRef := *metav1.NewControllerRef(s.scope.ControlPlane, infrav1exp.GroupVersion.WithKind("AWSManagedControlPlane"))

	clusterName := *cluster.Name
	userName := fmt.Sprintf("%s-admin", clusterName)
	contextName := fmt.Sprintf("%s@%s", userName, clusterName)

	cfg := &api.Config{
		APIVersion: api.SchemeGroupVersion.Version,
		Clusters: map[string]*api.Cluster{
			clusterName: {
				Server:                   *cluster.Endpoint,
				CertificateAuthorityData: []byte(*cluster.CertificateAuthority.Data),
			},
		},
		Contexts: map[string]*api.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		CurrentContext: contextName,
	}

	execConfig := &api.ExecConfig{APIVersion: "client.authentication.k8s.io/v1alpha1"}
	switch s.scope.TokenMethod() {
	case infrav1exp.EKSTokenMethodIAMAuthenticator:
		execConfig.Command = "aws-iam-authenticator"
		execConfig.Args = []string{
			"token",
			"-i",
			clusterName,
		}
	case infrav1exp.EKSTokenMethodAWSCli:
		execConfig.Command = "aws"
		execConfig.Args = []string{
			"eks",
			"get-token",
			"--cluster-name",
			clusterName,
		}
	default:
		return fmt.Errorf("unknown token method %s", s.scope.TokenMethod())
	}
	cfg.AuthInfos = map[string]*api.AuthInfo{
		userName: {
			Exec: execConfig,
		},
	}

	out, err := clientcmd.Write(*cfg)
	if err != nil {
		return errors.Wrap(err, "failed to serialize config to yaml")
	}

	kubeconfigSecret := kubeconfig.GenerateSecretWithOwner(*clusterRef, out, controllerOwnerRef)
	if err := s.scope.Client.Create(ctx, kubeconfigSecret); err != nil {
		return errors.Wrap(err, "failed to create kubeconfig secret")
	}

	record.Eventf(s.scope.ControlPlane, "SucessfulCreateKubeconfig", "Created kubeconfig for cluster %q", s.scope.Name())
	return nil
}
