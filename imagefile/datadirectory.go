package imagefile

import (
	"encoding/binary"
	"fmt"
)

// DirectoryID identifies a data directory by its index within the directory
// table.
type DirectoryID int

// Known data directories.
const (
	ExportTableID           = DirectoryID(iota) // The export table address and size.
	ImportTableID                               // The import table address and size.
	ResourceTableID                             // The resource table address and size.
	ExceptionTableID                            // The exception table address and size.
	CertificateTableID                          // The attribute certificate table address and size. Uses a file offset instead of a virtual address.
	BaseRelocationTableID                       // The base relocation table address and size.
	DebugID                                     // The debug data starting address and size.
	ArchitectureID                              // Reserved, must be 0.
	GlobalPtrID                                 // The RVA of the value to be stored in the global pointer register. The size member of this structure must be set to zero.
	TLSTableID                                  // The thread local storage (TLS) table address and size.
	LoadConfigTableID                           // The load configuration table address and size.
	BoundImportID                               // The bound import table address and size.
	IATID                                       // The import address table address and size.
	DelayImportDescriptorID                     // The delay import descriptor address and size.
	CLRRuntimeHeaderID                          // The CLR runtime header address and size.
	ReservedID                                  // Reserved, must be zero.
)

// IsVirtual returns true if the data directory's address is a relative
// virtual address.
func (id DirectoryID) IsVirtual() bool {
	switch id {
	case CertificateTableID:
		return false
	default:
		return true
	}
}

// IsPointer returns true if the data directory has an address but no size.
func (id DirectoryID) IsPointer() bool {
	switch id {
	case GlobalPtrID:
		return true
	default:
		return false
	}
}

// String returns a string representation of the data directory ID.
func (id DirectoryID) String() string {
	switch id {
	case ExportTableID:
		return "Export Table"
	case ImportTableID:
		return "Import Table"
	case ResourceTableID:
		return "Resource Table"
	case ExceptionTableID:
		return "Exception Table"
	case CertificateTableID:
		return "Certificate Table"
	case BaseRelocationTableID:
		return "Base Relocation Table"
	case DebugID:
		return "Debug"
	case ArchitectureID:
		return "Architecture"
	case GlobalPtrID:
		return "Global Ptr"
	case TLSTableID:
		return "TLS Table"
	case LoadConfigTableID:
		return "Load Config Table"
	case BoundImportID:
		return "Bound Import"
	case IATID:
		return "IAT"
	case DelayImportDescriptorID:
		return "Delay Import Descriptor"
	case CLRRuntimeHeaderID:
		return "CLR Runtime Header"
	case ReservedID:
		return "Reserved"
	default:
		return fmt.Sprintf("Directory %d", id)
	}
}

// DataDirectory a describes a data directory entry within an image file.
type DataDirectory []byte

// IsZero returns true if the data directory entry contains zeroed data.
func (dir DataDirectory) IsZero() bool {
	for _, value := range dir {
		if value != 0 {
			return false
		}
	}
	return true
}

// Address returns the address of the data directory, which is usually
// a relative virtual address, but may be a file offset in some circumstances.
func (dir DataDirectory) Address() uint {
	return uint(binary.LittleEndian.Uint32(dir[0:4]))
}

// Size returns the size of the data directory.
func (dir DataDirectory) Size() uint32 {
	return binary.LittleEndian.Uint32(dir[4:8])
}
