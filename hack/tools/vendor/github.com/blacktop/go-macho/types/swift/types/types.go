package types

import (
	"fmt"
	"strings"

	"github.com/blacktop/go-macho/types"
	"github.com/blacktop/go-macho/types/swift"
	"github.com/blacktop/go-macho/types/swift/fields"
)

//go:generate stringer -type ContextDescriptorKind,TypeReferenceKind,MetadataInitializationKind -linecomment -output types_string.go

// __TEXT.__swift5_types
// This section contains an array of 32-bit signed integers.
// Each integer is a relative offset that points to a nominal type descriptor in the __TEXT.__const section.

type ContextDescriptorKind uint8

const (
	// This context descriptor represents a module.
	CDKindModule ContextDescriptorKind = 0 // module

	/// This context descriptor represents an extension.
	CDKindExtension ContextDescriptorKind = 1 // extension

	/// This context descriptor represents an anonymous possibly-generic context
	/// such as a function body.
	CDKindAnonymous ContextDescriptorKind = 2 // anonymous

	/// This context descriptor represents a protocol context.
	CDKindProtocol ContextDescriptorKind = 3 // protocol

	/// This context descriptor represents an opaque type alias.
	CDKindOpaqueType ContextDescriptorKind = 4 // opaque_type

	/// First kind that represents a type of any sort.
	CDKindTypeFirst = 16 // type_first

	/// This context descriptor represents a class.
	CDKindClass ContextDescriptorKind = CDKindTypeFirst // class

	/// This context descriptor represents a struct.
	CDKindStruct ContextDescriptorKind = CDKindTypeFirst + 1 // struct

	/// This context descriptor represents an enum.
	CDKindEnum ContextDescriptorKind = CDKindTypeFirst + 2 // enum

	/// Last kind that represents a type of any sort.
	CDKindTypeLast = 31 // type_last
)

type TypeContextDescriptorFlags uint16

const (
	// All of these values are bit offsets or widths.
	// Generic flags build upwards from 0.
	// Type-specific flags build downwards from 15.

	/// Whether there's something unusual about how the metadata is
	/// initialized.
	///
	/// Meaningful for all type-descriptor kinds.
	MetadataInitialization       TypeContextDescriptorFlags = 0
	MetadataInitialization_width TypeContextDescriptorFlags = 2

	/// Set if the type has extended import information.
	///
	/// If true, a sequence of strings follow the null terminator in the
	/// descriptor, terminated by an empty string (i.e. by two null
	/// terminators in a row).  See TypeImportInfo for the details of
	/// these strings and the order in which they appear.
	///
	/// Meaningful for all type-descriptor kinds.
	HasImportInfo TypeContextDescriptorFlags = 2

	/// Set if the type descriptor has a pointer to a list of canonical
	/// prespecializations.
	HasCanonicalMetadataPrespecializations TypeContextDescriptorFlags = 3

	// Type-specific flags:

	/// Set if the class is an actor.
	///
	/// Only meaningful for class descriptors.
	Class_IsActor TypeContextDescriptorFlags = 7

	/// Set if the class is a default actor class.  Note that this is
	/// based on the best knowledge available to the class; actor
	/// classes with resilient superclassess might be default actors
	/// without knowing it.
	///
	/// Only meaningful for class descriptors.
	Class_IsDefaultActor TypeContextDescriptorFlags = 8

	/// The kind of reference that this class makes to its resilient superclass
	/// descriptor.  A TypeReferenceKind.
	///
	/// Only meaningful for class descriptors.
	Class_ResilientSuperclassReferenceKind       TypeContextDescriptorFlags = 9
	Class_ResilientSuperclassReferenceKind_width TypeContextDescriptorFlags = 3

	/// Whether the immediate class members in this metadata are allocated
	/// at negative offsets.  For now, we don't use this.
	Class_AreImmediateMembersNegative TypeContextDescriptorFlags = 12

	/// Set if the context descriptor is for a class with resilient ancestry.
	///
	/// Only meaningful for class descriptors.
	Class_HasResilientSuperclass TypeContextDescriptorFlags = 13

	/// Set if the context descriptor includes metadata for dynamically
	/// installing method overrides at metadata instantiation time.
	Class_HasOverrideTable TypeContextDescriptorFlags = 14

	/// Set if the context descriptor includes metadata for dynamically
	/// constructing a class's vtables at metadata instantiation time.
	///
	/// Only meaningful for class descriptors.
	Class_HasVTable TypeContextDescriptorFlags = 15
)

func (f TypeContextDescriptorFlags) MetadataInitialization() MetadataInitializationKind {
	return MetadataInitializationKind(types.ExtractBits(uint64(f), int32(MetadataInitialization), int32(MetadataInitialization_width)))
}
func (f TypeContextDescriptorFlags) HasImportInfo() bool {
	return types.ExtractBits(uint64(f), int32(HasImportInfo), 1) != 0
}
func (f TypeContextDescriptorFlags) IsActor() bool {
	return types.ExtractBits(uint64(f), int32(Class_IsActor), 1) != 0
}
func (f TypeContextDescriptorFlags) IsDefaultActor() bool {
	return types.ExtractBits(uint64(f), int32(Class_IsDefaultActor), 1) != 0
}
func (f TypeContextDescriptorFlags) ResilientSuperclassReferenceKind() TypeReferenceKind {
	return TypeReferenceKind(types.ExtractBits(uint64(f), int32(Class_ResilientSuperclassReferenceKind), int32(Class_ResilientSuperclassReferenceKind_width)))
}
func (f TypeContextDescriptorFlags) AreImmediateMembersNegative() bool {
	return types.ExtractBits(uint64(f), int32(Class_AreImmediateMembersNegative), 1) != 0
}
func (f TypeContextDescriptorFlags) HasResilientSuperclass() bool {
	return types.ExtractBits(uint64(f), int32(Class_HasResilientSuperclass), 1) != 0
}
func (f TypeContextDescriptorFlags) HasOverrideTable() bool {
	return types.ExtractBits(uint64(f), int32(Class_HasOverrideTable), 1) != 0
}
func (f TypeContextDescriptorFlags) HasVTable() bool {
	return types.ExtractBits(uint64(f), int32(Class_HasVTable), 1) != 0
}
func (f TypeContextDescriptorFlags) String() string {
	var flags []string
	if f.MetadataInitialization() != MetadataInitNone {
		flags = append(flags, fmt.Sprintf("metadata_init:%s", f.MetadataInitialization()))
	}
	if f.HasImportInfo() {
		flags = append(flags, "import_info")
	}
	if f.IsActor() {
		flags = append(flags, "actor")
	}
	if f.IsDefaultActor() {
		flags = append(flags, "default_actor")
	}
	if f.AreImmediateMembersNegative() {
		flags = append(flags, "negative_immediate_members")
	}
	if f.HasResilientSuperclass() {
		flags = append(flags, "resilient_superclass")
		flags = append(flags, fmt.Sprintf("resilient_superclass_ref:%s", f.ResilientSuperclassReferenceKind()))
	}
	if f.HasOverrideTable() {
		flags = append(flags, "override_table")
	}
	if f.HasVTable() {
		flags = append(flags, "vtable")
	}
	return strings.Join(flags, "|")
}

// TypeReferenceKind kinds of type metadata/protocol conformance records.
type TypeReferenceKind uint8

const (
	//The conformance is for a nominal type referenced directly; getTypeDescriptor() points to the type context descriptor.
	DirectTypeDescriptor TypeReferenceKind = 0x00
	// The conformance is for a nominal type referenced indirectly; getTypeDescriptor() points to the type context descriptor.
	IndirectTypeDescriptor TypeReferenceKind = 0x01
	// The conformance is for an Objective-C class that should be looked up by class name.
	DirectObjCClassName TypeReferenceKind = 0x02
	// The conformance is for an Objective-C class that has no nominal type descriptor.
	// getIndirectObjCClass() points to a variable that contains the pointer to
	// the class object, which then requires a runtime call to get metadata.
	//
	// On platforms without Objective-C interoperability, this case is unused.
	IndirectObjCClass TypeReferenceKind = 0x03
	// We only reserve three bits for this in the various places we store it.
	FirstKind = DirectTypeDescriptor
	LastKind  = IndirectObjCClass
)

type MetadataInitializationKind uint8

const (
	// There are either no special rules for initializing the metadata or the metadata is generic.
	// (Genericity is set in the non-kind-specific descriptor flags.)
	MetadataInitNone MetadataInitializationKind = 0 // none
	//The type requires non-trivial singleton initialization using the "in-place" code pattern.
	MetadataInitSingleton MetadataInitializationKind = 1 // singleton
	// The type requires non-trivial singleton initialization using the "foreign" code pattern.
	MetadataInitForeign MetadataInitializationKind = 2 // foreign
	// We only have two bits here, so if you add a third special kind, include more flag bits in its out-of-line storage.
)

type TargetSingletonMetadataInitialization struct {
	InitializationCacheOffset int32 // The initialization cache. Out-of-line because mutable.
	IncompleteMetadata        int32 // UNION: The incomplete metadata, for structs, enums and classes without resilient ancestry.
	// ResilientPattern
	// If the class descriptor's hasResilientSuperclass() flag is set,
	// this field instead points at a pattern used to allocate and
	// initialize metadata for this class, since it's size and contents
	// is not known at compile time.
	CompletionFunction int32 // The completion function. The pattern will always be null, even for a resilient class.
}

type TargetForeignMetadataInitialization struct {
	CompletionFunction int32 // The completion function. The pattern will always be null.
}

type ContextDescriptorFlags uint32

func (f ContextDescriptorFlags) Kind() ContextDescriptorKind {
	return ContextDescriptorKind(f & 0x1F)
}
func (f ContextDescriptorFlags) IsGeneric() bool {
	return (f & 0x80) != 0
}
func (f ContextDescriptorFlags) IsUnique() bool {
	return (f & 0x40) != 0
}
func (f ContextDescriptorFlags) Version() uint8 {
	return uint8(f >> 8 & 0xFF)
}
func (f ContextDescriptorFlags) KindSpecific() TypeContextDescriptorFlags {
	return TypeContextDescriptorFlags((f >> 16) & 0xFFFF)
}
func (f ContextDescriptorFlags) String() string {
	return fmt.Sprintf("kind: %s, generic: %t, unique: %t, version: %d, kind_flags: %s",
		f.Kind(),
		f.IsGeneric(),
		f.IsUnique(),
		f.Version(),
		f.KindSpecific())
}

type Type struct {
	Address uint64
	Name    string
}

type TypeDescriptor struct {
	Address        uint64
	Parent         Type
	Name           string
	Kind           ContextDescriptorKind
	AccessFunction uint64
	FieldOffsets   []int32
	Generic        *TargetTypeGenericContextDescriptorHeader
	VTable         *VTable
	Fields         []*fields.Field
	Type           any
}

func (t TypeDescriptor) String() string {
	switch t.Kind {
	case CDKindModule:
		return fmt.Sprintf("module %s", t.Name)
	case CDKindExtension:
		return fmt.Sprintf("extension %s", t.Name)
	case CDKindAnonymous:
		return fmt.Sprintf("anonymous %s", t.Name)
	case CDKindProtocol:
		return fmt.Sprintf("protocol %s", t.Name)
	case CDKindOpaqueType:
		return fmt.Sprintf("opaque type %s", t.Name)
	case CDKindClass:
		return fmt.Sprintf("class %s", t.Name)
	case CDKindStruct:
		return fmt.Sprintf("struct %s", t.Name)
	case CDKindEnum:
		return fmt.Sprintf("enum %s", t.Name)
	default:
		return fmt.Sprintf("unknown type %s", t.Name)
	}
}

type TargetContextDescriptor struct {
	Flags        ContextDescriptorFlags
	ParentOffset int32
}

type TargetModuleContext struct {
	Name string
	TargetModuleContextDescriptor
}

type TargetModuleContextDescriptor struct {
	TargetContextDescriptor
	NameOffset int32
}

func (t TypeDescriptor) IsCImportedModuleName() bool {
	if t.Kind == CDKindModule {
		return t.Name == swift.MANGLING_MODULE_OBJC
	}
	return false
}

type TargetTypeContextDescriptor struct {
	TargetContextDescriptor
	NameOffset        int32 // The name of the type.
	AccessFunctionPtr int32 // The access function for the type.
	FieldsOffset      int32 // A pointer to the field descriptor for the type, if any.
}

type TargetExtensionContextDescriptor struct {
	TargetContextDescriptor
	ExtendedContext int32
}

type TargetAnonymousContextDescriptor struct {
	TargetContextDescriptor
}

type TargetEnumDescriptor struct {
	TargetTypeContextDescriptor
	NumPayloadCasesAndPayloadSizeOffset uint32
	NumEmptyCases                       uint32
}

func (e TargetEnumDescriptor) GetNumPayloadCases() uint32 {
	return e.NumPayloadCasesAndPayloadSizeOffset & 0x00FFFFFF
}
func (e TargetEnumDescriptor) GetNumCases() uint32 {
	return e.GetNumPayloadCases() + e.NumEmptyCases
}
func (e TargetEnumDescriptor) GetPayloadSizeOffset() uint32 {
	return (e.NumPayloadCasesAndPayloadSizeOffset & 0xFF000000) >> 24
}

type TargetStructDescriptor struct {
	TargetTypeContextDescriptor
	NumFields               uint32
	FieldOffsetVectorOffset uint32
}

type TargetProtocolDescriptor struct {
	TargetContextDescriptor
	NameOffset                 int32  // The name of the protocol.
	NumRequirementsInSignature uint32 // The number of generic requirements in the requirement signature of the protocol.
	NumRequirements            uint32 /* The number of requirements in the protocol. If any requirements beyond MinimumWitnessTableSizeInWords are present
	 * in the witness table template, they will be not be overwritten with defaults. */
	AssociatedTypeNamesOffset int32 // Associated type names, as a space-separated list in the same order as the requirements.
}

type TargetOpaqueTypeDescriptor struct {
	TargetContextDescriptor
}

type TargetClassDescriptor struct {
	TargetTypeContextDescriptor
	SuperclassType              int32
	MetadataNegativeSizeInWords uint32
	MetadataPositiveSizeInWords uint32
	NumImmediateMembers         uint32
	NumFields                   uint32
	FieldOffsetVectorOffset     uint32
}

type TargetTypeGenericContextDescriptorHeader struct {
	InstantiationCache          int32
	DefaultInstantiationPattern int32
	Base                        TargetGenericContextDescriptorHeader
}

type TargetGenericContextDescriptorHeader struct {
	NumParams         uint16
	NumRequirements   uint16
	NumKeyArguments   uint16
	NumExtraArguments uint16
}

func (g TargetGenericContextDescriptorHeader) GetNumArguments() uint32 {
	return uint32(g.NumKeyArguments + g.NumExtraArguments)
}
func (g TargetGenericContextDescriptorHeader) GetArgumentLayoutSizeInWords() uint32 {
	return g.GetNumArguments()
}
func (g TargetGenericContextDescriptorHeader) HasArguments() bool {
	return g.GetNumArguments() > 0
}

type VTable struct {
	TargetVTableDescriptorHeader
	MethodListOffset int64
	Methods          []TargetMethodDescriptor
}

type method struct {
	Flags   MethodDescriptorFlags
	Address uint64
}

func (v *VTable) GetMethods(base uint64) []method {
	var methods []method
	for idx, m := range v.Methods {
		methods = append(methods, method{
			Flags:   m.Flags,
			Address: base + uint64(v.MethodListOffset+int64(m.Impl)+int64((idx*8)+4)),
		})
	}
	return methods
}

type TargetVTableDescriptorHeader struct {
	VTableOffset uint32
	VTableSize   uint32
}

type TargetMethodDescriptor struct {
	Flags MethodDescriptorFlags
	Impl  int32
}

type mdKind uint8

const (
	Method mdKind = iota
	Init
	Getter
	Setter
	ModifyCoroutine
	ReadCoroutine
)

const (
	ExtraDiscriminatorShift = 16
	ExtraDiscriminatorMask  = 0xFFFF0000
)

type MethodDescriptorFlags uint32

func (f MethodDescriptorFlags) Kind() string {
	switch mdKind(f & 0x0F) {
	case Method:
		return "method"
	case Init:
		return "init"
	case Getter:
		return "getter"
	case Setter:
		return "setter"
	case ModifyCoroutine:
		return "modify coroutine"
	case ReadCoroutine:
		return "read coroutine"
	default:
		return fmt.Sprintf("unknown kind %d", mdKind(f&0x0F))
	}
}
func (f MethodDescriptorFlags) IsInstance() bool {
	return (f & 0x10) != 0
}
func (f MethodDescriptorFlags) IsDynamic() bool {
	return (f & 0x20) != 0
}
func (f MethodDescriptorFlags) IsAsync() bool {
	return (f & 0x40) != 0
}
func (f MethodDescriptorFlags) ExtraDiscriminator() uint16 {
	return uint16(f >> ExtraDiscriminatorShift)
}
func (f MethodDescriptorFlags) String() string {
	var flags []string
	if f.IsInstance() {
		flags = append(flags, "instance")
	}
	if f.IsDynamic() {
		flags = append(flags, "dynamic")
	}
	if f.IsAsync() {
		flags = append(flags, "async")
	}
	if f.ExtraDiscriminator() != 0 {
		flags = append(flags, fmt.Sprintf("extra discriminator %#x", f.ExtraDiscriminator()))
	}
	if len(strings.Join(flags, "|")) == 0 {
		return f.Kind()
	}
	return fmt.Sprintf("%s (%s)", f.Kind(), strings.Join(flags, "|"))
}

// __swift5_acfuncs

type AccessibleFunctionsSection struct {
	Begin uint64 // AccessibleFunctionRecord
	End   uint64 // AccessibleFunctionRecord
}

type AccessibleFunctionFlags uint32

const (
	Distributed AccessibleFunctionFlags = 0
)

type TargetAccessibleFunctionRecord struct {
	Name               int32 // char *
	GenericEnvironment int32 // TargetGenericEnvironment
	FunctionType       int32 // mangled name
	Function           int32 // void *
	Flags              AccessibleFunctionFlags
}

type GenericEnvironmentFlags uint32

func (f GenericEnvironmentFlags) GetNumGenericParameterLevels() uint32 {
	return uint32(f & 0xFFF)
}
func (f GenericEnvironmentFlags) GetNumGenericRequirements() uint32 {
	return uint32((f & (0xFFFF << 12)) >> 12)
}

type TargetGenericEnvironment struct {
	Flags GenericEnvironmentFlags
}

type AccessibleFunctionCacheEntry struct {
	Name    string
	NameLen uint32
	R       uint64 // AccessibleFunctionRecord

}

type AccessibleFunctionsState struct {
	Cache          AccessibleFunctionCacheEntry
	SectionsToScan AccessibleFunctionsSection
}
