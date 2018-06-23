package xorm

import (
	"database/sql"
	"strings"
	"time"
)

type SqlsExecutor struct {
	session *Session
	sqls    interface{}
	parmas  interface{}
	err     error
}

func (sqlsExecutor *SqlsExecutor) Execute() ([][]map[string]interface{}, map[string][]map[string]interface{}, error) {
	defer sqlsExecutor.session.resetStatement()
	defer sqlsExecutor.session.Close()

	if sqlsExecutor.err != nil {
		return nil, nil, sqlsExecutor.err
	}
	var model_1_results *ResultMap
	var model_2_results sql.Result
	var err error

	sqlModel := 1

	if sqlsExecutor.session.isSqlFunc == true {
		err := sqlsExecutor.session.Begin()
		if err != nil {
			return nil, nil, err
		}
	}

	switch sqlsExecutor.sqls.(type) {
	case string:
		sqlStr := strings.TrimSpace(sqlsExecutor.sqls.(string))
		sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])

		if sqlsExecutor.parmas == nil {
			switch sqlCmd {
			case "select":
				model_1_results = sqlsExecutor.session.Sql(sqlStr).Query()
				sqlModel = 1
			case "insert", "delete", "update", "create", "drop":
				model_2_results, err = sqlsExecutor.session.Sql(sqlStr).Execute()
				sqlModel = 2
			default:
				sqlModel = 3
			}
		} else {
			switch sqlsExecutor.parmas.(type) {
			case []map[string]interface{}:
				parmaMap, _ := sqlsExecutor.parmas.([]map[string]interface{})

				key := NewV4().String() + time.Now().String()
				sqlsExecutor.session.engine.AddSql(key, sqlStr)
				switch sqlCmd {
				case "select":
					model_1_results = sqlsExecutor.session.SqlMapClient(key, &parmaMap[0]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlsExecutor.session.SqlMapClient(key, &parmaMap[0]).Execute()
					sqlModel = 2
				default:
					sqlModel = 3

				}
				sqlsExecutor.session.engine.RemoveSql(key)
			case map[string]interface{}:
				parmaMap, _ := sqlsExecutor.parmas.(map[string]interface{})

				key := NewV4().String() + time.Now().String()
				sqlsExecutor.session.engine.AddSql(key, sqlStr)
				switch sqlCmd {
				case "select":
					model_1_results = sqlsExecutor.session.SqlMapClient(key, &parmaMap).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlsExecutor.session.SqlMapClient(key, &parmaMap).Execute()
					sqlModel = 2
				default:
					sqlModel = 3
				}
				sqlsExecutor.session.engine.RemoveSql(key)
			default:
				if sqlsExecutor.session.isSqlFunc == true {
					err1 := sqlsExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, ErrParamsType
			}
		}
		sqlsExecutor.session.isSqlFunc = true
		resultSlice := make([][]map[string]interface{}, 1)

		if sqlModel == 1 {
			if model_1_results.Error != nil {
				if sqlsExecutor.session.isSqlFunc == true {
					err1 := sqlsExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, model_1_results.Error
			}

			resultSlice[0] = make([]map[string]interface{}, len(model_1_results.Result))
			resultSlice[0] = model_1_results.Result
			if sqlsExecutor.session.isSqlFunc == true {
				err1 := sqlsExecutor.session.Commit()
				if err1 != nil {
					return nil, nil, err1
				}
			}
			return resultSlice, nil, nil
		} else if sqlModel == 2 {
			if err != nil {
				if sqlsExecutor.session.isSqlFunc == true {
					err1 := sqlsExecutor.session.Rollback()
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
				return nil, nil, err
			}
			resultMap[0]["RowsAffected"] = RowsAffected
			resultSlice[0] = resultMap
			if sqlsExecutor.session.isSqlFunc == true {
				err1 := sqlsExecutor.session.Commit()
				if err1 != nil {
					return nil, nil, err1
				}
			}
			return resultSlice, nil, nil
		} else {
			resultSlice[0] = nil
		}
	case []string:

		sqlsSlice := sqlsExecutor.sqls.([]string)
		n := len(sqlsSlice)
		resultSlice := make([][]map[string]interface{}, n)
		parmaSlice := make([]map[string]interface{}, n)

		if sqlsExecutor.parmas == nil {
			for i, _ := range sqlsSlice {
				sqlsExecutor.session.isSqlFunc = true
				sqlStr := strings.TrimSpace(sqlsSlice[i])
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				switch sqlCmd {
				case "select":
					model_1_results = sqlsExecutor.session.Sql(sqlStr).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create", "drop":
					model_2_results, err = sqlsExecutor.session.Sql(sqlStr).Execute()
					sqlModel = 2
				default:
					sqlModel = 3
				}
				sqlsExecutor.session.isSqlFunc = true
				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultSlice[i] = model_1_results.Result

				} else if sqlModel == 2 {
					if err != nil {
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
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
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
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
			switch sqlsExecutor.parmas.(type) {
			case []map[string]interface{}:
				parmaSlice = sqlsExecutor.parmas.([]map[string]interface{})

			default:
				if sqlsExecutor.session.isSqlFunc == true {
					err1 := sqlsExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, ErrParamsType
			}

			for i, _ := range sqlsSlice {
				sqlsExecutor.session.isSqlFunc = true
				sqlStr := strings.TrimSpace(sqlsSlice[i])
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				if parmaSlice[i] == nil {
					switch sqlCmd {
					case "select":
						model_1_results = sqlsExecutor.session.Sql(sqlStr).Query()
						sqlModel = 1
					case "insert", "delete", "update", "create", "drop":
						model_2_results, err = sqlsExecutor.session.Sql(sqlStr).Execute()
						sqlModel = 2
					default:
						sqlModel = 3
					}
				} else {
					key := NewV4().String() + time.Now().String()
					sqlsExecutor.session.engine.AddSql(key, sqlStr)
					switch sqlCmd {
					case "select":
						model_1_results = sqlsExecutor.session.SqlMapClient(key, &parmaSlice[i]).Query()
						sqlModel = 1
					case "insert", "delete", "update", "create", "drop":
						model_2_results, err = sqlsExecutor.session.SqlMapClient(key, &parmaSlice[i]).Execute()
						sqlModel = 2
					default:
						sqlModel = 3
					}
					sqlsExecutor.session.engine.RemoveSql(key)
				}
				sqlsExecutor.session.isSqlFunc = true
				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultSlice[i] = model_1_results.Result

				} else if sqlModel == 2 {
					if err != nil {
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
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
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
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

		if sqlsExecutor.session.isSqlFunc == true {
			err1 := sqlsExecutor.session.Commit()
			if err1 != nil {
				return nil, nil, err1
			}
		}
		return resultSlice, nil, nil

	case map[string]string:

		sqlsMap := sqlsExecutor.sqls.(map[string]string)
		n := len(sqlsMap)
		resultsMap := make(map[string][]map[string]interface{}, n)
		parmasMap := make(map[string]map[string]interface{}, n)

		if sqlsExecutor.parmas == nil {
			for k, _ := range sqlsMap {
				sqlsExecutor.session.isSqlFunc = true
				sqlStr := strings.TrimSpace(sqlsMap[k])
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])

				switch sqlCmd {
				case "select":
					sqlModel = 1
					model_1_results = sqlsExecutor.session.Sql(sqlStr).Query()

				case "insert", "delete", "update", "create", "drop":
					sqlModel = 2
					model_2_results, err = sqlsExecutor.session.Sql(sqlStr).Execute()

				default:
					sqlModel = 3
				}
				sqlsExecutor.session.isSqlFunc = true
				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultsMap[k] = model_1_results.Result

				} else if sqlModel == 2 {
					if err != nil {
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
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

						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
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
			switch sqlsExecutor.parmas.(type) {
			case map[string]map[string]interface{}:
				parmasMap = sqlsExecutor.parmas.(map[string]map[string]interface{})

			default:
				if sqlsExecutor.session.isSqlFunc == true {
					err1 := sqlsExecutor.session.Rollback()
					if err1 != nil {
						return nil, nil, err1
					}
				}
				return nil, nil, ErrParamsType
			}

			for k, _ := range sqlsMap {
				sqlsExecutor.session.isSqlFunc = true
				sqlStr := strings.TrimSpace(sqlsMap[k])
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				if parmasMap[k] == nil {
					switch sqlCmd {
					case "select":
						sqlModel = 1
						model_1_results = sqlsExecutor.session.Sql(sqlStr).Query()

					case "insert", "delete", "update", "create", "drop":
						sqlModel = 2
						model_2_results, err = sqlsExecutor.session.Sql(sqlStr).Execute()

					default:
						sqlModel = 3
					}
				} else {
					key := NewV4().String() + time.Now().String()
					sqlsExecutor.session.engine.AddSql(key, sqlStr)
					parmaMap := parmasMap[k]
					switch sqlCmd {
					case "select":
						sqlModel = 1
						model_1_results = sqlsExecutor.session.SqlMapClient(key, &parmaMap).Query()

					case "insert", "delete", "update", "create", "drop":
						sqlModel = 2
						model_2_results, err = sqlsExecutor.session.SqlMapClient(key, &parmaMap).Execute()

					default:
						sqlModel = 3
					}
					sqlsExecutor.session.engine.RemoveSql(key)
				}
				sqlsExecutor.session.isSqlFunc = true
				if sqlModel == 1 {
					if model_1_results.Error != nil {
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
							if err1 != nil {
								return nil, nil, err1
							}
						}
						return nil, nil, model_1_results.Error
					}

					resultsMap[k] = model_1_results.Result

				} else if sqlModel == 2 {
					if err != nil {
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
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
						if sqlsExecutor.session.isSqlFunc == true {
							err1 := sqlsExecutor.session.Rollback()
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

		if sqlsExecutor.session.isSqlFunc == true {
			err1 := sqlsExecutor.session.Commit()
			if err1 != nil {
				return nil, nil, err1
			}
		}
		return nil, resultsMap, nil

	}

	return nil, nil, nil
}
