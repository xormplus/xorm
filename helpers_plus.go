package xorm

import (
	"fmt"
	"reflect"
	"time"

	"github.com/xormplus/core"
)

func reflect2objectWithDateFormat(rawValue *reflect.Value, dateFormat string) (value interface{}, err error) {
	aa := reflect.TypeOf((*rawValue).Interface())
	vv := reflect.ValueOf((*rawValue).Interface())
	switch aa.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = vv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = vv.Uint()
	case reflect.Float32, reflect.Float64:
		value = vv.Float()
	case reflect.String:
		value = vv.String()
	case reflect.Array, reflect.Slice:
		switch aa.Elem().Kind() {
		case reflect.Uint8:
			data := rawValue.Interface().([]byte)
			value = string(data)
		default:
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	// time type
	case reflect.Struct:
		if aa.ConvertibleTo(core.TimeType) {
			value = vv.Convert(core.TimeType).Interface().(time.Time).Format(dateFormat)
		} else {
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	case reflect.Bool:
		value = vv.Bool()
	case reflect.Complex128, reflect.Complex64:
		value = vv.Complex()
	/* TODO: unsupported types below
	   case reflect.Map:
	   case reflect.Ptr:
	   case reflect.Uintptr:
	   case reflect.UnsafePointer:
	   case reflect.Chan, reflect.Func, reflect.Interface:
	*/
	default:
		err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
	}
	return
}

func value2ObjectWithDateFormat(rawValue *reflect.Value, dateFormat string) (data interface{}, err error) {
	data, err = reflect2objectWithDateFormat(rawValue, dateFormat)
	if err != nil {
		return
	}
	return
}

func rows2mapObjectsWithDateFormat(rows *core.Rows, dateFormat string) (resultsSlice []map[string]interface{}, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := rows2mapObjectWithDateFormat(rows, dateFormat, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func rows2mapObjectWithDateFormat(rows *core.Rows, dateFormat string, fields []string) (resultsMap map[string]interface{}, err error) {
	result := make(map[string]interface{})
	scanResultContainers := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return nil, err
	}

	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		//if row is null then ignore
		if rawValue.Interface() == nil {
			continue
		}

		if data, err := value2ObjectWithDateFormat(&rawValue, dateFormat); err == nil {
			result[key] = data
		} else {
			return nil, err // !nashtsai! REVIEW, should return err or just error log?
		}

	}
	return result, nil
}

func txQueryByMap(tx *core.Tx, sqlStr string, params interface{}) (resultsSlice []map[string]interface{}, err error) {
	rows, err := tx.QueryMap(sqlStr, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rows2mapObjects(rows)
}

func txQuery3WithDateFormat(tx *core.Tx, dateFormat string, sqlStr string, params ...interface{}) (resultsSlice []map[string]interface{}, err error) {
	rows, err := tx.Query(sqlStr, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rows2mapObjectsWithDateFormat(rows, dateFormat)
}

func queryByMap(db *core.DB, sqlStr string, params interface{}) (resultsSlice []map[string]interface{}, err error) {
	s, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer s.Close()

	rows, err := s.QueryMap(params)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rows2mapObjects(rows)
}

func query3WithDateFormat(db *core.DB, dateFormat string, sqlStr string, params ...interface{}) (resultsSlice []map[string]interface{}, err error) {
	s, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer s.Close()
	rows, err := s.Query(params...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rows2mapObjectsWithDateFormat(rows, dateFormat)
}

func reflect2object(rawValue *reflect.Value) (value interface{}, err error) {
	aa := reflect.TypeOf((*rawValue).Interface())
	vv := reflect.ValueOf((*rawValue).Interface())
	switch aa.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = vv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = vv.Uint()
	case reflect.Float32, reflect.Float64:
		value = vv.Float()
	case reflect.String:
		value = vv.String()
	case reflect.Array, reflect.Slice:
		switch aa.Elem().Kind() {
		case reflect.Uint8:
			data := rawValue.Interface().([]byte)
			value = string(data)
		default:
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	// time type
	case reflect.Struct:
		if aa.ConvertibleTo(core.TimeType) {
			value = vv.Convert(core.TimeType).Interface().(time.Time)
		} else {
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	case reflect.Bool:
		value = vv.Bool()
	case reflect.Complex128, reflect.Complex64:
		value = vv.Complex()
	/* TODO: unsupported types below
	   case reflect.Map:
	   case reflect.Ptr:
	   case reflect.Uintptr:
	   case reflect.UnsafePointer:
	   case reflect.Chan, reflect.Func, reflect.Interface:
	*/
	default:
		err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
	}
	return
}

func value2Object(rawValue *reflect.Value) (data interface{}, err error) {
	data, err = reflect2object(rawValue)
	if err != nil {
		return
	}
	return
}

func rows2mapObjects(rows *core.Rows) (resultsSlice []map[string]interface{}, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := rows2mapObject(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func rows2mapObject(rows *core.Rows, fields []string) (resultsMap map[string]interface{}, err error) {
	result := make(map[string]interface{})
	scanResultContainers := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return nil, err
	}

	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		//if row is null then ignore
		if rawValue.Interface() == nil {
			continue
		}

		if data, err := value2Object(&rawValue); err == nil {
			result[key] = data
		} else {
			return nil, err // !nashtsai! REVIEW, should return err or just error log?
		}

	}
	return result, nil
}

func txQuery3(tx *core.Tx, sqlStr string, params ...interface{}) (resultsSlice []map[string]interface{}, err error) {
	rows, err := tx.Query(sqlStr, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rows2mapObjects(rows)
}

func query3(db *core.DB, sqlStr string, params ...interface{}) (resultsSlice []map[string]interface{}, err error) {
	s, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer s.Close()
	rows, err := s.Query(params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rows2mapObjects(rows)
}
