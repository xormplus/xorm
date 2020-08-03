package xorm

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/flosch/pongo2.v3"
)

type Pongo2Template struct {
	SqlTemplateRootDir string
	Template           map[string]*pongo2.Template
	extension          string
	Capacity           uint
	Cipher             Cipher
	Type               int
	Funcs              map[string]FuncMap
}

func (sqlTemplate *Pongo2Template) RootDir() string {
	return sqlTemplate.SqlTemplateRootDir
}

func (sqlTemplate *Pongo2Template) SetFuncs(key string, funcMap FuncMap) {
	sqlTemplate.Funcs[key] = funcMap
}

func (sqlTemplate *Pongo2Template) Extension() string {
	return sqlTemplate.extension
}

func (sqlTemplate *Pongo2Template) SetSqlTemplateCipher(cipher Cipher) {
	sqlTemplate.Cipher = cipher
}

func (sqlTemplate *Pongo2Template) WalkFunc(path string, info os.FileInfo, err error) error {
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

func (sqlTemplate *Pongo2Template) paresSqlTemplate(filename string, filepath string) error {
	var sqlt *pongo2.Template
	var err error
	var content []byte

	if sqlTemplate.Cipher == nil {
		buf, err2 := ioutil.ReadFile(filepath)
		if err2 != nil {
			return err2
		}
		var sqlId []byte
		var sqlValue []byte
		var loaded string
		addSqlValue := false
		addSqlId := false
		for _, b := range buf {
			if len(loaded) >3 && "#id:" == loaded[len(loaded)-4:]{
				addSqlId = true
				addSqlValue = false
				if len(sqlId)>0 && len(sqlValue)>0 {
					sqlt, err = pongo2.FromString(string(sqlValue)[0:len(sqlValue)-4])
					if err != nil {
						return err
					}
					sqlTemplate.checkNilAndInit()
					sqlTemplate.Template[filename+"."+string(sqlId)] = sqlt
					sqlId = nil
					sqlValue = nil
				}

			}
			if addSqlId && string(b) == "\n"{
				addSqlId = false
				addSqlValue = true
			}
			if addSqlId {
				sqlId = append(sqlId, b)
			}
			if addSqlValue {
				sqlValue = append(sqlValue, b)
			}
			loaded = loaded+string(b)
		}
		if len(sqlId)>0 && len(sqlValue)>0 {
			sqlt, err = pongo2.FromString(string(sqlValue))
			if err != nil {
				return err
			}
			sqlTemplate.checkNilAndInit()
			sqlTemplate.Template[filename+"."+string(sqlId)] = sqlt
		} else {
			sqlt, err = pongo2.FromFile(filepath)
			if err != nil {
				return err
			}
			sqlTemplate.checkNilAndInit()
			sqlTemplate.Template[filename] = sqlt
		}
	} else {
		content, err = sqlTemplate.ReadTemplate(filepath)
		if err != nil {
			return err
		}
		sqlt, err = pongo2.FromString(string(content))
		if err != nil {
			return err
		}
		
		sqlTemplate.checkNilAndInit()

		sqlTemplate.Template[filename] = sqlt
	}

	return nil

}

func (sqlTemplate *Pongo2Template) ReadTemplate(filepath string) ([]byte, error) {
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
func (sqlTemplate *Pongo2Template) LoadSqlTemplate(filepath string) error {

	if strings.HasSuffix(filepath, sqlTemplate.extension) {
		err := sqlTemplate.loadSqlTemplate(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sqlTemplate *Pongo2Template) BatchLoadSqlTemplate(filepathSlice []string) error {

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

func (sqlTemplate *Pongo2Template) ReloadSqlTemplate(filepath string) error {

	if strings.HasSuffix(filepath, sqlTemplate.extension) {
		err := sqlTemplate.reloadSqlTemplate(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sqlTemplate *Pongo2Template) BatchReloadSqlTemplate(filepathSlice []string) error {

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

func (sqlTemplate *Pongo2Template) loadSqlTemplate(filepath string) error {
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

func (sqlTemplate *Pongo2Template) reloadSqlTemplate(filepath string) error {
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

func (sqlTemplate *Pongo2Template) checkNilAndInit() {
	if sqlTemplate.Template == nil {
		sqlTemplate.Template = make(map[string]*pongo2.Template, 100)
	}
}

func (sqlTemplate *Pongo2Template) AddSqlTemplate(key string, sqlTemplateStr string) error {

	template, err := pongo2.FromString(sqlTemplateStr)
	if err != nil {
		return err
	}

	sqlTemplate.checkNilAndInit()
	sqlTemplate.Template[key] = template

	return nil

}

func (sqlTemplate *Pongo2Template) UpdateSqlTemplate(key string, sqlTemplateStr string) error {

	template, err := pongo2.FromString(sqlTemplateStr)
	if err != nil {
		return err
	}
	sqlTemplate.checkNilAndInit()
	sqlTemplate.Template[key] = template

	return nil

}

func (sqlTemplate *Pongo2Template) RemoveSqlTemplate(key string) {
	sqlTemplate.checkNilAndInit()
	delete(sqlTemplate.Template, key)
}

func (sqlTemplate *Pongo2Template) BatchAddSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {
	sqlTemplate.checkNilAndInit()
	for k, v := range sqlTemplateStrMap {
		template, err := pongo2.FromString(v)
		if err != nil {
			return err
		}

		sqlTemplate.Template[k] = template
	}

	return nil

}

func (sqlTemplate *Pongo2Template) BatchUpdateSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {
	sqlTemplate.checkNilAndInit()
	for k, v := range sqlTemplateStrMap {
		template, err := pongo2.FromString(v)
		if err != nil {
			return err
		}

		sqlTemplate.Template[k] = template
	}

	return nil

}

func (sqlTemplate *Pongo2Template) BatchRemoveSqlTemplate(key []string) {
	sqlTemplate.checkNilAndInit()
	for _, v := range key {
		delete(sqlTemplate.Template, v)
	}
}

func (sqlTemplate *Pongo2Template) GetSqlTemplate(key string) *pongo2.Template {
	return sqlTemplate.Template[key]
}

func (sqlTemplate *Pongo2Template) GetSqlTemplates(keys ...interface{}) map[string]*pongo2.Template {

	var resultSqlTemplates map[string]*pongo2.Template
	i := len(keys)
	if i == 0 {
		return sqlTemplate.Template
	}

	if i == 1 {
		switch keys[0].(type) {
		case string:
			resultSqlTemplates = make(map[string]*pongo2.Template, 1)
		case []string:
			ks := keys[0].([]string)
			n := len(ks)
			resultSqlTemplates = make(map[string]*pongo2.Template, n)
		}
	} else {
		resultSqlTemplates = make(map[string]*pongo2.Template, i)
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

func (sqlTemplate *Pongo2Template) Execute(key string, args ...interface{}) (string, error) {
	if sqlTemplate.Template[key] == nil {
		return "", nil
	}

	if len(args) == 0 {
		parmap := &pongo2.Context{"1": 1}
		sql, err := sqlTemplate.Template[key].Execute(*parmap)
		return sql, err
	} else {
		map1 := args[0].(*map[string]interface{})
		sql, err := sqlTemplate.Template[key].Execute(*map1)
		return sql, err
	}
}
