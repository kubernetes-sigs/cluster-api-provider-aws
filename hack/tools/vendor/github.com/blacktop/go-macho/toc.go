package macho

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/blacktop/go-macho/types"
)

type FileTOC struct {
	types.FileHeader
	ByteOrder binary.ByteOrder
	Loads     loads
	Sections  []*types.Section
	functions []types.Function
}

func (t *FileTOC) AddLoad(l Load) uint32 {
	loadsz := t.LoadSize()
	t.Loads = append(t.Loads, l)
	t.NCommands++
	t.SizeCommands += l.LoadSize()
	return t.LoadSize() - loadsz // delta
}

func (t *FileTOC) ModifySizeCommands(prev, curr int32) int32 {
	t.SizeCommands = uint32(int32(t.SizeCommands) + (curr - prev))
	return curr - prev
}

func (t *FileTOC) RemoveLoad(l Load) error {
	if len(t.Loads) == 0 {
		return fmt.Errorf("no loads to remove")
	}
	for i, load := range t.Loads {
		if load == l {
			t.Loads = append(t.Loads[:i], t.Loads[i+1:]...)
			t.NCommands--
			t.SizeCommands -= l.LoadSize()
			break
		}
	}
	return nil
}

// AddSegment adds segment s to the file table of contents,
// and also zeroes out the segment information with the expectation
// that this will be added next.
func (t *FileTOC) AddSegment(s *Segment) {
	s.Nsect = 0
	s.Firstsect = 0
	t.AddLoad(s)
}

// AddSection adds section to the most recently added Segment
func (t *FileTOC) AddSection(s *types.Section) {
	g := t.Loads[len(t.Loads)-1].(*Segment)
	if g.Nsect == 0 {
		g.Firstsect = uint32(len(t.Sections))
	}
	g.Nsect++
	t.Sections = append(t.Sections, s)
	sectionsize := uint32(unsafe.Sizeof(types.Section32{}))
	if g.Command() == types.LC_SEGMENT_64 {
		sectionsize = uint32(unsafe.Sizeof(types.Section64{}))
	}
	t.SizeCommands += sectionsize
	g.Len += sectionsize
}

// DerivedCopy returns a modified copy of the TOC, with empty loads and sections,
// and with the specified header type and flags.
func (t *FileTOC) DerivedCopy(Type types.HeaderFileType, Flags types.HeaderFlag) *FileTOC {
	h := t.FileHeader
	h.NCommands, h.SizeCommands, h.Type, h.Flags = 0, 0, Type, Flags

	return &FileTOC{FileHeader: h, ByteOrder: t.ByteOrder}
}

// TOCSize returns the size in bytes of the object file representation
// of the header and Load Commands (including Segments and Sections, but
// not their contents) at the beginning of a Mach-O file.  This typically
// overlaps the text segment in the object file.
func (t *FileTOC) TOCSize() uint32 {
	return t.HdrSize() + t.LoadSize()
}

// LoadAlign returns the required alignment of Load commands in a binary.
// This is used to add padding for necessary alignment.
func (t *FileTOC) LoadAlign() uint64 {
	if t.Magic == types.Magic64 {
		return 8
	}
	return 4
}

// HdrSize returns the size in bytes of the Macho header for a given
// magic number (where the magic number has been appropriately byte-swapped).
func (t *FileTOC) HdrSize() uint32 {
	switch t.Magic {
	case types.Magic32:
		return types.FileHeaderSize32
	case types.Magic64:
		return types.FileHeaderSize64
	case types.MagicFat:
		panic("MagicFat not handled yet")
	default:
		panic(fmt.Sprintf("Unexpected magic number %#x, expected Mach-O object file", t.Magic))
	}
}

// LoadSize returns the size of all the load commands in a file's table-of contents
// (but not their associated data, e.g., sections and symbol tables)
func (t *FileTOC) LoadSize() uint32 {
	cmdsz := uint32(0)
	for _, l := range t.Loads {
		s := l.LoadSize()
		cmdsz += s
	}
	return cmdsz
}

// FileSize returns the size in bytes of the header, load commands, and the
// in-file contents of all the segments and sections included in those
// load commands, accounting for their offsets within the file.
func (t *FileTOC) FileSize() uint64 {
	sz := uint64(t.LoadSize()) // ought to be contained in text segment, but just in case.
	for _, l := range t.Loads {
		if s, ok := l.(*Segment); ok {
			if m := s.Offset + s.Filesz; m > sz {
				sz = m
			}
		}
	}
	return sz
}

func (t *FileTOC) String() string {
	return fmt.Sprintf("%s\n%s\n", t.FileHeader.String(), t.Loads.String())
}

func (t *FileTOC) Print(printer func(t *FileTOC) string) string {
	return printer(t)
}

func (t *FileTOC) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Header types.FileHeader `json:"header"`
		Loads  loads            `json:"loads"`
	}{
		Header: t.FileHeader,
		Loads:  t.Loads,
	})
}

type loads []Load

// LoadsString returns a string representation of all the MachO's load commands
func (ls loads) String() string {
	var loadsStr string
	for i, l := range ls {
		if sg, ok := l.(*Segment); ok {
			loadsStr += fmt.Sprintf("%03d: %s\n", i, sg)
			for _, sc := range sg.sections {
				loadsStr += fmt.Sprintf("%s\n", sc)
			}
		} else {
			if l != nil {
				loadsStr += fmt.Sprintf("%03d: %-28s%s\n", i, l.Command(), l.String())
			}
		}
	}
	return loadsStr
}
