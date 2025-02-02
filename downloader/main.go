package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/no-yan/tmp/downloader/internal/backoff"
	"github.com/no-yan/tmp/downloader/internal/pubsub"
)

var defaultPolicy = backoff.Policy{
	DelayMin:   10 * time.Millisecond,
	DelayMax:   50 * time.Millisecond,
	RetryLimit: 10,
}

func main() {
	config := NewConfigFromFlags()
	args := flag.Args()

	ctx, cancel := context.WithTimeout(context.Background(), config.timeout)
	defer cancel()

	ctx, stop := setupSignalContext(ctx)
	defer stop()

	pub := pubsub.NewPublisher[Event]()
	bar := NewMultiProgressBar()
	printer := NewResult(os.Stdout, config.outputDir)
	nop := NopSubscriber{}
	pub.Register(bar, nop, printer)

	tasks := NewTasks(args...)
	saver := NewFileSaver(config.outputDir)
	dc := NewDownloadController(tasks, &defaultPolicy, pub, saver, config.workers)
	c := dc.Run(ctx)
	<-c

	bar.Flush()
	printer.Show()
}

func setupSignalContext(parent context.Context) (ctx context.Context, stop context.CancelFunc) {
	ctx, stop = signal.NotifyContext(parent, os.Interrupt)
	return
}
