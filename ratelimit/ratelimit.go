//go:build !solution
// +build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	buckets  chan struct{}
	interval time.Duration
	stop     chan struct{}
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	return &Limiter{
		buckets:  make(chan struct{}, maxCount),
		interval: interval,
		stop:     make(chan struct{}, 1),
	}
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-l.stop:
		return ErrStopped
	default:

	}
	select {
	case <-l.stop:
		return ErrStopped
	case <-ctx.Done():
		return ctx.Err()
	case l.buckets <- struct{}{}:
		go func() {
			select {
			case <-l.stop:
			case <-time.After(l.interval):
				<-l.buckets
			}
		}()
		return nil
	}
}

func (l *Limiter) Stop() {
	close(l.stop)
}
