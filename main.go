package main

import (
	"log"

	"github.com/pkg/errors"
)

func main() {
	bc := NewBlockchain()
	defer bc.db.Close()

	// cli := CLI{bc}
	// cli.Run()

	// communication := COMMUNICATION{bc}
	err := server()
	// err := server()

	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
}
