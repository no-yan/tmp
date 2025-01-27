package main

import (
	"flag"
	"fmt"
	"io"
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
	pub.Register(NopSubscriber{})

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

	b, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Body: \n%s", string(b))

	pub.Publish(News{EventEnd})
}
