// Copyright 2020 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package caches

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelDBStore(t *testing.T) {
	store, err := NewLevelDBStore("./level.db")
	assert.NoError(t, err)

	var kvs = map[string]interface{}{
		"a": "b",
	}
	for k, v := range kvs {
		assert.NoError(t, store.Put(k, v))
	}

	for k, v := range kvs {
		val, err := store.Get(k)
		assert.NoError(t, err)
		assert.EqualValues(t, v, val)
	}

	for k := range kvs {
		err := store.Del(k)
		assert.NoError(t, err)
	}

	for k := range kvs {
		_, err := store.Get(k)
		assert.EqualValues(t, ErrNotExist, err)
	}
}
