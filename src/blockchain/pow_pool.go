package blockchain

import (
	"math"
	"runtime"

	"github.com/apex/log"
)

// thread pool of proof of work
type PowPool struct {
	poolSize  int
	batchSize int
}

func newPowPool() *PowPool {
	return &PowPool{
		poolSize:  int(math.Max(float64(runtime.NumCPU()), 4)),
		batchSize: 10000, // maybe this can be adjusted based on pow's difficulty level
	}
}

type nonceHashSig struct {
	nonce int
	hash  []byte
}
type sig struct{}

func (pp *PowPool) run(process func(start, end int) (int, []byte)) (int, []byte) {
	foundChan := make(chan *nonceHashSig, pp.poolSize)
	workerCompleteChan := make(chan sig, pp.poolSize)
	defer close(foundChan)
	defer close(workerCompleteChan)

	log.Infof("proof of work concurrent size %d, batch size %d", pp.poolSize, pp.batchSize)
	var foundObj *nonceHashSig
	nonce := -1
	var data []byte
	nonceCounter := 0
	threadCounter := 0 // no need to be thread safe since only it's updated by this thread
	//
	for i := 0; i < pp.poolSize; i++ {
		go processHelper(nonceCounter, nonceCounter+pp.batchSize, foundChan, workerCompleteChan, process)
		nonceCounter = nonceCounter + pp.batchSize + 1 // next iteration starting point
		threadCounter++
	}

	for threadCounter > 0 {
		select {
		case foundObj = <-foundChan:
			if nonce == -1 || foundObj.nonce < nonce {
				// multiple nonce can be found if the difficulty level is low
				// so we will process all the possible nonce but only take the smallest one
				nonce = foundObj.nonce
				data = foundObj.hash
			}
			threadCounter--
		case <-workerCompleteChan:
			threadCounter--
			if nonceCounter < math.MaxInt64 && nonce < 1 {
				endCounter := nonceCounter + pp.batchSize
				if math.MaxInt64-pp.batchSize < nonceCounter {
					endCounter = math.MaxInt64
				}
				threadCounter++
				go processHelper(nonceCounter, endCounter, foundChan, workerCompleteChan, process)
				nonceCounter = endCounter + 1 // next iteration starting point
			}
		}
	}

	return nonce, data
}

func processHelper(start, end int, successChan chan *nonceHashSig, nextWorkerChan chan sig, process func(start, end int) (int, []byte)) {
	nonce, data := process(start, end)
	if nonce > 1 {
		successChan <- &nonceHashSig{nonce, data}
	} else {
		nextWorkerChan <- sig{}
	}
}
