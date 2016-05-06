package xorm

import (
	"encoding/json"
)

func (engine *Engine) SqlMapClient(sqlTagName string, args ...interface{}) *Session {
	session := engine.NewSession()
	session.IsAutoClose = true
	return session.Sql(engine.SqlMap.Sql[sqlTagName], args...)
}

func (engine *Engine) SqlTemplateClient(sqlTagName string, args ...interface{}) *Session {
	session := engine.NewSession()
	session.IsAutoClose = true
	map1 := args[0].(map[string]interface{})
	if engine.SqlTemplate.Template[sqlTagName] == nil {
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

	if string(result) == "null" {
		return "", nil
	}
	return string(result), nil
}
