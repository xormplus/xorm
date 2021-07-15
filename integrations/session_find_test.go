// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrations

import (
	"testing"
	"time"

	"github.com/xormplus/xorm/internal/utils"
	"github.com/xormplus/xorm/names"

	"github.com/stretchr/testify/assert"
)

func TestJoinLimit(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type Salary struct {
		Id  int64
		Lid int64
	}

	type CheckList struct {
		Id  int64
		Eid int64
	}

	type Empsetting struct {
		Id   int64
		Name string
	}

	assert.NoError(t, testEngine.Sync2(new(Salary), new(CheckList), new(Empsetting)))

	var emp Empsetting
	cnt, err := testEngine.Insert(&emp)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var checklist = CheckList{
		Eid: emp.Id,
	}
	cnt, err = testEngine.Insert(&checklist)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var salary = Salary{
		Lid: checklist.Id,
	}
	cnt, err = testEngine.Insert(&salary)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var salaries []Salary
	err = testEngine.Table("salary").
		Join("INNER", "check_list", "check_list.id = salary.lid").
		Join("LEFT", "empsetting", "empsetting.id = check_list.eid").
		Limit(10, 0).
		Find(&salaries)
	assert.NoError(t, err)
}

func TestWhere(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.Where("id > ?", 2).Find(&users)
	assert.NoError(t, err)

	err = testEngine.Where("id > ?", 2).And("id < ?", 10).Find(&users)
	assert.NoError(t, err)
}

func TestFind(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)

	err := testEngine.Find(&users)
	assert.NoError(t, err)

	users2 := make([]Userinfo, 0)
	var tbName = testEngine.Quote(testEngine.TableName(new(Userinfo), true))
	err = testEngine.SQL("select * from " + tbName).Find(&users2)
	assert.NoError(t, err)
}

func TestFind2(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	users := make([]*Userinfo, 0)

	assertSync(t, new(Userinfo))

	err := testEngine.Find(&users)
	assert.NoError(t, err)
}

type Team struct {
	Id int64
}

type TeamUser struct {
	OrgId  int64
	Uid    int64
	TeamId int64
}

func (TeamUser) TableName() string {
	return "team_user"
}

func TestFind3(t *testing.T) {
	var teamUser = new(TeamUser)
	assert.NoError(t, PrepareEngine())
	err := testEngine.Sync2(new(Team), teamUser)
	assert.NoError(t, err)

	var teams []Team
	err = testEngine.Cols("`team`.id").
		Where("`team_user`.org_id=?", 1).
		And("`team_user`.uid=?", 2).
		Join("INNER", "`team_user`", "`team_user`.team_id=`team`.id").
		Find(&teams)
	assert.NoError(t, err)

	teams = make([]Team, 0)
	err = testEngine.Cols("`team`.id").
		Where("`team_user`.org_id=?", 1).
		And("`team_user`.uid=?", 2).
		Join("INNER", teamUser, "`team_user`.team_id=`team`.id").
		Find(&teams)
	assert.NoError(t, err)

	teams = make([]Team, 0)
	err = testEngine.Cols("`team`.id").
		Where("`team_user`.org_id=?", 1).
		And("`team_user`.uid=?", 2).
		Join("INNER", []interface{}{teamUser}, "`team_user`.team_id=`team`.id").
		Find(&teams)
	assert.NoError(t, err)

	teams = make([]Team, 0)
	err = testEngine.Cols("`team`.id").
		Where("`tu`.org_id=?", 1).
		And("`tu`.uid=?", 2).
		Join("INNER", []string{"team_user", "tu"}, "`tu`.team_id=`team`.id").
		Find(&teams)
	assert.NoError(t, err)

	teams = make([]Team, 0)
	err = testEngine.Cols("`team`.id").
		Where("`tu`.org_id=?", 1).
		And("`tu`.uid=?", 2).
		Join("INNER", []interface{}{"team_user", "tu"}, "`tu`.team_id=`team`.id").
		Find(&teams)
	assert.NoError(t, err)

	teams = make([]Team, 0)
	err = testEngine.Cols("`team`.id").
		Where("`tu`.org_id=?", 1).
		And("`tu`.uid=?", 2).
		Join("INNER", []interface{}{teamUser, "tu"}, "`tu`.team_id=`team`.id").
		Find(&teams)
	assert.NoError(t, err)
}

func TestFindMap(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	cnt, err := testEngine.Insert(&Userinfo{
		Username:   "lunny",
		Departname: "depart1",
		IsMan:      true,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	users := make(map[int64]Userinfo)
	err = testEngine.Find(&users)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(users))
	assert.EqualValues(t, "lunny", users[1].Username)
	assert.EqualValues(t, "depart1", users[1].Departname)
	assert.True(t, users[1].IsMan)

	users = make(map[int64]Userinfo)
	err = testEngine.Cols("username, departname").Find(&users)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(users))
	assert.EqualValues(t, "lunny", users[1].Username)
	assert.EqualValues(t, "depart1", users[1].Departname)
	assert.False(t, users[1].IsMan)
}

func TestFindMap2(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	users := make(map[int64]*Userinfo)
	err := testEngine.Find(&users)
	assert.NoError(t, err)
}

func TestDistinct(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	_, err := testEngine.Insert(&Userinfo{
		Username: "lunny",
	})
	assert.NoError(t, err)

	users := make([]Userinfo, 0)
	departname := testEngine.GetTableMapper().Obj2Table("Departname")
	err = testEngine.Distinct(departname).Find(&users)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(users))

	type Depart struct {
		Departname string
	}

	users2 := make([]Depart, 0)
	err = testEngine.Distinct(departname).Table(new(Userinfo)).Find(&users2)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(users2))
}

func TestOrder(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.OrderBy("id desc").Find(&users)
	assert.NoError(t, err)

	users2 := make([]Userinfo, 0)
	err = testEngine.Asc("id", "username").Desc("height").Find(&users2)
	assert.NoError(t, err)
}

func TestGroupBy(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.GroupBy("id, username").Find(&users)
	assert.NoError(t, err)
}

func TestHaving(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.GroupBy("username").Having("username='xlw'").Find(&users)
	assert.NoError(t, err)
}

func TestOrderSameMapper(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	testEngine.UnMapType(utils.ReflectValue(new(Userinfo)).Type())

	mapper := testEngine.GetTableMapper()
	testEngine.SetMapper(names.SameMapper{})

	defer func() {
		testEngine.UnMapType(utils.ReflectValue(new(Userinfo)).Type())
		testEngine.SetMapper(mapper)
	}()

	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.OrderBy("id desc").Find(&users)
	assert.NoError(t, err)

	users2 := make([]Userinfo, 0)
	err = testEngine.Asc("id", "Username").Desc("Height").Find(&users2)
	assert.NoError(t, err)
}

func TestHavingSameMapper(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	testEngine.UnMapType(utils.ReflectValue(new(Userinfo)).Type())

	mapper := testEngine.GetTableMapper()
	testEngine.SetMapper(names.SameMapper{})
	defer func() {
		testEngine.UnMapType(utils.ReflectValue(new(Userinfo)).Type())
		testEngine.SetMapper(mapper)
	}()
	assertSync(t, new(Userinfo))

	users := make([]Userinfo, 0)
	err := testEngine.GroupBy("`Username`").Having("`Username`='xlw'").Find(&users)
	assert.NoError(t, err)
}

func TestFindInts(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	var idsInt64 []int64
	err := testEngine.Table(userinfo).Cols("id").Desc("id").Find(&idsInt64)
	assert.NoError(t, err)

	var idsInt32 []int32
	err = testEngine.Table(userinfo).Cols("id").Desc("id").Find(&idsInt32)
	assert.NoError(t, err)

	var idsInt []int
	err = testEngine.Table(userinfo).Cols("id").Desc("id").Find(&idsInt)
	assert.NoError(t, err)

	var idsUint []uint
	err = testEngine.Table(userinfo).Cols("id").Desc("id").Find(&idsUint)
	assert.NoError(t, err)

	type MyInt int
	var idsMyInt []MyInt
	err = testEngine.Table(userinfo).Cols("id").Desc("id").Find(&idsMyInt)
	assert.NoError(t, err)
}

func TestFindStrings(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))
	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	username := testEngine.GetColumnMapper().Obj2Table("Username")
	var idsString []string
	err := testEngine.Table(userinfo).Cols(username).Desc("id").Find(&idsString)
	assert.NoError(t, err)
}

func TestFindMyString(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))
	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	username := testEngine.GetColumnMapper().Obj2Table("Username")

	var idsMyString []MyString
	err := testEngine.Table(userinfo).Cols(username).Desc("id").Find(&idsMyString)
	assert.NoError(t, err)
}

func TestFindInterface(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	username := testEngine.GetColumnMapper().Obj2Table("Username")
	var idsInterface []interface{}
	err := testEngine.Table(userinfo).Cols(username).Desc("id").Find(&idsInterface)
	assert.NoError(t, err)
}

func TestFindSliceBytes(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	var ids [][][]byte
	err := testEngine.Table(userinfo).Desc("id").Find(&ids)
	assert.NoError(t, err)
}

func TestFindSlicePtrString(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	var ids [][]*string
	err := testEngine.Table(userinfo).Desc("id").Find(&ids)
	assert.NoError(t, err)
}

func TestFindMapBytes(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	var ids []map[string][]byte
	err := testEngine.Table(userinfo).Desc("id").Find(&ids)
	assert.NoError(t, err)
}

func TestFindMapPtrString(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	userinfo := testEngine.GetTableMapper().Obj2Table("Userinfo")
	var ids []map[string]*string
	err := testEngine.Table(userinfo).Desc("id").Find(&ids)
	assert.NoError(t, err)
}

func TestFindBit(t *testing.T) {
	type FindBitStruct struct {
		Id  int64
		Msg bool `xorm:"bit"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(FindBitStruct))

	cnt, err := testEngine.Insert([]FindBitStruct{
		{
			Msg: false,
		},
		{
			Msg: true,
		},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cnt)

	var results = make([]FindBitStruct, 0, 2)
	err = testEngine.Find(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(results))
}

func TestFindMark(t *testing.T) {

	type Mark struct {
		Mark1 string `xorm:"VARCHAR(1)"`
		Mark2 string `xorm:"VARCHAR(1)"`
		MarkA string `xorm:"VARCHAR(1)"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Mark))

	cnt, err := testEngine.Insert([]Mark{
		{
			Mark1: "1",
			Mark2: "2",
			MarkA: "A",
		},
		{
			Mark1: "1",
			Mark2: "2",
			MarkA: "A",
		},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cnt)

	var results = make([]Mark, 0, 2)
	err = testEngine.Find(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(results))
}

func TestFindAndCountOneFunc(t *testing.T) {
	type FindAndCountStruct struct {
		Id      int64
		Content string
		Msg     bool `xorm:"bit"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(FindAndCountStruct))

	cnt, err := testEngine.Insert([]FindAndCountStruct{
		{
			Content: "111",
			Msg:     false,
		},
		{
			Content: "222",
			Msg:     true,
		},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cnt)

	var results = make([]FindAndCountStruct, 0, 2)
	cnt, err = testEngine.Limit(1).FindAndCount(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	assert.EqualValues(t, 2, cnt)

	results = make([]FindAndCountStruct, 0, 2)
	cnt, err = testEngine.FindAndCount(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(results))
	assert.EqualValues(t, 2, cnt)

	results = make([]FindAndCountStruct, 0, 1)
	cnt, err = testEngine.Where("msg = ?", true).FindAndCount(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	assert.EqualValues(t, 1, cnt)

	results = make([]FindAndCountStruct, 0, 1)
	cnt, err = testEngine.Where("1=1").Limit(1).FindAndCount(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	assert.EqualValues(t, 2, cnt)
	assert.EqualValues(t, FindAndCountStruct{
		Id:      1,
		Content: "111",
		Msg:     false,
	}, results[0])

	results = make([]FindAndCountStruct, 0, 1)
	cnt, err = testEngine.Where("1=1").Limit(1).FindAndCount(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	assert.EqualValues(t, 2, cnt)
	assert.EqualValues(t, FindAndCountStruct{
		Id:      1,
		Content: "111",
		Msg:     false,
	}, results[0])

	results = make([]FindAndCountStruct, 0, 1)
	cnt, err = testEngine.Where("1=1").Limit(1, 1).FindAndCount(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	assert.EqualValues(t, 2, cnt)
	assert.EqualValues(t, FindAndCountStruct{
		Id:      2,
		Content: "222",
		Msg:     true,
	}, results[0])

	results = make([]FindAndCountStruct, 0, 1)
	cnt, err = testEngine.Where("1=1").Limit(1, 1).FindAndCount(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	assert.EqualValues(t, 2, cnt)
	assert.EqualValues(t, FindAndCountStruct{
		Id:      2,
		Content: "222",
		Msg:     true,
	}, results[0])

	results = make([]FindAndCountStruct, 0, 1)
	cnt, err = testEngine.Where("msg = ?", true).Select("id, content, msg").
		Limit(1).FindAndCount(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	assert.EqualValues(t, 1, cnt)

	results = make([]FindAndCountStruct, 0, 1)
	cnt, err = testEngine.Where("msg = ?", true).Desc("id").
		Limit(1).Cols("content").FindAndCount(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	assert.EqualValues(t, 1, cnt)

	ids := make([]int64, 0, 2)
	tableName := testEngine.GetTableMapper().Obj2Table("FindAndCountStruct")
	cnt, err = testEngine.Table(tableName).Limit(1).Cols("id").FindAndCount(&ids)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(ids))
	assert.EqualValues(t, 2, cnt)
}

func TestFindAndCountOneFuncWithDeleted(t *testing.T) {
	type CommentWithDeleted struct {
		Id        int   `xorm:"pk autoincr"`
		DeletedAt int64 `xorm:"deleted notnull default(0) index"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(CommentWithDeleted))

	var comments []CommentWithDeleted
	cnt, err := testEngine.FindAndCount(&comments)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, cnt)
}

func TestFindAndCount2(t *testing.T) {
	// User
	type TestFindAndCountUser struct {
		Id   int64  `xorm:"bigint(11) pk autoincr"`
		Name string `xorm:"'name'"`
	}

	// Hotel
	type TestFindAndCountHotel struct {
		Id       int64                 `xorm:"bigint(11) pk autoincr"`
		Name     string                `xorm:"'name'"`
		Code     string                `xorm:"'code'"`
		Region   string                `xorm:"'region'"`
		CreateBy *TestFindAndCountUser `xorm:"'create_by'"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(TestFindAndCountUser), new(TestFindAndCountHotel))

	var u = TestFindAndCountUser{
		Name: "myname",
	}
	cnt, err := testEngine.Insert(&u)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var hotel = TestFindAndCountHotel{
		Name:     "myhotel",
		Code:     "111",
		Region:   "222",
		CreateBy: &u,
	}
	cnt, err = testEngine.Insert(&hotel)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	hotels := make([]*TestFindAndCountHotel, 0)
	cnt, err = testEngine.
		Alias("t").
		Limit(10, 0).
		FindAndCount(&hotels)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	hotels = make([]*TestFindAndCountHotel, 0)
	cnt, err = testEngine.
		Table(new(TestFindAndCountHotel)).
		Alias("t").
		Limit(10, 0).
		FindAndCount(&hotels)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	hotels = make([]*TestFindAndCountHotel, 0)
	cnt, err = testEngine.
		Table(new(TestFindAndCountHotel)).
		Alias("t").
		Where("t.region like '6501%'").
		Limit(10, 0).
		FindAndCount(&hotels)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, cnt)
}

type FindMapDevice struct {
	Deviceid string `xorm:"pk"`
	Status   int
}

func (device *FindMapDevice) TableName() string {
	return "devices"
}

func TestFindMapStringId(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(FindMapDevice))

	cnt, err := testEngine.Insert(&FindMapDevice{
		Deviceid: "1",
		Status:   1,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	deviceIDs := []string{"1"}

	deviceMaps := make(map[string]*FindMapDevice, len(deviceIDs))
	err = testEngine.
		Where("status = ?", 1).
		In("deviceid", deviceIDs).
		Find(&deviceMaps)
	assert.NoError(t, err)

	deviceMaps2 := make(map[string]FindMapDevice, len(deviceIDs))
	err = testEngine.
		Where("status = ?", 1).
		In("deviceid", deviceIDs).
		Find(&deviceMaps2)
	assert.NoError(t, err)

	devices := make([]*FindMapDevice, 0, len(deviceIDs))
	err = testEngine.Find(&devices)
	assert.NoError(t, err)

	devices2 := make([]FindMapDevice, 0, len(deviceIDs))
	err = testEngine.Find(&devices2)
	assert.NoError(t, err)

	var device FindMapDevice
	has, err := testEngine.Get(&device)
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = testEngine.Exist(&FindMapDevice{})
	assert.NoError(t, err)
	assert.True(t, has)

	cnt, err = testEngine.Count(new(FindMapDevice))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.ID("1").Update(&FindMapDevice{
		Status: 2,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	sum, err := testEngine.SumInt(new(FindMapDevice), "status")
	assert.NoError(t, err)
	assert.EqualValues(t, 2, sum)

	cnt, err = testEngine.ID("1").Delete(new(FindMapDevice))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func TestFindExtends(t *testing.T) {
	type FindExtendsB struct {
		ID int64
	}

	type FindExtendsA struct {
		FindExtendsB `xorm:"extends"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(FindExtendsA))

	cnt, err := testEngine.Insert(&FindExtendsA{
		FindExtendsB: FindExtendsB{},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Insert(&FindExtendsA{
		FindExtendsB: FindExtendsB{},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var results []FindExtendsA
	err = testEngine.Find(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(results))

	results = make([]FindExtendsA, 0, 2)
	err = testEngine.Find(&results, &FindExtendsB{
		ID: 1,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
}

func TestFindExtends3(t *testing.T) {
	type FindExtendsCC struct {
		ID   int64
		Name string
	}

	type FindExtendsBB struct {
		FindExtendsCC `xorm:"extends"`
	}

	type FindExtendsAA struct {
		FindExtendsBB `xorm:"extends"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(FindExtendsAA))

	cnt, err := testEngine.Insert(&FindExtendsAA{
		FindExtendsBB: FindExtendsBB{
			FindExtendsCC: FindExtendsCC{
				Name: "cc1",
			},
		},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Insert(&FindExtendsAA{
		FindExtendsBB: FindExtendsBB{
			FindExtendsCC: FindExtendsCC{
				Name: "cc2",
			},
		},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var results []FindExtendsAA
	err = testEngine.Find(&results)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(results))
}

func TestFindCacheLimit(t *testing.T) {
	type InviteCode struct {
		ID      int64     `xorm:"pk autoincr 'id'"`
		Code    string    `xorm:"unique"`
		Created time.Time `xorm:"created"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(InviteCode))

	cnt, err := testEngine.Insert(&InviteCode{
		Code: "123456",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	cnt, err = testEngine.Insert(&InviteCode{
		Code: "234567",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	for i := 0; i < 8; i++ {
		var beans []InviteCode
		err = testEngine.Limit(1, 0).Find(&beans)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(beans))
	}

	for i := 0; i < 8; i++ {
		var beans2 []*InviteCode
		err = testEngine.Limit(1, 0).Find(&beans2)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, len(beans2))
	}
}

func TestFindJoin(t *testing.T) {
	type SceneItem struct {
		Type     int
		DeviceId int64
	}

	type DeviceUserPrivrels struct {
		UserId   int64
		DeviceId int64
	}

	type Order struct {
		Id int64
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(SceneItem), new(DeviceUserPrivrels), new(Order))

	var scenes []SceneItem
	err := testEngine.Join("LEFT OUTER", "device_user_privrels", "device_user_privrels.device_id=scene_item.device_id").
		Where("scene_item.type=?", 3).Or("device_user_privrels.user_id=?", 339).Find(&scenes)
	assert.NoError(t, err)

	scenes = make([]SceneItem, 0)
	err = testEngine.Join("LEFT OUTER", new(DeviceUserPrivrels), "device_user_privrels.device_id=scene_item.device_id").
		Where("scene_item.type=?", 3).Or("device_user_privrels.user_id=?", 339).Find(&scenes)
	assert.NoError(t, err)

	scenes = make([]SceneItem, 0)
	err = testEngine.Join("INNER", "order", "`scene_item`.device_id=`order`.id").Find(&scenes)
	assert.NoError(t, err)
}

func TestJoinFindLimit(t *testing.T) {
	type JoinFindLimit1 struct {
		Id   int64
		Name string
	}

	type JoinFindLimit2 struct {
		Id   int64
		Eid  int64 `xorm:"index"`
		Name string
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(JoinFindLimit1), new(JoinFindLimit2))

	var finds []JoinFindLimit1
	err := testEngine.Join("INNER", new(JoinFindLimit2), "join_find_limit2.eid=join_find_limit1.id").
		Limit(10, 10).Find(&finds)
	assert.NoError(t, err)
}

func TestMoreExtends(t *testing.T) {
	type MoreExtendsUsers struct {
		ID        int64     `xorm:"id autoincr pk" json:"id"`
		Name      string    `xorm:"name not null" json:"name"`
		CreatedAt time.Time `xorm:"created not null" json:"created_at"`
		UpdatedAt time.Time `xorm:"updated not null" json:"updated_at"`
		DeletedAt time.Time `xorm:"deleted" json:"deleted_at"`
	}

	type MoreExtendsBooks struct {
		ID        int64     `xorm:"id autoincr pk" json:"id"`
		Name      string    `xorm:"name not null" json:"name"`
		UserID    int64     `xorm:"user_id not null" json:"user_id"`
		CreatedAt time.Time `xorm:"created not null" json:"created_at"`
		UpdatedAt time.Time `xorm:"updated not null" json:"updated_at"`
		DeletedAt time.Time `xorm:"deleted" json:"deleted_at"`
	}

	type MoreExtendsBooksExtend struct {
		MoreExtendsBooks `xorm:"extends"`
		Users            MoreExtendsUsers `xorm:"extends" json:"users"`
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(MoreExtendsUsers), new(MoreExtendsBooks))

	var books []MoreExtendsBooksExtend
	err := testEngine.Table("more_extends_books").Select("more_extends_books.*, more_extends_users.*").
		Join("INNER", "more_extends_users", "more_extends_books.user_id = more_extends_users.id").
		Where("more_extends_books.name LIKE ?", "abc").
		Limit(10, 10).
		Find(&books)
	assert.NoError(t, err)

	books = make([]MoreExtendsBooksExtend, 0, len(books))
	err = testEngine.Table("more_extends_books").
		Alias("m").
		Select("m.*, more_extends_users.*").
		Join("INNER", "more_extends_users", "m.user_id = more_extends_users.id").
		Where("m.name LIKE ?", "abc").
		Limit(10, 10).
		Find(&books)
	assert.NoError(t, err)
}

func TestDistinctAndCols(t *testing.T) {
	type DistinctAndCols struct {
		Id   int64
		Name string
	}

	assert.NoError(t, PrepareEngine())
	assertSync(t, new(DistinctAndCols))

	cnt, err := testEngine.Insert(&DistinctAndCols{
		Name: "test",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var names []string
	err = testEngine.Table("distinct_and_cols").Cols("name").Distinct("name").Find(&names)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(names))
	assert.EqualValues(t, "test", names[0])
}
