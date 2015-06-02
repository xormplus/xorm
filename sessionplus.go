// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"encoding/json"
	"errors"
//	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/xormplus/core"
	"github.com/Chronokeeper/anyxml"
)

type ResultBean struct {
	Has bool
	Result interface{}
	Error    error
}

func (resultBean ResultBean) Json() (bool,string, error) {
	if resultBean.Error != nil {
		return resultBean.Has,"", resultBean.Error
	}
	if !resultBean.Has{
		return resultBean.Has,"", nil
	}
	result,err:= JSONString(resultBean.Result, true)
	return resultBean.Has,result,err
}

func (session *Session) GetFirst(bean interface{}) ResultBean {
	has, err := session.Get(bean)
	r := ResultBean{Has: has,Result:bean, Error: err}
	return r
}

func (resultBean ResultBean) Xml() (bool,string, error) {
	
	if resultBean.Error != nil {
		return false,"", resultBean.Error
	}
	if !resultBean.Has{
		return resultBean.Has,"", nil
	}
	has,result,err:=resultBean.Json()
	if err != nil {
		return false,"", err
	}
	if !has{
		return has,"", nil
	}
	var anydata = []byte(result)
	var i interface{}
	err = json.Unmarshal(anydata, &i)
	if err != nil {
		return false,"", err
	}
	resultByte, err := anyxml.Xml(i)
	if err != nil {
		return false,"", err
	}

	return resultBean.Has,string(resultByte),err
}

func (resultBean ResultBean) XmlIndent(prefix string, indent string, recordTag string) (bool,string, error) {
	if resultBean.Error != nil {
		return false,"", resultBean.Error
	}
	if !resultBean.Has{
		return resultBean.Has,"", nil
	}
	has,result,err:=resultBean.Json()
	if err != nil {
		return false,"", err
	}
	if !has{
		return has,"", nil
	}
	var anydata = []byte(result)
	var i interface{}
	err = json.Unmarshal(anydata, &i)
	if err != nil {
		return false,"", err
	}
	resultByte, err := anyxml.XmlIndent(i,prefix,indent,recordTag)
	if err != nil {
		return false,"", err
	}

	return resultBean.Has,string(resultByte),err
}

type ResultMap struct {
	Result []map[string]interface{}
	Error    error
}

func (resultMap ResultMap) Json() (string, error) {

	if resultMap.Error != nil {
		return "", resultMap.Error
	}
	return JSONString(resultMap.Result, true)
}

func (resultMap ResultMap) Xml() (string, error) {
	if resultMap.Error != nil {
		return "", resultMap.Error
	}
	results, err := anyxml.Xml(resultMap.Result)
	if err != nil {
		return "", err
	}
	return string(results), nil
}

func (resultMap ResultMap) XmlIndent(prefix string, indent string, recordTag string) (string, error) {
	if resultMap.Error != nil {
		return "", resultMap.Error
	}

	results, err := anyxml.XmlIndent(resultMap.Result, prefix, indent, recordTag)
	if err != nil {
		return "", err
	}
	return string(results), nil
}

// Exec a raw sql and return records as []map[string]interface{}
func (session *Session) Query() ResultMap {
	sql := session.Statement.RawSQL
	params := session.Statement.RawParams
	result, err := session.queryAll(sql, params...)
	r := ResultMap{Result: result, Error: err}
	return r
}

// Exec a raw sql and return records as []map[string]interface{}
func (session *Session) QueryWithDateFormat(dateFormat string) ResultMap {
	sql := session.Statement.RawSQL
	params := session.Statement.RawParams
	result, err := session.queryAllWithDateFormat(dateFormat, sql, params...)
	r := ResultMap{Result: result, Error: err}
	return r
}

// Exec a raw sql and return records as []map[string]interface{}
func (session *Session) QueryByParamMap() ResultMap {
	sql := session.Statement.RawSQL
	params := session.Statement.RawParams
	result, err := session.queryAllByMap(sql, params[0])
	r := ResultMap{Result: result, Error: err}
	return r
}

func (session *Session) QueryByParamMapWithDateFormat(dateFormat string) ResultMap {
	sql := session.Statement.RawSQL
	params := session.Statement.RawParams
	results, err := session.queryAllByMapWithDateFormat(dateFormat, sql, params[0])
	r := ResultMap{Result: results, Error: err}
	return r
}

// =============================
// for Object
// =============================
func (session *Session) queryAll(sqlStr string, paramStr ...interface{}) (resultsSlice []map[string]interface{}, err error) {
	session.queryPreprocess(&sqlStr, paramStr...)

	if session.IsAutoCommit {
		return query3(session.DB(), sqlStr, paramStr...)
	}
	return txQuery3(session.Tx, sqlStr, paramStr...)
}

func (session *Session) queryAllByMap(sqlStr string, paramMap interface{}) (resultsSlice []map[string]interface{}, err error) {
	sqlStr1, param, _ := core.MapToSlice(sqlStr, paramMap)

	session.queryPreprocess(&sqlStr1, param...)

	if session.IsAutoCommit {
		return query3(session.DB(), sqlStr1, param...)
	}
	return txQuery3(session.Tx, sqlStr1, param...)
}

func (session *Session) queryAllByMapWithDateFormat(dateFormat string, sqlStr string, paramMap interface{}) (resultsSlice []map[string]interface{}, err error) {
	sqlStr1, param, _ := core.MapToSlice(sqlStr, paramMap)
	session.queryPreprocess(&sqlStr1, param...)

	if session.IsAutoCommit {
		return query3WithDateFormat(session.DB(), dateFormat, sqlStr1, param...)
	}
	return txQuery3WithDateFormat(session.Tx, dateFormat, sqlStr1, param...)
}

func (session *Session) queryAllWithDateFormat(dateFormat string, sqlStr string, paramStr ...interface{}) (resultsSlice []map[string]interface{}, err error) {
	session.queryPreprocess(&sqlStr, paramStr...)

	if session.IsAutoCommit {
		return query3WithDateFormat(session.DB(), dateFormat, sqlStr, paramStr...)
	}
	return txQuery3WithDateFormat(session.Tx, dateFormat, sqlStr, paramStr...)
}

func (session *Session) queryAllToJsonString(sql string, paramStr ...interface{}) (string, error) {
	results, err := session.queryAll(sql, paramStr...)
	if err != nil {
		return "", err
	}
	return JSONString(results, true)
}

func (session *Session) queryAllToXmlString(sql string, paramStr ...interface{}) (string, error) {
	resultMap, err := session.queryAll(sql, paramStr...)

	if err != nil {
		return "", err
	}
	results, err := anyxml.Xml(resultMap)
	if err != nil {
		return "", err
	}
	return string(results), nil
}

func (session *Session) queryAllToXmlIndentString(sql string, prefix string, indent string, paramStr ...interface{}) (string, error) {
	resultSlice, err := session.queryAll(sql, paramStr...)
	if err != nil {
		return "", err
	}
	results, err := anyxml.XmlIndent(resultSlice, prefix, indent, "result")
	if err != nil {
		return "", err
	}
	return string(results), nil
}

func (session *Session) queryAllToXmlStringWithDateFormat(dateFormat string, sql string, paramStr ...interface{}) (string, error) {
	resultSlice, err := session.queryAll(sql, paramStr...)
	if err != nil {
		return "", err
	}
	results, err := anyxml.XmlWithDateFormat(dateFormat, resultSlice)
	if err != nil {
		return "", err
	}
	return string(results), nil
}

func (session *Session) queryAllToXmlIndentStringWithDateFormat(dateFormat string, sql string, prefix string, indent string, paramStr ...interface{}) (string, error) {
	resultSlice, err := session.queryAll(sql, paramStr...)
	if err != nil {
		return "", err
	}
	results, err := anyxml.XmlIndentWithDateFormat(dateFormat, resultSlice, prefix, indent, "results")

	if err != nil {
		return "", err
	}
	return string(results), nil
}

func (session *Session) queryAllByMapToJsonString(sql string, paramMap interface{}) (string, error) {
	results, err := session.queryAllByMap(sql, paramMap)
	if err != nil {
		return "", err
	}
	return JSONString(results, true)
}

func (session *Session) queryAllByMapToJsonStringWithDateFormat(dateFormat string, sql string, paramMap interface{}) (string, error) {
	results, err := session.queryAllByMapWithDateFormat(dateFormat, sql, paramMap)
	if err != nil {
		return "", err
	}
	return JSONString(results, true)
}

func (session *Session) queryAllToJsonStringWithDateFormat(dateFormat string, sql string, paramStr ...interface{}) (string, error) {
	results, err := session.queryAllWithDateFormat(dateFormat, sql, paramStr...)
	if err != nil {
		return "", err
	}
	return JSONString(results, true)
}

func (session *Session) row2BeanWithDateFormat(dateFormat string, rows *core.Rows, fields []string, fieldsCount int, bean interface{}) error {
	dataStruct := rValue(bean)
	if dataStruct.Kind() != reflect.Struct {
		return errors.New("Expected a pointer to a struct")
	}

	table := session.Engine.autoMapType(dataStruct)
	return session._row2BeanWithDateFormat(dateFormat, rows, fields, fieldsCount, bean, &dataStruct, table)
}

func (session *Session) _row2BeanWithDateFormat(dateFormat string, rows *core.Rows, fields []string, fieldsCount int, bean interface{}, dataStruct *reflect.Value, table *core.Table) error {
	scanResults := make([]interface{}, fieldsCount)
	for i := 0; i < len(fields); i++ {
		var cell interface{}
		scanResults[i] = &cell
	}
	if err := rows.Scan(scanResults...); err != nil {
		return err
	}

	if b, hasBeforeSet := bean.(BeforeSetProcessor); hasBeforeSet {
		for ii, key := range fields {
			b.BeforeSet(key, Cell(scanResults[ii].(*interface{})))
		}
	}

	var tempMap = make(map[string]int)
	for ii, key := range fields {
		var idx int
		var ok bool
		if idx, ok = tempMap[strings.ToLower(key)]; !ok {
			idx = 0
		} else {
			idx = idx + 1
		}
		tempMap[strings.ToLower(key)] = idx

		if fieldValue := session.getField(dataStruct, key, table, idx); fieldValue != nil {
			rawValue := reflect.Indirect(reflect.ValueOf(scanResults[ii]))

			//if row is null then ignore
			if rawValue.Interface() == nil {
				continue
			}

			if fieldValue.CanAddr() {
				if structConvert, ok := fieldValue.Addr().Interface().(core.Conversion); ok {
					if data, err := value2Bytes(&rawValue); err == nil {
						structConvert.FromDB(data)
					} else {
						session.Engine.LogError(err)
					}
					continue
				}
			}

			if _, ok := fieldValue.Interface().(core.Conversion); ok {
				if data, err := value2Bytes(&rawValue); err == nil {
					if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
						fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
					}
					fieldValue.Interface().(core.Conversion).FromDB(data)
				} else {
					session.Engine.LogError(err)
				}
				continue
			}

			rawValueType := reflect.TypeOf(rawValue.Interface())
			vv := reflect.ValueOf(rawValue.Interface())

			fieldType := fieldValue.Type()
			hasAssigned := false

			switch fieldType.Kind() {

			case reflect.Complex64, reflect.Complex128:
				if rawValueType.Kind() == reflect.String {
					hasAssigned = true
					x := reflect.New(fieldType)
					err := json.Unmarshal([]byte(vv.String()), x.Interface())
					if err != nil {
						session.Engine.LogError(err)
						return err
					}
					fieldValue.Set(x.Elem())
				}
			case reflect.Slice, reflect.Array:
				switch rawValueType.Kind() {
				case reflect.Slice, reflect.Array:
					switch rawValueType.Elem().Kind() {
					case reflect.Uint8:
						if fieldType.Elem().Kind() == reflect.Uint8 {
							hasAssigned = true
							fieldValue.Set(vv)
						}
					}
				}
			case reflect.String:
				if rawValueType.Kind() == reflect.String {
					hasAssigned = true
					fieldValue.SetString(vv.String())
				}
			case reflect.Bool:
				if rawValueType.Kind() == reflect.Bool {
					hasAssigned = true
					fieldValue.SetBool(vv.Bool())
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch rawValueType.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					hasAssigned = true
					fieldValue.SetInt(vv.Int())
				}
			case reflect.Float32, reflect.Float64:
				switch rawValueType.Kind() {
				case reflect.Float32, reflect.Float64:
					hasAssigned = true
					fieldValue.SetFloat(vv.Float())
				}
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				switch rawValueType.Kind() {
				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
					hasAssigned = true
					fieldValue.SetUint(vv.Uint())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					hasAssigned = true
					fieldValue.SetUint(uint64(vv.Int()))
				}
			case reflect.Struct:
				if fieldType.ConvertibleTo(core.TimeType) {
					if rawValueType == core.TimeType {
						hasAssigned = true

						t := vv.Convert(core.TimeType).Interface().(time.Time)
						z, _ := t.Zone()
						if len(z) == 0 || t.Year() == 0 { // !nashtsai! HACK tmp work around for lib/pq doesn't properly time with location
							session.Engine.LogDebug("empty zone key[%v] : %v | zone: %v | location: %+v\n", key, t, z, *t.Location())
							t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),
								t.Minute(), t.Second(), t.Nanosecond(), time.Local)
						}
						// !nashtsai! convert to engine location
						t = t.In(session.Engine.TZLocation)
						// dateFormat to string
						loc, _ := time.LoadLocation("Local") //重要：获取时区  rawValue.Interface().(time.Time).Format(dateFormat)
						t, _ = time.ParseInLocation(dateFormat, t.Format(dateFormat), loc)
						//						fieldValue.Set(reflect.ValueOf(t).Convert(core.StringType))
						fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))

						// t = fieldValue.Interface().(time.Time)
						// z, _ = t.Zone()
						// session.Engine.LogDebug("fieldValue key[%v]: %v | zone: %v | location: %+v\n", key, t, z, *t.Location())
					} else if rawValueType == core.IntType || rawValueType == core.Int64Type ||
						rawValueType == core.Int32Type {
						hasAssigned = true
						t := time.Unix(vv.Int(), 0).In(session.Engine.TZLocation)
						vv = reflect.ValueOf(t)
						fieldValue.Set(vv)
					}
				} else if session.Statement.UseCascade {
					table := session.Engine.autoMapType(*fieldValue)
					if table != nil {
						if len(table.PrimaryKeys) > 1 {
							panic("unsupported composited primary key cascade")
						}
						var pk = make(core.PK, len(table.PrimaryKeys))
						switch rawValueType.Kind() {
						case reflect.Int64:
							pk[0] = vv.Int()
						case reflect.Int:
							pk[0] = int(vv.Int())
						case reflect.Int32:
							pk[0] = int32(vv.Int())
						case reflect.Int16:
							pk[0] = int16(vv.Int())
						case reflect.Int8:
							pk[0] = int8(vv.Int())
						case reflect.Uint64:
							pk[0] = vv.Uint()
						case reflect.Uint:
							pk[0] = uint(vv.Uint())
						case reflect.Uint32:
							pk[0] = uint32(vv.Uint())
						case reflect.Uint16:
							pk[0] = uint16(vv.Uint())
						case reflect.Uint8:
							pk[0] = uint8(vv.Uint())
						case reflect.String:
							pk[0] = vv.String()
						default:
							panic("unsupported primary key type cascade")
						}

						if !isPKZero(pk) {
							// !nashtsai! TODO for hasOne relationship, it's preferred to use join query for eager fetch
							// however, also need to consider adding a 'lazy' attribute to xorm tag which allow hasOne
							// property to be fetched lazily
							structInter := reflect.New(fieldValue.Type())
							newsession := session.Engine.NewSession()
							defer newsession.Close()
							has, err := newsession.Id(pk).NoCascade().Get(structInter.Interface())
							if err != nil {
								return err
							}
							if has {
								v := structInter.Elem().Interface()
								fieldValue.Set(reflect.ValueOf(v))
							} else {
								return errors.New("cascade obj is not exist!")
							}
						}
					} else {
						session.Engine.LogError("unsupported struct type in Scan: ", fieldValue.Type().String())
					}
				}
			case reflect.Ptr:
				// !nashtsai! TODO merge duplicated codes above
				//typeStr := fieldType.String()
				switch fieldType {
				// following types case matching ptr's native type, therefore assign ptr directly
				case core.PtrStringType:
					if rawValueType.Kind() == reflect.String {
						x := vv.String()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrBoolType:
					if rawValueType.Kind() == reflect.Bool {
						x := vv.Bool()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrTimeType:
					if rawValueType == core.PtrTimeType {
						hasAssigned = true
						var x time.Time = rawValue.Interface().(time.Time)
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrFloat64Type:
					if rawValueType.Kind() == reflect.Float64 {
						x := vv.Float()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUint64Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x uint64 = uint64(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt64Type:
					if rawValueType.Kind() == reflect.Int64 {
						x := vv.Int()
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrFloat32Type:
					if rawValueType.Kind() == reflect.Float64 {
						var x float32 = float32(vv.Float())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrIntType:
					if rawValueType.Kind() == reflect.Int64 {
						var x int = int(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt32Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x int32 = int32(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt8Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x int8 = int8(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt16Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x int16 = int16(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUintType:
					if rawValueType.Kind() == reflect.Int64 {
						var x uint = uint(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUint32Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x uint32 = uint32(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Uint8Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x uint8 = uint8(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Uint16Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x uint16 = uint16(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Complex64Type:
					var x complex64
					err := json.Unmarshal([]byte(vv.String()), &x)
					if err != nil {
						session.Engine.LogError(err)
					} else {
						fieldValue.Set(reflect.ValueOf(&x))
					}
					hasAssigned = true
				case core.Complex128Type:
					var x complex128
					err := json.Unmarshal([]byte(vv.String()), &x)
					if err != nil {
						session.Engine.LogError(err)
					} else {
						fieldValue.Set(reflect.ValueOf(&x))
					}
					hasAssigned = true
				} // switch fieldType
				// default:
				// 	session.Engine.LogError("unsupported type in Scan: ", reflect.TypeOf(v).String())
			} // switch fieldType.Kind()

			// !nashtsai! for value can't be assigned directly fallback to convert to []byte then back to value
			if !hasAssigned {
				data, err := value2Bytes(&rawValue)
				if err == nil {
					session.bytes2Value(table.GetColumn(key), fieldValue, data)
				} else {
					session.Engine.LogError(err.Error())
				}
			}
		}
	}
	return nil

}

func (session *Session) queryPreprocessByMap(sqlStr *string, paramMap interface{}) {
	re := regexp.MustCompile(`[?](\w+)`)
	query := *sqlStr
	names := make(map[string]int)
	var i int
	query = re.ReplaceAllStringFunc(query, func(src string) string {
		names[src[1:]] = i
		i += 1
		return "?"
	})

	for _, filter := range session.Engine.dialect.Filters() {
		query = filter.Do(query, session.Engine.dialect, session.Statement.RefTable)
	}

	*sqlStr = query
	session.Engine.logSQL(*sqlStr, paramMap)
}
