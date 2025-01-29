package main

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type (
	bars             map[string]*mpb.Bar
	MultiProgressBar struct {
		p    *mpb.Progress
		bars bars
	}
)

func NewMultiProgressBar() *MultiProgressBar {
	p := mpb.New(mpb.WithWidth(64))
	bars := make(bars)

	return &MultiProgressBar{
		p:    p,
		bars: bars,
	}
}

func (p *MultiProgressBar) Flush() {
	p.p.Wait()
}

func (p *MultiProgressBar) CreateBar(title string) *mpb.Bar {
	name := title
	return p.p.New(int64(100),

		mpb.BarStyle().Lbound("╢").Filler("▌").Tip("▌").Padding("░").Rbound("╟"),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			decor.Name("downloading", decor.WCSyncSpaceR),
			decor.OnComplete(decor.Percentage(decor.WC{W: 5}), "done"),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)
}

func (p MultiProgressBar) HandleEvent(news News) {
	switch news.Event {
	case EventStart:
		bar := p.CreateBar(news.URL)
		p.bars[news.URL] = bar

	case EventProgress:
		b, ok := p.findBar(news.URL)
		if !ok {
			panic("bar not found")
		}
		b.Increment()
	case EventEnd:
		b, ok := p.findBar(news.URL)
		if !ok {
			panic("bar not found")
		}
		b.IncrBy(100)
	}
}

func (p *MultiProgressBar) findBar(url string) (bar *mpb.Bar, ok bool) {
	bar, ok = p.bars[url]
	return
}
