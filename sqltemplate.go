package xorm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
	"gopkg.in/flosch/pongo2.v3"
)

type SqlTemplate struct {
	SqlTemplateRootDir string
	Template           map[string]*pongo2.Template
	Extension          string
	Capacity           uint
	Cipher             Cipher
}

type SqlTemplateOptions struct {
	Capacity  uint
	Extension string
	Cipher    Cipher
}

func (engine *Engine) SetSqlTemplateCipher(cipher Cipher) {
	engine.sqlTemplate.Cipher = cipher
}

func (engine *Engine) ClearSqlTemplateCipher() {
	engine.sqlTemplate.Cipher = nil
}

func (sqlTemplate *SqlTemplate) checkNilAndInit() {
	if sqlTemplate.Template == nil {
		if sqlTemplate.Capacity == 0 {
			sqlTemplate.Template = make(map[string]*pongo2.Template, 100)
		} else {
			sqlTemplate.Template = make(map[string]*pongo2.Template, sqlTemplate.Capacity)
		}
	}
}

func (engine *Engine) InitSqlTemplate(options ...SqlTemplateOptions) error {
	var opt SqlTemplateOptions

	if len(options) > 0 {
		opt = options[0]
	}

	if len(opt.Extension) == 0 {
		opt.Extension = ".stpl"
	}
	engine.sqlTemplate.Extension = opt.Extension
	engine.sqlTemplate.Capacity = opt.Capacity

	engine.sqlTemplate.Cipher = opt.Cipher

	var err error
	if engine.sqlTemplate.SqlTemplateRootDir == "" {
		cfg, err := goconfig.LoadConfigFile("./sql/xormcfg.ini")
		if err != nil {
			return err
		}
		engine.sqlTemplate.SqlTemplateRootDir, err = cfg.GetValue("", "SqlTemplateRootDir")
		if err != nil {
			return err
		}
	}

	err = filepath.Walk(engine.sqlTemplate.SqlTemplateRootDir, engine.sqlTemplate.walkFunc)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) LoadSqlTemplate(filepath string) error {
	if len(engine.sqlTemplate.Extension) == 0 {
		engine.sqlTemplate.Extension = ".stpl"
	}
	if strings.HasSuffix(filepath, engine.sqlTemplate.Extension) {
		err := engine.loadSqlTemplate(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) BatchLoadSqlTemplate(filepathSlice []string) error {
	if len(engine.sqlTemplate.Extension) == 0 {
		engine.sqlTemplate.Extension = ".stpl"
	}
	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, engine.sqlTemplate.Extension) {
			err := engine.loadSqlTemplate(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (engine *Engine) ReloadSqlTemplate(filepath string) error {
	if len(engine.sqlTemplate.Extension) == 0 {
		engine.sqlTemplate.Extension = ".stpl"
	}
	if strings.HasSuffix(filepath, engine.sqlTemplate.Extension) {
		err := engine.reloadSqlTemplate(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) BatchReloadSqlTemplate(filepathSlice []string) error {
	if len(engine.sqlTemplate.Extension) == 0 {
		engine.sqlTemplate.Extension = ".stpl"
	}
	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, engine.sqlTemplate.Extension) {
			err := engine.loadSqlTemplate(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (engine *Engine) loadSqlTemplate(filepath string) error {
	info, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	err = engine.sqlTemplate.paresSqlTemplate(info.Name(), filepath)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) reloadSqlTemplate(filepath string) error {
	info, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	err = engine.sqlTemplate.paresSqlTemplate(info.Name(), filepath)
	if err != nil {
		return err
	}

	return nil
}

func (sqlTemplate *SqlTemplate) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if strings.HasSuffix(path, sqlTemplate.Extension) {
		err = sqlTemplate.paresSqlTemplate(info.Name(), path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sqlTemplate *SqlTemplate) paresSqlTemplate(filename string, filepath string) error {
	var template *pongo2.Template
	var err error
	var content []byte

	if sqlTemplate.Cipher == nil {
		template, err = pongo2.FromFile(filepath)
		if err != nil {
			return err
		}
	} else {
		content, err = ioutil.ReadFile(filepath)

		if err != nil {
			return err
		}
		content, err = sqlTemplate.Cipher.Decrypt(content)
		if err != nil {
			return err
		}
		template, err = pongo2.FromString(string(content))
		if err != nil {
			return err
		}
	}

	sqlTemplate.checkNilAndInit()
	sqlTemplate.Template[filename] = template

	return nil
}

func (engine *Engine) AddSqlTemplate(key string, sqlTemplateStr string) error {
	return engine.sqlTemplate.addSqlTemplate(key, sqlTemplateStr)
}

func (sqlTemplate *SqlTemplate) addSqlTemplate(key string, sqlTemplateStr string) error {

	template, err := pongo2.FromString(sqlTemplateStr)
	if err != nil {
		return err
	}

	sqlTemplate.checkNilAndInit()
	sqlTemplate.Template[key] = template

	return nil

}

func (engine *Engine) UpdateSqlTemplate(key string, sqlTemplateStr string) error {
	return engine.sqlTemplate.updateSqlTemplate(key, sqlTemplateStr)
}

func (sqlTemplate *SqlTemplate) updateSqlTemplate(key string, sqlTemplateStr string) error {

	template, err := pongo2.FromString(sqlTemplateStr)
	if err != nil {
		return err
	}
	sqlTemplate.checkNilAndInit()
	sqlTemplate.Template[key] = template

	return nil

}

func (engine *Engine) RemoveSqlTemplate(key string) {
	engine.sqlTemplate.removeSqlTemplate(key)
}

func (sqlTemplate *SqlTemplate) removeSqlTemplate(key string) {
	sqlTemplate.checkNilAndInit()
	delete(sqlTemplate.Template, key)
}

func (engine *Engine) BatchAddSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {
	return engine.sqlTemplate.batchAddSqlTemplate(key, sqlTemplateStrMap)

}

func (sqlTemplate *SqlTemplate) batchAddSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {
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

func (engine *Engine) BatchUpdateSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {
	return engine.sqlTemplate.batchAddSqlTemplate(key, sqlTemplateStrMap)

}

func (sqlTemplate *SqlTemplate) batchUpdateSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {
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

func (engine *Engine) BatchRemoveSqlTemplate(key []string) {
	engine.sqlTemplate.batchRemoveSqlTemplate(key)
}

func (sqlTemplate *SqlTemplate) batchRemoveSqlTemplate(key []string) {
	sqlTemplate.checkNilAndInit()
	for _, v := range key {
		delete(sqlTemplate.Template, v)
	}
}

func (engine *Engine) GetSqlTemplate(key string) *pongo2.Template {
	return engine.sqlTemplate.getSqlTemplate(key)
}

func (sqlTemplate *SqlTemplate) getSqlTemplate(key string) *pongo2.Template {
	return sqlTemplate.Template[key]
}

func (engine *Engine) GetSqlTemplates(keys ...interface{}) map[string]*pongo2.Template {
	return engine.sqlTemplate.getSqlTemplates(keys...)
}

func (sqlTemplate *SqlTemplate) getSqlTemplates(keys ...interface{}) map[string]*pongo2.Template {

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
