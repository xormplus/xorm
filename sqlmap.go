package xorm

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"os"
	"path/filepath"
	"strings"
)

type SqlMap struct {
	SqlMapRootDir string
	Sql           map[string]string
	Extension     map[string]string
	Capacity      uint
	Cipher        Cipher
}

type Result struct {
	Sql []Sql `xml:"sql"`
}

type Sql struct {
	Value string `xml:",chardata"`
	Id    string `xml:"id,attr"`
}

func (engine *Engine) SetSqlMapCipher(cipher Cipher) {
	engine.SqlMap.Cipher = cipher
}

func (engine *Engine) ClearSqlMapCipher() {
	engine.SqlMap.Cipher = nil
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

type SqlM interface {
	RootDir() string
	Extension() string
}

type XmlSqlMap struct {
	sqlMapRootDir string
	extension     string
}

type JsonSqlMap struct {
	sqlMapRootDir string
	extension     string
}

func Xml(directory, extension string) *XmlSqlMap {
	return &XmlSqlMap{
		sqlMapRootDir: directory,
		extension:     extension,
	}
}

func Json(directory, extension string) *JsonSqlMap {
	return &JsonSqlMap{
		sqlMapRootDir: directory,
		extension:     extension,
	}
}

func (sqlMap *XmlSqlMap) RootDir() string {
	return sqlMap.sqlMapRootDir
}

func (sqlMap *JsonSqlMap) RootDir() string {
	return sqlMap.sqlMapRootDir
}

func (sqlMap *XmlSqlMap) Extension() string {
	return sqlMap.extension
}

func (sqlMap *JsonSqlMap) Extension() string {
	return sqlMap.extension
}

func (engine *Engine) RegisterSqlMap(sqlm SqlM, Cipher ...Cipher) error {
	switch sqlm.(type) {
	case *XmlSqlMap:
		if len(engine.SqlMap.Extension) == 0 {
			engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json"}
		}
		engine.SqlMap.Extension["xml"] = sqlm.Extension()
	case *JsonSqlMap:
		if len(engine.SqlMap.Extension) == 0 {
			engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json"}
		}
		engine.SqlMap.Extension["json"] = sqlm.Extension()
	default:
		return ErrParamsType
	}

	if len(Cipher) > 0 {
		engine.SqlMap.Cipher = Cipher[0]
	}

	engine.SqlMap.SqlMapRootDir = sqlm.RootDir()

	err := filepath.Walk(engine.SqlMap.SqlMapRootDir, engine.SqlMap.walkFunc)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) LoadSqlMap(filepath string) error {

	if len(engine.SqlMap.Extension) == 0 {
		engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json"}
	} else {
		if engine.SqlMap.Extension["xml"] == "" || len(engine.SqlMap.Extension["xml"]) == 0 {
			engine.SqlMap.Extension["xml"] = ".xml"
		}
		if engine.SqlMap.Extension["json"] == "" || len(engine.SqlMap.Extension["json"]) == 0 {
			engine.SqlMap.Extension["json"] = ".json"
		}
	}

	if strings.HasSuffix(filepath, engine.SqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["json"]) {
		err := engine.loadSqlMap(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) BatchLoadSqlMap(filepathSlice []string) error {
	if len(engine.SqlMap.Extension) == 0 {
		engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json"}
	} else {
		if engine.SqlMap.Extension["xml"] == "" || len(engine.SqlMap.Extension["xml"]) == 0 {
			engine.SqlMap.Extension["xml"] = ".xml"
		}
		if engine.SqlMap.Extension["json"] == "" || len(engine.SqlMap.Extension["json"]) == 0 {
			engine.SqlMap.Extension["json"] = ".json"
		}
	}

	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, engine.SqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["json"]) {
			err := engine.loadSqlMap(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (engine *Engine) ReloadSqlMap(filepath string) error {
	if len(engine.SqlMap.Extension) == 0 {
		engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json"}
	} else {
		if engine.SqlMap.Extension["xml"] == "" || len(engine.SqlMap.Extension["xml"]) == 0 {
			engine.SqlMap.Extension["xml"] = ".xml"
		}
		if engine.SqlMap.Extension["json"] == "" || len(engine.SqlMap.Extension["json"]) == 0 {
			engine.SqlMap.Extension["json"] = ".json"
		}
	}

	if strings.HasSuffix(filepath, engine.SqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["json"]) {
		err := engine.reloadSqlMap(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) BatchReloadSqlMap(filepathSlice []string) error {
	if len(engine.SqlMap.Extension) == 0 {
		engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json"}
	} else {
		if engine.SqlMap.Extension["xml"] == "" || len(engine.SqlMap.Extension["xml"]) == 0 {
			engine.SqlMap.Extension["xml"] = ".xml"
		}
		if engine.SqlMap.Extension["json"] == "" || len(engine.SqlMap.Extension["json"]) == 0 {
			engine.SqlMap.Extension["json"] = ".json"
		}
	}

	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, engine.SqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["json"]) {
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

	if strings.HasSuffix(path, sqlMap.Extension["xml"]) || strings.HasSuffix(path, sqlMap.Extension["json"]) {
		err = sqlMap.paresSql(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sqlMap *SqlMap) paresSql(filepath string) error {

	content, err := ioutil.ReadFile(filepath)
	fmt.Println("filepath:", filepath)
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

	if strings.HasSuffix(filepath, sqlMap.Extension["xml"]) {
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

	if strings.HasSuffix(filepath, sqlMap.Extension["json"]) {
		var result map[string]string
		err = json.Unmarshal(content, &result)
		if err != nil {
			return err
		}
		for k := range result {
			sqlMap.Sql[k] = result[k]
		}

		return nil
	}
	return nil

}

func (engine *Engine) AddSql(key string, sql string) {
	engine.SqlMap.addSql(key, sql)
}

func (sqlMap *SqlMap) addSql(key string, sql string) {
	sqlMap.checkNilAndInit()
	sqlMap.Sql[key] = sql
}

func (engine *Engine) UpdateSql(key string, sql string) {
	engine.SqlMap.updateSql(key, sql)
}

func (sqlMap *SqlMap) updateSql(key string, sql string) {
	sqlMap.checkNilAndInit()
	sqlMap.Sql[key] = sql
}

func (engine *Engine) RemoveSql(key string) {
	engine.SqlMap.removeSql(key)
}

func (sqlMap *SqlMap) removeSql(key string) {
	sqlMap.checkNilAndInit()
	delete(sqlMap.Sql, key)
}

func (engine *Engine) BatchAddSql(sqlStrMap map[string]string) {
	engine.SqlMap.batchAddSql(sqlStrMap)
}

func (sqlMap *SqlMap) batchAddSql(sqlStrMap map[string]string) {
	sqlMap.checkNilAndInit()
	for k, v := range sqlStrMap {
		sqlMap.Sql[k] = v
	}
}

func (engine *Engine) BatchUpdateSql(sqlStrMap map[string]string) {
	engine.SqlMap.batchUpdateSql(sqlStrMap)
}

func (sqlMap *SqlMap) batchUpdateSql(sqlStrMap map[string]string) {
	sqlMap.checkNilAndInit()
	for k, v := range sqlStrMap {
		sqlMap.Sql[k] = v
	}
}

func (engine *Engine) BatchRemoveSql(key []string) {
	engine.SqlMap.batchRemoveSql(key)
}

func (sqlMap *SqlMap) batchRemoveSql(key []string) {
	sqlMap.checkNilAndInit()
	for _, v := range key {
		delete(sqlMap.Sql, v)
	}
}

func (engine *Engine) GetSql(key string) string {
	return engine.SqlMap.getSql(key)
}

func (sqlMap *SqlMap) getSql(key string) string {
	return sqlMap.Sql[key]
}

func (engine *Engine) GetSqlMap(keys ...interface{}) map[string]string {
	return engine.SqlMap.getSqlMap(keys...)
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
