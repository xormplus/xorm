# xorm
xorm是一个简单而强大的Go语言ORM库. 通过它可以使数据库操作非常简便。

## 说明

* 本库是基于原版 <b>xorm</b>：[https://github.com/go-xorm/xorm](https://github.com/go-xorm/xorm) 的定制增强版本，由于本定制版有第三方库依赖（原版xorm无任何第三方库依赖），所以只好单独开了本Github库。
* 本库的相关定制功能是为了解决更简单的进行复杂SQL调用和一些特殊业务需求场景而开发的。
* 本定制版ORM相关核心功能和原版保持一致，会跟随原版xorm更新。
* 定制功能采用针对原版弱侵入性代码实现。

## 特性

* 支持Struct和数据库表之间的灵活映射，并支持自动同步

* 事务支持

* 同时支持原始SQL语句和ORM操作的混合执行

* 支持类ibatis方式配置SQL语句（支持xml配置文件和pogon2模板2种方式）

* 支持动态SQL功能

* 使用连写来简化调用

* 支持使用Id, In, Where, Limit, Join, Having, Table, Sql, Cols等函数和结构体等方式作为条件

* 支持级联加载Struct

* 支持缓存

* 支持根据数据库自动生成xorm的结构体

* 支持记录版本（即乐观锁）

## 驱动支持

目前支持的Go数据库驱动和对应的数据库如下：

* Mysql: [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)

* MyMysql: [github.com/ziutek/mymysql/godrv](https://github.com/ziutek/mymysql/godrv)

* Postgres: [github.com/lib/pq](https://github.com/lib/pq)

* Tidb: [github.com/pingcap/tidb](https://github.com/pingcap/tidb)

* SQLite: [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

* MsSql: [github.com/denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb)

* MsSql: [github.com/lunny/godbc](https://github.com/lunny/godbc)

* Oracle: [github.com/mattn/go-oci8](https://github.com/mattn/go-oci8) (试验性支持)

* ql: [github.com/cznic/ql](https://github.com/cznic/ql) (试验性支持)

## 安装

推荐使用 [gopm](https://github.com/gpmgo/gopm) 进行安装：

	gopm get github.com/xormplus/xorm

或者您也可以使用go工具进行安装：

	go get -u github.com/xormplus/xorm

# 快速开始

* 第一步创建引擎，driverName, dataSourceName和database/sql接口相同

```go
//xorm原版标准方式创建引擎
engine, err := xorm.NewEngine(driverName, dataSourceName)

//您也可以针对特定数据库及数据库驱动使用类似下面的快捷方式创建引擎
engine, err = xorm.NewPostgreSQL(dataSourceName)
engine, err = xorm.NewSqlite3(dataSourceName)
```

* 创建引擎，且需要使用类ibatis方式配置SQL语句，参考如下方式

```go
var err error
engine, err = xorm.NewPostgreSQL("postgres://postgres:root@localhost:5432/testdb?sslmode=disable")

if err != nil {
	t.Fatal(err)
}

engine.SqlMap.SqlMapRootDir="./sql/oracle" //SqlMap配置文件总根目录，可代码指定，也可在配置文件中配置，代码指定优先级高于配置
engine.SqlTemplate.SqlTemplateRootDir="./sql/oracle" //SqlTemplate模板配置文件总根目录，可代码指定，也可在配置文件中配置，代码指定优先级高于配置

err = engine.InitSqlMap() //初始化SqlMap配置，可选功能，如应用中无需使用SqlMap，可无需初始化
if err != nil {
	t.Fatal(err)
}
err = engine.InitSqlTemplate() //初始化动态SQL模板配置，可选功能，如应用中无需使用SqlTemplate，可无需初始化
if err != nil {
	t.Fatal(err)
}

err = engine.StartMonitorFs() //开启SqlMap配置文件和SqlTemplate配置文件更新监控功能，将配置文件更新内容实时更新到内存，如无需要可以不调用该方法
if err != nil {
	t.Fatal(err)
}

```

* <b>db.InitSqlMap()过程</b>
    * 如指定db.SqlMap.SqlMapRootDir，则err = db.InitSqlMap()按指定目录遍历SqlMapRootDir所配置的目录及其子目录下的所有xml配置文件（<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/studygolang.xml">配置文件样例 </a>）
    * 如未指定db.SqlMap.SqlMapRootDir，err = db.InitSqlMap()则读取程序所在目下的sql/xormcfg.ini配置文件(<a href="https://github.com/xormplus/xorm/blob/master/test/sql/xormcfg.ini">样例</a>)中的SqlMapRootDir配置项，遍历SqlMapRootDir所配置的目录及其子目录下的所有xml配置文件（<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/studygolang.xml">配置文件样例 </a>）
    * 解析所有配置SqlMap的xml配置文件
    * 配置文件中sql标签的id属性值作为SqlMap的key，如有重名id，则后加载的覆盖之前加载的配置sql条目

* <b>db.InitSqlTemplate()过程</b>
    * 如指定db.SqlTemplate.SqlTemplateRootDir，err = db.InitSqlTemplate()按指定目录遍历SqlTemplateRootDir所配置的目录及其子目录下的所有stpl模板文件（<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/select.example.stpl">模板文件样例</a>）
    * 如指未定db.SqlTemplate.SqlTemplateRootDir，err = db.InitSqlTemplate()则读取程序所在目下的sql/xormcfg.ini配置文件(<a href="https://github.com/xormplus/xorm/blob/master/test/sql/xormcfg.ini">样例</a>)中的SqlTemplateRootDir配置项，遍历SqlTemplateRootDir所配置的目录及其子目录下的所有stpl模板文件（<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/select.example.stpl">模板文件样例</a>）
    * 解析stpl模板文件
    * stpl模板文件名作为SqlTemplate存储的key（不包含目录路径），如有不同路径下出现同名文件，则后加载的覆盖之前加载的配置模板内容

* 支持最原始的SQL语句查询

```go
//第1种方式，返回的结果类型为 []map[string][]byte
sql_1 := "select * from user"
results, err := engine.Query(sql_1)

//第2种方式，返回的结果类型为 []map[string]interface{}
sql_2_1 := "select * from user"
results := db.Sql(sql_2_1).Query().Result

sql_2_2 := "select * from user where id = ? and age = ?"
results := engine.Sql(sql_2_2, 7, 17).Query().Result


//第3种方式，执行SqlMap配置文件中的Sql语句，返回的结果类型为 []map[string]interface{}
sql_id_3_1 := "sql_3_1" //配置文件中sql标签的id属性,SqlMap的key
results := db.SqlMapClient(sql_3_1).Query().Result

sql_id_3_2 := "sql_3_2"
results := db.SqlMapClient(sql_id_3_2, 7, 17).Query().Result

sql_id_3_3 := "sql_3_3"
paramMap_3_3 := map[string]interface{}{"id": 7, "name": "xormplus"}
results1 := engine.SqlMapClient(sql_id_3_3, &paramMap_3_3).Query().Result

//第4种方式，执行SqlTemplate配置文件中的Sql语句，返回的结果类型为 []map[string]interface{}
sql_key_4_1 := "select.example.stpl" //配置文件名,SqlTemplate的key

//执行的 sql：select * from user where id=7
//如部分参数未使用，请记得使用对应类型0值，如此处name参数值为空字符串，模板使用指南请详见pogon2
paramMap_4_1 := map[string]interface{}{"count": 1, "id": 7, "name": ""}
results := db.SqlTemplateClient(sql_key_4_1, paramMap_4_1).Query().Result

//执行的 sql：select * from user where name='xormplus'
//如部分参数未使用，请记得使用对应类型0值，如此处id参数值为0，模板使用指南请详见pogon2
paramMap_4_2 := map[string]interface{}{"id": 0, "count": 2, "name": "xormplus"}
results := db.SqlTemplateClient(sql_key_4_1, paramMap_4_2).Query().Result
```

* 第3种方式所使用的SqlMap配置文件内容如下

```xml
<sqlMap>
	<sql id="sql_3_1">
		select * from user
	</sql>
	<sql id="sql_3_2">
		select * from user where id=? and age=?
	</sql>
    <sql id="sql_3_3">
		select * from user where id=?id and name=?name
	</sql>
</sqlMap>
```

* 第4种方式所使用的SqlTemplate配置文件内容如下，文件名：select.example.stpl，路径为engine.SqlMap.SqlMapRootDir配置目录下的任意子目录中。使用模板方式配置Sql较为灵活，可以使用pogon2引擎的相关功能灵活组织Sql语句以及动态SQL拼装。

```java
select * from user
where
{% if count>1%}
id=?id
{% else%}
name=?name
{% endif %}
```

* 执行一个SQL语句

```go
//第1种方式
affected, err := engine.Exec("update user set age = ? where name = ?", age, name)

//第2种方式
sql_2 := "INSERT INTO config(key,value) VALUES (?, ?)"
affected, err := engine.Sql(sql_4, "OSCHINA", "OSCHINA").Execute()

//第3种方式
sql_i_1 := "sql_i_1" //SqlMap中key为 "sql_i_1" 配置的Sql语句为：INSERT INTO config(key,value) VALUES (?, ?)
affected, err := engine.SqlMapClient(sql_i_1, "config_1", "1").Execute()

sql_i_2 := "sql_i_2" //SqlMap中key为 "sql_i_2" 配置的Sql语句为：INSERT INTO config(key,value) VALUES (?key, ?value)
paramMap_i := map[string]interface{}{"key": "config_2", "value": "2"}
affected, err := engine.SqlMapClient(sql_i_2, &paramMap_i).Execute()

//第4种方式
sql_i_3 := "insert.example.stpl"
paramMap_i_t := map[string]interface{}{"key": "config_3", "value": "3"}
affected, err := engine.SqlTemplateClient(sql_i_3, paramMap_i_t).Execute()
```

* 支持链式读取数据操作查询返回json或xml字符串

```go
//第1种方式
users := make([]User, 0)
results,err := db.Where("id=?", 6).Find(&users).Xml() //返回查询结果的xml字符串
results,err := db.Where("id=?", 6).Find(&users).Json() //返回查询结果的json字符串

//第2种方式
sql := "select * from user where id = ?"
results, err := db.Sql(sql, 2).Query().Json() //返回查询结果的json字符串
results, err := db.Sql(sql, 2).QueryWithDateFormat("20060102").Json() //返回查询结果的json字符串，并支持格式化日期
results, err := db.Sql(sql, 2).QueryWithDateFormat("20060102").Xml() //返回查询结果的xml字符串，并支持格式化日期

sql := "select * from user where id = ?id and userid=?userid"
paramMap := map[string]interface{}{"id": 6, "userid": 1} //支持参数使用map存放
results, err := db.Sql(sql, &paramMap).Query().XmlIndent("", "  ", "article") //返回查询结果格式化后的xml字符串

//第3种方式
sql_id_3_1 := "sql_3_1" //配置文件中sql标签的id属性,SqlMap的key
results, err := db.SqlMapClient(sql_id_3_1, 7, 17).Query().Json() //返回查询结果的json字符串

sql_id_3_2 := "sql_3_2" //配置文件中sql标签的id属性,SqlMap的key
paramMap := map[string]interface{}{"id": 6, "userid": 1} //支持参数使用map存放
results, err := db.SqlMapClient(sql_id_3_2, &paramMap).Query().Xml() //返回查询结果的xml字符串

//第4种方式
sql_key_4_1 := "select.example.stpl"
paramMap_4_1 := map[string]interface{}{"id": 6, "userid": 1}
results := db.SqlTemplateClient(sql_key_4_1, paramMap_4_1).Query().Json()
```

* 支持链式读取数据操作查询返回某条记录的某个字段的值

```go
//第1种方式
id := db.Sql(sql, 2).Query().Result[0]["id"] //返回查询结果的第一条数据的id列的值

//第2种方式
id := db.SqlMapClient(key, 2).Query().Result[0]["id"] //返回查询结果的第一条数据的id列的值
id := db.SqlMapClient(key, &paramMap).Query().Result[0]["id"] //返回查询结果的第一条数据的id列的值

//第3种方式
id := db.SqlTemplateClient(key, paramMap).Query().Result[0]["id"] //返回查询结果的第一条数据的id列的值
```

* SqlMap及SqlTemplate相关功能API

```go
//SqlMap配置文件总根目录，可代码指定，也可在配置文件中配置，代码指定优先级高于配置
engine.SqlMap.SqlMapRootDir="./sql/oracle"
//SqlTemplate模板配置文件总根目录，可代码指定，也可在配置文件中配置，代码指定优先级高于配置
engine.SqlTemplate.SqlTemplateRootDir="./sql/oracle"

err := engine.InitSqlMap()//初始化加载SqlMap配置文件，默认初始化后缀为".xml"
err := engine.InitSqlTemplate()//初始化加载SqlTemplate配置文件，默认初始化后缀为".stpl"

//SqlMap配置文件和SqlTemplate配置文件后缀不要相同
option := xorm.SqlMapOptions{Extension: ".xx"} //指定SqlMap配置文件后缀为".xx"，但配置内容必须为样例的xml格式
err := engine.InitSqlMap(option) //按指定SqlMap配置文件后缀为".xx"初始化

option := xorm.SqlTemplateOptions{Extension: ".yy"} //指定SqlTemplate配置文件后缀为".yy"
err = engine.InitSqlTemplate(option) //按指定SqlMap配置文件后缀为".xx"初始化

//开启SqlMap配置文件和SqlTemplate配置文件更新监控功能，将配置文件更新内容实时更新到内存，如无需要可以不调用该方法
//该监控模式下，如删除配置文件，内存中不会删除相关配置
engine.StartFSWatcher()
//停止SqlMap配置文件和SqlTemplate配置文件更新监控功能
engine.StopFSWatcher()

engine.LoadSqlMap(filepath) //加载指定文件的SqlMap配置
engine.ReloadSqlMap(filepath) //重新加载指定文件的SqlMap配置

engine.AddSql(key, sql) //新增一条SqlMap配置
engine.UpdateSql(key, sql) //更新一条SqlMap配置
engine.RemoveSql(key) //删除一条SqlMap配置

engine.LoadSqlTemplate(filepath) //加载指定文件的SqlTemplate模板
engine.ReloadSqlTemplate(filepath) //重新加载指定文件的SqlTemplate模板

engine.AddSqlTemplate(key, sql) //新增一条SqlTemplate模板
engine.UpdateSqlTemplate(key, sql) //更新一条SqlTemplate模板
engine.RemoveSqlTemplate(key) //删除一条SqlTemplate模板
```
* 插入一条或者多条记录

```Go
affected, err := engine.Insert(&user)
// INSERT INTO struct () values ()
affected, err := engine.Insert(&user1, &user2)
// INSERT INTO struct1 () values ()
// INSERT INTO struct2 () values ()
affected, err := engine.Insert(&users)
// INSERT INTO struct () values (),(),()
affected, err := engine.Insert(&user1, &users)
// INSERT INTO struct1 () values ()
// INSERT INTO struct2 () values (),(),()
```

* 查询单条记录

```Go
has, err := engine.Get(&user)
// SELECT * FROM user LIMIT 1
has, err := engine.Where("name = ?", name).Desc("id").Get(&user)
// SELECT * FROM user WHERE name = ? ORDER BY id DESC LIMIT 1
```

* 查询多条记录，当然可以使用Join和extends来组合使用

```Go
var users []User
err := engine.Where("name = ?", name).And("age > 10").Limit(10, 0).Find(&users)
// SELECT * FROM user WHERE name = ? AND age > 10 limit 0 offset 10

type Detail struct {
    Id int64
    UserId int64 `xorm:"index"`
}

type UserDetail struct {
    User `xorm:"extends"`
    Detail `xorm:"extends"`
}

var users []UserDetail
err := engine.Table("user").Select("user.*, detail.*")
    Join("INNER", "detail", "detail.user_id = user.id").
    Where("user.name = ?", name).Limit(10, 0).
    Find(&users)
// SELECT user.*, detail.* FROM user INNER JOIN detail WHERE user.name = ? limit 0 offset 10
```

* 根据条件遍历数据库，可以有两种方式: Iterate and Rows

```Go
err := engine.Iterate(&User{Name:name}, func(idx int, bean interface{}) error {
    user := bean.(*User)
    return nil
})
// SELECT * FROM user

rows, err := engine.Rows(&User{Name:name})
// SELECT * FROM user
defer rows.Close()
bean := new(Struct)
for rows.Next() {
    err = rows.Scan(bean)
}
```

* 更新数据，除非使用Cols,AllCols函数指明，默认只更新非空和非0的字段

```Go
affected, err := engine.Id(1).Update(&user)
// UPDATE user SET ... Where id = ?

affected, err := engine.Update(&user, &User{Name:name})
// UPDATE user SET ... Where name = ?

var ids = []int64{1, 2, 3}
affected, err := engine.In(ids).Update(&user)
// UPDATE user SET ... Where id IN (?, ?, ?)

// force update indicated columns by Cols
affected, err := engine.Id(1).Cols("age").Update(&User{Name:name, Age: 12})
// UPDATE user SET age = ?, updated=? Where id = ?

// force NOT update indicated columns by Omit
affected, err := engine.Id(1).Omit("name").Update(&User{Name:name, Age: 12})
// UPDATE user SET age = ?, updated=? Where id = ?

affected, err := engine.Id(1).AllCols().Update(&user)
// UPDATE user SET name=?,age=?,salt=?,passwd=?,updated=? Where id = ?
```

* 删除记录，需要注意，删除必须至少有一个条件，否则会报错。要清空数据库可以用EmptyTable

```Go
affected, err := engine.Where(...).Delete(&user)
// DELETE FROM user Where ...
```

* 获取记录条数

```Go
counts, err := engine.Count(&user)
// SELECT count(*) AS total FROM user
```

## 部分测试用例


<a href="https://github.com/xormplus/xorm/blob/master/test/xorm_test.go">测试用例</a>，<a href="https://github.com/xormplus/xorm/blob/master/test/测试结果.txt">测试结果</a>

## 文档

* 原版核心ORM功能请详见原版[《xorm操作指南》](http://xorm.io/docs)

## 讨论
请加入QQ群：280360085 进行讨论。API设计相关建议可联系本人QQ：50892683
