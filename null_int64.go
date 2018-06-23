package xorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type NullInt64 struct {
	Int64 int64
	Valid bool
}

func (ni NullInt64) Ptr() *int64 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int64
}

func (ni NullInt64) ValueOrZero() int64 {
	if !ni.Valid {
		return 0
	}
	return ni.Int64
}

func (ni NullInt64) IsNil() bool {
	return !ni.Valid
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &ni.Int64)
	case string:
		str := string(x)
		if len(str) == 0 {
			ni.Valid = false
			return nil
		}
		ni.Int64, err = strconv.ParseInt(str, 10, 64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &ni)
	case nil:
		ni.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullInt64", reflect.TypeOf(v).Name())
	}
	ni.Valid = err == nil
	return err
}

func (ni *NullInt64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		ni.Valid = false
		return nil
	}
	var err error
	ni.Int64, err = strconv.ParseInt(string(text), 10, 64)
	ni.Valid = err == nil
	return err
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(ni.Int64, 10)), nil
}

func (ni NullInt64) MarshalText() ([]byte, error) {
	if !ni.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(ni.Int64, 10)), nil
}
