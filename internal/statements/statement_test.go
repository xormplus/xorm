// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statements

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xormplus/xorm/caches"
	"github.com/xormplus/xorm/dialects"
	"github.com/xormplus/xorm/names"
	"github.com/xormplus/xorm/schemas"
	"github.com/xormplus/xorm/tags"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dialect   dialects.Dialect
	tagParser *tags.Parser
)

func TestMain(m *testing.M) {
	var err error
	dialect, err = dialects.OpenDialect("sqlite3", "./test.db")
	if err != nil {
		panic("unknow dialect")
	}

	tagParser = tags.NewParser("xorm", dialect, names.SnakeMapper{}, names.SnakeMapper{}, caches.NewManager())
	if tagParser == nil {
		panic("tags parser is nil")
	}
	m.Run()
}

var colStrTests = []struct {
	omitColumn        string
	onlyToDBColumnNdx int
	expected          string
}{
	{"", -1, "`ID`, `IsDeleted`, `Caption`, `Code1`, `Code2`, `Code3`, `ParentID`, `Latitude`, `Longitude`"},
	{"Code2", -1, "`ID`, `IsDeleted`, `Caption`, `Code1`, `Code3`, `ParentID`, `Latitude`, `Longitude`"},
	{"", 1, "`ID`, `Caption`, `Code1`, `Code2`, `Code3`, `ParentID`, `Latitude`, `Longitude`"},
	{"Code3", 1, "`ID`, `Caption`, `Code1`, `Code2`, `ParentID`, `Latitude`, `Longitude`"},
	{"Longitude", 1, "`ID`, `Caption`, `Code1`, `Code2`, `Code3`, `ParentID`, `Latitude`"},
	{"", 8, "`ID`, `IsDeleted`, `Caption`, `Code1`, `Code2`, `Code3`, `ParentID`, `Latitude`"},
}

func TestColumnsStringGeneration(t *testing.T) {
	for ndx, testCase := range colStrTests {
		statement, err := createTestStatement()
		assert.NoError(t, err)

		if testCase.omitColumn != "" {
			statement.Omit(testCase.omitColumn)
		}

		columns := statement.RefTable.Columns()
		if testCase.onlyToDBColumnNdx >= 0 {
			columns[testCase.onlyToDBColumnNdx].MapType = schemas.ONLYTODB
		}

		actual := statement.genColumnStr()

		if actual != testCase.expected {
			t.Errorf("[test #%d] Unexpected columns string:\nwant:\t%s\nhave:\t%s", ndx, testCase.expected, actual)
		}
		if testCase.onlyToDBColumnNdx >= 0 {
			columns[testCase.onlyToDBColumnNdx].MapType = schemas.TWOSIDES
		}
	}
}

func BenchmarkGetFlagForColumnWithICKey_ContainsKey(b *testing.B) {

	b.StopTimer()

	mapCols := make(map[string]bool)
	cols := []*schemas.Column{
		{Name: `ID`},
		{Name: `IsDeleted`},
		{Name: `Caption`},
		{Name: `Code1`},
		{Name: `Code2`},
		{Name: `Code3`},
		{Name: `ParentID`},
		{Name: `Latitude`},
		{Name: `Longitude`},
	}

	for _, col := range cols {
		mapCols[strings.ToLower(col.Name)] = true
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {

		for _, col := range cols {

			if _, ok := getFlagForColumn(mapCols, col); !ok {
				b.Fatal("Unexpected result")
			}
		}
	}
}

func BenchmarkGetFlagForColumnWithICKey_EmptyMap(b *testing.B) {

	b.StopTimer()

	mapCols := make(map[string]bool)
	cols := []*schemas.Column{
		{Name: `ID`},
		{Name: `IsDeleted`},
		{Name: `Caption`},
		{Name: `Code1`},
		{Name: `Code2`},
		{Name: `Code3`},
		{Name: `ParentID`},
		{Name: `Latitude`},
		{Name: `Longitude`},
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {

		for _, col := range cols {

			if _, ok := getFlagForColumn(mapCols, col); ok {
				b.Fatal("Unexpected result")
			}
		}
	}
}

type TestType struct {
	ID        int64   `xorm:"ID PK"`
	IsDeleted bool    `xorm:"IsDeleted"`
	Caption   string  `xorm:"Caption"`
	Code1     string  `xorm:"Code1"`
	Code2     string  `xorm:"Code2"`
	Code3     string  `xorm:"Code3"`
	ParentID  int64   `xorm:"ParentID"`
	Latitude  float64 `xorm:"Latitude"`
	Longitude float64 `xorm:"Longitude"`
}

func (TestType) TableName() string {
	return "TestTable"
}

func createTestStatement() (*Statement, error) {
	statement := NewStatement(dialect, tagParser, time.Local)
	if err := statement.SetRefValue(reflect.ValueOf(TestType{})); err != nil {
		return nil, err
	}
	return statement, nil
}

func BenchmarkColumnsStringGeneration(b *testing.B) {
	b.StopTimer()

	statement, err := createTestStatement()
	if err != nil {
		panic(err)
	}

	testCase := colStrTests[0]

	if testCase.omitColumn != "" {
		statement.Omit(testCase.omitColumn) // !nemec784! Column must be skipped
	}

	if testCase.onlyToDBColumnNdx >= 0 {
		columns := statement.RefTable.Columns()
		columns[testCase.onlyToDBColumnNdx].MapType = schemas.ONLYTODB // !nemec784! Column must be skipped
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		actual := statement.genColumnStr()

		if actual != testCase.expected {
			b.Errorf("Unexpected columns string:\nwant:\t%s\nhave:\t%s", testCase.expected, actual)
		}
	}
}
