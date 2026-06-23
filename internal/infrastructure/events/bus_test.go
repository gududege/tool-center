package events

import (
	"sync"
	"testing"
	"time"

	domainEvents "github.com/cli-tool-center/tool-center/internal/domain/events"
)

func TestInMemoryBus_Publish(t *testing.T) {
	bus := NewInMemoryBus()
	received := make([]domainEvents.DomainEvent, 0)
	var mu sync.Mutex

	bus.Subscribe("test.event", func(event domainEvents.DomainEvent) {
		mu.Lock()
		received = append(received, event)
		mu.Unlock()
	})

	bus.Publish(testEvent{id: "1"})

	if len(received) != 1 {
		t.Fatalf("expected 1 event, got %d", len(received))
	}
}

func TestInMemoryBus_Unsubscribe(t *testing.T) {
	bus := NewInMemoryBus()
	count := 0

	cancel := bus.Subscribe("test.event", func(event domainEvents.DomainEvent) {
		count++
	})

	cancel()
	bus.Publish(testEvent{id: "2"})

	if count != 0 {
		t.Error("expected handler to be removed after cancel")
	}
}

func TestInMemoryBus_MultipleSubscribers(t *testing.T) {
	bus := NewInMemoryBus()
	var mu sync.Mutex
	count := 0

	bus.Subscribe("test.event", func(event domainEvents.DomainEvent) {
		mu.Lock()
		count++
		mu.Unlock()
	})
	bus.Subscribe("test.event", func(event domainEvents.DomainEvent) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	bus.Publish(testEvent{id: "3"})

	mu.Lock()
	if count != 2 {
		t.Errorf("expected 2 handler calls, got %d", count)
	}
	mu.Unlock()
}

func TestInMemoryBus_TopicFiltering(t *testing.T) {
	bus := NewInMemoryBus()
	received := make([]string, 0)
	var mu sync.Mutex

	bus.Subscribe("topic.a", func(event domainEvents.DomainEvent) {
		mu.Lock()
		received = append(received, "a")
		mu.Unlock()
	})
	bus.Subscribe("topic.b", func(event domainEvents.DomainEvent) {
		mu.Lock()
		received = append(received, "b")
		mu.Unlock()
	})

	bus.Publish(testEvent{topic: "topic.a"})

	mu.Lock()
	if len(received) != 1 || received[0] != "a" {
		t.Errorf("expected only 'a', got %v", received)
	}
	mu.Unlock()
}

func TestInMemoryBus_PublishAsync(t *testing.T) {
	bus := NewInMemoryBus()
	received := make(chan domainEvents.DomainEvent, 10)

	bus.Subscribe("test.event", func(event domainEvents.DomainEvent) {
		received <- event
	})

	bus.PublishAsync(testEvent{id: "async"})

	var e domainEvents.DomainEvent
	select {
	case e = <-received:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for async event")
	}
	if e.Topic() != "test.event" {
		t.Errorf("expected topic 'test.event', got '%s'", e.Topic())
	}
}

type testEvent struct {
	id    string
	topic string
}

func (e testEvent) Topic() string {
	if e.topic != "" {
		return e.topic
	}
	return "test.event"
}
