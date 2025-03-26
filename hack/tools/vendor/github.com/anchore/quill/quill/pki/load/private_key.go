package load

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/anchore/quill/internal/log"
)

func PrivateKey(path string, password string) (crypto.PrivateKey, error) {
	log.Debug("loading private key")

	b, err := BytesFromFileOrEnv(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %w", err)
	}
	pemObj, _ := pem.Decode(b)

	if pemObj == nil {
		return nil, fmt.Errorf("unable to decode PEM formatted private key")
	}

	switch pemObj.Type {
	case "RSA PRIVATE KEY", "PRIVATE KEY", "ENCRYPTED PRIVATE KEY":
		// pass
	default:
		return nil, fmt.Errorf("RSA private key is of the wrong type: %q", pemObj.Type)
	}

	var privPemBytes []byte

	//nolint: staticcheck // we have no other alternatives
	if x509.IsEncryptedPEMBlock(pemObj) {
		// why is this deprecated?
		//	> "Legacy PEM encryption as specified in RFC 1423 is insecure by
		//  > design. Since it does not authenticate the ciphertext, it is vulnerable to
		//  > padding oracle attacks that can let an attacker recover the plaintext."
		//
		// This method of encrypting the key isn't recommended anymore.
		// See https://github.com/golang/go/issues/8860 for more discussion

		log.Trace("decrypting private key")

		//nolint: staticcheck // we have no other alternatives
		privPemBytes, err = x509.DecryptPEMBlock(pemObj, []byte(password))
		if err != nil {
			return nil, fmt.Errorf("unable to decrypt PEM block: %w", err)
		}
	} else {
		log.Trace("using unencrypted private key")
		privPemBytes = pemObj.Bytes
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil {
			return nil, fmt.Errorf("unable to parse RSA private key: %w", err)
		}
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unable to find RSA private key after parsing: %w", err)
	}

	return privateKey, nil
}
