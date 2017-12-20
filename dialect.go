package xorm

import "github.com/xormplus/core"

type MutableFilter interface {
	AddFilter(filters ...core.Filter)
}
