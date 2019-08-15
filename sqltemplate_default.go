package xorm

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"html/template"
)

type HTMLTemplate struct {
	SqlTemplateRootDir string
	Template           map[string]*template.Template
	extension          string
	Capacity           uint
	Cipher             Cipher
	Type               int
	Funcs              map[string]FuncMap
}

func (sqlTemplate *HTMLTemplate) RootDir() string {
	return sqlTemplate.SqlTemplateRootDir
}

func (sqlTemplate *HTMLTemplate) SetFuncs(key string, funcMap FuncMap) {
	sqlTemplate.Funcs[key] = funcMap
}

func (sqlTemplate *HTMLTemplate) Extension() string {
	return sqlTemplate.extension
}

func (sqlTemplate *HTMLTemplate) SetSqlTemplateCipher(cipher Cipher) {
	sqlTemplate.Cipher = cipher
}

func (sqlTemplate *HTMLTemplate) WalkFunc(path string, info os.FileInfo, err error) error {
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

func (sqlTemplate *HTMLTemplate) paresSqlTemplate(filename string, filepath string) error {
	var sqlt *template.Template
	var err error
	var content []byte

	fmap := sqlTemplate.Funcs[filename]
	if fmap != nil {
		funcMap := make(template.FuncMap, 20)
		for k := range fmap {
			funcMap[k] = fmap[k]
		}
		sqlt = template.New(filename)
		sqlt = sqlt.Funcs(funcMap)
	}

	if fmap == nil {
		if sqlTemplate.Cipher == nil {
			sqlt = template.Must(template.ParseFiles(filepath))
		} else {
			content, err = sqlTemplate.ReadTemplate(filepath)
			if err != nil {
				return err
			}
			sqlt = template.Must(template.New(filename).Parse(string(content)))
		}
	} else {

		if sqlTemplate.Cipher == nil {
			sqlt, err = sqlt.ParseFiles(filepath)
		} else {
			content, err = sqlTemplate.ReadTemplate(filepath)
			if err != nil {
				return err
			}
			sqlt, err = sqlt.Parse(string(content))
		}
		if err != nil {
			return err
		}

	}

	sqlTemplate.checkNilAndInit()
	sqlTemplate.Template[filename] = sqlt
	return nil

}

func (sqlTemplate *HTMLTemplate) ReadTemplate(filepath string) ([]byte, error) {
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
func (sqlTemplate *HTMLTemplate) LoadSqlTemplate(filepath string) error {

	if strings.HasSuffix(filepath, sqlTemplate.extension) {
		err := sqlTemplate.loadSqlTemplate(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sqlTemplate *HTMLTemplate) BatchLoadSqlTemplate(filepathSlice []string) error {

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

func (sqlTemplate *HTMLTemplate) ReloadSqlTemplate(filepath string) error {

	if strings.HasSuffix(filepath, sqlTemplate.extension) {
		err := sqlTemplate.reloadSqlTemplate(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sqlTemplate *HTMLTemplate) BatchReloadSqlTemplate(filepathSlice []string) error {

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

func (sqlTemplate *HTMLTemplate) loadSqlTemplate(filepath string) error {
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

func (sqlTemplate *HTMLTemplate) reloadSqlTemplate(filepath string) error {
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

func (sqlTemplate *HTMLTemplate) checkNilAndInit() {
	if sqlTemplate.Template == nil {
		sqlTemplate.Template = make(map[string]*template.Template, 100)
	}
}

func (sqlTemplate *HTMLTemplate) AddSqlTemplate(key string, sqlTemplateStr string) error {

	var sqlt *template.Template
	var err error
	sqlt = template.New(key)
	sqlt, err = sqlt.Parse(sqlTemplateStr)
	if err != nil {
		return err
	}

	sqlTemplate.checkNilAndInit()
	sqlTemplate.Template[key] = sqlt

	return nil

}

func (sqlTemplate *HTMLTemplate) UpdateSqlTemplate(key string, sqlTemplateStr string) error {

	var sqlt *template.Template
	var err error
	sqlt = template.New(key)
	sqlt, err = sqlt.Parse(sqlTemplateStr)
	if err != nil {
		return err
	}

	sqlTemplate.checkNilAndInit()
	sqlTemplate.Template[key] = sqlt

	return nil

}

func (sqlTemplate *HTMLTemplate) RemoveSqlTemplate(key string) {
	sqlTemplate.checkNilAndInit()
	delete(sqlTemplate.Template, key)
}

func (sqlTemplate *HTMLTemplate) BatchAddSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {

	sqlTemplate.checkNilAndInit()

	for k, v := range sqlTemplateStrMap {
		sqlt := template.New(key)
		sqlt, err := sqlt.Parse(v)
		if err != nil {
			return err
		}

		sqlTemplate.Template[k] = sqlt
	}

	return nil

}

func (sqlTemplate *HTMLTemplate) BatchUpdateSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {
	sqlTemplate.checkNilAndInit()
	for k, v := range sqlTemplateStrMap {
		sqlt := template.New(key)
		sqlt, err := sqlt.Parse(v)
		if err != nil {
			return err
		}

		sqlTemplate.Template[k] = sqlt
	}

	return nil

}

func (sqlTemplate *HTMLTemplate) BatchRemoveSqlTemplate(key []string) {
	sqlTemplate.checkNilAndInit()
	for _, v := range key {
		delete(sqlTemplate.Template, v)
	}
}

func (sqlTemplate *HTMLTemplate) GetSqlTemplate(key string) *template.Template {
	return sqlTemplate.Template[key]
}

func (sqlTemplate *HTMLTemplate) GetSqlTemplates(keys ...interface{}) map[string]*template.Template {

	var resultSqlTemplates map[string]*template.Template
	i := len(keys)
	if i == 0 {
		return sqlTemplate.Template
	}

	if i == 1 {
		switch keys[0].(type) {
		case string:
			resultSqlTemplates = make(map[string]*template.Template, 1)
		case []string:
			ks := keys[0].([]string)
			n := len(ks)
			resultSqlTemplates = make(map[string]*template.Template, n)
		}
	} else {
		resultSqlTemplates = make(map[string]*template.Template, i)
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

func (sqlTemplate *HTMLTemplate) Execute(key string, args ...interface{}) (string, error) {
	var buf bytes.Buffer
	if sqlTemplate.Template[key] == nil {
		return "", nil
	}

	if len(args) == 0 {
		err := sqlTemplate.Template[key].Execute(&buf, nil)
		return buf.String(), err
	} else {
		map1 := args[0].(*map[string]interface{})
		err := sqlTemplate.Template[key].Execute(&buf, *map1)
		return buf.String(), err
	}
}
