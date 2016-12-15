package tablib

// HTML returns the HTML representation of the Dataset as an Exportable.
func (d *Dataset) HTML() *Exportable {
	back := d.Records()
	b := newBuffer()

	b.WriteString("<table class=\"table table-striped\">\n\t<thead>")
	for i, r := range back {
		b.WriteString("\n\t\t<tr>")
		for _, c := range r {
			tag := "td"
			if i == 0 {
				tag = "th"
			}
			b.WriteString("\n\t\t\t<" + tag + ">")
			b.WriteString(c)
			b.WriteString("</" + tag + ">")
		}
		b.WriteString("\n\t\t</tr>")
		if i == 0 {
			b.WriteString("\n\t</thead>\n\t<tbody>")
		}
	}
	b.WriteString("\n\t</tbody>\n</table>")

	return newExportable(b)
}

// HTML returns a HTML representation of the Databook as an Exportable.
func (d *Databook) HTML() *Exportable {
	b := newBuffer()

	for _, s := range d.sheets {
		b.WriteString("<h1>" + s.title + "</h1>\n")
		b.Write(s.dataset.HTML().Bytes())
		b.WriteString("\n\n")
	}

	return newExportable(b)
}
