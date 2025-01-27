package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

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
	pub := NewPublisher()
	progressBar := NewProgressBar(args[0], os.Stdout)
	pub.Register(NopSubscriber{}, progressBar)

	downloadAll(args, c, &defaultPolicy, pub)

	for result := range c {
		print(result, pub)
	}
}

func print(r Result, pub *Publisher) {
	if r.Err != nil {
		fmt.Printf("Error: %v\n", r.Err)
		return
	}

	_, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Body: \n%s", string(b))

	pub.Publish(News{EventEnd, 100, 100})
}
