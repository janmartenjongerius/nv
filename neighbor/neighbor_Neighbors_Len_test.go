package neighbor

import "testing"

func TestNeighbors_Len(t *testing.T) {
	cases := []struct {
		hood *Neighbors
		expected int
	}{
		{
			hood: &Neighbors{},
			expected: 0,
		},
		{
			hood: &Neighbors{
				&Neighbor{
					Name:     "Foo",
					Distance: 0,
				},
			},
			expected: 1,
		},
		{
			hood: &Neighbors{
				&Neighbor{
					Name:     "Foo",
					Distance: 0,
				},
				&Neighbor{
					Name:     "Bar",
					Distance: 1,
				},
			},
			expected: 2,
		},
	}

	for _, c := range cases {
		got := len(*c.hood)
		if got != c.expected {
			t.Errorf("Len(%#v) == %v, expected %v", c.hood, got, c.expected)
		}
	}
}
