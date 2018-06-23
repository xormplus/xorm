package xorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type NullInt8 struct {
	Int8  int8
	Valid bool
}

func (ni NullInt8) Ptr() *int8 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int8
}

func (ni NullInt8) ValueOrZero() int8 {
	if !ni.Valid {
		return 0
	}
	return ni.Int8
}

func (ni NullInt8) IsNil() bool {
	return !ni.Valid
}

func (ni *NullInt8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &ni.Int8)
	case string:
		str := string(x)
		if len(str) == 0 {
			ni.Valid = false
			return nil
		}
		var i int64
		i, err = strconv.ParseInt(str, 10, 8)
		if err == nil {
			ni.Int8 = int8(i)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &ni)
	case nil:
		ni.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullInt8", reflect.TypeOf(v).Name())
	}
	ni.Valid = err == nil
	return err
}

func (ni *NullInt8) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		ni.Valid = false
		return nil
	}
	var err error
	var i int64
	i, err = strconv.ParseInt(string(text), 10, 8)
	if err == nil {
		ni.Int8 = int8(i)
	}
	ni.Valid = err == nil
	return err
}

func (ni NullInt8) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(ni.Int8), 10)), nil
}

func (ni NullInt8) MarshalText() ([]byte, error) {
	if !ni.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(ni.Int8), 10)), nil
}
