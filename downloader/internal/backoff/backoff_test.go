package backoff

import (
	"testing"
	"time"
)

func TestPolicy_Next(t *testing.T) {
	tests := map[string]struct {
		cnt    uint
		policy *Policy
		want   time.Duration
	}{
		"first retry": {
			cnt:    0,
			policy: &Policy{10 * time.Millisecond, 100 * time.Millisecond, 10},
			want:   0,
		},
		"second retry": {
			cnt:    1,
			policy: &Policy{10 * time.Millisecond, 100 * time.Millisecond, 10},
			want:   10 * time.Millisecond,
		},
		"third retry": {
			cnt:    3,
			policy: &Policy{10 * time.Millisecond, 100 * time.Millisecond, 10},
			want:   40 * time.Millisecond,
		},
		"reach DelayMax": {
			cnt:    10,
			policy: &Policy{10 * time.Millisecond, 100 * time.Millisecond, 10},
			want:   100 * time.Millisecond,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.policy.Next(tt.cnt); got != tt.want {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackoff_LimitExceeded(t *testing.T) {
	defaultPolicy := Policy{1, 100, 100}
	tests := []struct {
		name string
		p    Policy
		cnt  uint
		want bool
	}{
		{"0", defaultPolicy, 0, false},
		{"1", defaultPolicy, 1, false},
		{"100", defaultPolicy, 100, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backoff{
				p:   tt.p,
				cnt: tt.cnt,
			}
			if got := b.LimitExceeded(); got != tt.want {
				t.Errorf("LimitExceeded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackoff_NextTick(t *testing.T) {
	p := Policy{1 * time.Second, 100 * time.Second, 200}

	tests := []struct {
		name string
		p    Policy
		cnt  uint
		want time.Duration
	}{
		{"0", p, 0, 0},
		{"1", p, 1, 1 * time.Second},
		{"2", p, 2, 2 * time.Second},
		{"5", p, 5, 16 * time.Second},
		{"max", p, 400, 100 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backoff{
				p:   p,
				cnt: tt.cnt,
			}
			if got := b.NextTick(); got != tt.want {
				t.Errorf("Backoff.NextTick() = %v, want %v", got, tt.want)
			}
		})
	}
}
