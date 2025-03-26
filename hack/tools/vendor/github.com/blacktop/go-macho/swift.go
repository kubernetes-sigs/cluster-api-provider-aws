package macho

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/blacktop/go-macho/pkg/fixupchains"
	"github.com/blacktop/go-macho/types/swift"
	"github.com/blacktop/go-macho/types/swift/fields"
	"github.com/blacktop/go-macho/types/swift/types"
)

const sizeOfInt32 = 4
const sizeOfInt64 = 8

var ErrSwiftSectionError = fmt.Errorf("missing swift section")

// GetSwiftProtocols parses all the protocols in the __TEXT.__swift5_protos section
func (f *File) GetSwiftProtocols() ([]types.Protocol, error) {
	var protos []types.Protocol

	if sec := f.Section("__TEXT", "__swift5_protos"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		relOffsets := make([]int32, len(dat)/sizeOfInt32)

		if err := binary.Read(bytes.NewReader(dat), f.ByteOrder, &relOffsets); err != nil {
			return nil, fmt.Errorf("failed to read relative offsets: %v", err)
		}

		for idx, relOff := range relOffsets {
			offset := int64(sec.Offset+uint32(idx*sizeOfInt32)) + int64(relOff)

			f.cr.Seek(offset, io.SeekStart)

			var proto types.Protocol
			if err := binary.Read(f.cr, f.ByteOrder, &proto.Descriptor); err != nil {
				return nil, fmt.Errorf("failed to read protocols descriptor: %v", err)
			}

			if proto.NumRequirementsInSignature > 0 {
				currentOffset, _ := f.cr.Seek(0, io.SeekCurrent)
				proto.SignatureRequirements = make([]types.TargetGenericRequirement, proto.NumRequirementsInSignature)
				for i := 0; i < int(proto.NumRequirementsInSignature); i++ {
					if err := binary.Read(f.cr, f.ByteOrder, &proto.SignatureRequirements[i].TargetGenericRequirementDescriptor); err != nil {
						return nil, fmt.Errorf("failed to read protocols requirements in signature : %v", err)
					}
				}

				for idx, req := range proto.SignatureRequirements {
					proto.SignatureRequirements[idx].Name, err = f.makeSymbolicMangledNameStringRef(currentOffset + 4 + int64(req.Param))
					if err != nil {
						return nil, fmt.Errorf("failed to read protocols requirements in signature : %v", err)
					}
					switch req.Flags.Kind() {
					case types.GRKindProtocol:
						off := currentOffset + 8 + int64(req.TypeOrProtocolOrConformanceOrLayout)
						var ptr uint64
						if (off & 1) == 1 {
							off = off &^ 1
							ptr, _ = f.GetPointer(uint64(off))
						} else {
							ptr = uint64(offset + int64(off))
						}
						if fixupchains.Generic64IsBind(ptr) {
							proto.SignatureRequirements[idx].Kind, err = f.GetBindName(ptr)
							if err != nil {
								return nil, fmt.Errorf("failed to read protocol name: %v", err)
							}
						} else {
							proto.SignatureRequirements[idx].Kind, err = f.GetCString(f.vma.Convert(ptr))
							if err != nil {
								return nil, fmt.Errorf("failed to read protocol name: %v", err)
							}
						}
					case types.GRKindSameType:
						fmt.Println("same type")
					case types.GRKindSameConformance:
						fmt.Println("same conformance")
					case types.GRKindLayout:
						fmt.Println("layout")
					}
					fmt.Printf("%s (%s): %s\n", proto.SignatureRequirements[idx].Name, proto.SignatureRequirements[idx].Kind, req.Flags)
				}

				f.cr.Seek(currentOffset+int64(binary.Size(types.TargetGenericRequirementDescriptor{})*int(proto.NumRequirementsInSignature)), io.SeekStart)
			}

			if proto.NumRequirements > 0 {
				proto.Requirements = make([]types.TargetProtocolRequirement, proto.NumRequirements)
				if err := binary.Read(f.cr, f.ByteOrder, &proto.Requirements); err != nil {
					return nil, fmt.Errorf("failed to read protocols requirements: %v", err)
				}
				for _, req := range proto.Requirements {
					addr, err := f.GetOffset(uint64(req.DefaultImplementation))
					if err != nil {
						return nil, fmt.Errorf("failed to read protocols requirements: %v", err)
					}
					fmt.Printf("%#x: flags: %s\n", addr, req.Flags)
				}
			}

			proto.Name, err = f.GetCStringAtOffset(offset + int64(sizeOfInt32*2) + int64(proto.NameOffset))
			if err != nil {
				return nil, fmt.Errorf("failed to read cstring: %v", err)
			}

			if proto.ParentOffset != 0 {
				parentOffset := offset + 4 + int64(proto.Descriptor.ParentOffset)
				f.cr.Seek(parentOffset, io.SeekStart)
				parentAddr, err := f.GetVMAddress(uint64(parentOffset))
				if err != nil {
					return nil, fmt.Errorf("failed to read protocols parent address: %v", err)
				}
				_ = parentAddr

				var pcDesc types.TargetModuleContextDescriptor
				if err := binary.Read(f.cr, f.ByteOrder, &pcDesc); err != nil {
					return nil, fmt.Errorf("failed to read protocols parent descriptor: %v", err)
				}

				proto.Parent, err = f.GetCStringAtOffset(parentOffset + int64(sizeOfInt32*2) + int64(pcDesc.NameOffset))
				if err != nil {
					return nil, fmt.Errorf("failed to read protocols parent name: %v", err)
				}

				if pcDesc.ParentOffset != 0 { // TODO: what if parent has parent ?
					fmt.Printf("found a grand parent while parsing %s", proto.Parent) // FIXME: if this happens this should be recursive
				}
			}

			if proto.AssociatedTypeNamesOffset != 0 {
				proto.AssociatedType, err = f.GetCStringAtOffset(offset + int64(sizeOfInt32*5) + int64(proto.AssociatedTypeNamesOffset))
				if err != nil {
					return nil, fmt.Errorf("failed to read protocols assocated type names: %v", err)
				}
			}

			protos = append(protos, proto)
		}

		return protos, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_protos' section: %w", ErrSwiftSectionError)
}

// GetSwiftProtocolConformances parses all the protocol conformances in the __TEXT.__swift5_proto section
func (f *File) GetSwiftProtocolConformances() ([]types.ConformanceDescriptor, error) {
	var protoConfDescs []types.ConformanceDescriptor

	if sec := f.Section("__TEXT", "__swift5_proto"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		relOffsets := make([]int32, len(dat)/sizeOfInt32)

		if err := binary.Read(bytes.NewReader(dat), f.ByteOrder, &relOffsets); err != nil {
			return nil, fmt.Errorf("failed to read relative offsets: %v", err)
		}

		for idx, relOff := range relOffsets {
			offset := int64(sec.Offset+uint32(idx*sizeOfInt32)) + int64(relOff)

			f.cr.Seek(offset, io.SeekStart)

			var pcd types.ConformanceDescriptor
			if err := binary.Read(f.cr, f.ByteOrder, &pcd.TargetProtocolConformanceDescriptor); err != nil {
				return nil, fmt.Errorf("failed to read swift ProtocolDescriptor: %v", err)
			}

			pcd.Address = uint64(int64(sec.Addr+uint64(idx*sizeOfInt32)) + int64(relOff))

			var ptr uint64
			if (pcd.ProtocolOffsest & 1) == 1 {
				pcd.ProtocolOffsest = pcd.ProtocolOffsest &^ 1
				f.cr.Seek(offset+int64(pcd.ProtocolOffsest), io.SeekStart)
				if err := binary.Read(f.cr, f.ByteOrder, &ptr); err != nil {
					return nil, fmt.Errorf("failed to read protocol name offset: %v", err)
				}
			} else {
				ptr = uint64(offset + int64(pcd.ProtocolOffsest))
			}

			if f.HasFixups() {
				dcf, err := f.DyldChainedFixups()
				if err != nil {
					return nil, fmt.Errorf("failed to get dyld chained fixups: %v", err)
				}
				if _, _, ok := dcf.IsBind(ptr); ok {
					pcd.Protocol, err = f.GetBindName(ptr)
					if err != nil {
						return nil, fmt.Errorf("failed to read protocol name: %v", err)
					}
				} else {
					pcd.Protocol, err = f.GetCString(f.vma.Convert(ptr))
					if err != nil {
						return nil, fmt.Errorf("failed to read protocol name: %v", err)
					}
				}
			} else { // TODO: fix this (redundant???)
				pcd.Protocol, err = f.GetCString(f.vma.Convert(ptr))
				if err != nil {
					return nil, fmt.Errorf("failed to read protocol name: %v", err)
				}
			}

			// Parse the TargetTypeReference
			switch pcd.Flags.GetTypeReferenceKind() {
			case types.DirectTypeDescriptor:
				pcd.NominalType, err = f.readType(offset + sizeOfInt32 + int64(pcd.TypeRefOffsest))
				if err != nil {
					return nil, fmt.Errorf("failed to read type: %v", err)
				}
			case types.IndirectTypeDescriptor:
				addr, err := f.GetVMAddress(uint64(int64(offset) + sizeOfInt32 + int64(pcd.TypeRefOffsest)))
				if err != nil {
					return nil, fmt.Errorf("failed to get vmaddr for indirect nominal type descriptor: %v", err)
				}
				ptr, err := f.GetPointerAtAddress(addr)
				if err != nil {
					return nil, fmt.Errorf("failed to read type pointer: %v", err)
				}
				off, err := f.vma.GetOffset(f.vma.Convert(ptr))
				if err != nil {
					return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
				}
				pcd.NominalType, err = f.readType(int64(off))
				if err != nil {
					return nil, fmt.Errorf("failed to read type: %v", err)
				}
			case types.DirectObjCClassName:
				fmt.Println("SUP - DirectObjCClassName")
			case types.IndirectObjCClass:
				fmt.Println("SUP - IndirectObjCClass")
			}

			protoConfDescs = append(protoConfDescs, pcd)
		}

		return protoConfDescs, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_proto' section: %w", ErrSwiftSectionError)
}

// GetSwiftTypes parses all the types in the __TEXT.__swift5_types section
func (f *File) GetSwiftTypes() ([]*types.TypeDescriptor, error) {
	var typs []*types.TypeDescriptor

	if sec := f.Section("__TEXT", "__swift5_types"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		relOffsets := make([]int32, len(dat)/sizeOfInt32)
		if err := binary.Read(bytes.NewReader(dat), f.ByteOrder, &relOffsets); err != nil {
			return nil, fmt.Errorf("failed to read relative offsets: %v", err)
		}

		for idx, relOff := range relOffsets {
			offset := int64(sec.Offset+uint32(idx*sizeOfInt32)) + int64(relOff)

			typ, err := f.readType(offset)
			if err != nil {
				return nil, fmt.Errorf("failed to read type: %v", err)
			}

			typs = append(typs, typ)
		}

		return typs, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_types' section: %w", ErrSwiftSectionError)
}

func (f *File) readType(offset int64) (*types.TypeDescriptor, error) {
	var err error
	var typ types.TypeDescriptor

	f.cr.Seek(offset, io.SeekStart)

	var tDesc types.TargetTypeContextDescriptor
	if err := binary.Read(f.cr, f.ByteOrder, &tDesc); err != nil {
		return nil, fmt.Errorf("failed to read swift type context descriptor: %v", err)
	}

	f.cr.Seek(-int64(binary.Size(tDesc)), io.SeekCurrent) // rewind

	typ.Address, err = f.GetVMAddress(uint64(offset))
	if err != nil {
		return nil, fmt.Errorf("failed to get swift type context descriptor address: %v", err)
	}

	typ.Kind = tDesc.Flags.Kind()

	fmt.Println(tDesc.Flags)

	var metadataInitSize int

	switch tDesc.Flags.KindSpecific().MetadataInitialization() {
	case types.MetadataInitNone:
		metadataInitSize = 0
	case types.MetadataInitSingleton:
		metadataInitSize = binary.Size(types.TargetSingletonMetadataInitialization{})
	case types.MetadataInitForeign:
		metadataInitSize = binary.Size(types.TargetForeignMetadataInitialization{})
	}
	fmt.Println("metadataInitSize: ", metadataInitSize)
	_ = metadataInitSize // TODO: use this in size/offset calculations

	switch typ.Kind {
	case types.CDKindModule:
		var mod types.TargetModuleContextDescriptor
		if err := binary.Read(f.cr, f.ByteOrder, &mod); err != nil {
			return nil, fmt.Errorf("failed to read swift module descriptor: %v", err)
		}
		typ.Type = &mod
	case types.CDKindExtension:
		var ext types.TargetExtensionContextDescriptor
		if err := binary.Read(f.cr, f.ByteOrder, &ext); err != nil {
			return nil, fmt.Errorf("failed to read swift extension descriptor: %v", err)
		}
		typ.Type = &ext
	case types.CDKindAnonymous:
		var anon types.TargetAnonymousContextDescriptor
		if err := binary.Read(f.cr, f.ByteOrder, &anon); err != nil {
			return nil, fmt.Errorf("failed to read swift anonymous descriptor: %v", err)
		}
		typ.Type = &anon
	case types.CDKindClass:
		var cD types.TargetClassDescriptor
		if err := binary.Read(f.cr, f.ByteOrder, &cD); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", cD, err)
		}
		if cD.Flags.IsGeneric() {
			var g types.TargetTypeGenericContextDescriptorHeader
			if err := binary.Read(f.cr, f.ByteOrder, &g); err != nil {
				return nil, fmt.Errorf("failed to read generic header: %v", err)
			}
			typ.Generic = &g
		}
		if cD.FieldOffsetVectorOffset != 0 {
			if cD.Flags.KindSpecific().HasResilientSuperclass() {
				cD.FieldOffsetVectorOffset += cD.MetadataNegativeSizeInWords
			}
			fmt.Printf("FieldOffsetVectorOffset: %d\n", offset+int64(cD.FieldOffsetVectorOffset*8))
			typ.FieldOffsets = make([]int32, cD.NumFields)
			if err := binary.Read(f.cr, f.ByteOrder, &typ.FieldOffsets); err != nil {
				return nil, fmt.Errorf("failed to read field offset vector: %v", err)
			}
		}
		if cD.Flags.KindSpecific().HasVTable() {
			var v types.VTable
			if err := binary.Read(f.cr, f.ByteOrder, &v.TargetVTableDescriptorHeader); err != nil {
				return nil, fmt.Errorf("failed to read vtable header: %v", err)
			}
			v.MethodListOffset, _ = f.cr.Seek(0, io.SeekCurrent)
			v.Methods = make([]types.TargetMethodDescriptor, v.VTableSize)
			if err := binary.Read(f.cr, f.ByteOrder, &v.Methods); err != nil {
				return nil, fmt.Errorf("failed to read vtable method descriptors: %v", err)
			}
			typ.VTable = &v
		}
		typ.Type = &cD
	case types.CDKindEnum:
		var eD types.TargetEnumDescriptor
		if err := binary.Read(f.cr, f.ByteOrder, &eD); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", eD, err)
		}
		if eD.Flags.IsGeneric() {
			var g types.TargetTypeGenericContextDescriptorHeader
			if err := binary.Read(f.cr, f.ByteOrder, &g); err != nil {
				return nil, fmt.Errorf("failed to read generic header: %v", err)
			}
			typ.Generic = &g
		}
		if eD.NumPayloadCasesAndPayloadSizeOffset != 0 {
			fmt.Println("NumPayloadCasesAndPayloadSizeOffset: ", eD.NumPayloadCasesAndPayloadSizeOffset)
		}
		typ.Type = &eD
	case types.CDKindStruct:
		var sD types.TargetStructDescriptor
		if err := binary.Read(f.cr, f.ByteOrder, &sD); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", sD, err)
		}
		current, _ := f.cr.Seek(0, io.SeekCurrent)
		if sD.Flags.IsGeneric() {
			var g types.TargetTypeGenericContextDescriptorHeader
			if err := binary.Read(f.cr, f.ByteOrder, &g); err != nil {
				return nil, fmt.Errorf("failed to read generic header: %v", err)
			}
			typ.Generic = &g
		}
		current, _ = f.cr.Seek(0, io.SeekCurrent)
		if sD.FieldOffsetVectorOffset != 0 {
			typ.FieldOffsets = make([]int32, sD.NumFields)
			if err := binary.Read(f.cr, f.ByteOrder, &typ.FieldOffsets); err != nil {
				return nil, fmt.Errorf("failed to read field offset vector: %v", err)
			}
		}
		current, _ = f.cr.Seek(0, io.SeekCurrent)
		_ = current
		if sD.Flags.KindSpecific().MetadataInitialization() == types.MetadataInitSingleton {
			var md types.TargetSingletonMetadataInitialization
			if err := binary.Read(f.cr, f.ByteOrder, &md); err != nil {
				return nil, fmt.Errorf("failed to read singleton metadata initialization: %v", err)
			}
			fmt.Println(md)
		}
		typ.Type = &sD
	case types.CDKindProtocol:
		var pD types.TargetProtocolDescriptor
		if err := binary.Read(f.cr, f.ByteOrder, &pD); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", pD, err)
		}
		typ.Type = &pD
	case types.CDKindOpaqueType:
		var oD types.TargetOpaqueTypeDescriptor
		if err := binary.Read(f.cr, f.ByteOrder, &oD); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", oD, err)
		}
		typ.Type = &oD
	}

	typ.Name, err = f.GetCStringAtOffset(offset + int64(sizeOfInt32*2) + int64(tDesc.NameOffset))
	if err != nil {
		return nil, fmt.Errorf("failed to read cstring: %v", err)
	}

	if tDesc.ParentOffset < 0 {
		typ.Parent.Address = uint64(int64(typ.Address) + sizeOfInt32 + int64(tDesc.ParentOffset))
		parent, err := f.getParent(offset + sizeOfInt32 + int64(tDesc.ParentOffset))
		if err != nil {
			return nil, fmt.Errorf("failed to get parent: %v", err)
		}
		typ.Parent.Name = parent.Name
	}

	typ.AccessFunction = uint64(int64(typ.Address) + int64(sizeOfInt32*3) + int64(tDesc.AccessFunctionPtr))

	if typ.VTable != nil {
		fmt.Println("METHODS")
		for idx, m := range typ.VTable.GetMethods(f.preferredLoadAddress()) {
			fmt.Printf("%2d)  flags:   %s\n", idx, m.Flags)
			if m.Flags.IsAsync() {
				fmt.Println("ASYNC")
			}
			var sym string
			syms, _ := f.FindAddressSymbols(m.Address)
			if len(syms) > 0 {
				for _, s := range syms {
					if !s.Type.IsDebugSym() {
						sym = s.Name
						break
					}
				}
			}
			// fmt.Printf("      impl:    %d\n", m.Impl)
			fmt.Printf("      address: %#x\nsym: %s\n", m.Address, sym)
		}
	}

	if tDesc.FieldsOffset != 0 {
		offset += int64(sizeOfInt32*4) + int64(tDesc.FieldsOffset)
		fd, err := f.readField(offset, typ.FieldOffsets...)
		if err != nil {
			return nil, fmt.Errorf("failed to read swift field: %v", err)
		}
		typ.Fields = append(typ.Fields, fd)
	}

	return &typ, nil
}

// GetSwiftFields parses all the fields in the __TEXT.__swift5_fields section
func (f *File) GetSwiftFields() ([]*fields.Field, error) {
	var fds []*fields.Field

	if sec := f.Section("__TEXT", "__swift5_fieldmd"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}

		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		r := bytes.NewReader(dat)

		for {
			currentOffset, _ := r.Seek(0, io.SeekCurrent)
			currentOffset += int64(sec.Offset)

			var header fields.FDHeader
			err = binary.Read(r, f.ByteOrder, &header)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("failed to read swift FieldDescriptor header: %v", err)
			}

			field, err := f.readField(currentOffset)
			if err != nil {
				return nil, fmt.Errorf("failed to read field at offset %#x: %v", currentOffset, err)
			}

			r.Seek(int64(uint32(header.FieldRecordSize)*header.NumFields), io.SeekCurrent)

			fds = append(fds, field)
		}

		return fds, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_fieldmd' section: %w", ErrSwiftSectionError)
}

func (f *File) readField(offset int64, fieldOffsets ...int32) (*fields.Field, error) {
	var field fields.Field

	currOffset, err := f.cr.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to swift field offset %#x: %v", offset, err)
	}

	field.Address, err = f.GetVMAddress(uint64(currOffset))
	if err != nil {
		return nil, fmt.Errorf("failed to get swift field address from offset: %v", err)
	}

	if err := binary.Read(f.cr, f.ByteOrder, &field.Descriptor.FDHeader); err != nil {
		return nil, fmt.Errorf("failed to read swift field descriptor header: %v", err)
	}

	field.Kind = field.Descriptor.Kind.String()

	field.MangledType, err = f.makeSymbolicMangledNameStringRef(currOffset + int64(field.Descriptor.MangledTypeNameOffset))
	if err != nil {
		return nil, fmt.Errorf("failed to read swift field mangled type name at %#x: %v", currOffset+int64(field.Descriptor.MangledTypeNameOffset), err)
	}

	if field.Descriptor.SuperclassOffset != 0 {
		field.SuperClass, err = f.makeSymbolicMangledNameStringRef(currOffset + sizeOfInt32 + int64(field.Descriptor.SuperclassOffset))
		if err != nil {
			return nil, fmt.Errorf("failed to read swift field super class mangled name: %v", err)
		}
	}

	currOffset, err = f.cr.Seek(offset+int64(binary.Size(fields.FDHeader{})), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to swift field records offset: %v", err)
	}

	if field.Descriptor.FieldRecordSize != uint16(binary.Size(fields.FieldRecordType{})) {
		return nil, fmt.Errorf("invalid swift field record size: got %d, want %d", field.Descriptor.FieldRecordSize, binary.Size(fields.FieldRecordType{}))
	}

	field.Descriptor.FieldRecords = make([]fields.FieldRecordType, field.Descriptor.NumFields)
	if err := binary.Read(f.cr, f.ByteOrder, &field.Descriptor.FieldRecords); err != nil {
		return nil, fmt.Errorf("failed to read swift field record headers: %v", err)
	}

	for _, record := range field.Descriptor.FieldRecords {
		rec := fields.FieldRecord{
			Flags: record.Flags.String(),
		}

		if record.MangledTypeNameOffset != 0 {
			rec.MangledType, err = f.makeSymbolicMangledNameStringRef(currOffset + sizeOfInt32 + int64(record.MangledTypeNameOffset))
			if err != nil {
				return nil, fmt.Errorf("failed to read swift field record mangled type name at %#x; %v", currOffset+sizeOfInt32+int64(record.MangledTypeNameOffset), err)
			}
		}

		rec.Name, err = f.GetCStringAtOffset(currOffset + int64(sizeOfInt32*2) + int64(record.FieldNameOffset))
		if err != nil {
			return nil, fmt.Errorf("failed to read swift field record name cstring: %v", err)
		}

		field.Records = append(field.Records, rec)
		currOffset += int64(binary.Size(record))
	}

	return &field, nil
}

// GetSwiftAssociatedTypes parses all the associated types in the __TEXT.__swift5_assocty section
func (f *File) GetSwiftAssociatedTypes() ([]swift.AssociatedTypeDescriptor, error) {
	var accocTypes []swift.AssociatedTypeDescriptor

	if sec := f.Section("__TEXT", "__swift5_assocty"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		r := bytes.NewReader(dat)

		for {
			currentOffset, _ := r.Seek(0, io.SeekCurrent)

			var aType swift.AssociatedTypeDescriptor
			err := binary.Read(r, f.ByteOrder, &aType.ATDHeader)

			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("failed to read swift AssociatedTypeDescriptor header: %v", err)
			}

			aType.Address = sec.Addr + uint64(currentOffset)

			off, err := f.GetOffset(aType.Address)
			if err != nil {
				return nil, fmt.Errorf("failed to get offset for associated type at addr %#x: %v", aType.Address, err)
			}

			// AssociatedTypeDescriptor.ConformingTypeName
			aType.ConformingTypeName, err = f.makeSymbolicMangledNameStringRef(int64(off) + int64(aType.ConformingTypeNameOffset))
			if err != nil {
				return nil, fmt.Errorf("failed to read conforming type for associated type at addr %#x: %v", aType.Address, err)
			}

			// AssociatedTypeDescriptor.ProtocolTypeName
			addr := uint64(int64(aType.Address) + sizeOfInt32 + int64(aType.ProtocolTypeNameOffset))
			aType.ProtocolTypeName, err = f.GetCString(addr)
			if err != nil {
				return nil, fmt.Errorf("failed to read swift assocated type protocol type name at addr %#x: %v", addr, err)
			}

			// AssociatedTypeRecord
			aType.AssociatedTypeRecords = make([]swift.AssociatedTypeRecord, aType.ATDHeader.NumAssociatedTypes)
			for i := uint32(0); i < aType.ATDHeader.NumAssociatedTypes; i++ {
				if err := binary.Read(r, f.ByteOrder, &aType.AssociatedTypeRecords[i].ATRecordType); err != nil {
					return nil, fmt.Errorf("failed to read %T: %v", aType.AssociatedTypeRecords[i].ATRecordType, err)
				}
				// AssociatedTypeRecord.Name
				addr := int64(aType.Address) + int64(binary.Size(aType.ATDHeader)) + int64(aType.AssociatedTypeRecords[i].NameOffset)
				aType.AssociatedTypeRecords[i].Name, err = f.GetCString(uint64(addr))
				if err != nil {
					return nil, fmt.Errorf("failed to read associated type record name: %v", err)
				}
				// AssociatedTypeRecord.SubstitutedTypeName
				symMangOff := int64(off) + int64(binary.Size(aType.ATDHeader)) + int64(aType.AssociatedTypeRecords[i].SubstitutedTypeNameOffset) + sizeOfInt32
				aType.AssociatedTypeRecords[i].SubstitutedTypeName, err = f.makeSymbolicMangledNameStringRef(symMangOff)
				if err != nil {
					return nil, fmt.Errorf("failed to read associated type substituted type symbolic ref at offset %#x: %v", symMangOff, err)
				}
			}

			accocTypes = append(accocTypes, aType)
		}

		return accocTypes, nil
	}
	return nil, fmt.Errorf("MachO has no '__swift5_assocty' section: %w", ErrSwiftSectionError)
}

func (f *File) getParent(offset int64) (*types.TargetModuleContext, error) {
	var parent types.TargetModuleContext

	if _, err := f.cr.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to swift context descriptor parent offset: %v", err)
	}
	if err := binary.Read(f.cr, f.ByteOrder, &parent.TargetModuleContextDescriptor); err != nil {
		return nil, fmt.Errorf("failed to read type swift context descriptor parent type context descriptor: %v", err)
	}
	if parent.Flags.Kind() != types.CDKindAnonymous {
		var err error
		parent.Name, err = f.GetCStringAtOffset(offset + int64(sizeOfInt32*2) + int64(parent.NameOffset))
		if err != nil {
			return nil, fmt.Errorf("failed to read swift context descriptor name: %v", err)
		}
	}

	return &parent, nil
}

// ref: https://github.com/apple/swift/blob/1a7146fb04665e2434d02bada06e6296f966770b/lib/Demangling/Demangler.cpp#L155
// ref: https://github.com/apple/swift/blob/main/docs/ABI/Mangling.rst#symbolic-references
func (f *File) makeSymbolicMangledNameStringRef(offset int64) (string, error) {

	var name string
	var symbolic bool

	controlData := make([]byte, 9)
	f.cr.ReadAt(controlData, offset)

	if controlData[0] >= 0x01 && controlData[0] <= 0x17 {
		var reference int32
		if err := binary.Read(bytes.NewReader(controlData[1:]), f.ByteOrder, &reference); err != nil {
			return "", fmt.Errorf("failed to read swift symbolic reference: %v", err)
		}
		symbolic = true
		offset += 1 + int64(reference)
	} else if controlData[0] >= 0x18 && controlData[0] <= 0x1f {
		var reference uint64
		if err := binary.Read(bytes.NewReader(controlData[1:]), f.ByteOrder, &reference); err != nil {
			return "", fmt.Errorf("failed to read swift symbolic reference: %v", err)
		}
		symbolic = true
		offset = int64(reference)
	} else {
		name, err := f.GetCStringAtOffset(offset)
		if err != nil {
			return "", fmt.Errorf("failed to read swift symbolic reference @ %#x: %v", offset, err)
		}
		if strings.HasPrefix(name, "S") {
			return "_$s" + name, nil
		}
		return name, nil
	}

	f.cr.Seek(offset, io.SeekStart)

	switch uint8(controlData[0]) {
	case 1: // Reference points directly to context descriptor
		var err error
		var tDesc types.TargetModuleContextDescriptor
		if err := binary.Read(f.cr, f.ByteOrder, &tDesc); err != nil {
			return "", fmt.Errorf("failed to read swift context descriptor: %v", err)
		}
		name, err = f.GetCStringAtOffset(offset + int64(sizeOfInt32*2) + int64(tDesc.NameOffset))
		if err != nil {
			return "", fmt.Errorf("failed to read swift context descriptor descriptor name: %v", err)
		}
		if tDesc.ParentOffset < 0 {
			parentOffset := offset + sizeOfInt32 + int64(tDesc.ParentOffset)
			for { // walk the family tree
				parent, err := f.getParent(parentOffset)
				if err != nil {
					return "", fmt.Errorf("failed to read swift context descriptor parent: %v", err)
				}
				if len(parent.Name) > 0 {
					name = parent.Name + "." + name
				}
				if parent.ParentOffset >= 0 {
					break
				}
				parentOffset += sizeOfInt32 + int64(parent.ParentOffset)
			}
		}
	case 2: // Reference points indirectly to context descriptor
		addr, err := f.GetVMAddress(uint64(offset))
		if err != nil {
			return "", fmt.Errorf("failed to get vmaddr for indirect context descriptor: %v", err)
		}
		ptr, err := f.GetPointerAtAddress(addr)
		if err != nil {
			return "", fmt.Errorf("failed to get pointer for indirect context descriptor: %v", err)
		}
		if f.HasFixups() {
			dcf, err := f.DyldChainedFixups()
			if err != nil {
				return "", fmt.Errorf("failed to get dyld chained fixups: %v", err)
			}
			if _, _, ok := dcf.IsBind(ptr); ok {
				name, err = f.GetBindName(ptr)
				if err != nil {
					return "", fmt.Errorf("failed to read protocol name: %v", err)
				}
			} else {
				if err := f.cr.SeekToAddr(f.vma.Convert(ptr)); err != nil {
					return "", fmt.Errorf("failed to seek to indirect context descriptor: %v", err)
				}
				var tDesc types.TargetModuleContextDescriptor
				if err := binary.Read(f.cr, f.ByteOrder, &tDesc); err != nil {
					return "", fmt.Errorf("failed to read indirect context descriptor: %v", err)
				}
				name, err = f.GetCString(ptr + uint64(sizeOfInt32*2) + uint64(tDesc.NameOffset))
				if err != nil {
					return "", fmt.Errorf("failed to read indirect context descriptor name: %v", err)
				}
				if tDesc.ParentOffset != 0 {
					parentAddr := f.vma.Convert(ptr) + sizeOfInt32 + uint64(tDesc.ParentOffset)
					if err := f.cr.SeekToAddr(parentAddr); err != nil {
						return "", fmt.Errorf("failed to seek to indirect context descriptor parent: %v", err)
					}
					var parentDesc types.TargetModuleContextDescriptor
					if err := binary.Read(f.cr, f.ByteOrder, &parentDesc); err != nil {
						return "", fmt.Errorf("failed to read type swift indirect context descriptor parent type context descriptor: %v", err)
					}
					parent, err := f.GetCString(parentAddr + uint64(sizeOfInt32*2) + uint64(parentDesc.NameOffset))
					if err != nil {
						return "", fmt.Errorf("failed to read indirect context descriptor name: %v", err)
					}
					if len(parent) > 0 {
						name = parent + "." + name
					}
				}
			}
		} else { // TODO: fix this (redundant???)
			name, err = f.GetCString(f.vma.Convert(ptr))
			if err != nil {
				return "", fmt.Errorf("failed to read protocol name: %v", err)
			}
		}
	case 3: // Reference points directly to protocol conformance descriptor (NOT IMPLEMENTED)
		return "", fmt.Errorf("symbolic reference control character %#x is not implemented", controlData[0])
	case 4: // Reference points indirectly to protocol conformance descriptor (NOT IMPLEMENTED)
		fallthrough
	case 5: // Reference points directly to associated conformance descriptor (NOT IMPLEMENTED)
		fallthrough
	case 6: // Reference points indirectly to associated conformance descriptor (NOT IMPLEMENTED)
		fallthrough
	case 7: // Reference points directly to associated conformance access function relative to the protocol
		fallthrough
	case 8: // Reference points indirectly to associated conformance access function relative to the protocol
		fallthrough
	case 9: // Reference points directly to metadata access function that can be invoked to produce referenced object
		// kind = SymbolicReferenceKind::AccessorFunctionReference; TODO: implement
		// direct = Directness::Direct;
		fallthrough
	case 10: // Reference points directly to an ExtendedExistentialTypeShape
		// kind = SymbolicReferenceKind::UniqueExtendedExistentialTypeShape;  TODO: implement
		// direct = Directness::Direct;
		fallthrough
	case 11: // Reference points directly to a NonUniqueExtendedExistentialTypeShape
		// kind = SymbolicReferenceKind::NonUniqueExtendedExistentialTypeShape;
		// direct = Directness::Direct;
		fallthrough
	default:
		// return "", fmt.Errorf("symbolic reference control character %#x is not implemented", controlData[0])
		return "(error)", nil
	}

	if symbolic {
		return "symbolic " + name, nil
	} else {
		return name, nil
	}
}

// GetSwiftBuiltinTypes parses all the built-in types in the __TEXT.__swift5_builtin section
func (f *File) GetSwiftBuiltinTypes() ([]swift.BuiltinType, error) {
	var builtins []swift.BuiltinType

	if sec := f.Section("__TEXT", "__swift5_builtin"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		builtInTypes := make([]swift.BuiltinTypeDescriptor, int(sec.Size)/binary.Size(swift.BuiltinTypeDescriptor{}))

		if err := binary.Read(bytes.NewReader(dat), f.ByteOrder, &builtInTypes); err != nil {
			return nil, fmt.Errorf("failed to read []swift.BuiltinTypeDescriptor: %v", err)
		}

		for idx, bType := range builtInTypes {
			currOffset := int64(sec.Offset) + int64(idx*binary.Size(swift.BuiltinTypeDescriptor{}))
			currAddr := sec.Addr + uint64(idx*binary.Size(swift.BuiltinTypeDescriptor{}))
			name, err := f.makeSymbolicMangledNameStringRef(currOffset + int64(bType.TypeName))
			if err != nil {
				return nil, fmt.Errorf("failed to read record.MangledTypeName; %v", err)
			}

			builtins = append(builtins, swift.BuiltinType{
				Address:             currAddr,
				Name:                name,
				Size:                bType.Size,
				Alignment:           bType.AlignmentAndFlags.Alignment(),
				BitwiseTakable:      bType.AlignmentAndFlags.IsBitwiseTakable(),
				Stride:              bType.Stride,
				NumExtraInhabitants: bType.NumExtraInhabitants,
			})
		}

		return builtins, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_builtin' section: %w", ErrSwiftSectionError)
}

// GetSwiftClosures parses all the closure context objects in the __TEXT.__swift5_capture section
func (f *File) GetSwiftClosures() ([]swift.CaptureDescriptor, error) {
	var closures []swift.CaptureDescriptor

	if sec := f.Section("__TEXT", "__swift5_capture"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		r := bytes.NewReader(dat)

		for {
			currOffset, _ := r.Seek(0, io.SeekCurrent)
			currAddr := sec.Addr + uint64(currOffset)
			currOffset += int64(sec.Offset)

			var capture swift.CaptureDescriptor
			err := binary.Read(r, f.ByteOrder, &capture.CaptureDescriptorHeader)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("failed to read swift %T: %v", capture.CaptureDescriptorHeader, err)
			}

			capture.Address = currAddr

			if capture.NumCaptureTypes > 0 {
				numCapsOffset := currOffset + int64(binary.Size(capture.CaptureDescriptorHeader))
				captureTypeRecords := make([]swift.CaptureTypeRecord, capture.NumCaptureTypes)
				if err := binary.Read(r, f.ByteOrder, &captureTypeRecords); err != nil {
					return nil, fmt.Errorf("failed to read %T: %v", captureTypeRecords, err)
				}
				for _, capRecord := range captureTypeRecords {
					name, err := f.makeSymbolicMangledNameStringRef(numCapsOffset + int64(capRecord.MangledTypeName))
					if err != nil {
						return nil, fmt.Errorf("failed to read mangled type name at offset %#x: %v", numCapsOffset+int64(capRecord.MangledTypeName), err)
					}
					capture.CaptureTypes = append(capture.CaptureTypes, name)
					numCapsOffset += int64(binary.Size(capRecord))
				}
			}

			if capture.NumMetadataSources > 0 {
				metadataSourceRecords := make([]swift.MetadataSourceRecord, capture.NumMetadataSources)
				if err := binary.Read(r, f.ByteOrder, &metadataSourceRecords); err != nil {
					return nil, fmt.Errorf("failed to read %T: %v", metadataSourceRecords, err)
				}
				for idx, metasource := range metadataSourceRecords {
					currOffset += int64(idx * binary.Size(swift.MetadataSourceRecord{}))
					typeName, err := f.makeSymbolicMangledNameStringRef(currOffset + int64(metasource.MangledTypeName))
					if err != nil {
						return nil, fmt.Errorf("failed to read mangled type name at offset %#x: %v", currOffset+int64(metasource.MangledTypeName), err)
					}
					metaSource, err := f.makeSymbolicMangledNameStringRef(currOffset + sizeOfInt32 + int64(metasource.MangledMetadataSource))
					if err != nil {
						return nil, fmt.Errorf("failed to read mangled metadata source at offset %#x: %v", currOffset+int64(metasource.MangledMetadataSource), err)
					}
					capture.MetadataSources = append(capture.MetadataSources, swift.MetadataSource{
						MangledType:           typeName,
						MangledMetadataSource: metaSource,
					})
				}
			}

			if capture.NumBindings > 0 {
				capture.Bindings = make([]swift.NecessaryBindings, capture.NumBindings)
				if err := binary.Read(r, f.ByteOrder, &capture.Bindings); err != nil {
					return nil, fmt.Errorf("failed to read %T: %v", capture.Bindings, err)
				}
			}

			closures = append(closures, capture)
		}

		return closures, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_capture' section: %w", ErrSwiftSectionError)
}

func (f *File) GetSwiftEntry() (uint64, error) {
	if sec := f.Section("__TEXT", "__swift5_entry"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return 0, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return 0, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		var swiftEntry int32
		if err := binary.Read(bytes.NewReader(dat), f.ByteOrder, &swiftEntry); err != nil {
			return 0, fmt.Errorf("failed to read __swift5_entry data: %v", err)
		}

		return sec.Addr + uint64(swiftEntry), nil
	}

	return 0, fmt.Errorf("MachO has no '__swift5_entry' section: %w", ErrSwiftSectionError)
}

func (f *File) GetSwiftDynamicReplacementInfo() (*types.AutomaticDynamicReplacements, error) {
	if sec := f.Section("__TEXT", "__swift5_replace"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		var rep types.AutomaticDynamicReplacements
		if err := binary.Read(bytes.NewReader(dat), f.ByteOrder, &rep); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", rep, err)
		}

		f.cr.Seek(int64(off)+int64(sizeOfInt32*2)+int64(rep.ReplacementScope), io.SeekStart)

		var rscope types.DynamicReplacementScope
		if err := binary.Read(f.cr, f.ByteOrder, &rscope); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", rscope, err)
		}

		return &rep, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_replace' section: %w", ErrSwiftSectionError)
}

func (f *File) GetSwiftDynamicReplacementInfoForOpaqueTypes() (*types.AutomaticDynamicReplacementsSome, error) {
	if sec := f.Section("__TEXT", "__swift5_replac2"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		var rep2 types.AutomaticDynamicReplacementsSome
		if err := binary.Read(bytes.NewReader(dat), f.ByteOrder, &rep2.Flags); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", rep2.Flags, err)
		}
		if err := binary.Read(bytes.NewReader(dat), f.ByteOrder, &rep2.NumReplacements); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", rep2.NumReplacements, err)
		}
		rep2.Replacements = make([]types.DynamicReplacementSomeDescriptor, rep2.NumReplacements)
		if err := binary.Read(bytes.NewReader(dat), f.ByteOrder, &rep2.Replacements); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", rep2.Replacements, err)
		}

		return &rep2, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_replac2' section: %w", ErrSwiftSectionError)
}

func (f *File) GetSwiftAccessibleFunctions() (*types.AccessibleFunctionsSection, error) {
	if sec := f.Section("__TEXT", "__swift5_acfuncs"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		var afsec types.AccessibleFunctionsSection
		if err := binary.Read(bytes.NewReader(dat), f.ByteOrder, &afsec); err != nil {
			return nil, fmt.Errorf("failed to read %T: %v", afsec, err)
		}

		return &afsec, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_acfuncs' section: %w", ErrSwiftSectionError)
}

func (f *File) GetSwiftTypeRefs() ([]string, error) {
	var typeRefs []string
	if sec := f.Section("__TEXT", "__swift5_typeref"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		r := bytes.NewBuffer(dat)

		for {
			s, err := r.ReadString('\x00')
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("failed to read from type ref string pool: %v", err)
			}

			s = strings.TrimSpace(strings.Trim(s, "\x00"))

			if len(s) > 0 {
				typeRefs = append(typeRefs, s)
			}
		}

		return typeRefs, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_typeref' section: %w", ErrSwiftSectionError)
}

func (f *File) GetSwiftReflectionStrings() (map[uint64]string, error) {
	reflStrings := make(map[uint64]string)
	if sec := f.Section("__TEXT", "__swift5_reflstr"); sec != nil {
		off, err := f.vma.GetOffset(f.vma.Convert(sec.Addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert vmaddr: %v", err)
		}
		f.cr.Seek(int64(off), io.SeekStart)

		dat := make([]byte, sec.Size)
		if err := binary.Read(f.cr, f.ByteOrder, dat); err != nil {
			return nil, fmt.Errorf("failed to read %s.%s data: %v", sec.Seg, sec.Name, err)
		}

		r := bytes.NewBuffer(dat)

		for {
			s, err := r.ReadString('\x00')
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("failed to read from class name string pool: %v", err)
			}

			if len(strings.Trim(s, "\x00")) > 0 {
				reflStrings[sec.Addr+(sec.Size-uint64(r.Len()+len(s)))] = strings.Trim(s, "\x00")
			}
		}

		return reflStrings, nil
	}

	return nil, fmt.Errorf("MachO has no '__swift5_reflstr' section: %w", ErrSwiftSectionError)
}
