package xorm

import (
	"fmt"
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
	has, rows, err := db.Where("userid =?", 3).GetFirst(&article).Xml()
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
	has, rows, err := db.Where("userid =?", 3).GetFirst(&article).XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("[Test_GetFirst_XmlIndent]->rows: not exist\n")
	}
	t.Log("[Test_GetFirst_XmlIndent]->rows:\n" + rows)
}

func Test_FindAll_Json(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).FindAll().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAll_Json]->rows:\n" + rows)
}

func Test_FindAll_ID(t *testing.T) {
	rows := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).FindAll()
	if rows.Error != nil {
		t.Fatal(rows.Error)
	}
	t.Log("[Test_FindAll_Json]->rows[0][\"id\"]:\n", rows.Result[0]["id"])
}

func Test_FindAll_Xml(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).FindAll().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAll_Xml]->rows:\n" + rows)
}

func Test_FindAll_XmlIndent(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).FindAll().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAll_XmlIndent]->rows:\n" + rows)
}

func Test_FindAllWithDateFormat_Json(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).FindAllWithDateFormat("20060102").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllWithDateFormat_Json]->rows:\n" + rows)
}

func Test_FindAllWithDateFormat_Xml(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?", 2).FindAllWithDateFormat("20060102").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_FindAllWithDateFormat_XmlIndent(t *testing.T) {
	rows, err := db.Sql("select id,title,createdatetime,content from article where id in (?,?)", 2, 5).FindAllWithDateFormat("20060102").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_FindAllByParamMap_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 4, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).FindAllByParamMap().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllByParamMap_Json]->rows:\n" + rows)
}

func Test_FindAllByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 6, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).FindAllByParamMap().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllByParamMap_Xml]->rows:\n" + rows)
}

func Test_FindAllByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 6, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).FindAllByParamMap().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_FindAllByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 5, "userid": 1}
	rows, err := db.Sql("select id,title,createdatetime,content from article where id = ?id and userid=?userid", &paramMap).FindAllByParamMapWithDateFormat("2006/01/02").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_FindAllByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMap_Json(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).FindAllByParamMap().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMap_Json]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMapWithDateFormat_Json(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).FindAllByParamMapWithDateFormat("2006-01-02 15:04").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMapWithDateFormat_Json]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).FindAllByParamMap().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMap_Xml]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMapWithDateFormat_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).FindAllByParamMapWithDateFormat("2006-01-02 15:04").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMapWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).FindAllByParamMap().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_SqlMapClient_FindAllByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"1": 2, "2": 5}
	rows, err := db.SqlMapClient("selectAllArticle", &paramMap).FindAllByParamMapWithDateFormat("2006-01-02 15:04").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlMapClient_FindAllByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMap_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).FindAllByParamMap().Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMap_Json]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_Json(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 1}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).FindAllByParamMapWithDateFormat("01/02/2006").Json()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_Json]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMap_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).FindAllByParamMap().Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMap_Xml]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_Xml(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).FindAllByParamMapWithDateFormat("01/02/2006").Xml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_Xml]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMap_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).FindAllByParamMap().XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMap_XmlIndent]->rows:\n" + rows)
}

func Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_XmlIndent(t *testing.T) {
	paramMap := map[string]interface{}{"id": 2, "userid": 3, "count": 2}
	rows, err := db.SqlTemplateClient("select.example.stpl", paramMap).FindAllByParamMapWithDateFormat("01/02/2006").XmlIndent("", "  ", "article")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("[Test_SqlTemplateClient_FindAllByParamMapWithDateFormat_XmlIndent]->rows:\n" + rows)
}
