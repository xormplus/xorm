// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrations

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xormplus/xorm/internal/utils"
	"github.com/xormplus/xorm/names"
	"github.com/xormplus/xorm/schemas"
)

type tempUser struct {
	Id       int64
	Username string
}

type tempUser2 struct {
	TempUser   tempUser `xorm:"extends"`
	Departname string
}

type tempUser3 struct {
	Temp       *tempUser `xorm:"extends"`
	Departname string
}

type tempUser4 struct {
	TempUser2 tempUser2 `xorm:"extends"`
}

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

type UserAndDetail struct {
	Userinfo   `xorm:"extends"`
	Userdetail `xorm:"extends"`
}

func TestExtends(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&tempUser2{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&tempUser2{})
	assert.NoError(t, err)

	tu := &tempUser2{tempUser{0, "extends"}, "dev depart"}
	_, err = testEngine.Insert(tu)
	assert.NoError(t, err)

	tu2 := &tempUser2{}
	_, err = testEngine.Get(tu2)
	assert.NoError(t, err)

	tu3 := &tempUser2{tempUser{0, "extends update"}, ""}
	_, err = testEngine.ID(tu2.TempUser.Id).Update(tu3)
	assert.NoError(t, err)

	err = testEngine.DropTables(&tempUser4{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&tempUser4{})
	assert.NoError(t, err)

	tu8 := &tempUser4{tempUser2{tempUser{0, "extends"}, "dev depart"}}
	_, err = testEngine.Insert(tu8)
	assert.NoError(t, err)

	tu9 := &tempUser4{}
	_, err = testEngine.Get(tu9)
	assert.NoError(t, err)
	assert.EqualValues(t, tu8.TempUser2.TempUser.Username, tu9.TempUser2.TempUser.Username)
	assert.EqualValues(t, tu8.TempUser2.Departname, tu9.TempUser2.Departname)

	tu10 := &tempUser4{tempUser2{tempUser{0, "extends update"}, ""}}
	_, err = testEngine.ID(tu9.TempUser2.TempUser.Id).Update(tu10)
	assert.NoError(t, err)

	err = testEngine.DropTables(&tempUser3{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&tempUser3{})
	assert.NoError(t, err)

	tu4 := &tempUser3{&tempUser{0, "extends"}, "dev depart"}
	_, err = testEngine.Insert(tu4)
	assert.NoError(t, err)

	tu5 := &tempUser3{}
	_, err = testEngine.Get(tu5)
	assert.NoError(t, err)

	assert.NotNil(t, tu5.Temp)
	assert.EqualValues(t, 1, tu5.Temp.Id)
	assert.EqualValues(t, "extends", tu5.Temp.Username)
	assert.EqualValues(t, "dev depart", tu5.Departname)

	tu6 := &tempUser3{&tempUser{0, "extends update"}, ""}
	_, err = testEngine.ID(tu5.Temp.Id).Update(tu6)
	assert.NoError(t, err)

	users := make([]tempUser3, 0)
	err = testEngine.Find(&users)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(users), "error get data not 1")

	assertSync(t, new(Userinfo), new(Userdetail))

	detail := Userdetail{
		Intro: "I'm in China",
	}
	_, err = testEngine.Insert(&detail)
	assert.NoError(t, err)

	_, err = testEngine.Insert(&Userinfo{
		Username: "lunny",
		Detail:   detail,
	})
	assert.NoError(t, err)

	var info UserAndDetail
	qt := testEngine.Quote
	ui := testEngine.TableName(new(Userinfo), true)
	ud := testEngine.TableName(&detail, true)
	uiid := testEngine.GetColumnMapper().Obj2Table("Id")
	udid := "detail_id"
	sql := fmt.Sprintf("select * from %s, %s where %s.%s = %s.%s",
		qt(ui), qt(ud), qt(ui), qt(udid), qt(ud), qt(uiid))
	b, err := testEngine.SQL(sql).NoCascade().Get(&info)
	assert.NoError(t, err)
	assert.True(t, b, "should has lest one record")
	assert.True(t, info.Userinfo.Uid > 0, "all of the id should has value")
	assert.True(t, info.Userdetail.Id > 0, "all of the id should has value")

	var info2 UserAndDetail
	b, err = testEngine.Table(&Userinfo{}).
		Join("LEFT", qt(ud), qt(ui)+"."+qt("detail_id")+" = "+qt(ud)+"."+qt(uiid)).
		NoCascade().Get(&info2)
	assert.NoError(t, err)
	assert.True(t, b)
	assert.True(t, info2.Userinfo.Uid > 0, "all of the id should has value")
	assert.True(t, info2.Userdetail.Id > 0, "all of the id should has value")

	var infos2 = make([]UserAndDetail, 0)
	err = testEngine.Table(&Userinfo{}).
		Join("LEFT", qt(ud), qt(ui)+"."+qt("detail_id")+" = "+qt(ud)+"."+qt(uiid)).
		NoCascade().
		Find(&infos2)
	assert.NoError(t, err)
}

type MessageBase struct {
	Id     int64 `xorm:"int(11) pk autoincr"`
	TypeId int64 `xorm:"int(11) notnull"`
}

type Message struct {
	MessageBase `xorm:"extends"`
	Title       string    `xorm:"varchar(100) notnull"`
	Content     string    `xorm:"text notnull"`
	Uid         int64     `xorm:"int(11) notnull"`
	ToUid       int64     `xorm:"int(11) notnull"`
	CreateTime  time.Time `xorm:"datetime notnull created"`
}

type MessageUser struct {
	Id   int64
	Name string
}

type MessageType struct {
	Id   int64
	Name string
}

type MessageExtend3 struct {
	Message  `xorm:"extends"`
	Sender   MessageUser `xorm:"extends"`
	Receiver MessageUser `xorm:"extends"`
	Type     MessageType `xorm:"extends"`
}

type MessageExtend4 struct {
	Message     `xorm:"extends"`
	MessageUser `xorm:"extends"`
	MessageType `xorm:"extends"`
}

func TestExtends2(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&Message{}, &MessageUser{}, &MessageType{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&Message{}, &MessageUser{}, &MessageType{})
	assert.NoError(t, err)

	var sender = MessageUser{Name: "sender"}
	var receiver = MessageUser{Name: "receiver"}
	var msgtype = MessageType{Name: "type"}
	_, err = testEngine.Insert(&sender, &receiver, &msgtype)
	assert.NoError(t, err)

	msg := Message{
		MessageBase: MessageBase{
			Id: msgtype.Id,
		},
		Title:   "test",
		Content: "test",
		Uid:     sender.Id,
		ToUid:   receiver.Id,
	}

	session := testEngine.NewSession()
	defer session.Close()

	// MSSQL deny insert identity column excep declare as below
	if testEngine.Dialect().URI().DBType == schemas.MSSQL {
		err = session.Begin()
		assert.NoError(t, err)
		_, err = session.Exec("SET IDENTITY_INSERT message ON")
		assert.NoError(t, err)
	}
	cnt, err := session.Insert(&msg)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	if testEngine.Dialect().URI().DBType == schemas.MSSQL {
		err = session.Commit()
		assert.NoError(t, err)
	}

	var mapper = testEngine.GetTableMapper().Obj2Table
	var quote = testEngine.Quote
	userTableName := quote(testEngine.TableName(mapper("MessageUser"), true))
	typeTableName := quote(testEngine.TableName(mapper("MessageType"), true))
	msgTableName := quote(testEngine.TableName(mapper("Message"), true))

	list := make([]Message, 0)
	err = session.Table(msgTableName).Join("LEFT", []string{userTableName, "sender"}, "`sender`.`"+mapper("Id")+"`="+msgTableName+".`"+mapper("Uid")+"`").
		Join("LEFT", []string{userTableName, "receiver"}, "`receiver`.`"+mapper("Id")+"`="+msgTableName+".`"+mapper("ToUid")+"`").
		Join("LEFT", []string{typeTableName, "type"}, "`type`.`"+mapper("Id")+"`="+msgTableName+".`"+mapper("Id")+"`").
		Find(&list)
	assert.NoError(t, err)

	assert.EqualValues(t, 1, len(list), fmt.Sprintln("should have 1 message, got", len(list)))
	assert.EqualValues(t, msg.Id, list[0].Id, fmt.Sprintln("should message equal", list[0], msg))
}

func TestExtends3(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&Message{}, &MessageUser{}, &MessageType{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&Message{}, &MessageUser{}, &MessageType{})
	assert.NoError(t, err)

	var sender = MessageUser{Name: "sender"}
	var receiver = MessageUser{Name: "receiver"}
	var msgtype = MessageType{Name: "type"}
	_, err = testEngine.Insert(&sender, &receiver, &msgtype)
	assert.NoError(t, err)

	msg := Message{
		MessageBase: MessageBase{
			Id: msgtype.Id,
		},
		Title:   "test",
		Content: "test",
		Uid:     sender.Id,
		ToUid:   receiver.Id,
	}

	session := testEngine.NewSession()
	defer session.Close()

	// MSSQL deny insert identity column excep declare as below
	if testEngine.Dialect().URI().DBType == schemas.MSSQL {
		err = session.Begin()
		assert.NoError(t, err)
		_, err = session.Exec("SET IDENTITY_INSERT message ON")
		assert.NoError(t, err)
	}
	_, err = session.Insert(&msg)
	assert.NoError(t, err)

	if testEngine.Dialect().URI().DBType == schemas.MSSQL {
		err = session.Commit()
		assert.NoError(t, err)
	}

	var mapper = testEngine.GetTableMapper().Obj2Table
	var quote = testEngine.Quote
	userTableName := quote(testEngine.TableName(mapper("MessageUser"), true))
	typeTableName := quote(testEngine.TableName(mapper("MessageType"), true))
	msgTableName := quote(testEngine.TableName(mapper("Message"), true))

	list := make([]MessageExtend3, 0)
	err = session.Table(msgTableName).Join("LEFT", []string{userTableName, "sender"}, "`sender`.`"+mapper("Id")+"`="+msgTableName+".`"+mapper("Uid")+"`").
		Join("LEFT", []string{userTableName, "receiver"}, "`receiver`.`"+mapper("Id")+"`="+msgTableName+".`"+mapper("ToUid")+"`").
		Join("LEFT", []string{typeTableName, "type"}, "`type`.`"+mapper("Id")+"`="+msgTableName+".`"+mapper("Id")+"`").
		Find(&list)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(list))
	assert.EqualValues(t, list[0].Message.Id, msg.Id)
	assert.EqualValues(t, list[0].Sender.Id, sender.Id)
	assert.EqualValues(t, list[0].Sender.Name, sender.Name)
	assert.EqualValues(t, list[0].Receiver.Id, receiver.Id)
	assert.EqualValues(t, list[0].Receiver.Name, receiver.Name)
	assert.EqualValues(t, list[0].Type.Id, msgtype.Id)
	assert.EqualValues(t, list[0].Type.Name, msgtype.Name)
}

func TestExtends4(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(&Message{}, &MessageUser{}, &MessageType{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&Message{}, &MessageUser{}, &MessageType{})
	assert.NoError(t, err)

	var sender = MessageUser{Name: "sender"}
	var msgtype = MessageType{Name: "type"}
	_, err = testEngine.Insert(&sender, &msgtype)
	assert.NoError(t, err)

	msg := Message{
		MessageBase: MessageBase{
			Id: msgtype.Id,
		},
		Title:   "test",
		Content: "test",
		Uid:     sender.Id,
	}

	session := testEngine.NewSession()
	defer session.Close()

	// MSSQL deny insert identity column excep declare as below
	if testEngine.Dialect().URI().DBType == schemas.MSSQL {
		err = session.Begin()
		assert.NoError(t, err)
		_, err = session.Exec("SET IDENTITY_INSERT message ON")
		assert.NoError(t, err)
	}
	_, err = session.Insert(&msg)
	assert.NoError(t, err)

	if testEngine.Dialect().URI().DBType == schemas.MSSQL {
		err = session.Commit()
		assert.NoError(t, err)
	}

	var mapper = testEngine.GetTableMapper().Obj2Table
	var quote = testEngine.Quote
	userTableName := quote(testEngine.TableName(mapper("MessageUser"), true))
	typeTableName := quote(testEngine.TableName(mapper("MessageType"), true))
	msgTableName := quote(testEngine.TableName(mapper("Message"), true))

	list := make([]MessageExtend4, 0)
	err = session.Table(msgTableName).Join("LEFT", userTableName, userTableName+".`"+mapper("Id")+"`="+msgTableName+".`"+mapper("Uid")+"`").
		Join("LEFT", typeTableName, typeTableName+".`"+mapper("Id")+"`="+msgTableName+".`"+mapper("Id")+"`").
		Find(&list)
	assert.NoError(t, err)
	assert.EqualValues(t, len(list), 1)
	assert.EqualValues(t, list[0].Message.Id, msg.Id)
	assert.EqualValues(t, list[0].MessageUser.Id, sender.Id)
	assert.EqualValues(t, list[0].MessageUser.Name, sender.Name)
	assert.EqualValues(t, list[0].MessageType.Id, msgtype.Id)
	assert.EqualValues(t, list[0].MessageType.Name, msgtype.Name)
}

type Size struct {
	ID     int64   `xorm:"int(4) 'id' pk autoincr"`
	Width  float32 `json:"width" xorm:"float 'Width'"`
	Height float32 `json:"height" xorm:"float 'Height'"`
}

type Book struct {
	ID         int64 `xorm:"int(4) 'id' pk autoincr"`
	SizeOpen   *Size `xorm:"extends('Open')"`
	SizeClosed *Size `xorm:"extends('Closed')"`
	Size       *Size `xorm:"extends('')"`
}

func TestExtends5(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	err := testEngine.DropTables(&Book{}, &Size{})
	assert.NoError(t, err)

	err = testEngine.CreateTables(&Size{}, &Book{})
	assert.NoError(t, err)

	var sc = Size{Width: 0.2, Height: 0.4}
	var so = Size{Width: 0.2, Height: 0.8}
	var s = Size{Width: 0.15, Height: 1.5}
	var bk1 = Book{
		SizeOpen:   &so,
		SizeClosed: &sc,
		Size:       &s,
	}
	var bk2 = Book{
		SizeOpen: &so,
	}
	var bk3 = Book{
		SizeClosed: &sc,
		Size:       &s,
	}
	var bk4 = Book{}
	var bk5 = Book{Size: &s}
	_, err = testEngine.Insert(&sc, &so, &s, &bk1, &bk2, &bk3, &bk4, &bk5)
	if err != nil {
		t.Fatal(err)
	}

	var books = map[int64]Book{
		bk1.ID: bk1,
		bk2.ID: bk2,
		bk3.ID: bk3,
		bk4.ID: bk4,
		bk5.ID: bk5,
	}

	session := testEngine.NewSession()
	defer session.Close()

	var mapper = testEngine.GetTableMapper().Obj2Table
	var quote = testEngine.Quote
	bookTableName := quote(testEngine.TableName(mapper("Book"), true))
	sizeTableName := quote(testEngine.TableName(mapper("Size"), true))

	list := make([]Book, 0)
	err = session.
		Select(fmt.Sprintf(
			"%s.%s, sc.%s AS %s, sc.%s AS %s, s.%s, s.%s",
			quote(bookTableName),
			quote("id"),
			quote("Width"),
			quote("ClosedWidth"),
			quote("Height"),
			quote("ClosedHeight"),
			quote("Width"),
			quote("Height"),
		)).
		Table(bookTableName).
		Join(
			"LEFT",
			sizeTableName+" AS `sc`",
			bookTableName+".`SizeClosed`=sc.`id`",
		).
		Join(
			"LEFT",
			sizeTableName+" AS `s`",
			bookTableName+".`Size`=s.`id`",
		).
		Find(&list)
	assert.NoError(t, err)

	for _, book := range list {
		if ok := assert.Equal(t, books[book.ID].SizeClosed.Width, book.SizeClosed.Width); !ok {
			t.Error("Not bounded size closed")
			panic("Not bounded size closed")
		}

		if ok := assert.Equal(t, books[book.ID].SizeClosed.Height, book.SizeClosed.Height); !ok {
			t.Error("Not bounded size closed")
			panic("Not bounded size closed")
		}

		if books[book.ID].Size != nil || book.Size != nil {
			if ok := assert.Equal(t, books[book.ID].Size.Width, book.Size.Width); !ok {
				t.Error("Not bounded size")
				panic("Not bounded size")
			}

			if ok := assert.Equal(t, books[book.ID].Size.Height, book.Size.Height); !ok {
				t.Error("Not bounded size")
				panic("Not bounded size")
			}
		}
	}
}

func TestCacheTag(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type CacheDomain struct {
		Id   int64 `xorm:"pk cache"`
		Name string
	}

	assert.NoError(t, testEngine.CreateTables(&CacheDomain{}))
	assert.True(t, testEngine.GetCacher(testEngine.TableName(&CacheDomain{})) != nil)
}

func TestNoCacheTag(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type NoCacheDomain struct {
		Id   int64 `xorm:"pk nocache"`
		Name string
	}

	assert.NoError(t, testEngine.CreateTables(&NoCacheDomain{}))
	assert.True(t, testEngine.GetCacher(testEngine.TableName(&NoCacheDomain{})) == nil)
}

type IDGonicMapper struct {
	ID int64
}

func TestGonicMapperID(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	oldMapper := testEngine.GetColumnMapper()
	testEngine.UnMapType(utils.ReflectValue(new(IDGonicMapper)).Type())
	testEngine.SetMapper(names.LintGonicMapper)
	defer func() {
		testEngine.UnMapType(utils.ReflectValue(new(IDGonicMapper)).Type())
		testEngine.SetMapper(oldMapper)
	}()

	err := testEngine.CreateTables(new(IDGonicMapper))
	if err != nil {
		t.Fatal(err)
	}

	tables, err := testEngine.DBMetas()
	if err != nil {
		t.Fatal(err)
	}

	for _, tb := range tables {
		if tb.Name == "id_gonic_mapper" {
			if len(tb.PKColumns()) != 1 || tb.PKColumns()[0].Name != "id" {
				t.Fatal(tb)
			}
			return
		}
	}

	t.Fatal("not table id_gonic_mapper")
}

type IDSameMapper struct {
	ID int64
}

func TestSameMapperID(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	oldMapper := testEngine.GetColumnMapper()
	testEngine.UnMapType(utils.ReflectValue(new(IDSameMapper)).Type())
	testEngine.SetMapper(names.SameMapper{})
	defer func() {
		testEngine.UnMapType(utils.ReflectValue(new(IDSameMapper)).Type())
		testEngine.SetMapper(oldMapper)
	}()

	err := testEngine.CreateTables(new(IDSameMapper))
	if err != nil {
		t.Fatal(err)
	}

	tables, err := testEngine.DBMetas()
	if err != nil {
		t.Fatal(err)
	}

	for _, tb := range tables {
		if tb.Name == "IDSameMapper" {
			if len(tb.PKColumns()) != 1 || tb.PKColumns()[0].Name != "ID" {
				t.Fatalf("tb %s tb.PKColumns() is %d not 1, tb.PKColumns()[0].Name is %s not ID", tb.Name, len(tb.PKColumns()), tb.PKColumns()[0].Name)
			}
			return
		}
	}
	t.Fatal("not table IDSameMapper")
}

type UserCU struct {
	Id      int64
	Name    string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func TestCreatedAndUpdated(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	u := new(UserCU)
	err := testEngine.DropTables(u)
	assert.NoError(t, err)

	err = testEngine.CreateTables(u)
	assert.NoError(t, err)

	u.Name = "sss"
	cnt, err := testEngine.Insert(u)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	u.Name = "xxx"
	cnt, err = testEngine.ID(u.Id).Update(u)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	u.Id = 0
	u.Created = time.Now().Add(-time.Hour * 24 * 365)
	u.Updated = u.Created
	cnt, err = testEngine.NoAutoTime().Insert(u)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

type StrangeName struct {
	Id_t int64 `xorm:"pk autoincr"`
	Name string
}

func TestStrangeName(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(new(StrangeName))
	assert.NoError(t, err)

	err = testEngine.CreateTables(new(StrangeName))
	assert.NoError(t, err)

	_, err = testEngine.Insert(&StrangeName{Name: "sfsfdsfds"})
	assert.NoError(t, err)

	beans := make([]StrangeName, 0)
	err = testEngine.Find(&beans)
	assert.NoError(t, err)
}

func TestCreatedUpdated(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type CreatedUpdated struct {
		Id       int64
		Name     string
		Value    float64   `xorm:"numeric"`
		Created  time.Time `xorm:"created"`
		Created2 time.Time `xorm:"created"`
		Updated  time.Time `xorm:"updated"`
	}

	err := testEngine.Sync2(&CreatedUpdated{})
	assert.NoError(t, err)

	c := &CreatedUpdated{Name: "test"}
	_, err = testEngine.Insert(c)
	assert.NoError(t, err)

	c2 := new(CreatedUpdated)
	has, err := testEngine.ID(c.Id).Get(c2)
	assert.NoError(t, err)

	assert.True(t, has)

	c2.Value--
	_, err = testEngine.ID(c2.Id).Update(c2)
	assert.NoError(t, err)
}

func TestCreatedUpdatedInt64(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type CreatedUpdatedInt64 struct {
		Id       int64
		Name     string
		Value    float64 `xorm:"numeric"`
		Created  int64   `xorm:"created"`
		Created2 int64   `xorm:"created"`
		Updated  int64   `xorm:"updated"`
	}

	assertSync(t, &CreatedUpdatedInt64{})

	c := &CreatedUpdatedInt64{Name: "test"}
	_, err := testEngine.Insert(c)
	assert.NoError(t, err)

	c2 := new(CreatedUpdatedInt64)
	has, err := testEngine.ID(c.Id).Get(c2)
	assert.NoError(t, err)
	assert.True(t, has)

	c2.Value--
	_, err = testEngine.ID(c2.Id).Update(c2)
	assert.NoError(t, err)
}

type Lowercase struct {
	Id    int64
	Name  string
	ended int64 `xorm:"-"`
}

func TestLowerCase(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.Sync2(&Lowercase{})
	assert.NoError(t, err)
	_, err = testEngine.Where("id > 0").Delete(&Lowercase{})
	assert.NoError(t, err)

	_, err = testEngine.Insert(&Lowercase{ended: 1})
	assert.NoError(t, err)

	ls := make([]Lowercase, 0)
	err = testEngine.Find(&ls)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(ls))
}

func TestAutoIncrTag(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type TestAutoIncr1 struct {
		Id int64
	}

	tb, err := testEngine.TableInfo(new(TestAutoIncr1))
	assert.NoError(t, err)

	cols := tb.Columns()
	assert.EqualValues(t, 1, len(cols))
	assert.True(t, cols[0].IsAutoIncrement)
	assert.True(t, cols[0].IsPrimaryKey)
	assert.Equal(t, "id", cols[0].Name)

	type TestAutoIncr2 struct {
		Id int64 `xorm:"id"`
	}

	tb, err = testEngine.TableInfo(new(TestAutoIncr2))
	assert.NoError(t, err)

	cols = tb.Columns()
	assert.EqualValues(t, 1, len(cols))
	assert.False(t, cols[0].IsAutoIncrement)
	assert.False(t, cols[0].IsPrimaryKey)
	assert.Equal(t, "id", cols[0].Name)

	type TestAutoIncr3 struct {
		Id int64 `xorm:"'ID'"`
	}

	tb, err = testEngine.TableInfo(new(TestAutoIncr3))
	assert.NoError(t, err)

	cols = tb.Columns()
	assert.EqualValues(t, 1, len(cols))
	assert.False(t, cols[0].IsAutoIncrement)
	assert.False(t, cols[0].IsPrimaryKey)
	assert.Equal(t, "ID", cols[0].Name)

	type TestAutoIncr4 struct {
		Id int64 `xorm:"pk"`
	}

	tb, err = testEngine.TableInfo(new(TestAutoIncr4))
	assert.NoError(t, err)

	cols = tb.Columns()
	assert.EqualValues(t, 1, len(cols))
	assert.False(t, cols[0].IsAutoIncrement)
	assert.True(t, cols[0].IsPrimaryKey)
	assert.Equal(t, "id", cols[0].Name)
}

func TestTagComment(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	// FIXME: only support mysql
	if testEngine.Dialect().URI().DBType != schemas.MYSQL {
		return
	}

	type TestComment1 struct {
		Id int64 `xorm:"comment(主键)"`
	}

	assert.NoError(t, testEngine.Sync2(new(TestComment1)))

	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, 1, len(tables[0].Columns()))
	assert.EqualValues(t, "主键", tables[0].Columns()[0].Comment)

	assert.NoError(t, testEngine.DropTables(new(TestComment1)))

	type TestComment2 struct {
		Id int64 `xorm:"comment('主键')"`
	}

	assert.NoError(t, testEngine.Sync2(new(TestComment2)))

	tables, err = testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, 1, len(tables[0].Columns()))
	assert.EqualValues(t, "主键", tables[0].Columns()[0].Comment)
}

func TestTagDefault(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type DefaultStruct struct {
		Id   int64
		Name string
		Age  int `xorm:"default(10)"`
	}

	assertSync(t, new(DefaultStruct))

	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)

	var defaultVal string
	var isDefaultExist bool
	tableName := testEngine.GetColumnMapper().Obj2Table("DefaultStruct")
	for _, table := range tables {
		if table.Name == tableName {
			col := table.GetColumn("age")
			assert.NotNil(t, col)
			defaultVal = col.Default
			isDefaultExist = !col.DefaultIsEmpty
			break
		}
	}
	assert.True(t, isDefaultExist)
	assert.EqualValues(t, "10", defaultVal)

	cnt, err := testEngine.Omit("age").Insert(&DefaultStruct{
		Name: "test",
		Age:  20,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var s DefaultStruct
	has, err := testEngine.ID(1).Get(&s)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, 10, s.Age)
	assert.EqualValues(t, "test", s.Name)
}

func TestTagDefault2(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type DefaultStruct2 struct {
		Id   int64
		Name string
	}

	assertSync(t, new(DefaultStruct2))

	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)

	var defaultVal string
	var isDefaultExist bool
	tableName := testEngine.GetColumnMapper().Obj2Table("DefaultStruct2")
	for _, table := range tables {
		if table.Name == tableName {
			col := table.GetColumn("name")
			assert.NotNil(t, col)
			defaultVal = col.Default
			isDefaultExist = !col.DefaultIsEmpty
			break
		}
	}
	assert.False(t, isDefaultExist, fmt.Sprintf("default value is --%v--", defaultVal))
	assert.EqualValues(t, "", defaultVal)
}

func TestTagDefault3(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type DefaultStruct3 struct {
		Id   int64
		Name string `xorm:"default('myname')"`
	}

	assertSync(t, new(DefaultStruct3))

	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)

	var defaultVal string
	var isDefaultExist bool
	tableName := testEngine.GetColumnMapper().Obj2Table("DefaultStruct3")
	for _, table := range tables {
		if table.Name == tableName {
			col := table.GetColumn("name")
			assert.NotNil(t, col)
			defaultVal = col.Default
			isDefaultExist = !col.DefaultIsEmpty
			break
		}
	}
	assert.True(t, isDefaultExist)
	assert.EqualValues(t, "'myname'", defaultVal)
}

func TestTagDefault4(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type DefaultStruct4 struct {
		Id      int64
		Created time.Time `xorm:"default(CURRENT_TIMESTAMP)"`
	}

	assertSync(t, new(DefaultStruct4))

	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)

	var defaultVal string
	var isDefaultExist bool
	tableName := testEngine.GetColumnMapper().Obj2Table("DefaultStruct4")
	for _, table := range tables {
		if table.Name == tableName {
			col := table.GetColumn("created")
			assert.NotNil(t, col)
			defaultVal = col.Default
			isDefaultExist = !col.DefaultIsEmpty
			break
		}
	}
	assert.True(t, isDefaultExist)
	assert.True(t, "CURRENT_TIMESTAMP" == defaultVal ||
		"current_timestamp()" == defaultVal || // for cockroach
		"now()" == defaultVal ||
		"getdate" == defaultVal, defaultVal)
}

func TestTagDefault5(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type DefaultStruct5 struct {
		Id      int64
		Created time.Time `xorm:"default('2006-01-02 15:04:05')"`
	}

	assertSync(t, new(DefaultStruct5))
	table, err := testEngine.TableInfo(new(DefaultStruct5))
	assert.NoError(t, err)

	createdCol := table.GetColumn("created")
	assert.NotNil(t, createdCol)
	assert.EqualValues(t, "'2006-01-02 15:04:05'", createdCol.Default)
	assert.False(t, createdCol.DefaultIsEmpty)

	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)

	var defaultVal string
	var isDefaultExist bool
	tableName := testEngine.GetColumnMapper().Obj2Table("DefaultStruct5")
	for _, table := range tables {
		if table.Name == tableName {
			col := table.GetColumn("created")
			assert.NotNil(t, col)
			defaultVal = col.Default
			isDefaultExist = !col.DefaultIsEmpty
			break
		}
	}
	assert.True(t, isDefaultExist)
	assert.EqualValues(t, "'2006-01-02 15:04:05'", defaultVal)
}

func TestTagDefault6(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type DefaultStruct6 struct {
		Id    int64
		IsMan bool `xorm:"default(true)"`
	}

	assertSync(t, new(DefaultStruct6))

	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)

	var defaultVal string
	var isDefaultExist bool
	tableName := testEngine.GetColumnMapper().Obj2Table("DefaultStruct6")
	for _, table := range tables {
		if table.Name == tableName {
			col := table.GetColumn("is_man")
			assert.NotNil(t, col)
			defaultVal = col.Default
			isDefaultExist = !col.DefaultIsEmpty
			break
		}
	}
	assert.True(t, isDefaultExist)
	if defaultVal == "1" {
		defaultVal = "true"
	} else if defaultVal == "0" {
		defaultVal = "false"
	}
	assert.EqualValues(t, "true", defaultVal)
}

func TestTagsDirection(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type OnlyFromDBStruct struct {
		Id   int64
		Name string
		Uuid string `xorm:"<- default '1'"`
	}

	assertSync(t, new(OnlyFromDBStruct))

	cnt, err := testEngine.Insert(&OnlyFromDBStruct{
		Name: "test",
		Uuid: "2",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var s OnlyFromDBStruct
	has, err := testEngine.ID(1).Get(&s)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, "1", s.Uuid)
	assert.EqualValues(t, "test", s.Name)

	cnt, err = testEngine.ID(1).Update(&OnlyFromDBStruct{
		Uuid: "3",
		Name: "test1",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var s3 OnlyFromDBStruct
	has, err = testEngine.ID(1).Get(&s3)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, "1", s3.Uuid)
	assert.EqualValues(t, "test1", s3.Name)

	type OnlyToDBStruct struct {
		Id   int64
		Name string
		Uuid string `xorm:"->"`
	}

	assertSync(t, new(OnlyToDBStruct))

	cnt, err = testEngine.Insert(&OnlyToDBStruct{
		Name: "test",
		Uuid: "2",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var s2 OnlyToDBStruct
	has, err = testEngine.ID(1).Get(&s2)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, "", s2.Uuid)
	assert.EqualValues(t, "test", s2.Name)
}

func TestTagTime(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type TagUTCStruct struct {
		Id      int64
		Name    string
		Created time.Time `xorm:"created utc"`
	}

	assertSync(t, new(TagUTCStruct))

	assert.EqualValues(t, time.Local.String(), testEngine.GetTZLocation().String())

	s := TagUTCStruct{
		Name: "utc",
	}
	cnt, err := testEngine.Insert(&s)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var u TagUTCStruct
	has, err := testEngine.ID(1).Get(&u)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, s.Created.Format("2006-01-02 15:04:05"), u.Created.Format("2006-01-02 15:04:05"))

	var tm string
	has, err = testEngine.Table("tag_u_t_c_struct").Cols("created").Get(&tm)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, s.Created.UTC().Format("2006-01-02 15:04:05"),
		strings.Replace(strings.Replace(tm, "T", " ", -1), "Z", "", -1))
}

func TestTagAutoIncr(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type TagAutoIncr struct {
		Id   int64
		Name string
	}

	assertSync(t, new(TagAutoIncr))

	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, tableMapper.Obj2Table("TagAutoIncr"), tables[0].Name)
	col := tables[0].GetColumn(colMapper.Obj2Table("Id"))
	assert.NotNil(t, col)
	assert.True(t, col.IsPrimaryKey)
	assert.True(t, col.IsAutoIncrement)

	col2 := tables[0].GetColumn(colMapper.Obj2Table("Name"))
	assert.NotNil(t, col2)
	assert.False(t, col2.IsPrimaryKey)
	assert.False(t, col2.IsAutoIncrement)
}

func TestTagPrimarykey(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type TagPrimaryKey struct {
		Id   int64  `xorm:"pk"`
		Name string `xorm:"VARCHAR(20) pk"`
	}

	assertSync(t, new(TagPrimaryKey))

	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, tableMapper.Obj2Table("TagPrimaryKey"), tables[0].Name)
	col := tables[0].GetColumn(colMapper.Obj2Table("Id"))
	assert.NotNil(t, col)
	assert.True(t, col.IsPrimaryKey)
	assert.False(t, col.IsAutoIncrement)

	col2 := tables[0].GetColumn(colMapper.Obj2Table("Name"))
	assert.NotNil(t, col2)
	assert.True(t, col2.IsPrimaryKey)
	assert.False(t, col2.IsAutoIncrement)
}

type VersionS struct {
	Id      int64
	Name    string
	Ver     int       `xorm:"version"`
	Created time.Time `xorm:"created"`
}

func TestVersion1(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(new(VersionS))
	assert.NoError(t, err)

	err = testEngine.CreateTables(new(VersionS))
	assert.NoError(t, err)

	ver := &VersionS{Name: "sfsfdsfds"}
	_, err = testEngine.Insert(ver)
	assert.NoError(t, err)
	assert.EqualValues(t, ver.Ver, 1)

	newVer := new(VersionS)
	has, err := testEngine.ID(ver.Id).Get(newVer)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, newVer.Ver, 1)

	newVer.Name = "-------"
	_, err = testEngine.ID(ver.Id).Update(newVer)
	assert.NoError(t, err)
	assert.EqualValues(t, newVer.Ver, 2)

	newVer = new(VersionS)
	has, err = testEngine.ID(ver.Id).Get(newVer)
	assert.NoError(t, err)
	assert.EqualValues(t, newVer.Ver, 2)
}

func TestVersion2(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(new(VersionS))
	assert.NoError(t, err)

	err = testEngine.CreateTables(new(VersionS))
	assert.NoError(t, err)

	var vers = []VersionS{
		{Name: "sfsfdsfds"},
		{Name: "xxxxx"},
	}
	_, err = testEngine.Insert(vers)
	assert.NoError(t, err)
	for _, v := range vers {
		assert.EqualValues(t, v.Ver, 1)
	}
}

type VersionUintS struct {
	Id      int64
	Name    string
	Ver     uint      `xorm:"version"`
	Created time.Time `xorm:"created"`
}

func TestVersion3(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(new(VersionUintS))
	assert.NoError(t, err)

	err = testEngine.CreateTables(new(VersionUintS))
	assert.NoError(t, err)

	ver := &VersionUintS{Name: "sfsfdsfds"}
	_, err = testEngine.Insert(ver)
	assert.NoError(t, err)
	assert.EqualValues(t, ver.Ver, 1)

	newVer := new(VersionUintS)
	has, err := testEngine.ID(ver.Id).Get(newVer)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, newVer.Ver, 1)

	newVer.Name = "-------"
	_, err = testEngine.ID(ver.Id).Update(newVer)
	assert.NoError(t, err)
	assert.EqualValues(t, newVer.Ver, 2)

	newVer = new(VersionUintS)
	has, err = testEngine.ID(ver.Id).Get(newVer)
	assert.NoError(t, err)
	assert.EqualValues(t, newVer.Ver, 2)
}

func TestVersion4(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	err := testEngine.DropTables(new(VersionUintS))
	assert.NoError(t, err)

	err = testEngine.CreateTables(new(VersionUintS))
	assert.NoError(t, err)

	var vers = []VersionUintS{
		{Name: "sfsfdsfds"},
		{Name: "xxxxx"},
	}
	_, err = testEngine.Insert(vers)
	assert.NoError(t, err)
	for _, v := range vers {
		assert.EqualValues(t, v.Ver, 1)
	}
}

func TestIndexes(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type TestIndexesStruct struct {
		Id    int64
		Name  string `xorm:"index unique(s)"`
		Email string `xorm:"index unique(s)"`
	}

	assertSync(t, new(TestIndexesStruct))

	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, 3, len(tables[0].Columns()))
	slice1 := []string{
		testEngine.GetColumnMapper().Obj2Table("Id"),
		testEngine.GetColumnMapper().Obj2Table("Name"),
		testEngine.GetColumnMapper().Obj2Table("Email"),
	}
	slice2 := []string{
		tables[0].Columns()[0].Name,
		tables[0].Columns()[1].Name,
		tables[0].Columns()[2].Name,
	}
	sort.Strings(slice1)
	sort.Strings(slice2)
	assert.EqualValues(t, slice1, slice2)
	assert.EqualValues(t, 3, len(tables[0].Indexes))
}
