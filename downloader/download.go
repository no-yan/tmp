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

type Downloader struct {
	url    string
	policy *backoff.Policy
	pub    *Publisher
}

func NewDownloader(url string, policy *backoff.Policy, publisher *Publisher) *Downloader {
	return &Downloader{url: url, policy: policy, pub: publisher}
}

func (d *Downloader) Run() Result {
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

func downloadAll(urls []string, c chan Result, policy *backoff.Policy, publisher *Publisher) {
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			d := NewDownloader(url, policy, publisher)
			c <- d.Run()
		}(url)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
}
