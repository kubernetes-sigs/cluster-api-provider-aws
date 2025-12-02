package certchain

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
)

func ToPEM(certs ...*x509.Certificate) ([]byte, error) {
	var pemBytes bytes.Buffer
	for _, cert := range certs {
		if err := pem.Encode(&pemBytes, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}); err != nil {
			return nil, err
		}
	}
	return pemBytes.Bytes(), nil
}

func ToPEMs(certs ...*x509.Certificate) ([][]byte, error) {
	var pems [][]byte
	for _, cert := range certs {
		p, err := ToPEM(cert)
		if err != nil {
			return nil, err
		}
		pems = append(pems, p)
	}
	return pems, nil
}
