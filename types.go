package xorm

import (
	"reflect"

	"github.com/xormplus/core"
)

var (
	ptrPkType = reflect.TypeOf(&core.PK{})
	pkType    = reflect.TypeOf(core.PK{})
)
