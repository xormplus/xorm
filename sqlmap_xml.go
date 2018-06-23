package xorm

type XmlSqlMap struct {
	sqlMapRootDir string
	extension     string
}

type XmlSql struct {
	Sql []Sql `xml:"sql"`
}

type Sql struct {
	Value string `xml:",chardata"`
	Id    string `xml:"id,attr"`
}

func Xml(directory, extension string) *XmlSqlMap {
	return &XmlSqlMap{
		sqlMapRootDir: directory,
		extension:     extension,
	}
}

func (sqlMap *XmlSqlMap) RootDir() string {
	return sqlMap.sqlMapRootDir
}

func (sqlMap *XmlSqlMap) Extension() string {
	return sqlMap.extension
}
