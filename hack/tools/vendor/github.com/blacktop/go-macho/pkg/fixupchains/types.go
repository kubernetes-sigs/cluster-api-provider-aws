package fixupchains

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/blacktop/go-macho/types"
)

type DyldChainedFixups struct {
	DyldChainedFixupsHeader
	PointerFormat DCPtrKind
	Starts        []DyldChainedStarts
	Imports       []DcfImport
	r             *bytes.Reader
	sr            types.MachoReader
	bo            binary.ByteOrder
}

type Fixup interface {
	Offset() uint64
	String(baseAddr ...uint64) string
}

type Rebase interface {
	Offset() uint64
	Target() uint64
	Raw() uint64
	String(baseAddr ...uint64) string
}

type Bind interface {
	Name() string
	Offset() uint64
	Ordinal() uint64
	Addend() uint64
	Raw() uint64
	String(baseAddr ...uint64) string
}

type DCSymbolsFormat uint32

const (
	DC_SFORMAT_UNCOMPRESSED    DCSymbolsFormat = 0
	DC_SFORMAT_ZLIB_COMPRESSED DCSymbolsFormat = 1
)

// DyldChainedFixupsHeader object is the header of the LC_DYLD_CHAINED_FIXUPS payload
type DyldChainedFixupsHeader struct {
	FixupsVersion uint32          // 0
	StartsOffset  uint32          // offset of DyldChainedStartsInImage in chain_data
	ImportsOffset uint32          // offset of imports table in chain_data
	SymbolsOffset uint32          // offset of symbol strings in chain_data
	ImportsCount  uint32          // number of imported symbol names
	ImportsFormat ImportFormat    // DYLD_CHAINED_IMPORT*
	SymbolsFormat DCSymbolsFormat // 0 => uncompressed, 1 => zlib compressed
}

// DyldChainedStartsInImage this struct is embedded in LC_DYLD_CHAINED_FIXUPS payload
type DyldChainedStartsInImage struct {
	SegCount       uint32
	SegInfoOffsets []uint32 // []uint32 ARRAY each entry is offset into this struct for that segment
	// followed by pool of dyld_chain_starts_in_segment data
}

// DCPtrKind are values for dyld_chained_starts_in_segment.pointer_format
type DCPtrKind uint16

const (
	DYLD_CHAINED_PTR_ARM64E              DCPtrKind = 1 // stride 8, unauth target is vmaddr
	DYLD_CHAINED_PTR_64                  DCPtrKind = 2 // target is vmaddr
	DYLD_CHAINED_PTR_32                  DCPtrKind = 3
	DYLD_CHAINED_PTR_32_CACHE            DCPtrKind = 4
	DYLD_CHAINED_PTR_32_FIRMWARE         DCPtrKind = 5
	DYLD_CHAINED_PTR_64_OFFSET           DCPtrKind = 6 // target is vm offset
	DYLD_CHAINED_PTR_ARM64E_OFFSET       DCPtrKind = 7 // old name
	DYLD_CHAINED_PTR_ARM64E_KERNEL       DCPtrKind = 7 // stride 4, unauth target is vm offset
	DYLD_CHAINED_PTR_64_KERNEL_CACHE     DCPtrKind = 8
	DYLD_CHAINED_PTR_ARM64E_USERLAND     DCPtrKind = 9  // stride 8, unauth target is vm offset
	DYLD_CHAINED_PTR_ARM64E_FIRMWARE     DCPtrKind = 10 // stride 4, unauth target is vmaddr
	DYLD_CHAINED_PTR_X86_64_KERNEL_CACHE DCPtrKind = 11 // stride 1, x86_64 kernel caches
	DYLD_CHAINED_PTR_ARM64E_USERLAND24   DCPtrKind = 12 // stride 8, unauth target is vm offset, 24-bit bind
)

type DyldChainedStarts struct {
	DyldChainedStartsInSegment
	PageStarts  []DCPtrStart
	ChainStarts []uint16
	Fixups      []Fixup
}

// Rebases filters fixups to only rebases
func (s *DyldChainedStarts) Rebases() []Rebase {
	var rebases []Rebase
	for _, fixup := range s.Fixups {
		if r, ok := fixup.(Rebase); ok {
			rebases = append(rebases, r)
		}
	}
	return rebases
}

// Binds filters fixups to only binds
func (s *DyldChainedStarts) Binds() []Bind {
	var binds []Bind
	for _, fixup := range s.Fixups {
		if b, ok := fixup.(Bind); ok {
			binds = append(binds, b)
		}
	}
	return binds
}

func stride(pointerFormat DCPtrKind) uint64 {
	switch pointerFormat {
	case DYLD_CHAINED_PTR_ARM64E:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_USERLAND:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_USERLAND24:
		return uint64(8)
	case DYLD_CHAINED_PTR_ARM64E_KERNEL:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_FIRMWARE:
		fallthrough
	case DYLD_CHAINED_PTR_32_FIRMWARE:
		fallthrough
	case DYLD_CHAINED_PTR_64:
		fallthrough
	case DYLD_CHAINED_PTR_64_OFFSET:
		fallthrough
	case DYLD_CHAINED_PTR_32:
		fallthrough
	case DYLD_CHAINED_PTR_32_CACHE:
		fallthrough
	case DYLD_CHAINED_PTR_64_KERNEL_CACHE:
		return uint64(4)
	case DYLD_CHAINED_PTR_X86_64_KERNEL_CACHE:
		return uint64(1)
	default:
		panic(fmt.Sprintf("unsupported pointer chain format: %d", pointerFormat))
	}
}

// DyldChainedStartsInSegment object is embedded in dyld_chain_starts_in_image
// and passed down to the kernel for page-in linking
type DyldChainedStartsInSegment struct {
	Size            uint32    // size of this (amount kernel needs to copy)
	PageSize        uint16    // 0x1000 or 0x4000
	PointerFormat   DCPtrKind // DYLD_CHAINED_PTR_*
	SegmentOffset   uint64    // offset in memory to start of segment
	MaxValidPointer uint32    // for 32-bit OS, any value beyond this is not a pointer
	PageCount       uint16    // how many pages are in array
	// uint16_t    page_start[1]      // each entry is offset in each page of first element in chain
	//                                 // or DYLD_CHAINED_PTR_START_NONE if no fixups on page
	// uint16_t    chain_starts[1];    // some 32-bit formats may require multiple starts per page.
	// for those, if high bit is set in page_starts[], then it
	// is index into chain_starts[] which is a list of starts
	// the last of which has the high bit set
}

type DCPtrStart uint16

const (
	DYLD_CHAINED_PTR_START_NONE  DCPtrStart = 0xFFFF // used in page_start[] to denote a page with no fixups
	DYLD_CHAINED_PTR_START_MULTI DCPtrStart = 0x8000 // used in page_start[] to denote a page which has multiple starts
	DYLD_CHAINED_PTR_START_LAST  DCPtrStart = 0x8000 // used in chain_starts[] to denote last start in list for page
)

func DcpArm64eIsBind(ptr uint64) bool {
	return types.ExtractBits(ptr, 62, 1) != 0
}
func DcpArm64eIsRebase(ptr uint64) bool {
	return types.ExtractBits(ptr, 62, 1) == 0
}
func DcpArm64eIsAuth(ptr uint64) bool {
	return types.ExtractBits(ptr, 63, 1) != 0
}
func DcpArm64eNext(ptr uint64) uint64 {
	return types.ExtractBits(uint64(ptr), 51, 11)
}

func Generic64Next(ptr uint64) uint64 {
	return types.ExtractBits(uint64(ptr), 51, 12)
}
func Generic64IsBind(ptr uint64) bool {
	return types.ExtractBits(uint64(ptr), 63, 1) != 0
}

func Generic32Next(ptr uint32) uint64 {
	return types.ExtractBits(uint64(ptr), 26, 5)
}
func Generic32IsBind(ptr uint32) bool {
	return types.ExtractBits(uint64(ptr), 31, 1) != 0
}

// KeyName returns the chained pointer's key name
func KeyName(keyVal uint64) string {
	name := []string{"IA", "IB", "DA", "DB"}
	if keyVal >= 4 {
		return "ERROR"
	}
	return name[keyVal]
}

// DYLD_CHAINED_PTR_ARM64E
type DyldChainedPtrArm64eRebase struct {
	Fixup   uint64
	Pointer uint64
}

func (d DyldChainedPtrArm64eRebase) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtrArm64eRebase) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtrArm64eRebase) Target() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 43) // runtimeOffset
}
func (d DyldChainedPtrArm64eRebase) High8() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 43, 8)
}
func (d DyldChainedPtrArm64eRebase) UnpackTarget() uint64 {
	return d.High8()<<56 | d.Target()
}
func (d DyldChainedPtrArm64eRebase) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 11) // 4 or 8-byte stide
}
func (d DyldChainedPtrArm64eRebase) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 62, 1) // == 0
}
func (d DyldChainedPtrArm64eRebase) Auth() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1) // == 0
}
func (d DyldChainedPtrArm64eRebase) Kind() string {
	return "rebase"
}
func (d DyldChainedPtrArm64eRebase) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, target: %#x, high8: 0x%02x)",
		d.Fixup,
		d.Pointer,
		d.Kind(),
		d.Next(),
		d.UnpackTarget(),
		d.High8(),
	)
}

// DYLD_CHAINED_PTR_ARM64E
type DyldChainedPtrArm64eBind struct {
	Fixup   uint64
	Pointer uint64
	Import  string
}

func (d DyldChainedPtrArm64eBind) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtrArm64eBind) Ordinal() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 16)
}
func (d DyldChainedPtrArm64eBind) Zero() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 16, 16)
}
func (d DyldChainedPtrArm64eBind) Addend() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 32, 19) // +/-256K
}
func (d DyldChainedPtrArm64eBind) SignExtendedAddend() int64 {
	addend := d.Addend()
	if (addend & 0x40000) != 0 {
		return int64(addend | 0xFFFFFFFFFFFC0000)
	}
	return int64(addend)
}
func (d DyldChainedPtrArm64eBind) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 11) // 4 or 8-byte stide
}
func (d DyldChainedPtrArm64eBind) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 62, 1) // == 1
}
func (d DyldChainedPtrArm64eBind) Auth() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1) // == 0
}
func (d DyldChainedPtrArm64eBind) Kind() string {
	return "bind"
}
func (d DyldChainedPtrArm64eBind) Name() string {
	return d.Import
}
func (d DyldChainedPtrArm64eBind) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtrArm64eBind) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, ordinal: %03d, addend: %d)",
		d.Fixup,
		d.Raw(),
		d.Kind(),
		d.Next(),
		d.Ordinal(),
		d.SignExtendedAddend(),
	)
}

// DYLD_CHAINED_PTR_ARM64E
type DyldChainedPtrArm64eAuthRebase struct {
	Fixup   uint64
	Pointer uint64
}

func (d DyldChainedPtrArm64eAuthRebase) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtrArm64eAuthRebase) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtrArm64eAuthRebase) Target() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 32) // target
}
func (d DyldChainedPtrArm64eAuthRebase) Diversity() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 32, 16)
}
func (d DyldChainedPtrArm64eAuthRebase) AddrDiv() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 48, 1)
}
func (d DyldChainedPtrArm64eAuthRebase) Key() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 49, 2)
}
func (d DyldChainedPtrArm64eAuthRebase) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 11) // 4 or 8-byte stide
}
func (d DyldChainedPtrArm64eAuthRebase) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 62, 1) // == 0
}
func (d DyldChainedPtrArm64eAuthRebase) Auth() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1) // == 1
}
func (d DyldChainedPtrArm64eAuthRebase) Kind() string {
	return "auth-rebase"
}
func (d DyldChainedPtrArm64eAuthRebase) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, target: %#x, key: %s, addrDiv: %d, diversity: 0x%04x)",
		d.Fixup,
		d.Pointer,
		d.Kind(),
		d.Next(),
		d.Target(),
		KeyName(d.Key()),
		d.AddrDiv(),
		d.Diversity(),
	)
}

// DYLD_CHAINED_PTR_ARM64E
type DyldChainedPtrArm64eAuthBind struct {
	Fixup   uint64
	Pointer uint64
	Import  string
}

func (d DyldChainedPtrArm64eAuthBind) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtrArm64eAuthBind) Addend() uint64 {
	return 0
}
func (d DyldChainedPtrArm64eAuthBind) Ordinal() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 16)
}
func (d DyldChainedPtrArm64eAuthBind) Zero() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 16, 16)
}
func (d DyldChainedPtrArm64eAuthBind) Diversity() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 32, 16)
}
func (d DyldChainedPtrArm64eAuthBind) AddrDiv() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 48, 1)
}
func (d DyldChainedPtrArm64eAuthBind) Key() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 49, 2)
}
func (d DyldChainedPtrArm64eAuthBind) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 11) // 4 or 8-byte stide
}
func (d DyldChainedPtrArm64eAuthBind) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 62, 1) // == 1
}
func (d DyldChainedPtrArm64eAuthBind) Auth() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1) // == 1
}
func (d DyldChainedPtrArm64eAuthBind) Kind() string {
	return "auth-bind"
}
func (d DyldChainedPtrArm64eAuthBind) Name() string {
	return d.Import
}
func (d DyldChainedPtrArm64eAuthBind) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtrArm64eAuthBind) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, ordinal: %03d, key: %s, addrDiv: %d, diversity: 0x%04x)",
		d.Fixup,
		d.Pointer,
		d.Kind(),
		d.Next(),
		d.Ordinal(),
		KeyName(d.Key()),
		d.AddrDiv(),
		d.Diversity(),
	)
}

// DYLD_CHAINED_PTR_64
type DyldChainedPtr64Rebase struct {
	Fixup   uint64
	Pointer uint64
}

func (d DyldChainedPtr64Rebase) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtr64Rebase) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtr64Rebase) Target() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 36) // runtimeOffset 64GB max image size
}
func (d DyldChainedPtr64Rebase) High8() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 36, 8) // after slide added
}
func (d DyldChainedPtr64Rebase) UnpackedTarget() uint64 {
	return d.High8()<<56 | d.Target()
}
func (d DyldChainedPtr64Rebase) Reserved() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 44, 7) // all zeros
}
func (d DyldChainedPtr64Rebase) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 12) // 4-byte stride
}
func (d DyldChainedPtr64Rebase) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1) // == 0
}
func (d DyldChainedPtr64Rebase) Kind() string {
	return "ptr64-rebase"
}
func (d DyldChainedPtr64Rebase) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, target: %#x, high8: 0x%02x)",
		d.Fixup,
		d.Pointer,
		d.Kind(),
		d.Next(),
		d.Target(),
		d.High8(),
	)
}

// DYLD_CHAINED_PTR_64_OFFSET
type DyldChainedPtr64RebaseOffset struct {
	Fixup   uint64
	Pointer uint64
}

func (d DyldChainedPtr64RebaseOffset) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtr64RebaseOffset) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtr64RebaseOffset) Target() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 36) // vmAddr 64GB max image size
}
func (d DyldChainedPtr64RebaseOffset) High8() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 36, 8) // before slide added)
}
func (d DyldChainedPtr64RebaseOffset) UnpackedTarget() uint64 {
	return d.High8()<<56 | d.Target()
}
func (d DyldChainedPtr64RebaseOffset) Reserved() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 44, 7) // all zeros
}
func (d DyldChainedPtr64RebaseOffset) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 12) // 4-byte stride
}
func (d DyldChainedPtr64RebaseOffset) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1) // == 0
}
func (d DyldChainedPtr64RebaseOffset) Kind() string {
	return "rebase-offset"
}
func (d DyldChainedPtr64RebaseOffset) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, target: %#x, high8: 0x%02x)",
		d.Fixup,
		d.Pointer,
		d.Kind(),
		d.Next(),
		d.Target(),
		d.High8(),
	)
}

// DYLD_CHAINED_PTR_ARM64E_USERLAND24
type DyldChainedPtrArm64eRebase24 struct {
	Fixup   uint64
	Pointer uint64
}

func (d DyldChainedPtrArm64eRebase24) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtrArm64eRebase24) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtrArm64eRebase24) Target() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 24) // runtimeOffset
}
func (d DyldChainedPtrArm64eRebase24) High8() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 24, 8) // TODO: check that this is correct when src is released
}
func (d DyldChainedPtrArm64eRebase24) UnpackTarget() uint64 {
	return d.High8()<<56 | d.Target()
}
func (d DyldChainedPtrArm64eRebase24) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 11) // 8-byte stide
}
func (d DyldChainedPtrArm64eRebase24) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 62, 1) // == 0
}
func (d DyldChainedPtrArm64eRebase24) Auth() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1) // == 0
}
func (d DyldChainedPtrArm64eRebase24) Kind() string {
	return "rebase24"
}
func (d DyldChainedPtrArm64eRebase24) String(baseAddr ...uint64) string {
	var baddr uint64
	if len(baseAddr) > 0 {
		baddr = baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, target: %#x, high8: 0x%02x)",
		d.Fixup+baddr,
		d.Pointer,
		d.Kind(),
		d.Next(),
		d.UnpackTarget()+baddr, // TODO: check that this is correct when src is released
		d.High8(),
	)
}

// DYLD_CHAINED_PTR_ARM64E_USERLAND24
type DyldChainedPtrArm64eAuthRebase24 struct {
	Fixup   uint64
	Pointer uint64
}

func (d DyldChainedPtrArm64eAuthRebase24) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtrArm64eAuthRebase24) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtrArm64eAuthRebase24) Target() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 24) // target
}
func (d DyldChainedPtrArm64eAuthRebase24) Diversity() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 24, 16)
}
func (d DyldChainedPtrArm64eAuthRebase24) AddrDiv() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 40, 1)
}
func (d DyldChainedPtrArm64eAuthRebase24) Key() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 41, 2)
}
func (d DyldChainedPtrArm64eAuthRebase24) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 11) // 8-byte stide
}
func (d DyldChainedPtrArm64eAuthRebase24) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 62, 1) // == 0
}
func (d DyldChainedPtrArm64eAuthRebase24) Auth() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1) // == 1
}
func (d DyldChainedPtrArm64eAuthRebase24) Kind() string {
	return "auth-rebase24"
}
func (d DyldChainedPtrArm64eAuthRebase24) String(baseAddr ...uint64) string {
	var baddr uint64
	if len(baseAddr) > 0 {
		baddr = baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, target: %#x, key: %s, addrDiv: %d, diversity: 0x%04x)",
		d.Fixup+baddr,
		d.Pointer,
		d.Kind(),
		d.Next(),
		d.Target()+baddr,
		KeyName(d.Key()),
		d.AddrDiv(),
		d.Diversity(),
	)
}

// DYLD_CHAINED_PTR_ARM64E_USERLAND24
type DyldChainedPtrArm64eBind24 struct {
	Fixup   uint64
	Pointer uint64
	Import  string
}

func (d DyldChainedPtrArm64eBind24) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtrArm64eBind24) Ordinal() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 24)
}
func (d DyldChainedPtrArm64eBind24) Zero() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 24, 8)
}
func (d DyldChainedPtrArm64eBind24) Addend() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 32, 19)
}
func (d DyldChainedPtrArm64eBind24) SignExtendedAddend() int64 {
	addend := d.Addend()
	if (addend & 0x40000) != 0 {
		return int64(addend | 0xFFFFFFFFFFFC0000)
	}
	return int64(addend)
}
func (d DyldChainedPtrArm64eBind24) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 11)
}
func (d DyldChainedPtrArm64eBind24) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 62, 1)
}
func (d DyldChainedPtrArm64eBind24) Auth() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1)
}
func (d DyldChainedPtrArm64eBind24) Kind() string {
	return "bind24"
}
func (d DyldChainedPtrArm64eBind24) Name() string {
	return d.Import
}
func (d DyldChainedPtrArm64eBind24) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtrArm64eBind24) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, ordinal: %03d, addend: %d)", d.Fixup, d.Pointer, d.Kind(), d.Next(), d.Ordinal(), d.SignExtendedAddend())
}

// DYLD_CHAINED_PTR_ARM64E_USERLAND24
type DyldChainedPtrArm64eAuthBind24 struct {
	Fixup   uint64
	Pointer uint64
	Import  string
}

func (d DyldChainedPtrArm64eAuthBind24) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtrArm64eAuthBind24) Addend() uint64 {
	return 0
}
func (d DyldChainedPtrArm64eAuthBind24) Ordinal() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 24)
}
func (d DyldChainedPtrArm64eAuthBind24) Zero() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 24, 8)
}
func (d DyldChainedPtrArm64eAuthBind24) Diversity() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 32, 16)
}
func (d DyldChainedPtrArm64eAuthBind24) AddrDiv() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 48, 1)
}
func (d DyldChainedPtrArm64eAuthBind24) Key() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 49, 2)
}
func (d DyldChainedPtrArm64eAuthBind24) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 11)
}
func (d DyldChainedPtrArm64eAuthBind24) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 62, 1)
}
func (d DyldChainedPtrArm64eAuthBind24) Auth() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1)
}
func (d DyldChainedPtrArm64eAuthBind24) Kind() string {
	return "auth-bind24"
}
func (d DyldChainedPtrArm64eAuthBind24) Name() string {
	return d.Import
}
func (d DyldChainedPtrArm64eAuthBind24) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtrArm64eAuthBind24) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, key: %s, addrDiv: %d, diversity: 0x%04x, ordinal: %03d)",
		d.Fixup,
		d.Pointer,
		d.Kind(),
		d.Next(),
		KeyName(d.Key()),
		d.AddrDiv(),
		d.Diversity(),
		d.Ordinal(),
	)
}

// DYLD_CHAINED_PTR_64
type DyldChainedPtr64Bind struct {
	Fixup   uint64
	Pointer uint64
	Import  string
}

func (d DyldChainedPtr64Bind) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtr64Bind) Ordinal() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 24)
}
func (d DyldChainedPtr64Bind) Addend() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 24, 8) // 0 thru 255
}
func (d DyldChainedPtr64Bind) Reserved() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 32, 19) // all zeros
}
func (d DyldChainedPtr64Bind) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 12) // 4-byte stride
}
func (d DyldChainedPtr64Bind) Bind() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1) // == 1
}
func (d DyldChainedPtr64Bind) Kind() string {
	return "ptr64-bind"
}
func (d DyldChainedPtr64Bind) Name() string {
	return d.Import
}
func (d DyldChainedPtr64Bind) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtr64Bind) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, ordinal: %06x, addend: %d)",
		d.Fixup,
		d.Pointer,
		d.Kind(),
		d.Next(),
		d.Ordinal(),
		d.Addend(),
	)
}

// DYLD_CHAINED_PTR_64_KERNEL_CACHE, DYLD_CHAINED_PTR_X86_64_KERNEL_CACHE
type DyldChainedPtr64KernelCacheRebase struct {
	Fixup   uint64
	Pointer uint64
}

func (d DyldChainedPtr64KernelCacheRebase) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtr64KernelCacheRebase) Raw() uint64 {
	return d.Pointer
}
func (d DyldChainedPtr64KernelCacheRebase) Target() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 30) // basePointers[cacheLevel] + target
}
func (d DyldChainedPtr64KernelCacheRebase) CacheLevel() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 30, 2) // what level of cache to bind to (indexes a mach_header array)
}
func (d DyldChainedPtr64KernelCacheRebase) Diversity() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 32, 16)
}
func (d DyldChainedPtr64KernelCacheRebase) AddrDiv() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 48, 1)
}
func (d DyldChainedPtr64KernelCacheRebase) Key() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 49, 2)
}
func (d DyldChainedPtr64KernelCacheRebase) Next() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 51, 12) // 1 or 4-byte stide
}
func (d DyldChainedPtr64KernelCacheRebase) IsAuth() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 63, 1)
}
func (d DyldChainedPtr64KernelCacheRebase) Kind() string {
	return "kcache-rebase"
}
func (d DyldChainedPtr64KernelCacheRebase) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	if d.IsAuth() == 1 {
		return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, key: %s, addrDiv: %d, diversity: 0x%04x, target: 0x%08x, cacheLevel: %d)",
			d.Fixup,
			d.Pointer,
			d.Kind(),
			d.Next(),
			KeyName(d.Key()),
			d.AddrDiv(),
			d.Diversity(),
			d.Target(),
			d.CacheLevel(),
		)
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%016x %16s: (next: %03d, target: 0x%08x, cacheLevel: %d)",
		d.Fixup,
		d.Pointer,
		d.Kind(),
		d.Next(),
		d.Target(),
		d.CacheLevel(),
	)
}

// DYLD_CHAINED_PTR_32
// Note: for DYLD_CHAINED_PTR_32 some non-pointer values are co-opted into the chain
// as out of range rebases.  If an entry in the chain is > max_valid_pointer, then it
// is not a pointer.  To restore the value, subtract off the bias, which is
// (64MB+max_valid_pointer)/2.
type DyldChainedPtr32Rebase struct {
	Fixup   uint64
	Pointer uint32
}

func (d DyldChainedPtr32Rebase) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtr32Rebase) Target() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 26) // vmaddr, 64MB max image size
}
func (d DyldChainedPtr32Rebase) Next() uint32 {
	return uint32(types.ExtractBits(uint64(d.Pointer), 26, 5)) // 4-byte stride
}
func (d DyldChainedPtr32Rebase) Bind() uint32 {
	return uint32(types.ExtractBits(uint64(d.Pointer), 31, 1)) // == 0
}
func (d DyldChainedPtr32Rebase) Kind() string {
	return "ptr32-rebase"
}
func (d DyldChainedPtr32Rebase) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%08x %16s: (next:%02d target: 0x%07x)", d.Fixup, d.Pointer, d.Kind(), d.Next(), d.Target())
}

// DYLD_CHAINED_PTR_32
type DyldChainedPtr32Bind struct {
	Fixup   uint64
	Pointer uint32
	Import  string
}

func (d DyldChainedPtr32Bind) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtr32Bind) Ordinal() uint64 {
	return uint64(types.ExtractBits(uint64(d.Pointer), 0, 20))
}
func (d DyldChainedPtr32Bind) Addend() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 20, 6) // 0 thru 63
}
func (d DyldChainedPtr32Bind) Next() uint32 {
	return uint32(types.ExtractBits(uint64(d.Pointer), 26, 5)) // 4-byte stride
}
func (d DyldChainedPtr32Bind) Bind() uint32 {
	return uint32(types.ExtractBits(uint64(d.Pointer), 31, 1)) // == 1
}
func (d DyldChainedPtr32Bind) Kind() string {
	return "ptr32-bind"
}
func (d DyldChainedPtr32Bind) Name() string {
	return d.Import
}
func (d DyldChainedPtr32Bind) Raw() uint64 {
	return uint64(d.Pointer)
}
func (d DyldChainedPtr32Bind) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%08x %16s: (next:%02d ordinal:%05x addend:%d)", d.Fixup, d.Pointer, d.Kind(), d.Next(), d.Ordinal(), d.Addend())
}

// DYLD_CHAINED_PTR_32_CACHE
type DyldChainedPtr32CacheRebase struct {
	Fixup   uint64
	Pointer uint32
}

func (d DyldChainedPtr32CacheRebase) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtr32CacheRebase) Raw() uint64 {
	return uint64(d.Pointer)
}
func (d DyldChainedPtr32CacheRebase) Target() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 30) // 1GB max dyld cache TEXT and DATA
}
func (d DyldChainedPtr32CacheRebase) Next() uint32 {
	return uint32(types.ExtractBits(uint64(d.Pointer), 30, 2)) // 4-byte stride
}
func (d DyldChainedPtr32CacheRebase) Kind() string {
	return "cache-rebase"
}
func (d DyldChainedPtr32CacheRebase) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%08x %16s: (next:%02d target: 0x%07x)", d.Fixup, d.Pointer, d.Kind(), d.Next(), d.Target())
}

// DYLD_CHAINED_PTR_32_FIRMWARE
type DyldChainedPtr32FirmwareRebase struct {
	Fixup   uint64
	Pointer uint32
}

func (d DyldChainedPtr32FirmwareRebase) Offset() uint64 {
	return d.Fixup
}
func (d DyldChainedPtr32FirmwareRebase) Raw() uint64 {
	return uint64(d.Pointer)
}
func (d DyldChainedPtr32FirmwareRebase) Target() uint64 {
	return types.ExtractBits(uint64(d.Pointer), 0, 26) // 64MB max firmware TEXT and DATA
}
func (d DyldChainedPtr32FirmwareRebase) Next() uint32 {
	return uint32(types.ExtractBits(uint64(d.Pointer), 26, 6)) // 4-byte stride
}
func (d DyldChainedPtr32FirmwareRebase) Kind() string {
	return "firmware-rebase"
}
func (d DyldChainedPtr32FirmwareRebase) String(baseAddr ...uint64) string {
	if len(baseAddr) > 0 {
		d.Fixup += baseAddr[0]
	}
	return fmt.Sprintf("0x%08x:  raw: 0x%08x %16s: (next:%02d target: 0x%07x)", d.Fixup, d.Pointer, d.Kind(), d.Next(), d.Target())
}
