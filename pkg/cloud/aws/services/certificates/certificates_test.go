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

package certificates

import (
	"crypto/x509"
	"testing"

	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

func TestGenerateCACert(t *testing.T) {
	testCases := []struct {
		name             string
		inputKeyPair     *v1alpha1.KeyPair
		inputUser        string
		expectKeyPairGen bool
		expectedError    error
	}{
		{
			name:             "should generate keypair when inputKeyPair==nil",
			inputKeyPair:     nil,
			inputUser:        "foo-ca",
			expectKeyPairGen: true,
			expectedError:    nil,
		},
		{
			name:             "should generate keypair when inputKeyPair has no cert",
			inputKeyPair:     &v1alpha1.KeyPair{Key: []byte("foo-key")},
			inputUser:        "foo-ca",
			expectKeyPairGen: true,
			expectedError:    nil,
		},
		{
			name:             "should generate keypair when inputKeyPair has no key",
			inputKeyPair:     &v1alpha1.KeyPair{Cert: []byte("foo-cert")},
			inputUser:        "foo-ca",
			expectKeyPairGen: true,
			expectedError:    nil,
		},
		{
			name:             "should generate keypair when inputKeyPair has no cert and nokey",
			inputKeyPair:     &v1alpha1.KeyPair{},
			inputUser:        "foo-ca",
			expectKeyPairGen: true,
			expectedError:    nil,
		},
		{
			name:             "should not generate keypair when inputKeyPair has cert and key",
			inputKeyPair:     &v1alpha1.KeyPair{Cert: []byte("foo-cert"), Key: []byte("foo-key")},
			inputUser:        "foo-ca",
			expectKeyPairGen: false,
			expectedError:    nil,
		},
	}
	for _, tc := range testCases {
		actualKeyPair, actualError := generateCACert(tc.inputKeyPair, tc.inputUser)
		if tc.expectedError != nil {
			if tc.expectedError.Error() != actualError.Error() {
				t.Fatalf("[%s], Unexpected error, Want [%v], Got: [%v]", tc.name, tc.expectedError, actualError)
			}
			continue
		}
		if !tc.expectKeyPairGen {
			if len(tc.inputKeyPair.Cert) != len(actualKeyPair.Cert) || string(tc.inputKeyPair.Cert) != string(actualKeyPair.Cert) {
				t.Fatalf("[%s] Want cert=%q, Got cert=%q", tc.name, string(tc.inputKeyPair.Cert), string(actualKeyPair.Cert))
			}
			if len(tc.inputKeyPair.Key) != len(actualKeyPair.Key) || string(tc.inputKeyPair.Key) != string(actualKeyPair.Key) {
				t.Fatalf("[%s] Want key=%q, Got key=%q", tc.name, string(tc.inputKeyPair.Key), string(actualKeyPair.Key))
			}
		} else {
			_, decodeErr := DecodeCertPEM(actualKeyPair.Cert)
			if decodeErr != nil {
				t.Fatalf("[%s], Expected to decode generated cert, Got decode failure %v", tc.name, decodeErr)
			}
			_, decodeErr = DecodePrivateKeyPEM(actualKeyPair.Key)
			if decodeErr != nil {
				t.Fatalf("[%s], Expected to decode generated private key, Got decode failure failed %v", tc.name, decodeErr)
			}
		}
	}
}

func TestGenerateServiceAccountKeys(t *testing.T) {
	testCases := []struct {
		name             string
		inputKeyPair     *v1alpha1.KeyPair
		inputUser        string
		expectKeyPairGen bool
		expectedError    error
	}{
		{
			name:             "should generate keypair when inputKeyPair==nil",
			inputKeyPair:     nil,
			inputUser:        "foo-sa",
			expectKeyPairGen: true,
			expectedError:    nil,
		},
		{
			name:             "should generate keypair when inputKeyPair has no cert",
			inputKeyPair:     &v1alpha1.KeyPair{Key: []byte("foo-key")},
			inputUser:        "foo-sa",
			expectKeyPairGen: true,
			expectedError:    nil,
		},
		{
			name:             "should generate keypair when inputKeyPair has no key",
			inputKeyPair:     &v1alpha1.KeyPair{Cert: []byte("foo-cert")},
			inputUser:        "foo-sa",
			expectKeyPairGen: true,
			expectedError:    nil,
		},
		{
			name:             "should generate keypair when inputKeyPair has no cert and nokey",
			inputKeyPair:     &v1alpha1.KeyPair{},
			inputUser:        "foo-ca",
			expectKeyPairGen: true,
			expectedError:    nil,
		},
		{
			name:             "should not generate keypair when inputKeyPair has cert and key",
			inputKeyPair:     &v1alpha1.KeyPair{Cert: []byte("foo-cert"), Key: []byte("foo-key")},
			inputUser:        "foo-sa",
			expectKeyPairGen: false,
			expectedError:    nil,
		},
	}
	for _, tc := range testCases {
		actualKeyPair, actualError := generateServiceAccountKeys(tc.inputKeyPair, tc.inputUser)
		if tc.expectedError != nil {
			if tc.expectedError.Error() != actualError.Error() {
				t.Fatalf("[%s], Unexpected error, Want [%v], Got: [%v]", tc.name, tc.expectedError, actualError)
			}
			continue
		}
		if !tc.expectKeyPairGen {
			if len(tc.inputKeyPair.Cert) != len(actualKeyPair.Cert) || string(tc.inputKeyPair.Cert) != string(actualKeyPair.Cert) {
				t.Fatalf("[%s] Want cert=%q, Got cert=%q", tc.name, string(tc.inputKeyPair.Cert), string(actualKeyPair.Cert))
			}
			if len(tc.inputKeyPair.Key) != len(actualKeyPair.Key) || string(tc.inputKeyPair.Key) != string(actualKeyPair.Key) {
				t.Fatalf("[%s] Want key=%q, Got key=%q", tc.name, string(tc.inputKeyPair.Key), string(actualKeyPair.Key))
			}
		} else {
			_, decodeErr := DecodePrivateKeyPEM(actualKeyPair.Key)
			if decodeErr != nil {
				t.Fatalf("[%s], Expected to decode generated private key, Got decode failure failed %v", tc.name, decodeErr)
			}

			// TODO: find a stronger check
			if len(actualKeyPair.Key) <= 0 {
				t.Fatalf("[%s], Expected to public key of length > 0, Got public key of length %d", tc.name, len(actualKeyPair.Key))
			}
		}
	}
}

func TestNewCertificateAuthority(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func() error
	}{
		{
			name: "should generate conformant ca certificate",
			testFunc: func() error {
				cert, _, err := NewCertificateAuthority()
				if err != nil {
					return err
				}

				if !cert.MaxPathLenZero {
					return errors.Errorf("Unexpected value for MaxPathLenZero, Want: [true]; Got: [%t]", cert.MaxPathLenZero)
				}

				if cert.MaxPathLen != 0 {
					return errors.Errorf("Unexpected value for MaxPathLen, Want: [0]; Got: [%d]", cert.MaxPathLen)
				}

				if !cert.BasicConstraintsValid {
					return errors.Errorf("Unexpected value for BasicConstraintsValid, Want: [true]; Got: [%t]", cert.BasicConstraintsValid)
				}

				expectedUsage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
				actualUsage := cert.KeyUsage
				if expectedUsage != actualUsage {
					return errors.Errorf("Unexpected value for KeyUsage, Want: [%d]; Got: [%d]", expectedUsage, actualUsage)
				}

				expectedCommonName := "kubernetes"
				actualCommonName := cert.Subject.CommonName
				if expectedCommonName != actualCommonName {
					return errors.Errorf("Unexpected CommonName, Want: [%s]; Got: [%s]", expectedCommonName, actualCommonName)
				}

				if !cert.IsCA {
					return errors.Errorf("Unexpected value for IsCA, Want: [true]; Got: [%t]", cert.IsCA)
				}

				return nil
			},
		},
	}

	for _, tc := range testCases {
		if err := tc.testFunc(); err != nil {
			t.Fatalf("[%s] failed: %v", tc.name, err)
		}
	}
}
