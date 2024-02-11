// package mdrender renders Markdown documents into HTML
package mdrender

import (
	"unsafe"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const (
	htmlFlags  = html.CommonFlags | html.HrefTargetBlank | html.TOC
	extensions = parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
)

// Since we'll be potentially be converting a lot of markdown, we want to avoid
// unnecessary copying.
//
// CAUTION: Mutating the []byte provided from this function will SEGSEV!
func toString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func mdToHtml(md []byte) string {
	p := parser.NewWithExtensions(extensions)
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})

	doc := p.Parse(md)

	return toString(markdown.Render(doc, renderer))
}

// Render converts a markdown snippet into HTML
func Render(md string) string {
	// Parse evidently modifies the []byte provided to it. Can't use our hack :(
	mdBytes := []byte(md)

	return mdToHtml(mdBytes)
}
