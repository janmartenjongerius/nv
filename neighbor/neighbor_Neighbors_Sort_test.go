package neighbor

import (
	"sort"
	"testing"
)

func TestNeighbors_Sort(t *testing.T) {
	cases := []struct{
		in Neighbors
		want Neighbors
	}{
		{
			in: Neighbors{},
			want: Neighbors{},
		},
		{
			in: Neighbors{
				&Neighbor{
					Name:     "Foo",
					Distance: 0,
				},
				&Neighbor{
					Name:     "Bar",
					Distance: 15,
				},
				&Neighbor{
					Name:     "Baz",
					Distance: 2,
				},
			},
			want: Neighbors{
				&Neighbor{
					Name:     "Foo",
					Distance: 0,
				},
				&Neighbor{
					Name:     "Baz",
					Distance: 2,
				},
				&Neighbor{
					Name:     "Bar",
					Distance: 15,
				},
			},
		},
	}

	for _, c := range cases {
		got := make(Neighbors, len(c.in))
		copy(got, c.in)

		sort.Sort(got)

		for n := range c.want {
			if got[n].Name != c.want[n].Name {
				t.Errorf("Neighbor at position %d is expected to be %s, got %s",
					n,
					c.want[n].Name,
					got[n].Name)
			}
		}
	}
}
