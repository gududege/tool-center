package events

import (
	"sync"
	"sync/atomic"

	"github.com/cli-tool-center/tool-center/internal/domain/events"
)

type InMemoryBus struct {
	mu        sync.RWMutex
	handlers  map[string]map[uint64]events.Handler
	idCounter uint64
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string]map[uint64]events.Handler),
	}
}

func (b *InMemoryBus) Subscribe(topic string, handler events.Handler) func() {
	id := atomic.AddUint64(&b.idCounter, 1)

	b.mu.Lock()
	if b.handlers[topic] == nil {
		b.handlers[topic] = make(map[uint64]events.Handler)
	}
	b.handlers[topic][id] = handler
	b.mu.Unlock()

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		delete(b.handlers[topic], id)
	}
}

func (b *InMemoryBus) Publish(event events.DomainEvent) {
	b.mu.RLock()
	topic := event.Topic()
	var handlers []events.Handler
	if m, ok := b.handlers[topic]; ok {
		handlers = make([]events.Handler, 0, len(m))
		for _, h := range m {
			handlers = append(handlers, h)
		}
	}
	b.mu.RUnlock()

	for _, h := range handlers {
		h(event)
	}
}

func (b *InMemoryBus) PublishAsync(event events.DomainEvent) {
	go b.Publish(event)
}
