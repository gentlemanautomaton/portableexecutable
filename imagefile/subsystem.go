package imagefile

import "fmt"

// Subsystem identifies the subsystem that is responsible for executing an
// image.
type Subsystem uint16

// Known subsystem types.
//
// https://learn.microsoft.com/en-us/windows/win32/debug/pe-format#windows-subsystem
var (
	SubsystemUnknown                Subsystem = 0  // IMAGE_SUBSYSTEM_UNKNOWN, An unknown subsystem
	SubsystemNative                 Subsystem = 1  // IMAGE_SUBSYSTEM_NATIVE, Device drivers and native Windows processes
	SubsystemWindowsGUI             Subsystem = 2  // IMAGE_SUBSYSTEM_WINDOWS_GUI, The Windows graphical user interface (GUI) subsystem
	SubsystemWindowsCUI             Subsystem = 3  // IMAGE_SUBSYSTEM_WINDOWS_CUI, The Windows character subsystem
	SubsystemOS2CUI                 Subsystem = 5  // IMAGE_SUBSYSTEM_OS2_CUI, The OS/2 character subsystem
	SubsystemPosixCUI               Subsystem = 7  // IMAGE_SUBSYSTEM_POSIX_CUI, The Posix character subsystem
	SubsystemNativeWindows          Subsystem = 8  // IMAGE_SUBSYSTEM_NATIVE_WINDOWS, Native Win9x driver
	SubsystemWindowsCEGUI           Subsystem = 9  // IMAGE_SUBSYSTEM_WINDOWS_CE_GUI, Windows CE
	SubsystemEFIApplication         Subsystem = 10 // IMAGE_SUBSYSTEM_EFI_APPLICATION, An Extensible Firmware Interface (EFI) application
	SubsystemEFIBootServiceDriver   Subsystem = 11 // IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER, An EFI driver with boot services
	SubsystemEFIRuntimeDriver       Subsystem = 12 // IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER, An EFI driver with run-time services
	SubsystemEFIROM                 Subsystem = 13 // IMAGE_SUBSYSTEM_EFI_ROM, An EFI ROM image
	SubsystemXBOX                   Subsystem = 14 // IMAGE_SUBSYSTEM_XBOX, XBOX
	SubsystemWindowsBootApplication Subsystem = 16 // IMAGE_SUBSYSTEM_WINDOWS_BOOT_APPLICATION, Windows boot application.
)

// String returns a string representation of the subsystem.
func (subsystem Subsystem) String() string {
	switch subsystem {
	case SubsystemUnknown:
		return "Unknown"
	case SubsystemNative:
		return "Native"
	case SubsystemWindowsGUI:
		return "Windows GUI"
	case SubsystemWindowsCUI:
		return "Windows CUI"
	case SubsystemOS2CUI:
		return "OS2 CUI"
	case SubsystemPosixCUI:
		return "POSIX CUI"
	case SubsystemNativeWindows:
		return "Native Windows"
	case SubsystemWindowsCEGUI:
		return "Windows CE GUI"
	case SubsystemEFIApplication:
		return "EFI Application"
	case SubsystemEFIBootServiceDriver:
		return "EFI Boot Service Driver"
	case SubsystemEFIRuntimeDriver:
		return "EFI Runtime Driver"
	case SubsystemEFIROM:
		return "EFI ROM"
	case SubsystemXBOX:
		return "XBOX"
	case SubsystemWindowsBootApplication:
		return "Windows Boot Application"
	default:
		return fmt.Sprintf("<unrecognized subsystem code: %x>", uint16(subsystem))
	}
}
