package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/no-yan/tmp/downloader/internal/backoff"
	"github.com/no-yan/tmp/downloader/internal/pubsub"
)

var defaultPolicy = backoff.Policy{
	DelayMin:   10 * time.Millisecond,
	DelayMax:   50 * time.Millisecond,
	RetryLimit: 10,
}

const outDir = "out"

func main() {
	flag.Parse()
	args := flag.Args()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pub := pubsub.NewPublisher[Event]()
	bar := NewMultiProgressBar()
	printer := NewResult(os.Stdout, outDir)
	nop := NopSubscriber{}
	pub.Register(bar, nop, printer)

	tasks := NewTasks(args...)
	saver := NewFileSaver(outDir)
	dc := NewDownloadController(tasks, &defaultPolicy, pub, saver)
	c := dc.Run(ctx)
	<-c

	bar.Flush()
	printer.Show()
}
