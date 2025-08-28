package versioninfo

import (
	"encoding/binary"

	"github.com/gentlemanautomaton/portableexecutable/internal/bytesconv"
)

// Value data types
const (
	dataTypeString = 1
)

type Value struct {
	data     []byte
	dataType uint16
}

func (v Value) Data() []byte {
	return v.data
}

func (v Value) String() string {
	if len(v.data) < 2 {
		return ""
	}
	if v.dataType != dataTypeString {
		return ""
	}
	data := v.data[:len(v.data)-2] // Exclude null terminator
	return bytesconv.DecodeUTF16(data, binary.LittleEndian)
}
