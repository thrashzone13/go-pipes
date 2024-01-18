package main

import (
	"fmt"
	"github.com/thrashzone13/go-pipes/model"
)

func main() {
	pipe1 := &model.Pipe{
		ID: 1,
		Function: func() (interface{}, error) {
			return fmt.Print("First pipeline")
		},
	}

	pipe2 := &model.Pipe{
		ID: 2,
		Function: func() (interface{}, error) {
			return fmt.Println("Second pipeline")
		},
	}

	pipeline := model.NewPipeline(pipe1, pipe2)
	pipeline.Process()
}
