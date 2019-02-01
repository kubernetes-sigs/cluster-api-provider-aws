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
	"testing"

	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/services/certificates"
)

func TestGetOrGenerateKeyPair(t *testing.T) {
	testCases := []struct {
		name                string
		inputKeyPair        *v1alpha1.KeyPair
		inputUser           string
		expectKeyPairGen    bool
		expectCACertKeyPair bool
		expectedError       error
	}{
		{
			name:                "should return generated \"cluster-ca\" keypair when inputKeyPair==nil",
			inputKeyPair:        nil,
			inputUser:           "cluster-ca",
			expectKeyPairGen:    true,
			expectCACertKeyPair: true,
			expectedError:       nil,
		},
		{
			name:                "should return generated \"cluster-ca\" keypair when inputKeyPair has no cert",
			inputKeyPair:        &v1alpha1.KeyPair{Key: []byte("foo-key")},
			inputUser:           "cluster-ca",
			expectKeyPairGen:    true,
			expectCACertKeyPair: true,
			expectedError:       nil,
		},
		{
			name:                "should return generated \"cluster-ca\" keypair when inputKeyPair has no key",
			inputKeyPair:        &v1alpha1.KeyPair{Cert: []byte("foo-cert")},
			inputUser:           "cluster-ca",
			expectKeyPairGen:    true,
			expectCACertKeyPair: true,
			expectedError:       nil,
		},
		{
			name:                "should generate \"cluster-ca\" keypair",
			inputKeyPair:        nil,
			inputUser:           "cluster-ca",
			expectKeyPairGen:    true,
			expectCACertKeyPair: true,
			expectedError:       nil,
		},
		{
			name:                "should generate \"etcd-ca\" keypair",
			inputKeyPair:        nil,
			inputUser:           "etcd-ca",
			expectKeyPairGen:    true,
			expectCACertKeyPair: true,
			expectedError:       nil,
		},
		{
			name:                "should generate \"front-proxy-ca\" keypair",
			inputKeyPair:        nil,
			inputUser:           "front-proxy-ca",
			expectKeyPairGen:    true,
			expectCACertKeyPair: true,
			expectedError:       nil,
		},
		{
			name:                "should generate \"service-account\" keypair",
			inputKeyPair:        nil,
			inputUser:           "service-account",
			expectKeyPairGen:    true,
			expectCACertKeyPair: false,
			expectedError:       nil,
		},
		{
			name:                "should return error for unknown user",
			inputKeyPair:        nil,
			inputUser:           "foo-ca",
			expectKeyPairGen:    true,
			expectCACertKeyPair: false,
			expectedError:       errors.Errorf("Unknown user \"foo-ca\", skipping generating keyPair"),
		},
		{
			name:                "should not generate keypair when inputKeyPair has cert and key",
			inputKeyPair:        &v1alpha1.KeyPair{Cert: []byte("foo-cert"), Key: []byte("foo-key")},
			inputUser:           "cluster-ca",
			expectKeyPairGen:    false,
			expectCACertKeyPair: false,
			expectedError:       nil,
		},
	}

	for _, tc := range testCases {
		actualCert, actualKey, actualError := GetOrGenerateKeyPair(tc.inputKeyPair, tc.inputUser)
		if tc.expectedError != nil {
			if tc.expectedError.Error() != actualError.Error() {
				t.Fatalf("[%s], Unexpected error, Want [%v], Got: [%v]", tc.name, tc.expectedError, actualError)
			} else {
				continue
			}
		}

		if !tc.expectKeyPairGen {
			if len(tc.inputKeyPair.Cert) != len(actualCert) || string(tc.inputKeyPair.Cert) != string(actualCert) {
				t.Fatalf("[%s] Want cert=%q, Got cert=%q", tc.name, string(tc.inputKeyPair.Cert), string(actualCert))
			}
			if len(tc.inputKeyPair.Key) != len(actualKey) || string(tc.inputKeyPair.Key) != string(actualKey) {
				t.Fatalf("[%s] Want key=%q, Got key=%q", tc.name, string(tc.inputKeyPair.Key), string(actualKey))
			}
		} else {
			if tc.expectCACertKeyPair {
				_, decodeErr := certificates.DecodeCertPEM(actualCert)
				if decodeErr != nil {
					t.Fatalf("[%s], Expected to decode generated cert, Got decode failure %v", tc.name, decodeErr)
				}
				_, decodeErr = certificates.DecodePrivateKeyPEM(actualKey)
				if decodeErr != nil {
					t.Fatalf("[%s], Expected to decode generated private key, Got decode failure failed %v", tc.name, decodeErr)
				}
			} else {
				_, decodeErr := certificates.DecodePrivateKeyPEM(actualKey)
				if decodeErr != nil {
					t.Fatalf("[%s], Expected to decode generated private key, Got decode failure failed %v", tc.name, decodeErr)
				}

				// TODO: find a stronger check
				if len(actualCert) <= 0 {
					t.Fatalf("[%s], Expected to public key of length > 0, Got public key of lenght %d", tc.name, len(actualCert))
				}
			}
		}
	}
}
