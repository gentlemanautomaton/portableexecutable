package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gentlemanautomaton/portableexecutable"
	"github.com/gentlemanautomaton/portableexecutable/imagefile"
	"github.com/gentlemanautomaton/portableexecutable/tables/resourcedirectory"
	"github.com/gentlemanautomaton/portableexecutable/tables/resourcedirectory/resourcetype"
	"github.com/gentlemanautomaton/portableexecutable/tables/resourcedirectory/resourcetype/versioninfo"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Please provide the path to an portable executable file (.exe or .dll).\n")
		os.Exit(1)
	}

	path := os.Args[1]
	fmt.Printf("Path: %s\n", path)

	var elapsed time.Duration
	defer func() {
		fmt.Printf("Time Elapsed: %s\n", elapsed)
	}()

	start := time.Now()
	file, err := os.Open(path)
	elapsed += time.Since(start)
	if err != nil {
		fmt.Printf("Failed to open \"%s\": %v\n", path, err)
		os.Exit(1)
	}
	defer file.Close()

	start = time.Now()
	reader, err := portableexecutable.NewReader(file)
	elapsed += time.Since(start)
	if err != nil {
		fmt.Printf("Failed to read information for \"%s\": %v\n", path, err)
		os.Exit(1)
	}

	fmt.Printf("Machine: %s\n", reader.Machine())
	fmt.Printf("Format: %s\n", reader.Format())
	fmt.Printf("Subsystem: %s\n", reader.Subsystem())

	layout := reader.Layout()
	{
		table := layout.SymbolTable()
		fmt.Printf("Symbol Table (%d %s)\n", layout.NumberOfSymbols, plural(layout.NumberOfSymbols, "entry", "entries"))
		fmt.Printf("  Address Range: %s (%d bytes)\n", table, table.Length)
	}
	{
		table := layout.StringTable()
		fmt.Printf("String Table\n")
		fmt.Printf("  Address Range: %s (%d bytes)\n", table, table.Length)
	}
	{
		table := layout.SectionTable()
		fmt.Printf("Section Table (%d %s)\n", layout.NumberOfSections, plural(layout.NumberOfSections, "section", "sections"))
		fmt.Printf("  Address Range: %s (%d bytes)\n", table, table.Length)

		for i, section := range reader.Sections() {
			var name string
			{
				if isReference, offset := section.Name.Reference(); isReference {
					if resolved, err := reader.ReadString(offset); err == nil {
						name = resolved
					}
				} else {
					name = string(section.Name)
				}
			}
			fmt.Printf("  Section %2d: %-16s (Virtual Range: %s, File Range: %s)\n", i, name, section.RelativeVirtualAddressRange, section.FileRange)
		}
	}

	{
		dirs := reader.DataDirectories()
		fmt.Printf("Data Directory Table (%d %s)\n", len(dirs), plural(len(dirs), "entry", "entries"))
		for i, dir := range reader.DataDirectories() {
			id := portableexecutable.DirectoryID(i)
			if dir.IsZero() {
				fmt.Printf("  Data Directory %2d: %-23s\n", i, id)
			} else {
				fmt.Printf("  Data Directory %2d: %-23s (File Range: %s, %d %s)\n", i, id, dir.Location, dir.Location.Length, plural(dir.Location.Length, "byte", "bytes"))
			}
		}

		if resources := dirs.Get(imagefile.ResourceTableID); !resources.IsZero() {
			fmt.Printf("Resource Directory Table\n")
			reader, err := resourcedirectory.NewReader(reader)
			if err != nil {
				fmt.Printf("Failed to prepare a reader for the resource directory: %v\n", err)
				os.Exit(1)
			}
			root, err := reader.ReadRoot()
			if err != nil {
				fmt.Printf("Failed to read the resource directory: %v\n", err)
				os.Exit(1)
			}
			printResourceDirectory(reader, 0, root)

			start = time.Now()
			versions, err := reader.ReadType(resourcetype.Version)
			elapsed += time.Since(start)
			if err != nil {
				fmt.Printf("Failed to read the file version resource: %v\n", err)
				os.Exit(1)
			}
			if len(versions) > 0 {
				fmt.Printf("File Version Information\n")
				for _, version := range versions {
					if version.Reference.IsTable() {
						fmt.Printf("  Version Resource %s\n", version.ID)
						start = time.Now()
						languages, err := reader.ReadTable(version.Reference.Table())
						elapsed += time.Since(start)
						if err != nil {
							fmt.Printf("Failed to read the file version language table: %v\n", err)
							os.Exit(1)
						}
						for _, language := range languages {
							if !language.Reference.IsTable() {
								fmt.Printf("    Language %s\n", language.ID)
								start = time.Now()
								data, err := reader.ReadData(language.Reference.Data())
								elapsed += time.Since(start)
								if err != nil {
									fmt.Printf("      Failed to read the file version data: %v\n", err)
									os.Exit(1)
								}
								fmt.Printf("      Version Data: %d bytes\n", len(data))
								printVersionInfo(data)
							}
						}
					}
				}
			}
		}
	}
}

func printResourceDirectory(reader *resourcedirectory.Reader, depth int, table resourcedirectory.Table) {
	indent := strings.Repeat("  ", depth+1)
	for _, entry := range table {
		if depth == 0 && entry.ID.IsNumeric() {
			fmt.Printf("%s%s\n", indent, resourcetype.ID(entry.ID.Number()))
		} else {
			fmt.Printf("%s%s\n", indent, entry.ID)
		}
		if entry.Reference.IsTable() {
			next, err := reader.ReadTable(entry.Reference.Table())
			if err != nil {
				fmt.Printf("%s  Error: %v\n", indent, err)
				continue
			}
			printResourceDirectory(reader, depth+1, next)
		}
	}
}

func printVersionInfo(data []byte) {
	root, err := versioninfo.NewRoot(data)
	if err != nil {
		fmt.Printf("      Failed to collect file version information: %v\n", err)
	}
	info := root.FileInfo()
	if info.Valid() {
		fmt.Printf("      File Version: %s\n", info.FileVersion())
		fmt.Printf("      Product Version: %s\n", info.ProductVersion())
	} else {
		fmt.Printf("      Invalid File Info\n")
	}
	fmt.Printf("      %s\n", root.Key())
	for node, err := range root.Children() {
		if err != nil {
			fmt.Printf("      Failed to read version info: %v\n", err)
			break
		}
		printVersionInfoTree(2, node)
	}
}

func printVersionInfoTree(depth int, node versioninfo.Node) {
	indent := strings.Repeat("  ", depth+1)
	value := node.Value().String()
	if value != "" {
		fmt.Printf("  %s%s: %s\n", indent, node.Key(), value)
	} else {
		fmt.Printf("  %s%s\n", indent, node.Key())
	}
	for node, err := range node.Children() {
		if err != nil {
			fmt.Printf("%sFailed to read version info: %v\n", indent, err)
			break
		}
		printVersionInfoTree(depth+1, node)
	}
}
