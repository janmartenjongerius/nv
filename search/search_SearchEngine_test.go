package search

import (
	"context"
	"fmt"
	"janmarten.name/nv/config"
)

func ExampleEngine_QueryAll() {
	engine := New(
		context.Background(),
		config.Variables{
			&config.Variable{
				Key:   "HOME",
				Value: "/home/gopher",
			},
			&config.Variable{
				Key:   "USER",
				Value: "gopher",
			},
		},
	)

	engine.QueryAll([]string{"HOME", "USER"}, 0)

	for _, r := range engine.Results() {
		fmt.Printf("%s -> %q\n", r.Request.Query, r.Match)
	}

	// Unordered output:
	// HOME -> &{"HOME" "/home/gopher"}
	// USER -> &{"USER" "gopher"}
}

func ExampleEngine_Query_interrupted() {
	ctx, cancel := context.WithCancel(context.Background())

	engine := New(
		ctx,
		config.Variables{
			&config.Variable{
				Key:   "HOME",
				Value: "/home/gopher",
			},
			&config.Variable{
				Key:   "USER",
				Value: "gopher",
			},
		},
	)

	engine.Query("HOME", 0)
	cancel()
	engine.Query("USER", 0)

	for _, r := range engine.Results() {
		fmt.Printf("%s -> %q", r.Request.Query, r.Match)
	}

	// Output: HOME -> &{"HOME" "/home/gopher"}
}
