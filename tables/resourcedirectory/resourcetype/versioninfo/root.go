package versioninfo

import "iter"

// Root is the root node of a version information tree.
type Root Node

// NewRoot interprets the given data as a version info tree and returns the
// root for it.
//
// If the data is too small or malformed it returns an error.
func NewRoot(data []byte) (Root, error) {
	node, err := NewNode(data)
	return Root(node), err
}

// Key returns the key that identifies the root node.
func (r Root) Key() string {
	return Node(r).Key()
}

// FileInfo returns the fixed file information contained in the root.
func (r Root) FileInfo() FixedFileInfo {
	value := Node(r).Value()
	return FixedFileInfo{data: value.Data()}
}

// Children returns an iterator for the root's child nodes.
func (r Root) Children() iter.Seq2[Node, error] {
	return Node(r).Children()
}
