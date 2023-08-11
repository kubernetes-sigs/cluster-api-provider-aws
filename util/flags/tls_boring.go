//go:build boringcrypto

package flags

import "crypto/tls"

func InsecureSkipVerify(insecureSkipVerify bool) bool {
	return false
}

func GetTlsMaxVersion() uint16 {
	return tls.VersionTLS12
}
