package xorm

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/xormplus/xorm"

	_ "github.com/lib/pq"
)

type Article struct {
	Id             int       `xorm:"not null pk autoincr unique INTEGER"`
	Content        string    `xorm:"not null TEXT"`
	Title          string    `xorm:"not null VARCHAR(255)"`
	Categorysubid  int       `xorm:"not null INTEGER"`
	Remark         string    `xorm:"not null VARCHAR(2555)"`
	Userid         int       `xorm:"not null INTEGER"`
	Viewcount      int       `xorm:"not null default 0 INTEGER"`
	Replycount     int       `xorm:"not null default 0 INTEGER"`
	Tags           string    `xorm:"not null VARCHAR(300)"`
	Createdatetime JSONTime  `xorm:"not null default 'now()' DATETIME"`
	Isdraft        int       `xorm:"SMALLINT"`
	Lastupdatetime time.Time `xorm:"not null default 'now()' DATETIME"`
}

type Category struct {
	Id             int       `xorm:"not null pk autoincr unique INTEGER"`
	Name           string    `xorm:"not null VARCHAR(200)"`
	Counts         int       `xorm:"not null default 0 INTEGER"`
	Orders         int       `xorm:"not null default 0 INTEGER"`
	Createtime     time.Time `xorm:"not null default 'now()' created DATETIME"`
	Pid            int       `xorm:"not null default 0 INTEGER"`
	Lastupdatetime time.Time `xorm:"not null default 'now()' updated  DATETIME"`
	Status         int       `xorm:"not null default 1 SMALLINT"`
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006/01/08 15:04:05"))
	return []byte(stamp), nil
}

var db *xorm.Engine

func Test_InitDB(t *testing.T) {
	var err error
	db, err = xorm.NewPostgreSQL("postgres://postgres:root@localhost:5432/mblog?sslmode=disable")

	if err != nil {
		t.Fatal(err)
	}

	err = db.RegisterSqlMap(xorm.Xml("./sql/oracle", ".xml"))
	if err != nil {
		t.Fatal(err)
	}

	err = db.RegisterSqlMap(xorm.Json("./sql/oracle", ".json"))
	if err != nil {
		t.Fatal(err)
	}

	err = db.RegisterSqlTemplate(xorm.Default("./sql/oracle", ".tpl"))
	if err != nil {
		t.Fatal(err)
	}

	err = db.StartFSWatcher()
	if err != nil {
		t.Fatal(err)
	}

	db.ShowSQL(true)

	log.Println(db)
	//	db.NewSession().SqlMapClient().Execute()
	log.Println(db.GetSqlMap("json_category-16-17"))
}

func Test_Get_Struct(t *testing.T) {
	var article Article
	has, err := db.Id(2).Get(&article)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_Get_Struct]->rows: not exist\n")
	}

	t.Log("[Test_Get_Struct]->rows:\n", article)
}

func Test_Sql_Get_Struct(t *testing.T) {
	var article Article
	has, err := db.Sql("select * from article where id=?", 2).Get(&article)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_Sql_Get_Struct]->rows: not exist\n")
	}

	t.Log("[Test_Sql_Get_Struct]->rows:\n", article)
}

func Test_Sql_Get_String(t *testing.T) {
	var title string
	has, err := db.Sql("select title from article where id=?", 2).Get(&title)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_Sql_Get_String]->rows: not exist\n")
	}

	t.Log("[Test_Sql_Get_String]->titles:\n", title)
}

func Test_Sql_Get_Int(t *testing.T) {
	var id int
	has, err := db.Sql("select id from article where id=?", 2).Get(&id)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_Sql_Get_Int]->rows: not exist\n")
	}

	t.Log("[Test_Sql_Get_Int]->id:\n", id)
}

func Test_Sql_Get_Map(t *testing.T) {
	var valuesMap = make(map[string]string)
	has, err := db.Sql("select * from article where id=?", 2).Get(&valuesMap)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_Sql_Get_Map]->rows: not exist\n")
	}

	t.Log("[Test_Sql_Get_Map]->id:\n", valuesMap)
}

func Test_Sql_Get_Map2(t *testing.T) {
	var valuesMap = make(map[string]interface{})
	has, err := db.Sql("select * from article where id=?", 2).Get(&valuesMap)
	if err != nil {
		t.Fatal(err)
	}

	if !has {
		t.Log("[Test_Sql_Get_Map2]->rows: not exist\n")
	}

	t.Log("[Test_Sql_Get_Map2]->id:\n", valuesMap)
	id, ok := valuesMap["id"].(int64)
	if ok {
		t.Log("[Test_Sql_Get_Map2]->id:\n", id)
	}

}

func Test_Sql_Get_Slice(t *testing.T) {
	var valuesSlice = make([]interface{}, 2)
	has, err := db.Sql("select id,title from article where id=?", 2).Get(&valuesSlice)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_Sql_Get_Slice]->rows: not exist\n")
	}

	t.Log("[Test_Sql_Get_Slice]->id:\n", valuesSlice)
}

func Test_GetFirst_Json(t *testing.T) {

	var article Article
	has, rows, err := db.Id(2).
		GetFirst(&article).
		Json()
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_GetFirst_Json]->rows: not exist\n")
	}
	t.Log("[Test_GetFirst_Json]->rows:\n" + rows)
}

func Test_GetFirst_Xml(t *testing.T) {
	var article Article
	has, rows, err := db.Where("userid =?", 2).GetFirst(&article).Xml()
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_GetFirst_Xml]->rows: not exist\n")
	}
	t.Log("[Test_GetFirst_Xml]->rows:\n" + rows)
}

func Test_GetFirst_XmlIndent(t *testing.T) {
	var article Article
	has, rows, err := db.Where("userid =?", 2).GetFirst(&article).XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_GetFirst_XmlIndent]->rows: not exist\n")
	}
	t.Log("[Test_GetFirst_XmlIndent]->rows:\n" + rows)
}

func Test_Search(t *testing.T) {
	var article []Article
	result := db.Sql("select id,title,createdatetime,content from article where id = ?", 25).Search(&article)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	t.Log("[Test_Find]->article[0].Id:\n", article[0].Id)
	t.Log("[Test_Find]->article[0].Content:\n", article[0].Content)
	t.Log("[Test_Find]->article[0].Title:\n", article[0].Title)
	t.Log("[Test_Find]->article[0].Categorysubid:\n", article[0].Categorysubid)
	t.Log("[Test_Find]->article[0].Createdatetime:\n", article[0].Createdatetime)
	t.Log("[Test_Find]->article[0].Isdraft:\n", article[0].Isdraft)
	t.Log("[Test_Find]->article[0].Lastupdatetime:\n", article[0].Lastupdatetime)
	t.Log("[Test_Find]->article[0].Remark:\n", article[0].Remark)
	t.Log("[Test_Find]->article[0].Replycount:\n", article[0].Replycount)
	t.Log("[Test_Find]->article[0].Tags:\n", article[0].Tags)
	t.Log("[Test_Find]->article[0].Userid:\n", article[0].Userid)
	t.Log("[Test_Find]->article[0].Viewcount:\n", article[0].Viewcount)
	t.Log("[Test_Find]-> result.Result:\n", result.Result)

	resultJson, err := result.Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_Search]-> result.Json():\n", resultJson)
}

func Test_Session_Search2(t *testing.T) {
	session := db.NewSession()
	defer session.Close()
	var article []Article
	result := session.Sql("select id,title,createdatetime,content from article where id = ?", 25).Search(&article)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	resultJson, err := result.Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_Search]-> result.Json():\n", resultJson)
}

func Test_Session_Search(t *testing.T) {
	session := db.NewSession()
	defer session.Close()
	var article []Article
	result := session.Sql("select id,title,createdatetime,content from article where id = ?id", &map[string]interface{}{"id": 25}).Search(&article)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	resultJson, err := result.Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_Search]-> result.Json():\n", resultJson)
}

func Test_Session_SqlMapClient_Search(t *testing.T) {
	session := db.NewSession()
	defer session.Close()
	var article []Article
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	result := session.SqlMapClient("selectAllArticle", &paramMap).Search(&article)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMap_Xml]->rows:\n", result)
	t.Log("[Test_SqlMapClient_QueryByParamMap_Xml]->article:\n", article)

}

func Test_Query_Json(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 139).Query().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_Query_Json]->rows:\n" + rows)
}

func Test_Query_Result(t *testing.T) {
	rows := db.Sql("select id,title,createdatetime,content from article where id = ?", 139).Query()
	if rows.Error != nil {
		t.Fatal(rows.Error)
	}

	t.Log("[Test_Query_Result]->rows[0][\"id\"]:\n", rows.Results[0]["id"])
	t.Log("[Test_Query_Result]->reflect.TypeOf(rows.Result[0][\"id\"]):\n", reflect.TypeOf(rows.Results[0]["id"]))
	t.Log("[Test_Query_Result]->rows[0][\"title\"]:\n", rows.Results[0]["title"])
	t.Log("[Test_Query_Result]->reflect.TypeOf(rows.Result[0][\"title\"]):\n", reflect.TypeOf(rows.Results[0]["title"]))
	t.Log("[Test_Query_Result]->rows[0][\"createdatetime\"]:\n", rows.Results[0]["createdatetime"])
	t.Log("[Test_Query_Result]->reflect.TypeOf(rows.Result[0][\"createdatetime\"]):\n", reflect.TypeOf(rows.Results[0]["createdatetime"]))

}

func Test_Query_Xml(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 139).Query().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_Query_Xml]->rows:\n" + rows)
}

func Test_Query_XmlIndent(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 139).Query().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_Query_XmlIndent]->rows:\n" + rows)
}

func Test_QueryWithDateFormat_Json(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 139).QueryWithDateFormat("20060102").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_QueryWithDateFormat_Json]->rows:\n" + rows)
}

func Test_QueryWithDateFormat_Xml(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 33).QueryWithDateFormat("20060102").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_QueryWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_QueryWithDateFormat_XmlIndent(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id in (?,?)", 139, 33).QueryWithDateFormat("20060102").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_QueryWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_QueryByParamMap_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 32, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).Query().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_QueryByParamMap_Json]->rows:\n" + rows)
}

func Test_QueryByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 6, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).Query().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_QueryByParamMap_Xml]->rows:\n" + rows)
}

func Test_QueryByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 6, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).Query().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_QueryByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_QueryByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 6, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).QueryWithDateFormat("2006/01/02").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_QueryByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMap_Json(t *testing.T) {
	paramMap := map[string]interface{}{"1": 14, "2": 139}
	rows, err := db.SqlMapClient("json_selectAllArticle", &paramMap).Query().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMap_Json]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMapWithDateFormat_Json(t *testing.T) {
	paramMap := map[string]interface{}{"1": 14, "2": 139}
	rows, err := db.SqlMapClient("json_selectAllArticle", &paramMap).QueryWithDateFormat("2006-01-02 15:04").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMapWithDateFormat_Json]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"1": 14, "2": 139}
	rows, err := db.SqlMapClient("json_selectAllArticle", &paramMap).Query().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMap_Xml]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMapWithDateFormat_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"1": 14, "2": 139}
	rows, err := db.SqlMapClient("json_selectAllArticle", &paramMap).QueryWithDateFormat("2006-01-02 15:04").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMapWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"1": 14, "2": 139}
	rows, err := db.SqlMapClient("json_selectAllArticle", &paramMap).Query().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"1": 14, "2": 139}
	rows, err := db.SqlMapClient("json_selectAllArticle", &paramMap).QueryWithDateFormat("2006-01-02 15:04").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMap_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}
	rows, err := db.SqlTemplateClient("select.example.tpl", &paramMap).Query().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMap_Json]->rows:\n" + rows)
}

func Test_Session_SqlTemplateClient_QueryByParamMap_Json(t *testing.T) {
	session := db.NewSession()
	defer session.Close()
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}
	rows, err := session.SqlTemplateClient("select.example.tpl", &paramMap).Query().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMap_Json]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMapWithDateFormat_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}
	rows, err := db.SqlTemplateClient("select.example.tpl", &paramMap).QueryWithDateFormat("01/02/2006").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMapWithDateFormat_Json]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.tpl", &paramMap).Query().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMap_Xml]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMapWithDateFormat_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.tpl", &paramMap).QueryWithDateFormat("01/02/2006").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMapWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.tpl", &paramMap).Query().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", &paramMap).QueryWithDateFormat("01/02/2006").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_Where_Search_Structs_Json(t *testing.T) {
	var articles []Article
	json, err := db.Where("id=?", 6).Search(&articles).Json()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Where_Search_Structs_Json]->rows:\n" + json)
}

func Test_Search_Structs_Xml(t *testing.T) {
	var articles []Article
	xml, err := db.Where("id=?", 6).Search(&articles).Xml()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Search_Structs_Xml]->rows:\n" + xml)
}

func Test_Search_Structs_XmlIndent(t *testing.T) {
	var articles []Article
	xml, err := db.Where("id=?", 6).Search(&articles).XmlIndent("", "  ", "Article")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Search_Structs_XmlIndent]->rows:\n" + xml)
}

func Test_Search_Structs_Json(t *testing.T) {
	var categories []Category
	Json, err := db.Select("id").Search(&categories).Json()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Search_Structs_Json]->rows:\n", Json)
}

func Test_Sql_Find_Structs(t *testing.T) {

	var categories2 []Category
	err := db.Sql("select * from category where id =?", 16).Find(&categories2)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Sql_Find_Structs]->rows:\n", categories2)
}

func Test_SqlMapClient_Find_Structs(t *testing.T) {

	var categories2 []Category
	db.AddSql("1", "select * from category where id =?")
	err := db.SqlMapClient("1", 16).Find(&categories2)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_SqlMapClient_Find_Structs]->rows:\n", categories2)
}

func Test_SqlTemplateClient_Find_Structs(t *testing.T) {

	var categories2 []Category
	db.AddSqlTemplate("1", "select * from category where id =?id")
	err := db.SqlTemplateClient("1", &map[string]interface{}{"id": 25}).Find(&categories2)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_SqlTemplateClient_Find_Structs]->rows:\n", categories2)
}

func Test_Session_SqlTemplateClient_Find_Structs(t *testing.T) {
	session := db.NewSession()
	defer session.Close()
	var categories2 []Category
	db.AddSqlTemplate("1", "select * from category where id =?id")
	err := session.SqlTemplateClient("1", &map[string]interface{}{"id": 25}).Find(&categories2)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_SqlTemplateClient_Find_Structs]->rows:\n", categories2)
}

func Test_Sql_Search_Json(t *testing.T) {

	var categories2 []Category
	json, err := db.Sql("select * from category where id =?", 16).Search(&categories2).Json()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Sql_Search_Json]->rows:\n", json)
}

func Test_SqlMapClient_Search_Json(t *testing.T) {

	var categories2 []Category
	db.AddSql("1", "select * from category where id =?")
	json, err := db.SqlMapClient("1", 16).Search(&categories2).Json()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_SqlMapClient_Search_Json]->rows:\n", json)
}

func Test_SqlTemplateClient_Search_Json(t *testing.T) {

	var categories2 []Category
	db.AddSqlTemplate("1", "select * from category where id =?id")
	json, err := db.SqlTemplateClient("1", &map[string]interface{}{"id": 25}).Search(&categories2).Json()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_SqlTemplateClient_Search_Json]->rows:\n", json)
}

func Test_Query(t *testing.T) {

	result, err := db.QueryInterface("select * from category where id =25")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Query]->rows:\n", result)
}

func STest_Sql_Execute(t *testing.T) {

	result, err := db.Sql("INSERT INTO categories VALUES (?, ?, ?, ?, ?)", 148, "xiaozhang", 1, 1, 1).Execute()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Sql_Execute]->rows:\n", result)
}

func STest_SqlMapClient_Execute(t *testing.T) {
	db.AddSql("Test_SqlMapClient_Execute", "INSERT INTO categories VALUES (?id, ?name, ?counts, ?orders, ?pid)")
	result, err := db.SqlMapClient("Test_SqlMapClient_Execute", &map[string]interface{}{"id": 149, "name": "xiaowang", "counts": 1, "orders": 1, "pid": 1}).Execute()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_SqlMapClient_Execute]->rows:\n", result)
}

func STest_SqlTemplateClientt_Execute(t *testing.T) {
	db.AddSqlTemplate("Test_SqlTemplateClientt_Execute", "INSERT INTO categories VALUES (?id, ?name, ?counts, ?orders, ?pid)")
	result, err := db.SqlTemplateClient("Test_SqlTemplateClientt_Execute", &map[string]interface{}{"id": 240, "name": "laowang", "counts": 1, "orders": 1, "pid": 1}).Execute()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_SqlTemplateClientt_Execute]->rows:\n", result)
}

func Test_GetSQL(t *testing.T) {
	db.AddSql("Test_GetSQL_1", "select * from Test_GetSQL_1")
	t.Log("[Test_GetSQL]->Test_GetSQL_1:\n", db.GetSql("Test_GetSQL_1"))
	t.Log("[Test_GetSQL]->Test_GetSQL_2:\n", db.GetSql("Test_GetSQL_1"))
}

func Test_GetSqlMap(t *testing.T) {

	t.Log("[Test_GetSqlMap]->3:\n")
	t.Log(xorm.JSONString(db.GetSqlMap(3), true))
	sqlmap := db.GetSqlMap(3)
	t.Log("[Test_GetSqlMap]->len(sqlmap):\n", len(sqlmap))

	db.AddSql("Test_GetSqlMap_1", "select * from Test_GetSqlMap_1")
	db.AddSql("Test_GetSqlMap_2", "select * from Test_GetSqlMap_2")
	db.AddSql("Test_GetSqlMap_3", "select * from Test_GetSqlMap_3")
	db.AddSql("Test_GetSqlMap_4", "select * from Test_GetSqlMap_4")
	db.AddSql("Test_GetSqlMap_5", "select * from Test_GetSqlMap_5")
	t.Log("[Test_GetSqlMap]->init->3:\n")
	t.Log(xorm.JSONString(db.GetSqlMap(3), true))
	sqlmap = db.GetSqlMap(3)
	t.Log("[Test_GetSqlMap]->init->len(sqlmap):\n", len(sqlmap))
	t.Log("[Test_GetSqlMap]->Test_GetSqlMap_null:\n")
	t.Log(xorm.JSONString(db.GetSqlMap("Test_GetSqlMap_null"), true))
	t.Log("[Test_GetSqlMap]->Test_GetSqlMap_1:\n")
	t.Log(xorm.JSONString(db.GetSqlMap("Test_GetSqlMap_1"), true))
	t.Log("[Test_GetSqlMap]->Test_GetSqlMap_1,Test_GetSqlMap_3:\n")
	t.Log(xorm.JSONString(db.GetSqlMap("Test_GetSqlMap_1", "Test_GetSqlMap_3"), true))
	t.Log("[Test_GetSqlMap]->Test_GetSqlMap_1,Test_GetSqlMap_3,3:\n")
	t.Log(xorm.JSONString(db.GetSqlMap("Test_GetSqlMap_1", "Test_GetSqlMap_3", 3), true))
	t.Log("[Test_GetSqlMap]->Test_GetSqlMap_1,Test_GetSqlMap_3,3,Test_GetSqlMap_null:\n")
	t.Log(xorm.JSONString(db.GetSqlMap("Test_GetSqlMap_1", "Test_GetSqlMap_3", 3, "Test_GetSqlMap_null"), true))
	t.Log("[Test_GetSqlMap]->Test_GetSqlMap_1,Test_GetSqlMap_3,[]string{Test_GetSqlMap_2, Test_GetSqlMap_4}:\n")
	t.Log(xorm.JSONString(db.GetSqlMap("Test_GetSqlMap_1", "Test_GetSqlMap_3", []string{"Test_GetSqlMap_2", "Test_GetSqlMap_4"}), true))
	t.Log("[Test_GetSqlMap]->Test_GetSqlMap_1,Test_GetSqlMap_3,[]string{Test_GetSqlMap_2, Test_GetSqlMap_4},2:\n")
	t.Log(xorm.JSONString(db.GetSqlMap("Test_GetSqlMap_1", "Test_GetSqlMap_3", []string{"Test_GetSqlMap_2", "Test_GetSqlMap_4"}, 2), true))
	t.Log("[Test_GetSqlMap]->Test_GetSqlMap_1,Test_GetSqlMap_3,[]string{Test_GetSqlMap_2, Test_GetSqlMap_4},2 ,Test_GetSqlMap_null:\n")
	t.Log(xorm.JSONString(db.GetSqlMap("Test_GetSqlMap_1", "Test_GetSqlMap_3", []string{"Test_GetSqlMap_2", "Test_GetSqlMap_4"}, 2, "Test_GetSqlMap_null"), true))
}

func Test_Limit_Func(t *testing.T) {
	res, _ := db.Sql("SELECT b.id,a.name,b.title FROM category a,article b where a.id=b.categorysubid ORDER BY b.id").Limit(10, 3).Query().List()
	t.Log(res)
}

func Test_Find(t *testing.T) {
	var category []Category
	err := db.Find(&category)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(category)
}

func Test_EngineGroup(t *testing.T) {
	conns := []string{
		"postgres://postgres:root@localhost:5432/mblog?sslmode=disable;",
		"postgres://postgres:root@localhost:5432/mblog1?sslmode=disable;",
		"postgres://postgres:root@localhost:5432/mblog2?sslmode=disable",
	}

	var err error
	xg, err := xorm.NewEngineGroup("postgres", conns)

	if err != nil {
		t.Fatal(err)
	}

	xg.ShowSQL(true)
	xg.SetPolicy(xorm.WeightRoundRobinPolicy([]int{2, 3}))

	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())
	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())
	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())
	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())
	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())
	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())
	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())
	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())
	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())
	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())
	t.Log("[xg.Slave():\n", xg.Slave().DataSourceName())

	var category []Category
	err = xg.Find(&category)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(category)
}
