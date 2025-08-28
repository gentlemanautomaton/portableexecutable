package resourcedirectory

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/gentlemanautomaton/portableexecutable"
	"github.com/gentlemanautomaton/portableexecutable/imagefile"
	"github.com/gentlemanautomaton/portableexecutable/internal/bytesconv"
	"github.com/gentlemanautomaton/portableexecutable/tables/resourcedirectory/resourcetype"
)

var (
	ErrMissingResourceTable = errors.New("the portable executable does not have a resource table")
)

// Reader reads resource table data for a portable executable image file
// from an underlying [portableexecutable.Reader].
type Reader struct {
	// source is an [io.SectionReader] that is limited to the range of bytes
	// that belong to the resource table, so the zero address is the start of
	// the resource table.
	source io.ReaderAt

	// pe is used to translate relative virtual addresses and retrieve
	// resource data.
	pe *portableexecutable.Reader
}

// NewReader creates and initializes a new resource directory [Reader] that
// reads from portable executable [portableexecutable.Reader] pe. It returns
// ErrMissingResourceTable if the portable executable does not have a resource
// table.
//
// TODO: Try to develop some sort of interface that removes the dependency
// on the portableexecutable package. This is challenging because the
// resource directory stores relative virtual addresses that must be
// translated with the help of the portable executable's section table.
func NewReader(pe *portableexecutable.Reader) (*Reader, error) {
	resources := pe.DataDirectories().Get(imagefile.ResourceTableID)
	if resources.IsZero() {
		return nil, ErrMissingResourceTable
	}

	return &Reader{
		source: io.NewSectionReader(pe.Source(), int64(resources.Location.Start), int64(resources.Location.Length)),
		pe:     pe,
	}, nil
}

// ReadRoot returns the root table of the resource directory.
func (r *Reader) ReadRoot() (Table, error) {
	return r.ReadTable(0)
}

// ReadTable returns the resource table with the given offset from the
// resource directory.
func (r *Reader) ReadTable(offset TableOffset) (Table, error) {
	layout, data, err := r.readTableData(offset)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, layout.Total())
	for i := range layout.Total() {
		start := i * tableEntrySize
		end := start + tableEntrySize
		if i < layout.Named {
			entry := namedEntry(data[start:end])
			name, err := r.ReadString(entry.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to read the ID string for resource directory table entry %d: %w", i, err)
			}
			entries = append(entries, Entry{
				ID:        ID{str: name},
				Reference: entry.Ref(),
			})
		} else {
			entry := numberedEntry(data[start:end])
			entries = append(entries, Entry{
				ID:        ID{number: entry.Number()},
				Reference: entry.Ref(),
			})
		}
	}

	return entries, nil
}

// ReadType returns the resource table for the given resource type.
//
// If the requested resource type does not exist, it returns a nil table.
func (r *Reader) ReadType(id resourcetype.ID) (Table, error) {
	layout, data, err := r.readTableData(0)
	if err != nil {
		return nil, err
	}

	for i := layout.Named; i < layout.Total(); i++ {
		start := i * tableEntrySize
		end := start + tableEntrySize
		entry := numberedEntry(data[start:end])
		if entry.Number() == uint32(id) {
			if ref := entry.Ref(); ref.IsTable() {
				return r.ReadTable(ref.Table())
			}
		}
	}

	return nil, nil
}

// ReadString returns the resource string with the given offset from the
// resource directory.
func (r *Reader) ReadString(offset StringOffset) (string, error) {
	// There is a two byte unsigned integer before the string that indicates
	// its length.
	var buf [2]byte
	if _, err := r.source.ReadAt(buf[:], int64(offset)); err != nil {
		return "", fmt.Errorf("failed to read resource directory string header at offset %d: %w", offset, err)
	}

	// The length is encoded as a little endian value.
	length := binary.LittleEndian.Uint16(buf[:])

	// Allocate an array of appropriate size and read the data into it.
	data := make([]byte, length)
	offset += 2
	if _, err := r.source.ReadAt(data, int64(offset)); err != nil {
		return "", fmt.Errorf("failed to read resource directory string data at offset %d: %w", offset, err)
	}

	// Convert the UTF-16 unicode data into a string.
	return bytesconv.DecodeUTF16(data, binary.LittleEndian), nil
}

// ReadData returns the resource data with the given offset from the resource
// directory.
func (r *Reader) ReadData(offset DataOffset) ([]byte, error) {
	// Read the data descriptor, which tells us where the actual data is
	// stored.
	var buf [dataDescriptorSize]byte
	if _, err := r.source.ReadAt(buf[:], int64(offset)); err != nil {
		return nil, fmt.Errorf("failed to read resource directory data descriptor at offset %d: %w", offset, err)
	}

	// Extract a relative virtual address range from the descriptor.
	descriptor := dataDescriptor(buf[:])
	addressRange := imagefile.RelativeVirtualAddressRange{
		Start:  descriptor.Address(),
		Length: uint(descriptor.Size()),
	}

	// Translate the relative virtual address range into a file range.
	ok, location := r.pe.Sections().TranslateRange(addressRange)
	if !ok {
		return nil, fmt.Errorf("the resource data descriptor at offset \"%d\" has a virtual address range (%s) that is not mapped to any section within the image file", offset, addressRange)
	}

	// Read the data from the file.
	data, err := r.pe.ReadRange(location)
	if err != nil {
		return nil, fmt.Errorf("failed to read resource directory data at location %s: %w", location, err)
	}

	return data, nil
}

// readTableData reads the table header and entry data for the table at the
// given offset. It returns the entry data without interpretation.
func (r *Reader) readTableData(offset TableOffset) (layout tableLayout, data []byte, err error) {
	var buf [tableHeaderSize]byte
	if _, err := r.source.ReadAt(buf[:], int64(offset)); err != nil {
		return layout, nil, fmt.Errorf("failed to read resource directory table header at offset %d: %w", offset, err)
	}

	header := tableHeader(buf[:])
	layout = tableLayout{
		Named:    int(header.NumberOfNamedEntries()),
		Numbered: int(header.NumberOfNumeberedEntries()),
	}

	data = make([]byte, layout.Total()*tableEntrySize)
	offset += tableHeaderSize
	if _, err := r.source.ReadAt(data, int64(offset)); err != nil {
		return layout, nil, fmt.Errorf("failed to read resource directory table entries at offset %d: %w", offset, err)
	}
	return
}
