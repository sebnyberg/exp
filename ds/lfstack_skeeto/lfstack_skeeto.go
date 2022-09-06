package lfstack_skeeto

import (
	"sync"
	"sync/atomic"
)

// LFStack implements a concurrency-safe, lock-free LIFO stack similar to
// skeeto's C11 version.
type LFStack[T any] struct {
	initOnce   sync.Once
	free, head atomic.Pointer[lfhead[T]]
	size       atomic.Uint64
}

type lfhead[T any] struct {
	node *lfnode[T]
	aba  int
}

type lfnode[T any] struct {
	value T
	next  *lfnode[T]
}

func pop[T any](head *atomic.Pointer[lfhead[T]]) *lfnode[T] {
	var next lfhead[T]
	for orig := head.Load(); orig.node != nil; orig = head.Load() {
		next.aba = orig.aba + 1
		next.node = orig.node.next
		if head.CompareAndSwap(orig, &next) {
			return orig.node
		}
	}
	return nil
}

func push[T any](head *atomic.Pointer[lfhead[T]], node *lfnode[T]) {
	var nextHead lfhead[T]
	nextHead.node = node
	for orig := head.Load(); ; orig = head.Load() {
		nextHead.aba = orig.aba + 1
		node.next = orig.node
		if head.CompareAndSwap(orig, &nextHead) {
			break
		}
	}
}

func (s *LFStack[T]) newnode() *lfnode[T] {
	if node := pop(&s.free); node != nil {
		return node
	}
	return new(lfnode[T])
}

func (s *LFStack[T]) init() {
	s.initOnce.Do(func() {
		s.head.Store(&lfhead[T]{
			node: nil,
			aba:  0,
		})
		s.free.Store(&lfhead[T]{
			node: nil,
			aba:  0,
		})
	})
}

// Push pushes an item onto the stack.
func (s *LFStack[T]) Push(x T) {
	s.init()
	node := s.newnode()
	node.value = x
	push(&s.head, node)
	s.size.Add(1)
}

// Pop pops an item from the stack. If the stack is empty, the second return
// value is false.
func (s *LFStack[T]) Pop() (item T, ok bool) {
	s.init()
	node := pop(&s.head)
	if node == nil {
		return item, false
	}
	// Put item into freelist
	item = node.value
	push(&s.free, node)
	s.size.Add(^uint64(0))
	return item, true
}

// Peek removes the top-most item from the stack. If the stack is empty, the
// second return value is false.
func (s *LFStack[T]) Peek() (item T, ok bool) {
	s.init()
	old := s.head.Load()
	if old.node != nil {
		return old.node.value, true
	}
	var empty T
	return empty, false
}

// Len returns the length of the stack.
func (s *LFStack[T]) Len() int {
	return int(s.size.Load())
}
