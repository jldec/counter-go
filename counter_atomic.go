package counter

import "sync/atomic"

// Thread-safe counter
// Uses sync.atomic to coordinate reads and writes.
// No initialization required.
type CounterAtomic struct {
	count uint64
}

// Increment counter
func (c *CounterAtomic) Inc() {
	atomic.AddUint64(&c.count, 1)
}

// Get current counter value
func (c *CounterAtomic) Get() uint64 {
	return atomic.LoadUint64(&c.count)
}
