package blockchain

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/apex/log"
)

type Block struct {
	Timestamp    int64
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
	Height       int
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

func CreateBlock(txs []*Transaction, prevHash []byte, height, powDifficulty int) *Block {
	block := &Block{time.Now().Unix(), []byte{}, txs, prevHash, 0, height}
	pow := NewProof(block, powDifficulty)
	startTime := time.Now()
	nonce, hash := pow.Run()
	log.Infof("took %v run prood of work of a new block, nonce is %d", time.Since(startTime), nonce)

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis(coinbase *Transaction, powDifficulty int) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{}, 0, powDifficulty)
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Fatalf("panic on error %v", err)
	}
}
