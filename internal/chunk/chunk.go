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
	BLOCK ChunkType = iota
	INLINE
	MD
	EOF
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
func checkType(md *bufio.Reader) (ChunkType, error) {
	c, _, err := md.ReadRune()
	if err != nil && err != io.EOF {
		return 0, err
	}

	switch c {
	case '$':
		return BLOCK, nil
	case ' ':
		// The space is technically markdown
		md.UnreadRune()

		return MD, nil
	default:
		md.UnreadRune()
		// We assume it's an inline block here, and determine whether it's actually
		// markdown during readInline
		return INLINE, nil
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
	latex, err := tex.ReadString('$')

	// Read off the second '$'
	tex.ReadRune()

	return Chunk{BLOCK, strings.Trim(latex, "$")}, err
}

// Read valid Inline LaTeX or treat as Markdown.
func readInline(tex *bufio.Reader) (Chunk, error) {
	content, err := tex.ReadString('$')
	if err != nil && err != io.EOF {
		return Chunk{}, err
	}

	// Check whether to return markdown or inline latex
	if n := len(content); n > 1 {
		// When we reach the next '$', decide whether the content forms a legitimate
		// inline latex segment.
		// Example: $x + y = 10 $
		//                     ^
		// We only know for sure that the characters up until the 2nd '$' are MD. The
		// 2nd '$' may start a valid inline block, etc.
		if unicode.IsSpace(rune(content[n-1])) {
			return Chunk{MD, content[:n-1]}, nil
		}

		return Chunk{INLINE, content[:n-1]}, nil
	}

	return Chunk{MD, content}, nil
}

func readCodeBlock(md *bufio.Reader) (Chunk, error) {
	var b strings.Builder

	// We need 5 ticks total (6 - 1 read off during checkType)
	for ticks := 0; ticks < 5; {
		c, _, err := md.ReadRune()
		if err != nil {
			if err == io.EOF {
				// Non-terminated code fence
				return Chunk{MD, "`" + b.String()}, nil
			}

			return Chunk{}, err
		}

		b.WriteRune(c)

		if c == '`' {
			ticks++
		}
	}

	return Chunk{MD, "`" + b.String()}, nil
}

// If a fence marker is found, all content until the matching delimiter will be
// treated as markdown. If no terminated delimiter is found, read until the end
// of the document.
//
//	`inline fence`
//
//	```code
//	code here
//	```
func readFence(md *bufio.Reader) (Chunk, error) {
	// If the first two runes are backticks, we read until find '```'
	chars, err := md.Peek(2)
	if err != nil && err != io.EOF {
		return Chunk{}, err
	}

	if len(chars) > 1 {
		if chars[0] == '`' && chars[1] == '`' {
			return readCodeBlock(md)
		}
	}

	fence, err := md.ReadString('`')
	if err != nil && err != io.EOF {
		return Chunk{}, err
	}

	// Replace the backtick we read off during checkType
	return Chunk{MD, "`" + fence}, nil
}

// Check the first character to determine which type of content to read
func lex(md *bufio.Reader) (Chunk, error) {
	c, _, err := md.ReadRune()
	if err != nil {
		if err != io.EOF {
			return Chunk{}, err
		}

		return Chunk{}, io.EOF
	}

	switch c {
	case '$':
		t, err := checkType(md)
		if err != nil {
			return Chunk{}, err
		}

		// Check for LaTeX
		switch t {
		case BLOCK:
			return readBlock(md)
		case INLINE:
			return readInline(md)
		default:
			return readMd(md)
		}
	case '`':
		// Check for fence. This is almost identical to reading latex sections.
		return readFence(md)
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
func ChunkDoc(md io.Reader) (Chunk, error) {
	stream := bufio.NewReader(md)

	return lex(stream)
}
