package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/no-yan/multierr"
	"github.com/no-yan/tmp/downloader/internal/backoff"
)

var defaultPolicy = backoff.Policy{
	DelayMin:   time.Millisecond,
	DelayMax:   20 * time.Millisecond,
	Timeout:    10 * time.Second,
	RetryLimit: 10,
}

func main() {
	flag.Parse()
	args := flag.Args()

	fmt.Printf("URL: %s\n", args)

	c := make(chan Result)

	downloadAll(args, c, defaultPolicy)

	for result := range c {
		fmt.Println("=============================")
		if result.Err != nil {
			fmt.Printf("Error: %v\n", result.Err)
			continue
		}

		b, err := io.ReadAll(result.Body)
		result.Body.Close()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Body: \n%s", string(b))
	}
}

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

	return Result{nil, fmt.Errorf("retry failed; got error\n%v", m.Err())}
}
