package lfstackbuff

import (
	"sync"
	"sync/atomic"
)

type LFStack[T any] struct {
	head atomic.Pointer[LFNode[T]]

	// buf generates sets of LFNode[T] in batches to reduce individual
	// allocations.
	buf   []LFNode[T]
	bufmu sync.Mutex
	capsz uint64
}

// acquirenode has an off-by-one bug where one item is lost per allocation.
func (s *LFStack[T]) acquirenode() *LFNode[T] {
	for {
		x := atomic.LoadUint64(&s.capsz)
		next := x + 1
		if atomic.CompareAndSwapUint64(&s.capsz, x, next) {
			// we're guaranteed to 'own' the index sz-1
			// However, the capacity may not have enough elements to fetch that item.
			// So if sz-1 >= capacity, we lock bufmu, check whether the current
			// capacity has been increased or not, then re-allocate buf
			sz, currCap := x&0xffffffff, x>>32
			if sz < currCap {
				return &s.buf[sz]
			}
			s.bufmu.Lock()
			x := atomic.LoadUint64(&s.capsz)
			newCap := x >> 32
			if newCap == currCap {
				// Growslice
				s.buf = append(s.buf, LFNode[T]{})
				s.buf = s.buf[:cap(s.buf)]
				// Update capacity
				for {
					orig := atomic.LoadUint64(&s.capsz)
					new := x&0xffffffff | (uint64(cap(s.buf)) << 32)
					if atomic.CompareAndSwapUint64(&s.capsz, orig, new) {
						break
					}
				}
			}
			s.bufmu.Unlock()
			return s.acquirenode() // try again
		}
	}
}

func (s *LFStack[T]) Push(v T) {
	new := s.acquirenode()
	new.Value = v
	new.Pushcnt++
	for {
		old := s.head.Load()
		new.Next = old
		if s.head.CompareAndSwap(old, new) {
			return
		}
	}
}

func (s *LFStack[T]) Pop() (v T, ok bool) {
	for {
		old := s.head.Load()
		if old == nil {
			return v, false
		}
		next := old.Next
		if s.head.CompareAndSwap(old, next) {
			return old.Value, true
		}
	}
}

type LFStackHead[T any] atomic.Pointer[T]

type LFNode[T any] struct {
	Next    *LFNode[T]
	Pushcnt uintptr
	Value   T
}
