package imagefile

import "fmt"

// StringOffset is the offset of a string within the strings table.
type StringOffset uint

// FileOffset is an absolute offset within an image file.
type FileOffset uint

// String returns a hexadecimal representation of the file offset.
func (offset FileOffset) String() string {
	return fmt.Sprintf("0x%x", uint(offset))
}

// FileRange describes an address range within an executable image file.
type FileRange = Range[FileOffset]

// VirtualAddress is a virtual address within an image file. It is an absolute
// address within the virtual memory space.
//type VirtualAddress uint

// RelativeVirtualAddress is an address within the virtual address space of a
// mapped image file. It is relative to the image's base address.
type RelativeVirtualAddress uint

// String returns a hexadecimal representation of the address.
func (rva RelativeVirtualAddress) String() string {
	return fmt.Sprintf("0x%x", uint(rva))
}

// RelativeVirtualAddressRange describes an address range within virtual memory.
type RelativeVirtualAddressRange = Range[RelativeVirtualAddress]

// Address is a type constraint for any sort of address, either virtual or
// real.
type Address interface {
	~uint
}

// Range describes an address range for addresses of type T.
type Range[T Address] struct {
	Start  T
	Length uint
}

// IsZero returns true if the address range is zero.
func (r Range[T]) IsZero() bool {
	return r.Start == 0 && r.Length == 0
}

// End returns the last address within the range.
func (r Range[T]) End() T {
	if r.Length == 0 {
		return r.Start
	}
	return r.Start + T(r.Length) - 1
}

// Contains returns true if the range contains the given address.
func (r Range[T]) Contains(address T) bool {
	if r.Length == 0 {
		return false
	}
	if address < r.Start {
		return false
	}
	if address >= r.Start+T(r.Length) {
		return false
	}
	return true
}

// ContainsRange returns true if the range contains other.
func (r Range[T]) ContainsRange(other Range[T]) bool {
	if r.Length == 0 {
		return false
	}
	if other.Length == 0 {
		return false
	}
	if other.Start < r.Start {
		return false
	}
	if other.Start+T(other.Length) > r.Start+T(r.Length) {
		return false
	}
	return true
}

// String returns a string representation of the address range.
func (r Range[T]) String() string {
	return fmt.Sprintf("0x%x-0x%x", uint(r.Start), uint(r.End()))
}
