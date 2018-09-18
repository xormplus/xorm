// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/Chronokeeper/anyxml"
	"github.com/xormplus/core"
)

type Record map[string]Value
type Result []Record

type ResultValue struct {
	Result Result
	Error  error
}

func (resultValue *ResultValue) List() (Result, error) {
	return resultValue.Result, resultValue.Error
}

func (resultValue *ResultValue) Count() (int, error) {
	if resultValue.Error != nil {
		return 0, resultValue.Error
	}
	if resultValue.Result == nil {
		return 0, nil
	}
	return len(resultValue.Result), nil
}

func (resultValue *ResultValue) ListPage(firstResult int, maxResults int) (Result, error) {
	if resultValue.Error != nil {
		return nil, resultValue.Error
	}
	if resultValue.Result == nil {
		return nil, nil
	}
	if firstResult > maxResults {
		return nil, ErrParamsFormat
	}
	if firstResult < 0 {
		return nil, ErrParamsFormat
	}
	if maxResults < 0 {
		return nil, ErrParamsFormat
	}
	if maxResults > len(resultValue.Result) {
		return nil, ErrParamsFormat
	}
	return resultValue.Result[(firstResult - 1):maxResults], resultValue.Error
}

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
	Result []map[string]interface{}
	Error  error
}

func (resultMap *ResultMap) List() ([]map[string]interface{}, error) {
	return resultMap.Result, resultMap.Error
}

func (resultMap *ResultMap) Count() (int, error) {
	if resultMap.Error != nil {
		return 0, resultMap.Error
	}
	if resultMap.Result == nil {
		return 0, nil
	}
	return len(resultMap.Result), nil
}

func (resultMap *ResultMap) ListPage(firstResult int, maxResults int) ([]map[string]interface{}, error) {
	if resultMap.Error != nil {
		return nil, resultMap.Error
	}
	if resultMap.Result == nil {
		return nil, nil
	}
	if firstResult > maxResults {
		return nil, ErrParamsFormat
	}
	if firstResult < 0 {
		return nil, ErrParamsFormat
	}
	if maxResults < 0 {
		return nil, ErrParamsFormat
	}
	if maxResults > len(resultMap.Result) {
		return nil, ErrParamsFormat
	}
	return resultMap.Result[(firstResult - 1):maxResults], resultMap.Error
}

func (resultMap *ResultMap) Json() (string, error) {

	if resultMap.Error != nil {
		return "", resultMap.Error
	}
	return JSONString(resultMap.Result, true)
}

func (resultMap *ResultMap) Xml() (string, error) {
	if resultMap.Error != nil {
		return "", resultMap.Error
	}
	results, err := anyxml.Xml(resultMap.Result)
	if err != nil {
		return "", err
	}
	return string(results), nil
}

func (resultMap *ResultMap) XmlIndent(prefix string, indent string, recordTag string) (string, error) {
	if resultMap.Error != nil {
		return "", resultMap.Error
	}

	results, err := anyxml.XmlIndent(resultMap.Result, prefix, indent, recordTag)
	if err != nil {
		return "", err
	}
	return string(results), nil
}

func (resultMap *ResultMap) SaveAsCSV(filename string, headers []string, perm os.FileMode) error {
	if resultMap.Error != nil {
		return resultMap.Error
	}

	dataset, err := NewDatasetWithData(headers, resultMap.Result, true)
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

	dataset, err := NewDatasetWithData(headers, resultMap.Result, true)
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

	dataset, err := NewDatasetWithData(headers, resultMap.Result, true)
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

	dataset, err := NewDatasetWithData(headers, resultMap.Result, false)
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

	dataset, err := NewDatasetWithData(headers, resultMap.Result, false)
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

	dataset, err := NewDatasetWithData(headers, resultMap.Result, false)
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

	dataset, err := NewDatasetWithData(headers, resultMap.Result, false)
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

	dataset, err := NewDatasetWithData(headers, resultMap.Result, true)
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
	return session.Sql(session.engine.SqlMap.Sql[sqlTagName], args...)
}

func (session *Session) SqlTemplateClient(sqlTagName string, args ...interface{}) *Session {
	session.isSqlFunc = true
	sql, err := session.engine.SqlTemplate.Execute(sqlTagName, args...)
	if err != nil {
		session.engine.logger.Error(err)
	}

	if len(args) == 0 {
		return session.Sql(sql)
	} else {
		map1 := args[0].(*map[string]interface{})
		return session.Sql(sql, map1)
	}

}

func (session *Session) Search(rowsSlicePtr interface{}, condiBean ...interface{}) *ResultStructs {
	err := session.Find(rowsSlicePtr, condiBean...)
	r := &ResultStructs{Result: rowsSlicePtr, Error: err}
	return r
}

func (session *Session) genSelectSql(dialect core.Dialect, rownumber string) string {

	var sql = session.statement.RawSQL
	var orderBys = session.statement.OrderStr

	if dialect.DBType() != core.MSSQL && dialect.DBType() != core.ORACLE {
		if session.statement.Start > 0 {
			sql = fmt.Sprintf("%v LIMIT %v OFFSET %v", sql, session.statement.LimitN, session.statement.Start)
		} else if session.statement.LimitN > 0 {
			sql = fmt.Sprintf("%v LIMIT %v", sql, session.statement.LimitN)
		}
	} else if dialect.DBType() == core.ORACLE {
		if session.statement.Start != 0 || session.statement.LimitN != 0 {
			sql = fmt.Sprintf("SELECT aat.* FROM (SELECT at.*,ROWNUM %v FROM (%v) at WHERE ROWNUM <= %d) aat WHERE %v > %d",
				rownumber, sql, session.statement.Start+session.statement.LimitN, rownumber, session.statement.Start)
		}
	} else {
		keepSelect := false
		var fullQuery string
		if session.statement.Start > 0 {
			fullQuery = fmt.Sprintf("SELECT sq.* FROM (SELECT ROW_NUMBER() OVER (ORDER BY %v) AS %v,", orderBys, rownumber)
		} else if session.statement.LimitN > 0 {
			fullQuery = fmt.Sprintf("SELECT TOP %d", session.statement.LimitN)
		} else {
			keepSelect = true
		}

		if !keepSelect {
			expr := `^\s*SELECT\s*`
			reg, err := regexp.Compile(expr)
			if err != nil {
				fmt.Println(err)
			}
			sql = strings.ToUpper(sql)
			if reg.MatchString(sql) {
				str := reg.FindAllString(sql, -1)
				fullQuery = fmt.Sprintf("%v %v", fullQuery, sql[len(str[0]):])
			}
		}

		if session.statement.Start > 0 {
			// T-SQL offset starts with 1, not like MySQL with 0;
			if session.statement.LimitN > 0 {
				fullQuery = fmt.Sprintf("%v) AS sq WHERE %v BETWEEN %d AND %d", fullQuery, rownumber, session.statement.Start+1, session.statement.Start+session.statement.LimitN)
			} else {
				fullQuery = fmt.Sprintf("%v) AS sq WHERE %v >= %d", fullQuery, rownumber, session.statement.Start+1)
			}
		} else {
			fullQuery = fmt.Sprintf("%v ORDER BY %v", fullQuery, orderBys)
		}

		if keepSelect {
			if len(orderBys) > 0 {
				sql = fmt.Sprintf("%v ORDER BY %v", sql, orderBys)
			}
		} else {
			sql = fullQuery
		}
	}

	return sql
}

// Exec a raw sql and return records as ResultMap
func (session *Session) Query() *ResultMap {
	defer session.resetStatement()
	if session.isAutoClose {
		defer session.Close()
	}

	var dialect = session.statement.Engine.Dialect()
	rownumber := "xorm" + NewShortUUID().String()
	sql := session.genSelectSql(dialect, rownumber)

	params := session.statement.RawParams
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

	if dialect.DBType() == core.MSSQL {
		if session.statement.Start > 0 {
			for i, _ := range result {
				delete(result[i], rownumber)
			}
		}
	} else if dialect.DBType() == core.ORACLE {
		if session.statement.Start != 0 || session.statement.LimitN != 0 {
			for i, _ := range result {
				delete(result[i], rownumber)
			}
		}
	}
	r := &ResultMap{Result: result, Error: err}
	return r
}

// Exec a raw sql and return records as ResultMap
func (session *Session) QueryWithDateFormat(dateFormat string) *ResultMap {
	defer session.resetStatement()
	if session.isAutoClose {
		defer session.Close()
	}

	var dialect = session.statement.Engine.Dialect()
	rownumber := "xorm" + NewShortUUID().String()
	sql := session.genSelectSql(dialect, rownumber)

	params := session.statement.RawParams
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

	if dialect.DBType() == core.MSSQL {
		if session.statement.Start > 0 {
			for i, _ := range result {
				delete(result[i], rownumber)
			}
		}
	} else if dialect.DBType() == core.ORACLE {
		if session.statement.Start != 0 || session.statement.LimitN != 0 {
			for i, _ := range result {
				delete(result[i], rownumber)
			}
		}
	}
	r := &ResultMap{Result: result, Error: err}
	return r
}

// Execute raw sql
func (session *Session) Execute() (sql.Result, error) {
	defer session.resetStatement()
	if session.isAutoClose {
		defer session.Close()
	}

	sqlStr := session.statement.RawSQL
	params := session.statement.RawParams

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

	if session.isAutoCommit {
		return query3(session.DB(), sqlStr, paramStr...)
	}
	return txQuery3(session.tx, sqlStr, paramStr...)
}

func (session *Session) queryAllByMap(sqlStr string, paramMap interface{}) (resultsSlice []map[string]interface{}, err error) {
	sqlStr1, param, _ := core.MapToSlice(sqlStr, paramMap)

	session.queryPreprocess(&sqlStr1, param...)

	if session.isAutoCommit {
		return query3(session.DB(), sqlStr1, param...)
	}
	return txQuery3(session.tx, sqlStr1, param...)
}

func (session *Session) queryAllByMapWithDateFormat(dateFormat string, sqlStr string, paramMap interface{}) (resultsSlice []map[string]interface{}, err error) {
	sqlStr1, param, _ := core.MapToSlice(sqlStr, paramMap)
	session.queryPreprocess(&sqlStr1, param...)

	if session.isAutoCommit {
		return query3WithDateFormat(session.DB(), dateFormat, sqlStr1, param...)
	}
	return txQuery3WithDateFormat(session.tx, dateFormat, sqlStr1, param...)
}

func (session *Session) queryAllWithDateFormat(dateFormat string, sqlStr string, paramStr ...interface{}) (resultsSlice []map[string]interface{}, err error) {
	session.queryPreprocess(&sqlStr, paramStr...)

	if session.isAutoCommit {
		return query3WithDateFormat(session.DB(), dateFormat, sqlStr, paramStr...)
	}
	return txQuery3WithDateFormat(session.tx, dateFormat, sqlStr, paramStr...)
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

	for _, filter := range session.engine.dialect.Filters() {
		query = filter.Do(query, session.engine.dialect, session.statement.RefTable)
	}

	*sqlStr = query
	session.engine.logSQL(session, *sqlStr, paramMap)
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
