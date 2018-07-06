// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

var errNilPtr = errors.New("destination pointer is nil") // embedded in descriptive error

func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	} else {
		c := make([]byte, len(b))
		copy(c, b)
		return c
	}
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}
	return
}

// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
func convertAssign(dest, src interface{}) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s)
			return nil
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = string(s)
			return nil
		case *interface{}:
			if d == nil {
				return errNilPtr
			}
			*d = cloneBytes(s)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = cloneBytes(s)
			return nil
		}

	case time.Time:
		switch d := dest.(type) {
		case *string:
			*d = s.Format(time.RFC3339Nano)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s.Format(time.RFC3339Nano))
			return nil
		}
	case nil:
		switch d := dest.(type) {
		case *interface{}:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		}
	}

	var sv reflect.Value

	switch d := dest.(type) {
	case *string:
		sv = reflect.ValueOf(src)
		switch sv.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			*d = asString(src)
			return nil
		}
	case *[]byte:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes(nil, sv); ok {
			*d = b
			return nil
		}
	case *bool:
		bv, err := driver.Bool.ConvertValue(src)
		if err == nil {
			*d = bv.(bool)
		}
		return err
	case *interface{}:
		*d = src
		return nil
	}

	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		return errors.New("destination not a pointer")
	}
	if dpv.IsNil() {
		return errNilPtr
	}

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		switch b := src.(type) {
		case []byte:
			dv.Set(reflect.ValueOf(cloneBytes(b)))
		default:
			dv.Set(sv)
		}
		return nil
	}

	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	switch dv.Kind() {
	case reflect.Ptr:
		if src == nil {
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		}

		dv.Set(reflect.New(dv.Type().Elem()))
		return convertAssign(dv.Interface(), src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := asString(src)
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := asString(src)
		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		s := asString(src)
		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetFloat(f64)
		return nil
	case reflect.String:
		dv.SetString(asString(src))
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, dest)
}

func asKind(vv reflect.Value, tp reflect.Type) (interface{}, error) {
	switch tp.Kind() {
	case reflect.Int64:
		return vv.Int(), nil
	case reflect.Int:
		return int(vv.Int()), nil
	case reflect.Int32:
		return int32(vv.Int()), nil
	case reflect.Int16:
		return int16(vv.Int()), nil
	case reflect.Int8:
		return int8(vv.Int()), nil
	case reflect.Uint64:
		return vv.Uint(), nil
	case reflect.Uint:
		return uint(vv.Uint()), nil
	case reflect.Uint32:
		return uint32(vv.Uint()), nil
	case reflect.Uint16:
		return uint16(vv.Uint()), nil
	case reflect.Uint8:
		return uint8(vv.Uint()), nil
	case reflect.String:
		return vv.String(), nil
	case reflect.Slice:
		if tp.Elem().Kind() == reflect.Uint8 {
			v, err := strconv.ParseInt(string(vv.Interface().([]byte)), 10, 64)
			if err != nil {
				return nil, err
			}
			return v, nil
		}

	}
	return nil, fmt.Errorf("unsupported primary key type: %v, %v", tp, vv)
}

func convertFloat(v interface{}) (float64, error) {
	switch v.(type) {
	case float32:
		return float64(v.(float32)), nil
	case float64:
		return v.(float64), nil
	case string:
		i, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	case []byte:
		i, err := strconv.ParseFloat(string(v.([]byte)), 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, fmt.Errorf("unsupported type: %v", v)
}

func convertInt(v interface{}) (int64, error) {
	switch v.(type) {
	case int:
		return int64(v.(int)), nil
	case int8:
		return int64(v.(int8)), nil
	case int16:
		return int64(v.(int16)), nil
	case int32:
		return int64(v.(int32)), nil
	case int64:
		return v.(int64), nil
	case []byte:
		i, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	case string:
		i, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, fmt.Errorf("unsupported type: %v", v)
}

func asBool(bs []byte) (bool, error) {
	if len(bs) == 0 {
		return false, nil
	}
	if bs[0] == 0x00 {
		return false, nil
	} else if bs[0] == 0x01 {
		return true, nil
	}
	return strconv.ParseBool(string(bs))
}

func EncodeString(s string) []byte {
	return []byte(s)
}

func DecodeToString(b []byte) string {
	return string(b)
}

func EncodeBool(b bool) []byte {
	if b == true {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

func EncodeInt(i int) []byte {
	if i <= math.MaxInt8 {
		return EncodeInt8(int8(i))
	} else if i <= math.MaxInt16 {
		return EncodeInt16(int16(i))
	} else if i <= math.MaxInt32 {
		return EncodeInt32(int32(i))
	} else {
		return EncodeInt64(int64(i))
	}
}

func EncodeUint(i uint) []byte {
	if i <= math.MaxUint8 {
		return EncodeUint8(uint8(i))
	} else if i <= math.MaxUint16 {
		return EncodeUint16(uint16(i))
	} else if i <= math.MaxUint32 {
		return EncodeUint32(uint32(i))
	} else {
		return EncodeUint64(uint64(i))
	}
}

func EncodeInt8(i int8) []byte {
	return []byte{byte(i)}
}

func EncodeUint8(i uint8) []byte {
	return []byte{byte(i)}
}

func EncodeInt16(i int16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(i))
	return bytes
}

func EncodeUint16(i uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, i)
	return bytes
}

func EncodeInt32(i int32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(i))
	return bytes
}

func EncodeUint32(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, i)
	return bytes
}

func EncodeInt64(i int64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(i))
	return bytes
}

func EncodeUint64(i uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, i)
	return bytes
}

func EncodeFloat32(f float32) []byte {
	bits := math.Float32bits(f)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func EncodeFloat64(f float64) []byte {
	bits := math.Float64bits(f)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func Encode(vs ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for i := 0; i < len(vs); i++ {
		switch value := vs[i].(type) {
		case int:
			buf.Write(EncodeInt(value))
		case int8:
			buf.Write(EncodeInt8(value))
		case int16:
			buf.Write(EncodeInt16(value))
		case int32:
			buf.Write(EncodeInt32(value))
		case int64:
			buf.Write(EncodeInt64(value))
		case uint:
			buf.Write(EncodeUint(value))
		case uint8:
			buf.Write(EncodeUint8(value))
		case uint16:
			buf.Write(EncodeUint16(value))
		case uint32:
			buf.Write(EncodeUint32(value))
		case uint64:
			buf.Write(EncodeUint64(value))
		case bool:
			buf.Write(EncodeBool(value))
		case string:
			buf.Write(EncodeString(value))
		case []byte:
			buf.Write(value)
		case float32:
			buf.Write(EncodeFloat32(value))
		case float64:
			buf.Write(EncodeFloat64(value))
		default:
			if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
				buf.Write(EncodeString(fmt.Sprintf("%v", value)))
			}
		}
	}
	return buf.Bytes()
}

func IsNumeric(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < byte('0') || s[i] > byte('9') {
			return false
		}
	}
	return true
}

func Time(i interface{}, format string, TZLocation ...*time.Location) time.Time {
	s := String(i)
	t, _ := StrToTime(s, format, TZLocation...)
	return t
}

func StrToTime(str string, format string, TZLocation ...*time.Location) (time.Time, error) {
	if len(TZLocation) > 0 {
		if t, err := time.ParseInLocation(format, str, TZLocation[0]); err == nil {
			return t, nil
		} else {
			return time.Time{}, err
		}
	} else {
		if t, err := time.ParseInLocation(format, str, time.Local); err == nil {
			return t, nil
		} else {
			return time.Time{}, err
		}
	}
}

func TimeDuration(i interface{}) time.Duration {
	return time.Duration(Int64(i))
}

func Bytes(i interface{}) []byte {
	if i == nil {
		return nil
	}
	if r, ok := i.([]byte); ok {
		return r
	} else {
		return Encode(i)
	}
}

func String(i interface{}) string {
	if i == nil {
		return ""
	}
	switch value := i.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.Itoa(int(value))
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(uint64(value), 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func Strings(i interface{}) []string {
	if i == nil {
		return nil
	}
	if r, ok := i.([]string); ok {
		return r
	} else if r, ok := i.([]interface{}); ok {
		strs := make([]string, len(r))
		for k, v := range r {
			strs[k] = String(v)
		}
		return strs
	}
	return []string{fmt.Sprintf("%v", i)}
}

//false: "", 0, false, off
func Bool(i interface{}) bool {
	if i == nil {
		return false
	}
	if v, ok := i.(bool); ok {
		return v
	}
	if s := String(i); s != "" && s != "0" && s != "false" && s != "off" {
		return true
	}
	return false
}

func Int(i interface{}) int {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case int:
		return value
	case int8:
		return int(value)
	case int16:
		return int(value)
	case int32:
		return int(value)
	case int64:
		return int(value)
	case uint:
		return int(value)
	case uint8:
		return int(value)
	case uint16:
		return int(value)
	case uint32:
		return int(value)
	case uint64:
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	case bool:
		if value {
			return 1
		}
		return 0
	default:
		v, _ := strconv.Atoi(String(value))
		return v
	}
}

func Int8(i interface{}) int8 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int8); ok {
		return v
	}
	return int8(Int(i))
}

func Int16(i interface{}) int16 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int16); ok {
		return v
	}
	return int16(Int(i))
}

func Int32(i interface{}) int32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int32); ok {
		return v
	}
	return int32(Int(i))
}

func Int64(i interface{}) int64 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int64); ok {
		return v
	}
	return int64(Int(i))
}

func Uint(i interface{}) uint {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case int:
		return uint(value)
	case int8:
		return uint(value)
	case int16:
		return uint(value)
	case int32:
		return uint(value)
	case int64:
		return uint(value)
	case uint:
		return value
	case uint8:
		return uint(value)
	case uint16:
		return uint(value)
	case uint32:
		return uint(value)
	case uint64:
		return uint(value)
	case float32:
		return uint(value)
	case float64:
		return uint(value)
	case bool:
		if value {
			return 1
		}
		return 0
	default:
		v, _ := strconv.ParseUint(String(value), 10, 64)
		return uint(v)
	}
}

func Uint8(i interface{}) uint8 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint8); ok {
		return v
	}
	return uint8(Uint(i))
}

func Uint16(i interface{}) uint16 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint16); ok {
		return v
	}
	return uint16(Uint(i))
}

func Uint32(i interface{}) uint32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint32); ok {
		return v
	}
	return uint32(Uint(i))
}

func Uint64(i interface{}) uint64 {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint64); ok {
		return v
	}
	return uint64(Uint(i))
}

func Float32(i interface{}) float32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(float32); ok {
		return v
	}
	v, _ := strconv.ParseFloat(String(i), 32)
	return float32(v)
}

func Float64(i interface{}) float64 {
	if i == nil {
		return 0
	}
	if v, ok := i.(float64); ok {
		return v
	}
	v, _ := strconv.ParseFloat(String(i), 64)
	return v
}
