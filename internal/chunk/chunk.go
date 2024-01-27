// package chunk partitions text into markdown and LaTeX blocks.
package chunk

import (
	"fmt"
	"strings"
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
	if str == "" || len(str) < 2 {
		return MD
	}

	switch str[0] {
	case '$':
		if c := str[1]; c == ' ' {
			return MD
		}

		return BLOCK
	case ' ':
		return MD
	default:
		return INLINE
	}
}

// Read a valid Block LaTeX or treat as Markdown.
func readBlock(str string) (Chunk, string) {
	end := strings.Index(str, "$$")

	// No delimiter or space before ending $$
	if end < 0 {
		// Add back the sliced off '$' now that we know it's not LaTeX
		return Chunk{MD, "$" + str}, ""
	}

	// Invalid block, read up to and including the 'delimiters'
	// TODO: Decide precisely what invalidates a block
	if str[end-1] == ' ' {
		fmt.Println(str)
		return Chunk{MD, "$" + str[:end+2]}, str[end+2:]
	}

	// +2 to skip past the delimiter
	return Chunk{BLOCK, str[:end]}, str[end+2:]
}

// Read valid Inline LaTeX or treat as Markdown.
func readInline(str string) (Chunk, string) {
	return Chunk{}, ""
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

// Due to limitations in the partitioning scheme, we may obtain chunks of the
// same type adjacent in the same list. It is safe to merge their contents into
// a single chunk and continue.
//
// TODO: Find a way to avoid adding unnecessary blocks.
func mergeChunks(chunks []Chunk) []Chunk {
	if len(chunks) < 2 {
		return chunks
	}

	a, b := chunks[0], chunks[1]

	if a.T == b.T {
		merged := Chunk{a.T, a.Content + b.Content}
		newChunks := append([]Chunk{merged}, chunks[2:]...)

		return mergeChunks(newChunks)
	}

	return append([]Chunk{a}, mergeChunks(chunks[1:])...)
}

func ChunkDoc(md string) []Chunk {
	if md == "" {
		return []Chunk{}
	}

	// First '$'
	start := strings.Index(md, "$")
	if start < 0 {
		return []Chunk{Chunk{MD, md}}
	}

	// Everything before '$' is markdown
	markdown, rest := Chunk{MD, md[:start]}, md[start:]

	chunk, rem := readChunk(rest)

	return mergeChunks(append([]Chunk{markdown, chunk}, ChunkDoc(rem)...))
}
