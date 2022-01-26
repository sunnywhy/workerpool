package pool

import "sync"

type Pool struct {
	capacity int // the size of the worker pool
	active chan struct{} // active channel
	tasks chan Task // tasks channel

	wg sync.WaitGroup // used for waiting all workers exit when destroy the pool
	quit chan struct{} // used to notify all workers exit
}

type Task func()