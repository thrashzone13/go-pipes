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

func (p *Pipeline) execute(pipe *pipe, payload interface{}) {
	p.mu.Lock()
	if _, executed := p.executedPipeIDs[pipe.ID]; executed {
		p.mu.Unlock()
		return
	}

	p.executedPipeIDs[pipe.ID] = true
	p.mu.Unlock()

	for _, depId := range pipe.Dependencies {
		dep := p.findPipeById(depId)
		if dep != nil {
			p.execute(dep, payload)
		}
	}

	result, err := pipe.Function(payload)
	pipe.Result = result
	pipe.Err = err
}

func (p *Pipeline) findPipeById(id uuid.UUID) *pipe {
	for _, pipe := range p.Pipes {
		if pipe.ID == id {
			return pipe
		}
	}

	return nil
}

func (p *Pipeline) Reset() {
	p.mu.Lock()
	p.executedPipeIDs = make(map[uuid.UUID]bool)
	p.mu.Unlock()
}
