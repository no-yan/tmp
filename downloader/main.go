package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func main() {
	flag.Parse()
	args := flag.Args()

	fmt.Printf("URL: %s\n", args)

	c := make(chan Result)

	downloadAll(args, c)

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

func downloadAll(urls []string, c chan Result) {
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			c <- download(url)
		}(url)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
}

func download(url string) Result {
	resp, err := http.Get(url)
	if err != nil {
		return NewErrorResult(err)
	}

	return Result{Body: resp.Body, Err: nil}
}
