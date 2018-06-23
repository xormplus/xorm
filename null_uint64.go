package xorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

func (nu NullUint64) Ptr() *uint64 {
	if !nu.Valid {
		return nil
	}
	return &nu.Uint64
}

func (nu NullUint64) ValueOrZero() uint64 {
	if !nu.Valid {
		return 0
	}
	return nu.Uint64
}

func (nu NullUint64) IsNil() bool {
	return !nu.Valid
}

func (nu *NullUint64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &nu.Uint64)
	case string:
		str := string(x)
		if len(str) == 0 {
			nu.Valid = false
			return nil
		}
		nu.Uint64, err = strconv.ParseUint(str, 10, 64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &nu)
	case nil:
		nu.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullUint64", reflect.TypeOf(v).Name())
	}
	nu.Valid = err == nil
	return err
}

func (nu *NullUint64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		nu.Valid = false
		return nil
	}
	var err error
	nu.Uint64, err = strconv.ParseUint(string(text), 10, 64)
	nu.Valid = err == nil
	return err
}

func (nu NullUint64) MarshalJSON() ([]byte, error) {
	if !nu.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(nu.Uint64, 10)), nil
}

func (nu NullUint64) MarshalText() ([]byte, error) {
	if !nu.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(nu.Uint64, 10)), nil
}
