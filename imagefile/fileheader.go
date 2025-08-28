package imagefile

import (
	"encoding/binary"
)

// FileHeader is a standard Common Object File Format header that is used by
// PE/COFF files.
type FileHeader []byte

// Machine returns the machine that the image file was built for.
func (header FileHeader) Machine() Machine {
	return Machine(binary.LittleEndian.Uint16(header[:2]))
}

// NumberOfSections returns the number of sections present in the image file.
func (header FileHeader) NumberOfSections() uint16 {
	return binary.LittleEndian.Uint16(header[2:4])
}

// PointerToSymbolTable returns the number of symbols present in the image file.
func (header FileHeader) PointerToSymbolTable() FileOffset {
	return FileOffset(binary.LittleEndian.Uint32(header[8:12]))
}

// NumberOfSymbols returns the number of symbols present in the image file.
func (header FileHeader) NumberOfSymbols() uint32 {
	return binary.LittleEndian.Uint32(header[12:16])
}

// SizeOfOptionalHeader returns the size of the optional header.
func (header FileHeader) SizeOfOptionalHeader() uint16 {
	return binary.LittleEndian.Uint16(header[16:18])
}
