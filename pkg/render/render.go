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

func renderBlock(c chunk.Chunk) (string, error) {
	if c.T != chunk.BLOCK {
		panic("Implementation error. Expected LaTeX block")
	}

	return texrender.RenderBlock(c.Content)
}

func renderInline(c chunk.Chunk) (string, error) {
	if c.T != chunk.INLINE {
		panic("Implementation error. Expected inline LaTeX")
	}

	return texrender.RenderInline(c.Content)
}

func processChunk(c chunk.Chunk) (string, error) {
	switch c.T {
	case chunk.MD:
		return renderMd(c), nil
	case chunk.INLINE:
		return renderInline(c)
	case chunk.BLOCK:
		return renderBlock(c)
	}

	return "", nil
}

// RenderDoc accepts a string containing an individual markdown document and
// writes an HTML document with the rendered content of [md] to [out].
func RenderDoc(md string, out io.Writer) error {
	var b strings.Builder

	chunks := chunk.ChunkDoc(md)

	// We can stream the output of processChunk directly to out.
	for _, c := range chunks {
		if svg, err := processChunk(c); err != nil {
			b.WriteString("Error")
		} else {
			b.WriteString(svg)
		}
	}

	return b.String()
}
