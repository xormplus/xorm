// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrations

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xormplus/xorm"
	"github.com/xormplus/xorm/internal/statements"
	"github.com/xormplus/xorm/internal/utils"
	"github.com/xormplus/xorm/names"
)

func TestUpdateMap(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UpdateTable struct {
		Id   int64
		Name string
		Age  int
	}

	assert.NoError(t, testEngine.Sync2(new(UpdateTable)))
	var tb = UpdateTable{
		Name: "test",
		Age:  35,
	}
	_, err := testEngine.Insert(&tb)
	assert.NoError(t, err)

	cnt, err := testEngine.Table("update_table").Where("id = ?", tb.Id).Update(map[string]interface{}{
		"name": "test2",
		"age":  36,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Table("update_table").ID(tb.Id).Update(map[string]interface{}{
		"name": "test2",
		"age":  36,
	})
	assert.Error(t, err)
	assert.True(t, statements.IsIDConditionWithNoTableErr(err))
	assert.EqualValues(t, 0, cnt)
}

func TestUpdateLimit(t *testing.T) {
	if *ingoreUpdateLimit {
		t.Skip()
		return
	}

	assert.NoError(t, PrepareEngine())

	type UpdateTable2 struct {
		Id   int64
		Name string
		Age  int
	}

	assert.NoError(t, testEngine.Sync2(new(UpdateTable2)))
	var tb = UpdateTable2{
		Name: "test1",
		Age:  35,
	}
	cnt, err := testEngine.Insert(&tb)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	tb.Name = "test2"
	tb.Id = 0
	cnt, err = testEngine.Insert(&tb)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.OrderBy("name desc").Limit(1).Update(&UpdateTable2{
		Age: 30,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var uts []UpdateTable2
	err = testEngine.Find(&uts)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(uts))
	assert.EqualValues(t, 35, uts[0].Age)
	assert.EqualValues(t, 30, uts[1].Age)
}

type ForUpdate struct {
	Id   int64 `xorm:"pk"`
	Name string
}

func setupForUpdate(engine xorm.EngineInterface) error {
	v := new(ForUpdate)
	err := testEngine.DropTables(v)
	if err != nil {
		return err
	}
	err = testEngine.CreateTables(v)
	if err != nil {
		return err
	}

	list := []ForUpdate{
		{1, "data1"},
		{2, "data2"},
		{3, "data3"},
	}

	for _, f := range list {
		_, err = testEngine.Insert(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestForUpdate(t *testing.T) {
	if *ignoreSelectUpdate {
		return
	}

	err := setupForUpdate(testEngine)
	if err != nil {
		t.Error(err)
		return
	}

	session1 := testEngine.NewSession()
	session2 := testEngine.NewSession()
	session3 := testEngine.NewSession()
	defer session1.Close()
	defer session2.Close()
	defer session3.Close()

	// start transaction
	err = session1.Begin()
	if err != nil {
		t.Error(err)
		return
	}

	// use lock
	fList := make([]ForUpdate, 0)
	session1.ForUpdate()
	session1.Where("id = ?", 1)
	err = session1.Find(&fList)
	switch {
	case err != nil:
		t.Error(err)
		return
	case len(fList) != 1:
		t.Errorf("find not returned single row")
		return
	case fList[0].Name != "data1":
		t.Errorf("for_update.name must be `data1`")
		return
	}

	// wait for lock
	wg := &sync.WaitGroup{}

	// lock is used
	wg.Add(1)
	go func() {
		f2 := new(ForUpdate)
		session2.Where("id = ?", 1).ForUpdate()
		has, err := session2.Get(f2) // wait release lock
		switch {
		case err != nil:
			t.Error(err)
		case !has:
			t.Errorf("cannot find target row. for_update.id = 1")
		case f2.Name != "updated by session1":
			t.Errorf("read lock failed")
		}
		wg.Done()
	}()

	// lock is NOT used
	wg.Add(1)

	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		f3 := new(ForUpdate)
		session3.Where("id = ?", 1)
		has, err := session3.Get(f3) // wait release lock
		switch {
		case err != nil:
			t.Error(err)
		case !has:
			t.Errorf("cannot find target row. for_update.id = 1")
		case f3.Name != "data1":
			t.Errorf("read lock failed")
		}
		wg.Done()
		wg2.Done()
	}()

	wg2.Wait()

	f := new(ForUpdate)
	f.Name = "updated by session1"
	session1.Where("id = ?", 1)
	session1.Update(f)

	// release lock
	err = session1.Commit()
	if err != nil {
		t.Error(err)
		return
	}

	wg.Wait()
}

func TestWithIn(t *testing.T) {
	type temp3 struct {
		Id   int64  `xorm:"Id pk autoincr"`
		Name string `xorm:"Name"`
		Test bool   `xorm:"Test"`
	}

	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Sync(new(temp3)))

	testEngine.Insert(&[]temp3{
		{
			Name: "user1",
		},
		{
			Name: "user1",
		},
		{
			Name: "user1",
		},
	})

	cnt, err := testEngine.In("Id", 1, 2, 3, 4).Update(&temp3{Name: "aa"}, &temp3{Name: "user1"})
	assert.NoError(t, err)
	assert.EqualValues(t, 3, cnt)
}

type Condi map[string]interface{}

type UpdateAllCols struct {
	Id     int64
	Bool   bool
	String string
	Ptr    *string
}

type UpdateMustCols struct {
	Id     int64
	Bool   bool
	String string
}

type UpdateIncr struct {
	Id   int64
	Cnt  int
	Name string
}

type Article struct {
	Id      int32  `xorm:"pk INT autoincr"`
	Name    string `xorm:"VARCHAR(45)"`
	Img     string `xorm:"VARCHAR(100)"`
	Aside   string `xorm:"VARCHAR(200)"`
	Desc    string `xorm:"VARCHAR(200)"`
	Content string `xorm:"TEXT"`
	Status  int8   `xorm:"TINYINT(4)"`
}

func TestUpdateMap2(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(UpdateMustCols))

	_, err := testEngine.Table("update_must_cols").Where("id =?", 1).Update(map[string]interface{}{
		"bool": true,
	})
	assert.NoError(t, err)
}

func TestUpdate1(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	_, err := testEngine.Insert(&Userinfo{
		Username: "user1",
	})

	var ori Userinfo
	has, err := testEngine.Get(&ori)
	assert.NoError(t, err)
	assert.True(t, has)

	// update by id
	user := Userinfo{Username: "xxx", Height: 1.2}
	cnt, err := testEngine.ID(ori.Uid).Update(&user)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	condi := Condi{"username": "zzz", "departname": ""}
	cnt, err = testEngine.Table(&user).ID(ori.Uid).Update(&condi)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Update(&Userinfo{Username: "yyy"}, &user)
	assert.NoError(t, err)

	total, err := testEngine.Count(&user)
	assert.NoError(t, err)
	assert.EqualValues(t, cnt, total)

	// nullable update
	{
		user := &Userinfo{Username: "not null data", Height: 180.5}
		_, err := testEngine.Insert(user)
		assert.NoError(t, err)
		userID := user.Uid

		has, err := testEngine.ID(userID).
			And("username = ?", user.Username).
			And("height = ?", user.Height).
			And("departname = ?", "").
			And("detail_id = ?", 0).
			And("is_man = ?", 0).
			Get(&Userinfo{})
		assert.NoError(t, err)
		assert.True(t, has, "cannot insert properly")

		updatedUser := &Userinfo{Username: "null data"}
		cnt, err = testEngine.ID(userID).
			Nullable("height", "departname", "is_man", "created").
			Update(updatedUser)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, cnt, "update not returned 1")

		has, err = testEngine.ID(userID).
			And("username = ?", updatedUser.Username).
			And("height IS NULL").
			And("departname IS NULL").
			And("is_man IS NULL").
			And("created IS NULL").
			And("detail_id = ?", 0).
			Get(&Userinfo{})
		assert.NoError(t, err)
		assert.True(t, has, "cannot update with null properly")

		cnt, err = testEngine.ID(userID).Delete(&Userinfo{})
		assert.NoError(t, err)
		assert.EqualValues(t, 1, cnt, "delete not returned 1")
	}

	err = testEngine.StoreEngine("Innodb").Sync2(&Article{})
	assert.NoError(t, err)

	defer func() {
		err = testEngine.DropTables(&Article{})
		assert.NoError(t, err)
	}()

	a := &Article{0, "1", "2", "3", "4", "5", 2}
	cnt, err = testEngine.Insert(a)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt, fmt.Sprintf("insert not returned 1 but %d", cnt))
	assert.Greater(t, a.Id, int32(0), "insert returned id is 0")

	cnt, err = testEngine.ID(a.Id).Update(&Article{Name: "6"})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var s = "test"

	col1 := &UpdateAllCols{Ptr: &s}
	err = testEngine.Sync(col1)
	assert.NoError(t, err)

	_, err = testEngine.Insert(col1)
	assert.NoError(t, err)

	col2 := &UpdateAllCols{col1.Id, true, "", nil}
	_, err = testEngine.ID(col2.Id).AllCols().Update(col2)
	assert.NoError(t, err)

	col3 := &UpdateAllCols{}
	has, err = testEngine.ID(col2.Id).Get(col3)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, *col2, *col3)

	{
		col1 := &UpdateMustCols{}
		err = testEngine.Sync(col1)
		assert.NoError(t, err)

		_, err = testEngine.Insert(col1)
		assert.NoError(t, err)

		col2 := &UpdateMustCols{col1.Id, true, ""}
		boolStr := testEngine.GetColumnMapper().Obj2Table("Bool")
		stringStr := testEngine.GetColumnMapper().Obj2Table("String")
		_, err = testEngine.ID(col2.Id).MustCols(boolStr, stringStr).Update(col2)
		assert.NoError(t, err)

		col3 := &UpdateMustCols{}
		has, err := testEngine.ID(col2.Id).Get(col3)
		assert.NoError(t, err)
		assert.True(t, has)
		assert.EqualValues(t, *col2, *col3)
	}
}

func TestUpdateIncrDecr(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	col1 := &UpdateIncr{
		Name: "test",
	}
	assert.NoError(t, testEngine.Sync(col1))

	_, err := testEngine.Insert(col1)
	assert.NoError(t, err)

	colName := testEngine.GetColumnMapper().Obj2Table("Cnt")

	cnt, err := testEngine.ID(col1.Id).Incr(colName).Update(col1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	newCol := new(UpdateIncr)
	has, err := testEngine.ID(col1.Id).Get(newCol)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, 1, newCol.Cnt)

	cnt, err = testEngine.ID(col1.Id).Decr(colName).Update(col1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	newCol = new(UpdateIncr)
	has, err = testEngine.ID(col1.Id).Get(newCol)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, 0, newCol.Cnt)

	cnt, err = testEngine.ID(col1.Id).Cols(colName).Incr(colName).Update(col1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

type UpdatedUpdate struct {
	Id      int64
	Updated time.Time `xorm:"updated"`
}

type UpdatedUpdate2 struct {
	Id      int64
	Updated int64 `xorm:"updated"`
}

type UpdatedUpdate3 struct {
	Id      int64
	Updated int `xorm:"updated bigint"`
}

type UpdatedUpdate4 struct {
	Id      int64
	Updated int `xorm:"updated"`
}

type UpdatedUpdate5 struct {
	Id      int64
	Updated time.Time `xorm:"updated bigint"`
}

func TestUpdateUpdated(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	di := new(UpdatedUpdate)
	err := testEngine.Sync2(di)
	assert.NoError(t, err)

	_, err = testEngine.Insert(&UpdatedUpdate{})
	assert.NoError(t, err)

	ci := &UpdatedUpdate{}
	_, err = testEngine.ID(1).Update(ci)
	assert.NoError(t, err)

	has, err := testEngine.ID(1).Get(di)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, ci.Updated.Unix(), di.Updated.Unix())

	di2 := new(UpdatedUpdate2)
	err = testEngine.Sync2(di2)
	assert.NoError(t, err)

	now := time.Now()
	var di20 UpdatedUpdate2
	cnt, err := testEngine.Insert(&di20)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	assert.True(t, now.Unix() <= di20.Updated)

	var di21 UpdatedUpdate2
	has, err = testEngine.ID(di20.Id).Get(&di21)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, di20.Updated, di21.Updated)

	ci2 := &UpdatedUpdate2{}
	_, err = testEngine.ID(1).Update(ci2)
	assert.NoError(t, err)

	has, err = testEngine.ID(1).Get(di2)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, ci2.Updated, di2.Updated)
	assert.True(t, ci2.Updated >= di21.Updated)

	di3 := new(UpdatedUpdate3)
	err = testEngine.Sync2(di3)
	assert.NoError(t, err)

	_, err = testEngine.Insert(&UpdatedUpdate3{})
	assert.NoError(t, err)

	ci3 := &UpdatedUpdate3{}
	_, err = testEngine.ID(1).Update(ci3)
	assert.NoError(t, err)

	has, err = testEngine.ID(1).Get(di3)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, ci3.Updated, di3.Updated)

	di4 := new(UpdatedUpdate4)
	err = testEngine.Sync2(di4)
	assert.NoError(t, err)

	_, err = testEngine.Insert(&UpdatedUpdate4{})
	assert.NoError(t, err)

	ci4 := &UpdatedUpdate4{}
	_, err = testEngine.ID(1).Update(ci4)
	assert.NoError(t, err)

	has, err = testEngine.ID(1).Get(di4)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, ci4.Updated, di4.Updated)

	di5 := new(UpdatedUpdate5)
	err = testEngine.Sync2(di5)
	assert.NoError(t, err)

	_, err = testEngine.Insert(&UpdatedUpdate5{})
	assert.NoError(t, err)

	ci5 := &UpdatedUpdate5{}
	_, err = testEngine.ID(1).Update(ci5)
	assert.NoError(t, err)

	has, err = testEngine.ID(1).Get(di5)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, ci5.Updated.Unix(), di5.Updated.Unix())
}

func TestUpdateSameMapper(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	oldMapper := testEngine.GetTableMapper()
	testEngine.UnMapType(utils.ReflectValue(new(Userinfo)).Type())
	testEngine.UnMapType(utils.ReflectValue(new(Condi)).Type())
	testEngine.UnMapType(utils.ReflectValue(new(Article)).Type())
	testEngine.UnMapType(utils.ReflectValue(new(UpdateAllCols)).Type())
	testEngine.UnMapType(utils.ReflectValue(new(UpdateMustCols)).Type())
	testEngine.UnMapType(utils.ReflectValue(new(UpdateIncr)).Type())
	testEngine.SetMapper(names.SameMapper{})
	defer func() {
		testEngine.UnMapType(utils.ReflectValue(new(Userinfo)).Type())
		testEngine.UnMapType(utils.ReflectValue(new(Condi)).Type())
		testEngine.UnMapType(utils.ReflectValue(new(Article)).Type())
		testEngine.UnMapType(utils.ReflectValue(new(UpdateAllCols)).Type())
		testEngine.UnMapType(utils.ReflectValue(new(UpdateMustCols)).Type())
		testEngine.UnMapType(utils.ReflectValue(new(UpdateIncr)).Type())
		testEngine.SetMapper(oldMapper)
	}()

	assertSync(t, new(Userinfo))

	_, err := testEngine.Insert(&Userinfo{
		Username: "user1",
	})
	assert.NoError(t, err)

	var ori Userinfo
	has, err := testEngine.Get(&ori)
	assert.NoError(t, err)
	assert.True(t, has)

	// update by id
	user := Userinfo{Username: "xxx", Height: 1.2}
	cnt, err := testEngine.ID(ori.Uid).Update(&user)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	condi := Condi{"Username": "zzz", "Departname": ""}
	cnt, err = testEngine.Table(&user).ID(ori.Uid).Update(&condi)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Update(&Userinfo{Username: "yyy"}, &user)
	assert.NoError(t, err)

	total, err := testEngine.Count(&user)
	assert.NoError(t, err)
	assert.EqualValues(t, cnt, total)

	err = testEngine.Sync(&Article{})
	assert.NoError(t, err)

	defer func() {
		err = testEngine.DropTables(&Article{})
		assert.NoError(t, err)
	}()

	a := &Article{0, "1", "2", "3", "4", "5", 2}
	cnt, err = testEngine.Insert(a)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	assert.Greater(t, a.Id, int32(0))

	cnt, err = testEngine.ID(a.Id).Update(&Article{Name: "6"})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	col1 := &UpdateAllCols{}
	err = testEngine.Sync(col1)
	assert.NoError(t, err)

	_, err = testEngine.Insert(col1)
	assert.NoError(t, err)

	col2 := &UpdateAllCols{col1.Id, true, "", nil}
	_, err = testEngine.ID(col2.Id).AllCols().Update(col2)
	assert.NoError(t, err)

	col3 := &UpdateAllCols{}
	has, err = testEngine.ID(col2.Id).Get(col3)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, *col2, *col3)

	{
		col1 := &UpdateMustCols{}
		err = testEngine.Sync(col1)
		assert.NoError(t, err)

		_, err = testEngine.Insert(col1)
		assert.NoError(t, err)

		col2 := &UpdateMustCols{col1.Id, true, ""}
		boolStr := testEngine.GetColumnMapper().Obj2Table("Bool")
		stringStr := testEngine.GetColumnMapper().Obj2Table("String")
		_, err = testEngine.ID(col2.Id).MustCols(boolStr, stringStr).Update(col2)
		assert.NoError(t, err)

		col3 := &UpdateMustCols{}
		has, err := testEngine.ID(col2.Id).Get(col3)
		assert.NoError(t, err)
		assert.True(t, has)
		assert.EqualValues(t, *col2, *col3)
	}

	{
		col1 := &UpdateIncr{}
		err = testEngine.Sync(col1)
		assert.NoError(t, err)

		_, err = testEngine.Insert(col1)
		assert.NoError(t, err)

		cnt, err := testEngine.ID(col1.Id).Incr("`Cnt`").Update(col1)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, cnt)

		newCol := new(UpdateIncr)
		has, err := testEngine.ID(col1.Id).Get(newCol)
		assert.NoError(t, err)
		assert.True(t, has)
		assert.EqualValues(t, 1, newCol.Cnt)
	}
}

func TestUseBool(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	cnt1, err := testEngine.Count(&Userinfo{})
	assert.NoError(t, err)

	users := make([]Userinfo, 0)
	err = testEngine.Find(&users)
	assert.NoError(t, err)
	var fNumber int64
	for _, u := range users {
		if u.IsMan == false {
			fNumber++
		}
	}

	cnt2, err := testEngine.UseBool().Update(&Userinfo{IsMan: true})
	assert.NoError(t, err)
	if fNumber != cnt2 {
		fmt.Println("cnt1", cnt1, "fNumber", fNumber, "cnt2", cnt2)
		/*err = errors.New("Updated number is not corrected.")
		  t.Error(err)
		  panic(err)*/
	}

	_, err = testEngine.Update(&Userinfo{IsMan: true})
	assert.Error(t, err)
}

func TestBool(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	_, err := testEngine.UseBool().Update(&Userinfo{IsMan: true})
	assert.NoError(t, err)
	users := make([]Userinfo, 0)
	err = testEngine.Find(&users)
	assert.NoError(t, err)
	for _, user := range users {
		assert.True(t, user.IsMan)
	}

	_, err = testEngine.UseBool().Update(&Userinfo{IsMan: false})
	assert.NoError(t, err)
	users = make([]Userinfo, 0)
	err = testEngine.Find(&users)
	assert.NoError(t, err)
	for _, user := range users {
		assert.True(t, user.IsMan)
	}
}

func TestNoUpdate(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type NoUpdate struct {
		Id      int64
		Content string
	}

	assertSync(t, new(NoUpdate))

	cnt, err := testEngine.Insert(&NoUpdate{
		Content: "test",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	_, err = testEngine.ID(1).Update(&NoUpdate{})
	assert.Error(t, err)
	assert.EqualValues(t, "No content found to be updated", err.Error())
}

func TestNewUpdate(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type TbUserInfo struct {
		Id       int64       `xorm:"pk autoincr unique BIGINT" json:"id"`
		Phone    string      `xorm:"not null unique VARCHAR(20)" json:"phone"`
		UserName string      `xorm:"VARCHAR(20)" json:"user_name"`
		Gender   int         `xorm:"default 0 INTEGER" json:"gender"`
		Pw       string      `xorm:"VARCHAR(100)" json:"pw"`
		Token    string      `xorm:"TEXT" json:"token"`
		Avatar   string      `xorm:"TEXT" json:"avatar"`
		Extras   interface{} `xorm:"JSON" json:"extras"`
		Created  time.Time   `xorm:"DATETIME created"`
		Updated  time.Time   `xorm:"DATETIME updated"`
		Deleted  time.Time   `xorm:"DATETIME deleted"`
	}

	assertSync(t, new(TbUserInfo))

	targetUsr := TbUserInfo{Phone: "13126564922"}
	changeUsr := TbUserInfo{Token: "ABCDEFG"}
	af, err := testEngine.Update(&changeUsr, &targetUsr)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, af)

	af, err = testEngine.Table(new(TbUserInfo)).Where("phone=?", 13126564922).Update(&changeUsr)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, af)
}

func TestUpdateUpdate(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type PublicKeyUpdate struct {
		Id          int64
		UpdatedUnix int64 `xorm:"updated"`
	}

	assertSync(t, new(PublicKeyUpdate))

	cnt, err := testEngine.ID(1).Cols("updated_unix").Update(&PublicKeyUpdate{
		UpdatedUnix: time.Now().Unix(),
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 0, cnt)
}

func TestCreatedUpdated2(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type CreatedUpdatedStruct struct {
		Id       int64
		Name     string
		CreateAt time.Time `xorm:"created" json:"create_at"`
		UpdateAt time.Time `xorm:"updated" json:"update_at"`
	}

	assertSync(t, new(CreatedUpdatedStruct))

	var s = CreatedUpdatedStruct{
		Name: "test",
	}
	cnt, err := testEngine.Insert(&s)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	assert.EqualValues(t, s.UpdateAt.Unix(), s.CreateAt.Unix())

	time.Sleep(time.Second)

	var s1 = CreatedUpdatedStruct{
		Name:     "test1",
		CreateAt: s.CreateAt,
		UpdateAt: s.UpdateAt,
	}

	cnt, err = testEngine.ID(1).Update(&s1)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	assert.EqualValues(t, s.CreateAt.Unix(), s1.CreateAt.Unix())
	assert.True(t, s1.UpdateAt.Unix() > s.UpdateAt.Unix())

	var s2 CreatedUpdatedStruct
	has, err := testEngine.ID(1).Get(&s2)
	assert.NoError(t, err)
	assert.True(t, has)

	assert.EqualValues(t, s.CreateAt.Unix(), s2.CreateAt.Unix())
	assert.True(t, s2.UpdateAt.Unix() > s.UpdateAt.Unix())
	assert.True(t, s2.UpdateAt.Unix() > s2.CreateAt.Unix())
}

func TestDeletedUpdate(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type DeletedUpdatedStruct struct {
		Id        int64
		Name      string
		DeletedAt time.Time `xorm:"deleted"`
	}

	assertSync(t, new(DeletedUpdatedStruct))

	var s = DeletedUpdatedStruct{
		Name: "test",
	}
	cnt, err := testEngine.Insert(&s)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.ID(s.Id).Delete(&DeletedUpdatedStruct{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	s.DeletedAt = time.Time{}
	cnt, err = testEngine.Unscoped().Nullable("deleted_at").Update(&s)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var s1 DeletedUpdatedStruct
	has, err := testEngine.ID(s.Id).Get(&s1)
	assert.EqualValues(t, true, has)

	cnt, err = testEngine.ID(s.Id).Delete(&DeletedUpdatedStruct{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.ID(s.Id).Cols("deleted_at").Update(&DeletedUpdatedStruct{})
	assert.EqualValues(t, "No content found to be updated", err.Error())
	assert.EqualValues(t, 0, cnt)

	cnt, err = testEngine.ID(s.Id).Unscoped().Cols("deleted_at").Update(&DeletedUpdatedStruct{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var s2 DeletedUpdatedStruct
	has, err = testEngine.ID(s.Id).Get(&s2)
	assert.EqualValues(t, true, has)
}

func TestUpdateMapCondition(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UpdateMapCondition struct {
		Id     int64
		String string
	}

	assertSync(t, new(UpdateMapCondition))

	var c = UpdateMapCondition{
		String: "string",
	}
	_, err := testEngine.Insert(&c)
	assert.NoError(t, err)

	cnt, err := testEngine.Update(&UpdateMapCondition{
		String: "string1",
	}, map[string]interface{}{
		"id": c.Id,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var c2 UpdateMapCondition
	has, err := testEngine.ID(c.Id).Get(&c2)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, "string1", c2.String)
}

func TestUpdateMapContent(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UpdateMapContent struct {
		Id     int64
		Name   string
		IsMan  bool
		Age    int
		Gender int // 1 is man, 2 is woman
	}

	assertSync(t, new(UpdateMapContent))

	var c = UpdateMapContent{
		Name:   "lunny",
		IsMan:  true,
		Gender: 1,
		Age:    18,
	}
	_, err := testEngine.Insert(&c)
	assert.NoError(t, err)
	assert.EqualValues(t, 18, c.Age)

	cnt, err := testEngine.Table(new(UpdateMapContent)).ID(c.Id).Update(map[string]interface{}{"age": 0})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var c1 UpdateMapContent
	has, err := testEngine.ID(c.Id).Get(&c1)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, 0, c1.Age)

	cnt, err = testEngine.Table(new(UpdateMapContent)).ID(c.Id).Update(map[string]interface{}{
		"age":    16,
		"is_man": false,
		"gender": 2,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var c2 UpdateMapContent
	has, err = testEngine.ID(c.Id).Get(&c2)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, 16, c2.Age)
	assert.EqualValues(t, false, c2.IsMan)
	assert.EqualValues(t, 2, c2.Gender)

	cnt, err = testEngine.Table(new(UpdateMapContent)).ID(c.Id).Update(map[string]interface{}{
		"age":    15,
		"is_man": true,
		"gender": 1,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var c3 UpdateMapContent
	has, err = testEngine.ID(c.Id).Get(&c3)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, 15, c3.Age)
	assert.EqualValues(t, true, c3.IsMan)
	assert.EqualValues(t, 1, c3.Gender)
}

func TestUpdateCondiBean(t *testing.T) {
	type NeedUpdateBean struct {
		Id   int64
		Name string
	}

	type NeedUpdateCondiBean struct {
		Name string
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(NeedUpdateBean))

	cnt, err := testEngine.Insert(&NeedUpdateBean{
		Name: "name1",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	has, err := testEngine.Exist(&NeedUpdateBean{
		Name: "name1",
	})
	assert.NoError(t, err)
	assert.True(t, has)

	cnt, err = testEngine.Update(&NeedUpdateBean{
		Name: "name2",
	}, &NeedUpdateCondiBean{
		Name: "name1",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	has, err = testEngine.Exist(&NeedUpdateBean{
		Name: "name2",
	})
	assert.NoError(t, err)
	assert.True(t, has)

	cnt, err = testEngine.Update(&NeedUpdateBean{
		Name: "name1",
	}, NeedUpdateCondiBean{
		Name: "name2",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	has, err = testEngine.Exist(&NeedUpdateBean{
		Name: "name1",
	})
	assert.NoError(t, err)
	assert.True(t, has)
}

func TestWhereCondErrorWhenUpdate(t *testing.T) {
	type AuthRequestError struct {
		ChallengeToken string
		RequestToken   string
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(AuthRequestError))

	_, err := testEngine.Cols("challenge_token", "request_token", "challenge_agent", "status").
		Where(&AuthRequestError{ChallengeToken: "1"}).
		Update(&AuthRequestError{
			ChallengeToken: "2",
		})
	assert.Error(t, err)
	assert.EqualValues(t, xorm.ErrConditionType, err)
}

func TestUpdateDeleted(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UpdateDeletedStruct struct {
		Id        int64
		Name      string
		DeletedAt time.Time `xorm:"deleted"`
	}

	assertSync(t, new(UpdateDeletedStruct))

	var s = UpdateDeletedStruct{
		Name: "test",
	}
	cnt, err := testEngine.Insert(&s)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.ID(s.Id).Delete(&UpdateDeletedStruct{})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.ID(s.Id).Update(&UpdateDeletedStruct{
		Name: "test1",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 0, cnt)

	cnt, err = testEngine.Table(&UpdateDeletedStruct{}).ID(s.Id).Update(map[string]interface{}{
		"name": "test1",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 0, cnt)

	cnt, err = testEngine.ID(s.Id).Unscoped().Update(&UpdateDeletedStruct{
		Name: "test1",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestUpdateExprs(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UpdateExprs struct {
		Id        int64
		NumIssues int
		Name      string
	}

	assertSync(t, new(UpdateExprs))

	_, err := testEngine.Insert(&UpdateExprs{
		NumIssues: 1,
		Name:      "lunny",
	})
	assert.NoError(t, err)

	_, err = testEngine.SetExpr("num_issues", "num_issues+1").AllCols().Update(&UpdateExprs{
		NumIssues: 3,
		Name:      "lunny xiao",
	})
	assert.NoError(t, err)

	var ue UpdateExprs
	has, err := testEngine.Get(&ue)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, 2, ue.NumIssues)
	assert.EqualValues(t, "lunny xiao", ue.Name)
}

func TestUpdateAlias(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UpdateAlias struct {
		Id        int64
		NumIssues int
		Name      string
	}

	assertSync(t, new(UpdateAlias))

	_, err := testEngine.Insert(&UpdateAlias{
		NumIssues: 1,
		Name:      "lunny",
	})
	assert.NoError(t, err)

	_, err = testEngine.Alias("ua").Where("ua.id = ?", 1).Update(&UpdateAlias{
		NumIssues: 2,
		Name:      "lunny xiao",
	})
	assert.NoError(t, err)

	var ue UpdateAlias
	has, err := testEngine.Get(&ue)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, 2, ue.NumIssues)
	assert.EqualValues(t, "lunny xiao", ue.Name)
}

func TestUpdateExprs2(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UpdateExprsRelease struct {
		Id         int64
		RepoId     int
		IsTag      bool
		IsDraft    bool
		NumCommits int
		Sha1       string
	}

	assertSync(t, new(UpdateExprsRelease))

	var uer = UpdateExprsRelease{
		RepoId:     1,
		IsTag:      false,
		IsDraft:    false,
		NumCommits: 1,
		Sha1:       "sha1",
	}
	inserted, err := testEngine.Insert(&uer)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, inserted)

	updated, err := testEngine.
		Where("repo_id = ? AND is_tag = ?", 1, false).
		SetExpr("is_draft", true).
		SetExpr("num_commits", 0).
		SetExpr("sha1", "").
		Update(new(UpdateExprsRelease))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, updated)

	var uer2 UpdateExprsRelease
	has, err := testEngine.ID(uer.Id).Get(&uer2)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, 1, uer2.RepoId)
	assert.EqualValues(t, false, uer2.IsTag)
	assert.EqualValues(t, true, uer2.IsDraft)
	assert.EqualValues(t, 0, uer2.NumCommits)
	assert.EqualValues(t, "", uer2.Sha1)
}

func TestUpdateMap3(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UpdateMapUser struct {
		Id   uint64 `xorm:"PK autoincr"`
		Name string `xorm:""`
		Ver  uint64 `xorm:"version"`
	}

	oldMapper := testEngine.GetColumnMapper()
	defer func() {
		testEngine.SetColumnMapper(oldMapper)
	}()

	mapper := names.NewPrefixMapper(names.SnakeMapper{}, "F")
	testEngine.SetColumnMapper(mapper)

	assertSync(t, new(UpdateMapUser))

	_, err := testEngine.Table(new(UpdateMapUser)).Insert(map[string]interface{}{
		"Fname": "first user name",
		"Fver":  1,
	})
	assert.NoError(t, err)

	update := map[string]interface{}{
		"Fname": "user name",
		"Fver":  1,
	}
	rows, err := testEngine.Table(new(UpdateMapUser)).ID(1).Update(update)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, rows)

	update = map[string]interface{}{
		"Name": "user name",
		"Ver":  1,
	}
	rows, err = testEngine.Table(new(UpdateMapUser)).ID(1).Update(update)
	assert.Error(t, err)
	assert.EqualValues(t, 0, rows)
}

func TestUpdateIgnoreOnlyFromDBFields(t *testing.T) {
	type TestOnlyFromDBField struct {
		Id              int64  `xorm:"PK"`
		OnlyFromDBField string `xorm:"<-"`
		OnlyToDBField   string `xorm:"->"`
		IngoreField     string `xorm:"-"`
	}

	assertGetRecord := func() *TestOnlyFromDBField {
		var record TestOnlyFromDBField
		has, err := testEngine.Where("id = ?", 1).Get(&record)
		assert.NoError(t, err)
		assert.EqualValues(t, true, has)
		assert.EqualValues(t, "", record.OnlyFromDBField)
		return &record

	}
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(TestOnlyFromDBField))

	_, err := testEngine.Insert(&TestOnlyFromDBField{
		Id:              1,
		OnlyFromDBField: "a",
		OnlyToDBField:   "b",
		IngoreField:     "c",
	})
	assert.NoError(t, err)

	assertGetRecord()

	_, err = testEngine.ID(1).Update(&TestOnlyFromDBField{
		OnlyToDBField:   "b",
		OnlyFromDBField: "test",
	})
	assert.NoError(t, err)
	assertGetRecord()
}

func TestUpdateMultiplePK(t *testing.T) {
	type TestUpdateMultiplePKStruct struct {
		Id    string `xorm:"notnull pk" description:"唯一ID号"`
		Name  string `xorm:"notnull pk" description:"名称"`
		Value string `xorm:"notnull varchar(4000)" description:"值"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(TestUpdateMultiplePKStruct))

	test := &TestUpdateMultiplePKStruct{
		Id:    "ID1",
		Name:  "Name1",
		Value: "1",
	}
	_, err := testEngine.Insert(test)
	assert.NoError(t, err)

	test.Value = "2"
	_, err = testEngine.Where("`id` = ? And `name` = ?", test.Id, test.Name).Cols("Value").Update(test)
	assert.NoError(t, err)

	test.Value = "3"
	num, err := testEngine.Where("`id` = ? And `name` = ?", test.Id, test.Name).Update(test)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, num)

	test.Value = "4"
	_, err = testEngine.ID([]interface{}{test.Id, test.Name}).Update(test)
	assert.NoError(t, err)

	type MySlice []interface{}
	test.Value = "5"
	_, err = testEngine.ID(&MySlice{test.Id, test.Name}).Update(test)
	assert.NoError(t, err)
}

type TestFieldType1 struct {
	cb []byte
}

func (a *TestFieldType1) FromDB(src []byte) error {
	a.cb = src
	return nil
}

func (a TestFieldType1) ToDB() ([]byte, error) {
	return a.cb, nil
}

type TestTable1 struct {
	Id         int64
	Field1     *TestFieldType1 `xorm:"text"`
	UpdateTime time.Time
}

func TestNilFromDB(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(TestTable1))

	cnt, err := testEngine.Insert(&TestTable1{
		Field1: &TestFieldType1{
			cb: []byte("string"),
		},
		UpdateTime: time.Now(),
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Update(TestTable1{
		UpdateTime: time.Now().Add(time.Second),
	}, TestTable1{
		Id: 1,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Insert(&TestTable1{
		UpdateTime: time.Now(),
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}
