package tablib

import (
	"bytes"
	"github.com/agrison/mxj"
)

// XML returns a XML representation of the Dataset as an Exportable.
func (d *Dataset) XML() (*Exportable, error) {
	return d.XMLWithTagNamePrefixIndent("row", "  ", "  ")
}

// XML returns a XML representation of the Databook as an Exportable.
func (d *Databook) XML() (*Exportable, error) {
	b := newBuffer()
	b.WriteString("<databook>\n")
	for _, s := range d.sheets {
		b.WriteString("  <sheet>\n    <title>" + s.title + "</title>\n    ")
		row, err := s.dataset.XMLWithTagNamePrefixIndent("row", "      ", "  ")
		if err != nil {
			return nil, err
		}
		b.Write(row.Bytes())
		b.WriteString("\n  </sheet>")
	}
	b.WriteString("\n</databook>")
	return newExportable(b), nil
}

// XMLWithTagNamePrefixIndent returns a XML representation with custom tag, prefix and indent.
func (d *Dataset) XMLWithTagNamePrefixIndent(tagName, prefix, indent string) (*Exportable, error) {
	back := d.Dict()

	exportable := newExportable(newBuffer())
	exportable.buffer.WriteString("<dataset>\n")
	for _, r := range back {
		m := mxj.Map(r.(map[string]interface{}))
		if err := m.XmlIndentWriter(exportable.buffer, prefix, indent, tagName); err != nil {
			return nil, err
		}
	}
	exportable.buffer.WriteString("\n" + prefix + "</dataset>")

	return exportable, nil
}

// LoadXML loads a Dataset from an XML source.
func LoadXML(input []byte) (*Dataset, error) {
	m, _, err := mxj.NewMapXmlReaderRaw(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	// this seems quite a bit hacky
	datasetNode, _ := m.ValueForPath("dataset")
	rowNode := datasetNode.(map[string]interface{})["row"].([]interface{})

	back := make([]map[string]interface{}, 0, len(rowNode))
	for _, r := range rowNode {
		back = append(back, r.(map[string]interface{}))
	}

	return internalLoadFromDict(back)
}
