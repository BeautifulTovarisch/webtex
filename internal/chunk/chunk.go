// package chunk lexs text into markdown and LaTeX blocks.
package chunk

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// ChunkType represents the nature of the contiguous block of content contained
// in the Chunk.
type ChunkType uint8

const (
	NULL ChunkType = iota
	MD
	CODE
	BLOCK
	FENCE
	INLINE
)

func (c ChunkType) String() string {
	t := map[ChunkType]string{
		NULL:   "Null",
		MD:     "Markdown",
		CODE:   "Code",
		BLOCK:  "Block",
		FENCE:  "Fence",
		INLINE: "Inline",
	}

	return t[c]
}

// Chunk is a contiguous block of either Markdown or LaTeX content.
type Chunk struct {
	T       ChunkType // Indicates whether the Chunk is markdown or LaTeX
	Content string    // The raw contents of the chunk of text.
}

func (c Chunk) String() string {
	return fmt.Sprintf("{%s %s}", c.T.String(), c.Content)
}

// Peek at the next character to determine whether we have one of (possible)
// delimiters:
//
//   - $
//   - `
//   - $$
//   - ```
func checkType(md *bufio.Reader) (ChunkType, error) {
	peek, err := md.Peek(1)
	if err != nil {
		return NULL, err
	}

	switch {
	case peek[0] == '$':
		peek, err = md.Peek(2)
		// Dollar sign at the very end of the document
		if err != nil {
			return MD, err
		}

		if string(peek) == "$$" {
			return BLOCK, nil
		}

		return INLINE, nil

	case peek[0] == '`':
		// If we reach EOF here, the only two possibilities are an unterminated '`'
		// or a closed fence with no content ('``')
		peek, err = md.Peek(3)
		if err != nil {
			return FENCE, err
		}

		if string(peek) == "```" {
			return CODE, nil
		}

		return FENCE, nil
	default:
		return MD, nil
	}
}

// Read until '$' or '`'
func readMd(md *bufio.Reader) (Chunk, error) {
	var b strings.Builder

	for {
		c, _, err := md.ReadRune()
		if err != nil {
			return Chunk{MD, b.String()}, err
		}

		if c == '$' || c == '`' {
			md.UnreadRune()

			return Chunk{MD, b.String()}, nil
		}

		b.WriteRune(c)
	}
}

// Read until terminating '$$' or end of document. Anything after a '$$' is a
// block.
//
// TODO: Consider supporting escaping dollar signs.
func readBlock(tex *bufio.Reader) (Chunk, error) {
	var b strings.Builder

	// Read off the first two '$$'
	tex.Discard(2)

	for {
		c, _, err := tex.ReadRune()
		if err != nil {
			return Chunk{BLOCK, strings.TrimSuffix(b.String(), "$")}, err
		}

		b.WriteRune(c)

		peek, err := tex.Peek(1)
		if err != nil {
			return Chunk{BLOCK, strings.TrimSuffix(b.String(), "$")}, err
		}

		if c == '$' && peek[0] == '$' {
			tex.Discard(1)

			return Chunk{BLOCK, strings.TrimSuffix(b.String(), "$")}, err
		}
	}
}

// Read valid Inline LaTeX or treat as Markdown.
func readInline(tex *bufio.Reader) (Chunk, error) {
	tex.Discard(1)

	content, err := tex.ReadString('$')

	if err != nil && err != io.EOF {
		return Chunk{}, err
	}

	n := len(content)
	if n > 1 {
		if unicode.IsSpace(rune(content[n-2])) {
			// Because we use ReadString() above, we cannot use UnreadRune()
			tex.UnreadByte()

			return Chunk{MD, strings.TrimSuffix("$"+content[:n-1], "$")}, err
		}

		return Chunk{INLINE, strings.TrimSuffix(content, "$")}, err
	}

	return Chunk{MD, "$" + content}, nil
}

// Fenced code block delimited with ```
//
//	 Example
//	   ```code
//		  code here
//		  ```
func readCodeBlock(md *bufio.Reader) (Chunk, error) {
	md.Discard(3)

	block, err := md.ReadString('`')
	if err != nil {
		return Chunk{MD, "```" + block}, err
	}

	// Read off the remaining backticks
	md.Discard(2)

	// Only need to append two backticks here, since ReadString includes one
	return Chunk{MD, "```" + block + "``"}, nil
}

// If a fence marker is found, all content until the matching delimiter will be
// treated as markdown. If no terminated delimiter is found, read until the end
// of the document.
//
//	`inline fence`
func readFence(md *bufio.Reader) (Chunk, error) {
	md.Discard(1)

	fence, err := md.ReadString('`')
	if err != nil {
		return Chunk{MD, fence}, err
	}

	return Chunk{MD, "`" + fence}, nil
}

// Check the first character to determine which type of content to read
func lex(md *bufio.Reader) (Chunk, error) {
	t, err := checkType(md)
	if err != nil {
		return Chunk{}, err
	}

	switch t {
	case BLOCK:
		return readBlock(md)
	case INLINE:
		return readInline(md)
	case FENCE:
		return readFence(md)
	case CODE:
		return readCodeBlock(md)
	default:
		return readMd(md)
	}
}

// ChunkDoc lexs markdown content into three distinct types of "chunks":
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
func ChunkDoc(md *bufio.Reader) (Chunk, error) {
	return lex(md)
}
