package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"

	"github.com/pkg/errors"
)

type COMMUNICATION struct {
	bc *Blockchain
}

type transaction struct {
	OperationType string
	DoctorHash    []byte
	PacientHash   []byte
	data          []byte
}

// listening port for incoming wallet requests
const port = ":4400"
const ADD_RECORD = "Add record"
const GET_RECORD = "Get record"
const GET_RECORDS = "Get records"

func handleConnection(conn net.Conn) {
	log.Print("Receive transaction data:")
	var incomingData transaction
	var network bytes.Buffer

	dec := gob.NewDecoder(&network)
	err := dec.Decode(&incomingData)
	if err != nil {
		log.Println("Error decoding transaction data:", err)
		return
	}

	switch {
	case incomingData.OperationType == ADD_RECORD:
		log.Println("Addind a record...\n")

	case incomingData.OperationType == GET_RECORD:
		log.Println("Getting a record...\n")

	case incomingData.OperationType == GET_RECORDS:
		log.Println("Getting all the records...\n")
	}
	log.Printf("Incoming complexData struct, doc: \n%#v\n", incomingData.OperationType)
	log.Printf("Incoming complexData struct, doc: \n%#v\n", string(incomingData.DoctorHash))
	log.Printf("Incoming complexData struct, data: \n%#v\n", string(incomingData.data))
	log.Printf("Incoming complexData struct, pacient: \n%#v\n", string(incomingData.PacientHash))
}

func server() error {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", port)
	}
	log.Println("Listen on", port)

	for {
		log.Println("Accept a connection request.")
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go handleConnection(conn)
	}
}
