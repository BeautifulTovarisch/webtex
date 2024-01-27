package chunk

import (
	"os"
	"path/filepath"
	"testing"
)

// TODO: Find or generate large MD files and read them from disk for tests.

func chunksEq(a, b []Chunk) bool {
	if len(a) != len(b) {
		return false
	}

	for i, c := range a {
		if c.T != b[i].T || c.Content != b[i].Content {
			return false
		}
	}

	return true
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
		validBlocks := []string{
			// Valid LaTeX in Obsidian, although breaks formatting.
			// `$$$$`,
			`$$a + b = c$$`,
			`$$
      \begin{tabular}{c c c}
      a & b & c
      \end{tabular}
      $$`,
		}

		expected := [][]Chunk{
			// []Chunk{Chunk{BLOCK, ""}},
			[]Chunk{Chunk{BLOCK, "a + b = c"}},
			[]Chunk{Chunk{BLOCK, "\n\\begin{tabular}{c c c}\na & b & c\n\\end{tabular}\n"}},
		}

		for i, block := range validBlocks {
			if actual := ChunkDoc(block); !chunksEq(actual, expected[i]) {
				t.Errorf("Expected: %v. Got %v", expected[i], actual)
			}
		}
	})

	t.Run("Malformed", func(t *testing.T) {
		matches, err := filepath.Glob("testdata/malformed-*")
		if err != nil {
			t.Error(err)
		}

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
