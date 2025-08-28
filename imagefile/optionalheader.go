package imagefile

import "encoding/binary"

// OptionalHeader is a common interface implemented by both the PE32 and P32+
// optional header formats.
type OptionalHeader interface {
	Subsystem() Subsystem
	DataDirectories() []DataDirectory
}

// MinOptionalHeaderSize32 is the minimum size of a PE32 optional header.
const MinOptionalHeaderSize32 = 96

// OptionalHeader32 maps values within the Optional Header section of a
// PE/COFF file that uses the PE32 format.
type OptionalHeader32 []byte

// Subsystem returns the subsystem responsible for executing the image.
func (header OptionalHeader32) Subsystem() Subsystem {
	if len(header) < MinOptionalHeaderSize32 {
		return 0
	}
	return Subsystem(binary.LittleEndian.Uint16(header[68:70]))
}

// NumberOfDataDirectories returns the the number of data directories
// declared by the header.
func (header OptionalHeader32) NumberOfDataDirectories() uint32 {
	if len(header) < MinOptionalHeaderSize32 {
		return 0
	}
	return binary.LittleEndian.Uint32(header[92:96])
}

// DataDirectories returns the set of data directories.
func (header OptionalHeader32) DataDirectories() []DataDirectory {
	if len(header) <= MinOptionalHeaderSize32 {
		return nil
	}
	return makeDataDirectories(header[MinOptionalHeaderSize32:], header.NumberOfDataDirectories())
}

// MinOptionalHeaderSize64 is the minimum size of a PE32+ optional header.
const MinOptionalHeaderSize64 = 112

// OptionalHeader32 maps values within the Optional Header section of a
// PE/COFF file that uses the PE32+ format.
type OptionalHeader64 []byte

// Subsystem returns the subsystem responsible for executing the image.
func (header OptionalHeader64) Subsystem() Subsystem {
	if len(header) < MinOptionalHeaderSize64 {
		return 0
	}
	return Subsystem(binary.LittleEndian.Uint16(header[68:70]))
}

// NumberOfDataDirectories returns the the number of data directories
// declared by the header.
func (header OptionalHeader64) NumberOfDataDirectories() uint32 {
	if len(header) < MinOptionalHeaderSize64 {
		return 0
	}
	return binary.LittleEndian.Uint32(header[108:112])
}

// DataDirectories returns the set of data directories.
func (header OptionalHeader64) DataDirectories() []DataDirectory {
	if len(header) <= MinOptionalHeaderSize64 {
		return nil
	}
	return makeDataDirectories(header[MinOptionalHeaderSize64:], header.NumberOfDataDirectories())
}

func makeDataDirectories(data []byte, count uint32) []DataDirectory {
	if count == 0 {
		return nil
	}
	dirs := make([]DataDirectory, 0, count)
	for i := range count {
		start := i * DataDirectorySize
		end := start + DataDirectorySize
		if int(end) > len(data) {
			break
		}
		dirs = append(dirs, DataDirectory(data[start:end]))
	}
	return dirs
}
