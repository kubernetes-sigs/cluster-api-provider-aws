package objc

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/blacktop/go-macho/types"
)

const IsDyldPreoptimized = 1 << 7

type Toc struct {
	ClassList        uint64
	NonLazyClassList uint64
	CatList          uint64
	NonLazyCatList   uint64
	ProtoList        uint64
	ClassRefs        uint64
	SuperRefs        uint64
	SelRefs          uint64
	Stubs            uint64
}

func (i Toc) String() string {
	return fmt.Sprintf(
		"ObjC TOC\n"+
			"--------\n"+
			"  __objc_classlist  = %d\n"+
			"  __objc_nlclslist  = %d\n"+
			"  __objc_catlist    = %d\n"+
			"  __objc_nlcatlist  = %d\n"+
			"  __objc_protolist  = %d\n"+
			"  __objc_classrefs  = %d\n"+
			"  __objc_superrefs  = %d\n"+
			"  __objc_selrefs    = %d\n",
		// "  __objc_stubs      = %d\n",
		i.ClassList,
		i.NonLazyClassList,
		i.CatList,
		i.NonLazyCatList,
		i.ProtoList,
		i.ClassRefs,
		i.SuperRefs,
		i.SelRefs,
		// i.Stubs,
	)
}

type ImageInfoFlag uint32

const (
	IsReplacement              ImageInfoFlag = 1 << 0 // used for Fix&Continue, now ignored
	SupportsGC                 ImageInfoFlag = 1 << 1 // image supports GC
	RequiresGC                 ImageInfoFlag = 1 << 2 // image requires GC
	OptimizedByDyld            ImageInfoFlag = 1 << 3 // image is from an optimized shared cache
	SignedClassRO              ImageInfoFlag = 1 << 4 // class_ro_t pointers are signed
	IsSimulated                ImageInfoFlag = 1 << 5 // image compiled for a simulator platform
	HasCategoryClassProperties ImageInfoFlag = 1 << 6 // class properties in category_t
	OptimizedByDyldClosure     ImageInfoFlag = 1 << 7 // dyld (not the shared cache) optimized this.
	// 1 byte Swift unstable ABI version number
	SwiftUnstableVersionMaskShift = 8
	SwiftUnstableVersionMask      = 0xff << SwiftUnstableVersionMaskShift

	// 2 byte Swift stable ABI version number
	SwiftStableVersionMaskShift = 16
	SwiftStableVersionMask      = 0xffff << SwiftStableVersionMaskShift
)

func (f ImageInfoFlag) IsReplacement() bool {
	return f&IsReplacement != 0
}
func (f ImageInfoFlag) SupportsGC() bool {
	return f&SupportsGC != 0
}
func (f ImageInfoFlag) RequiresGC() bool {
	return f&RequiresGC != 0
}
func (f ImageInfoFlag) OptimizedByDyld() bool {
	return f&OptimizedByDyld != 0
}
func (f ImageInfoFlag) SignedClassRO() bool {
	return f&SignedClassRO != 0
}
func (f ImageInfoFlag) IsSimulated() bool {
	return f&IsSimulated != 0
}
func (f ImageInfoFlag) HasCategoryClassProperties() bool {
	return f&HasCategoryClassProperties != 0
}
func (f ImageInfoFlag) OptimizedByDyldClosure() bool {
	return f&OptimizedByDyldClosure != 0
}

func (f ImageInfoFlag) List() []string {
	var flags []string
	if (f & IsReplacement) != 0 {
		flags = append(flags, "IsReplacement")
	}
	if (f & SupportsGC) != 0 {
		flags = append(flags, "SupportsGC")
	}
	if (f & RequiresGC) != 0 {
		flags = append(flags, "RequiresGC")
	}
	if (f & OptimizedByDyld) != 0 {
		flags = append(flags, "OptimizedByDyld")
	}
	if (f & SignedClassRO) != 0 {
		flags = append(flags, "SignedClassRO")
	}
	if (f & IsSimulated) != 0 {
		flags = append(flags, "IsSimulated")
	}
	if (f & HasCategoryClassProperties) != 0 {
		flags = append(flags, "HasCategoryClassProperties")
	}
	if (f & OptimizedByDyldClosure) != 0 {
		flags = append(flags, "OptimizedByDyldClosure")
	}
	return flags
}

func (f ImageInfoFlag) String() string {
	return fmt.Sprintf(
		"Flags = %s\n"+
			"Swift = %s\n",
		strings.Join(f.List(), ", "),
		f.SwiftVersion(),
	)
}

func (f ImageInfoFlag) SwiftVersion() string {
	// TODO: I noticed there is some flags higher than swift version
	// (Console has 84019008, which is a version of 0x502)
	swiftVersion := (f >> 8) & 0xff
	if swiftVersion != 0 {
		switch swiftVersion {
		case 1:
			return "Swift 1.0"
		case 2:
			return "Swift 1.2"
		case 3:
			return "Swift 2.0"
		case 4:
			return "Swift 3.0"
		case 5:
			return "Swift 4.0"
		case 6:
			return "Swift 4.1/4.2"
		case 7:
			return "Swift 5 or later"
		default:
			return fmt.Sprintf("Unknown future Swift version: %d", swiftVersion)
		}
	}
	return "not swift"
}

const dyldPreoptimized = 1 << 7

type ImageInfo struct {
	Version uint32
	Flags   ImageInfoFlag
}

func (i ImageInfo) IsDyldPreoptimized() bool {
	return (i.Flags & dyldPreoptimized) != 0
}

func (i ImageInfo) HasSwift() bool {
	return (i.Flags>>8)&0xff != 0
}

const (
	bigSignedMethodListFlag              uint64 = 0x8000000000000000
	relativeMethodSelectorsAreDirectFlag uint32 = 0x40000000
	smallMethodListFlag                  uint32 = 0x80000000
	METHOD_LIST_FLAGS_MASK               uint32 = 0xffff0003
	// The size is bits 2 through 16 of the entsize field
	// The low 2 bits are uniqued/sorted as above.  The upper 16-bits
	// are reserved for other flags
	METHOD_LIST_SIZE_MASK uint32 = 0x0000FFFC
)

type MLFlags uint32

const (
	METHOD_LIST_IS_UNIQUED MLFlags = 1
	METHOD_LIST_IS_SORTED  MLFlags = 2
	METHOD_LIST_FIXED_UP   MLFlags = 3
)

type MLKind uint32

const (
	kindMask = 3
	// Note: method_invoke detects small methods by detecting 1 in the low
	// bit. Any change to that will require a corresponding change to
	// method_invoke.
	big MLKind = 0
	// `small` encompasses both small and small direct methods. We
	// distinguish those cases by doing a range check against the shared
	// cache.
	small       MLKind = 1
	bigSigned   MLKind = 2
	bigStripped MLKind = 3 // ***HACK: This is a TEMPORARY HACK FOR EXCLAVEKIT. It MUST go away.
)

type methodPtr uint64

func (m methodPtr) Kind() MLKind {
	return MLKind(m & kindMask)
}
func (m methodPtr) Pointer() uint64 {
	return uint64(m & ^methodPtr(kindMask))
}

type EntryList struct {
	Entsize uint32
	Count   uint32
}

func (el EntryList) String() string {
	return fmt.Sprintf("ent_size: %d, count: %d", el.Entsize, el.Count)
}

type Entry int64

func (e Entry) ImageIndex() uint16 {
	return uint16(e & 0xFFFF)
}
func (e Entry) MethodListOffset() int64 {
	return int64(e >> 16)
}
func (e Entry) String() string {
	return fmt.Sprintf("image_index: %d, method_list_offset: %d", e.ImageIndex(), e.MethodListOffset())
}

type MethodList struct {
	EntSizeAndFlags uint32
	Count           uint32
	// Space           uint32
	// MethodArrayBase uint64
}

func (ml MethodList) IsUniqued() bool {
	return (ml.Flags() & METHOD_LIST_IS_UNIQUED) == 1
}
func (ml MethodList) Sorted() bool {
	return (ml.Flags() & METHOD_LIST_IS_SORTED) == 1
}
func (ml MethodList) FixedUp() bool {
	return (ml.Flags() & METHOD_LIST_FIXED_UP) == 1
}
func (ml MethodList) UsesDirectOffsetsToSelectors() bool {
	return (ml.EntSizeAndFlags & relativeMethodSelectorsAreDirectFlag) != 0
}
func (ml MethodList) UsesRelativeOffsets() bool {
	return (ml.EntSizeAndFlags & smallMethodListFlag) != 0
}
func (ml MethodList) EntSize() uint32 {
	return ml.EntSizeAndFlags & METHOD_LIST_SIZE_MASK
}
func (ml MethodList) Flags() MLFlags {
	return MLFlags(ml.EntSizeAndFlags & METHOD_LIST_FLAGS_MASK)
}
func (ml MethodList) String() string {
	offType := "direct"
	if ml.UsesRelativeOffsets() {
		offType = "relative"
	}
	return fmt.Sprintf("count=%d, entsiz_flags=%#x, entrysize=%d, flags=%#x, fixed_up=%t, sorted=%t, uniqued=%t, type=%s",
		ml.Count,
		ml.EntSizeAndFlags,
		ml.EntSize(),
		ml.Flags(),
		ml.FixedUp(),
		ml.Sorted(),
		ml.IsUniqued(),
		offType)
}

type MethodT struct {
	NameVMAddr  uint64 // SEL
	TypesVMAddr uint64 // const char *
	ImpVMAddr   uint64 // IMP
}

type RelativeMethodT struct {
	NameOffset  int32 // SEL
	TypesOffset int32 // const char *
	ImpOffset   int32 // IMP
}

type Method struct {
	NameVMAddr  uint64 // & SEL
	TypesVMAddr uint64 // & const char *
	ImpVMAddr   uint64 // & IMP

	// We also need to know where the reference to the nameVMAddr was
	// This is so that we know how to rebind that location
	NameLocationVMAddr uint64
	Name               string
	Types              string
}

// NumberOfArguments returns the number of method arguments
func (m *Method) NumberOfArguments() int {
	if m == nil {
		return 0
	}
	return getNumberOfArguments(m.Types)
}

// ReturnType returns the method's return type
func (m *Method) ReturnType() string {
	return getReturnType(m.Types)
}

func (m *Method) ArgumentType(index int) string {
	args := getArguments(m.Types)
	if 0 < len(args) && index <= len(args) {
		return args[index].DecType
	}
	return "<error>"
}

type PropertyList struct {
	EntSize uint32
	Count   uint32
}

type PropertyT struct {
	NameVMAddr       uint64
	AttributesVMAddr uint64
}

type Property struct {
	PropertyT
	Name       string
	Attributes string
}

type CategoryT struct {
	NameVMAddr               uint64
	ClsVMAddr                uint64
	InstanceMethodsVMAddr    uint64
	ClassMethodsVMAddr       uint64
	ProtocolsVMAddr          uint64
	InstancePropertiesVMAddr uint64
}

type Category struct {
	Name            string
	VMAddr          uint64
	Class           *Class
	Protocols       []Protocol
	ClassMethods    []Method
	InstanceMethods []Method
	Properties      []Property
	CategoryT
}

func (c *Category) dump(verbose bool) string {
	var cMethods string
	var iMethods string
	var isSwift string

	var protos string
	if len(c.Protocols) > 0 {
		var prots []string
		for _, prot := range c.Protocols {
			prots = append(prots, prot.Name)
		}
		protos += fmt.Sprintf(" <%s>", strings.Join(prots, ", "))
	}

	var className string
	if c.Class != nil {
		className = c.Class.Name + " "
		if c.Class.IsSwift() {
			isSwift = " (Swift)"
		}
	}
	cat := fmt.Sprintf("@interface %s(%s)%s // %#x%s", className, c.Name, protos, c.VMAddr, isSwift)

	if len(c.ClassMethods) > 0 {
		s := bytes.NewBufferString("/* class methods */\n")
		w := tabwriter.NewWriter(s, 0, 0, 1, ' ', 0)
		for _, meth := range c.ClassMethods {
			if verbose {
				rtype, args := decodeMethodTypes(meth.Types)
				// fmt.Fprintf(w, "+ %s;\t// %#x\n", getMethodWithArgs(meth.Name, rtype, args), meth.ImpVMAddr)
				s.WriteString(fmt.Sprintf("// %#x\n", meth.ImpVMAddr))
				s.WriteString(fmt.Sprintf("+ %s\n", getMethodWithArgs(meth.Name, rtype, args)))
				// s.WriteString(fmt.Sprintf("+ %-80s // %#x\n", getMethodWithArgs(meth.Name, rtype, args), meth.ImpVMAddr))
			} else {
				fmt.Fprintf(w, "+[%s %s];\t// %#x\n", c.Name, meth.Name, meth.ImpVMAddr)
				// s.WriteString(fmt.Sprintf("+[%s %s]; %-40s\n", c.Name, meth.Name, fmt.Sprintf("// %#x", meth.ImpVMAddr)))
			}
		}
		w.Flush()
		cMethods = s.String()
	}
	if len(c.InstanceMethods) > 0 {
		var s *bytes.Buffer
		if len(c.ClassMethods) > 0 {
			s = bytes.NewBufferString("\n/* instance methods */\n\n")
		} else {
			s = bytes.NewBufferString("/* instance methods */\n\n")
		}
		w := tabwriter.NewWriter(s, 0, 0, 1, ' ', 0)
		for _, meth := range c.InstanceMethods {
			if verbose {
				rtype, args := decodeMethodTypes(meth.Types)
				// fmt.Fprintf(w, "- %s;\t// %#x\n", getMethodWithArgs(meth.Name, rtype, args), meth.ImpVMAddr)
				s.WriteString(fmt.Sprintf("// %#x\n", meth.ImpVMAddr))
				s.WriteString(fmt.Sprintf("- %s\n", getMethodWithArgs(meth.Name, rtype, args)))
				// s.WriteString(fmt.Sprintf("- %-80s // %#x\n", getMethodWithArgs(meth.Name, rtype, args), meth.ImpVMAddr))
			} else {
				fmt.Fprintf(w, "-[%s %s];\t// %#x\n", c.Name, meth.Name, meth.ImpVMAddr)
				// s.WriteString(fmt.Sprintf("-[%s %s]; %-40sn", c.Name, meth.Name, fmt.Sprintf("// %#x", meth.ImpVMAddr)))
			}
		}
		w.Flush()
		iMethods = s.String()
	}

	return fmt.Sprintf(
		"%s\n%s%s@end\n",
		cat,
		cMethods,
		iMethods)
}

func (c *Category) String() string {
	return c.dump(false)
}

func (c *Category) Verbose() string {
	return c.dump(true)
}

const (
	// Values for protocol_t->flags
	PROTOCOL_FIXED_UP_2   = (1 << 31) // must never be set by compiler
	PROTOCOL_FIXED_UP_1   = (1 << 30) // must never be set by compiler
	PROTOCOL_IS_CANONICAL = (1 << 29) // must never be set by compiler
	// Bits 0..15 are reserved for Swift's use.
	PROTOCOL_FIXED_UP_MASK = (PROTOCOL_FIXED_UP_1 | PROTOCOL_FIXED_UP_2)
)

type ProtocolList struct {
	Count     uint64
	Protocols []uint64
}

type ProtocolT struct {
	IsaVMAddr                     uint64
	NameVMAddr                    uint64
	ProtocolsVMAddr               uint64
	InstanceMethodsVMAddr         uint64
	ClassMethodsVMAddr            uint64
	OptionalInstanceMethodsVMAddr uint64
	OptionalClassMethodsVMAddr    uint64
	InstancePropertiesVMAddr      uint64
	Size                          uint32
	Flags                         uint32
	// Fields below this point are not always present on disk.
	ExtendedMethodTypesVMAddr uint64
	DemangledNameVMAddr       uint64
	ClassPropertiesVMAddr     uint64
}

type Protocol struct {
	Name                    string
	Ptr                     uint64
	Isa                     *Class
	Prots                   []Protocol
	InstanceMethods         []Method
	InstanceProperties      []Property
	ClassMethods            []Method
	OptionalInstanceMethods []Method
	OptionalClassMethods    []Method
	ExtendedMethodTypes     string
	DemangledName           string
	ProtocolT
}

func (p *Protocol) dump(verbose bool) string {
	var props string
	var cMethods string
	var iMethods string
	var optMethods string

	protocol := fmt.Sprintf("@protocol %s ", p.Name)

	if len(p.Prots) > 0 {
		var subProts []string
		for _, prot := range p.Prots {
			subProts = append(subProts, prot.Name)
		}
		protocol += fmt.Sprintf("<%s>", strings.Join(subProts, ", "))
	}
	protocol += fmt.Sprintf(" // %#x", p.Ptr)
	if len(p.InstanceProperties) > 0 {
		for _, prop := range p.InstanceProperties {
			if verbose {
				props += fmt.Sprintf("@property %s%s;\n", getPropertyAttributeTypes(prop.Attributes), prop.Name)
			} else {
				props += fmt.Sprintf("@property (%s) %s;\n", prop.Attributes, prop.Name)
			}
		}
	}
	if len(p.ClassMethods) > 0 {
		cMethods = "/* class methods */\n"
		for _, meth := range p.ClassMethods {
			if verbose {
				rtype, args := decodeMethodTypes(meth.Types)
				cMethods += fmt.Sprintf("+ %s\n", getMethodWithArgs(meth.Name, rtype, args))
			} else {
				cMethods += fmt.Sprintf("+[%s %s];\n", p.Name, meth.Name)
			}
		}
	}
	if len(p.InstanceMethods) > 0 {
		iMethods = "/* instance methods */\n"
		for _, meth := range p.InstanceMethods {
			if verbose {
				rtype, args := decodeMethodTypes(meth.Types)
				iMethods += fmt.Sprintf("- %s\n", getMethodWithArgs(meth.Name, rtype, args))
			} else {
				iMethods += fmt.Sprintf("-[%s %s];\n", p.Name, meth.Name)
			}
		}
	}
	if len(p.OptionalInstanceMethods) > 0 {
		optMethods = "@optional\n/* instance methods */\n"
		for _, meth := range p.OptionalInstanceMethods {
			if verbose {
				rtype, args := decodeMethodTypes(meth.Types)
				optMethods += fmt.Sprintf("- %s\n", getMethodWithArgs(meth.Name, rtype, args))
			} else {
				optMethods += fmt.Sprintf("-[%s %s];\n", p.Name, meth.Name)
			}
		}
	}
	return fmt.Sprintf(
		"%s\n"+
			"%s%s%s%s"+
			"@end\n",
		protocol,
		props,
		cMethods,
		iMethods,
		optMethods,
	)
}

func (p *Protocol) String() string {
	return p.dump(false)
}
func (p *Protocol) Verbose() string {
	return p.dump(true)
}

// CFString object in a 64-bit MachO file
type CFString struct {
	Name    string
	ISA     string
	Address uint64
	Class   *Class
	CFString64Type
}

// CFString64Type object in a 64-bit MachO file
type CFString64Type struct {
	IsaVMAddr uint64 // class64_t * (64-bit pointer)
	Info      uint64 // flag bits
	Data      uint64 // char * (64-bit pointer)
	Length    uint64 // number of non-NULL characters in above
}

type Class struct {
	Name                  string
	SuperClass            string
	Isa                   string
	InstanceMethods       []Method
	ClassMethods          []Method
	Ivars                 []Ivar
	Props                 []Property
	Protocols             []Protocol
	ClassPtr              uint64
	IsaVMAddr             uint64
	SuperclassVMAddr      uint64
	MethodCacheBuckets    uint64
	MethodCacheProperties uint64
	DataVMAddr            uint64
	IsSwiftLegacy         bool
	IsSwiftStable         bool
	ReadOnlyData          ClassRO64
}

func (c *Class) dump(verbose bool) string {
	var iVars string
	var props string
	var isSwift string
	var cMethods string
	var iMethods string

	var subClass string
	if c.ReadOnlyData.Flags.IsRoot() {
		subClass = "<ROOT>"
	} else if len(c.SuperClass) > 0 {
		subClass = c.SuperClass
	}

	if c.IsSwift() {
		isSwift = " (Swift)"
	}

	class := fmt.Sprintf("@interface %s : %s", c.Name, subClass)

	if len(c.Protocols) > 0 {
		var subProts []string
		for _, prot := range c.Protocols {
			subProts = append(subProts, prot.Name)
		}
		class += fmt.Sprintf("<%s>", strings.Join(subProts, ", "))
	}
	if len(c.Ivars) > 0 {
		s := bytes.NewBufferString("")
		w := tabwriter.NewWriter(s, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, " { // %#x%s\n  /* instance variables */\t// +size   offset\n", c.ClassPtr, isSwift)
		// s.WriteString(fmt.Sprintf(" { // %#x\n  // instance variables\t   +size   offset\n", c.ClassPtr))
		for _, ivar := range c.Ivars {
			if verbose {
				fmt.Fprintf(w, "  %s\n", ivar.Verbose())
				// s.WriteString(fmt.Sprintf("  %s\n", ivar.Verbose()))
			} else {
				fmt.Fprintf(w, "  %s\n", &ivar)
				// s.WriteString(fmt.Sprintf("  %s\n", &ivar))
			}
		}
		w.Flush()
		s.WriteString("}\n\n")
		iVars = s.String()
	} else {
		iVars = fmt.Sprintf(" { // %#x%s\n", c.ClassPtr, isSwift)
	}
	if len(c.Props) > 0 {
		for _, prop := range c.Props {
			if verbose {
				props += fmt.Sprintf("@property %s%s;\n", getPropertyAttributeTypes(prop.Attributes), prop.Name)
			} else {
				props += fmt.Sprintf("@property (%s) %s;\n", prop.Attributes, prop.Name)
			}
		}
		props += "\n"
	}
	if len(c.ClassMethods) > 0 {
		s := bytes.NewBufferString("/* class methods */\n\n")
		w := tabwriter.NewWriter(s, 0, 0, 1, ' ', 0)
		for _, meth := range c.ClassMethods {
			if verbose {
				rtype, args := decodeMethodTypes(meth.Types)
				// fmt.Fprintf(w, "+ %s;\t// %#x\n", getMethodWithArgs(meth.Name, rtype, args), meth.ImpVMAddr)
				s.WriteString(fmt.Sprintf("// %#x\n", meth.ImpVMAddr))
				s.WriteString(fmt.Sprintf("+ %s\n", getMethodWithArgs(meth.Name, rtype, args)))
				// s.WriteString(fmt.Sprintf("+ %-80s // %#x\n", getMethodWithArgs(meth.Name, rtype, args), meth.ImpVMAddr))
			} else {
				fmt.Fprintf(w, "+[%s %s];\t// %#x\n", c.Name, meth.Name, meth.ImpVMAddr)
				// s.WriteString(fmt.Sprintf("+[%s %s]; %-40s\n", c.Name, meth.Name, fmt.Sprintf("// %#x", meth.ImpVMAddr)))
			}
		}
		w.Flush()
		cMethods = s.String()
	}
	if len(c.InstanceMethods) > 0 {
		var s *bytes.Buffer
		if len(c.ClassMethods) > 0 {
			s = bytes.NewBufferString("\n/* instance methods */\n\n")
		} else {
			s = bytes.NewBufferString("/* instance methods */\n\n")
		}
		w := tabwriter.NewWriter(s, 0, 0, 1, ' ', 0)
		for _, meth := range c.InstanceMethods {
			if verbose {
				rtype, args := decodeMethodTypes(meth.Types)
				// fmt.Fprintf(w, "- %s;\t// %#x\n", getMethodWithArgs(meth.Name, rtype, args), meth.ImpVMAddr)
				s.WriteString(fmt.Sprintf("// %#x\n", meth.ImpVMAddr))
				s.WriteString(fmt.Sprintf("- %s\n", getMethodWithArgs(meth.Name, rtype, args)))
				// s.WriteString(fmt.Sprintf("- %-80s // %#x\n", getMethodWithArgs(meth.Name, rtype, args), meth.ImpVMAddr))
			} else {
				fmt.Fprintf(w, "-[%s %s];\t// %#x\n", c.Name, meth.Name, meth.ImpVMAddr)
				// s.WriteString(fmt.Sprintf("-[%s %s]; %-40s\n", c.Name, meth.Name, fmt.Sprintf("// %#x", meth.ImpVMAddr)))
			}
		}
		w.Flush()
		iMethods = s.String()
	}

	return fmt.Sprintf(
		"%s%s%s%s%s@end\n",
		class,
		iVars,
		props,
		cMethods,
		iMethods)
}

func (c *Class) IsSwift() bool {
	return c.IsSwiftLegacy || c.IsSwiftStable
}
func (c *Class) String() string {
	return c.dump(false)
}
func (c *Class) Verbose() string {
	return c.dump(true)
}

type ObjcClassT struct {
	IsaVMAddr              uint32
	SuperclassVMAddr       uint32
	MethodCacheBuckets     uint32
	MethodCacheProperties  uint32
	DataVMAddrAndFastFlags uint32
}

type SwiftClassMetadata struct {
	ObjcClassT
	SwiftClassFlags uint32
}

const (
	FAST_IS_SWIFT_LEGACY = 1 << 0 // < 5
	FAST_IS_SWIFT_STABLE = 1 << 1 // 5.X
	FAST_HAS_DEFAULT_RR  = 1 << 2
	IsSwiftPreStableABI  = 0x1
)

const (
	FAST_DATA_MASK  = 0xfffffffc
	FAST_FLAGS_MASK = 0x00000003

	FAST_DATA_MASK64_IPHONE = 0x0000007ffffffff8
	FAST_DATA_MASK64        = 0x00007ffffffffff8
	FAST_FLAGS_MASK64       = 0x0000000000000007
	FAST_IS_RW_POINTER64    = 0x8000000000000000
)

type ClassRoFlags uint32

const (
	// class is a metaclass
	RO_META ClassRoFlags = (1 << 0)
	// class is a root class
	RO_ROOT ClassRoFlags = (1 << 1)
	// class has .cxx_construct/destruct implementations
	RO_HAS_CXX_STRUCTORS ClassRoFlags = (1 << 2)
	// class has +load implementation
	RO_HAS_LOAD_METHOD ClassRoFlags = (1 << 3)
	// class has visibility=hidden set
	RO_HIDDEN ClassRoFlags = (1 << 4)
	// class has attributeClassRoFlags = (objc_exception): OBJC_EHTYPE_$_ThisClass is non-weak
	RO_EXCEPTION ClassRoFlags = (1 << 5)
	// class has ro field for Swift metadata initializer callback
	RO_HAS_SWIFT_INITIALIZER ClassRoFlags = (1 << 6)
	// class compiled with ARC
	RO_IS_ARC ClassRoFlags = (1 << 7)
	// class has .cxx_destruct but no .cxx_construct ClassRoFlags = (with RO_HAS_CXX_STRUCTORS)
	RO_HAS_CXX_DTOR_ONLY ClassRoFlags = (1 << 8)
	// class is not ARC but has ARC-style weak ivar layout
	RO_HAS_WEAK_WITHOUT_ARC ClassRoFlags = (1 << 9)
	// class does not allow associated objects on instances
	RO_FORBIDS_ASSOCIATED_OBJECTS ClassRoFlags = (1 << 10)

	// class is in an unloadable bundle - must never be set by compiler
	RO_FROM_BUNDLE ClassRoFlags = (1 << 29)
	// class is unrealized future class - must never be set by compiler
	RO_FUTURE ClassRoFlags = (1 << 30)
	// class is realized - must never be set by compiler
	RO_REALIZED ClassRoFlags = (1 << 31)
)

func (f ClassRoFlags) IsMeta() bool {
	return (f & RO_META) != 0
}
func (f ClassRoFlags) IsRoot() bool {
	return (f & RO_ROOT) != 0
}
func (f ClassRoFlags) HasCxxStructors() bool {
	return (f & RO_HAS_CXX_STRUCTORS) != 0
}
func (f ClassRoFlags) HasFuture() bool {
	return (f & RO_FUTURE) != 0
}

type ClassRO struct {
	Flags                ClassRoFlags
	InstanceStart        uint32
	InstanceSize         uint32
	_                    uint32
	IvarLayoutVMAddr     uint32
	NameVMAddr           uint32
	BaseMethodsVMAddr    uint32
	BaseProtocolsVMAddr  uint32
	IvarsVMAddr          uint32
	WeakIvarLayoutVMAddr uint32
	BasePropertiesVMAddr uint32
}

type ObjcClass64 struct {
	IsaVMAddr              uint64
	SuperclassVMAddr       uint64
	MethodCacheBuckets     uint64
	MethodCacheProperties  uint64
	DataVMAddrAndFastFlags uint64
}

type SwiftClassMetadata64 struct {
	ObjcClass64
	SwiftClassFlags uint64
}

type ClassRO64 struct {
	Flags         ClassRoFlags
	InstanceStart uint32
	InstanceSize  uint64
	// _                    uint32
	IvarLayoutVMAddr     uint64
	NameVMAddr           uint64
	BaseMethodsVMAddr    uint64
	BaseProtocolsVMAddr  uint64
	IvarsVMAddr          uint64
	WeakIvarLayoutVMAddr uint64
	BasePropertiesVMAddr uint64
}

type IvarList struct {
	EntSize uint32
	Count   uint32
}

type IvarT struct {
	Offset      uint64 // uint32_t*  (uint64_t* on x86_64)
	NameVMAddr  uint64 // const char*
	TypesVMAddr uint64 // const char*
	Alignment   uint32
	Size        uint32
}

type Ivar struct {
	Name   string
	Type   string
	Offset uint32
	IvarT
}

func (i *Ivar) dump(verbose bool) string {
	if verbose {
		ivtype := getIVarType(i.Type)
		if strings.ContainsAny(ivtype, "[]") { // array special case
			ivtype = strings.TrimSpace(strings.Replace(ivtype, "x", i.Name, 1))
			return fmt.Sprintf("%s;\t// %-7s %#x", ivtype, fmt.Sprintf("+%#x", i.Size), i.Offset)
		}
		return fmt.Sprintf("%s%s;\t// %-7s %#x", ivtype, i.Name, fmt.Sprintf("+%#x", i.Size), i.Offset)
	}
	return fmt.Sprintf("%s %s;\t// %-7s %#x", i.Type, i.Name, fmt.Sprintf("+%#x", i.Size), i.Offset)
}

func (i *Ivar) String() string {
	return i.dump(false)
}
func (i *Ivar) Verbose() string {
	return i.dump(true)
}

type Selector struct {
	VMAddr uint64
	Name   string
}

type OptOffsets struct {
	MethodNameStart     uint64
	MethodNameEnd       uint64
	InlinedMethodsStart uint64
	InlinedMethodsEnd   uint64
}

type OptOffsets2 struct {
	Version             uint64
	MethodNameStart     uint64
	MethodNameEnd       uint64
	InlinedMethodsStart uint64
	InlinedMethodsEnd   uint64
}

type ImpCache struct {
	PreoptCacheT
	Entries []PreoptCacheEntryT
}
type PreoptCacheEntryT struct {
	SelOffset uint32
	ImpOffset uint32
}

type PreoptCacheT struct {
	FallbackClassOffset int32
	Info                uint32
	// uint32_t cache_shift :  5
	// uint32_t cache_mask  : 11
	// uint32_t occupied    : 14
	// uint32_t has_inlines :  1
	// uint32_t bit_one     :  1
}

func (p PreoptCacheT) CacheShift() uint32 {
	return uint32(types.ExtractBits(uint64(p.Info), 0, 5))
}
func (p PreoptCacheT) CacheMask() uint32 {
	return uint32(types.ExtractBits(uint64(p.Info), 5, 11))
}
func (p PreoptCacheT) Occupied() uint32 {
	return uint32(types.ExtractBits(uint64(p.Info), 16, 14))
}
func (p PreoptCacheT) HasInlines() bool {
	return types.ExtractBits(uint64(p.Info), 30, 1) != 0
}
func (p PreoptCacheT) BitOne() bool {
	return types.ExtractBits(uint64(p.Info), 31, 1) != 0
}
func (p PreoptCacheT) Capacity() uint32 {
	return p.CacheMask() + 1
}
func (p PreoptCacheT) String() string {
	return fmt.Sprintf("cache_shift: %d, cache_mask: %d, occupied: %d, has_inlines: %t, bit_one: %t",
		p.CacheShift(),
		p.CacheMask(),
		p.Occupied(),
		p.HasInlines(),
		p.BitOne())
}

type Stub struct {
	Name        string
	SelectorRef uint64
}

type IntObj struct {
	ISA          uint64
	EncodingAddr uint64
	Number       uint64
}

type ImpCache2 struct {
	PreoptCache2T
	Entries []PreoptCacheEntry2T
}
type PreoptCacheEntry2T struct {
	ImpOffset int64
	SelOffset uint64
}

func (e PreoptCacheEntry2T) GetImpOffset() int64 {
	return int64(types.ExtractBits(uint64(e.ImpOffset), 0, 38))
}
func (e PreoptCacheEntry2T) GetSelOffset() uint32 {
	return uint32(types.ExtractBits(uint64(e.SelOffset), 0, 26))
}

type PreoptCache2T struct { // FIXME: 64bit new version
	FallbackClassOffset int64
	Info                uint64
	// int64_t  fallback_class_offset;
	// union {
	//     struct {
	//         uint16_t shift       :  5;
	//         uint16_t mask        : 11;
	//     };
	//     uint16_t hash_params;
	// };
	// uint16_t occupied    : 14;
	// uint16_t has_inlines :  1;
	// uint16_t padding     :  1;
	// uint32_t unused      : 31;
	// uint32_t bit_one     :  1;
	// preopt_cache_entry_t entries[];
}

func (p PreoptCache2T) CacheShift() uint32 {
	return uint32(types.ExtractBits(uint64(p.Info), 0, 5))
}
func (p PreoptCache2T) CacheMask() uint32 {
	return uint32(types.ExtractBits(uint64(p.Info), 5, 11))
}
func (p PreoptCache2T) Occupied() uint32 {
	return uint32(types.ExtractBits(uint64(p.Info), 16, 14))
}
func (p PreoptCache2T) HasInlines() bool {
	return types.ExtractBits(uint64(p.Info), 30, 1) != 0
}
func (p PreoptCache2T) BitOne() bool {
	return types.ExtractBits(uint64(p.Info), 63, 1) != 0
}
func (p PreoptCache2T) Capacity() uint32 {
	return p.CacheMask() + 1
}
func (p PreoptCache2T) String() string {
	return fmt.Sprintf("cache_shift: %d, cache_mask: %d, occupied: %d, has_inlines: %t, bit_one: %t",
		p.CacheShift(),
		p.CacheMask(),
		p.Occupied(),
		p.HasInlines(),
		p.BitOne())
}
