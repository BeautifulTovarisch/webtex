// package texrender converts TeX code into SVGs. The host machine must have:
//
//   - pdflatex (and required packages)
//   - pdf2svg
//
// in order to function.
package texrender

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/beautifultovarisch/webtex/internal/logger"
)

func texWrapper(tex string) string {
	var b strings.Builder

	b.WriteString("\\documentclass{standalone}\n")
	b.WriteString("\\usepackage{tikz}\n")
	b.WriteString("\\usepackage{pgfplots}\n")
	b.WriteString("\\usepackage{graphicx}\n")
	b.WriteString("\\usepackage{xcolor}\n")
	b.WriteString("\\begin{document}\n")

	b.WriteString(tex)

	b.WriteString("\n\\end{document}")

	return b.String()
}

func toPDF(tex, dir string) error {
	pdflatex, err := exec.LookPath("pdflatex")
	if err != nil {
		return err
	}

	cmd := exec.Command(pdflatex, "-output-directory", dir)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, texWrapper(tex))
	}()

	out, _ := cmd.CombinedOutput()
	if str := string(out); strings.Contains(str, "!") || strings.Contains(str, "Emergency") {
		logger.Error(str)

		return errors.New(str)
	}

	return nil
}

func toSVG(dir string) error {
	pdf2svg, err := exec.LookPath("pdf2svg")
	if err != nil {
		return err
	}

	cmd := exec.Command(pdf2svg, filepath.Join(dir, "texput.pdf"), filepath.Join(dir, "texput.svg"))

	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

// Render accepts [tex] markup and produces an SVG represented as a string.
func Render(tex string) (string, error) {
	tmp, err := os.MkdirTemp("", "tex")
	if err != nil {
		logger.Error("Error creating temp directory: %s", err)

		return "", err
	}

	defer os.RemoveAll(tmp)

	if err := toPDF(tex, tmp); err != nil {
		logger.Error("Error processing TeX: %s", err)

		return "", err
	}

	if err := toSVG(tmp); err != nil {
		logger.Error("Error converting to SVG: %s", err)

		return "", err
	}

	svg, err := os.ReadFile(filepath.Join(tmp, "texput.svg"))
	if err != nil {
		logger.Error("Error reading SVG file: %s", err)

		return "", err
	}

	return string(svg), nil
}
