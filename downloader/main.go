package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/no-yan/tmp/downloader/internal/backoff"
	"github.com/no-yan/tmp/downloader/internal/pubsub"
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

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	pub := pubsub.NewPublisher[News]()
	progressBar := NewProgressBar(args[0], os.Stdout)
	pub.Register(NopSubscriber{}, progressBar)

	tasks := NewTasks(args...)
	dc := NewDownloadController(tasks, &defaultPolicy, pub)
	downloadChannel := dc.Run(ctx)

	for result := range downloadChannel {
		print(result, pub)
	}
}

func print(r Result, pub *pubsub.Publisher[News]) {
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

	pub.Publish(News{EventEnd, 100, 100})
}
