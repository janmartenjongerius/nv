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

type Response struct {
	Match       *config.Variable
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
	ctx        context.Context
	targets    map[string]*config.Variable
	requests   requestChan
	responses  responseChan
	processing chan bool
}

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
		engine.processing <- true
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

func (engine searchEngine) Result() *Response {
	return <-engine.responses
}

func (engine searchEngine) Results() []*Response {
	responses := make([]*Response, 0)

	for engine.busy() {
		responses = append(responses, <-engine.responses)
	}

	return responses
}

func (engine searchEngine) busy() bool {
	return len(engine.requests) > 0 || len(engine.responses) > 0 || len(engine.processing) > 0
}

type Service struct {
	Suggestions uint
	Targets     config.Variables
}

func NewService(variables config.Variables) Service {
	return Service{Targets: variables}
}

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
