// Copyright 2018 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dialects

import (
	"testing"

	"github.com/xormplus/xorm/names"

	"github.com/stretchr/testify/assert"
)

type MCC struct {
	ID          int64  `xorm:"pk 'id'"`
	Code        string `xorm:"'code'"`
	Description string `xorm:"'description'"`
}

func (mcc *MCC) TableName() string {
	return "mcc"
}

func TestFullTableName(t *testing.T) {
	dialect := QueryDialect("mysql")

	assert.EqualValues(t, "mcc", FullTableName(dialect, names.SnakeMapper{}, &MCC{}))
	assert.EqualValues(t, "mcc", FullTableName(dialect, names.SnakeMapper{}, "mcc"))
}
