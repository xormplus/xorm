package xorm

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Unknwon/goconfig"
	"gopkg.in/flosch/pongo2.v3"
)

type SqlTemplate struct {
	Template map[string]*pongo2.Template
}

func (engine *Engine) InitSqlTemplate() error {
	var err error
	cfg, err := goconfig.LoadConfigFile("./sql/xormcfg.ini")
	if err != nil {
		return err
	}
	var sqlMapRootDir string
	sqlMapRootDir, err = cfg.GetValue("", "sqlMapRootDir")
	if err != nil {
		return err
	}
	
	engine.SqlTemplate.Template = make(map[string]*pongo2.Template)
	err = filepath.Walk(sqlMapRootDir, engine.SqlTemplate.walkFunc)
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

	if strings.HasSuffix(path, ".stpl") {
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
