package util_test

import (
	"time"

	. "github.com/hanindo/util/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mock Clock", func() {
	var c *MockClock
	t := time.Date(2021, time.February, 1, 23, 24, 25, 0, time.Local)
	BeforeEach(func() {
		c = NewMockClock(t)
	})

	Describe("Now", func() {
		Context("unscripted", func() {
			It("append now with 1 duration to ops", func() {
				Expect(c.Now()).To(Equal(t.Add(1)))
				Expect(c.Ops).To(Equal([]TimeOp{{"now", 1}}))
			})
		})

		Context("scripted", func() {
			It("append now with the durations to ops", func() {
				c.NowScript = []time.Duration{100, 0}
				Expect(c.Now()).To(Equal(t.Add(100)), "#1")
				Expect(c.Now()).To(Equal(t.Add(101)), "#2")
				Expect(c.Now()).To(Equal(t.Add(102)), "#3")
				Expect(c.Ops).To(Equal([]TimeOp{
					{"now", 100},
					{"now", 1},
					{"now", 1},
				}))
			})
		})
	})

	Describe("Ticker", func() {
		Context("unscripted", func() {
			It("append ticker to ops", func() {
				t1 := c.NewTicker(time.Second)
				t2 := c.NewTicker(2 * time.Second)
				t2.Reset(time.Second / 2)
				t1.Stop()

				Eventually(t2.C).Should(Receive())
				Expect(c.Ops).To(Equal([]TimeOp{
					{"ticker", time.Second},
					{"ticker", 2 * time.Second},
					{"ticker-2.reset", time.Second / 2},
					{"ticker-1.stop", 0},
				}))
			})
		})

		Context("scripted", func() {
			It("append ticker to ops", func() {
				c.TickerScript = [][]time.Duration{
					nil,
					{MOCK_TICKER_SPEED / 2, MOCK_TICKER_SPEED * 2},
				}
				t1 := c.NewTicker(time.Second)
				t2 := c.NewTicker(2 * time.Second)
				t2.Reset(time.Second / 2)
				t2.Reset(1)
				t1.Stop()

				Eventually(t2.C).Should(Receive())
				Expect(c.Ops).To(Equal([]TimeOp{
					{"ticker", time.Second},
					{"ticker", 2 * time.Second},
					{"ticker-2.reset", time.Second / 2},
					{"ticker-2.reset", 1},
					{"ticker-1.stop", 0},
				}))
			})
		})
	})

	Describe("Timer", func() {
		Context("unscripted", func() {
			It("append ticker to ops", func() {
				t1 := c.NewTimer(time.Second)
				t2 := c.NewTimer(2 * time.Second)
				Expect(t2.Reset(time.Second / 2)).To(BeTrue())
				Expect(t1.Stop()).To(BeTrue())

				Eventually(t2.C).Should(Receive())
				Expect(c.Ops).To(Equal([]TimeOp{
					{"timer", time.Second},
					{"timer", 2 * time.Second},
					{"timer-2.reset", time.Second / 2},
					{"timer-1.stop", 0},
				}))
			})
		})

		Context("scripted", func() {
			It("append ticker to ops", func() {
				c.TimerScript = [][]time.Duration{
					nil,
					{MOCK_TIMER_SPEED / 2, MOCK_TIMER_SPEED * 2},
				}
				t1 := c.NewTimer(time.Second)
				t2 := c.NewTimer(2 * time.Second)
				Expect(t2.Reset(time.Second / 2)).To(BeTrue())
				Expect(t2.Reset(1)).To(BeTrue())
				Expect(t1.Stop()).To(BeTrue())

				Eventually(t2.C).Should(Receive())
				Expect(c.Ops).To(Equal([]TimeOp{
					{"timer", time.Second},
					{"timer", 2 * time.Second},
					{"timer-2.reset", time.Second / 2},
					{"timer-2.reset", 1},
					{"timer-1.stop", 0},
				}))
			})
		})
	})
})
