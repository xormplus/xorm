package xorm

import (
	"os"
	"path/filepath"
	"strings"
)

func (engine *Engine) SetSqlMapCipher(cipher Cipher) {
	engine.SqlMap.Cipher = cipher
}

func (engine *Engine) ClearSqlMapCipher() {
	engine.SqlMap.Cipher = nil
}

func (engine *Engine) RegisterSqlMap(sqlm SqlM, Cipher ...Cipher) error {
	switch sqlm.(type) {
	case *XmlSqlMap:
		if len(engine.SqlMap.Extension) == 0 {
			engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json", "xsql": ".sql"}
		}
		engine.SqlMap.Extension["xml"] = sqlm.Extension()
	case *JsonSqlMap:
		if len(engine.SqlMap.Extension) == 0 {
			engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json", "xsql": ".sql"}
		}
		engine.SqlMap.Extension["json"] = sqlm.Extension()
	case *XSqlMap:
		if len(engine.SqlMap.Extension) == 0 {
			engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json", "xsql": ".sql"}
		}
		engine.SqlMap.Extension["xsql"] = sqlm.Extension()
	default:
		return ErrParamsType
	}

	if len(Cipher) > 0 {
		engine.SqlMap.Cipher = Cipher[0]
	}

	engine.SqlMap.SqlMapRootDir = sqlm.RootDir()

	err := filepath.Walk(engine.SqlMap.SqlMapRootDir, engine.SqlMap.walkFunc)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) LoadSqlMap(filepath string) error {

	if len(engine.SqlMap.Extension) == 0 {
		engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json", "xsql": ".sql"}
	} else {
		if engine.SqlMap.Extension["xml"] == "" || len(engine.SqlMap.Extension["xml"]) == 0 {
			engine.SqlMap.Extension["xml"] = ".xml"
		}
		if engine.SqlMap.Extension["json"] == "" || len(engine.SqlMap.Extension["json"]) == 0 {
			engine.SqlMap.Extension["json"] = ".json"
		}
		if engine.SqlMap.Extension["xsql"] == "" || len(engine.SqlMap.Extension["xsql"]) == 0 {
			engine.SqlMap.Extension["xsql"] = ".sql"
		}
	}

	if strings.HasSuffix(filepath, engine.SqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["json"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["xsql"]) {
		err := engine.loadSqlMap(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) BatchLoadSqlMap(filepathSlice []string) error {
	if len(engine.SqlMap.Extension) == 0 {
		engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json", "xsql": ".sql"}
	} else {
		if engine.SqlMap.Extension["xml"] == "" || len(engine.SqlMap.Extension["xml"]) == 0 {
			engine.SqlMap.Extension["xml"] = ".xml"
		}
		if engine.SqlMap.Extension["json"] == "" || len(engine.SqlMap.Extension["json"]) == 0 {
			engine.SqlMap.Extension["json"] = ".json"
		}
		if engine.SqlMap.Extension["xsql"] == "" || len(engine.SqlMap.Extension["xsql"]) == 0 {
			engine.SqlMap.Extension["xsql"] = ".sql"
		}
	}

	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, engine.SqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["json"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["xsql"]) {
			err := engine.loadSqlMap(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (engine *Engine) ReloadSqlMap(filepath string) error {
	if len(engine.SqlMap.Extension) == 0 {
		engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json", "xsql": ".sql"}
	} else {
		if engine.SqlMap.Extension["xml"] == "" || len(engine.SqlMap.Extension["xml"]) == 0 {
			engine.SqlMap.Extension["xml"] = ".xml"
		}
		if engine.SqlMap.Extension["json"] == "" || len(engine.SqlMap.Extension["json"]) == 0 {
			engine.SqlMap.Extension["json"] = ".json"
		}
		if engine.SqlMap.Extension["xsql"] == "" || len(engine.SqlMap.Extension["xsql"]) == 0 {
			engine.SqlMap.Extension["xsql"] = ".sql"
		}
	}

	if strings.HasSuffix(filepath, engine.SqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["json"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["xsql"]) {
		err := engine.reloadSqlMap(filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) BatchReloadSqlMap(filepathSlice []string) error {
	if len(engine.SqlMap.Extension) == 0 {
		engine.SqlMap.Extension = map[string]string{"xml": ".xml", "json": ".json", "xsql": ".sql"}
	} else {
		if engine.SqlMap.Extension["xml"] == "" || len(engine.SqlMap.Extension["xml"]) == 0 {
			engine.SqlMap.Extension["xml"] = ".xml"
		}
		if engine.SqlMap.Extension["json"] == "" || len(engine.SqlMap.Extension["json"]) == 0 {
			engine.SqlMap.Extension["json"] = ".json"
		}
		if engine.SqlMap.Extension["xsql"] == "" || len(engine.SqlMap.Extension["xsql"]) == 0 {
			engine.SqlMap.Extension["xsql"] = ".sql"
		}
	}

	for _, filepath := range filepathSlice {
		if strings.HasSuffix(filepath, engine.SqlMap.Extension["xml"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["json"]) || strings.HasSuffix(filepath, engine.SqlMap.Extension["xsql"]) {
			err := engine.loadSqlMap(filepath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (engine *Engine) loadSqlMap(filepath string) error {

	info, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	err = engine.SqlMap.paresSql(filepath)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) reloadSqlMap(filepath string) error {

	info, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}
	err = engine.SqlMap.paresSql(filepath)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) AddSql(key string, sql string) {
	engine.SqlMap.addSql(key, sql)
}

func (engine *Engine) UpdateSql(key string, sql string) {
	engine.SqlMap.updateSql(key, sql)
}

func (engine *Engine) RemoveSql(key string) {
	engine.SqlMap.removeSql(key)
}

func (engine *Engine) BatchAddSql(sqlStrMap map[string]string) {
	engine.SqlMap.batchAddSql(sqlStrMap)
}

func (engine *Engine) BatchUpdateSql(sqlStrMap map[string]string) {
	engine.SqlMap.batchUpdateSql(sqlStrMap)
}

func (engine *Engine) BatchRemoveSql(key []string) {
	engine.SqlMap.batchRemoveSql(key)
}

func (engine *Engine) GetSql(key string) string {
	return engine.SqlMap.getSql(key)
}

func (engine *Engine) GetSqlMap(keys ...interface{}) map[string]string {
	return engine.SqlMap.getSqlMap(keys...)
}
