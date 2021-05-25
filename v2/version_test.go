package util_test

import (
	"encoding/json"
	"fmt"
	. "github.com/hanindo/util/v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	var v Version
	Describe("Make", func() {
		It("should create version", func() {
			v = MakeVersion(1, 2, 3)
			Expect(v).To(Equal(Version{Major: 1, Minor: 2, Patch: 3}))
		})
	})

	Describe("Before", func() {
		DescribeTable("True",
			func(major, minor, patch int) {
				v = MakeVersion(1, 2, 3)
				b := MakeVersion(uint8(major), uint8(minor), uint8(patch))
				Expect(v.Before(b)).To(BeTrue())
			},
			Entry("< major", 2, 1, 2),
			Entry("< minor", 1, 3, 2),
			Entry("< patch", 1, 2, 4),
		)

		DescribeTable("False",
			func(major, minor, patch int) {
				v = MakeVersion(1, 2, 3)
				b := MakeVersion(uint8(major), uint8(minor), uint8(patch))
				Expect(v.Before(b)).To(BeFalse())
			},
			Entry("equal", 1, 2, 3),
			Entry("> major", 0, 3, 4),
			Entry("> minor", 1, 1, 4),
			Entry("> patch", 1, 2, 2),
		)
	})

	Describe("String", func() {
		It("should create version", func() {
			v = MakeVersion(1, 2, 3)
			Expect(v.String()).To(Equal("1.2.3"))
		})
	})

	Describe("JSON encode", func() {
		It("should be encoded", func() {
			v = MakeVersion(1, 2, 3)
			Expect(json.Marshal(v)).To(Equal([]byte(`"1.2.3"`)))
		})
	})

	Describe("JSON decode", func() {
		DescribeTable("valid",
			func(s string, x Version) {
				Expect(json.Unmarshal([]byte(s), &v)).To(Succeed())
				Expect(v).To(Equal(x))
			},
			Entry("sample", `"1.2.3"`, MakeVersion(1, 2, 3)),
			Entry("min", `"0.0.0"`, MakeVersion(0, 0, 0)),
			Entry("max", `"255.255.255"`, MakeVersion(255, 255, 255)),
		)
		Describe("error", func() {
			Context("invalid text", func() {
				It("should not decode", func() {
					b := []byte(`"1"`)
					e := fmt.Sprintf("Invalid text for %T: %q", &v, "1")
					Expect(json.Unmarshal(b, &v).Error()).To(Equal(e))
				})
			})
			Context("invalid major", func() {
				It("should not decode", func() {
					b := []byte(`"256.1.2"`)
					e := fmt.Sprintf("Invalid major version: %s", "256")
					Expect(json.Unmarshal(b, &v).Error()).To(Equal(e))
				})
			})
			Context("invalid minor", func() {
				It("should not decode", func() {
					b := []byte(`"1.256.2"`)
					e := fmt.Sprintf("Invalid minor version: %s", "256")
					Expect(json.Unmarshal(b, &v).Error()).To(Equal(e))
				})
			})
			Context("invalid patch", func() {
				It("should not decode", func() {
					b := []byte(`"1.2.256"`)
					e := fmt.Sprintf("Invalid patch version: %s", "256")
					Expect(json.Unmarshal(b, &v).Error()).To(Equal(e))
				})
			})
		})
	})
})
