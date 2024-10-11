package counter

// Thread-safe counter
// Uses 2 Channels to coordinate reads and writes.
// Must be initialized with NewCounterChannel().
type CounterChannel struct {
	readCh  chan uint64
	writeCh chan int
}

// NewCounterChannel() is required to initialize a Counter.
func NewCounterChannel() *CounterChannel {
	c := &CounterChannel{
		readCh:  make(chan uint64),
		writeCh: make(chan int),
	}

	// The actual counter value lives inside this goroutine.
	// It can only be accessed for R/W via one of the channels.
	go func() {
		var count uint64 = 0
		for {
			select {
			// Reading from readCh is equivalent to reading count.
			case c.readCh <- count:
			case delta := <-c.writeCh:
				if delta < 0 && count > 0 {
					count--
				} else if delta > 0 {
					count++
				}
			}
		}
	}()

	return c
}

// Increment counter by pushing 1 to the write channel.
func (c *CounterChannel) Inc() {
	c.check()
	c.writeCh <- 1
}

// Decrement counter by pushing -1 to the write channel.
func (c *CounterChannel) Dec() {
	c.check()
	c.writeCh <- -1
}

// Get current counter value from the read channel.
func (c *CounterChannel) Get() uint64 {
	c.check()
	return <-c.readCh
}

func (c *CounterChannel) check() {
	if c.readCh == nil {
		panic("Uninitialized Counter, requires NewCounterChannel()")
	}
}
