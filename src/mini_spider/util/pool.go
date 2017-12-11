package util

import (
	"sync"
)

type Pool struct {
	lock     *sync.Mutex
	capacity uint
	used     uint
}

func NewPool(capacity uint) *Pool {
	return &Pool{lock: &sync.Mutex{}, capacity: capacity, used: 0}
}

func (pool *Pool) Allot(count uint) bool {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if pool.used+count <= pool.capacity {
		pool.used += count
		return true
	}

	return false
}

func (pool *Pool) Release(count uint) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	used := pool.used - count
	if used < 0 {
		used = 0
	}
	pool.used = used
}

func (pool *Pool) Used() uint {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	return pool.used
}

func (pool *Pool) Left() uint {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	return pool.capacity - pool.used
}
