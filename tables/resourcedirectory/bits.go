package resourcedirectory

const (
	// controlBit32 is the uppermost bit in a 32-bit value. It is sometimes
	// used as a control bit to indicate how the lower 32 bits within the
	// value should be interpreted.
	controlBit32 = uint32(1) << 31

	// valueMask31 is a mask with the first 31 bits set. It is used to exclude
	// bit 32, which is used as a control bit for some values.
	valueMask31 = ^controlBit32
)
