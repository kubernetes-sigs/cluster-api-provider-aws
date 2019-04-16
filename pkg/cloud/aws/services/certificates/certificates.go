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

package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1alpha1"
)

const (
	rsaKeySize     = 2048
	duration365d   = time.Hour * 24 * 365
	clusterCA      = "cluster-ca"
	etcdCA         = "etcd-ca"
	frontProxyCA   = "front-proxy-ca"
	serviceAccount = "service-account"
)

// NewPrivateKey creates an RSA private key
func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, rsaKeySize)
}

// AltNames contains the domain names and IP addresses that will be added
// to the API Server's x509 certificate SubAltNames field. The values will
// be passed directly to the x509.Certificate object.
type AltNames struct {
	DNSNames []string
	IPs      []net.IP
}

// Config contains the basic fields required for creating a certificate
type Config struct {
	CommonName   string
	Organization []string
	AltNames     AltNames
	Usages       []x509.ExtKeyUsage
}

// ReconcileCertificates generate certificates if none exists.
func (s *Service) ReconcileCertificates() error {
	s.scope.V(2).Info("Reconciling certificates", "cluster-name", s.scope.Cluster.Name, "cluster-namespace", s.scope.Cluster.Namespace)

	if !s.scope.ClusterConfig.CAKeyPair.HasCertAndKey() {
		s.scope.V(2).Info("Generating keypair for", "user", clusterCA)
		clusterCAKeyPair, err := generateCACert(&s.scope.ClusterConfig.CAKeyPair, clusterCA)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate certs for %q", clusterCA)
		}
		s.scope.ClusterConfig.CAKeyPair = clusterCAKeyPair
	}

	if !s.scope.ClusterConfig.EtcdCAKeyPair.HasCertAndKey() {
		s.scope.V(2).Info("Generating keypair", "user", etcdCA)
		etcdCAKeyPair, err := generateCACert(&s.scope.ClusterConfig.EtcdCAKeyPair, etcdCA)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate certs for %q", etcdCA)
		}
		s.scope.ClusterConfig.EtcdCAKeyPair = etcdCAKeyPair
	}
	if !s.scope.ClusterConfig.FrontProxyCAKeyPair.HasCertAndKey() {
		s.scope.V(2).Info("Generating keypair", "user", frontProxyCA)
		fpCAKeyPair, err := generateCACert(&s.scope.ClusterConfig.FrontProxyCAKeyPair, frontProxyCA)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate certs for %q", frontProxyCA)
		}
		s.scope.ClusterConfig.FrontProxyCAKeyPair = fpCAKeyPair
	}

	if !s.scope.ClusterConfig.SAKeyPair.HasCertAndKey() {
		s.scope.V(2).Info("Generating service account keys", "user", serviceAccount)
		saKeyPair, err := generateServiceAccountKeys(&s.scope.ClusterConfig.SAKeyPair, serviceAccount)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate keyPair for %q", serviceAccount)
		}
		s.scope.ClusterConfig.SAKeyPair = saKeyPair
	}
	return nil
}

func generateCACert(kp *v1alpha1.KeyPair, user string) (v1alpha1.KeyPair, error) {
	x509Cert, privKey, err := NewCertificateAuthority()
	if err != nil {
		return v1alpha1.KeyPair{}, errors.Wrapf(err, "failed to generate CA cert for %q", user)
	}
	if kp == nil {
		return v1alpha1.KeyPair{
			Cert: EncodeCertPEM(x509Cert),
			Key:  EncodePrivateKeyPEM(privKey),
		}, nil
	}
	kp.Cert = EncodeCertPEM(x509Cert)
	kp.Key = EncodePrivateKeyPEM(privKey)
	return *kp, nil
}

func generateServiceAccountKeys(kp *v1alpha1.KeyPair, user string) (v1alpha1.KeyPair, error) {
	saCreds, err := NewPrivateKey()
	if err != nil {
		return v1alpha1.KeyPair{}, errors.Wrapf(err, "failed to create service account public and private keys")
	}
	saPub, err := EncodePublicKeyPEM(&saCreds.PublicKey)
	if err != nil {
		return v1alpha1.KeyPair{}, errors.Wrapf(err, "failed to encode service account public key to PEM")
	}
	if kp == nil {
		return v1alpha1.KeyPair{
			Cert: saPub,
			Key:  EncodePrivateKeyPEM(saCreds),
		}, nil
	}
	kp.Cert = saPub
	kp.Key = EncodePrivateKeyPEM(saCreds)
	return *kp, nil
}

// NewSignedCert creates a signed certificate using the given CA certificate and key
func (cfg *Config) NewSignedCert(key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random integer for signed cerficate")
	}

	if len(cfg.CommonName) == 0 {
		return nil, errors.New("must specify a CommonName")
	}

	if len(cfg.Usages) == 0 {
		return nil, errors.New("must specify at least one ExtKeyUsage")
	}

	tmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		DNSNames:     cfg.AltNames.DNSNames,
		IPAddresses:  cfg.AltNames.IPs,
		SerialNumber: serial,
		NotBefore:    caCert.NotBefore,
		NotAfter:     time.Now().Add(duration365d).UTC(),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  cfg.Usages,
	}

	b, err := x509.CreateCertificate(rand.Reader, &tmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create signed certificate: %+v", tmpl)
	}

	return x509.ParseCertificate(b)
}

// NewCertificateAuthority creates new certificate and private key for the certificate authority
func NewCertificateAuthority() (*x509.Certificate, *rsa.PrivateKey, error) {
	key, err := NewPrivateKey()
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create private key")
	}

	cert, err := NewSelfSignedCACert(key)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create self-signed certificate")
	}

	return cert, key, nil
}

// NewSelfSignedCACert creates a CA certificate.
func NewSelfSignedCACert(key *rsa.PrivateKey) (*x509.Certificate, error) {
	cfg := Config{
		CommonName: "kubernetes",
	}

	now := time.Now().UTC()

	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		NotBefore:             now,
		NotAfter:              now.Add(duration365d * 10),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		MaxPathLenZero:        true,
		BasicConstraintsValid: true,
		MaxPathLen:            0,
		IsCA:                  true,
	}

	b, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create self signed CA certificate: %+v", tmpl)
	}

	return x509.ParseCertificate(b)
}

// NewKubeconfig creates a new Kubeconfig where endpoint is the ELB endpoint.
func NewKubeconfig(clusterName, endpoint string, caCert *x509.Certificate, caKey *rsa.PrivateKey) (*api.Config, error) {
	cfg := &Config{
		CommonName:   "kubernetes-admin",
		Organization: []string{"system:masters"},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	clientKey, err := NewPrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create private key")
	}

	clientCert, err := cfg.NewSignedCert(clientKey, caCert, caKey)
	if err != nil {
		return nil, errors.Wrap(err, "unable to sign certificate")
	}

	userName := "kubernetes-admin"
	contextName := fmt.Sprintf("%s@%s", userName, clusterName)

	return &api.Config{
		Clusters: map[string]*api.Cluster{
			clusterName: {
				Server:                   endpoint,
				CertificateAuthorityData: EncodeCertPEM(caCert),
			},
		},
		Contexts: map[string]*api.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			userName: {
				ClientKeyData:         EncodePrivateKeyPEM(clientKey),
				ClientCertificateData: EncodeCertPEM(clientCert),
			},
		},
		CurrentContext: contextName,
	}, nil
}

// EncodeCertPEM returns PEM-endcoded certificate data.
func EncodeCertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

// EncodePrivateKeyPEM returns PEM-encoded private key data.
func EncodePrivateKeyPEM(key *rsa.PrivateKey) []byte {
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	return pem.EncodeToMemory(&block)
}

// EncodePublicKeyPEM returns PEM-encoded public key data.
func EncodePublicKeyPEM(key *rsa.PublicKey) ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return []byte{}, err
	}
	block := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}
	return pem.EncodeToMemory(&block), nil
}

// DecodeCertPEM attempts to return a decoded certificate or nil
// if the encoded input does not contain a certificate.
func DecodeCertPEM(encoded []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(encoded)
	if block == nil {
		return nil, nil
	}

	return x509.ParseCertificate(block.Bytes)
}

// DecodePrivateKeyPEM attempts to return a decoded key or nil
// if the encoded input does not contain a private key.
func DecodePrivateKeyPEM(encoded []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(encoded)
	if block == nil {
		return nil, nil
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// GenerateCertificateHash returns the encoded sha256 hash for the certificate provided
func GenerateCertificateHash(encoded []byte) (string, error) {
	cert, err := DecodeCertPEM(encoded)
	if err != nil || cert == nil {
		return "", errors.Errorf("failed to parse PEM block containing the public key")
	}

	certHash := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
	return "sha256:" + strings.ToLower(hex.EncodeToString(certHash[:])), nil
}
