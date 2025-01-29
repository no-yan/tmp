package main

type NopSubscriber struct{}

func (n NopSubscriber) HandleEvent(event News) {}
