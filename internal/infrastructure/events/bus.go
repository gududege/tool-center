package events

import (
	"sync"

	"github.com/cli-tool-center/tool-center/internal/domain/events"
)

// InMemoryBus 基于内存的事件总线实现
type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]events.Handler
}

// NewInMemoryBus 创建事件总线
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string][]events.Handler),
	}
}

// Subscribe 订阅事件，返回取消订阅函数
func (b *InMemoryBus) Subscribe(topic string, handler events.Handler) func() {
	b.mu.Lock()
	b.handlers[topic] = append(b.handlers[topic], handler)
	// 记录索引以便取消
	idx := len(b.handlers[topic]) - 1
	b.mu.Unlock()

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		handlers := b.handlers[topic]
		if idx < len(handlers) {
			b.handlers[topic] = append(handlers[:idx], handlers[idx+1:]...)
		}
	}
}

// Publish 同步发布事件
func (b *InMemoryBus) Publish(event events.DomainEvent) {
	b.mu.RLock()
	handlers := b.handlers[event.Topic()]
	b.mu.RUnlock()

	for _, h := range handlers {
		h(event)
	}
}

// PublishAsync 异步发布事件
func (b *InMemoryBus) PublishAsync(event events.DomainEvent) {
	go b.Publish(event)
}
