package macho

import "encoding/binary"

const (
	// Magic numbers used by Code Signing
	MagicRequirement             Magic = 0xfade0c00 // single Requirement blob
	MagicRequirements            Magic = 0xfade0c01 // Requirements vector (internal requirements)
	MagicCodedirectory           Magic = 0xfade0c02 // CodeDirectory blob
	MagicEmbeddedSignature       Magic = 0xfade0cc0 // embedded form of signature data
	MagicEmbeddedSignatureOld    Magic = 0xfade0b02 /* XXX */
	MagicLibraryDependencyBlob   Magic = 0xfade0c05
	MagicEmbeddedEntitlements    Magic = 0xfade7171 /* embedded entitlements */
	MagicEmbeddedEntitlementsDer Magic = 0xfade7172 /* embedded entitlements */
	MagicDetachedSignature       Magic = 0xfade0cc1 // multi-arch collection of embedded signatures
	MagicBlobwrapper             Magic = 0xfade0b01 // used for the cms blob
)

type Magic uint32

var SigningOrder = binary.BigEndian
