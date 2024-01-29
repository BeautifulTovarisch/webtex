package render

import (
	"testing"
)

// TODO: These types of tests are a good fit for snapshot testing. Look into
// writing a diffing alg and potentially working up a lightweight snapshot tool
func TestRenderDoc(t *testing.T) {
	t.Run("Smoke", func(t *testing.T) {
		doc := `
    # Heading

    Some markdown text here

    **bolded** _italics_

    - Item 1
    - Item 2
      - SubItem 1
      - SubItem 2
    - Item 3

    ## Subheading
    `

		RenderDoc(doc)
	})
}
