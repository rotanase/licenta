package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"

	"github.com/pkg/errors"
)

type transaction struct {
	OperationType int
	DoctorHash    []byte
	PacientHash   []byte
	Data          []byte
	Signature     []byte
}

// listening port for incoming requests
const port = ":4400"
const localNode = "127.0.0.1:4500"
const remoteNode1 = "127.0.0.1:4600"
const remoteNode2 = "127.0.0.1:4700"

const addRecord = 0
const getRecord = 1
const getRecords = 2

func handleAddRecord(data transaction) error {

	nodeConn, err := net.Dial("tcp", localNode)
	nodeConn1, err1 := net.Dial("tcp", remoteNode1)
	nodeConn2, err2 := net.Dial("tcp", remoteNode2)

	defer nodeConn.Close()
	defer nodeConn1.Close()
	defer nodeConn2.Close()

	if err != nil {
		log.Println("Dial problem at local node, aborting...")
		return err
	}

	if err1 != nil {
		log.Println("Dial problem at node 1, aborting...")
		return err1
	}

	if err2 != nil {
		log.Println("Dial problem at node 2, aborting...")
		return err2
	}

	rw := bufio.NewReadWriter(bufio.NewReader(nodeConn), bufio.NewWriter(nodeConn))
	enc := gob.NewEncoder(rw)
	err = enc.Encode(data)

	rw1 := bufio.NewReadWriter(bufio.NewReader(nodeConn1), bufio.NewWriter(nodeConn1))
	enc1 := gob.NewEncoder(rw1)
	err1 = enc1.Encode(data)

	rw2 := bufio.NewReadWriter(bufio.NewReader(nodeConn2), bufio.NewWriter(nodeConn2))
	enc2 := gob.NewEncoder(rw2)
	err2 = enc2.Encode(data)

	if err != nil {
		log.Println("Failed to encode the gob")
		return err
	}

	if err1 != nil {
		log.Println("Failed to encode the gob")
		return err1
	}

	if err2 != nil {
		log.Println("Failed to encode the gob")
		return err2
	}

	err = rw.Flush()
	err1 = rw1.Flush()
	err2 = rw2.Flush()

	if err != nil {
		log.Println("Failed to flush (send) to the local node")
		return err
	}

	if err1 != nil {
		log.Println("Failed to flush (send) to the first remote node")
		return err1
	}

	if err2 != nil {
		log.Println("Failed to flush (send) to the second remote node")
		return err2
	}

	log.Println("Sended the data to all nodes")

	return nil
}

func handleGetRecord(data transaction) transaction {

	nodeConn, err := net.Dial("tcp", localNode)

	defer nodeConn.Close()

	if err != nil {
		log.Println("Dial problem, aborting...")
		return transaction{}
	}

	rw := bufio.NewReadWriter(bufio.NewReader(nodeConn), bufio.NewWriter(nodeConn))
	enc := gob.NewEncoder(rw)
	err = enc.Encode(data)

	if err != nil {
		log.Println("Failed to encode the gob")
		return transaction{}
	}

	err = rw.Flush()

	if err != nil {
		log.Println("Failed to flush (send)")
		return transaction{}
	}

	var receive transaction

	dec := gob.NewDecoder(rw)
	err = dec.Decode(&receive)

	log.Println("Received the data from the local node")
	return receive
}

func handleGetRecords(data transaction) []transaction {

	nodeConn, err := net.Dial("tcp", localNode)

	defer nodeConn.Close()

	if err != nil {
		log.Println("Dial problem, aborting...")
		return []transaction{}
	}

	rw := bufio.NewReadWriter(bufio.NewReader(nodeConn), bufio.NewWriter(nodeConn))
	enc := gob.NewEncoder(rw)
	err = enc.Encode(data)

	if err != nil {
		log.Println("Failed to encode the gob")
		return []transaction{}
	}

	err = rw.Flush()

	if err != nil {
		log.Println("Failed to flush (send)")
		return []transaction{}
	}

	var receive []transaction

	dec := gob.NewDecoder(rw)
	err = dec.Decode(&receive)

	log.Println("Received the data from the local node")
	return receive
}

func handleGOB(rw *bufio.ReadWriter, conn net.Conn) {
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
		log.Println("All the nodes will add this record.")

		err = handleAddRecord(incomingData)

	case incomingData.OperationType == getRecord:
		log.Println("The local node will return the record, based on the given hash.")
		res := handleGetRecord(incomingData)

		enc := gob.NewEncoder(rw)
		err = enc.Encode(res)
		if err != nil {
			log.Printf("Encode failed for struct: %#v", err)
		}
		err = rw.Flush()
		if err != nil {
			log.Printf("Flush failed: %#v", err)
		}

	case incomingData.OperationType == getRecords:
		log.Println("The local node will return all the records.")

		res := handleGetRecords(incomingData)

		enc := gob.NewEncoder(rw)
		err = enc.Encode(res)
		if err != nil {
			log.Printf("Encode failed for struct: %#v", err)
		}
		err = rw.Flush()
		if err != nil {
			log.Printf("Flush failed: %#v", err)
		}
	}
}

func server() error {

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
		go handleGOB(rw, conn)
	}
}
