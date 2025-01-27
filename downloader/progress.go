package main

import (
	"fmt"
	"io"
)

func NewProgressReader(r io.ReadCloser, current int, total int, pub *Publisher) *ProgressReader {
	return &ProgressReader{r: r, current: current, total: total, pub: pub}
}

type ProgressReader struct {
	r       io.ReadCloser
	current int
	total   int
	pub     *Publisher
}

func (p ProgressReader) Read(n []byte) (int, error) {
	nn, err := p.r.Read(n)
	p.current += nn

	p.pub.Publish(News{
		Event:       EventProgress,
		TotalSize:   int64(p.total),
		CurrentSize: int64(p.current),
	})

	return nn, err
}

func (p ProgressReader) Close() error {
	return p.r.Close()
}

type ProgressBar struct {
	w io.Writer
}

func (p ProgressBar) HandleEvent(news News) {
	switch news.Event {
	case EventStart:
		fmt.Println("Start")
	case EventProgress:
		p.Render(int(news.CurrentSize), int(news.TotalSize))
	case EventEnd:
		p.Clear()
	}
}

func (p ProgressBar) Clear() {
}

func (p ProgressBar) Render(current, total int) {
	var percent float64
	if total > 0 {
		percent = float64(current) / float64(total) * 100
	} else {
		percent = 0
	}
	text := fmt.Sprintf("\r%s [%d/%d MB] (%f%%)", "title", current, total, percent)
	io.WriteString(p.w, text)
}
