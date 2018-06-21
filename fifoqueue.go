package go_utils

import (
	"container/list"
)

// FifoQuque is a head-tail linked list data structure implementation.
// It is based on a doubly linked list container, so that every
// operations time complexity is O(1).
//
// every operations over an instiated Deque are synchronized and
// safe for concurrent usage.
type FifoQueue struct {
	container *list.List
	capacity  int
}

// NewCappedDeque creates a Deque with the specified capacity limit.
func NewFifoQueue(capacity int) *FifoQueue {
	return &FifoQueue{
		container: list.New(),
		capacity:  capacity,
	}
}

// Append inserts element at the back of the Deque in a O(1) time complexity,
// returning true if successful or false if the deque is at capacity.
func (s *FifoQueue) Append(item interface{}) bool {
	if s.container.Len() >= s.capacity {
		s.Shift()
	}

	s.container.PushBack(item)
	return true
}

// Pop removes the last element of the deque in a O(1) time complexity
func (s *FifoQueue) Pop() interface{} {
	var item interface{} = nil
	var lastContainerItem *list.Element = nil

	lastContainerItem = s.container.Back()
	if lastContainerItem != nil {
		item = s.container.Remove(lastContainerItem)
	}

	return item
}

// Shift removes the first element of the deque in a O(1) time complexity
func (s *FifoQueue) Shift() interface{} {
	var item interface{} = nil
	var firstContainerItem *list.Element = nil

	firstContainerItem = s.container.Front()
	if firstContainerItem != nil {
		item = s.container.Remove(firstContainerItem)
	}

	return item
}

// First returns the first value stored in the deque in a O(1) time complexity
func (s *FifoQueue) First() interface{} {
	item := s.container.Front()
	if item != nil {
		return item.Value
	} else {
		return nil
	}
}

// Last returns the last value stored in the deque in a O(1) time complexity
func (s *FifoQueue) Last() interface{} {
	item := s.container.Back()
	if item != nil {
		return item.Value
	} else {
		return nil
	}
}

// Size returns the actual deque size
func (s *FifoQueue) Capacity() int {
	return s.capacity
}

func (s *FifoQueue) Size() int {
	return s.container.Len()
}

// Empty checks if the deque is empty
func (s *FifoQueue) Empty() bool {
	return s.container.Len() == 0
}

// Full checks if the deque is full
func (s *FifoQueue) Full() bool {
	return s.container.Len() >= s.capacity
}
