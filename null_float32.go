package xorm

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type NullFloat32 struct {
	Float32 float32
	Valid   bool
}

func (nf NullFloat32) Ptr() *float32 {
	if !nf.Valid {
		return nil
	}
	return &nf.Float32
}

func (nf NullFloat32) ValueOrZero() float32 {
	if !nf.Valid {
		return 0
	}
	return nf.Float32
}

func (nf NullFloat32) IsNil() bool {
	return !nf.Valid
}

func (nf *NullFloat32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		nf.Float32 = float32(x)
	case string:
		str := string(x)
		if len(str) == 0 {
			nf.Valid = false
			return nil
		}
		var f float64
		f, err = strconv.ParseFloat(str, 32)
		if err == nil {
			nf.Float32 = float32(f)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &nf)
	case nil:
		nf.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type xorm.NullFloat32", reflect.TypeOf(v).Name())
	}
	nf.Valid = err == nil
	return err
}

func (nf *NullFloat32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		nf.Valid = false
		return nil
	}
	var err error
	var f float64
	f, err = strconv.ParseFloat(string(text), 32)
	if err == nil {
		nf.Float32 = float32(f)
	}
	nf.Valid = err == nil
	return err
}

func (nf NullFloat32) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	if math.IsInf(float64(nf.Float32), 0) || math.IsNaN(float64(nf.Float32)) {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(nf.Float32),
			Str:   strconv.FormatFloat(float64(nf.Float32), 'g', -1, 32),
		}
	}
	return []byte(strconv.FormatFloat(float64(nf.Float32), 'f', -1, 32)), nil
}

func (nf NullFloat32) MarshalText() ([]byte, error) {
	if !nf.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(float64(nf.Float32), 'f', -1, 32)), nil
}
