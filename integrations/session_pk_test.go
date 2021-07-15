// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrations

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xormplus/xorm/schemas"
)

type IntId struct {
	Id   int `xorm:"pk autoincr"`
	Name string
}

type Int16Id struct {
	Id   int16 `xorm:"pk autoincr"`
	Name string
}

type Int32Id struct {
	Id   int32 `xorm:"pk autoincr"`
	Name string
}

type UintId struct {
	Id   uint `xorm:"pk autoincr"`
	Name string
}

type Uint16Id struct {
	Id   uint16 `xorm:"pk autoincr"`
	Name string
}

type Uint32Id struct {
	Id   uint32 `xorm:"pk autoincr"`
	Name string
}

type Uint64Id struct {
	Id   uint64 `xorm:"pk autoincr"`
	Name string
}

type StringPK struct {
	Id   string `xorm:"pk notnull"`
	Name string
}

type ID int64
type MyIntPK struct {
	ID   ID `xorm:"pk autoincr"`
	Name string
}

type StrID string
type MyStringPK struct {
	ID   StrID `xorm:"pk notnull"`
	Name string
}

func TestIntId(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&IntId{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&IntId{})
	assert.NoError(t, err)

	cnt, err := testEngine.Insert(&IntId{Name: "test"})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	bean := new(IntId)
	has, err := testEngine.Get(bean)
	assert.NoError(t, err)
	assert.True(t, has)

	beans := make([]IntId, 0)
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans))

	beans2 := make(map[int]IntId)
	err = testEngine.Find(&beans2)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans2))

	cnt, err = testEngine.ID(bean.Id).Delete(&IntId{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestInt16Id(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&Int16Id{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&Int16Id{})
	assert.NoError(t, err)

	cnt, err := testEngine.Insert(&Int16Id{Name: "test"})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	bean := new(Int16Id)
	has, err := testEngine.Get(bean)
	assert.NoError(t, err)
	assert.True(t, has)

	beans := make([]Int16Id, 0)
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans))

	beans2 := make(map[int16]Int16Id, 0)
	err = testEngine.Find(&beans2)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans2))

	cnt, err = testEngine.ID(bean.Id).Delete(&Int16Id{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestInt32Id(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&Int32Id{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&Int32Id{})
	assert.NoError(t, err)

	cnt, err := testEngine.Insert(&Int32Id{Name: "test"})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	bean := new(Int32Id)
	has, err := testEngine.Get(bean)
	assert.NoError(t, err)
	assert.True(t, has)

	beans := make([]Int32Id, 0)
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans))

	beans2 := make(map[int32]Int32Id, 0)
	err = testEngine.Find(&beans2)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans2))

	cnt, err = testEngine.ID(bean.Id).Delete(&Int32Id{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestUintId(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&UintId{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&UintId{})
	assert.NoError(t, err)

	cnt, err := testEngine.Insert(&UintId{Name: "test"})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var inserts = []UintId{
		{Name: "test1"},
		{Name: "test2"},
	}
	cnt, err = testEngine.Insert(&inserts)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cnt)

	bean := new(UintId)
	has, err := testEngine.Get(bean)
	assert.NoError(t, err)
	assert.True(t, has)

	beans := make([]UintId, 0)
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, len(beans))

	beans2 := make(map[uint]UintId, 0)
	err = testEngine.Find(&beans2)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, len(beans2))

	cnt, err = testEngine.ID(bean.Id).Delete(&UintId{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestUint16Id(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&Uint16Id{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&Uint16Id{})
	assert.NoError(t, err)

	cnt, err := testEngine.Insert(&Uint16Id{Name: "test"})
	assert.NoError(t, err)

	assert.EqualValues(t, 1, cnt)

	bean := new(Uint16Id)
	has, err := testEngine.Get(bean)
	assert.NoError(t, err)
	assert.True(t, has)

	beans := make([]Uint16Id, 0)
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans))

	beans2 := make(map[uint16]Uint16Id, 0)
	err = testEngine.Find(&beans2)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans2))

	cnt, err = testEngine.ID(bean.Id).Delete(&Uint16Id{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestUint32Id(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&Uint32Id{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&Uint32Id{})
	assert.NoError(t, err)

	cnt, err := testEngine.Insert(&Uint32Id{Name: "test"})
	assert.NoError(t, err)

	assert.EqualValues(t, 1, cnt)

	bean := new(Uint32Id)
	has, err := testEngine.Get(bean)
	assert.NoError(t, err)
	assert.True(t, has)

	beans := make([]Uint32Id, 0)
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans))

	beans2 := make(map[uint32]Uint32Id, 0)
	err = testEngine.Find(&beans2)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans2))

	cnt, err = testEngine.ID(bean.Id).Delete(&Uint32Id{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestUint64Id(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&Uint64Id{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&Uint64Id{})
	assert.NoError(t, err)

	idbean := &Uint64Id{Name: "test"}
	cnt, err := testEngine.Insert(idbean)
	assert.NoError(t, err)

	assert.EqualValues(t, 1, cnt)

	bean := new(Uint64Id)
	has, err := testEngine.Get(bean)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, bean.Id, idbean.Id)

	beans := make([]Uint64Id, 0)
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans))
	assert.EqualValues(t, *bean, beans[0])

	beans2 := make(map[uint64]Uint64Id, 0)
	err = testEngine.Find(&beans2)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans2))
	assert.EqualValues(t, *bean, beans2[bean.Id])

	cnt, err = testEngine.ID(bean.Id).Delete(&Uint64Id{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestStringPK(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&StringPK{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&StringPK{})
	assert.NoError(t, err)

	cnt, err := testEngine.Insert(&StringPK{Id: "1-1-2", Name: "test"})
	assert.NoError(t, err)

	assert.EqualValues(t, 1, cnt)

	bean := new(StringPK)
	has, err := testEngine.Get(bean)
	assert.NoError(t, err)
	assert.True(t, has)

	beans := make([]StringPK, 0)
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans))

	beans2 := make(map[string]StringPK)
	err = testEngine.Find(&beans2)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans2))

	cnt, err = testEngine.ID(bean.Id).Delete(&StringPK{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

type CompositeKey struct {
	Id1       int64 `xorm:"id1 pk"`
	Id2       int64 `xorm:"id2 pk"`
	UpdateStr string
}

func TestCompositeKey(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&CompositeKey{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&CompositeKey{})
	assert.NoError(t, err)

	cnt, err := testEngine.Insert(&CompositeKey{11, 22, ""})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Insert(&CompositeKey{11, 22, ""})
	assert.Error(t, err)
	assert.NotEqual(t, int64(1), cnt)

	var compositeKeyVal CompositeKey
	has, err := testEngine.ID(schemas.PK{11, 22}).Get(&compositeKeyVal)
	assert.NoError(t, err)
	assert.True(t, has)

	var compositeKeyVal2 CompositeKey
	// test passing PK ptr, this test seem failed withCache
	has, err = testEngine.ID(&schemas.PK{11, 22}).Get(&compositeKeyVal2)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, compositeKeyVal, compositeKeyVal2)

	var cps = make([]CompositeKey, 0)
	err = testEngine.Find(&cps)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(cps))
	assert.EqualValues(t, cps[0], compositeKeyVal)

	cnt, err = testEngine.Insert(&CompositeKey{22, 22, ""})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cps = make([]CompositeKey, 0)
	err = testEngine.Find(&cps)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(cps), "should has two record")
	assert.EqualValues(t, compositeKeyVal, cps[0], "should be equeal")

	compositeKeyVal = CompositeKey{UpdateStr: "test1"}
	cnt, err = testEngine.ID(schemas.PK{11, 22}).Update(&compositeKeyVal)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.ID(schemas.PK{11, 22}).Delete(&CompositeKey{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestCompositeKey2(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type User struct {
		UserId   string `xorm:"varchar(19) not null pk"`
		NickName string `xorm:"varchar(19) not null"`
		GameId   uint32 `xorm:"integer pk"`
		Score    int32  `xorm:"integer"`
	}

	err := testEngine.DropTables(&User{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&User{})
	assert.NoError(t, err)

	cnt, err := testEngine.Insert(&User{"11", "nick", 22, 5})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Insert(&User{"11", "nick", 22, 6})
	assert.Error(t, err)
	assert.NotEqual(t, 1, cnt)

	var user User
	has, err := testEngine.ID(schemas.PK{"11", 22}).Get(&user)
	assert.NoError(t, err)
	assert.True(t, has)

	// test passing PK ptr, this test seem failed withCache
	has, err = testEngine.ID(&schemas.PK{"11", 22}).Get(&user)
	assert.NoError(t, err)
	assert.True(t, has)

	user = User{NickName: "test1"}
	cnt, err = testEngine.ID(schemas.PK{"11", 22}).Update(&user)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.ID(schemas.PK{"11", 22}).Delete(&User{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

type MyString string
type UserPK2 struct {
	UserId   MyString `xorm:"varchar(19) not null pk"`
	NickName string   `xorm:"varchar(19) not null"`
	GameId   uint32   `xorm:"integer pk"`
	Score    int32    `xorm:"integer"`
}

func TestCompositeKey3(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&UserPK2{})

	assert.NoError(t, err)

	err = testEngine.CreateTables(&UserPK2{})
	assert.NoError(t, err)

	cnt, err := testEngine.Insert(&UserPK2{"11", "nick", 22, 5})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Insert(&UserPK2{"11", "nick", 22, 6})
	assert.Error(t, err)
	assert.NotEqual(t, 1, cnt)

	var user UserPK2
	has, err := testEngine.ID(schemas.PK{"11", 22}).Get(&user)
	assert.NoError(t, err)
	assert.True(t, has)

	// test passing PK ptr, this test seem failed withCache
	has, err = testEngine.ID(&schemas.PK{"11", 22}).Get(&user)
	assert.NoError(t, err)
	assert.True(t, has)

	user = UserPK2{NickName: "test1"}
	cnt, err = testEngine.ID(schemas.PK{"11", 22}).Update(&user)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.ID(schemas.PK{"11", 22}).Delete(&UserPK2{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestMyIntId(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&MyIntPK{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&MyIntPK{})
	assert.NoError(t, err)

	idbean := &MyIntPK{Name: "test"}
	cnt, err := testEngine.Insert(idbean)
	assert.NoError(t, err)

	assert.EqualValues(t, 1, cnt)

	bean := new(MyIntPK)
	has, err := testEngine.Get(bean)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, bean.ID, idbean.ID)

	var beans []MyIntPK
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans))
	assert.EqualValues(t, *bean, beans[0])

	beans2 := make(map[ID]MyIntPK, 0)
	err = testEngine.Find(&beans2)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans2))
	assert.EqualValues(t, *bean, beans2[bean.ID])

	cnt, err = testEngine.ID(bean.ID).Delete(&MyIntPK{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestMyStringId(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&MyStringPK{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&MyStringPK{})
	assert.NoError(t, err)

	idbean := &MyStringPK{ID: "1111", Name: "test"}
	cnt, err := testEngine.Insert(idbean)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	bean := new(MyStringPK)
	has, err := testEngine.Get(bean)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, bean.ID, idbean.ID)

	var beans []MyStringPK
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans))
	assert.EqualValues(t, *bean, beans[0])

	beans2 := make(map[StrID]MyStringPK, 0)
	err = testEngine.Find(&beans2)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(beans2))
	assert.EqualValues(t, *bean, beans2[bean.ID])

	cnt, err = testEngine.ID(bean.ID).Delete(&MyStringPK{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestSingleAutoIncrColumn(t *testing.T) {
	type Account struct {
		Id int64 `xorm:"pk autoincr"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Account))

	_, err := testEngine.Insert(&Account{})
	assert.NoError(t, err)
}

func TestCompositePK(t *testing.T) {
	type TaskSolution struct {
		UID     string    `xorm:"notnull pk UUID 'uid'"`
		TID     string    `xorm:"notnull pk UUID 'tid'"`
		Created time.Time `xorm:"created"`
		Updated time.Time `xorm:"updated"`
	}

	assert.NoError(t, PrepareEngine())

	tables1, err := testEngine.DBMetas()
	assert.NoError(t, err)

	assertSync(t, new(TaskSolution))
	assert.NoError(t, testEngine.Sync2(new(TaskSolution)))

	tables2, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1+len(tables1), len(tables2))

	var table *schemas.Table
	for _, t := range tables2 {
		if t.Name == testEngine.GetTableMapper().Obj2Table("TaskSolution") {
			table = t
			break
		}
	}

	assert.NotEqual(t, nil, table)

	pkCols := table.PKColumns()
	assert.EqualValues(t, 2, len(pkCols))

	names := []string{pkCols[0].Name, pkCols[1].Name}
	sort.Strings(names)
	assert.EqualValues(t, []string{"tid", "uid"}, names)
}

func TestNoPKIdQueryUpdate(t *testing.T) {
	type NoPKTable struct {
		Username string
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(NoPKTable))

	cnt, err := testEngine.Insert(&NoPKTable{
		Username: "test",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var res NoPKTable
	has, err := testEngine.ID("test").Get(&res)
	assert.Error(t, err)
	assert.False(t, has)

	cnt, err = testEngine.ID("test").Update(&NoPKTable{
		Username: "test1",
	})
	assert.Error(t, err)
	assert.EqualValues(t, 0, cnt)

	type UnvalidPKTable struct {
		ID       int `xorm:"id"`
		Username string
	}

	assertSync(t, new(UnvalidPKTable))

	cnt, err = testEngine.Insert(&UnvalidPKTable{
		ID:       1,
		Username: "test",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var res2 UnvalidPKTable
	has, err = testEngine.ID(1).Get(&res2)
	assert.Error(t, err)
	assert.False(t, has)

	cnt, err = testEngine.ID(1).Update(&UnvalidPKTable{
		Username: "test1",
	})
	assert.Error(t, err)
	assert.EqualValues(t, 0, cnt)
}
