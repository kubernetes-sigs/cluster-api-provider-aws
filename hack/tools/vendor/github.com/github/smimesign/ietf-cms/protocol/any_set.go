package protocol

import (
	"encoding/asn1"
	"fmt"
)

// AnySet is a helper for dealing with SET OF ANY types.
type AnySet struct {
	Elements []asn1.RawValue `asn1:"set"`
}

// NewAnySet creates a new AnySet.
func NewAnySet(elts ...asn1.RawValue) AnySet {
	return AnySet{elts}
}

// DecodeAnySet manually decodes a SET OF ANY type, since Go's parser can't
// handle them.
func DecodeAnySet(rv asn1.RawValue) (as AnySet, err error) {
	// Make sure it's really a SET.
	if rv.Class != asn1.ClassUniversal {
		err = ASN1Error{fmt.Sprintf("Bad class. Expecting %d, got %d", asn1.ClassUniversal, rv.Class)}
		return
	}
	if rv.Tag != asn1.TagSet {
		err = ASN1Error{fmt.Sprintf("Bad tag. Expecting %d, got %d", asn1.TagSet, rv.Tag)}
		return
	}

	// Decode each element.
	der := rv.Bytes
	for len(der) > 0 {
		if der, err = asn1.Unmarshal(der, &rv); err != nil {
			return
		}

		as.Elements = append(as.Elements, rv)
	}

	return
}

// Encode manually encodes a SET OF ANY type, since Go's parser can't handle
// them.
func (as AnySet) Encode(dst *asn1.RawValue) (err error) {
	dst.Class = asn1.ClassUniversal
	dst.Tag = asn1.TagSet
	dst.IsCompound = true

	var der []byte
	for _, elt := range as.Elements {
		if der, err = asn1.Marshal(elt); err != nil {
			return
		}

		dst.Bytes = append(dst.Bytes, der...)
	}

	dst.FullBytes, err = asn1.Marshal(*dst)

	return
}
