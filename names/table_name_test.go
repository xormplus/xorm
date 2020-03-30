// Copyright 2020 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package names

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Userinfo struct {
	Uid        int64  `xorm:"id pk not null autoincr"`
	Username   string `xorm:"unique"`
	Departname string
	Alias      string `xorm:"-"`
	Created    time.Time
	Detail     Userdetail `xorm:"detail_id int(11)"`
	Height     float64
	Avatar     []byte
	IsMan      bool
}

type Userdetail struct {
	Id      int64
	Intro   string `xorm:"text"`
	Profile string `xorm:"varchar(2000)"`
}

type MyGetCustomTableImpletation struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

const getCustomTableName = "GetCustomTableInterface"

func (MyGetCustomTableImpletation) TableName() string {
	return getCustomTableName
}

type TestTableNameStruct struct{}

const getTestTableName = "my_test_table_name_struct"

func (t *TestTableNameStruct) TableName() string {
	return getTestTableName
}

func TestGetTableName(t *testing.T) {
	var kases = []struct {
		mapper            Mapper
		v                 reflect.Value
		expectedTableName string
	}{
		{
			SnakeMapper{},
			reflect.ValueOf(new(Userinfo)),
			"userinfo",
		},
		{
			SnakeMapper{},
			reflect.ValueOf(Userinfo{}),
			"userinfo",
		},
		{
			SameMapper{},
			reflect.ValueOf(new(Userinfo)),
			"Userinfo",
		},
		{
			SameMapper{},
			reflect.ValueOf(Userinfo{}),
			"Userinfo",
		},
		{
			SnakeMapper{},
			reflect.ValueOf(new(MyGetCustomTableImpletation)),
			getCustomTableName,
		},
		{
			SnakeMapper{},
			reflect.ValueOf(MyGetCustomTableImpletation{}),
			getCustomTableName,
		},
		{
			SnakeMapper{},
			reflect.ValueOf(new(TestTableNameStruct)),
			new(TestTableNameStruct).TableName(),
		},
		{
			SnakeMapper{},
			reflect.ValueOf(new(TestTableNameStruct)),
			getTestTableName,
		},
		{
			SnakeMapper{},
			reflect.ValueOf(TestTableNameStruct{}),
			getTestTableName,
		},
	}

	for _, kase := range kases {
		assert.EqualValues(t, kase.expectedTableName, GetTableName(kase.mapper, kase.v))
	}
}

type OAuth2Application struct {
}

// TableName sets the table name to `oauth2_application`
func (app *OAuth2Application) TableName() string {
	return "oauth2_application"
}

func TestGonicMapperCustomTable(t *testing.T) {
	assert.EqualValues(t, "oauth2_application",
		GetTableName(LintGonicMapper, reflect.ValueOf(new(OAuth2Application))))
	assert.EqualValues(t, "oauth2_application",
		GetTableName(LintGonicMapper, reflect.ValueOf(OAuth2Application{})))
}

type MyTable struct {
	Idx int
}

func (t *MyTable) TableName() string {
	return fmt.Sprintf("mytable_%d", t.Idx)
}

func TestMyTable(t *testing.T) {
	var table MyTable
	for i := 0; i < 10; i++ {
		table.Idx = i
		assert.EqualValues(t, fmt.Sprintf("mytable_%d", i), GetTableName(SameMapper{}, reflect.ValueOf(&table)))
	}
}
