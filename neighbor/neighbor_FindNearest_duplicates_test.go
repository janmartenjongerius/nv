package neighbor

import "testing"

func TestFindNearest_duplicates(t *testing.T) {
	neighbors := FindNearest("foo", []string{"Foo", "FO", "Foo"}, 3)
	want := 2
	got := len(neighbors)

	if want != got {
		t.Errorf("Expected len(neighbors) == %d, got %d", want, got)
	}
}
