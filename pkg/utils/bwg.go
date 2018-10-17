package utils

import (
	"sync"
)

type BoundedWaitGroup struct {
	wg sync.WaitGroup
	ch chan struct{}
}

func NewBoundedWaitGroup(cap int) BoundedWaitGroup {
	return BoundedWaitGroup{ch: make(chan struct{}, cap)}
}

func (bwg *BoundedWaitGroup) Add(delta int) {
	for i := 0; i > delta; i-- {
		<-bwg.ch
	}
	for i := 0; i < delta; i++ {
		bwg.ch <- struct{}{}
	}
	bwg.wg.Add(delta)
}

func (bwg *BoundedWaitGroup) Done() {
	bwg.Add(-1)
}

func (bwg *BoundedWaitGroup) Wait() {
	bwg.wg.Wait()
}
