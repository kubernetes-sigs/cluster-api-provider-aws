package fixupchains

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/blacktop/go-macho/types"
)

// NewChainedFixups creates a new DyldChainedFixups instance
func NewChainedFixups(lcdat *bytes.Reader, sr *types.MachoReader, bo binary.ByteOrder) *DyldChainedFixups {
	return &DyldChainedFixups{
		r:  lcdat,
		sr: *sr,
		bo: bo,
	}
}

// Parse parses a LC_DYLD_CHAINED_FIXUPS load command
func (dcf *DyldChainedFixups) Parse() (*DyldChainedFixups, error) {

	if dcf.Starts == nil {
		if err := dcf.ParseStarts(); err != nil {
			return nil, err
		}
	}

	// Parse Imports
	dcf.parseImports()

	for segIdx, start := range dcf.Starts {

		if start.PageStarts == nil {
			continue
		}

		for pageIndex := uint16(0); pageIndex < start.DyldChainedStartsInSegment.PageCount; pageIndex++ {
			offsetInPage := start.PageStarts[pageIndex]

			if offsetInPage == DYLD_CHAINED_PTR_START_NONE {
				continue
			}

			if offsetInPage&DYLD_CHAINED_PTR_START_MULTI != 0 {
				// 32-bit chains which may need multiple starts per page
				overflowIndex := offsetInPage & ^DYLD_CHAINED_PTR_START_MULTI
				chainEnd := false
				for !chainEnd {
					chainEnd = (start.PageStarts[overflowIndex]&DYLD_CHAINED_PTR_START_LAST != 0)
					offsetInPage = (start.PageStarts[overflowIndex] & ^DYLD_CHAINED_PTR_START_LAST)
					if err := dcf.walkDcFixupChain(segIdx, pageIndex, offsetInPage); err != nil {
						return nil, err
					}
					overflowIndex++
				}

			} else {
				// one chain per page
				if err := dcf.walkDcFixupChain(segIdx, pageIndex, offsetInPage); err != nil {
					return nil, err
				}
			}
		}
	}

	return dcf, nil
}

// ParseStarts parses the DyldChainedStartsInSegment(s)
func (dcf *DyldChainedFixups) ParseStarts() error {

	if err := binary.Read(dcf.r, dcf.bo, &dcf.DyldChainedFixupsHeader); err != nil {
		return err
	}

	dcf.r.Seek(int64(dcf.DyldChainedFixupsHeader.StartsOffset), io.SeekStart)

	var segCount uint32
	if err := binary.Read(dcf.r, dcf.bo, &segCount); err != nil {
		return err
	}

	dcf.Starts = make([]DyldChainedStarts, segCount)
	segInfoOffsets := make([]uint32, segCount)
	if err := binary.Read(dcf.r, dcf.bo, &segInfoOffsets); err != nil {
		return err
	}

	for segIdx, segInfoOffset := range segInfoOffsets {
		if segInfoOffset == 0 {
			continue
		}

		dcf.r.Seek(int64(dcf.DyldChainedFixupsHeader.StartsOffset+segInfoOffset), io.SeekStart)
		if err := binary.Read(dcf.r, dcf.bo, &dcf.Starts[segIdx].DyldChainedStartsInSegment); err != nil {
			return err
		}

		dcf.Starts[segIdx].PageStarts = make([]DCPtrStart, dcf.Starts[segIdx].DyldChainedStartsInSegment.PageCount)
		if err := binary.Read(dcf.r, dcf.bo, &dcf.Starts[segIdx].PageStarts); err != nil {
			return err
		}

		dcf.PointerFormat = dcf.Starts[segIdx].DyldChainedStartsInSegment.PointerFormat
	}

	return nil
}

func (dcf *DyldChainedFixups) walkDcFixupChain(segIdx int, pageIndex uint16, offsetInPage DCPtrStart) error {

	var dcPtr uint32
	var dcPtr64 uint64
	var next uint64

	chainEnd := false
	segOffset := dcf.Starts[segIdx].DyldChainedStartsInSegment.SegmentOffset
	pageContentStart := segOffset + uint64(pageIndex)*uint64(dcf.Starts[segIdx].DyldChainedStartsInSegment.PageSize)

	for !chainEnd {
		fixupLocation := pageContentStart + uint64(offsetInPage) + next
		dcf.sr.Seek(int64(fixupLocation), io.SeekStart)

		pointerFormat := dcf.Starts[segIdx].DyldChainedStartsInSegment.PointerFormat

		switch pointerFormat {
		case DYLD_CHAINED_PTR_32:
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr); err != nil {
				return err
			}
			if Generic32IsBind(dcPtr) {
				bind := DyldChainedPtr32Bind{Pointer: dcPtr, Fixup: fixupLocation}
				bind.Import = dcf.Imports[bind.Ordinal()].Name
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, bind)
			} else {
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtr32Rebase{
					Pointer: dcPtr,
					Fixup:   fixupLocation,
				})
			}
			if Generic32Next(dcPtr) == 0 {
				chainEnd = true
			}
			next += Generic32Next(dcPtr) * stride(pointerFormat)
		case DYLD_CHAINED_PTR_32_CACHE:
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr); err != nil {
				return err
			}
			dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtr32CacheRebase{
				Pointer: dcPtr,
				Fixup:   fixupLocation,
			})
			if Generic32Next(dcPtr) == 0 {
				chainEnd = true
			}
			next += Generic32Next(dcPtr) * stride(pointerFormat)
		case DYLD_CHAINED_PTR_32_FIRMWARE:
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr); err != nil {
				return err
			}
			dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtr32FirmwareRebase{
				Pointer: dcPtr,
				Fixup:   fixupLocation,
			})
			if Generic32Next(dcPtr) == 0 {
				chainEnd = true
			}
			next += Generic32Next(dcPtr) * stride(pointerFormat)
		case DYLD_CHAINED_PTR_64: // target is vmaddr
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr64); err != nil {
				return err
			}
			if Generic64IsBind(dcPtr64) {
				bind := DyldChainedPtr64Bind{Pointer: dcPtr64, Fixup: fixupLocation}
				bind.Import = dcf.Imports[bind.Ordinal()].Name
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, bind)
			} else {
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtr64Rebase{
					Pointer: dcPtr64,
					Fixup:   fixupLocation,
				})
			}
			if Generic64Next(dcPtr64) == 0 {
				chainEnd = true
			}
			next += Generic64Next(dcPtr64) * stride(pointerFormat)
		case DYLD_CHAINED_PTR_64_OFFSET: // target is vm offset
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr64); err != nil {
				return err
			}
			dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtr64RebaseOffset{
				Pointer: dcPtr64,
				Fixup:   fixupLocation,
			})
			if Generic64Next(dcPtr64) == 0 {
				chainEnd = true
			}
			next += Generic64Next(dcPtr64) * stride(pointerFormat)
		case DYLD_CHAINED_PTR_64_KERNEL_CACHE:
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr64); err != nil {
				return err
			}
			dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtr64KernelCacheRebase{
				Pointer: dcPtr64,
				Fixup:   fixupLocation,
			})
			if Generic64Next(dcPtr64) == 0 {
				chainEnd = true
			}
			next += Generic64Next(dcPtr64) * stride(pointerFormat)
		case DYLD_CHAINED_PTR_X86_64_KERNEL_CACHE: // stride 1, x86_64 kernel caches
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr64); err != nil {
				return err
			}
			dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtr64KernelCacheRebase{
				Pointer: dcPtr64,
				Fixup:   fixupLocation,
			})
			if Generic64Next(dcPtr64) == 0 {
				chainEnd = true
			}
			next += Generic64Next(dcPtr64) * stride(pointerFormat)
		case DYLD_CHAINED_PTR_ARM64E_KERNEL: // stride 4, unauth target is vm offset
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr64); err != nil {
				return err
			}
			if !DcpArm64eIsBind(dcPtr64) && !DcpArm64eIsAuth(dcPtr64) {
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtrArm64eRebase{
					Pointer: dcPtr64,
					Fixup:   fixupLocation,
				})
			} else if DcpArm64eIsBind(dcPtr64) && !DcpArm64eIsAuth(dcPtr64) {
				bind := DyldChainedPtrArm64eBind{Pointer: dcPtr64, Fixup: fixupLocation}
				bind.Import = dcf.Imports[bind.Ordinal()].Name
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, bind)
			} else if !DcpArm64eIsBind(dcPtr64) && DcpArm64eIsAuth(dcPtr64) {
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtrArm64eAuthRebase{
					Pointer: dcPtr64,
					Fixup:   fixupLocation,
				})
			} else {
				bind := DyldChainedPtrArm64eAuthBind{Pointer: dcPtr64, Fixup: fixupLocation}
				bind.Import = dcf.Imports[bind.Ordinal()].Name
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, bind)
			}
			if DcpArm64eNext(dcPtr64) == 0 {
				chainEnd = true
			}
			next += DcpArm64eNext(dcPtr64) * stride(pointerFormat)
		case DYLD_CHAINED_PTR_ARM64E_FIRMWARE: // stride 4, unauth target is vmaddr
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr64); err != nil {
				return err
			}
			if !DcpArm64eIsBind(dcPtr64) && !DcpArm64eIsAuth(dcPtr64) {
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtrArm64eRebase{
					Pointer: dcPtr64,
					Fixup:   fixupLocation,
				})
			} else if DcpArm64eIsBind(dcPtr64) && !DcpArm64eIsAuth(dcPtr64) {
				bind := DyldChainedPtrArm64eBind{Pointer: dcPtr64, Fixup: fixupLocation}
				bind.Import = dcf.Imports[bind.Ordinal()].Name
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, bind)
			} else if !DcpArm64eIsBind(dcPtr64) && DcpArm64eIsAuth(dcPtr64) {
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtrArm64eAuthRebase{
					Pointer: dcPtr64,
					Fixup:   fixupLocation,
				})
			} else {
				bind := DyldChainedPtrArm64eAuthBind{Pointer: dcPtr64, Fixup: fixupLocation}
				bind.Import = dcf.Imports[bind.Ordinal()].Name
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, bind)
			}
			if DcpArm64eNext(dcPtr64) == 0 {
				chainEnd = true
			}
			next += DcpArm64eNext(dcPtr64) * stride(pointerFormat)
		case DYLD_CHAINED_PTR_ARM64E: // stride 8, unauth target is vmaddr
			fallthrough
		case DYLD_CHAINED_PTR_ARM64E_USERLAND: // stride 8, unauth target is vm offset
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr64); err != nil {
				return err
			}
			if !DcpArm64eIsBind(dcPtr64) && !DcpArm64eIsAuth(dcPtr64) {
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtrArm64eRebase{
					Pointer: dcPtr64,
					Fixup:   fixupLocation,
				})
			} else if DcpArm64eIsBind(dcPtr64) && !DcpArm64eIsAuth(dcPtr64) {
				bind := DyldChainedPtrArm64eBind{Pointer: dcPtr64, Fixup: fixupLocation}
				bind.Import = dcf.Imports[bind.Ordinal()].Name
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, bind)
			} else if !DcpArm64eIsBind(dcPtr64) && DcpArm64eIsAuth(dcPtr64) {
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtrArm64eAuthRebase{
					Pointer: dcPtr64,
					Fixup:   fixupLocation,
				})
			} else {
				bind := DyldChainedPtrArm64eAuthBind{Pointer: dcPtr64, Fixup: fixupLocation}
				bind.Import = dcf.Imports[bind.Ordinal()].Name
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, bind)
			}
			if DcpArm64eNext(dcPtr64) == 0 {
				chainEnd = true
			}
			next += DcpArm64eNext(dcPtr64) * stride(pointerFormat)
		case DYLD_CHAINED_PTR_ARM64E_USERLAND24: // stride 8, unauth target is vm offset, 24-bit bind
			if err := binary.Read(dcf.sr, dcf.bo, &dcPtr64); err != nil {
				return err
			}
			if !DcpArm64eIsBind(dcPtr64) && !DcpArm64eIsAuth(dcPtr64) {
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtrArm64eRebase24{
					Pointer: dcPtr64,
					Fixup:   fixupLocation,
				})
			} else if DcpArm64eIsBind(dcPtr64) && DcpArm64eIsAuth(dcPtr64) {
				bind := DyldChainedPtrArm64eAuthBind24{Pointer: dcPtr64, Fixup: fixupLocation}
				bind.Import = dcf.Imports[bind.Ordinal()].Name
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, bind)
			} else if !DcpArm64eIsBind(dcPtr64) && DcpArm64eIsAuth(dcPtr64) {
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, DyldChainedPtrArm64eAuthRebase24{
					Pointer: dcPtr64,
					Fixup:   fixupLocation,
				})
			} else if DcpArm64eIsBind(dcPtr64) && !DcpArm64eIsAuth(dcPtr64) {
				bind := DyldChainedPtrArm64eBind24{Pointer: dcPtr64, Fixup: fixupLocation}
				bind.Import = dcf.Imports[bind.Ordinal()].Name
				dcf.Starts[segIdx].Fixups = append(dcf.Starts[segIdx].Fixups, bind)
			}
			if DcpArm64eNext(dcPtr64) == 0 {
				chainEnd = true
			}
			next += DcpArm64eNext(dcPtr64) * stride(pointerFormat)
		default:
			return fmt.Errorf("unknown pointer format %#04X", dcf.Starts[segIdx].DyldChainedStartsInSegment.PointerFormat)
		}
	}

	return nil
}

func (dcf *DyldChainedFixups) parseImports() error {

	var imports []Import

	dcf.r.Seek(int64(dcf.ImportsOffset), io.SeekStart)

	switch dcf.DyldChainedFixupsHeader.ImportsFormat {
	case DC_IMPORT:
		ii := make([]DyldChainedImport, dcf.ImportsCount)
		if err := binary.Read(dcf.r, dcf.bo, &ii); err != nil {
			return err
		}
		for _, i := range ii {
			imports = append(imports, i)
		}
	case DC_IMPORT_ADDEND:
		ii := make([]DyldChainedImportAddend, dcf.ImportsCount)
		if err := binary.Read(dcf.r, dcf.bo, &ii); err != nil {
			return err
		}
		for _, i := range ii {
			imports = append(imports, i)
		}
	case DC_IMPORT_ADDEND64:
		ii := make([]DyldChainedImportAddend64, dcf.ImportsCount)
		if err := binary.Read(dcf.r, dcf.bo, &ii); err != nil {
			return err
		}
		for _, i := range ii {
			imports = append(imports, i)
		}
	}

	symbolsPool := io.NewSectionReader(dcf.r, int64(dcf.SymbolsOffset), dcf.r.Size()-int64(dcf.SymbolsOffset))
	for _, i := range imports {
		symbolsPool.Seek(int64(i.NameOffset()), io.SeekStart)
		s, err := bufio.NewReader(symbolsPool).ReadString('\x00')
		if err != nil {
			return fmt.Errorf("failed to read string at: %d: %v", uint64(dcf.SymbolsOffset)+i.NameOffset(), err)
		}
		dcf.Imports = append(dcf.Imports, DcfImport{
			Name:   strings.Trim(s, "\x00"),
			Import: i,
		})
	}

	return nil
}

func (dcf *DyldChainedFixups) IsRebase(addr, preferredLoadAddress uint64) (uint64, bool) {
	var targetRuntimeOffset uint64
	switch dcf.PointerFormat {
	case DYLD_CHAINED_PTR_ARM64E:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_USERLAND:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_USERLAND24:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_KERNEL:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_FIRMWARE:
		if DcpArm64eIsBind(addr) {
			return 0, false
		}
		if DcpArm64eIsAuth(addr) {
			return DyldChainedPtrArm64eAuthRebase{Pointer: addr}.Target(), true
		}
		if DcpArm64eIsRebase(addr) {
			targetRuntimeOffset = DyldChainedPtrArm64eRebase{Pointer: addr}.UnpackTarget()
			if (dcf.PointerFormat == DYLD_CHAINED_PTR_ARM64E) || (dcf.PointerFormat == DYLD_CHAINED_PTR_ARM64E_USERLAND24) || (dcf.PointerFormat == DYLD_CHAINED_PTR_ARM64E_FIRMWARE) {
				targetRuntimeOffset -= preferredLoadAddress
			}
			return targetRuntimeOffset, true
		}
		return 0, false
	case DYLD_CHAINED_PTR_64, DYLD_CHAINED_PTR_64_OFFSET:
		if Generic64IsBind(addr) {
			return targetRuntimeOffset, false
		}
		targetRuntimeOffset = DyldChainedPtr64Rebase{Pointer: addr}.UnpackedTarget()
		if dcf.PointerFormat == DYLD_CHAINED_PTR_64 || dcf.PointerFormat == DYLD_CHAINED_PTR_64_OFFSET {
			targetRuntimeOffset -= preferredLoadAddress
		}
		return targetRuntimeOffset, true
	case DYLD_CHAINED_PTR_64_KERNEL_CACHE, DYLD_CHAINED_PTR_X86_64_KERNEL_CACHE:
		targetRuntimeOffset = DyldChainedPtr64KernelCacheRebase{Pointer: addr}.Target()
		return targetRuntimeOffset, true
	case DYLD_CHAINED_PTR_32:
		if Generic32IsBind(uint32(addr)) {
			return targetRuntimeOffset, false
		}
		targetRuntimeOffset = uint64(DyldChainedPtr32Rebase{Pointer: uint32(addr)}.Target()) - preferredLoadAddress
		return targetRuntimeOffset, true
	case DYLD_CHAINED_PTR_32_FIRMWARE:
		targetRuntimeOffset = uint64(DyldChainedPtr32FirmwareRebase{Pointer: uint32(addr)}.Target()) - preferredLoadAddress
		return targetRuntimeOffset, true
	default:
		return 0, false
	}
}

func (dcf *DyldChainedFixups) IsBind(addr uint64) (*DcfImport, int64, bool) {
	if len(dcf.Imports) == 0 {
		return nil, 0, false
	}

	switch dcf.PointerFormat {
	case DYLD_CHAINED_PTR_ARM64E:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_USERLAND:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_USERLAND24:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_KERNEL:
		fallthrough
	case DYLD_CHAINED_PTR_ARM64E_FIRMWARE:
		if !DcpArm64eIsBind(addr) {
			return nil, 0, false
		}
		if DcpArm64eIsAuth(addr) { // is auth-bind
			if dcf.PointerFormat == DYLD_CHAINED_PTR_ARM64E_USERLAND24 {
				return &dcf.Imports[DyldChainedPtrArm64eAuthBind24{Pointer: addr}.Ordinal()], 0, true
			}
			return &dcf.Imports[DyldChainedPtrArm64eAuthBind{Pointer: addr}.Ordinal()], 0, true
		}
		if dcf.PointerFormat == DYLD_CHAINED_PTR_ARM64E_USERLAND24 {
			return &dcf.Imports[DyldChainedPtrArm64eAuthBind24{Pointer: addr}.Ordinal()], DyldChainedPtrArm64eBind{Pointer: addr}.SignExtendedAddend(), true
		}
		return &dcf.Imports[DyldChainedPtrArm64eAuthBind{Pointer: addr}.Ordinal()], DyldChainedPtrArm64eBind{Pointer: addr}.SignExtendedAddend(), true
	case DYLD_CHAINED_PTR_64, DYLD_CHAINED_PTR_64_OFFSET:
		if !Generic64IsBind(addr) {
			return nil, 0, false
		}
		return &dcf.Imports[DyldChainedPtr64Bind{Pointer: addr}.Ordinal()], int64(DyldChainedPtr64Bind{Pointer: addr}.Addend()), true
	case DYLD_CHAINED_PTR_32:
		if !Generic32IsBind(uint32(addr)) {
			return nil, 0, false
		}
		return &dcf.Imports[DyldChainedPtr32Bind{Pointer: uint32(addr)}.Ordinal()], int64(DyldChainedPtr32Bind{Pointer: uint32(addr)}.Addend()), true
	case DYLD_CHAINED_PTR_64_KERNEL_CACHE, DYLD_CHAINED_PTR_X86_64_KERNEL_CACHE:
		return nil, 0, false
	default:
		return nil, 0, false
	}
}
