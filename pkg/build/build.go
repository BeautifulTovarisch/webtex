// package build provides utilities to help build a website from markdown files
package build

import (
	"io/fs"
	"maps"
	"os"
	"path/filepath"

	"github.com/beautifultovarisch/webtex/pkg/render"

	"github.com/beautifultovarisch/webtex/internal/sitebuilder"
)

// Nav is an adjacency list of the file organization of markdown files. Entries
// are represented as [os.DirEntry] for convenience.
type Nav map[string][]os.DirEntry

// SiteNav constructs a directed graph representing the navigation links of the
// source files located under [file]. This information is useful for generating
// a global site map and navigation components.
//
// TODO: Optimize by eliminating the call to os.Stat and traversing on DirEntry
// types instead.
func SiteNav(file string) (Nav, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return Nav{}, err
	}

	// Base case
	if !stat.IsDir() {
		return Nav{}, nil
	}

	files, err := os.ReadDir(file)
	if err != nil {
		return Nav{}, err
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

// Build reads the markdown files under the [src] directory and writes HTML to
// the [out] directory.
func Build(src string, dst string) error {
	_, err := SiteNav(src)
	if err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		doc := sitebuilder.Document{
			Title:   d.Name(),
			Content: render.RenderDoc(string(content)),
		}

		sitebuilder.HTMLDoc(os.Stdout, doc)

		return nil
	})
}
