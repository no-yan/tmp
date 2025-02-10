package main

import (
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), config.timeout)
	defer cancel()

	ctx, stop := setupSignalContext(ctx)
	defer stop()

	pub := pubsub.NewPublisher[Event]()
	bar := NewMultiProgressBar(ctx)
	printer := NewPrinter(os.Stdout, config.outputDir)
	pub.Register(bar, printer)

	saver := NewFileSaver(config.outputDir)
	dc := NewDownloadController(config.tasks, &defaultPolicy, pub, saver, config.workers)
	c := dc.Run(ctx)
	<-c

	bar.Flush()
	printer.Print()
}

func setupSignalContext(parent context.Context) (ctx context.Context, stop context.CancelFunc) {
	ctx, stop = signal.NotifyContext(parent, os.Interrupt)
	return
}
