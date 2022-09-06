package ds_test

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/sebnyberg/exp/ds"
	"github.com/sebnyberg/exp/ds/lfstack"
	"github.com/sebnyberg/exp/ds/lfstack_skeeto"
	"github.com/sebnyberg/exp/ds/lfstackbuff"
)

var v int
var ok bool

type Stacker[T any] interface {
	Push(item T)
	Pop() (T, bool)
}
type stackDef[T any] struct {
	name   string
	create func() Stacker[T]
}

func getStackDefs[T any]() []stackDef[T] {
	defs := []stackDef[T]{
		{"LFStack", func() Stacker[T] {
			return new(lfstack.LFStack[T])
		}},
		{"LFStackGC2", func() Stacker[T] {
			var stack lfstackbuff.LFStack[T]
			return &stack
		}},
		{"LFStack2", func() Stacker[T] {
			return new(lfstack_skeeto.LFStack[T])
		}},
		{"SyncStack", func() Stacker[T] {
			return new(ds.SyncStack[T])
		}},
	}
	// Randomize order (the teeniest Fisher Yates' ever)
	j := rand.Intn(2)
	defs[1], defs[j] = defs[j], defs[1]
	return defs
}

func BenchmarkStackSequential(b *testing.B) {
	for _, def := range getStackDefs[int]() {
		for _, sz := range []int{1000, 10000} {
			b.Run(fmt.Sprintf("%s_%d", def.name, sz), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					stack := def.create()
					for i := 0; i < sz; i++ {
						stack.Push(i)
					}
					for i := 0; i < sz; i++ {
						v, ok = stack.Pop()
					}
				}
			})
		}
	}
}

func TestStack(t *testing.T) {
	const nitems = 1e5
	const maxVal = 100

	for _, def := range getStackDefs[int]() {
		t.Run(def.name, func(t *testing.T) {
			t.Parallel()

			// Perform a bunch of random pushes/pops for the stack, with a slight bias
			// towards pushing.
			var pushed [maxVal + 1]uint64
			var wg sync.WaitGroup

			stack := def.create()

			wg.Add(runtime.NumCPU())
			for c := 0; c < runtime.NumCPU(); c++ {
				go func() {
					defer wg.Done()
					for i := 0; i < nitems; i++ {
						if rand.Intn(10) >= 4 {
							v := rand.Intn(maxVal + 1)
							atomic.AddUint64(&pushed[v], 1)
							stack.Push(v)
						} else {
							v, ok := stack.Pop()
							if ok {
								atomic.AddUint64(&pushed[v], ^uint64(0))
							}
						}
					}
				}()
			}
			wg.Wait()

			// empty the stack
			for {
				v, ok := stack.Pop()
				if !ok {
					break
				}
				atomic.AddUint64(&pushed[v], ^uint64(0))
			}

			for x, npush := range pushed {
				if npush != 0 {
					t.Fatalf("non-zero remaining items with value %v, was: %v",
						x, npush,
					)
				}
			}
		})
	}
}

func BenchmarkStackConcurrent(b *testing.B) {
	for _, sz := range []int{1000, 10000, 100000} {
		m := runtime.NumCPU() * 10

		// Pre-determine random actions that the goroutines will take
		actions := make([][]byte, m)
		for j := range actions {
			actions[j] = make([]byte, sz)
			n, err := rand.Read(actions[j])
			if n != sz || err != nil {
				b.Fail()
			}
		}

		for _, def := range getStackDefs[int]() {
			b.Run(fmt.Sprintf("%s_%d", def.name, sz), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					stack := def.create()
					var wg sync.WaitGroup
					wg.Add(m)
					for j := 0; j < m; j++ {
						go func(j int) {
							defer wg.Done()
							for k := 0; k < sz; k++ {
								if actions[j][k]&1 == 0 {
									stack.Push(k)
								} else {
									stack.Pop()
								}
							}
						}(j)
					}
					wg.Wait()
				}
			})
		}
	}
}
