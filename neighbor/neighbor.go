/*
Package neighbor holds structs to describe the results of FindNearest, which performs nearest neighbor checks against a
provided subject and options.

Each Neighbor holds a Neighbor.Name, which identifies it, and a Neighbor.Distance, describing the distance between the
resulting Neighbor and the original subject.
*/
package neighbor

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"strings"
)

// Neighbor is a struct representing a neighbor, with a Distance to its relation
// and a Name to use as identification.
type Neighbor struct {
	Name     string
	Distance int
}

// Neighbors is a slice of Neighbor objects.
type Neighbors []*Neighbor

// Len provides the number of Neighbor object in the current slice.
func (n Neighbors) Len() int {
	return len(n)
}

// Swap Neighbor at offset i with Neighbor at offset j.
func (n Neighbors) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// Less will determine if the Neighbor at offset i has a lesser
// Neighbor.Distance than Neighbor at offset j.
func (n Neighbors) Less(i, j int) bool {
	return n[i].Distance < n[j].Distance
}

// FindNearest returns the nearest, up to the given limit, Neighbors for the
// given subject.
func FindNearest(subject string, options []string, limit int) Neighbors {
	var (
		buffer  Neighbors
		targets []string
	)

	ranked := map[string]fuzzy.Rank{}
	mapping := map[string][]string{}

	for _, o := range options {
		k := strings.ToLower(o)
		mapping[k] = append(mapping[k], o)

		targets = append(targets, k)
	}

	for _, rank := range fuzzy.RankFind(strings.ToLower(subject), targets) {
		if _, ok := ranked[rank.Target]; ok {
			continue
		}

		ranked[rank.Target] = rank

		for _, name := range mapping[rank.Target] {
			buffer = append(
				buffer,
				&Neighbor{name, rank.Distance},
			)
		}
	}

	sort.Sort(buffer)

	if len(buffer) > limit {
		buffer = buffer[0:limit]
	}

	return buffer
}
