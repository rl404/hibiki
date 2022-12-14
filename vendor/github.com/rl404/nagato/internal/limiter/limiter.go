package limiter

import (
	"sync"
	"time"
)

// Limiter is interface for rate limiter.
type Limiter interface {
	Take()
}

type mutexLimiter struct {
	sync.Mutex
	last       time.Time
	sleepFor   time.Duration
	perRequest time.Duration
	maxSlack   time.Duration
	clock      clocker
}

type clocker interface {
	Now() time.Time
	Sleep(time.Duration)
}

// New returns a new atomic based limiter.
func New(rate int, interval time.Duration) Limiter {
	perRequest := interval / time.Duration(rate)
	return &mutexLimiter{
		perRequest: perRequest,
		maxSlack:   -1 * time.Duration(10) * perRequest,
		clock:      newClock(),
	}
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Second/rate.
func (t *mutexLimiter) Take() {
	t.Lock()
	defer t.Unlock()

	now := t.clock.Now()

	// If this is our first request, then we allow it.
	if t.last.IsZero() {
		t.last = now
		return
	}

	// sleepFor calculates how much time we should sleep based on
	// the perRequest budget and how long the last request took.
	// Since the request may take longer than the budget, this number
	// can get negative, and is summed across requests.
	t.sleepFor += t.perRequest - now.Sub(t.last)

	// We shouldn't allow sleepFor to get too negative, since it would mean that
	// a service that slowed down a lot for a short period of time would get
	// a much higher RPS following that.
	if t.sleepFor < t.maxSlack {
		t.sleepFor = t.maxSlack
	}

	// If sleepFor is positive, then we should sleep now.
	if t.sleepFor > 0 {
		t.clock.Sleep(t.sleepFor)
		t.last = now.Add(t.sleepFor)
		t.sleepFor = 0
	} else {
		t.last = now
	}
}
