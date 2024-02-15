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
	"sync"

	"github.com/beautifultovarisch/webtex/internal/chunk"
	"github.com/beautifultovarisch/webtex/internal/mdrender"
	"github.com/beautifultovarisch/webtex/internal/texrender"
)

const MAX_ROUTINES = 10

var wg sync.WaitGroup

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

func processChunk(idx int, c chunk.Chunk, out []string) {
	defer wg.Done()

	switch c.T {
	case chunk.MD:
		out[idx] = renderMd(c)
	case chunk.INLINE:
		return
	case chunk.BLOCK:
		svg, err := renderBlock(c)
		if err != nil {
			return
		}

		out[idx] = svg
	}

	return
}

func assembleDoc(out []string) string {
	var b strings.Builder

	for _, s := range out {
		b.WriteString(s)
	}

	return b.String()
}

// RenderDoc accepts a string containing an individual markdown document and
// returns an HTML document with the rendered content of [md].
func RenderDoc(md string) string {
	// Buffer the number of active goroutines
	maxRoutines := make(chan struct{}, MAX_ROUTINES)
	defer close(maxRoutines)

	chunks := chunk.ChunkDoc(md)

	n := len(chunks)

	wg.Add(n)
	out := make([]string, n)

	// Process chunks concurrently.
	for i, c := range chunks {
		go func(i int) {
			<-maxRoutines

			processChunk(i, c, out)
		}(i)

		maxRoutines <- struct{}{}
	}

	wg.Wait()

	return assembleDoc(out)
}
