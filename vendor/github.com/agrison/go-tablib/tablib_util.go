package tablib

import (
	"fmt"
	"strconv"
	"time"
)

// internalLoadFromDict creates a Dataset from an array of map representing columns.
func internalLoadFromDict(input []map[string]interface{}) (*Dataset, error) {
	// retrieve columns
	headers := make([]string, 0, 10)
	for h := range input[0] {
		headers = append(headers, h)
	}

	ds := NewDataset(headers)
	for _, e := range input {
		row := make([]interface{}, 0, len(headers))
		for _, h := range headers {
			row = append(row, e[h])
		}
		ds.AppendValues(row...)
	}

	return ds, nil
}

// isTagged checks if a tag is in an array of tags.
func isTagged(tag string, tags []string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

// asString returns a value as a string.
func (d *Dataset) asString(vv interface{}) string {
	var v string
	switch vv.(type) {
	case string:
		v = vv.(string)
	case int:
		v = strconv.Itoa(vv.(int))
	case int64:
		v = strconv.FormatInt(vv.(int64), 10)
	case uint64:
		v = strconv.FormatUint(vv.(uint64), 10)
	case bool:
		v = strconv.FormatBool(vv.(bool))
	case float64:
		v = strconv.FormatFloat(vv.(float64), 'G', -1, 32)
	case time.Time:
		v = vv.(time.Time).Format(time.RFC3339)
	default:
		if d.EmptyValue != "" {
			v = d.EmptyValue
		} else {
			v = fmt.Sprintf("%s", v)
		}
	}
	return v
}
