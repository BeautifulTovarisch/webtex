package build

import (
	"testing"
)

func TestSiteNav(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		if _, err := SiteNav("testdata/"); err != nil {
			t.Error(err)
		}
	})

	t.Run("ErrorPath", func(t *testing.T) {
		_, err := SiteNav("NON_EXISTENT")
		if err == nil {
			t.Errorf("Failed to return error on non-existent directory")
		}
	})
}

func TestBuild(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		if err := Build("testdata/source/single", ""); err != nil {
			t.Errorf("Failed to build site: %s", err)
		}
	})
}
