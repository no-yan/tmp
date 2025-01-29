package main

import (
	"context"
	"flag"
	"io"
	"time"

	"github.com/no-yan/tmp/downloader/internal/backoff"
	"github.com/no-yan/tmp/downloader/internal/pubsub"
)

var defaultPolicy = backoff.Policy{
	DelayMin:   time.Millisecond,
	DelayMax:   20 * time.Millisecond,
	RetryLimit: 10,
}

func main() {
	flag.Parse()
	args := flag.Args()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pub := pubsub.NewPublisher[News]()
	progressBar := NewMultiProgressBar()
	pub.Register(NopSubscriber{}, progressBar)

	tasks := NewTasks(args...)
	dc := NewDownloadController(tasks, &defaultPolicy, pub)
	c := dc.Run(ctx)

	for range c {
	}
}

type Printer struct {
	results []Result
	r       io.Reader
	w       io.Writer
}
