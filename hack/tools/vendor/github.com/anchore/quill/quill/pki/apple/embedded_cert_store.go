package apple

import (
	"bytes"
	"crypto/x509"
	"embed"
	"fmt"
	"path"

	"github.com/anchore/quill/quill/pki/certchain"
	"github.com/anchore/quill/quill/pki/load"
)

//go:generate go run ./internal/generate/

//go:embed certs
var content embed.FS

var _ certchain.Store = (*embeddedCertStore)(nil)

type embeddedCertStore struct {
	rootCerts         []*x509.Certificate
	rootPEMs          [][]byte
	intermediateCerts []*x509.Certificate
	intermediatePEMs  [][]byte
	certsByCN         map[string][]*x509.Certificate
}

const certDir = "certs"

var singleton *embeddedCertStore

func GetEmbeddedCertStore() certchain.Store {
	return singleton
}

func init() {
	var err error
	singleton, err = newEmbeddedCertStore()
	if err != nil {
		panic(err)
	}
}

func newEmbeddedCertStore() (*embeddedCertStore, error) {
	var err error

	store := &embeddedCertStore{}

	store.rootPEMs, err = getPEMs(path.Join(certDir, "root"))
	if err != nil {
		return nil, fmt.Errorf("unable to load root certificates: %w", err)
	}

	store.intermediatePEMs, err = getPEMs(path.Join(certDir, "intermediate"))
	if err != nil {
		return nil, fmt.Errorf("unable to load root certificates: %w", err)
	}

	store.rootCerts, err = load.CertificatesFromPEMs(store.rootPEMs)
	if err != nil {
		return nil, fmt.Errorf("unable to parse root certificates: %w", err)
	}

	store.intermediateCerts, err = load.CertificatesFromPEMs(store.intermediatePEMs)
	if err != nil {
		return nil, fmt.Errorf("unable to parse intermediate certificates: %w", err)
	}

	store.certsByCN = make(map[string][]*x509.Certificate)
	for _, cert := range store.intermediateCerts {
		store.certsByCN[cert.Subject.CommonName] = append(store.certsByCN[cert.Subject.CommonName], cert)
	}
	for _, cert := range store.rootCerts {
		store.certsByCN[cert.Subject.CommonName] = append(store.certsByCN[cert.Subject.CommonName], cert)
	}

	return store, nil
}

func (s embeddedCertStore) CertificatesByCN(commonName string) ([]*x509.Certificate, error) {
	return s.certsByCN[commonName], nil
}

func (s embeddedCertStore) RootPEMs() [][]byte {
	return s.rootPEMs
}

func (s embeddedCertStore) IntermediatePEMs() [][]byte {
	return s.intermediatePEMs
}

func getPEMs(certsDir string) ([][]byte, error) {
	files, err := content.ReadDir(certsDir)
	if err != nil {
		return nil, err
	}

	var result [][]byte
	for _, f := range files {
		// read the file contents
		b, err := content.ReadFile(path.Join(certsDir, f.Name()))
		if err != nil {
			return nil, err
		}

		// store the file contents
		result = append(result, bytes.TrimRight(b, "\n"))
	}
	return result, nil
}
