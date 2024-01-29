// package render accepts a Markdown document potentially containing LaTeX code
// and renders the components into HTML.
//
// render does not consider the overall structure of the document, instead only
// rendering the snippets of Markdown and LaTeX into corresponding HTML code
package render

// TODO: Our chunking strategy potentially allows documents to be streamed.
// Consider an option to allow documents to be provided to STDIN directly and
// processed in a pipeline.

import (
	"fmt"
	"strings"

	"github.com/beautifultovarisch/webtex/internal/chunk"
	"github.com/beautifultovarisch/webtex/internal/mdrender"
	"github.com/beautifultovarisch/webtex/internal/texrender"
)

func renderMd(c chunk.Chunk) string {
	if c.T != chunk.MD {
		panic("Implementation error. Expected markdown chunk")
	}

	return mdrender.Render(c.Content)
}

func renderTex(c chunk.Chunk) string {
	if c.T == chunk.MD {
		panic("Implementation error. Expected LaTeX chunk")
	}

	return texrender.Render(c.Content)
}

func assembleDoc(chunks []string) string {
	var b strings.Builder
	for _, c := range chunks {
		b.WriteString(c)
	}

	return b.String()
}

// RenderDoc accepts a string containing an individual markdown document and
// returns an HTML document with the rendered content of [md].
func RenderDoc(md string) string {
	blocks := chunk.ChunkDoc(md)

	fmt.Println(blocks)

	return ""
}
