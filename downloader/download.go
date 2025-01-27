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

type Result struct {
	Body io.ReadCloser
	Err  error
}

func NewErrorResult(err error) Result {
	return Result{Body: nil, Err: err}
}

func downloadAll(urls []string, c chan Result, policy backoff.Policy) {
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			c <- download(url, policy)
		}(url)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
}

func download(url string, p backoff.Policy) Result {
	b, ctx, cancel := p.NewBackoff(context.Background())
	defer cancel()

	m := multierr.New()
	for backoff.Continue(ctx, b) {
		resp, err := http.Get(url)
		if err != nil {
			m.Add(err)
			continue
		}

		// サーバーエラーはリトライを行う
		if resp.StatusCode >= http.StatusInternalServerError {
			body, _ := io.ReadAll(resp.Body)
			err := fmt.Errorf("server error (%d): %s: %s", resp.StatusCode, url, body)
			m.Add(err)
			continue
		}

		return Result{resp.Body, nil}

	}

	return Result{nil, fmt.Errorf("retry failed; got error:\n%v", m.Err())}
}
