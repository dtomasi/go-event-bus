package helpers

import "sync"

// A concurrency safe counter
type SafeCounter struct {
	mu sync.Mutex
	v  int
}

// Create a new counter
func NewSafeCounter() *SafeCounter {
	return &SafeCounter{
		v: 0,
	}
}

// Increment the counter
func (c *SafeCounter) Inc() {
	c.mu.Lock()
	c.v++
	c.mu.Unlock()
}

// Get the counter value
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v
}

