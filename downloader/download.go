package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/no-yan/multierr"
	"github.com/no-yan/tmp/downloader/internal/backoff"
	"github.com/no-yan/tmp/downloader/internal/pubsub"
)

type Event int

const (
	EventStart Event = iota
	EventProgress
	EventRetry
	EventEnd
	EventAbort
)

type News struct {
	Event       Event
	TotalSize   int64
	CurrentSize int64
	URL         string
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
	pub    *pubsub.Publisher[News]
	c      chan bool
	sem    chan int
	wg     *sync.WaitGroup
}

func NewDownloadController(tasks Tasks, policy *backoff.Policy, publisher *pubsub.Publisher[News]) *DownloadController {
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
				dc.pub.Publish(News{
					Event:       EventAbort,
					TotalSize:   int64(size),
					CurrentSize: int64(size),
					URL:         d.url,
				})
				return
			}
			defer body.Close()

			b, err := io.ReadAll(body)
			if err != nil {
				panic(err)
			}

			d.pub.Publish(News{
				Event:       EventEnd,
				TotalSize:   int64(size),
				CurrentSize: int64(len(b)),
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
	pub    *pubsub.Publisher[News]
}

func NewDownloadWorker(url string, policy *backoff.Policy, publisher *pubsub.Publisher[News]) *DownloadWorker {
	return &DownloadWorker{url: url, policy: policy, pub: publisher}
}

func (d *DownloadWorker) Run(ctx context.Context) (body io.ReadCloser, contentLength int, err error) {
	b := d.policy.NewBackoff()

	m := multierr.New()

	d.pub.Publish(News{
		Event:       EventStart,
		TotalSize:   0,
		CurrentSize: 0,
		URL:         d.url,
	})

	for backoff.Continue(ctx, b) {
		resp, err := http.Get(d.url)
		if err != nil {
			m.Add(err)

			d.pub.Publish(News{
				Event:       EventRetry,
				TotalSize:   0,
				CurrentSize: 0,
				URL:         d.url,
			})
			continue
		}

		// サーバーエラーはリトライを行う
		if resp.StatusCode >= http.StatusInternalServerError {
			body, _ := io.ReadAll(resp.Body)
			err := fmt.Errorf("server error (%d): %s: %s", resp.StatusCode, d.url, body)
			m.Add(err)

			d.pub.Publish(News{
				Event:       EventRetry,
				TotalSize:   0,
				CurrentSize: 0,
				URL:         d.url,
			})
			continue
		}

		return resp.Body, int(resp.ContentLength), nil
	}

	return nil, 0, m.Err()
}
