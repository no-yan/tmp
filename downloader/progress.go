package main

import (
	"fmt"
	"io"

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
	// TODO: if content-size is unknown, let bar will be spinner.
	return p.p.New(
		int64(100),
		mpb.BarStyle().Lbound("╢").Filler("▌").Tip("▌").Padding("░").Rbound("╟"),
		ClearBarFilerOnFinish(),
		mpb.PrependDecorators(
			decor.Name(title, decor.WC{C: decor.DSyncWidthR | decor.DextraSpace}),
			decor.OnAbort(
				decor.OnComplete(
					decor.Name("downloading", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
					"completed",
				),
				"aborted",
			),
			decor.OnAbort(
				decor.OnComplete(
					decor.Percentage(), "",
				),
				"",
			),
		),
	)
}

func (p MultiProgressBar) HandleEvent(news News) {
	switch news.Event {
	case EventStart:
		bar := p.CreateBar(news.URL)
		p.bars[news.URL] = bar
	case EventProgress:
		b := p.findBar(news.URL)
		b.Increment()
	case EventRetry:
		b := p.findBar(news.URL)
		b.SetCurrent(0)
	case EventEnd:
		b := p.findBar(news.URL)
		b.IncrBy(100)
	case EventAbort:
		b := p.findBar(news.URL)
		b.Abort(false)
	default:
		panic(fmt.Sprintf("unexpected main.Event: %#v", news.Event))
	}
}

func (p *MultiProgressBar) findBar(url string) *mpb.Bar {
	bar, ok := p.bars[url]
	if !ok {
		panic("bar not found")
	}
	return bar
}

func ClearBarFilerOnFinish() mpb.BarOption {
	return barFilterOnFinish("")
}

func barFilterOnFinish(message string) mpb.BarOption {
	return mpb.BarFillerMiddleware(func(base mpb.BarFiller) mpb.BarFiller {
		return mpb.BarFillerFunc(func(w io.Writer, st decor.Statistics) error {
			if st.Completed || st.Aborted {
				_, err := io.WriteString(w, message)
				return err
			}

			return base.Fill(w, st)
		})
	})
}
