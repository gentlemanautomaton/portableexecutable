package resourcetype

// Predefined resource types.
//
// https://learn.microsoft.com/en-us/windows/win32/menurc/resource-types
const (
	Cursor       = 1  // RT_CURSOR, Hardware-dependent cursor resource.
	Bitmap       = 2  // RT_BITMAP, Bitmap resource.
	Icon         = 3  // RT_ICON, Hardware-dependent icon resource.
	Menu         = 4  // RT_MENU, Menu resource.
	Dialog       = 5  // RT_DIALOG, Dialog box.
	String       = 6  // RT_STRING, String-table entry.
	FontDir      = 7  // RT_FONTDIR, Font directory resource.
	Font         = 8  // RT_FONT, Font resource.
	Accelerator  = 9  // RT_ACCELERATOR, Accelerator table.
	RCDATA       = 10 // RT_RCDATA, Application-defined resource (raw data).
	MessageTable = 11 // RT_MESSAGETABLE, Message-table entry.
	CursorGroup  = 12 // RT_GROUP_CURSOR, Hardware-independent cursor resource.
	IconGroup    = 14 // RT_GROUP_ICON, Hardware-independent icon resource.
	Version      = 16 // RT_VERSION, Version resource.
	DlgInclude   = 17 // RT_DLGINCLUDE, Allows a resource editing tool to associate a string with an .rc file.
	PlugAndPlay  = 19 // RT_PLUGPLAY, Plug and Play resource.
	VXD          = 20 // RT_VXD, VXD.
	AniCursor    = 21 // RT_ANICURSOR, Animated cursor.
	AniIcon      = 22 // RT_ANIICON, Animated icon.
	HTML         = 23 // RT_HTML, HTML resource.
	Manifest     = 24 // RT_MANIFEST, Side-by-Side Assembly Manifest.
)

var predefinedNames = []string{
	"",
	"Cursor",
	"Bitmap",
	"Icon",
	"Menu",
	"Dialog",
	"String",
	"FontDir",
	"Font",
	"Accelerator",
	"RCDATA",
	"MessageTable",
	"CursorGroup",
	"",
	"IconGroup",
	"",
	"Version",
	"DlgInclude",
	"",
	"PlugAndPlay",
	"VXD",
	"AniCursor",
	"AniIcon",
	"HTML",
	"Manifest",
}
