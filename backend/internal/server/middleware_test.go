package server

import (
	"testing"
	"time"
)

func TestDelayFromEnv(t *testing.T) {
	cases := []struct {
		value             string
		wantOk            bool
		shortest, longest time.Duration
	}{
		{"", false, 0, 0},
		{"800ms", true, 800 * time.Millisecond, 800 * time.Millisecond},
		{" 1s ", true, time.Second, time.Second},
		{"300ms-2s", true, 300 * time.Millisecond, 2 * time.Second},
		{"2s-300ms", false, 0, 0},
		{"slow", false, 0, 0},
		{"1s-", false, 0, 0},
	}
	for _, c := range cases {
		t.Setenv("API_DELAY", c.value)
		shortest, longest, ok := delayFromEnv()
		if ok != c.wantOk || shortest != c.shortest || longest != c.longest {
			t.Errorf("API_DELAY=%q: got %v-%v ok %v, want %v-%v ok %v", c.value, shortest, longest, ok, c.shortest, c.longest, c.wantOk)
		}
	}
}
