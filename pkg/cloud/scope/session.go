/*
Copyright 2019 The Kubernetes Authors.

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

package scope

import (
	"context"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	region         = "AWS_REGION"
	awsAccessKeyID = "AWS_ACCESS_KEY_ID"
	awsAccessKey   = "AWS_SECRET_ACCESS_KEY"
)

var (
	sessionCache sync.Map
)

func parseCredentialsSecret(config *aws.Config, secret *corev1.Secret) error {
	if v, ok := secret.Data[region]; ok {
		config.WithRegion(string(v))
	}

	keyID, ok := secret.Data[awsAccessKeyID]
	if !ok || len(keyID) == 0 {
		return errors.Errorf("missing %q", awsAccessKeyID)
	}

	keySecret, ok := secret.Data[awsAccessKey]
	if !ok || len(keySecret) == 0 {
		return errors.Errorf("missing %q", awsAccessKey)
	}

	config.WithCredentials(
		credentials.NewStaticCredentials(string(keyID), string(keySecret), ""),
	)

	return nil
}

func sessionFromCluster(ctx context.Context, c client.Client, cluster *infrav1.AWSCluster) (*session.Session, error) {
	sessionKey := path.Join(cluster.Spec.Region, cluster.Namespace, cluster.Spec.CredentialsSecretName)
	s, ok := sessionCache.Load(sessionKey)
	if ok {
		return s.(*session.Session), nil
	}

	// Create a new AWS configuration.
	config := aws.NewConfig().WithRegion(cluster.Spec.Region)

	// Gather credentials from the secret if one is specified.
	//
	// Note that if the credentials change, we won't know or get notified,
	// users are expected to restart the manager, or change the name of the secret and reference.
	if cluster.Spec.CredentialsSecretName != "" {
		secret := &corev1.Secret{}
		secretKey := types.NamespacedName{Namespace: cluster.Namespace, Name: cluster.Spec.CredentialsSecretName}
		if err := c.Get(ctx, secretKey, secret); err != nil {
			return nil, errors.Wrapf(err, "failed to get credentials secret %q for AWSCluster %q", secretKey, cluster.Name)
		}
		if err := parseCredentialsSecret(config, secret); err != nil {
			return nil, errors.Wrapf(err, "failed to parse credentials secret %q for AWSCluster %q", secretKey, cluster.Name)
		}
	}

	ns, err := session.NewSession(aws.NewConfig().WithRegion(region))
	if err != nil {
		return nil, err
	}

	sessionCache.Store(region, ns)
	return ns, nil
}
