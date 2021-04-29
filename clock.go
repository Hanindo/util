package util

import (
	"fmt"
	"time"
)

var (
	MockClockSpeed time.Duration = time.Millisecond / 2
)

type Clock interface {
	Now() time.Time
	NewTimer(d time.Duration) *Timer
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

//====================================================================

type MockClock struct {
	NowScript   []time.Duration
	TimerScript [][]time.Duration
	Ops         []TimeOp
	time        time.Time
	iNow        int
	iTimer      int
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
	rd := MockClockSpeed
	if m.iTimer <= len(m.TimerScript) {
		no = m.iTimer
		rd = m.TimerScript[m.iTimer-1][0]
	}

	m.Ops = append(m.Ops, TimeOp{"timer", d})
	ch := make(chan time.Time, 1)
	return &Timer{
		C:    ch,
		no:   no,
		mock: m,
		fake: time.AfterFunc(rd, func() {
			ch <- m.now()
		}),
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

	rd := MockClockSpeed
	if t.no > 0 {
		t.i++
		if t.i <= len(t.mock.TimerScript[t.no-1]) {
			rd = t.mock.TimerScript[t.no-1][t.i-1]
		}
	}
	t.mock.Ops = append(t.mock.Ops, TimeOp{
		fmt.Sprintf("timer-%d.reset", t.no), d,
	})
	return t.fake.Reset(rd)
}
