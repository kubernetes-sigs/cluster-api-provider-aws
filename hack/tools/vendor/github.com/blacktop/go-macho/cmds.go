package macho

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/blacktop/go-macho/pkg/codesign"
	"github.com/blacktop/go-macho/types"
)

// A Load represents any Mach-O load command.
type Load interface {
	Command() types.LoadCmd
	LoadSize() uint32 // Need the TOC for alignment, sigh.
	Raw() []byte
	Write(buf *bytes.Buffer, o binary.ByteOrder) error
	String() string
	MarshalJSON() ([]byte, error)
}

// LoadCmdBytes is a command-tagged sequence of bytes.
// This is used for Load Commands that are not (yet)
// interesting to us, and to common up this behavior for
// all those that are.
type LoadCmdBytes struct {
	types.LoadCmd
	LoadBytes
}

func (s LoadCmdBytes) String() string {
	return s.LoadCmd.String() + ": " + s.LoadBytes.String()
}
func (s LoadCmdBytes) Copy() LoadCmdBytes {
	return LoadCmdBytes{LoadCmd: s.LoadCmd, LoadBytes: s.LoadBytes.Copy()}
}

// A LoadBytes is the uninterpreted bytes of a Mach-O load command.
type LoadBytes []byte

func (b LoadBytes) String() string {
	s := "["
	for i, a := range b {
		if i > 0 {
			s += " "
			if len(b) > 48 && i >= 16 {
				s += fmt.Sprintf("... (%d bytes)", len(b))
				break
			}
		}
		s += fmt.Sprintf("%x", a)
	}
	s += "]"
	return s
}
func (b LoadBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		LoadCmd string `json:"load_cmd"`
		Data    []byte `json:"data,omitempty"`
	}{
		LoadCmd: "unknown",
		Data:    b,
	})
}
func (b LoadBytes) Raw() []byte      { return b }
func (b LoadBytes) Copy() LoadBytes  { return LoadBytes(append([]byte{}, b...)) }
func (b LoadBytes) LoadSize() uint32 { return uint32(len(b)) }
func (b LoadBytes) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	_, err := buf.Write(b)
	return err
}

func pointerAlign(sz uint32) uint32 {
	if (sz % 8) != 0 {
		sz += 8 - (sz % 8)
	}
	return sz
}

/*******************************************************************************
 * SEGMENT
 *******************************************************************************/

// A SegmentHeader is the header for a Mach-O 32-bit or 64-bit load segment command.
type SegmentHeader struct {
	types.LoadCmd
	Len       uint32
	Name      string
	Addr      uint64
	Memsz     uint64
	Offset    uint64
	Filesz    uint64
	Maxprot   types.VmProtection
	Prot      types.VmProtection
	Nsect     uint32
	Flag      types.SegFlag
	Firstsect uint32
}

func (s *SegmentHeader) String() string {
	return fmt.Sprintf(
		"[SegmentHeader] %s, len=%#x, addr=%#x, memsz=%#x, offset=%#x, filesz=%#x, maxprot=%#x, prot=%#x, nsect=%d, flag=%#x, firstsect=%d",
		s.Name, s.Len, s.Addr, s.Memsz, s.Offset, s.Filesz, s.Maxprot, s.Prot, s.Nsect, s.Flag, s.Firstsect)
}

// A Segment represents a Mach-O 32-bit or 64-bit load segment command.
type Segment struct {
	SegmentHeader
	LoadBytes

	sections []*types.Section

	// Embed ReaderAt for ReadAt method.
	// Do not embed SectionReader directly
	// to avoid having Read and Seek.
	// If a client wants Read and Seek it must use
	// Open() to avoid fighting over the seek offset
	// with other clients.
	io.ReaderAt
	sr *io.SectionReader
}

func (s *Segment) Put32(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0*4:], uint32(s.LoadCmd))
	o.PutUint32(b[1*4:], s.Len)
	types.PutAtMost16Bytes(b[2*4:], s.Name)
	o.PutUint32(b[6*4:], uint32(s.Addr))
	o.PutUint32(b[7*4:], uint32(s.Memsz))
	o.PutUint32(b[8*4:], uint32(s.Offset))
	o.PutUint32(b[9*4:], uint32(s.Filesz))
	o.PutUint32(b[10*4:], uint32(s.Maxprot))
	o.PutUint32(b[11*4:], uint32(s.Prot))
	o.PutUint32(b[12*4:], s.Nsect)
	o.PutUint32(b[13*4:], uint32(s.Flag))
	return 14 * 4
}

func (s *Segment) Put64(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0*4:], uint32(s.LoadCmd))
	o.PutUint32(b[1*4:], s.Len)
	types.PutAtMost16Bytes(b[2*4:], s.Name)
	o.PutUint64(b[6*4+0*8:], s.Addr)
	o.PutUint64(b[6*4+1*8:], s.Memsz)
	o.PutUint64(b[6*4+2*8:], s.Offset)
	o.PutUint64(b[6*4+3*8:], s.Filesz)
	o.PutUint32(b[6*4+4*8:], uint32(s.Maxprot))
	o.PutUint32(b[7*4+4*8:], uint32(s.Prot))
	o.PutUint32(b[8*4+4*8:], s.Nsect)
	o.PutUint32(b[9*4+4*8:], uint32(s.Flag))
	return 10*4 + 4*8
}

func (s *Segment) LessThan(o *Segment) bool {
	return s.Addr < o.Addr
}

func (s *Segment) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	var name [16]byte
	copy(name[:], s.Name)

	switch s.Command() {
	case types.LC_SEGMENT:
		if err := binary.Write(buf, o, types.Segment32{
			LoadCmd: s.LoadCmd,        //              /* LC_SEGMENT */
			Len:     s.Len,            // uint32       /* includes sizeof section_64 structs */
			Name:    name,             // [16]byte     /* segment name */
			Addr:    uint32(s.Addr),   // uint32       /* memory address of this segment */
			Memsz:   uint32(s.Memsz),  // uint32       /* memory size of this segment */
			Offset:  uint32(s.Offset), // uint32       /* file offset of this segment */
			Filesz:  uint32(s.Filesz), // uint32       /* amount to map from the file */
			Maxprot: s.Maxprot,        // VmProtection /* maximum VM protection */
			Prot:    s.Prot,           // VmProtection /* initial VM protection */
			Nsect:   s.Nsect,          // uint32       /* number of sections in segment */
			Flag:    s.Flag,           // SegFlag      /* flags */
		}); err != nil {
			return fmt.Errorf("failed to write LC_SEGMENT to buffer: %v", err)
		}
	case types.LC_SEGMENT_64:
		if err := binary.Write(buf, o, types.Segment64{
			LoadCmd: s.LoadCmd, //              /* LC_SEGMENT_64 */
			Len:     s.Len,     // uint32       /* includes sizeof section_64 structs */
			Name:    name,      // [16]byte     /* segment name */
			Addr:    s.Addr,    // uint64       /* memory address of this segment */
			Memsz:   s.Memsz,   // uint64       /* memory size of this segment */
			Offset:  s.Offset,  // uint64       /* file offset of this segment */
			Filesz:  s.Filesz,  // uint64       /* amount to map from the file */
			Maxprot: s.Maxprot, // VmProtection /* maximum VM protection */
			Prot:    s.Prot,    // VmProtection /* initial VM protection */
			Nsect:   s.Nsect,   // uint32       /* number of sections in segment */
			Flag:    s.Flag,    // SegFlag      /* flags */
		}); err != nil {
			return fmt.Errorf("failed to write LC_SEGMENT to buffer: %v", err)
		}
	default:
		return fmt.Errorf("found unknown segment command: %s", s.Command().String())
	}

	return nil
}

// Data reads and returns the contents of the segment.
func (s *Segment) Data() ([]byte, error) {
	dat := make([]byte, s.Filesz)
	n, err := s.ReadAt(dat, int64(s.Offset))
	if n == len(dat) {
		err = nil
	}
	return dat[0:n], err
}

// Open returns a new ReadSeeker reading the segment.
func (s *Segment) Open() io.ReadSeeker { return io.NewSectionReader(s.sr, 0, 1<<63-1) }

// UncompressedSize returns the size of the segment with its sections uncompressed, ignoring
// its offset within the file.  The returned size is rounded up to the power of two in align.
func (s *Segment) UncompressedSize(t *FileTOC, align uint64) uint64 {
	sz := uint64(0)
	for j := uint32(0); j < s.Nsect; j++ {
		c := t.Sections[j+s.Firstsect]
		sz += c.UncompressedSize()
	}
	return (sz + align - 1) & uint64(-int64(align))
}

func (s *Segment) Copy() *Segment {
	r := &Segment{SegmentHeader: s.SegmentHeader}
	return r
}
func (s *Segment) CopyZeroed() *Segment {
	r := s.Copy()
	r.Filesz = 0
	r.Offset = 0
	r.Nsect = 0
	r.Firstsect = 0
	if s.Command() == types.LC_SEGMENT_64 {
		r.Len = uint32(unsafe.Sizeof(types.Segment64{}))
	} else {
		r.Len = uint32(unsafe.Sizeof(types.Segment32{}))
	}
	return r
}
func (s *Segment) LoadSize() uint32 {
	if s.Command() == types.LC_SEGMENT_64 {
		return uint32(unsafe.Sizeof(types.Segment64{})) + uint32(s.Nsect)*uint32(unsafe.Sizeof(types.Section64{}))
	}
	return uint32(unsafe.Sizeof(types.Segment32{})) + uint32(s.Nsect)*uint32(unsafe.Sizeof(types.Section32{}))
}

func (s *Segment) String() string {
	return fmt.Sprintf("%s sz=0x%08x off=0x%08x-0x%08x addr=0x%09x-0x%09x %s/%s   %-18s%s",
		s.Command(),
		s.Filesz,
		s.Offset,
		s.Offset+s.Filesz,
		s.Addr,
		s.Addr+s.Memsz,
		s.Prot,
		s.Maxprot,
		s.Name,
		s.Flag)
}

func (s *Segment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		LoadCmd  string           `json:"load_cmd"`
		Len      uint32           `json:"len,omitempty"`
		Name     string           `json:"name,omitempty"`
		Addr     uint64           `json:"addr,omitempty"`
		Memsz    uint64           `json:"memsz,omitempty"`
		Offset   uint64           `json:"offset,omitempty"`
		Filesz   uint64           `json:"filesz,omitempty"`
		Maxprot  string           `json:"maxprot,omitempty"`
		Prot     string           `json:"prot,omitempty"`
		Nsect    uint32           `json:"nsect,omitempty"`
		Flags    []string         `json:"flags,omitempty"`
		Sections []*types.Section `json:"sections,omitempty"`
	}{
		LoadCmd:  s.SegmentHeader.LoadCmd.String(),
		Len:      s.SegmentHeader.Len,
		Name:     s.SegmentHeader.Name,
		Addr:     s.SegmentHeader.Addr,
		Memsz:    s.SegmentHeader.Memsz,
		Offset:   s.SegmentHeader.Offset,
		Filesz:   s.SegmentHeader.Filesz,
		Maxprot:  s.SegmentHeader.Maxprot.String(),
		Prot:     s.SegmentHeader.Prot.String(),
		Nsect:    s.SegmentHeader.Nsect,
		Flags:    s.SegmentHeader.Flag.List(),
		Sections: s.sections,
	})
}

type Segments []*Segment

func (v Segments) Len() int {
	return len(v)
}

func (v Segments) Less(i, j int) bool {
	return v[i].LessThan(v[j])
}

func (v Segments) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

/*******************************************************************************
 * LC_SYMTAB
 *******************************************************************************/

// A Symtab represents a Mach-O LC_SYMTAB command.
type Symtab struct {
	LoadBytes
	types.SymtabCmd
	Syms []Symbol
}

func (s *Symtab) LoadSize() uint32 {
	return uint32(binary.Size(s.SymtabCmd))
}
func (s *Symtab) Put(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0*4:], uint32(s.LoadCmd))
	o.PutUint32(b[1*4:], s.Len)
	o.PutUint32(b[2*4:], s.Symoff)
	o.PutUint32(b[3*4:], s.Nsyms)
	o.PutUint32(b[4*4:], s.Stroff)
	o.PutUint32(b[5*4:], s.Strsize)
	return 6 * 4
}
func (s *Symtab) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, s.SymtabCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", s.Command(), err)
	}
	return nil
}
func (s *Symtab) Search(name string) (*Symbol, error) {
	for _, sym := range s.Syms {
		if sym.Name == name {
			return &sym, nil
		}
	}
	return nil, fmt.Errorf("%s not found in symtab", name)
}
func (s *Symtab) String() string {
	if s.Nsyms == 0 && s.Strsize == 0 {
		return "Symbols stripped"
	}
	return fmt.Sprintf("Symbol offset=0x%08X, Num Syms: %d, String offset=0x%08X-0x%08X", s.Symoff, s.Nsyms, s.Stroff, s.Stroff+s.Strsize)
}
func (s *Symtab) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"len,omitempty"`
		Symoff  uint32 `json:"symoff,omitempty"`
		Nsyms   uint32 `json:"nsyms,omitempty"`
		Stroff  uint32 `json:"stroff,omitempty"`
		Strsize uint32 `json:"strsize,omitempty"`
		Count   int    `json:"count,omitempty"`
	}{
		LoadCmd: s.LoadCmd.String(),
		Len:     s.Len,
		Symoff:  s.Symoff,
		Nsyms:   s.Nsyms,
		Stroff:  s.Stroff,
		Strsize: s.Strsize,
		Count:   len(s.Syms),
	})
}

// A Symbol is a Mach-O 32-bit or 64-bit symbol table entry.
type Symbol struct {
	Name  string
	Type  types.NType
	Sect  uint8
	Desc  types.NDescType
	Value uint64
}

func (s Symbol) GetType(m *File) string {
	var typ string

	if s.Type.IsUndefinedSym() {
		if s.Value != 0 {
			typ += "(common)  "
			if s.Desc.GetCommAlign() != 0 {
				typ += fmt.Sprintf("(alignment 2^%d)", s.Desc.GetCommAlign())
			}
		} else {
			if s.Type.IsPreboundUndefinedSym() {
				typ += "(prebound  "
			} else {
				typ += "("
			}
			if s.Desc.IsUndefinedLazy() {
				typ += "undefined [lazy bound]) "
			} else if s.Desc.IsPrivateUndefinedLazy() {
				typ += "undefined [private lazy bound]) "
			} else if s.Desc.IsPrivateUndefinedNonLazy() {
				typ += "undefined [private]) "
			} else {
				typ += "undefined) "
			}
		}
	} else if s.Type.IsAbsoluteSym() {
		typ += "(absolute) "
	} else if s.Type.IsIndirectSym() {
		typ += "(indirect) "
	} else if s.Type.IsDefinedInSection() {
		if s.Sect != types.NO_SECT && int(s.Sect) <= len(m.Sections) {
			typ += fmt.Sprintf("(%s,%s) ", m.Sections[s.Sect-1].Seg, m.Sections[s.Sect-1].Name)
		}
	} else if s.Type.IsDebugSym() {
		typ += "(debug) "
	} else {
		typ += "(?) "
	}

	if s.Type.IsExternalSym() {
		if s.Desc.IsReferencedDynamically() {
			typ += "[referenced dynamically] "
		}
		if s.Type.IsPrivateExternalSym() {
			if s.Desc.IsWeakDefintion() {
				typ += "weak private external "
			} else {
				typ += "private external "
			}
		} else {
			if s.Desc.IsWeakReferenced() || s.Desc.IsWeakDefintion() {
				if s.Desc.IsWeakDefintionOrReferenced() {
					typ += "weak external automatically hidden "
				} else {
					typ += "weak external "
				}
			} else {
				typ += "external "
			}
		}
	} else {
		if s.Type.IsPrivateExternalSym() {
			typ += "non-external (was a private external) "
		} else {
			typ += "private "
		}
	}

	if m.Type == types.MH_OBJECT {
		if s.Desc.IsNoDeadStrip() {
			typ += "[no dead strip] "
		}
		if !s.Type.IsUndefinedSym() && s.Desc.IsSymbolResolver() {
			typ += "[symbol resolver] "
		}
		if !s.Type.IsUndefinedSym() && s.Desc.IsAltEntry() {
			typ += "[alt entry] "
		}
		if !s.Type.IsUndefinedSym() && s.Desc.IsColdFunc() {
			typ += "[cold func] "
		}
	}

	if s.Desc.IsArmThumbDefintion() {
		typ += "[Thumb] "
	}

	if s.Type.IsIndirectSym() {
		// typ += fmt.Sprintf("(for %s)", s.Name) FIXME: find indirect symbol example
	}

	return strings.TrimSpace(typ)
}
func (s Symbol) GetLib(m *File) string {
	if (m.Flags.TwoLevel() && s.Type.IsUndefinedSym() && s.Value == 0) || s.Type.IsPreboundUndefinedSym() {
		if s.Desc.GetLibraryOrdinal() > types.SELF_LIBRARY_ORDINAL {
			switch s.Desc.GetLibraryOrdinal() {
			case types.EXECUTABLE_ORDINAL:
				return " (from executable)"
			case types.DYNAMIC_LOOKUP_ORDINAL:
				return " (dynamically looked up)"
			default:
				if s.Desc.GetLibraryOrdinal() <= uint16(len(m.ImportedLibraries())) {
					return fmt.Sprintf(" (from %s)", filepath.Base(m.ImportedLibraries()[s.Desc.GetLibraryOrdinal()-1]))
				} else {
					return fmt.Sprintf(" (from bad library ordinal %d)", s.Desc.GetLibraryOrdinal())
				}
			}
		}
	}
	return ""
}
func (s Symbol) String(m *File) string {
	return fmt.Sprintf("0x%016X\t%s\t%s%s", s.Value, s.GetType(m), s.Name, s.GetLib(m))
}
func (s *Symbol) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Sect  uint8  `json:"sect"`
		Desc  string `json:"desc"`
		Value uint64 `json:"value"`
	}{
		Name:  s.Name,
		Type:  s.Type.String(fmt.Sprintf("sect_num=%d", s.Sect)),
		Sect:  s.Sect,
		Desc:  s.Desc.String(),
		Value: s.Value,
	})
}

/*******************************************************************************
 * LC_SYMSEG - link-edit gdb symbol table info (obsolete)
 *******************************************************************************/

// A SymSeg represents a Mach-O LC_SYMSEG command.
type SymSeg struct {
	LoadBytes
	types.SymsegCmd
}

func (s *SymSeg) LoadSize() uint32 {
	return uint32(binary.Size(s.SymsegCmd))
}
func (s *SymSeg) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, s.SymsegCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", s.Command(), err)
	}
	return nil
}
func (s *SymSeg) String() string {
	return fmt.Sprintf("offset=0x%08x-0x%08x size=%5d", s.Offset, s.Offset+s.Size, s.Size)
}
func (s *SymSeg) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"len,omitempty"`
		Offset  uint32 `json:"offset,omitempty"`
		Size    uint32 `json:"size,omitempty"`
	}{
		LoadCmd: s.LoadCmd.String(),
		Len:     s.Len,
		Offset:  s.Offset,
		Size:    s.Size,
	})
}

/*******************************************************************************
 * LC_THREAD
 *******************************************************************************/

// A Thread represents a Mach-O LC_THREAD command.
type Thread struct {
	LoadBytes
	types.ThreadCmd
	bo      binary.ByteOrder
	Threads []types.ThreadState
}

func (t *Thread) LoadSize() uint32 {
	return uint32(binary.Size(t.ThreadCmd))
}
func (t *Thread) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, t.ThreadCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", t.Command(), err)
	}
	for _, thread := range t.Threads {
		if err := binary.Write(buf, o, thread.Flavor); err != nil {
			return fmt.Errorf("failed to write thread_state flavor to %s buffer: %v", t.Command(), err)
		}
		if err := binary.Write(buf, o, thread.Count); err != nil {
			return fmt.Errorf("failed to write thread_state count to %s buffer: %v", t.Command(), err)
		}
		if err := binary.Write(buf, o, thread.Data[:]); err != nil {
			return fmt.Errorf("failed to write thread_state data to %s buffer: %v", t.Command(), err)
		}
	}
	return nil
}
func (t *Thread) String() string {
	for _, thread := range t.Threads {
		if thread.Flavor == types.ARM_THREAD_STATE64 {
			regs := make([]uint64, thread.Count/2)
			binary.Read(bytes.NewReader(thread.Data), t.bo, &regs)
			return fmt.Sprintf("Threads: %d, ARM64 EntryPoint: %#016x", len(t.Threads), regs[len(regs)-2])
		}
	}
	return fmt.Sprintf("Threads: %d", len(t.Threads))
}
func (t *Thread) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		LoadCmd string              `json:"load_cmd"`
		Len     uint32              `json:"len,omitempty"`
		Threads []types.ThreadState `json:"threads,omitempty"`
	}{
		LoadCmd: t.LoadCmd.String(),
		Len:     t.Len,
		Threads: t.Threads,
	})
}

/*******************************************************************************
 * LC_UNIXTHREAD
 *******************************************************************************/

// A UnixThread represents a Mach-O LC_UNIXTHREAD command.
type UnixThread struct {
	Thread
}

/*******************************************************************************
 * LC_LOADFVMLIB - load a specified fixed VM shared library
 *******************************************************************************/

// A LoadFvmlib represents a Mach-O LC_LOADFVMLIB command.
type LoadFvmlib struct {
	LoadBytes
	types.LoadFvmLibCmd
	Name string
}

func (l *LoadFvmlib) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(l.LoadFvmLibCmd) + len(l.Name) + 1))
}
func (l *LoadFvmlib) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.LoadFvmLibCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	if _, err := buf.WriteString(l.Name + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", l.Name, l.Command(), err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", l.Command(), err)
		}
	}
	return nil
}
func (l *LoadFvmlib) String() string {
	return fmt.Sprintf("%s (%s), Header Addr: %#08x", l.Name, l.MinorVersion, l.HeaderAddr)
}
func (l *LoadFvmlib) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		LoadCmd       string `json:"load_cmd"`
		Len           uint32 `json:"len,omitempty"`
		Name          string `json:"name,omitempty"`
		MinorVersion  string `json:"minor_version,omitempty"`
		HeaderAddress uint32 `json:"header_addr,omitempty"`
	}{
		LoadCmd:       l.LoadCmd.String(),
		Len:           l.Len,
		Name:          l.Name,
		MinorVersion:  l.MinorVersion.String(),
		HeaderAddress: l.HeaderAddr,
	})
}

/*******************************************************************************
 * LC_IDFVMLIB - fixed VM shared library identification
 *******************************************************************************/

// A IDFvmlib represents a Mach-O LC_IDFVMLIB command.
type IDFvmlib struct {
	LoadFvmlib
}

/*******************************************************************************
 * LC_IDENT - object identification info (obsolete)
 *******************************************************************************/

// A Ident represents a Mach-O LC_IDENT command.
type Ident struct {
	LoadBytes
	types.IdentCmd
	StrTable []string
}

func (i *Ident) LoadSize() uint32 {
	sz := uint32(binary.Size(i.IdentCmd))
	for _, str := range i.StrTable {
		sz += uint32(len(str) + 1)
	}
	if (sz % 4) != 0 {
		sz = 4 - (sz % 4)
	}
	return sz
}
func (i *Ident) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, i.IdentCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", i.Command(), err)
	}
	for _, str := range i.StrTable {
		if _, err := buf.WriteString(str + "\x00"); err != nil {
			return fmt.Errorf("failed to write %s to %s buffer: %v", str, i.Command(), err)
		}
	}
	if (buf.Len() % 4) != 0 {
		pad := 4 - (buf.Len() % 4)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", i.Command(), err)
		}
	}
	return nil
}
func (i *Ident) String() string {
	return fmt.Sprintf("str_count=%d", len(i.StrTable))
}
func (i *Ident) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		LoadCmd string   `json:"load_cmd"`
		Len     uint32   `json:"len,omitempty"`
		Strings []string `json:"strings,omitempty"`
	}{
		LoadCmd: i.LoadCmd.String(),
		Len:     i.Len,
		Strings: i.StrTable,
	})
}

/*******************************************************************************
 * LC_FVMFILE - fixed VM file inclusion (internal use)
 *******************************************************************************/

// A FvmFile represents a Mach-O LC_FVMFILE command.
type FvmFile struct {
	LoadBytes
	types.FvmFileCmd
	Name string
}

func (l *FvmFile) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(l.FvmFileCmd) + len(l.Name) + 1))
}
func (l *FvmFile) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.FvmFileCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	if _, err := buf.WriteString(l.Name + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", l.Name, l.Command(), err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", l.Command(), err)
		}
	}
	return nil
}
func (l *FvmFile) String() string {
	return fmt.Sprintf("%s, Header Addr: %#08x", l.Name, l.HeaderAddr)
}
func (l *FvmFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		LoadCmd       string `json:"load_cmd"`
		Len           uint32 `json:"len,omitempty"`
		Name          string `json:"name,omitempty"`
		HeaderAddress uint32 `json:"header_addr,omitempty"`
	}{
		LoadCmd:       l.LoadCmd.String(),
		Len:           l.Len,
		Name:          l.Name,
		HeaderAddress: l.HeaderAddr,
	})
}

/*******************************************************************************
 * LC_PREPAGE - prepage command (internal use)
 *******************************************************************************/

// A Prepage represents a Mach-O LC_PREPAGE command.
type Prepage struct {
	LoadBytes
	types.PrePageCmd
}

func (c *Prepage) LoadSize() uint32 {
	return uint32(binary.Size(c.PrePageCmd))
}
func (c *Prepage) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, c.PrePageCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", c.Command(), err)
	}
	return nil
}
func (c *Prepage) String() string {
	return fmt.Sprintf("size=%d", c.Len)
}
func (c *Prepage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"len,omitempty"`
	}{
		LoadCmd: c.LoadCmd.String(),
		Len:     c.Len,
	})
}

/*******************************************************************************
 * LC_DYSYMTAB
 *******************************************************************************/

// A Dysymtab represents a Mach-O LC_DYSYMTAB command.
type Dysymtab struct {
	LoadBytes
	types.DysymtabCmd
	IndirectSyms []uint32 // indices into Symtab.Syms
}

func (d *Dysymtab) LoadSize() uint32 {
	return uint32(binary.Size(d.DysymtabCmd))
}
func (d *Dysymtab) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, d.DysymtabCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", d.Command(), err)
	}
	return nil
}
func (d *Dysymtab) String() string {
	var tocStr, modStr, extSymStr, indirSymStr, extRelStr, locRelStr string
	if d.Ntoc == 0 {
		tocStr = "No"
	} else {
		tocStr = fmt.Sprintf("%d at 0x%08x", d.Ntoc, d.Tocoffset)
	}
	if d.Nmodtab == 0 {
		modStr = "No"
	} else {
		modStr = fmt.Sprintf("%d at 0x%08x", d.Nmodtab, d.Modtaboff)
	}
	if d.Nextrefsyms == 0 {
		extSymStr = "None"
	} else {
		extSymStr = fmt.Sprintf("%d at 0x%08x", d.Nextrefsyms, d.Extrefsymoff)
	}
	if d.Nindirectsyms == 0 {
		indirSymStr = "None"
	} else {
		indirSymStr = fmt.Sprintf("%d at 0x%08x", d.Nindirectsyms, d.Indirectsymoff)
	}
	if d.Nextrel == 0 {
		extRelStr = "None"
	} else {
		extRelStr = fmt.Sprintf("%d at 0x%08x", d.Nextrel, d.Extreloff)
	}
	if d.Nlocrel == 0 {
		locRelStr = "None"
	} else {
		locRelStr = fmt.Sprintf("%d at 0x%08x", d.Nlocrel, d.Locreloff)
	}
	return fmt.Sprintf(
		"\n"+
			"\t             Local Syms: %d at %d\n"+
			"\t          External Syms: %d at %d\n"+
			"\t         Undefined Syms: %d at %d\n"+
			"\t                    TOC: %s\n"+
			"\t                 Modtab: %s\n"+
			"\tExternal symtab Entries: %s\n"+
			"\tIndirect symtab Entries: %s\n"+
			"\t External Reloc Entries: %s\n"+
			"\t    Local Reloc Entries: %s",
		d.Nlocalsym, d.Ilocalsym,
		d.Nextdefsym, d.Iextdefsym,
		d.Nundefsym, d.Iundefsym,
		tocStr,
		modStr,
		extSymStr,
		indirSymStr,
		extRelStr,
		locRelStr)
}
func (d *Dysymtab) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd        string `json:"load_cmd"`
		Len            uint32 `json:"length"`
		Ilocalsym      uint32 `json:"ilocalsym"`
		Nlocalsym      uint32 `json:"nlocalsym"`
		Iextdefsym     uint32 `json:"iextdefsym"`
		Nextdefsym     uint32 `json:"nextdefsym"`
		Iundefsym      uint32 `json:"iundefsym"`
		Nundefsym      uint32 `json:"nundefsym"`
		Tocoffset      uint32 `json:"tocoffset"`
		Ntoc           uint32 `json:"ntoc"`
		Modtaboff      uint32 `json:"modtaboff"`
		Nmodtab        uint32 `json:"nmodtab"`
		Extrefsymoff   uint32 `json:"extrefsymoff"`
		Nextrefsyms    uint32 `json:"nextrefsyms"`
		Indirectsymoff uint32 `json:"indirectsymoff"`
		Nindirectsyms  uint32 `json:"nindirectsyms"`
		Extreloff      uint32 `json:"extreloff"`
		Nextrel        uint32 `json:"nextrel"`
		Locreloff      uint32 `json:"locreloff"`
		Nlocrel        uint32 `json:"nlocrel"`
	}{
		LoadCmd:        d.Command().String(),
		Len:            d.LoadSize(),
		Ilocalsym:      d.Ilocalsym,
		Nlocalsym:      d.Nlocalsym,
		Iextdefsym:     d.Iextdefsym,
		Nextdefsym:     d.Nextdefsym,
		Iundefsym:      d.Iundefsym,
		Nundefsym:      d.Nundefsym,
		Tocoffset:      d.Tocoffset,
		Ntoc:           d.Ntoc,
		Modtaboff:      d.Modtaboff,
		Nmodtab:        d.Nmodtab,
		Extrefsymoff:   d.Extrefsymoff,
		Nextrefsyms:    d.Nextrefsyms,
		Indirectsymoff: d.Indirectsymoff,
		Nindirectsyms:  d.Nindirectsyms,
		Extreloff:      d.Extreloff,
		Nextrel:        d.Nextrel,
		Locreloff:      d.Locreloff,
		Nlocrel:        d.Nlocrel,
	})
}

/*******************************************************************************
 * LC_LOAD_DYLIB
 *******************************************************************************/

// A LoadDylib represents a Mach-O LC_LOAD_DYLIB command.
type LoadDylib struct {
	Dylib
}

/*******************************************************************************
 * LC_ID_DYLIB
 *******************************************************************************/

// A IDDylib represents a Mach-O LC_ID_DYLIB command.
type IDDylib struct {
	Dylib
}

/*******************************************************************************
 * LC_LOAD_DYLINKER
 *******************************************************************************/

// A LoadDylinker represents a Mach-O LC_LOAD_DYLINKER command.
type LoadDylinker struct {
	Dylinker
}

/*******************************************************************************
 * LC_ID_DYLINKER
 *******************************************************************************/

// DylinkerID represents a Mach-O LC_ID_DYLINKER command.
type DylinkerID struct {
	Dylinker
}

/*******************************************************************************
 * LC_PREBOUND_DYLIB - modules prebound for a dynamically linked shared library
 *******************************************************************************/

// PreboundDylib represents a Mach-O LC_PREBOUND_DYLIB command.
type PreboundDylib struct {
	LoadBytes
	types.PreboundDylibCmd
	Name                   string
	LinkedModulesBitVector string
}

func (d *PreboundDylib) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(d.PreboundDylibCmd) + len(d.Name) + 1 + len(d.LinkedModulesBitVector) + 1))
}
func (d *PreboundDylib) Put(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0*4:], uint32(d.LoadCmd))
	o.PutUint32(b[1*4:], d.Len)
	o.PutUint32(b[2*4:], d.NameOffset)
	o.PutUint32(b[3*4:], d.NumModules)
	o.PutUint32(b[4*4:], d.LinkedModulesOffset)
	return 5 * binary.Size(uint32(0))
}
func (d *PreboundDylib) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, d.PreboundDylibCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", d.Command(), err)
	}
	if _, err := buf.WriteString(d.Name + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", d.Name, d.Command(), err)
	}
	if _, err := buf.WriteString(d.LinkedModulesBitVector + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", d.LinkedModulesBitVector, d.Command(), err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", d.Command(), err)
		}
	}
	return nil
}
func (d *PreboundDylib) String() string {
	return fmt.Sprintf("%s, NumModules=%d, LinkedModulesBitVector=%v", d.Name, d.NumModules, []byte(d.LinkedModulesBitVector))
}
func (d *PreboundDylib) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd                string `json:"load_cmd"`
		Len                    uint32 `json:"length"`
		Name                   string `json:"name"`
		NumModules             uint32 `json:"num_modules"`
		LinkedModulesBitVector string `json:"linked_modules"`
	}{
		LoadCmd:                d.Command().String(),
		Len:                    d.Len,
		Name:                   d.Name,
		NumModules:             d.NumModules,
		LinkedModulesBitVector: d.LinkedModulesBitVector,
	})
}

/*******************************************************************************
 * LC_ROUTINES - image routines
 *******************************************************************************/

// A Routines is a Mach-O LC_ROUTINES command.
type Routines struct {
	LoadBytes
	types.RoutinesCmd
}

func (l *Routines) LoadSize() uint32 {
	return uint32(binary.Size(l.RoutinesCmd))
}
func (l *Routines) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.RoutinesCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	return nil
}
func (l *Routines) String() string {
	return fmt.Sprintf("Address: %#08x, Module: %d", l.InitAddress, l.InitModule)
}
func (l *Routines) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd     string `json:"load_cmd"`
		Len         uint32 `json:"length"`
		InitAddress uint32 `json:"init_address"`
		InitModule  uint32 `json:"init_module"`
	}{
		LoadCmd:     l.Command().String(),
		Len:         l.Len,
		InitAddress: l.InitAddress,
		InitModule:  l.InitModule,
	})
}

/*******************************************************************************
 * LC_SUB_FRAMEWORK
 *******************************************************************************/

// A SubFramework is a Mach-O LC_SUB_FRAMEWORK command.
type SubFramework struct {
	LoadBytes
	types.SubFrameworkCmd
	Framework string
}

func (l *SubFramework) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(l.SubFrameworkCmd) + len(l.Framework) + 1))
}
func (l *SubFramework) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.SubFrameworkCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	if _, err := buf.WriteString(l.Framework + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", l.Framework, l.Command(), err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", l.Command(), err)
		}
	}
	return nil
}
func (l *SubFramework) String() string { return l.Framework }
func (l *SubFramework) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd   string `json:"load_cmd"`
		Len       uint32 `json:"length"`
		Framework string `json:"framework"`
	}{
		LoadCmd:   l.Command().String(),
		Len:       l.Len,
		Framework: l.Framework,
	})
}

/*******************************************************************************
 * LC_SUB_UMBRELLA - sub umbrella
 *******************************************************************************/

// A SubUmbrella is a Mach-O LC_SUB_UMBRELLA command.
type SubUmbrella struct {
	LoadBytes
	types.SubUmbrellaCmd
	Umbrella string
}

func (l *SubUmbrella) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(l.SubUmbrellaCmd)) + uint32(len(l.Umbrella)) + 1)
}
func (l *SubUmbrella) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.SubUmbrellaCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	if _, err := buf.WriteString(l.Umbrella + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", l.Umbrella, l.Command(), err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", l.Command(), err)
		}
	}
	return nil
}
func (l *SubUmbrella) String() string { return l.Umbrella }
func (l *SubUmbrella) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd  string `json:"load_cmd"`
		Len      uint32 `json:"length"`
		Umbrella string `json:"umbrella"`
	}{
		LoadCmd:  l.Command().String(),
		Len:      l.Len,
		Umbrella: l.Umbrella,
	})
}

/*******************************************************************************
 * LC_SUB_CLIENT
 *******************************************************************************/

// A SubClient is a Mach-O LC_SUB_CLIENT command.
type SubClient struct {
	LoadBytes
	types.SubClientCmd
	Name string
}

func (l *SubClient) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(l.SubClientCmd)) + uint32(len(l.Name)) + 1)
}
func (l *SubClient) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.SubClientCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	if _, err := buf.WriteString(l.Name + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", l.Name, l.Command(), err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", l.Command(), err)
		}
	}
	return nil
}
func (l *SubClient) String() string {
	return l.Name
}
func (l *SubClient) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Name    string `json:"name"`
	}{
		LoadCmd: l.Command().String(),
		Len:     l.Len,
		Name:    l.Name,
	})
}

/*******************************************************************************
 * LC_SUB_LIBRARY - sub library
 *******************************************************************************/

// A SubLibrary is a Mach-O LC_SUB_LIBRARY command.
type SubLibrary struct {
	LoadBytes
	types.SubLibraryCmd
	Library string
}

func (l *SubLibrary) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(l.SubLibraryCmd)) + uint32(len(l.Library)) + 1)
}
func (l *SubLibrary) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.SubLibraryCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	if _, err := buf.WriteString(l.Library + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", l.Library, l.Command(), err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", l.Command(), err)
		}
	}
	return nil
}
func (l *SubLibrary) String() string { return l.Library }
func (l *SubLibrary) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Library string `json:"library"`
	}{
		LoadCmd: l.Command().String(),
		Len:     l.Len,
		Library: l.Library,
	})
}

/*******************************************************************************
 * LC_TWOLEVEL_HINTS - two-level namespace lookup hints
 *******************************************************************************/

// A TwolevelHints  is a Mach-O LC_TWOLEVEL_HINTS command.
type TwolevelHints struct {
	LoadBytes
	types.TwolevelHintsCmd
	Hints []types.TwolevelHint
}

func (l *TwolevelHints) LoadSize() uint32 {
	return uint32(binary.Size(l.TwolevelHintsCmd) + binary.Size(l.Hints))
}
func (l *TwolevelHints) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.TwolevelHintsCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	if err := binary.Write(buf, o, l.Hints); err != nil {
		return fmt.Errorf("failed to write hints to %s buffer: %v", l.Command(), err)
	}
	return nil
}
func (l *TwolevelHints) String() string {
	return fmt.Sprintf("Offset: %#08x, Num of Hints: %d", l.Offset, len(l.Hints))
}
func (l *TwolevelHints) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string               `json:"load_cmd"`
		Len     uint32               `json:"length"`
		Offset  uint32               `json:"offset"`
		Hints   []types.TwolevelHint `json:"hints,omitempty"`
	}{
		LoadCmd: l.Command().String(),
		Len:     l.Len,
		Offset:  l.Offset,
		Hints:   l.Hints,
	})
}

/*******************************************************************************
 * LC_PREBIND_CKSUM - prebind checksum
 *******************************************************************************/

// A PrebindCheckSum  is a Mach-O LC_PREBIND_CKSUM command.
type PrebindCheckSum struct {
	LoadBytes
	types.PrebindCksumCmd
}

func (l *PrebindCheckSum) LoadSize() uint32 {
	return uint32(binary.Size(l.PrebindCksumCmd))
}
func (l *PrebindCheckSum) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.PrebindCksumCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	return nil
}
func (l *PrebindCheckSum) String() string {
	return fmt.Sprintf("CheckSum: %#08x", l.CheckSum)
}
func (l *PrebindCheckSum) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd  string `json:"load_cmd"`
		Len      uint32 `json:"length"`
		CheckSum uint32 `json:"checksum"`
	}{
		LoadCmd:  l.Command().String(),
		Len:      l.Len,
		CheckSum: l.CheckSum,
	})
}

/*******************************************************************************
 * LC_LOAD_WEAK_DYLIB
 *******************************************************************************/

// A WeakDylib represents a Mach-O LC_LOAD_WEAK_DYLIB command.
type WeakDylib struct {
	Dylib
}

/*******************************************************************************
 * LC_ROUTINES_64
 *******************************************************************************/

// A Routines64 is a Mach-O LC_ROUTINES_64 command.
type Routines64 struct {
	LoadBytes
	types.Routines64Cmd
}

func (l *Routines64) LoadSize() uint32 {
	return uint32(binary.Size(l.Routines64Cmd))
}
func (l *Routines64) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.Routines64Cmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	return nil
}
func (l *Routines64) String() string {
	return fmt.Sprintf("Address: %#016x, Module: %d", l.InitAddress, l.InitModule)
}
func (l *Routines64) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd     string `json:"load_cmd"`
		Len         uint32 `json:"length"`
		InitAddress uint64 `json:"init_address"`
		InitModule  uint64 `json:"init_module"`
	}{
		LoadCmd:     l.Command().String(),
		Len:         l.Len,
		InitAddress: l.InitAddress,
		InitModule:  l.InitModule,
	})
}

/*******************************************************************************
 * LC_UUID
 *******************************************************************************/

// UUID represents a Mach-O LC_UUID command.
type UUID struct {
	LoadBytes
	types.UUIDCmd
}

func (l *UUID) LoadSize() uint32 {
	return uint32(binary.Size(l.UUIDCmd))
}
func (l *UUID) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.UUIDCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	return nil
}
func (l *UUID) String() string {
	return l.UUID.String()
}
func (l *UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		UUID    string `json:"uuid"`
	}{
		LoadCmd: l.Command().String(),
		Len:     l.Len,
		UUID:    l.UUID.String(),
	})
}

/*******************************************************************************
 * LC_RPATH
 *******************************************************************************/

// A Rpath represents a Mach-O LC_RPATH command.
type Rpath struct {
	LoadBytes
	types.RpathCmd
	Path string
}

func (r *Rpath) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(r.RpathCmd)) + uint32(len(r.Path)) + 1)
}
func (r *Rpath) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, r.RpathCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", r.Command(), err)
	}
	if _, err := buf.WriteString(r.Path + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to Dylib buffer: %v", r.Path, err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", r.Command(), err)
		}
	}
	return nil
}
func (r *Rpath) String() string {
	return r.Path
}
func (r *Rpath) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Path    string `json:"path"`
	}{
		LoadCmd: r.Command().String(),
		Len:     r.Len,
		Path:    r.Path,
	})
}

/*******************************************************************************
 * LC_CODE_SIGNATURE
 *******************************************************************************/

// A CodeSignature represents a Mach-O LC_CODE_SIGNATURE command.
type CodeSignature struct {
	LoadBytes
	types.CodeSignatureCmd
	codesign.CodeSignature
}

func (l *CodeSignature) LoadSize() uint32 {
	return uint32(binary.Size(l.CodeSignatureCmd))
}
func (l *CodeSignature) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.CodeSignatureCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	return nil
}
func (l *CodeSignature) String() string { // TODO: add more info
	return fmt.Sprintf("offset=0x%09x  size=%#x", l.Offset, l.Size)
}
func (l *CodeSignature) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd       string                  `json:"load_cmd"`
		Len           uint32                  `json:"length"`
		Offset        uint32                  `json:"offset"`
		Size          uint32                  `json:"size"`
		CodeSignature *codesign.CodeSignature `json:"code_signature,omitempty"`
	}{
		LoadCmd:       l.Command().String(),
		Len:           l.Len,
		Offset:        l.Offset,
		Size:          l.Size,
		CodeSignature: nil, // TODO: add MarshalJSON for CodeSignature
	})
}

/*******************************************************************************
 * LC_SEGMENT_SPLIT_INFO
 *******************************************************************************/

// A SplitInfo represents a Mach-O LC_SEGMENT_SPLIT_INFO command.
type SplitInfo struct {
	LoadBytes
	types.SegmentSplitInfoCmd
	Version uint8
}

func (l *SplitInfo) LoadSize() uint32 {
	return uint32(binary.Size(l.SegmentSplitInfoCmd))
}
func (l *SplitInfo) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.SegmentSplitInfoCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	return nil
}
func (s *SplitInfo) String() string {
	version := "format=v1"
	if s.Version == types.DYLD_CACHE_ADJ_V2_FORMAT {
		version = "format=v2"
	} else {
		version = fmt.Sprintf("kind=%#x", s.Version)
	}
	return fmt.Sprintf("offset=0x%08x-0x%08x size=%5d, %s", s.Offset, s.Offset+s.Size, s.Size, version)
}
func (l *SplitInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Offset  uint32 `json:"offset"`
		Size    uint32 `json:"size"`
		Version uint8  `json:"version"`
	}{
		LoadCmd: l.Command().String(),
		Len:     l.Len,
		Offset:  l.Offset,
		Size:    l.Size,
		Version: l.Version,
	})
}

/*******************************************************************************
 * LC_REEXPORT_DYLIB
 *******************************************************************************/

// A ReExportDylib represents a Mach-O LC_REEXPORT_DYLIB command.
type ReExportDylib struct {
	Dylib
}

/*******************************************************************************
 * LC_LAZY_LOAD_DYLIB - delay load of dylib until first use
 *******************************************************************************/

// A LazyLoadDylib represents a Mach-O LC_LAZY_LOAD_DYLIB command.
type LazyLoadDylib struct {
	Dylib
}

/*******************************************************************************
 * LC_ENCRYPTION_INFO
 *******************************************************************************/

// A EncryptionInfo represents a Mach-O 32-bit encrypted segment information
type EncryptionInfo struct {
	LoadBytes
	types.EncryptionInfoCmd
}

func (e *EncryptionInfo) LoadSize() uint32 {
	return uint32(binary.Size(e.EncryptionInfoCmd))
}
func (e *EncryptionInfo) Put(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0*4:], uint32(e.LoadCmd))
	o.PutUint32(b[1*4:], e.Len)
	o.PutUint32(b[2*4:], e.Offset)
	o.PutUint32(b[3*4:], e.Size)
	o.PutUint32(b[3*4:], uint32(e.CryptID))

	return int(e.Len)
}
func (l *EncryptionInfo) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.EncryptionInfoCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	return nil
}
func (e *EncryptionInfo) String() string {
	if e.CryptID == 0 {
		return fmt.Sprintf("offset=%#x size=%#x (not-encrypted yet)", e.Offset, e.Size)
	}
	return fmt.Sprintf("offset=%#x size=%#x CryptID: %#x", e.Offset, e.Size, e.CryptID)
}
func (l *EncryptionInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Offset  uint32 `json:"offset"`
		Size    uint32 `json:"size"`
		CryptID uint32 `json:"crypt_id"`
	}{
		LoadCmd: l.Command().String(),
		Len:     l.Len,
		Offset:  l.Offset,
		Size:    l.Size,
		CryptID: uint32(l.CryptID),
	})
}

/*******************************************************************************
 * LC_DYLD_INFO
 *******************************************************************************/

// A DyldInfo represents a Mach-O LC_DYLD_INFO command.
type DyldInfo struct {
	LoadBytes
	types.DyldInfoCmd
}

func (d *DyldInfo) LoadSize() uint32 {
	return uint32(binary.Size(d.DyldInfoCmd))
}
func (d *DyldInfo) Put(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0*4:], uint32(d.LoadCmd))
	o.PutUint32(b[1*4:], d.Len)
	o.PutUint32(b[2*4:], d.RebaseOff)
	o.PutUint32(b[3*4:], d.RebaseSize)
	o.PutUint32(b[4*4:], d.BindOff)
	o.PutUint32(b[5*4:], d.BindSize)
	o.PutUint32(b[6*4:], d.WeakBindOff)
	o.PutUint32(b[7*4:], d.WeakBindSize)
	o.PutUint32(b[8*4:], d.LazyBindOff)
	o.PutUint32(b[9*4:], d.LazyBindSize)
	o.PutUint32(b[10*4:], d.ExportOff)
	o.PutUint32(b[11*4:], d.ExportSize)
	return int(d.Len)
}
func (l *DyldInfo) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.DyldInfoCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	return nil
}

func (d *DyldInfo) String() string {
	return fmt.Sprintf(
		"\n"+
			"\t\tRebase info: %5d bytes at offset:  0x%08X -> 0x%08X\n"+
			"\t\tBind info:   %5d bytes at offset:  0x%08X -> 0x%08X\n"+
			"\t\tWeak info:   %5d bytes at offset:  0x%08X -> 0x%08X\n"+
			"\t\tLazy info:   %5d bytes at offset:  0x%08X -> 0x%08X\n"+
			"\t\tExport info: %5d bytes at offset:  0x%08X -> 0x%08X",
		d.RebaseSize, d.RebaseOff, d.RebaseOff+d.RebaseSize,
		d.BindSize, d.BindOff, d.BindOff+d.BindSize,
		d.WeakBindSize, d.WeakBindOff, d.WeakBindOff+d.WeakBindSize,
		d.LazyBindSize, d.LazyBindOff, d.LazyBindOff+d.LazyBindSize,
		d.ExportSize, d.ExportOff, d.ExportOff+d.ExportSize,
	)
}
func (l *DyldInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd      string `json:"load_cmd"`
		Len          uint32 `json:"length"`
		RebaseOff    uint32 `json:"rebase_offset"`
		RebaseSize   uint32 `json:"rebase_size"`
		BindOff      uint32 `json:"bind_offset"`
		BindSize     uint32 `json:"bind_size"`
		WeakBindOff  uint32 `json:"weak_bind_offset"`
		WeakBindSize uint32 `json:"weak_bind_size"`
		LazyBindOff  uint32 `json:"lazy_bind_offset"`
		LazyBindSize uint32 `json:"lazy_bind_size"`
		ExportOff    uint32 `json:"export_offset"`
		ExportSize   uint32 `json:"export_size"`
	}{
		LoadCmd:      l.Command().String(),
		Len:          l.Len,
		RebaseOff:    l.RebaseOff,
		RebaseSize:   l.RebaseSize,
		BindOff:      l.BindOff,
		BindSize:     l.BindSize,
		WeakBindOff:  l.WeakBindOff,
		WeakBindSize: l.WeakBindSize,
		LazyBindOff:  l.LazyBindOff,
		LazyBindSize: l.LazyBindSize,
		ExportOff:    l.ExportOff,
		ExportSize:   l.ExportSize,
	})
}

/*******************************************************************************
 * LC_DYLD_INFO_ONLY
 *******************************************************************************/

// DyldInfoOnly is compressed dyld information only
type DyldInfoOnly struct {
	DyldInfo
}

/*******************************************************************************
 * LC_LOAD_UPWARD_DYLIB
 *******************************************************************************/

// A UpwardDylib represents a Mach-O LC_LOAD_UPWARD_DYLIB load command.
type UpwardDylib struct {
	Dylib
}

/*******************************************************************************
 * LC_VERSION_MIN_MACOSX
 *******************************************************************************/

// VersionMinMacOSX build for MacOSX min OS version
type VersionMinMacOSX struct {
	VersionMin
}

/*******************************************************************************
 * LC_VERSION_MIN_IPHONEOS
 *******************************************************************************/

// VersionMiniPhoneOS build for iPhoneOS min OS version
type VersionMiniPhoneOS struct {
	VersionMin
}

/*******************************************************************************
 * LC_FUNCTION_STARTS
 *******************************************************************************/

// A FunctionStarts represents a Mach-O function starts command.
type FunctionStarts struct {
	LinkEditData
}

/*******************************************************************************
 * LC_DYLD_ENVIRONMENT
 *******************************************************************************/

// DyldEnvironment represents a Mach-O LC_DYLD_ENVIRONMENT command.
type DyldEnvironment struct {
	Dylinker
}

/*******************************************************************************
 * LC_MAIN
 *******************************************************************************/

// EntryPoint represents a Mach-O LC_MAIN command.
type EntryPoint struct {
	LoadBytes
	types.EntryPointCmd
}

func (e *EntryPoint) LoadSize() uint32 {
	return uint32(binary.Size(e.EntryPointCmd))
}
func (e *EntryPoint) Put(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0*4:], uint32(e.LoadCmd))
	o.PutUint32(b[1*4:], e.Len)
	o.PutUint64(b[2*8:], e.EntryOffset)
	o.PutUint64(b[3*8:], e.StackSize)
	return int(e.Len)
}
func (e *EntryPoint) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, e.EntryPointCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", e.Command(), err)
	}
	return nil
}
func (e *EntryPoint) String() string {
	return fmt.Sprintf("Entry Point: 0x%016x, Stack Size: %#x", e.EntryOffset, e.StackSize)
}
func (e *EntryPoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Entry   uint64 `json:"entry_offset"`
		Stack   uint64 `json:"stack_size"`
	}{
		LoadCmd: e.Command().String(),
		Len:     e.Len,
		Entry:   e.EntryOffset,
		Stack:   e.StackSize,
	})
}

/*******************************************************************************
 * LC_DATA_IN_CODE
 *******************************************************************************/

// A DataInCode represents a Mach-O LC_DATA_IN_CODE command.
type DataInCode struct {
	LoadBytes
	types.DataInCodeCmd
	Entries []types.DataInCodeEntry
}

func (l *DataInCode) LoadSize() uint32 {
	return uint32(binary.Size(l.DataInCodeCmd))
}
func (l *DataInCode) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.DataInCodeCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	return nil
}
func (d *DataInCode) String() string {
	var ents string
	if len(d.Entries) > 0 {
		ents = "\n"
	}
	for _, e := range d.Entries {
		ents += fmt.Sprintf("\toffset: %#08x length: %d kind: %s\n", e.Offset, e.Length, e.Kind)
	}
	ents = strings.TrimSuffix(ents, "\n")
	return fmt.Sprintf(
		"offset=0x%08x-0x%08x size=%5d entries=%d%s",
		d.Offset, d.Offset+d.Size, d.Size, len(d.Entries), ents)
}
func (l *DataInCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Offset  uint32 `json:"offset"`
		Size    uint32 `json:"size"`
		Entries []struct {
			Offset uint32 `json:"offset"`
			Length uint16 `json:"length"`
			Kind   string `json:"kind"`
		} `json:"entries"`
	}{
		LoadCmd: l.Command().String(),
		Len:     l.Len,
		Offset:  l.Offset,
		Size:    l.Size,
		Entries: func() []struct {
			Offset uint32 `json:"offset"`
			Length uint16 `json:"length"`
			Kind   string `json:"kind"`
		} {
			ents := make([]struct {
				Offset uint32 `json:"offset"`
				Length uint16 `json:"length"`
				Kind   string `json:"kind"`
			}, len(l.Entries))
			for i, e := range l.Entries {
				ents[i].Offset = e.Offset
				ents[i].Length = e.Length
				ents[i].Kind = e.Kind.String()
			}
			return ents
		}(),
	})
}

/*******************************************************************************
 * LC_SOURCE_VERSION
 *******************************************************************************/

// A SourceVersion represents a Mach-O LC_SOURCE_VERSION command.
type SourceVersion struct {
	LoadBytes
	types.SourceVersionCmd
}

func (s *SourceVersion) LoadSize() uint32 {
	return uint32(binary.Size(s.SourceVersionCmd))
}
func (s *SourceVersion) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, s.SourceVersionCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", s.Command(), err)
	}
	return nil
}
func (s *SourceVersion) String() string {
	return s.Version.String()
}
func (s *SourceVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Version string `json:"version"`
	}{
		LoadCmd: s.Command().String(),
		Len:     s.Len,
		Version: s.Version.String(),
	})
}

/*******************************************************************************
 * LC_DYLIB_CODE_SIGN_DRS Code signing DRs copied from linked dylibs
 *******************************************************************************/

type DylibCodeSignDrs struct {
	LinkEditData
}

/*******************************************************************************
 * LC_ENCRYPTION_INFO_64
 *******************************************************************************/

// A EncryptionInfo64 represents a Mach-O 64-bit encrypted segment information
type EncryptionInfo64 struct {
	LoadBytes
	types.EncryptionInfo64Cmd
}

func (e *EncryptionInfo64) LoadSize() uint32 {
	return uint32(binary.Size(e.EncryptionInfo64Cmd))
}
func (e *EncryptionInfo64) Put(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0*4:], uint32(e.LoadCmd))
	o.PutUint32(b[1*4:], e.Len)
	o.PutUint32(b[2*4:], e.Offset)
	o.PutUint32(b[3*4:], e.Size)
	o.PutUint32(b[3*4:], uint32(e.CryptID))
	o.PutUint32(b[3*4:], e.Pad)

	return int(e.Len)
}
func (e *EncryptionInfo64) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, e.EncryptionInfo64Cmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", e.Command(), err)
	}
	return nil
}
func (e *EncryptionInfo64) String() string {
	if e.CryptID == 0 {
		return fmt.Sprintf("offset=0x%09x  size=%#x (not-encrypted yet)", e.Offset, e.Size)
	}
	return fmt.Sprintf("offset=0x%09x  size=%#x CryptID: %#x", e.Offset, e.Size, e.CryptID)
}
func (e *EncryptionInfo64) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Offset  uint32 `json:"offset"`
		Size    uint32 `json:"size"`
		CryptID uint32 `json:"crypt_id"`
		Pad     uint32 `json:"pad"`
	}{
		LoadCmd: e.Command().String(),
		Len:     e.Len,
		Offset:  e.Offset,
		Size:    e.Size,
		CryptID: uint32(e.CryptID),
		Pad:     e.Pad,
	})
}

/*******************************************************************************
 * LC_LINKER_OPTION - linker options in MH_OBJECT files
 *******************************************************************************/

// A LinkerOption represents a Mach-O LC_LINKER_OPTION command.
type LinkerOption struct {
	LoadBytes
	types.LinkerOptionCmd
	Options []string
}

func (l *LinkerOption) LoadSize() uint32 {
	var totalStrLen uint32
	for _, opt := range l.Options {
		totalStrLen += uint32(len(opt) + 1)
	}
	return uint32(binary.Size(l.LinkerOptionCmd)) + totalStrLen
}
func (l *LinkerOption) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.LinkerOptionCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	for _, opt := range l.Options {
		if _, err := buf.WriteString(opt + "\x00"); err != nil {
			return fmt.Errorf("failed to write %s to %s buffer: %v", opt, l.Command(), err)
		}
	}
	return nil
}
func (l *LinkerOption) String() string {
	return fmt.Sprintf("Options=%s", strings.Join(l.Options, ","))
}
func (l *LinkerOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string   `json:"load_cmd"`
		Len     uint32   `json:"length"`
		Options []string `json:"options"`
	}{
		LoadCmd: l.Command().String(),
		Len:     l.Len,
		Options: l.Options,
	})
}

/*******************************************************************************
 * LC_LINKER_OPTIMIZATION_HINT - linker options in MH_OBJECT files
 *******************************************************************************/

type LinkerOptimizationHint struct {
	LinkEditData
}

/*******************************************************************************
 * LC_VERSION_MIN_TVOS
 *******************************************************************************/

// VersionMinTvOS build for AppleTV min OS version
type VersionMinTvOS struct {
	VersionMin
}

/*******************************************************************************
 * LC_VERSION_MIN_WATCHOS
 *******************************************************************************/

// VersionMinWatchOS build for Watch min OS version
type VersionMinWatchOS struct {
	VersionMin
}

/*******************************************************************************
 * LC_NOTE - arbitrary data included within a Mach-O file
 *******************************************************************************/

// A Note represents a Mach-O LC_NOTE command.
type Note struct {
	LoadBytes
	types.NoteCmd
}

func (n *Note) LoadSize() uint32 {
	return uint32(binary.Size(n.NoteCmd))
}
func (n *Note) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, n.NoteCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", n.Command(), err)
	}
	return nil
}
func (n *Note) String() string {
	return fmt.Sprintf("DataOwner=%s, offset=0x%08x-0x%08x size=%5d", string(n.DataOwner[:]), n.Offset, n.Offset+n.Size, n.Size)
}
func (n *Note) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd   string `json:"load_cmd"`
		Len       uint32 `json:"length"`
		DataOwner string `json:"data_owner"`
		Offset    uint64 `json:"offset"`
		Size      uint64 `json:"size"`
	}{
		LoadCmd:   n.Command().String(),
		Len:       n.Len,
		DataOwner: string(n.DataOwner[:]),
		Offset:    n.Offset,
		Size:      n.Size,
	})
}

/*******************************************************************************
 * LC_BUILD_VERSION
 *******************************************************************************/

// A BuildVersion represents a Mach-O build for platform min OS version.
type BuildVersion struct {
	LoadBytes
	types.BuildVersionCmd
	Tools []types.BuildVersionTool
}

func (b *BuildVersion) LoadSize() uint32 {
	return uint32(binary.Size(b.BuildVersionCmd) + binary.Size(b.Tools))
}
func (b *BuildVersion) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, b.BuildVersionCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", b.Command(), err)
	}
	if err := binary.Write(buf, o, b.Tools); err != nil {
		return fmt.Errorf("failed to write build tools to buffer: %v", err)
	}
	return nil
}
func (b *BuildVersion) String() string {
	if b.NumTools > 0 {
		if b.NumTools == 1 {
			return fmt.Sprintf("Platform: %s, SDK: %s, Tool: %s (%s)",
				b.Platform,
				b.Sdk,
				b.Tools[0].Tool,
				b.Tools[0].Version)
		} else {
			var tools []string
			for _, t := range b.Tools {
				tools = append(tools, fmt.Sprintf("%s (%s)", t.Tool, t.Version))
			}
			return fmt.Sprintf("Platform: %s, SDK: %s, Tools: [%s]",
				b.Platform,
				b.Sdk,
				strings.Join(tools, ", "))
		}
	}
	return fmt.Sprintf("Platform: %s, SDK: %s", b.Platform, b.Sdk)
}
func (b *BuildVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd  string                   `json:"load_cmd"`
		Len      uint32                   `json:"length"`
		Platform string                   `json:"platform"`
		Minos    string                   `json:"min_os"`
		Sdk      string                   `json:"sdk"`
		NumTools uint32                   `json:"num_tools"`
		Tools    []types.BuildVersionTool `json:"tools"`
	}{
		LoadCmd:  b.Command().String(),
		Len:      b.Len,
		Platform: b.Platform.String(),
		Minos:    b.Minos.String(),
		Sdk:      b.Sdk.String(),
		NumTools: b.NumTools,
		Tools:    b.Tools,
	})
}

/*******************************************************************************
 * LC_DYLD_EXPORTS_TRIE
 *******************************************************************************/

// A DyldExportsTrie used with linkedit_data_command, payload is trie
type DyldExportsTrie struct {
	LinkEditData
}

/*******************************************************************************
 * LC_DYLD_CHAINED_FIXUPS
 *******************************************************************************/

// A DyldChainedFixups used with linkedit_data_command
type DyldChainedFixups struct {
	LinkEditData
}

/*******************************************************************************
 * LC_FILESET_ENTRY
 *******************************************************************************/

// FilesetEntry used with fileset_entry_command
type FilesetEntry struct {
	LoadBytes
	types.FilesetEntryCmd
	EntryID string // contained entry id
}

func (l *FilesetEntry) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(l.FilesetEntryCmd)) + uint32(len(l.EntryID)) + 1)
}
func (l *FilesetEntry) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.FilesetEntryCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	if _, err := buf.WriteString(l.EntryID + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", l.EntryID, l.Command(), err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", l.Command(), err)
		}
	}
	return nil
}
func (f *FilesetEntry) String() string {
	return fmt.Sprintf("offset=0x%09x addr=0x%016x %s", f.FileOffset, f.Addr, f.EntryID)
}
func (l *FilesetEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Offset  uint64 `json:"offset"`
		Addr    uint64 `json:"address"`
		EntryID string `json:"entry_id"`
	}{
		LoadCmd: l.Command().String(),
		Len:     l.Len,
		Offset:  l.FileOffset,
		Addr:    l.Addr,
		EntryID: l.EntryID,
	})
}

/*******************************************************************************
 * LC_ATOM_INFO
 *******************************************************************************/

type AtomInfo struct {
	LinkEditData
}

/*******************************************************************************
 * COMMON COMMANDS
 *******************************************************************************/

// A Dylib represents a Mach-O LC_ID_DYLIB, LC_LOAD_{,WEAK_}DYLIB,LC_REEXPORT_DYLIB load command.
type Dylib struct {
	LoadBytes
	types.DylibCmd
	Name string
}

func (d *Dylib) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(d.DylibCmd)) + uint32(len(d.Name)) + 1)
}
func (d *Dylib) Put(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0*4:], uint32(d.LoadCmd))
	o.PutUint32(b[1*4:], d.Len)
	o.PutUint32(b[2*4:], d.NameOffset)
	o.PutUint32(b[3*4:], d.Timestamp)
	o.PutUint32(b[4*4:], uint32(d.CurrentVersion))
	o.PutUint32(b[5*4:], uint32(d.CompatVersion))
	return 6 * binary.Size(uint32(0))
}
func (d *Dylib) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, d.DylibCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", d.Command(), err)
	}
	if _, err := buf.WriteString(d.Name + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", d.Name, d.Command(), err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", d.Command(), err)
		}
	}
	return nil
}
func (d *Dylib) String() string {
	return fmt.Sprintf("%s (%s)", d.Name, d.CurrentVersion)
}
func (d *Dylib) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd   string `json:"load_cmd"`
		Len       uint32 `json:"length"`
		Name      string `json:"name"`
		Timestamp uint32 `json:"timestamp"`
		Current   string `json:"current_version"`
		Compat    string `json:"compatibility_version"`
	}{
		LoadCmd:   d.Command().String(),
		Len:       d.Len,
		Name:      d.Name,
		Timestamp: d.Timestamp,
		Current:   d.CurrentVersion.String(),
		Compat:    d.CompatVersion.String(),
	})
}

// A Dylinker represents a Mach-O LC_ID_DYLINKER, LC_LOAD_DYLINKER or LC_DYLD_ENVIRONMENT load command.
type Dylinker struct {
	LoadBytes
	types.DylinkerCmd
	Name string
}

func (d *Dylinker) LoadSize() uint32 {
	return pointerAlign(uint32(binary.Size(d.DylinkerCmd)) + uint32(len(d.Name)) + 1)
}
func (d *Dylinker) Put(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0*4:], uint32(d.LoadCmd))
	o.PutUint32(b[1*4:], d.Len)
	o.PutUint32(b[2*4:], d.NameOffset)
	return 3 * binary.Size(uint32(0))
}
func (d *Dylinker) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, d.DylinkerCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", d.Command(), err)
	}
	if _, err := buf.WriteString(d.Name + "\x00"); err != nil {
		return fmt.Errorf("failed to write %s to %s buffer: %v", d.Name, d.Command(), err)
	}
	if (buf.Len() % 8) != 0 {
		pad := 8 - (buf.Len() % 8)
		if _, err := buf.Write(make([]byte, pad)); err != nil {
			return fmt.Errorf("failed to write %s padding: %v", d.Command(), err)
		}
	}
	return nil
}
func (d *Dylinker) String() string {
	return d.Name
}
func (d *Dylinker) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Name    string `json:"name"`
	}{
		LoadCmd: d.Command().String(),
		Len:     d.Len,
		Name:    d.Name,
	})
}

// A VersionMin represents a Mach-O LC_VERSION_MIN_* command.
type VersionMin struct {
	LoadBytes
	types.VersionMinCmd
}

func (v *VersionMin) LoadSize() uint32 {
	return uint32(binary.Size(v.VersionMinCmd))
}
func (v *VersionMin) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, v.VersionMinCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", v.Command(), err)
	}
	return nil
}
func (v *VersionMin) String() string {
	return fmt.Sprintf("Version=%s, SDK=%s", v.Version, v.Sdk)
}
func (v *VersionMin) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Version string `json:"version"`
		SDK     string `json:"sdk"`
	}{
		LoadCmd: v.Command().String(),
		Len:     v.Len,
		Version: v.Version.String(),
		SDK:     v.Sdk.String(),
	})
}

// A LinkEditData represents a Mach-O linkedit data command.
/* LC_CODE_SIGNATURE, LC_SEGMENT_SPLIT_INFO, LC_FUNCTION_STARTS, LC_DATA_IN_CODE,
   LC_DYLIB_CODE_SIGN_DRS, LC_LINKER_OPTIMIZATION_HINT, LC_DYLD_EXPORTS_TRIE, or LC_DYLD_CHAINED_FIXUPS. */
type LinkEditData struct {
	LoadBytes
	types.LinkEditDataCmd
}

func (l *LinkEditData) LoadSize() uint32 {
	return uint32(binary.Size(l.LinkEditDataCmd))
}
func (l *LinkEditData) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, l.LinkEditDataCmd); err != nil {
		return fmt.Errorf("failed to write %s to buffer: %v", l.Command(), err)
	}
	return nil
}
func (l *LinkEditData) String() string {
	return fmt.Sprintf("offset=0x%09x  size=%#x", l.Offset, l.Size)
}
func (l *LinkEditData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LoadCmd string `json:"load_cmd"`
		Len     uint32 `json:"length"`
		Offset  uint32 `json:"offset"`
		Size    uint32 `json:"size"`
	}{
		LoadCmd: l.Command().String(),
		Len:     l.Len,
		Offset:  l.Offset,
		Size:    l.Size,
	})
}
