package parallel

import (
	"sync"
)

type Pool struct {
	ingestionPoolChannel chan func()
	executionPoolChannel chan func()
	waitGroup sync.WaitGroup
}

//taskQueueSize needs to be bigger than poolSize if you want to saturate the pool
func NewPool(poolSize int, taskQueueSize int) *Pool {
	newPool := &Pool{
		ingestionPoolChannel: make(chan func(), taskQueueSize),
		executionPoolChannel: make(chan func(), poolSize),
	}

	go newPool.submitIngestionGoroutine()
	for workerIndex := 0; workerIndex < poolSize; workerIndex++ {
		go newPool.submitExecutionGoroutine()
	}

	return newPool
}

func (pool *Pool) Submit(task func()) {
	pool.ingestionPoolChannel <- task
}

func (pool *Pool) Release() {
	close(pool.ingestionPoolChannel)
	close(pool.executionPoolChannel)
}

func (pool *Pool) submitIngestionGoroutine() {
	for submittedTask := range pool.ingestionPoolChannel {
		pool.executionPoolChannel <- submittedTask
	}
}

func (pool *Pool) submitExecutionGoroutine() {
	for submittedTask := range pool.executionPoolChannel {
		submittedTask()
	}
}