package xorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type NullUint16 struct {
	Uint16 uint16
	Valid  bool
}

func (nu NullUint16) Ptr() *uint16 {
	if !nu.Valid {
		return nil
	}
	return &nu.Uint16
}

func (nu NullUint16) ValueOrZero() uint16 {
	if !nu.Valid {
		return 0
	}
	return nu.Uint16
}

func (nu NullUint16) IsNil() bool {
	return !nu.Valid
}

func (nu *NullUint16) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &nu.Uint16)
	case string:
		str := string(x)
		if len(str) == 0 {
			nu.Valid = false
			return nil
		}
		var u uint64
		u, err = strconv.ParseUint(str, 10, 16)
		if err == nil {
			nu.Uint16 = uint16(u)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &nu)
	case nil:
		nu.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullUint16", reflect.TypeOf(v).Name())
	}
	nu.Valid = err == nil
	return err
}

func (nu *NullUint16) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		nu.Valid = false
		return nil
	}
	var err error
	var u uint64
	u, err = strconv.ParseUint(string(text), 10, 16)
	if err == nil {
		nu.Uint16 = uint16(u)
	}
	nu.Valid = err == nil
	return err
}

func (nu NullUint16) MarshalJSON() ([]byte, error) {
	if !nu.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(nu.Uint16), 10)), nil
}

func (nu NullUint16) MarshalText() ([]byte, error) {
	if !nu.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(nu.Uint16), 10)), nil
}
