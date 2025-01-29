package pubsub

import (
	"slices"
	"sync"
)

type Subscriber[T any] interface {
	HandleEvent(event T)
}

type Publisher[T any] struct {
	mu  sync.Mutex
	sub []Subscriber[T]
}

func NewPublisher[T any]() *Publisher[T] {
	return &Publisher[T]{}
}

func (p *Publisher[T]) Register(s ...Subscriber[T]) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sub = append(p.sub, s...)
}

func (p *Publisher[T]) Cancel(s Subscriber[T]) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if i := slices.Index(p.sub, s); i >= 0 {
		p.sub = slices.Delete(p.sub, i, i+1)
	}
}

func (p *Publisher[T]) Publish(event T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, obs := range p.sub {
		obs.HandleEvent(event)
	}
}
