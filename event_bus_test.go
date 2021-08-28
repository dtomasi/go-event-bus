package eventbus_test

import (
	eb "github.com/dtomasi/go-event-bus/v2"
	"github.com/dtomasi/helpers"
	"sync"
	"testing"
)

func TestNewEventBus(t *testing.T) {
	if ebi := eb.NewEventBus(); ebi == nil {
		t.Fail()
	}

	if seb := eb.DefaultBus(); seb == nil {
		t.Fail()
	}
}

func TestEventBus_Subscribe(t *testing.T) {
	inst := eb.NewEventBus()
	_ = inst.Subscribe("foo")

	if !inst.HasSubscribers("foo") {
		t.Error("subscriber topic was not registered")
	}
}

func TestEventBus_SubscribeChannel(t *testing.T) {
	inst := eb.NewEventBus()
	ch := eb.NewEventChannel()
	inst.SubscribeChannel("foo", ch)

	if !inst.HasSubscribers("foo") {
		t.Error("subscriber topic was not registered")
	}
}

func TestEventBus_PublishAsync(t *testing.T) {
	inst := eb.NewEventBus()
	ch1 := eb.NewEventChannel()
	inst.SubscribeChannel("foo:baz", ch1)
	ch2 := inst.Subscribe("foo:*")

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		evt := <-ch1
		if evt.Topic != "foo:baz" { // nolint:goconst
			t.Fail()
		}

		if evt.Data != "bar" { // nolint:goconst
			t.Fail()
		}

		wg.Done()
	}() //nolint:wsl,nolintlint

	go func() {
		evt := <-ch2
		if evt.Topic != "foo:baz" {
			t.Fail()
		}

		if evt.Data != "bar" {
			t.Fail()
		}

		wg.Done()
	}() //nolint:wsl,nolintlint

	inst.PublishAsync("foo:baz", "bar")

	wg.Wait()
}

func TestEventBus_Publish(t *testing.T) {
	inst := eb.NewEventBus()
	ch1 := eb.NewEventChannel()
	inst.SubscribeChannel("foo:baz", ch1)
	ch2 := inst.Subscribe("foo:*")

	callCounter := helpers.NewSafeCounter()

	go func() {
		evt := <-ch1
		if evt.Topic != "foo:baz" {
			t.Fail()
		}

		if evt.Data != "bar" {
			t.Fail()
		}

		callCounter.Inc()

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

		callCounter.Inc()

		evt.Done()
	}()

	inst.Publish("foo:baz", "bar")

	if callCounter.Value() != 2 {
		t.Fail()
	}
}

func TestEventBus_SubscribeCallback(t *testing.T) {
	inst := eb.NewEventBus()

	callCounter := helpers.NewSafeCounter()

	inst.SubscribeCallback("foo:baz", func(topic string, data interface{}) {
		if topic != "foo:baz" {
			t.Fail()
		}

		if data != "bar" {
			t.Fail()
		}
		callCounter.Inc()
	})

	inst.SubscribeCallback("foo:*", func(topic string, data interface{}) {
		if topic != "foo:baz" {
			t.Fail()
		}

		if data != "bar" {
			t.Fail()
		}
		callCounter.Inc()
	})

	inst.Publish("foo:baz", "bar")

	if callCounter.Value() != 2 {
		t.Fail()
	}
}
