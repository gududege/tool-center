package events

// Handler 事件处理器
type Handler func(event DomainEvent)

// Bus 事件总线接口
type Bus interface {
	Publish(event DomainEvent)
	Subscribe(topic string, handler Handler) func()
	PublishAsync(event DomainEvent)
}
