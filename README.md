[![GoDoc](https://img.shields.io/badge/go.dev-reference-blue)](https://pkg.go.dev/github.com/chmike/ratelimiter)
[![Build](https://github.com/chmike/ratelimiter/actions/workflows/audit.yml/badge.svg)
[![Coverage](https://codecov.io/gh/chmike/ratelimiter/graph/badge.svg?token=06TJPZ1S5J)](https://codecov.io/gh/chmike/ratelimiter)
[![Go Report](https://goreportcard.com/badge/github.com/chmike/ratelimiter)](https://goreportcard.com/report/github.com/chmike/ratelimiter)
![Status](https://img.shields.io/badge/status-beta-orange)
![release](https://img.shields.io/github/release/chmike/ratelimiter.svg)

# Event rate limiter

This event rate limiter limits the rate of accepted events to *n* per 
*deltaT* where *deltaT* is a sliding window. It is intended to be used 
to precisely limit the rate of some operations like for instance loging 
attempts, mail sending, or request processing. 

This rate limiter allows to dynamically modify *n* to allow throttling 
the accepted event rate. This might be needed in some application for 
defensive measure. Setting *n* to 0 result in rejecting all events. 

It is implemented by using a circular FIFO queue containing the accepted
event time stamps. When a new event is received, outdated time stamps are
removed. If there is room in the queue the time stamp is appended and the 
event is accepted, Otherwise it is rejected.

This algorithm will thus accept bursts of request until *n* is reached, 
and require to wait the *deltaT* duration to accept new events. This may 
result in periodic bursts of accepted events when the input rate is higher 
than the accepted event rate. This is a feature of this rate limiting
method.

For a rate limiter working more like a random sampler, consider using 
the token bucket or leaky bucket rate limiting algorithm. Their 
limitation is that they may display temporary overshoots of the limited 
rate in some conditions. This is unacceptable when a strict rate limit is 
required. These rate limiter are perfect when the input rate is steadily 
higher than the output rate.  


## Example

The following is a simple usage example of the rate limiter.

```go
import "github.com/chmike/ratelimiter"

// taskRateLimiter limits tasks to at most 10 per seconds.
var taskRateLimiter = ratelimiter.New(10, time.Second)

// process a new task. Returns an error if rejected or an error occured.
func process(task *Task) error {
  if taskRateLimiter.Reject() {
    return errors.New("too many tasks, retry later")
  }
  // ... process task ...
  return nil
}
```

The following is a request rate limiting web server middleware. 

```go
import (
  "github.com/chmike/ratelimiter"
  "net/http"
)

// some example server action
func hello(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "hello !")
}

// rateLimiter middleware
func rateLimiter(n int, deltaT time.Duration, f http.HandlerFunc) http.HandlerFunc {
  requestRateLimiter = ratelimiter.New(n, deltaT)
  return func(w http.ResponseWriter, r *http.Request) {
    if requestRateLimiter.Accept() {
      f(w, r)
    } else {
      http.Error(w, "too many requests, retry later", http.StatusTooManyRequests)
    }
  }
}

func main() {
  http.HandleFunc("/", rateLimiter(2, time.Second, hello)) // accept calling hello at most twice per second
  http.ListenAndServe(":8080", nil)
}
```


