package main

import (
	"context"
	"fmt"
	"io"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type bars map[string]*mpb.Bar

type MultiProgressBar struct {
	p    *mpb.Progress
	bars bars
}

func NewMultiProgressBar(ctx context.Context) *MultiProgressBar {
	p := mpb.NewWithContext(ctx, mpb.WithWidth(64))
	bars := make(bars)

	return &MultiProgressBar{
		p:    p,
		bars: bars,
	}
}

func (p *MultiProgressBar) Flush() {
	p.p.Wait()
	p.clear()
}

func (p *MultiProgressBar) clear() {
	linesToDelete := len(p.bars)

	for range linesToDelete {
		fmt.Printf("\033[F\033[K")
	}
}

func (p *MultiProgressBar) CreateBar(title string) *mpb.Bar {
	// TODO: if content-size is unknown, let bar will be spinner.
	return p.p.New(
		int64(100),
		mpb.BarStyle().Lbound("╢").Filler("▌").Tip("▌").Padding("░").Rbound("╟"),
		clearBarFillerOnFinish(),
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

func (p *MultiProgressBar) HandleEvent(event Event) {
	switch e := event.(type) {
	case EventStart:
		bar := p.CreateBar(e.URL)
		p.bars[e.URL] = bar
	case EventProgress:
		b := p.findBar(e.URL)
		if e.Total > 0 {
			b.SetTotal(e.Total, false)
		}
		b.SetCurrent(e.Current)
	case EventRetry:
		b := p.findBar(e.URL)
		b.SetCurrent(0)
	case EventEnd:
		b := p.findBar(e.URL)
		b.SetCurrent(e.CurrentSize)
	case EventAbort:
		b := p.findBar(e.URL)
		b.Abort(false)
	default:
		panic(fmt.Sprintf("unexpected main.Event: %#v", e))
	}
}

func (p *MultiProgressBar) findBar(url string) *mpb.Bar {
	bar, ok := p.bars[url]
	if !ok {
		panic("bar not found")
	}
	return bar
}

func clearBarFillerOnFinish() mpb.BarOption {
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
