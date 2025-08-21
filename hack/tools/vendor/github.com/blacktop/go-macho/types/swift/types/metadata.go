package types

//go:generate stringer -type MetadataKind -linecomment

const (
	// Non-type metadata kinds have this bit set.
	MetadataKindIsNonType = 0x400
	// Non-heap metadata kinds have this bit set.
	MetadataKindIsNonHeap = 0x200
	// The above two flags are negative because the "class" kind has to be zero,
	// and class metadata is both type and heap metadata.
	// Runtime-private metadata has this bit set. The compiler must not statically
	// generate metadata objects with these kinds, and external tools should not
	// rely on the stability of these values or the precise binary layout of
	// their associated data structures.
	MetadataKindIsRuntimePrivate = 0x100
)

type MetadataKind uint32

const (
	ClassMetadataKind                    MetadataKind = 0                                                        // class
	StructMetadataKind                   MetadataKind = 0 | MetadataKindIsNonHeap                                // struct
	EnumMetadataKind                     MetadataKind = 1 | MetadataKindIsNonHeap                                // enum
	OptionalMetadataKind                 MetadataKind = 2 | MetadataKindIsNonHeap                                // optional
	ForeignClassMetadataKind             MetadataKind = 3 | MetadataKindIsNonHeap                                // foreign class
	ForeignReferenceTypeMetadataKind     MetadataKind = 4 | MetadataKindIsNonHeap                                // foreign reference type
	OpaqueMetadataKind                   MetadataKind = 0 | MetadataKindIsRuntimePrivate | MetadataKindIsNonHeap // opaque
	TupleMetadataKind                    MetadataKind = 1 | MetadataKindIsRuntimePrivate | MetadataKindIsNonHeap // tuple
	FunctionMetadataKind                 MetadataKind = 2 | MetadataKindIsRuntimePrivate | MetadataKindIsNonHeap // function
	ExistentialMetadataKind              MetadataKind = 3 | MetadataKindIsRuntimePrivate | MetadataKindIsNonHeap // existential
	MetatypeMetadataKind                 MetadataKind = 4 | MetadataKindIsRuntimePrivate | MetadataKindIsNonHeap // metatype
	ObjCClassWrapperMetadataKind         MetadataKind = 5 | MetadataKindIsRuntimePrivate | MetadataKindIsNonHeap // objc class wrapper
	ExistentialMetatypeMetadataKind      MetadataKind = 6 | MetadataKindIsRuntimePrivate | MetadataKindIsNonHeap // existential metatype
	ExtendedExistentialMetadataKind      MetadataKind = 7 | MetadataKindIsRuntimePrivate | MetadataKindIsNonHeap // extended existential type
	HeapLocalVariableMetadataKind        MetadataKind = 0 | MetadataKindIsNonType                                // heap local variable
	HeapGenericLocalVariableMetadataKind MetadataKind = 0 | MetadataKindIsNonType | MetadataKindIsRuntimePrivate // heap generic local variable
	ErrorObjectMetadataKind              MetadataKind = 1 | MetadataKindIsNonType | MetadataKindIsRuntimePrivate // error object
	TaskMetadataKind                     MetadataKind = 2 | MetadataKindIsNonType | MetadataKindIsRuntimePrivate // task
	JobMetadataKind                      MetadataKind = 3 | MetadataKindIsNonType | MetadataKindIsRuntimePrivate // job
	// The largest possible non-isa-pointer metadata kind value.
	LastEnumerated MetadataKind = 0x7FF
	// This is included in the enumeration to prevent against attempts to
	// exhaustively match metadata kinds. Future Swift runtimes or compilers
	// may introduce new metadata kinds, so for forward compatibility, the
	// runtime must tolerate metadata with unknown kinds.
	// This specific value is not mapped to a valid metadata kind at this time,
	// however.
)
