package xorm

import (
	"strings"
	"time"
)

type SqlExecutor struct {
	session *Session
	sqls    interface{}
	parmas  interface{}
	err     error
}

func (sqlExecutor *SqlExecutor) Execute() ([][]map[string]interface{}, map[string][]map[string]interface{}, error) {
	if sqlExecutor.err != nil {
		return nil, nil, sqlExecutor.err
	}

	switch sqlExecutor.sqls.(type) {
	case string:
		sqlstr := strings.TrimLeft(sqlExecutor.sqls.(string), " \n")
		sqlCmd := strings.ToLower(strings.Split(sqlstr, " ")[0])

		if sqlExecutor.parmas == nil {
			switch sqlCmd {
			case "select", "desc":
				rsults := sqlExecutor.session.Sql(sqlstr).Query()
				if rsults.Error != nil {
					return nil, nil, rsults.Error
				}
				resultSlice := make([][]map[string]interface{}, 1)
				resultSlice[0] = rsults.Results
				return resultSlice, nil, nil
			case "insert", "delete", "update":
				rsults, err := sqlExecutor.session.Sql(sqlstr).Execute()
				if err != nil {
					return nil, nil, err
				}

				resultSlice := make([][]map[string]interface{}, 1)
				resultMap := make([]map[string]interface{}, 1)
				resultMap[0] = make(map[string]interface{})

				//todo all database support LastInsertId
				LastInsertId, _ := rsults.LastInsertId()

				resultMap[0]["LastInsertId"] = LastInsertId
				RowsAffected, err := rsults.RowsAffected()
				if err != nil {
					return nil, nil, err
				}
				resultMap[0]["RowsAffected"] = RowsAffected
				resultSlice[0] = resultMap
				return resultSlice, nil, nil
			}
		} else {
			switch sqlExecutor.parmas.(type) {
			case []map[string]interface{}:
				parmaMap, ok := sqlExecutor.parmas.([]map[string]interface{})
				if !ok {
					return nil, nil, ErrParamsType
				}
				switch sqlCmd {
				case "select", "desc":
					rsults := sqlExecutor.session.Sql(sqlstr, &parmaMap[0]).Query()
					if rsults.Error != nil {
						return nil, nil, rsults.Error
					}
					resultSlice := make([][]map[string]interface{}, 1)
					resultSlice[0] = rsults.Results
					return resultSlice, nil, nil

				case "insert", "delete", "update":
					rsults, err := sqlExecutor.session.Sql(sqlstr, &parmaMap[0]).Execute()
					if err != nil {
						return nil, nil, err
					}

					resultSlice := make([][]map[string]interface{}, 1)
					resultMap := make([]map[string]interface{}, 1)
					resultMap[0] = make(map[string]interface{})
					LastInsertId, _ := rsults.LastInsertId()

					resultMap[0]["LastInsertId"] = LastInsertId
					RowsAffected, err := rsults.RowsAffected()
					if err != nil {
						return nil, nil, err
					}
					resultMap[0]["RowsAffected"] = RowsAffected
					resultSlice[0] = resultMap
					return resultSlice, nil, nil
				}
			case map[string]interface{}:
				parmaMap, ok := sqlExecutor.parmas.(map[string]interface{})
				if !ok {
					return nil, nil, ErrParamsType
				}
				switch sqlCmd {
				case "select", "desc":
					key := NewV4().String() + time.Now().String()
					sqlExecutor.session.Engine.AddSql(key, sqlstr)

					rsults := sqlExecutor.session.SqlMapClient(key, &parmaMap).Query()
					sqlExecutor.session.Engine.RemoveSql(key)

					if rsults.Error != nil {
						return nil, nil, rsults.Error
					}

					resultSlice := make([][]map[string]interface{}, 1)
					resultSlice[0] = rsults.Results
					return resultSlice, nil, nil
				case "insert", "delete", "update":
					rsults, err := sqlExecutor.session.Sql(sqlstr, &parmaMap).Execute()
					if err != nil {
						return nil, nil, err
					}

					resultSlice := make([][]map[string]interface{}, 1)
					resultMap := make([]map[string]interface{}, 1)
					resultMap[0] = make(map[string]interface{})
					LastInsertId, _ := rsults.LastInsertId()

					resultMap[0]["LastInsertId"] = LastInsertId
					RowsAffected, err := rsults.RowsAffected()
					if err != nil {
						return nil, nil, err
					}
					resultMap[0]["RowsAffected"] = RowsAffected
					resultSlice[0] = resultMap
					return resultSlice, nil, nil
				}
			default:
				return nil, nil, ErrParamsType
			}
		}
	case []string:
		if sqlExecutor.session.IsSqlFuc == true {
			err := sqlExecutor.session.Begin()
			if err != nil {
				return nil, nil, err
			}
		}

	case map[string]string:
		if sqlExecutor.session.IsSqlFuc == true {
			err := sqlExecutor.session.Begin()
			if err != nil {
				return nil, nil, err
			}
		}

	}

	return nil, nil, nil
}

//func (sqlExecutor *SqlExecutor) Execute() (interface{}, error) {

//	switch sqlExecutor.sqls.(type) {
//	case string:
//		sqlstr := strings.TrimLeft(sqlExecutor.sqls.(string), " \n")
//		sqlCmd := strings.ToLower(strings.Split(sqlstr, " ")[0])
//		switch sqlExecutor.parmasCount {
//		case 0:
//			switch sqlCmd {
//			case "select", "desc":
//				rsults := sqlExecutor.session.Sql(sqlstr).Query()
//				if rsults.Error != nil {
//					return nil, rsults.Error
//				}
//				return rsults.Results, nil
//			case "insert", "delete", "update":
//				return sqlExecutor.session.Sql(sqlstr).Execute()
//			}
//		case 1:
//			switch sqlExecutor.parmas.(type) {
//			case map[string]interface{}:
//				switch sqlCmd {
//				case "select", "desc":

//				case "insert", "delete", "update":

//				}
//			default:
//				switch sqlCmd {
//				case "select", "desc":

//				case "insert", "delete", "update":

//				}
//			}
//		default:
//			switch sqlCmd {
//			case "select", "desc":

//			case "insert", "delete", "update":

//			}
//		}
//	case []string:
//		if sqlExecutor.session.IsSqlFuc == true {
//			err := sqlExecutor.session.Begin()
//			if err != nil {
//				return nil, err
//			}
//		}

//	case map[string]string:
//		if sqlExecutor.session.IsSqlFuc == true {
//			err := sqlExecutor.session.Begin()
//			if err != nil {
//				return nil, err
//			}
//		}

//	}

//	return 1, nil
//}
