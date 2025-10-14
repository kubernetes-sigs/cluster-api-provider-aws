package sign

import (
	"fmt"

	cms "github.com/github/smimesign/ietf-cms"

	"github.com/anchore/quill/quill/macho"
	"github.com/anchore/quill/quill/pki"
)

func generateCMS(signingMaterial pki.SigningMaterial, cdBlob *macho.Blob) (*macho.Blob, error) {
	cdBlobBytes, err := cdBlob.Pack()
	if err != nil {
		return nil, err
	}

	var cmsBytes []byte
	if signingMaterial.Signer != nil {
		cmsBytes, err = signDetached(cdBlobBytes, signingMaterial)
		if err != nil {
			return nil, fmt.Errorf("unable to sign code directory: %w", err)
		}
	}

	blob := macho.NewBlob(macho.MagicBlobwrapper, cmsBytes)

	return &blob, nil
}

func signDetached(data []byte, signingMaterial pki.SigningMaterial) ([]byte, error) {
	sd, err := cms.NewSignedData(data)
	if err != nil {
		return nil, err
	}

	if err = sd.Sign(signingMaterial.Certs, signingMaterial.Signer); err != nil {
		return nil, err
	}

	sd.Detached()

	if signingMaterial.TimestampServer != "" {
		if err = sd.AddTimestamps(signingMaterial.TimestampServer); err != nil {
			return nil, fmt.Errorf("unable to add timestamps (RFC3161): %w", err)
		}
	}

	return sd.ToDER()
}
