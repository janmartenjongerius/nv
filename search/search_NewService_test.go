package search

import (
	"fmt"
	"janmarten.name/nv/config"
)

func ExampleNewService() {
	// This is a mock of config.Environment
	environment := config.Variables{
		&config.Variable{
			Key:   "HOME",
			Value: "/home/gopher",
		},
		&config.Variable{
			Key:   "USER",
			Value: "gopher",
		},
		&config.Variable{
			Key:   "SHELL",
			Value: "/bin/gopher",
		},
		&config.Variable{
			Key:   "SHLVL",
			Value: "1",
		},
	}

	svc := NewService(environment)
	svc.Suggestions = 5

	for _, r := range svc.Search("HOME", "user", "SHL", "HOME") {
		fmt.Printf(
			"%s -> %q\n",
			r.Request.Query,
			r.Match)

		for _, s := range r.Suggestions {
			fmt.Printf("\t%q\n", s)
		}

		fmt.Println("")
	}

	// Output:
	// HOME -> &{"HOME" "/home/gopher"}
	//
	// user -> %!q(*config.Variable=<nil>)
	// 	"USER"
	//
	// SHL -> %!q(*config.Variable=<nil>)
	// 	"SHELL"
	// 	"SHLVL"
}
