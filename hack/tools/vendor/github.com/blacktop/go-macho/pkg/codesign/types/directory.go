package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"

	mtypes "github.com/blacktop/go-macho/types"
)

var (
	EmptySha256Slot    = bytes.Repeat([]byte{0}, sha256.New().Size())
	EmptySha256ReqSlot = []byte{
		0x98, 0x79, 0x20, 0x90, 0x4E, 0xAB, 0x65, 0x0E,
		0x75, 0x78, 0x8C, 0x05, 0x4A, 0xA0, 0xB0, 0x52,
		0x4E, 0x6A, 0x80, 0xBF, 0xC7, 0x1A, 0xA3, 0x2D,
		0xF8, 0xD2, 0x37, 0xA6, 0x17, 0x43, 0xF9, 0x86,
	}
)

// CodeDirectory object
type CodeDirectory struct {
	BlobHeader
	ID             string            `json:"id,omitempty"`
	TeamID         string            `json:"team_id,omitempty"`
	Scatter        Scatter           `json:"scatter,omitempty"`
	CDHash         string            `json:"cd_hash,omitempty"`
	SpecialSlots   []SpecialSlot     `json:"special_slots,omitempty"`
	CodeSlots      []CodeSlot        `json:"code_slots,omitempty"`
	Header         CodeDirectoryType `json:"header,omitempty"`
	RuntimeVersion string            `json:"runtime_version,omitempty"`
	CodeLimit      uint64            `json:"code_limit,omitempty"`

	PreEncryptSlots [][]byte `json:"pre_encrypt_slots,omitempty"`
	LinkageData     []byte   `json:"linkage_data,omitempty"`
}

type SpecialSlot struct {
	Index uint32 `json:"index,omitempty"`
	Hash  []byte `json:"hash,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

type CodeSlot struct {
	Index uint32 `json:"index,omitempty"`
	Page  uint32 `json:"page,omitempty"`
	Hash  []byte `json:"hash,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

type hashType uint8

const (
	PAGE_SIZE_BITS = 12
	PAGE_SIZE      = 1 << PAGE_SIZE_BITS

	HASHTYPE_NOHASH           hashType = 0
	HASHTYPE_SHA1             hashType = 1
	HASHTYPE_SHA256           hashType = 2
	HASHTYPE_SHA256_TRUNCATED hashType = 3
	HASHTYPE_SHA384           hashType = 4
	HASHTYPE_SHA512           hashType = 5

	HASH_SIZE_SHA1             = 20
	HASH_SIZE_SHA256           = 32
	HASH_SIZE_SHA256_TRUNCATED = 20

	CDHASH_LEN    = 20 /* always - larger hashes are truncated */
	HASH_MAX_SIZE = 48 /* max size of the hash we'll support */
)

func (c hashType) String() string {
	switch c {
	case HASHTYPE_NOHASH:
		return "No Hash"
	case HASHTYPE_SHA1:
		return "Sha1"
	case HASHTYPE_SHA256:
		return "Sha256"
	case HASHTYPE_SHA256_TRUNCATED:
		return "Sha256 (Truncated)"
	case HASHTYPE_SHA384:
		return "Sha384"
	case HASHTYPE_SHA512:
		return "Sha512"
	default:
		return fmt.Sprintf("hashType(%d)", c)
	}
}

type cdVersion uint32

const (
	EARLIEST_VERSION     cdVersion = 0x20001
	SUPPORTS_SCATTER     cdVersion = 0x20100
	SUPPORTS_TEAMID      cdVersion = 0x20200
	SUPPORTS_CODELIMIT64 cdVersion = 0x20300
	SUPPORTS_EXECSEG     cdVersion = 0x20400
	SUPPORTS_RUNTIME     cdVersion = 0x20500
	SUPPORTS_LINKAGE     cdVersion = 0x20600
	COMPATIBILITY_LIMIT  cdVersion = 0x2F000 // "version 3 with wiggle room"
)

func (v cdVersion) String() string {
	switch v {
	case SUPPORTS_SCATTER:
		return "Scatter"
	case SUPPORTS_TEAMID:
		return "TeamID"
	case SUPPORTS_CODELIMIT64:
		return "Codelimit64"
	case SUPPORTS_EXECSEG:
		return "ExecSeg"
	case SUPPORTS_RUNTIME:
		return "Runtime"
	case SUPPORTS_LINKAGE:
		return "Linkage"
	default:
		return fmt.Sprintf("cdVersion(%d)", v)
	}
}

type CDFlag uint32

const (
	/* code signing attributes of a process */
	NONE           CDFlag = 0x00000000 /* no flags */
	VALID          CDFlag = 0x00000001 /* dynamically valid */
	ADHOC          CDFlag = 0x00000002 /* ad hoc signed */
	GET_TASK_ALLOW CDFlag = 0x00000004 /* has get-task-allow entitlement */
	INSTALLER      CDFlag = 0x00000008 /* has installer entitlement */

	FORCED_LV       CDFlag = 0x00000010 /* Library Validation required by Hardened System Policy */
	INVALID_ALLOWED CDFlag = 0x00000020 /* (macOS Only) Page invalidation allowed by task port policy */

	HARD             CDFlag = 0x00000100 /* don't load invalid pages */
	KILL             CDFlag = 0x00000200 /* kill process if it becomes invalid */
	CHECK_EXPIRATION CDFlag = 0x00000400 /* force expiration checking */
	RESTRICT         CDFlag = 0x00000800 /* tell dyld to treat restricted */

	ENFORCEMENT            CDFlag = 0x00001000 /* require enforcement */
	REQUIRE_LV             CDFlag = 0x00002000 /* require library validation */
	ENTITLEMENTS_VALIDATED CDFlag = 0x00004000 /* code signature permits restricted entitlements */
	NVRAM_UNRESTRICTED     CDFlag = 0x00008000 /* has com.apple.rootless.restricted-nvram-variables.heritable entitlement */

	RUNTIME CDFlag = 0x00010000 /* Apply hardened runtime policies */

	LINKER_SIGNED CDFlag = 0x20000 // type property

	ALLOWED_MACHO CDFlag = (ADHOC | HARD | KILL | CHECK_EXPIRATION | RESTRICT | ENFORCEMENT | REQUIRE_LV | RUNTIME)

	EXEC_SET_HARD        CDFlag = 0x00100000 /* set HARD on any exec'ed process */
	EXEC_SET_KILL        CDFlag = 0x00200000 /* set KILL on any exec'ed process */
	EXEC_SET_ENFORCEMENT CDFlag = 0x00400000 /* set ENFORCEMENT on any exec'ed process */
	EXEC_INHERIT_SIP     CDFlag = 0x00800000 /* set INSTALLER on any exec'ed process */

	KILLED          CDFlag = 0x01000000 /* was killed by kernel for invalidity */
	DYLD_PLATFORM   CDFlag = 0x02000000 /* dyld used to load this is a platform binary */
	PLATFORM_BINARY CDFlag = 0x04000000 /* this is a platform binary */
	PLATFORM_PATH   CDFlag = 0x08000000 /* platform binary by the fact of path (osx only) */

	DEBUGGED             CDFlag = 0x10000000 /* process is currently or has previously been debugged and allowed to run with invalid pages */
	SIGNED               CDFlag = 0x20000000 /* process has a signature (may have gone invalid) */
	DEV_CODE             CDFlag = 0x40000000 /* code is dev signed, cannot be loaded into prod signed code (will go away with rdar://problem/28322552) */
	DATAVAULT_CONTROLLER CDFlag = 0x80000000 /* has Data Vault controller entitlement */

	ENTITLEMENT_FLAGS CDFlag = (GET_TASK_ALLOW | INSTALLER | DATAVAULT_CONTROLLER | NVRAM_UNRESTRICTED)
)

func (f CDFlag) String() string {
	var out []string
	if f == NONE {
		out = append(out, "none")
	}
	if (f & VALID) != 0 {
		out = append(out, "valid")
	}
	if (f & ADHOC) != 0 {
		out = append(out, "adhoc")
	}
	if (f & GET_TASK_ALLOW) != 0 {
		out = append(out, "get-task-allow")
	}
	if (f & INSTALLER) != 0 {
		out = append(out, "installer")
	}
	if (f & FORCED_LV) != 0 {
		out = append(out, "forced-lv")
	}
	if (f & INVALID_ALLOWED) != 0 {
		out = append(out, "invalid-allowed")
	}
	if (f & HARD) != 0 {
		out = append(out, "hard")
	}
	if (f & KILL) != 0 {
		out = append(out, "kill")
	}
	if (f & CHECK_EXPIRATION) != 0 {
		out = append(out, "check-expiration")
	}
	if (f & RESTRICT) != 0 {
		out = append(out, "restrict")
	}
	if (f & ENFORCEMENT) != 0 {
		out = append(out, "enforcement")
	}
	if (f & REQUIRE_LV) != 0 {
		out = append(out, "require-lv")
	}
	if (f & ENTITLEMENTS_VALIDATED) != 0 {
		out = append(out, "entitlements-validated")
	}
	if (f & NVRAM_UNRESTRICTED) != 0 {
		out = append(out, "nvram-unrestricted")
	}
	if (f & RUNTIME) != 0 {
		out = append(out, "runtime")
	}
	if (f & LINKER_SIGNED) != 0 {
		out = append(out, "linker-signed")
	}
	// if (f & ALLOWED_MACHO) != 0 {
	// 	out = append(out, "allowed-macho")
	// }
	if (f & EXEC_SET_HARD) != 0 {
		out = append(out, "exec-set-hard")
	}
	if (f & EXEC_SET_KILL) != 0 {
		out = append(out, "exec-set-kill")
	}
	if (f & EXEC_SET_ENFORCEMENT) != 0 {
		out = append(out, "exec-set-enforcement")
	}
	if (f & EXEC_INHERIT_SIP) != 0 {
		out = append(out, "exec-inherit-sip")
	}
	if (f & KILLED) != 0 {
		out = append(out, "killed")
	}
	if (f & DYLD_PLATFORM) != 0 {
		out = append(out, "dyld-platform")
	}
	if (f & PLATFORM_BINARY) != 0 {
		out = append(out, "platform-binary")
	}
	if (f & PLATFORM_PATH) != 0 {
		out = append(out, "platform-path")
	}
	if (f & DEBUGGED) != 0 {
		out = append(out, "debugged")
	}
	if (f & SIGNED) != 0 {
		out = append(out, "signed")
	}
	if (f & DEV_CODE) != 0 {
		out = append(out, "dev-code")
	}
	if (f & DATAVAULT_CONTROLLER) != 0 {
		out = append(out, "datavault-controller")
	}
	return strings.Join(out, ", ")
}

type cdPlatform uint8

const (
	// A signature with a nonzero platform identifier value, when endorsed as originated by Apple,
	// identifies code as belonging to a particular operating system deliverable set. Some system
	// components restrict functionality to platform binaries. The actual values are arbitrary.
	NON_PLATFORM_BINARY cdPlatform = 0
	// PLATFORM_PLATFORM_BINARY cdPlatform = 0xE // TODO: this is what /bin/ls's platform is, but are there other values?
)

func (p cdPlatform) String() string {
	if p == NON_PLATFORM_BINARY {
		return "non-platform-binary"
	}
	return "platform-binary"
}

// CodeDirectoryType header
type CodeDirectoryType struct {
	CdEarliest
	CdScatter
	CdTeamID
	CdCodeLimit64
	CdExecSeg
	CdRuntime
	CdLinkage
	/* followed by dynamic content as located by offset fields above */
}

type CdEarliest struct {
	Version       cdVersion  `json:"version,omitempty"`         // compatibility version
	Flags         CDFlag     `json:"flags,omitempty"`           // setup and mode flags
	HashOffset    uint32     `json:"hash_offset,omitempty"`     // offset of hash slot element at index zero
	IdentOffset   uint32     `json:"ident_offset,omitempty"`    // offset of identifier string
	NSpecialSlots uint32     `json:"n_special_slots,omitempty"` // number of special hash slots
	NCodeSlots    uint32     `json:"n_code_slots,omitempty"`    // number of ordinary (code) hash slots
	CodeLimit     uint32     `json:"code_limit,omitempty"`      // limit to main image signature range
	HashSize      uint8      `json:"hash_size,omitempty"`       // size of each hash in bytes
	HashType      hashType   `json:"hash_type,omitempty"`       // type of hash (cdHashType* constants)
	Platform      cdPlatform `json:"platform,omitempty"`        // platform identifier zero if not platform binary
	PageSize      uint8      `json:"page_size,omitempty"`       // log2(page size in bytes) 0 => infinite
	_             uint32     // unused (must be zero)
}

type CdScatter struct {
	/* Version 0x20100 */
	ScatterOffset uint32 `json:"scatter_offset,omitempty"` /* offset of optional scatter vector */
}

type CdTeamID struct {
	/* Version 0x20200 */
	TeamOffset uint32 `json:"team_offset,omitempty"` /* offset of optional team identifier */
}

type CdCodeLimit64 struct {
	/* Version 0x20300 */
	_           uint32 /* unused (must be zero) */
	CodeLimit64 uint64 `json:"code_limit_64,omitempty"` /* limit to main image signature range, 64 bits */
}

type CdExecSeg struct {
	/* Version 0x20400 */
	ExecSegBase  uint64      `json:"exec_seg_base,omitempty"`  /* offset of executable segment */
	ExecSegLimit uint64      `json:"exec_seg_limit,omitempty"` /* limit of executable segment */
	ExecSegFlags execSegFlag `json:"exec_seg_flags,omitempty"` /* exec segment flags */
}

type CdRuntime struct {
	/* Version 0x20500 */
	Runtime          mtypes.Version `json:"runtime,omitempty"`            // Runtime version
	PreEncryptOffset uint32         `json:"pre_encrypt_offset,omitempty"` // offset of pre-encrypt hash slots
}

type CdLinkage struct {
	/* Version 0x20600 */
	LinkageHashType           uint8  `json:"linkage_hash_type,omitempty"`
	LinkageApplicationType    uint8  `json:"linkage_application_type,omitempty"`
	LinkageApplicationSubType uint16 `json:"linkage_application_sub_type,omitempty"`
	LinkageOffset             uint32 `json:"linkage_offset,omitempty"`
	LinkageSize               uint32 `json:"linkage_size,omitempty"`
}

// Scatter object
type Scatter struct {
	Count        uint32 `json:"count,omitempty"`         // number of pages zero for sentinel (only)
	Base         uint32 `json:"base,omitempty"`          // first page number
	TargetOffset uint64 `json:"target_offset,omitempty"` // byte offset in target
	_            uint64 // reserved (must be zero)
}

type execSegFlag uint64

/* executable segment flags */
const (
	EXECSEG_MAIN_BINARY     execSegFlag = 0x01  /* executable segment denotes main binary */
	EXECSEG_ALLOW_UNSIGNED  execSegFlag = 0x10  /* allow unsigned pages (for debugging) */
	EXECSEG_DEBUGGER        execSegFlag = 0x20  /* main binary is debugger */
	EXECSEG_JIT             execSegFlag = 0x40  /* JIT enabled */
	EXECSEG_SKIP_LV         execSegFlag = 0x80  /* OBSOLETE: skip library validation */
	EXECSEG_CAN_LOAD_CDHASH execSegFlag = 0x100 /* can bless cdhash for execution */
	EXECSEG_CAN_EXEC_CDHASH execSegFlag = 0x200 /* can execute blessed cdhash */
)

func (f execSegFlag) String() string {
	switch f {
	case EXECSEG_MAIN_BINARY:
		return "Main Binary"
	case EXECSEG_ALLOW_UNSIGNED:
		return "Allow Unsigned"
	case EXECSEG_DEBUGGER:
		return "Debugger"
	case EXECSEG_JIT:
		return "JIT"
	case EXECSEG_SKIP_LV:
		return "Skip LV"
	case EXECSEG_CAN_LOAD_CDHASH:
		return "Can Load CDHash"
	case EXECSEG_CAN_EXEC_CDHASH:
		return "Can Exec CDHash"
	default:
		return fmt.Sprintf("execSegFlag(%#x)", uint64(f))
	}
}
