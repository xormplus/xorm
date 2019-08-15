package xorm

import (
	"html/template"
	"os"

	"github.com/CloudyKit/jet"
	"gopkg.in/flosch/pongo2.v3"
)

type FuncMap map[string]interface{}

type SqlTemplate interface {
	WalkFunc(path string, info os.FileInfo, err error) error
	paresSqlTemplate(filename string, filepath string) error
	ReadTemplate(filepath string) ([]byte, error)
	Execute(key string, args ...interface{}) (string, error)
	RootDir() string
	Extension() string
	SetSqlTemplateCipher(cipher Cipher)
	LoadSqlTemplate(filepath string) error
	BatchLoadSqlTemplate(filepathSlice []string) error
	ReloadSqlTemplate(filepath string) error
	BatchReloadSqlTemplate(filepathSlice []string) error
	AddSqlTemplate(key string, sqlTemplateStr string) error
	UpdateSqlTemplate(key string, sqlTemplateStr string) error
	RemoveSqlTemplate(key string)
	BatchAddSqlTemplate(key string, sqlTemplateStrMap map[string]string) error
	BatchUpdateSqlTemplate(key string, sqlTemplateStrMap map[string]string) error
	BatchRemoveSqlTemplate(key []string)
	SetFuncs(key string, funcMap FuncMap)
}

func Pongo2(directory, extension string) *Pongo2Template {
	template := make(map[string]*pongo2.Template, 100)
	funcs := make(map[string]FuncMap, 20)
	return &Pongo2Template{
		SqlTemplateRootDir: directory,
		extension:          extension,
		Template:           template,
		Funcs:              funcs,
	}
}

func Default(directory, extension string) *HTMLTemplate {
	template := make(map[string]*template.Template, 100)
	funcs := make(map[string]FuncMap, 20)
	return &HTMLTemplate{
		SqlTemplateRootDir: directory,
		extension:          extension,
		Template:           template,
		Funcs:              funcs,
	}
}

func Jet(directory, extension string) *JetTemplate {
	template := make(map[string]*jet.Template, 100)
	funcs := make(map[string]FuncMap, 20)
	return &JetTemplate{
		SqlTemplateRootDir: directory,
		extension:          extension,
		Template:           template,
		Funcs:              funcs,
	}
}
