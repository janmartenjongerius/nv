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
			"%s	-> %q\t(Suggestions: %d)\n",
			r.Request.Query,
			r.Match,
			len(r.Suggestions))
	}

	// Unordered output:
	// HOME	-> &{"HOME" "/home/gopher"}	(Suggestions: 0)
	// user	-> %!q(*config.Variable=<nil>)	(Suggestions: 1)
	// SHL	-> %!q(*config.Variable=<nil>)	(Suggestions: 2)
}
