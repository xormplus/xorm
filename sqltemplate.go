package xorm

import (
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
}

type SqlTemplateOptions struct {
	Extension string
}

func (engine *Engine) InitSqlTemplate(options ...SqlTemplateOptions) error {
	var opt SqlTemplateOptions
	if len(options) > 0 {
		opt = options[0]
	}

	if len(opt.Extension) == 0 {
		opt.Extension = ".stpl"
	}

	var err error
	if engine.SqlTemplate.SqlTemplateRootDir == "" {
		cfg, err := goconfig.LoadConfigFile("./sql/xormcfg.ini")
		if err != nil {
			return err
		}
		engine.SqlTemplate.SqlTemplateRootDir, err = cfg.GetValue("", "SqlTemplateRootDir")
		if err != nil {
			return err
		}
	}

	engine.SqlTemplate.Template = make(map[string]*pongo2.Template)
	err = filepath.Walk(engine.SqlTemplate.SqlTemplateRootDir, engine.SqlTemplate.walkFunc)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) LoadSqlTemplate(filepath string) error {
	if strings.HasSuffix(filepath, engine.SqlTemplate.Extension) {
		err := engine.loadSqlTemplate(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) ReloadSqlTemplate(filepath string) error {
	if strings.HasSuffix(filepath, engine.SqlTemplate.Extension) {
		err := engine.reloadSqlTemplate(filepath)
		if err != nil {
			return err
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

	err = engine.SqlTemplate.paresSqlTemplate(info.Name(), filepath)
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

	err = engine.SqlTemplate.paresSqlTemplate(info.Name(), filepath)
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
	template, err := pongo2.FromFile(filepath)
	if err != nil {
		return err
	}

	sqlTemplate.Template[filename] = template

	return nil
}

func (engine *Engine) AddSqlTemplate(key string, sqlTemplateStr string) error {
	return engine.SqlTemplate.addSqlTemplate(key, sqlTemplateStr)

}

func (sqlTemplate *SqlTemplate) addSqlTemplate(key string, sqlTemplateStr string) error {

	template, err := pongo2.FromString(sqlTemplateStr)
	if err != nil {
		return err
	}

	sqlTemplate.Template[key] = template

	return nil

}

func (engine *Engine) UpdateSqlTemplate(key string, sqlTemplateStr string) error {
	return engine.SqlTemplate.updateSqlTemplate(key, sqlTemplateStr)
}

func (sqlTemplate *SqlTemplate) updateSqlTemplate(key string, sqlTemplateStr string) error {

	template, err := pongo2.FromString(sqlTemplateStr)
	if err != nil {
		return err
	}

	sqlTemplate.Template[key] = template

	return nil

}

func (engine *Engine) RemoveSqlTemplate(key string) {
	engine.SqlTemplate.removeSqlTemplate(key)
}

func (sqlTemplate *SqlTemplate) removeSqlTemplate(key string) {
	delete(sqlTemplate.Template, key)
}
