# xorm

优化xorm的查询API，并提供类似ibatis的配置文件及动态SQL功能

支持类似这样的链式操作
```go
sql：="select id,title,createdatetime,content from article where id = ?"
rows, err := db.Sql(sql, 2).FindAll().Json()

id := db.Sql(sql, 2).FindAll().Result[0]["id"]
title := db.Sql(sql, 2).FindAll().Result[0]["title"]
createdatetime := db.Sql(sql, 2).FindAll().Result[0]["createdatetime"]
content := db.Sql(sql, 2).FindAll().Result[0]["content"]
```

也支持SqlMap配置，<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/studygolang.xml">配置文件样例 </a>
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
paramMap := map[string]interface{}{"1": 2, "2": 5}
rows, err := db.SqlMapClient("selectAllArticle", &amp;paramMap).FindAllByParamMap().Xml()
```
同时提供动态SQL支持，使用pongo2模板引擎</br></br>
例如配置文件名：select.example.stpl</br>
配置<a href="https://github.com/xormplus/xorm/blob/master/test/sql/oracle/select.example.stpl">样例</a>内容如下：
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
rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).FindAllByParamMap().Json()
```
