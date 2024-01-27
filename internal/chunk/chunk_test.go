package chunk

import (
	"os"
	"path/filepath"
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
			t.Fatalf("No input file corresponding to %s", f)
		}

		if !chunksEq(v, actual) {
			t.Errorf("Expected: %v. Got: %v", v, actual)
		}
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
			"block-1.md": []Chunk{},
			"block-2.md": []Chunk{Chunk{BLOCK, "a + b = c"}},
			"block-3.md": []Chunk{Chunk{BLOCK, "\n\\begin{tabular}{c c c}\na & b & c \\\\\n\\end{tabular}\n"}},
			"block-4.md": []Chunk{Chunk{BLOCK, "\n$x = 10$\n"}, Chunk{MD, "\n\n# Heading\n"}},
		}

		testFiles(files, expected, t)
	})

	t.Run("Malformed", func(t *testing.T) {
		matches, _ := filepath.Glob("testdata/malformed-*")

		for _, f := range matches {
			md, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}

			input := string(md)

			actual := ChunkDoc(input)
			expected := []Chunk{Chunk{MD, input}}

			if !chunksEq(expected, actual) {
				t.Errorf("Expected %v. Got %v", expected, actual)
			}
		}
	})

	t.Run("Inline", func(t *testing.T) {
		files, _ := filepath.Glob("testdata/inline-*")

		expected := map[string][]Chunk{
			"inline-1.md": []Chunk{Chunk{INLINE, "x + y = 10"}},
		}

		testFiles(files, expected, t)
	})

	t.Run("Heterogeneous", func(t *testing.T) {
		md := `
    # Heading

    Some text here

    ## SubHeading

    $x + y = z$

    ### Pythagorean equation: $x^2 + y^2 = z^2$

    Some notes on the Pythagorean Theorem.

    $$
    \begin{equation}
    E_n(x) = \frac 1 {n!} \int_1^x (x - t)^n f^(n+1)(t) \; dt
    \end{equation}
    $$

    $$
    \begin{tabular}{c c c}
    a & b & c \\
    d & e & f
    \end{tabular}
    $$
    `

		ChunkDoc(md)
	})
}
