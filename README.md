# xorm

###优化xorm的查询API，并提供类似ibatis的配置文件及动态SQL功能

```go
var err error
db, err = xorm.NewPostgreSQL("postgres://postgres:root@localhost:5432/testdb?sslmode=disable")

if err != nil {
	t.Fatal(err)
}

err = db.InitSqlMap() //初始化SqlMap配置，可选功能
if err != nil {
	t.Fatal(err)
}
err = db.InitSqlTemplate() //初始化动态SQL模板配置，可选功能
if err != nil {
	t.Fatal(err)
}
```
####db.InitSqlMap()过程
* err = db.InitSqlMap()读取程序所在目下的sql/xormcfg.ini配置文件中的sqlMapRootDir配置项
* 遍历sqlMapRootDir所配置的目录及其子目录下的所有xml配置文件，<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/studygolang.xml">配置文件样例 </a>
* 解析xml配置文件

####db.InitSqlTemplate()过程
* err = db.InitSqlTemplate()读取程序所在目下的sql/xormcfg.ini配置文件中的SqlTemplateRootDir配置项
* 遍历SqlTemplateRootDir所配置的目录及其子目录下的所有stpl模板文件，<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/select.example.stpl">模板文件样例</a>
* 解析stpl模板文件


###支持类似这样的链式读取数据操作
```go
sql：="select id,title,createdatetime,content from article where id = ?"
rows, err := db.Sql(sql, 2).Query().Json() //返回查询数据的json字符串

id := db.Sql(sql, 2).Query().Result[0]["id"] //返回查询数据的第一条数据的id列的值
title := db.Sql(sql, 2).Query().Result[0]["title"]
createdatetime := db.Sql(sql, 2).Query().Result[0]["createdatetime"]
content := db.Sql(sql, 2).Query().Result[0]["content"]
```

###支持SqlMap配置，<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/studygolang.xml">配置文件样例 </a>
```xml
<sqlMap>
	<sql id="selectAllArticle">
		select id,title,createdatetime,content 
		from article where id in (?1,?2)
	</sql>
	<sql id="selectStudentById1">
		select * from article where id=?id
	</sql>
</sqlMap>
```

```go
paramMap := map[string]interface{}{"1": 2, "2": 5} //支持参数使用map存放
rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryByParamMap().Xml() //返回查询结果的xml字符串
```
###提供动态SQL支持，使用pongo2模板引擎
例如配置文件名：select.example.stpl</br>
<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/select.example.stpl">配置样例</a>内容如下：
```java
select id,userid,title,createdatetime,content 
from article where  
{% if count>1%}
id=?id
{% else%}
userid=?userid
{% endif %}
```
```go
paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}
rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryByParamMap().Json()
```
