package sign

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/go-restruct/restruct"

	"github.com/anchore/quill/quill/macho"
	"github.com/anchore/quill/quill/pki"
)

func GenerateSigningSuperBlob(id string, m *macho.File, signingMaterial pki.SigningMaterial, paddingTarget int) (int, []byte, error) {
	var cdFlags macho.CdFlag
	if signingMaterial.Signer != nil {
		// TODO: add options to enable more strict rules (such as macho.Hard)
		// note: we must at least support the runtime option for notarization (requirement introduced in macOS 10.14 / Mojave).
		// cdFlags = macho.Runtime | macho.Hard
		cdFlags = macho.Runtime
	} else {
		cdFlags = macho.Adhoc
	}

	requirementsBlob, requirementsHashBytes, err := generateRequirements(id, sha256.New(), signingMaterial)
	if err != nil {
		return 0, nil, fmt.Errorf("unable to create requirements: %w", err)
	}

	// TODO: add entitlements, for the meantime, don't include it
	entitlementsHashBytes := bytes.Repeat([]byte{0}, sha256.New().Size())

	cdBlob, err := generateCodeDirectory(id, sha256.New(), m, cdFlags, requirementsHashBytes, entitlementsHashBytes)
	if err != nil {
		return 0, nil, fmt.Errorf("unable to create code directory: %w", err)
	}

	cmsBlob, err := generateCMS(signingMaterial, cdBlob)
	if err != nil {
		return 0, nil, fmt.Errorf("unable to create signature block: %w", err)
	}

	sb := macho.NewSuperBlob(macho.MagicEmbeddedSignature)

	sb.Add(macho.CsSlotCodedirectory, cdBlob)
	sb.Add(macho.CsSlotRequirements, requirementsBlob)
	sb.Add(macho.CsSlotCmsSignature, cmsBlob)

	sb.Finalize(paddingTarget)

	sbBytes, err := restruct.Pack(macho.SigningOrder, &sb)
	if err != nil {
		return 0, nil, fmt.Errorf("unable to encode super blob: %w", err)
	}

	return int(sb.Length), sbBytes, nil
}

func UpdateSuperBlobOffsetReferences(m *macho.File, numSbBytes uint64) error {
	// (patch) patch  LcCodeSignature loader referencing the superblob offset
	if err := m.UpdateCodeSigningCmdDataSize(int(numSbBytes)); err != nil {
		return fmt.Errorf("unable to update code signature loader data size: %w", err)
	}

	// (patch) update the __LINKEDIT segment sizes to be "oldsize + newsuperblobsize"
	linkEditSegment := m.Segment("__LINKEDIT")

	linkEditSegment.Filesz += numSbBytes
	for linkEditSegment.Filesz > linkEditSegment.Memsz {
		linkEditSegment.Memsz *= 2
	}
	if err := m.UpdateSegmentHeader(linkEditSegment.SegmentHeader); err != nil {
		return fmt.Errorf("failed to update linkedit segment size: %w", err)
	}
	return nil
}
