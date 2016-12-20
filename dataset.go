// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	tablib "github.com/agrison/go-tablib"
)

// NewDataset creates a new Dataset.
func NewDataset(headers []string) *tablib.Dataset {
	return tablib.NewDataset(headers)
}

// NewDatasetWithData creates a new Dataset.
func NewDatasetWithData(headers []string, data interface{}, mustMatch bool) (*tablib.Dataset, error) {
	if data == nil {
		return tablib.NewDatasetWithData(headers, nil), nil
	}
	n := len(headers)

	switch data.(type) {
	case [][]interface{}:
		return tablib.NewDatasetWithData(headers, data.([][]interface{})), nil
	case []map[string]interface{}:
		dataSlice := data.([]map[string]interface{})
		if len(dataSlice) > 0 {
			if len(dataSlice[0]) == 0 {
				return tablib.NewDatasetWithData(headers, make([][]interface{}, len(dataSlice))), nil
			} else {
				if n != len(dataSlice[0]) && mustMatch {
					return nil, ErrParamsType
				}
				mapHeaders := make(map[string]int, n)
				for i := 0; i < n; i++ {
					mapHeaders[headers[i]] = i
				}

				for k, _ := range dataSlice[0] {
					if _, ok := mapHeaders[k]; !ok {
						return nil, ErrParamsType
					}
				}

				d := tablib.NewDataset(headers)
				var row []interface{}
				for i, _ := range dataSlice {
					row = nil
					for j := 0; j < n; j++ {
						row = append(row, dataSlice[i][headers[j]])
					}
					d.Append(row)
				}
				return d, nil

			}

		} else {
			return tablib.NewDatasetWithData(headers, nil), nil
		}
	default:
		return nil, ErrParamsType
	}

}
