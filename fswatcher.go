package xorm

import (
	"strings"

	"github.com/fsnotify/fsnotify"
)

//start filesystem watcher
func (engine *Engine) StartFSWatcher() error {
	var err error
	engine.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {

		for {

			select {
			case event := <-engine.watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if strings.HasSuffix(event.Name, engine.sqlTemplate.Extension) {
						err = engine.reloadSqlTemplate(event.Name)
						if err != nil {
							engine.logger.Error(err)
						}
					}

					if strings.HasSuffix(event.Name, engine.sqlMap.Extension) {
						err = engine.reloadSqlMap(event.Name)
						if err != nil {
							engine.logger.Error(err)
						}
					}
				}

			case err := <-engine.watcher.Errors:
				if err != nil {
					engine.logger.Error(err)
				}
			}
		}
	}()

	if engine.sqlMap.SqlMapRootDir != "" {
		err = engine.watcher.Add(engine.sqlMap.SqlMapRootDir)
		if err != nil {
			return err
		}
	}

	if engine.sqlTemplate.SqlTemplateRootDir != "" {
		err = engine.watcher.Add(engine.sqlTemplate.SqlTemplateRootDir)
		if err != nil {
			return err
		}
	}

	return nil
}

//stop filesystem watcher
func (engine *Engine) StopFSWatcher() error {
	return engine.watcher.Close()
}
