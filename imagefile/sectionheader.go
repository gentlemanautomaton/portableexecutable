package imagefile

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

// SectionName holds the full name of a section if the name is less than
// eight characters long.
//
// If the section name exceeds eight characters, it instead encodes the
// address of a string within the COFF strings table.
type SectionName string

// Reference returns true if the section name is a reference to a string in
// the strings table. If it is, it also returns the offset of the string
// within the table.
func (name SectionName) Reference() (isReference bool, offset StringOffset) {
	if len(name) < 2 {
		return false, 0
	}
	if name[0] != '/' {
		return false, 0
	}
	number, err := strconv.ParseUint(string(name[1:]), 10, 32)
	if err != nil {
		return false, 0
	}
	return true, StringOffset(number)
}

// SectionHeader is the header for a standard Common Object File Format
// section within PE/COFF image files.
type SectionHeader []byte

// Name returns the name of the section header.
func (header SectionHeader) Name() SectionName {
	if cutoff := bytes.IndexByte(header[0:8], 0); cutoff >= 0 {
		return SectionName(header[0:cutoff])
	}
	return SectionName(header[0:8])
}

// VirtualSize returns the size of the section's virtual memory region.
func (header SectionHeader) VirtualSize() uint {
	return uint(binary.LittleEndian.Uint32(header[8:12]))
}

// VirtualAddress returns the start of the section's virtual memory region.
func (header SectionHeader) VirtualAddress() RelativeVirtualAddress {
	return RelativeVirtualAddress(binary.LittleEndian.Uint32(header[12:16]))
}

// SizeOfRawData returns the size of the section's image file region.
func (header SectionHeader) SizeOfRawData() uint {
	return uint(binary.LittleEndian.Uint32(header[16:20]))
}

// PointerToRawData returns the start of the section's image file region.
func (header SectionHeader) PointerToRawData() FileOffset {
	return FileOffset(binary.LittleEndian.Uint32(header[20:24]))
}
