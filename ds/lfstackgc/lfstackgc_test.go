package lfstackgc

import "testing"

var global any

// func TestLFStackGC(t *testing.T) {
// 	stack := new(uint64)
// 	global = stack

// 	var nodes []*LFNode[int]

// 	if LFStackPop(stack) != nil {
// 		t.Fatalf("stack is not empty")
// 	}

// 	// Push one element
// 	node := &LFNode[int]{Value: 42}
// 	nodes = append(nodes, node)
// 	LFStackPush(stack, AsHdr(node))

// 	// Push another
// 	node = &LFNode[int]{Value: 52}
// 	nodes = append(nodes, node)
// 	LFStackPush(stack, AsHdr(node))

// 	// Pop one element
// 	node = AsNode[LFNode[int]](LFStackPop(stack))
// 	if node == nil {
// 		t.Fatalf("empty stack")
// 	}
// 	if node.Value != 52 {
// 		t.Fatalf("no lifo")
// 	}

// 	// Pop another
// 	node = AsNode[LFNode[int]](LFStackPop(stack))
// 	if node == nil {
// 		t.Fatalf("empty stack")
// 	}
// 	if node.Value != 42 {
// 		t.Fatalf("no lifo")
// 	}
// }

func TestLFStackGC(t *testing.T) {
	var stack LFStack[int]
	stack.head = new(uint64)
	global = stack

	_, ok := stack.Pop()
	if ok {
		t.Fatalf("stack is not empty")
	}

	// Push one element
	stack.Push(42)

	// Push another
	stack.Push(52)

	// Pop one element
	v, ok := stack.Pop()
	if !ok || v != 52 {
		t.Fatalf("empty stack")
	}

	// Pop another
	v, ok = stack.Pop()
	if !ok || v != 42 {
		t.Fatalf("empty stack")
	}
}
