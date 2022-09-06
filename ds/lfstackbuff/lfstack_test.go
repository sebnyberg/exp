package lfstackbuff

import "testing"

func TestLFStackGC(t *testing.T) {
	var stack LFStack[int]
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
