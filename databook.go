// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"os"

	tablib "github.com/agrison/go-tablib"
)

type Databook struct {
	XDatabook *tablib.Databook
}

func NewDatabook() *Databook {
	db := &Databook{XDatabook: tablib.NewDatabook()}
	return db
}

func NewDatabookWithData(sheetName map[string]string, data interface{}, mustMatch bool, headers ...map[string][]string) (*Databook, error) {
	s := len(sheetName)

	switch data.(type) {
	case map[string]*tablib.Dataset:
		dataModel1 := data.(map[string]*tablib.Dataset)
		d1 := len(dataModel1)
		if s != d1 && mustMatch {
			return nil, ErrParamsType
		}

		databook := tablib.NewDatabook()
		for k, _ := range dataModel1 {
			if _, ok := sheetName[k]; !ok {
				return nil, ErrParamsType
			}
			databook.AddSheet(sheetName[k], dataModel1[k])
		}
		db := &Databook{XDatabook: databook}
		return db, nil
	case map[string][]map[string]interface{}:
		dataModel2 := data.(map[string][]map[string]interface{})
		d2 := len(dataModel2)

		if len(headers) != 1 {
			return nil, ErrParamsType
		}

		h := len(headers[0])

		if (s != h || s != d2) && mustMatch {
			return nil, ErrParamsType
		}
		databook := tablib.NewDatabook()
		for k, _ := range dataModel2 {
			if _, ok := sheetName[k]; !ok {
				return nil, ErrParamsType
			}
			if _, ok := headers[0][k]; !ok {
				return nil, ErrParamsType
			}

			dataset, err := NewDatasetWithData(headers[0][k], dataModel2[k], mustMatch)
			if err != nil {
				return nil, err
			}
			databook.AddSheet(sheetName[k], dataset)
		}
		db := &Databook{XDatabook: databook}
		return db, nil

	default:
		return nil, ErrParamsType
	}
}

func (databook *Databook) AddSheet(title string, data interface{}, mustMatch bool, headers ...[]string) error {

	switch data.(type) {
	case *tablib.Dataset:
		dataset := data.(*tablib.Dataset)
		databook.XDatabook.AddSheet(title, dataset)
		return nil
	case []map[string]interface{}:
		dataSlice := data.([]map[string]interface{})
		if len(headers) != 1 {
			return ErrParamsType
		}
		dataset, err := NewDatasetWithData(headers[0], dataSlice, mustMatch)
		if err != nil {
			return err
		}

		databook.XDatabook.AddSheet(title, dataset)
		return nil
	default:
		return ErrParamsType
	}
}

func (databook *Databook) HTML() *tablib.Exportable {
	return databook.XDatabook.HTML()
}

func (databook *Databook) SaveAsHTML(filename string, perm os.FileMode) error {
	html := databook.XDatabook.HTML()
	return html.WriteFile(filename, perm)
}

func (databook *Databook) JSON() (*tablib.Exportable, error) {
	return databook.XDatabook.JSON()
}

func (databook *Databook) SaveAsJSON(filename string, perm os.FileMode) error {
	json, err := databook.XDatabook.JSON()
	if err != nil {
		return err
	}
	return json.WriteFile(filename, perm)
}

func (databook *Databook) XLSX() (*tablib.Exportable, error) {
	return databook.XDatabook.XLSX()
}

func (databook *Databook) SaveAsXLSX(filename string, perm os.FileMode) error {
	xlsx, err := databook.XDatabook.XLSX()
	if err != nil {
		return err
	}
	return xlsx.WriteFile(filename, perm)
}

func (databook *Databook) XML() (*tablib.Exportable, error) {
	return databook.XDatabook.XML()
}

func (databook *Databook) SaveAsXML(filename string, perm os.FileMode) error {
	xml, err := databook.XDatabook.XML()
	if err != nil {
		return err
	}
	return xml.WriteFile(filename, perm)
}

func (databook *Databook) YAML() (*tablib.Exportable, error) {
	return databook.XDatabook.YAML()
}

func (databook *Databook) SaveAsYAML(filename string, perm os.FileMode) error {
	yaml, err := databook.XDatabook.YAML()
	if err != nil {
		return err
	}
	return yaml.WriteFile(filename, perm)
}

func (databook *Databook) Sheet(title string) tablib.Sheet {
	return databook.XDatabook.Sheet(title)
}

func (databook *Databook) Sheets() map[string]tablib.Sheet {
	return databook.XDatabook.Sheets()
}

func (databook *Databook) Size() int {
	return databook.XDatabook.Size()
}

func (databook *Databook) Wipe() {
	databook.XDatabook.Wipe()
}
