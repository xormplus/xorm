package xorm

import (
	"encoding/json"
)

func (engine *Engine) SetSqlMapRootDir(sqlMapRootDir string) *Engine {
	engine.sqlMap.SqlMapRootDir = sqlMapRootDir
	return engine
}

func (engine *Engine) SetSqlTemplateRootDir(sqlTemplateRootDir string) *Engine {
	engine.sqlTemplate.SqlTemplateRootDir = sqlTemplateRootDir
	return engine
}

func (engine *Engine) SqlMapClient(sqlTagName string, args ...interface{}) *Session {
	session := engine.NewSession()
	session.IsAutoClose = true
	session.IsSqlFuc = true
	return session.Sql(engine.sqlMap.Sql[sqlTagName], args...)
}

func (engine *Engine) SqlTemplateClient(sqlTagName string, args ...interface{}) *Session {
	session := engine.NewSession()
	session.IsAutoClose = true
	session.IsSqlFuc = true
	return session.SqlTemplateClient(sqlTagName, args...)

}

func (engine *Engine) Search(beans interface{}, condiBeans ...interface{}) *ResultStructs {
	session := engine.NewSession()
	defer session.Close()
	return session.Search(beans, condiBeans...)
}

// Get retrieve one record from table, bean's non-empty fields
// are conditions
func (engine *Engine) GetFirst(bean interface{}) *ResultBean {
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

func (engine *Engine) Sqls(sqls interface{}, parmas ...interface{}) *SqlsExecutor {
	session := engine.NewSession()
	session.IsAutoClose = true
	session.IsSqlFuc = true
	return session.Sqls(sqls, parmas...)
}

func (engine *Engine) SqlMapsClient(sqlkeys interface{}, parmas ...interface{}) *SqlMapsExecutor {
	session := engine.NewSession()
	session.IsAutoClose = true
	session.IsSqlFuc = true
	return session.SqlMapsClient(sqlkeys, parmas...)
}

func (engine *Engine) SqlTemplatesClient(sqlkeys interface{}, parmas ...interface{}) *SqlTemplatesExecutor {
	session := engine.NewSession()
	session.IsAutoClose = true
	session.IsSqlFuc = true
	return session.SqlTemplatesClient(sqlkeys, parmas...)
}
