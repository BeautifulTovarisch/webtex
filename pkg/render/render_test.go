package render

import (
	"os"
	"testing"
)

// TODO: These types of tests are a good fit for snapshot testing. Look into
// writing a diffing alg and potentially working up a lightweight snapshot tool
func TestRenderDoc(t *testing.T) {
	t.Run("Smoke", func(t *testing.T) {
		doc, err := os.ReadFile("testdata/sample.md")
		if err != nil {
			t.Fatal(err)
		}

		RenderDoc(string(doc))
	})
}
