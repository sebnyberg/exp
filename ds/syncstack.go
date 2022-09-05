package ds

import "sync"

// SyncStack implements a concurrency-safe, LIFO stack of items. Unlike lfstack,
// it uses a mutex to serialize access to its contents.
type SyncStack[T any] struct {
	items []T
	mu    sync.RWMutex
}

// Init initializes the stack with some items.
func (s *SyncStack[T]) Init(items ...T) {
	s.items = append(s.items, items...)
}

// Push pushes an item onto the stack.
func (s *SyncStack[T]) Push(x T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, x)
}

// Pop pops an item from the stack. If the stack is empty, the second return
// value is false.
func (s *SyncStack[T]) Pop() (item T, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := len(s.items)
	if n == 0 {
		return item, false
	}
	item = s.items[n-1]
	s.items = s.items[:n-1]
	return item, true
}

// Peek removes the top-most item from the stack. If the stack is empty, the
// second return value is false.
func (s *SyncStack[T]) Peek() (item T, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := len(s.items)
	if n == 0 {
		return item, false
	}
	return s.items[n-1], true
}

// Len returns the length of the stack.
func (s *SyncStack[T]) Len() int {
	return len(s.items)
}
