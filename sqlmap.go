package xorm

import (
	"encoding/xml"
	"io/ioutil"

	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
)

type SqlMap struct {
	SqlMapRootDir string
	Sql           map[string]string
	Extension     string
	Capacity      uint
	Cipher        Cipher
}

type SqlMapOptions struct {
	Capacity  uint
	Extension string
	Cipher    Cipher
}

type Result struct {
	Sql []Sql `xml:"sql"`
}

type Sql struct {
	Value string `xml:",chardata"`
	Id    string `xml:"id,attr"`
}

func (engine *Engine) SetSqlMapCipher(cipher Cipher) {
	engine.sqlMap.Cipher = cipher
}

func (engine *Engine) ClearSqlMapCipher() {
	engine.sqlMap.Cipher = nil
}

func (sqlMap *SqlMap) checkNilAndInit() {
	if sqlMap.Sql == nil {
		if sqlMap.Capacity == 0 {
			sqlMap.Sql = make(map[string]string, 100)
		} else {
			sqlMap.Sql = make(map[string]string, sqlMap.Capacity)
		}

	}
}

func (engine *Engine) InitSqlMap(options ...SqlMapOptions) error {
	var opt SqlMapOptions

	if len(options) > 0 {
		opt = options[0]
	}

	if len(opt.Extension) == 0 {
		opt.Extension = ".xml"
	}

	engine.sqlMap.Extension = opt.Extension
	engine.sqlMap.Capacity = opt.Capacity

	engine.sqlMap.Cipher = opt.Cipher

	var err error
	if engine.sqlMap.SqlMapRootDir == "" {
		cfg, err := goconfig.LoadConfigFile("./sql/xormcfg.ini")
		if err != nil {
			return err
		}
		engine.sqlMap.SqlMapRootDir, err = cfg.GetValue("", "SqlMapRootDir")
		if err != nil {
			return err
		}
	}

	err = filepath.Walk(engine.sqlMap.SqlMapRootDir, engine.sqlMap.walkFunc)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) LoadSqlMap(filepath string) error {
	if len(engine.sqlMap.Extension) == 0 {
		engine.sqlMap.Extension = ".xml"
	}
	if strings.HasSuffix(filepath, engine.sqlMap.Extension) {
		err := engine.loadSqlMap(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) BatchLoadSqlMap(filepathSlice []string) error {
	if len(engine.sqlMap.Extension) == 0 {
		engine.sqlMap.Extension = ".xml"
	}
	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, engine.sqlMap.Extension) {
			err := engine.loadSqlMap(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (engine *Engine) ReloadSqlMap(filepath string) error {
	if len(engine.sqlMap.Extension) == 0 {
		engine.sqlMap.Extension = ".xml"
	}
	if strings.HasSuffix(filepath, engine.sqlMap.Extension) {
		err := engine.reloadSqlMap(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) BatchReloadSqlMap(filepathSlice []string) error {
	if len(engine.sqlMap.Extension) == 0 {
		engine.sqlMap.Extension = ".xml"
	}
	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, engine.sqlMap.Extension) {
			err := engine.loadSqlMap(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (engine *Engine) loadSqlMap(filepath string) error {

	info, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	err = engine.sqlMap.paresSql(filepath)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) reloadSqlMap(filepath string) error {

	info, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}
	err = engine.sqlMap.paresSql(filepath)
	if err != nil {
		return err
	}

	return nil
}

func (sqlMap *SqlMap) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if strings.HasSuffix(path, sqlMap.Extension) {
		err = sqlMap.paresSql(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sqlMap *SqlMap) paresSql(filepath string) error {

	content, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}
	enc := sqlMap.Cipher
	if enc != nil {
		content, err = enc.Decrypt(content)

		if err != nil {
			return err
		}
	}

	sqlMap.checkNilAndInit()
	var result Result
	err = xml.Unmarshal(content, &result)
	if err != nil {
		return err
	}

	for _, sql := range result.Sql {
		sqlMap.Sql[sql.Id] = sql.Value
	}

	return nil

}

func (engine *Engine) AddSql(key string, sql string) {
	engine.sqlMap.addSql(key, sql)
}

func (sqlMap *SqlMap) addSql(key string, sql string) {
	sqlMap.checkNilAndInit()
	sqlMap.Sql[key] = sql
}

func (engine *Engine) UpdateSql(key string, sql string) {
	engine.sqlMap.updateSql(key, sql)
}

func (sqlMap *SqlMap) updateSql(key string, sql string) {
	sqlMap.checkNilAndInit()
	sqlMap.Sql[key] = sql
}

func (engine *Engine) RemoveSql(key string) {
	engine.sqlMap.removeSql(key)
}

func (sqlMap *SqlMap) removeSql(key string) {
	sqlMap.checkNilAndInit()
	delete(sqlMap.Sql, key)
}

func (engine *Engine) BatchAddSql(sqlStrMap map[string]string) {
	engine.sqlMap.batchAddSql(sqlStrMap)
}

func (sqlMap *SqlMap) batchAddSql(sqlStrMap map[string]string) {
	sqlMap.checkNilAndInit()
	for k, v := range sqlStrMap {
		sqlMap.Sql[k] = v
	}
}

func (engine *Engine) BatchUpdateSql(sqlStrMap map[string]string) {
	engine.sqlMap.batchUpdateSql(sqlStrMap)
}

func (sqlMap *SqlMap) batchUpdateSql(sqlStrMap map[string]string) {
	sqlMap.checkNilAndInit()
	for k, v := range sqlStrMap {
		sqlMap.Sql[k] = v
	}
}

func (engine *Engine) BatchRemoveSql(key []string) {
	engine.sqlMap.batchRemoveSql(key)
}

func (sqlMap *SqlMap) batchRemoveSql(key []string) {
	sqlMap.checkNilAndInit()
	for _, v := range key {
		delete(sqlMap.Sql, v)
	}
}

func (engine *Engine) GetSql(key string) string {
	return engine.sqlMap.getSql(key)
}

func (sqlMap *SqlMap) getSql(key string) string {
	return sqlMap.Sql[key]
}

func (engine *Engine) GetSqlMap(keys ...interface{}) map[string]string {
	return engine.sqlMap.getSqlMap(keys...)
}

func (sqlMap *SqlMap) getSqlMap(keys ...interface{}) map[string]string {
	var resultSqlMap map[string]string
	i := len(keys)
	if i == 0 {
		return sqlMap.Sql
	}

	if i == 1 {
		switch keys[0].(type) {
		case string:
			resultSqlMap = make(map[string]string, 1)
		case []string:
			ks := keys[0].([]string)
			n := len(ks)
			resultSqlMap = make(map[string]string, n)
		}
	} else {
		resultSqlMap = make(map[string]string, i)
	}

	for k, _ := range keys {
		switch keys[k].(type) {
		case string:
			key := keys[k].(string)
			resultSqlMap[key] = sqlMap.Sql[key]
		case []string:
			ks := keys[k].([]string)
			for _, v := range ks {
				resultSqlMap[v] = sqlMap.Sql[v]
			}
		}
	}

	return resultSqlMap
}
