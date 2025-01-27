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
)

type News struct {
	Message string
	Event   Event
}

type Subscriber interface {
	Update(news News)
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
	for _, sub := range s {
		p.sub = append(p.sub, sub)
	}
}

func (p *Publisher) Cancel(s Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if i := slices.Index(p.sub, s); i > 0 {
		slices.Delete(p.sub, i, i+1)
	}
}

func (p *Publisher) Publish(news News) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, obs := range p.sub {
		obs.Update(news)
	}
}

type NopSubscriber struct{}

func (n NopSubscriber) Update(News) {}
