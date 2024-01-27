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

// Determine whether we should read an inline, block, or markdown section
// based on the sequence characters following the first '$'.
//
// Examples:
//
// $$:
// $$a + b + c$$  => Block
// $$ a + b + c$$ => Markdown (space after $$)
// $$a + b + c $$ => Markdown (space before $$)
// $$             => Markdown (no terminating $$)
//
// $:
// $a + b = c$    => Inline
// $ a + b = c$   => Markdown
// $a + b = c $   => Markdown
// $              => Markdown (no terminating $)
func chunkType(str string) ChunkType {
	if len(str) < 2 {
		return MD
	}

	fmt.Println(str)

	if str[:2] == "$$" {
		return BLOCK
	}

	if str[0] == '$' && !unicode.IsSpace(rune(str[1])) {
		return INLINE
	}

	return MD
}

// Read a valid Block LaTeX or treat as Markdown.
// This function should always receive a string beginning with a '$', which is
// then sliced off to form the chunk.
func readBlock(str string) (Chunk, string) {
	end := strings.Index(str, "$$")

	// Empty block ($$$$)
	if end == 0 {
		// Skip past the extra '$'
		return Chunk{BLOCK, ""}, str[end+3:]
	}

	// No matching delimiter found. Treating as Markdown
	if end < 0 {
		return Chunk{MD, "$" + str}, ""
	}

	// The character immediately preceding the end delimiter
	preceding := str[end-1]

	// Whitespace before '$$'. Treating as Markdown. Newlines allowed.
	if preceding != '\n' && unicode.IsSpace(rune(preceding)) {
		return Chunk{MD, "$" + str}, str[end+2:]
	}

	// +2 to skip past the delimiter
	return Chunk{BLOCK, str[1:end]}, str[end+2:]
}

// Read valid Inline LaTeX or treat as Markdown.
func readInline(str string) (Chunk, string) {
	end := strings.Index(str, "$")

	if end == 0 {
		// This is an implementation bug, if two consecutive '$' appear, we should be
		// reading a block.
		panic(fmt.Sprintf("Implementation error: %s", str))
	}

	// Rest of document is markdown.
	if end < 0 {
		return Chunk{MD, str}, ""
	}

	preceding := str[end-1]

	// Example: $x + y = 10 $
	//                     ^
	if unicode.IsSpace(rune(preceding)) {
		return Chunk{MD, str}, str[end+1:]
	}

	return Chunk{INLINE, str[:end]}, str[end+1:]
}

// Read [str] into a Chunk labeled with the appropriate chunk type.
func readChunk(str string) (Chunk, string) {
	if str == "" {
		return Chunk{MD, ""}, ""
	}

	switch rem := str[1:]; chunkType(str) {
	case BLOCK:
		return readBlock(rem)
	case INLINE:
		return readInline(rem)
	default:
		// If there is no matching delimiter, the rest of the file is markdown.
		return Chunk{MD, str}, ""
	}
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

func partition(md string) []Chunk {
	if strings.TrimSpace(md) == "" {
		return []Chunk{}
	}

	// First '$'
	start := strings.Index(md, "$")
	if start < 0 {
		return []Chunk{Chunk{MD, md}}
	}

	// Everything before '$' is markdown
	markdown, candidate := Chunk{MD, md[:start]}, md[start:]

	// '$' detected, so we have a candidate for a LaTeX block.
	chunk, rem := readChunk(candidate)

	var chunks []Chunk

	if markdown.Content != "" {
		chunks = append(chunks, markdown)
	}

	if chunk.Content != "" {
		chunks = append(chunks, chunk)
	}

	return mergeChunks(append(chunks, ChunkDoc(rem)...))
}

// ChunkDoc partitions markdown content into three distinct types of "chunks":
//
// - Markdown
// - Inline LaTeX
// - Block LaTeX
//
// Individual LaTeX chunks will include the contents of a properly formed block
// or inline LaTeX, e.g:
//
// $$\begin{equation} a^2 + b^2 = c^2 \end{equation}$$
//
// $$
// x + y = z
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
