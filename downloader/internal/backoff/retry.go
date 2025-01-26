package backoff

import (
	"context"
	"time"
)

type Policy struct {
	DelayMin   time.Duration
	DelayMax   time.Duration
	Timeout    time.Duration
	RetryLimit uint
}

func (p Policy) Next(cnt uint) time.Duration {
	delayMS := int64(1) << cnt

	switch {
	case delayMS < p.DelayMin.Milliseconds():
		return p.DelayMin
	case delayMS > p.DelayMax.Milliseconds():
		return p.DelayMax
	default:
		return time.Duration(delayMS) * time.Millisecond
	}
}

func (p Policy) NewBackoff(c context.Context) (backoff *Backoff, ctx context.Context, cancelFunc func()) {
	ctx, cancelFunc = context.WithTimeout(c, p.Timeout)
	backoff = &Backoff{p, 0}
	return
}

type Backoff struct {
	p   Policy
	cnt uint
}

func (b *Backoff) NextTick() time.Duration {
	delay := b.p.Next(b.cnt)
	return delay
}

func (b *Backoff) LimitExceeded() bool {
	return b.cnt >= b.p.RetryLimit
}

func Continue(ctx context.Context, b *Backoff) bool {
	if b.LimitExceeded() {
		return false
	}

	select {
	case <-ctx.Done():
		return false
	case <-time.After(b.NextTick()):
		b.cnt++
		return true
	}
}
