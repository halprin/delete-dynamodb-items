package parallel

import "sync"

type Pool struct {
	ingestionPoolChannel  chan func()
	executionPoolChannel  chan func()
	shutdownSignalChannel chan bool
	executorsAreStopped   sync.WaitGroup
}

//taskQueueSize needs to be bigger than poolSize if you want to saturate the pool
func NewPool(poolSize int, taskQueueSize int) *Pool {
	newPool := &Pool{
		ingestionPoolChannel:  make(chan func(), taskQueueSize),
		executionPoolChannel:  make(chan func(), poolSize),
		shutdownSignalChannel: make(chan bool, 1), //buffered so the Release method doesn't block when sending the shutdown signal
		executorsAreStopped:   sync.WaitGroup{},
	}

	newPool.executorsAreStopped.Add(poolSize)
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
	pool.shutdownSignalChannel <- true
	close(pool.ingestionPoolChannel)
	close(pool.shutdownSignalChannel)
	pool.executorsAreStopped.Wait() //wait for all the executors to finish their current task
}

func (pool *Pool) submitIngestionGoroutine() {
	for submittedTask := range pool.ingestionPoolChannel {
		pool.executionPoolChannel <- submittedTask

		//check for the shutdown signal instead of waiting to drain the ingestionPoolChannel because there could be many long running tasks remaining
		if pool.shutdownSignalRequested() {
			break
		}
	}

	close(pool.executionPoolChannel) //this goroutine sends on the executionPoolChannel, so it is in charge of it, and so it closes it
}

func (pool *Pool) submitExecutionGoroutine() {
	for submittedTask := range pool.executionPoolChannel {
		submittedTask()
	}
	pool.executorsAreStopped.Done()
}

func (pool *Pool) shutdownSignalRequested() bool {
	shutdownRequested := false

	select {
	case <-pool.shutdownSignalChannel:
		shutdownRequested = true
	default:
		shutdownRequested = false
	}

	return shutdownRequested
}
