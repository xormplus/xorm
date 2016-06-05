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
	if sqlTemplatesExecutor.err != nil {
		return nil, nil, sqlTemplatesExecutor.err
	}

	var model_1_results ResultMap
	var model_2_results sql.Result
	var err error

	sqlModel := 1

	switch sqlTemplatesExecutor.sqlkeys.(type) {
	case string:
		sqlkey := strings.TrimSpace(sqlTemplatesExecutor.sqlkeys.(string))
		if sqlTemplatesExecutor.parmas == nil {
			sqlStr, err := sqlTemplatesExecutor.session.Engine.GetSqlTemplate(sqlkey).Execute(nil)
			if err != nil {
				if sqlTemplatesExecutor.session.IsSqlFuc == true {
					err = sqlTemplatesExecutor.session.Rollback()
					if err != nil {
						return nil, nil, err
					}
				}
				return nil, nil, err
			}
			sqlStr = strings.TrimSpace(sqlStr)

			sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
			switch sqlCmd {
			case "select", "desc":
				model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey).Query()
				sqlModel = 1
			case "insert", "delete", "update", "create":
				model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey).Execute()
				sqlModel = 2
			}
		} else {
			switch sqlTemplatesExecutor.parmas.(type) {
			case []map[string]interface{}:
				parmaMap, ok := sqlTemplatesExecutor.parmas.([]map[string]interface{})
				if !ok {
					return nil, nil, ErrParamsType
				}
				sqlStr, err := sqlTemplatesExecutor.session.Engine.GetSqlTemplate(sqlkey).Execute(parmaMap[0])
				if err != nil {
					if sqlTemplatesExecutor.session.IsSqlFuc == true {
						err = sqlTemplatesExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, err
				}
				sqlStr = strings.TrimSpace(sqlStr)

				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey, &parmaMap[0]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey, &parmaMap[0]).Execute()
					sqlModel = 2
				}

			case map[string]interface{}:
				parmaMap, ok := sqlTemplatesExecutor.parmas.(map[string]interface{})
				if !ok {
					return nil, nil, ErrParamsType
				}
				sqlStr, err := sqlTemplatesExecutor.session.Engine.GetSqlTemplate(sqlkey).Execute(parmaMap)
				if err != nil {
					if sqlTemplatesExecutor.session.IsSqlFuc == true {
						err = sqlTemplatesExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, err
				}
				sqlStr = strings.TrimSpace(sqlStr)
				sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey, &parmaMap).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkey, &parmaMap).Execute()
					sqlModel = 2
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

			resultSlice[0] = make([]map[string]interface{}, len(model_1_results.Results))
			resultSlice[0] = model_1_results.Results
			return resultSlice, nil, nil
		} else {
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
			return resultSlice, nil, nil
		}
	case []string:
		if sqlTemplatesExecutor.session.IsSqlFuc == true {
			err := sqlTemplatesExecutor.session.Begin()
			if err != nil {
				return nil, nil, err
			}
		}
		sqlkeysSlice := sqlTemplatesExecutor.sqlkeys.([]string)
		n := len(sqlkeysSlice)
		resultSlice := make([][]map[string]interface{}, n)
		parmaSlice := make([]map[string]interface{}, n)
		switch sqlTemplatesExecutor.parmas.(type) {
		case []map[string]interface{}:
			parmaSlice = sqlTemplatesExecutor.parmas.([]map[string]interface{})

		default:
			if sqlTemplatesExecutor.session.IsSqlFuc == true {
				err = sqlTemplatesExecutor.session.Rollback()
				if err != nil {
					return nil, nil, err
				}
			}
			return nil, nil, ErrParamsType
		}

		for i, _ := range sqlkeysSlice {
			sqlStr, err := sqlTemplatesExecutor.session.Engine.GetSqlTemplate(sqlkeysSlice[i]).Execute(parmaSlice[i])
			if err != nil {
				if sqlTemplatesExecutor.session.IsSqlFuc == true {
					err = sqlTemplatesExecutor.session.Rollback()
					if err != nil {
						return nil, nil, err
					}
				}
				return nil, nil, err
			}
			sqlStr = strings.TrimSpace(sqlStr)

			sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
			if parmaSlice[i] == nil {
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysSlice[i]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysSlice[i]).Execute()
					sqlModel = 2
				}
			} else {
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysSlice[i], &parmaSlice[i]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysSlice[i], &parmaSlice[i]).Execute()
					sqlModel = 2
				}
			}

			if sqlModel == 1 {
				if model_1_results.Error != nil {
					if sqlTemplatesExecutor.session.IsSqlFuc == true {
						err = sqlTemplatesExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, model_1_results.Error
				}

				resultSlice[i] = make([]map[string]interface{}, len(model_1_results.Results))
				resultSlice[i] = model_1_results.Results

			} else {
				if err != nil {
					if sqlTemplatesExecutor.session.IsSqlFuc == true {
						err = sqlTemplatesExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
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
					if sqlTemplatesExecutor.session.IsSqlFuc == true {
						err = sqlTemplatesExecutor.session.Commit()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, err
				}
				resultMap[0]["RowsAffected"] = RowsAffected
				resultSlice[i] = make([]map[string]interface{}, 1)
				resultSlice[i] = resultMap

			}
		}

		if sqlTemplatesExecutor.session.IsSqlFuc == true {
			err = sqlTemplatesExecutor.session.Commit()
			if err != nil {
				return nil, nil, err
			}
		}
		return resultSlice, nil, nil

	case map[string]string:
		if sqlTemplatesExecutor.session.IsSqlFuc == true {
			err = sqlTemplatesExecutor.session.Begin()
			if err != nil {
				return nil, nil, err
			}
		}
		sqlkeysMap := sqlTemplatesExecutor.sqlkeys.(map[string]string)
		n := len(sqlkeysMap)
		resultsMap := make(map[string][]map[string]interface{}, n)
		parmasMap := make(map[string]map[string]interface{}, n)

		switch sqlTemplatesExecutor.parmas.(type) {
		case map[string]map[string]interface{}:
			parmasMap = sqlTemplatesExecutor.parmas.(map[string]map[string]interface{})

		default:
			if sqlTemplatesExecutor.session.IsSqlFuc == true {
				err = sqlTemplatesExecutor.session.Rollback()
				if err != nil {
					return nil, nil, err
				}
			}
			return nil, nil, ErrParamsType
		}

		for k, _ := range sqlkeysMap {
			sqlStr, err := sqlTemplatesExecutor.session.Engine.GetSqlTemplate(sqlkeysMap[k]).Execute(parmasMap[k])
			if err != nil {
				if sqlTemplatesExecutor.session.IsSqlFuc == true {
					err = sqlTemplatesExecutor.session.Rollback()
					if err != nil {
						return nil, nil, err
					}
				}
				return nil, nil, err
			}
			sqlStr = strings.TrimSpace(sqlStr)
			sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
			if parmasMap[k] == nil {
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysMap[k]).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysMap[k]).Execute()
					sqlModel = 2
				}
			} else {
				parmaMap := parmasMap[k]
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysMap[k], &parmaMap).Query()
					sqlModel = 1
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlTemplatesExecutor.session.SqlTemplateClient(sqlkeysMap[k], &parmaMap).Execute()
					sqlModel = 2
				}
			}

			if sqlModel == 1 {
				if model_1_results.Error != nil {
					if sqlTemplatesExecutor.session.IsSqlFuc == true {
						err = sqlTemplatesExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, model_1_results.Error
				}

				resultsMap[k] = make([]map[string]interface{}, len(model_1_results.Results))
				resultsMap[k] = model_1_results.Results

			} else {
				if err != nil {
					if sqlTemplatesExecutor.session.IsSqlFuc == true {
						err = sqlTemplatesExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
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
					if sqlTemplatesExecutor.session.IsSqlFuc == true {
						err = sqlTemplatesExecutor.session.Rollback()
						if err != nil {
							return nil, nil, err
						}
					}
					return nil, nil, err
				}
				resultMap[0]["RowsAffected"] = RowsAffected
				resultsMap[k] = make([]map[string]interface{}, 1)
				resultsMap[k] = resultMap

			}
		}
		if sqlTemplatesExecutor.session.IsSqlFuc == true {
			err = sqlTemplatesExecutor.session.Commit()
			if err != nil {
				return nil, nil, err
			}
		}
		return nil, resultsMap, nil

	}

	return nil, nil, nil
}
