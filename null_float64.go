package xorm

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type NullFloat64 struct {
	Float64 float64
	Valid   bool
}

func (nf NullFloat64) Ptr() *float64 {
	if !nf.Valid {
		return nil
	}
	return &nf.Float64
}

func (nf NullFloat64) ValueOrZero() float64 {
	if !nf.Valid {
		return 0
	}
	return nf.Float64
}

func (nf NullFloat64) IsNil() bool {
	return !nf.Valid
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		nf.Float64 = float64(x)
	case string:
		str := string(x)
		if len(str) == 0 {
			nf.Valid = false
			return nil
		}
		nf.Float64, err = strconv.ParseFloat(str, 64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &nf)
	case nil:
		nf.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullFloat64", reflect.TypeOf(v).Name())
	}
	nf.Valid = err == nil
	return err
}

func (nf *NullFloat64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		nf.Valid = false
		return nil
	}
	var err error
	nf.Float64, err = strconv.ParseFloat(string(text), 64)
	nf.Valid = err == nil
	return err
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	if math.IsInf(nf.Float64, 0) || math.IsNaN(nf.Float64) {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(nf.Float64),
			Str:   strconv.FormatFloat(nf.Float64, 'g', -1, 64),
		}
	}
	return []byte(strconv.FormatFloat(nf.Float64, 'f', -1, 64)), nil
}

func (nf NullFloat64) MarshalText() ([]byte, error) {
	if !nf.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(nf.Float64, 'f', -1, 64)), nil
}
