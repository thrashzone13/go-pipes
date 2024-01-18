package model

type Pipe struct {
	ID           int
	Function     func(interface{}) (interface{}, error)
	Dependencies []int
	Result       interface{}
	Err          error
}
