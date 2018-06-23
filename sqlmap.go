package xorm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
)

type SqlMap struct {
	SqlMapRootDir string
	Sql           map[string]string
	Extension     map[string]string
	Capacity      uint
	Cipher        Cipher
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

func (sqlMap *SqlMap) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if strings.HasSuffix(path, sqlMap.Extension["xml"]) || strings.HasSuffix(path, sqlMap.Extension["json"]) || strings.HasSuffix(path, sqlMap.Extension["xsql"]) {
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

	if strings.HasSuffix(filepath, sqlMap.Extension["xml"]) {
		var result XmlSql
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

	if strings.HasSuffix(filepath, sqlMap.Extension["xsql"]) {
		scanner := &Scanner{}
		result := scanner.Run(bufio.NewScanner(bytes.NewReader(content)))
		for k := range result {
			sqlMap.Sql[k] = result[k]
		}

		return nil
	}
	return nil

}

func (sqlMap *SqlMap) addSql(key string, sql string) {
	sqlMap.checkNilAndInit()
	sqlMap.Sql[key] = sql
}

func (sqlMap *SqlMap) updateSql(key string, sql string) {
	sqlMap.checkNilAndInit()
	sqlMap.Sql[key] = sql
}

func (sqlMap *SqlMap) removeSql(key string) {
	sqlMap.checkNilAndInit()
	delete(sqlMap.Sql, key)
}

func (sqlMap *SqlMap) batchAddSql(sqlStrMap map[string]string) {
	sqlMap.checkNilAndInit()
	for k, v := range sqlStrMap {
		sqlMap.Sql[k] = v
	}
}

func (sqlMap *SqlMap) batchUpdateSql(sqlStrMap map[string]string) {
	sqlMap.checkNilAndInit()
	for k, v := range sqlStrMap {
		sqlMap.Sql[k] = v
	}
}

func (sqlMap *SqlMap) batchRemoveSql(key []string) {
	sqlMap.checkNilAndInit()
	for _, v := range key {
		delete(sqlMap.Sql, v)
	}
}

func (sqlMap *SqlMap) getSql(key string) string {
	return sqlMap.Sql[key]
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
