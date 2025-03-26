package load

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"

	"software.sslmate.com/src/go-pkcs12"

	"github.com/anchore/quill/internal/log"
)

var ErrNeedPassword = errors.New("need password to decode file")

type P12Contents struct {
	PrivateKey   crypto.PrivateKey
	Certificate  *x509.Certificate
	Certificates []*x509.Certificate
}

func P12(path, password string) (*P12Contents, error) {
	by, err := BytesFromFileOrEnv(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read p12 bytes: %w", err)
	}

	key, cert, certs, err := pkcs12.DecodeChain(by, password)
	if err != nil {
		if errors.Is(err, pkcs12.ErrIncorrectPassword) && password == "" {
			log.Debug("p12 file requires a password but none provided")
			return nil, ErrNeedPassword
		}
		return nil, fmt.Errorf("unable to decode p12 file: %w", err)
	}

	return &P12Contents{
		PrivateKey:   key,
		Certificate:  cert,
		Certificates: certs,
	}, nil
}
