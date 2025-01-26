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

	wg := sync.WaitGroup{}
	c := make(chan Result)

	go downloadAndSend(args, c, &wg)
	wg.Wait()

	for result := range c {
		fmt.Println("=============================")
		if result.Err != nil {
			fmt.Printf("Error: %v\n", result.Err)
			continue
		}

		go func() {
			defer result.Body.Close()

			b, err := io.ReadAll(result.Body)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Body: \n%s", string(b))
		}()
	}
}

type Result struct {
	Body io.ReadCloser
	Err  error
}

func NewErrorResult(err error) Result {
	return Result{Body: nil, Err: err}
}

func downloadAndSend(urls []string, c chan Result, wg *sync.WaitGroup) {
	wg.Add(len(urls))
	for _, url := range urls {
		go func() {
			c <- download(url)
			defer wg.Done()
		}()
	}
}

func download(url string) Result {
	resp, err := http.Get(url)
	if err != nil {
		return NewErrorResult(err)
	}

	return Result{Body: resp.Body, Err: nil}
}
