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
	Sql map[string]string
}

type Result struct {
	Sql []Sql `xml:"sql"`
}

type Sql struct {
	Value string `xml:",chardata"`
	Id    string `xml:"id,attr"`
}

func (engine *Engine) InitSqlMap() error {
	var err error
	cfg, err := goconfig.LoadConfigFile("./sql/xormcfg.ini")
	if err != nil {
		return err
	}
	var sqlMapRootDir string
	sqlMapRootDir, err = cfg.GetValue("", "sqlMapRootDir")
	if err != nil {
		return err
	}
	
	engine.SqlMap.Sql = make(map[string]string)
	err = filepath.Walk(sqlMapRootDir, engine.SqlMap.walkFunc)
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

	if strings.HasSuffix(path, ".xml") {
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
