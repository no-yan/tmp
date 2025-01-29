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
	pub.Register(NopSubscriber{}, progressBar)

	tasks := NewTasks(args...)
	dc := NewDownloadController(tasks, &defaultPolicy, pub)
	c := dc.Run(ctx)

	<-c
}

type Printer struct {
	results []Result
	r       io.Reader
	w       io.Writer
	pub     *pubsub.Publisher[News]
}

func NewPrinter(dst io.Writer, pub *pubsub.Publisher[News]) *Printer {
	return &Printer{w: dst, pub: pub}
}

func (p *Printer) Print(b []byte) {
	p.w.Write(b)
}
