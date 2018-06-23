package xorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type NullInt struct {
	Int   int
	Valid bool
}

func (ni NullInt) Ptr() *int {
	if !ni.Valid {
		return nil
	}
	return &ni.Int
}

func (ni NullInt) ValueOrZero() int {
	if !ni.Valid {
		return 0
	}
	return ni.Int
}

func (ni NullInt) IsNil() bool {
	return !ni.Valid
}

func (ni *NullInt) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &ni.Int)
	case string:
		str := string(x)
		if len(str) == 0 {
			ni.Valid = false
			return nil
		}
		ni.Int, err = strconv.Atoi(string(str))
	case map[string]interface{}:
		err = json.Unmarshal(data, &ni)
	case nil:
		ni.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullInt", reflect.TypeOf(v).Name())
	}
	ni.Valid = err == nil
	return err
}

func (ni *NullInt) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		ni.Valid = false
		return nil
	}
	var err error
	ni.Int, err = strconv.Atoi(string(text))
	ni.Valid = err == nil
	return err
}

func (ni NullInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.Itoa(ni.Int)), nil
}

func (ni NullInt) MarshalText() ([]byte, error) {
	if !ni.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.Itoa(ni.Int)), nil
}
