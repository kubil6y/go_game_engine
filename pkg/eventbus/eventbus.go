package eventbus

import (
	"reflect"
	"sync"
)

type EventID int
type EventCallback func(any)

type EventBus struct {
	callbacks map[EventID][]EventCallback
	mu        sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		callbacks: make(map[EventID][]EventCallback),
	}
}

func (b *EventBus) On(eventID EventID, callback EventCallback) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, exists := b.callbacks[eventID]
	if !exists {
		b.callbacks[eventID] = make([]EventCallback, 0)
	}

	for _, v := range b.callbacks[eventID] {
		if reflect.ValueOf(v).Pointer() == reflect.ValueOf(callback).Pointer() {
			return false
		}
	}

	b.callbacks[eventID] = append(b.callbacks[eventID], callback)
	return true
}

func (b *EventBus) Off(eventID EventID, callback EventCallback) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	callbacks, exists := b.callbacks[eventID]
	if !exists {
		return false
	}

	for i, v := range callbacks {
		if reflect.ValueOf(v).Pointer() == reflect.ValueOf(callback).Pointer() {
			b.callbacks[eventID] = append(callbacks[:i], callbacks[i+1:]...)
			return true
		}
	}

	return false
}

func (b *EventBus) Emit(eventID EventID, payload any) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	callbacks, exists := b.callbacks[eventID]
	if !exists {
		return false
	}

	for _, callback := range callbacks {
		callback(payload)
	}

	return true
}
