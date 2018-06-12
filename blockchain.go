package main

import (
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

// we should take a look at err
// we should log every operation

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

/*

   Open a DB file.
   Check if there’s a blockchain stored in it.
   If there’s a blockchain:
       Create a new Blockchain instance.
       Set the tip of the Blockchain instance to the last block hash stored in the DB.
   If there’s no existing blockchain:
       Create the genesis block.
       Store in the DB.
       Save the genesis block’s hash as the last block hash.
       Create a new Blockchain instance with its tip pointing at the genesis block.
*/
func NewBlockchain() *Blockchain {
	// return &Blockchain{[]*Block{NewGenesisBlock()}}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Println("Error at opening the database.")
	}
	// running a write-read DB transaction (adding the genesis block)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Println("Error at creating a bucket.")
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("lastHash"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("lastHash"))
		}

		return nil
	})

	// we store only the tip of the chain
	bc := Blockchain{tip, db}

	return &bc
}

func (bc *Blockchain) AddBlock(data []byte, doctorHash []byte, pacientHash []byte, signature []byte) {
	var lastHash []byte

	// running a read DB transaction (read the last block)
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("lastHash"))

		return nil
	})

	if err != nil {
		log.Println("Error at reading from database.")
	}

	newBlock := NewBlock(data, doctorHash, pacientHash, lastHash, signature)

	// update the blockchain with the new block
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Println("Error at adding a new block to database.")
		}
		err = b.Put([]byte("lastHash"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Println("Error reading from the blockchain")
	}

	i.currentHash = block.PrevBlockHash

	return block
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	// starting from the tip
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}
