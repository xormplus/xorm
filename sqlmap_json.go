package xorm

type JsonSqlMap struct {
	sqlMapRootDir string
	extension     string
}

func Json(directory, extension string) *JsonSqlMap {
	return &JsonSqlMap{
		sqlMapRootDir: directory,
		extension:     extension,
	}
}

func (sqlMap *JsonSqlMap) RootDir() string {
	return sqlMap.sqlMapRootDir
}

func (sqlMap *JsonSqlMap) Extension() string {
	return sqlMap.extension
}
