package xorm

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
//	"fmt"

	"github.com/Unknwon/goconfig"
	"gopkg.in/flosch/pongo2.v3"
)

const (
	MSSQL_DRIVER      string = "mssql"
	MSSQL_ODBC_DRIVER string = "odbc"
	MYSQL_DRIVER      string = "mysql"
	MYMYSQL_DRIVER    string = "mymysql"
	POSTGRESQL_DRIVER string = "postgres"
	OCI8_DRIVER       string = "oci8"
	GORACLE_DRIVER    string = "goracle"
	SQLITE3_DRIVER    string = "sqlite3"
)

type SqlMap struct {
	Sql         map[string]string
	SqlTemplate map[string]*pongo2.Template
	Cfg goconfig.ConfigFile
}

type Result struct {
	Sql []Sql `xml:"sql"`
}

type Sql struct {
	Value string `xml:",chardata"`
	Id    string `xml:"id,attr"`
}

func (sqlMap *SqlMap) Init() error {
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

	sqlMap.Sql = make(map[string]string)
	sqlMap.SqlTemplate = make(map[string]*pongo2.Template)
	err=filepath.Walk(sqlMapRootDir, sqlMap.walkFunc)
	if err!=nil{
		return err
	}

	return nil
}

func (sqlMap *SqlMap) walkFunc(path string, info os.FileInfo, err error) error {
	if err!=nil{
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

	if strings.HasSuffix(path, ".stpl") {
		err = sqlMap.paresSqlTemplate(info.Name(), path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sqlMap *SqlMap) paresSqlTemplate(filename string, filepath string) error {
	template, err := pongo2.FromFile(filepath)
	if err != nil {
		return err
	}

	sqlMap.SqlTemplate[filename] = template
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

func NewOracle(driverName string, dataSourceName string) (*Engine, error) {
	return NewEngine(driverName, dataSourceName)
}

func NewMSSQL(driverName string, dataSourceName string) (*Engine, error) {
	return NewEngine(driverName, dataSourceName)
}

func NewMySQL(driverName string, dataSourceName string) (*Engine, error) {
	return NewEngine(driverName, dataSourceName)
}

func NewPostgreSQL(dataSourceName string) (*Engine, error) {
	return NewEngine(POSTGRESQL_DRIVER, dataSourceName)
}

func NewSqlite3(dataSourceName string) (*Engine, error) {
	return NewEngine(SQLITE3_DRIVER, dataSourceName)
}

func NewDB(driverName string, dataSourceName string) (*Engine, error) {
	return NewEngine(driverName, dataSourceName)
}
