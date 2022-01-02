package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requirements:
// The First few bytes must contain 0s

type ProofOfWork struct {
	Block      *Block
	Target     *big.Int
	difficulty int
}

func NewProof(b *Block, difficulty int) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	pow := &ProofOfWork{b, target, difficulty}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(pow.difficulty)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) RunSingleThread() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	for nonce := 0; nonce < math.MaxInt64; nonce++ {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			fmt.Printf("complete with nonce %d", nonce)
			return nonce, hash[:]
		}
	}
	return -1, nil
}

func (pow *ProofOfWork) Run() (int, []byte) {
	threadPool := newPowPool()
	return threadPool.run(pow.RunRange)
}

func (pow *ProofOfWork) RunRange(start, end int) (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	for nonce := start; nonce < end; nonce++ {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			return nonce, hash[:]
		}
	}
	return -1, nil
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)

	}

	return buff.Bytes()
}
