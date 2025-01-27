package main

import (
	"slices"
	"sync"
)

type Event int

const (
	EventStart Event = iota
	EventProgress
	EventEnd
	EventAbort
)

type News struct {
	Event       Event
	TotalSize   int64
	CurrentSize int64
}

type Subscriber interface {
	HandleEvent(news News)
}

type Publisher struct {
	mu  sync.Mutex
	sub []Subscriber
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Register(s ...Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sub = append(p.sub, s...)
}

func (p *Publisher) Cancel(s Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if i := slices.Index(p.sub, s); i >= 0 {
		p.sub = slices.Delete(p.sub, i, i+1)
	}
}

func (p *Publisher) Publish(news News) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, obs := range p.sub {
		obs.HandleEvent(news)
	}
}

type NopSubscriber struct{}

func (n NopSubscriber) HandleEvent(News) {}
