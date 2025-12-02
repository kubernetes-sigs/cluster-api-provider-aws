package types

import (
	"fmt"
	"strings"
)

const (
	/* The following are used to encode rebasing information */
	REBASE_TYPE_POINTER                              = 1
	REBASE_TYPE_TEXT_ABSOLUTE32                      = 2
	REBASE_TYPE_TEXT_PCREL32                         = 3
	REBASE_OPCODE_MASK                               = 0xF0
	REBASE_IMMEDIATE_MASK                            = 0x0F
	REBASE_OPCODE_DONE                               = 0x00
	REBASE_OPCODE_SET_TYPE_IMM                       = 0x10
	REBASE_OPCODE_SET_SEGMENT_AND_OFFSET_ULEB        = 0x20
	REBASE_OPCODE_ADD_ADDR_ULEB                      = 0x30
	REBASE_OPCODE_ADD_ADDR_IMM_SCALED                = 0x40
	REBASE_OPCODE_DO_REBASE_IMM_TIMES                = 0x50
	REBASE_OPCODE_DO_REBASE_ULEB_TIMES               = 0x60
	REBASE_OPCODE_DO_REBASE_ADD_ADDR_ULEB            = 0x70
	REBASE_OPCODE_DO_REBASE_ULEB_TIMES_SKIPPING_ULEB = 0x80
)

type Rebase struct {
	Type    uint8
	Segment string
	Section string
	Start   uint64
	Offset  uint64
	Value   uint64
}

func (r Rebase) String() string {
	return fmt.Sprintf(
		"%-7s %-16s\t%#x  %s  %#x",
		r.Segment,
		r.Section,
		r.Start+r.Offset,
		getBindType(r.Type),
		r.Value,
	)
}

const (
	/* The following are used to encode binding information */
	BIND_TYPE_POINTER                                        = 1
	BIND_TYPE_TEXT_ABSOLUTE32                                = 2
	BIND_TYPE_TEXT_PCREL32                                   = 3
	BIND_TYPE_THREADED_BIND                                  = 100
	BIND_TYPE_THREADED_REBASE                                = 102
	BIND_SPECIAL_DYLIB_SELF                                  = 0
	BIND_SPECIAL_DYLIB_MAIN_EXECUTABLE                       = -1
	BIND_SPECIAL_DYLIB_FLAT_LOOKUP                           = -2
	BIND_SPECIAL_DYLIB_WEAK_LOOKUP                           = -3
	BIND_SYMBOL_FLAGS_WEAK_IMPORT                            = 0x1
	BIND_SYMBOL_FLAGS_NON_WEAK_DEFINITION                    = 0x8
	BIND_OPCODE_MASK                                         = 0xF0
	BIND_IMMEDIATE_MASK                                      = 0x0F
	BIND_OPCODE_DONE                                         = 0x00
	BIND_OPCODE_SET_DYLIB_ORDINAL_IMM                        = 0x10
	BIND_OPCODE_SET_DYLIB_ORDINAL_ULEB                       = 0x20
	BIND_OPCODE_SET_DYLIB_SPECIAL_IMM                        = 0x30
	BIND_OPCODE_SET_SYMBOL_TRAILING_FLAGS_IMM                = 0x40
	BIND_OPCODE_SET_TYPE_IMM                                 = 0x50
	BIND_OPCODE_SET_ADDEND_SLEB                              = 0x60
	BIND_OPCODE_SET_SEGMENT_AND_OFFSET_ULEB                  = 0x70
	BIND_OPCODE_ADD_ADDR_ULEB                                = 0x80
	BIND_OPCODE_DO_BIND                                      = 0x90
	BIND_OPCODE_DO_BIND_ADD_ADDR_ULEB                        = 0xA0
	BIND_OPCODE_DO_BIND_ADD_ADDR_IMM_SCALED                  = 0xB0
	BIND_OPCODE_DO_BIND_ULEB_TIMES_SKIPPING_ULEB             = 0xC0
	BIND_OPCODE_THREADED                                     = 0xD0
	BIND_SUBOPCODE_THREADED_SET_BIND_ORDINAL_TABLE_SIZE_ULEB = 0x00
	BIND_SUBOPCODE_THREADED_APPLY                            = 0x01
)

type BindKind uint8

const (
	BIND_KIND BindKind = iota
	WEAK_KIND
	LAZY_KIND
)

func (k BindKind) String() string {
	switch k {
	case BIND_KIND:
		return "BIND"
	case WEAK_KIND:
		return "WEAK"
	case LAZY_KIND:
		return "LAZY"
	}
	return ""
}

type Binds []Bind

func (bs Binds) Search(name string) (*Bind, error) {
	for _, b := range bs {
		if b.Name == name {
			return &b, nil
		}
	}
	return nil, fmt.Errorf("%s not found in bind info", name)
}

type Bind struct {
	Name    string
	Type    uint8
	Kind    BindKind
	Flags   uint8
	Addend  int64
	Segment string
	Section string
	Start   uint64
	Offset  uint64
	Dylib   string
	Value   uint64
}

func (b Bind) String() string {
	return fmt.Sprintf(
		"%-7s %-16s  %#x  %-4s  %-10s  %5d %-25s\t%s%s",
		b.Segment,
		b.Section,
		b.Start+b.Offset,
		b.Kind,
		getBindType(b.Type),
		b.Addend,
		b.Dylib,
		b.Name,
		getBindFlag(b.Flags, b.Kind),
	)
}

func getBindType(t uint8) string {
	switch t {
	case 0:
		return ""
	case BIND_TYPE_POINTER:
		return "pointer"
	case BIND_TYPE_TEXT_ABSOLUTE32:
		return "text abs32"
	case BIND_TYPE_TEXT_PCREL32:
		return "text rel32"
	case BIND_TYPE_THREADED_BIND:
		return "BIND_TYPE_THREADED_BIND"
	case BIND_TYPE_THREADED_REBASE:
		return "BIND_TYPE_THREADED_REBASE"
	}
	return fmt.Sprintf(" bad bind type %#02x", t)
}
func getBindFlag(f uint8, k BindKind) string {
	if f&BIND_SYMBOL_FLAGS_WEAK_IMPORT != 0 {
		return " (weak import)"
	} else if f&BIND_SYMBOL_FLAGS_NON_WEAK_DEFINITION != 0 {
		if k == WEAK_KIND {
			return " (strong)"
		} else {
			return ""
		}
	} else if f == 0 {
		return ""
	}

	return fmt.Sprintf("bad bind flag %#02x", f)
}

type SplitInfoKind uint64

const (
	DYLD_CACHE_ADJ_V2_FORMAT = 0x7F

	DYLD_CACHE_ADJ_V2_POINTER_32          SplitInfoKind = 0x01
	DYLD_CACHE_ADJ_V2_POINTER_64          SplitInfoKind = 0x02
	DYLD_CACHE_ADJ_V2_DELTA_32            SplitInfoKind = 0x03
	DYLD_CACHE_ADJ_V2_DELTA_64            SplitInfoKind = 0x04
	DYLD_CACHE_ADJ_V2_ARM64_ADRP          SplitInfoKind = 0x05
	DYLD_CACHE_ADJ_V2_ARM64_OFF12         SplitInfoKind = 0x06
	DYLD_CACHE_ADJ_V2_ARM64_BR26          SplitInfoKind = 0x07
	DYLD_CACHE_ADJ_V2_ARM_MOVW_MOVT       SplitInfoKind = 0x08
	DYLD_CACHE_ADJ_V2_ARM_BR24            SplitInfoKind = 0x09
	DYLD_CACHE_ADJ_V2_THUMB_MOVW_MOVT     SplitInfoKind = 0x0A
	DYLD_CACHE_ADJ_V2_THUMB_BR22          SplitInfoKind = 0x0B
	DYLD_CACHE_ADJ_V2_IMAGE_OFF_32        SplitInfoKind = 0x0C
	DYLD_CACHE_ADJ_V2_THREADED_POINTER_64 SplitInfoKind = 0x0D
)

func (k SplitInfoKind) String() string {
	switch k {
	case DYLD_CACHE_ADJ_V2_POINTER_32:
		return "pointer_32"
	case DYLD_CACHE_ADJ_V2_POINTER_64:
		return "pointer_64"
	case DYLD_CACHE_ADJ_V2_DELTA_32:
		return "delta_32"
	case DYLD_CACHE_ADJ_V2_DELTA_64:
		return "delta_64"
	case DYLD_CACHE_ADJ_V2_ARM64_ADRP:
		return "arm64_adrp"
	case DYLD_CACHE_ADJ_V2_ARM64_OFF12:
		return "arm64_off_12"
	case DYLD_CACHE_ADJ_V2_ARM64_BR26:
		return "arm64_br_26"
	case DYLD_CACHE_ADJ_V2_ARM_MOVW_MOVT:
		return "arm_movw_movt"
	case DYLD_CACHE_ADJ_V2_ARM_BR24:
		return "arm_br_24"
	case DYLD_CACHE_ADJ_V2_THUMB_MOVW_MOVT:
		return "thumb_movw_movt"
	case DYLD_CACHE_ADJ_V2_THUMB_BR22:
		return "thumb_br_22"
	case DYLD_CACHE_ADJ_V2_IMAGE_OFF_32:
		return "image_off_32"
	case DYLD_CACHE_ADJ_V2_THREADED_POINTER_64:
		return "threaded_pointer_64"
	default:
		return fmt.Sprintf("unknown kind %#02x", k)
	}

}

type ExportFlag int

const (
	/*
	 * The following are used on the flags byte of a terminal node
	 * in the export information.
	 */
	EXPORT_SYMBOL_FLAGS_KIND_MASK         ExportFlag = 0x03
	EXPORT_SYMBOL_FLAGS_KIND_REGULAR      ExportFlag = 0x00
	EXPORT_SYMBOL_FLAGS_KIND_THREAD_LOCAL ExportFlag = 0x01
	EXPORT_SYMBOL_FLAGS_KIND_ABSOLUTE     ExportFlag = 0x02
	EXPORT_SYMBOL_FLAGS_WEAK_DEFINITION   ExportFlag = 0x04
	EXPORT_SYMBOL_FLAGS_REEXPORT          ExportFlag = 0x08
	EXPORT_SYMBOL_FLAGS_STUB_AND_RESOLVER ExportFlag = 0x10
	EXPORT_SYMBOL_FLAGS_STATIC_RESOLVER   ExportFlag = 0x20
)

func (f ExportFlag) Regular() bool {
	return (f & EXPORT_SYMBOL_FLAGS_KIND_MASK) == EXPORT_SYMBOL_FLAGS_KIND_REGULAR
}
func (f ExportFlag) ThreadLocal() bool {
	return (f & EXPORT_SYMBOL_FLAGS_KIND_MASK) == EXPORT_SYMBOL_FLAGS_KIND_THREAD_LOCAL
}
func (f ExportFlag) Absolute() bool {
	return (f & EXPORT_SYMBOL_FLAGS_KIND_MASK) == EXPORT_SYMBOL_FLAGS_KIND_ABSOLUTE
}
func (f ExportFlag) WeakDefinition() bool {
	return f == EXPORT_SYMBOL_FLAGS_WEAK_DEFINITION
}
func (f ExportFlag) ReExport() bool {
	return f == EXPORT_SYMBOL_FLAGS_REEXPORT
}
func (f ExportFlag) StubAndResolver() bool {
	return f == EXPORT_SYMBOL_FLAGS_STUB_AND_RESOLVER
}
func (f ExportFlag) StaticResolver() bool {
	return f == EXPORT_SYMBOL_FLAGS_STATIC_RESOLVER
}

func (f ExportFlag) String() string {
	var fStr string
	if f.Regular() {
		fStr += "regular"
		if f.StubAndResolver() {
			fStr += "|has_resolver"
		} else if f.StaticResolver() {
			fStr += "|static_resolver"
		} else if f.WeakDefinition() {
			fStr += "|weak_def"
		}
	} else if f.ThreadLocal() {
		fStr += "per-thread"
	} else if f.Absolute() {
		fStr += "absolute"
	} else if f.ReExport() {
		fStr += "[re-export]"
	}
	return strings.TrimSpace(fStr)
}
