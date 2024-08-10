package lib

import (
	"fmt"
	"sync"
)

type Value struct {
	v   int
	mux *sync.Mutex
}

func NewValue(v int) *Value {
	var mux sync.Mutex
	return &Value{v: v, mux: &mux}
}

func Multi(ch chan<- string) {
	var wg sync.WaitGroup
	for n := range 10 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			ch <- fmt.Sprintf("Hello %d", n)
		}(n)
	}
	wg.Wait()
	defer close(ch)
}

func (v *Value) Increment() {
	v.mux.Lock()
	v.v++
	v.mux.Unlock()
}
