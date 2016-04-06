package xorm

import (
	"encoding/json"

	"github.com/Chronokeeper/anyxml"
)

func (engine *Engine) SqlMapClient(sqlTagName string, args ...interface{}) *Session {
	session := engine.NewSession()
	session.IsAutoClose = true
	return session.Sql(engine.SqlMap.Sql[sqlTagName], args...)
}

func (engine *Engine) SqlTemplateClient(sqlTagName string, args ...interface{}) *Session {
	session := engine.NewSession()
	session.IsAutoClose = true
	map1:=args[0].(map[string]interface{})
	if engine.SqlTemplate.Template[sqlTagName]==nil{
		return session.Sql("", &map1)
	}
	sql, err := engine.SqlTemplate.Template[sqlTagName].Execute(map1)
	if err != nil {
		engine.logger.Error(err)
	}

	return session.Sql(sql, &map1)
}

// Get retrieve one record from table, bean's non-empty fields
// are conditions
func (engine *Engine) GetFirst(bean interface{}) ResultBean {
	session := engine.NewSession()
	defer session.Close()
	return session.GetFirst(bean)
}

// Exec a raw sql and return records as []map[string]interface{}
func (engine *Engine) QueryAll(sql string, paramStr ...interface{}) (resultsSlice []map[string]interface{}, err error) {
	session := engine.NewSession()
	defer session.Close()
	return session.queryAll(sql, paramStr...)
}

// Exec a raw sql and return records as []map[string]interface{}
func (engine *Engine) QueryAllByMap(sql string, paramMap interface{}) (resultsSlice []map[string]interface{}, err error) {
	session := engine.NewSession()
	defer session.Close()
	return session.queryAllByMap(sql, paramMap)
}

func (engine *Engine) QueryAllByMapToJsonString(sql string, paramMap interface{}) (string, error) {
	session := engine.NewSession()
	defer session.Close()
	results, err := session.queryAllByMap(sql, paramMap)
	if err != nil {
		return "", err
	}
	return JSONString(results, true)
}

func (engine *Engine) QueryAllByMapToJsonStringWithDateFormat(dateFormat string, sql string, paramMap interface{}) (string, error) {
	session := engine.NewSession()
	defer session.Close()
	results, err := session.queryAllByMapWithDateFormat(dateFormat, sql, paramMap)
	if err != nil {
		return "", err
	}
	return JSONString(results, true)
}

func (engine *Engine) QueryAllToJsonString(sql string, paramStr ...interface{}) (string, error) {
	session := engine.NewSession()
	defer session.Close()
	results, err := session.queryAll(sql, paramStr...)
	if err != nil {
		return "", err
	}
	return JSONString(results, true)
}

func (engine *Engine) QueryAllToJsonStringWithDateFormat(dateFormat string, sql string, paramStr ...interface{}) (string, error) {
	session := engine.NewSession()
	defer session.Close()
	results, err := session.queryAllWithDateFormat(dateFormat, sql, paramStr...)
	if err != nil {
		return "", err
	}
	return JSONString(results, true)
}

func (engine *Engine) queryAllToXmlString(sql string, paramStr ...interface{}) (string, error) {
	session := engine.NewSession()
	defer session.Close()
	resultSlice, err := session.queryAll(sql, paramStr...)
	if err != nil {
		return "", err
	}

	results, err := anyxml.Xml(resultSlice, "result")
	if err != nil {
		return "", err
	}
	return string(results), nil
}

func (engine *Engine) queryAllToXmlIndentString(sql string, prefix string, indent string, paramStr ...interface{}) (string, error) {
	session := engine.NewSession()
	defer session.Close()
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

func (engine *Engine) queryAllToXmlStringWithDateFormat(dateFormat string, sql string, paramStr ...interface{}) (string, error) {
	session := engine.NewSession()
	defer session.Close()
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

func (engine *Engine) QueryAllToXmlIndentStringWithDateFormat(dateFormat string, sql string, prefix string, indent string, paramStr ...interface{}) (string, error) {
	session := engine.NewSession()
	defer session.Close()
	resultSlice, err := session.queryAll(sql, paramStr...)
	if err != nil {
		return "", err
	}
	results, err := anyxml.XmlIndentWithDateFormat(dateFormat, resultSlice, "", "  ", "results")

	if err != nil {
		return "", err
	}
	return string(results), nil
}

func JSONString(v interface{}, IndentJSON bool) (string, error) {
	var result []byte
	var err error
	if IndentJSON {
		result, err = json.MarshalIndent(v, "", "  ")
	} else {
		result, err = json.Marshal(v)
	}
	if err != nil {
		return "", err
	}

	if string(result)=="null"{
		return "", nil
	}
	return string(result), nil
}
