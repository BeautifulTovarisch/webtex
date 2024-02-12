// package build provides utilities to help build a website from markdown files
package build

import (
	"io"
	"maps"
	"os"
	"path/filepath"
)

// Nav is an adjacency list of the file organization of markdown files. Entries
// are represented as [os.DirEntry] for convenience.
type Nav map[string][]os.DirEntry

// SiteNav constructs a directed graph representing the navigation links of the
// source files located under [file]. This information is useful for generating
// a global site map and navigation components.
func SiteNav(file string) (Nav, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return Nav{}, err
	}

	if !stat.IsDir() {
		return Nav{}, nil
	}

	files, err := os.ReadDir(file)
	if err != nil {
		return Nav{}, err
	}

	// Not a directory.
	if len(files) == 0 {
		return Nav{}, nil
	}

	nav := make(Nav)

	nav[file] = files
	for _, f := range files {
		tree, err := SiteNav(filepath.Join(file, f.Name()))

		if err != nil {
			return Nav{}, err
		}

		maps.Copy(nav, tree)
	}

	return nav, nil
}

func readSource(src io.Reader) {}

// Build reads the markdown files under the [src] directory and writes HTML to
// the [out] directory.
func Build(src string, out string) {
}
