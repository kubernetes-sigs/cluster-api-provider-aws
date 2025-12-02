package macho

// Definitions From: https://github.com/Apple-FOSS-Mirror/Security/blob/5bcad85836c8bbb383f660aaf25b555a805a48e4/OSX/sec/Security/Tool/codesign.c#L53-L89

const (
	EarliestVersion     CdVersion = 0x20001
	SupportsScatter     CdVersion = 0x20100
	SupportsTeamid      CdVersion = 0x20200
	SupportsCodelimit64 CdVersion = 0x20300
	SupportsExecseg     CdVersion = 0x20400
	SupportsRuntime     CdVersion = 0x20500
	SupportsLinkage     CdVersion = 0x20600
	CompatibilityLimit  CdVersion = 0x2F000 // "version 3 with wiggle room"
)

const (
	// code signing attributes of a process
	None         CdFlag = 0x00000000 // no flags
	Valid        CdFlag = 0x00000001 // dynamically valid
	Adhoc        CdFlag = 0x00000002 // ad hoc signed
	GetTaskAllow CdFlag = 0x00000004 // has get-task-allow entitlement
	Installer    CdFlag = 0x00000008 // has installer entitlement

	ForcedLv       CdFlag = 0x00000010 // Library Validation required by Hardened System Policy
	InvalidAllowed CdFlag = 0x00000020 // (macOS Only) Page invalidation allowed by task port policy

	Hard            CdFlag = 0x00000100 // don't load invalid pages
	Kill            CdFlag = 0x00000200 // kill process if it becomes invalid
	CheckExpiration CdFlag = 0x00000400 // force expiration checking
	Restrict        CdFlag = 0x00000800 // tell dyld to treat restricted

	Enforcement           CdFlag = 0x00001000 // require enforcement
	RequireLv             CdFlag = 0x00002000 // require library validation
	EntitlementsValidated CdFlag = 0x00004000 // code signature permits restricted entitlements
	NvramUnrestricted     CdFlag = 0x00008000 // has com.apple.rootless.restricted-nvram-variables.heritable entitlement

	Runtime CdFlag = 0x00010000 // Apply hardened runtime policies

	LinkerSigned CdFlag = 0x20000 // type property

	AllowedMacho CdFlag = (Adhoc | Hard | Kill | CheckExpiration | Restrict | Enforcement | RequireLv | Runtime)

	ExecSetHard        CdFlag = 0x00100000 // set HARD on any exec'ed process
	ExecSetKill        CdFlag = 0x00200000 // set KILL on any exec'ed process
	ExecSetEnforcement CdFlag = 0x00400000 // set ENFORCEMENT on any exec'ed process
	ExecInheritSIP     CdFlag = 0x00800000 // set INSTALLER on any exec'ed process

	Killed         CdFlag = 0x01000000 // was killed by kernel for invalidity
	DyldPlatform   CdFlag = 0x02000000 // dyld used to load this is a platform binary
	PlatformBinary CdFlag = 0x04000000 // this is a platform binary
	PlatformPath   CdFlag = 0x08000000 // platform binary by the fact of path (osx only)

	Debugged            CdFlag = 0x10000000 // process is currently or has previously been debugged and allowed to run with invalid pages
	Signed              CdFlag = 0x20000000 // process has a signature (may have gone invalid)
	DevCode             CdFlag = 0x40000000 // code is dev signed, cannot be loaded into prod signed code (will go away with rdar://problem/28322552)
	DatavaultController CdFlag = 0x80000000 // has Data Vault controller entitlement

	EntitlementFlags CdFlag = (GetTaskAllow | Installer | DatavaultController | NvramUnrestricted)
)

// executable segment flags
const (
	ExecsegMainBinary    ExecSegFlag = 0x1   // executable segment denotes main binary
	ExecsegAllowUnsigned ExecSegFlag = 0x10  // allow unsigned pages (for debugging)
	ExecsegDebugger      ExecSegFlag = 0x20  // main binary is debugger
	ExecsegJit           ExecSegFlag = 0x40  // JIT enabled
	ExecsegSkipLv        ExecSegFlag = 0x80  // OBSOLETE: skip library validation
	ExecsegCanLoadCdhash ExecSegFlag = 0x100 // can bless cdhash for execution
	ExecsegCanExecCdhash ExecSegFlag = 0x200 // can execute blessed cdhash
)

type CdVersion uint32
type CdFlag uint32
type ExecSegFlag uint64

type CodeDirectory struct {
	CodeDirectoryHeader
	// followed by dynamic content as located by offset fields above
	Payload []byte
}

type CodeDirectoryHeader struct {
	Version       CdVersion // compatibility version
	Flags         CdFlag    // setup and mode flags
	HashOffset    uint32    // offset of hash slot element at index zero
	IdentOffset   uint32    // offset of identifier string
	NSpecialSlots uint32    // number of special hash slots
	NCodeSlots    uint32    // number of ordinary (code) hash slots
	CodeLimit     uint32    // limit to main image signature range
	HashSize      uint8     // size of each hash in bytes
	HashType      HashType  // type of hash (cdHashType* constants)
	Platform      uint8     // platform identifier zero if not platform binary
	PageSize      uint8     // log2(page size in bytes) 0 => infinite
	Spare2        uint32    // unused (must be zero)

	EndEarliest [0]uint8

	// Version 0x20100
	ScatterOffset  uint32 // offset of optional scatter vector
	EndWithScatter [0]uint8

	// Version 0x20200
	TeamOffset  uint32 // offset of optional team identifier
	EndWithTeam [0]uint8

	// Version 0x20300
	Spare3             uint32 // unused (must be zero)
	CodeLimit64        uint64 // limit to main image signature range, 64 bits
	EndWithCodeLimit64 [0]uint8

	// Version 0x20400
	ExecSegBase  uint64      // offset of executable segment
	ExecSegLimit uint64      // limit of executable segment
	ExecSegFlags ExecSegFlag // exec segment flags

	// Version 0x20500
	Runtime          uint32 // Runtime version encoded as an unsigned int
	PreEncryptOffset uint32 // offset of pre-encrypt hash slots

	// Version 0x20600
	// TODO: linkage options
}
