package sign

import (
	"encoding/asn1"
	"fmt"
	"hash"
	"unsafe"

	"github.com/go-restruct/restruct"

	"github.com/anchore/quill/internal/log"
	"github.com/anchore/quill/quill/macho"
	"github.com/anchore/quill/quill/pki"
)

const (
	sizeOfUint32 = 4
)

type reqStatement []uint32

const (
	opFalse              uint32 = iota // unconditionally false
	opTrue                             // unconditionally true
	opIdent                            // match canonical code [string]
	opAppleAnchor                      // signed by Apple as Apple's product
	opAnchorHash                       // match anchor [cert hash]
	opInfoKeyValue                     // *legacy* - use opInfoKeyField [key; value]
	opAnd                              // binary prefix expr AND expr [expr; expr]
	opOr                               // binary prefix expr OR expr [expr; expr]
	opCDHash                           // match hash of CodeDirectory directly [cd hash]
	opNot                              // logical inverse [expr]
	opInfoKeyField                     // Info.plist key field [string; match suffix]
	opCertField                        // Certificate field [cert index; field name; match suffix]
	opTrustedCert                      // require trust settings to approve one particular cert [cert index]
	opTrustedCerts                     // require trust settings to approve the cert chain
	opCertGeneric                      // Certificate component by OID [cert index; oid; match suffix]
	opAppleGenericAnchor               // signed by Apple in any capacity
	opEntitlementField                 // entitlement dictionary field [string; match suffix]
	opCertPolicy                       // Certificate policy by OID [cert index; oid; match suffix]
	opNamedAnchor                      // named anchor type
	opNamedCode                        // named subroutine
	exprOpCount                        // (total opcode count in use)
)

const (
	matchExists       uint32 = iota // anything but explicit "false" - no value stored
	matchEqual                      // equal (CFEqual)
	matchContains                   // partial match (substring)
	matchBeginsWith                 // partial match (initial substring)
	matchEndsWith                   // partial match (terminal substring)
	matchLessThan                   // less than (string with numeric comparison)
	matchGreaterThan                // greater than (string with numeric comparison)
	matchLessEqual                  // less or equal (string with numeric comparison)
	matchGreaterEqual               // greater or equal (string with numeric comparison)
)

const (
	// certificate positions (within a standard certificate chain)
	leafCertIndex   uint32 = 0          // index for leaf (first in chain)
	anchorCertIndex        = ^uint32(0) // index for anchor (last in chain), equiv to -1
)

func generateRequirements(id string, h hash.Hash, signingMaterial pki.SigningMaterial) (*macho.Blob, []byte, error) {
	var reqBytes []byte
	if signingMaterial.Signer == nil {
		log.Trace("skipping adding designated requirement because no signer was found")
		reqBytes = []byte{0, 0, 0, 0}
	} else {
		req, err := newRequirements(id, signingMaterial)
		if err != nil {
			return nil, nil, err
		}

		reqBytes, err = restruct.Pack(macho.SigningOrder, req)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to encode requirement: %w", err)
		}
	}

	blob := macho.NewBlob(macho.MagicRequirements, reqBytes)

	blobBytes, err := restruct.Pack(macho.SigningOrder, &blob)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to encode requirements blob: %w", err)
	}

	// the requirements hash is against the entire blob, not just the payload
	_, err = h.Write(blobBytes)
	if err != nil {
		return nil, nil, err
	}

	return &blob, h.Sum(nil), nil
}

func newRequirements(id string, signingMaterial pki.SigningMaterial) (*macho.Requirements, error) {
	requirementBytes, err := buildRequirementStatements(id, signingMaterial)
	if err != nil {
		return nil, fmt.Errorf("unable to build requirement statements from signing material: %w", err)
	}

	// TODO: what is this field??
	requirementBytes = append([]byte{0, 0, 0, 1}, requirementBytes...)

	reqBlob := macho.NewBlob(macho.MagicRequirement, requirementBytes)

	reqBlobBytes, err := restruct.Pack(macho.SigningOrder, &reqBlob)
	if err != nil {
		return nil, fmt.Errorf("unable to encode requirements blob: %w", err)
	}

	offsetFromStartOfBlob := uint32(unsafe.Sizeof(macho.BlobHeader{}) + unsafe.Sizeof(macho.RequirementsHeader{}))

	return &macho.Requirements{

		RequirementsHeader: macho.RequirementsHeader{
			Count:  1, // we only support a single requirement at this time
			Type:   macho.DesignatedRequirementType,
			Offset: offsetFromStartOfBlob,
		},
		Payload: reqBlobBytes,
	}, nil
}

func buildRequirementStatements(id string, signingMaterial pki.SigningMaterial) ([]byte, error) {
	var statements []reqStatement

	// add on the identifier
	if id != "" {
		var ops []uint32
		ops = append(ops, opIdent)
		ops = append(ops, encodeBytes([]byte(id))...)
		statements = append(statements, ops)
	}

	// add on "anchor apple generic"
	if signingMaterial.HasCertWithOrg("Apple Inc.") {
		statements = append(statements, []uint32{opAppleGenericAnchor})
	}

	// add on "appleCertificateExtensions cert extension check", usually on the intermediate cert
	appleCertificateExtensionsOID := asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 2, 6}
	index, certWithAppleCertificateExtensionsOID := signingMaterial.CertWithExtension(appleCertificateExtensionsOID)
	if index != -1 && certWithAppleCertificateExtensionsOID != nil && certWithAppleCertificateExtensionsOID.IsCA {
		slotIndex := uint32(index)
		if index == 0 {
			slotIndex = anchorCertIndex
		}
		var ops []uint32
		ops = append(ops, opCertGeneric)
		ops = append(ops, slotIndex)
		ops = append(ops, encodeBytes(encodeOID(appleCertificateExtensionsOID))...)
		ops = append(ops, matchExists)

		statements = append(statements, ops)
	}

	// add on subject OU check
	leafCert := signingMaterial.Leaf()
	if leafCert != nil && len(leafCert.Subject.OrganizationalUnit) > 0 {
		var ops []uint32
		ops = append(ops, opCertField)
		ops = append(ops, leafCertIndex)
		ops = append(ops, encodeBytes([]byte("subject.OU"))...)
		ops = append(ops, matchEqual)
		ops = append(ops, encodeBytes([]byte(leafCert.Subject.OrganizationalUnit[0]))...)

		statements = append(statements, ops)
	}

	// note: must be the latest operand added
	if len(statements) == 0 {
		// this is an empty requirements set
		statements = []reqStatement{{opFalse}}
	}

	// and-conjoin all statements

	var finalOps []uint32

	// note: this is polish notation
	for i := 0; i < len(statements); i++ {
		if i < len(statements)-1 {
			finalOps = append(finalOps, opAnd)
		}
		finalOps = append(finalOps, statements[i]...)
	}

	return restruct.Pack(macho.SigningOrder, finalOps)
}

func encodeOID(oid asn1.ObjectIdentifier) []byte {
	var res []byte
	b := addBase128Int(int64(oid[0])*40 + int64(oid[1]))
	res = append(res, b...)
	for _, v := range oid[2:] {
		b = addBase128Int(int64(v))
		res = append(res, b...)
	}
	return res
}

func addBase128Int(n int64) []byte {
	var length int
	if n == 0 {
		length = 1
	} else {
		for i := n; i > 0; i >>= 7 {
			length++
		}
	}

	var b []byte
	for i := length - 1; i >= 0; i-- {
		o := byte(n >> uint(i*7))
		o &= 0x7f
		if i != 0 {
			o |= 0x80
		}

		b = append(b, o)
	}
	return b
}

func encodeBytes(by []byte) []uint32 {
	var ops []uint32

	ops = append(ops, uint32(len(by)))

	encoded := byteAlignRequirementData(by)
	data := make([]uint32, len(encoded)/sizeOfUint32)
	for i := range data {
		data[i] = macho.SigningOrder.Uint32(encoded[i*sizeOfUint32 : (i+1)*sizeOfUint32])
	}
	ops = append(ops, data...)

	return ops
}

func byteAlignRequirementData(by []byte) []byte {
	alignmentSize := 4
	alignmentLength := roundUp(uint64(len(by)), uint64(alignmentSize))
	if alignmentLength > uint64(len(by)) {
		by = append(by, make([]byte, alignmentLength-uint64(len(by)))...)
	}
	return by
}

func roundUp(x, align uint64) uint64 {
	return (x + align - 1) & -align
}
