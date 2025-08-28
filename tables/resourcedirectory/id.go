package resourcedirectory

import "strconv"

// ID stores a resource identifier, which can be either a string or a 31-bit
// unsigned integer.
type ID struct {
	str    string
	number uint32
}

// IsNumeric returns true if the ID is numeric.
func (id ID) IsNumeric() bool {
	return id.str == ""
}

// Number returns the value of the ID if it is numeric.
func (id ID) Number() uint32 {
	return id.number
}

// String returns a string representation of the ID. If the ID is numeric
// it converts the number to a string and returns it.
func (id ID) String() string {
	if id.str != "" {
		return id.str
	}
	return strconv.FormatUint(uint64(id.number), 10)
}
