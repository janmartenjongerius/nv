package main

import (
	"fmt"
	"janmarten.name/nv/config"
	"reflect"
	"testing"
)

var encodeTests = []struct {
	in config.Variables
	want []byte
}{
	{
		in: config.Variables{
			&config.Variable{
				Key: "HOME",
				Value: "C:\\Users\\Gopher",
			},
			&config.Variable{
				Key: "USERNAME",
				Value: "Gopher",
			},
			&config.Variable{
				Key: "EMAIL",
				Value: "Gopher <gopher@golang.goph>",
			},
		},
		want: []byte(`<Config>
	<Variable key="HOME">C:\Users\Gopher</Variable>
	<Variable key="USERNAME">Gopher</Variable>
	<Variable key="EMAIL">Gopher &lt;gopher@golang.goph&gt;</Variable>
</Config>`),
	},
	{
		in: config.Variables{},
		want: []byte(`<Config></Config>`),
	},
}

func TestXmlEncoder_Encode(t *testing.T) {
	enc, err := config.NewEncoding("xml")

	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	for _, tt := range encodeTests {
		var keys = make([]string, 0)

		for _, v := range tt.in {
			keys = append(keys, v.Key)
		}

		t.Run(fmt.Sprintf("%q", keys), func(t *testing.T) {
			got, err := enc.Encode(tt.in...)

			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf(
					"exp.Export(%v) does not output '%v'. Got '%v'",
					tt.in,
					string(tt.want),
					string(got))
			}
		})
	}
}

var decodeTests = []struct {
	in []byte
	want config.Variables
}{
	{
		in: []byte(`
		<Config>
	        <Variable key="HOME">C:\Users\Gopher</Variable>
    	    <Variable key="USERNAME">Gopher</Variable>
			<Variable key="EMAIL">Gopher &lt;gopher@golang.goph&gt;</Variable>
		</Config>

`),
		want: config.Variables{
			&config.Variable{
				Key: "HOME",
				Value: "C:\\Users\\Gopher",
			},
			&config.Variable{
				Key: "USERNAME",
				Value: "Gopher",
			},
			&config.Variable{
				Key: "EMAIL",
				Value: "Gopher <gopher@golang.goph>",
			},
		},
	},
	{
		in: []byte(`<Config></Config>`),
		want: config.Variables{},
	},
}

func TestXmlDecoder_Decode(t *testing.T) {
	enc, err := config.NewEncoding("xml")

	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	for _, tt := range decodeTests {
		var keys = make([]string, 0)

		for _, v := range tt.want {
			keys = append(keys, v.Key)
		}

		t.Run(fmt.Sprintf("%q", keys), func(t *testing.T) {
			got, err := enc.Decode(tt.in)

			if err != nil {
				fmt.Println(err)
				t.Fail()
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Expect '%v' for '%v', got '%v'", tt.want, tt.in, got)
			}
		})
	}
}
