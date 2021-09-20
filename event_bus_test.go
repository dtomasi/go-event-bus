package eventbus_test

import (
	eb "github.com/dtomasi/go-event-bus/v3"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewEventBus(t *testing.T) {
	if ebi := eb.NewEventBus(); ebi == nil {
		t.Fail()
	}
}

func TestEventBus_Subscribe(t *testing.T) {
	ebi := eb.NewEventBus()
	_ = ebi.Subscribe("foo")

	assert.True(t, ebi.HasSubscribers("foo"))
}

func TestEventBus_SubscribeChannel(t *testing.T) {
	ebi := eb.NewEventBus()
	ch := eb.NewEventChannel()
	ebi.SubscribeChannel("foo", ch)

	assert.True(t, ebi.HasSubscribers("foo"))
}

func TestEventBus_PublishAsync(t *testing.T) {
	const testTopicName = "foo:bar"

	ebi := eb.NewEventBus()
	ch1 := eb.NewEventChannel()
	ebi.SubscribeChannel(testTopicName, ch1)
	ch2 := ebi.Subscribe("foo:*")

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		evt := <-ch1
		if evt.Topic != testTopicName {
			t.Fail()
		}

		if evt.Data != "bar" { // nolint:goconst
			t.Fail()
		}

		wg.Done()
	}() //nolint:wsl,nolintlint

	go func() {
		evt := <-ch2
		if evt.Topic != testTopicName {
			t.Fail()
		}

		if evt.Data != "bar" {
			t.Fail()
		}

		wg.Done()
	}() //nolint:wsl,nolintlint

	ebi.PublishAsync(testTopicName, "bar")

	wg.Wait()

	assert.Equal(t, 1, ebi.Stats().GetPublishedCountByTopic(testTopicName))
}

func TestEventBus_Publish(t *testing.T) {
	const testTopicName = "foo:bar"

	ebi := eb.NewEventBus()
	ch1 := eb.NewEventChannel()
	ebi.SubscribeChannel(testTopicName, ch1)
	ch2 := ebi.Subscribe("foo:*")

	callCounter := eb.NewSafeCounter()

	go func() {
		evt := <-ch1
		if evt.Topic != testTopicName {
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
		if evt.Topic != testTopicName {
			t.Fail()
		}

		if evt.Data != "bar" {
			t.Fail()
		}

		callCounter.Inc()

		evt.Done()
	}()

	ebi.Publish(testTopicName, "bar")

	if callCounter.Value() != 2 {
		t.Fail()
	}

	assert.Equal(t, 1, ebi.Stats().GetPublishedCountByTopic(testTopicName))

	// Try to republish with publish once
	ebi.PublishOnce(testTopicName, "bar")

	// Count should be still 1
	assert.Equal(t, 1, ebi.Stats().GetPublishedCountByTopic(testTopicName))
}

func TestEventBus_SubscribeCallback(t *testing.T) {
	const testTopicName = "foo:bar"

	ebi := eb.NewEventBus()

	callCounter := eb.NewSafeCounter()

	ebi.SubscribeCallback(testTopicName, func(topic string, data interface{}) {
		if topic != testTopicName {
			t.Fail()
		}

		if data != "bar" {
			t.Fail()
		}
		callCounter.Inc()
	})

	ebi.SubscribeCallback("foo:*", func(topic string, data interface{}) {
		if topic != testTopicName {
			t.Fail()
		}

		if data != "bar" {
			t.Fail()
		}
		callCounter.Inc()
	})

	ebi.Publish(testTopicName, "bar")

	if callCounter.Value() != 2 {
		t.Fail()
	}

	// Try to republish with publish once
	ebi.PublishOnce(testTopicName, "bar")

	// Count should be still 1
	assert.Equal(t, 1, ebi.Stats().GetPublishedCountByTopic(testTopicName))
}
