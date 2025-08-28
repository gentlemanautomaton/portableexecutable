package imagefile

import "bytes"

// Signature is 4 byte value that identifies a file as a portable executable
// image file.
//
// A valid signature holds the well-known value 0x50450000 in big-endian byte
// order.
//
// https://learn.microsoft.com/en-us/windows/win32/debug/pe-format#signature-image-only
type Signature []byte

// Valid returns true if the signature contains the expected value 0x50450000.
func (s Signature) Valid() bool {
	return bytes.Equal(s, []byte{0x50, 0x45, 0x00, 0x00})
}
