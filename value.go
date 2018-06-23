package xorm

import (
	"time"
)

type Value []byte

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

func (ns NullString) IsZero() bool {
	return !ns.Valid
}

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

func (nb NullBool) IsZero() bool {
	return !nb.Valid
}

type NullInt struct {
	Int   int
	Valid bool
}

func (ni NullInt) Ptr() *int {
	if !ni.Valid {
		return nil
	}
	return &ni.Int
}

func (ni NullInt) ValueOrZero() int {
	if !ni.Valid {
		return 0
	}
	return ni.Int
}

func (ni NullInt) IsZero() bool {
	return !ni.Valid
}

type NullInt8 struct {
	Int8  int8
	Valid bool
}

func (ni NullInt8) Ptr() *int8 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int8
}

func (ni NullInt8) ValueOrZero() int8 {
	if !ni.Valid {
		return 0
	}
	return ni.Int8
}

func (ni NullInt8) IsZero() bool {
	return !ni.Valid
}

type NullInt16 struct {
	Int16 int16
	Valid bool
}

func (ni NullInt16) Ptr() *int16 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int16
}

func (ni NullInt16) ValueOrZero() int16 {
	if !ni.Valid {
		return 0
	}
	return ni.Int16
}

func (ni NullInt16) IsZero() bool {
	return !ni.Valid
}

type NullInt32 struct {
	Int32 int32
	Valid bool
}

func (ni NullInt32) Ptr() *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}

func (ni NullInt32) ValueOrZero() int32 {
	if !ni.Valid {
		return 0
	}
	return ni.Int32
}

func (ni NullInt32) IsZero() bool {
	return !ni.Valid
}

type NullInt64 struct {
	Int64 int64
	Valid bool
}

func (ni NullInt64) Ptr() *int64 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int64
}

func (ni NullInt64) ValueOrZero() int64 {
	if !ni.Valid {
		return 0
	}
	return ni.Int64
}

func (ni NullInt64) IsZero() bool {
	return !ni.Valid
}

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

func (nu NullUint) IsZero() bool {
	return !nu.Valid
}

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

func (nu NullUint8) IsZero() bool {
	return !nu.Valid
}

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

func (nu NullUint16) IsZero() bool {
	return !nu.Valid
}

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

func (nu NullUint32) IsZero() bool {
	return !nu.Valid
}

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

func (nu NullUint64) IsZero() bool {
	return !nu.Valid
}

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

func (nf NullFloat32) IsZero() bool {
	return !nf.Valid
}

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

func (nf NullFloat64) IsZero() bool {
	return !nf.Valid
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt NullTime) Ptr() *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func (nt NullTime) ValueOrZero() time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}

func (v Value) Bytes() []byte {
	return []byte(v)
}

func (v Value) String() string {
	return string(v)
}

func (v Value) NullString() NullString {
	if v == nil {
		return NullString{
			String: "",
			Valid:  false,
		}
	} else {
		return NullString{
			String: string(v),
			Valid:  true,
		}
	}
}

func (v Value) Bool() bool {
	return Bool(v)
}

func (v Value) NullBool() NullBool {
	if v == nil {
		return NullBool{
			Bool:  false,
			Valid: false,
		}
	} else {
		return NullBool{
			Bool:  Bool(v),
			Valid: true,
		}
	}
}

func (v Value) Int() int {
	return Int(v)
}

func (v Value) NullInt() NullInt {
	if v == nil {
		return NullInt{
			Int:   0,
			Valid: false,
		}
	} else {
		return NullInt{
			Int:   Int(v),
			Valid: true,
		}
	}
}

func (v Value) Int8() int8 {
	return Int8(v)
}

func (v Value) NullInt8() NullInt8 {
	if v == nil {
		return NullInt8{
			Int8:  0,
			Valid: false,
		}
	} else {
		return NullInt8{
			Int8:  Int8(v),
			Valid: true,
		}
	}
}

func (v Value) Int16() int16 {
	return Int16(v)
}

func (v Value) NullInt16() NullInt16 {
	if v == nil {
		return NullInt16{
			Int16: 0,
			Valid: false,
		}
	} else {
		return NullInt16{
			Int16: Int16(v),
			Valid: true,
		}
	}
}

func (v Value) Int32() int32 {
	return Int32(v)
}

func (v Value) NullInt32() NullInt32 {
	if v == nil {
		return NullInt32{
			Int32: 0,
			Valid: false,
		}
	} else {
		return NullInt32{
			Int32: Int32(v),
			Valid: true,
		}
	}
}

func (v Value) Int64() int64 {
	return Int64(v)
}

func (v Value) NullInt64() NullInt64 {
	if v == nil {
		return NullInt64{
			Int64: 0,
			Valid: false,
		}
	} else {
		return NullInt64{
			Int64: Int64(v),
			Valid: true,
		}
	}
}

func (v Value) Uint() uint {
	return Uint(v)
}

func (v Value) NullUint() NullUint {
	if v == nil {
		return NullUint{
			Uint:  0,
			Valid: false,
		}
	} else {
		return NullUint{
			Uint:  Uint(v),
			Valid: true,
		}
	}
}

func (v Value) Uint8() uint8 {
	return Uint8(v)
}

func (v Value) NullUint8() NullUint8 {
	if v == nil {
		return NullUint8{
			Uint8: 0,
			Valid: false,
		}
	} else {
		return NullUint8{
			Uint8: Uint8(v),
			Valid: true,
		}
	}
}

func (v Value) Uint16() uint16 {
	return Uint16(v)
}

func (v Value) NullUint16() NullUint16 {
	if v == nil {
		return NullUint16{
			Uint16: 0,
			Valid:  false,
		}
	} else {
		return NullUint16{
			Uint16: Uint16(v),
			Valid:  true,
		}
	}
}

func (v Value) Uint32() uint32 {
	return Uint32(v)
}

func (v Value) NullUint32() NullUint32 {
	if v == nil {
		return NullUint32{
			Uint32: 0,
			Valid:  false,
		}
	} else {
		return NullUint32{
			Uint32: Uint32(v),
			Valid:  true,
		}
	}
}

func (v Value) Uint64() uint64 {
	return Uint64(v)
}

func (v Value) NullUint64() NullUint64 {
	if v == nil {
		return NullUint64{
			Uint64: 0,
			Valid:  false,
		}
	} else {
		return NullUint64{
			Uint64: Uint64(v),
			Valid:  true,
		}
	}
}

func (v Value) Float32() float32 {
	return Float32(v)
}

func (v Value) NullFloat32() NullFloat32 {
	if v == nil {
		return NullFloat32{
			Float32: 0,
			Valid:   false,
		}
	} else {
		return NullFloat32{
			Float32: Float32(v),
			Valid:   true,
		}
	}
}

func (v Value) Float64() float64 {
	return Float64(v)
}

func (v Value) NullFloat64() NullFloat64 {
	if v == nil {
		return NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	} else {
		return NullFloat64{
			Float64: Float64(v),
			Valid:   true,
		}
	}
}

func (v Value) Time(format string, TZLocation ...*time.Location) time.Time {
	return Time(v, format, TZLocation...)
}

func (v Value) NullTime(format string, TZLocation ...*time.Location) NullTime {
	if v == nil {
		return NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	} else {
		return NullTime{
			Time:  Time(v, format, TZLocation...),
			Valid: true,
		}
	}

}

func (v Value) TimeDuration() time.Duration {
	return TimeDuration(v)
}
