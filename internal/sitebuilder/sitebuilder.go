// package sitebuilder provides utilities for constructing HTML documents given
// rendered markdown and LaTeX (SVGs). The package is responsible for:
//
//   - Constructing well-formed HTML documents
//   - Including any static dependencies (CSS, JavaScript, etc.)
//
// TODO: Determine how to process Obsidian backlinks into refs.
package sitebuilder

import (
	"embed"
	"io"
	"text/template"

	"github.com/beautifultovarisch/webtex/internal/logger"
)

const tmplPath = "templates/doc.tmpl"

var (
	//go:embed templates/doc.tmpl
	docTemplateFile embed.FS
	docTemplate     *template.Template
)

func init() {
	docTemplate = template.Must(template.New("doc").ParseFS(docTemplateFile, tmplPath))
}

// Href represents a navigation link in a web document.
type Href struct {
	Ref     string // Ref is the URI or local reference to the target resource (e.g #heading)
	Display string // Display is the human readable text displayed to represent the underlying Ref
}

// Document encapsulates the metadata required to render a document.
type Document struct {
	Title      string // Title is the title of the document (for use in a <title> tag).
	Content    string // Content is the main content of the page.
	Navigation []Href // Navigation is a list of outgoing links from the current document
}

// HTMLDoc produces a complete HTML document with [content] as its body. The
// [content] is escaped using Go's html templating.
func HTMLDoc(out io.Writer, doc Document) error {
	if err := docTemplate.ExecuteTemplate(out, "doc.tmpl", doc); err != nil {
		logger.Error("Error rendering HTML: %s", err)

		return err
	}

	return nil
}
