/*
TODO:
	> Write unit tests
	> Document symbols
	> Convert type of Request.Suggestions:
		> uint -> uint8
			255 suggestions is enough
	> Add additional commands to manage environment variables
		> set
		> unset
		> export
		> import
	> Rename KeyParallel -> CtxParallel
*/
package search

import (
	"context"
	"janmarten.name/env/neighbor"
)

type requestChan chan *Request
type responseChan chan *Response

const KeyParallel contextKey = "parallel"

type Response struct {
	Match       *interface{}
	Suggestions []string
	Request     *Request
}

type Request struct {
	Query       string
	Suggestions uint
}

type Engine interface {
	Query(query string, suggestions uint) Engine
	QueryAll(queries []string, suggestions uint) Engine
	Result() *Response
	Results() []*Response
}

type searchEngine struct {
	Engine
	ctx       context.Context
	targets   map[string]*interface{}
	requests  requestChan
	responses responseChan
}

type contextKey string

func New(ctx context.Context, targets map[string]*interface{}) Engine {
	numParallel, ok := ctx.Value(KeyParallel).(uint)

	if ok == false || numParallel < 1 {
		numParallel = 1
	}

	engine := &searchEngine{
		ctx:       ctx,
		targets:   targets,
		requests:  make(requestChan, numParallel),
		responses: make(responseChan, numParallel),
	}

	return engine
}

func (engine searchEngine) Query(query string, suggestions uint) Engine {
	select {
	case <-engine.ctx.Done():
		break
	default:
		if suggestions < 0 {
			suggestions = 0
		}

		engine.requests <- &Request{
			Query:       query,
			Suggestions: suggestions,
		}

		go engine.processNextRequest()
	}

	return &engine
}

func (engine searchEngine) QueryAll(queries []string, suggestions uint) Engine {
	for _, q := range queries {
		engine.Query(q, suggestions)
	}

	return &engine
}

func (engine searchEngine) processNextRequest() {
	select {
	case <-engine.ctx.Done():
		break
	default:
		request := <-engine.requests
		response := &Response{
			Match:       engine.targets[request.Query],
			Suggestions: nil,
			Request:     request,
		}

		if response.Match == nil && request.Suggestions > 0 {
			response.Suggestions = append(response.Suggestions, engine.suggestions(request)...)
		}

		engine.responses <- response
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

	for k, _ := range engine.targets {
		keys = append(keys, k)
	}

	return keys
}

func (engine searchEngine) Result() *Response {
	return <-engine.responses
}

func (engine searchEngine) Results() []*Response {
	responses := make([]*Response, 0)

	for len(engine.requests) > 0 || len(engine.responses) > 0 {
		responses = append(responses, <-engine.responses)
	}

	return responses
}
