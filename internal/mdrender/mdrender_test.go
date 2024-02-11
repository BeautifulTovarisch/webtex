package mdrender

import (
	"testing"
)

func TestRender(t *testing.T) {
	// These are just here to make sure no panics occur.
	t.Run("Smoke", func(t *testing.T) {
		Render("# Heading")
		Render("## Subheading")
		Render("```python\n[x for x in range(1, 11)]\n```")
	})
}
