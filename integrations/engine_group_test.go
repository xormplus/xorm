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

	main := testEngine.(*xorm.Engine)
	if main.Dialect().URI().DBType == schemas.SQLITE {
		t.Skip()
		return
	}

	eg, err := xorm.NewEngineGroup(main, []*xorm.Engine{main})
	assert.NoError(t, err)

	eg.SetMaxIdleConns(10)
	eg.SetMaxOpenConns(100)
	eg.SetTableMapper(main.GetTableMapper())
	eg.SetColumnMapper(main.GetColumnMapper())
	eg.SetLogLevel(log.LOG_INFO)
	eg.ShowSQL(true)
}
