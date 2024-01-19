package model

import "github.com/google/uuid"

type pipe struct {
	ID           uuid.UUID
	Function     func(interface{}) (interface{}, error)
	Dependencies []*pipe
	Result       interface{}
	Err          error
}

func NewPipe(function func(payload interface{}) (interface{}, error)) *pipe {
	return &pipe{
		ID:       uuid.New(),
		Function: function,
	}
}

func (p *pipe) AddDependency(pipe *pipe) {
	p.Dependencies = append(p.Dependencies, pipe)
}
