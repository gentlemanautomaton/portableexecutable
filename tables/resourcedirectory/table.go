package resourcedirectory

import (
	"encoding/binary"

	"github.com/gentlemanautomaton/portableexecutable/imagefile"
)

// Table is a resource directory table that holds a set of resource directory
// entries.
type Table []Entry

// Index returns the index of the table entry with the given ID.
// It returns -1 if an entry with the requested ID does not exist.
func (table Table) Index(id ID) int {
	for i := range table {
		if table[i].ID == id {
			return i
		}
	}
	return -1
}

const tableHeaderSize = 16

type tableHeader []byte

func (header tableHeader) NumberOfNamedEntries() uint16 {
	return binary.LittleEndian.Uint16(header[12:14])
}

func (header tableHeader) NumberOfNumeberedEntries() uint16 {
	return binary.LittleEndian.Uint16(header[14:16])
}

type tableLayout struct {
	Named    int
	Numbered int
}

func (layout tableLayout) Total() int {
	return layout.Named + layout.Numbered
}

// Entry is a resource directory entry within a resource directory table.
type Entry struct {
	ID        ID
	Reference Reference
}

const tableEntrySize = 8

type namedEntry []byte

func (entry namedEntry) Name() StringOffset {
	offset := binary.LittleEndian.Uint32(entry[0:4])
	return StringOffset(offset & valueMask31)
}

func (entry namedEntry) Ref() Reference {
	return Reference{data: binary.LittleEndian.Uint32(entry[4:8])}
}

type numberedEntry []byte

func (entry numberedEntry) Number() uint32 {
	value := binary.LittleEndian.Uint32(entry[0:4])
	return value & valueMask31
}

func (entry numberedEntry) Ref() Reference {
	return Reference{data: binary.LittleEndian.Uint32(entry[4:8])}
}

const dataDescriptorSize = 16

type dataDescriptor []byte

func (descriptor dataDescriptor) Address() imagefile.RelativeVirtualAddress {
	return imagefile.RelativeVirtualAddress(binary.LittleEndian.Uint32(descriptor[0:4]))
}

func (descriptor dataDescriptor) Size() uint32 {
	return binary.LittleEndian.Uint32(descriptor[4:8])
}

func (descriptor dataDescriptor) Codepage() uint32 {
	return binary.LittleEndian.Uint32(descriptor[8:12])
}
