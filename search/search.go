package search

import (
	"context"
	"janmarten.name/env/neighbor"
	"reflect"
)

type requestChan chan *Request
type responseChan chan *Response

type Response struct {
	Match       *interface{}
	Suggestions []string
	Request     *Request
}

type Request struct {
	Query       string
	Suggestions int
}

type Engine struct {
	ctx       context.Context
	targets   map[string]*interface{}
	requests  requestChan
	responses responseChan
}

type ContextKey string

func NewEngine(ctx context.Context, targets map[string]*interface{}) *Engine {
	numParallel, ok := ctx.Value(ContextKey("parallel")).(int)

	if ok == false || numParallel < 1 {
		numParallel = 1
	}

	engine := &Engine{
		ctx:       ctx,
		targets:   targets,
		requests:  make(requestChan, numParallel),
		responses: make(responseChan, numParallel),
	}

	return engine
}

func (engine Engine) Query(query string, suggestions int) *Engine {
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

func (engine Engine) processNextRequest() {
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

func (engine Engine) suggestions(request *Request) []string {
	var suggestions []string

	if neighbors := neighbor.FindNearest(request.Query, engine.targetKeys(), request.Suggestions); neighbors != nil {
		for _, n := range neighbors {
			suggestions = append(suggestions, n.Name)
		}
	}

	return suggestions
}

func (engine Engine) targetKeys() []string {
	var keys []string

	for _, k := range reflect.ValueOf(engine.targets).MapKeys() {
		keys = append(keys, k.String())
	}

	return keys
}

func (engine Engine) Result() *Response {
	return <-engine.responses
}
