package main

import (
	"github.com/llitto372/todo-app"
	"log"
	"todo"
)

func main() {
	srv := new(todo.Server)
	if err := srv.Run("8000"); err != nil {
		log.Fatalf(err)
	}
}
