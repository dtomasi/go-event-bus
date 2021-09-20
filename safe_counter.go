package eventbus

import (
	"sync/atomic"
)

// SafeCounter is a concurrency safe counter.
type SafeCounter struct {
	v *uint64
}

// NewSafeCounter creates a new counter.
func NewSafeCounter() *SafeCounter {
	return &SafeCounter{
		v: new(uint64),
	}
}

// Value returns the current value.
func (c *SafeCounter) Value() int {
	return int(atomic.LoadUint64(c.v))
}

// IncBy increments the counter by given delta.
func (c *SafeCounter) IncBy(add uint) {
	atomic.AddUint64(c.v, uint64(add))
}

// Inc increments the counter by 1.
func (c *SafeCounter) Inc() {
	c.IncBy(1)
}

// DecBy decrements the counter by given delta.
func (c *SafeCounter) DecBy(dec uint) {
	atomic.AddUint64(c.v, ^uint64(dec-1))
}

// Dec decrements the counter by 1.
func (c *SafeCounter) Dec() {
	c.DecBy(1)
}
