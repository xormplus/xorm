package gotabulate

import "strconv"
import "fmt"

// Create normalized Array from strings
func createFromString(data [][]string) []*TabulateRow {
	rows := make([]*TabulateRow, len(data))

	for index, el := range data {
		rows[index] = &TabulateRow{Elements: el}
	}
	return rows
}

// Create normalized array of rows from mixed data (interface{})
func createFromMixed(data [][]interface{}, format byte) []*TabulateRow {
	rows := make([]*TabulateRow, len(data))
	for index_1, element := range data {
		normalized := make([]string, len(element))
		for index, el := range element {
			switch el.(type) {
			case int32:
				quoted := strconv.QuoteRuneToASCII(el.(int32))
				normalized[index] = quoted[1 : len(quoted)-1]
			case int:
				normalized[index] = strconv.Itoa(el.(int))
			case int64:
				normalized[index] = strconv.FormatInt(el.(int64), 10)
			case bool:
				normalized[index] = strconv.FormatBool(el.(bool))
			case float64:
				normalized[index] = strconv.FormatFloat(el.(float64), format, -1, 64)
			case uint64:
				normalized[index] = strconv.FormatUint(el.(uint64), 10)
			case nil:
				normalized[index] = "nil"
			default:
				normalized[index] = fmt.Sprintf("%s", el)
			}
		}
		rows[index_1] = &TabulateRow{Elements: normalized}
	}
	return rows
}

// Create normalized array from ints
func createFromInt(data [][]int) []*TabulateRow {
	rows := make([]*TabulateRow, len(data))
	for index_1, arr := range data {
		row := make([]string, len(arr))
		for index, el := range arr {
			row[index] = strconv.Itoa(el)
		}
		rows[index_1] = &TabulateRow{Elements: row}
	}
	return rows
}

// Create normalized array from float64
func createFromFloat64(data [][]float64, format byte) []*TabulateRow {
	rows := make([]*TabulateRow, len(data))
	for index_1, arr := range data {
		row := make([]string, len(arr))
		for index, el := range arr {
			row[index] = strconv.FormatFloat(el, format, -1, 64)
		}
		rows[index_1] = &TabulateRow{Elements: row}
	}
	return rows
}

// Create normalized array from ints32
func createFromInt32(data [][]int32) []*TabulateRow {
	rows := make([]*TabulateRow, len(data))
	for index_1, arr := range data {
		row := make([]string, len(arr))
		for index, el := range arr {
			quoted := strconv.QuoteRuneToASCII(el)
			row[index] = quoted[1 : len(quoted)-1]
		}
		rows[index_1] = &TabulateRow{Elements: row}
	}
	return rows
}

// Create normalized array from ints64
func createFromInt64(data [][]int64) []*TabulateRow {
	rows := make([]*TabulateRow, len(data))
	for index_1, arr := range data {
		row := make([]string, len(arr))
		for index, el := range arr {
			row[index] = strconv.FormatInt(el, 10)
		}
		rows[index_1] = &TabulateRow{Elements: row}
	}
	return rows
}

// Create normalized array from bools
func createFromBool(data [][]bool) []*TabulateRow {
	rows := make([]*TabulateRow, len(data))
	for index_1, arr := range data {
		row := make([]string, len(arr))
		for index, el := range arr {
			row[index] = strconv.FormatBool(el)
		}
		rows[index_1] = &TabulateRow{Elements: row}
	}
	return rows
}

// Create normalized array from a map of mixed elements (interface{})
// Keys will be used as header
func createFromMapMixed(data map[string][]interface{}, format byte) (headers []string, tData []*TabulateRow) {

	var dataslice [][]interface{}
	for key, value := range data {
		headers = append(headers, key)
		dataslice = append(dataslice, value)
	}
	return headers, createFromMixed(dataslice, format)
}

// Create normalized array from Map of strings
// Keys will be used as header
func createFromMapString(data map[string][]string) (headers []string, tData []*TabulateRow) {
	var dataslice [][]string
	for key, value := range data {
		headers = append(headers, key)
		dataslice = append(dataslice, value)
	}
	return headers, createFromString(dataslice)
}

// Check if element is present in a slice.
func inSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
