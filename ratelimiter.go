package ratelimiter

import (
	"sync"
	"time"
)

// A RateLimiter accepts at most n events per deltatT. It uses the
// sliding window algorithm applied on a circular FIFO queue of event
// time stamps.
type RateLimiter struct {
	mtx    sync.Mutex
	log    []int64 // UnixNano time stamp circular FIFO queue
	end    int     // end of queue holding the oldest time stamp
	len    int     // number of time stamps in the queue
	cap    int     // number of events accepted per deltaT
	deltaT int64   // time frame in nano seconds
}

// New instantiates a new rate limiter accepting at most maxN events per deltaT.
// N may be later changed to a value in the rage 0 to maxN to throttle the rate.
func New(maxN int, deltaT time.Duration) *RateLimiter {
	return &RateLimiter{
		log:    make([]int64, maxN),
		cap:    maxN,
		deltaT: deltaT.Nanoseconds(),
	}
}

// Len returns the number of event time stamps in the queue.
// Call Purge() to remove outdated events.
func (r *RateLimiter) Len() int {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.len
}

// Reset empties the event time stamp queue.
func (r *RateLimiter) Reset() {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.len = 0
}

// Accept return true when the event that occured now is accepted.
// It purges the queue from outdated event time stamps.
func (r *RateLimiter) Accept() bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.accept(time.Now().UnixNano())
}

// Reject return true when the event that occured now is not accepted.
// It purges the queue from outdated event time stamps.
func (r *RateLimiter) Reject() bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return !r.accept(time.Now().UnixNano())
}

func (r *RateLimiter) accept(t int64) bool {
	r.purge(t)
	if r.len >= r.cap {
		return false
	}
	r.log[(r.end+r.len)%len(r.log)] = t
	r.len++
	return true
}

// Purge removes outdated time stamps relative to the current time.
func (r *RateLimiter) Purge() {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.purge(time.Now().UnixNano())
}

func (r *RateLimiter) purge(t int64) {
	oldest := t - r.deltaT
	for r.len > 0 && r.log[r.end] <= oldest {
		r.len--
		r.end = (r.end + 1) % len(r.log)
	}
}

// N return the maxmimum number of accepted events per deltaT.
func (r *RateLimiter) N() int {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.cap
}

// SetN sets the maximum number of accepted events per deltaT.
// It doesn't modify the actual number of events in the queue.
func (r *RateLimiter) SetN(cap int) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	switch {
	case cap < 0:
		r.cap = 0
	case cap > len(r.log):
		r.cap = len(r.log)
	default:
		r.cap = cap
	}
}

// ResetN resets N to the maximum value.
func (r *RateLimiter) ResetN() {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.cap = len(r.log)
}
