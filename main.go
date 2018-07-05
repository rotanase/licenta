package main

import (
	"log"

	"github.com/pkg/errors"
)

func main() {

	err := server()

	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
}
