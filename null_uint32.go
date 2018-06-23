package xorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type NullUint32 struct {
	Uint32 uint32
	Valid  bool
}

func (nu NullUint32) Ptr() *uint32 {
	if !nu.Valid {
		return nil
	}
	return &nu.Uint32
}

func (nu NullUint32) ValueOrZero() uint32 {
	if !nu.Valid {
		return 0
	}
	return nu.Uint32
}

func (nu NullUint32) IsNil() bool {
	return !nu.Valid
}

func (nu *NullUint32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &nu.Uint32)
	case string:
		str := string(x)
		if len(str) == 0 {
			nu.Valid = false
			return nil
		}
		var u uint64
		u, err = strconv.ParseUint(str, 10, 32)
		if err == nil {
			nu.Uint32 = uint32(u)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &nu)
	case nil:
		nu.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullUint32", reflect.TypeOf(v).Name())
	}
	nu.Valid = err == nil
	return err
}

func (nu *NullUint32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		nu.Valid = false
		return nil
	}
	var err error
	var u uint64
	u, err = strconv.ParseUint(string(text), 10, 32)
	if err == nil {
		nu.Uint32 = uint32(u)
	}
	nu.Valid = err == nil
	return err
}

func (nu NullUint32) MarshalJSON() ([]byte, error) {
	if !nu.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(nu.Uint32), 10)), nil
}

func (nu NullUint32) MarshalText() ([]byte, error) {
	if !nu.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(nu.Uint32), 10)), nil
}
