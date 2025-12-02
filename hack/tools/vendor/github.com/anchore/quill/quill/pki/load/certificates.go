package load

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/anchore/quill/internal/log"
)

func Certificates(path string) ([]*x509.Certificate, error) {
	log.WithFields("path", path).Trace("reading certificate(s)")
	certPEM, err := BytesFromFileOrEnv(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read signing certificate: %w", err)
	}

	chainBlockBytes := decodeChainFromPEM(certPEM)

	if len(chainBlockBytes) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}

	var certs []*x509.Certificate

	for i, certBytes := range chainBlockBytes {
		c, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, fmt.Errorf("unable to parse certificate %d of %d: %w", i+1, len(chainBlockBytes), err)
		}

		certs = append(certs, c)
	}

	return certs, nil
}

func CertificatesFromPEMs(pems [][]byte) ([]*x509.Certificate, error) {
	var result []*x509.Certificate
	for _, pemBytes := range pems {
		certs, err := CertificatesFromPEM(pemBytes)
		if err != nil {
			return nil, err
		}
		result = append(result, certs...)
	}
	return result, nil
}

func CertificatesFromPEM(pemBytes []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	chainBlockBytes := decodeChainFromPEM(pemBytes)

	if len(chainBlockBytes) == 0 {
		return nil, fmt.Errorf("no PEM blocks found")
	}

	for i, certBytes := range chainBlockBytes {
		c, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, fmt.Errorf("unable to parse certificate %d of %d: %w", i+1, len(chainBlockBytes), err)
		}
		certs = append(certs, c)
	}

	return certs, nil
}

func decodeChainFromPEM(certInput []byte) (blocks [][]byte) {
	var certDERBlock *pem.Block
	for {
		certDERBlock, certInput = pem.Decode(certInput)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			blocks = append(blocks, certDERBlock.Bytes)
		}
	}
	return blocks
}
