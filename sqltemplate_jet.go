package xorm

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/CloudyKit/jet"
)

type JetTemplate struct {
	SqlTemplateRootDir string
	Template           map[string]*jet.Template
	extension          string
	Capacity           uint
	Cipher             Cipher
	Funcs              map[string]FuncMap
}

func (sqlTemplate *JetTemplate) RootDir() string {
	return sqlTemplate.SqlTemplateRootDir
}

func (sqlTemplate *JetTemplate) SetFuncs(key string, funcMap FuncMap) {
	sqlTemplate.Funcs[key] = funcMap
}

func (sqlTemplate *JetTemplate) Extension() string {
	return sqlTemplate.extension
}

func (sqlTemplate *JetTemplate) SetSqlTemplateCipher(cipher Cipher) {
	sqlTemplate.Cipher = cipher
}

func (sqlTemplate *JetTemplate) WalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if strings.HasSuffix(path, sqlTemplate.extension) {
		err = sqlTemplate.paresSqlTemplate(info.Name(), path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sqlTemplate *JetTemplate) paresSqlTemplate(filename string, filepath string) error {
	var sqlt *jet.Template
	var err error
	var content []byte
	if sqlTemplate.Cipher == nil {
		templates := jet.NewHTMLSet(sqlTemplate.SqlTemplateRootDir)
		sqlt, err = templates.GetTemplate(filename)
		if err != nil {
			return err
		}
	} else {
		content, err = sqlTemplate.ReadTemplate(filepath)
		if err != nil {
			return err
		}
		templates := jet.NewHTMLSet(sqlTemplate.SqlTemplateRootDir)
		sqlt, err = templates.LoadTemplate(filename, string(content))
		if err != nil {
			return err
		}
	}

	sqlTemplate.checkNilAndInit()

	sqlTemplate.Template[filename] = sqlt

	return nil

}

func (sqlTemplate *JetTemplate) ReadTemplate(filepath string) ([]byte, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	content, err = sqlTemplate.Cipher.Decrypt(content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

//-------------------------------------------------------------------------------------------------------------
func (sqlTemplate *JetTemplate) LoadSqlTemplate(filepath string) error {

	if strings.HasSuffix(filepath, sqlTemplate.extension) {
		err := sqlTemplate.loadSqlTemplate(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sqlTemplate *JetTemplate) BatchLoadSqlTemplate(filepathSlice []string) error {

	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, sqlTemplate.extension) {
			err := sqlTemplate.loadSqlTemplate(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (sqlTemplate *JetTemplate) ReloadSqlTemplate(filepath string) error {

	if strings.HasSuffix(filepath, sqlTemplate.extension) {
		err := sqlTemplate.reloadSqlTemplate(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sqlTemplate *JetTemplate) BatchReloadSqlTemplate(filepathSlice []string) error {

	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, sqlTemplate.extension) {
			err := sqlTemplate.loadSqlTemplate(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (sqlTemplate *JetTemplate) loadSqlTemplate(filepath string) error {
	info, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	err = sqlTemplate.paresSqlTemplate(info.Name(), filepath)
	if err != nil {
		return err
	}

	return nil
}

func (sqlTemplate *JetTemplate) reloadSqlTemplate(filepath string) error {
	info, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	err = sqlTemplate.paresSqlTemplate(info.Name(), filepath)
	if err != nil {
		return err
	}

	return nil
}

func (sqlTemplate *JetTemplate) checkNilAndInit() {
	if sqlTemplate.Template == nil {
		sqlTemplate.Template = make(map[string]*jet.Template, 100)
	}
}

func (sqlTemplate *JetTemplate) AddSqlTemplate(key string, sqlTemplateStr string) error {

	templates := jet.NewHTMLSet(sqlTemplate.SqlTemplateRootDir)
	sqlt, err := templates.LoadTemplate(key, sqlTemplateStr)
	if err != nil {
		return err
	}

	sqlTemplate.checkNilAndInit()
	sqlTemplate.Template[key] = sqlt

	return nil

}

func (sqlTemplate *JetTemplate) UpdateSqlTemplate(key string, sqlTemplateStr string) error {

	templates := jet.NewHTMLSet(sqlTemplate.SqlTemplateRootDir)
	sqlt, err := templates.LoadTemplate(key, sqlTemplateStr)
	if err != nil {
		return err
	}

	sqlTemplate.checkNilAndInit()
	sqlTemplate.Template[key] = sqlt

	return nil

}

func (sqlTemplate *JetTemplate) RemoveSqlTemplate(key string) {
	sqlTemplate.checkNilAndInit()
	delete(sqlTemplate.Template, key)
}

func (sqlTemplate *JetTemplate) BatchAddSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {

	templates := jet.NewHTMLSet(sqlTemplate.SqlTemplateRootDir)

	sqlTemplate.checkNilAndInit()

	for k, v := range sqlTemplateStrMap {
		sqlt, err := templates.LoadTemplate(key, v)
		if err != nil {
			return err
		}

		sqlTemplate.Template[k] = sqlt
	}

	return nil

}

func (sqlTemplate *JetTemplate) BatchUpdateSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {
	templates := jet.NewHTMLSet(sqlTemplate.SqlTemplateRootDir)
	sqlTemplate.checkNilAndInit()
	for k, v := range sqlTemplateStrMap {
		sqlt, err := templates.LoadTemplate(key, v)
		if err != nil {
			return err
		}

		sqlTemplate.Template[k] = sqlt
	}

	return nil

}

func (sqlTemplate *JetTemplate) BatchRemoveSqlTemplate(key []string) {
	sqlTemplate.checkNilAndInit()
	for _, v := range key {
		delete(sqlTemplate.Template, v)
	}
}

func (sqlTemplate *JetTemplate) GetSqlTemplate(key string) *jet.Template {
	return sqlTemplate.Template[key]
}

func (sqlTemplate *JetTemplate) GetSqlTemplates(keys ...interface{}) map[string]*jet.Template {

	var resultSqlTemplates map[string]*jet.Template
	i := len(keys)
	if i == 0 {
		return sqlTemplate.Template
	}

	if i == 1 {
		switch keys[0].(type) {
		case string:
			resultSqlTemplates = make(map[string]*jet.Template, 1)
		case []string:
			ks := keys[0].([]string)
			n := len(ks)
			resultSqlTemplates = make(map[string]*jet.Template, n)
		}
	} else {
		resultSqlTemplates = make(map[string]*jet.Template, i)
	}

	for k, _ := range keys {
		switch keys[k].(type) {
		case string:
			key := keys[k].(string)
			resultSqlTemplates[key] = sqlTemplate.Template[key]
		case []string:
			ks := keys[k].([]string)
			for _, v := range ks {
				resultSqlTemplates[v] = sqlTemplate.Template[v]
			}
		}
	}

	return resultSqlTemplates
}

func (sqlTemplate *JetTemplate) Execute(key string, args ...interface{}) (string, error) {
	var buf bytes.Buffer
	if sqlTemplate.Template[key] == nil {
		return "", nil
	}

	if len(args) == 0 {
		err := sqlTemplate.Template[key].Execute(&buf, nil, nil)
		return buf.String(), err
	} else {
		map1 := args[0].(*map[string]interface{})
		vars := make(jet.VarMap)
		vars.Set("data", map1)
		err := sqlTemplate.Template[key].Execute(&buf, vars, nil)
		return buf.String(), err
	}
}
