package resourcetype

import (
	"fmt"
	"strconv"
)

// ID is a resource type identifier.
type ID uint32

// String returns a string representation of the resource type ID.
//
// If there is a predefined name for the ID, it is included in parentheses.
func (id ID) String() string {
	name := ""
	if int(id) < len(predefinedNames) {
		name = predefinedNames[id]
	}
	if name == "" {
		return strconv.FormatUint(uint64(id), 10)
	}
	return fmt.Sprintf("%d (%s)", id, name)
}
