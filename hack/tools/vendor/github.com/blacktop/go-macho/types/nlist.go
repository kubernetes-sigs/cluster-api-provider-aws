package types

import (
	"encoding/binary"
	"strings"
)

// An Nlist is a Mach-O generic symbol table entry.
type Nlist struct {
	Name uint32
	Type NType
	Sect uint8
	Desc NDescType
}

// An Nlist32 is a Mach-O 32-bit symbol table entry.
type Nlist32 struct {
	Nlist
	Value uint32
}

func (n *Nlist32) Put32(b []byte, o binary.ByteOrder) uint32 {
	o.PutUint32(b[0:], n.Name)
	b[4] = byte(n.Type)
	b[5] = byte(n.Sect)
	o.PutUint16(b[6:], uint16(n.Desc))
	o.PutUint32(b[8:], uint32(n.Value))
	return 8 + 4
}

// An Nlist64 is a Mach-O 64-bit symbol table entry.
type Nlist64 struct {
	Nlist
	Value uint64
}

func (n *Nlist64) Put64(b []byte, o binary.ByteOrder) uint32 {
	o.PutUint32(b[0:], n.Name)
	b[4] = byte(n.Type)
	b[5] = byte(n.Sect)
	o.PutUint16(b[6:], uint16(n.Desc))
	o.PutUint64(b[8:], n.Value)
	return 8 + 8
}

type NType uint8

/*
 * The n_type field really contains four fields:
 *	unsigned char N_STAB:3,
 *		      N_PEXT:1,
 *		      N_TYPE:3,
 *		      N_EXT:1;
 * which are used via the following masks.
 */
const (
	N_STAB NType = 0xe0 /* if any of these bits set, a symbolic debugging entry */
	N_PEXT NType = 0x10 /* private external symbol bit */
	N_TYPE NType = 0x0e /* mask for the type bits */
	N_EXT  NType = 0x01 /* external symbol bit, set for external symbols */
)

/*
 * Values for N_TYPE bits of the n_type field.
 */
const (
	N_UNDF NType = 0x0 /* undefined, n_sect == NO_SECT */
	N_ABS  NType = 0x2 /* absolute, n_sect == NO_SECT */
	N_SECT NType = 0xe /* defined in section number n_sect */
	N_PBUD NType = 0xc /* prebound undefined (defined in a dylib) */
	N_INDR NType = 0xa /* indirect */
)

const (
	/*
	 * If the type is N_INDR then the symbol is defined to be the same as another
	 * symbol.  In this case the n_value field is an index into the string table
	 * of the other symbol's name.  When the other symbol is defined then they both
	 * take on the defined type and value.
	 */

	/*
	 * If the type is N_SECT then the n_sect field contains an ordinal of the
	 * section the symbol is defined in.  The sections are numbered from 1 and
	 * refer to sections in order they appear in the load commands for the file
	 * they are in.  This means the same ordinal may very well refer to different
	 * sections in different files.
	 *
	 * The n_value field for all symbol table entries (including N_STAB's) gets
	 * updated by the link editor based on the value of it's n_sect field and where
	 * the section n_sect references gets relocated.  If the value of the n_sect
	 * field is NO_SECT then it's n_value field is not changed by the link editor.
	 */
	NO_SECT  = 0   /* symbol is not in any section */
	MAX_SECT = 255 /* 1 thru 255 inclusive */
)

func (t NType) IsDebugSym() bool {
	return (t & N_STAB) != 0
}

func (t NType) IsPrivateExternalSym() bool {
	return (t & N_PEXT) != 0
}

func (t NType) IsExternalSym() bool {
	return (t & N_EXT) != 0
}

func (t NType) IsUndefinedSym() bool {
	return (t & N_TYPE) == N_UNDF
}
func (t NType) IsAbsoluteSym() bool {
	return (t & N_TYPE) == N_ABS
}
func (t NType) IsDefinedInSection() bool {
	return (t & N_TYPE) == N_SECT
}
func (t NType) IsPreboundUndefinedSym() bool {
	return (t & N_TYPE) == N_PBUD
}
func (t NType) IsIndirectSym() bool {
	return (t & N_TYPE) == N_INDR
}

func (t NType) String(secName string) string {
	var out []string
	if t.IsDebugSym() {
		switch {
		case t.IsGlobal():
			out = append(out, "debug(global)")
		case t.IsProcedureName():
			out = append(out, "debug(procedure name)")
		case t.IsProcedure():
			out = append(out, "debug(procedure)")
		case t.IsStatic():
			out = append(out, "debug(static)")
		case t.IsLcommSym():
			out = append(out, "debug(.lcomm)")
		case t.IsBeginNsectSym():
			out = append(out, "debug(begin nsect sym)")
		case t.IsAstFilePath():
			out = append(out, "debug(ast file path)")
		case t.IsGccCompiled():
			out = append(out, "debug(gcc compiled)")
		case t.IsRegisterSym():
			out = append(out, "debug(register)")
		case t.IsSourceLine():
			out = append(out, "debug(source line)")
		case t.IsEndNsectSym():
			out = append(out, "debug(end nsect sym)")
		case t.IsStructure():
			out = append(out, "debug(struct_offset)")
		case t.IsSourceFile():
			out = append(out, "debug(source file)")
		case t.IsObjectFile():
			out = append(out, "debug(object file)")
		case t.IsLib():
			out = append(out, "debug(dylib)")
		case t.IsLocalSym():
			out = append(out, "debug(local)")
		case t.IsIncludeFileBegin():
			out = append(out, "debug(include file beginning)")
		case t.IsIncludedFile():
			out = append(out, "debug(#included file name)")
		case t.IsCompilerParams():
			out = append(out, "debug(compiler parameters)")
		case t.IsCompilerVersion():
			out = append(out, "debug(compiler version)")
		case t.IsCompilerOLevel():
			out = append(out, "debug(compiler -O level)")
		case t.IsParameter():
			out = append(out, "debug(parameter)")
		case t.IsIncludeFileEnd():
			out = append(out, "debug(include file end)")
		case t.IsAlternateEntry():
			out = append(out, "debug(alternate entry)")
		case t.IsLeftBracket():
			out = append(out, "debug(left bracket)")
		case t.IsDeletedIncludeFile():
			out = append(out, "debug(deleted include file)")
		case t.IsRightBracket():
			out = append(out, "debug(right bracket)")
		case t.IsBeginCommon():
			out = append(out, "debug(begin common)")
		case t.IsEndCommon():
			out = append(out, "debug(end common)")
		case t.IsEndCommonLocal():
			out = append(out, "debug(end common - local)")
		case t.IsSecondStabEntry():
			out = append(out, "debug(second stab entry)")
		case t.IsPascalSymbol():
			out = append(out, "debug(global pascal symbol)")
		default:
			out = append(out, "debug")
		}
	}
	if t.IsPrivateExternalSym() {
		out = append(out, "private")
	}
	if t.IsExternalSym() {
		out = append(out, "external")
	}
	if t.IsUndefinedSym() {
		out = append(out, "undefined")
	}
	if t.IsAbsoluteSym() {
		out = append(out, "absolute")
	}
	if t.IsDefinedInSection() {
		out = append(out, secName)
	}
	if t.IsPreboundUndefinedSym() {
		out = append(out, "prebound")
	}
	if t.IsIndirectSym() {
		out = append(out, "indirect")
	}
	return strings.Join(out, "|")
}

type NDescType uint16

/*
 * Common symbols are represented by undefined (N_UNDF) external (N_EXT) types
 * who's values (n_value) are non-zero.  In which case the value of the n_value
 * field is the size (in bytes) of the common symbol.  The n_sect field is set
 * to NO_SECT.  The alignment of a common symbol may be set as a power of 2
 * between 2^1 and 2^15 as part of the n_desc field using the macros below. If
 * the alignment is not set (a value of zero) then natural alignment based on
 * the size is used.
 */
func (d NDescType) GetCommAlign() NDescType { // TODO: apply this to common symbol's value
	return (d >> 8) & 0x0f
}

const REFERENCE_TYPE_MASK NDescType = 0x7

const (
	/* types of references */
	REFERENCE_FLAG_UNDEFINED_NON_LAZY         NDescType = 0
	REFERENCE_FLAG_UNDEFINED_LAZY             NDescType = 1
	REFERENCE_FLAG_DEFINED                    NDescType = 2
	REFERENCE_FLAG_PRIVATE_DEFINED            NDescType = 3
	REFERENCE_FLAG_PRIVATE_UNDEFINED_NON_LAZY NDescType = 4
	REFERENCE_FLAG_PRIVATE_UNDEFINED_LAZY     NDescType = 5
)

func (d NDescType) IsUndefinedNonLazy() bool {
	return (d & REFERENCE_TYPE_MASK) == REFERENCE_FLAG_UNDEFINED_NON_LAZY
}
func (d NDescType) IsUndefinedLazy() bool {
	return (d & REFERENCE_TYPE_MASK) == REFERENCE_FLAG_UNDEFINED_LAZY
}
func (d NDescType) IsDefined() bool {
	return (d & REFERENCE_TYPE_MASK) == REFERENCE_FLAG_DEFINED
}
func (d NDescType) IsPrivateDefined() bool {
	return (d & REFERENCE_TYPE_MASK) == REFERENCE_FLAG_PRIVATE_DEFINED
}
func (d NDescType) IsPrivateUndefinedNonLazy() bool {
	return (d & REFERENCE_TYPE_MASK) == REFERENCE_FLAG_PRIVATE_UNDEFINED_NON_LAZY
}
func (d NDescType) IsPrivateUndefinedLazy() bool {
	return (d & REFERENCE_TYPE_MASK) == REFERENCE_FLAG_PRIVATE_UNDEFINED_LAZY
}

const (
	SELF_LIBRARY_ORDINAL   = 0x0
	MAX_LIBRARY_ORDINAL    = 0xfd
	DYNAMIC_LOOKUP_ORDINAL = 0xfe
	EXECUTABLE_ORDINAL     = 0xff
)

func (d NDescType) GetLibraryOrdinal() uint16 {
	return (uint16(d) >> 8) & 0xff
}

const (
	/*
	 * To simplify stripping of objects that use are used with the dynamic link
	 * editor, the static link editor marks the symbols defined an object that are
	 * referenced by a dynamicly bound object (dynamic shared libraries, bundles).
	 * With this marking strip knows not to strip these symbols.
	 */
	REFERENCED_DYNAMICALLY NDescType = 0x0010
	/*
	 * The N_NO_DEAD_STRIP bit of the n_desc field only ever appears in a
	 * relocatable .o file (MH_OBJECT filetype). And is used to indicate to the
	 * static link editor it is never to dead strip the symbol.
	 */
	NO_DEAD_STRIP NDescType = 0x0020 /* symbol is not to be dead stripped */

	/*
	 * The N_DESC_DISCARDED bit of the n_desc field never appears in linked image.
	 * But is used in very rare cases by the dynamic link editor to mark an in
	 * memory symbol as discared and longer used for linking.
	 */
	DESC_DISCARDED NDescType = 0x0020 /* symbol is discarded */

	/*
	 * The N_WEAK_REF bit of the n_desc field indicates to the dynamic linker that
	 * the undefined symbol is allowed to be missing and is to have the address of
	 * zero when missing.
	 */
	WEAK_REF NDescType = 0x0040 /* symbol is weak referenced */

	/*
	 * The N_WEAK_DEF bit of the n_desc field indicates to the static and dynamic
	 * linkers that the symbol definition is weak, allowing a non-weak symbol to
	 * also be used which causes the weak definition to be discared.  Currently this
	 * is only supported for symbols in coalesed sections.
	 */
	WEAK_DEF NDescType = 0x0080 /* coalesed symbol is a weak definition */

	/*
	 * The N_REF_TO_WEAK bit of the n_desc field indicates to the dynamic linker
	 * that the undefined symbol should be resolved using flat namespace searching.
	 */
	REF_TO_WEAK NDescType = 0x0080 /* reference to a weak symbol */

	/*
	 * The N_ARM_THUMB_DEF bit of the n_desc field indicates that the symbol is
	 * a defintion of a Thumb function.
	 */
	ARM_THUMB_DEF NDescType = 0x0008 /* symbol is a Thumb function (ARM) */

	/*
	 * The N_SYMBOL_RESOLVER bit of the n_desc field indicates that the
	 * that the function is actually a resolver function and should
	 * be called to get the address of the real function to use.
	 * This bit is only available in .o files (MH_OBJECT filetype)
	 */
	SYMBOL_RESOLVER NDescType = 0x0100

	/*
	 * The N_ALT_ENTRY bit of the n_desc field indicates that the
	 * symbol is pinned to the previous content.
	 */
	ALT_ENTRY NDescType = 0x0200

	/*
	 * The N_COLD_FUNC bit of the n_desc field indicates that the symbol is used
	 * infrequently and the linker should order it towards the end of the section.
	 */
	N_COLD_FUNC NDescType = 0x0400
)

func (d NDescType) IsReferencedDynamically() bool {
	return (d & REFERENCED_DYNAMICALLY) != 0
}
func (d NDescType) IsNoDeadStrip() bool {
	return (d & NO_DEAD_STRIP) != 0
}
func (d NDescType) IsDescDiscarded() bool {
	return (d & DESC_DISCARDED) != 0
}
func (d NDescType) IsWeakReferenced() bool {
	return (d & WEAK_REF) != 0
}
func (d NDescType) IsWeakDefintion() bool {
	return (d & WEAK_DEF) != 0
}
func (d NDescType) IsWeakDefintionOrReferenced() bool {
	return (d & (WEAK_DEF | WEAK_REF)) != 0
}
func (d NDescType) IsReferenceToWeak() bool {
	return (d & REF_TO_WEAK) != 0
}
func (d NDescType) IsArmThumbDefintion() bool {
	return (d & ARM_THUMB_DEF) != 0
}
func (d NDescType) IsSymbolResolver() bool {
	return (d & SYMBOL_RESOLVER) != 0
}
func (d NDescType) IsAltEntry() bool {
	return (d & ALT_ENTRY) != 0
}
func (d NDescType) IsColdFunc() bool {
	return (d & N_COLD_FUNC) != 0
}

func (t NDescType) String() string {
	var out []string
	if t.IsUndefinedNonLazy() {
		out = append(out, "undefined")
	}
	if t.IsUndefinedLazy() {
		out = append(out, "undefined_lazy")
	}
	if t.IsDefined() {
		out = append(out, "def")
	}
	if t.IsPrivateDefined() {
		out = append(out, "priv_def")
	}
	if t.IsPrivateUndefinedNonLazy() {
		out = append(out, "priv_undef_nonlazy")
	}
	if t.IsPrivateUndefinedLazy() {
		out = append(out, "priv_undef_lazy")
	}
	// if t.GetLibraryOrdinal() != SELF_LIBRARY_ORDINAL {
	// 	out = append(out, fmt.Sprintf("libord=%s", libName))
	// }
	if t.IsReferencedDynamically() {
		out = append(out, "referenced_dynamically")
	}
	if t.IsNoDeadStrip() {
		out = append(out, "no_dead_strip")
	}
	if t.IsDescDiscarded() {
		out = append(out, "discarded")
	}
	if t.IsWeakReferenced() {
		out = append(out, "weak_ref")
	}
	if t.IsWeakDefintion() {
		out = append(out, "weak_def")
	}
	if t.IsReferenceToWeak() {
		out = append(out, "ref_to_weak")
	}
	if t.IsArmThumbDefintion() {
		out = append(out, "arm_thumb_def")
	}
	if t.IsSymbolResolver() {
		out = append(out, "symbol_resolver")
	}
	if t.IsAltEntry() {
		out = append(out, "alt_entry")
	}
	if t.IsColdFunc() {
		out = append(out, "cold_func")
	}
	return strings.Join(out, "|")
}

/*
 * Symbolic debugger symbols.
 */
const (
	N_GSYM  NType = 0x20 /* global symbol: name,,NO_SECT,type,0 */
	N_FNAME NType = 0x22 /* procedure name (f77 kludge): name,,NO_SECT,0,0 */
	N_FUN   NType = 0x24 /* procedure: name,,n_sect,linenumber,address */
	N_STSYM NType = 0x26 /* static symbol: name,,n_sect,type,address */
	N_LCSYM NType = 0x28 /* .lcomm symbol: name,,n_sect,type,address */
	N_BNSYM NType = 0x2e /* begin nsect sym: 0,,n_sect,0,address */
	N_AST   NType = 0x32 /* AST file path: name,,NO_SECT,0,0 */
	N_OPT   NType = 0x3c /* emitted with gcc2_compiled and in gcc source */
	N_RSYM  NType = 0x40 /* register sym: name,,NO_SECT,type,register */
	N_SLINE NType = 0x44 /* src line: 0,,n_sect,linenumber,address */
	N_ENSYM NType = 0x4e /* end nsect sym: 0,,n_sect,0,address */
	N_SSYM  NType = 0x60 /* structure elt: name,,NO_SECT,type,struct_offset */
	N_SO    NType = 0x64 /* source file name: name,,n_sect,0,address */
	N_OSO   NType = 0x66 /* object file name: name,,(see below),1,st_mtime */
	/*   historically N_OSO set n_sect to 0. The N_OSO
	 *   n_sect may instead hold the low byte of the
	 *   cpusubtype value from the Mach-O header. */
	N_LIB     NType = 0x68 /* dynamic library file name: name,,NO_SECT,0,0 */
	N_LSYM    NType = 0x80 /* local sym: name,,NO_SECT,type,offset */
	N_BINCL   NType = 0x82 /* include file beginning: name,,NO_SECT,0,sum */
	N_SOL     NType = 0x84 /* #included file name: name,,n_sect,0,address */
	N_PARAMS  NType = 0x86 /* compiler parameters: name,,NO_SECT,0,0 */
	N_VERSION NType = 0x88 /* compiler version: name,,NO_SECT,0,0 */
	N_OLEVEL  NType = 0x8A /* compiler -O level: name,,NO_SECT,0,0 */
	N_PSYM    NType = 0xa0 /* parameter: name,,NO_SECT,type,offset */
	N_EINCL   NType = 0xa2 /* include file end: name,,NO_SECT,0,0 */
	N_ENTRY   NType = 0xa4 /* alternate entry: name,,n_sect,linenumber,address */
	N_LBRAC   NType = 0xc0 /* left bracket: 0,,NO_SECT,nesting level,address */
	N_EXCL    NType = 0xc2 /* deleted include file: name,,NO_SECT,0,sum */
	N_RBRAC   NType = 0xe0 /* right bracket: 0,,NO_SECT,nesting level,address */
	N_BCOMM   NType = 0xe2 /* begin common: name,,NO_SECT,0,0 */
	N_ECOMM   NType = 0xe4 /* end common: name,,n_sect,0,0 */
	N_ECOML   NType = 0xe8 /* end common (local name): 0,,n_sect,0,address */
	N_LENG    NType = 0xfe /* second stab entry with length information */
	/*
	 * for the berkeley pascal compiler, pc(1):
	 */
	N_PC NType = 0x30 /* global pascal symbol: name,,NO_SECT,subtype,line */
)

func (t NType) IsGlobal() bool {
	return t == N_GSYM
}
func (t NType) IsProcedureName() bool {
	return t == N_FNAME
}
func (t NType) IsProcedure() bool {
	return t == N_FUN
}
func (t NType) IsStatic() bool {
	return t == N_STSYM
}
func (t NType) IsLcommSym() bool {
	return t == N_LCSYM
}
func (t NType) IsBeginNsectSym() bool {
	return t == N_BNSYM
}
func (t NType) IsAstFilePath() bool {
	return t == N_AST
}
func (t NType) IsGccCompiled() bool {
	return t == N_OPT
}
func (t NType) IsRegisterSym() bool {
	return t == N_RSYM
}
func (t NType) IsSourceLine() bool {
	return t == N_SLINE
}
func (t NType) IsEndNsectSym() bool {
	return t == N_ENSYM
}
func (t NType) IsStructure() bool {
	return t == N_SSYM
}
func (t NType) IsSourceFile() bool {
	return t == N_SO
}
func (t NType) IsObjectFile() bool {
	return t == N_OSO
}
func (t NType) IsLib() bool {
	return t == N_LIB
}
func (t NType) IsLocalSym() bool {
	return t == N_LSYM
}
func (t NType) IsIncludeFileBegin() bool {
	return t == N_BINCL
}
func (t NType) IsIncludedFile() bool {
	return t == N_SOL
}
func (t NType) IsCompilerParams() bool {
	return t == N_PARAMS
}
func (t NType) IsCompilerVersion() bool {
	return t == N_VERSION
}
func (t NType) IsCompilerOLevel() bool {
	return t == N_OLEVEL
}
func (t NType) IsParameter() bool {
	return t == N_PSYM
}
func (t NType) IsIncludeFileEnd() bool {
	return t == N_EINCL
}
func (t NType) IsAlternateEntry() bool {
	return t == N_ENTRY
}
func (t NType) IsLeftBracket() bool {
	return t == N_LBRAC
}
func (t NType) IsDeletedIncludeFile() bool {
	return t == N_EXCL
}
func (t NType) IsRightBracket() bool {
	return t == N_RBRAC
}
func (t NType) IsBeginCommon() bool {
	return t == N_BCOMM
}
func (t NType) IsEndCommon() bool {
	return t == N_ECOMM
}
func (t NType) IsEndCommonLocal() bool {
	return t == N_ECOML
}
func (t NType) IsSecondStabEntry() bool {
	return t == N_LENG
}
func (t NType) IsPascalSymbol() bool {
	return t == N_PC
}
