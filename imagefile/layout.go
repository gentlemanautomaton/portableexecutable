package imagefile

// Image file structure sizes.
const (
	SignatureSize     = 4
	FileHeaderSize    = 20
	DataDirectorySize = 8
	SectionHeaderSize = 40
	SymbolSize        = 18
)

// Layout keeps track of data used to determine the layout of an image file.
type Layout struct {
	Start                FileOffset
	NumberOfSections     uint
	SizeOfOptionalHeader uint
	StartOfSymbolTable   FileOffset
	NumberOfSymbols      uint
	SizeOfStringTable    uint
}

// Signature returns the address range of the PE signature.
func (layout Layout) Signature() FileRange {
	return FileRange{Start: layout.Start, Length: SignatureSize}
}

// FileHeader returns the address range of the COFF file header.
func (layout Layout) FileHeader() FileRange {
	return FileRange{Start: layout.Start + SignatureSize, Length: FileHeaderSize}
}

// SignatureAndFileHeader returns an address range that includes the PE
// signature and COFF file header.
func (layout Layout) SignatureAndFileHeader() FileRange {
	return FileRange{Start: layout.Start, Length: SignatureSize + FileHeaderSize}
}

// OptionalHeader returns the address range of the COFF optional header.
func (layout Layout) OptionalHeader() FileRange {
	return FileRange{Start: layout.Start + SignatureSize + FileHeaderSize, Length: layout.SizeOfOptionalHeader}
}

// SymbolTable returns the address range of the COFF symbol table.
func (layout Layout) SymbolTable() FileRange {
	return FileRange{Start: layout.StartOfSymbolTable, Length: layout.NumberOfSymbols * SymbolSize}
}

// SectionTable returns the address range of the COFF section table.
func (layout Layout) SectionTable() FileRange {
	return FileRange{Start: layout.Start + SignatureSize + FileHeaderSize + FileOffset(layout.SizeOfOptionalHeader), Length: layout.NumberOfSections * SectionHeaderSize}
}

// StringTable returns the address range of the COFF string table.
func (layout Layout) StringTable() FileRange {
	return FileRange{Start: layout.StartOfSymbolTable + FileOffset(layout.NumberOfSymbols)*SymbolSize, Length: layout.SizeOfStringTable}
}
