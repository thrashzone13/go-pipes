package main

import (
	"errors"
	"fmt"

	"github.com/thrashzone13/go-pipes/model"
)

func main() {
	pipe1 := model.NewPipe(func(payload interface{}) (interface{}, error) {
		return fmt.Print("First pipeline")
	})

	pipe2 := model.NewPipe(func(payload interface{}) (interface{}, error) {
		if str, ok := payload.(string); ok {
			return fmt.Println(str)
		} else {
			return payload, errors.New("Value should be string")
		}
	})

	pipeline := model.NewPipeline(pipe1, pipe2)
	pipeline.Process("hello")
}
