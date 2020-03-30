// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tags

import (
	"testing"

	"github.com/xormplus/xorm/internal/utils"
)

func TestSplitTag(t *testing.T) {
	var cases = []struct {
		tag  string
		tags []string
	}{
		{"not null default '2000-01-01 00:00:00' TIMESTAMP", []string{"not", "null", "default", "'2000-01-01 00:00:00'", "TIMESTAMP"}},
		{"TEXT", []string{"TEXT"}},
		{"default('2000-01-01 00:00:00')", []string{"default('2000-01-01 00:00:00')"}},
		{"json  binary", []string{"json", "binary"}},
	}

	for _, kase := range cases {
		tags := splitTag(kase.tag)
		if !utils.SliceEq(tags, kase.tags) {
			t.Fatalf("[%d]%v is not equal [%d]%v", len(tags), tags, len(kase.tags), kase.tags)
		}
	}
}
