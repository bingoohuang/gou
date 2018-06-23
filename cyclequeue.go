package go_utils

import (
	"sync"
)

// https://gist.github.com/moraes/2141121

func NewCycleQueue(capacity int) *CycleQueue {
	return &CycleQueue{
		Nodes:    make([]interface{}, capacity),
		Capacity: capacity,
	}
}

// CycleQueue is a basic FIFO queue based on a circular list.
type CycleQueue struct {
	Mutex    sync.Mutex
	Attach   interface{}
	Nodes    []interface{}
	Capacity int
	Counting int
}

// adds a node to the cycle queue.
func (q *CycleQueue) Add(n interface{}) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	q.Nodes[q.Counting%q.Capacity] = n
	q.Counting++
}

func (q *CycleQueue) FetchAll(index int) ([]interface{}, int) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if index < 0 {
		index = IntMax(q.Counting-1, 0)
	}

	nodes := make([]interface{}, 0)
	for index < q.Counting && len(nodes) < q.Capacity {
		nodes = append(nodes, q.Nodes[index%q.Capacity])
		index++
	}

	return nodes, index
}
