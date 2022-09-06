package lfstackgc

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type LFStack[T any] struct {
	head   *uint64
	pinned []*LFNode[T]
}

func (s *LFStack[T]) Push(v T) {
	x := &LFNode[T]{Value: v}
	s.pinned = append(s.pinned, x)
	(*LFStackHead)(s.head).push((*LFNodeHdr)(unsafe.Pointer(x)))
}

func (s *LFStack[T]) Pop() (v T, ok bool) {
	ptr := (*LFNodeHdr)(unsafe.Pointer((*LFStackHead)(s.head).pop()))
	if ptr == nil {
		return v, false
	}
	res := (*LFNode[T])(unsafe.Pointer(ptr))
	return res.Value, true
}

func (s *LFStack[T]) Init() {
	s.head = new(uint64)
}

type LFStackHead uint64

type LFNodeHdr struct {
	Next    uint64
	Pushcnt uintptr
}

type LFNode[T any] struct {
	LFNodeHdr
	Value T
}

func AsHdr[T any](val *T) *LFNodeHdr {
	return (*LFNodeHdr)(unsafe.Pointer(val))
}

func AsNode[T any](hdr *LFNodeHdr) *T {
	return (*T)(unsafe.Pointer(hdr))
}

func (head *LFStackHead) push(node *LFNodeHdr) {
	node.Pushcnt++
	new := lfstackPack(node, node.Pushcnt)
	if node1 := lfstackUnpack(new); node1 != node {
		panic("push: invalid node packing")
	}
	for {
		old := atomic.LoadUint64((*uint64)(head))
		node.Next = old
		if atomic.CompareAndSwapUint64((*uint64)(head), old, new) {
			break
		}
	}
}

func (head *LFStackHead) pop() unsafe.Pointer {
	for {
		old := atomic.LoadUint64((*uint64)(head))
		if old == 0 {
			return nil
		}
		node := lfstackUnpack(old)
		next := atomic.LoadUint64(&node.Next)
		if atomic.CompareAndSwapUint64((*uint64)(head), old, next) {
			return unsafe.Pointer(node)
		}
	}
}

func (head *LFStackHead) empty() bool {
	return atomic.LoadUint64((*uint64)(head)) == 0
}

// lfnodeValidate panics if node is not a valid address for use with
// lfstack.push. This only needs to be called when node is allocated.
func lfnodeValidate(node *LFNodeHdr) {
	if lfstackUnpack(lfstackPack(node, ^uintptr(0))) != node {
		fmt.Println("bad lfnode address")
	}
}
