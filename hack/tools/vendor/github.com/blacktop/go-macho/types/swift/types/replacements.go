package types

// __TEXT.__swift5_replace
// This section contains dynamic replacement information.
// This is essentially the Swift equivalent of Objective-C method swizzling.

type AutomaticDynamicReplacements struct {
	Flags     uint32
	NumScopes uint32
	AutomaticDynamicReplacementEntry
}

type AutomaticDynamicReplacementEntry struct {
	ReplacementScope int32 // DynamicReplacementScope
	Flags            uint32
}

type DynamicReplacementKey struct {
	Root  int32
	Flags uint32
}

func (d DynamicReplacementKey) ExtraDiscriminator() uint16 {
	return uint16(d.Flags & 0x0000FFFF)
}
func (d DynamicReplacementKey) IsAsync() bool {
	return ((d.Flags >> 16) & 0x1) != 0
}

type DynamicReplacementScope struct {
	Flags           uint32
	NumReplacements uint32 // hard coded to 1
	DynamicReplacementDescriptor
}

const EnableChainingMask = 0x1

type DynamicReplacementDescriptor struct {
	ReplacedFunctionKey int32 // DynamicReplacementKey
	ReplacementFunction int32 // UNION w/ ReplacementAsyncFunction - TargetCompactFunctionPointer|TargetRelativeDirectPointer
	ChainEntry          int32 // DynamicReplacementChainEntry
	Flags               uint32
}

func (d DynamicReplacementDescriptor) ShouldChain() bool {
	return (d.Flags & EnableChainingMask) != 0
}

type DynamicReplacementChainEntry struct {
	ImplementationFunction uint64 // void *
	Next                   uint64 // DynamicReplacementChainEntry *
}

// __TEXT.__swift5_replac2
// This section contains dynamic replacement information for opaque types.

type DynamicReplacementSomeDescriptor struct {
	OriginalOpaqueTypeDesc    int32 // OpaqueTypeDescriptor -> TargetContextDescriptor
	ReplacementOpaqueTypeDesc int32 // OpaqueTypeDescriptor -> TargetContextDescriptor
}

type AutomaticDynamicReplacementsSome struct {
	Flags           uint32
	NumReplacements uint32
	Replacements    []DynamicReplacementSomeDescriptor
}
