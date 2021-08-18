package util

import (
	"time"
)

// Clock is interface for time operation, so time sensitive application
// can be easily unit tested by mocking time.
type Clock interface {
	Now() time.Time
	NewTicker(d time.Duration) *Ticker
	NewTimer(d time.Duration) *Timer
}

// Tickerable is interface for time.Ticker.
type Tickerable interface {
	Reset(d time.Duration)
	Stop()
}

// Timerable is interface for time.Timer.
type Timerable interface {
	Reset(d time.Duration) bool
	Stop() bool
}

// Ticker is time.Ticker drop-in replacement.
type Ticker struct {
	// The real ticker implementation
	Tickerable
	// The channel on which the ticks are delivered.
	C <-chan time.Time
}

// Timer is time.Timer drop-in replacement.
type Timer struct {
	// The real timer implementation
	Timerable
	C <-chan time.Time
}

//====================================================================

type clock struct{}

// NewClock returns a new real-time Clock.
func NewClock() Clock {
	return &clock{}
}

func (c *clock) Now() time.Time {
	return time.Now()
}

func (c *clock) NewTicker(d time.Duration) *Ticker {
	t := time.NewTicker(d)
	return &Ticker{
		Tickerable: t,
		C:          t.C,
	}
}

func (c *clock) NewTimer(d time.Duration) *Timer {
	t := time.NewTimer(d)
	return &Timer{
		Timerable: t,
		C:         t.C,
	}
}
