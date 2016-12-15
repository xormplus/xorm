package tablib

import "gopkg.in/yaml.v2"

// LoadYAML loads a dataset from a YAML source.
func LoadYAML(yamlContent []byte) (*Dataset, error) {
	var input []map[string]interface{}
	if err := yaml.Unmarshal(yamlContent, &input); err != nil {
		return nil, err
	}

	return internalLoadFromDict(input)
}

// LoadDatabookYAML loads a Databook from a YAML source.
func LoadDatabookYAML(yamlContent []byte) (*Databook, error) {
	var input []map[string]interface{}
	var internalInput []map[string]interface{}
	if err := yaml.Unmarshal(yamlContent, &input); err != nil {
		return nil, err
	}

	db := NewDatabook()
	for _, d := range input {
		b, err := yaml.Marshal(d["data"])
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(b, &internalInput); err != nil {
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

// YAML returns a YAML representation of the Dataset as an Exportable.
func (d *Dataset) YAML() (*Exportable, error) {
	back := d.Dict()

	b, err := yaml.Marshal(back)
	if err != nil {
		return nil, err
	}
	return newExportableFromBytes(b), nil
}

// YAML returns a YAML representation of the Databook as an Exportable.
func (d *Databook) YAML() (*Exportable, error) {
	y := make([]map[string]interface{}, len(d.sheets))
	i := 0
	for _, s := range d.sheets {
		y[i] = make(map[string]interface{})
		y[i]["title"] = s.title
		y[i]["data"] = s.dataset.Dict()
		i++
	}
	b, err := yaml.Marshal(y)
	if err != nil {
		return nil, err
	}
	return newExportableFromBytes(b), nil
}
