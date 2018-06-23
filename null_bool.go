package xorm

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type NullBool struct {
	Bool  bool
	Valid bool
}

func (nb NullBool) Ptr() *bool {
	if !nb.Valid {
		return nil
	}
	return &nb.Bool
}

func (nb NullBool) ValueOrZero() bool {
	if !nb.Valid {
		return false
	}
	return nb.Bool
}

func (nb NullBool) IsNil() bool {
	return !nb.Valid
}

func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		nb.Bool = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &nb)
	case nil:
		nb.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullBool", reflect.TypeOf(v).Name())
	}
	nb.Valid = err == nil
	return err
}

func (nb *NullBool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		nb.Valid = false
		return nil
	case "true":
		nb.Bool = true
	case "false":
		nb.Bool = false
	default:
		nb.Valid = false
		return errors.New("invalid input:" + str)
	}
	nb.Valid = true
	return nil
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	if !nb.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

func (nb NullBool) MarshalText() ([]byte, error) {
	if !nb.Valid {
		return []byte{}, nil
	}
	if !nb.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}
