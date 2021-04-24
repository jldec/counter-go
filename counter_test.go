package counter_test

import (
	"testing"

	"github.com/jldec/counter-go"
)

func TestCounter(t *testing.T) {
	t.Run("Atomic", func(t *testing.T) {
		test1(new(counter.CounterAtomic), t)
	})
	t.Run("Mutex", func(t *testing.T) {
		test1(new(counter.CounterMutex), t)
	})
	t.Run("Channel", func(t *testing.T) {
		test1(counter.NewCounterChannel(), t)
	})
}

func TestCounterChannelPanicGet(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("CounterChannel failed to panic on uninitialized Get().")
		}
	}()
	cnt := new(counter.CounterChannel)
	cnt.Get()
}

func TestCounterChannelPanicInc(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("CounterChannel failed to panic on uninitialized Inc().")
		}
	}()
	cnt := new(counter.CounterChannel)
	cnt.Inc()
}

// basic counter test - not concurrent
func test1(cnt counter.Counter, t *testing.T) {
	val := cnt.Get()

	if val != 0 {
		t.Errorf("Inital value: %d\nexpected: 0", val)
	}

	cnt.Inc()
	cnt.Inc()
	cnt.Inc()

	val = cnt.Get()

	if val != 3 {
		t.Errorf("After 3x Inc() value: %d\nexpected: 3", val)
	}

	val = cnt.Get()

	if val != 3 {
		t.Errorf("After 3x Inc() + 1 Get() value: %d\nexpected: 3", val)
	}
}

// Simple: 1 op = 1 Inc() in same thread
func BenchmarkCounter_1(b *testing.B) {
	b.Run("Atomic", func(b *testing.B) {
		bench1(new(counter.CounterAtomic), b)
	})
	b.Run("Mutex", func(b *testing.B) {
		bench1(new(counter.CounterMutex), b)
	})
	b.Run("Channel", func(b *testing.B) {
		bench1(counter.NewCounterChannel(), b)
	})
}

func bench1(cnt counter.Counter, b *testing.B) {
	i := 0
	for i < b.N {
		cnt.Inc()
		i++
	}
	// b.Logf("%d iterations, counter = %d", i, cnt.Get())
}

// 1 op = 1 Inc() x 10 goroutines
// 1 op = [ 1 Inc() + 10 Get() ] x 10 goroutines
func BenchmarkCounter_2(b *testing.B) {
	b.Run("Atomic no reads", func(b *testing.B) {
		bench2(new(counter.CounterAtomic), 0, b)
	})
	b.Run("Mutex no reads", func(b *testing.B) {
		bench2(new(counter.CounterMutex), 0, b)
	})
	b.Run("Channel no reads", func(b *testing.B) {
		bench2(counter.NewCounterChannel(), 0, b)
	})
	b.Run("Atomic 10 reads", func(b *testing.B) {
		bench2(new(counter.CounterAtomic), 10, b)
	})
	b.Run("Mutex 10 reads", func(b *testing.B) {
		bench2(new(counter.CounterMutex), 10, b)
	})
	b.Run("Channel 10 reads", func(b *testing.B) {
		bench2(counter.NewCounterChannel(), 10, b)
	})
}

// Concurrent counter benchmark
func bench2(cnt counter.Counter, readRatio int, b *testing.B) {
	done := make(chan int)

	const GONUM = 10 // number of goroutines

	for n := 0; n < GONUM; n++ {
		go func(n int) {
			var sumReads uint64 = 0
			for i := 0; i < b.N; i++ {
				cnt.Inc()
				if readRatio > 0 {
					for j := 0; j < readRatio; j++ {
						sumReads += cnt.Get()
					}
				}
			}
			done <- n
		}(n)
	}

	// Collect goroutine completion order and counter values.
	type result struct {
		n     int
		count uint64
	}
	var order [GONUM]result
	for n := 0; n < GONUM; n++ {
		order[n] = result{<-done, cnt.Get()}
	}

	val := cnt.Get()
	if uint64(b.N*GONUM) != val {
		b.Errorf("Concurrent counter expected %d\ngot %d", b.N*GONUM, val)
	}
	// b.Logf("%d iterations, order %v", b.N, order)
}
