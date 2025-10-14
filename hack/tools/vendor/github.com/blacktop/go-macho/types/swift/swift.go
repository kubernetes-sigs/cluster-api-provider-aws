package swift

import (
	"fmt"
	"strings"
)

// credit: https://knight.sc/reverse%20engineering/2019/07/17/swift-metadata.html

const (
	/// The name of the standard library, which is a reserved module name.
	STDLIB_NAME = "Swift"
	/// The name of the Onone support library, which is a reserved module name.
	SWIFT_ONONE_SUPPORT = "SwiftOnoneSupport"
	/// The name of the SwiftShims module, which contains private stdlib decls.
	SWIFT_SHIMS_NAME = "SwiftShims"
	/// The name of the Builtin module, which contains Builtin functions.
	BUILTIN_NAME = "Builtin"
	/// The name of the clang imported header module.
	CLANG_HEADER_MODULE_NAME = "__ObjC"
	/// The prefix of module names used by LLDB to capture Swift expressions
	LLDB_EXPRESSIONS_MODULE_NAME_PREFIX = "__lldb_expr_"

	/// The name of the fake module used to hold imported Objective-C things.
	MANGLING_MODULE_OBJC = "__C"
	/// The name of the fake module used to hold synthesized ClangImporter things.
	MANGLING_MODULE_CLANG_IMPORTER = "__C_Synthesized"

	/// The name prefix for C++ template instantiation imported as a Swift struct.
	CXX_TEMPLATE_INST_PREFIX = "__CxxTemplateInst"

	SEMANTICS_PROGRAMTERMINATION_POINT = "programtermination_point"

	/// The name of the Builtin type prefix
	BUILTIN_TYPE_NAME_PREFIX = "Builtin."
)

const (
	/// The name of the Builtin type for Int
	BUILTIN_TYPE_NAME_INT = "Builtin.Int"
	/// The name of the Builtin type for Int8
	BUILTIN_TYPE_NAME_INT8 = "Builtin.Int8"
	/// The name of the Builtin type for Int16
	BUILTIN_TYPE_NAME_INT16 = "Builtin.Int16"
	/// The name of the Builtin type for Int32
	BUILTIN_TYPE_NAME_INT32 = "Builtin.Int32"
	/// The name of the Builtin type for Int64
	BUILTIN_TYPE_NAME_INT64 = "Builtin.Int64"
	/// The name of the Builtin type for Int128
	BUILTIN_TYPE_NAME_INT128 = "Builtin.Int128"
	/// The name of the Builtin type for Int256
	BUILTIN_TYPE_NAME_INT256 = "Builtin.Int256"
	/// The name of the Builtin type for Int512
	BUILTIN_TYPE_NAME_INT512 = "Builtin.Int512"
	/// The name of the Builtin type for IntLiteral
	BUILTIN_TYPE_NAME_INTLITERAL = "Builtin.IntLiteral"
	/// The name of the Builtin type for IEEE Floating point types.
	BUILTIN_TYPE_NAME_FLOAT = "Builtin.FPIEEE"
	// The name of the builtin type for power pc specific floating point types.
	BUILTIN_TYPE_NAME_FLOAT_PPC = "Builtin.FPPPC"
	/// The name of the Builtin type for NativeObject
	BUILTIN_TYPE_NAME_NATIVEOBJECT = "Builtin.NativeObject"
	/// The name of the Builtin type for BridgeObject
	BUILTIN_TYPE_NAME_BRIDGEOBJECT = "Builtin.BridgeObject"
	/// The name of the Builtin type for RawPointer
	BUILTIN_TYPE_NAME_RAWPOINTER = "Builtin.RawPointer"
	/// The name of the Builtin type for UnsafeValueBuffer
	BUILTIN_TYPE_NAME_UNSAFEVALUEBUFFER = "Builtin.UnsafeValueBuffer"
	/// The name of the Builtin type for UnknownObject
	///
	/// This no longer exists as an AST-accessible type, but it's still used for
	/// fields shaped like AnyObject when ObjC interop is enabled.
	BUILTIN_TYPE_NAME_UNKNOWNOBJECT = "Builtin.UnknownObject"
	/// The name of the Builtin type for Vector
	BUILTIN_TYPE_NAME_VEC = "Builtin.Vec"
	/// The name of the Builtin type for SILToken
	BUILTIN_TYPE_NAME_SILTOKEN = "Builtin.SILToken"
	/// The name of the Builtin type for Word
	BUILTIN_TYPE_NAME_WORD = "Builtin.Word"
)

//go:generate stringer -type SpecialPointerAuthDiscriminators,NecessaryBindingsKind -output swift_string.go

type SpecialPointerAuthDiscriminators uint16

const (
	// All of these values are the stable string hash of the corresponding
	// variable name:
	//   (computeStableStringHash % 65535 + 1)

	/// HeapMetadataHeader::destroy
	HeapDestructor SpecialPointerAuthDiscriminators = 0xbbbf

	/// Type descriptor data pointers.
	TypeDescriptor SpecialPointerAuthDiscriminators = 0xae86

	/// Runtime function variables exported by the runtime.
	RuntimeFunctionEntry SpecialPointerAuthDiscriminators = 0x625b

	/// Protocol conformance descriptors.
	ProtocolConformanceDescriptor SpecialPointerAuthDiscriminators = 0xc6eb

	/// Pointer to value witness table stored in type metadata.
	///
	/// Computed with ptrauth_string_discriminator("value_witness_table_t").
	ValueWitnessTable SpecialPointerAuthDiscriminators = 0x2e3f

	/// Extended existential type shapes.
	ExtendedExistentialTypeShape          SpecialPointerAuthDiscriminators = 0x5a3d // SpecialPointerAuthDiscriminators = 23101
	NonUniqueExtendedExistentialTypeShape SpecialPointerAuthDiscriminators = 0xe798 // SpecialPointerAuthDiscriminators = 59288

	/// Value witness functions.
	InitializeBufferWithCopyOfBuffer   SpecialPointerAuthDiscriminators = 0xda4a
	Destroy                            SpecialPointerAuthDiscriminators = 0x04f8
	InitializeWithCopy                 SpecialPointerAuthDiscriminators = 0xe3ba
	AssignWithCopy                     SpecialPointerAuthDiscriminators = 0x8751
	InitializeWithTake                 SpecialPointerAuthDiscriminators = 0x48d8
	AssignWithTake                     SpecialPointerAuthDiscriminators = 0xefda
	DestroyArray                       SpecialPointerAuthDiscriminators = 0x2398
	InitializeArrayWithCopy            SpecialPointerAuthDiscriminators = 0xa05c
	InitializeArrayWithTakeFrontToBack SpecialPointerAuthDiscriminators = 0x1c3e
	InitializeArrayWithTakeBackToFront SpecialPointerAuthDiscriminators = 0x8dd3
	StoreExtraInhabitant               SpecialPointerAuthDiscriminators = 0x79c5
	GetExtraInhabitantIndex            SpecialPointerAuthDiscriminators = 0x2ca8
	GetEnumTag                         SpecialPointerAuthDiscriminators = 0xa3b5
	DestructiveProjectEnumData         SpecialPointerAuthDiscriminators = 0x041d
	DestructiveInjectEnumTag           SpecialPointerAuthDiscriminators = 0xb2e4
	GetEnumTagSinglePayload            SpecialPointerAuthDiscriminators = 0x60f0
	StoreEnumTagSinglePayload          SpecialPointerAuthDiscriminators = 0xa0d1

	/// KeyPath metadata functions.
	KeyPathDestroy           SpecialPointerAuthDiscriminators = 0x7072
	KeyPathCopy              SpecialPointerAuthDiscriminators = 0x6f66
	KeyPathEquals            SpecialPointerAuthDiscriminators = 0x756e
	KeyPathHash              SpecialPointerAuthDiscriminators = 0x6374
	KeyPathGetter            SpecialPointerAuthDiscriminators = 0x6f72
	KeyPathNonmutatingSetter SpecialPointerAuthDiscriminators = 0x6f70
	KeyPathMutatingSetter    SpecialPointerAuthDiscriminators = 0x7469
	KeyPathGetLayout         SpecialPointerAuthDiscriminators = 0x6373
	KeyPathInitializer       SpecialPointerAuthDiscriminators = 0x6275
	KeyPathMetadataAccessor  SpecialPointerAuthDiscriminators = 0x7474

	/// ObjC bridging entry points.
	ObjectiveCTypeDiscriminator                    SpecialPointerAuthDiscriminators = 0x31c3 // SpecialPointerAuthDiscriminators = 12739
	bridgeToObjectiveCDiscriminator                SpecialPointerAuthDiscriminators = 0xbca0 // SpecialPointerAuthDiscriminators = 48288
	forceBridgeFromObjectiveCDiscriminator         SpecialPointerAuthDiscriminators = 0x22fb // SpecialPointerAuthDiscriminators = 8955
	conditionallyBridgeFromObjectiveCDiscriminator SpecialPointerAuthDiscriminators = 0x9a9b // SpecialPointerAuthDiscriminators = 39579

	/// Dynamic replacement pointers.
	DynamicReplacementScope SpecialPointerAuthDiscriminators = 0x48F0 // SpecialPointerAuthDiscriminators = 18672
	DynamicReplacementKey   SpecialPointerAuthDiscriminators = 0x2C7D // SpecialPointerAuthDiscriminators = 11389

	/// Resume functions for yield-once coroutines that yield a single
	/// opaque borrowed/inout value.  These aren't actually hard-coded, but
	/// they're important enough to be worth writing in one place.
	OpaqueReadResumeFunction   SpecialPointerAuthDiscriminators = 56769
	OpaqueModifyResumeFunction SpecialPointerAuthDiscriminators = 3909

	/// ObjC class pointers.
	ObjCISA        SpecialPointerAuthDiscriminators = 0x6AE1
	ObjCSuperclass SpecialPointerAuthDiscriminators = 0xB5AB

	/// Resilient class stub initializer callback
	ResilientClassStubInitCallback SpecialPointerAuthDiscriminators = 0xC671

	/// Jobs, tasks, and continuations.
	JobInvokeFunction                SpecialPointerAuthDiscriminators = 0xcc64 // SpecialPointerAuthDiscriminators = 52324
	TaskResumeFunction               SpecialPointerAuthDiscriminators = 0x2c42 // SpecialPointerAuthDiscriminators = 11330
	TaskResumeContext                SpecialPointerAuthDiscriminators = 0x753a // SpecialPointerAuthDiscriminators = 30010
	AsyncRunAndBlockFunction         SpecialPointerAuthDiscriminators = 0x0f08 // 3848
	AsyncContextParent               SpecialPointerAuthDiscriminators = 0xbda2 // SpecialPointerAuthDiscriminators = 48546
	AsyncContextResume               SpecialPointerAuthDiscriminators = 0xd707 // SpecialPointerAuthDiscriminators = 55047
	AsyncContextYield                SpecialPointerAuthDiscriminators = 0xe207 // SpecialPointerAuthDiscriminators = 57863
	CancellationNotificationFunction SpecialPointerAuthDiscriminators = 0x1933 // SpecialPointerAuthDiscriminators = 6451
	EscalationNotificationFunction   SpecialPointerAuthDiscriminators = 0x5be4 // SpecialPointerAuthDiscriminators = 23524
	AsyncThinNullaryFunction         SpecialPointerAuthDiscriminators = 0x0f08 // SpecialPointerAuthDiscriminators = 3848
	AsyncFutureFunction              SpecialPointerAuthDiscriminators = 0x720f // SpecialPointerAuthDiscriminators = 29199

	/// Swift async context parameter stored in the extended frame info.
	SwiftAsyncContextExtendedFrameEntry SpecialPointerAuthDiscriminators = 0xc31a // SpecialPointerAuthDiscriminators = 49946

	// C type TaskContinuationFunction* descriminator.
	ClangTypeTaskContinuationFunction SpecialPointerAuthDiscriminators = 0x2abe // SpecialPointerAuthDiscriminators = 10942

	/// Dispatch integration.
	DispatchInvokeFunction SpecialPointerAuthDiscriminators = 0xf493 // SpecialPointerAuthDiscriminators = 62611

	/// Functions accessible at runtime (i.e. distributed method accessors).
	AccessibleFunctionRecord SpecialPointerAuthDiscriminators = 0x438c // = 17292
)

// __TEXT.__swift5_assocty
// This section contains an array of associated type descriptors.
// An associated type descriptor contains a collection of associated type records for a conformance.
// An associated type records describe the mapping from an associated type to the type witness of a conformance.

type AssociatedTypeRecord struct {
	Name                string
	SubstitutedTypeName string
	SubstitutedTypeAddr uint64
	ATRecordType
}
type ATRecordType struct {
	NameOffset                int32
	SubstitutedTypeNameOffset int32
}

type ATDHeader struct {
	ConformingTypeNameOffset int32
	ProtocolTypeNameOffset   int32
	NumAssociatedTypes       uint32
	AssociatedTypeRecordSize uint32
}
type AssociatedTypeDescriptor struct {
	ATDHeader
	Address               uint64
	ConformingTypeAddr    uint64
	ConformingTypeName    string
	ProtocolTypeName      string
	AssociatedTypeRecords []AssociatedTypeRecord
}

func (a AssociatedTypeDescriptor) String() string {
	var vars []string
	for _, v := range a.AssociatedTypeRecords {
		vars = append(vars, fmt.Sprintf("\t%s: %s", v.Name, v.SubstitutedTypeName))
	}
	return fmt.Sprintf(
		"extension %s: %s {\n"+
			"%s\n"+
			"}",
		a.ConformingTypeName,
		a.ProtocolTypeName,
		strings.Join(vars, "\n"),
	)
}

// __TEXT.__swift5_builtin
// This section contains an array of builtin type descriptors.
// A builtin type descriptor describes the basic layout information about any builtin types referenced from other sections.

type builtinTypeFlag uint32

func (f builtinTypeFlag) IsBitwiseTakable() bool {
	return ((f >> 16) & 1) != 0
}
func (f builtinTypeFlag) Alignment() uint16 {
	return uint16(f & 0xffff)
}

const MaxNumExtraInhabitants = 0x7FFFFFFF

type BuiltinTypeDescriptor struct {
	TypeName            int32
	Size                uint32
	AlignmentAndFlags   builtinTypeFlag
	Stride              uint32
	NumExtraInhabitants uint32
}

// BuiltinType builtin swift type
type BuiltinType struct {
	Address             uint64
	Name                string
	Size                uint32
	Alignment           uint16
	BitwiseTakable      bool
	Stride              uint32
	NumExtraInhabitants uint32
}

func (b BuiltinType) String() string {
	var numExtraInhabitants string
	if b.NumExtraInhabitants == MaxNumExtraInhabitants {
		numExtraInhabitants = "max"
	} else {
		numExtraInhabitants = fmt.Sprintf("%d", b.NumExtraInhabitants)
	}
	return fmt.Sprintf(
		"%#x:         %s\n"+
			"  size:              %d\n"+
			"  alignment:         %d\n"+
			"  bitwise-takable:   %t\n"+
			"  stride:            %d\n"+
			"  extra-inhabitants: %s\n",
		b.Address, b.Name, b.Size, b.Alignment, b.BitwiseTakable, b.Stride, numExtraInhabitants)
}

// __TEXT.__swift5_capture
// Capture descriptors describe the layout of a closure context object.
// Unlike nominal types, the generic substitutions for a closure context come from the object, and not the metadata.

type CaptureTypeRecord struct {
	MangledTypeName int32
}

type MetadataSourceRecord struct {
	MangledTypeName       int32
	MangledMetadataSource int32
}

type MetadataSource struct {
	MangledType           string
	MangledMetadataSource string
}

type NecessaryBindingsKind uint32

const (
	PartialApply NecessaryBindingsKind = iota
	AsyncFunction
)

type NecessaryBindings struct {
	Kind               NecessaryBindingsKind
	RequirementsSet    int32
	RequirementsVector int32
	Conformances       int32
}

type CaptureDescriptorHeader struct {
	NumCaptureTypes    uint32 // The number of captures in the closure and the number of typerefs that immediately follow this struct.
	NumMetadataSources uint32 // The number of sources of metadata available in the MetadataSourceMap directly following the list of capture's typerefs.
	NumBindings        uint32 // The number of items in the NecessaryBindings structure at the head of the closure.
}

type CaptureDescriptor struct {
	Address uint64
	CaptureDescriptorHeader
	CaptureTypes    []string
	MetadataSources []MetadataSource
	Bindings        []NecessaryBindings
}

func (c CaptureDescriptor) String() string {
	var captureTypes string
	if len(c.CaptureTypes) > 0 {
		captureTypes += "\t/* capture types */\n"
		for _, t := range c.CaptureTypes {
			captureTypes += fmt.Sprintf("\t%s\n", t)
		}
	}
	var metadataSources string
	if len(c.MetadataSources) > 0 {
		metadataSources += "\t/* metadata sources */\n"
		for _, m := range c.MetadataSources {
			metadataSources += fmt.Sprintf("\t%s: %s\n", m.MangledType, m.MangledMetadataSource)
		}
	}
	var bindings string
	if len(c.Bindings) > 0 {
		bindings += "\t/* necessary bindings */\n"
		for _, b := range c.Bindings {
			bindings += fmt.Sprintf("\t// Kind: %d, RequirementsSet: %d, RequirementsVector: %d, Conformances: %d\n", b.Kind, b.RequirementsSet, b.RequirementsVector, b.Conformances)
		}
	}
	return fmt.Sprintf(
		"block /* %#x */ {\n"+
			"%s"+
			"%s"+
			"%s"+
			"}",
		c.Address,
		captureTypes,
		metadataSources,
		bindings,
	)
}

// __TEXT.__swift5_typeref
// This section contains a list of mangled type names that are referenced from other sections.
// This is essentially all the different types that are used in the application.
// The Swift docs and code are the best places to find out more information about mangled type names.

// __TEXT.__swift5_reflstr
// This section contains an array of C strings. The strings are field names for the properties of the metadata defined in other sections.
