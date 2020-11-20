package eventbus

import (
	"sync"
	"testing"
)

func TestNewEventBus(t *testing.T) {
	eb := NewEventBus()
	if eb == nil {
		t.Fail()
	}

	seb := GetBus()
	if seb == nil {
		t.Fail()
	}
}

func TestEventBus_Subscribe(t *testing.T) {
	eb := NewEventBus()
	ch := NewEventChannel()
	eb.Subscribe("foo", ch)

	sbs, ok := eb.subscribers["foo"]
	if !ok {
		t.Error("subscriber topic was not registered")
	}

	if len(sbs) != 1 {
		t.Error("subscriber was registered correctly")
	}
}

func TestEventBus_PublishAsync(t *testing.T) {
	eb := NewEventBus()
	ch1 := NewEventChannel()
	ch2 := NewEventChannel()
	eb.Subscribe("foo:baz", ch1)
	eb.Subscribe("foo:*", ch2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		evt := <-ch1
		if evt.Topic != "foo:baz" {
			t.Fail()
		}

		if evt.Data != "bar" {
			t.Fail()
		}
		wg.Done()
	}()

	go func() {
		evt := <-ch2
		if evt.Topic != "foo:baz" {
			t.Fail()
		}

		if evt.Data != "bar" {
			t.Fail()
		}
		wg.Done()
	}()

	eb.PublishAsync("foo:baz", "bar")

	wg.Wait()
}

func TestEventBus_Publish(t *testing.T) {
	eb := NewEventBus()
	ch1 := NewEventChannel()
	ch2 := NewEventChannel()
	eb.Subscribe("foo:baz", ch1)
	eb.Subscribe("foo:*", ch2)

	var callCount int

	go func() {
		evt := <-ch1
		if evt.Topic != "foo:baz" {
			t.Fail()
		}

		if evt.Data != "bar" {
			t.Fail()
		}
		callCount++
		evt.Done()
	}()

	go func() {
		evt := <-ch2
		if evt.Topic != "foo:baz" {
			t.Fail()
		}

		if evt.Data != "bar" {
			t.Fail()
		}
		callCount++
		evt.Done()
	}()

	eb.Publish("foo:baz", "bar")

	if callCount != 2 {
		t.Fail()
	}
}

func TestEventBus_SubscribeCallback(t *testing.T) {
	eb := NewEventBus()

	var callCount int

	eb.SubscribeCallback("foo:baz", func(topic string, data interface{}) {
		if topic != "foo:baz" {
			t.Fail()
		}

		if data != "bar" {
			t.Fail()
		}
		callCount++
	})

	eb.SubscribeCallback("foo:*", func(topic string, data interface{}) {
		if topic != "foo:baz" {
			t.Fail()
		}

		if data != "bar" {
			t.Fail()
		}
		callCount++
	})

	eb.Publish("foo:baz", "bar")

	if callCount != 2 {
		t.Fail()
	}
}
