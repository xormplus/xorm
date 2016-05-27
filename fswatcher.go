package xorm

import (
	"strings"

	"github.com/fsnotify/fsnotify"
)

//start filesytem watcher
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
					if strings.HasSuffix(event.Name, engine.SqlTemplate.Extension) {
						err = engine.reloadSqlTemplate(event.Name)
						if err != nil {
							engine.logger.Error(err)
						}
					}

					if strings.HasSuffix(event.Name, engine.SqlMap.Extension) {
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

	if engine.SqlTemplate.SqlTemplateRootDir != "" {
		err = engine.watcher.Add(engine.SqlTemplate.SqlTemplateRootDir)
		if err != nil {
			return err
		}
	}

	return nil
}

//stop filesytem watcher
func (engine *Engine) StopFSWatcher() error {
	return engine.watcher.Close()
}
