package main

import (
	"context"
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

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	pub := NewPublisher[News]()
	progressBar := NewProgressBar(args[0], os.Stdout)
	pub.Register(NopSubscriber{}, progressBar)

	downloadChannel := downloadAll(ctx, args, &defaultPolicy, pub)

	for result := range downloadChannel {
		print(result, pub)
	}
}

func print(r Result, pub *Publisher[News]) {
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
