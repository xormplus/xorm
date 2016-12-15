package tablib

import (
	"github.com/bndr/gotabulate"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// TabularGrid is the value to be passed to gotabulate to render the table
	// as ASCII table with grid format
	TabularGrid = "grid"
	// TabularSimple is the value to be passed to gotabulate to render the table
	// as ASCII table with simple format
	TabularSimple = "simple"
	// TabularCondensed is the value to be passed to gotabulate to render the table
	// as ASCII table with condensed format
	TabularCondensed = "condensed"
	// TabularMarkdown is the value to be passed to gotabulate to render the table
	// as ASCII table with Markdown format
	TabularMarkdown = "markdown"
)

// Markdown returns a Markdown table Exportable representation of the Dataset.
func (d *Dataset) Markdown() *Exportable {
	return d.Tabular(TabularMarkdown)
}

// Tabular returns a tabular Exportable representation of the Dataset.
// format is either grid, simple, condensed or markdown.
func (d *Dataset) Tabular(format string) *Exportable {
	back := d.Records()
	t := gotabulate.Create(back)

	if format == TabularCondensed || format == TabularMarkdown {
		rendered := regexp.MustCompile("\n\n\\s").ReplaceAllString(t.Render("simple"), "\n ")
		if format == TabularMarkdown {
			firstLine := regexp.MustCompile("-\\s+-").ReplaceAllString(strings.Split(rendered, "\n")[0], "- | -")
			// now just locate the position of pipe characterds, and set them
			positions := make([]int, 0, d.cols-1)
			x := 0
			for _, c := range firstLine {
				if c == '|' {
					positions = append(positions, x)
				}
				x += utf8.RuneLen(c)
			}

			b := newBuffer()
			lines := strings.Split(rendered, "\n")
			for _, line := range lines[1 : len(lines)-2] {
				ipos := 0
				b.WriteString("| ")
				for _, pos := range positions {
					if ipos < len(line) && pos < len(line) {
						b.WriteString(line[ipos:pos])
						b.WriteString(" | ")
						ipos = pos + 1
					}
				}
				if ipos < len(line) {
					b.WriteString(line[ipos:])
				}
				b.WriteString(" | \n")
			}
			return newExportable(b)
		}
		return newExportableFromString(rendered)
	}
	return newExportableFromString(t.Render(format))
}
