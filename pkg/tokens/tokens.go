/*
Copyright 2018 The Kubernetes Authors.

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

package tokens

import (
	"time"

	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
)

// NewBootstrap attempts to create a token with the given ID.
func NewBootstrap(client corev1.SecretsGetter, ttl time.Duration) (string, error) {
	token, err := bootstraputil.GenerateBootstrapToken()
	if err != nil {
		return "", errors.Wrap(err, "unable to generate bootstrap token")
	}

	substrs := bootstraputil.BootstrapTokenRegexp.FindStringSubmatch(token)
	if len(substrs) != 3 {
		return "", errors.Wrapf(err, "the bootstrap token %q was not of the form %q", token, bootstrapapi.BootstrapTokenPattern)
	}
	tokenID := substrs[1]
	tokenSecret := substrs[2]

	secretName := bootstraputil.BootstrapTokenSecretName(tokenID)
	secretToken := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: metav1.NamespaceSystem,
		},
		Type: bootstrapapi.SecretTypeBootstrapToken,
		Data: map[string][]byte{
			bootstrapapi.BootstrapTokenIDKey:               []byte(tokenID),
			bootstrapapi.BootstrapTokenSecretKey:           []byte(tokenSecret),
			bootstrapapi.BootstrapTokenExpirationKey:       []byte(time.Now().UTC().Add(ttl).Format(time.RFC3339)),
			bootstrapapi.BootstrapTokenUsageSigningKey:     []byte("true"),
			bootstrapapi.BootstrapTokenUsageAuthentication: []byte("true"),
			bootstrapapi.BootstrapTokenExtraGroupsKey:      []byte("system:bootstrappers:kubeadm:default-node-token"),
		},
	}

	_, err = client.Secrets(secretToken.ObjectMeta.Namespace).Create(secretToken)
	return token, err
}
