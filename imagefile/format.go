package imagefile

import "fmt"

// Format describes the format of the image file's header.
type Format uint16

// Supported formats for the image file header.
const (
	PE32     Format = 0x10b
	PE32Plus Format = 0x20b
)

// Supported returns true if the format is supported.
func (format Format) Supported() bool {
	switch format {
	case PE32, PE32Plus:
		return true
	default:
		return false
	}
}

// String returns a string representation of the format.
func (format Format) String() string {
	switch format {
	case PE32:
		return "PE32 (32-bit)"
	case PE32Plus:
		return "PE32+ (64-bit)"
	default:
		return fmt.Sprintf("<unsupported PE header format \"%x\">", uint16(format))
	}
}
