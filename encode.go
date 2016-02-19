// Package vpack implement velocypack encoding
package vpack

import (
	"math"
	"reflect"
)

func Marshal(v interface{}) []byte {
	return encode(nil, v)

}

func encode(b []byte, v interface{}) []byte {

	rv := reflect.ValueOf(v)

	switch rv.Kind() {

	case reflect.Array, reflect.Slice:
		if rv.Len() == 0 {
			b = append(b, 0x01)
			break
		}

		if rv.Type().Kind() == reflect.Uint8 {
			// blob
			l := rv.Len()
			i := len(b)
			b = append(b, 0xc0)
			for l != 0 {
				b = append(b, byte(l))
				l >>= 8
			}
			b[i] += byte(len(b) - i - 1)
			b = append(b, rv.Bytes()...)
			break
		}

		panic("unsupported item")

	case reflect.Ptr:
		if rv.IsNil() {
			b = append(b, 0x18)
			break
		}

	case reflect.Bool:
		if rv.Bool() {
			b = append(b, 0x1a)
		} else {
			b = append(b, 0x19)
		}

	case reflect.Float32, reflect.Float64:
		u := math.Float64bits(rv.Float())
		b = append(b, 0x1b, byte(u), byte(u>>8), byte(u>>16), byte(u>>24), byte(u>>32), byte(u>>40), byte(u>>48), byte(u>>56))

	case reflect.String:
		l := rv.Len()
		if l < 127 {
			b = append(b, 0x40+byte(l))
		} else {
			u := uint64(l)
			b = append(b, 0xbf, byte(u), byte(u>>8), byte(u>>16), byte(u>>24), byte(u>>32), byte(u>>40), byte(u>>48), byte(u>>56))
		}
		b = append(b, rv.String()...)
	default:
		panic("unsupported item")
	}

	return b
}

func encodeInt(b []byte, rv reflect.Value) []byte {

	// TODO(dgryski): optimize all numbers for size of value rather than size of variable, but include an input parameter to override

	// currently unused
	// - 0x30-0x39 : small integers 0, 1, ... 9
	// - 0x3a-0x3f : small negative integers -6, -5, ..., -1

	switch rv.Kind() {
	case reflect.Int8:
		b = append(b, 0x20, byte(rv.Int()))
	case reflect.Int16:
		u := rv.Int()
		b = append(b, 0x21, byte(u), byte(u>>8))
	case reflect.Int32:
		u := rv.Int()
		b = append(b, 0x23, byte(u), byte(u>>8), byte(u>>16), byte(u>>24))
	case reflect.Int64:
		u := rv.Int()
		b = append(b, 0x27, byte(u), byte(u>>8), byte(u>>16), byte(u>>24), byte(u>>32), byte(u>>40), byte(u>>48), byte(u>>56))
	case reflect.Uint8:
		b = append(b, 0x28, byte(rv.Uint()))
	case reflect.Uint16:
		u := rv.Uint()
		b = append(b, 0x29, byte(u), byte(u>>8))
	case reflect.Uint32:
		u := rv.Uint()
		b = append(b, 0x2b, byte(u), byte(u>>8), byte(u>>16), byte(u>>24))
	case reflect.Uint64:
		u := rv.Uint()
		b = append(b, 0x2f, byte(u), byte(u>>8), byte(u>>16), byte(u>>24), byte(u>>32), byte(u>>40), byte(u>>48), byte(u>>56))
	}

	return b
}
