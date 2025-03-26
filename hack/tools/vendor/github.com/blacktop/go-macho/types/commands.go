package types

//go:generate stringer -type=LoadCmd -output commands_string.go

import (
	"encoding/json"
	"fmt"
	"strings"
)

// A LoadCmd is a Mach-O load command.
type LoadCmd uint32

func (c LoadCmd) Command() LoadCmd { return c }

const (
	LC_REQ_DYLD       LoadCmd = 0x80000000
	LC_SEGMENT        LoadCmd = 0x1  // segment of this file to be mapped
	LC_SYMTAB         LoadCmd = 0x2  // link-edit stab symbol table info
	LC_SYMSEG         LoadCmd = 0x3  // link-edit gdb symbol table info (obsolete)
	LC_THREAD         LoadCmd = 0x4  // thread
	LC_UNIXTHREAD     LoadCmd = 0x5  // thread+stack
	LC_LOADFVMLIB     LoadCmd = 0x6  // load a specified fixed VM shared library
	LC_IDFVMLIB       LoadCmd = 0x7  // fixed VM shared library identification
	LC_IDENT          LoadCmd = 0x8  // object identification info (obsolete)
	LC_FVMFILE        LoadCmd = 0x9  // fixed VM file inclusion (internal use)
	LC_PREPAGE        LoadCmd = 0xa  // prepage command (internal use)
	LC_DYSYMTAB       LoadCmd = 0xb  // dynamic link-edit symbol table info
	LC_LOAD_DYLIB     LoadCmd = 0xc  // load dylib command
	LC_ID_DYLIB       LoadCmd = 0xd  // id dylib command
	LC_LOAD_DYLINKER  LoadCmd = 0xe  // load a dynamic linker
	LC_ID_DYLINKER    LoadCmd = 0xf  // id dylinker command (not load dylinker command)
	LC_PREBOUND_DYLIB LoadCmd = 0x10 // modules prebound for a dynamically linked shared library
	LC_ROUTINES       LoadCmd = 0x11 // image routines
	LC_SUB_FRAMEWORK  LoadCmd = 0x12 // sub framework
	LC_SUB_UMBRELLA   LoadCmd = 0x13 // sub umbrella
	LC_SUB_CLIENT     LoadCmd = 0x14 // sub client
	LC_SUB_LIBRARY    LoadCmd = 0x15 // sub library
	LC_TWOLEVEL_HINTS LoadCmd = 0x16 // two-level namespace lookup hints
	LC_PREBIND_CKSUM  LoadCmd = 0x17 // prebind checksum
	/*
	 * load a dynamically linked shared library that is allowed to be missing
	 * (all symbols are weak imported).
	 */
	LC_LOAD_WEAK_DYLIB          LoadCmd = (0x18 | LC_REQ_DYLD)
	LC_SEGMENT_64               LoadCmd = 0x19                 // 64-bit segment of this file to be mapped
	LC_ROUTINES_64              LoadCmd = 0x1a                 // 64-bit image routines
	LC_UUID                     LoadCmd = 0x1b                 // the uuid
	LC_RPATH                    LoadCmd = (0x1c | LC_REQ_DYLD) // runpath additions
	LC_CODE_SIGNATURE           LoadCmd = 0x1d                 // local of code signature
	LC_SEGMENT_SPLIT_INFO       LoadCmd = 0x1e                 // local of info to split segments
	LC_REEXPORT_DYLIB           LoadCmd = (0x1f | LC_REQ_DYLD) // load and re-export dylib
	LC_LAZY_LOAD_DYLIB          LoadCmd = 0x20                 // delay load of dylib until first use
	LC_ENCRYPTION_INFO          LoadCmd = 0x21                 // encrypted segment information
	LC_DYLD_INFO                LoadCmd = 0x22                 // compressed dyld information
	LC_DYLD_INFO_ONLY           LoadCmd = (0x22 | LC_REQ_DYLD) // compressed dyld information only
	LC_LOAD_UPWARD_DYLIB        LoadCmd = (0x23 | LC_REQ_DYLD) // load upward dylib
	LC_VERSION_MIN_MACOSX       LoadCmd = 0x24                 // build for MacOSX min OS version
	LC_VERSION_MIN_IPHONEOS     LoadCmd = 0x25                 // build for iPhoneOS min OS version
	LC_FUNCTION_STARTS          LoadCmd = 0x26                 // compressed table of function start addresses
	LC_DYLD_ENVIRONMENT         LoadCmd = 0x27                 // string for dyld to treat like environment variable
	LC_MAIN                     LoadCmd = (0x28 | LC_REQ_DYLD) // replacement for LC_UNIXTHREAD
	LC_DATA_IN_CODE             LoadCmd = 0x29                 // table of non-instructions in __text
	LC_SOURCE_VERSION           LoadCmd = 0x2A                 // source version used to build binary
	LC_DYLIB_CODE_SIGN_DRS      LoadCmd = 0x2B                 // Code signing DRs copied from linked dylibs
	LC_ENCRYPTION_INFO_64       LoadCmd = 0x2C                 // 64-bit encrypted segment information
	LC_LINKER_OPTION            LoadCmd = 0x2D                 // linker options in MH_OBJECT files
	LC_LINKER_OPTIMIZATION_HINT LoadCmd = 0x2E                 // optimization hints in MH_OBJECT files
	LC_VERSION_MIN_TVOS         LoadCmd = 0x2F                 // build for AppleTV min OS version
	LC_VERSION_MIN_WATCHOS      LoadCmd = 0x30                 // build for Watch min OS version
	LC_NOTE                     LoadCmd = 0x31                 // arbitrary data included within a Mach-O file
	LC_BUILD_VERSION            LoadCmd = 0x32                 // build for platform min OS version
	LC_DYLD_EXPORTS_TRIE        LoadCmd = (0x33 | LC_REQ_DYLD) // used with linkedit_data_command, payload is trie
	LC_DYLD_CHAINED_FIXUPS      LoadCmd = (0x34 | LC_REQ_DYLD) // used with linkedit_data_command
	LC_FILESET_ENTRY            LoadCmd = (0x35 | LC_REQ_DYLD) /* used with fileset_entry_command */
	LC_ATOM_INFO                LoadCmd = 0x36                 /* used with linkedit_data_command */
)

type SegFlag uint32

/* Constants for the flags field of the segment_command */
const (
	HighVM SegFlag = 0x1 /* the file contents for this segment is for
	   the high part of the VM space, the low part
	   is zero filled (for stacks in core files) */
	FvmLib SegFlag = 0x2 /* this segment is the VM that is allocated by
	   a fixed VM library, for overlap checking in
	   the link editor */
	NoReLoc SegFlag = 0x4 /* this segment has nothing that was relocated
	   in it and nothing relocated to it, that is
	   it maybe safely replaced without relocation*/
	ProtectedVersion1 SegFlag = 0x8 /* This segment is protected.  If the
	   segment starts at file offset 0, the
	   first page of the segment is not
	   protected.  All other pages of the
	   segment are protected. */
	ReadOnly SegFlag = 0x10 /* This segment is made read-only after fixups */
)

func (s SegFlag) List() []string {
	var flags []string
	if (s & HighVM) != 0 {
		flags = append(flags, "HighVM")
	}
	if (s & FvmLib) != 0 {
		flags = append(flags, "FixedVMLib")
	}
	if (s & NoReLoc) != 0 {
		flags = append(flags, "NoReLoc")
	}
	if (s & ProtectedVersion1) != 0 {
		flags = append(flags, "ProtectedV1")
	}
	if (s & ReadOnly) != 0 {
		flags = append(flags, "ReadOnly")
	}
	return flags
}

func (s SegFlag) String() string {
	var fStr string
	for _, attr := range s.List() {
		fStr += fmt.Sprintf("%s|", attr)
	}
	return strings.TrimSuffix(fStr, "|")
}

/*
 * Segment32 is a 32-bit segment load command indicates that a part of this file is to be
 * mapped into the task's address space.  The size of this segment in memory,
 * vmsize, maybe equal to or larger than the amount to map from this file,
 * filesize.  The file is mapped starting at fileoff to the beginning of
 * the segment in memory, vmaddr.  The rest of the memory of the segment,
 * if any, is allocated zero fill on demand.  The segment's maximum virtual
 * memory protection and initial virtual memory protection are specified
 * by the maxprot and initprot fields.  If the segment has sections then the
 * section structures directly follow the segment command and their size is
 * reflected in cmdsize.
 */
type Segment32 struct {
	LoadCmd              /* LC_SEGMENT */
	Len     uint32       /* includes sizeof section structs */
	Name    [16]byte     /* segment name */
	Addr    uint32       /* memory address of this segment */
	Memsz   uint32       /* memory size of this segment */
	Offset  uint32       /* file offset of this segment */
	Filesz  uint32       /* amount to map from the file */
	Maxprot VmProtection /* maximum VM protection */
	Prot    VmProtection /* initial VM protection */
	Nsect   uint32       /* number of sections in segment */
	Flag    SegFlag      /* flags */
}

/*
 * Segment64 is a 64-bit segment load command indicates that a part of this file is to be
 * mapped into a 64-bit task's address space.  If the 64-bit segment has
 * sections then section_64 structures directly follow the 64-bit segment
 * command and their size is reflected in cmdsize.
 */
type Segment64 struct {
	LoadCmd              /* LC_SEGMENT_64 */
	Len     uint32       /* includes sizeof section_64 structs */
	Name    [16]byte     /* segment name */
	Addr    uint64       /* memory address of this segment */
	Memsz   uint64       /* memory size of this segment */
	Offset  uint64       /* file offset of this segment */
	Filesz  uint64       /* amount to map from the file */
	Maxprot VmProtection /* maximum VM protection */
	Prot    VmProtection /* initial VM protection */
	Nsect   uint32       /* number of sections in segment */
	Flag    SegFlag      /* flags */
}

/*
 * LoadFvmLibCmd a fixed virtual shared library (filetype == MH_FVMLIB in the mach header)
 * contains a fvmlib_command (cmd == LC_IDFVMLIB) to identify the library.
 * An object that uses a fixed virtual shared library also contains a
 * fvmlib_command (cmd == LC_LOADFVMLIB) for each library it uses.
 * (THIS IS OBSOLETE and no longer supported).
 */
type LoadFvmLibCmd struct {
	LoadCmd        // LC_IDFVMLIB or LC_LOADFVMLIB
	Len     uint32 /* includes pathname string */
	/*
	 * Fixed virtual memory shared libraries are identified by two things.  The
	 * target pathname (the name of the library as found for execution), and the
	 * minor version number.  The address of where the headers are loaded is in
	 * header_addr. (THIS IS OBSOLETE and no longer supported).
	 */
	NameOffset   uint32  // library's target pathname
	MinorVersion Version /* library's minor version number */
	HeaderAddr   uint32  /* library's header address */
}

// A IDFvmLibCmd is a Mach-O fixed VM shared library identification command.
type IDFvmLibCmd LoadFvmLibCmd // LC_IDFVMLIB

/*
 * DylibCmd a dynamically linked shared library (filetype == MH_DYLIB in the mach header)
 * contains a dylib_command (cmd == LC_ID_DYLIB) to identify the library.
 * An object that uses a dynamically linked shared library also contains a
 * dylib_command (cmd == LC_LOAD_DYLIB, LC_LOAD_WEAK_DYLIB, or
 * LC_REEXPORT_DYLIB) for each library it uses.
 */
type DylibCmd struct {
	LoadCmd        /* LC_ID_DYLIB, LC_LOAD_{,WEAK_}DYLIB, LC_REEXPORT_DYLIB */
	Len     uint32 /* includes pathname string */
	/*
	 * Dynamicly linked shared libraries are identified by two things.  The
	 * pathname (the name of the library as found for execution), and the
	 * compatibility version number.  The pathname must match and the compatibility
	 * number in the user of the library must be greater than or equal to the
	 * library being used.  The time stamp is used to record the time a library was
	 * built and copied into user so it can be use to determined if the library used
	 * at runtime is exactly the same as used to built the program.
	 */
	NameOffset     uint32
	Timestamp      uint32
	CurrentVersion Version
	CompatVersion  Version
}

// A LoadDylibCmd load a dynamically linked shared library.
type LoadDylibCmd DylibCmd // LC_LOAD_DYLIB
// A IDDylibCmd represents a Mach-O load dynamic library ident command.
type IDDylibCmd DylibCmd // LC_ID_DYLIB
// A LoadWeakDylibCmd is a Mach-O load a dynamically linked shared library that is allowed to be missing (all symbols are weak imported) command.
type LoadWeakDylibCmd DylibCmd // LC_LOAD_WEAK_DYLIB
// A ReExportDylibCmd is a Mach-O load and re-export dylib command.
type ReExportDylibCmd DylibCmd // LC_REEXPORT_DYLIB
// A LazyLoadDylibCmd is a Mach-O delay load of dylib until first use command.
type LazyLoadDylibCmd DylibCmd // LC_LAZY_LOAD_DYLIB
// A LoadUpwardDylibCmd is a Mach-O load upward dylibcommand.
type LoadUpwardDylibCmd DylibCmd // LC_LOAD_UPWARD_DYLIB

/*
 * SubFrameworkCmd a dynamically linked shared library may be a subframework of an umbrella
 * framework.  If so it will be linked with "-umbrella umbrella_name" where
 * Where "umbrella_name" is the name of the umbrella framework. A subframework
 * can only be linked against by its umbrella framework or other subframeworks
 * that are part of the same umbrella framework.  Otherwise the static link
 * editor produces an error and states to link against the umbrella framework.
 * The name of the umbrella framework for subframeworks is recorded in the
 * following structure.
 */
type SubFrameworkCmd struct {
	LoadCmd                // LC_SUB_FRAMEWORK
	Len             uint32 /* includes umbrella string */
	FrameworkOffset uint32 /* the umbrella framework name */
}

/*
 * SubClientCmd for dynamically linked shared libraries that are subframework of an umbrella
 * framework they can allow clients other than the umbrella framework or other
 * subframeworks in the same umbrella framework.  To do this the subframework
 * is built with "-allowable_client client_name" and an LC_SUB_CLIENT load
 * command is created for each -allowable_client flag.  The client_name is
 * usually a framework name.  It can also be a name used for bundles clients
 * where the bundle is built with "-client_name client_name".
 */
type SubClientCmd struct {
	LoadCmd             // LC_SUB_CLIENT
	Len          uint32 /* includes client string */
	ClientOffset uint32 /* the client name */
}

/*
 * SubUmbrellaCmd a dynamically linked shared library may be a sub_umbrella of an umbrella
 * framework.  If so it will be linked with "-sub_umbrella umbrella_name" where
 * Where "umbrella_name" is the name of the sub_umbrella framework.  When
 * staticly linking when -twolevel_namespace is in effect a twolevel namespace
 * umbrella framework will only cause its subframeworks and those frameworks
 * listed as sub_umbrella frameworks to be implicited linked in.  Any other
 * dependent dynamic libraries will not be linked it when -twolevel_namespace
 * is in effect.  The primary library recorded by the static linker when
 * resolving a symbol in these libraries will be the umbrella framework.
 * Zero or more sub_umbrella frameworks may be use by an umbrella framework.
 * The name of a sub_umbrella framework is recorded in the following structure.
 */
type SubUmbrellaCmd struct {
	LoadCmd               // LC_SUB_UMBRELLA
	Len            uint32 /* includes sub_umbrella string */
	UmbrellaOffset uint32 /* the sub_umbrella framework name */
}

/*
 * SubLibraryCmd a dynamically linked shared library may be a sub_library of another shared
 * library.  If so it will be linked with "-sub_library library_name" where
 * Where "library_name" is the name of the sub_library shared library.  When
 * staticly linking when -twolevel_namespace is in effect a twolevel namespace
 * shared library will only cause its subframeworks and those frameworks
 * listed as sub_umbrella frameworks and libraries listed as sub_libraries to
 * be implicited linked in.  Any other dependent dynamic libraries will not be
 * linked it when -twolevel_namespace is in effect.  The primary library
 * recorded by the static linker when resolving a symbol in these libraries
 * will be the umbrella framework (or dynamic library). Zero or more sub_library
 * shared libraries may be use by an umbrella framework or (or dynamic library).
 * The name of a sub_library framework is recorded in the following structure.
 * For example /usr/lib/libobjc_profile.A.dylib would be recorded as "libobjc".
 */
type SubLibraryCmd struct {
	LoadCmd              // LC_SUB_LIBRARY
	Len           uint32 /* includes sub_library string */
	LibraryOffset uint32 /* the sub_library name */
}

/*
 * PreboundDylibCmd a program (filetype == MH_EXECUTE) that is
 * prebound to its dynamic libraries has one of these for each library that
 * the static linker used in prebinding.  It contains a bit vector for the
 * modules in the library.  The bits indicate which modules are bound (1) and
 * which are not (0) from the library.  The bit for module 0 is the low bit
 * of the first byte.  So the bit for the Nth module is:
 * (linked_modules[N/8] >> N%8) & 1
 */
type PreboundDylibCmd struct {
	LoadCmd                    // LC_PREBOUND_DYLIB
	Len                 uint32 /* includes strings */
	NameOffset          uint32 // library's path name
	NumModules          uint32 // number of modules in library
	LinkedModulesOffset uint32 // bit vector of linked modules
}

/*
 * DylinkerCmd a program that uses a dynamic linker contains a dylinker_command to identify
 * the name of the dynamic linker (LC_LOAD_DYLINKER).  And a dynamic linker
 * contains a dylinker_command to identify the dynamic linker (LC_ID_DYLINKER).
 * A file can have at most one of these.
 * This struct is also used for the LC_DYLD_ENVIRONMENT load command and
 * contains string for dyld to treat like environment variable.
 */
type DylinkerCmd struct {
	LoadCmd           // LC_ID_DYLINKER, LC_LOAD_DYLINKER or LC_DYLD_ENVIRONMENT
	Len        uint32 // includes pathname string
	NameOffset uint32 // dynamic linker's path name
}

// A IDDylinkerCmd is a Mach-O dynamic linker identification command.
type IDDylinkerCmd DylinkerCmd // LC_ID_DYLINKER
// A DyldEnvironmentCmd is a Mach-O string for dyld to treat like environment variable command.
type DyldEnvironmentCmd DylinkerCmd // LC_DYLD_ENVIRONMENT

type ThreadFlavor uint32

const (
	//
	// x86 flavors
	//
	X86_THREAD_STATE32    ThreadFlavor = 1
	X86_FLOAT_STATE32     ThreadFlavor = 2
	X86_EXCEPTION_STATE32 ThreadFlavor = 3
	X86_THREAD_STATE64    ThreadFlavor = 4
	X86_FLOAT_STATE64     ThreadFlavor = 5
	X86_EXCEPTION_STATE64 ThreadFlavor = 6
	X86_THREAD_STATE      ThreadFlavor = 7
	X86_FLOAT_STATE       ThreadFlavor = 8
	X86_EXCEPTION_STATE   ThreadFlavor = 9
	X86_DEBUG_STATE32     ThreadFlavor = 10
	X86_DEBUG_STATE64     ThreadFlavor = 11
	X86_DEBUG_STATE       ThreadFlavor = 12
	X86_THREAD_STATE_NONE ThreadFlavor = 13
	/* 14 and 15 are used for the internal X86_SAVED_STATE flavours */
	/* Arrange for flavors to take sequential values, 32-bit, 64-bit, non-specific */
	X86_AVX_STATE32         ThreadFlavor = 16
	X86_AVX_STATE64         ThreadFlavor = (X86_AVX_STATE32 + 1)
	X86_AVX_STATE           ThreadFlavor = (X86_AVX_STATE32 + 2)
	X86_AVX512_STATE32      ThreadFlavor = 19
	X86_AVX512_STATE64      ThreadFlavor = (X86_AVX512_STATE32 + 1)
	X86_AVX512_STATE        ThreadFlavor = (X86_AVX512_STATE32 + 2)
	X86_PAGEIN_STATE        ThreadFlavor = 22
	X86_THREAD_FULL_STATE64 ThreadFlavor = 23
	X86_INSTRUCTION_STATE   ThreadFlavor = 24
	X86_LAST_BRANCH_STATE   ThreadFlavor = 25
	//
	// arm flavors
	//
	ARM_THREAD_STATE         ThreadFlavor = 1
	ARM_UNIFIED_THREAD_STATE ThreadFlavor = ARM_THREAD_STATE
	ARM_VFP_STATE            ThreadFlavor = 2
	ARM_EXCEPTION_STATE      ThreadFlavor = 3
	ARM_DEBUG_STATE          ThreadFlavor = 4 /* pre-armv8 */
	ARM_THREAD_STATE_NONE    ThreadFlavor = 5
	ARM_THREAD_STATE64       ThreadFlavor = 6
	ARM_EXCEPTION_STATE64    ThreadFlavor = 7
	//      ARM_THREAD_STATE_LAST    8 /* legacy */
	ARM_THREAD_STATE32 ThreadFlavor = 9
	/* API */
	ARM_DEBUG_STATE32 ThreadFlavor = 14
	ARM_DEBUG_STATE64 ThreadFlavor = 15
	ARM_NEON_STATE    ThreadFlavor = 16
	ARM_NEON_STATE64  ThreadFlavor = 17
	ARM_CPMU_STATE64  ThreadFlavor = 18
	ARM_PAGEIN_STATE  ThreadFlavor = 27
)

type ThreadState struct {
	Flavor ThreadFlavor // flavor of thread state
	Count  uint32       // count of 's in thread state
	Data   []byte       // thread state for this flavor
}

/*
 * ThreadCmd contain machine-specific data structures suitable for
 * use in the thread state primitives.  The machine specific data structures
 * follow the struct thread_command as follows.
 * Each flavor of machine specific data structure is preceded by an
 * constant for the flavor of that data structure, an  that is the
 * count of 's of the size of the state data structure and then
 * the state data structure follows.  This triple may be repeated for many
 * flavors.  The constants for the flavors, counts and state data structure
 * definitions are expected to be in the header file <machine/thread_status.h>.
 * These machine specific data structures sizes must be multiples of
 * 4 bytes.  The cmdsize reflects the total size of the thread_command
 * and all of the sizes of the constants for the flavors, counts and state
 * data structures.
 *
 * For executable objects that are unix processes there will be one
 * thread_command (cmd == LC_UNIXTHREAD) created for it by the link-editor.
 * This is the same as a LC_THREAD, except that a stack is automatically
 * created (based on the shell's limit for the stack size).  Command arguments
 * and environment variables are copied onto that stack.
 */
type ThreadCmd struct { // FIXME: handle all flavors ?
	LoadCmd        // LC_THREAD or  LC_UNIXTHREAD
	Len     uint32 // total size of this command
	/*  flavor		   flavor of thread state */
	/*  count		   count of 's in thread state */
	/* struct XXX_thread_state state   thread state for this flavor */
	/* ... */
}

// A UnixThreadCmd is a Mach-O unix thread command.
type UnixThreadCmd ThreadCmd

/*
 * RoutinesCmd contains the address of the dynamic shared library
 * initialization routine and an index into the module table for the module
 * that defines the routine.  Before any modules are used from the library the
 * dynamic linker fully binds the module that defines the initialization routine
 * and then calls it.  This gets called before any module initialization
 * routines (used for C++ static constructors) in the library.
 */
type RoutinesCmd struct { // for 32-bit architectures
	LoadCmd            // LC_ROUTINES
	Len         uint32 // total size of this command
	InitAddress uint32 //  address of initialization routine
	InitModule  uint32 // index into the module table that the init routine is defined in
	Reserved1   uint32
	Reserved2   uint32
	Reserved3   uint32
	Reserved4   uint32
	Reserved5   uint32
	Reserved6   uint32
}

// A Routines64Cmd is a Mach-O 64-bit version of RoutinesCmd
type Routines64Cmd struct {
	LoadCmd            // LC_ROUTINES_64
	Len         uint32 // total size of this command
	InitAddress uint64 // address of initialization routine
	InitModule  uint64 // index into the module table that the init routine is defined in
	Reserved1   uint64
	Reserved2   uint64
	Reserved3   uint64
	Reserved4   uint64
	Reserved5   uint64
	Reserved6   uint64
}

/*
 * SymtabCmd contains the offsets and sizes of the link-edit 4.3BSD
 * "stab" style symbol table information as described in the header files
 * <nlist.h> and <stab.h>.
 */
type SymtabCmd struct {
	LoadCmd        // LC_SYMTAB
	Len     uint32 // sizeof(struct symtab_command)
	Symoff  uint32 // symbol table offset
	Nsyms   uint32 // number of symbol table entries
	Stroff  uint32 // string table offset
	Strsize uint32 // string table size in bytes
}

/*
 * DysymtabCmd is the second set of the symbolic information which is used to support
 * the data structures for the dynamically link editor.
 *
 * The original set of symbolic information in the symtab_command which contains
 * the symbol and string tables must also be present when this load command is
 * present.  When this load command is present the symbol table is organized
 * into three groups of symbols:
 *	local symbols (static and debugging symbols) - grouped by module
 *	defined external symbols - grouped by module (sorted by name if not lib)
 *	undefined external symbols (sorted by name if MH_BINDATLOAD is not set,
 *	     			    and in order the were seen by the static
 *				    linker if MH_BINDATLOAD is set)
 * In this load command there are offsets and counts to each of the three groups
 * of symbols.
 *
 * This load command contains a the offsets and sizes of the following new
 * symbolic information tables:
 *	table of contents
 *	module table
 *	reference symbol table
 *	indirect symbol table
 * The first three tables above (the table of contents, module table and
 * reference symbol table) are only present if the file is a dynamically linked
 * shared library.  For executable and object modules, which are files
 * containing only one module, the information that would be in these three
 * tables is determined as follows:
 * 	table of contents - the defined external symbols are sorted by name
 *	module table - the file contains only one module so everything in the
 *		       file is part of the module.
 *	reference symbol table - is the defined and undefined external symbols
 *
 * For dynamically linked shared library files this load command also contains
 * offsets and sizes to the pool of relocation entries for all sections
 * separated into two groups:
 *	external relocation entries
 *	local relocation entries
 * For executable and object modules the relocation entries continue to hang
 * off the section structures.
 */
type DysymtabCmd struct {
	LoadCmd        // LC_DYSYMTAB
	Len     uint32 // sizeof(struct dysymtab_command)
	/*
	 * The symbols indicated by symoff and nsyms of the LC_SYMTAB load command
	 * are grouped into the following three groups:
	 *    local symbols (further grouped by the module they are from)
	 *    defined external symbols (further grouped by the module they are from)
	 *    undefined symbols
	 *
	 * The local symbols are used only for debugging.  The dynamic binding
	 * process may have to use them to indicate to the debugger the local
	 * symbols for a module that is being bound.
	 *
	 * The last two groups are used by the dynamic binding process to do the
	 * binding (indirectly through the module table and the reference symbol
	 * table when this is a dynamically linked shared library file).
	 */
	Ilocalsym uint32 // index to local symbols
	Nlocalsym uint32 // number of local symbols

	Iextdefsym uint32 // index to externally defined symbols
	Nextdefsym uint32 // number of externally defined symbols

	Iundefsym uint32 // index to undefined symbols
	Nundefsym uint32 // number of undefined symbols
	/*
	 * For the for the dynamic binding process to find which module a symbol
	 * is defined in the table of contents is used (analogous to the ranlib
	 * structure in an archive) which maps defined external symbols to modules
	 * they are defined in.  This exists only in a dynamically linked shared
	 * library file.  For executable and object modules the defined external
	 * symbols are sorted by name and is use as the table of contents.
	 */
	Tocoffset uint32 // file offset to table of contents
	Ntoc      uint32 // number of entries in table of contents
	/*
	 * To support dynamic binding of "modules" (whole object files) the symbol
	 * table must reflect the modules that the file was created from.  This is
	 * done by having a module table that has indexes and counts into the merged
	 * tables for each module.  The module structure that these two entries
	 * refer to is described below.  This exists only in a dynamically linked
	 * shared library file.  For executable and object modules the file only
	 * contains one module so everything in the file belongs to the module.
	 */
	Modtaboff uint32 // file offset to module table
	Nmodtab   uint32 // number of module table entries
	/*
	 * To support dynamic module binding the module structure for each module
	 * indicates the external references (defined and undefined) each module
	 * makes.  For each module there is an offset and a count into the
	 * reference symbol table for the symbols that the module references.
	 * This exists only in a dynamically linked shared library file.  For
	 * executable and object modules the defined external symbols and the
	 * undefined external symbols indicates the external references.
	 */
	Extrefsymoff uint32 // offset to referenced symbol table
	Nextrefsyms  uint32 // number of referenced symbol table entries
	/*
	 * The sections that contain "symbol pointers" and "routine stubs" have
	 * indexes and (implied counts based on the size of the section and fixed
	 * size of the entry) into the "indirect symbol" table for each pointer
	 * and stub.  For every section of these two types the index into the
	 * indirect symbol table is stored in the section header in the field
	 * reserved1.  An indirect symbol table entry is simply a 32bit index into
	 * the symbol table to the symbol that the pointer or stub is referring to.
	 * The indirect symbol table is ordered to match the entries in the section.
	 */
	Indirectsymoff uint32 // file offset to the indirect symbol table
	Nindirectsyms  uint32 // number of indirect symbol table entries
	/*
	 * To support relocating an individual module in a library file quickly the
	 * external relocation entries for each module in the library need to be
	 * accessed efficiently.  Since the relocation entries can't be accessed
	 * through the section headers for a library file they are separated into
	 * groups of local and external entries further grouped by module.  In this
	 * case the presents of this load command who's extreloff, nextrel,
	 * locreloff and nlocrel fields are non-zero indicates that the relocation
	 * entries of non-merged sections are not referenced through the section
	 * structures (and the reloff and nreloc fields in the section headers are
	 * set to zero).
	 *
	 * Since the relocation entries are not accessed through the section headers
	 * this requires the r_address field to be something other than a section
	 * offset to identify the item to be relocated.  In this case r_address is
	 * set to the offset from the vmaddr of the first LC_SEGMENT command.
	 * For MH_SPLIT_SEGS images r_address is set to the the offset from the
	 * vmaddr of the first read-write LC_SEGMENT command.
	 *
	 * The relocation entries are grouped by module and the module table
	 * entries have indexes and counts into them for the group of external
	 * relocation entries for that the module.
	 *
	 * For sections that are merged across modules there must not be any
	 * remaining external relocation entries for them (for merged sections
	 * remaining relocation entries must be local).
	 */
	Extreloff uint32 // offset to external relocation entries
	Nextrel   uint32 // number of external relocation entries
	/*
	 * All the local relocation entries are grouped together (they are not
	 * grouped by their module since they are only used if the object is moved
	 * from it staticly link edited address).
	 */
	Locreloff uint32 // offset to local relocation entries
	Nlocrel   uint32 // number of local relocation entries
}

const (
	/*
	 * An indirect symbol table entry is simply a 32bit index into the symbol table
	 * to the symbol that the pointer or stub is refering to.  Unless it is for a
	 * non-lazy symbol pointer section for a defined symbol which strip(1) as
	 * removed.  In which case it has the value INDIRECT_SYMBOL_LOCAL.  If the
	 * symbol was also absolute INDIRECT_SYMBOL_ABS is or'ed with that.
	 */
	INDIRECT_SYMBOL_LOCAL = 0x80000000 // TODO: use this ?
	INDIRECT_SYMBOL_ABS   = 0x40000000
)

/* a table of contents entry */
type DylibTableOfContents struct {
	SymbolIndex uint32 /* the defined external symbol (index into the symbol table) */
	ModuleIndex uint32 /* index into the module table this symbol is defined in */
}

/* a module table entry */
type DylibModule struct {
	ModuleName uint32 /* the module name (index into string table) */

	Iextdefsym uint32 /* index into externally defined symbols */
	Nextdefsym uint32 /* number of externally defined symbols */
	Irefsym    uint32 /* index into reference symbol table */
	Nrefsym    uint32 /* number of reference symbol table entries */
	Ilocalsym  uint32 /* index into symbols for local symbols */
	Nlocalsym  uint32 /* number of local symbols */

	Iextrel uint32 /* index into external relocation entries */
	Nextrel uint32 /* number of external relocation entries */

	IinitIterm uint32 /* low 16 bits are the index into the init
		   section, high 16 bits are the index into
	           the term section */
	NinitNterm uint32 /* low 16 bits are the number of init section
	   entries, high 16 bits are the number of
	   term section entries */

	/* for this module address of the start of */
	ObjcModuleInfoAddr uint32 /*  the (__OBJC,__module_info) section */
	/* for this module size of */
	ObjcModuleInfoSize uint32 /*  the (__OBJC,__module_info) section */
}

/* a 64-bit module table entry */
type DylibModule64 struct {
	ModuleName uint32 /* the module name (index into string table) */

	Iextdefsym uint32 /* index into externally defined symbols */
	Nextdefsym uint32 /* number of externally defined symbols */
	Irefsym    uint32 /* index into reference symbol table */
	Nrefsym    uint32 /* number of reference symbol table entries */
	Ilocalsym  uint32 /* index into symbols for local symbols */
	Nlocalsym  uint32 /* number of local symbols */

	Iextrel uint32 /* index into external relocation entries */
	Nextrel uint32 /* number of external relocation entries */

	IinitIterm uint32 /* low 16 bits are the index into the init
	   section, high 16 bits are the index into the term section */
	NinitNterm uint32 /* low 16 bits are the number of init section
	entries, high 16 bits are the number of term section entries */

	/* for this module size of */
	ObjcModuleInfoSize uint32 /*  the (__OBJC,__module_info) section */
	/* for this module address of the start of */
	ObjcModuleInfoAddr uint64 /*  the (__OBJC,__module_info) section */
}

/*
 * The entries in the reference symbol table are used when loading the module
 * (both by the static and dynamic link editors) and if the module is unloaded
 * or replaced.  Therefore all external symbols (defined and undefined) are
 * listed in the module's reference table.  The flags describe the type of
 * reference that is being made.  The constants for the flags are defined in
 * <mach-o/nlist.h> as they are also used for symbol table entries.
 */
// isym:24, /* index into the symbol table */
// flags:8; /* flags to indicate the type of reference */
type DylibReference uint32

func (d DylibReference) SymIndex() uint32 {
	return uint32(d >> 8)
}
func (d DylibReference) Flags() uint8 {
	return uint8(d & 0xff)
}

/*
 * TwolevelHintsCmd contains the offset and number of hints in the
 * two-level namespace lookup hints table.
 */
type TwolevelHintsCmd struct {
	LoadCmd         // LC_TWOLEVEL_HINTS
	Len      uint32 // sizeof(struct twolevel_hints_command)
	Offset   uint32 // offset to the hint table
	NumHints uint32 // number of hints in the hint table
}

/*
 * TwolevelHint provide hints to the dynamic link editor where to start
 * looking for an undefined symbol in a two-level namespace image.  The
 * isub_image field is an index into the sub-images (sub-frameworks and
 * sub-umbrellas list) that made up the two-level image that the undefined
 * symbol was found in when it was built by the static link editor.  If
 * isub-image is 0 the the symbol is expected to be defined in library and not
 * in the sub-images.  If isub-image is non-zero it is an index into the array
 * of sub-images for the umbrella with the first index in the sub-images being
 * 1. The array of sub-images is the ordered list of sub-images of the umbrella
 * that would be searched for a symbol that has the umbrella recorded as its
 * primary library.  The table of contents index is an index into the
 * library's table of contents.  This is used as the starting point of the
 * binary search or a directed linear search.
 */
type TwolevelHint uint32

// SubImageIndex index into the sub images
func (t TwolevelHint) SubImageIndex() uint8 {
	return uint8(t & 0xf)
}

// TableOfContentsIndex index into the table of contents
func (t TwolevelHint) TableOfContentsIndex() uint32 {
	return uint32(t >> 8)
}

func (t TwolevelHint) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SubImageIndex        uint8  `json:"subimage_index"`
		TableOfContentsIndex uint32 `json:"toc_index"`
	}{
		SubImageIndex:        t.SubImageIndex(),
		TableOfContentsIndex: t.TableOfContentsIndex(),
	})
}

/*
 * PrebindCksumCmd contains the value of the original check sum for
 * prebound files or zero.  When a prebound file is first created or modified
 * for other than updating its prebinding information the value of the check sum
 * is set to zero.  When the file has it prebinding re-done and if the value of
 * the check sum is zero the original check sum is calculated and stored in
 * cksum field of this load command in the output file.  If when the prebinding
 * is re-done and the cksum field is non-zero it is left unchanged from the
 * input file.
 */
type PrebindCksumCmd struct {
	LoadCmd         // LC_PREBIND_CKSUM
	Len      uint32 // sizeof(struct prebind_cksum_command)
	CheckSum uint32 // the check sum or zero
}

/*
 * UUIDCmd contains a single 128-bit unique random number that
 * identifies an object produced by the static link editor.
 */
type UUIDCmd struct {
	LoadCmd        // LC_UUID
	Len     uint32 // sizeof(struct uuid_command)
	UUID    UUID   // the 128-bit uuid
}

/*
 * RpathCmd contains a path which at runtime should be added to
 * the current run path used to find @rpath prefixed dylibs.
 */
type RpathCmd struct {
	LoadCmd           // LC_RPATH
	Len        uint32 // includes string
	PathOffset uint32 // path to add to run path
}

/*
 * LinkEditDataCmd contains the offsets and sizes of a blob
 * of data in the __LINKEDIT segment.
 */
type LinkEditDataCmd struct {
	LoadCmd /* LC_CODE_SIGNATURE, LC_SEGMENT_SPLIT_INFO,
	   LC_FUNCTION_STARTS, LC_DATA_IN_CODE,
	   LC_DYLIB_CODE_SIGN_DRS,
	   LC_ATOM_INFO,
	   LC_LINKER_OPTIMIZATION_HINT,
	   LC_DYLD_EXPORTS_TRIE, or
	   LC_DYLD_CHAINED_FIXUPS. */
	Len    uint32 // sizeof(struct linkedit_data_command)
	Offset uint32 // file offset of data in __LINKEDIT segment
	Size   uint32 // file size of data in __LINKEDIT segment
}

// A CodeSignatureCmd is a Mach-O code signature command.
type CodeSignatureCmd LinkEditDataCmd // LC_CODE_SIGNATURE
// A SegmentSplitInfoCmd is a Mach-O code info to split segments command.
type SegmentSplitInfoCmd LinkEditDataCmd // LC_SEGMENT_SPLIT_INFO
// A FunctionStartsCmd is a Mach-O compressed table of function start addresses command.
type FunctionStartsCmd LinkEditDataCmd // LC_FUNCTION_STARTS
// A DataInCodeCmd is a Mach-O data in code command.
type DataInCodeCmd LinkEditDataCmd // LC_DATA_IN_CODE
// A DylibCodeSignDrsCmd is a Mach-O code signing DRs copied from linked dylibs command.
type DylibCodeSignDrsCmd LinkEditDataCmd // LC_DYLIB_CODE_SIGN_DRS
// A LinkerOptimizationHintCmd is a Mach-O optimization hints command.
type LinkerOptimizationHintCmd LinkEditDataCmd // LC_LINKER_OPTIMIZATION_HINT
// A DyldExportsTrieCmd is used with linkedit_data_command, payload is trie command.
type DyldExportsTrieCmd LinkEditDataCmd // LC_DYLD_EXPORTS_TRIE
// A DyldChainedFixupsCmd is used with linkedit_data_command command.
type DyldChainedFixupsCmd LinkEditDataCmd // LC_DYLD_CHAINED_FIXUPS

type EncryptionSystem uint32

const NOT_ENCRYPTED_YET EncryptionSystem = 0

/*
 * EncryptionInfoCmd contains the file offset and size of an
 * of an encrypted segment.
 */
type EncryptionInfoCmd struct {
	LoadCmd                  // LC_ENCRYPTION_INFO
	Len     uint32           // sizeof(struct encryption_info_command)
	Offset  uint32           // file offset of encrypted range
	Size    uint32           // file size of encrypted range
	CryptID EncryptionSystem // which enryption system, 0 means not-encrypted yet
}

/*
 * EncryptionInfo64Cmd contains the file offset and size of an
 * of an encrypted segment (for use in x86_64 targets).
 */
type EncryptionInfo64Cmd struct {
	LoadCmd                  // LC_ENCRYPTION_INFO_64
	Len     uint32           // sizeof(struct encryption_info_command_64)
	Offset  uint32           // file offset of encrypted range
	Size    uint32           // file size of encrypted range
	CryptID EncryptionSystem // which enryption system, 0 means not-encrypted yet
	Pad     uint32           // padding to make this struct's size a multiple of 8 bytes
}

/*
 * VersionMinCmd contains the min OS version on which this
 * binary was built to run.
 */
type VersionMinCmd struct {
	LoadCmd /* LC_VERSION_MIN_MACOSX or
	   LC_VERSION_MIN_IPHONEOS or
	   LC_VERSION_MIN_WATCHOS or
	   LC_VERSION_MIN_TVOS */
	Len     uint32  // sizeof(struct min_version_command)
	Version Version // X.Y.Z is encoded in nibbles xxxx.yy.zz
	Sdk     Version // X.Y.Z is encoded in nibbles xxxx.yy.zz
}

// A VersionMinMacOSCmd is a Mach-O build for macOS min OS version command.
type VersionMinMacOSCmd VersionMinCmd // LC_VERSION_MIN_MACOSX
// A VersionMinIPhoneOSCmd is a Mach-O build for iPhoneOS min OS version command.
type VersionMinIPhoneOSCmd VersionMinCmd // LC_VERSION_MIN_IPHONEOS
// A VersionMinWatchOSCmd is a Mach-O build for watchOS min OS version command.
type VersionMinWatchOSCmd VersionMinCmd // LC_VERSION_MIN_WATCHOS
// A VersionMinTvOSCmd is a Mach-O build for tvOS min OS version command.
type VersionMinTvOSCmd VersionMinCmd // LC_VERSION_MIN_TVOS

/*
* BuildVersionCmd contains the min OS version on which this
* binary was built to run for its platform.  The list of known platforms and
* tool values following it.
 */
type BuildVersionCmd struct {
	LoadCmd        /* LC_BUILD_VERSION */
	Len     uint32 /* sizeof(struct build_version_command) plus */
	/* ntools * sizeof(struct build_tool_version) */
	Platform Platform /* platform */
	Minos    Version  /* X.Y.Z is encoded in nibbles xxxx.yy.zz */
	Sdk      Version  /* X.Y.Z is encoded in nibbles xxxx.yy.zz */
	NumTools uint32   /* number of tool entries following this */
}

/*
 * DyldInfoCmd contains the file offsets and sizes of
 * the new compressed form of the information dyld needs to
 * load the image.  This information is used by dyld on Mac OS X
 * 10.6 and later.  All information pointed to by this command
 * is encoded using byte streams, so no endian swapping is needed
 * to interpret it.
 */
type DyldInfoCmd struct {
	LoadCmd        //  LC_DYLD_INFO or LC_DYLD_INFO_ONLY
	Len     uint32 // sizeof(struct dyld_info_command)
	/*
	 * Dyld rebases an image whenever dyld loads it at an address different
	 * from its preferred address.  The rebase information is a stream
	 * of byte sized opcodes whose symbolic names start with REBASE_OPCODE_.
	 * Conceptually the rebase information is a table of tuples:
	 *    <seg-index, seg-offset, type>
	 * The opcodes are a compressed way to encode the table by only
	 * encoding when a column changes.  In addition simple patterns
	 * like "every n'th offset for m times" can be encoded in a few
	 * bytes.
	 */
	RebaseOff  uint32 // file offset to rebase info
	RebaseSize uint32 //  size of rebase info
	/*
	 * Dyld binds an image during the loading process, if the image
	 * requires any pointers to be initialized to symbols in other images.
	 * The bind information is a stream of byte sized
	 * opcodes whose symbolic names start with BIND_OPCODE_.
	 * Conceptually the bind information is a table of tuples:
	 *    <seg-index, seg-offset, type, symbol-library-ordinal, symbol-name, addend>
	 * The opcodes are a compressed way to encode the table by only
	 * encoding when a column changes.  In addition simple patterns
	 * like for runs of pointers initialzed to the same value can be
	 * encoded in a few bytes.
	 */
	BindOff  uint32 // file offset to binding info
	BindSize uint32 // size of binding info
	/*
	 * Some C++ programs require dyld to unique symbols so that all
	 * images in the process use the same copy of some code/data.
	 * This step is done after binding. The content of the weak_bind
	 * info is an opcode stream like the bind_info.  But it is sorted
	 * alphabetically by symbol name.  This enable dyld to walk
	 * all images with weak binding information in order and look
	 * for collisions.  If there are no collisions, dyld does
	 * no updating.  That means that some fixups are also encoded
	 * in the bind_info.  For instance, all calls to "operator new"
	 * are first bound to libstdc++.dylib using the information
	 * in bind_info.  Then if some image overrides operator new
	 * that is detected when the weak_bind information is processed
	 * and the call to operator new is then rebound.
	 */
	WeakBindOff  uint32 // file offset to weak binding info
	WeakBindSize uint32 //  size of weak binding info
	/*
	 * Some uses of external symbols do not need to be bound immediately.
	 * Instead they can be lazily bound on first use.  The lazy_bind
	 * are contains a stream of BIND opcodes to bind all lazy symbols.
	 * Normal use is that dyld ignores the lazy_bind section when
	 * loading an image.  Instead the static linker arranged for the
	 * lazy pointer to initially point to a helper function which
	 * pushes the offset into the lazy_bind area for the symbol
	 * needing to be bound, then jumps to dyld which simply adds
	 * the offset to lazy_bind_off to get the information on what
	 * to bind.
	 */
	LazyBindOff  uint32 // file offset to lazy binding info
	LazyBindSize uint32 //  size of lazy binding info
	/*
	 * The symbols exported by a dylib are encoded in a trie.  This
	 * is a compact representation that factors out common prefixes.
	 * It also reduces LINKEDIT pages in RAM because it encodes all
	 * information (name, address, flags) in one small, contiguous range.
	 * The export area is a stream of nodes.  The first node sequentially
	 * is the start node for the trie.
	 *
	 * Nodes for a symbol start with a uleb128 that is the length of
	 * the exported symbol information for the string so far.
	 * If there is no exported symbol, the node starts with a zero byte.
	 * If there is exported info, it follows the length.
	 *
	 * First is a uleb128 containing flags. Normally, it is followed by
	 * a uleb128 encoded offset which is location of the content named
	 * by the symbol from the mach_header for the image.  If the flags
	 * is EXPORT_SYMBOL_FLAGS_REEXPORT, then following the flags is
	 * a uleb128 encoded library ordinal, then a zero terminated
	 * UTF8 string.  If the string is zero length, then the symbol
	 * is re-export from the specified dylib with the same name.
	 * If the flags is EXPORT_SYMBOL_FLAGS_STUB_AND_RESOLVER, then following
	 * the flags is two uleb128s: the stub offset and the resolver offset.
	 * The stub is used by non-lazy pointers.  The resolver is used
	 * by lazy pointers and must be called to get the actual address to use.
	 *
	 * After the optional exported symbol information is a byte of
	 * how many edges (0-255) that this node has leaving it,
	 * followed by each edge.
	 * Each edge is a zero terminated UTF8 of the addition chars
	 * in the symbol, followed by a uleb128 offset for the node that
	 * edge points to.
	 *
	 */
	ExportOff  uint32 // file offset to export info
	ExportSize uint32 //  size of export info
}

// A DyldInfoOnlyCmd is a Mach-O compressed dyld information only command.
type DyldInfoOnlyCmd DyldInfoCmd // LC_DYLD_INFO_ONLY

/*
 * LinkerOptionCmd contains linker options embedded in object files.
 */
type LinkerOptionCmd struct {
	LoadCmd        // LC_LINKER_OPTION only used in MH_OBJECT filetypes
	Len     uint32 // sizeof(struct linker_option_command)
	Count   uint32 // number of strings concatenation of zero terminated UTF8 strings. Zero filled at end to align
}

/*
 * SymsegCmd contains the offset and size of the GNU style
 * symbol table information as described in the header file <symseg.h>.
 * The symbol roots of the symbol segments must also be aligned properly
 * in the file.  So the requirement of keeping the offsets aligned to a
 * multiple of a 4 bytes translates to the length field of the symbol
 * roots also being a multiple of a long.  Also the padding must again be
 * zeroed. (THIS IS OBSOLETE and no longer supported).
 */
type SymsegCmd struct {
	LoadCmd        /* LC_SYMSEG */
	Len     uint32 /* sizeof(struct symseg_command) */
	Offset  uint32 /* symbol segment offset */
	Size    uint32 /* symbol segment size in bytes */
}

/*
 * IdentCmd contains a free format string table following the
 * ident_command structure.  The strings are null terminated and the size of
 * the command is padded out with zero bytes to a multiple of 4 bytes/
 * (THIS IS OBSOLETE and no longer supported).
 */
type IdentCmd struct {
	LoadCmd        // LC_IDENT
	Len     uint32 // strings that follow this command
}

/*
 * FvmFileCmdcontains a reference to a file to be loaded at the
 * specified virtual address.  (Presently, this command is reserved for
 * internal use.  The kernel ignores this command when loading a program into
 * memory).
 */
type FvmFileCmd struct {
	LoadCmd           // LC_FVMFILE
	Len        uint32 // includes pathname string
	NameOffset uint32 // files pathname
	HeaderAddr uint32 // files virtual address
}

// A PrePageCmd is a fixed VM file inclusion (internal use).
type PrePageCmd struct {
	LoadCmd // LC_PREPAGE
	Len     uint32
}

/*
 * EntryPointCmd is a replacement for thread_command.
 * It is used for main executables to specify the location (file offset)
 * of main().  If -stack_size was used at link time, the stacksize
 * field will contain the stack size need for the main thread.
 */
type EntryPointCmd struct {
	LoadCmd            // LC_MAIN only used in MH_EXECUTE filetypes
	Len         uint32 // 24
	EntryOffset uint64 // file (__TEXT) offset of main()
	StackSize   uint64 // if not zero, initial stack size
}

/*
 * SourceVersionCmd is an optional load command containing
 * the version of the sources used to build the binary.
 */
type SourceVersionCmd struct {
	LoadCmd            // LC_SOURCE_VERSION
	Len     uint32     // 16
	Version SrcVersion // A.B.C.D.E packed as a24.b10.c10.d10.e10
}

/*
 * NoteCmd describe a region of arbitrary data included in a Mach-O
 * file.  Its initial use is to record extra data in MH_CORE files.
 */
type NoteCmd struct {
	LoadCmd            // LC_NOTE
	Len       uint32   // sizeof(struct note_command)
	DataOwner [16]byte // owner name for this LC_NOTE
	Offset    uint64   // file offset of this data
	Size      uint64   // length of data region
}

// FilesetEntryCmd commands describe constituent Mach-O files that are part
// of a fileset. In one implementation, entries are dylibs with individual
// mach headers and repositionable text and data segments. Each entry is
// further described by its own mach header.
type FilesetEntryCmd struct {
	LoadCmd       // LC_FILESET_ENTRY
	Len           uint32
	Addr          uint64 // memory address of the entry
	FileOffset    uint64 // file offset of the entry
	EntryIdOffset uint32 // contained entry id
	Reserved      uint32 // reserved
}
