package macho

// High level access to low level data structures.

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blacktop/go-dwarf"

	"github.com/blacktop/go-macho/pkg/codesign"
	"github.com/blacktop/go-macho/pkg/fixupchains"
	"github.com/blacktop/go-macho/pkg/trie"
	"github.com/blacktop/go-macho/pkg/xar"
	"github.com/blacktop/go-macho/types"
)

// A File represents an open Mach-O file.
type File struct {
	FileTOC

	Symtab   *Symtab
	Dysymtab *Dysymtab

	vma         *types.VMAddrConverter
	dcf         *fixupchains.DyldChainedFixups
	exp         []trie.TrieExport
	exptrieData []byte
	binds       types.Binds
	objc        map[uint64]any
	ledata      *bytes.Buffer // tmp storage of linkedit data

	sharedCacheRelativeSelectorBaseVMAddress uint64 // objc_opt version 16

	mu     sync.Mutex
	sr     types.MachoReader
	cr     types.MachoReader
	closer io.Closer
}

/*
 * Mach-O reader
 */

var ErrMachOSectionNotFound = errors.New("MachO missing required section")
var ErrMachODyldInfoNotFound = errors.New("LC_DYLD_INFO(_ONLY) not found")

// FormatError is returned by some operations if the data does
// not have the correct format for an object file.
type FormatError struct {
	off int64
	msg string
	val any
}

func (e *FormatError) Error() string {
	msg := e.msg
	if e.val != nil {
		msg += fmt.Sprintf(" '%v'", e.val)
	}
	msg += fmt.Sprintf(" in record at byte %#x", e.off)
	return msg
}

func loadInSlice(c types.LoadCmd, list []types.LoadCmd) bool {
	for _, b := range list {
		if b == c {
			return true
		}
	}
	return false
}

// FileConfig is a MachO file config object
type FileConfig struct {
	Offset               int64
	LoadIncluding        []types.LoadCmd
	LoadExcluding        []types.LoadCmd
	VMAddrConverter      types.VMAddrConverter
	SectionReader        types.MachoReader
	CacheReader          types.MachoReader
	RelativeSelectorBase uint64
}

// Open opens the named file using os.Open and prepares it for use as a Mach-O binary.
func Open(name string) (*File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	ff, err := NewFile(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	ff.closer = f
	return ff, nil
}

// Close closes the File.
// If the File was created using NewFile directly instead of Open,
// Close has no effect.
func (f *File) Close() error {
	var err error
	if f.closer != nil {
		err = f.closer.Close()
		f.closer = nil
	}
	return err
}

// NewFile creates a new File for accessing a Mach-O binary in an underlying reader.
// The Mach-O binary is expected to start at position 0 in the ReaderAt.
func NewFile(r io.ReaderAt, config ...FileConfig) (*File, error) {
	var loadIncluding []types.LoadCmd
	var loadExcluding []types.LoadCmd

	f := new(File)

	f.objc = make(map[uint64]any)

	if config != nil {
		if config[0].SectionReader != nil {
			f.sr = config[0].SectionReader
			f.sr.Seek(config[0].Offset, io.SeekStart)
			f.cr = f.sr
		}
		if config[0].CacheReader != nil {
			f.cr = config[0].CacheReader
		}
		f.vma = &config[0].VMAddrConverter
		loadIncluding = config[0].LoadIncluding
		loadExcluding = config[0].LoadExcluding
		f.sharedCacheRelativeSelectorBaseVMAddress = config[0].RelativeSelectorBase
	} else {
		f.vma = &types.VMAddrConverter{
			Converter:    f.convertToVMAddr,
			VMAddr2Offet: f.getOffset,
			Offet2VMAddr: f.getVMAddress,
		}
		f.sr = types.NewCustomSectionReader(r, f.vma, 0, 1<<63-1)
		f.cr = f.sr
	}

	// Read and decode Mach magic to determine byte order, size.
	// Magic32 and Magic64 differ only in the bottom bit.
	var ident [4]byte
	if _, err := r.ReadAt(ident[0:], 0); err != nil {
		return nil, fmt.Errorf("failed to parse magic: %v", err)
	}
	be := binary.BigEndian.Uint32(ident[0:])
	le := binary.LittleEndian.Uint32(ident[0:])
	switch types.Magic32.Int() &^ 1 {
	case be &^ 1:
		f.ByteOrder = binary.BigEndian
		f.Magic = types.Magic(be)
	case le &^ 1:
		f.ByteOrder = binary.LittleEndian
		f.Magic = types.Magic(le)
	default:
		return nil, &FormatError{0, "invalid magic number", nil}
	}

	// Read entire file header.
	if err := binary.Read(f.sr, f.ByteOrder, &f.FileHeader); err != nil {
		return nil, fmt.Errorf("failed to parse header: %v", err)
	}

	// Then load commands.
	offset := int64(types.FileHeaderSize32)
	if f.Magic == types.Magic64 {
		offset = types.FileHeaderSize64
	}
	dat := make([]byte, f.SizeCommands)
	if _, err := r.ReadAt(dat, offset); err != nil {
		return nil, fmt.Errorf("failed to parse command dat: %v", err)
	}
	f.Loads = []Load{}
	bo := f.ByteOrder
	for i := uint32(0); i < f.NCommands; i++ {
		// Each load command begins with uint32 command and length.
		if len(dat) < 8 {
			return nil, &FormatError{offset, "command block too small", nil}
		}
		cmd, siz := types.LoadCmd(bo.Uint32(dat[0:4])), bo.Uint32(dat[4:8])
		if siz < 8 || siz > uint32(len(dat)) {
			return nil, &FormatError{offset, "invalid command block size", nil}
		}

		var cmddat []byte
		cmddat, dat = dat[0:siz], dat[siz:]
		offset += int64(siz)
		var s *Segment

		// skip unwanted load commands
		if len(loadIncluding) > 0 && !loadInSlice(cmd, loadIncluding) {
			continue
		} else if loadInSlice(cmd, loadExcluding) {
			continue
		}

		switch cmd {
		default:
			log.Printf("found NEW load command: %s (please let the author know via https://github.com/blacktop/go-macho/issues)", cmd)
			f.Loads = append(f.Loads, LoadCmdBytes{types.LoadCmd(cmd), LoadBytes(cmddat)})
		case types.LC_SEGMENT:
			var seg32 types.Segment32
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &seg32); err != nil {
				return nil, fmt.Errorf("failed to read LC_SEGMENT: %v", err)
			}
			s = new(Segment)
			s.LoadBytes = cmddat
			s.LoadCmd = cmd
			s.Len = siz
			s.Name = cstring(seg32.Name[0:])
			s.Addr = uint64(seg32.Addr)
			s.Memsz = uint64(seg32.Memsz)
			s.Offset = uint64(seg32.Offset)
			s.Filesz = uint64(seg32.Filesz)
			s.Maxprot = seg32.Maxprot
			s.Prot = seg32.Prot
			s.Nsect = seg32.Nsect
			s.Flag = seg32.Flag
			s.Firstsect = uint32(len(f.Sections))
			for i := 0; i < int(s.Nsect); i++ {
				var sh32 types.Section32
				if err := binary.Read(b, bo, &sh32); err != nil {
					return nil, fmt.Errorf("failed to read Section32: %v", err)
				}
				sh := new(types.Section)
				sh.Type = 32
				sh.Name = cstring(sh32.Name[0:])
				sh.Seg = cstring(sh32.Seg[0:])
				sh.Addr = uint64(sh32.Addr)
				sh.Size = uint64(sh32.Size)
				sh.Offset = sh32.Offset
				sh.Align = sh32.Align
				sh.Reloff = sh32.Reloff
				sh.Nreloc = sh32.Nreloc
				sh.Flags = sh32.Flags
				sh.Reserved1 = sh32.Reserve1
				sh.Reserved2 = sh32.Reserve2
				sh.SetReaders(f.cr, io.NewSectionReader(f.cr, int64(sh32.Offset), int64(sh32.Size)))
				if err := f.pushSection(sh, f.cr); err != nil {
					return nil, fmt.Errorf("failed to pushSection32: %v", err)
				}
				s.sections = append(s.sections, sh)
			}
			f.Loads = append(f.Loads, s)
		case types.LC_SEGMENT_64:
			var seg64 types.Segment64
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &seg64); err != nil {
				return nil, fmt.Errorf("failed to read LC_SEGMENT_64: %v", err)
			}
			s = new(Segment)
			s.LoadBytes = cmddat
			s.LoadCmd = cmd
			s.Len = siz
			s.Name = cstring(seg64.Name[0:])
			s.Addr = seg64.Addr
			s.Memsz = seg64.Memsz
			s.Offset = seg64.Offset
			s.Filesz = seg64.Filesz
			s.Maxprot = seg64.Maxprot
			s.Prot = seg64.Prot
			s.Nsect = seg64.Nsect
			s.Flag = seg64.Flag
			s.Firstsect = uint32(len(f.Sections))
			for i := 0; i < int(s.Nsect); i++ {
				var sh64 types.Section64
				if err := binary.Read(b, bo, &sh64); err != nil {
					return nil, fmt.Errorf("failed to read Section64: %v", err)
				}
				sh := new(types.Section)
				sh.Type = 64
				sh.Name = cstring(sh64.Name[0:])
				sh.Seg = cstring(sh64.Seg[0:])
				sh.Addr = sh64.Addr
				sh.Size = sh64.Size
				sh.Offset = sh64.Offset
				sh.Align = sh64.Align
				sh.Reloff = sh64.Reloff
				sh.Nreloc = sh64.Nreloc
				sh.Flags = sh64.Flags
				sh.Reserved1 = sh64.Reserve1
				sh.Reserved2 = sh64.Reserve2
				sh.Reserved3 = sh64.Reserve3
				sh.SetReaders(f.cr, io.NewSectionReader(f.cr, int64(sh64.Offset), int64(sh64.Size)))
				if err := f.pushSection(sh, f.cr); err != nil {
					return nil, fmt.Errorf("failed to pushSection64: %v", err)
				}
				s.sections = append(s.sections, sh)
			}
			f.Loads = append(f.Loads, s)
		case types.LC_SYMTAB:
			var hdr types.SymtabCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_SYMTAB: %v", err)
			}
			symdat := make([]byte, int(hdr.Nsyms)*f.symbolSize())
			if _, err := f.cr.ReadAt(symdat, int64(hdr.Symoff)); err != nil {
				return nil, fmt.Errorf("failed to read data at Symoff=%#x; %v", int64(hdr.Symoff), err)
			}
			strtab := make([]byte, hdr.Strsize)
			if _, err := f.cr.ReadAt(strtab, int64(hdr.Stroff)); err != nil {
				return nil, fmt.Errorf("failed to read data at Stroff=%#x; %v", int64(hdr.Stroff), err)
			}
			st, err := f.parseSymtab(symdat, strtab, cmddat, &hdr, offset)
			if err != nil {
				return nil, fmt.Errorf("failed to read parseSymtab: %v", err)
			}
			st.LoadBytes = cmddat
			st.LoadCmd = cmd
			st.Len = siz
			f.Loads = append(f.Loads, st)
			f.Symtab = st
		case types.LC_SYMSEG:
			var led types.SymsegCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &led); err != nil {
				return nil, fmt.Errorf("failed to read LC_SYMSEG: %v", err)
			}

			l := new(SymSeg)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = led.Offset
			l.Size = led.Size
			f.Loads = append(f.Loads, l)
		case types.LC_THREAD:
			var t types.ThreadCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &t); err != nil {
				return nil, fmt.Errorf("failed to read LC_THREAD: %v", err)
			}
			l := new(Thread)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			for {
				var thread types.ThreadState
				err := binary.Read(b, bo, &thread.Flavor)
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, fmt.Errorf("failed to read LC_THREAD flavor: %v", err)
				}
				if err := binary.Read(b, bo, &thread.Count); err != nil {
					return nil, fmt.Errorf("failed to read LC_THREAD count: %v", err)
				}
				thread.Data = make([]byte, thread.Count*uint32(binary.Size(uint32(0))))
				if err := binary.Read(b, bo, &thread.Data); err != nil {
					return nil, fmt.Errorf("failed to read LC_THREAD state struct data: %v", err)
				}
				l.Threads = append(l.Threads, thread)
			}
			f.Loads = append(f.Loads, l)
		case types.LC_UNIXTHREAD:
			var ut types.UnixThreadCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &ut); err != nil {
				return nil, fmt.Errorf("failed to read LC_UNIXTHREAD: %v", err)
			}
			l := new(UnixThread)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.bo = bo
			for {
				var thread types.ThreadState
				err := binary.Read(b, bo, &thread.Flavor)
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, fmt.Errorf("failed to read LC_UNIXTHREAD flavor: %v", err)
				}
				if err := binary.Read(b, bo, &thread.Count); err != nil {
					return nil, fmt.Errorf("failed to read LC_UNIXTHREAD count: %v", err)
				}
				thread.Data = make([]byte, thread.Count*uint32(binary.Size(uint32(0))))
				if err := binary.Read(b, bo, &thread.Data); err != nil {
					return nil, fmt.Errorf("failed to read LC_UNIXTHREAD state struct data: %v", err)
				}
				l.Threads = append(l.Threads, thread)
			}
			f.Loads = append(f.Loads, l)
		case types.LC_LOADFVMLIB:
			var hdr types.LoadFvmLibCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_LOADFVMLIB: %v", err)
			}
			l := new(LoadFvmlib)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			l.MinorVersion = hdr.MinorVersion
			l.HeaderAddr = hdr.HeaderAddr
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in LC_LOADFVMLIB command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_IDFVMLIB:
			var hdr types.IDFvmLibCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_IDFVMLIB: %v", err)
			}
			l := new(IDFvmlib)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			l.MinorVersion = hdr.MinorVersion
			l.HeaderAddr = hdr.HeaderAddr
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in LC_IDFVMLIB command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_IDENT:
			var hdr types.IdentCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_IDENT: %v", err)
			}
			l := new(Ident)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			br := bufio.NewReader(b)
			for {
				o, err := br.ReadString('\x00')
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, fmt.Errorf("failed to read LC_IDENT options: %v", err)
				}
				l.StrTable = append(l.StrTable, o)
			}
			f.Loads = append(f.Loads, l)
		case types.LC_FVMFILE:
			var hdr types.FvmFileCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_FVMFILE: %v", err)
			}
			l := new(FvmFile)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			l.HeaderAddr = hdr.HeaderAddr
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in LC_FVMFILE command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_PREPAGE:
			var hdr types.PrePageCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_PREPAGE: %v", err)
			}
			l := new(Prepage)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			f.Loads = append(f.Loads, l)
		case types.LC_DYSYMTAB:
			var hdr types.DysymtabCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_DYSYMTAB: %v", err)
			}
			if f.Symtab != nil && hdr.Iundefsym > uint32(len(f.Symtab.Syms)) {
				return nil, &FormatError{offset, fmt.Sprintf(
					"undefined symbols index in dynamic symbol table command is greater than symbol table length (%d > %d)",
					hdr.Iundefsym, len(f.Symtab.Syms)), nil}
			} else if f.Symtab != nil && hdr.Iundefsym+hdr.Nundefsym > uint32(len(f.Symtab.Syms)) {
				return nil, &FormatError{offset, fmt.Sprintf(
					"number of undefined symbols after index in dynamic symbol table command is greater than symbol table length (%d > %d)",
					hdr.Iundefsym+hdr.Nundefsym, len(f.Symtab.Syms)), nil}
			}
			dat := make([]byte, hdr.Nindirectsyms*4)
			if _, err := f.cr.ReadAt(dat, int64(hdr.Indirectsymoff)); err != nil {
				return nil, fmt.Errorf("failed to read data at Indirectsymoff=%#x; %v", int64(hdr.Indirectsymoff), err)
			}
			x := make([]uint32, hdr.Nindirectsyms)
			if err := binary.Read(bytes.NewReader(dat), bo, x); err != nil {
				return nil, fmt.Errorf("failed to read Nindirectsyms: %v", err)
			}
			// TODO: parse DylibTableOfContents if Ntoc > 0
			// TODO: parse DylibModule if Nmodtab > 0
			// TODO: parse DylibReference if Nextrefsyms > 0
			// TODO: parse RelocInfo if Nlocrel > 0
			st := new(Dysymtab)
			st.LoadBytes = cmddat
			st.LoadCmd = cmd
			st.Len = siz
			st.DysymtabCmd = hdr
			st.IndirectSyms = x
			f.Loads = append(f.Loads, st)
			f.Dysymtab = st
		case types.LC_LOAD_DYLIB:
			var hdr types.DylibCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_LOAD_DYLIB: %v", err)
			}
			l := new(LoadDylib)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in dynamic library command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			l.Timestamp = hdr.Timestamp
			l.CurrentVersion = hdr.CurrentVersion
			l.CompatVersion = hdr.CompatVersion
			f.Loads = append(f.Loads, l)
		case types.LC_ID_DYLIB:
			var hdr types.DylibCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_ID_DYLIB: %v", err)
			}
			l := new(IDDylib)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in dynamic library ident command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			l.Timestamp = hdr.Timestamp
			l.CurrentVersion = hdr.CurrentVersion
			l.CompatVersion = hdr.CompatVersion
			f.Loads = append(f.Loads, l)
		case types.LC_LOAD_DYLINKER:
			var hdr types.DylinkerCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_LOAD_DYLINKER: %v", err)
			}
			l := new(LoadDylinker)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in load dylinker command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_ID_DYLINKER:
			var hdr types.IDDylinkerCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_ID_DYLINKER: %v", err)
			}
			l := new(DylinkerID)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in load dylinker command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_PREBOUND_DYLIB:
			var hdr types.PreboundDylibCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_PREBOUND_DYLIB: %v", err)
			}
			l := new(PreboundDylib)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in LC_PREBOUND_DYLIB command", hdr.NameOffset}
			}
			l.NumModules = hdr.NumModules
			l.Name = cstring(cmddat[hdr.NameOffset:])
			if hdr.LinkedModulesOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid linked modules in LC_PREBOUND_DYLIB command", hdr.NameOffset}
			}
			l.LinkedModulesBitVector = cstring(cmddat[hdr.LinkedModulesOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_ROUTINES:
			var rt types.RoutinesCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &rt); err != nil {
				return nil, fmt.Errorf("failed to read LC_ROUTINES: %v", err)
			}
			l := new(Routines)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.InitAddress = rt.InitAddress
			l.InitModule = rt.InitModule
			f.Loads = append(f.Loads, l)
		case types.LC_SUB_FRAMEWORK:
			var sf types.SubFrameworkCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &sf); err != nil {
				return nil, fmt.Errorf("failed to read LC_SUB_FRAMEWORK: %v", err)
			}
			l := new(SubFramework)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.FrameworkOffset = sf.FrameworkOffset
			if sf.FrameworkOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid framework in sub-framework command", sf.FrameworkOffset}
			}
			l.Framework = cstring(cmddat[sf.FrameworkOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_SUB_UMBRELLA:
			var su types.SubUmbrellaCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &su); err != nil {
				return nil, fmt.Errorf("failed to read LC_SUB_UMBRELLA: %v", err)
			}
			l := new(SubUmbrella)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.UmbrellaOffset = su.UmbrellaOffset
			if su.UmbrellaOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid framework in sub-umbrella command", su.UmbrellaOffset}
			}
			l.Umbrella = cstring(cmddat[su.UmbrellaOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_SUB_CLIENT:
			var sc types.SubClientCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &sc); err != nil {
				return nil, fmt.Errorf("failed to read LC_SUB_CLIENT: %v", err)
			}
			l := new(SubClient)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.ClientOffset = sc.ClientOffset
			if sc.ClientOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid path in sub client command", sc.ClientOffset}
			}
			l.Name = cstring(cmddat[sc.ClientOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_SUB_LIBRARY:
			var s types.SubLibraryCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &s); err != nil {
				return nil, fmt.Errorf("failed to read LC_SUB_LIBRARY: %v", err)
			}
			l := new(SubLibrary)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.LibraryOffset = s.LibraryOffset
			if s.LibraryOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid framework in sub-library command", s.LibraryOffset}
			}
			l.Library = cstring(cmddat[s.LibraryOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_TWOLEVEL_HINTS:
			var t types.TwolevelHintsCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &t); err != nil {
				return nil, fmt.Errorf("failed to read LC_TWOLEVEL_HINTS: %v", err)
			}
			l := new(TwolevelHints)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = t.Offset
			l.Hints = make([]types.TwolevelHint, t.NumHints)
			if err := binary.Read(b, bo, &l.Hints); err != nil {
				return nil, fmt.Errorf("failed to read hints data: %v", err)
			}
			f.Loads = append(f.Loads, l)

		case types.LC_PREBIND_CKSUM:
			var p types.PrebindCksumCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &p); err != nil {
				return nil, fmt.Errorf("failed to read LC_PREBIND_CKSUM: %v", err)
			}
			l := new(PrebindCheckSum)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.CheckSum = p.CheckSum
			f.Loads = append(f.Loads, l)
		case types.LC_LOAD_WEAK_DYLIB:
			var hdr types.DylibCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_LOAD_WEAK_DYLIB: %v", err)
			}
			l := new(WeakDylib)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in weak dynamic library command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			l.Timestamp = hdr.Timestamp
			l.CurrentVersion = hdr.CurrentVersion
			l.CompatVersion = hdr.CompatVersion
			f.Loads = append(f.Loads, l)
		case types.LC_ROUTINES_64:
			var r64 types.Routines64Cmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &r64); err != nil {
				return nil, fmt.Errorf("failed to read LC_ROUTINES_64: %v", err)
			}
			l := new(Routines64)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.InitAddress = r64.InitAddress
			l.InitModule = r64.InitModule
			f.Loads = append(f.Loads, l)
		case types.LC_UUID:
			var u types.UUIDCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &u); err != nil {
				return nil, fmt.Errorf("failed to read LC_UUID: %v", err)
			}
			l := new(UUID)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.UUID = u.UUID
			f.Loads = append(f.Loads, l)
		case types.LC_RPATH:
			var hdr types.RpathCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_RPATH: %v", err)
			}
			l := new(Rpath)
			if hdr.PathOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid path in rpath command", hdr.PathOffset}
			}
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.PathOffset = hdr.PathOffset
			if hdr.PathOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid path in rpath command", hdr.PathOffset}
			}
			l.Path = cstring(cmddat[hdr.PathOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_CODE_SIGNATURE:
			var hdr types.CodeSignatureCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_CODE_SIGNATURE: %v", err)
			}

			l := new(CodeSignature)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = hdr.Offset
			l.Size = hdr.Size
			csdat := make([]byte, hdr.Size)
			if _, err := f.cr.ReadAt(csdat, int64(hdr.Offset)); err != nil {
				return nil, fmt.Errorf("failed to read CS data at offset=%#x; %v", int64(hdr.Offset), err)
			}
			cs, err := codesign.ParseCodeSignature(csdat)
			if err != nil {
				return nil, fmt.Errorf("failed to ParseCodeSignature: %v", err)
			}
			l.CodeSignature = *cs
			f.Loads = append(f.Loads, l)
		case types.LC_SEGMENT_SPLIT_INFO:
			var hdr types.SegmentSplitInfoCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO: %v", err)
			}
			l := new(SplitInfo)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = hdr.Offset
			l.Size = hdr.Size
			if l.Size > 0 {
				ldat := make([]byte, l.Size)
				if _, err := f.cr.ReadAt(ldat, int64(l.Offset)); err != nil {
					return nil, fmt.Errorf("failed to read SplitInfo data at offset=%#x; %v", int64(hdr.Offset), err)
				}
				fsr := bytes.NewReader(ldat)
				if err := binary.Read(fsr, bo, &l.Version); err != nil {
					return nil, fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO Version: %v", err)
				}
			}
			f.Loads = append(f.Loads, l)
		case types.LC_REEXPORT_DYLIB:
			var hdr types.ReExportDylibCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_REEXPORT_DYLIB: %v", err)
			}
			l := new(ReExportDylib)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in dynamic library command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			l.Timestamp = hdr.Timestamp
			l.CurrentVersion = hdr.CurrentVersion
			l.CompatVersion = hdr.CompatVersion
			f.Loads = append(f.Loads, l)
		case types.LC_LAZY_LOAD_DYLIB:
			var hdr types.LazyLoadDylibCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_LAZY_LOAD_DYLIB: %v", err)
			}
			l := new(LazyLoadDylib)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in load upwardl dylib command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			l.Timestamp = hdr.Timestamp
			l.CurrentVersion = hdr.CurrentVersion
			l.CompatVersion = hdr.CompatVersion
			f.Loads = append(f.Loads, l)
		case types.LC_ENCRYPTION_INFO:
			var ei types.EncryptionInfoCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &ei); err != nil {
				return nil, fmt.Errorf("failed to read LC_ENCRYPTION_INFO: %v", err)
			}

			l := new(EncryptionInfo)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = ei.Offset
			l.Size = ei.Size
			l.CryptID = ei.CryptID
			f.Loads = append(f.Loads, l)
		case types.LC_DYLD_INFO:
			var info types.DyldInfoCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &info); err != nil {
				return nil, fmt.Errorf("failed to read LC_DYLD_INFO: %v", err)
			}
			l := new(DyldInfo)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.RebaseOff = info.RebaseOff
			l.RebaseSize = info.RebaseSize
			l.BindOff = info.BindOff
			l.BindSize = info.BindSize
			l.WeakBindOff = info.WeakBindOff
			l.WeakBindSize = info.WeakBindSize
			l.LazyBindOff = info.LazyBindOff
			l.LazyBindSize = info.LazyBindSize
			l.ExportOff = info.ExportOff
			l.ExportSize = info.ExportSize
			f.Loads = append(f.Loads, l)
		case types.LC_DYLD_INFO_ONLY:
			var info types.DyldInfoOnlyCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &info); err != nil {
				return nil, fmt.Errorf("failed to read LC_DYLD_INFO_ONLY: %v", err)
			}
			l := new(DyldInfoOnly)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.RebaseOff = info.RebaseOff
			l.RebaseSize = info.RebaseSize
			l.BindOff = info.BindOff
			l.BindSize = info.BindSize
			l.WeakBindOff = info.WeakBindOff
			l.WeakBindSize = info.WeakBindSize
			l.LazyBindOff = info.LazyBindOff
			l.LazyBindSize = info.LazyBindSize
			l.ExportOff = info.ExportOff
			l.ExportSize = info.ExportSize
			f.Loads = append(f.Loads, l)
		case types.LC_LOAD_UPWARD_DYLIB:
			var hdr types.LoadUpwardDylibCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_LOAD_UPWARD_DYLIB: %v", err)
			}
			l := new(UpwardDylib)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in load upwardl dylib command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			l.Timestamp = hdr.Timestamp
			l.CurrentVersion = hdr.CurrentVersion
			l.CompatVersion = hdr.CompatVersion
			f.Loads = append(f.Loads, l)
		case types.LC_VERSION_MIN_MACOSX:
			var verMin types.VersionMinMacOSCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &verMin); err != nil {
				return nil, fmt.Errorf("failed to read LC_VERSION_MIN_MACOSX: %v", err)
			}
			l := new(VersionMinMacOSX)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Version = verMin.Version
			l.Sdk = verMin.Sdk
			f.Loads = append(f.Loads, l)
		case types.LC_VERSION_MIN_IPHONEOS:
			var verMin types.VersionMinIPhoneOSCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &verMin); err != nil {
				return nil, fmt.Errorf("failed to read LC_VERSION_MIN_IPHONEOS: %v", err)
			}
			l := new(VersionMiniPhoneOS)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Version = verMin.Version
			l.Sdk = verMin.Sdk
			f.Loads = append(f.Loads, l)
		case types.LC_FUNCTION_STARTS:
			var led types.LinkEditDataCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &led); err != nil {
				return nil, fmt.Errorf("failed to read LC_FUNCTION_STARTS: %v", err)
			}

			l := new(FunctionStarts)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = led.Offset
			l.Size = led.Size
			f.Loads = append(f.Loads, l)
		case types.LC_DYLD_ENVIRONMENT:
			var hdr types.DyldEnvironmentCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_DYLD_ENVIRONMENT: %v", err)
			}
			l := new(DyldEnvironment)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.NameOffset = hdr.NameOffset
			if hdr.NameOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in dyld environment command", hdr.NameOffset}
			}
			l.Name = cstring(cmddat[hdr.NameOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_MAIN:
			var hdr types.EntryPointCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_MAIN: %v", err)
			}
			l := new(EntryPoint)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.EntryOffset = hdr.EntryOffset
			l.StackSize = hdr.StackSize
			f.Loads = append(f.Loads, l)
		case types.LC_DATA_IN_CODE:
			var led types.LinkEditDataCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &led); err != nil {
				return nil, fmt.Errorf("failed to read LC_DATA_IN_CODE: %v", err)
			}
			l := new(DataInCode)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = led.Offset
			l.Size = led.Size
			ldat := make([]byte, l.Size)
			if _, err := f.cr.ReadAt(ldat, int64(l.Offset)); err != nil {
				return nil, fmt.Errorf("failed to read DataInCode data at offset=%#x; %v", int64(led.Offset), err)
			}
			l.Entries = make([]types.DataInCodeEntry, len(ldat)/binary.Size(types.DataInCodeEntry{}))
			if err := binary.Read(bytes.NewReader(ldat), bo, &l.Entries); err != nil {
				return nil, fmt.Errorf("failed to read LC_DATA_IN_CODE entries: %v", err)
			}
			f.Loads = append(f.Loads, l)
		case types.LC_SOURCE_VERSION:
			var sv types.SourceVersionCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &sv); err != nil {
				return nil, fmt.Errorf("failed to read LC_SOURCE_VERSION: %v", err)
			}
			l := new(SourceVersion)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Version = sv.Version
			f.Loads = append(f.Loads, l)
		case types.LC_DYLIB_CODE_SIGN_DRS:
			var led types.LinkEditDataCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &led); err != nil {
				return nil, fmt.Errorf("failed to read LC_DYLIB_CODE_SIGN_DRS: %v", err)
			}

			l := new(DylibCodeSignDrs)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = led.Offset
			l.Size = led.Size
			f.Loads = append(f.Loads, l)
		case types.LC_ENCRYPTION_INFO_64:
			var ei types.EncryptionInfo64Cmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &ei); err != nil {
				return nil, fmt.Errorf("failed to read LC_ENCRYPTION_INFO_64: %v", err)
			}
			l := new(EncryptionInfo64)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = ei.Offset
			l.Size = ei.Size
			l.CryptID = ei.CryptID
			f.Loads = append(f.Loads, l)
		case types.LC_LINKER_OPTION:
			var lo types.LinkerOptionCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &lo); err != nil {
				return nil, fmt.Errorf("failed to read LC_LINKER_OPTION: %v", err)
			}
			l := new(LinkerOption)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			for i := 0; i < int(lo.Count); i++ {
				o, err := bufio.NewReader(b).ReadString('\x00')
				if err != nil {
					break // FIXME: should this error?
				}
				l.Options = append(l.Options, o)
			}
			f.Loads = append(f.Loads, l)
		case types.LC_LINKER_OPTIMIZATION_HINT:
			var led types.LinkEditDataCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &led); err != nil {
				return nil, fmt.Errorf("failed to read LC_LINKER_OPTIMIZATION_HINT: %v", err)
			}

			l := new(LinkerOptimizationHint)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = led.Offset
			l.Size = led.Size
			f.Loads = append(f.Loads, l)
		case types.LC_VERSION_MIN_TVOS:
			var verMin types.VersionMinMacOSCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &verMin); err != nil {
				return nil, fmt.Errorf("failed to read LC_VERSION_MIN_TVOS: %v", err)
			}
			l := new(VersionMinTvOS)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Version = verMin.Version
			l.Sdk = verMin.Sdk
			f.Loads = append(f.Loads, l)
		case types.LC_VERSION_MIN_WATCHOS:
			var verMin types.VersionMinWatchOSCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &verMin); err != nil {
				return nil, fmt.Errorf("failed to read LC_VERSION_MIN_WATCHOS: %v", err)
			}
			l := new(VersionMinWatchOS)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Version = verMin.Version
			l.Sdk = verMin.Sdk
			f.Loads = append(f.Loads, l)
		case types.LC_NOTE:
			var n types.NoteCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &n); err != nil {
				return nil, fmt.Errorf("failed to read LC_NOTE: %v", err)
			}
			l := new(Note)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.DataOwner = n.DataOwner
			l.Offset = n.Offset
			l.Size = n.Size
			f.Loads = append(f.Loads, l)
		case types.LC_BUILD_VERSION:
			var build types.BuildVersionCmd
			var buildTool types.BuildVersionTool
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &build); err != nil {
				return nil, fmt.Errorf("failed to read LC_BUILD_VERSION: %v", err)
			}
			l := new(BuildVersion)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Platform = build.Platform
			l.Minos = build.Minos
			l.Sdk = build.Sdk
			l.NumTools = build.NumTools
			for i := uint32(0); i < build.NumTools; i++ {
				if err := binary.Read(b, bo, &buildTool); err != nil {
					return nil, fmt.Errorf("failed to read LC_BUILD_VERSION buildTool: %v", err)
				}
				l.Tools = append(l.Tools, types.BuildVersionTool{
					Tool:    buildTool.Tool,
					Version: buildTool.Version,
				})
			}
			f.Loads = append(f.Loads, l)
		case types.LC_DYLD_EXPORTS_TRIE:
			var led types.LinkEditDataCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &led); err != nil {
				return nil, fmt.Errorf("failed to read LC_DYLD_EXPORTS_TRIE: %v", err)
			}

			l := new(DyldExportsTrie)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = led.Offset
			l.Size = led.Size
			f.Loads = append(f.Loads, l)
		case types.LC_DYLD_CHAINED_FIXUPS:
			var led types.DyldChainedFixupsCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &led); err != nil {
				return nil, fmt.Errorf("failed to read LC_DYLD_CHAINED_FIXUPS: %v", err)
			}

			l := new(DyldChainedFixups)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = led.Offset
			l.Size = led.Size
			f.Loads = append(f.Loads, l)
		case types.LC_FILESET_ENTRY:
			var hdr types.FilesetEntryCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &hdr); err != nil {
				return nil, fmt.Errorf("failed to read LC_FILESET_ENTRY: %v", err)
			}
			l := new(FilesetEntry)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Addr = hdr.Addr
			l.FileOffset = hdr.FileOffset
			l.EntryIdOffset = hdr.EntryIdOffset
			if hdr.EntryIdOffset >= uint32(len(cmddat)) {
				return nil, &FormatError{offset, "invalid name in load fileset entry command", hdr.EntryIdOffset}
			}
			l.EntryID = cstring(cmddat[hdr.EntryIdOffset:])
			f.Loads = append(f.Loads, l)
		case types.LC_ATOM_INFO:
			var led types.LinkEditDataCmd
			b := bytes.NewReader(cmddat)
			if err := binary.Read(b, bo, &led); err != nil {
				return nil, fmt.Errorf("failed to read LC_ATOM_INFO: %v", err)
			}
			l := new(AtomInfo)
			l.LoadBytes = cmddat
			l.LoadCmd = cmd
			l.Len = siz
			l.Offset = led.Offset
			l.Size = led.Size
			f.Loads = append(f.Loads, l)
		}
		if s != nil {
			// s.sr = io.NewSectionReader(r, int64(s.Offset), int64(s.Filesz))
			s.ReaderAt = f.sr
		}
	}
	return f, nil
}

func (f *File) parseSymtab(symdat, strtab, cmddat []byte, hdr *types.SymtabCmd, offset int64) (*Symtab, error) {
	bo := f.ByteOrder
	symtab := make([]Symbol, hdr.Nsyms)
	b := bytes.NewReader(symdat)
	for i := range symtab {
		var n types.Nlist64
		if f.Magic == types.Magic64 {
			if err := binary.Read(b, bo, &n); err != nil {
				return nil, fmt.Errorf("failed to read Symtab magic: %v", err)
			}
		} else {
			var n32 types.Nlist32
			if err := binary.Read(b, bo, &n32); err != nil {
				return nil, fmt.Errorf("failed to read Symtab nlist32: %v", err)
			}
			n.Name = n32.Name
			n.Type = n32.Type
			n.Sect = n32.Sect
			n.Desc = n32.Desc
			n.Value = uint64(n32.Value)
		}
		sym := &symtab[i]
		if n.Name >= uint32(len(strtab)) {
			return nil, &FormatError{offset, "invalid name in symbol table", n.Name}
		}
		// We add "_" to Go symbols. Strip it here. See issue 33808.
		name := cstring(strtab[n.Name:])
		if strings.Contains(name, ".") && name[0] == '_' {
			name = name[1:]
		}
		sym.Name = name
		sym.Type = n.Type
		sym.Sect = n.Sect
		sym.Desc = n.Desc
		sym.Value = n.Value
	}
	st := new(Symtab)
	st.LoadBytes = LoadBytes(cmddat)
	st.Symoff = hdr.Symoff
	st.Nsyms = hdr.Nsyms
	st.Stroff = hdr.Stroff
	st.Strsize = hdr.Strsize
	st.Len = hdr.Len
	st.Syms = symtab
	return st, nil
}

func (f *File) pushSection(sh *types.Section, r io.ReaderAt) error {
	f.Sections = append(f.Sections, sh)

	if sh.Nreloc > 0 {
		reldat := make([]byte, int(sh.Nreloc)*8)
		if _, err := r.ReadAt(reldat, int64(sh.Reloff)); err != nil {
			return fmt.Errorf("failed to read data at Reloff=%#x; %v", int64(sh.Reloff), err)
		}
		b := bytes.NewReader(reldat)

		bo := f.ByteOrder

		sh.Relocs = make([]types.Reloc, sh.Nreloc)
		for i := range sh.Relocs {
			rel := &sh.Relocs[i]

			var ri types.RelocInfo
			if err := binary.Read(b, bo, &ri); err != nil {
				return fmt.Errorf("failed to read relocInfo; %v", err)
			}

			if ri.Addr&(1<<31) != 0 { // scattered
				rel.Addr = ri.Addr & (1<<24 - 1)
				rel.Type = uint8((ri.Addr >> 24) & (1<<4 - 1))
				rel.Len = uint8((ri.Addr >> 28) & (1<<2 - 1))
				rel.Pcrel = ri.Addr&(1<<30) != 0
				rel.Value = ri.Symnum
				rel.Scattered = true
			} else {
				switch bo {
				case binary.LittleEndian:
					rel.Addr = ri.Addr
					rel.Value = ri.Symnum & (1<<24 - 1)
					rel.Pcrel = ri.Symnum&(1<<24) != 0
					rel.Len = uint8((ri.Symnum >> 25) & (1<<2 - 1))
					rel.Extern = ri.Symnum&(1<<27) != 0
					rel.Type = uint8((ri.Symnum >> 28) & (1<<4 - 1))
				case binary.BigEndian:
					rel.Addr = ri.Addr
					rel.Value = ri.Symnum >> 8
					rel.Pcrel = ri.Symnum&(1<<7) != 0
					rel.Len = uint8((ri.Symnum >> 5) & (1<<2 - 1))
					rel.Extern = ri.Symnum&(1<<4) != 0
					rel.Type = uint8(ri.Symnum & (1<<4 - 1))
				default:
					panic("unreachable")
				}
			}
		}
	}

	return nil
}

func cstring(b []byte) string {
	i := bytes.IndexByte(b, 0)
	if i == -1 {
		i = len(b)
	}
	return string(b[0:i])
}

func readString(r io.Reader) (string, error) {
	var b byte
	var str string

	for {
		err := binary.Read(r, binary.BigEndian, &b)

		if err != nil {
			return str, err
		}

		if b == '\x00' {
			return str, nil
		}

		str += string(b)
	}
}

func (f *File) is64bit() bool { return f.FileHeader.Magic == types.Magic64 }

func (f *File) pointerSize() uint64 {
	if f.is64bit() {
		return 8
	}
	return 4
}

func (f *File) symbolSize() int {
	if f.is64bit() {
		return binary.Size(types.Nlist64{})
	}
	return binary.Size(types.Nlist32{})
}

func (f *File) has16KPages() bool {
	switch f.CPU {
	case types.CPUArm64, types.CPUArm6432:
		return true
	case types.CPUArm:
		if f.Type != types.MH_KEXT_BUNDLE {
			return false
		}
		return f.SubCPU == types.CPUSubtypeArmV7K
	default:
		return false
	}
}

func (f *File) preferredLoadAddress() uint64 {
	if text := f.Segment("__TEXT"); text != nil {
		return text.Addr
	}
	return 0
}

// ReadAt reads data at offset within MachO
func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	return f.cr.ReadAt(p, off) // TODO: should this be f.cr  or f.sr?
}

// GetOffset returns the file offset for a given virtual address
func (f *File) GetOffset(address uint64) (uint64, error) {
	return f.vma.GetOffset(address)
}

func (f *File) getOffset(address uint64) (uint64, error) {
	for _, seg := range f.Segments() {
		if seg.Addr <= address && address < seg.Addr+seg.Memsz {
			return (address - seg.Addr) + seg.Offset, nil
		}
	}
	return 0, fmt.Errorf("address %#x not within any segment's adress range", address)
}

// GetVMAddress returns the virtal address for a given file offset
func (f *File) GetVMAddress(offset uint64) (uint64, error) {
	return f.vma.GetVMAddress(offset)
}

func (f *File) getVMAddress(offset uint64) (uint64, error) {
	for _, seg := range f.Segments() {
		if seg.Offset <= offset && offset < seg.Offset+seg.Filesz {
			return (offset - seg.Offset) + seg.Addr, nil
		}
	}
	return 0, fmt.Errorf("offset %#x not within any segment's file offset range", offset)
}

// GetBaseAddress returns the MachO's preferred load address
func (f *File) GetBaseAddress() uint64 {
	return f.preferredLoadAddress()
}

// GetPointer returns pointer at a given offset
func (f *File) GetPointer(offset uint64) (uint64, error) {
	if _, err := f.cr.Seek(int64(offset), io.SeekStart); err != nil {
		return 0, fmt.Errorf("failed to Seek to offset %#x: %v", offset, err)
	}
	var ptr uint64
	if err := binary.Read(f.cr, binary.LittleEndian, &ptr); err != nil {
		return 0, fmt.Errorf("failed to read pointer at offset %#x: %v", offset, err)
	}
	return f.vma.Convert(ptr), nil
}

// GetPointerAtAddress returns pointer at a given virtual address
func (f *File) GetPointerAtAddress(address uint64) (uint64, error) {
	if err := f.cr.SeekToAddr(address); err != nil {
		return 0, fmt.Errorf("failed to Seek to address %#x: %v", address, err)
	}
	var ptr uint64
	if err := binary.Read(f.cr, binary.LittleEndian, &ptr); err != nil {
		return 0, fmt.Errorf("failed to read pointer @ %#x: %v", address, err)
	}
	return f.vma.Convert(ptr), nil
}

// SlidePointer returns slid or un-chained pointer
func (f *File) SlidePointer(ptr uint64) uint64 {
	return f.vma.Convert(ptr)
}

func (f *File) convertToVMAddr(value uint64) uint64 {
	if value == 0 {
		return 0
	}

	if f.HasFixups() {
		dcf, err := f.DyldChainedFixups()
		if err != nil {
			return value
		}
		if target, ok := dcf.IsRebase(value, f.GetBaseAddress()); ok {
			return target + f.preferredLoadAddress()
		} else if bind, addend, ok := dcf.IsBind(value); ok {
			if bind.Import.LibOrdinal() == types.BIND_SPECIAL_DYLIB_SELF {
				symAddr, err := f.FindSymbolAddress(bind.Name)
				if err != nil {
					return 0
				}
				return uint64(int64(symAddr) + addend)
			}
		}
	}

	return value
}

// GetBindName returns the import name for a given dyld chained pointer
func (f *File) GetBindName(pointer uint64) (string, error) {
	var err error
	if f.HasFixups() {
		if f.dcf == nil {
			f.dcf, err = f.DyldChainedFixups()
			if err != nil {
				return "", fmt.Errorf("failed to parse dyld chained fixups: %v", err)
			}
		}
		if len(f.dcf.Imports) > 0 {
			if bind, _, ok := f.dcf.IsBind(pointer); ok {
				return bind.Name, nil
			}
			return "", fmt.Errorf("pointer %#x is not a bind", pointer)
		}
		return "", fmt.Errorf("MachO does not contain dyld chained fixups importts")
	} else {
		binds, err := f.GetBindInfo()
		if err != nil {
			return "", fmt.Errorf("failed to parse LC_DYLD_INFO_ONLY bind info: %v", err)
		}
		for _, bind := range binds {
			if (bind.Start + bind.Offset) == pointer {
				return bind.Name, nil
			}
		}
		return "", fmt.Errorf("pointer %#x is not a bind", pointer)
	}
}

// GetCString returns a c-string at a given virtual address in the MachO
func (f *File) GetCString(addr uint64) (string, error) {
	if err := f.cr.SeekToAddr(addr); err != nil {
		return "", fmt.Errorf("failed to Seek to address %#x: %v", addr, err)
	}

	s, err := bufio.NewReader(f.cr).ReadString('\x00')
	if err != nil {
		return "", fmt.Errorf("failed to read strubg at address %#x, %v", addr, err)
	}

	if len(s) > 0 {
		return strings.Trim(s, "\x00"), nil
	}

	return "", fmt.Errorf("string not found at address %#x", addr)
}

// GetCStringAtOffset returns a c-string at a given offset into the MachO
func (f *File) GetCStringAtOffset(strOffset int64) (string, error) {
	if _, err := f.cr.Seek(strOffset, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to Seek to offset %#x: %v", strOffset, err)
	}

	s, err := bufio.NewReader(f.cr).ReadString('\x00')
	if err != nil {
		return "", fmt.Errorf("failed to ReadString as offset %#x, %v", strOffset, err)
	}

	if len(s) > 0 {
		return strings.Trim(s, "\x00"), nil
	}

	return "", fmt.Errorf("string not found at offset %#x", strOffset)
}

// IsCString returns cstring at given virtual address if is in a CstringLiterals section
func (f *File) IsCString(addr uint64) (string, bool) {
	for _, sec := range f.Sections {
		if sec.Flags.IsCstringLiterals() {
			if sec.Addr <= addr && addr < sec.Addr+sec.Size {
				str, err := f.GetCString(addr)
				if err != nil {
					return "", false
				}
				return str, true
			}
		}
	}
	return "", false
}

func (f *File) GetLoadsByName(name string) []Load {
	var loads []Load
	for _, l := range f.Loads {
		if l.Command().Command().String() == name {
			loads = append(loads, l)
		}
	}
	return loads
}

// Segment returns the first Segment with the given name, or nil if no such segment exists.
func (f *File) Segment(name string) *Segment {
	for _, l := range f.Loads {
		if s, ok := l.(*Segment); ok && s.Name == name {
			return s
		}
	}
	return nil
}

// Segments returns all Segments.
func (f *File) Segments() Segments {
	var segs Segments
	for _, l := range f.Loads {
		if s, ok := l.(*Segment); ok {
			segs = append(segs, s)
		}
	}
	// sort.Sort(segs)
	return segs
}

// GetSectionsForSegment returns all the segment's sections or nil if it doesn't have any
func (f *File) GetSectionsForSegment(name string) []*types.Section {
	var secs []*types.Section
	if seg := f.Segment(name); seg != nil {
		if seg.Nsect > 0 {
			for i := uint32(0); i < seg.Nsect; i++ {
				if int(i+seg.Firstsect) < len(f.Sections) {
					secs = append(secs, f.Sections[i+seg.Firstsect])
				}
			}
			return secs
		}
	}
	return nil
}

// Section returns the section with the given name in the given segment,
// or nil if no such section exists.
func (f *File) Section(segment, section string) *types.Section {
	for _, sec := range f.Sections {
		if sec.Seg == segment && sec.Name == section {
			return sec
		}
	}
	return nil
}

// FindSegmentForVMAddr returns the segment containing a given virtual memory ddress.
func (f *File) FindSegmentForVMAddr(vmAddr uint64) *Segment {
	for _, seg := range f.Segments() {
		if seg.Addr <= vmAddr && vmAddr < seg.Addr+seg.Memsz {
			return seg
		}
	}
	return nil
}

// FindSectionForVMAddr returns the section containing a given virtual memory ddress.
func (f *File) FindSectionForVMAddr(vmAddr uint64) *types.Section {
	for _, sec := range f.Sections {
		if sec.Addr <= vmAddr && vmAddr < sec.Addr+sec.Size {
			return sec
		}
	}
	return nil
}

// UUID returns the UUID load command, or nil if no UUID exists.
func (f *File) UUID() *UUID {
	for _, l := range f.Loads {
		if u, ok := l.(*UUID); ok {
			return u
		}
	}
	return nil
}

// DylibID returns the dylib ID load command, or nil if no dylib ID exists.
func (f *File) DylibID() *IDDylib {
	for _, l := range f.Loads {
		if s, ok := l.(*IDDylib); ok {
			return s
		}
	}
	return nil
}

// DyldInfo returns the dyld info load command, or nil if no dyld info exists.
func (f *File) DyldInfo() *DyldInfo {
	for _, l := range f.Loads {
		if s, ok := l.(*DyldInfo); ok {
			return s
		}
	}
	return nil
}

// DyldInfoOnly returns the dyld info only load command, or nil if no dyld info only exists.
func (f *File) DyldInfoOnly() *DyldInfoOnly {
	for _, l := range f.Loads {
		if s, ok := l.(*DyldInfoOnly); ok {
			return s
		}
	}
	return nil
}

// SourceVersion returns the source version load command, or nil if no source version exists.
func (f *File) SourceVersion() *SourceVersion {
	for _, l := range f.Loads {
		if s, ok := l.(*SourceVersion); ok {
			return s
		}
	}
	return nil
}

// BuildVersion returns the build version load command, or nil if no build version exists.
func (f *File) BuildVersion() *BuildVersion {
	for _, l := range f.Loads {
		if s, ok := l.(*BuildVersion); ok {
			return s
		}
	}
	return nil
}

// FileSets returns an array of Fileset entries.
func (f *File) FileSets() []*FilesetEntry {
	var fsets []*FilesetEntry
	for _, l := range f.Loads {
		if fs, ok := l.(*FilesetEntry); ok {
			fsets = append(fsets, fs)
		}
	}
	return fsets
}

// GetFileSetFileByName returns the Fileset MachO for a given name.
func (f *File) GetFileSetFileByName(name string) (*File, error) {
	for _, l := range f.Loads {
		if fs, ok := l.(*FilesetEntry); ok {
			if strings.EqualFold(fs.EntryID, name) || strings.HasSuffix(strings.ToLower(fs.EntryID), strings.ToLower(name)) {
				return NewFile(io.NewSectionReader(f.sr, int64(fs.FileOffset), 1<<63-1), FileConfig{
					Offset:        int64(fs.FileOffset),
					SectionReader: f.sr,
					CacheReader:   f.cr,
					VMAddrConverter: types.VMAddrConverter{
						Converter:    f.convertToVMAddr,
						VMAddr2Offet: f.GetOffset,
						Offet2VMAddr: f.GetVMAddress,
					},
				})
			}
		}
	}
	return nil, fmt.Errorf("fileset does NOT contain %s", name)
}

// DataInCode returns the LC_DATA_IN_CODE, or nil if none exists.
func (f *File) DataInCode() *DataInCode {
	for _, l := range f.Loads {
		if s, ok := l.(*DataInCode); ok {
			return s
		}
	}
	return nil
}

// FunctionStarts returns the function starts array, or nil if none exists.
func (f *File) FunctionStarts() *FunctionStarts {
	for _, l := range f.Loads {
		if s, ok := l.(*FunctionStarts); ok {
			return s
		}
	}
	return nil
}

// GetFunctions returns the function array, or nil if none exists.
func (f *File) GetFunctions(data ...byte) []types.Function {

	if len(f.functions) > 0 {
		return f.functions
	}

	var funcs []types.Function

	fs := f.FunctionStarts()
	if fs == nil {
		return nil
	}

	var fsr *bytes.Reader
	if len(data) > 0 {
		fsr = bytes.NewReader(data)
	} else {
		ldat := make([]byte, fs.Size)
		if _, err := f.cr.ReadAt(ldat, int64(fs.Offset)); err != nil {
			return nil
		}
		fsr = bytes.NewReader(ldat)
	}

	offset, err := trie.ReadUleb128(fsr)
	if err != nil {
		return nil
	}

	startVMA := offset + f.GetBaseAddress()

	for {
		offset, err = trie.ReadUleb128(fsr)
		if err == io.EOF {
			break
		}
		if offset == 0 {
			break
		}
		if err != nil {
			return nil
		}

		funcs = append(funcs, types.Function{
			StartAddr: startVMA,
			EndAddr:   startVMA + offset,
		})

		startVMA += offset
	}

	// get last function
	if s := f.FindSectionForVMAddr(startVMA); s != nil {
		funcs = append(funcs, types.Function{
			StartAddr: startVMA,
			EndAddr:   s.Addr + s.Size,
		})
	}

	// cache parsed functions
	f.functions = funcs

	return funcs
}

// GetFunctionForVMAddr returns the function containing a given virual address
func (f *File) GetFunctionForVMAddr(addr uint64) (types.Function, error) {
	for _, fn := range f.GetFunctions() {
		if addr >= fn.StartAddr && addr < fn.EndAddr {
			return fn, nil
		}
	}
	return types.Function{}, fmt.Errorf("address %#016x not in any function", addr)
}

// GetFunctionsForRange returns the functions contained in a given virual address range
func (f *File) GetFunctionsForRange(start, end uint64) ([]types.Function, error) {
	var funcs []types.Function
	for _, fn := range f.GetFunctions() {
		if start >= fn.StartAddr && fn.StartAddr < end {
			funcs = append(funcs, fn)
		}
	}
	return funcs, nil
}

func (f *File) GetFunctionData(fn types.Function) ([]byte, error) {
	data := make([]byte, fn.EndAddr-fn.StartAddr)
	if _, err := f.cr.ReadAtAddr(data, fn.StartAddr); err != nil {
		return nil, fmt.Errorf("failed to read data at address %#x: %v", fn.StartAddr, err)
	}
	return data, nil
}

// CodeSignature returns the code signature, or nil if none exists.
func (f *File) CodeSignature() *CodeSignature {
	for _, l := range f.Loads {
		if s, ok := l.(*CodeSignature); ok {
			return s
		}
	}
	return nil
}

// DyldExportsTrie returns the dyld export trie load command, or nil if no dyld info exists.
func (f *File) DyldExportsTrie() *DyldExportsTrie {
	for _, l := range f.Loads {
		if s, ok := l.(*DyldExportsTrie); ok {
			return s
		}
	}
	return nil
}

// DyldExports returns the dyld export trie symbols
func (f *File) GetDyldExport(symbol string) (*trie.TrieExport, error) {
	if dxt := f.DyldExportsTrie(); dxt == nil {
		return nil, fmt.Errorf("macho does not contain LC_DYLD_EXPORTS_TRIE")
	} else {
		var err error
		var r *bytes.Reader
		if f.exptrieData != nil {
			r = bytes.NewReader(f.exptrieData)
		} else {
			f.exptrieData = make([]byte, dxt.Size)
			if _, err := f.cr.ReadAt(f.exptrieData, int64(dxt.Offset)); err != nil {
				f.exptrieData = nil
				return nil, fmt.Errorf("failed to read %s data at offset=%#x; %v", types.LC_DYLD_EXPORTS_TRIE, int64(dxt.Offset), err)
			}
			r = bytes.NewReader(f.exptrieData)
		}
		if _, err = trie.WalkTrie(r, symbol); err != nil {
			return nil, err
		}
		return trie.ReadExport(r, symbol, f.preferredLoadAddress())
	}
}

// DyldExports returns the dyld export trie symbols
func (f *File) DyldExports() ([]trie.TrieExport, error) {
	var err error
	if f.exp != nil {
		return f.exp, nil
	}
	if dxt := f.DyldExportsTrie(); dxt != nil {
		if dxt.Size == 0 {
			return []trie.TrieExport{}, nil
		}
		data := make([]byte, dxt.Size)
		if _, err := f.cr.ReadAt(data, int64(dxt.Offset)); err != nil {
			return nil, fmt.Errorf("failed to read %s data at offset=%#x; %v", types.LC_DYLD_EXPORTS_TRIE, int64(dxt.Offset), err)
		}
		f.exp, err = trie.ParseTrieExports(bytes.NewReader(data), f.GetBaseAddress())
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %v", types.LC_DYLD_EXPORTS_TRIE, err)
		}
		return f.exp, nil
	}

	return nil, fmt.Errorf("macho does not contain LC_DYLD_EXPORTS_TRIE")
}

// HasFixups does macho contain a LC_DYLD_CHAINED_FIXUPS load command
func (f *File) HasFixups() bool {
	for _, l := range f.Loads {
		if _, ok := l.(*DyldChainedFixups); ok {
			return true
		}
	}
	return false
}

// DyldChainedFixups returns the dyld chained fixups.
func (f *File) DyldChainedFixups() (*fixupchains.DyldChainedFixups, error) {
	if f.dcf != nil { // is cached
		return f.dcf, nil
	}

	for _, l := range f.Loads {
		if dcfLC, ok := l.(*DyldChainedFixups); ok {
			data := make([]byte, dcfLC.Size)
			if _, err := f.cr.ReadAt(data, int64(dcfLC.Offset)); err != nil {
				return nil, fmt.Errorf("failed to read DyldChainedFixups data at offset=%#x; %v", int64(dcfLC.Offset), err)
			}
			dcf := fixupchains.NewChainedFixups(bytes.NewReader(data), &f.cr, f.ByteOrder)
			if err := dcf.ParseStarts(); err != nil {
				return nil, fmt.Errorf("failed to parse dyld chained fixup starts: %v", err)
			}
			segs := f.Segments()
			for idx, start := range dcf.Starts {
				if start.PageStarts != nil {
					// Replacing SegmentOffset(vmaddr) with FileOffset
					// (for static analysis of binaries with split segs
					// since we aren't actually loading the MachO
					// ref: void Adjustor<P>::adjustChainedFixups() in
					// dyld-750.6/dyld3/shared-cache/AdjustDylibSegments.cpp
					dcf.Starts[idx].SegmentOffset = segs[idx].Offset
				}
			}

			dcf, err := dcf.Parse()
			if err != nil {
				return nil, fmt.Errorf("failed to parse dyld chained fixups: %v", err)
			}

			f.dcf = dcf // cache

			return dcf, nil
		}
	}
	return nil, fmt.Errorf("macho does not contain LC_DYLD_CHAINED_FIXUPS")
}

func (f *File) ForEachV2SplitSegReference(handler func(fromSectionIndex, fromSectionOffset, toSectionIndex, toSectionOffset uint64, kind types.SplitInfoKind)) error {
	for _, l := range f.Loads {
		if si, ok := l.(*SplitInfo); ok {
			if si.Size == 0 {
				return nil
			}
			data := make([]byte, si.Size)
			if _, err := f.cr.ReadAt(data, int64(si.Offset)); err != nil {
				return fmt.Errorf("failed to read %s data at offset=%#x; %v", types.LC_SEGMENT_SPLIT_INFO, int64(si.Offset), err)
			}

			r := bytes.NewReader(data)

			var version uint8
			if err := binary.Read(r, f.ByteOrder, &version); err != nil {
				return fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO Version: %v", err)
			}
			if version != types.DYLD_CACHE_ADJ_V2_FORMAT {
				return nil
			}

			sectionCount, err := trie.ReadUleb128(r)
			if err != nil {
				return fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO SectionCount: %v", err)
			}

			for i := uint64(0); i < sectionCount; i++ {
				fromSectionIndex, err := trie.ReadUleb128(r)
				if err != nil {
					return fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO fromSectionIndex: %v", err)
				}
				toSectionIndex, err := trie.ReadUleb128(r)
				if err != nil {
					return fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO toSectionIndex: %v", err)
				}
				toOffsetCount, err := trie.ReadUleb128(r)
				if err != nil {
					return fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO toOffsetCount: %v", err)
				}

				var toSectionOffset uint64
				for j := uint64(0); j < toOffsetCount; j++ {
					toSectionDelta, err := trie.ReadUleb128(r)
					if err != nil {
						return fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO toSectionDelta: %v", err)
					}
					fromOffsetCount, err := trie.ReadUleb128(r)
					if err != nil {
						return fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO fromOffsetCount: %v", err)
					}

					toSectionOffset += toSectionDelta
					for k := uint64(0); k < fromOffsetCount; k++ {
						kind, err := trie.ReadUleb128(r)
						if err != nil {
							return fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO kind: %v", err)
						}
						if kind > 13 {
							return fmt.Errorf("invalid LC_SEGMENT_SPLIT_INFO kind: %d", kind)
						}

						fromSectDeltaCount, err := trie.ReadUleb128(r)
						if err != nil {
							return fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO fromSectDeltaCount: %v", err)
						}

						var fromSectionOffset uint64
						for l := uint64(0); l < fromSectDeltaCount; l++ {
							delta, err := trie.ReadUleb128(r)
							if err != nil {
								return fmt.Errorf("failed to read LC_SEGMENT_SPLIT_INFO delta: %v", err)
							}

							fromSectionOffset += delta

							handler(fromSectionIndex, fromSectionOffset, toSectionIndex, toSectionOffset, types.SplitInfoKind(kind))
						}
					}
				}
			}
		}
	}
	return nil
}

func (f *File) GetEmbeddedInfoPlist() ([]byte, error) {
	infoSec := f.Section("__TEXT", "__info_plist")
	if infoSec == nil {
		return nil, fmt.Errorf("no __TEXT.__info_plist section")
	}
	data, err := infoSec.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to read __TEXT.__info_plist section data: %v", err)
	}
	return data, nil
}

const (
	BitcodeWrapperMagic = 0x0b17c0de
	RawBitcodeMagic     = 0xdec04342 // 'BC' 0xc0de
)

type BitstreamWrapperHeader struct {
	Magic   uint32
	Version uint32
	Offset  uint32
	Size    uint32
	CPUType uint32
}

func (f *File) GetEmbeddedLLVMBitcode() (*xar.Reader, error) {
	llvmBundle := f.Section("__LLVM", "__bundle")
	if llvmBundle == nil {
		return nil, fmt.Errorf("no %s.%s section", llvmBundle.Seg, llvmBundle.Name)
	}
	data, err := llvmBundle.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to read %s.%s section data: %v", llvmBundle.Seg, llvmBundle.Name, err)
	}
	return xar.NewReader(bytes.NewReader(data), int64(len(data)))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create xar reader: %v", err)
	// }

	// if xr.HasSignature() {
	// 	fmt.Println(xr.Certificates)
	// }

	// for _, xf := range xr.File {
	// 	fmt.Printf("name: %s, type: %v, info: %v, valid_checksum: %t\n", xf.Name, xf.Type, xf.Info, xf.VerifyChecksum())
	// 	f, err := xf.Open()
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to open xar file: %v", err)
	// 	}
	// 	data, err := io.ReadAll(f)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to read xar file: %v", err)
	// 	}
	// 	var header BitstreamWrapperHeader
	// 	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &header); err != nil {
	// 		return nil, fmt.Errorf("failed to read bitstream wrapper header: %v", err)
	// 	}
	// 	_ = header
	// 	os.WriteFile(xf.Name, data, 0644)
	// 	f.Close()
	// }
}

// DWARF returns the DWARF debug information for the Mach-O file.
func (f *File) DWARF() (*dwarf.Data, error) {
	dwarfSuffix := func(s *types.Section) string {
		switch {
		case strings.HasPrefix(s.Name, "__debug_"):
			return s.Name[8:]
		case strings.HasPrefix(s.Name, "__zdebug_"):
			return s.Name[9:]
		default:
			return ""
		}
	}
	appleSuffix := func(s *types.Section) string {
		switch {
		case strings.HasPrefix(s.Name, "__apple_"):
			return s.Name[8:]
		default:
			return ""
		}
	}
	sectionData := func(s *types.Section) ([]byte, error) {
		b, err := s.Data()
		if err != nil && uint64(len(b)) < s.Size {
			return nil, err
		}

		if len(b) >= 12 && string(b[:4]) == "ZLIB" {
			dlen := binary.BigEndian.Uint64(b[4:12])
			dbuf := make([]byte, dlen)
			r, err := zlib.NewReader(bytes.NewBuffer(b[12:]))
			if err != nil {
				return nil, err
			}
			if _, err := io.ReadFull(r, dbuf); err != nil {
				return nil, err
			}
			if err := r.Close(); err != nil {
				return nil, err
			}
			b = dbuf
		}
		return b, nil
	}

	// There are many other DWARF sections, but these
	// are the ones the debug/dwarf package uses.
	// Don't bother loading others.
	var dat = map[string][]byte{"abbrev": nil, "info": nil, "str": nil, "line": nil, "ranges": nil}
	for _, s := range f.Sections {
		suffix := dwarfSuffix(s)
		if suffix == "" {
			continue
		}
		if _, ok := dat[suffix]; !ok {
			continue
		}
		b, err := sectionData(s)
		if err != nil {
			return nil, err
		}
		dat[suffix] = b
	}

	d, err := dwarf.New(dat["abbrev"], nil, nil, dat["info"], dat["line"], nil, dat["ranges"], dat["str"])
	if err != nil {
		return nil, err
	}

	// Look for DWARF4 .debug_types sections and DWARF5 sections.
	for i, s := range f.Sections {
		suffix := dwarfSuffix(s)
		if suffix == "" {
			continue
		}
		if _, ok := dat[suffix]; ok {
			// Already handled.
			continue
		}

		b, err := sectionData(s)
		if err != nil {
			return nil, err
		}

		if suffix == "types" {
			err = d.AddTypes(fmt.Sprintf("types-%d", i), b)
		} else {
			err = d.AddSection(".debug_"+suffix, b)
		}
		if err != nil {
			return nil, err
		}
	}

	// Look for Apple HASH table .apple_names, .apple_types, .apple_namespaces or .apple_objc sections.
	for _, s := range f.Sections {
		suffix := appleSuffix(s)
		if suffix == "" {
			continue
		}
		if _, ok := dat[suffix]; ok {
			// Already handled.
			continue
		}

		b, err := sectionData(s)
		if err != nil {
			return nil, err
		}
		dat[suffix] = b
		// TODO: finish implementing this.
		if err := d.AddHashes(suffix, b); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (f *File) GetBindInfo() (types.Binds, error) {
	if f.binds != nil {
		return f.binds, nil
	}
	if dinfo := f.DyldInfo(); dinfo != nil {
		if dinfo.BindSize > 0 {
			dat := make([]byte, dinfo.BindSize)
			if _, err := f.cr.ReadAt(dat, int64(dinfo.BindOff)); err != nil {
				return nil, fmt.Errorf("failed to read bind info: %v", err)
			}
			bs, err := f.parseBinds(bytes.NewReader(dat), types.BIND_KIND)
			if err != nil {
				return nil, err
			}
			f.binds = append(f.binds, bs...)
		}
		if dinfo.WeakBindSize > 0 {
			dat := make([]byte, dinfo.WeakBindSize)
			if _, err := f.cr.ReadAt(dat, int64(dinfo.WeakBindOff)); err != nil {
				return nil, fmt.Errorf("failed to read weak bind info: %v", err)
			}
			bs, err := f.parseBinds(bytes.NewReader(dat), types.WEAK_KIND)
			if err != nil {
				return nil, err
			}
			f.binds = append(f.binds, bs...)
		}
		if dinfo.LazyBindSize > 0 {
			dat := make([]byte, dinfo.LazyBindSize)
			if _, err := f.cr.ReadAt(dat, int64(dinfo.LazyBindOff)); err != nil {
				return nil, fmt.Errorf("failed to read lazy bind info: %v", err)
			}
			bs, err := f.parseBinds(bytes.NewReader(dat), types.LAZY_KIND)
			if err != nil {
				return nil, err
			}
			f.binds = append(f.binds, bs...)
		}
	} else if dinfo := f.DyldInfoOnly(); dinfo != nil {
		if dinfo.BindSize > 0 {
			dat := make([]byte, dinfo.BindSize)
			if _, err := f.cr.ReadAt(dat, int64(dinfo.BindOff)); err != nil {
				return nil, fmt.Errorf("failed to read bind info: %v", err)
			}
			bs, err := f.parseBinds(bytes.NewReader(dat), types.BIND_KIND)
			if err != nil {
				return nil, err
			}
			f.binds = append(f.binds, bs...)
		}
		if dinfo.WeakBindSize > 0 {
			dat := make([]byte, dinfo.WeakBindSize)
			if _, err := f.cr.ReadAt(dat, int64(dinfo.WeakBindOff)); err != nil {
				return nil, fmt.Errorf("failed to read weak bind info: %v", err)
			}
			bs, err := f.parseBinds(bytes.NewReader(dat), types.WEAK_KIND)
			if err != nil {
				return nil, err
			}
			f.binds = append(f.binds, bs...)
		}
		if dinfo.LazyBindSize > 0 {
			dat := make([]byte, dinfo.LazyBindSize)
			if _, err := f.cr.ReadAt(dat, int64(dinfo.LazyBindOff)); err != nil {
				return nil, fmt.Errorf("failed to read lazy bind info: %v", err)
			}
			bs, err := f.parseBinds(bytes.NewReader(dat), types.LAZY_KIND)
			if err != nil {
				return nil, err
			}
			f.binds = append(f.binds, bs...)
		}
	} else {
		return nil, ErrMachODyldInfoNotFound
	}

	return f.binds, nil
}

func (f *File) GetRebaseInfo() ([]types.Rebase, error) {
	if dinfo := f.DyldInfo(); dinfo != nil {
		if dinfo.RebaseSize > 0 {
			dat := make([]byte, dinfo.RebaseSize)
			if _, err := f.cr.ReadAt(dat, int64(dinfo.RebaseOff)); err != nil {
				return nil, fmt.Errorf("failed to read rebase info: %v", err)
			}
			return f.parseRebase(bytes.NewReader(dat))
		}
	} else if dinfo := f.DyldInfoOnly(); dinfo != nil {
		if dinfo.RebaseSize > 0 {
			dat := make([]byte, dinfo.RebaseSize)
			if _, err := f.cr.ReadAt(dat, int64(dinfo.RebaseOff)); err != nil {
				return nil, fmt.Errorf("failed to read rebase info: %v", err)
			}
			return f.parseRebase(bytes.NewReader(dat))
		}
	} else {
		return nil, ErrMachODyldInfoNotFound
	}
	return nil, nil
}

func (f *File) GetExports() ([]trie.TrieExport, error) {
	if dinfo := f.DyldInfo(); dinfo != nil {
		if dinfo.ExportSize > 0 {
			dat := make([]byte, dinfo.ExportSize)
			if _, err := f.cr.ReadAt(dat, int64(dinfo.ExportOff)); err != nil {
				return nil, fmt.Errorf("failed to read bind info: %v", err)
			}
			return trie.ParseTrieExports(bytes.NewReader(dat), f.GetBaseAddress())
		}
	} else if dinfo := f.DyldInfoOnly(); dinfo != nil {
		if dinfo.ExportSize > 0 {
			// addr := linkedit.Addr + (uint64(dinfo.ExportOff) - linkedit.Offset)
			dat := make([]byte, dinfo.ExportSize)
			if _, err := f.cr.ReadAt(dat, int64(dinfo.ExportOff)); err != nil {
				return nil, fmt.Errorf("failed to read bind info: %v", err)
			}
			return trie.ParseTrieExports(bytes.NewReader(dat), f.GetBaseAddress())
		}
	} else {
		return nil, ErrMachODyldInfoNotFound
	}
	return nil, nil
}

func (f *File) parseBinds(r *bytes.Reader, kind types.BindKind) ([]types.Bind, error) {
	var binds []types.Bind
	var ordinalTable []types.Bind
	var ordinalTableSize uint64
	var segOffset uint64

	useThreadedRebaseBind := false
	bind := types.Bind{Kind: kind}

	for {
		ptr, err := r.ReadByte()

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		imm := ptr & types.BIND_IMMEDIATE_MASK
		opcode := ptr & types.BIND_OPCODE_MASK

		switch opcode {
		case types.BIND_OPCODE_DONE:
			if kind != types.LAZY_KIND {
				return binds, nil
			}
			bind = types.Bind{Kind: kind}
		case types.BIND_OPCODE_SET_DYLIB_ORDINAL_IMM:
			bind.Dylib = f.LibraryOrdinalName(int(imm))
		case types.BIND_OPCODE_SET_DYLIB_ORDINAL_ULEB:
			i, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			bind.Dylib = f.LibraryOrdinalName(int(i))
		case types.BIND_OPCODE_SET_DYLIB_SPECIAL_IMM:
			if imm == 0 {
				bind.Dylib = f.LibraryOrdinalName(int(imm))
			} else {
				bind.Dylib = f.LibraryOrdinalName(int(types.BIND_OPCODE_MASK | imm))
			}
		case types.BIND_OPCODE_SET_SYMBOL_TRAILING_FLAGS_IMM:
			s, err := readString(r)
			if err != nil {
				return nil, err
			}
			bind.Name = strings.Trim(s, "\x00")
			bind.Flags = imm
		case types.BIND_OPCODE_SET_TYPE_IMM:
			bind.Type = imm
		case types.BIND_OPCODE_SET_ADDEND_SLEB:
			add, err := trie.ReadSleb128(r)
			if err != nil {
				return nil, err
			}
			bind.Addend = add
		case types.BIND_OPCODE_SET_SEGMENT_AND_OFFSET_ULEB:
			segOffset, err = trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			bind.Start = f.Segments()[imm].Addr
			bind.Segment = f.Segments()[imm].Name
		case types.BIND_OPCODE_ADD_ADDR_ULEB:
			out, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			segOffset += out
		case types.BIND_OPCODE_DO_BIND:
			if useThreadedRebaseBind {
				ordinalTable = append(ordinalTable, bind)
			} else {
				if sec := f.FindSectionForVMAddr(f.Segment(bind.Segment).Addr + segOffset); sec != nil {
					bind.Section = sec.Name
				}
				bind.Offset = segOffset
				binds = append(binds, bind)
				segOffset += f.pointerSize()
			}
		case types.BIND_OPCODE_DO_BIND_ADD_ADDR_ULEB:
			if sec := f.FindSectionForVMAddr(f.Segment(bind.Segment).Addr + segOffset); sec != nil {
				bind.Section = sec.Name
			}
			bind.Offset = segOffset
			binds = append(binds, bind)
			off, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			segOffset += off + f.pointerSize()
		case types.BIND_OPCODE_DO_BIND_ADD_ADDR_IMM_SCALED:
			if sec := f.FindSectionForVMAddr(f.Segment(bind.Segment).Addr + segOffset); sec != nil {
				bind.Section = sec.Name
			}
			bind.Offset = segOffset
			binds = append(binds, bind)
			segOffset += uint64(imm)*f.pointerSize() + f.pointerSize()
		case types.BIND_OPCODE_DO_BIND_ULEB_TIMES_SKIPPING_ULEB:
			count, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			skip, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			for i := uint64(0); i < count; i++ {
				if sec := f.FindSectionForVMAddr(f.Segment(bind.Segment).Addr + segOffset); sec != nil {
					bind.Section = sec.Name
				}
				bind.Offset = segOffset
				binds = append(binds, bind)
				segOffset += skip + f.pointerSize()
			}
		case types.BIND_OPCODE_THREADED:
			switch imm {
			case types.BIND_SUBOPCODE_THREADED_SET_BIND_ORDINAL_TABLE_SIZE_ULEB:
				ordinalTableSize, err = trie.ReadUleb128(r)
				if err != nil {
					return nil, err
				}
				useThreadedRebaseBind = true
			case types.BIND_SUBOPCODE_THREADED_APPLY: // parse chain
				delta := uint64(0)
				for {
					var ptr uint64
					f.cr.Seek(int64(f.Segment(bind.Segment).Offset+segOffset), io.SeekStart)
					if err := binary.Read(f.cr, f.ByteOrder, &ptr); err != nil {
						return nil, fmt.Errorf("failed to read pointer: %v", err)
					}
					if (ptr & (1 << 62)) == 0 { // isRebase TODO: handle rebases
						if (ptr & (1 << 63)) != 0 { // isAuthenticated
							// uint16_t diversity = (uint16_t)(value >> 32);
							// bool hasAddressDiversity = (value & (1ULL << 48)) != 0;
							// uint8_t key = (uint8_t)((value >> 49) & 0x3);
							// static const char* keyNames[] = {
							// 	"IA", "IB", "DA", "DB"
							// };
							// printf("%-7s %-16s 0x%08llX %10s  %5lld %-16s %s%s with value 0x%016llX (JOP: diversity %d, address %s, %s)\n", segName, sectionName(segIndex, segStartAddr+segOffset), segStartAddr+segOffset, typeName, addend, fromDylib, symbolName, weak_import, value, diversity, hasAddressDiversity ? "true" : "false", keyNames[key]);
						} else {
							// // Regular pointer which needs to fit in 51-bits of value.
							// // C++ RTTI uses the top bit, so we'll allow the whole top-byte
							// // and the signed-extended bottom 43-bits to be fit in to 51-bits.
							// uint64_t top8Bits = value & 0x0007F80000000000ULL;
							// uint64_t bottom43Bits = value & 0x000007FFFFFFFFFFULL;
							// uint64_t targetValue = ( top8Bits << 13 ) | (((intptr_t)(bottom43Bits << 21) >> 21) & 0x00FFFFFFFFFFFFFF);
							// targetValue = targetValue + slide;
							// *(uint64_t*)address = targetValue;
						}
					} else {
						// the ordinal is bits [0..15]
						ord := ptr & 0xFFFF
						// XXX: there is a note in LIEF source about adding +1 to ordinalTableSize
						// when parsing it because of ld64 - I haven't investigated yet - fG!
						if ord > ordinalTableSize { // TODO: make sure this is right
							return nil, fmt.Errorf("bind ordinal is out of range")
						}
						ordinalTable[ord].Value = ptr
						ordinalTable[ord].Offset = segOffset
						binds = append(binds, ordinalTable[ord])
					}
					// The delta is bits [51..61]
					// And bit 62 is to tell us if we are a rebase (0) or bind (1)
					ptr &= ^uint64(1 << 62)
					delta = (ptr & 0x3FF8000000000000) >> 51
					segOffset += delta * f.pointerSize()
					if delta == 0 {
						break
					}
				}
			default:
				return nil, fmt.Errorf("bad threaded bind subopcode %#02x", imm)
			}
		default:
			return nil, fmt.Errorf("bad bind opcode %#02x", opcode)
		}
	}

	return binds, nil
}

func (f *File) parseRebase(r *bytes.Reader) ([]types.Rebase, error) {
	var rebase types.Rebase
	var rebases []types.Rebase

	for {
		ptr, err := r.ReadByte()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		imm := ptr & types.BIND_IMMEDIATE_MASK
		opcode := ptr & types.BIND_OPCODE_MASK

		switch opcode {
		case types.REBASE_OPCODE_DONE:
			return rebases, nil
		case types.REBASE_OPCODE_SET_TYPE_IMM:
			rebase.Type = imm
		case types.REBASE_OPCODE_SET_SEGMENT_AND_OFFSET_ULEB:
			off, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			rebase.Offset += off
			rebase.Start = f.Segments()[imm].Addr
			rebase.Segment = f.Segments()[imm].Name
		case types.REBASE_OPCODE_ADD_ADDR_ULEB:
			off, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			rebase.Offset += off
		case types.REBASE_OPCODE_ADD_ADDR_IMM_SCALED:
			rebase.Offset += uint64(imm) * f.pointerSize()
		case types.REBASE_OPCODE_DO_REBASE_IMM_TIMES:
			for i := byte(0); i < imm; i++ {
				f.cr.Seek(int64(f.Segment(rebase.Segment).Offset+rebase.Offset), io.SeekStart)
				if err := binary.Read(f.cr, f.ByteOrder, &rebase.Value); err != nil {
					return nil, fmt.Errorf("failed to read pointer: %v", err)
				}
				if sec := f.FindSectionForVMAddr(f.Segment(rebase.Segment).Addr + rebase.Offset); sec != nil {
					rebase.Section = sec.Name
				}
				rebases = append(rebases, rebase)
				rebase.Offset += f.pointerSize()
			}
		case types.REBASE_OPCODE_DO_REBASE_ULEB_TIMES:
			count, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			for i := uint64(0); i < count; i++ {
				f.cr.Seek(int64(f.Segment(rebase.Segment).Offset+rebase.Offset), io.SeekStart)
				if err := binary.Read(f.cr, f.ByteOrder, &rebase.Value); err != nil {
					return nil, fmt.Errorf("failed to read pointer: %v", err)
				}
				if sec := f.FindSectionForVMAddr(f.Segment(rebase.Segment).Addr + rebase.Offset); sec != nil {
					rebase.Section = sec.Name
				}
				rebases = append(rebases, rebase)
				rebase.Offset += f.pointerSize()
			}
		case types.REBASE_OPCODE_DO_REBASE_ADD_ADDR_ULEB:
			f.cr.Seek(int64(f.Segment(rebase.Segment).Offset+rebase.Offset), io.SeekStart)
			if err := binary.Read(f.cr, f.ByteOrder, &rebase.Value); err != nil {
				return nil, fmt.Errorf("failed to read pointer: %v", err)
			}
			off, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			if sec := f.FindSectionForVMAddr(f.Segment(rebase.Segment).Addr + rebase.Offset); sec != nil {
				rebase.Section = sec.Name
			}
			rebases = append(rebases, rebase)
			rebase.Offset += off + f.pointerSize()
		case types.REBASE_OPCODE_DO_REBASE_ULEB_TIMES_SKIPPING_ULEB:
			count, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			skip, err := trie.ReadUleb128(r)
			if err != nil {
				return nil, err
			}
			for i := uint64(0); i < count; i++ {
				f.cr.Seek(int64(f.Segment(rebase.Segment).Offset+rebase.Offset), io.SeekStart)
				if err := binary.Read(f.cr, f.ByteOrder, &rebase.Value); err != nil {
					return nil, fmt.Errorf("failed to read pointer: %v", err)
				}
				if sec := f.FindSectionForVMAddr(f.Segment(rebase.Segment).Addr + rebase.Offset); sec != nil {
					rebase.Section = sec.Name
				}
				rebases = append(rebases, rebase)
				rebase.Offset += skip + f.pointerSize()
			}
		default:
			return nil, fmt.Errorf("bad rebase opcode %#02x", opcode)
		}
	}

	return rebases, nil
}

// ImportedSymbols returns the names of all symbols
// referred to by the binary f that are expected to be
// satisfied by other libraries at dynamic load time.
func (f *File) ImportedSymbols() ([]Symbol, error) {
	if f.Dysymtab == nil || f.Symtab == nil {
		return nil, &FormatError{0, "missing symbol table", nil}
	}

	st := f.Symtab
	dt := f.Dysymtab
	var all []Symbol
	all = append(all, st.Syms[dt.Iundefsym:dt.Iundefsym+dt.Nundefsym]...)
	return all, nil
}

// ImportedSymbolNames returns the names of all symbols
// referred to by the binary f that are expected to be
// satisfied by other libraries at dynamic load time.
func (f *File) ImportedSymbolNames() ([]string, error) {
	var all []string

	syms, err := f.ImportedSymbols()
	if err != nil {
		return nil, fmt.Errorf("failed to get imported symbols: %v", err)
	}

	for _, s := range syms {
		all = append(all, s.Name)
	}

	return all, nil
}

// ImportedLibraries returns the paths of all libraries
// referred to by the binary f that are expected to be
// linked with the binary at dynamic link time.
func (f *File) ImportedLibraries() []string {
	var all []string
	for _, l := range f.Loads {
		switch v := l.(type) {
		case *LoadDylib:
			all = append(all, v.Name)
		case *WeakDylib:
			all = append(all, v.Name)
		case *ReExportDylib:
			all = append(all, v.Name)
		case *UpwardDylib:
			all = append(all, v.Name)
		case *LazyLoadDylib:
			all = append(all, v.Name)
		}
	}
	return all
}

// LibraryOrdinalName returns the depancy library oridinal's name
func (f *File) LibraryOrdinalName(libraryOrdinal int) string {
	dylibs := f.ImportedLibraries()

	if libraryOrdinal > 0 {
		if libraryOrdinal > len(dylibs) {
			return "ordinal-too-large"
		}
		return filepath.Base(dylibs[libraryOrdinal-1])
	}

	switch libraryOrdinal {
	case types.BIND_SPECIAL_DYLIB_SELF:
		return "this-image"
	case types.BIND_SPECIAL_DYLIB_MAIN_EXECUTABLE:
		return "main-executable"
	case types.BIND_SPECIAL_DYLIB_FLAT_LOOKUP:
		return "flat-namespace"
	case types.BIND_SPECIAL_DYLIB_WEAK_LOOKUP:
		return "weak-coalesce"
	default:
		return "unknown-ordinal"
	}
}

func (f *File) FindSymbolAddress(symbol string) (uint64, error) {
	if f.Symtab == nil {
		return 0, &FormatError{0, "missing symbol table", nil}
	}
	for _, sym := range f.Symtab.Syms {
		if strings.EqualFold(sym.Name, symbol) {
			return sym.Value, nil
		}
	}
	return 0, fmt.Errorf("symbol not found in macho symtab")
}

func (f *File) FindAddressSymbols(addr uint64) ([]Symbol, error) {
	if f.Symtab == nil {
		return nil, &FormatError{0, "missing symbol table", nil}
	}
	var syms []Symbol
	for _, sym := range f.Symtab.Syms {
		if sym.Value == addr {
			syms = append(syms, sym)
		}
	}
	if len(syms) > 0 {
		return syms, nil
	}
	return nil, fmt.Errorf("symbol(s) not found in macho symtab for addr 0x%016x", addr)
}
