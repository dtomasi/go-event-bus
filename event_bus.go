package eventbus

import (
	"sync"
)

// Event holds topic name and data.
type Event struct {
	Data  interface{}
	Topic string
	wg    *sync.WaitGroup
}

// Done calls Done on sync.WaitGroup if set.
func (e *Event) Done() {
	if e.wg != nil {
		e.wg.Done()
	}
}

// CallbackFunc Defines a CallbackFunc.
type CallbackFunc func(topic string, data interface{})

// EventChannel is a channel which can accept an Event.
type EventChannel chan Event

// NewEventChannel Creates a new EventChannel.
func NewEventChannel() EventChannel {
	return make(EventChannel)
}

// dataChannelSlice is a slice of DataChannels.
type eventChannelSlice []EventChannel

// EventBus stores the information about subscribers interested for a particular topic.
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string]eventChannelSlice
	stats       *Stats
}

// NewEventBus returns a new EventBus instance.
func NewEventBus() *EventBus {
	return &EventBus{ //nolint:exhaustivestruct
		subscribers: map[string]eventChannelSlice{},
		stats:       newStats(),
	}
}

// getSubscribingChannels returns all subscribing channels including wildcard matches.
func (eb *EventBus) getSubscribingChannels(topic string) eventChannelSlice {
	subChannels := eventChannelSlice{}

	for topicName := range eb.subscribers {
		if topicName == topic || matchWildcard(topicName, topic) {
			subChannels = append(subChannels, eb.subscribers[topicName]...)
		}
	}

	return subChannels
}

// doPublish is publishing events to channels internally.
func (eb *EventBus) doPublish(channels eventChannelSlice, evt Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	go func(channels eventChannelSlice, evt Event) {
		for _, ch := range channels {
			ch <- evt
		}
	}(channels, evt)
}

// Code from https://github.com/minio/minio/blob/master/pkg/wildcard/match.go
func matchWildcard(pattern, name string) bool {
	if pattern == "" {
		return name == pattern
	}

	if pattern == "*" {
		return true
	}
	// Does only wildcard '*' match.
	return deepMatchRune([]rune(name), []rune(pattern), true)
}

// Code from https://github.com/minio/minio/blob/master/pkg/wildcard/match.go
func deepMatchRune(str, pattern []rune, simple bool) bool { //nolint:unparam
	for len(pattern) > 0 {
		switch pattern[0] {
		default:
			if len(str) == 0 || str[0] != pattern[0] {
				return false
			}
		case '*':
			return deepMatchRune(str, pattern[1:], simple) ||
				(len(str) > 0 && deepMatchRune(str[1:], pattern, simple))
		}

		str = str[1:]

		pattern = pattern[1:]
	}

	return len(str) == 0 && len(pattern) == 0
}

// PublishAsync data to a topic asynchronously
// This function returns a bool channel which indicates that all subscribers where called.
func (eb *EventBus) PublishAsync(topic string, data interface{}) {
	eb.doPublish(
		eb.getSubscribingChannels(topic),
		Event{
			Data:  data,
			Topic: topic,
			wg:    nil,
		})

	eb.stats.incPublishedCountByTopic(topic)
}

// PublishAsyncOnce same as PublishAsync but makes sure that topic is only published once.
func (eb *EventBus) PublishAsyncOnce(topic string, data interface{}) {
	if eb.stats.GetPublishedCountByTopic(topic) > 0 {
		return
	}

	eb.PublishAsync(topic, data)
}

// Publish data to a topic and wait for all subscribers to finish
// This function creates a waitGroup internally. All subscribers must call Done() function on Event.
func (eb *EventBus) Publish(topic string, data interface{}) interface{} {
	wg := sync.WaitGroup{}
	channels := eb.getSubscribingChannels(topic)
	wg.Add(len(channels))
	eb.doPublish(
		channels,
		Event{
			Data:  data,
			Topic: topic,
			wg:    &wg,
		})
	wg.Wait()

	eb.stats.incPublishedCountByTopic(topic)

	return data
}

// PublishOnce same as Publish but makes sure only published once on topic.
func (eb *EventBus) PublishOnce(topic string, data interface{}) interface{} {
	if eb.stats.GetPublishedCountByTopic(topic) > 0 {
		return nil
	}

	return eb.Publish(topic, data)
}

// Subscribe to a topic passing a EventChannel.
func (eb *EventBus) Subscribe(topic string) EventChannel {
	ch := make(EventChannel)
	eb.SubscribeChannel(topic, ch)

	eb.stats.incSubscriberCountByTopic(topic)

	return ch
}

// SubscribeChannel subscribes to a given Channel.
func (eb *EventBus) SubscribeChannel(topic string, ch EventChannel) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]EventChannel{}, ch)
	}

	eb.stats.incSubscriberCountByTopic(topic)
}

// SubscribeCallback provides a simple wrapper that allows to directly register CallbackFunc instead of channels.
func (eb *EventBus) SubscribeCallback(topic string, callable CallbackFunc) {
	ch := NewEventChannel()
	eb.SubscribeChannel(topic, ch)

	go func(callable CallbackFunc) {
		evt := <-ch
		callable(evt.Topic, evt.Data)
		evt.Done()
	}(callable)

	eb.stats.incSubscriberCountByTopic(topic)
}

// HasSubscribers Check if a topic has subscribers.
func (eb *EventBus) HasSubscribers(topic string) bool {
	return len(eb.getSubscribingChannels(topic)) > 0
}

// Stats returns the stats map.
func (eb *EventBus) Stats() *Stats {
	return eb.stats
}
