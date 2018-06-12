package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

// Block keeps block headers
type Block struct {
	Timestamp     int64
	DoctorHash    []byte
	PacientHash   []byte
	Data          []byte
	Signature     []byte
	PrevBlockHash []byte
	Hash          []byte
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func NewBlock(data []byte, doctorHash []byte, pacientHash []byte, prevBlockHash []byte, signature []byte) *Block {
	block := &Block{time.Now().Unix(), doctorHash, pacientHash, data, signature, prevBlockHash, []byte{}}
	block.SetHash()

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock([]byte("Genesis Block"), nil, nil, nil, []byte{})
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
