package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// adapterInfo holds the extracted information about a Hive TypeAdapter class.
type adapterInfo struct {
	adapterClass string
	genericType  string
	typeID       int
	file         string
}

var (
	classRe  = regexp.MustCompile(`class\s+(\w+)\s+extends\s+TypeAdapter<(\w+)>`)
	typeIDRe = regexp.MustCompile(`final\s+typeId\s*=\s*(\d+)`)
)

func scanFile(path string) (*adapterInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var info adapterInfo
	info.file = path
	foundClass := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if !foundClass {
			if m := classRe.FindStringSubmatch(line); m != nil {
				info.adapterClass = m[1]
				info.genericType = m[2]
				foundClass = true
			}
			continue
		}

		// Once we've found the class, look for the typeId within the next lines.
		if m := typeIDRe.FindStringSubmatch(line); m != nil {
			id, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, err
			}
			info.typeID = id
			return &info, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Found class declaration but no typeId — not a hive adapter we can parse.
	return nil, nil
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var adapters []adapterInfo

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".dart" {
			return nil
		}

		result, err := scanFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not parse %s: %v\n", path, err)
			return nil
		}
		if result != nil {
			adapters = append(adapters, *result)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error walking directory: %v\n", err)
		os.Exit(1)
	}

	if len(adapters) == 0 {
		fmt.Println("No Hive TypeAdapter classes found.")
		return
	}

	// Sort by typeId for a predictable, readable output.
	sort.Slice(adapters, func(i, j int) bool {
		return adapters[i].typeID < adapters[j].typeID
	})

	// relPath returns the path relative to root, falling back to the absolute path.
	relPath := func(abs string) string {
		rel, err := filepath.Rel(root, abs)
		if err != nil {
			return abs
		}
		return rel
	}

	// Compute column widths.
	colAdapter := len("Adapter Class")
	colType := len("Generic Type")
	colID := len("Type ID")
	colFile := len("File")

	for _, a := range adapters {
		if l := len(a.adapterClass); l > colAdapter {
			colAdapter = l
		}
		if l := len(a.genericType); l > colType {
			colType = l
		}
		if l := len(relPath(a.file)); l > colFile {
			colFile = l
		}
	}

	// Build separator and header.
	sep := fmt.Sprintf("+-%s-+-%s-+-%s-+-%s-+",
		strings.Repeat("-", colAdapter),
		strings.Repeat("-", colType),
		strings.Repeat("-", colID),
		strings.Repeat("-", colFile),
	)

	row := func(adapter, typ, id, file string) string {
		return fmt.Sprintf("| %-*s | %-*s | %-*s | %-*s |",
			colAdapter, adapter,
			colType, typ,
			colID, id,
			colFile, file,
		)
	}

	fmt.Println(sep)
	fmt.Println(row("Adapter Class", "Generic Type", "Type ID", "File"))
	fmt.Println(sep)

	for _, a := range adapters {
		fmt.Println(row(a.adapterClass, a.genericType, strconv.Itoa(a.typeID), relPath(a.file)))
	}

	fmt.Println(sep)
	fmt.Printf("\n%d adapter(s) found.\n", len(adapters))
}
