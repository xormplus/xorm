// Copyright 2020 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MyInt int
type ZeroStruct struct{}

func TestZero(t *testing.T) {
	var zeroValues = []interface{}{
		int8(0),
		int16(0),
		int(0),
		int32(0),
		int64(0),
		uint8(0),
		uint16(0),
		uint(0),
		uint32(0),
		uint64(0),
		MyInt(0),
		reflect.ValueOf(0),
		nil,
		time.Time{},
		&time.Time{},
		nilTime,
		ZeroStruct{},
		&ZeroStruct{},
	}

	for _, v := range zeroValues {
		t.Run(fmt.Sprintf("%#v", v), func(t *testing.T) {
			assert.True(t, IsZero(v))
		})
	}
}

func TestIsValueZero(t *testing.T) {
	var zeroReflectValues = []reflect.Value{
		reflect.ValueOf(int8(0)),
		reflect.ValueOf(int16(0)),
		reflect.ValueOf(int(0)),
		reflect.ValueOf(int32(0)),
		reflect.ValueOf(int64(0)),
		reflect.ValueOf(uint8(0)),
		reflect.ValueOf(uint16(0)),
		reflect.ValueOf(uint(0)),
		reflect.ValueOf(uint32(0)),
		reflect.ValueOf(uint64(0)),
		reflect.ValueOf(MyInt(0)),
		reflect.ValueOf(time.Time{}),
		reflect.ValueOf(&time.Time{}),
		reflect.ValueOf(nilTime),
		reflect.ValueOf(ZeroStruct{}),
		reflect.ValueOf(&ZeroStruct{}),
	}

	for _, v := range zeroReflectValues {
		t.Run(fmt.Sprintf("%#v", v), func(t *testing.T) {
			assert.True(t, IsValueZero(v))
		})
	}
}
