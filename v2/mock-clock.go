package util

import (
	"fmt"
	"time"
)

// Default speed for fake ticker/timer.
const (
	MOCK_TICKER_SPEED time.Duration = time.Millisecond
	MOCK_TIMER_SPEED  time.Duration = time.Millisecond
)

// A TimeOp represents time operation.
type TimeOp struct {
	Op       string
	Duration time.Duration
}

// A MockClock represents simple mock clock.
//
// The clock operation can be directed by various *Script fields. All scripts
// are optional, the clock can run fine without any script. All clock
// operation will be recorded on Ops field.
//
// This clock runs on fake timer/ticker speed so unit testing don't have to
// wait that long. The default maximum wait can be adjusted on
// TickerSpeed/TimerSpeed field, or it can be scripted in
// TickerScript/TimerScript field.
//
// For simplicity this clock don't send the correct mocked time value on
// Ticker/Timer channel, so DON'T USE THIS CLOCK if your program depend on the
// value sent from those channels. Patches are welcome to fix this limitation.
type MockClock struct {
	NowScript    []time.Duration
	TickerScript [][]time.Duration
	TimerScript  [][]time.Duration
	Ops          []TimeOp

	TickerSpeed time.Duration
	TimerSpeed  time.Duration

	time    time.Time
	iNow    int
	iTicker int
	iTimer  int
}

// NewMockClock creates a new MockClock with the time initialized to t.
func NewMockClock(t time.Time) *MockClock {
	return &MockClock{
		TickerSpeed: MOCK_TICKER_SPEED,
		TimerSpeed:  MOCK_TIMER_SPEED,
		time:        t,
	}
}

// Now returns the current mocked time.
// Please note this always advance the time.
func (m *MockClock) Now() time.Time {
	m.iNow++
	var d time.Duration = 1
	if m.iNow <= len(m.NowScript) && m.NowScript[m.iNow-1] > 0 {
		d = m.NowScript[m.iNow-1]
	}
	m.time = m.time.Add(d)

	m.Ops = append(m.Ops, TimeOp{"now", d})
	return m.time
}

// NewTicker returns a new time.Ticker compatible ticker.
func (m *MockClock) NewTicker(d time.Duration) *Ticker {
	fd := d
	if fd > m.TickerSpeed {
		fd = m.TickerSpeed
	}
	m.iTicker++
	if m.iTicker <= len(m.TickerScript) &&
		len(m.TickerScript[m.iTicker-1]) > 0 {
		fd = m.TickerScript[m.iTicker-1][0]
	}

	m.Ops = append(m.Ops, TimeOp{"ticker", d})

	fake := time.NewTicker(fd)
	return &Ticker{
		Tickerable: &mockTicker{
			no:   m.iTicker,
			mock: m,
			fake: fake,
		},
		C: fake.C,
	}
}

// NewTimer returns a new time.Timer compatible timer.
func (m *MockClock) NewTimer(d time.Duration) *Timer {
	fd := d
	if fd > m.TimerSpeed {
		fd = m.TimerSpeed
	}
	m.iTimer++
	if m.iTimer <= len(m.TimerScript) && len(m.TimerScript[m.iTimer-1]) > 0 {
		fd = m.TimerScript[m.iTimer-1][0]
	}

	m.Ops = append(m.Ops, TimeOp{"timer", d})

	fake := time.NewTimer(fd)
	return &Timer{
		Timerable: &mockTimer{
			no:   m.iTimer,
			mock: m,
			fake: fake,
		},
		C: fake.C,
	}
}

//============================================================================

type mockTicker struct {
	i    int
	no   int
	mock *MockClock
	fake *time.Ticker
}

func (t *mockTicker) Stop() {
	t.mock.Ops = append(t.mock.Ops, TimeOp{
		fmt.Sprintf("ticker-%d.stop", t.no),
		0,
	})
	t.fake.Stop()
	return
}

func (t *mockTicker) Reset(d time.Duration) {
	fd := d
	if d > t.mock.TimerSpeed {
		fd = t.mock.TickerSpeed
	}
	if t.no <= len(t.mock.TickerScript) {
		t.i++
		if t.i < len(t.mock.TickerScript[t.no-1]) {
			fd = t.mock.TickerScript[t.no-1][t.i]
		}
	}

	t.mock.Ops = append(t.mock.Ops, TimeOp{
		fmt.Sprintf("ticker-%d.reset", t.no),
		d,
	})
	t.fake.Reset(fd)
	return
}

//============================================================================

type mockTimer struct {
	i    int
	no   int
	mock *MockClock
	fake *time.Timer
}

func (t *mockTimer) Stop() bool {
	t.mock.Ops = append(t.mock.Ops, TimeOp{
		fmt.Sprintf("timer-%d.stop", t.no),
		0,
	})
	return t.fake.Stop()
}

func (t *mockTimer) Reset(d time.Duration) bool {
	fd := d
	if d > t.mock.TimerSpeed {
		fd = t.mock.TimerSpeed
	}
	if t.no <= len(t.mock.TimerScript) {
		t.i++
		if t.i < len(t.mock.TimerScript[t.no-1]) {
			fd = t.mock.TimerScript[t.no-1][t.i]
		}
	}

	t.mock.Ops = append(t.mock.Ops, TimeOp{
		fmt.Sprintf("timer-%d.reset", t.no),
		d,
	})
	return t.fake.Reset(fd)
}
