// package chunk partitions text into markdown and LaTeX blocks.
package chunk

import (
	"fmt"
	"strings"
	"unicode"
)

// ChunkType represents the nature of the contiguous block of content contained
// in the Chunk.
type ChunkType uint8

const (
	BLOCK  ChunkType = 0
	INLINE ChunkType = 1
	MD     ChunkType = 2
)

// Chunk is a contiguous block of either Markdown or LaTeX content.
type Chunk struct {
	T       ChunkType // Indicates whether the Chunk is markdown or LaTeX
	Content string    // The raw contents of the chunk of text.
}

func (c Chunk) String() string {
	t := map[ChunkType]string{
		MD:     "Markdown",
		BLOCK:  "Block",
		INLINE: "Inline",
	}

	return fmt.Sprintf("{%s %s}", t[c.T], c.Content)
}

// Check whether we have INLINE/BLOCK LaTeX or Markdown. In other words, if the
// dollar sign is truly a delimiter.
func checkType(str string) ChunkType {
	// Cannot possibly have a delimited block
	if len(str) < 2 {
		return MD
	}

	switch str[1] {
	case '$':
		return BLOCK
	case ' ':
		return MD
	default:
		// We check for a matching delimiter.
		if strings.Index(str[1:], "$") > 0 {
			return INLINE
		}

		return MD
	}
}

// Read until '$'
func readMd(str string) (Chunk, string) {
	end := strings.Index(str[1:], "$")
	if end < 0 {
		return Chunk{MD, str}, ""
	}

	return Chunk{MD, str[:end+1]}, str[end+1:]
}

// Read until terminating '$$' or end of document. Anything after a '$$' is a
// block.
//
// TODO: Consider supporting escaping dollar signs.
func readBlock(str string) (Chunk, string) {
	// The index will be the index of the next '$$' + 2 to accommodate the start
	end := strings.Index(str[2:], "$$")

	// No matching delimiter found. Read the rest of the document.
	if end < 0 {
		return Chunk{BLOCK, str[2:]}, ""
	}

	// +2 to skip past the delimiter
	return Chunk{BLOCK, str[2 : end+2]}, str[end+4:]
}

// Read valid Inline LaTeX or treat as Markdown.
func readInline(str string) (Chunk, string) {
	end := strings.Index(str[1:], "$")

	// Rest of document is markdown.
	if end < 0 {
		return Chunk{MD, str}, ""
	}

	preceding := str[end]

	// Example: $x + y = 10 $
	//                     ^
	if unicode.IsSpace(rune(preceding)) {
		return Chunk{MD, str[:end+2]}, str[end+2:]
	}

	return Chunk{INLINE, str[1 : end+1]}, str[end+2:]
}

// merge adjacent markdown chunks.
func mergeChunks(chunks []Chunk) []Chunk {
	if len(chunks) < 2 {
		return chunks
	}

	a, b := chunks[0], chunks[1]

	if a.T == MD && b.T == MD {
		merged := Chunk{a.T, a.Content + b.Content}
		newChunks := append([]Chunk{merged}, chunks[2:]...)

		return mergeChunks(newChunks)
	}

	return append([]Chunk{a}, mergeChunks(chunks[1:])...)
}

// Each recursive call, consider the first non-whitespace character.
// If we encounter a '$', we check whether we have LaTeX, and then whether its
// inline or a block.
func partition(md string) []Chunk {
	if strings.TrimSpace(md) == "" {
		return []Chunk{}
	}

	var (
		c   Chunk
		rem string
	)

	switch md[0] {
	case '$':
		// Check for LaTeX
		switch checkType(md) {
		case BLOCK:
			c, rem = readBlock(md)
		case INLINE:
			c, rem = readInline(md)
		default:
			c, rem = readMd(md)
		}
	case '`':
		// Check for fence
		c, rem = Chunk{}, ""
	default:
		// Markdown
		c, rem = readMd(md)
	}

	return append([]Chunk{c}, partition(rem)...)
}

// ChunkDoc partitions markdown content into three distinct types of "chunks":
//
//   - Markdown
//   - Inline LaTeX
//   - Block LaTeX
//
// Individual LaTeX chunks will include the contents of a properly formed block
// or inline LaTeX, e.g:
//
// $$\begin{equation} a^2 + b^2 = c^2 \end{equation}$$
//
// $$
//
//	x + y = z
//
// $$
//
// or
//
// $\int_1^x x \; dx$
//
// While markdown blocks are contiguous blocks of non-LaTeX content.
func ChunkDoc(md string) []Chunk {
	return mergeChunks(partition(md))
}
