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
		tmp := t.TempDir()

		if err := Build("testdata/single", tmp); err != nil {
			t.Errorf("Failed to build site: %s", err)
		}
	})

	t.Run("Small", func(t *testing.T) {
		tmp := t.TempDir()

		if err := Build("testdata/Calculus/Exponents and Logarithms", tmp); err != nil {
			t.Errorf("Failed to build site: %s", err)
		}
	})

	t.Run("Medium", func(t *testing.T) {
		tmp := t.TempDir()

		if err := Build("testdata/Calculus/Integration", tmp); err != nil {
			t.Errorf("Failed to build site: %s", err)
		}
	})

	t.Run("Big", func(t *testing.T) {
		tmp := t.TempDir()

		if err := Build("testdata/Calculus", tmp); err != nil {
			t.Errorf("Failed to build site: %s", err)
		}
	})
}
