package main

import (
	"log"

	"github.com/pkg/errors"
)

func main() {
	bc := NewBlockchain()
	defer bc.db.Close()

	communication := COMMUNICATION{bc}
	err := communication.server()

	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}
}
