package blockchain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProofOfWorkParallelismResult(t *testing.T) {
	address := "18kn66B4vJjg7R27RX68LmbABeYrbx8kbr"
	difficulty := 8
	cbtx := CoinbaseTx(address, genesisData)

	block := &Block{time.Now().Unix(), []byte{}, []*Transaction{cbtx}, []byte{}, 0, difficulty}
	pow := NewProof(block, difficulty)
	singleThreadedNonce, singleThreadedHash := pow.RunSingleThread()
	parallelNonce, parallelHash := pow.Run()
	assert.Equal(t, singleThreadedNonce, parallelNonce)
	assert.Equal(t, singleThreadedHash, parallelHash)
}
