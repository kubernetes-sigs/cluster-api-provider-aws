package types

//go:generate stringer -type=HeaderFileType -trimprefix=MH_ -output header_string.go

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
)

type Magic uint32

const (
	Magic32  Magic = 0xfeedface
	Magic64  Magic = 0xfeedfacf
	MagicFat Magic = 0xcafebabe
)

var magicStrings = []IntName{
	{uint32(Magic32), "32-bit MachO"},
	{uint32(Magic64), "64-bit MachO"},
	{uint32(MagicFat), "Fat MachO"},
}

func (i Magic) Int() uint32      { return uint32(i) }
func (i Magic) String() string   { return StringName(uint32(i), magicStrings, false) }
func (i Magic) GoString() string { return StringName(uint32(i), magicStrings, true) }

const (
	FileHeaderSize32 = 7 * 4
	FileHeaderSize64 = 8 * 4
)

// A FileHeader represents a Mach-O file header.
type FileHeader struct {
	Magic        Magic
	CPU          CPU
	SubCPU       CPUSubtype
	Type         HeaderFileType
	NCommands    uint32
	SizeCommands uint32
	Flags        HeaderFlag
	Reserved     uint32
}

func (h *FileHeader) Put(b []byte, o binary.ByteOrder) int {
	o.PutUint32(b[0:], uint32(h.Magic))
	o.PutUint32(b[4:], uint32(h.CPU))
	o.PutUint32(b[8:], uint32(h.SubCPU))
	o.PutUint32(b[12:], uint32(h.Type))
	o.PutUint32(b[16:], h.NCommands)
	o.PutUint32(b[20:], h.SizeCommands)
	o.PutUint32(b[24:], uint32(h.Flags))
	if h.Magic == Magic32 {
		return 28
	}
	o.PutUint32(b[28:], 0)
	return 32
}
func (h *FileHeader) Write(buf *bytes.Buffer, o binary.ByteOrder) error {
	if err := binary.Write(buf, o, h); err != nil {
		return fmt.Errorf("failed to write file header: %v", err)
	}
	return nil
}
func (h *FileHeader) String() string {
	var caps string
	if len(h.SubCPU.Caps(h.CPU)) > 0 {
		caps = fmt.Sprintf(" caps: %s", h.SubCPU.Caps(h.CPU))
	}
	return fmt.Sprintf(
		"Magic         = %s\n"+
			"Type          = %s\n"+
			"CPU           = %s, %s%s\n"+
			"Commands      = %d (Size: %d)\n"+
			"Flags         = %s",
		h.Magic,
		h.Type,
		h.CPU, h.SubCPU.String(h.CPU),
		caps,
		h.NCommands,
		h.SizeCommands,
		h.Flags,
	)
}
func (h *FileHeader) Print(printer func(h *FileHeader) string) string {
	return printer(h)
}
func (h *FileHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Magic        string   `json:"magic"`
		Type         string   `json:"type"`
		CPU          string   `json:"cpu"`
		SubCPU       string   `json:"subcpu"`
		SubCPUCaps   string   `json:"subcpu_caps"`
		Commands     int      `json:"commands"`
		SizeCommands int      `json:"commands_size"`
		Flags        []string `json:"flags"`
	}{
		Magic:        h.Magic.String(),
		Type:         h.Type.String(),
		CPU:          h.CPU.String(),
		SubCPU:       h.SubCPU.String(h.CPU),
		SubCPUCaps:   h.SubCPU.Caps(h.CPU),
		Commands:     int(h.NCommands),
		SizeCommands: int(h.SizeCommands),
		Flags:        h.Flags.Flags(),
	})
}

// A HeaderFileType is the Mach-O file type, e.g. an object file, executable, or dynamic library.
type HeaderFileType uint32

const (
	MH_OBJECT      HeaderFileType = 0x1 /* relocatable object file */
	MH_EXECUTE     HeaderFileType = 0x2 /* demand paged executable file */
	MH_FVMLIB      HeaderFileType = 0x3 /* fixed VM shared library file */
	MH_CORE        HeaderFileType = 0x4 /* core file */
	MH_PRELOAD     HeaderFileType = 0x5 /* preloaded executable file */
	MH_DYLIB       HeaderFileType = 0x6 /* dynamically bound shared library */
	MH_DYLINKER    HeaderFileType = 0x7 /* dynamic link editor */
	MH_BUNDLE      HeaderFileType = 0x8 /* dynamically bound bundle file */
	MH_DYLIB_STUB  HeaderFileType = 0x9 /* shared library stub for static linking only, no section contents */
	MH_DSYM        HeaderFileType = 0xa /* companion file with only debug sections */
	MH_KEXT_BUNDLE HeaderFileType = 0xb /* x86_64 kexts */
	MH_FILESET     HeaderFileType = 0xc /* a file composed of other Mach-Os to be run in the same userspace sharing a single linkedit. */
	MH_GPU_EXECUTE HeaderFileType = 0xd /* gpu program */
	MH_GPU_DYLIB   HeaderFileType = 0xe /* gpu support functions */
)

type HeaderFlag uint32

const (
	None                       HeaderFlag = 0x0
	NoUndefs                   HeaderFlag = 0x1
	IncrLink                   HeaderFlag = 0x2
	DyldLink                   HeaderFlag = 0x4
	BindAtLoad                 HeaderFlag = 0x8
	Prebound                   HeaderFlag = 0x10
	SplitSegs                  HeaderFlag = 0x20
	LazyInit                   HeaderFlag = 0x40
	TwoLevel                   HeaderFlag = 0x80
	ForceFlat                  HeaderFlag = 0x100
	NoMultiDefs                HeaderFlag = 0x200
	NoFixPrebinding            HeaderFlag = 0x400
	Prebindable                HeaderFlag = 0x800
	AllModsBound               HeaderFlag = 0x1000
	SubsectionsViaSymbols      HeaderFlag = 0x2000
	Canonical                  HeaderFlag = 0x4000
	WeakDefines                HeaderFlag = 0x8000
	BindsToWeak                HeaderFlag = 0x10000
	AllowStackExecution        HeaderFlag = 0x20000
	RootSafe                   HeaderFlag = 0x40000
	SetuidSafe                 HeaderFlag = 0x80000
	NoReexportedDylibs         HeaderFlag = 0x100000
	PIE                        HeaderFlag = 0x200000
	DeadStrippableDylib        HeaderFlag = 0x400000
	HasTLVDescriptors          HeaderFlag = 0x800000
	NoHeapExecution            HeaderFlag = 0x1000000
	AppExtensionSafe           HeaderFlag = 0x2000000
	NlistOutofsyncWithDyldinfo HeaderFlag = 0x4000000
	SimSupport                 HeaderFlag = 0x8000000
	DylibInCache               HeaderFlag = 0x80000000
)

func (f HeaderFlag) None() bool {
	return f == 0
}
func (f HeaderFlag) NoUndefs() bool {
	return (f & NoUndefs) != 0
}
func (f HeaderFlag) IncrLink() bool {
	return (f & IncrLink) != 0
}
func (f HeaderFlag) DyldLink() bool {
	return (f & DyldLink) != 0
}
func (f HeaderFlag) BindAtLoad() bool {
	return (f & BindAtLoad) != 0
}
func (f HeaderFlag) Prebound() bool {
	return (f & Prebound) != 0
}
func (f HeaderFlag) SplitSegs() bool {
	return (f & SplitSegs) != 0
}
func (f HeaderFlag) LazyInit() bool {
	return (f & LazyInit) != 0
}
func (f HeaderFlag) TwoLevel() bool {
	return (f & TwoLevel) != 0
}
func (f HeaderFlag) ForceFlat() bool {
	return (f & ForceFlat) != 0
}
func (f HeaderFlag) NoMultiDefs() bool {
	return (f & NoMultiDefs) != 0
}
func (f HeaderFlag) NoFixPrebinding() bool {
	return (f & NoFixPrebinding) != 0
}
func (f HeaderFlag) Prebindable() bool {
	return (f & Prebindable) != 0
}
func (f HeaderFlag) AllModsBound() bool {
	return (f & AllModsBound) != 0
}
func (f HeaderFlag) SubsectionsViaSymbols() bool {
	return (f & SubsectionsViaSymbols) != 0
}
func (f HeaderFlag) Canonical() bool {
	return (f & Canonical) != 0
}
func (f HeaderFlag) WeakDefines() bool {
	return (f & WeakDefines) != 0
}
func (f HeaderFlag) BindsToWeak() bool {
	return (f & BindsToWeak) != 0
}
func (f HeaderFlag) AllowStackExecution() bool {
	return (f & AllowStackExecution) != 0
}
func (f HeaderFlag) RootSafe() bool {
	return (f & RootSafe) != 0
}
func (f HeaderFlag) SetuidSafe() bool {
	return (f & SetuidSafe) != 0
}
func (f HeaderFlag) NoReexportedDylibs() bool {
	return (f & NoReexportedDylibs) != 0
}
func (f HeaderFlag) PIE() bool {
	return (f & PIE) != 0
}
func (f HeaderFlag) DeadStrippableDylib() bool {
	return (f & DeadStrippableDylib) != 0
}
func (f HeaderFlag) HasTLVDescriptors() bool {
	return (f & HasTLVDescriptors) != 0
}
func (f HeaderFlag) NoHeapExecution() bool {
	return (f & NoHeapExecution) != 0
}
func (f HeaderFlag) AppExtensionSafe() bool {
	return (f & AppExtensionSafe) != 0
}
func (f HeaderFlag) NlistOutofsyncWithDyldinfo() bool {
	return (f & NlistOutofsyncWithDyldinfo) != 0
}
func (f HeaderFlag) SimSupport() bool {
	return (f & SimSupport) != 0
}
func (f HeaderFlag) DylibInCache() bool {
	return (f & DylibInCache) != 0
}

func (f *HeaderFlag) Set(flag HeaderFlag, set bool) {
	if set {
		*f = (*f | flag)
	} else {
		*f = (*f ^ flag)
	}
}

func (f HeaderFlag) Flags() []string {
	var flags []string
	if f.None() {
		flags = append(flags, "None")
	}
	if f.NoUndefs() {
		flags = append(flags, "NoUndefs")
	}
	if f.IncrLink() {
		flags = append(flags, "IncrLink")
	}
	if f.DyldLink() {
		flags = append(flags, "DyldLink")
	}
	if f.BindAtLoad() {
		flags = append(flags, "BindAtLoad")
	}
	if f.Prebound() {
		flags = append(flags, "Prebound")
	}
	if f.SplitSegs() {
		flags = append(flags, "SplitSegs")
	}
	if f.LazyInit() {
		flags = append(flags, "LazyInit")
	}
	if f.TwoLevel() {
		flags = append(flags, "TwoLevel")
	}
	if f.ForceFlat() {
		flags = append(flags, "ForceFlat")
	}
	if f.NoMultiDefs() {
		flags = append(flags, "NoMultiDefs")
	}
	if f.NoFixPrebinding() {
		flags = append(flags, "NoFixPrebinding")
	}
	if f.Prebindable() {
		flags = append(flags, "Prebindable")
	}
	if f.AllModsBound() {
		flags = append(flags, "AllModsBound")
	}
	if f.SubsectionsViaSymbols() {
		flags = append(flags, "SubsectionsViaSymbols")
	}
	if f.Canonical() {
		flags = append(flags, "Canonical")
	}
	if f.WeakDefines() {
		flags = append(flags, "WeakDefines")
	}
	if f.BindsToWeak() {
		flags = append(flags, "BindsToWeak")
	}
	if f.AllowStackExecution() {
		flags = append(flags, "AllowStackExecution")
	}
	if f.RootSafe() {
		flags = append(flags, "RootSafe")
	}
	if f.SetuidSafe() {
		flags = append(flags, "SetuidSafe")
	}
	if f.NoReexportedDylibs() {
		flags = append(flags, "NoReexportedDylibs")
	}
	if f.PIE() {
		flags = append(flags, "PIE")
	}
	if f.DeadStrippableDylib() {
		flags = append(flags, "DeadStrippableDylib")
	}
	if f.HasTLVDescriptors() {
		flags = append(flags, "HasTLVDescriptors")
	}
	if f.NoHeapExecution() {
		flags = append(flags, "NoHeapExecution")
	}
	if f.AppExtensionSafe() {
		flags = append(flags, "AppExtensionSafe")
	}
	if f.NlistOutofsyncWithDyldinfo() {
		flags = append(flags, "NlistOutofsyncWithDyldinfo")
	}
	if f.SimSupport() {
		flags = append(flags, "SimSupport")
	}
	if f.DylibInCache() {
		flags = append(flags, "DylibInCache")
	}
	return flags
}

func (f HeaderFlag) String() string {
	return strings.Join(f.Flags(), ", ")
}
