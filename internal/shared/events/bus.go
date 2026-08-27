package events

import "sync"

type Handler func(payload any)

type subscription struct {
	id      uint64
	handler Handler
}

type Bus struct {
	mu   sync.RWMutex
	subs map[string][]subscription
	next uint64
}

func NewBus() *Bus {
	return &Bus{subs: make(map[string][]subscription)}
}

func (b *Bus) Subscribe(topic string, handler Handler) func() {
	b.mu.Lock()
	b.next++
	id := b.next
	b.subs[topic] = append(b.subs[topic], subscription{id: id, handler: handler})
	b.mu.Unlock()

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		list := b.subs[topic]
		for i := range list {
			if list[i].id == id {
				b.subs[topic] = append(list[:i], list[i+1:]...)
				return
			}
		}
	}
}

func (b *Bus) Publish(topic string, payload any) {
	b.mu.RLock()
	list := append([]subscription(nil), b.subs[topic]...)
	b.mu.RUnlock()
	for _, sub := range list {
		sub.handler(payload)
	}
}
