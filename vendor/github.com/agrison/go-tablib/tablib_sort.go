package tablib

import "time"

// entryPair represents a pair of a value and its row index in the dataset
// which is used while sorting the dataset using a colum.
type entryPair struct {
	index int
	value interface{}
}

type byIntValue []entryPair

func (p byIntValue) Len() int           { return len(p) }
func (p byIntValue) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byIntValue) Less(i, j int) bool { return p[i].value.(int) < p[j].value.(int) }

type byInt64Value []entryPair

func (p byInt64Value) Len() int           { return len(p) }
func (p byInt64Value) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byInt64Value) Less(i, j int) bool { return p[i].value.(int64) < p[j].value.(int64) }

type byUint64Value []entryPair

func (p byUint64Value) Len() int           { return len(p) }
func (p byUint64Value) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byUint64Value) Less(i, j int) bool { return p[i].value.(uint64) < p[j].value.(uint64) }

type byFloatValue []entryPair

func (p byFloatValue) Len() int           { return len(p) }
func (p byFloatValue) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byFloatValue) Less(i, j int) bool { return p[i].value.(float64) < p[j].value.(float64) }

type byTimeValue []entryPair

func (p byTimeValue) Len() int      { return len(p) }
func (p byTimeValue) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p byTimeValue) Less(i, j int) bool {
	return p[i].value.(time.Time).UnixNano() < p[j].value.(time.Time).UnixNano()
}

type byStringValue []entryPair

func (p byStringValue) Len() int           { return len(p) }
func (p byStringValue) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p byStringValue) Less(i, j int) bool { return p[i].value.(string) < p[j].value.(string) }
