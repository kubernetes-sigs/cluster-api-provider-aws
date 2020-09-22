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

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// EC2NodeUserName is the username required for EC2 nodes
	EC2NodeUserName = "system:node:{{EC2PrivateDNSName}}"

	configMapName = "aws-auth"
	configMapNS   = metav1.NamespaceSystem
)

var (
	// NodeGroups is the groups that are required for a node
	NodeGroups = []string{"system:bootstrappers", "system:nodes"}
)

type configMapBackend struct {
	client crclient.Client
}

func (b *configMapBackend) MapRole(roleARN string, groups []string, username string) error {
	_, err := b.getConfigMap()
	if err != nil {
		return fmt.Errorf("getting aws-iam-authenticator config map: %w", err)
	}

	return nil
}

func (b *configMapBackend) MapUser(userARN string, groups []string, username string) error {
	return nil
}

// func (b *configMapBackend) getMappedRoles() ([]RoleMapping, error) {
// 	return nil
// }

func (b *configMapBackend) getConfigMap() (*corev1.ConfigMap, error) {
	ctx := context.Background()

	configMapRef := types.NamespacedName{
		Name:      configMapName,
		Namespace: configMapNS,
	}

	authConfigMap := &corev1.ConfigMap{}

	err := b.client.Get(ctx, configMapRef, authConfigMap)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, fmt.Errorf("getting %s/%s config map: %w", configMapName, configMapNS, err)
	}

	return authConfigMap, nil
}
