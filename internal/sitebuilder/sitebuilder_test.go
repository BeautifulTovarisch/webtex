package sitebuilder

import (
	"os"
	"testing"
)

func TestHTMLDoc(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		HTMLDoc(os.Stdout, Document{})
	})

	t.Run("Basic", func(t *testing.T) {
		HTMLDoc(os.Stdout, Document{Title: "Some Title", Content: "<p>hello, world!</p>"})
	})
}
