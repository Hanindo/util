package util

import (
	"fmt"
	"time"
)

var (
	MockTimerSpeed  time.Duration = time.Millisecond / 2
	MockTickerSpeed time.Duration = time.Millisecond
)

type Clock interface {
	Now() time.Time
	NewTimer(d time.Duration) *Timer
	NewTicker(d time.Duration) *Ticker
}

type TimeOp struct {
	Op       string
	Duration time.Duration
}

//====================================================================

type clock struct{}

func NewClock() Clock {
	return &clock{}
}

func (c *clock) Now() time.Time {
	return time.Now()
}

func (c *clock) NewTimer(d time.Duration) *Timer {
	t := time.NewTimer(d)
	return &Timer{
		C:     t.C,
		timer: t,
	}
}

func (c *clock) NewTicker(d time.Duration) *Ticker {
	t := time.NewTicker(d)
	return &Ticker{
		C:      t.C,
		ticker: t,
	}
}

//====================================================================

type MockClock struct {
	NowScript    []time.Duration
	TimerScript  [][]time.Duration
	TickerScript [][]time.Duration
	Ops          []TimeOp
	time         time.Time
	iNow         int
	iTimer       int
	iTicker      int
}

func NewMockClock(t time.Time) *MockClock {
	return &MockClock{time: t}
}

func (m *MockClock) Now() time.Time {
	/*
	   _, file, line, _ := runtime.Caller(1)
	   debugLog("now:%s %d", file, line)
	*/
	m.iNow++
	var d time.Duration = 1 // time.Now() always advance time
	if m.iNow <= len(m.NowScript) && m.NowScript[m.iNow-1] > 0 {
		d = m.NowScript[m.iNow-1]
	}
	m.time = m.time.Add(d)

	m.Ops = append(m.Ops, TimeOp{"now", d})
	return m.time
}

func (m *MockClock) now() time.Time {
	return m.time
}

func (m *MockClock) NewTimer(d time.Duration) *Timer {
	m.iTimer++
	var no int
	fd := MockTimerSpeed
	if m.iTimer <= len(m.TimerScript) {
		no = m.iTimer
		fd = m.TimerScript[m.iTimer-1][0]
	}

	m.Ops = append(m.Ops, TimeOp{"timer", d})
	ch := make(chan time.Time, 1)
	return &Timer{
		C:    ch,
		no:   no,
		mock: m,
		fake: time.AfterFunc(fd, func() {
			ch <- m.now()
		}),
	}
}

func (m *MockClock) NewTicker(d time.Duration) *Ticker {
	m.iTicker++
	var no int
	fd := MockTickerSpeed
	if m.iTicker <= len(m.TickerScript) {
		no = m.iTicker
		fd = m.TickerScript[m.iTicker-1][0]
	}

	m.Ops = append(m.Ops, TimeOp{"ticker", d})
	ch := make(chan time.Time, 1)
	fake := time.NewTicker(fd)
	go func() {
		for _ = range fake.C {
			ch <- m.now()
		}
	}()

	return &Ticker{
		C:    ch,
		no:   no,
		mock: m,
		fake: fake,
	}
}

//====================================================================

type Timer struct {
	C     <-chan time.Time
	timer *time.Timer
	no    int
	i     int
	mock  *MockClock
	fake  *time.Timer
}

func (t *Timer) Stop() bool {
	if t.timer != nil {
		return t.timer.Stop()
	}

	t.mock.Ops = append(t.mock.Ops, TimeOp{
		fmt.Sprintf("timer-%d.stop", t.no),
		0,
	})
	return t.fake.Stop()
}

func (t *Timer) Reset(d time.Duration) bool {
	if t.timer != nil {
		return t.timer.Reset(d)
	}

	fd := MockTimerSpeed
	if t.no > 0 {
		t.i++
		if t.i <= len(t.mock.TimerScript[t.no-1]) {
			fd = t.mock.TimerScript[t.no-1][t.i-1]
		}
	}
	t.mock.Ops = append(t.mock.Ops, TimeOp{
		fmt.Sprintf("timer-%d.reset", t.no), d,
	})
	return t.fake.Reset(fd)
}

//====================================================================

type Ticker struct {
	C      <-chan time.Time
	ticker *time.Ticker
	no     int
	i      int
	mock   *MockClock
	fake   *time.Ticker
}

func (t *Ticker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		return
	}

	t.mock.Ops = append(t.mock.Ops, TimeOp{
		fmt.Sprintf("ticker-%d.stop", t.no),
		0,
	})
	t.fake.Stop()
	return
}

func (t *Ticker) Reset(d time.Duration) {
	if t.ticker != nil {
		t.ticker.Reset(d)
		return
	}

	fd := MockTickerSpeed
	if t.no > 0 {
		t.i++
		if t.i <= len(t.mock.TickerScript[t.no-1]) {
			fd = t.mock.TickerScript[t.no-1][t.i-1]
		}
	}
	t.mock.Ops = append(t.mock.Ops, TimeOp{
		fmt.Sprintf("ticker-%d.reset", t.no), d,
	})
	t.fake.Reset(fd)
	return
}
