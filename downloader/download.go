package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/no-yan/multierr"
	"github.com/no-yan/tmp/downloader/internal/backoff"
)

func downloadAll(urls []string, policy *backoff.Policy, publisher *Publisher) chan Result {
	c := make(chan Result, 1)
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			d := NewDownloadWorker(url, policy, publisher)
			c <- d.Run()
		}(url)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	return c
}

type DownloadWorker struct {
	url    string
	policy *backoff.Policy
	pub    *Publisher
}

func NewDownloadWorker(url string, policy *backoff.Policy, publisher *Publisher) *DownloadWorker {
	return &DownloadWorker{url: url, policy: policy, pub: publisher}
}

func (d *DownloadWorker) Run() Result {
	b, ctx, cancel := d.policy.NewBackoff(context.Background())
	defer cancel()

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

		d.pub.Publish(News{EventStart, resp.ContentLength, 0})

		return Result{resp.Body, nil}
	}

	return Result{nil, fmt.Errorf("retry failed; got error:\n%v", m.Err())}
}

type Result struct {
	Body io.ReadCloser
	Err  error
}

func NewErrorResult(err error) Result {
	return Result{Body: nil, Err: err}
}
