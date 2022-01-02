package blockchain

import (
	"math"
	"runtime"
	"sync"

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
	foundChan := make(chan *nonceHashSig)
	workerCompleteChan := make(chan sig, pp.poolSize)
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(foundChan)
		close(workerCompleteChan)
	}()

	log.Infof("proof of work concurrent size %d, batch size %d", pp.poolSize, pp.batchSize)
	var foundObj *nonceHashSig
	nonce := -1
	var data []byte
	nonceCounter := 0
	//
	for i := 0; i < pp.poolSize; i++ {
		go processHelper(nonceCounter, nonceCounter+pp.batchSize, &wg, foundChan, workerCompleteChan, process)
		nonceCounter = nonceCounter + pp.batchSize + 1 // next iteration starting point
	}

	for nonce == -1 {
		select {
		case foundObj = <-foundChan:
			nonce = foundObj.nonce
			data = foundObj.hash
		case <-workerCompleteChan:
			if nonceCounter < math.MaxInt64 {
				endCounter := nonceCounter + pp.batchSize
				if math.MaxInt64-pp.batchSize < nonceCounter {
					endCounter = math.MaxInt64
				}
				go processHelper(nonceCounter, endCounter, &wg, foundChan, workerCompleteChan, process)
				nonceCounter = endCounter + 1 // next iteration starting point

			} else {
				nonce = -2 //break out with nothing found
			}
		}
	}

	return nonce, data
}

func processHelper(start, end int, wg *sync.WaitGroup, successChan chan *nonceHashSig, nextWorkerChan chan sig, process func(start, end int) (int, []byte)) {
	wg.Add(1)
	defer wg.Done()
	nonce, data := process(start, end)
	if nonce > 1 {
		successChan <- &nonceHashSig{nonce, data}
	} else {
		nextWorkerChan <- sig{}
	}
}
