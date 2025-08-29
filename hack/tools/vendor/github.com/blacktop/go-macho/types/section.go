package types

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

/*
 * A segment is made up of zero or more sections.  Non-MH_OBJECT files have
 * all of their segments with the proper sections in each, and padded to the
 * specified segment alignment when produced by the link editor.  The first
 * segment of a MH_EXECUTE and MH_FVMLIB format file contains the mach_header
 * and load commands of the object file before its first section.  The zero
 * fill sections are always last in their segment (in all formats).  This
 * allows the zeroed segment padding to be mapped into memory where zero fill
 * sections might be. The gigabyte zero fill sections, those with the section
 * type S_GB_ZEROFILL, can only be in a segment with sections of this type.
 * These segments are then placed after all other segments.
 *
 * The MH_OBJECT format has all of its sections in one segment for
 * compactness.  There is no padding to a specified segment boundary and the
 * mach_header and load commands are not part of the segment.
 *
 * Sections with the same section name, sectname, going into the same segment,
 * segname, are combined by the link editor.  The resulting section is aligned
 * to the maximum alignment of the combined sections and is the new section's
 * alignment.  The combined sections are aligned to their original alignment in
 * the combined section.  Any padded bytes to get the specified alignment are
 * zeroed.
 *
 * The format of the relocation entries referenced by the reloff and nreloc
 * fields of the section structure for mach object files is described in the
 * header file <reloc.h>.
 */

// A Section32 is a 32-bit Mach-O section header.
type Section32 struct {
	Name     [16]byte
	Seg      [16]byte
	Addr     uint32
	Size     uint32
	Offset   uint32
	Align    uint32
	Reloff   uint32
	Nreloc   uint32
	Flags    SectionFlag
	Reserve1 uint32
	Reserve2 uint32
}

// A Section64 is a 64-bit Mach-O section header.
type Section64 struct {
	Name     [16]byte
	Seg      [16]byte
	Addr     uint64
	Size     uint64
	Offset   uint32
	Align    uint32
	Reloff   uint32
	Nreloc   uint32
	Flags    SectionFlag
	Reserve1 uint32
	Reserve2 uint32
	Reserve3 uint32
}

type SectionFlag uint32

const (
	SectionType       SectionFlag = 0x000000ff /* 256 section types */
	SectionAttributes SectionFlag = 0xffffff00 /*  24 section attributes */
)

/*
 * The flags field of a section structure is separated into two parts a section
 * type and section attributes.  The section types are mutually exclusive (it
 * can only have one type) but the section attributes are not (it may have more
 * than one attribute).
 */
const (
	/* Constants for the type of a section */
	Regular         SectionFlag = 0x0 /* regular section */
	Zerofill        SectionFlag = 0x1 /* zero fill on demand section */
	CstringLiterals SectionFlag = 0x2 /* section with only literal C strings*/
	ByteLiterals4   SectionFlag = 0x3 /* section with only 4 byte literals */
	ByteLiterals8   SectionFlag = 0x4 /* section with only 8 byte literals */
	LiteralPointers SectionFlag = 0x5 /* section with only pointers to literals */
	/*
	 * For the two types of symbol pointers sections and the symbol stubs section
	 * they have indirect symbol table entries.  For each of the entries in the
	 * section the indirect symbol table entries, in corresponding order in the
	 * indirect symbol table, start at the index stored in the reserved1 field
	 * of the section structure.  Since the indirect symbol table entries
	 * correspond to the entries in the section the number of indirect symbol table
	 * entries is inferred from the size of the section divided by the size of the
	 * entries in the section.  For symbol pointers sections the size of the entries
	 * in the section is 4 bytes and for symbol stubs sections the byte size of the
	 * stubs is stored in the reserved2 field of the section structure.
	 */
	NonLazySymbolPointers   SectionFlag = 0x6  /* section with only non-lazy symbol pointers */
	LazySymbolPointers      SectionFlag = 0x7  /* section with only lazy symbol pointers */
	SymbolStubs             SectionFlag = 0x8  /* section with only symbol stubs, byte size of stub in the reserved2 field */
	ModInitFuncPointers     SectionFlag = 0x9  /* section with only function pointers for initialization*/
	ModTermFuncPointers     SectionFlag = 0xa  /* section with only function pointers for termination */
	Coalesced               SectionFlag = 0xb  /* section contains symbols that are to be coalesced */
	GbZerofill              SectionFlag = 0xc  /* zero fill on demand section (that can be larger than 4 gigabytes) */
	Interposing             SectionFlag = 0xd  /* section with only pairs of function pointers for interposing */
	ByteLiterals16          SectionFlag = 0xe  /* section with only 16 byte literals */
	DtraceDof               SectionFlag = 0xf  /* section contains DTrace Object Format */
	LazyDylibSymbolPointers SectionFlag = 0x10 /* section with only lazy symbol pointers to lazy loaded dylibs */
	/*
	 * Section types to support thread local variables
	 */
	ThreadLocalRegular              SectionFlag = 0x11 /* template of initial values for TLVs */
	ThreadLocalZerofill             SectionFlag = 0x12 /* template of initial values for TLVs */
	ThreadLocalVariables            SectionFlag = 0x13 /* TLV descriptors */
	ThreadLocalVariablePointers     SectionFlag = 0x14 /* pointers to TLV descriptors */
	ThreadLocalInitFunctionPointers SectionFlag = 0x15 /* functions to call to initialize TLV values */
	InitFuncOffsets                 SectionFlag = 0x16 /* 32-bit offsets to initializers */
)

func (t SectionFlag) IsRegular() bool {
	return (t & SectionType) == Regular
}

func (t SectionFlag) IsZerofill() bool {
	return (t & SectionType) == Zerofill
}

func (t SectionFlag) IsCstringLiterals() bool {
	return (t & SectionType) == CstringLiterals
}

func (t SectionFlag) Is4ByteLiterals() bool {
	return (t & SectionType) == ByteLiterals4
}

func (t SectionFlag) Is8ByteLiterals() bool {
	return (t & SectionType) == ByteLiterals8
}

func (t SectionFlag) IsLiteralPointers() bool {
	return (t & SectionType) == LiteralPointers
}

func (t SectionFlag) IsNonLazySymbolPointers() bool {
	return (t & SectionType) == NonLazySymbolPointers
}

func (t SectionFlag) IsLazySymbolPointers() bool {
	return (t & SectionType) == LazySymbolPointers
}

func (t SectionFlag) IsSymbolStubs() bool {
	return (t & SectionType) == SymbolStubs
}

func (t SectionFlag) IsModInitFuncPointers() bool {
	return (t & SectionType) == ModInitFuncPointers
}

func (t SectionFlag) IsModTermFuncPointers() bool {
	return (t & SectionType) == ModTermFuncPointers
}

func (t SectionFlag) IsCoalesced() bool {
	return (t & SectionType) == Coalesced
}

func (t SectionFlag) IsGbZerofill() bool {
	return (t & SectionType) == GbZerofill
}

func (t SectionFlag) IsInterposing() bool {
	return (t & SectionType) == Interposing
}

func (t SectionFlag) Is16ByteLiterals() bool {
	return (t & SectionType) == ByteLiterals16
}

func (t SectionFlag) IsDtraceDof() bool {
	return (t & SectionType) == DtraceDof
}

func (t SectionFlag) IsLazyDylibSymbolPointers() bool {
	return (t & SectionType) == LazyDylibSymbolPointers
}

func (t SectionFlag) IsThreadLocalRegular() bool {
	return (t & SectionType) == ThreadLocalRegular
}

func (t SectionFlag) IsThreadLocalZerofill() bool {
	return (t & SectionType) == ThreadLocalZerofill
}

func (t SectionFlag) IsThreadLocalVariables() bool {
	return (t & SectionType) == ThreadLocalVariables
}

func (t SectionFlag) IsThreadLocalVariablePointers() bool {
	return (t & SectionType) == ThreadLocalVariablePointers
}

func (t SectionFlag) IsThreadLocalInitFunctionPointers() bool {
	return (t & SectionType) == ThreadLocalInitFunctionPointers
}

func (t SectionFlag) IsInitFuncOffsets() bool {
	return (t & SectionType) == InitFuncOffsets
}

func (f SectionFlag) List() []string {
	var flags []string
	// if f.IsRegular() {
	// 	flags = append(flags, "Regular")
	// }
	if f.IsZerofill() {
		flags = append(flags, "Zerofill")
	}
	if f.IsCstringLiterals() {
		flags = append(flags, "CstringLiterals")
	}
	if f.Is4ByteLiterals() {
		flags = append(flags, "4ByteLiterals")
	}
	if f.Is8ByteLiterals() {
		flags = append(flags, "8ByteLiterals")
	}
	if f.IsLiteralPointers() {
		flags = append(flags, "LiteralPointers")
	}
	if f.IsNonLazySymbolPointers() {
		flags = append(flags, "NonLazySymbolPointers")
	}
	if f.IsLazySymbolPointers() {
		flags = append(flags, "LazySymbolPointers")
	}
	if f.IsSymbolStubs() {
		flags = append(flags, "SymbolStubs")
	}
	if f.IsModInitFuncPointers() {
		flags = append(flags, "ModInitFuncPointers")
	}
	if f.IsModTermFuncPointers() {
		flags = append(flags, "ModTermFuncPointers")
	}
	if f.IsCoalesced() {
		flags = append(flags, "Coalesced")
	}
	if f.IsGbZerofill() {
		flags = append(flags, "GbZerofill")
	}
	if f.IsInterposing() {
		flags = append(flags, "Interposing")
	}
	if f.Is16ByteLiterals() {
		flags = append(flags, "16ByteLiterals")
	}
	if f.IsDtraceDof() {
		flags = append(flags, "DtraceDOF")
	}
	if f.IsLazyDylibSymbolPointers() {
		flags = append(flags, "LazyDylibSymbolPointers")
	}
	if f.IsThreadLocalRegular() {
		flags = append(flags, "ThreadLocalRegular")
	}
	if f.IsThreadLocalZerofill() {
		flags = append(flags, "ThreadLocalZerofill")
	}
	if f.IsThreadLocalVariables() {
		flags = append(flags, "ThreadLocalVariables")
	}
	if f.IsThreadLocalVariablePointers() {
		flags = append(flags, "ThreadLocalVariablePointers")
	}
	if f.IsThreadLocalInitFunctionPointers() {
		flags = append(flags, "ThreadLocalInitFunctionPointers")
	}
	if f.IsInitFuncOffsets() {
		flags = append(flags, "InitFuncOffsets")
	}
	return flags
}

func (f SectionFlag) String() string {
	return strings.Join(f.List(), ", ")
}

const (
	/*
	 * Constants for the section attributes part of the flags field of a section
	 * structure.
	 */
	SECTION_ATTRIBUTES_USR SectionFlag = 0xff000000 /* User setable attributes */
	SECTION_ATTRIBUTES_SYS SectionFlag = 0x00ffff00 /* system setable attributes */

	PURE_INSTRUCTIONS   SectionFlag = 0x80000000 /* section contains only true machine instructions */
	NO_TOC              SectionFlag = 0x40000000 /* section contains coalesced symbols that are not to be in a ranlib table of contents */
	STRIP_STATIC_SYMS   SectionFlag = 0x20000000 /* ok to strip static symbols in this section in files with the MH_DYLDLINK flag */
	NoDeadStrip         SectionFlag = 0x10000000 /* no dead stripping */
	LIVE_SUPPORT        SectionFlag = 0x08000000 /* blocks are live if they reference live blocks */
	SELF_MODIFYING_CODE SectionFlag = 0x04000000 /* Used with i386 code stubs written on by dyld */
	/*
	 * If a segment contains any sections marked with DEBUG then all
	 * sections in that segment must have this attribute.  No section other than
	 * a section marked with this attribute may reference the contents of this
	 * section.  A section with this attribute may contain no symbols and must have
	 * a section type S_REGULAR.  The static linker will not copy section contents
	 * from sections with this attribute into its output file.  These sections
	 * generally contain DWARF debugging info.
	 */
	DEBUG             SectionFlag = 0x02000000 /* a debug section */
	SOME_INSTRUCTIONS SectionFlag = 0x00000400 /* section contains some machine instructions */
	EXT_RELOC         SectionFlag = 0x00000200 /* section has external relocation entries */
	LOC_RELOC         SectionFlag = 0x00000100 /* section has local relocation entries */
)

func (t SectionFlag) IsPureInstructions() bool {
	return ((t & SectionAttributes) & PURE_INSTRUCTIONS) != 0
}
func (t SectionFlag) IsNoToc() bool {
	return ((t & SectionAttributes) & NO_TOC) != 0
}
func (t SectionFlag) IsStripStaticSyms() bool {
	return ((t & SectionAttributes) & STRIP_STATIC_SYMS) != 0
}
func (t SectionFlag) IsNoDeadStrip() bool {
	return ((t & SectionAttributes) & NoDeadStrip) != 0
}
func (t SectionFlag) IsLiveSupport() bool {
	return ((t & SectionAttributes) & LIVE_SUPPORT) != 0
}
func (t SectionFlag) IsSelfModifyingCode() bool {
	return ((t & SectionAttributes) & SELF_MODIFYING_CODE) != 0
}
func (t SectionFlag) IsDebug() bool {
	return ((t & SectionAttributes) & DEBUG) != 0
}
func (t SectionFlag) IsSomeInstructions() bool {
	return ((t & SectionAttributes) & SOME_INSTRUCTIONS) != 0
}
func (t SectionFlag) IsExtReloc() bool {
	return ((t & SectionAttributes) & EXT_RELOC) != 0
}
func (t SectionFlag) IsLocReloc() bool {
	return ((t & SectionAttributes) & LOC_RELOC) != 0
}

func (f SectionFlag) AttributesList() []string {
	var attrs []string
	if f.IsPureInstructions() {
		attrs = append(attrs, "PureInstructions")
	}
	if f.IsNoToc() {
		attrs = append(attrs, "NoToc")
	}
	if f.IsStripStaticSyms() {
		attrs = append(attrs, "StripStaticSyms")
	}
	if f.IsNoDeadStrip() {
		attrs = append(attrs, "NoDeadStrip")
	}
	if f.IsLiveSupport() {
		attrs = append(attrs, "LiveSupport")
	}
	if f.IsSelfModifyingCode() {
		attrs = append(attrs, "SelfModifyingCode")
	}
	if f.IsDebug() {
		attrs = append(attrs, "Debug")
	}
	if f.IsSomeInstructions() {
		attrs = append(attrs, "SomeInstructions")
	}
	if f.IsExtReloc() {
		attrs = append(attrs, "ExtReloc")
	}
	if f.IsLocReloc() {
		attrs = append(attrs, "LocReloc")
	}
	return attrs
}

func (f SectionFlag) Attributes() string {
	return strings.Join(f.AttributesList(), "|")
}

/*
 * Sections of type S_THREAD_LOCAL_VARIABLES contain an array
 * of tlv_descriptor structures.
 */
type TlvDescriptor struct { // TODO: implement this
	Thunk  uint64
	Key    uint64
	Offset uint64
}

/*******************************************************************************
 * SECTION
 *******************************************************************************/

type SectionHeader struct {
	Name      string
	Seg       string
	Addr      uint64
	Size      uint64
	Offset    uint32
	Align     uint32
	Reloff    uint32
	Nreloc    uint32
	Flags     SectionFlag
	Reserved1 uint32
	Reserved2 uint32
	Reserved3 uint32 // only present if original was 64-bit
	Type      uint8
}

// A Reloc represents a Mach-O relocation.
type Reloc struct {
	Addr  uint32
	Value uint32
	// when Scattered == false && Extern == true, Value is the symbol number.
	// when Scattered == false && Extern == false, Value is the section number.
	// when Scattered == true, Value is the value that this reloc refers to.
	Type      uint8
	Len       uint8 // 0=byte, 1=word, 2=long, 3=quad
	Pcrel     bool
	Extern    bool // valid if Scattered == false
	Scattered bool
}

func (r *Reloc) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Addr      uint32 `json:"addr"`
		Value     uint32 `json:"value"`
		Type      uint8  `json:"type"`
		Len       uint8  `json:"len"`
		Pcrel     bool   `json:"pcrel"`
		Extern    bool   `json:"extern"`
		Scattered bool   `json:"scattered"`
	}{
		Addr:      r.Addr,
		Value:     r.Value,
		Type:      r.Type,
		Len:       r.Len,
		Pcrel:     r.Pcrel,
		Extern:    r.Extern,
		Scattered: r.Scattered,
	})
}

type RelocInfo struct {
	Addr   uint32
	Symnum uint32
}

type Section struct {
	SectionHeader
	Relocs []Reloc

	// Embed ReaderAt for ReadAt method.
	// Do not embed SectionReader directly
	// to avoid having Read and Seek.
	// If a client wants Read and Seek it must use
	// Open() to avoid fighting over the seek offset
	// with other clients.
	io.ReaderAt
	sr *io.SectionReader
}

func (s *Section) SetReaders(r io.ReaderAt, sr *io.SectionReader) {
	s.ReaderAt = r
	s.sr = sr
}

// Data reads and returns the contents of the Mach-O section.
func (s *Section) Data() ([]byte, error) {
	dat := make([]byte, s.Size)
	n, err := s.ReadAt(dat, int64(s.Offset))
	if n == len(dat) {
		err = nil
	}
	return dat[0:n], err
}

// Open returns a new ReadSeeker reading the Mach-O section.
func (s *Section) Open() io.ReadSeeker { return io.NewSectionReader(s.sr, 0, 1<<63-1) }

func (s *Section) Put32(b []byte, o binary.ByteOrder) int {
	PutAtMost16Bytes(b[0:], s.Name)
	PutAtMost16Bytes(b[16:], s.Seg)
	o.PutUint32(b[8*4:], uint32(s.Addr))
	o.PutUint32(b[9*4:], uint32(s.Size))
	o.PutUint32(b[10*4:], s.Offset)
	o.PutUint32(b[11*4:], s.Align)
	o.PutUint32(b[12*4:], s.Reloff)
	o.PutUint32(b[13*4:], s.Nreloc)
	o.PutUint32(b[14*4:], uint32(s.Flags))
	o.PutUint32(b[15*4:], s.Reserved1)
	o.PutUint32(b[16*4:], s.Reserved2)
	a := 17 * 4
	return a + s.PutRelocs(b[a:], o)
}

func (s *Section) Put64(b []byte, o binary.ByteOrder) int {
	PutAtMost16Bytes(b[0:], s.Name)
	PutAtMost16Bytes(b[16:], s.Seg)
	o.PutUint64(b[8*4+0*8:], s.Addr)
	o.PutUint64(b[8*4+1*8:], s.Size)
	o.PutUint32(b[8*4+2*8:], s.Offset)
	o.PutUint32(b[9*4+2*8:], s.Align)
	o.PutUint32(b[10*4+2*8:], s.Reloff)
	o.PutUint32(b[11*4+2*8:], s.Nreloc)
	o.PutUint32(b[12*4+2*8:], uint32(s.Flags))
	o.PutUint32(b[13*4+2*8:], s.Reserved1)
	o.PutUint32(b[14*4+2*8:], s.Reserved2)
	o.PutUint32(b[15*4+2*8:], s.Reserved3)
	a := 16*4 + 2*8
	return a + s.PutRelocs(b[a:], o)
}

func (s *Section) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	var name [16]byte
	var seg [16]byte
	copy(name[:], s.Name)
	copy(seg[:], s.Seg)

	if s.Type == 32 {
		if err := binary.Write(buf, o, Section32{
			Name:     name,           // [16]byte
			Seg:      seg,            // [16]byte
			Addr:     uint32(s.Addr), // uint32
			Size:     uint32(s.Size), // uint32
			Offset:   s.Offset,       // uint32
			Align:    s.Align,        // uint32
			Reloff:   s.Reloff,       // uint32
			Nreloc:   s.Nreloc,       // uint32
			Flags:    s.Flags,        // SectionFlag
			Reserve1: s.Reserved1,    // uint32
			Reserve2: s.Reserved2,    // uint32
		}); err != nil {
			return fmt.Errorf("failed to write 32bit Section %s data to buffer: %v", s.Name, err)
		}
	} else { // 64
		if err := binary.Write(buf, o, Section64{
			Name:     name,        // [16]byte
			Seg:      seg,         // [16]byte
			Addr:     s.Addr,      // uint64
			Size:     s.Size,      // uint64
			Offset:   s.Offset,    // uint32
			Align:    s.Align,     // uint32
			Reloff:   s.Reloff,    // uint32
			Nreloc:   s.Nreloc,    // uint32
			Flags:    s.Flags,     // SectionFlag
			Reserve1: s.Reserved1, // uint32
			Reserve2: s.Reserved2, // uint32
			Reserve3: s.Reserved3, // uint32
		}); err != nil {
			return fmt.Errorf("failed to write 64bit Section %s data to buffer: %v", s.Name, err)
		}
	}

	return nil
}

func (s *Section) PutRelocs(b []byte, o binary.ByteOrder) int {
	a := 0
	for _, r := range s.Relocs {
		var ri RelocInfo
		typ := uint32(r.Type) & (1<<4 - 1)
		len := uint32(r.Len) & (1<<2 - 1)
		pcrel := uint32(0)
		if r.Pcrel {
			pcrel = 1
		}
		ext := uint32(0)
		if r.Extern {
			ext = 1
		}
		switch {
		case r.Scattered:
			ri.Addr = r.Addr&(1<<24-1) | typ<<24 | len<<28 | 1<<31 | pcrel<<30
			ri.Symnum = r.Value
		case o == binary.LittleEndian:
			ri.Addr = r.Addr
			ri.Symnum = r.Value&(1<<24-1) | pcrel<<24 | len<<25 | ext<<27 | typ<<28
		case o == binary.BigEndian:
			ri.Addr = r.Addr
			ri.Symnum = r.Value<<8 | pcrel<<7 | len<<5 | ext<<4 | typ
		}
		o.PutUint32(b, ri.Addr)
		o.PutUint32(b[4:], ri.Symnum)
		a += 8
		b = b[8:]
	}
	return a
}

func (s *Section) UncompressedSize() uint64 {
	if !strings.HasPrefix(s.Name, "__z") {
		return s.Size
	}
	b := make([]byte, 12)
	n, err := s.sr.ReadAt(b, 0)
	if err != nil {
		panic("Malformed object file")
	}
	if n != len(b) {
		return s.Size
	}
	if string(b[:4]) == "ZLIB" {
		return binary.BigEndian.Uint64(b[4:12])
	}
	return s.Size
}

func (s *Section) PutData(b []byte) {
	bb := b[0:s.Size]
	n, err := s.sr.ReadAt(bb, 0)
	if err != nil || uint64(n) != s.Size {
		panic("Malformed object file (ReadAt error)")
	}
}

func (s *Section) PutUncompressedData(b []byte) {
	if strings.HasPrefix(s.Name, "__z") {
		bb := make([]byte, 12)
		n, err := s.sr.ReadAt(bb, 0)
		if err != nil {
			panic("Malformed object file")
		}
		if n == len(bb) && string(bb[:4]) == "ZLIB" {
			size := binary.BigEndian.Uint64(bb[4:12])
			// Decompress starting at b[12:]
			r, err := zlib.NewReader(io.NewSectionReader(s, 12, int64(size)-12))
			if err != nil {
				panic("Malformed object file (zlib.NewReader error)")
			}
			n, err := io.ReadFull(r, b[0:size])
			if err != nil {
				panic("Malformed object file (ReadFull error)")
			}
			if uint64(n) != size {
				panic(fmt.Sprintf("PutUncompressedData, expected to read %d bytes, instead read %d", size, n))
			}
			if err := r.Close(); err != nil {
				panic("Malformed object file (Close error)")
			}
			return
		}
	}
	// Not compressed
	s.PutData(b)
}

func (s *Section) Copy() *Section {
	return &Section{SectionHeader: s.SectionHeader}
}

func (s *Section) String() string {
	secFlags := ""
	if !s.Flags.IsRegular() {
		secFlags = fmt.Sprintf("(%s)", s.Flags)
	}
	return fmt.Sprintf("\tsz=0x%08x off=0x%08x-0x%08x addr=0x%09x-0x%09x%s%s %s",
		s.Size,
		s.Offset,
		uint64(s.Offset)+s.Size,
		s.Addr,
		s.Addr+s.Size,
		fmt.Sprintf("%21s.%-18s", s.Seg, s.Name),
		s.Flags.Attributes(),
		secFlags)
}

func (s *Section) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name    string   `json:"name"`
		Segment string   `json:"segment"`
		Addr    uint64   `json:"addr"`
		Size    uint64   `json:"size"`
		Offset  uint32   `json:"offset"`
		Align   uint32   `json:"align"`
		Reloff  uint32   `json:"reloff"`
		Nreloc  uint32   `json:"nreloc"`
		Flags   []string `json:"flags,omitempty"`
		Type    uint8    `json:"type"`
		Relocs  []Reloc  `json:"relocs,omitempty"`
	}{
		Name:    s.Name,
		Segment: s.Seg,
		Addr:    s.Addr,
		Size:    s.Size,
		Offset:  s.Offset,
		Align:   s.Align,
		Reloff:  s.Reloff,
		Nreloc:  s.Nreloc,
		Flags:   s.Flags.List(),
		Type:    s.Type,
		Relocs:  s.Relocs,
	})
}
