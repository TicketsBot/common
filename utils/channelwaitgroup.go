package utils

import "sync"

type ChannelWaitGroup struct {
	internal sync.WaitGroup
}

func NewChannelWaitGroup() ChannelWaitGroup {
	return ChannelWaitGroup{}
}

func (wg *ChannelWaitGroup) Add(delta int) {
	wg.internal.Add(delta)
}

func (wg *ChannelWaitGroup) Done() {
	wg.internal.Done()
}

func (wg *ChannelWaitGroup) Wait() chan struct{} {
	ch := make(chan struct{})

	go func() {
		wg.internal.Wait()
		ch <- struct{}{}
	}()

	return ch
}
