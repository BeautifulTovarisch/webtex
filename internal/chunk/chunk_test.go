package chunk

// TODO: Fix messy testing utilities.
// TODO: Write proper diffing algorithm(s)

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TODO: Find or generate large MD files and read them from disk for tests.

func chunkEq(a, b Chunk) bool {
	return a.T == b.T && a.Content == b.Content
}

func chunksEq(a, b []Chunk) bool {
	if len(a) != len(b) {
		return false
	}

	for i, c := range a {
		if !chunkEq(c, b[i]) {
			return false
		}
	}

	return true
}

func diffTable(a, b string, t *testing.T) {
	var buf strings.Builder

	n := len(b)

	fmt.Fprintf(&buf, "\nDIFF\n")

	for i, c1 := range a {
		if i < n {
			if c2 := b[i]; c1 != rune(c2) {
				// Highlight row as red
				fmt.Fprintf(&buf, "%q", '>')
				fmt.Fprintf(&buf, "%+q\t%+q\n", c1, c2)
			}
		}
	}

	t.Log(buf.String())
}

// Chunk by chunk comparision
func cmpChunk(expected, actual []Chunk, t *testing.T) {
	if len(expected) == 0 && len(actual) == 0 {
		// Pass
		return
	}

	if len(expected) == 0 {
		t.Errorf("Expected fewer chunks than received. Actual: %v\n", actual)
		return
	}

	if len(actual) == 0 {
		t.Errorf("Recieved fewer chunks then expected. Expected: %v\n", expected)
		return
	}

	a, b := expected[0], actual[0]

	if !chunkEq(a, b) {
		t.Errorf("Expected: %v\n\nActual: %v\n\n", a, b)

		diffTable(a.Content, b.Content, t)

		return
	}

	cmpChunk(expected[1:], actual[1:], t)
}

func testFiles(files []string, expected map[string][]Chunk, t *testing.T) {
	for _, f := range files {
		md, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		input := string(md)

		actual := ChunkDoc(input)

		v, ok := expected[filepath.Base(f)]
		if !ok {
			t.Fatalf("No test case corresponding to %s", f)
		}

		cmpChunk(v, actual, t)
	}
}

func TestChunkDoc(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		md := ""

		actual := ChunkDoc(md)

		if len(actual) != 0 {
			t.Errorf("Failed to chunk empty document. Got %v", actual)
		}
	})

	t.Run("Block", func(t *testing.T) {
		files, _ := filepath.Glob("testdata/block-*")

		expected := map[string][]Chunk{
			"block-1.md": []Chunk{Chunk{BLOCK, ""}},
			"block-2.md": []Chunk{Chunk{BLOCK, "a + b = c"}},
			"block-3.md": []Chunk{Chunk{BLOCK, "\n\\begin{tabular}{c c c}\na & b & c \\\\\n\\end{tabular}\n"}},
			"block-4.md": []Chunk{
				Chunk{BLOCK, "a + b = c"},
				Chunk{MD, "\n"},
				Chunk{BLOCK, "\n\\begin{tabular}{c c c}\na & b & c\n\\end{tabular}\n"},
			},
			"block-5.md": []Chunk{
				Chunk{BLOCK, "\n\\begin{tabular}{c c c}\na & b & c\n\\end{tabular}\n"},
				Chunk{MD, "\n\n"},
				Chunk{BLOCK, "\n\\begin{equation}\n$x + y = z$\n\\end{equation}\n"},
			},
		}

		testFiles(files, expected, t)
	})

	t.Run("Malformed", func(t *testing.T) {
		files, _ := filepath.Glob("testdata/malformed-*")

		expected := map[string][]Chunk{
			"malformed-1.md": []Chunk{Chunk{BLOCK, "\n\\begin{equation}\n\nMore text\n"}},
			"malformed-2.md": []Chunk{Chunk{BLOCK, "\\begin{equation}x + y = z\\end{equation}$abc\n"}},
			"malformed-3.md": []Chunk{
				// Remember consecutive markdown blocks are merged!
				Chunk{MD, "$x + y = z $\n$ "},
				Chunk{INLINE, "100"},
				Chunk{MD, "\n$x = -b \\pm \\frac {\\sqrt{b^2 - 4ac}} {2a}\n$"},
			},
		}

		testFiles(files, expected, t)
	})

	t.Run("Inline", func(t *testing.T) {
		files, _ := filepath.Glob("testdata/inline-*")

		expected := map[string][]Chunk{
			"inline-1.md": []Chunk{Chunk{INLINE, "x + y = 10"}},
		}

		testFiles(files, expected, t)
	})

	t.Run("Fence", func(t *testing.T) {
		files, _ := filepath.Glob("testdata/fence-*")

		expected := map[string][]Chunk{
			"fence-1.md": []Chunk{Chunk{MD, "```python\ndef fib(n):\n    if n <= 1:\n        return 1\n\n    return fib(n-1) + fib(n-2)\n```\n"}},
		}

		testFiles(files, expected, t)
	})

	t.Run("Heterogeneous", func(t *testing.T) {
		files, _ := filepath.Glob("testdata/hetero-*")

		expected := map[string][]Chunk{
			"hetero-1.md": []Chunk{
				Chunk{MD, "# Heading\n\nSome text here\n\n## SubHeading\n\n"},
				Chunk{INLINE, "x + y = z"},
				Chunk{MD, "\n\n### Pythagorean equation: "},
				Chunk{INLINE, "x^2 + y^2 = z^2"},
				Chunk{MD, "\n\nSome notes on the Pythagorean Theorem.\n\n"},
				Chunk{BLOCK, "\n\\begin{equation}\nE_n(x) = \\frac 1 {n!} \\int_1^x (x - t)^n f^{(n+1)}(t) \\; dt\n\\end{equation}\n"},
				Chunk{MD, "\n\n"},
				Chunk{BLOCK, "\n\\begin{tabular}{c c c}\na & b & c \\\\\nd & e & f\n\\end{tabular}\n"},
			},
		}

		testFiles(files, expected, t)
	})
}
