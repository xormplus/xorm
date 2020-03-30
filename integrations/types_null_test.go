// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrations

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type NullType struct {
	Id           int `xorm:"pk autoincr"`
	Name         sql.NullString
	Age          sql.NullInt64
	Height       sql.NullFloat64
	IsMan        sql.NullBool `xorm:"null"`
	Nil          driver.Valuer
	CustomStruct CustomStruct `xorm:"varchar(64) null"`
}

type CustomStruct struct {
	Year  int
	Month int
	Day   int
}

func (CustomStruct) String() string {
	return "CustomStruct"
}

func (m *CustomStruct) Scan(value interface{}) error {
	if value == nil {
		m.Year, m.Month, m.Day = 0, 0, 0
		return nil
	}

	if s, ok := value.([]byte); ok {
		seps := strings.Split(string(s), "/")
		m.Year, _ = strconv.Atoi(seps[0])
		m.Month, _ = strconv.Atoi(seps[1])
		m.Day, _ = strconv.Atoi(seps[2])
		return nil
	}

	return errors.New("scan data not fit []byte")
}

func (m CustomStruct) Value() (driver.Value, error) {
	return fmt.Sprintf("%d/%d/%d", m.Year, m.Month, m.Day), nil
}

func TestCreateNullStructTable(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	err := testEngine.CreateTables(new(NullType))
	assert.NoError(t, err)
}

func TestDropNullStructTable(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	err := testEngine.DropTables(new(NullType))
	assert.NoError(t, err)
}

func TestNullStructInsert(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(NullType))

	item1 := new(NullType)
	_, err := testEngine.Insert(item1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, item1.Id)

	item := NullType{
		Name:   sql.NullString{String: "haolei", Valid: true},
		Age:    sql.NullInt64{Int64: 34, Valid: true},
		Height: sql.NullFloat64{Float64: 1.72, Valid: true},
		IsMan:  sql.NullBool{Bool: true, Valid: true},
		Nil:    nil,
	}
	_, err = testEngine.Insert(&item)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, item.Id)

	items := []NullType{}
	for i := 0; i < 5; i++ {
		item := NullType{
			Name:         sql.NullString{String: "haolei_" + fmt.Sprint(i+1), Valid: true},
			Age:          sql.NullInt64{Int64: 30 + int64(i), Valid: true},
			Height:       sql.NullFloat64{Float64: 1.5 + 1.1*float64(i), Valid: true},
			IsMan:        sql.NullBool{Bool: true, Valid: true},
			CustomStruct: CustomStruct{i, i + 1, i + 2},
			Nil:          nil,
		}
		items = append(items, item)
	}

	_, err = testEngine.Insert(&items)
	assert.NoError(t, err)

	items = make([]NullType, 0, 7)
	err = testEngine.Find(&items)
	assert.NoError(t, err)
	assert.EqualValues(t, 7, len(items))
}

func TestNullStructUpdate(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(NullType))

	_, err := testEngine.Insert([]NullType{
		{
			Name: sql.NullString{
				String: "name1",
				Valid:  true,
			},
		},
		{
			Name: sql.NullString{
				String: "name2",
				Valid:  true,
			},
		},
		{
			Name: sql.NullString{
				String: "name3",
				Valid:  true,
			},
		},
		{
			Name: sql.NullString{
				String: "name4",
				Valid:  true,
			},
		},
	})
	assert.NoError(t, err)

	if true { // 测试可插入NULL
		item := new(NullType)
		item.Age = sql.NullInt64{Int64: 23, Valid: true}
		item.Height = sql.NullFloat64{Float64: 0, Valid: false} // update to NULL

		affected, err := testEngine.ID(2).Cols("age", "height", "is_man").Update(item)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, affected)
	}

	if true { // 测试In update
		item := new(NullType)
		item.Age = sql.NullInt64{Int64: 23, Valid: true}
		affected, err := testEngine.In("id", 3, 4).Cols("age", "height", "is_man").Update(item)
		assert.NoError(t, err)
		assert.EqualValues(t, 2, affected)
	}

	if true { // 测试where
		item := new(NullType)
		item.Name = sql.NullString{String: "nullname", Valid: true}
		item.IsMan = sql.NullBool{Bool: true, Valid: true}
		item.Age = sql.NullInt64{Int64: 34, Valid: true}

		_, err := testEngine.Where("age > ?", 34).Update(item)
		assert.NoError(t, err)
	}

	if true { // 修改全部时，插入空值
		item := &NullType{
			Name:   sql.NullString{String: "winxxp", Valid: true},
			Age:    sql.NullInt64{Int64: 30, Valid: true},
			Height: sql.NullFloat64{Float64: 1.72, Valid: true},
			// IsMan:  sql.NullBool{true, true},
		}

		_, err := testEngine.AllCols().ID(6).Update(item)
		assert.NoError(t, err)
	}
}

func TestNullStructFind(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(NullType))

	_, err := testEngine.Insert([]NullType{
		{
			Name: sql.NullString{
				String: "name1",
				Valid:  false,
			},
		},
		{
			Name: sql.NullString{
				String: "name2",
				Valid:  true,
			},
		},
		{
			Name: sql.NullString{
				String: "name3",
				Valid:  true,
			},
		},
		{
			Name: sql.NullString{
				String: "name4",
				Valid:  true,
			},
		},
	})
	assert.NoError(t, err)

	if true {
		item := new(NullType)
		has, err := testEngine.ID(1).Get(item)
		assert.NoError(t, err)
		assert.True(t, has)
		assert.EqualValues(t, item.Id, 1)
		assert.False(t, item.Name.Valid)
		assert.False(t, item.Age.Valid)
		assert.False(t, item.Height.Valid)
		assert.False(t, item.IsMan.Valid)
	}

	if true {
		item := new(NullType)
		item.Id = 2
		has, err := testEngine.Get(item)
		assert.NoError(t, err)
		assert.True(t, has)
	}

	if true {
		item := make([]NullType, 0)
		err := testEngine.ID(2).Find(&item)
		assert.NoError(t, err)
	}

	if true {
		item := make([]NullType, 0)
		err := testEngine.Asc("age").Find(&item)
		assert.NoError(t, err)
	}
}

func TestNullStructIterate(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(NullType))

	if true {
		err := testEngine.Where("age IS NOT NULL").OrderBy("age").Iterate(new(NullType),
			func(i int, bean interface{}) error {
				nultype := bean.(*NullType)
				fmt.Println(i, nultype)
				return nil
			})
		assert.NoError(t, err)
	}
}

func TestNullStructCount(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(NullType))

	if true {
		item := new(NullType)
		_, err := testEngine.Where("age IS NOT NULL").Count(item)
		assert.NoError(t, err)
	}
}

func TestNullStructRows(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(NullType))

	item := new(NullType)
	rows, err := testEngine.Where("id > ?", 1).Rows(item)
	assert.NoError(t, err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(item)
		assert.NoError(t, err)
	}
}

func TestNullStructDelete(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(NullType))

	item := new(NullType)

	_, err := testEngine.ID(1).Delete(item)
	assert.NoError(t, err)

	_, err = testEngine.Where("id > ?", 1).Delete(item)
	assert.NoError(t, err)
}
