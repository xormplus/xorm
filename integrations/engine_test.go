// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrations

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/xormplus/xorm"
	"github.com/xormplus/xorm/schemas"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	_ "github.com/ziutek/mymysql/godrv"
)

func TestPing(t *testing.T) {
	if err := testEngine.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestPingContext(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	ctx, canceled := context.WithTimeout(context.Background(), time.Nanosecond)
	defer canceled()

	time.Sleep(time.Nanosecond)

	err := testEngine.(*xorm.Engine).PingContext(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestAutoTransaction(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type TestTx struct {
		Id      int64     `xorm:"autoincr pk"`
		Msg     string    `xorm:"varchar(255)"`
		Created time.Time `xorm:"created"`
	}

	assert.NoError(t, testEngine.Sync2(new(TestTx)))

	engine := testEngine.(*xorm.Engine)

	// will success
	engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		_, err := session.Insert(TestTx{Msg: "hi"})
		assert.NoError(t, err)

		return nil, nil
	})

	has, err := engine.Exist(&TestTx{Msg: "hi"})
	assert.NoError(t, err)
	assert.EqualValues(t, true, has)

	// will rollback
	_, err = engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		_, err := session.Insert(TestTx{Msg: "hello"})
		assert.NoError(t, err)

		return nil, fmt.Errorf("rollback")
	})
	assert.Error(t, err)

	has, err = engine.Exist(&TestTx{Msg: "hello"})
	assert.NoError(t, err)
	assert.EqualValues(t, false, has)
}

func assertSync(t *testing.T, beans ...interface{}) {
	for _, bean := range beans {
		t.Run(testEngine.TableName(bean, true), func(t *testing.T) {
			assert.NoError(t, testEngine.DropTables(bean))
			assert.NoError(t, testEngine.Sync2(bean))
		})
	}
}

func TestDump(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type TestDumpStruct struct {
		Id      int64
		Name    string
		IsMan   bool
		Created time.Time `xorm:"created"`
	}

	assertSync(t, new(TestDumpStruct))

	cnt, err := testEngine.Insert([]TestDumpStruct{
		{Name: "1", IsMan: true},
		{Name: "2\n"},
		{Name: "3;"},
		{Name: "4\n;\n''"},
		{Name: "5'\n"},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 5, cnt)

	fp := fmt.Sprintf("%v.sql", testEngine.Dialect().URI().DBType)
	os.Remove(fp)
	assert.NoError(t, testEngine.DumpAllToFile(fp))

	assert.NoError(t, PrepareEngine())

	sess := testEngine.NewSession()
	defer sess.Close()
	assert.NoError(t, sess.Begin())
	_, err = sess.ImportFile(fp)
	assert.NoError(t, err)
	assert.NoError(t, sess.Commit())

	for _, tp := range []schemas.DBType{schemas.SQLITE, schemas.MYSQL, schemas.POSTGRES, schemas.MSSQL} {
		name := fmt.Sprintf("dump_%v.sql", tp)
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, testEngine.DumpAllToFile(name, tp))
		})
	}
}

func TestDumpTables(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type TestDumpTableStruct struct {
		Id      int64
		Name    string
		IsMan   bool
		Created time.Time `xorm:"created"`
	}

	assertSync(t, new(TestDumpTableStruct))

	testEngine.Insert([]TestDumpTableStruct{
		{Name: "1", IsMan: true},
		{Name: "2\n"},
		{Name: "3;"},
		{Name: "4\n;\n''"},
		{Name: "5'\n"},
	})

	fp := fmt.Sprintf("%v-table.sql", testEngine.Dialect().URI().DBType)
	os.Remove(fp)
	tb, err := testEngine.TableInfo(new(TestDumpTableStruct))
	assert.NoError(t, err)
	assert.NoError(t, testEngine.(*xorm.Engine).DumpTablesToFile([]*schemas.Table{tb}, fp))

	assert.NoError(t, PrepareEngine())

	sess := testEngine.NewSession()
	defer sess.Close()
	assert.NoError(t, sess.Begin())
	_, err = sess.ImportFile(fp)
	assert.NoError(t, err)
	assert.NoError(t, sess.Commit())

	for _, tp := range []schemas.DBType{schemas.SQLITE, schemas.MYSQL, schemas.POSTGRES, schemas.MSSQL} {
		name := fmt.Sprintf("dump_%v-table.sql", tp)
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, testEngine.(*xorm.Engine).DumpTablesToFile([]*schemas.Table{tb}, name, tp))
		})
	}
}

func TestSetSchema(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	if testEngine.Dialect().URI().DBType == schemas.POSTGRES {
		oldSchema := testEngine.Dialect().URI().Schema
		testEngine.SetSchema("my_schema")
		assert.EqualValues(t, "my_schema", testEngine.Dialect().URI().Schema)
		testEngine.SetSchema(oldSchema)
		assert.EqualValues(t, oldSchema, testEngine.Dialect().URI().Schema)
	}
}
