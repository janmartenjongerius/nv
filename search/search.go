/*
Package search holds a search Engine, which is specialized in searching through config.Variables.

It is preferred for the end user to create a new Service using the NewService constructor. This allows the Service to
abstract away the management of the search context and parallelization of the Engine.
 */
package search

import (
	"context"
	"janmarten.name/nv/config"
	"janmarten.name/nv/neighbor"
	"runtime"
)

type requestChan chan *Request
type responseChan chan *Response

// Context key defining the number of parallel queries.
const CtxParallel contextKey = "parallel"

type contextKey string

// The Response of a search operation.
type Response struct {
	Match       *config.Variable
	Suggestions []string
	Request     *Request
}

// A search Request, for the given Query.
type Request struct {
	Query       string
	Suggestions uint
}

// Interface for a search Engine.
type Engine interface {
	Query(query string, suggestions uint) Engine
	QueryAll(queries []string, suggestions uint) Engine
	Results() []*Response
}

type searchEngine struct {
	Engine
	ctx        context.Context
	targets    map[string]*config.Variable
	requests   requestChan
	responses  responseChan
	processing chan bool
}

// Create a new search Engine, to search through the given config.Variables.
func New(ctx context.Context, targets config.Variables) Engine {
	numParallel, ok := ctx.Value(CtxParallel).(uint)

	if ok == false || numParallel < 1 {
		numParallel = 1
	}

	engine := &searchEngine{
		ctx: ctx,
		targets: func(targets config.Variables) map[string]*config.Variable {
			res := make(map[string]*config.Variable)

			for _, t := range targets {
				res[t.Key] = t
			}

			return res
		}(targets),
		requests:   make(requestChan, numParallel),
		responses:  make(responseChan, numParallel),
		processing: make(chan bool, numParallel),
	}

	return engine
}

// Query the config.Variables to find a match for the given query string.
// When no match is found, up to the given number of suggestions is appended to the Response.
func (engine searchEngine) Query(query string, suggestions uint) Engine {
	select {
	case <-engine.ctx.Done():
		break
	default:
		engine.processing <- true
		engine.requests <- &Request{
			Query:       query,
			Suggestions: suggestions,
		}

		go engine.processNextRequest()
	}

	return &engine
}

// Query all the given query strings at-once.
func (engine searchEngine) QueryAll(queries []string, suggestions uint) Engine {
	for _, q := range queries {
		engine.Query(q, suggestions)
	}

	return &engine
}

func (engine searchEngine) processNextRequest() {
	request := <-engine.requests
	response := &Response{
		Match:       engine.targets[request.Query],
		Suggestions: nil,
		Request:     request,
	}

	defer func() {
		<-engine.processing
		engine.responses <- response
	}()

	if response.Match == nil && request.Suggestions > 0 {
		response.Suggestions = append(response.Suggestions, engine.suggestions(request)...)
	}
}

func (engine searchEngine) suggestions(request *Request) []string {
	suggestions := make([]string, 0)
	neighbors := neighbor.FindNearest(
		request.Query,
		engine.targetKeys(),
		int(request.Suggestions),
	)

	if neighbors != nil {
		for _, n := range neighbors {
			suggestions = append(suggestions, n.Name)
		}
	}

	return suggestions
}

func (engine searchEngine) targetKeys() []string {
	keys := make([]string, 0)

	for k := range engine.targets {
		keys = append(keys, k)
	}

	return keys
}

// Return a Response list for pending Request objects and currently processing Response objects, until the engine is drained.
func (engine searchEngine) Results() []*Response {
	responses := make([]*Response, 0)

	for engine.busy() {
		responses = append(responses, <-engine.responses)
	}

	return responses
}

func (engine searchEngine) busy() bool {
	accepted := len(engine.responses) > 0 ||
		len(engine.processing) > 0

	select {
	case <-engine.ctx.Done():
		return accepted
	default:
		return len(engine.requests) > 0 || accepted
	}
}

// Describe a search Service struct.
type Service struct {
	Suggestions uint
	Targets     config.Variables
}

// Create a new search Service for the given config.Variables.
func NewService(variables config.Variables) Service {
	return Service{Targets: variables}
}

// Search the given query strings and wait for the engine to respond with a Response for all queries.
func (s Service) Search(query ...string) []*Response {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx = context.WithValue(ctx, CtxParallel, runtime.GOMAXPROCS(0)*5)

	defer cancel()

	seen := make(map[string]bool)
	engine := New(ctx, s.Targets)

	for _, q := range query {
		if seen[q] {
			continue
		}

		engine.Query(q, s.Suggestions)
		seen[q] = true
	}

	return engine.Results()
}
