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

package iamauth

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/iam/api/v1beta1"
)

// ReconcileIAMAuthenticator is used to create the aws-iam-authenticator in a cluster.
func (s *Service) ReconcileIAMAuthenticator(ctx context.Context) error {
	s.scope.Info("Reconciling aws-iam-authenticator configuration", "cluster-name", s.scope.Name())

	accountID, err := s.getAccountID()
	if err != nil {
		return fmt.Errorf("getting account id: %w", err)
	}

	remoteClient, err := s.scope.RemoteClient()
	if err != nil {
		s.scope.Error(err, "getting client for remote cluster")
		return fmt.Errorf("getting client for remote cluster: %w", err)
	}

	authBackend, err := NewBackend(s.backend, remoteClient)
	if err != nil {
		return fmt.Errorf("getting aws-iam-authenticator backend: %w", err)
	}

	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/nodes%s", accountID, iamv1.DefaultNameSuffix)
	nodesRoleMapping := ekscontrolplanev1.RoleMapping{
		RoleARN: roleARN,
		KubernetesMapping: ekscontrolplanev1.KubernetesMapping{
			UserName: EC2NodeUserName,
			Groups:   NodeGroups,
		},
	}
	s.scope.V(2).Info("Mapping node IAM role", "iam-role", nodesRoleMapping.RoleARN, "user", nodesRoleMapping.UserName)
	if err := authBackend.MapRole(nodesRoleMapping); err != nil {
		return fmt.Errorf("mapping iam node role: %w", err)
	}

	s.scope.V(2).Info("Mapping additional IAM roles and users")
	iamCfg := s.scope.IAMAuthConfig()
	for _, roleMapping := range iamCfg.RoleMappings {
		s.scope.V(2).Info("Mapping IAM role", "iam-role", roleMapping.RoleARN, "user", roleMapping.UserName)
		if err := authBackend.MapRole(roleMapping); err != nil {
			return fmt.Errorf("mapping iam role: %w", err)
		}
	}

	for _, userMapping := range iamCfg.UserMappings {
		s.scope.V(2).Info("Mapping IAM user", "iam-user", userMapping.UserARN, "user", userMapping.UserName)
		if err := authBackend.MapUser(userMapping); err != nil {
			return fmt.Errorf("mapping iam user: %w", err)
		}
	}

	s.scope.Info("Reconciled aws-iam-authenticator configuration", "cluster-name", s.scope.KubernetesClusterName())

	return nil
}

func (s *Service) getAccountID() (string, error) {
	input := &sts.GetCallerIdentityInput{}

	out, err := s.STSClient.GetCallerIdentity(input)
	if err != nil {
		return "", errors.Wrap(err, "unable to get caller identity")
	}

	return aws.StringValue(out.Account), nil
}
