package xorm

import (
	sql2 "database/sql"
	"fmt"
	"reflect"
	"time"
)

type sqlExpr struct {
	sqlExpr string
}

func noSQLQuoteNeeded(a interface{}) bool {
	switch a.(type) {
	case int, int8, int16, int32, int64:
		return true
	case uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return true
	case bool:
		return true
	case string:
		return false
	case time.Time, *time.Time:
		return false
	case sqlExpr, *sqlExpr:
		return true
	}

	t := reflect.TypeOf(a)
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Bool:
		return true
	case reflect.String:
		return false
	}

	return false
}

// ConvertToBoundSQL will convert SQL and args to a bound SQL
func ConvertToBoundSQL(sql string, args []interface{}) (string, error) {
	buf := StringBuilder{}
	var i, j, start int
	for ; i < len(sql); i++ {
		if sql[i] == '?' {
			_, err := buf.WriteString(sql[start:i])
			if err != nil {
				return "", err
			}
			start = i + 1

			if len(args) == j {
				return "", ErrNeedMoreArguments
			}

			arg := args[j]

			if exprArg, ok := arg.(sqlExpr); ok {
				_, err = fmt.Fprint(&buf, exprArg.sqlExpr)
				if err != nil {
					return "", err
				}

			} else {
				if namedArg, ok := arg.(sql2.NamedArg); ok {
					arg = namedArg.Value
				}

				if noSQLQuoteNeeded(arg) {
					_, err = fmt.Fprint(&buf, arg)
				} else {
					_, err = fmt.Fprintf(&buf, "'%v'", arg)
				}
				if err != nil {
					return "", err
				}
			}

			j = j + 1
		}
	}
	_, err := buf.WriteString(sql[start:])
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
