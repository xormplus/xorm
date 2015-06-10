##创建Orm引擎

在xorm里面，可以同时存在多个Orm引擎，一个Orm引擎称为Engine，一个Engine一般只对应一个数据库。Engine通过调用xorm.NewEngine生成，如：
```go
import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/xormplus/xorm"
)

var engine *xorm.Engine

func main() {
    var err error
    engine, err = xorm.NewEngine("mysql", "root:123@/test?charset=utf8")
}
```
or
```go
import (
    _ "github.com/mattn/go-sqlite3"
    "github.com/xormplus/xorm"
)

var engine *xorm.Engine

func main() {
    var err error
    engine, err = xorm.NewEngine("sqlite3", "./test.db")
}
```
一般情况下如果只操作一个数据库，只需要创建一个engine即可。engine是GoRutine安全的。<br>

为了方便代码编写你也可以依据数据库类型创建指定的engine，例如：
```go
import (
    _ "github.com/lib/pq"
	"github.com/xormplus/xorm"
)

var engine *xorm.Engine

func main() {
    var err error
    engine, err = xorm.NewPostgreSQL("postgres://postgres:root@localhost:5432/test?sslmode=disable")
}
```
or
```go
import (
    _ "github.com/go-sql-driver/mysql"
    "github.com/xormplus/xorm"
)

var engine *xorm.Engine

func main() {
    var err error
    engine, err = xorm.NewMySQL(xorm.MYSQL_DRIVER,"root:123@/test?charset=utf8")
}
```

创建完成engine之后，并没有立即连接数据库，此时可以通过engine.Ping()来进行数据库的连接测试是否可以连接到数据库。另外对于某些数据库有连接超时设置的，可以通过起一个定期Ping的Go程来保持连接鲜活。

对于有大量数据并且需要分区的应用，也可以根据规则来创建多个Engine，比如：
```go
var err error
for i:=0;i<5;i++ {
    engines[i], err = xorm.NewEngine("sqlite3", fmt.Sprintf("./test%d.db", i))
}
```
engine可以通过engine.Close来手动关闭，但是一般情况下可以不用关闭，在程序退出时会自动关闭。

NewEngine传入的参数和sql.Open传入的参数完全相同，因此，在使用某个驱动前，请查看此驱动中关于传入参数的说明文档。以下为各个驱动的连接符对应的文档链接：

- <a href="http://godoc.org/github.com/mattn/go-sqlite3#SQLiteDriver.Open">sqlite3</a>
- <a href="https://github.com/go-sql-driver/mysql#dsn-data-source-name">mysql dsn</a>
- <a href="http://godoc.org/github.com/ziutek/mymysql/godrv#Driver.Open">mymysql</a>
- <a href="http://godoc.org/github.com/lib/pq">postgres</a>

在engine创建完成后可以进行一些设置，如：


1.调试，警告以及错误等显示设置，默认如下均为false
- engine.ShowSQL = true，则会在控制台打印出生成的SQL语句；
- engine.ShowDebug = true，则会在控制台打印调试信息；
- engine.ShowError = true，则会在控制台打印错误信息；
- engine.ShowWarn = true，则会在控制台打印警告信息；

2.如果希望将信息不仅打印到控制台，而是保存为文件，那么可以通过类似如下的代码实现，NewSimpleLogger(w io.Writer)接收一个io.Writer接口来将数据写入到对应的设施中。
```go
f, err := os.Create("sql.log")
    if err != nil {
        println(err.Error())
        return
    }
defer f.Close()
engine.Logger = xorm.NewSimpleLogger(f)
```
3.engine内部支持连接池接口和对应的函数。

如果需要设置连接池的空闲数大小，可以使用engine.SetMaxIdleConns()来实现。
如果需要设置最大打开连接数，则可以使用engine.SetMaxOpenConns()来实现。