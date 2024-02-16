// package build provides utilities to help build a website from markdown files
package build

import (
	// "io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"strings"

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

// Strip the source file of its extension and produce a path with the desired
// output extension.
func outputPath(path string) string {
	// Should only ever be one
	return strings.Replace(path, ".md", ".html", 1)
}

func processDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && path != src {
			// Mirror the directory structure of the source files.
			if err := os.MkdirAll(filepath.Join(dst, path), os.ModePerm); err != nil {
				return err
			}

			return nil
		}

		fi, err := d.Info()
		if err != nil {
			return err
		}

		if fi.Mode().IsRegular() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			doc := sitebuilder.Document{
				Title:   d.Name(),
				Content: render.RenderDoc(string(content)),
			}

			// Output file.
			file, err := os.Create(filepath.Join(dst, outputPath(path)))
			defer file.Close()

			if err != nil {
				return err
			}

			if err := sitebuilder.HTMLDoc(file, doc); err != nil {
				return err
			}
		}

		return nil
	})
}

// Build reads the markdown files under the [src] directory and writes HTML to
// the [dst] directory.
func Build(src string, dst string) error {
	_, err := SiteNav(src)
	if err != nil {
		return err
	}

	// Create output directory
	if err := os.MkdirAll(filepath.Join(dst, src), os.ModePerm); err != nil {
		return err
	}

	if err := processDir(src, dst); err != nil {
		return err
	}

	return nil
}
