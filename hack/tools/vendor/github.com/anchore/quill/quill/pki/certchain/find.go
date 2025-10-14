package certchain

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"

	"github.com/scylladb/go-set/strset"

	"github.com/anchore/quill/internal/log"
)

// Find will look for the full certificate chain for the given certificate from the given cert store.
//
// Implementation note: this searches for a set of candidates to include in the chain, however, this does NOT
// follow strict binding rules. Instead, the best candidate matches are found based on the key ID followed by CN.
//
// If this function were following strict binding rules the key identifier, issuer name, and certificate
// serial numbers must match (which again, this function is NOT doing).
//
// Specifically: "Issuer certificate must match all these values in the Subject Key Identifier
// (SKI) extension, Subject and Serial Number fields respectively. In other words: KeyID value in the particular
// certificate AKI extension must match the value in the Subject Key Identifier (SKI) extension of the
// issuer certificate. Certificate Issuer value in the particular certificate must match the value in the Subject
// field of the issuer certificate. And Serial Number in the particular certificate must match the value in the
// Serial Number of the issuer certificate. If one of them doesn't match, certificate binding will fail and CCE
// will attempt to find another certificate that can be considered as a particular certificate issuer."
//
// source: https://www.sysadmins.lv/blog-en/certificate-chaining-engine-how-this-works.aspx
func Find(store Store, cert *x509.Certificate) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate

	if cert == nil {
		return nil, fmt.Errorf("not certificate provided")
	}

	visitedCerts := strset.New()
	nextKeyIDs := strset.New()
	if len(cert.SubjectKeyId) > 0 {
		nextKeyIDs.Add(hex.EncodeToString(cert.AuthorityKeyId))
	}
	nextCNs := strset.New(cert.Issuer.CommonName)

	for !nextCNs.IsEmpty() {
		parentCN := nextCNs.Pop()

		log.WithFields("cn", fmt.Sprintf("%q", parentCN)).Debug("querying certificate store")

		parentCerts, err := store.CertificatesByCN(parentCN)
		if err != nil {
			return nil, fmt.Errorf("unable to get certificate chain cert CN=%q (from keychain or embedded quill store): %w", parentCN, err)
		}

		if len(parentCerts) == 0 {
			// the search has been working OK, we just don't have any more certs to find. In this case we want to
			// error out to indicate the full chain was not found, but we also want to return the certs we have found
			// so far (show your work!).
			return certs, fmt.Errorf("no certificates found for CN=%q", parentCN)
		}

		log.WithFields("cn", fmt.Sprintf("%q", parentCN), "count", len(parentCerts)).Trace("certificates found")

		for _, c := range parentCerts {
			currentKeyID := hex.EncodeToString(c.SubjectKeyId)

			log.WithFields("cn", fmt.Sprintf("%q", c.Issuer.CommonName), "key-id", currentKeyID).Trace("capturing certificate in chain")

			if visitedCerts.Has(string(c.Raw)) {
				continue
			}

			visitedCerts.Add(string(c.Raw))

			// use the key ID of this parent... if it exists, use it over just the CN match
			if len(c.SubjectKeyId) > 0 {
				if nextKeyIDs.Has(currentKeyID) {
					// the key ID matches, so use it!
					nextKeyIDs.Remove(currentKeyID)
					certs = append(certs, c)
					nextCNs.Add(c.Issuer.CommonName)
				}

				if len(c.AuthorityKeyId) > 0 {
					nextKeyIDs.Add(hex.EncodeToString(c.AuthorityKeyId))
				}

				// we might have found a match by ID, maybe not... either way, there was an ID we could
				// have matched on, so at this point that should have been used for a match. Don't attempt
				// to match purely on CN.
				continue
			}

			// no key ID to use... but the CNs match so use it anyway (this is a "best" match)
			certs = append(certs, c)
			nextCNs.Add(c.Issuer.CommonName)
		}
	}
	return certs, nil
}
