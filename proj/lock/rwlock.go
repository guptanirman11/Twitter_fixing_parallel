// Package lock provides an implementation of a read-write lock
// that uses condition variables and mutexes.
package lock

import "sync"

type ReadWriteLock struct {
	mu         sync.Mutex
	condition  *sync.Cond
	readCount  int
	writeCount int
	maxReaders int
}

func NewReadWriteLock() *ReadWriteLock {
	return NewReadWriteLockWithMaxReaders(32)
}

func NewReadWriteLockWithMaxReaders(maxReaders int) *ReadWriteLock {
	rwLock := &ReadWriteLock{
		maxReaders: maxReaders,
	}
	rwLock.condition = sync.NewCond(&rwLock.mu)
	return rwLock
}

func (rw *ReadWriteLock) Lock() {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	for rw.readCount > 0 || rw.writeCount > 0 {
		rw.condition.Wait()
	}
	rw.writeCount++
}

func (rw *ReadWriteLock) Unlock() {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	rw.writeCount--

	rw.condition.Broadcast()

}

func (rw *ReadWriteLock) RLock() {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	for rw.writeCount > 0 || rw.readCount >= rw.maxReaders {
		rw.condition.Wait()
	}
	rw.readCount++
}

func (rw *ReadWriteLock) RUnlock() {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	rw.readCount--

	if rw.readCount == 0 {
		rw.condition.Broadcast()
	}
}
