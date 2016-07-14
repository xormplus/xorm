// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"

	"strconv"
	"strings"
	"time"

	"github.com/Chronokeeper/anyxml"
	"github.com/xormplus/core"
)

type ResultBean struct {
	Has    bool
	Result interface{}
	Error  error
}

func (resultBean *ResultBean) Json() (bool, string, error) {
	if resultBean.Error != nil {
		return resultBean.Has, "", resultBean.Error
	}
	if !resultBean.Has {
		return resultBean.Has, "", nil
	}
	result, err := JSONString(resultBean.Result, true)
	return resultBean.Has, result, err
}

func (resultBean *ResultBean) GetResult() (bool, interface{}, error) {
	return resultBean.Has, resultBean.Result, resultBean.Error
}

func (session *Session) GetFirst(bean interface{}) *ResultBean {
	has, err := session.Get(bean)
	r := &ResultBean{Has: has, Result: bean, Error: err}
	return r
}

func (resultBean *ResultBean) Xml() (bool, string, error) {

	if resultBean.Error != nil {
		return false, "", resultBean.Error
	}
	if !resultBean.Has {
		return resultBean.Has, "", nil
	}
	has, result, err := resultBean.Json()
	if err != nil {
		return false, "", err
	}
	if !has {
		return has, "", nil
	}
	var anydata = []byte(result)
	var i interface{}
	err = json.Unmarshal(anydata, &i)
	if err != nil {
		return false, "", err
	}
	resultByte, err := anyxml.Xml(i)
	if err != nil {
		return false, "", err
	}

	return resultBean.Has, string(resultByte), err
}

func (resultBean *ResultBean) XmlIndent(prefix string, indent string, recordTag string) (bool, string, error) {
	if resultBean.Error != nil {
		return false, "", resultBean.Error
	}
	if !resultBean.Has {
		return resultBean.Has, "", nil
	}
	has, result, err := resultBean.Json()
	if err != nil {
		return false, "", err
	}
	if !has {
		return has, "", nil
	}
	var anydata = []byte(result)
	var i interface{}
	err = json.Unmarshal(anydata, &i)
	if err != nil {
		return false, "", err
	}
	resultByte, err := anyxml.XmlIndent(i, prefix, indent, recordTag)
	if err != nil {
		return false, "", err
	}

	return resultBean.Has, string(resultByte), err
}

type ResultMap struct {
	Results []map[string]interface{}
	Error   error
}

func (resultMap *ResultMap) List() ([]map[string]interface{}, error) {
	return resultMap.Results, resultMap.Error
}

func (resultMap *ResultMap) Count() (int, error) {
	if resultMap.Error != nil {
		return 0, resultMap.Error
	}
	if resultMap.Results == nil {
		return 0, nil
	}
	return len(resultMap.Results), nil
}

func (resultMap *ResultMap) ListPage(firstResult int, maxResults int) ([]map[string]interface{}, error) {
	if resultMap.Error != nil {
		return nil, resultMap.Error
	}
	if resultMap.Results == nil {
		return nil, nil
	}
	if firstResult >= maxResults {
		return nil, ErrParamsFormat
	}
	if firstResult < 0 {
		return nil, ErrParamsFormat
	}
	if maxResults < 0 {
		return nil, ErrParamsFormat
	}
	if maxResults > len(resultMap.Results) {
		return nil, ErrParamsFormat
	}
	return resultMap.Results[(firstResult - 1):maxResults], resultMap.Error
}

func (resultMap *ResultMap) Json() (string, error) {

	if resultMap.Error != nil {
		return "", resultMap.Error
	}
	return JSONString(resultMap.Results, true)
}

func (resultMap *ResultMap) Xml() (string, error) {
	if resultMap.Error != nil {
		return "", resultMap.Error
	}
	results, err := anyxml.Xml(resultMap.Results)
	if err != nil {
		return "", err
	}
	return string(results), nil
}

func (resultMap *ResultMap) XmlIndent(prefix string, indent string, recordTag string) (string, error) {
	if resultMap.Error != nil {
		return "", resultMap.Error
	}

	results, err := anyxml.XmlIndent(resultMap.Results, prefix, indent, recordTag)
	if err != nil {
		return "", err
	}
	return string(results), nil
}

func (resultMap *ResultMap) SaveAsCSV(filename string, headers []string, perm os.FileMode) error {
	if resultMap.Error != nil {
		return resultMap.Error
	}

	dataset, err := NewDatasetWithData(headers, resultMap.Results)
	if err != nil {
		return err
	}
	csv, err := dataset.CSV()
	if err != nil {
		return err
	}

	return csv.WriteFile(filename, perm)

}

func (resultMap *ResultMap) SaveAsTSV(filename string, headers []string, perm os.FileMode) error {
	if resultMap.Error != nil {
		return resultMap.Error
	}

	dataset, err := NewDatasetWithData(headers, resultMap.Results)
	if err != nil {
		return err
	}
	tsv, err := dataset.TSV()
	if err != nil {
		return err
	}

	return tsv.WriteFile(filename, perm)

}

func (resultMap *ResultMap) SaveAsHTML(filename string, headers []string, perm os.FileMode) error {
	if resultMap.Error != nil {
		return resultMap.Error
	}

	dataset, err := NewDatasetWithData(headers, resultMap.Results)
	if err != nil {
		return err
	}
	html := dataset.HTML()

	return html.WriteFile(filename, perm)

}

func (resultMap *ResultMap) SaveAsXML(filename string, headers []string, perm os.FileMode) error {
	if resultMap.Error != nil {
		return resultMap.Error
	}

	dataset, err := NewDatasetWithData(headers, resultMap.Results)
	if err != nil {
		return err
	}
	xml, err := dataset.XML()
	if err != nil {
		return err
	}

	return xml.WriteFile(filename, perm)

}

func (resultMap *ResultMap) SaveAsXMLWithTagNamePrefixIndent(tagName string, prifix string, indent string, filename string, headers []string, perm os.FileMode) error {
	if resultMap.Error != nil {
		return resultMap.Error
	}

	dataset, err := NewDatasetWithData(headers, resultMap.Results)
	if err != nil {
		return err
	}
	xml, err := dataset.XMLWithTagNamePrefixIndent(tagName, prifix, indent)
	if err != nil {
		return err
	}

	return xml.WriteFile(filename, perm)

}

func (resultMap *ResultMap) SaveAsYAML(filename string, headers []string, perm os.FileMode) error {
	if resultMap.Error != nil {
		return resultMap.Error
	}

	dataset, err := NewDatasetWithData(headers, resultMap.Results)
	if err != nil {
		return err
	}
	yaml, err := dataset.YAML()
	if err != nil {
		return err
	}

	return yaml.WriteFile(filename, perm)

}

func (resultMap *ResultMap) SaveAsJSON(filename string, headers []string, perm os.FileMode) error {
	if resultMap.Error != nil {
		return resultMap.Error
	}

	dataset, err := NewDatasetWithData(headers, resultMap.Results)
	if err != nil {
		return err
	}
	json, err := dataset.JSON()
	if err != nil {
		return err
	}

	return json.WriteFile(filename, perm)

}

func (resultMap *ResultMap) SaveAsXLSX(filename string, headers []string, perm os.FileMode) error {
	if resultMap.Error != nil {
		return resultMap.Error
	}

	dataset, err := NewDatasetWithData(headers, resultMap.Results)
	if err != nil {
		return err
	}
	xlsx, err := dataset.XLSX()
	if err != nil {
		return err
	}

	return xlsx.WriteFile(filename, perm)

}

type ResultStructs struct {
	Result interface{}
	Error  error
}

func (resultStructs *ResultStructs) Json() (string, error) {

	if resultStructs.Error != nil {
		return "", resultStructs.Error
	}
	return JSONString(resultStructs.Result, true)
}

func (resultStructs *ResultStructs) Xml() (string, error) {
	if resultStructs.Error != nil {
		return "", resultStructs.Error
	}

	result, err := resultStructs.Json()
	if err != nil {
		return "", err
	}

	var anydata = []byte(result)
	var i interface{}
	err = json.Unmarshal(anydata, &i)
	if err != nil {
		return "", err
	}
	resultByte, err := anyxml.Xml(i)
	if err != nil {
		return "", err
	}

	return string(resultByte), nil
}

func (resultStructs *ResultStructs) XmlIndent(prefix string, indent string, recordTag string) (string, error) {
	if resultStructs.Error != nil {
		return "", resultStructs.Error
	}

	result, err := resultStructs.Json()
	if err != nil {
		return "", err
	}

	var anydata = []byte(result)
	var i interface{}
	err = json.Unmarshal(anydata, &i)
	if err != nil {
		return "", err
	}
	resultByte, err := anyxml.XmlIndent(i, prefix, indent, recordTag)
	if err != nil {
		return "", err
	}

	return string(resultByte), nil
}

func (session *Session) SqlMapClient(sqlTagName string, args ...interface{}) *Session {
	return session.Sql(session.Engine.sqlMap.Sql[sqlTagName], args...)
}

func (session *Session) SqlTemplateClient(sqlTagName string, args ...interface{}) *Session {
	map1 := args[0].(map[string]interface{})
	if session.Engine.sqlTemplate.Template[sqlTagName] == nil {
		return session.Sql("", &map1)
	}
	sql, err := session.Engine.sqlTemplate.Template[sqlTagName].Execute(map1)
	if err != nil {
		session.Engine.logger.Error(err)
	}

	return session.Sql(sql, &map1)
}

func (session *Session) Search(rowsSlicePtr interface{}, condiBean ...interface{}) *ResultStructs {
	err := session.Find(rowsSlicePtr, condiBean...)
	r := &ResultStructs{Result: rowsSlicePtr, Error: err}
	return r
}

// Exec a raw sql and return records as ResultMap
func (session *Session) Query() *ResultMap {
	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}
	sql := session.Statement.RawSQL
	params := session.Statement.RawParams
	i := len(params)

	var result []map[string]interface{}
	var err error
	if i == 1 {
		vv := reflect.ValueOf(params[0])
		if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Map {
			result, err = session.queryAll(sql, params...)
		} else {
			result, err = session.queryAllByMap(sql, params[0])
		}
	} else {
		result, err = session.queryAll(sql, params...)
	}
	r := &ResultMap{Results: result, Error: err}
	return r
}

// Exec a raw sql and return records as ResultMap
func (session *Session) QueryWithDateFormat(dateFormat string) *ResultMap {
	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}
	sql := session.Statement.RawSQL
	params := session.Statement.RawParams
	i := len(params)
	var result []map[string]interface{}
	var err error
	if i == 1 {
		vv := reflect.ValueOf(params[0])
		if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Map {
			result, err = session.queryAllWithDateFormat(dateFormat, sql, params...)
		} else {
			result, err = session.queryAllByMapWithDateFormat(dateFormat, sql, params[0])
		}
	} else {
		result, err = session.queryAllWithDateFormat(dateFormat, sql, params...)
	}
	r := &ResultMap{Results: result, Error: err}
	return r
}

// Execute raw sql
func (session *Session) Execute() (sql.Result, error) {
	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	sqlStr := session.Statement.RawSQL
	params := session.Statement.RawParams

	i := len(params)
	if i == 1 {
		vv := reflect.ValueOf(params[0])
		if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Map {
			return session.exec(sqlStr, params...)
		} else {
			sqlStr1, args, _ := core.MapToSlice(sqlStr, params[0])
			return session.exec(sqlStr1, args...)
		}
	} else {
		return session.exec(sqlStr, params...)
	}
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

	session.Statement.setRefValue(dataStruct)

	return session._row2BeanWithDateFormat(dateFormat, rows, fields, fieldsCount, bean, &dataStruct, session.Statement.RefTable)
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

	defer func() {
		if b, hasAfterSet := bean.(AfterSetProcessor); hasAfterSet {
			for ii, key := range fields {
				b.AfterSet(key, Cell(scanResults[ii].(*interface{})))
			}
		}
	}()

	var tempMap = make(map[string]int)
	for ii, key := range fields {
		var idx int
		var ok bool
		var lKey = strings.ToLower(key)
		if idx, ok = tempMap[lKey]; !ok {
			idx = 0
		} else {
			idx = idx + 1
		}
		tempMap[lKey] = idx

		if fieldValue := session.getField(dataStruct, key, table, idx); fieldValue != nil {
			rawValue := reflect.Indirect(reflect.ValueOf(scanResults[ii]))

			// if row is null then ignore
			if rawValue.Interface() == nil {
				continue
			}

			if fieldValue.CanAddr() {
				if structConvert, ok := fieldValue.Addr().Interface().(core.Conversion); ok {
					if data, err := value2Bytes(&rawValue); err == nil {
						structConvert.FromDB(data)
					} else {
						session.Engine.logger.Error(err)
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
					session.Engine.logger.Error(err)
				}
				continue
			}

			rawValueType := reflect.TypeOf(rawValue.Interface())
			vv := reflect.ValueOf(rawValue.Interface())

			fieldType := fieldValue.Type()
			hasAssigned := false
			col := table.GetColumnIdx(key, idx)

			if col.SQLType.IsJson() {
				var bs []byte
				if rawValueType.Kind() == reflect.String {
					bs = []byte(vv.String())
				} else if rawValueType.ConvertibleTo(core.BytesType) {
					bs = vv.Bytes()
				} else {
					return fmt.Errorf("unsupported database data type: %s %v", key, rawValueType.Kind())
				}

				hasAssigned = true

				if len(bs) > 0 {
					if fieldValue.CanAddr() {
						err := json.Unmarshal(bs, fieldValue.Addr().Interface())
						if err != nil {
							session.Engine.logger.Error(key, err)
							return err
						}
					} else {
						x := reflect.New(fieldType)
						err := json.Unmarshal(bs, x.Interface())
						if err != nil {
							session.Engine.logger.Error(key, err)
							return err
						}
						fieldValue.Set(x.Elem())
					}
				}

				continue
			}

			switch fieldType.Kind() {
			case reflect.Complex64, reflect.Complex128:
				// TODO: reimplement this
				var bs []byte
				if rawValueType.Kind() == reflect.String {
					bs = []byte(vv.String())
				} else if rawValueType.ConvertibleTo(core.BytesType) {
					bs = vv.Bytes()
				}

				hasAssigned = true
				if len(bs) > 0 {
					if fieldValue.CanAddr() {
						err := json.Unmarshal(bs, fieldValue.Addr().Interface())
						if err != nil {
							session.Engine.logger.Error(err)
							return err
						}
					} else {
						x := reflect.New(fieldType)
						err := json.Unmarshal(bs, x.Interface())
						if err != nil {
							session.Engine.logger.Error(err)
							return err
						}
						fieldValue.Set(x.Elem())
					}
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
							dbTZ := session.Engine.DatabaseTZ
							if dbTZ == nil {
								dbTZ = time.Local
							}
							session.Engine.logger.Debugf("empty zone key[%v] : %v | zone: %v | location: %+v\n", key, t, z, *t.Location())
							t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),
								t.Minute(), t.Second(), t.Nanosecond(), dbTZ)
						}
						// !nashtsai! convert to engine location
						var tz *time.Location
						if col.TimeZone == nil {
							t = t.In(session.Engine.TZLocation)
							tz = session.Engine.TZLocation
						} else {
							t = t.In(col.TimeZone)
							tz = col.TimeZone
						}
						// dateFormat to string
						//loc, _ := time.LoadLocation("Local") //重要：获取时区  rawValue.Interface().(time.Time).Format(dateFormat)
						t, _ = time.ParseInLocation(dateFormat, t.Format(dateFormat), tz)

						fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
					} else if rawValueType == core.IntType || rawValueType == core.Int64Type ||
						rawValueType == core.Int32Type {
						hasAssigned = true
						var tz *time.Location
						if col.TimeZone == nil {
							tz = session.Engine.TZLocation
						} else {
							tz = col.TimeZone
						}
						t := time.Unix(vv.Int(), 0).In(tz)
						//vv = reflect.ValueOf(t)
						fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
					} else {
						if d, ok := vv.Interface().([]uint8); ok {
							hasAssigned = true
							t, err := session.byte2Time(col, d)
							if err != nil {
								session.Engine.logger.Error("byte2Time error:", err.Error())
								hasAssigned = false
							} else {
								fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
							}
						} else if d, ok := vv.Interface().(string); ok {
							hasAssigned = true
							t, err := session.str2Time(col, d)
							if err != nil {
								session.Engine.logger.Error("byte2Time error:", err.Error())
								hasAssigned = false
							} else {
								fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
							}
						} else {
							panic(fmt.Sprintf("rawValueType is %v, value is %v", rawValueType, vv.Interface()))
						}
					}
				} else if nulVal, ok := fieldValue.Addr().Interface().(sql.Scanner); ok {
					// !<winxxp>! 增加支持sql.Scanner接口的结构，如sql.NullString
					hasAssigned = true
					if err := nulVal.Scan(vv.Interface()); err != nil {
						//fmt.Println("sql.Sanner error:", err.Error())
						session.Engine.logger.Error("sql.Sanner error:", err.Error())
						hasAssigned = false
					}
				} else if col.SQLType.IsJson() {
					if rawValueType.Kind() == reflect.String {
						hasAssigned = true
						x := reflect.New(fieldType)
						if len([]byte(vv.String())) > 0 {
							err := json.Unmarshal([]byte(vv.String()), x.Interface())
							if err != nil {
								session.Engine.logger.Error(err)
								return err
							}
							fieldValue.Set(x.Elem())
						}
					} else if rawValueType.Kind() == reflect.Slice {
						hasAssigned = true
						x := reflect.New(fieldType)
						if len(vv.Bytes()) > 0 {
							err := json.Unmarshal(vv.Bytes(), x.Interface())
							if err != nil {
								session.Engine.logger.Error(err)
								return err
							}
							fieldValue.Set(x.Elem())
						}
					}
				} else if session.Statement.UseCascade {
					table := session.Engine.autoMapType(*fieldValue)
					if table != nil {
						hasAssigned = true
						if len(table.PrimaryKeys) != 1 {
							panic("unsupported non or composited primary key cascade")
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
						case reflect.Slice:
							pk[0], _ = strconv.ParseInt(string(rawValue.Interface().([]byte)), 10, 64)
						default:
							panic(fmt.Sprintf("unsupported primary key type: %v, %v", rawValueType, fieldValue))
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
								//v := structInter.Elem().Interface()
								//fieldValue.Set(reflect.ValueOf(v))
								fieldValue.Set(structInter.Elem())
							} else {
								return errors.New("cascade obj is not exist!")
							}
						}
					} else {
						session.Engine.logger.Error("unsupported struct type in Scan: ", fieldValue.Type().String())
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
						var x = rawValue.Interface().(time.Time)
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
						var x = uint64(vv.Int())
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
						var x = float32(vv.Float())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrIntType:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt32Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int32(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt8Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int8(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrInt16Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = int16(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUintType:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.PtrUint32Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint32(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Uint8Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint8(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Uint16Type:
					if rawValueType.Kind() == reflect.Int64 {
						var x = uint16(vv.Int())
						hasAssigned = true
						fieldValue.Set(reflect.ValueOf(&x))
					}
				case core.Complex64Type:
					var x complex64
					if len([]byte(vv.String())) > 0 {
						err := json.Unmarshal([]byte(vv.String()), &x)
						if err != nil {
							session.Engine.logger.Error(err)
						} else {
							fieldValue.Set(reflect.ValueOf(&x))
						}
					}
					hasAssigned = true
				case core.Complex128Type:
					var x complex128
					if len([]byte(vv.String())) > 0 {
						err := json.Unmarshal([]byte(vv.String()), &x)
						if err != nil {
							session.Engine.logger.Error(err)
						} else {
							fieldValue.Set(reflect.ValueOf(&x))
						}
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
					session.bytes2Value(col, fieldValue, data)
				} else {
					session.Engine.logger.Error(err.Error())
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

func (session *Session) Sqls(sqls interface{}, parmas ...interface{}) *SqlsExecutor {

	sqlsExecutor := new(SqlsExecutor)
	switch sqls.(type) {
	case string:
		sqlsExecutor.sqls = sqls.(string)
	case []string:
		sqlsExecutor.sqls = sqls.([]string)
	case map[string]string:
		sqlsExecutor.sqls = sqls.(map[string]string)
	default:
		sqlsExecutor.sqls = nil
		sqlsExecutor.err = ErrParamsType
	}

	if len(parmas) == 0 {
		sqlsExecutor.parmas = nil
	}

	if len(parmas) > 1 {
		sqlsExecutor.parmas = nil
		sqlsExecutor.err = ErrParamsType
	}

	if len(parmas) == 1 {
		switch parmas[0].(type) {
		case map[string]interface{}:
			sqlsExecutor.parmas = parmas[0].(map[string]interface{})

		case []map[string]interface{}:
			sqlsExecutor.parmas = parmas[0].([]map[string]interface{})

		case map[string]map[string]interface{}:
			sqlsExecutor.parmas = parmas[0].(map[string]map[string]interface{})
		default:
			sqlsExecutor.parmas = nil
			sqlsExecutor.err = ErrParamsType
		}
	}

	sqlsExecutor.session = session

	return sqlsExecutor
}

func (session *Session) SqlMapsClient(sqlkeys interface{}, parmas ...interface{}) *SqlMapsExecutor {
	sqlMapsExecutor := new(SqlMapsExecutor)

	switch sqlkeys.(type) {
	case string:
		sqlMapsExecutor.sqlkeys = sqlkeys.(string)
	case []string:
		sqlMapsExecutor.sqlkeys = sqlkeys.([]string)
	case map[string]string:
		sqlMapsExecutor.sqlkeys = sqlkeys.(map[string]string)
	default:
		sqlMapsExecutor.sqlkeys = nil
		sqlMapsExecutor.err = ErrParamsType
	}

	if len(parmas) == 0 {
		sqlMapsExecutor.parmas = nil
	}

	if len(parmas) > 1 {
		sqlMapsExecutor.parmas = nil
		sqlMapsExecutor.err = ErrParamsType
	}

	if len(parmas) == 1 {
		switch parmas[0].(type) {
		case map[string]interface{}:
			sqlMapsExecutor.parmas = parmas[0].(map[string]interface{})

		case []map[string]interface{}:
			sqlMapsExecutor.parmas = parmas[0].([]map[string]interface{})

		case map[string]map[string]interface{}:
			sqlMapsExecutor.parmas = parmas[0].(map[string]map[string]interface{})
		default:
			sqlMapsExecutor.parmas = nil
			sqlMapsExecutor.err = ErrParamsType
		}
	}

	sqlMapsExecutor.session = session

	return sqlMapsExecutor
}

func (session *Session) SqlTemplatesClient(sqlkeys interface{}, parmas ...interface{}) *SqlTemplatesExecutor {
	sqlTemplatesExecutor := new(SqlTemplatesExecutor)

	switch sqlkeys.(type) {
	case string:
		sqlTemplatesExecutor.sqlkeys = sqlkeys.(string)
	case []string:
		sqlTemplatesExecutor.sqlkeys = sqlkeys.([]string)
	case map[string]string:
		sqlTemplatesExecutor.sqlkeys = sqlkeys.(map[string]string)
	default:
		sqlTemplatesExecutor.sqlkeys = nil
		sqlTemplatesExecutor.err = ErrParamsType
	}

	if len(parmas) == 0 {
		sqlTemplatesExecutor.parmas = nil
	}

	if len(parmas) > 1 {
		sqlTemplatesExecutor.parmas = nil
		sqlTemplatesExecutor.err = ErrParamsType
	}

	if len(parmas) == 1 {
		switch parmas[0].(type) {
		case map[string]interface{}:
			sqlTemplatesExecutor.parmas = parmas[0].(map[string]interface{})

		case []map[string]interface{}:
			sqlTemplatesExecutor.parmas = parmas[0].([]map[string]interface{})

		case map[string]map[string]interface{}:
			sqlTemplatesExecutor.parmas = parmas[0].(map[string]map[string]interface{})
		default:
			sqlTemplatesExecutor.parmas = nil
			sqlTemplatesExecutor.err = ErrParamsType
		}
	}

	sqlTemplatesExecutor.session = session

	return sqlTemplatesExecutor
}
