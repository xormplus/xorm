package tablib

import (
	"encoding/json"
)

// LoadJSON loads a dataset from a YAML source.
func LoadJSON(jsonContent []byte) (*Dataset, error) {
	var input []map[string]interface{}
	if err := json.Unmarshal(jsonContent, &input); err != nil {
		return nil, err
	}

	return internalLoadFromDict(input)
}

// LoadDatabookJSON loads a Databook from a JSON source.
func LoadDatabookJSON(jsonContent []byte) (*Databook, error) {
	var input []map[string]interface{}
	var internalInput []map[string]interface{}
	if err := json.Unmarshal(jsonContent, &input); err != nil {
		return nil, err
	}

	db := NewDatabook()
	for _, d := range input {
		b, err := json.Marshal(d["data"])
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &internalInput); err != nil {
			return nil, err
		}

		if ds, err := internalLoadFromDict(internalInput); err == nil {
			db.AddSheet(d["title"].(string), ds)
		} else {
			return nil, err
		}
	}

	return db, nil
}

// JSON returns a JSON representation of the Dataset as an Exportable.
func (d *Dataset) JSON() (*Exportable, error) {
	back := d.Dict()

	b, err := json.Marshal(back)
	if err != nil {
		return nil, err
	}
	return newExportableFromBytes(b), nil
}

// JSON returns a JSON representation of the Databook as an Exportable.
func (d *Databook) JSON() (*Exportable, error) {
	b := newBuffer()
	b.WriteString("[")
	for _, s := range d.sheets {
		b.WriteString("{\"title\": \"" + s.title + "\", \"data\": ")
		js, err := s.dataset.JSON()
		if err != nil {
			return nil, err
		}
		b.Write(js.Bytes())
		b.WriteString("},")
	}
	by := b.Bytes()
	by[len(by)-1] = ']'
	return newExportableFromBytes(by), nil
}
