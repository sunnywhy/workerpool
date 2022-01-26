package pool

import (
	"fmt"
	"sync"
)

const (
	defaultCapacity = 10
	maxCapacity     = 50
)

type Pool struct {
	capacity int           // the size of the worker pool
	active   chan struct{} // active channel
	tasks    chan Task     // tasks channel

	wg   sync.WaitGroup // used for waiting all workers exit when destroy the pool
	quit chan struct{}  // used to notify all workers exit
}

type Task func()

func New(capacity int) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
		active:   make(chan struct{}),
	}
	fmt.Printf("workerpool start\n")
	go p.run()
	return p
}

func (p *Pool) run() {
	idx := 0

	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			// create a new worker
			idx++
			p.newWorker(idx)
		}

	}
}

func (p *Pool) newWorker(id int) {
	
}