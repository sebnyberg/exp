package ds_test

import "testing"

var llroot *Node[int]

func BenchmarkListFromSlice(b *testing.B) {
	const nitems = 1e5
	getItems := func() []int {
		items := make([]int, nitems)
		for i := range items {
			items[i] = i
		}
		return items
	}
	b.Run("naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			items := getItems()
			b.StartTimer()
			llroot = NewListFromSliceNaive(items)
		}
	})
	b.Run("maybefaster", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			items := getItems()
			b.StartTimer()
			llroot = NewListFromSliceMaybeFaster(items)
		}
	})
}

type Node[T any] struct {
	Val  T
	Next *Node[T]
}

func NewListFromSliceNaive[T any](values []T) *Node[T] {
	dummy := new(Node[T]) // sentinel node to reduce if/else statements
	prev := dummy
	for _, v := range values {
		curr := new(Node[T])
		curr.Val = v
		prev.Next, prev = curr, curr
	}
	return dummy.Next
}

func NewListFromSliceMaybeFaster[T any](values []T) *Node[T] {
	dummy := new(Node[T])
	prev := dummy
	preAlloced := make([]Node[T], len(values))
	for i, v := range values {
		preAlloced[i].Val = v
		prev.Next, prev = &preAlloced[i], &preAlloced[i]
	}
	return dummy.Next
}
