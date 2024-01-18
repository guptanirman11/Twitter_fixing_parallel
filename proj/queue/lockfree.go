package queue

import (
	"sync/atomic"
	"unsafe"
)

type Node struct {
	task *Request
	next unsafe.Pointer
}

type Request struct {
	Command   string  `json:"command"`
	ID        int     `json:"id"`
	Body      string  `json:"body"`
	Timestamp float64 `json:"timestamp"`
}

// LockfreeQueue represents a FIFO structure with operations to enqueue
// and dequeue tasks represented as Request
type LockFreeQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func new_node() *Node {
	return &Node{next: nil}
}

func initialize(q *LockFreeQueue) {
	node := new_node()
	q.head = unsafe.Pointer(node)
	q.tail = unsafe.Pointer(node)
}

// NewQueue creates and initializes a LockFreeQueue
func NewLockFreeQueue() *LockFreeQueue {
	q := &LockFreeQueue{}
	initialize(q)
	return q
}

// Enqueue adds a series of Request to the queue
func (queue *LockFreeQueue) Enqueue(task *Request) {
	newNode := &Node{task: task, next: nil}

	for {

		tail := atomicLoad(&queue.tail)
		next := atomicLoad(&tail.next)

		if tail == atomicLoad(&queue.tail) {
			if next != nil {
				atomicCompareAndSwap(&queue.tail, tail, next)
			} else {
				if atomicCompareAndSwap(&tail.next, next, newNode) {
					atomicCompareAndSwap(&queue.tail, tail, newNode)
					return
				}
			}
		}

	}

}

// Dequeue removes a Request from the queue
func (queue *LockFreeQueue) Dequeue() *Request {

	for {
		head := atomicLoad(&queue.head)
		tail := atomicLoad(&queue.tail)
		next := atomicLoad(&head.next)

		if head == tail {
			if next == nil {
				return nil
			}

			// Trying to advamce the tail
			atomicCompareAndSwap(&queue.tail, tail, next)
		} else {
			task := next.task
			if atomicCompareAndSwap(&queue.head, head, next) {
				return task
			}

		}

	}
}

func atomicLoad(pointer *unsafe.Pointer) *Node {
	return (*Node)(atomic.LoadPointer(pointer))
}

func atomicCompareAndSwap(pointer *unsafe.Pointer, old *Node, new *Node) bool {
	return atomic.CompareAndSwapPointer(pointer, unsafe.Pointer(old), unsafe.Pointer(new))
}
