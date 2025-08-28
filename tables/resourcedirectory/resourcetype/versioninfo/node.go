package versioninfo

import (
	"encoding/binary"
	"fmt"
	"iter"

	"github.com/gentlemanautomaton/portableexecutable/internal/bytesconv"
)

const minNodeSize = 8

// Node is a node within a version info tree.
type Node struct {
	data      []byte
	keyLength uint16 // Does not include the key's 2-byte null terminator.
}

// NewNode interprets the given data as a version info node and returns it.
//
// If the data is too small or malformed it returns an error.
func NewNode(data []byte) (Node, error) {
	if len(data) < minNodeSize {
		return Node{}, fmt.Errorf("the version info structure must be at least %d bytes", minNodeSize)
	}
	data = trim(data)
	keyLength := utf16Null(data[6:])
	if keyLength < 0 {
		return Node{}, fmt.Errorf("the version info structure's key field is not null-terminated")
	}
	return Node{
		data:      data,
		keyLength: uint16(keyLength),
	}, nil
}

// Key returns the key that identifies the node.
func (n Node) Key() string {
	const start = 6
	return bytesconv.DecodeUTF16(n.data[start:start+n.keyLength], binary.LittleEndian)
}

// Value returns the value of the node, if it has one.
func (n Node) Value() Value {
	length := binary.LittleEndian.Uint16(n.data[2:4])
	if length == 0 {
		return Value{}
	}

	dataType := binary.LittleEndian.Uint16(n.data[4:6])

	// If this is a string, the value length is probably in characters and not
	// in bytes. This is due to an ancient bug that became a defacto standard.
	//
	// See this article by Raymond Chen:
	//
	// The evolution of version resources – corrupted 32-bit version resources
	// https://devblogs.microsoft.com/oldnewthing/20061222-00/?p=28623
	if dataType == dataTypeString {
		length *= 2
	}

	// The offset of "8" added below is the sum of the following: node length,
	// value length and the key's null terminator.
	start := align(n.keyLength + 8)
	return Value{
		data:     n.data[start : start+length],
		dataType: binary.LittleEndian.Uint16(n.data[4:6]),
	}
}

// Children returns an iterator for the child nodes contained in n.
func (n Node) Children() iter.Seq2[Node, error] {
	valueStart := align(n.keyLength + 8)
	valueLength := binary.LittleEndian.Uint16(n.data[2:4])
	dataType := binary.LittleEndian.Uint16(n.data[4:6])

	// If this is a string, the value length is probably in characters and not
	// in bytes. This is due to an ancient bug that became a defacto standard.
	//
	// See this article by Raymond Chen:
	//
	// The evolution of version resources – corrupted 32-bit version resources
	// https://devblogs.microsoft.com/oldnewthing/20061222-00/?p=28623
	if dataType == dataTypeString {
		valueLength *= 2
	}

	offset := align(valueStart + valueLength)

	return func(yield func(Node, error) bool) {
		for int(offset) < len(n.data) {
			node, err := NewNode(n.data[offset:])
			if !yield(node, err) {
				return
			}
			if err != nil {
				return
			}
			offset += uint16(len(node.data))
			offset = align(offset)
		}
	}
}

// trim reads node length information from the start of data and returns the
// data truncated to that length.
func trim(data []byte) []byte {
	length := int(binary.LittleEndian.Uint16(data[0:2]))
	if len(data) > length {
		data = data[:length]
	}
	return data
}

// align returns offset rounded up to the nearest value that with a 4 byte
// (32-bit) alignment.
func align(offset uint16) uint16 {
	remainder := offset % 4
	if remainder == 0 {
		return offset
	}
	return offset + 4 - remainder
}

// utf16Null returns the index of the first UTF-16 null terminator in data.
func utf16Null(data []byte) int {
	for i := 0; i+1 < len(data); i += 2 {
		if data[i] == 0 && data[i+1] == 0 {
			return i
		}
	}
	return -1
}
