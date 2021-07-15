// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrations

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xormplus/xorm/internal/utils"
	"github.com/xormplus/xorm/names"
)

func TestTransaction(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	counter := func(t *testing.T) {
		_, err := testEngine.Count(&Userinfo{})
		assert.NoError(t, err)
	}

	counter(t)
	//defer counter()

	session := testEngine.NewSession()
	defer session.Close()

	err := session.Begin()
	assert.NoError(t, err)

	user1 := Userinfo{Username: "xiaoxiao", Departname: "dev", Alias: "lunny", Created: time.Now()}
	_, err = session.Insert(&user1)
	assert.NoError(t, err)

	user2 := Userinfo{Username: "yyy"}
	_, err = session.Where("id = ?", 0).Update(&user2)
	assert.NoError(t, err)

	_, err = session.Delete(&user2)
	assert.NoError(t, err)

	err = session.Commit()
	assert.NoError(t, err)
}

func TestCombineTransaction(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))

	counter := func() {
		total, err := testEngine.Count(&Userinfo{})
		assert.NoError(t, err)
		fmt.Printf("----now total %v records\n", total)
	}

	counter()
	//defer counter()
	session := testEngine.NewSession()
	defer session.Close()

	err := session.Begin()
	assert.NoError(t, err)

	user1 := Userinfo{Username: "xiaoxiao2", Departname: "dev", Alias: "lunny", Created: time.Now()}
	_, err = session.Insert(&user1)
	assert.NoError(t, err)

	user2 := Userinfo{Username: "zzz"}
	_, err = session.Where("id = ?", 0).Update(&user2)
	assert.NoError(t, err)

	_, err = session.Exec("delete from "+testEngine.TableName("userinfo", true)+" where username = ?", user2.Username)
	assert.NoError(t, err)

	err = session.Commit()
	assert.NoError(t, err)
}

func TestCombineTransactionSameMapper(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	oldMapper := testEngine.GetColumnMapper()
	testEngine.UnMapType(utils.ReflectValue(new(Userinfo)).Type())
	testEngine.SetMapper(names.SameMapper{})
	defer func() {
		testEngine.UnMapType(utils.ReflectValue(new(Userinfo)).Type())
		testEngine.SetMapper(oldMapper)
	}()

	assertSync(t, new(Userinfo))

	counter := func() {
		total, err := testEngine.Count(&Userinfo{})
		assert.NoError(t, err)
		fmt.Printf("----now total %v records\n", total)
	}

	counter()
	defer counter()

	session := testEngine.NewSession()
	defer session.Close()

	err := session.Begin()
	assert.NoError(t, err)

	user1 := Userinfo{Username: "xiaoxiao2", Departname: "dev", Alias: "lunny", Created: time.Now()}
	_, err = session.Insert(&user1)
	assert.NoError(t, err)

	user2 := Userinfo{Username: "zzz"}
	_, err = session.Where("id = ?", 0).Update(&user2)
	assert.NoError(t, err)

	_, err = session.Exec("delete from  "+testEngine.TableName("`Userinfo`", true)+" where `Username` = ?", user2.Username)
	assert.NoError(t, err)

	err = session.Commit()
	assert.NoError(t, err)
}

func TestMultipleTransaction(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type MultipleTransaction struct {
		Id   int64
		Name string
	}

	assertSync(t, new(MultipleTransaction))

	session := testEngine.NewSession()
	defer session.Close()

	err := session.Begin()
	assert.NoError(t, err)

	m1 := MultipleTransaction{Name: "xiaoxiao2"}
	_, err = session.Insert(&m1)
	assert.NoError(t, err)

	user2 := MultipleTransaction{Name: "zzz"}
	_, err = session.Where("id = ?", 0).Update(&user2)
	assert.NoError(t, err)

	err = session.Commit()
	assert.NoError(t, err)

	var ms []MultipleTransaction
	err = session.Find(&ms)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(ms))

	err = session.Begin()
	assert.NoError(t, err)

	_, err = session.Where("id=?", m1.Id).Delete(new(MultipleTransaction))
	assert.NoError(t, err)

	err = session.Commit()
	assert.NoError(t, err)

	ms = make([]MultipleTransaction, 0)
	err = session.Find(&ms)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, len(ms))

	err = session.Begin()
	assert.NoError(t, err)

	_, err = session.Insert(&MultipleTransaction{
		Name: "ssss",
	})
	assert.NoError(t, err)

	err = session.Rollback()
	assert.NoError(t, err)

	ms = make([]MultipleTransaction, 0)
	err = session.Find(&ms)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, len(ms))
}
