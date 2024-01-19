package model

import (
	"sync"

	"github.com/google/uuid"
)

type Pipeline struct {
	Pipes           []*pipe
	mu              sync.Mutex
	executedPipeIDs map[uuid.UUID]bool
}

func NewPipeline(pipes ...*pipe) *Pipeline {
	return &Pipeline{
		Pipes:           pipes,
		executedPipeIDs: make(map[uuid.UUID]bool),
	}
}

func (p *Pipeline) Pipe(pipe *pipe) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Pipes = append(p.Pipes, pipe)
}

func (p *Pipeline) Process(payload interface{}) {
	var wg sync.WaitGroup

	for _, pi := range p.Pipes {
		wg.Add(1)
		go func(pipe *pipe) {
			defer wg.Done()
			p.execute(pipe, payload)
		}(pi)
	}

	wg.Wait()
}

func (p *Pipeline) execute(pipe *pipe, payload interface{}) (interface{}, error) {
	p.mu.Lock()
	if _, executed := p.executedPipeIDs[pipe.ID]; executed {
		p.mu.Unlock()
		return pipe.Result, pipe.Err
	}

	p.executedPipeIDs[pipe.ID] = true
	p.mu.Unlock()

	for _, dep := range pipe.Dependencies {
		result, _ := p.execute(dep, payload)
		payload = result
	}

	result, err := pipe.Function(payload)
	pipe.Result = result
	pipe.Err = err

	return payload, err
}

func (p *Pipeline) Reset() {
	p.mu.Lock()
	p.executedPipeIDs = make(map[uuid.UUID]bool)
	p.mu.Unlock()
}
