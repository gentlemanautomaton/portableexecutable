package imagefile

import "fmt"

// Machine identifies the target machine type for an image file.
type Machine uint16

// Known image file machine types.
//
// https://learn.microsoft.com/en-us/windows/win32/debug/pe-format#machine-types
var (
	MachineUnknown     Machine = 0x0    // IMAGE_FILE_MACHINE_UNKNOWN, The content of this field is assumed to be applicable to any machine type
	MachineAlpha       Machine = 0x184  // IMAGE_FILE_MACHINE_ALPHA, Alpha AXP, 32-bit address space
	MachineAlpha64     Machine = 0x284  // IMAGE_FILE_MACHINE_ALPHA64, Alpha 64, 64-bit address space
	MachineAM33        Machine = 0x1D3  // IMAGE_FILE_MACHINE_AM33, Matsushita AM33
	MachineAMD64       Machine = 0x8664 // IMAGE_FILE_MACHINE_AMD64, x64
	MachineARM         Machine = 0x1C0  // IMAGE_FILE_MACHINE_ARM, ARM little endian
	MachineARM64       Machine = 0xAA64 // IMAGE_FILE_MACHINE_ARM64, ARM64 little endian
	MachineARM64EC     Machine = 0xA641 // IMAGE_FILE_MACHINE_ARM64EC, ABI that enables interoperability between native ARM64 and emulated x64 code.
	MachineARM64X      Machine = 0xA64E // IMAGE_FILE_MACHINE_ARM64X, Binary format that allows both native ARM64 and ARM64EC code to coexist in the same file.
	MachineARMNT       Machine = 0x1C4  // IMAGE_FILE_MACHINE_ARMNT, ARM Thumb-2 little endian
	MachineAXP64       Machine = 0x284  // IMAGE_FILE_MACHINE_AXP64, AXP 64 (Same as Alpha 64)
	MachineEBC         Machine = 0xEBC  // IMAGE_FILE_MACHINE_EBC, EFI byte code
	MachineX86         Machine = 0x14C  // IMAGE_FILE_MACHINE_I386, Intel 386 or later processors and compatible processors
	MachineIA64        Machine = 0x200  // IMAGE_FILE_MACHINE_IA64, Intel Itanium processor family
	MachineLoongArch32 Machine = 0x6232 // IMAGE_FILE_MACHINE_LOONGARCH32, LoongArch 32-bit processor family
	MachineLoongArch64 Machine = 0x6264 // IMAGE_FILE_MACHINE_LOONGARCH64, LoongArch 64-bit processor family
	MachineM32R        Machine = 0x9041 // IMAGE_FILE_MACHINE_M32R, Mitsubishi M32R little endian
	MachineMIPS16      Machine = 0x266  // IMAGE_FILE_MACHINE_MIPS16, MIPS16
	MachineMIPSFPU     Machine = 0x366  // IMAGE_FILE_MACHINE_MIPSFPU, MIPS with FPU
	MachineMIPSFPU16   Machine = 0x466  // IMAGE_FILE_MACHINE_MIPSFPU16, MIPS16 with FPU
	MachinePowerPC     Machine = 0x1F0  // IMAGE_FILE_MACHINE_POWERPC, Power PC little endian
	MachinePowerPCFP   Machine = 0x1F1  // IMAGE_FILE_MACHINE_POWERPCFP, Power PC with floating point support
	MachineR3000BE     Machine = 0x160  // IMAGE_FILE_MACHINE_R3000BE, MIPS I compatible 32-bit big endian
	MachineR3000       Machine = 0x162  // IMAGE_FILE_MACHINE_R3000, MIPS I compatible 32-bit little endian
	MachineR4000       Machine = 0x166  // IMAGE_FILE_MACHINE_R4000, MIPS III compatible 64-bit little endian
	MachineR10000      Machine = 0x168  // IMAGE_FILE_MACHINE_R10000, MIPS IV compatible 64-bit little endian
	MachineRISCV32     Machine = 0x5032 // IMAGE_FILE_MACHINE_RISCV32, RISC-V 32-bit address space
	MachineRISCV64     Machine = 0x5064 // IMAGE_FILE_MACHINE_RISCV64, RISC-V 64-bit address space
	MachineRISCV128    Machine = 0x5128 // IMAGE_FILE_MACHINE_RISCV128, RISC-V 128-bit address space
	MachineSH3         Machine = 0x1A2  // IMAGE_FILE_MACHINE_SH3, Hitachi SH3
	MachineSH3DSP      Machine = 0x1A3  // IMAGE_FILE_MACHINE_SH3DSP, Hitachi SH3 DSP
	MachineSH4         Machine = 0x1A6  // IMAGE_FILE_MACHINE_SH4, Hitachi SH4
	MachineSH5         Machine = 0x1A8  // IMAGE_FILE_MACHINE_SH5, Hitachi SH5
	MachineThumb       Machine = 0x1C2  // IMAGE_FILE_MACHINE_THUMB, Thumb
	MachineWCEMIPSV2   Machine = 0x169  // IMAGE_FILE_MACHINE_WCEMIPSV2, MIPS little-endian WCE v2
)

// Supported returns true if the machine type is supported.
func (machine Machine) Supported() bool {
	switch machine {
	case MachineX86, MachineAMD64:
		return true
	default:
		return false
	}
}

// String returns a string representation of the machine.
func (machine Machine) String() string {
	switch machine {
	case MachineUnknown:
		return "Unknown"
	case MachineAlpha:
		return "Alpha"
	case MachineAlpha64:
		return "Alpha64"
	case MachineAM33:
		return "AM33"
	case MachineAMD64:
		return "AMD64"
	case MachineARM:
		return "ARM"
	case MachineARM64:
		return "ARM64"
	case MachineARM64EC:
		return "ARM64EC"
	case MachineARM64X:
		return "ARM64X"
	case MachineARMNT:
		return "ARMNT"
	case MachineAXP64:
		return "AXP64"
	case MachineEBC:
		return "EBC"
	case MachineX86:
		return "X86"
	case MachineIA64:
		return "IA64"
	case MachineLoongArch32:
		return "LoongArch32"
	case MachineLoongArch64:
		return "LoongArch64"
	case MachineM32R:
		return "M32R"
	case MachineMIPS16:
		return "MIPS16"
	case MachineMIPSFPU:
		return "MIPSFPU"
	case MachineMIPSFPU16:
		return "MIPSFPU16"
	case MachinePowerPC:
		return "PowerPC"
	case MachinePowerPCFP:
		return "PowerPCFP"
	case MachineR3000BE:
		return "R3000BE"
	case MachineR3000:
		return "R3000"
	case MachineR4000:
		return "R4000"
	case MachineR10000:
		return "R10000"
	case MachineRISCV32:
		return "RISCV32"
	case MachineRISCV64:
		return "RISCV64"
	case MachineRISCV128:
		return "RISCV128"
	case MachineSH3:
		return "SH3"
	case MachineSH3DSP:
		return "SH3DSP"
	case MachineSH4:
		return "SH4"
	case MachineSH5:
		return "SH5"
	case MachineThumb:
		return "Thumb"
	case MachineWCEMIPSV2:
		return "WCEMIPSV2"
	default:
		return fmt.Sprintf("<unrecognized machine code: %x>", uint16(machine))
	}
}
