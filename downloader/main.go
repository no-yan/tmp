package main

import (
	"context"
	"flag"
	"time"

	"github.com/no-yan/tmp/downloader/internal/backoff"
	"github.com/no-yan/tmp/downloader/internal/pubsub"
)

var defaultPolicy = backoff.Policy{
	DelayMin:   100 * time.Millisecond,
	DelayMax:   500 * time.Millisecond,
	RetryLimit: 10,
}

func main() {
	flag.Parse()
	args := flag.Args()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pub := pubsub.NewPublisher[Event]()
	bar := NewMultiProgressBar()
	nop := NopSubscriber{}
	pub.Register(bar, nop)

	tasks := NewTasks(args...)
	dc := NewDownloadController(tasks, &defaultPolicy, pub)
	c := dc.Run(ctx)
	<-c

	bar.Flush()
}
