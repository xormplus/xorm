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
}

type SqlMapOptions struct {
	Extension string
}

type Result struct {
	Sql []Sql `xml:"sql"`
}

type Sql struct {
	Value string `xml:",chardata"`
	Id    string `xml:"id,attr"`
}

func (engine *Engine) InitSqlMap(options ...SqlMapOptions) error {
	var opt SqlMapOptions
	engine.SqlMap.Sql = make(map[string]string, 100)
	if len(options) > 0 {
		opt = options[0]
	}

	if len(opt.Extension) == 0 {
		opt.Extension = ".xml"
	}

	engine.SqlMap.Extension = opt.Extension

	var err error
	if engine.SqlMap.SqlMapRootDir == "" {
		cfg, err := goconfig.LoadConfigFile("./sql/xormcfg.ini")
		if err != nil {
			return err
		}
		engine.SqlMap.SqlMapRootDir, err = cfg.GetValue("", "SqlMapRootDir")
		if err != nil {
			return err
		}
	}

	err = filepath.Walk(engine.SqlMap.SqlMapRootDir, engine.SqlMap.walkFunc)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) LoadSqlMap(filepath string) error {
	if strings.HasSuffix(filepath, engine.SqlMap.Extension) {
		err := engine.loadSqlMap(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) ReloadSqlMap(filepath string) error {
	if strings.HasSuffix(filepath, engine.SqlMap.Extension) {
		err := engine.reloadSqlMap(filepath)
		if err != nil {
			return err
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

	err = engine.SqlMap.paresSql(filepath)
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
	err = engine.SqlMap.paresSql(filepath)
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
	engine.SqlMap.addSql(key, sql)
}

func (sqlMap *SqlMap) addSql(key string, sql string) {
	sqlMap.Sql[key] = sql
}

func (engine *Engine) UpdateSql(key string, sql string) {
	engine.SqlMap.updateSql(key, sql)
}

func (sqlMap *SqlMap) updateSql(key string, sql string) {
	sqlMap.Sql[key] = sql
}

func (engine *Engine) RemoveSql(key string) {
	engine.SqlMap.removeSql(key)
}

func (sqlMap *SqlMap) removeSql(key string) {
	delete(sqlMap.Sql, key)
}

func (engine *Engine) BatchAddSql(sqlStrMap map[string]string) {
	engine.SqlMap.batchAddSql(sqlStrMap)
}

func (sqlMap *SqlMap) batchAddSql(sqlStrMap map[string]string) {
	for k, v := range sqlStrMap {
		sqlMap.Sql[k] = v
	}
}

func (engine *Engine) BatchUpdateSql(sqlStrMap map[string]string) {
	engine.SqlMap.batchUpdateSql(sqlStrMap)
}

func (sqlMap *SqlMap) batchUpdateSql(sqlStrMap map[string]string) {
	for k, v := range sqlStrMap {
		sqlMap.Sql[k] = v
	}
}

func (engine *Engine) BatchRemoveSql(key []string) {
	engine.SqlMap.batchRemoveSql(key)
}

func (sqlMap *SqlMap) batchRemoveSql(key []string) {
	for _, v := range key {
		delete(sqlMap.Sql, v)
	}
}
