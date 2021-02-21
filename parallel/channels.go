package parallel

import "sync"

func MergeErrorChannels(inputChannels []chan error) <-chan error {
	outChannel := make(chan error)
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(inputChannels))

	for _, currentInputChannel := range inputChannels {
		go func(currentInputChannelCopy <-chan error) {
			for valueFromInputChannel := range currentInputChannelCopy {
				outChannel <- valueFromInputChannel
			}
			waitGroup.Done()
		}(currentInputChannel)
	}

	go func() {
		waitGroup.Wait()
		close(outChannel)
	}()

	return outChannel
}
