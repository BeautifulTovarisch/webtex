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
	"fmt"
	"os"
	"text/template"
)

const tmplPath = "templates/doc.tmpl"

//go:embed templates/doc.tmpl
var docTemplateFile embed.FS
var docTemplate *template.Template

func init() {
	docTemplate = template.Must(template.New("doc").ParseFS(docTemplateFile, tmplPath))
}

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
func HTMLDoc(doc Document) {
	if err := docTemplate.ExecuteTemplate(os.Stdout, "doc.tmpl", doc); err != nil {
		fmt.Println(err)
	}
}
