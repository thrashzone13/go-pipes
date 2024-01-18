package main

import (
	"errors"
	"fmt"

	"github.com/thrashzone13/go-pipes/model"
)

func main() {
	pipe1 := &model.Pipe{
		ID: 1,
		Function: func(interface{}) (interface{}, error) {
			return fmt.Print("First pipeline")
		},
	}

	pipe2 := &model.Pipe{
		ID: 2,
		Function: func(payload interface{}) (interface{}, error) {
			if str, ok := payload.(string); ok {
				return fmt.Println(str)
			} else {
				return payload, errors.New("Value should be string")
			}
		},
	}

	pipeline := model.NewPipeline(pipe1, pipe2)
	pipeline.Process("hello")
}
