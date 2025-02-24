package backoff

import (
	"context"
	"math/bits"
	"time"
)

type Policy struct {
	DelayMin   time.Duration
	DelayMax   time.Duration
	RetryLimit uint
}

func (p Policy) Next(cnt uint) time.Duration {
	if cnt < 1 {
		return 0
	}

	leadingZeros := bits.LeadingZeros(uint(p.DelayMin))
	// time.Durationはint(2の補数表現)なので、
	// 先頭1ビットは繰り上がりに 使用できない
	leadingZeros--
	if uint(leadingZeros) < cnt {
		return p.DelayMax
	}
	delay := p.DelayMin << (cnt - 1)
	switch {
	case delay < p.DelayMin:
		return p.DelayMin
	case delay > p.DelayMax:
		return p.DelayMax
	default:
		return delay
	}
}

func (p Policy) NewBackoff() *Backoff {
	return &Backoff{p, 0}
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
	select {
	case <-ctx.Done():
		return false
	case <-time.After(b.NextTick()):
		b.cnt++
		if b.LimitExceeded() {
			return false
		}
		return true
	}
}
