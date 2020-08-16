package config

import (
	"errors"
	"reflect"
	"testing"
)

func TestTextDecoder_Decode(t *testing.T) {
	cases := []struct {
		in   []byte
		want Variables
		err  error
	}{
		{
			in: []byte("HOME=/home/gopher"),
			want: Variables{
				&Variable{
					Key:   "HOME",
					Value: "/home/gopher",
				},
			},
			err: nil,
		},
		{
			in: []byte("HOME=/home/gopher\n\nUSER=gopher"),
			want: Variables{
				&Variable{
					Key:   "HOME",
					Value: "/home/gopher",
				},
				&Variable{
					Key:   "USER",
					Value: "gopher",
				},
			},
			err: nil,
		},
		{
			in:   []byte("HOME:/home/gopher"),
			want: Variables{},
			err:  IllFormattedVariable,
		},
	}

	decoder, err := NewEncoding(DefaultEncoding)

	if err != nil {
		t.Fatalf("Cannot initiate decoder: %v", err)
	}

	for _, c := range cases {
		t.Run(string(c.in), func(t *testing.T) {
			got, err := decoder.Decode(c.in)

			if !errors.Is(err, c.err) {
				t.Errorf("Unexpected error %q; expected %q", err, c.err)
			}

			if len(c.want)+len(got) > 0 && !reflect.DeepEqual(got, c.want) {
				t.Errorf("Unexpected variables. Got %v; Want %v", got, c.want)
			}
		})
	}
}
