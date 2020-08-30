package debug

import "sort"

var callbacks = make(map[Scope]Callback)

// Messages is a map of debug messages per group.
type Messages map[string]interface{}

// Callback describes a func that is called when the debug command runs.
type Callback func() Messages

// Scope identifies the origin of debugging information.
type Scope string

// GetMessages invokes the callback registered for the current scope and returns
// the resulting Messages.
func (s Scope) GetMessages() Messages {
	return callbacks[s]()
}

// RegisterCallback will register a debugging callback against the given scope.
// Calling this method multiple times with the same scope will override any
// previous registrations.
func RegisterCallback(scope Scope, cb Callback) {
	callbacks[scope] = cb
}

// ScopeWalker is a func that receives a call for each debug scope.
type ScopeWalker func(scope Scope) MessageWalker

// MessageWalker is a func that receives a call for each message and its group.
type MessageWalker func(group string, message interface{})

// Walk is a func that can be invoked with a walker, to walk through all debug
// messages.
func Walk(scopeWalker ScopeWalker) {
	scopes := make([]string, 0)

	for scope := range callbacks {
		scopes = append(scopes, string(scope))
	}

	sort.Strings(scopes)

	for _, scope := range scopes {
		scope := Scope(scope)
		result := callbacks[scope]()
		groups := make([]string, 0)
		messageWalker := scopeWalker(scope)

		for group := range result {
			groups = append(groups, group)
		}

		sort.Strings(groups)

		for _, group := range groups {
			messageWalker(group, result[group])
		}
	}
}

func init() {
	RegisterCallback("Debug", func() Messages {
		return Messages{
			"Callbacks": func(callbacks map[Scope]Callback) []string {
				result := make([]string, 0)

				for scope := range callbacks {
					result = append(result, string(scope))
				}

				sort.Strings(result)

				return result
			}(callbacks),
		}
	})
}
