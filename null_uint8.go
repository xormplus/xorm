package xorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type NullUint8 struct {
	Uint8 uint8
	Valid bool
}

func (nu NullUint8) Ptr() *uint8 {
	if !nu.Valid {
		return nil
	}
	return &nu.Uint8
}

func (nu NullUint8) ValueOrZero() uint8 {
	if !nu.Valid {
		return 0
	}
	return nu.Uint8
}

func (nu NullUint8) IsNil() bool {
	return !nu.Valid
}

func (nu *NullUint8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		err = json.Unmarshal(data, &nu.Uint8)
	case string:
		str := string(x)
		if len(str) == 0 {
			nu.Valid = false
			return nil
		}
		var u uint64
		u, err = strconv.ParseUint(str, 10, 8)
		if err == nil {
			nu.Uint8 = uint8(u)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &nu)
	case nil:
		nu.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullUint8", reflect.TypeOf(v).Name())
	}
	nu.Valid = err == nil
	return err
}

func (nu *NullUint8) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		nu.Valid = false
		return nil
	}
	var err error
	var u uint64
	u, err = strconv.ParseUint(string(text), 10, 8)
	if err == nil {
		nu.Uint8 = uint8(u)
	}
	nu.Valid = err == nil
	return err
}

func (nu NullUint8) MarshalJSON() ([]byte, error) {
	if !nu.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(nu.Uint8), 10)), nil
}

func (nu NullUint8) MarshalText() ([]byte, error) {
	if !nu.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(nu.Uint8), 10)), nil
}
