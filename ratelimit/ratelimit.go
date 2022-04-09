//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	buckets  chan time.Time
	interval time.Duration
	stop     chan struct{}
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	l := &Limiter{
		buckets:  make(chan time.Time, maxCount),
		interval: interval,
		stop:     make(chan struct{}, 1),
	}

	for i := 0; i < maxCount; i++ {
		l.buckets <- time.Now()
	}

	return l
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-l.stop:
		return ErrStopped
	default:

	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:

	}

	select {
	case <-l.stop:
		return ErrStopped
	case <-ctx.Done():
		return ctx.Err()
	case tm := <-l.buckets:
		now := time.Now()
		if d := tm.Sub(now); d > 0 {
			select {
			case <-time.After(d - 200*time.Millisecond):
			case <-l.stop:
				return ErrStopped
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		select {
		case l.buckets <- time.Now().Add(l.interval):
		case <-l.stop:
			return ErrStopped
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	}
}

func (l *Limiter) Stop() {
	close(l.stop)
}
