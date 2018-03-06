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
					if strings.HasSuffix(event.Name, engine.SqlTemplate.Extension()) {
						err = engine.ReloadSqlTemplate(event.Name)
						if err != nil {
							engine.logger.Error(err)
						}
					}

					if strings.HasSuffix(event.Name, engine.SqlMap.Extension["xml"]) || strings.HasSuffix(event.Name, engine.SqlMap.Extension["json"]) || strings.HasSuffix(event.Name, engine.SqlMap.Extension["xsql"]) {
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

	if engine.SqlMap.SqlMapRootDir != "" {
		err = engine.watcher.Add(engine.SqlMap.SqlMapRootDir)
		if err != nil {
			return err
		}
	}

	if engine.SqlTemplate.RootDir() != "" {
		err = engine.watcher.Add(engine.SqlTemplate.RootDir())
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
