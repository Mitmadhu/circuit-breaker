package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

const (
	OPEN      = "OPEN"
	CLOSE     = "CLOSE"
	HALF_OPEN = "HALFOPEN"
)

type circuitBreaker struct {
	lock             sync.Mutex
	lastFailed       time.Time
	successCount     int
	successThreshold int
	errorThreshold   int
	duration         time.Duration
	state            string
	errorCount       int
	errTimestamp     []time.Time
	lastStateChange  time.Time
}

func NewCB(duration time.Duration, errThresold, successThresold int) *circuitBreaker {
	return &circuitBreaker{
		lock:             sync.Mutex{},
		successCount:     0,
		successThreshold: successThresold,
		errorThreshold:   errThresold,
		duration:         duration,
		state:            CLOSE,
		errorCount:       0,
	}
}

func (cb *circuitBreaker) Run(apiCall func() error) error {
	cb.lock.Lock()
	defer cb.lock.Unlock()

	if cb.state == OPEN {
		if time.Since(cb.lastFailed) >= cb.duration {
			cb.state = HALF_OPEN
		} else {
			return errors.New("circuit is open")
		}
	}
	// close or half open
	err := apiCall()
	if err != nil {
		cb.successCount = 0
		cb.lastFailed = time.Now()
		cb.errTimestamp = append(cb.errTimestamp, time.Now())
		if cb.state == HALF_OPEN {
			cb.state = OPEN
			return err
		}

		// circuit was close
		cb.removeOldErrors()
		if len(cb.errTimestamp) >= cb.errorThreshold {
			cb.state = OPEN
			cb.lastStateChange = time.Now()
			println("[info] circuit is open now.")
		}

	}

	// error is not nil and CB is half open
	if cb.state == HALF_OPEN {
		cb.successCount += 1
		if cb.successCount == cb.successThreshold {
			cb.state = CLOSE
			println("[info] circuit is closed")
		}
	}

	return err
}

func (cb *circuitBreaker) removeOldErrors() {
	l := 0
	r := len(cb.errTimestamp)
	i := l
	now := time.Now()
	for l <= r {
		mid := l + (r-l)/2
		if now.Sub(cb.errTimestamp[mid]) > cb.duration {
			i = mid
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	cb.errTimestamp = cb.errTimestamp[i:]
}
