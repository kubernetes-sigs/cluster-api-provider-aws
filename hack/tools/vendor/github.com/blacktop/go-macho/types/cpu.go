package types

import "fmt"

// A CPU is a Mach-O cpu type.
type CPU uint32

const (
	cpuArchMask = 0xff000000 //  mask for architecture bits
	cpuArch64   = 0x01000000 // 64 bit ABI
	cpuArch6432 = 0x02000000 // ABI for 64-bit hardware with 32-bit types; LP32
)

const (
	CPUVax     CPU = 1
	CPUMC680x0 CPU = 6
	CPUX86     CPU = 7
	CPUI386    CPU = CPUX86 /* compatibility */
	CPUAmd64   CPU = CPUX86 | cpuArch64
	CPUMips    CPU = 8
	CPUMc98000 CPU = 10
	CPUHppa    CPU = 11
	CPUArm     CPU = 12
	CPUArm64   CPU = CPUArm | cpuArch64
	CPUArm6432 CPU = CPUArm | cpuArch6432
	CPUMc88000 CPU = 13
	CPUSparc   CPU = 14
	CPUI860    CPU = 15
	CPUPpc     CPU = 18
	CPUPpc64   CPU = CPUPpc | cpuArch64
)

var cpuStrings = []IntName{
	{uint32(CPUVax), "VAX"},
	{uint32(CPUMC680x0), "MC680x0"},
	{uint32(CPUI386), "i386"},
	{uint32(CPUAmd64), "Amd64"},
	{uint32(CPUMips), "MIPS"},
	{uint32(CPUMc98000), "MC98000"},
	{uint32(CPUHppa), "HPPA"},
	{uint32(CPUArm), "ARM"},
	{uint32(CPUArm64), "AARCH64"},
	{uint32(CPUArm6432), "ARM64_32"},
	{uint32(CPUMc88000), "MC88000"},
	{uint32(CPUSparc), "SPARC"},
	{uint32(CPUI860), "i860"},
	{uint32(CPUPpc), "PowerPC"},
	{uint32(CPUPpc64), "PowerPC 64"},
}

func (i CPU) String() string   { return StringName(uint32(i), cpuStrings, false) }
func (i CPU) GoString() string { return StringName(uint32(i), cpuStrings, true) }

type CPUSubtype uint32

// VAX subtypes
const (
	CPUSubtypeVaxAll  CPUSubtype = 0
	CPUSubtypeVax780  CPUSubtype = 1
	CPUSubtypeVax785  CPUSubtype = 2
	CPUSubtypeVax750  CPUSubtype = 3
	CPUSubtypeVax730  CPUSubtype = 4
	CPUSubtypeUVaxI   CPUSubtype = 5
	CPUSubtypeUVaxII  CPUSubtype = 6
	CPUSubtypeVax8200 CPUSubtype = 7
	CPUSubtypeVax8500 CPUSubtype = 8
	CPUSubtypeVax8600 CPUSubtype = 9
	CPUSubtypeVax8650 CPUSubtype = 10
	CPUSubtypeVax8800 CPUSubtype = 11
	CPUSubtypeUVaxIII CPUSubtype = 12
)

// 680x0 subtypes
const (
	CPUSubtypeMC680x0All  CPUSubtype = 1
	CPUSubtypeMC68030     CPUSubtype = 1
	CPUSubtypeMC68040     CPUSubtype = 2
	CPUSubtypeMC68030Only CPUSubtype = 3
)

// I386 subtypes
const (
	CPUSubtypeI386All           CPUSubtype = 3 + 0<<4
	CPUSubtypeI386386           CPUSubtype = 3 + 0<<4
	CPUSubtypeI386486           CPUSubtype = 4 + 0<<4
	CPUSubtypeI386486SX         CPUSubtype = 4 + 8<<4
	CPUSubtypeI386586           CPUSubtype = 5 + 0<<4
	CPUSubtypeI386Pent          CPUSubtype = 5 + 0<<4
	CPUSubtypeI386PentPro       CPUSubtype = 6 + 1<<4
	CPUSubtypeI386PentIIM3      CPUSubtype = 6 + 3<<4
	CPUSubtypeI386PentIIM5      CPUSubtype = 6 + 5<<4
	CPUSubtypeI386Celeron       CPUSubtype = 7 + 6<<4
	CPUSubtypeI386CeleronMobile CPUSubtype = 7 + 7<<4
	CPUSubtypeI386Pentium3      CPUSubtype = 8 + 0<<4
	CPUSubtypeI386Pentium3M     CPUSubtype = 8 + 1<<4
	CPUSubtypeI386Pentium3Xeon  CPUSubtype = 8 + 2<<4
	CPUSubtypeI386PentiumM      CPUSubtype = 9 + 0<<4
	CPUSubtypeI386Pentium4      CPUSubtype = 10 + 0<<4
	CPUSubtypeI386Pentium4M     CPUSubtype = 10 + 1<<4
	CPUSubtypeI386Itanium       CPUSubtype = 11 + 0<<4
	CPUSubtypeI386Itanium2      CPUSubtype = 11 + 1<<4
	CPUSubtypeI386Xeon          CPUSubtype = 12 + 0<<4
	CPUSubtypeI386XeonMP        CPUSubtype = 12 + 1<<4
)

// X86 subtypes
const (
	CPUSubtypeX86All   CPUSubtype = 3
	CPUSubtypeX8664All CPUSubtype = 3
	CPUSubtypeX86Arch1 CPUSubtype = 4
	CPUSubtypeX86_64H  CPUSubtype = 8
)

// Mips subtypes.
const (
	CPUSubtypeMipsAll    CPUSubtype = 0
	CPUSubtypeMipsR2300  CPUSubtype = 1
	CPUSubtypeMipsR2600  CPUSubtype = 2
	CPUSubtypeMipsR2800  CPUSubtype = 3
	CPUSubtypeMipsR2000a CPUSubtype = 4 // pmax
	CPUSubtypeMipsR2000  CPUSubtype = 5
	CPUSubtypeMipsR3000a CPUSubtype = 6 // 3max
	CPUSubtypeMipsR3000  CPUSubtype = 7
)

// MC98000 (PowerPC) subtypes
const (
	CPUSubtypeMc98000All CPUSubtype = 0
	CPUSubtypeMc98601    CPUSubtype = 1
)

// HPPA subtypes for Hewlett-Packard HP-PA family of risc processors. Port by NeXT to 700 series.
const (
	CPUSubtypeHppaAll    CPUSubtype = 0
	CPUSubtypeHppa7100   CPUSubtype = 0 // compat
	CPUSubtypeHppa7100LC CPUSubtype = 1
)

// MC88000 subtypes
const (
	CPUSubtypeMc88000All CPUSubtype = 0
	CPUSubtypeMc88100    CPUSubtype = 1
	CPUSubtypeMc88110    CPUSubtype = 2
)

// SPARC subtypes
const (
	CPUSubtypeSparcAll CPUSubtype = 0
)

// I860 subtypes
const (
	CPUSubtypeI860All  CPUSubtype = 0
	CPUSubtypeI860_860 CPUSubtype = 1
)

// PowerPC subtypes
const (
	CPUSubtypePowerPCAll   CPUSubtype = 0
	CPUSubtypePowerPC601   CPUSubtype = 1
	CPUSubtypePowerPC602   CPUSubtype = 2
	CPUSubtypePowerPC603   CPUSubtype = 3
	CPUSubtypePowerPC603e  CPUSubtype = 4
	CPUSubtypePowerPC603ev CPUSubtype = 5
	CPUSubtypePowerPC604   CPUSubtype = 6
	CPUSubtypePowerPC604e  CPUSubtype = 7
	CPUSubtypePowerPC620   CPUSubtype = 8
	CPUSubtypePowerPC750   CPUSubtype = 9
	CPUSubtypePowerPC7400  CPUSubtype = 10
	CPUSubtypePowerPC7450  CPUSubtype = 11
	CPUSubtypePowerPC970   CPUSubtype = 100
)

// ARM subtypes
const (
	CPUSubtypeArmAll    CPUSubtype = 0
	CPUSubtypeArmV4T    CPUSubtype = 5
	CPUSubtypeArmV6     CPUSubtype = 6
	CPUSubtypeArmV5Tej  CPUSubtype = 7
	CPUSubtypeArmXscale CPUSubtype = 8
	CPUSubtypeArmV7     CPUSubtype = 9
	CPUSubtypeArmV7F    CPUSubtype = 10
	CPUSubtypeArmV7S    CPUSubtype = 11
	CPUSubtypeArmV7K    CPUSubtype = 12
	CPUSubtypeArmV8     CPUSubtype = 13
	CPUSubtypeArmV6M    CPUSubtype = 14
	CPUSubtypeArmV7M    CPUSubtype = 15
	CPUSubtypeArmV7Em   CPUSubtype = 16
	CPUSubtypeArmV8M    CPUSubtype = 17
)

// ARM64 subtypes
const (
	CPUSubtypeArm64All CPUSubtype = 0
	CPUSubtypeArm64V8  CPUSubtype = 1
	CPUSubtypeArm64E   CPUSubtype = 2
)

// ARM64_32 subtypes
const (
	CPUSubtypeArm6432All CPUSubtype = 0
	CPUSubtypeArm6432V8  CPUSubtype = 1
)

// Capability bits used in the definition of cpu_subtype.
const (
	CpuSubtypeFeatureMask      CPUSubtype = 0xff000000                         /* mask for feature flags */
	CpuSubtypeMask                        = CPUSubtype(^CpuSubtypeFeatureMask) /* mask for cpu subtype */
	CpuSubtypeLib64                       = 0x80000000                         /* 64 bit libraries */
	CpuSubtypePtrauthAbi                  = 0x80000000                         /* pointer authentication with versioned ABI */
	CpuSubtypePtrauthAbiUser              = 0x40000000                         /* pointer authentication with userspace versioned ABI */
	CpuSubtypeArm64PtrAuthMask            = 0x0f000000
	/*
	 *      When selecting a slice, ANY will pick the slice with the best
	 *      grading for the selected cpu_type_t, unlike the "ALL" subtypes,
	 *      which are the slices that can run on any hardware for that cpu type.
	 */
	CpuSubtypeAny = -1
)

var cpuSubtypeX86Strings = []IntName{
	// {uint32(CPUSubtypeX86All), "x86"},
	{uint32(CPUSubtypeX8664All), "x86_64"},
	{uint32(CPUSubtypeX86Arch1), "x86 Arch1"},
	{uint32(CPUSubtypeX86_64H), "x86_64 (Haswell)"},
}
var cpuSubtypeArmStrings = []IntName{
	{uint32(CPUSubtypeArmAll), "ARM"},
	{uint32(CPUSubtypeArmV4T), "v4t"},
	{uint32(CPUSubtypeArmV6), "v6"},
	{uint32(CPUSubtypeArmV5Tej), "v5tej"},
	{uint32(CPUSubtypeArmXscale), "XScale"},
	{uint32(CPUSubtypeArmV7), "v7"},
	{uint32(CPUSubtypeArmV7F), "v7f"},
	{uint32(CPUSubtypeArmV7S), "v7s"},
	{uint32(CPUSubtypeArmV7K), "v7k"},
	{uint32(CPUSubtypeArmV8), "v8"},
	{uint32(CPUSubtypeArmV6M), "v6m"},
	{uint32(CPUSubtypeArmV7M), "v7m"},
	{uint32(CPUSubtypeArmV7Em), "v7em"},
	{uint32(CPUSubtypeArmV8M), "v8m"},
}
var cpuSubtypeArm64Strings = []IntName{
	{uint32(CPUSubtypeArm64All), "ARM64"},
	{uint32(CPUSubtypeArm64V8), "v8"},
	{uint32(CPUSubtypeArm64E), "ARM64e"},
	{uint32(CPUSubtypeArm6432All), "ARM64_32"},
	{uint32(CPUSubtypeArm6432V8), "v8"},
}

func (st CPUSubtype) String(cpu CPU) string {
	switch cpu {
	case CPUI386:
		fallthrough
	case CPUAmd64:
		return StringName(uint32(st&CpuSubtypeMask), cpuSubtypeX86Strings, false)
	case CPUArm:
		return StringName(uint32(st&CpuSubtypeMask), cpuSubtypeArmStrings, false)
	case CPUArm64:
		return StringName(uint32(st&CpuSubtypeMask), cpuSubtypeArm64Strings, false)
	case CPUArm6432:
		return StringName(uint32(st&CpuSubtypeMask), cpuSubtypeArm64Strings, false)
	}
	return "UNKNOWN"
}

func (st CPUSubtype) Caps(cpu CPU) string {
	switch cpu {
	case CPUArm64:
		caps := st & CpuSubtypeFeatureMask
		if caps > 0 {
			if caps&CpuSubtypePtrauthAbiUser == 0 {
				return fmt.Sprintf("USR%02d", (caps&CpuSubtypeArm64PtrAuthMask)>>24)
			} else {
				return fmt.Sprintf("KER%02d", (caps&CpuSubtypeArm64PtrAuthMask)>>24)
			}
		}
	}
	return ""
}

func (st CPUSubtype) GoString(cpu CPU) string {
	switch cpu {
	case CPUI386:
		fallthrough
	case CPUAmd64:
		return StringName(uint32(st&CpuSubtypeMask), cpuSubtypeX86Strings, true)
	case CPUArm:
		return StringName(uint32(st&CpuSubtypeMask), cpuSubtypeArmStrings, true)
	case CPUArm64:
		var feature string
		caps := st & CpuSubtypeFeatureMask
		if caps > 0 {
			if caps&CpuSubtypePtrauthAbiUser == 0 {
				feature = fmt.Sprintf(" caps: PAC%02d", (caps&CpuSubtypeArm64PtrAuthMask)>>24)
			} else {
				feature = fmt.Sprintf(" caps: PAK%02d", (caps&CpuSubtypeArm64PtrAuthMask)>>24)
			}
		}
		return StringName(uint32(st&CpuSubtypeMask), cpuSubtypeArm64Strings, true) + feature
	}
	return "UNKNOWN"
}
