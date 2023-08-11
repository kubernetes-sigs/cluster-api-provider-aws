//go:build !boringcrypto

package flags

import "crypto/tls"

func InsecureSkipVerify(insecureSkipVerify bool) bool {
	return insecureSkipVerify
}

func GetTlsMaxVersion() uint16 {
	return tls.VersionTLS13
}
