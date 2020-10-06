/*
Package debug exposes hooks and tools to provide and read out debugging
information about the application.
*/
package debug

import (
	"runtime"
	"sort"
	"strings"
)

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
	if cb, ok := callbacks[s]; ok == true {
		return cb()
	}

	return Messages{}
}

// newScope returns a new Scope referencing the scope in which newScope was called.
func newScope(skip int) *Scope {
	if pc, _, _, ok := runtime.Caller(1 + skip); ok == true {
		frames := runtime.CallersFrames([]uintptr{pc})

		for {
			frame, _ := frames.Next()
			scope := Scope(
				strings.TrimPrefix(frame.Function, "janmarten.name/"))

			return &scope
		}
	}

	return nil
}

// RegisterCallback will register a debugging callback against the current scope.
func RegisterCallback(cb Callback) Scope {
	scope := *newScope(1)

	callbacks[scope] = cb

	return scope
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
	RegisterCallback(func() Messages {
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
