package fixupchains

import (
	"fmt"

	"github.com/blacktop/go-macho/types"
)

// ImportFormat are values for dyld_chained_fixups_header.imports_format
type ImportFormat uint32

const (
	DC_IMPORT          ImportFormat = 1
	DC_IMPORT_ADDEND   ImportFormat = 2
	DC_IMPORT_ADDEND64 ImportFormat = 3
)

type Import interface {
	LibOrdinal() int
	WeakImport() bool
	NameOffset() uint64
	Addend() uint64
}

type DcfImport struct {
	Name string
	Import
}

func (i DcfImport) String() string {
	return fmt.Sprintf("%s, %s", i.Import, i.Name)
}

// DYLD_CHAINED_IMPORT
type DyldChainedImport uint32

func (d DyldChainedImport) LibOrdinal() int {
	return int(int8(types.ExtractBits(uint64(d), 0, 8)))
}
func (d DyldChainedImport) WeakImport() bool {
	return types.ExtractBits(uint64(d), 8, 1) == 1
}
func (d DyldChainedImport) NameOffset() uint64 {
	return types.ExtractBits(uint64(d), 9, 23)
}
func (d DyldChainedImport) Addend() uint64 {
	return 0
}
func (i DyldChainedImport) String() string {
	return fmt.Sprintf("lib ordinal: %2d, is_weak: %t", i.LibOrdinal(), i.WeakImport())
}

// DYLD_CHAINED_IMPORT_ADDEND
type DyldChainedImportAddend struct {
	Import    DyldChainedImport
	AddendVal int32
}

func (d DyldChainedImportAddend) LibOrdinal() int {
	return d.Import.LibOrdinal()
}
func (d DyldChainedImportAddend) WeakImport() bool {
	return d.Import.WeakImport()
}
func (d DyldChainedImportAddend) NameOffset() uint64 {
	return d.Import.NameOffset()
}
func (d DyldChainedImportAddend) Addend() uint64 {
	return uint64(d.AddendVal)
}

func (i DyldChainedImportAddend) String() string {
	return fmt.Sprintf("lib ordinal: %2d, is_weak: %t, addend: %#x", i.LibOrdinal(), i.WeakImport(), i.Addend())
}

type DyldChainedImport64 uint64

func (d DyldChainedImport64) LibOrdinal() int {
	return int(int16(types.ExtractBits(uint64(d), 0, 16)))
}
func (d DyldChainedImport64) WeakImport() bool {
	return types.ExtractBits(uint64(d), 16, 1) == 1
}
func (d DyldChainedImport64) Reserved() uint64 {
	return types.ExtractBits(uint64(d), 17, 15)
}
func (d DyldChainedImport64) NameOffset() uint64 {
	return types.ExtractBits(uint64(d), 32, 32)
}
func (d DyldChainedImport64) Addend() uint64 {
	return 0
}
func (i DyldChainedImport64) String() string {
	return fmt.Sprintf("lib ordinal: %2d, is_weak: %t", i.LibOrdinal(), i.WeakImport())
}

// DYLD_CHAINED_IMPORT_ADDEND64
type DyldChainedImportAddend64 struct {
	Import    DyldChainedImport64
	AddendVal uint64
}

func (d DyldChainedImportAddend64) LibOrdinal() int {
	return d.Import.LibOrdinal()
}
func (d DyldChainedImportAddend64) WeakImport() bool {
	return d.Import.WeakImport()
}
func (d DyldChainedImportAddend64) NameOffset() uint64 {
	return d.Import.NameOffset()
}
func (d DyldChainedImportAddend64) Addend() uint64 {
	return d.AddendVal
}
func (i DyldChainedImportAddend64) String() string {
	return fmt.Sprintf("lib ordinal: %2d, is_weak: %t, addend: %#x", i.LibOrdinal(), i.WeakImport(), i.Addend())
}
