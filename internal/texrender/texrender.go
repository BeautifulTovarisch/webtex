// package texrender converts TeX code into SVGs. The host machine must have:
//
//   - pdflatex (and required packages)
//   - pdf2svg
//
// in order to function.
package texrender

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Format proper latex document
func texDoc(tex string) string {
	var b strings.Builder

	b.WriteString("\\documentclass{standalone}\n")
	b.WriteString("\\usepackage{amsmath}\n")
	b.WriteString("\\usepackage{tikz}\n")
	b.WriteString("\\usepackage{pgfplots}\n")
	b.WriteString("\\usepackage{graphicx}\n")
	b.WriteString("\\usepackage{xcolor}\n")
	b.WriteString("\\begin{document}")
	b.WriteString(tex)
	b.WriteString("\\end{document}")

	return b.String()
}

func createSVG(dir string) error {
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

func createPDF(tex, dir string) error {
	pdflatex, err := exec.LookPath("pdflatex")
	if err != nil {
		return err
	}

	// TeX vomits out too much error output to reasonably convert into a Go error
	// here. Additionally, errors are reported on STDOUT. Attempting to convert a
	// missing file into an SVG will have to suffice as far for error reporting.
	cmd := exec.Command(pdflatex, "-file-line-error", "-output-directory", dir)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		fmt.Fprintln(stdin, texDoc(tex))
		stdin.Close()
	}()

	return cmd.Wait()
}

func render(tex string) (string, error) {
	tmp, err := os.MkdirTemp("", "tex")
	if err != nil {
		return "", err
	}

	if err := createPDF(tex, tmp); err != nil {
		return "", err
	}

	if err := createSVG(tmp); err != nil {
		return "", err
	}

	svg, err := os.ReadFile(filepath.Join(tmp, "texput.svg"))
	if err != nil {
		return "", err
	}

	return string(svg), nil
}

// RenderBlock accepts a block of [tex] and produces a corresponding SVG.
func RenderBlock(tex string) (string, error) {
	return render(tex)
}

// RenderInline accepts inline [tex] and produces a corresponding SVG.
func RenderInline(tex string) (string, error) {
	return render(fmt.Sprintf("$%s$", tex))
}
