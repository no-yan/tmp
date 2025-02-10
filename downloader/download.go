package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/no-yan/multierr"
	"github.com/no-yan/tmp/downloader/internal/backoff"
	"github.com/no-yan/tmp/downloader/internal/pubsub"
)

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

type Task struct {
	url string
}

func NewTask(url string) *Task {
	return &Task{url}
}

type Tasks map[string]Task

func NewTasks(urls ...string) Tasks {
	m := make(Tasks)

	for _, url := range urls {
		m[url] = *NewTask(url)
	}
	return m
}

type Saver interface {
	Save(r io.Reader, url string) (int64, error)
}

type DownloadController struct {
	tasks  map[string]Task
	policy *backoff.Policy
	pub    *pubsub.Publisher[Event]
	sem    chan int
	wg     *sync.WaitGroup
	saver  Saver
}

func NewDownloadController(tasks Tasks, policy *backoff.Policy, publisher *pubsub.Publisher[Event], saver Saver, maxWorkers uint) *DownloadController {
	sem := make(chan int, maxWorkers)
	wg := sync.WaitGroup{}

	return &DownloadController{
		sem:    sem,
		policy: policy,
		pub:    publisher,
		wg:     &wg,
		tasks:  tasks,
		saver:  saver,
	}
}

func (dc *DownloadController) Run(ctx context.Context) {
	for url := range dc.tasks {
		dc.wg.Add(1)
		go func(url string) {
			defer dc.wg.Done()

			// semaphore
			dc.sem <- 1

			d := NewDownloadWorker(url, dc.policy, dc.pub)
			body, size, err := d.Run(ctx)
			if err != nil {
				<-dc.sem
				dc.pub.PublishWithContext(ctx, NewEventAbort(d.url, err))
				return
			}
			defer body.Close()

			<-dc.sem

			tracker := NewProgressTracker(url, d.pub, int64(size))
			r := io.TeeReader(body, tracker)

			n, err := dc.saver.Save(r, d.url)
			if err != nil {
				dc.pub.PublishWithContext(ctx, NewEventAbort(d.url, err))
				return
			}

			d.pub.PublishWithContext(ctx, EventEnd{
				TotalSize:   int64(size),
				CurrentSize: n,
				URL:         d.url,
			})
		}(url)
	}

	dc.wg.Wait()
}

type DownloadWorker struct {
	url    string
	policy *backoff.Policy
	pub    *pubsub.Publisher[Event]
}

func NewDownloadWorker(url string, policy *backoff.Policy, publisher *pubsub.Publisher[Event]) *DownloadWorker {
	return &DownloadWorker{url: url, policy: policy, pub: publisher}
}

func (d *DownloadWorker) Run(ctx context.Context) (body io.ReadCloser, contentLength int, err error) {
	b := d.policy.NewBackoff()

	m := multierr.New()

	d.pub.PublishWithContext(ctx, EventStart{
		TotalSize:   0,
		CurrentSize: 0,
		URL:         d.url,
	})

	for backoff.Continue(ctx, b) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.url, nil)
		if err != nil {
			return nil, 0, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.Add(err)

			d.pub.PublishWithContext(ctx, EventRetry{
				TotalSize: 0,
				URL:       d.url,
			})
			continue
		}

		// サーバーエラーはリトライを行う
		if resp.StatusCode >= http.StatusInternalServerError {

			body, _ := io.ReadAll(io.LimitReader(resp.Body, 64))
			resp.Body.Close()
			err := fmt.Errorf("server error (%d):  %s", resp.StatusCode, body)
			m.Add(err)

			d.pub.PublishWithContext(ctx, EventRetry{
				TotalSize: 0,
				URL:       d.url,
			})
			continue
		}

		return resp.Body, int(resp.ContentLength), nil
	}

	err = m.Err()
	if err != nil {
		return nil, 0, err
	}
	// net/http同様、必ずBodyがCloseできるようにする
	return io.NopCloser(strings.NewReader("")), 0, nil
}

type ProgressTracker struct {
	current, total int64
	url            string
	pub            *pubsub.Publisher[Event]
}

func NewProgressTracker(url string, pub *pubsub.Publisher[Event], total int64) *ProgressTracker {
	return &ProgressTracker{
		current: 0,
		total:   total,
		url:     url,
		pub:     pub,
	}
}

func (p *ProgressTracker) Write(data []byte) (int, error) {
	n := len(data)
	p.current += int64(n)

	p.pub.Publish(EventProgress{Current: p.current, Total: p.total, URL: p.url})
	return n, nil
}
