package eventbus_test

import (
	eb "github.com/dtomasi/go-event-bus/v3"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
)

//nolint:gochecknoinits
func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestNewSafeCounter(t *testing.T) {
	counter := eb.NewSafeCounter()
	assert.Equal(t, 0, counter.Value())

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			for c := 0; c < 1000; c++ {
				// Put some noise into inc
				n := rand.Intn(10) //nolint:gosec
				time.Sleep(time.Duration(n) * time.Nanosecond)
				counter.Inc()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, 50000, counter.Value())

	counter.DecBy(10000)

	assert.Equal(t, 40000, counter.Value())

	counter.Dec()

	assert.Equal(t, 39999, counter.Value())
}
