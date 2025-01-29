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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pub := pubsub.NewPublisher[News]()
	progressBar := NewProgressBar(args[0], os.Stdout)
	printer := NewPrinter(os.Stdout)
	pub.Register(NopSubscriber{}, progressBar, printer)

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

func NewPrinter(dst io.Writer) *Printer {
	return &Printer{w: dst}
}

func (p *Printer) HandleEvent(news News) {
	var msg string
	switch news.Event {
	case EventStart:
		msg = fmt.Sprintf("Start downloading: %s %d\n", news.URL, news.CurrentSize)
	case EventProgress:
	case EventEnd:
		msg = fmt.Sprintf("Finished downloading: %s\n", news.URL)
	case EventAbort:
		msg = fmt.Sprintf("Aborted: %s\n", news.URL)
	default:
		panic(fmt.Sprintf("unexpected main.Event: %#v\n", news.Event))
	}

	io.WriteString(p.w, msg)
}
