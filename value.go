package xorm

import "time"

type Value []byte

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
