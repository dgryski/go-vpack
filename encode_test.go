package vpack

import (
	"reflect"
	"testing"
	"testing/quick"
)

func TestEncode(t *testing.T) {

	var b []byte

	t.Logf("%x", Marshal(b))
}

func TestDecoder(t *testing.T) {

	tests := []struct {
		b    []byte
		want interface{}
	}{
		{
			[]byte{0x01},
			[]interface{}{},
		},
		{
			[]byte{0x0a},
			map[string]interface{}{},
		},
		{
			[]byte{0x19},
			false,
		},
		{
			[]byte{0x1a},
			true,
		},
	}

	for _, tt := range tests {
		var i interface{}
		Unmarshal(tt.b, &i)
		if !reflect.DeepEqual(i, tt.want) {
			t.Errorf("Unmarshal(%x)=%v, want %v\n", tt.b, i, tt.want)
		}
	}
}

func TestDecodeString(t *testing.T) {

	f := func(s string) bool {

		b := Marshal(s)
		var i interface{}
		Unmarshal(b, &i)

		return i.(string) == s
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestDecodeFloat(t *testing.T) {

	f := func(f64 float64, f32 float32) bool {

		b64 := Marshal(f64)
		b32 := Marshal(f32)

		var i64 interface{}
		Unmarshal(b64, &i64)

		var i32 interface{}
		Unmarshal(b32, &i32)

		return i64.(float64) == f64 && i32.(float64) == float64(f32)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
