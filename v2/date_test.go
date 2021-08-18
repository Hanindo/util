package util_test

import (
	"encoding/json"
	. "github.com/hanindo/util/v2"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Date", func() {
	var d Date

	DescribeTable("Valid",
		func(bin, str, st, sx string) {
			t, _ := time.Parse(time.RFC3339, st)
			x, _ := time.Parse(time.RFC3339, sx)
			d = MakeDate(t)
			Expect(d.Time()).To(Equal(x), "Time")
			Expect(d.String()).To(Equal(str), "String")

			b, err := json.Marshal(d)
			Expect(err).To(Succeed(), "json encode err")
			Expect(string(b)).To(Equal(`"`+str+`"`), "json encode")

			var nd Date
			Expect(json.Unmarshal(b, &nd)).To(Succeed(), "json decode err")
			Expect(nd.Equal(d)).To(BeTrue(), "json decode")

			b, err = d.MarshalBinary()
			Expect(err).To(Succeed(), "binary encode err")
			Expect(FancyHex(b)).To(Equal(bin), "binary encode")

			nd = Date{}
			Expect(nd.UnmarshalBinary(b)).To(Succeed(), "binary decode err")
			Expect(nd.String()).To(Equal(d.String()), "binary decode")
		},
		Entry("UTC", "b. .. .. 4.",
			"0000-01-01 +00:00",
			"0000-01-01T12:13:14Z",
			"0000-01-01T00:00:00Z"),
		Entry("UTC+1", "b7 d. .. c4",
			"2000-01-02 +01:00",
			"2000-01-02T12:13:14+01:00",
			"2000-01-02T00:00:00+01:00"),
		Entry("UTC-1", "b7 dc .. bc",
			"2012-01-02 -01:00",
			"2012-01-02T12:13:14-01:00",
			"2012-01-02T00:00:00-01:00"),
		Entry("Nepal", "b7 e5 1. 57",
			"2021-02-01 +05:45",
			"2021-02-01T12:13:14+05:45",
			"2021-02-01T00:00:00+05:45"),
		Entry("UTC-12", "b7 e5 77 9.",
			"2021-08-16 -12:00",
			"2021-08-16T12:13:14-12:00",
			"2021-08-16T00:00:00-12:00"),
		Entry("UTC+14", "bf fe bf 78",
			"4094-12-31 +14:00",
			"4094-12-31T12:13:14+14:00",
			"4094-12-31T00:00:00+14:00"),
	)

	Describe("Local Binary", func() {
		It("should safely encoded and decoded", func() {
			t := time.Date(2021, time.February, 1, 23, 24, 25, 0, time.Local)
			d := MakeDate(t)

			b, err := d.MarshalBinary()
			Expect(err).To(Succeed(), "encode err")
			Expect(FancyHex(b[:3])).To(Equal("b7 e5 1."), "encode")
			Expect(b[3]&0x80).To(Equal(byte(0)), "encode last byte")

			var nd Date
			Expect(nd.UnmarshalBinary(b)).To(Succeed(), "decode err")
			Expect(nd.String()).To(Equal(d.String()), "decode")
			Expect(nd.Time().Location()).To(Equal(time.Local), "tz")
		})
	})

	Describe("UnmarshalText", func() {
		Context("invalid", func() {
			It("should error", func() {
				b := []byte(`"2012-01-02"`)
				Expect(json.Unmarshal(b, &d)).To(HaveOccurred())
			})
		})
	})

	Describe("MarshalBinary", func() {
		Context("minus year", func() {
			It("should error", func() {
				t := time.Date(-1, time.December, 31, 23, 59, 59, 0, time.UTC)
				b, err := MakeDate(t).MarshalBinary()
				Expect(b).To(BeEmpty())
				Expect(err).To(MatchError("year outside range 0-4094: -1"))
			})
		})

		Context("year > 4094", func() {
			It("should error", func() {
				t := time.Date(4095, time.January, 1, 0, 0, 0, 0, time.UTC)
				b, err := MakeDate(t).MarshalBinary()
				Expect(b).To(BeEmpty())
				Expect(err).To(MatchError("year outside range 0-4094: 4095"))
			})
		})
	})

	Describe("UnmarshalBinary", func() {
		DescribeTable("invalid",
			func(err string, b []byte) {
				var d Date
				Expect(d.UnmarshalBinary(b)).To(MatchError(err))
			},
			Entry("long", "invalid length: 5",
				[]byte{0xB7, 0xE5, 0x10, 0x57, 0x12}),
			Entry("short", "invalid length: 3", []byte{0xB7, 0xE5, 0x10}),
			Entry("type", "invalid type bits: 1010",
				[]byte{0xA7, 0xE5, 0x10, 0x57}),
			Entry("year", "invalid year: 4095",
				[]byte{0xBF, 0xFF, 0x10, 0x57}),
			Entry("month", "invalid month: 13",
				[]byte{0xB7, 0xE5, 0xC0, 0x57}),
			Entry("day", "invalid day: 32",
				[]byte{0xB7, 0xE5, 0x1F, 0xD7}),
			Entry("timezone", "invalid timezone: 15:30",
				[]byte{0xB7, 0xE5, 0x10, 0x7E}),
			Entry("date", "invalid date: 2021-02-29",
				[]byte{0xB7, 0xE5, 0x1E, 0x57}),
		)
	})
})
