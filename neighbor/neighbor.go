package neighbor

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sort"
	"strings"
)

// A struct representing a neighbor, with a Distance to its relation.
type Neighbor struct {
	Name     string
	Distance int
}

// A slice of Neighbor objects.
type Neighbors []*Neighbor

// The number of Neighbor object in the current slice.
func (n Neighbors) Len() int {
	return len(n)
}

// Swap Neighbor at offset i with Neighbor at offset j.
func (n Neighbors) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// Determine if the Neighbor at offset i has a lesser Neighbor.Distance than Neighbor at offset j.
func (n Neighbors) Less(i, j int) bool {
	return n[i].Distance < n[j].Distance
}

// Find the nearest, up to <limit>, Neighbors for the given subject.
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
		if len(rank.Target) == 0 {
			continue
		}

		if _, ok := ranked[rank.Target]; ok {
			continue
		}

		ranked[rank.Target] = rank

		for _, name := range mapping[rank.Target] {
			buffer = append(
				buffer,
				&Neighbor{name, rank.Distance})
		}
	}

	sort.Sort(buffer)

	if len(buffer) > limit {
		buffer = buffer[0:limit]
	}

	return buffer
}
