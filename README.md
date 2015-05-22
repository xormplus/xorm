# xorm

优化xorm的查询API，并提供类似ibatis的配置文件及动态SQL功能

支持类似这样的链式操作
<pre>
sql：="select id,title,createdatetime,content from article where id = ?"</br>
rows, err := db.Sql(sql, 2).FindAll().Json()
</pre>
或者
<pre>
id := db.Sql(sql, 2).FindAll().Result[0]["id"]</br>
title := db.Sql(sql, 2).FindAll().Result[0]["title"]</br>
createdatetime := db.Sql(sql, 2).FindAll().Result[0]["createdatetime"]</br>
content := db.Sql(sql, 2).FindAll().Result[0]["content"]</br>
</pre>

也支持SqlMa配置，配置文件样例
<pre>
<sqlMap>
	<sql id="selectAllArticle">
		select id,title,createdatetime,content 
		from article where id in (?1,?2)
	</sql>
	<sql id="selectStudentById1">
		select * from article where id=?id
	</sql>
</sqlMap>
</pre>
<pre>
paramMap := map[string]interface{}{"1": 2, "2": 5}</br>
rows, err := db.SqlMapClient("selectAllArticle", &amp;paramMap).FindAllByParamMap().Xml()
</pre>
同时提供动态SQL支持，使用pongo2模板引擎
例如配置文件名：select.example.stpl
配置内容如下：
<pre>
select id,userid,title,createdatetime,content 
from article where  
{% if count>1%}
id=?id
{% else%}
userid=?userid
{% endif %}
</pre>
<pre>
paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}</br>
rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).FindAllByParamMap().Json()
<pre>
