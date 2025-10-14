package certchain

import "crypto/x509"

type Store interface {
	Enumerator
	Searcher
}

type Enumerator interface {
	RootPEMs() [][]byte
	IntermediatePEMs() [][]byte
}

type Searcher interface {
	CertificatesByCN(commonName string) ([]*x509.Certificate, error)
}
