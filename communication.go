package main

import (
	"bufio"
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
	OperationType int
	DoctorHash    []byte
	PacientHash   []byte
	Data          []byte
	Signature     []byte
}

// listening port for incoming wallet requests
const port = ":4400"
const addRecord = 0
const getRecord = 1
const getRecords = 2

func (com *COMMUNICATION) handleGOB(rw *bufio.ReadWriter, conn net.Conn) {
	defer func() {
		log.Printf("closing connection from %v", conn.RemoteAddr())
		conn.Close()
	}()

	log.Println("Receive transaction data from " + conn.RemoteAddr().String())
	var incomingData transaction

	dec := gob.NewDecoder(rw)
	err := dec.Decode(&incomingData)
	if err != nil {
		log.Println("Error decoding transaction data:", err)
		return
	}

	switch {
	case incomingData.OperationType == addRecord:
		log.Println("Adding a record...")
		com.bc.AddBlock(incomingData.Data, incomingData.DoctorHash, incomingData.PacientHash, incomingData.Signature)
	case incomingData.OperationType == getRecord:
		log.Println("Getting a record by hash...")
		bci := com.bc.Iterator()
		enc := gob.NewEncoder(rw)
		var records []transaction
		var record transaction

		for {
			block := bci.Next()
			if bytes.Compare(block.PacientHash, incomingData.PacientHash) == 0 {
				record = transaction{incomingData.OperationType, block.DoctorHash,
					block.PacientHash, block.Data, block.Signature}
			}

			records = append(records, record)

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}

		err = enc.Encode(records)
		if err != nil {
			log.Printf("Encode failed for struct: %#v", err)
		}
		err = rw.Flush()
		if err != nil {
			log.Printf("Flush failed: %#v", err)
		}

	case incomingData.OperationType == getRecords:
		log.Println("Getting all the records...")
		bci := com.bc.Iterator()
		enc := gob.NewEncoder(rw)
		var records []transaction

		for {
			block := bci.Next()
			record := transaction{incomingData.OperationType, block.DoctorHash,
				block.PacientHash, block.Data, block.Signature}
			records = append(records, record)

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}

		err = enc.Encode(records)
		if err != nil {
			log.Printf("Encode failed for struct: %#v", err)
		}
		err = rw.Flush()
		if err != nil {
			log.Printf("Flush failed: %#v", err)
		}
		// conn.Write(byte[] (records))
	}
}

func (com *COMMUNICATION) server() error {

	listener, err := net.Listen("tcp", port)

	// defer listener.Close()

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

		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		go com.handleGOB(rw, conn)
	}
}
