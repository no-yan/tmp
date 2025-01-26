package backoff_test

import (
	"testing"
	"time"

	"github.com/no-yan/tmp/downloader/internal/backoff"
)

func TestPolicy_Next(t *testing.T) {
	tests := map[string]struct {
		cnt    uint
		policy *backoff.Policy
		want   time.Duration
	}{
		"first retry": {
			cnt:    0,
			policy: &backoff.Policy{10 * time.Millisecond, 100 * time.Millisecond, time.Second, 10},
			want:   0,
		},
		"second retry": {
			cnt:    1,
			policy: &backoff.Policy{10 * time.Millisecond, 100 * time.Millisecond, time.Second, 10},
			want:   10 * time.Millisecond,
		},
		"": {
			cnt:    3,
			policy: &backoff.Policy{10 * time.Millisecond, 100 * time.Millisecond, time.Second, 10},
			want:   40 * time.Millisecond,
		},
		"reach DelayMax": {
			cnt:    10,
			policy: &backoff.Policy{10 * time.Millisecond, 100 * time.Millisecond, time.Second, 10},
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
