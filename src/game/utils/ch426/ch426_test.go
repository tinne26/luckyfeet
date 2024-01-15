package ch426

import "testing"
import "slices"

func TestEncode(t *testing.T) {
	tests := []struct{ in []byte ; out string }{
		{nil, ""},
		{[]byte{0}, "AA"},
		{[]byte{1}, "A~"},
		{[]byte{2}, "BA"},
		{[]byte{0, 255}, "A}ó"},
		// 00000000 11111111
		// 0000000 0111111 1100000
		// 0       63      96
		{[]byte{1, 255}, "AÜó"},
		{[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, "AAggYQKGD@BQ"},
		// 00000000 00000001 00000010 00000011 00000100 00000101 00000110 00000111 00001000 00001001
		// 0000000 0000000 0100000 0100000 0011000 0010000 0001010 0000110 0000011 1000010 0000001 0010000
		// 0       0       32      32      24      16      10      6       3       66      1       16
		// A       A       g       g       Y       Q       K       G       D       @       B       Q
	}

	for i, test := range tests {
		result := Encode(test.in)
		if result != test.out {
			t.Fatalf("test #%d, expected %v to encode as '%s', but got '%s' instead", i, test.in, test.out, result)
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	tests := [][]byte{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, // "AAggYQKGD@BQ"
		{255},
		{255, 254},
		{255, 254, 253},
		{255, 254, 253, 252},
		{255, 254, 253, 252, 251},
		{255, 254, 253, 252, 251, 250},
		{255, 254, 253, 252, 251, 250, 249},
		{1, 254, 2, 252, 3, 250, 4, 248},
		{0, 128, 0, 128, 0, 128, 0, 128, 0, 128, 0, 128, 0, 128},
	}

	for i, test := range tests {
		result, err := Decode(Encode(test))
		if err != nil {
			t.Fatalf("test #%d, unexpected decoding error: %s", i, err)
		}
		if !slices.Equal(test, result) {
			t.Fatalf("test #%d, expected %v, got %v", i, test, result)
		}
	}
}
