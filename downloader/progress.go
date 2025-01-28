package main

import (
	"fmt"
	"io"
)

type ProgressBar struct {
	title string
	w     io.Writer
}

func NewProgressBar(title string, w io.Writer) *ProgressBar {
	return &ProgressBar{title, w}
}

func (p ProgressBar) HandleEvent(news News) {
	switch news.Event {
	case EventStart:
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
	text := fmt.Sprintf("\r%s [%d/%d MB] (%.1f%%)", p.title, current, total, percent)
	io.WriteString(p.w, text)
}
