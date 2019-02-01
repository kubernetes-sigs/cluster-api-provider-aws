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

package actuators

import (
	"github.com/pkg/errors"
	"k8s.io/klog"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/certificates"
)

const (
	ClusterCA      = "cluster-ca"
	EtcdCA         = "etcd-ca"
	FrontProxyCA   = "front-proxy-ca"
	ServiceAccount = "service-account"
)

// GetOrGenerateKeyPair returns a byte encoded cert and key pair if exists, generates one otherwise
func GetOrGenerateKeyPair(kp *v1alpha1.KeyPair, user string) ([]byte, []byte, error) {
	if kp == nil || !kp.HasCertAndKey() {
		klog.V(2).Infof("Generating key pair for %q", user)
		switch user {
		case EtcdCA, FrontProxyCA, ClusterCA:
			x509Cert, privKey, err := certificates.NewCertificateAuthority()
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to generate CA cert for %q", user)
			}
			return certificates.EncodeCertPEM(x509Cert), certificates.EncodePrivateKeyPEM(privKey), nil
		case ServiceAccount:
			saCreds, err := certificates.NewPrivateKey()
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to create service account public and private keys")
			}
			saPub, err := certificates.EncodePublicKeyPEM(&saCreds.PublicKey)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to encode service account public key to PEM")
			}

			return saPub, certificates.EncodePrivateKeyPEM(saCreds), nil
		default:
			return nil, nil, errors.Errorf("Unknown user %q, skipping generating keyPair", user)
		}
	}

	return kp.Cert, kp.Key, nil
}
