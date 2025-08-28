package dos

import (
	"encoding/binary"

	"github.com/gentlemanautomaton/portableexecutable/imagefile"
)

// Signature is the first 2 bytes of a DOS compatibility [Header], interpreted
// as a 16 bit unsigned integer in little-endian encoding. It's also known as
// the magic number.
//
// In a valid [Header], the signature is expected to hold the well-known
// value 0x5A4D.
type Signature uint16

// Valid returns true if the signature contains the expected value 0x5A4D.
func (s Signature) Valid() bool {
	return s == 0x5A4D
}

// Header is a 64 byte compatibility header at the start of PE image files
// which allows them to be handled gracefully by the DOS operating system.
type Header [64]byte

// Signature returns the two byte signature at the start of the header.
func (header *Header) Signature() Signature {
	return Signature(binary.LittleEndian.Uint16(header[:2]))
}

// NextHeader returns the address of the next header within an image file.
func (header *Header) NextHeader() imagefile.FileOffset {
	return imagefile.FileOffset(binary.LittleEndian.Uint32(header[60:]))
}
