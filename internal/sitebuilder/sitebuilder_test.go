package sitebuilder

import (
	"testing"
)

func TestHTMLDoc(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		HTMLDoc(Document{})
	})

	t.Run("Basic", func(t *testing.T) {
		HTMLDoc(Document{Title: "Some Title", Content: "<p>hello, world!</p>"})
	})
}
