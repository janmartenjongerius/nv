package debug

import (
	"reflect"
	"testing"
)

func TestWalk(t *testing.T) {
	got := make(map[Scope]map[string][]interface{})
	want := make(map[Scope]map[string][]interface{})

	Walk(func(scope Scope) MessageWalker {
		got[scope] = make(map[string][]interface{})

		return func(group string, message interface{}) {
			got[scope][group] = append(got[scope][group], message)
		}
	})

	scope := Scope("nv/debug.init.0")
	want[scope] = make(map[string][]interface{})
	want[scope]["Callbacks"] = []interface{}{
		[]string{"nv/debug.init.0"},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Unexpected debug entries. Want: %q, got: %q", want, got)
	}
}

func TestRegisterCallback(t *testing.T) {
	scope := RegisterCallback(func() Messages {
		return Messages{t.Name(): true}
	})

	want := Messages{t.Name(): true}
	got := scope.GetMessages()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Expected messages to be equal. Want: %v, Got: %v", want, got)
	}
}

func TestScope_GetMessages(t *testing.T) {
	for scope, cb := range callbacks {
		t.Run(string(scope), func(t *testing.T) {
			want := cb()
			got := scope.GetMessages()

			if !reflect.DeepEqual(want, got) {
				t.Errorf("Unexpected scope messages. Want %v, got %v", want, got)
			}
		})
	}
}

func TestScope_GetMessages_unknown_scope(t *testing.T) {
	want := Messages{}
	got := Scope("unknown").GetMessages()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Unexpected messages. Want %v, got %v", want, got)
	}
}

func TestNewScope(t *testing.T) {
	scope := newScope(200)

	if scope != nil {
		t.Errorf("Expected no scope, got: %v", scope)
	}
}
