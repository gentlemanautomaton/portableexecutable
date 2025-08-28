package portableexecutable

import "github.com/gentlemanautomaton/portableexecutable/imagefile"

// SectionTable holds a set of section definitions for a portable executable
// image file.
type SectionTable []Section

// Translate maps the given relative virtual address to a file offset within
// the image file that contains the corresponding data.
//
// If the given address is not contained within the virtual address range of
// any of the sections within the table, it returns false.
func (table SectionTable) Translate(address imagefile.RelativeVirtualAddress) (ok bool, offset imagefile.FileOffset) {
	for i := range table {
		if ok, offset = table[i].Translate(address); ok {
			return
		}
	}
	return
}

// TranslateRange maps the given relative virtual address range to a file
// offset range within the image file that contains the backing data.
//
// If the given range is not contained within the virtual address range of
// any of the sections within the table, it returns false.
func (table SectionTable) TranslateRange(addressRange imagefile.RelativeVirtualAddressRange) (ok bool, fileRange imagefile.FileRange) {
	for i := range table {
		if ok, fileRange = table[i].TranslateRange(addressRange); ok {
			return
		}
	}
	return
}

// Section describes a section within a portable executable image file.
type Section struct {
	Name                        imagefile.SectionName
	RelativeVirtualAddressRange imagefile.RelativeVirtualAddressRange
	FileRange                   imagefile.FileRange
}

// Translate maps the given relative virtual address to a file offset within
// the image file that contains the backing data.
//
// If the given address is not contained within the virtual address range of
// the section, it returns false.
func (section Section) Translate(address imagefile.RelativeVirtualAddress) (ok bool, offset imagefile.FileOffset) {
	if !section.RelativeVirtualAddressRange.Contains(address) {
		return false, 0
	}

	// The requested address and the section's starting address are both
	// relative to the same base address in the virtual memory space.
	// This means we can ignore the base address completely and just take
	// the difference between them. The calculation becomes very simple:
	//
	//   VirtualAddress - VirtualStart + FileStart
	//
	// Note that we're working with unsigned integers here, so we *must*
	// perform the addition first before we perform the subtraction. This
	// avoids a potential negative result that would wrap the unsigned
	// integer around to some huge positive value.
	return true, imagefile.FileOffset(address) + section.FileRange.Start - imagefile.FileOffset(section.RelativeVirtualAddressRange.Start)
}

// TranslateRange maps the given relative virtual address range to a file
// offset range within the image file that contains the backing data.
//
// If the given range is not contained within the virtual address range of
// the section, it returns false.
func (section Section) TranslateRange(addressRange imagefile.RelativeVirtualAddressRange) (ok bool, fileRange imagefile.FileRange) {
	if !section.RelativeVirtualAddressRange.ContainsRange(addressRange) {
		return false, imagefile.FileRange{}
	}
	return true, imagefile.FileRange{
		Start:  imagefile.FileOffset(addressRange.Start) + section.FileRange.Start - imagefile.FileOffset(section.RelativeVirtualAddressRange.Start),
		Length: addressRange.Length,
	}
}
