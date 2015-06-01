package xorm

//	"fmt"

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
