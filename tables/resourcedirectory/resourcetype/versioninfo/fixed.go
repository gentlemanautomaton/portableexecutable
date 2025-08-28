package versioninfo

import (
	"encoding/binary"
	"fmt"
)

const fixedFileInfoLength = 52

// Fixed holds the data for a fixed version info structure.
type FixedFileInfo struct {
	data []byte
}

// Valid returns true if the fixed file info is the expected size and has the
// expected signature.
func (info FixedFileInfo) Valid() bool {
	if len(info.data) < fixedFileInfoLength {
		return false
	}
	signature := binary.LittleEndian.Uint32(info.data[0:4])
	return signature == 0xFEEF04BD
}

// FileVersion returns the version of the file itself.
// It returns a zero value if the information is invalid.
func (info FixedFileInfo) FileVersion() FixedVersion {
	if len(info.data) < fixedFileInfoLength {
		return 0
	}
	return makeFixedVersion(info.data[8:16])
}

// ProductVersion returns the version of the product the file belongs to.
// It returns a zero value if the information is invalid.
func (info FixedFileInfo) ProductVersion() FixedVersion {
	if len(info.data) < fixedFileInfoLength {
		return 0
	}
	return makeFixedVersion(info.data[16:24])
}

// FixedVersion holds fixed version information.
type FixedVersion uint64

func makeFixedVersion(data []byte) FixedVersion {
	ms := FixedVersion(binary.LittleEndian.Uint32(data[0:4]))
	ls := FixedVersion(binary.LittleEndian.Uint32(data[4:8]))
	return (ms << 32) | ls
}

// Major returns the major version number.
func (v FixedVersion) Major() uint16 {
	return uint16((v >> 48) & 0xFFFF)
}

// Minor returns the minor version number.
func (v FixedVersion) Minor() uint16 {
	return uint16((v >> 32) & 0xFFFF)
}

// Build returns the build number.
func (v FixedVersion) Build() uint16 {
	return uint16((v >> 16) & 0xFFFF)
}

// Revision returns the revision number.
func (v FixedVersion) Revision() uint16 {
	return uint16(v & 0xFFFF)
}

// IsZero returns true if all fields of v are zero.
func (v FixedVersion) IsZero() bool {
	return v == 0
}

// String returns a string representation of the fixed version in the format
// "Major.Minor.Build.Revision".
func (v FixedVersion) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", v.Major(), v.Minor(), v.Build(), v.Revision())
}
