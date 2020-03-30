// Copyright 2020 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrations

import (
	"testing"

	"github.com/xormplus/xorm"
	"github.com/xormplus/xorm/log"
	"github.com/xormplus/xorm/schemas"

	"github.com/stretchr/testify/assert"
)

func TestEngineGroup(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	master := testEngine.(*xorm.Engine)
	if master.Dialect().URI().DBType == schemas.SQLITE {
		t.Skip()
		return
	}

	eg, err := xorm.NewEngineGroup(master, []*xorm.Engine{master})
	assert.NoError(t, err)

	eg.SetMaxIdleConns(10)
	eg.SetMaxOpenConns(100)
	eg.SetTableMapper(master.GetTableMapper())
	eg.SetColumnMapper(master.GetColumnMapper())
	eg.SetLogLevel(log.LOG_INFO)
	eg.ShowSQL(true)
}
