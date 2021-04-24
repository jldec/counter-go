package counter

import (
	"sync"
)

// Thread-safe counter
// Uses sync.RWMutex to coordinate reads and writes.
// No initialization required.
type CounterMutex struct {
	count uint64
	rw    sync.RWMutex
}

// Increment counter by pushing an arbitrary int to the write channel.
func (c *CounterMutex) Inc() {
	c.rw.Lock()
	c.count++
	c.rw.Unlock()
}

// Get current counter value from the read channel.
func (c *CounterMutex) Get() (ret uint64) {
	c.rw.RLock()
	ret = c.count
	c.rw.RUnlock()
	return
}
