package xorm

import (
	"database/sql"
	"strings"
)

type SqlMapsExecutor struct {
	session *Session
	sqlkeys interface{}
	parmas  interface{}
	err     error
}

func (sqlMapsExecutor *SqlMapsExecutor) Execute() ([][]map[string]interface{}, map[string][]map[string]interface{}, error) {
	if sqlMapsExecutor.err != nil {
		return nil, nil, sqlMapsExecutor.err
	}

	var model_1_results ResultMap
	var model_2_results sql.Result
	var err error

	sqlModel := 1

	switch sqlMapsExecutor.sqlkeys.(type) {
	case string:
		sqlkey := strings.TrimSpace(sqlMapsExecutor.sqlkeys.(string))
		sqlStr := sqlMapsExecutor.session.Engine.GetSql(sqlkey)
		sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])

		if sqlMapsExecutor.parmas == nil {
			switch sqlCmd {
			case "select", "desc":
				model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkey).Query()
			case "insert", "delete", "update", "create":
				model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkey).Execute()
				sqlModel = 2
			}
		} else {
			switch sqlMapsExecutor.parmas.(type) {
			case []map[string]interface{}:
				parmaMap, ok := sqlMapsExecutor.parmas.([]map[string]interface{})
				if !ok {
					return nil, nil, ErrParamsType
				}
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaMap[0]).Query()

				case "insert", "delete", "update", "create":
					model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaMap[0]).Execute()
					sqlModel = 2
				}

			case map[string]interface{}:
				parmaMap, ok := sqlMapsExecutor.parmas.(map[string]interface{})
				if !ok {
					return nil, nil, ErrParamsType
				}

				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaMap).Query()
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaMap).Execute()
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
		if sqlMapsExecutor.session.IsSqlFuc == true {
			err := sqlMapsExecutor.session.Begin()
			if err != nil {
				return nil, nil, err
			}
		}
		sqlkeysSlice := sqlMapsExecutor.sqlkeys.([]string)
		n := len(sqlkeysSlice)
		resultSlice := make([][]map[string]interface{}, n)
		parmaSlice := make([]map[string]interface{}, n)
		switch sqlMapsExecutor.parmas.(type) {
		case []map[string]interface{}:
			parmaSlice = sqlMapsExecutor.parmas.([]map[string]interface{})

		default:
			if sqlMapsExecutor.session.IsSqlFuc == true {
				err := sqlMapsExecutor.session.Rollback()
				if err != nil {
					return nil, nil, err
				}
			}
			return nil, nil, ErrParamsType
		}

		for i, _ := range sqlkeysSlice {
			sqlkey := strings.TrimSpace(sqlkeysSlice[i])
			sqlStr := sqlMapsExecutor.session.Engine.GetSql(sqlkey)
			sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
			if parmaSlice[i] == nil {
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkey).Query()
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkey).Execute()
					sqlModel = 2
				}
			} else {
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaSlice[i]).Query()
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaSlice[i]).Execute()
					sqlModel = 2
				}
			}

			if sqlModel == 1 {
				if model_1_results.Error != nil {
					if sqlMapsExecutor.session.IsSqlFuc == true {
						err := sqlMapsExecutor.session.Rollback()
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
					if sqlMapsExecutor.session.IsSqlFuc == true {
						err := sqlMapsExecutor.session.Rollback()
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
					return nil, nil, err
				}
				resultMap[0]["RowsAffected"] = RowsAffected
				resultSlice[i] = make([]map[string]interface{}, 1)
				resultSlice[i] = resultMap

			}
		}

		if sqlMapsExecutor.session.IsSqlFuc == true {
			err := sqlMapsExecutor.session.Commit()
			if err != nil {
				return nil, nil, err
			}
		}
		return resultSlice, nil, nil

	case map[string]string:
		if sqlMapsExecutor.session.IsSqlFuc == true {
			err := sqlMapsExecutor.session.Begin()
			if err != nil {
				return nil, nil, err
			}
		}
		sqlkeysMap := sqlMapsExecutor.sqlkeys.(map[string]string)
		n := len(sqlkeysMap)
		resultsMap := make(map[string][]map[string]interface{}, n)
		parmasMap := make(map[string]map[string]interface{}, n)
		switch sqlMapsExecutor.parmas.(type) {
		case map[string]map[string]interface{}:
			parmasMap = sqlMapsExecutor.parmas.(map[string]map[string]interface{})

		default:
			if sqlMapsExecutor.session.IsSqlFuc == true {
				err := sqlMapsExecutor.session.Rollback()
				if err != nil {
					return nil, nil, err
				}
			}
			return nil, nil, ErrParamsType
		}

		for k, _ := range sqlkeysMap {
			sqlkey := strings.TrimSpace(sqlkeysMap[k])
			sqlStr := sqlMapsExecutor.session.Engine.GetSql(sqlkey)
			sqlCmd := strings.ToLower(strings.Split(sqlStr, " ")[0])
			if parmasMap[k] == nil {
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkey).Query()

				case "insert", "delete", "update", "create":
					model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkey).Execute()
					sqlModel = 2
				}
			} else {
				parmaMap := parmasMap[k]
				switch sqlCmd {
				case "select", "desc":
					model_1_results = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaMap).Query()
				case "insert", "delete", "update", "create":
					model_2_results, err = sqlMapsExecutor.session.SqlMapClient(sqlkey, &parmaMap).Execute()
					sqlModel = 2
				}
			}

			if sqlModel == 1 {
				if model_1_results.Error != nil {
					if sqlMapsExecutor.session.IsSqlFuc == true {
						err := sqlMapsExecutor.session.Rollback()
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
					if sqlMapsExecutor.session.IsSqlFuc == true {
						err := sqlMapsExecutor.session.Rollback()
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
					if sqlMapsExecutor.session.IsSqlFuc == true {
						err := sqlMapsExecutor.session.Rollback()
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
		if sqlMapsExecutor.session.IsSqlFuc == true {
			err := sqlMapsExecutor.session.Commit()
			if err != nil {
				return nil, nil, err
			}
		}
		return nil, resultsMap, nil

	}

	return nil, nil, nil
}
