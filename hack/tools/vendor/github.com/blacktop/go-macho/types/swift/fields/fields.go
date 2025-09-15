package fields

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=FieldDescriptorKind,FieldRecordFlags -linecomment -output fields_string.go

// ref: swift/include/swift/Reflection/Records.h

// __TEXT.__swift5_fieldmd
// This section contains an array of field descriptors.
// A field descriptor contains a collection of field records for a single class,
// struct or enum declaration. Each field descriptor can be a different length depending on how many field records the type contains.

const SWIFT_REFLECTION_METADATA_VERSION = 3 // superclass field

type FieldRecordFlags uint32

const (
	// IsIndirectCase is this an indirect enum case?
	IsIndirectCase FieldRecordFlags = 0x1
	// IsVar is this a mutable `var` property?
	IsVar FieldRecordFlags = 0x2
	// IsArtificial is this an artificial field?
	IsArtificial FieldRecordFlags = 0x4
)

type FieldDescriptorKind uint16

const (
	// Swift nominal types.
	FDKindStruct FieldDescriptorKind = iota // struct
	FDKindClass                             // class
	FDKindEnum                              // enum

	// Fixed-size multi-payload enums have a special descriptor format that
	// encodes spare bits.
	//
	// FIXME: Actually implement this. For now, a descriptor with this kind
	// just means we also have a builtin descriptor from which we get the
	// size and alignment.
	FDKindMultiPayloadEnum // multi-payload enum

	// A Swift opaque protocol. There are no fields, just a record for the
	// type itself.
	FDKindProtocol // protocol

	// A Swift class-bound protocol.
	FDKindClassProtocol // class protocol

	// An Objective-C protocol, which may be imported or defined in Swift.
	FDKindObjCProtocol // objc protocol

	// An Objective-C class, which may be imported or defined in Swift.
	// In the former case, field type metadata is not emitted, and
	// must be obtained from the Objective-C runtime.
	FDKindObjCClass // objc class
)

type FDHeader struct {
	MangledTypeNameOffset int32
	SuperclassOffset      int32
	Kind                  FieldDescriptorKind
	FieldRecordSize       uint16
	NumFields             uint32
}

type FieldRecord struct {
	Name        string
	MangledType string
	Flags       string
}

type FieldRecordType struct {
	Flags                 FieldRecordFlags
	MangledTypeNameOffset int32
	FieldNameOffset       int32
}

type FieldDescriptor struct {
	FDHeader
	FieldRecords []FieldRecordType
}

type Field struct {
	Address     uint64
	MangledType string
	SuperClass  string
	Kind        string
	Records     []FieldRecord
	Offset      int64
	Descriptor  FieldDescriptor
}

func (f Field) IsEnum() bool {
	return f.Descriptor.Kind == FDKindEnum || f.Descriptor.Kind == FDKindMultiPayloadEnum
}
func (f Field) IsClass() bool {
	return f.Descriptor.Kind == FDKindClass || f.Descriptor.Kind == FDKindObjCClass
}
func (f Field) IsProtocol() bool {
	return f.Descriptor.Kind == FDKindProtocol || f.Descriptor.Kind == FDKindClassProtocol || f.Descriptor.Kind == FDKindObjCProtocol
}
func (f Field) String() string {
	var recs string
	if len(f.Records) > 0 {
		recs = "\n"
	}
	for _, r := range f.Records {
		var flags string
		var hasType string
		if f.Kind == "Enum" {
			flags = "case"
		} else {
			if r.Flags == "IsVar" {
				flags = "var"
			} else {
				flags = "let"
			}
		}
		if len(r.MangledType) > 0 {
			hasType = ": "
		}
		recs += fmt.Sprintf("        %s %s%s%s\n", flags, r.Name, hasType, r.MangledType)
	}
	if len(f.SuperClass) > 0 {
		return fmt.Sprintf("// %#x:\n%s %s.%s {%s}\n", f.Address, f.Kind, f.MangledType, f.SuperClass, recs)
	}
	return fmt.Sprintf("// %#x:\n%s %s {%s}\n", f.Address, strings.ToLower(f.Kind), f.MangledType, recs)
}
