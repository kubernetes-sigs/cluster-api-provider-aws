package notary

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/anchore/quill/internal/log"
	"github.com/anchore/quill/quill/pki/load"
)

type TokenConfig struct {
	Issuer        string
	PrivateKeyID  string
	TokenLifetime time.Duration
	PrivateKey    string
}

func NewSignedToken(cfg TokenConfig) (string, error) {
	method := jwt.SigningMethodES256 // TODO: add more methods
	token := &jwt.Token{
		Header: map[string]interface{}{
			"alg": method.Alg(),
			"kid": cfg.PrivateKeyID,
			"typ": "JWT",
		},
		Claims: jwt.MapClaims{
			"iss":   cfg.Issuer,                                     // issuer ID from Apple
			"iat":   time.Now().UTC().Unix(),                        // token’s creation timestamp (unix epoch)
			"exp":   time.Now().Add(cfg.TokenLifetime).UTC().Unix(), // token’s expiration timestamp (unix epoch).
			"aud":   "appstoreconnect-v1",                           // audience
			"scope": []string{"/notary/v2"},                         // list of operations you want App Store Connect to allow for this token
		},
		Method: method,
	}

	key, err := loadPrivateKey(cfg.PrivateKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}

func loadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	log.Debug("loading private key for notary")

	keyBytes, err := load.BytesFromFileOrEnv(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load JWT private key bytes: %w", err)
	}

	key, err := jwt.ParseECPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse EC private key for JWT: %w", err)
	}
	return key, nil
}
