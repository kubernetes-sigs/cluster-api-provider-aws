package codesign

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/blacktop/go-macho/pkg/codesign/types"
)

// CodeSignature object
type CodeSignature struct {
	CodeDirectories              []types.CodeDirectory `json:"code_directories,omitempty"`
	Requirements                 []types.Requirement   `json:"requirements,omitempty"`
	CMSSignature                 []byte                `json:"cms_signature,omitempty"`
	Entitlements                 string                `json:"entitlements,omitempty"`
	EntitlementsDER              []byte                `json:"entitlements_der,omitempty"`
	LaunchConstraintsSelf        []byte                `json:"launch_constraints_self,omitempty"`
	LaunchConstraintsParent      []byte                `json:"launch_constraints_parent,omitempty"`
	LaunchConstraintsResponsible []byte                `json:"launch_constraints_responsible,omitempty"`
	LibraryConstraints           []byte                `json:"library_constraints,omitempty"`
	Errors                       []error               `json:"errors,omitempty"`
}

// ParseCodeSignature parses the LC_CODE_SIGNATURE data
func ParseCodeSignature(cmddat []byte) (*CodeSignature, error) {
	r := bytes.NewReader(cmddat)
	cs := &CodeSignature{}

	csBlob := types.SuperBlob{}
	if err := binary.Read(r, binary.BigEndian, &csBlob.SbHeader); err != nil {
		return nil, err
	}

	csIndex := make([]types.BlobIndex, csBlob.Count)
	if err := binary.Read(r, binary.BigEndian, &csIndex); err != nil {
		return nil, err
	}

	for _, index := range csIndex {

		r.Seek(int64(index.Offset), io.SeekStart)

		switch index.Type {
		case types.CSSLOT_CODEDIRECTORY:
			fallthrough
		case types.CSSLOT_ALTERNATE_CODEDIRECTORIES:
			cd, err := parseCodeDirectory(r, index.Offset)
			if err != nil {
				return nil, err
			}
			cs.CodeDirectories = append(cs.CodeDirectories, *cd)
		case types.CSSLOT_REQUIREMENTS:
			req := types.Requirement{}
			if err := binary.Read(r, binary.BigEndian, &req.RequirementsBlob); err != nil {
				return nil, err
			}
			if req.RequirementsBlob.Magic != types.MAGIC_REQUIREMENT && req.RequirementsBlob.Magic != types.MAGIC_REQUIREMENTS {
				return nil, fmt.Errorf("invalid CSSLOT_REQUIREMENTS blob magic: %s", req.RequirementsBlob.Magic)
			}
			datLen := int(req.RequirementsBlob.Length) - binary.Size(types.RequirementsBlob{})
			if datLen > 0 {
				reqData := make([]byte, datLen)
				if err := binary.Read(r, binary.BigEndian, &reqData); err != nil {
					return nil, err
				}
				rqr := bytes.NewReader(reqData)
				if err := binary.Read(rqr, binary.BigEndian, &req.Requirements); err != nil {
					return nil, err
				}
				detail, err := types.ParseRequirements(rqr, req.Requirements)
				if err != nil {
					return nil, err
				}
				req.Detail = detail
			} else {
				req.Detail = "empty requirement set"
			}
			cs.Requirements = append(cs.Requirements, req)
		case types.CSSLOT_ENTITLEMENTS:
			entBlob := types.BlobHeader{}
			if err := binary.Read(r, binary.BigEndian, &entBlob); err != nil {
				return nil, err
			}
			if entBlob.Magic != types.MAGIC_EMBEDDED_ENTITLEMENTS {
				return nil, fmt.Errorf("invalid CSSLOT_ENTITLEMENTS blob magic: %s", entBlob.Magic)
			}
			plistData := make([]byte, int(entBlob.Length)-binary.Size(entBlob))
			if err := binary.Read(r, binary.BigEndian, &plistData); err != nil {
				return nil, err
			}
			cs.Entitlements = string(plistData)
		case types.CSSLOT_CMS_SIGNATURE:
			cmsBlob := types.BlobHeader{}
			if err := binary.Read(r, binary.BigEndian, &cmsBlob); err != nil {
				return nil, err
			}
			if cmsBlob.Magic != types.MAGIC_BLOBWRAPPER {
				return nil, fmt.Errorf("invalid CSSLOT_CMS_SIGNATURE blob magic: %s", cmsBlob.Magic)
			}
			cmsData := make([]byte, int(cmsBlob.Length)-binary.Size(cmsBlob))
			if err := binary.Read(r, binary.BigEndian, &cmsData); err != nil {
				return nil, err
			}
			// NOTE: openssl pkcs7 -inform DER -in <cmsData> -print_certs -text -noout
			cs.CMSSignature = cmsData
		case types.CSSLOT_ENTITLEMENTS_DER:
			entDerBlob := types.BlobHeader{}
			if err := binary.Read(r, binary.BigEndian, &entDerBlob); err != nil {
				return nil, err
			}
			if entDerBlob.Magic != types.MAGIC_EMBEDDED_ENTITLEMENTS_DER {
				return nil, fmt.Errorf("invalid CSSLOT_ENTITLEMENTS_DER blob magic: %s", entDerBlob.Magic)
			}
			entDerData := make([]byte, int(entDerBlob.Length)-binary.Size(entDerBlob))
			if err := binary.Read(r, binary.BigEndian, &entDerData); err != nil {
				return nil, err
			}
			cs.EntitlementsDER = entDerData
		case types.CSSLOT_REP_SPECIFIC:
			fallthrough // TODO 🤷‍♂️
		case types.CSSLOT_INFOSLOT:
			fallthrough // TODO 🤷‍♂️
		case types.CSSLOT_RESOURCEDIR:
			fallthrough // TODO 🤷‍♂️
		case types.CSSLOT_APPLICATION:
			fallthrough // TODO 🤷‍♂️
		case types.CSSLOT_IDENTIFICATIONSLOT:
			fallthrough // TODO 🤷‍♂️
		case types.CSSLOT_TICKETSLOT:
			fallthrough // TODO 🤷‍♂️
		case types.CSSLOT_LAUNCH_CONSTRAINT_SELF, types.CSSLOT_LAUNCH_CONSTRAINT_PARENT, types.CSSLOT_LAUNCH_CONSTRAINT_RESPONSIBLE, types.CSSLOT_LIBRARY_CONSTRAINT:
			lcBlob := types.BlobHeader{}
			if err := binary.Read(r, binary.BigEndian, &lcBlob); err != nil {
				return nil, err
			}
			if lcBlob.Magic != types.MAGIC_EMBEDDED_LAUNCH_CONSTRAINT {
				return nil, fmt.Errorf("invalid CSSLOT_LAUNCH_CONSTRAINT_SELF blob magic: %s", lcBlob.Magic)
			}
			lcData := make([]byte, int(lcBlob.Length)-binary.Size(lcBlob))
			if err := binary.Read(r, binary.BigEndian, &lcData); err != nil {
				return nil, err
			}
			switch index.Type {
			case types.CSSLOT_LAUNCH_CONSTRAINT_SELF:
				cs.LaunchConstraintsSelf = lcData
			case types.CSSLOT_LAUNCH_CONSTRAINT_PARENT:
				cs.LaunchConstraintsParent = lcData
			case types.CSSLOT_LAUNCH_CONSTRAINT_RESPONSIBLE:
				cs.LaunchConstraintsResponsible = lcData
			case types.CSSLOT_LIBRARY_CONSTRAINT:
				cs.LibraryConstraints = lcData
			}
		default:
			cs.Errors = append(cs.Errors, fmt.Errorf("unknown slot type: %s, please notify author", index.Type))
		}
	}
	return cs, nil
}

func parseCodeDirectory(r *bytes.Reader, offset uint32) (*types.CodeDirectory, error) {
	var cd types.CodeDirectory
	if err := binary.Read(r, binary.BigEndian, &cd.BlobHeader); err != nil {
		return nil, err
	}
	if cd.BlobHeader.Magic != types.MAGIC_CODEDIRECTORY {
		return nil, fmt.Errorf("invalid CSSLOT_(ALTERNATE_)CODEDIRECTORY blob magic: %#x", cd.BlobHeader.Magic)
	}
	if err := binary.Read(r, binary.BigEndian, &cd.Header.CdEarliest); err != nil {
		return nil, err
	}
	headoff, _ := r.Seek(0, io.SeekCurrent)
	// Calculate the cdhashs
	r.Seek(int64(offset), io.SeekStart)
	cdData := make([]byte, cd.BlobHeader.Length)
	if err := binary.Read(r, binary.LittleEndian, &cdData); err != nil {
		return nil, err
	}

	switch cd.Header.HashType {
	case types.HASHTYPE_SHA1:
		h := sha1.New()
		h.Write(cdData)
		cd.CDHash = fmt.Sprintf("%x", h.Sum(nil))
	case types.HASHTYPE_SHA256:
		h := sha256.New()
		h.Write(cdData)
		cd.CDHash = fmt.Sprintf("%x", h.Sum(nil))
	default:
		cd.CDHash = fmt.Sprintf("unsupported code directory hash type %s, please notify author", cd.Header.HashType)
	}

	// Parse version
	if cd.Header.Version < types.EARLIEST_VERSION {
		fmt.Printf("unsupported type or version of signature: %#x (too old)\n", cd.Header.Version)
	} else if cd.Header.Version > types.COMPATIBILITY_LIMIT {
		fmt.Printf("unsupported type or version of signature: %#x (too new)\n", cd.Header.Version)
	}

	// SUPPORTS_SCATTER
	if cd.Header.Version >= types.SUPPORTS_SCATTER {
		r.Seek(int64(headoff), io.SeekStart)
		if err := binary.Read(r, binary.BigEndian, &cd.Header.CdScatter); err != nil {
			return nil, err
		}
		if cd.Header.ScatterOffset > 0 {
			r.Seek(int64(offset+cd.Header.ScatterOffset), io.SeekStart)
			if err := binary.Read(r, binary.BigEndian, &cd.Scatter); err != nil {
				return nil, fmt.Errorf("failed to read SUPPORTS_SCATTER @ %#x: %v", offset+cd.Header.ScatterOffset, err)
			}
		}
	}
	// SUPPORTS_TEAMID
	if cd.Header.Version >= types.SUPPORTS_TEAMID {
		r.Seek(int64(headoff)+int64(binary.Size(cd.Header.CdScatter)), io.SeekStart)
		if err := binary.Read(r, binary.BigEndian, &cd.Header.CdTeamID); err != nil {
			return nil, err
		}
		if cd.Header.TeamOffset > 0 {
			r.Seek(int64(offset+cd.Header.TeamOffset), io.SeekStart)
			teamID, err := bufio.NewReader(r).ReadString('\x00')
			if err != nil {
				return nil, fmt.Errorf("failed to read SUPPORTS_TEAMID @ %#x: %v", offset+cd.Header.TeamOffset, err)
			}
			cd.TeamID = strings.Trim(teamID, "\x00")
		}
	}
	// SUPPORTS_CODELIMIT64
	if cd.Header.Version >= types.SUPPORTS_CODELIMIT64 {
		r.Seek(int64(headoff)+
			int64(binary.Size(cd.Header.CdScatter))+
			int64(binary.Size(cd.Header.CdTeamID)), io.SeekStart)
		if err := binary.Read(r, binary.BigEndian, &cd.Header.CdCodeLimit64); err != nil {
			return nil, err
		}
		cd.CodeLimit = uint64(cd.Header.CodeLimit)
		if cd.Header.CodeLimit64 > 0 {
			cd.CodeLimit = cd.Header.CodeLimit64
		}
	}
	// SUPPORTS_EXECSEG
	if cd.Header.Version >= types.SUPPORTS_EXECSEG {
		r.Seek(int64(headoff)+
			int64(binary.Size(cd.Header.CdScatter))+
			int64(binary.Size(cd.Header.CdTeamID))+
			int64(binary.Size(cd.Header.CdCodeLimit64)), io.SeekStart)
		if err := binary.Read(r, binary.BigEndian, &cd.Header.CdExecSeg); err != nil {
			return nil, err
		}
	}
	// SUPPORTS_RUNTIME
	if cd.Header.Version >= types.SUPPORTS_RUNTIME {
		r.Seek(int64(headoff)+
			int64(binary.Size(cd.Header.CdScatter))+
			int64(binary.Size(cd.Header.CdTeamID))+
			int64(binary.Size(cd.Header.CdCodeLimit64))+
			int64(binary.Size(cd.Header.CdExecSeg)), io.SeekStart)
		if err := binary.Read(r, binary.BigEndian, &cd.Header.CdRuntime); err != nil {
			return nil, err
		}
		cd.RuntimeVersion = cd.Header.Runtime.String()
		if cd.Header.PreEncryptOffset > 0 {
			r.Seek(int64(offset+cd.Header.PreEncryptOffset), io.SeekStart)
			for i := uint8(0); i < uint8(cd.Header.NCodeSlots); i++ {
				slot := make([]byte, cd.Header.HashSize)
				if err := binary.Read(r, binary.BigEndian, &slot); err != nil {
					return nil, fmt.Errorf("failed to read SUPPORTS_RUNTIME PreEncrypt hash slot #%d @ %#x: %v",
						i, offset+cd.Header.PreEncryptOffset+uint32(i*cd.Header.HashSize), err)
				}
				cd.PreEncryptSlots = append(cd.PreEncryptSlots, slot)
			}
		}
	}
	// SUPPORTS_LINKAGE
	if cd.Header.Version >= types.SUPPORTS_LINKAGE {
		r.Seek(int64(headoff)+
			int64(binary.Size(cd.Header.CdScatter))+
			int64(binary.Size(cd.Header.CdTeamID))+
			int64(binary.Size(cd.Header.CdCodeLimit64))+
			int64(binary.Size(cd.Header.CdExecSeg))+
			int64(binary.Size(cd.Header.CdRuntime)), io.SeekStart)
		if err := binary.Read(r, binary.BigEndian, &cd.Header.CdLinkage); err != nil {
			return nil, err
		}
		if cd.Header.LinkageOffset > 0 {
			r.Seek(int64(offset+cd.Header.LinkageOffset), io.SeekStart)
			cd.LinkageData = make([]byte, cd.Header.LinkageSize)
			if err := binary.Read(r, binary.BigEndian, &cd.LinkageData); err != nil {
				return nil, fmt.Errorf("failed to read SUPPORTS_LINKAGE @ %#x: %v", offset+cd.Header.LinkageOffset, err)
			}
			// TODO: what IS linkage
		}
	}
	// Parse Indentity
	r.Seek(int64(offset+cd.Header.IdentOffset), io.SeekStart)
	id, err := bufio.NewReader(r).ReadString('\x00')
	if err != nil {
		return nil, fmt.Errorf("failed to read CodeDirectory ID at: %d: %v", offset+cd.Header.IdentOffset, err)
	}
	cd.ID = strings.Trim(id, "\x00")
	// Parse Special Slots
	r.Seek(int64(offset+cd.Header.HashOffset-(cd.Header.NSpecialSlots*uint32(cd.Header.HashSize))), io.SeekStart)
	for slot := cd.Header.NSpecialSlots; slot > 0; slot-- {
		hash := make([]byte, cd.Header.HashSize)
		if err := binary.Read(r, binary.BigEndian, &hash); err != nil {
			return nil, err
		}
		sslot := types.SpecialSlot{
			Index: slot,
			Hash:  hash,
		}
		if bytes.Equal(hash, make([]byte, cd.Header.HashSize)) { // empty hash
			sslot.Desc = fmt.Sprintf("Special Slot   %d %-22v Not Bound", slot, types.SlotType(slot).String()+":")
		} else if bytes.Equal(hash, types.EmptySha256ReqSlot) && sslot.Index == 2 && cd.Header.HashType == types.HASHTYPE_SHA256 {
			sslot.Desc = fmt.Sprintf("Special Slot   %d %-22v Empty Requirement Set", slot, types.SlotType(slot).String()+":")
		} else {
			sslot.Desc = fmt.Sprintf("Special Slot   %d %-22v %x", slot, types.SlotType(slot).String()+":", hash)
		}
		cd.SpecialSlots = append(cd.SpecialSlots, sslot)
	}
	// Parse Slots
	pageSize := uint32(1 << cd.Header.PageSize)
	for slot := uint32(0); slot < cd.Header.NCodeSlots; slot++ {
		hash := make([]byte, cd.Header.HashSize)
		if err := binary.Read(r, binary.BigEndian, &hash); err != nil {
			return nil, err
		}
		cslot := types.CodeSlot{
			Index: slot,
			Page:  slot * pageSize,
			Hash:  hash,
		}
		if bytes.Equal(hash, types.NULL_PAGE_SHA256_HASH) && cd.Header.HashType == types.HASHTYPE_SHA256 {
			cslot.Desc = fmt.Sprintf("Slot   %d (File page @0x%04X):\tNULL PAGE HASH", slot, cslot.Page)
		} else {
			cslot.Desc = fmt.Sprintf("Slot   %d (File page @0x%04X):\t%x", slot, cslot.Page, hash)
		}
		cd.CodeSlots = append(cd.CodeSlots, cslot)
	}

	return &cd, nil
}

type slotHashes struct {
	InfoPlist       []byte
	Requirements    []byte
	ResourceDir     []byte
	Entitlements    []byte
	AppSpecific     []byte
	DmgSpecific     []byte
	EntitlementsDER []byte
}

type Config struct {
	ID                  string
	TeamID              string
	IsMain              bool
	Flags               types.CDFlag
	CodeSize            uint64
	TextOffset          uint64
	TextSize            uint64
	NSpecialSlots       uint32
	SpecialSlots        []types.SpecialSlot
	InfoPlist           []byte
	Entitlements        []byte
	EntitlementsDER     []byte
	ResourceDirSlotHash []byte
	SlotHashes          slotHashes
	CertChain           []*x509.Certificate
	SignerFunction      func([]byte) ([]byte, error)
}

func (c *Config) InitSlotHashes() {
	c.SlotHashes = slotHashes{
		InfoPlist:       types.EmptySha256Slot,
		Requirements:    types.EmptySha256ReqSlot,
		ResourceDir:     types.EmptySha256Slot,
		Entitlements:    types.EmptySha256Slot,
		AppSpecific:     types.EmptySha256Slot,
		DmgSpecific:     types.EmptySha256Slot,
		EntitlementsDER: types.EmptySha256Slot,
	}
}

func Sign(r io.Reader, config *Config) ([]byte, error) {
	var err error
	var buf bytes.Buffer
	var reqBlob types.Blob
	var entBlob types.Blob
	var entDerBlob types.Blob

	sb := types.NewSuperBlob(types.MAGIC_EMBEDDED_SIGNATURE)

	// Requirements /////////////////////////////////////////////
	reqBlob, err = types.CreateRequirements(config.ID, config.CertChain)
	if err != nil {
		return nil, fmt.Errorf("failed to create Requirements: %v", err)
	}
	config.SlotHashes.Requirements, err = reqBlob.Sha256Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash Requirements: %v", err)
	}

	config.NSpecialSlots = uint32(2)
	if len(config.SpecialSlots) > 0 {
		config.NSpecialSlots = uint32(len(config.SpecialSlots))
	}

	// Entitlements /////////////////////////////////////////////
	if len(config.Entitlements) > 0 {
		config.NSpecialSlots = 7
		entBlob = types.NewBlob(types.MAGIC_EMBEDDED_ENTITLEMENTS, config.Entitlements)
		config.SlotHashes.Entitlements, err = entBlob.Sha256Hash()
		if err != nil {
			return nil, fmt.Errorf("failed to hash entitlements plist blob: %v", err)
		}
		if len(config.SpecialSlots) >= 5 { // if we have previous entitlements plist hash, verify it against the new one
			if len(config.SpecialSlots[2].Hash) > 0 && !bytes.Equal(config.SpecialSlots[2].Hash, config.SlotHashes.Entitlements) {
				return nil, fmt.Errorf("previous and calulated entitlements plist hashes do not match")
			}
		}
		if len(config.EntitlementsDER) == 0 && bytes.Equal(config.SpecialSlots[0].Hash, types.EmptySha256Slot) {
			config.NSpecialSlots = 5
		} else {
			entDerBlob = types.NewBlob(types.MAGIC_EMBEDDED_ENTITLEMENTS_DER, config.EntitlementsDER)
			config.SlotHashes.EntitlementsDER, err = entDerBlob.Sha256Hash()
			if err != nil {
				return nil, fmt.Errorf("failed to hash entitlements asn1/der blob: %v", err)
			}
			if len(config.SpecialSlots) >= 7 { // if we have previous entitlements asn1/der hash, verify it against the new one
				if len(config.SpecialSlots[0].Hash) > 0 && !bytes.Equal(config.SpecialSlots[0].Hash, config.SlotHashes.EntitlementsDER) {
					return nil, fmt.Errorf("previous and calulated entitlements asn1/der hashes do not match")
				}
			}
		}
	}

	// CodeDirectory ////////////////////////////////////////////
	cdbuf, err := createCodeDirectory(r, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create CodeDirectory: %v", err)
	}
	// Blobs ////////////////////////////////////////////////////
	sb.AddBlob(types.CSSLOT_CODEDIRECTORY, types.NewBlob(types.MAGIC_CODEDIRECTORY, cdbuf.Bytes()))
	sb.AddBlob(types.CSSLOT_REQUIREMENTS, reqBlob)
	if len(config.Entitlements) > 0 {
		sb.AddBlob(types.CSSLOT_ENTITLEMENTS, entBlob)
	}
	if len(config.EntitlementsDER) > 0 {
		sb.AddBlob(types.CSSLOT_ENTITLEMENTS_DER, entDerBlob)
	}
	if config.SignerFunction != nil {
		cdblob, err := sb.GetBlob(types.CSSLOT_CODEDIRECTORY)
		if err != nil {
			return nil, fmt.Errorf("failed to get CodeDirectory blob: %v", err)
		}
		cddata, err := cdblob.Bytes()
		if err != nil {
			return nil, fmt.Errorf("failed to get CodeDirectory blob data: %v", err)
		}
		cert, err := config.SignerFunction(cddata)
		if err != nil {
			return nil, fmt.Errorf("failed to sign CodeDirectory blob: %v", err)
		}
		sb.AddBlob(types.CSSLOT_CMS_SIGNATURE, types.NewBlob(types.MAGIC_BLOBWRAPPER, cert))
	} else {
		sb.AddBlob(types.CSSLOT_CMS_SIGNATURE, types.NewBlob(types.MAGIC_BLOBWRAPPER, []byte{}))
	}

	if uint32(sb.Size()) < sb.Length { // TODO: should I remove this check?
		return nil, fmt.Errorf("SuperBlob size mismatch: calculated Size %d != Length %d", sb.Size(), sb.Length)
	}

	// write SuperBlob
	if err := sb.Write(&buf, binary.BigEndian); err != nil {
		return nil, fmt.Errorf("failed to write SuperBlob: %v", err)
	}

	return buf.Bytes(), nil
}

func createCodeDirectory(r io.Reader, config *Config) (*bytes.Buffer, error) {
	var cddelta int
	var cdbuf bytes.Buffer

	// Info.plist ///////////////////////////////////////////////
	if config.InfoPlist != nil {
		h := sha256.New()
		if _, err := h.Write(config.InfoPlist); err != nil {
			return nil, fmt.Errorf("failed to hash Info.plist: %v", err)
		}
		config.SlotHashes.InfoPlist = h.Sum(nil)
		if len(config.SpecialSlots) >= 1 {
			if len(config.SpecialSlots[len(config.SpecialSlots)-1].Hash) > 0 && !bytes.Equal(config.SpecialSlots[len(config.SpecialSlots)-1].Hash, config.SlotHashes.InfoPlist) {
				return nil, fmt.Errorf("previous and calulated Info.plist hashes do not match")
			}
		}
	}

	// Resource Directory ///////////////////////////////////////
	if len(config.SpecialSlots) >= 3 {
		// NOTE: this is sha256sum Some.app/Contents/_CodeSignature/CodeResources (which is a XML representation of the Resources directory)
		if bytes.Equal(config.SlotHashes.ResourceDir, types.EmptySha256Slot) {
			// if the slot is empty it was NOT set by the caller (try and reuse previous value)
			config.SlotHashes.ResourceDir = config.SpecialSlots[len(config.SpecialSlots)-3].Hash
		}
		// if the slot is NOT empty it was set by the caller (and was calculated from the created CodeResources file)
	}

	// Application Specific /////////////////////////////////////
	if len(config.SpecialSlots) >= 4 {
		if bytes.Equal(config.SlotHashes.AppSpecific, types.EmptySha256Slot) {
			// if the slot is empty it was NOT set by the caller (try and reuse previous value)
			config.SlotHashes.AppSpecific = config.SpecialSlots[len(config.SpecialSlots)-4].Hash
		}
	}

	// DMG Specific /////////////////////////////////////////////
	if len(config.SpecialSlots) >= 6 {
		if bytes.Equal(config.SlotHashes.DmgSpecific, types.EmptySha256Slot) {
			// if the slot is empty it was NOT set by the caller (try and reuse previous value)
			config.SlotHashes.DmgSpecific = config.SpecialSlots[len(config.SpecialSlots)-6].Hash
		}
	}

	// calculate the CodeDirectory offsets
	identOffset := uint32(binary.Size(types.BlobHeader{}) + binary.Size(types.CodeDirectoryType{}))
	hashOffset := identOffset + uint32(len(config.ID)+1+len(types.EmptySha256Slot)*int(config.NSpecialSlots))

	cdHeader := types.CodeDirectoryType{
		CdEarliest: types.CdEarliest{
			Version:       types.SUPPORTS_EXECSEG, // TODO: support other versions (e.g.SUPPORTS_RUNTIME)
			Flags:         config.Flags,
			HashOffset:    hashOffset,
			IdentOffset:   identOffset,
			NSpecialSlots: config.NSpecialSlots,
			NCodeSlots:    uint32((int(config.CodeSize) + types.PAGE_SIZE - 1) / types.PAGE_SIZE),
			CodeLimit:     uint32(config.CodeSize),
			HashSize:      sha256.Size,
			HashType:      types.HASHTYPE_SHA256,
			PageSize:      uint8(types.PAGE_SIZE_BITS),
		},
		CdExecSeg: types.CdExecSeg{
			ExecSegBase:  uint64(config.TextOffset),
			ExecSegLimit: uint64(config.TextSize),
		},
	}

	// CodeDirectoryType is a variable length struct based on the Version field
	if cdHeader.Version >= types.EARLIEST_VERSION {
		cddelta = binary.Size(types.CodeDirectoryType{}) - binary.Size(types.CdEarliest{})
	}
	if cdHeader.Version >= types.SUPPORTS_SCATTER {
		cddelta -= binary.Size(types.CdScatter{})
	}
	if cdHeader.Version >= types.SUPPORTS_TEAMID {
		cddelta -= binary.Size(types.CdTeamID{})
	}
	if cdHeader.Version >= types.SUPPORTS_CODELIMIT64 {
		cddelta -= binary.Size(types.CdCodeLimit64{})
	}
	if cdHeader.Version >= types.SUPPORTS_EXECSEG {
		cddelta -= binary.Size(types.CdExecSeg{})
	}
	if cdHeader.Version >= types.SUPPORTS_RUNTIME {
		cddelta -= binary.Size(types.CdRuntime{})
	}
	if cdHeader.Version >= types.SUPPORTS_LINKAGE {
		cddelta -= binary.Size(types.CdLinkage{})
	}
	// adjust CodeDirectory header offsets
	cdHeader.IdentOffset -= uint32(cddelta)
	cdHeader.HashOffset -= uint32(cddelta)

	if config.IsMain {
		cdHeader.ExecSegFlags = types.EXECSEG_MAIN_BINARY
	}

	// write CodeDirectory header
	if err := binary.Write(&cdbuf, binary.BigEndian, &cdHeader); err != nil {
		return nil, fmt.Errorf("failed to write CodeDirectory: %v", err)
	}
	// truncate CodeDirectory header to match Version length
	cdbuf.Truncate(cdbuf.Len() - cddelta)
	// write CodeDirectory identifier
	if _, err := cdbuf.WriteString(config.ID + "\x00"); err != nil {
		return nil, fmt.Errorf("failed to write identifier %s: %v", config.ID, err)
	}
	if len(config.Entitlements) > 0 {
		// write CodeDirectory Entitlements ASN1/DER slot hash
		if _, err := cdbuf.Write(config.SlotHashes.EntitlementsDER); err != nil {
			return nil, fmt.Errorf("failed to write entitlements asn1/der hash: %v", err)
		}
		// write CodeDirectory DMG Specific slot hash
		if _, err := cdbuf.Write(config.SlotHashes.DmgSpecific); err != nil {
			return nil, fmt.Errorf("failed to write dmg specific hash: %v", err)
		}
		// write CodeDirectory Entitlements Plist slot hash
		if _, err := cdbuf.Write(config.SlotHashes.Entitlements); err != nil {
			return nil, fmt.Errorf("failed to write entitlements plist hash: %v", err)
		}
		// write CodeDirectory Application Specific slot hash
		if _, err := cdbuf.Write(config.SlotHashes.AppSpecific); err != nil {
			return nil, fmt.Errorf("failed to write app specific hash: %v", err)
		}
		// write CodeDirectory Resource Directory slot hash
		if _, err := cdbuf.Write(config.SlotHashes.ResourceDir); err != nil {
			return nil, fmt.Errorf("failed to write rsc dir hash: %v", err)
		}
	}
	// write CodeDirectory Requirements Blob slot hash
	if _, err := cdbuf.Write(config.SlotHashes.Requirements); err != nil {
		return nil, fmt.Errorf("failed to write requirements hash: %v", err)
	}
	// write CodeDirectory Bound Info.plist slot hash
	if _, err := cdbuf.Write(config.SlotHashes.InfoPlist); err != nil {
		return nil, fmt.Errorf("failed to write info.plist hash: %v", err)
	}
	// write page hashes
	var hashCount int
	var hashes [types.PAGE_SIZE]byte
	h := sha256.New()
	p := 0
	for p < int(config.CodeSize) {
		n, err := io.ReadFull(r, hashes[:])
		if err == io.EOF {
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("failed to read file content without the signature: %v", err)
		}
		if p+n > int(config.CodeSize) {
			n = int(config.CodeSize) - p
		}
		p += n
		h.Reset()
		h.Write(hashes[:n])
		b := h.Sum(nil)
		if _, err := cdbuf.Write(b[:]); err != nil {
			return nil, fmt.Errorf("failed to write page %d hash: %v", hashCount, err)
		}
		hashCount++
	}

	return &cdbuf, nil
}
