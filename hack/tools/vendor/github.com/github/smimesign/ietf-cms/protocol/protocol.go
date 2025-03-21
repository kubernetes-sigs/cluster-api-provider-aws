// Package protocol implements low level CMS types, parsing and generation.
package protocol

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	_ "crypto/sha1" // for crypto.SHA1
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/github/smimesign/ietf-cms/oid"
)

// ASN1Error is an error from parsing ASN.1 structures.
type ASN1Error struct {
	Message string
}

// Error implements the error interface.
func (err ASN1Error) Error() string {
	return fmt.Sprintf("cms/protocol: ASN.1 Error — %s", err.Message)
}

var (
	// ErrWrongType is returned by methods that make assumptions about types.
	// Helper methods are defined for accessing CHOICE and  ANY feilds. These
	// helper methods get the value of the field, assuming it is of a given type.
	// This error is returned if that assumption is wrong and the field has a
	// different type.
	ErrWrongType = errors.New("cms/protocol: wrong choice or any type")

	// ErrNoCertificate is returned when a requested certificate cannot be found.
	ErrNoCertificate = errors.New("no certificate found")

	// ErrUnsupported is returned when an unsupported type or version
	// is encountered.
	ErrUnsupported = ASN1Error{"unsupported type or version"}

	// ErrTrailingData is returned when extra data is found after parsing an ASN.1
	// structure.
	ErrTrailingData = ASN1Error{"unexpected trailing data"}
)

// ContentInfo ::= SEQUENCE {
//   contentType ContentType,
//   content [0] EXPLICIT ANY DEFINED BY contentType }
//
// ContentType ::= OBJECT IDENTIFIER
type ContentInfo struct {
	ContentType asn1.ObjectIdentifier
	Content     asn1.RawValue `asn1:"explicit,tag:0"`
}

// ParseContentInfo parses a top-level ContentInfo type from BER encoded data.
func ParseContentInfo(ber []byte) (ci ContentInfo, err error) {
	var der []byte
	if der, err = BER2DER(ber); err != nil {
		return
	}

	var rest []byte
	if rest, err = asn1.Unmarshal(der, &ci); err != nil {
		return
	}
	if len(rest) > 0 {
		err = ErrTrailingData
	}

	return
}

// SignedDataContent gets the content assuming contentType is signedData.
func (ci ContentInfo) SignedDataContent() (*SignedData, error) {
	if !ci.ContentType.Equal(oid.ContentTypeSignedData) {
		return nil, ErrWrongType
	}

	sd := new(SignedData)
	if rest, err := asn1.Unmarshal(ci.Content.Bytes, sd); err != nil {
		return nil, err
	} else if len(rest) > 0 {
		return nil, ErrTrailingData
	}

	return sd, nil
}

// EncapsulatedContentInfo ::= SEQUENCE {
//   eContentType ContentType,
//   eContent [0] EXPLICIT OCTET STRING OPTIONAL }
//
// ContentType ::= OBJECT IDENTIFIER
type EncapsulatedContentInfo struct {
	EContentType asn1.ObjectIdentifier
	EContent     asn1.RawValue `asn1:"optional,explicit,tag:0"`
}

// NewDataEncapsulatedContentInfo creates a new EncapsulatedContentInfo of type
// id-data.
func NewDataEncapsulatedContentInfo(data []byte) (EncapsulatedContentInfo, error) {
	return NewEncapsulatedContentInfo(oid.ContentTypeData, data)
}

// NewEncapsulatedContentInfo creates a new EncapsulatedContentInfo.
func NewEncapsulatedContentInfo(contentType asn1.ObjectIdentifier, content []byte) (EncapsulatedContentInfo, error) {
	octets, err := asn1.Marshal(asn1.RawValue{
		Class:      asn1.ClassUniversal,
		Tag:        asn1.TagOctetString,
		Bytes:      content,
		IsCompound: false,
	})
	if err != nil {
		return EncapsulatedContentInfo{}, err
	}

	return EncapsulatedContentInfo{
		EContentType: contentType,
		EContent: asn1.RawValue{
			Class:      asn1.ClassContextSpecific,
			Tag:        0,
			Bytes:      octets,
			IsCompound: true,
		},
	}, nil
}

// EContentValue gets the OCTET STRING EContent value without tag or length.
// This is what the message digest is calculated over. A nil byte slice is
// returned if the OPTIONAL eContent field is missing.
func (eci EncapsulatedContentInfo) EContentValue() ([]byte, error) {
	if eci.EContent.Bytes == nil {
		return nil, nil
	}

	// The EContent is an `[0] EXPLICIT OCTET STRING`. EXPLICIT means that there
	// is another whole tag wrapping the OCTET STRING. When we decoded the
	// EContent into a asn1.RawValue we're just getting that outer tag, so the
	// EContent.Bytes is the encoded OCTET STRING, which is what we really want
	// the value of.
	var octets asn1.RawValue
	if rest, err := asn1.Unmarshal(eci.EContent.Bytes, &octets); err != nil {
		return nil, err
	} else if len(rest) > 0 {
		return nil, ErrTrailingData
	}
	if octets.Class != asn1.ClassUniversal || octets.Tag != asn1.TagOctetString {
		return nil, ASN1Error{"bad tag or class"}
	}

	// While we already tried converting BER to DER, we didn't take constructed
	// types into account. Constructed string types, as opposed to primitive
	// types, can encode indefinite length strings by including a bunch of
	// sub-strings that are joined together to get the actual value. Gpgsm uses
	// a constructed OCTET STRING for the EContent, so we have to manually decode
	// it here.
	var value []byte
	if octets.IsCompound {
		rest := octets.Bytes
		for len(rest) > 0 {
			var err error
			if rest, err = asn1.Unmarshal(rest, &octets); err != nil {
				return nil, err
			}

			// Don't allow further constructed types.
			if octets.Class != asn1.ClassUniversal || octets.Tag != asn1.TagOctetString || octets.IsCompound {
				return nil, ASN1Error{"bad class or tag"}
			}

			value = append(value, octets.Bytes...)
		}
	} else {
		value = octets.Bytes
	}

	return value, nil
}

// IsTypeData checks if the EContentType is id-data.
func (eci EncapsulatedContentInfo) IsTypeData() bool {
	return eci.EContentType.Equal(oid.ContentTypeData)
}

// DataEContent gets the EContent assuming EContentType is data.
func (eci EncapsulatedContentInfo) DataEContent() ([]byte, error) {
	if !eci.IsTypeData() {
		return nil, ErrWrongType
	}
	return eci.EContentValue()
}

// Attribute ::= SEQUENCE {
//   attrType OBJECT IDENTIFIER,
//   attrValues SET OF AttributeValue }
//
// AttributeValue ::= ANY
type Attribute struct {
	Type asn1.ObjectIdentifier

	// This should be a SET OF ANY, but Go's asn1 parser can't handle slices of
	// RawValues. Use value() to get an AnySet of the value.
	RawValue asn1.RawValue
}

// NewAttribute creates a single-value Attribute.
func NewAttribute(typ asn1.ObjectIdentifier, val interface{}) (attr Attribute, err error) {
	var der []byte
	if der, err = asn1.Marshal(val); err != nil {
		return
	}

	var rv asn1.RawValue
	if _, err = asn1.Unmarshal(der, &rv); err != nil {
		return
	}

	if err = NewAnySet(rv).Encode(&attr.RawValue); err != nil {
		return
	}

	attr.Type = typ

	return
}

// Value further decodes the attribute Value as a SET OF ANY, which Go's asn1
// parser can't handle directly.
func (a Attribute) Value() (AnySet, error) {
	return DecodeAnySet(a.RawValue)
}

// Attributes is a common Go type for SignedAttributes and UnsignedAttributes.
//
// SignedAttributes ::= SET SIZE (1..MAX) OF Attribute
//
// UnsignedAttributes ::= SET SIZE (1..MAX) OF Attribute
type Attributes []Attribute

// MarshaledForSigning DER encodes the Attributes as needed for signing
// SignedAttributes. RFC5652 explains this encoding:
//   A separate encoding of the signedAttrs field is performed for message
//   digest calculation. The IMPLICIT [0] tag in the signedAttrs is not used for
//   the DER encoding, rather an EXPLICIT SET OF tag is used.  That is, the DER
//   encoding of the EXPLICIT SET OF tag, rather than of the IMPLICIT [0] tag,
//   MUST be included in the message digest calculation along with the length
//   and content octets of the SignedAttributes value.
func (attrs Attributes) MarshaledForSigning() ([]byte, error) {
	seq, err := asn1.Marshal(struct {
		Attributes `asn1:"set"`
	}{attrs})

	if err != nil {
		return nil, err
	}

	// unwrap the outer SEQUENCE
	var raw asn1.RawValue
	if _, err = asn1.Unmarshal(seq, &raw); err != nil {
		return nil, err
	}

	return raw.Bytes, nil
}

// MarshaledForVerification DER encodes the Attributes as needed for
// verification of SignedAttributes. This is done differently than
// MarshaledForSigning because when verifying attributes, we need to
// use the received order.
func (attrs Attributes) MarshaledForVerification() ([]byte, error) {
	seq, err := asn1.Marshal(struct {
		Attributes `asn1:"sequence"`
	}{attrs})

	if err != nil {
		return nil, err
	}

	// unwrap the outer SEQUENCE
	var raw asn1.RawValue
	if _, err = asn1.Unmarshal(seq, &raw); err != nil {
		return nil, err
	}

	// Change SEQUENCE OF to SET OF.
	raw.Bytes[0] = 0x31
	return raw.Bytes, nil
}

// GetOnlyAttributeValueBytes gets an attribute value, returning an error if the
// attribute occurs multiple times or has multiple values.
func (attrs Attributes) GetOnlyAttributeValueBytes(oid asn1.ObjectIdentifier) (rv asn1.RawValue, err error) {
	var vals []AnySet
	if vals, err = attrs.GetValues(oid); err != nil {
		return
	}
	if len(vals) != 1 {
		err = ASN1Error{"bad attribute count"}
		return
	}
	if len(vals[0].Elements) != 1 {
		err = ASN1Error{"bad attribute element count"}
		return
	}

	return vals[0].Elements[0], nil
}

// GetValues retreives the attributes with the given OID. A nil value is
// returned if the OPTIONAL SET of Attributes is missing from the SignerInfo. An
// empty slice is returned if the specified attribute isn't in the set.
func (attrs Attributes) GetValues(oid asn1.ObjectIdentifier) ([]AnySet, error) {
	if attrs == nil {
		return nil, nil
	}

	vals := []AnySet{}
	for _, attr := range attrs {
		if attr.Type.Equal(oid) {
			val, err := attr.Value()
			if err != nil {
				return nil, err
			}

			vals = append(vals, val)
		}
	}

	return vals, nil
}

// HasAttribute checks if an attribute is present.
func (attrs Attributes) HasAttribute(oid asn1.ObjectIdentifier) bool {
	for _, attr := range attrs {
		if attr.Type.Equal(oid) {
			return true
		}
	}

	return false
}

// IssuerAndSerialNumber ::= SEQUENCE {
// 	issuer Name,
// 	serialNumber CertificateSerialNumber }
//
// CertificateSerialNumber ::= INTEGER
type IssuerAndSerialNumber struct {
	Issuer       asn1.RawValue
	SerialNumber *big.Int
}

// NewIssuerAndSerialNumber creates a IssuerAndSerialNumber SID for the given
// cert.
func NewIssuerAndSerialNumber(cert *x509.Certificate) (rv asn1.RawValue, err error) {
	sid := IssuerAndSerialNumber{
		SerialNumber: new(big.Int).Set(cert.SerialNumber),
	}

	if _, err = asn1.Unmarshal(cert.RawIssuer, &sid.Issuer); err != nil {
		return
	}

	var der []byte
	if der, err = asn1.Marshal(sid); err != nil {
		return
	}

	if _, err = asn1.Unmarshal(der, &rv); err != nil {
		return
	}

	return
}

// SignerInfo ::= SEQUENCE {
//   version CMSVersion,
//   sid SignerIdentifier,
//   digestAlgorithm DigestAlgorithmIdentifier,
//   signedAttrs [0] IMPLICIT SignedAttributes OPTIONAL,
//   signatureAlgorithm SignatureAlgorithmIdentifier,
//   signature SignatureValue,
//   unsignedAttrs [1] IMPLICIT UnsignedAttributes OPTIONAL }
//
// CMSVersion ::= INTEGER
//               { v0(0), v1(1), v2(2), v3(3), v4(4), v5(5) }
//
// SignerIdentifier ::= CHOICE {
//   issuerAndSerialNumber IssuerAndSerialNumber,
//   subjectKeyIdentifier [0] SubjectKeyIdentifier }
//
// DigestAlgorithmIdentifier ::= AlgorithmIdentifier
//
// SignedAttributes ::= SET SIZE (1..MAX) OF Attribute
//
// SignatureAlgorithmIdentifier ::= AlgorithmIdentifier
//
// SignatureValue ::= OCTET STRING
//
// UnsignedAttributes ::= SET SIZE (1..MAX) OF Attribute
type SignerInfo struct {
	Version            int
	SID                asn1.RawValue
	DigestAlgorithm    pkix.AlgorithmIdentifier
	SignedAttrs        Attributes `asn1:"optional,tag:0"`
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Signature          []byte
	UnsignedAttrs      Attributes `asn1:"set,optional,tag:1"`
}

// FindCertificate finds this SignerInfo's certificate in a slice of
// certificates.
func (si SignerInfo) FindCertificate(certs []*x509.Certificate) (*x509.Certificate, error) {
	switch si.Version {
	case 1: // SID is issuer and serial number
		isn, err := si.issuerAndSerialNumberSID()
		if err != nil {
			return nil, err
		}

		for _, cert := range certs {
			if bytes.Equal(cert.RawIssuer, isn.Issuer.FullBytes) && isn.SerialNumber.Cmp(cert.SerialNumber) == 0 {
				return cert, nil
			}
		}
	case 3: // SID is SubjectKeyIdentifier
		ski, err := si.subjectKeyIdentifierSID()
		if err != nil {
			return nil, err
		}

		for _, cert := range certs {
			for _, ext := range cert.Extensions {
				if oid.ExtensionSubjectKeyIdentifier.Equal(ext.Id) {
					if bytes.Equal(ski, ext.Value) {
						return cert, nil
					}
				}
			}
		}
	default:
		return nil, ErrUnsupported
	}

	return nil, ErrNoCertificate
}

// issuerAndSerialNumberSID gets the SID, assuming it is a issuerAndSerialNumber.
func (si SignerInfo) issuerAndSerialNumberSID() (isn IssuerAndSerialNumber, err error) {
	if si.SID.Class != asn1.ClassUniversal || si.SID.Tag != asn1.TagSequence {
		err = ErrWrongType
		return
	}

	var rest []byte
	if rest, err = asn1.Unmarshal(si.SID.FullBytes, &isn); err == nil && len(rest) > 0 {
		err = ErrTrailingData
	}

	return
}

// subjectKeyIdentifierSID gets the SID, assuming it is a subjectKeyIdentifier.
func (si SignerInfo) subjectKeyIdentifierSID() ([]byte, error) {
	if si.SID.Class != asn1.ClassContextSpecific || si.SID.Tag != 0 {
		return nil, ErrWrongType
	}

	return si.SID.Bytes, nil
}

// Hash gets the crypto.Hash associated with this SignerInfo's DigestAlgorithm.
// 0 is returned for unrecognized algorithms.
func (si SignerInfo) Hash() (crypto.Hash, error) {
	algo := si.DigestAlgorithm.Algorithm.String()
	hash := oid.DigestAlgorithmToCryptoHash[algo]
	if hash == 0 || !hash.Available() {
		return 0, ErrUnsupported
	}

	return hash, nil
}

// X509SignatureAlgorithm gets the x509.SignatureAlgorithm that should be used
// for verifying this SignerInfo's signature.
func (si SignerInfo) X509SignatureAlgorithm() x509.SignatureAlgorithm {
	var (
		sigOID    = si.SignatureAlgorithm.Algorithm.String()
		digestOID = si.DigestAlgorithm.Algorithm.String()
	)

	if sa := oid.SignatureAlgorithmToX509SignatureAlgorithm[sigOID]; sa != x509.UnknownSignatureAlgorithm {
		return sa
	}

	return oid.PublicKeyAndDigestAlgorithmToX509SignatureAlgorithm[sigOID][digestOID]
}

// GetContentTypeAttribute gets the signed ContentType attribute from the
// SignerInfo.
func (si SignerInfo) GetContentTypeAttribute() (asn1.ObjectIdentifier, error) {
	rv, err := si.SignedAttrs.GetOnlyAttributeValueBytes(oid.AttributeContentType)
	if err != nil {
		return nil, err
	}

	var ct asn1.ObjectIdentifier
	if rest, err := asn1.Unmarshal(rv.FullBytes, &ct); err != nil {
		return nil, err
	} else if len(rest) > 0 {
		return nil, ErrTrailingData
	}

	return ct, nil
}

// GetMessageDigestAttribute gets the signed MessageDigest attribute from the
// SignerInfo.
func (si SignerInfo) GetMessageDigestAttribute() ([]byte, error) {
	rv, err := si.SignedAttrs.GetOnlyAttributeValueBytes(oid.AttributeMessageDigest)
	if err != nil {
		return nil, err
	}
	if rv.Class != asn1.ClassUniversal || rv.Tag != asn1.TagOctetString {
		return nil, ASN1Error{"bad class or tag"}
	}

	return rv.Bytes, nil
}

// GetSigningTimeAttribute gets the signed SigningTime attribute from the
// SignerInfo.
func (si SignerInfo) GetSigningTimeAttribute() (time.Time, error) {
	var t time.Time

	if !si.SignedAttrs.HasAttribute(oid.AttributeSigningTime) {
		return t, nil
	}
	rv, err := si.SignedAttrs.GetOnlyAttributeValueBytes(oid.AttributeSigningTime)
	if err != nil {
		return t, err
	}
	if rv.Class != asn1.ClassUniversal || (rv.Tag != asn1.TagUTCTime && rv.Tag != asn1.TagGeneralizedTime) {
		return t, ASN1Error{"bad class or tag"}
	}

	if rest, err := asn1.Unmarshal(rv.FullBytes, &t); err != nil {
		return t, err
	} else if len(rest) > 0 {
		return t, ErrTrailingData
	}

	return t, nil
}

// SignedData ::= SEQUENCE {
//   version CMSVersion,
//   digestAlgorithms DigestAlgorithmIdentifiers,
//   encapContentInfo EncapsulatedContentInfo,
//   certificates [0] IMPLICIT CertificateSet OPTIONAL,
//   crls [1] IMPLICIT RevocationInfoChoices OPTIONAL,
//   signerInfos SignerInfos }
//
// CMSVersion ::= INTEGER
//               { v0(0), v1(1), v2(2), v3(3), v4(4), v5(5) }
//
// DigestAlgorithmIdentifiers ::= SET OF DigestAlgorithmIdentifier
//
// CertificateSet ::= SET OF CertificateChoices
//
// CertificateChoices ::= CHOICE {
//   certificate Certificate,
//   extendedCertificate [0] IMPLICIT ExtendedCertificate, -- Obsolete
//   v1AttrCert [1] IMPLICIT AttributeCertificateV1,       -- Obsolete
//   v2AttrCert [2] IMPLICIT AttributeCertificateV2,
//   other [3] IMPLICIT OtherCertificateFormat }
//
// OtherCertificateFormat ::= SEQUENCE {
//   otherCertFormat OBJECT IDENTIFIER,
//   otherCert ANY DEFINED BY otherCertFormat }
//
// RevocationInfoChoices ::= SET OF RevocationInfoChoice
//
// RevocationInfoChoice ::= CHOICE {
//   crl CertificateList,
//   other [1] IMPLICIT OtherRevocationInfoFormat }
//
// OtherRevocationInfoFormat ::= SEQUENCE {
//   otherRevInfoFormat OBJECT IDENTIFIER,
//   otherRevInfo ANY DEFINED BY otherRevInfoFormat }
//
// SignerInfos ::= SET OF SignerInfo
type SignedData struct {
	Version          int
	DigestAlgorithms []pkix.AlgorithmIdentifier `asn1:"set"`
	EncapContentInfo EncapsulatedContentInfo
	Certificates     []asn1.RawValue `asn1:"optional,set,tag:0"`
	CRLs             []asn1.RawValue `asn1:"optional,set,tag:1"`
	SignerInfos      []SignerInfo    `asn1:"set"`
}

// NewSignedData creates a new SignedData.
func NewSignedData(eci EncapsulatedContentInfo) (*SignedData, error) {
	// The version is picked based on which CMS features are used. We only use
	// version 1 features, except for supporting non-data econtent.
	version := 1
	if !eci.IsTypeData() {
		version = 3
	}

	return &SignedData{
		Version:          version,
		DigestAlgorithms: []pkix.AlgorithmIdentifier{},
		EncapContentInfo: eci,
		SignerInfos:      []SignerInfo{},
	}, nil
}

// AddSignerInfo adds a SignerInfo to the SignedData.
func (sd *SignedData) AddSignerInfo(chain []*x509.Certificate, signer crypto.Signer) error {
	// figure out which certificate is associated with signer.
	pub, err := x509.MarshalPKIXPublicKey(signer.Public())
	if err != nil {
		return err
	}

	var (
		cert    *x509.Certificate
		certPub []byte
	)

	for _, c := range chain {
		if err = sd.AddCertificate(c); err != nil {
			return err
		}

		if certPub, err = x509.MarshalPKIXPublicKey(c.PublicKey); err != nil {
			return err
		}

		if bytes.Equal(pub, certPub) {
			cert = c
		}
	}
	if cert == nil {
		return ErrNoCertificate
	}

	sid, err := NewIssuerAndSerialNumber(cert)
	if err != nil {
		return err
	}

	digestAlgorithmID := digestAlgorithmForPublicKey(pub)

	signatureAlgorithmOID, ok := oid.X509PublicKeyAndDigestAlgorithmToSignatureAlgorithm[cert.PublicKeyAlgorithm][digestAlgorithmID.Algorithm.String()]
	if !ok {
		return errors.New("unsupported certificate public key algorithm")
	}

	signatureAlgorithmID := pkix.AlgorithmIdentifier{Algorithm: signatureAlgorithmOID}

	si := SignerInfo{
		Version:            1,
		SID:                sid,
		DigestAlgorithm:    digestAlgorithmID,
		SignedAttrs:        nil,
		SignatureAlgorithm: signatureAlgorithmID,
		Signature:          nil,
		UnsignedAttrs:      nil,
	}

	// Get the message
	content, err := sd.EncapContentInfo.EContentValue()
	if err != nil {
		return err
	}
	if content == nil {
		return errors.New("already detached")
	}

	// Digest the message.
	hash, err := si.Hash()
	if err != nil {
		return err
	}
	md := hash.New()
	if _, err = md.Write(content); err != nil {
		return err
	}

	// Build our SignedAttributes
	stAttr, err := NewAttribute(oid.AttributeSigningTime, time.Now().UTC())
	if err != nil {
		return err
	}
	mdAttr, err := NewAttribute(oid.AttributeMessageDigest, md.Sum(nil))
	if err != nil {
		return err
	}
	ctAttr, err := NewAttribute(oid.AttributeContentType, sd.EncapContentInfo.EContentType)
	if err != nil {
		return err
	}

	// sort attributes to match required order in marshaled form
	si.SignedAttrs, err = sortAttributes(stAttr, mdAttr, ctAttr)
	if err != nil {
		return err
	}

	// Signature is over the marshaled signed attributes
	sm, err := si.SignedAttrs.MarshaledForSigning()
	if err != nil {
		return err
	}
	smd := hash.New()
	if _, errr := smd.Write(sm); errr != nil {
		return errr
	}
	if si.Signature, err = signer.Sign(rand.Reader, smd.Sum(nil), hash); err != nil {
		return err
	}

	sd.addDigestAlgorithm(si.DigestAlgorithm)

	sd.SignerInfos = append(sd.SignerInfos, si)

	return nil
}

func sortAttributes(attrs ...Attribute) ([]Attribute, error) {
	// Sort attrs by their encoded values (including tag and
	// lengths) as specified in X690 Section 11.6 and implemented
	// in go >= 1.15's asn1.Marshal().
	sort.Slice(attrs, func(i, j int) bool {
		return bytes.Compare(
			attrs[i].RawValue.FullBytes,
			attrs[j].RawValue.FullBytes) < 0
	})

	return attrs, nil
}

// algorithmsForPublicKey takes an opinionated stance on what algorithms to use
// for the given public key.
func digestAlgorithmForPublicKey(pub crypto.PublicKey) pkix.AlgorithmIdentifier {
	if ecPub, ok := pub.(*ecdsa.PublicKey); ok {
		switch ecPub.Curve {
		case elliptic.P384():
			return pkix.AlgorithmIdentifier{Algorithm: oid.DigestAlgorithmSHA384}
		case elliptic.P521():
			return pkix.AlgorithmIdentifier{Algorithm: oid.DigestAlgorithmSHA512}
		}
	}

	return pkix.AlgorithmIdentifier{Algorithm: oid.DigestAlgorithmSHA256}
}

// ClearCertificates removes all certificates.
func (sd *SignedData) ClearCertificates() {
	sd.Certificates = []asn1.RawValue{}
}

// AddCertificate adds a *x509.Certificate.
func (sd *SignedData) AddCertificate(cert *x509.Certificate) error {
	for _, existing := range sd.Certificates {
		if bytes.Equal(existing.Bytes, cert.Raw) {
			return errors.New("certificate already added")
		}
	}

	var rv asn1.RawValue
	if _, err := asn1.Unmarshal(cert.Raw, &rv); err != nil {
		return err
	}

	sd.Certificates = append(sd.Certificates, rv)

	return nil
}

// addDigestAlgorithm adds a new AlgorithmIdentifier if it doesn't exist yet.
func (sd *SignedData) addDigestAlgorithm(algo pkix.AlgorithmIdentifier) {
	for _, existing := range sd.DigestAlgorithms {
		if existing.Algorithm.Equal(algo.Algorithm) {
			return
		}
	}

	sd.DigestAlgorithms = append(sd.DigestAlgorithms, algo)
}

// X509Certificates gets the certificates, assuming that they're X.509 encoded.
func (sd *SignedData) X509Certificates() ([]*x509.Certificate, error) {
	// Certificates field is optional. Handle missing value.
	if sd.Certificates == nil {
		return nil, nil
	}

	// Empty set
	if len(sd.Certificates) == 0 {
		return []*x509.Certificate{}, nil
	}

	certs := make([]*x509.Certificate, 0, len(sd.Certificates))
	for _, raw := range sd.Certificates {
		if raw.Class != asn1.ClassUniversal || raw.Tag != asn1.TagSequence {
			return nil, ErrUnsupported
		}

		x509, err := x509.ParseCertificate(raw.FullBytes)
		if err != nil {
			return nil, err
		}

		certs = append(certs, x509)
	}

	return certs, nil
}

// ContentInfo returns the SignedData wrapped in a ContentInfo packet.
func (sd *SignedData) ContentInfo() (ContentInfo, error) {
	var nilCI ContentInfo

	der, err := asn1.Marshal(*sd)
	if err != nil {
		return nilCI, err
	}

	return ContentInfo{
		ContentType: oid.ContentTypeSignedData,
		Content: asn1.RawValue{
			Class:      asn1.ClassContextSpecific,
			Tag:        0,
			Bytes:      der,
			IsCompound: true,
		},
	}, nil

}

// ContentInfoDER returns the SignedData wrapped in a ContentInfo packet and DER
// encoded.
func (sd *SignedData) ContentInfoDER() ([]byte, error) {
	ci, err := sd.ContentInfo()
	if err != nil {
		return nil, err
	}

	return asn1.Marshal(ci)
}
