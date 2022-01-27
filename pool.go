package pool

import (
	"errors"
	"fmt"
	"sync"
)

const (
	defaultCapacity = 100
	maxCapacity     = 5000
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
		active:   make(chan struct{}, capacity),
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
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker [%03d]: recover panic[%s] and exit\n", id, err)
				<-p.active
			}
			p.wg.Done()
		}()
		fmt.Printf("worker [%03d]: start\n", id)

		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]:exit\n", id)
				<-p.active
				return
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", id)
				t()
			}
		}
	}()
}

var ErrorWorkerPoolFreed = errors.New("worker pool freed")

func (p * Pool) Schedule(t Task) error {
	select {
		case <-p.quit:
			return ErrorWorkerPoolFreed
		case p.tasks <- t:
			return nil
	}
}

func (p * Pool) Free() {
	close(p.quit)
	p.wg.Wait()
	fmt.Printf("worker pool freed\n")
}