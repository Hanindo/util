package util_test

import (
	"encoding/json"
	"fmt"
	. "github.com/hanindo/util/v2"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Date", func() {
	var d Date

	DescribeTable("Valid",
		func(str string, xbin []byte, st, sx string) {
			t, _ := time.Parse(time.RFC3339, st)
			x, _ := time.Parse(time.RFC3339, sx)
			d = MakeDate(t)
			Expect(d.Time()).To(Equal(x), "Time")
			Expect(d.String()).To(Equal(str), "String")

			b, err := json.Marshal(d)
			Expect(err).To(Succeed(), "json encode err")
			Expect(string(b)).To(Equal(`"`+str+`"`), "json encode")

			var dd Date
			Expect(json.Unmarshal(b, &dd)).To(Succeed(), "json decode err")
			Expect(dd.Equal(d)).To(BeTrue(), "json decode")

			b, err = d.MarshalBinary()
			Expect(err).To(Succeed(), "binary encode err")
			Expect(fmt.Sprintf("% X", b)).
				To(Equal(fmt.Sprintf("% X", xbin)), "binary encode")

			dd = Date{}
			Expect(dd.UnmarshalBinary(b)).To(Succeed(), "binary decode err")
			Expect(dd.String()).To(Equal(d.String()), "binary decode")
		},
		Entry("UTC", "0001-01-01 +00:00", []byte{0xB0, 0x01, 0x00, 0x40},
			"0001-01-01T12:13:14Z",
			"0001-01-01T00:00:00Z"),
		Entry("UTC+1", "2000-01-02 +01:00", []byte{0xB7, 0xD0, 0x00, 0xC4},
			"2000-01-02T12:13:14+01:00",
			"2000-01-02T00:00:00+01:00"),
		Entry("UTC-1", "2012-01-02 -01:00", []byte{0xB7, 0xDC, 0x00, 0xBC},
			"2012-01-02T12:13:14-01:00",
			"2012-01-02T00:00:00-01:00"),
		Entry("Nepal", "2021-02-01 +05:45", []byte{0xB7, 0xE5, 0x10, 0x57},
			"2021-02-01T12:13:14+05:45",
			"2021-02-01T00:00:00+05:45"),
		Entry("UTC-12", "2021-08-16 -12:00", []byte{0xB7, 0xE5, 0x77, 0x90},
			"2021-08-16T12:13:14-12:00",
			"2021-08-16T00:00:00-12:00"),
		Entry("UTC+14", "4094-12-31 +14:00", []byte{0xBF, 0xFE, 0xBF, 0x78},
			"4094-12-31T12:13:14+14:00",
			"4094-12-31T00:00:00+14:00"),
	)

	Describe("Invalid Text", func() {
		It("error on decode", func() {
			b := []byte(`"2012-01-02"`)
			Expect(json.Unmarshal(b, &d)).To(HaveOccurred())
		})
	})
})
