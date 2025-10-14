package certchain

import (
	"crypto/x509"
	"fmt"

	"github.com/anchore/quill/quill/pki/load"
)

var _ Store = (*Collection)(nil)

type Collection struct {
	enumerators      []Enumerator
	searchers        []Searcher
	rootPEMs         [][]byte
	intermediatePEMs [][]byte
	certsByCN        map[string][]*x509.Certificate
}

func NewCollection() *Collection {
	return &Collection{}
}

func (p *Collection) AddRoot(certs ...*x509.Certificate) error {
	if err := p.add(certs...); err != nil {
		return err
	}
	pems, err := ToPEMs(certs...)
	if err != nil {
		return fmt.Errorf("unable to convert certificates to PEMs: %w", err)
	}

	p.rootPEMs = append(p.rootPEMs, pems...)
	return nil
}

func (p *Collection) AddIntermediate(certs ...*x509.Certificate) error {
	if err := p.add(certs...); err != nil {
		return err
	}
	pems, err := ToPEMs(certs...)
	if err != nil {
		return fmt.Errorf("unable to convert certificates to PEMs: %w", err)
	}

	p.intermediatePEMs = append(p.intermediatePEMs, pems...)
	return nil
}

func (p *Collection) AddRootPEMs(pems ...[]byte) error {
	if err := p.addPEMs(pems...); err != nil {
		return err
	}
	p.rootPEMs = append(p.rootPEMs, pems...)
	return nil
}

func (p *Collection) AddIntermediatePEMs(pems ...[]byte) error {
	if err := p.addPEMs(pems...); err != nil {
		return err
	}
	p.intermediatePEMs = append(p.intermediatePEMs, pems...)
	return nil
}

func (p *Collection) addPEMs(pems ...[]byte) error {
	certs, err := load.CertificatesFromPEMs(pems)
	if err != nil {
		return err
	}
	return p.add(certs...)
}

func (p *Collection) add(certs ...*x509.Certificate) error {
	if len(certs) == 0 {
		return nil
	}

	if p.certsByCN == nil {
		p.certsByCN = make(map[string][]*x509.Certificate)
	}

	for _, cert := range certs {
		p.certsByCN[cert.Subject.CommonName] = append(p.certsByCN[cert.Subject.CommonName], cert)
	}
	return nil
}

func (p *Collection) WithEnumerator(enumerators ...Enumerator) *Collection {
	p.enumerators = append(p.enumerators, enumerators...)
	return p
}

func (p *Collection) WithSearchers(searchers ...Searcher) *Collection {
	p.searchers = append(p.searchers, searchers...)
	return p
}

func (p *Collection) WithStores(stores ...Store) *Collection {
	for _, store := range stores {
		p.enumerators = append(p.enumerators, store)
		p.searchers = append(p.searchers, store)
	}
	return p
}

func (p *Collection) RootPEMs() [][]byte {
	var rootPEMs [][]byte
	for _, enumerator := range p.enumerators {
		rootPEMs = append(rootPEMs, enumerator.RootPEMs()...)
	}
	rootPEMs = append(rootPEMs, p.rootPEMs...)
	return rootPEMs
}

func (p *Collection) IntermediatePEMs() [][]byte {
	var intermediatePEMs [][]byte
	for _, enumerator := range p.enumerators {
		intermediatePEMs = append(intermediatePEMs, enumerator.IntermediatePEMs()...)
	}
	intermediatePEMs = append(intermediatePEMs, p.intermediatePEMs...)
	return intermediatePEMs
}

func (p *Collection) CertificatesByCN(commonName string) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for _, searcher := range p.searchers {
		results, err := searcher.CertificatesByCN(commonName)
		if err != nil {
			return nil, err
		}
		certs = append(certs, results...)
	}
	res, exists := p.certsByCN[commonName]
	if exists {
		certs = append(certs, res...)
	}
	return certs, nil
}
