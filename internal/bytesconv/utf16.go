package bytesconv

import (
	"encoding/binary"
	"unicode/utf16"
)

// DecodeUTF16 interprets the given bytes as UTF-16 with the specified byte
// order and returns the value as a string.
//
// Any invalid characters will be replaced with the unicode replacement
// character.
func DecodeUTF16(p []byte, order binary.ByteOrder) string {
	// If the number of bytes is not an even number, drop the last byte.
	n := len(p)
	if n%2 != 0 {
		n--
	}

	// If there is no data, return an empty string.
	if n < 1 {
		return ""
	}

	// Prepare a buffer to receive the bytes.
	buf := make([]uint16, 0, n/2)

	// Parse the bytes in the desired order.
	for i := 0; i+1 < n; i += 2 {
		buf = append(buf, order.Uint16(p[i:]))
	}

	// Decode the runes and convert them to a string.
	return string(utf16.Decode(buf))
}
