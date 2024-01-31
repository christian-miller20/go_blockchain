package main

import (
	"crypto/sha256"
	"fmt"
	"encoding/json"
	"strconv"
	"time"
	"strings"
)

// define underlying structures of blockchain
type Block struct {
	data map[string]interface{}
	hash string
	previousHash string
	timestamp time.Time
	pow int
}

type Blockchain struct {
	genesisBlock Block
	chain []Block
	difficulty int
}

// calculate hash of a block
func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.data)
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

// mine a block, becomes more difficult as difficulty increases
func (b *Block) mineBlock(difficulty int) {
	// iterate until valid block is found
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}
}

// create empty genesis block to begin blockchain
func createBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hash: "0",
		timestamp: time.Now(),
	}
	return Blockchain{
		genesisBlock, 
		[]Block{genesisBlock},
		difficulty,
	}
}

// add a new block to the chain
func (b *Blockchain) addBlock(data map[string]interface{}) {
	previousBlock := b.chain[len(b.chain)-1] // get last block
	newBlock := Block{
		data: data,
		previousHash: previousBlock.hash,
		timestamp: time.Now(),
	}
	newBlock.mineBlock(b.difficulty) // generate hash + pow via mining
	b.chain = append(b.chain, newBlock) // add to chain
}

// check if blockchain is valid by checking all blocks and comparing hashes
func (b Blockchain) isValid() bool {
	for i := range b.chain[1:]{
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}

func main() {
	blockchain := createBlockchain(2)

	blockchain.addBlock(map[string]interface{}{"from": "A", "to": "B", "code": "0x938018"})
	blockchain.addBlock(map[string]interface{}{"from": "A", "to": "B", "amount": 10})

	fmt.Println(blockchain.isValid())
	for _, block := range blockchain.chain {
		fmt.Println(block)
	}
}