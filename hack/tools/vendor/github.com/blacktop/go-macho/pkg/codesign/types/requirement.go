package types

import (
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/blacktop/go-macho/pkg/trie"
	mtypes "github.com/blacktop/go-macho/types"
)

// Requirement object
type Requirement struct {
	RequirementsBlob
	Requirements
	Detail string `json:"detail,omitempty"`
}

// RequirementsBlob object
type RequirementsBlob struct {
	Magic  Magic  `json:"magic,omitempty"`  // magic number
	Length uint32 `json:"length,omitempty"` // total length of blob
	Data   uint32 `json:"data,omitempty"`   // zero for dyld shared cache
}

type RequirementType uint32

const (
	HostRequirementType       RequirementType = 1 /* what hosts may run us */
	GuestRequirementType      RequirementType = 2 /* what guests we may run */
	DesignatedRequirementType RequirementType = 3 /* designated requirement */
	LibraryRequirementType    RequirementType = 4 /* what libraries we may link against */
	PluginRequirementType     RequirementType = 5 /* what plug-ins we may load */
)

func (cm RequirementType) String() string {
	switch cm {
	case HostRequirementType:
		return "Host Requirement"
	case GuestRequirementType:
		return "Guest Requirement"
	case DesignatedRequirementType:
		return "Designated Requirement"
	case LibraryRequirementType:
		return "Library Requirement"
	case PluginRequirementType:
		return "Plugin Requirement"
	default:
		return fmt.Sprintf("RequirementType(%d)", cm)
	}
}

// Requirements object
type Requirements struct {
	Type   RequirementType `json:"type,omitempty"`   // type of entry
	Offset uint32          `json:"offset,omitempty"` // offset of entry
}

// NOTE: https://opensource.apple.com/source/libsecurity_codesigning/libsecurity_codesigning-36591/lib/requirement.h.auto.html

// exprForm opcodes.
//
// Opcodes are broken into flags in the (HBO) high byte, and an opcode value
// in the remaining 24 bits. Note that opcodes will remain fairly small
// (almost certainly <60000), so we have the third byte to play around with
// in the future, if needed. For now, small opcodes effective reserve this byte
// as zero.
// The flag byte allows for limited understanding of unknown opcodes. It allows
// the interpreter to use the known opcode parts of the program while semi-creatively
// disregarding the parts it doesn't know about. An unrecognized opcode with zero
// flag byte causes evaluation to categorically fail, since the semantics of such
// an opcode cannot safely be predicted.
const (
	// semantic bits or'ed into the opcode
	opFlagMask     exprOp = 0xFF000000 // high bit flags
	opGenericFalse exprOp = 0x80000000 // has size field; okay to default to false
	opGenericSkip  exprOp = 0x40000000 // has size field; skip and continue
)

type exprOp uint32

const (
	opFalse              exprOp = iota // unconditionally false
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

func (o exprOp) String() string {
	return [...]string{
		"False",
		"True",
		"Ident",
		"AppleAnchor",
		"AnchorHash",
		"InfoKeyValue",
		"And",
		"Or",
		"CDHash",
		"Not",
		"InfoKeyField",
		"CertField",
		"TrustedCert",
		"TrustedCerts",
		"CertGeneric",
		"AppleGenericAnchor",
		"EntitlementField",
		"CertPolicy",
		"NamedAnchor",
		"NamedCode",
		"exprOpCount",
	}[0]
}

type matchOp uint32

// match suffix opcodes
const (
	matchExists       matchOp = iota // anything but explicit "false" - no value stored
	matchEqual                       // equal (CFEqual)
	matchContains                    // partial match (substring)
	matchBeginsWith                  // partial match (initial substring)
	matchEndsWith                    // partial match (terminal substring)
	matchLessThan                    // less than (string with numeric comparison)
	matchGreaterThan                 // greater than (string with numeric comparison)
	matchLessEqual                   // less or equal (string with numeric comparison)
	matchGreaterEqual                // greater or equal (string with numeric comparison)
)

func (o matchOp) String() string {
	return [...]string{
		"Exists",
		"Equal",
		"Contains",
		"BeginsWith",
		"EndsWith",
		"LessThan",
		"GreaterThan",
		"LessEqual",
		"GreaterEqual",
	}[0]
}

const (
	// certificate positions (within a standard certificate chain)
	leafCertIndex   uint32 = 0          // index for leaf (first in chain)
	anchorCertIndex        = ^uint32(0) // index for anchor (last in chain), equiv to -1
)

func getData(r *bytes.Reader) ([]byte, error) {
	var idLength uint32

	err := binary.Read(r, binary.BigEndian, &idLength)
	if err != nil {
		return nil, err
	}

	// 4 byte align length
	alignedLength := uint32(mtypes.RoundUp(uint64(idLength), 4))

	data := make([]byte, alignedLength)

	_, err = r.Read(data)
	if err != nil {
		return nil, err
	}

	return data[:idLength], nil
}

func getMatch(r *bytes.Reader) (string, error) {
	var op matchOp
	err := binary.Read(r, binary.BigEndian, &op)
	if err != nil {
		return "", err
	}

	switch op {
	case matchExists:
		return " /* exists */", nil
	case matchEqual:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(" = \"%s\"", data), nil
	case matchContains:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(" ~ %s", data), nil
	case matchBeginsWith:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(" = %s*", data), nil
	case matchEndsWith:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(" = *%s", data), nil
	case matchLessThan:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(" < %s", data), nil
	case matchGreaterThan:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(" >= %s", data), nil
	case matchLessEqual:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(" <= %s", data), nil
	case matchGreaterEqual:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(" > %s", data), nil
	}
	return "", fmt.Errorf("MATCH OPCODE %d NOT UNDERSTOOD", op)
}

const (
	// certificate positions (within a standard certificate chain)
	leafCert   int32 = 0  // index for leaf (first in chain)
	anchorCert int32 = -1 // index for anchor (last in chain)
)

func getCertSlot(r *bytes.Reader) (string, error) {
	var slot int32

	err := binary.Read(r, binary.BigEndian, &slot)
	if err != nil {
		return "", err
	}

	switch slot {
	case leafCert:
		return "leaf", nil
	case anchorCert:
		return "root", nil
	default:
		return fmt.Sprintf("%d", slot), nil
	}
}

const (
	slPrimary = iota // syntax primary
	slAnd            // conjunctive
	slOr             // disjunctive
	slTop            // where we start
)

func getOid(r *bytes.Reader) (uint32, error) {
	var result uint32

	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			return 0, err
		}
		if err != nil {
			return 0, fmt.Errorf("could not parse OID value: %v", err)
		}

		result = uint32(result*128) + uint32(b&0x7f)

		// If high order bit is 1.
		if (b & 0x80) == 0 {
			break
		}
	}

	return result, nil
}

// NOTE:
// ref https://opensource.apple.com/source/Security/Security-59306.80.4/
// ref http://oid-info.com/get/1.2.840.113635.100.6.2.6
func toOID(data []byte) string {
	var oidStr string

	r := bytes.NewReader(data)

	oid1, err := getOid(r)
	if err != nil {
		return ""
	}

	q1 := uint32(math.Min(float64(oid1)/40, 2))
	oidStr += fmt.Sprintf("%d.%d", q1, oid1-q1*40)

	for {
		oid, err := getOid(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return ""
		}

		oidStr += fmt.Sprintf(".%d", uint32(oid))
	}

	return oidStr
}

func evalExpression(r *bytes.Reader, syntaxLevel int) (string, error) {
	var op exprOp

	err := binary.Read(r, binary.BigEndian, &op)
	if err != nil {
		return "", err
	}

	switch op {
	case opFalse:
		return "never", nil
	case opTrue:
		return "always", nil
	case opIdent:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("identifier \"%s\"", data), nil
	case opAppleAnchor:
		return "anchor apple", nil
	case opAppleGenericAnchor:
		return "anchor apple generic", nil
	case opAnchorHash:
		slot, err := getCertSlot(r)
		if err != nil {
			return "", err
		}
		data, err := getData(r)
		if err != nil {
			return "", err
		} // TODO data is a hashData str
		return fmt.Sprintf("certificate %s = %s", slot, data), nil
	case opInfoKeyValue:
		dot, err := getData(r)
		if err != nil {
			return "", err
		}
		data, err := getData(r)
		if err != nil {
			return "", err
		} // TODO dot is a dot str
		return fmt.Sprintf("info[%s] = %s", dot, data), nil
	case opAnd:
		var out string
		if syntaxLevel < slAnd {
			out += "("
		}
		part, err := evalExpression(r, slAnd)
		if err != nil {
			return "", err
		}
		out += part
		out += " and "
		part, err = evalExpression(r, slAnd)
		if err != nil {
			return "", err
		}
		out += part
		if syntaxLevel < slAnd {
			out += ")"
		}
		return out, nil
	case opOr:
		var out string
		if syntaxLevel < slOr {
			out += "("
		}
		part, err := evalExpression(r, slOr)
		if err != nil {
			return "", err
		}
		out += part
		out += " or "
		part, err = evalExpression(r, slOr)
		if err != nil {
			return "", err
		}
		out += part
		if syntaxLevel < slOr {
			out += ")"
		}
		return out, nil
	case opNot:
		part, err := evalExpression(r, slPrimary)
		if err != nil {
			return "", err
		}
		return "! " + part, nil
	case opCDHash:
		data, err := getData(r)
		if err != nil {
			return "", err
		} // TODO data is a hashData str
		return fmt.Sprintf(" cdhash %s", data), nil
	case opInfoKeyField:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		match, err := getMatch(r)
		if err != nil {
			return "", err
		} // TODO data is a dot str
		return fmt.Sprintf("info[%s] %s", data, match), nil
	case opEntitlementField:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		match, err := getMatch(r)
		if err != nil {
			return "", err
		} // TODO data is a dot str
		return fmt.Sprintf("entitlement[%s] %s", data, match), nil
	case opCertField:
		slot, err := getCertSlot(r)
		if err != nil {
			return "", err
		}
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		match, err := getMatch(r)
		if err != nil {
			return "", err
		} // TODO data is a dot str
		return fmt.Sprintf("certificate %s[%s] %s", slot, data, match), nil
	case opCertGeneric:
		slot, err := getCertSlot(r)
		if err != nil {
			return "", err
		}
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		match, err := getMatch(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("certificate %s[field.%s] %s", slot, toOID(data), match), nil
	case opCertPolicy:
		slot, err := getCertSlot(r)
		if err != nil {
			return "", err
		}
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		match, err := getMatch(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("certificate %s [policy.%s] %s", slot, toOID(data), match), nil
	case opTrustedCert:
		slot, err := getCertSlot(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("certificate %s trusted", slot), nil
	case opTrustedCerts:
		return "anchor trusted", nil
	case opNamedAnchor:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("anchor apple %s", string(data)), nil
	case opNamedCode:
		data, err := getData(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s)", string(data)), nil
	default:
		if (op & opGenericFalse) != 0 {
			return fmt.Sprintf(" false /* opcode %d */", op & ^opFlagMask), nil
		} else if (op & opGenericSkip) != 0 {
			return fmt.Sprintf(" /* opcode %d */", op & ^opFlagMask), nil
		}
		return fmt.Sprintf("OPCODE %d NOT UNDERSTOOD (ending print)", op), nil
	}
}

// ParseRequirements parses the requirements set bytes
func ParseRequirements(r *bytes.Reader, reqs Requirements) (string, error) {
	// NOTE: codesign -d -r- MACHO (to display requirement sets)
	r.Seek(int64(reqs.Offset), io.SeekStart)

	switch reqs.Type {
	case HostRequirementType:
		var reqSet []string
		for {
			rsPart, err := evalExpression(r, slTop)
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", err
			}
			reqSet = append(reqSet, rsPart)
		}
		return "host => " + strings.Join(reqSet, " "), nil
	case DesignatedRequirementType:
		var reqSet []string
		for {
			rsPart, err := evalExpression(r, slTop)
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", err
			}
			reqSet = append(reqSet, rsPart)
		}
		return strings.Join(reqSet, " "), nil
	default:
		return "", fmt.Errorf("failed to dump requirements set; found unsupported codesign requirement type '%s', please notify author", reqs.Type)
	}
}

// CreateRequirements creates a requirements set cs blob
// NOTE: /usr/bin/csreq -r="identifier com.foo.test" -t (to test it out)
func CreateRequirements(id string, certs []*x509.Certificate) (Blob, error) {

	if len(id) == 0 { // empty requirements set
		return NewBlob(MAGIC_REQUIREMENTS, make([]byte, 4)), nil
	}

	var ops []uint32
	var statements [][]uint32

	// add identifier
	ops = append(ops, uint32(opIdent))
	ops = append(ops, encodeBytes([]byte(id))...)
	statements = append(statements, ops)

	// add on "anchor apple generic"
	for _, cert := range certs {
		if len(cert.Subject.Organization) > 0 && cert.Subject.Organization[0] == "Apple Inc." {
			statements = append(statements, []uint32{uint32(opAppleGenericAnchor)})
			break
		}
	}

	// add appleCertificateExtensions cert extension check
	appleCertificateExtensionsOID := asn1.ObjectIdentifier{1, 2, 840, 113635, 100, 6, 2, 6}
	for idx, cert := range certs {
		for _, ext := range cert.Extensions {
			if ext.Id.Equal(appleCertificateExtensionsOID) {
				if cert.IsCA {
					slotIndex := uint32(idx)
					if idx == 0 {
						slotIndex = anchorCertIndex
					}
					var ops []uint32
					ops = append(ops, uint32(opCertGeneric))
					ops = append(ops, slotIndex)
					ops = append(ops, encodeBytes(encodeOID(appleCertificateExtensionsOID))...)
					ops = append(ops, uint32(matchExists))
					statements = append(statements, ops)
				}
			}
		}
	}

	if len(certs) > 1 {
		// add on subject OU check
		if !certs[len(certs)-1].IsCA {
			leafCert := certs[len(certs)-1]
			if leafCert != nil && len(leafCert.Subject.OrganizationalUnit) > 0 {
				var ops []uint32
				ops = append(ops, uint32(opCertField))
				ops = append(ops, leafCertIndex)
				ops = append(ops, encodeBytes([]byte("subject.OU"))...)
				ops = append(ops, uint32(matchEqual))
				ops = append(ops, encodeBytes([]byte(leafCert.Subject.OrganizationalUnit[0]))...)
				statements = append(statements, ops)
			}
		}
	}

	// and-conjoin all statements
	var finalOps []uint32
	for i := 0; i < len(statements); i++ {
		if i < len(statements)-1 {
			finalOps = append(finalOps, uint32(opAnd))
		}
		finalOps = append(finalOps, statements[i]...)
	}

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, uint32(1)); err != nil {
		return Blob{}, err
	}
	if err := binary.Write(&buf, binary.BigEndian, Requirements{
		Type:   DesignatedRequirementType,
		Offset: uint32(binary.Size(RequirementsBlob{}) + binary.Size(Requirements{})),
	}); err != nil {
		return Blob{}, err
	}
	if err := binary.Write(&buf, binary.BigEndian, RequirementsBlob{
		Magic:  MAGIC_REQUIREMENT,
		Length: uint32(binary.Size(RequirementsBlob{}) + binary.Size(finalOps)),
		Data:   1,
	}); err != nil {
		return Blob{}, err
	}
	if err := binary.Write(&buf, binary.BigEndian, finalOps); err != nil {
		return Blob{}, err
	}

	return Blob{
		BlobHeader: BlobHeader{
			Magic:  MAGIC_REQUIREMENTS,
			Length: uint32(binary.Size(BlobHeader{}) + buf.Len()),
		},
		Data: buf.Bytes(),
	}, nil
}

func encodeOID(oid asn1.ObjectIdentifier) []byte {
	var res bytes.Buffer
	trie.EncodeUleb128(&res, uint64(oid[0])*40+uint64(oid[1]))
	for _, v := range oid[2:] {
		trie.EncodeUleb128(&res, uint64(v))
	}
	return res.Bytes()
}

func encodeBytes(in []byte) []uint32 {
	var ops []uint32
	ops = append(ops, uint32(len(in)))
	if (len(in) % 4) != 0 {
		pad := make([]byte, 4-(len(in)%4))
		in = append(in, pad...)
	}
	data := make([]uint32, len(in)/binary.Size(uint32(0)))
	binary.Read(bytes.NewReader(in), binary.BigEndian, &data)
	ops = append(ops, data...)
	return ops
}
