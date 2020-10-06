package config

import (
	"janmarten.name/nv/debug"
	"reflect"
	"sort"
	"testing"
)

func TestConfig_init(t *testing.T) {
	got := debug.Scope("nv/config.init.0").GetMessages()
	vars, ok := got["Env"].([]string)

	if !ok {
		t.Fatalf("Unexpected or missing Env key in debug messages: %q", got)
	}

	want := make([]string, 0)

	for _, v := range Environment {
		want = append(want, v.Key)
	}

	sort.Strings(want)

	if !reflect.DeepEqual(want, vars) {
		t.Errorf(
			"Unexpected variables in debug output. Want %v, got %v",
			want,
			vars)
	}
}
