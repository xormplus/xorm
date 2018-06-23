package xorm

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type NullString struct {
	String string
	Valid  bool
}

func (ns NullString) Ptr() *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func (ns NullString) ValueOrZero() string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func (ns NullString) IsNil() bool {
	return !ns.Valid
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		ns.String = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &ns)
	case nil:
		ns.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullString", reflect.TypeOf(v).Name())
	}
	ns.Valid = err == nil
	return err
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns NullString) MarshalText() ([]byte, error) {
	if !ns.Valid {
		return []byte{}, nil
	}
	return []byte(ns.String), nil
}

func (ns *NullString) UnmarshalText(text []byte) error {
	ns.String = string(text)
	ns.Valid = ns.String != ""
	return nil
}
