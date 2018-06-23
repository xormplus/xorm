package xorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type NullUint struct {
	Uint  uint
	Valid bool
}

func (nu NullUint) Ptr() *uint {
	if !nu.Valid {
		return nil
	}
	return &nu.Uint
}

func (nu NullUint) ValueOrZero() uint {
	if !nu.Valid {
		return 0
	}
	return nu.Uint
}

func (nu NullUint) IsNil() bool {
	return !nu.Valid
}

func (nu *NullUint) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &nu.Uint)
	case string:
		str := string(x)
		if len(str) == 0 {
			nu.Valid = false
			return nil
		}
		var u uint64
		u, err = strconv.ParseUint(str, 10, 0)
		if err == nil {
			nu.Uint = uint(u)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &nu)
	case nil:
		nu.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullUint", reflect.TypeOf(v).Name())
	}
	nu.Valid = err == nil
	return err
}

func (nu *NullUint) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		nu.Valid = false
		return nil
	}
	var err error
	var u uint64
	u, err = strconv.ParseUint(string(text), 10, 0)
	if err == nil {
		nu.Uint = uint(u)
	}
	nu.Valid = err == nil
	return err
}

func (ni NullUint) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(ni.Uint), 10)), nil
}

func (ni NullUint) MarshalText() ([]byte, error) {
	if !ni.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(ni.Uint), 10)), nil
}
