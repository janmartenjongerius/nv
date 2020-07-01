package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"janmarten.name/env/config"
	"janmarten.name/env/neighbor"
	"log"
	"os"
	"strings"
	"time"
)

// Print the given neighbor.Neighbors as a list of suggestions.
func printSuggestions(w *color.Color, neighbors neighbor.Neighbors) {
	suggestion := "Suggestions:\n"

	for _, n := range neighbors {
		suggestion += fmt.Sprintf("  - %s\n", n.Name)
	}

	if _, e := w.Fprint(os.Stderr, suggestion); e != nil {
		log.Fatal(e)
	}
}

func main() {
	var (
		cfg            = config.Parse(os.Environ(), "=")
		stdin          = bufio.NewReader(os.Stdin)
		noAnsi         = color.NoColor
		numSuggestions = 5
	)

	flag.BoolVar(
		&noAnsi,
		"no-ansi",
		color.NoColor,
		"Whether to suppress decorating output with ANSI escape codes")
	flag.IntVar(
		&numSuggestions,
		"num-suggestions",
		numSuggestions,
		"Set the number of suggestions returned when an entry could not be found")
	flag.Parse()

	color.NoColor = noAnsi
	var (
		success = color.New(color.FgGreen, color.Bold)
		info    = color.New(color.FgCyan)
		comment = color.New(color.FgYellow)
	)

	if _, e := info.Fprint(os.Stdout, "Environment lookup"); e != nil {
		log.Fatal(e)
	}

	for {
		if _, e := fmt.Fprint(os.Stdout, "\n$ "); e != nil {
			log.Fatal(e)
		}

		var (
			line []byte
			e    error
		)

		if line, _, e = stdin.ReadLine(); e != nil || len(line) == 0 {
			break
		}

		v, e := cfg.Variable(string(line))

		if e != nil {
			if neighbors := neighbor.FindNearest(string(line), cfg.Keys(), numSuggestions); neighbors != nil {
				if len(neighbors) > 1 && strings.ToLower(neighbors[0].Name) != strings.ToLower(string(line)) {
					printSuggestions(comment, neighbors)
					time.Sleep(time.Millisecond * 100)
					continue
				}

				if v, e = cfg.Variable(neighbors[0].Name); e != nil {
					log.Fatal(e)
				}
			}
		}

		if v == nil {
			if e != nil {
				if _, e := fmt.Fprintln(os.Stderr, e.Error()); e != nil {
					log.Fatal(e)
				}
			}

			continue
		}

		if _, e := info.Fprintf(os.Stdout, "%s\n    %T ", v.Key, v.Value); e != nil {
			log.Fatal(e)
		}

		if _, e := success.Fprintf(os.Stdout, "%+v\n", v.Value); e != nil {
			log.Fatal(e)
		}
	}
}
