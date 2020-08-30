package config

import (
	"janmarten.name/nv/debug"
	"reflect"
	"testing"
)

func TestGetEncodings(t *testing.T) {
	for _, format := range GetEncodings() {
		if !HasEncoding(format) {
			t.Errorf("Got encoding %q, but it is not available.", format)
		}
	}
}

func TestHasEncoding(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{
			in:   "non-existent",
			want: false,
		},
		{
			in:   "foo",
			want: false,
		},
	}

	for _, format := range GetEncodings() {
		cases = append(
			cases,
			struct {
				in   string
				want bool
			}{in: format, want: true})
	}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			if HasEncoding(c.in) != c.want {
				t.Errorf("HasEncoding(%q) != %v", c.in, c.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	got := debug.Scope("Encoding").GetMessages()
	want := debug.Messages{
		"Default": DefaultEncoding,
		"Formats": GetEncodings(),
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Unexpected debug messages. Want: %q, got: %q", want, got)
	}
}
