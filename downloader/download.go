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
			d.Run(ctx) // TODO: エラーハンドリング
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

func (d *DownloadWorker) Run(ctx context.Context) error {
	b := d.policy.NewBackoff()

	m := multierr.New()

	for backoff.Continue(ctx, b) {
		resp, err := http.Get(d.url)
		if err != nil {
			m.Add(err)
			continue
		}

		d.pub.Publish(News{
			Event:       EventStart,
			TotalSize:   max(0, resp.ContentLength),
			CurrentSize: 0,
			URL:         d.url,
		})

		// サーバーエラーはリトライを行う
		if resp.StatusCode >= http.StatusInternalServerError {
			body, _ := io.ReadAll(resp.Body)
			err := fmt.Errorf("server error (%d): %s: %s", resp.StatusCode, d.url, body)
			m.Add(err)
			continue
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()

		d.pub.Publish(News{
			Event:       EventEnd,
			TotalSize:   resp.ContentLength,
			CurrentSize: int64(len(b)),
			URL:         d.url,
		})
		return nil
	}

	return m.Err()
}
