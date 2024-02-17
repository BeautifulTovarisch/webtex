// package render accepts a Markdown document potentially containing LaTeX code
// and renders the components into HTML.
//
// render does not consider the overall structure of the document, instead only
// rendering the snippets of Markdown and LaTeX into corresponding HTML code
package render

import (
	"bufio"
	"io"

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
func RenderDoc(md io.Reader, out io.Writer) error {
	buf := bufio.NewReader(md)

	// We can stream the output of processChunk directly to out.
	for {
		c, err := chunk.ChunkDoc(buf)
		if err != nil {
			return err
		}

		html, err := processChunk(c)
		if err != nil {
			return err
		}

		io.WriteString(out, html)
	}
}
