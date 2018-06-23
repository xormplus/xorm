package xorm

import (
	"database/sql"
	"strings"
)

type SqlTemplatesExecutor struct {
	session *Session
	sqlkeys interface{}
	parmas  interface{}
	err     error
}

func (sqlTemplatesExecutor *SqlTemplatesExecutor) Execute() ([][]map[string]interface{}, map[string][]map[string]interface{}, error) {
	defer sqlTemplatesExecutor.session.resetStatement()
	defer sqlTemplatesExecutor.session.Close()

	if sqlTemplatesExecutor.err != nil {
		return nil, nil, sqlTemplatesExecutor.err
	}

	var model_1_results *ResultMap
	var model_2_results sql.Result
	var err error
	var sqlStr string

	sqlModel := 1

	if sqlTemplatesExecutor.session.isSqlFunc == true {
		err := sqlTemplatesExecutor.session.Begin()
		if err != nil {
			return nil, nil, err
		}
	}

	switch sqlTemplatesExecutor.sqlkeys.(type) {
	case string:
		sqlkey := strings.TrimSpace(sqlTemplatesExecutor.sqlkeys.(string))
		if sqlTemplatesExecutor.parmas == nil {
			sqlStr, err = sqlTemplatesExecutor.session.engine.SqlTemplate.Execute(sqlkey)
			if err != nil {
				if sqlTemplatesExecutor.session.isSqlFunc == true {
					err1 := sqlTemplatesExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, err
			}
			sqlStr = strings.TrimSpace(sqlStr)

			sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
			switch sqlCmd {
			case "select":
				model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey).Query()
				sqlModel = 1
			case "insert", "delete", "update", "create", "drop":
				model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey).Execute()
				sqlModel = 2
			default:
				sqlModel = 3
			}
		} else {
			switch sqlTemplatesExecutor.parmas.(type) {
			case []map[string]interface{}:
				parmaMap, _ := sqlTemplatesExecutor.parmas.([]map[string]interface{})

				sqlStr, err = sqlTemplatesExecutor.session.engine.SqlTemplate.Execute(sqlkey, parmaMap[0])
				if err != nil {
					if sqlTemplatesExecutor.session.isSqlFunc == true {
						err1 := sqlTemplatesExecutor.session.Rollback()
						if err1 != nil {
							return nil, nil, err1
						}
					}
					return nil, nil, err
				}
				sqlStr = strings.TrimSpace(sqlStr)

				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				switch sqlCmd {
				case "select":
					model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey, &parmaMap[0]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey, &parmaMap[0]).Execute()
					sqlModel = 2
				default:
					sqlModel = 3
				}

			case map[string]interface{}:
				parmaMap, _ := sqlTemplatesExecutor.parmas.(map[string]interface{})

				sqlStr, err = sqlTemplatesExecutor.session.engine.SqlTemplate.Execute(sqlkey, parmaMap)
				if err != nil {
					if sqlTemplatesExecutor.session.isSqlFunc == true {
						err1 := sqlTemplatesExecutor.session.Rollback()
						if err1 != nil {
							return nil, nil, err1
						}
					}
					return nil, nil, err
				}
				sqlStr = strings.TrimSpace(sqlStr)
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				switch sqlCmd {
				case "select":
					model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey, &parmaMap).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey, &parmaMap).Execute()
					sqlModel = 2
				default:
					sqlModel = 3
				}

			default:
				return nil, nil, ErrParamsType
			}
		}

		resultSlice := make([][]map[string]interface{}, 1)

		if sqlModel == 1 {
			if model_1_results.Error != nil {
				return nil, nil, model_1_results.Error
			}

			resultSlice[0] = model_1_results.Result
			if sqlTemplatesExecutor.session.isSqlFunc == true {
				err1 := sqlTemplatesExecutor.session.Commit()
				if err1 != nil {
					return nil, nil, err1
				}
			}
			return resultSlice, nil, nil
		} else if sqlModel == 2 {
			if err != nil {
				return nil, nil, err
			}

			resultMap := make([]map[string]interface{}, 1)
			resultMap[0] = make(map[string]interface{})

			//todo all database support LastInsertId
			LastInsertId, _ := model_2_results.LastInsertId()

			resultMap[0]["LastInsertId"] = LastInsertId
			RowsAffected, err := model_2_results.RowsAffected()
			if err != nil {
				return nil, nil, err
			}
			resultMap[0]["RowsAffected"] = RowsAffected
			resultSlice[0] = resultMap
			if sqlTemplatesExecutor.session.isSqlFunc == true {
				err1 := sqlTemplatesExecutor.session.Commit()
				if err1 != nil {
					return nil, nil, err1
				}
			}
			return resultSlice, nil, nil
		} else {
			resultSlice[0] = nil
		}
	case []string:

		sqlkeysSlice := sqlTemplatesExecutor.sqlkeys.([]string)
		n := len(sqlkeysSlice)
		resultSlice := make([][]map[string]interface{}, n)
		parmaSlice := make([]map[string]interface{}, n)

		if sqlTemplatesExecutor.parmas == nil {
			for i, _ := range sqlkeysSlice {
				sqlTemplatesExecutor.session.isSqlFunc = true
				sqlStr, err := sqlTemplatesExecutor.session.engine.SqlTemplate.Execute(sqlkeysSlice[i])
				if err != nil {
					if sqlTemplatesExecutor.session.isSqlFunc == true {
						err1 := sqlTemplatesExecutor.session.Rollback()
						if err1 != nil {
							return nil, nil, err1
						}
					}
					return nil, nil, err
				}
				sqlStr = strings.TrimSpace(sqlStr)

				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				switch sqlCmd {
				case "select":
					model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysSlice[i]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysSlice[i]).Execute()
					sqlModel = 2
				default:
					sqlModel = 3
				}
				sqlTemplatesExecutor.session.isSqlFunc = true
				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultSlice[i] = model_1_results.Result

				} else if sqlModel == 2 {
					if err != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}

					resultMap := make([]map[string]interface{}, 1)
					resultMap[0] = make(map[string]interface{})

					//todo all database support LastInsertId
					LastInsertId, _ := model_2_results.LastInsertId()

					resultMap[0]["LastInsertId"] = LastInsertId
					RowsAffected, err := model_2_results.RowsAffected()
					if err != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}
					resultMap[0]["RowsAffected"] = RowsAffected
					resultSlice[i] = make([]map[string]interface{}, 1)
					resultSlice[i] = resultMap

				} else {
					resultSlice[i] = nil
				}

			}

		} else {

			switch sqlTemplatesExecutor.parmas.(type) {
			case []map[string]interface{}:
				parmaSlice = sqlTemplatesExecutor.parmas.([]map[string]interface{})

			default:
				if sqlTemplatesExecutor.session.isSqlFunc == true {
					err1 := sqlTemplatesExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, ErrParamsType
			}

			for i, _ := range sqlkeysSlice {
				sqlTemplatesExecutor.session.isSqlFunc = true
				sqlStr, err := sqlTemplatesExecutor.session.engine.SqlTemplate.Execute(sqlkeysSlice[i], parmaSlice[i])
				if err != nil {
					if sqlTemplatesExecutor.session.isSqlFunc == true {
						err1 := sqlTemplatesExecutor.session.Rollback()
						if err1 != nil {
							return nil, nil, err1
						}
					}
					return nil, nil, err
				}
				sqlStr = strings.TrimSpace(sqlStr)

				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				if parmaSlice[i] == nil {
					switch sqlCmd {
					case "select":
						model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysSlice[i]).Query()
						sqlModel = 1
					case "insert", "delete", "update", "create", "drop":
						model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysSlice[i]).Execute()
						sqlModel = 2
					default:
						sqlModel = 3
					}
				} else {
					switch sqlCmd {
					case "select":
						model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysSlice[i], &parmaSlice[i]).Query()
						sqlModel = 1
					case "insert", "delete", "update", "create", "drop":
						model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysSlice[i], &parmaSlice[i]).Execute()
						sqlModel = 2
					default:
						sqlModel = 3
					}
				}
				sqlTemplatesExecutor.session.isSqlFunc = true
				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultSlice[i] = model_1_results.Result

				} else if sqlModel == 2 {
					if err != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}

					resultMap := make([]map[string]interface{}, 1)
					resultMap[0] = make(map[string]interface{})

					//todo all database support LastInsertId
					LastInsertId, _ := model_2_results.LastInsertId()

					resultMap[0]["LastInsertId"] = LastInsertId
					RowsAffected, err := model_2_results.RowsAffected()
					if err != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}
					resultMap[0]["RowsAffected"] = RowsAffected
					resultSlice[i] = make([]map[string]interface{}, 1)
					resultSlice[i] = resultMap

				} else {
					resultSlice[i] = nil
				}
			}

		}

		if sqlTemplatesExecutor.session.isSqlFunc == true {
			err1 := sqlTemplatesExecutor.session.Commit()
			if err1 != nil {
				return nil, nil, err1
			}
		}
		return resultSlice, nil, nil

	case map[string]string:

		sqlkeysMap := sqlTemplatesExecutor.sqlkeys.(map[string]string)
		n := len(sqlkeysMap)
		resultsMap := make(map[string][]map[string]interface{}, n)
		parmasMap := make(map[string]map[string]interface{}, n)

		if sqlTemplatesExecutor.parmas == nil {

			for k, _ := range sqlkeysMap {
				sqlTemplatesExecutor.session.isSqlFunc = true
				sqlStr, err := sqlTemplatesExecutor.session.engine.SqlTemplate.Execute(sqlkeysMap[k])
				if err != nil {
					if sqlTemplatesExecutor.session.isSqlFunc == true {
						err1 := sqlTemplatesExecutor.session.Rollback()
						if err1 != nil {
							return nil, nil, err1
						}
					}
					return nil, nil, err
				}
				sqlStr = strings.TrimSpace(sqlStr)

				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				switch sqlCmd {
				case "select":
					model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysMap[k]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysMap[k]).Execute()
					sqlModel = 2
				default:
					sqlModel = 3
				}
				sqlTemplatesExecutor.session.isSqlFunc = true
				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultsMap[k] = model_1_results.Result

				} else if sqlModel == 2 {
					if err != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}

					resultMap := make([]map[string]interface{}, 1)
					resultMap[0] = make(map[string]interface{})

					//todo all database support LastInsertId
					LastInsertId, _ := model_2_results.LastInsertId()

					resultMap[0]["LastInsertId"] = LastInsertId
					RowsAffected, err := model_2_results.RowsAffected()
					if err != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}
					resultMap[0]["RowsAffected"] = RowsAffected
					resultsMap[k] = make([]map[string]interface{}, 1)
					resultsMap[k] = resultMap

				} else {
					resultsMap[k] = nil
				}

			}

		} else {

			switch sqlTemplatesExecutor.parmas.(type) {
			case map[string]map[string]interface{}:
				parmasMap = sqlTemplatesExecutor.parmas.(map[string]map[string]interface{})

			default:
				if sqlTemplatesExecutor.session.isSqlFunc == true {
					err1 := sqlTemplatesExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, ErrParamsType
			}

			for k, _ := range sqlkeysMap {
				sqlTemplatesExecutor.session.isSqlFunc = true
				sqlStr, err := sqlTemplatesExecutor.session.engine.SqlTemplate.Execute(sqlkeysMap[k], parmasMap[k])
				if err != nil {
					if sqlTemplatesExecutor.session.isSqlFunc == true {
						err1 := sqlTemplatesExecutor.session.Rollback()
						if err1 != nil {
							return nil, nil, err1
						}
					}
					return nil, nil, err
				}
				sqlStr = strings.TrimSpace(sqlStr)
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				if parmasMap[k] == nil {
					switch sqlCmd {
					case "select":
						model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysMap[k]).Query()
						sqlModel = 1
					case "insert", "delete", "update", "create", "drop":
						model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysMap[k]).Execute()
						sqlModel = 2
					default:
						sqlModel = 3
					}
				} else {
					parmaMap := parmasMap[k]
					switch sqlCmd {
					case "select":
						model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysMap[k], &parmaMap).Query()
						sqlModel = 1
					case "insert", "delete", "update", "create", "drop":
						model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysMap[k], &parmaMap).Execute()
						sqlModel = 2
					default:
						sqlModel = 3
					}
				}
				sqlTemplatesExecutor.session.isSqlFunc = true
				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultsMap[k] = model_1_results.Result

				} else if sqlModel == 2 {
					if err != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}

					resultMap := make([]map[string]interface{}, 1)
					resultMap[0] = make(map[string]interface{})

					//todo all database support LastInsertId
					LastInsertId, _ := model_2_results.LastInsertId()

					resultMap[0]["LastInsertId"] = LastInsertId
					RowsAffected, err := model_2_results.RowsAffected()
					if err != nil {
						if sqlTemplatesExecutor.session.isSqlFunc == true {
							err1 := sqlTemplatesExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, err
					}
					resultMap[0]["RowsAffected"] = RowsAffected
					resultsMap[k] = make([]map[string]interface{}, 1)
					resultsMap[k] = resultMap

				} else {
					resultsMap[k] = nil
				}
			}
		}

		if sqlTemplatesExecutor.session.isSqlFunc == true {
			err1 := sqlTemplatesExecutor.session.Commit()
			if err1 != nil {
				return nil, nil, err1
			}
		}
		return nil, resultsMap, nil

	}

	return nil, nil, nil
}
