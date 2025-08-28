package portableexecutable

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gentlemanautomaton/portableexecutable/dos"
	"github.com/gentlemanautomaton/portableexecutable/imagefile"
)

// Reader reads portable executable image file data from an underlying
// [io.ReaderAt].
type Reader struct {
	source io.ReaderAt

	layout      imagefile.Layout
	machine     imagefile.Machine
	format      imagefile.Format
	subsystem   imagefile.Subsystem
	sections    SectionTable
	directories DataDirectoryTable
}

// NewReader creates and initializes a new portable executable image file
// [Reader] that reads from source.
//
// It reads and validates a small amount of data in order to make future
// accesses quicker. If the validation fails it returns an error.
func NewReader(source io.ReaderAt) (*Reader, error) {
	reader := new(Reader)
	if err := reader.init(source); err != nil {
		return nil, err
	}
	return reader, nil
}

// Source returns the underlying source for the reader.
func (r *Reader) Source() io.ReaderAt {
	return r.source
}

// Machine returns the machine that the image file is targeting.
func (r *Reader) Machine() imagefile.Machine {
	return r.machine
}

// Format returns the format of the image file's header.
func (r *Reader) Format() imagefile.Format {
	return r.format
}

// Subsystem returns the subsystem that is responsible for executing the
// image.
func (r *Reader) Subsystem() imagefile.Subsystem {
	return r.subsystem
}

// Layout returns information about the layout of the image file.
func (r *Reader) Layout() imagefile.Layout {
	return r.layout
}

// Sections returns a table of sections that are present within the image
// file.
func (r *Reader) Sections() SectionTable {
	return r.sections
}

// DataDirectories returns the table of data directories for the the image
// file.
func (r *Reader) DataDirectories() DataDirectoryTable {
	return r.directories
}

// ReadRange reads data from the image file for the given file range.
func (r *Reader) ReadRange(fileRange imagefile.FileRange) ([]byte, error) {
	data := make([]byte, fileRange.Length)
	_, err := r.source.ReadAt(data, int64(fileRange.Start))
	return data, err
}

// ReadString returns a string from the image file's string table with the
// given offset. If the string is longer than 4096 bytes it will be
// truncated.
//
// TODO: Add support for longer strings and consider using some sort of
// string caching.
func (r *Reader) ReadString(offset imagefile.StringOffset) (string, error) {
	table := r.layout.StringTable()
	if uint(offset) >= table.Length {
		return "", fmt.Errorf("the string offset %d exceeds the %d byte length of the COFF string table", offset, table.Length)
	}

	// Unfortunately, the strings are null-terminated so we don't know how
	// much data we need to read ahead of time. For now, we use an arbitrary
	// limit of 4096 bytes, which is probably way more than is necessary for
	// typical uses of the string table.
	length := min(table.Length-uint(offset), 4096)

	data, err := r.ReadRange(imagefile.FileRange{
		Start:  table.Start + imagefile.FileOffset(offset),
		Length: length,
	})
	if err != nil {
		return "", err
	}

	// If there is a null terminator, use that as the cutoff.
	if cutoff := bytes.IndexByte(data, 0); cutoff >= 0 {
		return string(data[0:cutoff]), err
	}

	return string(data), nil
}

func (r *Reader) init(source io.ReaderAt) error {
	r.source = source

	// Read and verify the DOS compatibility header, and also get the address
	// of the PE header.
	{
		var dosHeader dos.Header
		if _, err := source.ReadAt(dosHeader[:], 0); err != nil {
			return fmt.Errorf("failed the read the DOS compatibility file header: %w", err)
		}
		if signature := dosHeader.Signature(); !signature.Valid() {
			return fmt.Errorf("the file does not have the expected DOS compatibility file header signature")
		}
		r.layout.Start = dosHeader.NextHeader()
	}

	// Read the PE signature and COFF file header.
	{
		data, err := r.ReadRange(r.layout.SignatureAndFileHeader())
		if err != nil {
			return fmt.Errorf("failed the read the portable executable file signature and COFF file header: %w", err)
		}

		// Verify the signature.
		signature := imagefile.Signature(data[0:imagefile.SignatureSize])
		if !signature.Valid() {
			return fmt.Errorf("the file does not have the expected portable executable signature")
		}

		// Extract data from the file header.
		fileHeader := imagefile.FileHeader(data[imagefile.SignatureSize:])
		r.machine = fileHeader.Machine()
		r.layout.NumberOfSections = uint(fileHeader.NumberOfSections())
		r.layout.SizeOfOptionalHeader = uint(fileHeader.SizeOfOptionalHeader())
		r.layout.StartOfSymbolTable = fileHeader.PointerToSymbolTable()
		r.layout.NumberOfSymbols = uint(fileHeader.NumberOfSymbols())
	}

	// Read the optional header and collect information about the data
	// directories.
	var dataDirs []imagefile.DataDirectory
	{
		if r.layout.SizeOfOptionalHeader < 2 {
			return fmt.Errorf("the portable executable file has an optional header section of %d byte(s), which is an insufficient size", r.layout.SizeOfOptionalHeader)
		}
		data, err := r.ReadRange(r.layout.OptionalHeader())
		if err != nil {
			return fmt.Errorf("failed the read the optional header section of the portable executable: %w", err)
		}

		r.format = imagefile.Format(binary.LittleEndian.Uint16(data[0:2]))

		var optionalHeader imagefile.OptionalHeader
		switch r.format {
		case imagefile.PE32:
			if len(data) < imagefile.MinOptionalHeaderSize32 {
				return fmt.Errorf("the portable executable file has an optional header section of %d byte(s), which is less than the mininium of %d bytes for 32-bit executables", r.layout.SizeOfOptionalHeader, imagefile.MinOptionalHeaderSize32)
			}
			optionalHeader = imagefile.OptionalHeader32(data)
		case imagefile.PE32Plus:
			if len(data) < imagefile.MinOptionalHeaderSize64 {
				return fmt.Errorf("the portable executable file has an optional header section of %d byte(s), which is less than the minimum size of %d bytes for 64-bit executables", r.layout.SizeOfOptionalHeader, imagefile.MinOptionalHeaderSize64)
			}
			optionalHeader = imagefile.OptionalHeader64(data)
		default:
			return fmt.Errorf("the optional header section of the portable executable has an unsupported format: %s", r.format)
		}

		r.subsystem = optionalHeader.Subsystem()
		dataDirs = optionalHeader.DataDirectories()
	}

	// Read the section table.
	{
		data, err := r.ReadRange(r.layout.SectionTable())
		if err != nil {
			return fmt.Errorf("failed the read the section table of the portable executable: %w", err)
		}
		count := len(data) / imagefile.SectionHeaderSize
		r.sections = make(SectionTable, 0, count)
		for i := range count {
			start := i * imagefile.SectionHeaderSize
			end := start + imagefile.SectionHeaderSize
			header := imagefile.SectionHeader(data[start:end])
			r.sections = append(r.sections, Section{
				Name: header.Name(),
				RelativeVirtualAddressRange: imagefile.RelativeVirtualAddressRange{
					Start:  header.VirtualAddress(),
					Length: header.VirtualSize(),
				},
				FileRange: imagefile.FileRange{
					Start:  header.PointerToRawData(),
					Length: header.SizeOfRawData(),
				},
			})
		}
	}

	// Process the data directories at the end of the optional header.
	//
	// As we process each one, translate virtual addresses to file offsets.
	{
		entries := make([]DataDirectory, 0, len(dataDirs))
		for i, entry := range dataDirs {
			id := imagefile.DirectoryID(i)

			// A zeroed entry indicates that the data directory is absent.
			// It's important that we still include it because data directory
			// IDs are based on well-known indices.
			if entry.IsZero() {
				entries = append(entries, DataDirectory{})
				continue
			}

			switch {
			case id.IsPointer():
				address := imagefile.RelativeVirtualAddress(entry.Address())
				ok, offset := r.sections.Translate(address)
				if !ok {
					return fmt.Errorf("data directory \"%s\" has a virtual address (%s) that is not mapped to any section within the image file", id, address)
				}
				entries = append(entries, DataDirectory{
					Location: imagefile.FileRange{
						Start: offset,
					},
				})
			case id.IsVirtual():
				addressRange := imagefile.RelativeVirtualAddressRange{
					Start:  imagefile.RelativeVirtualAddress(entry.Address()),
					Length: uint(entry.Size()),
				}
				ok, location := r.sections.TranslateRange(addressRange)
				if !ok {
					return fmt.Errorf("data directory \"%s\" has a virtual address range (%s) that is not mapped to any section within the image file", id, addressRange)
				}
				entries = append(entries, DataDirectory{
					Location: location,
				})
			default:
				entries = append(entries, DataDirectory{
					Location: imagefile.FileRange{
						Start:  imagefile.FileOffset(entry.Address()),
						Length: uint(entry.Size()),
					},
				})
			}
		}

		r.directories = DataDirectoryTable(entries)
	}

	// Read the string table header.
	if r.layout.StartOfSymbolTable != 0 {
		stringTable := r.layout.StringTable()
		stringTable.Length = 4
		data, err := r.ReadRange(stringTable)
		if err != nil {
			return fmt.Errorf("failed the read the string table header of the portable executable: %w", err)
		}
		r.layout.SizeOfStringTable = uint(binary.LittleEndian.Uint32(data))
	}

	return nil
}
