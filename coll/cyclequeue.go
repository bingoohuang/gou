package coll

import (
	"sync"

	"github.com/bingoohuang/gou/mat"
)

// https://gist.github.com/moraes/2141121

// NewCycleQueue ...
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

// Add adds a node to the cycle queue.
func (q *CycleQueue) Add(n interface{}) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	q.Nodes[q.Counting%q.Capacity] = n
	q.Counting++
}

// FetchAll ...
func (q *CycleQueue) FetchAll(index int) ([]interface{}, int) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if index < 0 {
		index = mat.MaxInt(q.Counting-1, 0) // nolint gomnd
	}

	nodes := make([]interface{}, 0)
	for index < q.Counting && len(nodes) < q.Capacity {
		nodes = append(nodes, q.Nodes[index%q.Capacity])
		index++
	}

	return nodes, index
}
