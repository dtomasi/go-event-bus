package eventbus

import (
	"github.com/minio/minio/pkg/wildcard"
	"sync"
)

// Event holds topic name and data
type Event struct {
	Data  interface{}
	Topic string
	wg    *sync.WaitGroup
}

// Done calls Done on sync.WaitGroup if set
func (e *Event) Done() {
	if e.wg != nil {
		e.wg.Done()
	}
}

// Defines a CallbackFunc
type CallbackFunc func(topic string, data interface{})

// EventChannel is a channel which can accept an Event
type EventChannel chan Event

// Creates a new EventChannel
func NewEventChannel() EventChannel {
	return make(EventChannel)
}

// dataChannelSlice is a slice of DataChannels
type eventChannelSlice []EventChannel

// EventBus stores the information about subscribers interested for a particular topic
type EventBus struct {
	subscribers map[string]eventChannelSlice
	rm          sync.RWMutex
}

// NewEventBus returns a new EventBus instance
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: map[string]eventChannelSlice{},
	}
}

// Singleton Bus instance
var defaultBus *EventBus

// GetBus returns the default EventBus instance
func GetBus() *EventBus {
	if defaultBus == nil {
		defaultBus = NewEventBus()
	}
	return defaultBus
}

// getSubscribingChannels returns all subscribing channels including wildcard matches
func (eb *EventBus) getSubscribingChannels(topic string) eventChannelSlice {
	subChannels := eventChannelSlice{}
	for topicName := range eb.subscribers {
		if topicName == topic || wildcard.MatchSimple(topicName, topic) {
			subChannels = append(subChannels, eb.subscribers[topicName]...)
		}
	}
	return subChannels
}

// doPublish is publishing events to channels internally
func (eb *EventBus) doPublish(channels eventChannelSlice, evt Event) {
	eb.rm.RLock()
	go func(channels eventChannelSlice, evt Event) {
		for _, ch := range channels {
			ch <- evt
		}
	}(channels, evt)
	eb.rm.RUnlock()
}

// PublishAsync data to a topic asynchronously
// This function returns a bool channel which indicates that all subscribers where called
func (eb *EventBus) PublishAsync(topic string, data interface{}) {
	eb.doPublish(
		eb.getSubscribingChannels(topic),
		Event{
			Data:  data,
			Topic: topic,
		})
}

// Publish data to a topic and wait for all subscribers to finish
// This function creates a waitGroup internally. All subscribers must call Done() function on Event
func (eb *EventBus) Publish(topic string, data interface{}) {
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
}

// Subscribe to a topic passing a EventChannel
func (eb *EventBus) Subscribe(topic string, ch EventChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]EventChannel{}, ch)
	}
	eb.rm.Unlock()
}

// SubscribeCallback provides a simple wrapper that allows to directly register CallbackFunc instead of channels
func (eb *EventBus) SubscribeCallback(topic string, callable CallbackFunc) {
	ch := NewEventChannel()
	eb.Subscribe(topic, ch)
	go func(callable CallbackFunc) {
		evt := <-ch
		callable(evt.Topic, evt.Data)
		evt.Done()
	}(callable)
}
