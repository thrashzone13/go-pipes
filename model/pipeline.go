package model

import "sync"

type Pipeline struct {
	Pipes           []*Pipe
	mu              sync.Mutex
	executedPipeIDs map[int]bool
}

func NewPipeline(pipes ...*Pipe) *Pipeline {
	return &Pipeline{
		Pipes:           pipes,
		executedPipeIDs: make(map[int]bool),
	}
}

func (p *Pipeline) Pipe(pipe *Pipe) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Pipes = append(p.Pipes, pipe)
}

func (p *Pipeline) Process(payload interface{}) {
	var wg sync.WaitGroup

	for _, pipe := range p.Pipes {
		wg.Add(1)
		go func(pipe *Pipe) {
			defer wg.Done()
			p.execute(pipe, payload)
		}(pipe)
	}

	wg.Wait()
}

func (p *Pipeline) execute(pipe *Pipe, payload interface{}) {
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

func (p *Pipeline) findPipeById(id int) *Pipe {
	for _, pipe := range p.Pipes {
		if pipe.ID == id {
			return pipe
		}
	}

	return nil
}

func (p *Pipeline) Reset() {
	p.mu.Lock()
	p.executedPipeIDs = make(map[int]bool)
	p.mu.Unlock()
}
