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
	c      chan Result
	sem    chan int
	wg     *sync.WaitGroup
}

func NewDownloadController(tasks Tasks, policy *backoff.Policy, publisher *pubsub.Publisher[News]) *DownloadController {
	c := make(chan Result)
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

func (dc *DownloadController) Run(ctx context.Context) chan Result {
	for url := range dc.tasks {
		dc.wg.Add(1)
		go func(url string) {
			defer dc.wg.Done()

			// semaphore
			dc.sem <- 1
			defer func() { <-dc.sem }()

			d := NewDownloadWorker(url, dc.policy, dc.pub)
			dc.c <- d.Run(ctx)
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

func (d *DownloadWorker) Run(ctx context.Context) Result {
	b, ctx, cancel := d.policy.NewBackoff(ctx)
	defer cancel()

	fmt.Println(d.url)
	d.pub.Publish(News{EventStart, 0, 0})
	m := multierr.New()

	for backoff.Continue(ctx, b) {
		resp, err := http.Get(d.url)
		if err != nil {
			m.Add(err)
			continue
		}

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

		d.pub.Publish(News{EventEnd, resp.ContentLength, 100})
		return Result{b, nil}
	}

	return Result{nil, fmt.Errorf("retry failed; got error:\n%v", m.Err())}
}

type Result struct {
	Body []byte
	Err  error
}

func NewErrorResult(err error) Result {
	return Result{Body: nil, Err: err}
}
