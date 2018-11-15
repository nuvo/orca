package utils

import (
	"sync"
)

// BoundedWaitGroup implements a sized WaitGroup
type BoundedWaitGroup struct {
	wg sync.WaitGroup
	ch chan struct{}
}

// NewBoundedWaitGroup initializes a new BoundedWaitGroup
func NewBoundedWaitGroup(cap int) BoundedWaitGroup {
	return BoundedWaitGroup{ch: make(chan struct{}, cap)}
}

// Add performs a WaitGroup Add of a specified delta
func (bwg *BoundedWaitGroup) Add(delta int) {
	for i := 0; i > delta; i-- {
		<-bwg.ch
	}
	for i := 0; i < delta; i++ {
		bwg.ch <- struct{}{}
	}
	bwg.wg.Add(delta)
}

// Done performs a WaitGroup Add of -1
func (bwg *BoundedWaitGroup) Done() {
	bwg.Add(-1)
}

// Wait performs a WaitGroup Wait
func (bwg *BoundedWaitGroup) Wait() {
	bwg.wg.Wait()
}
