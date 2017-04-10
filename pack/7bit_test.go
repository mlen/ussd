package pack

import (
	"bytes"
	"github.com/google/gofuzz"
	"testing"
)

func TestPack7Bit(t *testing.T) {
	cases := []struct {
		example  []byte
		expected []byte
	}{
		{
			example:  []byte{0x31, 0x32, 0x33, 0x34},
			expected: []byte{0x31, 0xd9, 0x8c, 0x06},
		},
		{
			example:  []byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39},
			expected: []byte{0xb0, 0x98, 0x6c, 0x46, 0xab, 0xd9, 0x6e, 0xb8, 0x1c},
		},
		{
			example:  []byte{0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x77, 0x78, 0x79, 0x7a},
			expected: []byte{0x61, 0xf1, 0x98, 0x5c, 0x36, 0x9f, 0xd1, 0x69, 0xf5, 0x9a, 0xdd, 0x76, 0xbf, 0xe1, 0x71, 0xf9, 0x9c, 0x5e, 0xbf, 0xe3, 0xf3, 0x7a},
		},
		{
			example:  []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37},
			expected: []byte{0x31, 0xd9, 0x8c, 0x56, 0xb3, 0xdd, 0x62, 0xb2, 0x19, 0xad, 0x66, 0xbb, 0xc5, 0x64, 0x33, 0x5a, 0xcd, 0x76, 0x03},
		},
		{
			example:  []byte{},
			expected: []byte{},
		},
	}

	for _, c := range cases {
		result := Pack7Bit(c.example)
		if !bytes.Equal(result, c.expected) {
			t.Errorf("\nExpected %08b\nGot      %08b", c.expected, result)
		}
	}
}

func TestUnpack7Bit(t *testing.T) {
	cases := []struct {
		example  []byte
		expected []byte
	}{
		{
			example:  []byte{0x31, 0xd9, 0x8c, 0x06},
			expected: []byte{0x31, 0x32, 0x33, 0x34},
		},
		{
			example:  []byte{0xb0, 0x98, 0x6c, 0x46, 0xab, 0xd9, 0x6e, 0xb8, 0x1c},
			expected: []byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39},
		},
		{
			example:  []byte{0x61, 0xf1, 0x98, 0x5c, 0x36, 0x9f, 0xd1, 0x69, 0xf5, 0x9a, 0xdd, 0x76, 0xbf, 0xe1, 0x71, 0xf9, 0x9c, 0x5e, 0xbf, 0xe3, 0xf3, 0x7a},
			expected: []byte{0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x77, 0x78, 0x79, 0x7a},
		},
		{
			example:  []byte{0x31, 0xd9, 0x8c, 0x56, 0xb3, 0xdd, 0x62, 0xb2, 0x19, 0xad, 0x66, 0xbb, 0xc5, 0x64, 0x33, 0x5a, 0xcd, 0x76, 0x03},
			expected: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37},
		},
		{
			example:  []byte{},
			expected: []byte{},
		},
	}

	for _, c := range cases {
		result := Unpack7Bit(c.example)
		if !bytes.Equal(result, c.expected) {
			t.Errorf("\nExpected %08b\nGot      %08b", c.expected, result)
		}
	}
}

func TestFuzz(t *testing.T) {
	var instance []byte
	f := fuzz.New().NumElements(10, 30).Funcs(
		func(b *byte, c fuzz.Continue) {
			c.FuzzNoCustom(b)
			*b = *b & 0x7f
		},
	)

	for i := 0; i < 1000; i++ {
		f.Fuzz(&instance)

		// instances with null byte at the end are invalid
		if len(instance) != 0 && instance[len(instance)-1] == 0 {
			instance = instance[:len(instance)-1]
		}

		result := Unpack7Bit(Pack7Bit(instance))
		if !bytes.Equal(result, instance) {
			t.Errorf("\nExpected %08b\nGot      %08b", instance, result)
		}
	}
}
