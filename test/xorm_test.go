package xorm

import (
	"fmt"
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
	//	db.SqlMap.SqlMapRootDir="./sql/oracle"
	//	db.SqlTemplate.SqlTemplateRootDir="./sql/oracle"
	if err != nil {
		t.Fatal(err)
	}

	err = db.InitSqlMap()
	if err != nil {
		t.Fatal(err)
	}
	err = db.InitSqlTemplate()
	if err != nil {
		t.Fatal(err)
	}

	db.ShowSQL(true)
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

func Test_GetFirst_Json(t *testing.T) {
	var article Article
	has, rows, err := db.Id(2).GetFirst(&article).Json()
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

func Test_Find(t *testing.T) {
	var article []Article
	result := db.Sql("select id,title,createdatetime,content from article where id = ?", 27).Find(&article)
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
	t.Log("[Test_Find]-> result.Json():\n", resultJson)
}

func Test_Query_Json(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 27).Query().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_Query_Json]->rows:\n" + rows)
}

func Test_Query_Result(t *testing.T) {
	rows := db.Sql("select id,title,createdatetime,content from article where id = ?", 27).Query()
	if rows.Error != nil {
		t.Fatal(rows.Error)
	}

	t.Log("[Test_Query_Result]->rows[0][\"id\"]:\n", rows.Result[0]["id"])
	t.Log("[Test_Query_Result]->reflect.TypeOf(rows.Result[0][\"id\"]):\n", reflect.TypeOf(rows.Result[0]["id"]))
	t.Log("[Test_Query_Result]->rows[0][\"title\"]:\n", rows.Result[0]["title"])
	t.Log("[Test_Query_Result]->reflect.TypeOf(rows.Result[0][\"title\"]):\n", reflect.TypeOf(rows.Result[0]["title"]))
	t.Log("[Test_Query_Result]->rows[0][\"createdatetime\"]:\n", rows.Result[0]["createdatetime"])
	t.Log("[Test_Query_Result]->reflect.TypeOf(rows.Result[0][\"createdatetime\"]):\n", reflect.TypeOf(rows.Result[0]["createdatetime"]))
}

func Test_Query_Xml(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 27).Query().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_Query_Xml]->rows:\n" + rows)
}

func Test_Query_XmlIndent(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 33).Query().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_Query_XmlIndent]->rows:\n" + rows)
}

func Test_QueryWithDateFormat_Json(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 33).QueryWithDateFormat("20060102").Json()
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
	rows, err := db.Sql("select id,title,createdatetime,content from article where id in (?,?)", 27, 33).QueryWithDateFormat("20060102").XmlIndent("", "  ", "article")
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
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).Query().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMap_Json]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMapWithDateFormat_Json(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryWithDateFormat("2006-01-02 15:04").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMapWithDateFormat_Json]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).Query().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMap_Xml]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMapWithDateFormat_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryWithDateFormat("2006-01-02 15:04").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMapWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).Query().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_SqlMapClient_QueryByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).QueryWithDateFormat("2006-01-02 15:04").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_QueryByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMap_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).Query().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMap_Json]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMapWithDateFormat_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryWithDateFormat("01/02/2006").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMapWithDateFormat_Json]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).Query().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMap_Xml]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMapWithDateFormat_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryWithDateFormat("01/02/2006").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMapWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).Query().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_SqlTemplateClient_QueryByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).QueryWithDateFormat("01/02/2006").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_QueryByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_Find_Structs_Json(t *testing.T) {
	articles := make([]Article, 0)
	json, err := db.Where("id=?", 6).Find(&articles).Json()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Find_Structs_Json]->rows:\n" + json)
}

func Test_Find_Structs_Xml(t *testing.T) {
	articles := make([]Article, 0)
	xml, err := db.Where("id=?", 6).Find(&articles).Xml()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Find_Structs_Xml]->rows:\n" + xml)
}

func Test_Find_Structs_XmlIndent(t *testing.T) {
	articles := make([]Article, 0)
	xml, err := db.Where("id=?", 6).Find(&articles).XmlIndent("", "  ", "Article")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[Test_Find_Structs_XmlIndent]->rows:\n" + xml)
}
