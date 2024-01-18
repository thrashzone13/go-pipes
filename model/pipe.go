package model

import "github.com/google/uuid"

type pipe struct {
	ID           uuid.UUID
	Function     func(interface{}) (interface{}, error)
	Dependencies []uuid.UUID
	Result       interface{}
	Err          error
}

func NewPipe(function func(payload interface{}) (interface{}, error)) *pipe {
	return &pipe{
		ID:       uuid.New(),
		Function: function,
	}
}

func (p *pipe) AddDependency(id uuid.UUID) {
	p.Dependencies = append(p.Dependencies, id)
}
