package resourcedirectory

// Reference is a reference to a either a leaf containing data or to another
// resource directory table that is its subdirectory.
type Reference struct {
	data uint32
}

// IsTable returns true if the reference refers to a table.
func (ref Reference) IsTable() bool {
	return ref.data&controlBit32 != 0
}

// Table returns the resource offset of the subdirectory the reference is
// pointing to if the reference is not a leaf. If the reference is a leaf,
// it returns 0.
func (ref Reference) Table() TableOffset {
	if ref.data&controlBit32 == 0 {
		return 0
	}
	return TableOffset(ref.data & valueMask31)
}

// Data returns the resource offset of the data the reference is pointing to
// if the reference is a leaf. If the reference is not a leaf, it returns 0.
func (ref Reference) Data() DataOffset {
	if ref.data&controlBit32 != 0 {
		return 0
	}
	return DataOffset(ref.data & valueMask31)
}
