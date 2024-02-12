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
	"github.com/beautifultovarisch/webtex/internal/logger"
	"github.com/beautifultovarisch/webtex/internal/mdrender"
	"github.com/beautifultovarisch/webtex/internal/texrender"
)

func renderMd(c chunk.Chunk) string {
	if c.T != chunk.MD {
		panic("Implementation error. Expected markdown chunk")
	}

	return mdrender.Render(c.Content)
}

func renderTex(c chunk.Chunk) (string, error) {
	if c.T == chunk.MD {
		panic("Implementation error. Expected LaTeX chunk")
	}

	return texrender.Render(c.Content)
}

// RenderDoc accepts a string containing an individual markdown document and
// returns an HTML document with the rendered content of [md].
func RenderDoc(md string) string {
	var doc strings.Builder

	chunks := chunk.ChunkDoc(md)

	for _, c := range chunks {
		if c.T == chunk.MD {
			doc.WriteString(renderMd(c))
		} else {
			svg, err := renderTex(c)
			if err != nil {
				logger.Error("Error rendering TeX: %s", err)

				continue
			}

			doc.WriteString(svg)
		}
	}

	return doc.String()
}
