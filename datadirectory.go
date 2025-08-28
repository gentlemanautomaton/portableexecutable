package portableexecutable

import (
	"github.com/gentlemanautomaton/portableexecutable/imagefile"
)

// DirectoryID identifies a data directory by its index within the directory
// table.
type DirectoryID = imagefile.DirectoryID

// DataDirectoryTable holds the set of data directory definitions for a
// portable executable file.
type DataDirectoryTable []DataDirectory

// Get returns information about the data directory with the given ID.
// If the directory is not present it returns a zeroed data directory.
func (table DataDirectoryTable) Get(id DirectoryID) DataDirectory {
	// Note that the table is sparse, and it's very common for an ID to be
	// present in the table but to have zeroed data, which indicates that
	// the table is missing.
	//
	// Returning a zeroed value here makes things simple for the caller,
	// because a zeroed structure indicates a missing directory in all cases.
	if int(id) >= len(table) {
		return DataDirectory{}
	}
	return table[id]
}

// DataDirectory describes the location of a data directory within an image
// file.
type DataDirectory struct {
	Location imagefile.FileRange
}

// IsZero returns true if the data directory is not present within the image
// file.
func (dir DataDirectory) IsZero() bool {
	return dir.Location.IsZero()
}
