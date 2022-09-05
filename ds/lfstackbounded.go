package ds

// import (
// 	"sync"
// 	"sync/atomic"
// )

// // LFStackBounded implements a lock-free, fixed-size LIFO stack of items. Its
// // zero-value can be used as-is, however the size should be set with SetSize(),
// // or it will default to 1000 items.
// type LFStackBounded[T any] struct {
// 	buf        *lfBoundedNode[T]
// 	head, free atomic.Pointer[lfBoundedHead[T]]
// 	size       atomic.Uint64
// 	initOnce   sync.Once
// 	maxSize    uint64
// }

// // lfBoundedNode wraps items and provides a pointer to the next node.
// type lfBoundedNode[T any] struct {
// 	item T
// 	next *lfnode[T]
// }

// type lfBoundedHead[T any] struct {
// 	aba  uint64
// 	node *lfBoundedNode[T]
// }

// // // Init initializes the stack with some items. Prior contents (if any) will be
// // // left for the GC to clean up.
// // func (s *LFStackBounded[T]) Init(items ...T) {
// // 	var head *lfnode[T]
// // 	var size uint64
// // 	for _, x := range items {
// // 		size++
// // 		newHead := new(lfnode[T])
// // 		newHead.next = head
// // 		newHead.item = x
// // 		head = newHead
// // 	}
// // 	s.head.Store(head)
// // 	s.size.Store(size)
// // }

// func (s *LFStackBounded[T]) init() {
// 	s.initOnce.Do(func() {
// 		items := make([]lfBoundedNode[T], s.maxSize)
// 	})
// }

// // Push pushes an item onto the stack.
// //
// // Note that this is aba-safe because new lfnodes are allocated on each Push.
// // In other words, it's not possible for the caller to re-insert the same node
// // twice.
// func (s *LFStackBounded[T]) Push(x T) {
// 	var next lfnode[T]
// 	next.item = x
// 	for old := s.head.Load(); ; old = s.head.Load() {
// 		next.next = old
// 		if s.head.CompareAndSwap(old, &next) {
// 			s.size.Add(1)
// 			return
// 		}
// 	}
// }

// // Pop pops an item from the stack. If the stack is empty, the second return
// // value is false.
// func (s *LFStackBounded[T]) Pop() (item T, ok bool) {
// 	// At this point, the pool must've been loaded, or else there would be no
// 	// items to return.
// 	for old := s.head.Load(); old != nil; old = s.head.Load() {
// 		if s.head.CompareAndSwap(old, old.next) {
// 			s.size.Add(^uint64(0))
// 			return old.item, true
// 		}
// 	}
// 	return item, false
// }

// // Peek removes the top-most item from the stack. If the stack is empty, the
// // second return value is false.
// func (s *LFStackBounded[T]) Peek() (item T, ok bool) {
// 	old := s.head.Load()
// 	if old != nil {
// 		return old.item, true
// 	}
// 	var empty T
// 	return empty, false
// }

// // Len returns the length of the stack.
// func (s *LFStackBounded[T]) Len() int {
// 	return int(s.size.Load())
// }
