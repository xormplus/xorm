# xorm

###优化xorm的查询API，并提供类似ibatis的配置文件及动态SQL功能，支持AcitveRecord操作

<pre>go get -u github.com/xormplus/xorm</pre>

<a href="https://github.com/xormplus/xorm/blob/master/test/xorm_test.go">测试用例</a>，<a href="https://github.com/xormplus/xorm/blob/master/test/测试结果.txt">测试结果</a>

```go
var err error
db, err = xorm.NewPostgreSQL("postgres://postgres:root@localhost:5432/testdb?sslmode=disable")

db.SqlMap.SqlMapRootDir="./sql/oracle" //SqlMap配置文件存根目录，可代码指定，可在配置文件中配置，代码指定优先级高于配置
db.SqlTemplate.SqlTemplateRootDir="./sql/oracle" //SqlTemplate配置文件存根目录，可代码指定，可在配置文件中配置，代码指定优先级高于配置
if err != nil {
	t.Fatal(err)
}

err = db.InitSqlMap() //初始化SqlMap配置，可选功能，如应用中无需使用SqlMap，可无需初始化
if err != nil {
	t.Fatal(err)
}
err = db.InitSqlTemplate() //初始化动态SQL模板配置，可选功能，如应用中无需使用SqlTemplate，可无需初始化
if err != nil {
	t.Fatal(err)
}
```
####db.InitSqlMap()过程
* 如指定db.SqlMap.SqlMapRootDir，则err = db.InitSqlMap()按指定目录遍历SqlMapRootDir所配置的目录及其子目录下的所有xml配置文件（<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/studygolang.xml">配置文件样例 </a>）
* 如未指定db.SqlMap.SqlMapRootDir，err = db.InitSqlMap()则读取程序所在目下的sql/xormcfg.ini配置文件(<a href="https://github.com/xormplus/xorm/blob/master/test/sql/xormcfg.ini">样例</a>)中的SqlMapRootDir配置项，遍历SqlMapRootDir所配置的目录及其子目录下的所有xml配置文件（<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/studygolang.xml">配置文件样例 </a>）
* 解析所有配置SqlMap的xml配置文件

####db.InitSqlTemplate()过程
* 如指定db.SqlTemplate.SqlTemplateRootDir，err = db.InitSqlTemplate()按指定目录遍历SqlTemplateRootDir所配置的目录及其子目录下的所有stpl模板文件（<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/select.example.stpl">模板文件样例</a>）
* 如指未定db.SqlTemplate.SqlTemplateRootDir，err = db.InitSqlTemplate()则读取程序所在目下的sql/xormcfg.ini配置文件(<a href="https://github.com/xormplus/xorm/blob/master/test/sql/xormcfg.ini">样例</a>)中的SqlTemplateRootDir配置项，遍历SqlTemplateRootDir所配置的目录及其子目录下的所有stpl模板文件（<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/select.example.stpl">模板文件样例</a>）
* 解析stpl模板文件


###支持类似这样的链式读取数据操作
```go
sql := "select id,title,createdatetime,content from article where id = ?"
rows, err := db.Sql(sql, 2).Query().Json() //返回查询结果的json字符串
rows, err := db.Sql("sql", 2).QueryWithDateFormat("20060102").Json() //返回查询结果的json字符串，并支持格式化日期
rows, err := db.Sql("sql", 2).QueryWithDateFormat("20060102").Xml() //返回查询结果的xml字符串，并支持格式化日期

id := db.Sql(sql, 2).Query().Result[0]["id"] //返回查询结果的第一条数据的id列的值
title := db.Sql(sql, 2).Query().Result[0]["title"]
createdatetime := db.Sql(sql, 2).Query().Result[0]["createdatetime"]
content := db.Sql(sql, 2).Query().Result[0]["content"]

articles := make([]Article, 0)
xml,err := db.Where("id=?", 6).Find(&articles).Xml() //返回查询结果的xml字符串
json,err := db.Where("id=?", 6).Find(&articles).Json() //返回查询结果的json字符串

sql := "select id,title,createdatetime,content from article where id = ?id and userid=?userid"
paramMap := map[string]interface{}{"id": 6, "userid": 1} //支持参数使用map存放
rows, err := db.Sql(sql, &paramMap).QueryByParamMap().XmlIndent("", "  ", "article")
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
rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryByParamMap().Json() //返回查询结果的json字符串
rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryByParamMapWithDateFormat("2006/01/02").XmlIndent("", "  ", "article") //返回查询结果格式化的xml字符串，并支持格式化日期
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
rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryByParamMapWithDateFormat("2006/01/02").Json()
rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryByParamMapWithDateFormat("2006/01/02").XmlIndent("", "  ", "article")
```

###讨论
请加入QQ群：280360085 进行讨论。API设计相关建议可联系本人QQ：50892683
