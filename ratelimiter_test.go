package ratelimiter

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	deltaT := 500 * time.Millisecond
	n := 10
	r := New(n, deltaT)

	for r.Accept() {
	}
	if r.Len() != n {
		t.Errorf("expect len %d, got %d", n, r.Len())
	}
	if r.Accept() {
		t.Error("expect Accept to fail")
	}
	time.Sleep(deltaT)
	if r.Reject() {
		t.Error("expect Accept to succeed")
	}
	if r.Len() != 1 {
		t.Errorf("expect len 1, got %d", r.Len())
	}
	r.Reset()
	if r.Len() != 0 {
		t.Errorf("expect len 0, got %d", r.Len())
	}

	// test modifying N
	if r.N() != n {
		t.Errorf("expect N is %d, got %d", n, r.N())
	}
	n = 5
	r.SetN(n)
	for r.Accept() {
	}
	if r.Len() != n {
		t.Errorf("expect len 5, got %d", r.Len())
	}
	if r.Accept() {
		t.Error("expect Accept to fail")
	}
	time.Sleep(deltaT)
	r.Purge()
	if r.Len() != 0 {
		t.Errorf("expect len 0, got %d", r.Len())
	}
	r.ResetN()
	n = 10
	if r.N() != n {
		t.Errorf("expect N is %d, got %d", n, r.N())
	}
	r.SetN(n + 10)
	if r.N() != n {
		t.Errorf("expect N is %d, got %d", n, r.N())
	}
	n = 0
	r.SetN(-10)
	if r.N() != n {
		t.Errorf("expect N is %d, got %d", n, r.N())
	}
}
