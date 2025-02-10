package main

type EventType int

type Event interface {
	Type() EventType
}

const (
	EventTypeStart EventType = iota
	EventTypeProgress
	EventTypeRetry
	EventTypeEnd
	EventTypeAbort
)

type EventStart struct {
	TotalSize   int64
	CurrentSize int64
	URL         string
}

func (e EventStart) Type() EventType {
	return EventTypeStart
}

type EventProgress struct {
	URL     string
	Current int64
	Total   int64
}

func (e EventProgress) Type() EventType {
	return EventTypeProgress
}

type EventRetry struct {
	TotalSize int64
	URL       string
}

func (e EventRetry) Type() EventType {
	return EventTypeRetry
}

type EventEnd struct {
	TotalSize   int64
	CurrentSize int64
	URL         string
}

func (e EventEnd) Type() EventType {
	return EventTypeEnd
}

type EventAbort struct {
	URL string
	Err error
}

func NewEventAbort(url string, err error) EventAbort {
	return EventAbort{
		URL: url,
		Err: err,
	}
}

func (e EventAbort) Type() EventType {
	return EventTypeAbort
}
