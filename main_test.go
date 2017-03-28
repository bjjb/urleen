package main

import "testing"

func TestUrleen(t *testing.T) {
	t.Run("encode", testEncode)
	t.Run("decode", testDecode)
}

func testEncode(t *testing.T) {
	cases := []struct {
		in   uint64
		want string
	}{
		{0, "0"},
		{9, "9"},
		{10, "a"},
		{35, "z"},
		{36, "A"},
		{61, "Z"},
		{62, "10"},
	}

	for _, c := range cases {
		got := encode(c.in)
		if c.want != got {
			t.Errorf("decode(%d); %q != %q", c.in, got, c.want)
		}
	}
}

func testDecode(t *testing.T) {
	cases := []struct {
		in   string
		want uint64
	}{
		{"0", 0},
		{"9", 9},
		{"a", 10},
		{"z", 35},
		{"A", 36},
		{"Z", 61},
		{"10", 62},
	}

	for _, c := range cases {
		got := decode(c.in)
		if c.want != got {
			t.Errorf("%d != %d", got, c.want)
		}
	}
}
