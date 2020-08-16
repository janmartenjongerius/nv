package config

import (
	"reflect"
	"testing"
)

func TestJsonDecoder_Decode(t *testing.T) {
	cases := []struct {
		in   []byte
		want Variables
		err  bool
	}{
		{
			in: []byte("{\"HOME\":\"\\/home\\/gopher\"}"),
			want: Variables{
				&Variable{
					Key:   "HOME",
					Value: "/home/gopher",
				},
			},
			err: false,
		},
		{
			in:   []byte("{\"HOME\":42}"),
			want: Variables{},
			err:  true,
		},
	}

	decoder, err := NewEncoding("json")

	if err != nil {
		t.Fatalf("Cannot initiate decoder: %v", err)
	}

	for _, c := range cases {
		t.Run(string(c.in), func(t *testing.T) {
			got, err := decoder.Decode(c.in)

			if c.err == false && err != nil {
				t.Errorf("Unexpected error %q", err)
			}

			if len(c.want)+len(got) > 0 && !reflect.DeepEqual(got, c.want) {
				t.Errorf("Unexpected variables. Got %v; Want %v", got, c.want)
			}
		})
	}
}
