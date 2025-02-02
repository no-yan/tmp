package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

type DownloadController struct {
	tasks  map[string]Task
	policy *backoff.Policy
	pub    *pubsub.Publisher[Event]
	c      chan bool
	sem    chan int
	wg     *sync.WaitGroup
}

func NewDownloadController(tasks Tasks, policy *backoff.Policy, publisher *pubsub.Publisher[Event]) *DownloadController {
	c := make(chan bool)
	sem := make(chan int, 4)
	wg := sync.WaitGroup{}

	return &DownloadController{
		c:      c,
		sem:    sem,
		policy: policy,
		pub:    publisher,
		wg:     &wg,
		tasks:  tasks,
	}
}

func (dc *DownloadController) Run(ctx context.Context) chan bool {
	for url := range dc.tasks {
		dc.wg.Add(1)
		go func(url string) {
			defer dc.wg.Done()

			// semaphore
			dc.sem <- 1
			defer func() { <-dc.sem }()

			d := NewDownloadWorker(url, dc.policy, dc.pub)
			body, size, err := d.Run(ctx)
			if err != nil {
				dc.pub.Publish(EventAbort{URL: d.url})
				return
			}
			defer body.Close()

			err = os.Mkdir("out", 0o777)
			if err != nil && !os.IsExist(err) {
				panic(err)
			}

			fName := createFileName(d.url)
			path := filepath.Join("out", fName)

			f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o777)
			if err != nil {
				d.pub.Publish(EventAbort{
					URL: d.url,
				})
				return
			}
			defer f.Close()

			tracker := NewProgressTracker(d.url, d.pub)
			r := io.TeeReader(body, tracker)

			n, err := io.Copy(f, r)
			if err != nil {
				dc.pub.Publish(EventAbort{URL: d.url})
				return
			}

			d.pub.Publish(EventEnd{
				TotalSize:   int64(size),
				CurrentSize: n,
				URL:         d.url,
			})
		}(url)
	}

	go func() {
		dc.wg.Wait()
		close(dc.c)
	}()

	return dc.c
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

	d.pub.Publish(EventStart{
		TotalSize:   0,
		CurrentSize: 0,
		URL:         d.url,
	})

	for backoff.Continue(ctx, b) {
		resp, err := http.Get(d.url)
		if err != nil {
			m.Add(err)

			d.pub.Publish(EventRetry{
				TotalSize: 0,
				URL:       d.url,
			})
			continue
		}

		// サーバーエラーはリトライを行う
		if resp.StatusCode >= http.StatusInternalServerError {
			body, _ := io.ReadAll(resp.Body)
			err := fmt.Errorf("server error (%d): %s: %s", resp.StatusCode, d.url, body)
			m.Add(err)

			d.pub.Publish(EventRetry{
				TotalSize: 0,
				URL:       d.url,
			})
			continue
		}

		return resp.Body, int(resp.ContentLength), nil
	}

	return nil, 0, m.Err()
}

func createFileName(url string) string {
	b := sha256.Sum256([]byte(url))
	hex := hex.EncodeToString(b[:])
	return fmt.Sprintf("%s.log", hex)
}

type ProgressTracker struct {
	total int64
	url   string
	size  int64
	pub   *pubsub.Publisher[Event]
}

func NewProgressTracker(url string, pub *pubsub.Publisher[Event]) *ProgressTracker {
	return &ProgressTracker{
		total: 0,
		url:   url,
		size:  0,
		pub:   pub,
	}
}

func (p *ProgressTracker) Write(data []byte) (int, error) {
	n := len(data)
	p.size += int64(n)

	p.pub.Publish(EventProgress{URL: p.url, Current: p.size, Total: p.total})
	return n, nil
}
