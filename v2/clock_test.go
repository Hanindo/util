package util_test

import (
	. "github.com/hanindo/util/v2"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clock", func() {
	var c Clock
	BeforeEach(func() {
		c = NewClock()
	})

	Describe("Now", func() {
		It("return increasing time.Time", func() {
			t1 := c.Now()
			time.Sleep(1)
			t2 := c.Now()
			Expect(t2.After(t1)).To(BeTrue())
		})
	})

	Describe("NewTicker", func() {
		It("tick in regular order", func() {
			d1 := 100 * time.Millisecond
			d2 := 50 * time.Millisecond

			t := c.NewTicker(d1)
			ct1 := <-t.C
			t1 := time.Now()
			ct2 := <-t.C
			t2 := time.Now()
			t.Reset(d2)
			ct3 := <-t.C
			t3 := time.Now()
			t.Stop()
			Consistently(t.C).ShouldNot(Receive())

			Expect(ct1).To(BeTemporally("~", t1), "ct1")
			Expect(ct2).To(BeTemporally("~", t2), "ct2")
			Expect(ct3).To(BeTemporally("~", t3), "ct3")

			Expect(ct1.Before(ct2)).To(BeTrue(), "ct1 < ct2")
			Expect(ct2).To(BeTemporally("~", ct1.Add(d1), 10*time.Millisecond),
				"ct2 - ct1")
			Expect(ct2.Before(ct3)).To(BeTrue(), "ct2 < ct3")
			Expect(ct3).To(BeTemporally("~", ct2.Add(d2), 10*time.Millisecond),
				"ct3 - ct2")
		})
	})

	Describe("NewTimer", func() {
		It("fired by the timer", func() {
			t := c.NewTimer(200 * time.Millisecond)
			Consistently(t.C, 80*time.Millisecond, 10*time.Millisecond).
				ShouldNot(Receive())

			d := 20 * time.Millisecond
			Expect(t.Reset(d)).To(BeTrue(), "1st reset")
			t1 := time.Now()
			ct := <-t.C
			t2 := time.Now()
			Expect(ct).To(BeTemporally("~", t2), "ct")
			Expect(t2).To(BeTemporally("~", t1.Add(d)), "t2 - t1")

			Expect(t.Reset(50 * time.Millisecond)).To(BeFalse())
			Expect(t.Stop()).To(BeTrue(), "stop")
			Consistently(t.C).ShouldNot(Receive())
		})
	})
})
