package vpack

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
)

func Unmarshal(b []byte, v interface{}) {
	rv := reflect.ValueOf(v)
	decode(b, rv.Elem())
}

var ErrTruncated = errors.New("vpack: truncated")

func decode(b []byte, rv reflect.Value) error {

	v := b[0]

	switch {
	case v == 0x01:
		rv.Set(reflect.ValueOf([]interface{}{}))

	case v == 0x0a:
		rv.Set(reflect.ValueOf(map[string]interface{}{}))

	case v == 0x19 || v == 0x1a:
		rv.Set(reflect.ValueOf(b[0] == 0x1a))

	case v == 0x1b:
		b = b[1:]
		if len(b) < 8 {
			return ErrTruncated
		}
		bits := binary.LittleEndian.Uint64(b)
		rv.Set(reflect.ValueOf(math.Float64frombits(bits)))

	case 0x40 <= v && v < 0xbf:
		l := v - 0x40
		b = b[1:]
		if len(b) < int(l) {
			return ErrTruncated
		}

		rv.Set(reflect.ValueOf(string(b[:l])))

	case v == 0xbf:
		b = b[1:]
		if len(b) < 8 {
			return ErrTruncated
		}
		l := binary.LittleEndian.Uint64(b)
		b = b[8:]
		if uint64(len(b)) < l {
			return ErrTruncated
		}

		rv.Set(reflect.ValueOf(string(b[:l])))
	}

	return nil

}
