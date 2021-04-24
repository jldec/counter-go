# counter-go

[![CI](https://github.com/jldec/counter-go/workflows/CI/badge.svg)](https://github.com/jldec/counter-go/actions)  
Go.dev [github.com/jldec/counter-go](https://pkg.go.dev/github.com/jldec/counter-go)

This package demonstrates 3 different implementations of a threadsafe global counter.

1. **CounterAtomic** uses `atomic.AddUint64` and `atomic.LoadUint64`.
2. **CounterMutex** uses `sync.RWMutex`.
3. **CounterChannel** serializes all reads and writes inside 1 goroutine with 2 channels.

This was written to accompany the blog post [Getting started with Goroutines and channels](https://jldec.me/getting-started-with-go-part-3-goroutines-and-channels).

All 3 types implement the Counter interface:

```go
type Counter interface {
    Get() uint32 // get current counter value
    Inc()        // increment by 1
}
```

### Initialization

Call `counter.NewCounterChannel()` to use CounterChannel.  
The other 2 do not require any initialization.


### Benchmarks

```
$ go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/jldec/counter-go
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
```

#### Simple: 1 op = 1 Inc() in same thread
```
BenchmarkCounter_1/Atomic-12                 195965660          6 ns/op
BenchmarkCounter_1/Mutex-12                   54177086         22 ns/op
BenchmarkCounter_1/Channel-12                  4499144        286 ns/op
```

#### Concurrent: 1 op = 1 Inc() across each of 10 goroutines
```
BenchmarkCounter_2/Atomic_no_reads-12          7298484        191 ns/op
BenchmarkCounter_2/Mutex_no_reads-12           1966656        621 ns/op
BenchmarkCounter_2/Channel_no_reads-12          256842       4771 ns/op
```

#### Concurrent: 1 op = [ 1 Inc() + 10 Get() ] across each of 10 goroutines
```
BenchmarkCounter_2/Atomic_10_reads-12          3922029        286 ns/op
BenchmarkCounter_2/Mutex_10_reads-12            416354       2844 ns/op
BenchmarkCounter_2/Channel_10_reads-12           21506      55733 ns/op
```

#### Constrained to single thread
```
$ GOMAXPROCS=1 go test -bench .

BenchmarkCounter_1/Atomic                    197135869          6 ns/op
BenchmarkCounter_1/Mutex                      55698454         22 ns/op
BenchmarkCounter_1/Channel                     5689788        214 ns/op

BenchmarkCounter_2/Atomic_no_reads            19519166         60 ns/op
BenchmarkCounter_2/Mutex_no_reads              4702759        254 ns/op
BenchmarkCounter_2/Channel_no_reads             530554       2197 ns/op

BenchmarkCounter_2/Atomic_10_reads             6269979        189 ns/op
BenchmarkCounter_2/Mutex_10_reads               927439       1354 ns/op
BenchmarkCounter_2/Channel_10_reads              47889      25054 ns/op
```