package xorm

import (
	"path/filepath"
)

func (engine *Engine) RegisterSqlTemplate(sqlt SqlTemplate, Cipher ...Cipher) error {
	engine.SqlTemplate = sqlt
	if len(Cipher) > 0 {
		engine.SqlTemplate.SetSqlTemplateCipher(Cipher[0])
	}
	err := filepath.Walk(engine.SqlTemplate.RootDir(), engine.SqlTemplate.WalkFunc)
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) LoadSqlTemplate(filepath string) error {
	return engine.SqlTemplate.LoadSqlTemplate(filepath)
}

func (engine *Engine) BatchLoadSqlTemplate(filepathSlice []string) error {
	return engine.SqlTemplate.BatchLoadSqlTemplate(filepathSlice)
}

func (engine *Engine) ReloadSqlTemplate(filepath string) error {
	return engine.SqlTemplate.ReloadSqlTemplate(filepath)
}

func (engine *Engine) BatchReloadSqlTemplate(filepathSlice []string) error {
	return engine.SqlTemplate.BatchReloadSqlTemplate(filepathSlice)
}

func (engine *Engine) AddSqlTemplate(key string, sqlTemplateStr string) error {
	return engine.SqlTemplate.AddSqlTemplate(key, sqlTemplateStr)
}

func (engine *Engine) UpdateSqlTemplate(key string, sqlTemplateStr string) error {
	return engine.SqlTemplate.UpdateSqlTemplate(key, sqlTemplateStr)
}

func (engine *Engine) RemoveSqlTemplate(key string) {
	engine.SqlTemplate.RemoveSqlTemplate(key)
}

func (engine *Engine) BatchAddSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {
	return engine.SqlTemplate.BatchAddSqlTemplate(key, sqlTemplateStrMap)

}

func (engine *Engine) BatchUpdateSqlTemplate(key string, sqlTemplateStrMap map[string]string) error {
	return engine.SqlTemplate.BatchUpdateSqlTemplate(key, sqlTemplateStrMap)

}

func (engine *Engine) BatchRemoveSqlTemplate(key []string) {
	engine.SqlTemplate.BatchRemoveSqlTemplate(key)
}
