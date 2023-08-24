package util_test

import (
	"math"
	"time"

	. "github.com/hanindo/util/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Json Enc", func() {
	DescribeTable("String",
		func(v interface{}, x string) {
			Expect(JsonEnc{v}.String()).To(Equal(x))
		},
		Entry("bool", true, "true"),
		Entry("number", 12.3, "12.3"),
		Entry("string", "abc", `"abc"`),
		Entry("array", []int{1, 2, 3}, "[1,2,3]"),
		Entry("object", struct {
			B bool    `json:"b"`
			N float64 `json:"n"`
			S string  `json:"s"`
			A []int   `json:"a"`
		}{
			B: false,
			N: 234.5678,
			S: "humpty dumpty",
			A: []int{3, 4, 5},
		}, `{"b":false,"n":234.5678,"s":"humpty dumpty","a":[3,4,5]}`),
	)

	Describe("String", func() {
		Context("error", func() {
			It("should panic", func() {
				Expect(func() {
					_ = JsonEnc{math.NaN()}.String()
				}).Should(Panic())
			})
		})
	})
})

var _ = DescribeTable("IsAcii",
	func(b []byte, x bool) {
		Expect(IsAscii(b)).To(Equal(x))
	},
	Entry("one", []byte("a"), true),
	Entry("two", []byte(" ~"), true),
	Entry("many", []byte(" !\"#$%&'()*+,-./0123456789:;<=>?@"+
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`"+
		"abcdefghijklmnopqrstuvwxyz{|}~"), true),
	Entry("false lower", []byte("abcd\x1F"), false),
	Entry("false upper", []byte("abcd\x7F"), false),
)

var _ = DescribeTable("IsHex",
	func(b []byte, x bool) {
		Expect(IsHex(b)).To(Equal(x))
	},
	Entry("one", []byte("0"), true),
	Entry("two", []byte("9F"), true),
	Entry("many", []byte("01234567890ABCDEF"), true),
	Entry("false", []byte("01234567890ABCDEFG"), false),
)

var _ = DescribeTable("BitString",
	func(b []byte, x string) {
		Expect(BitString(b)).To(Equal(x))
	},
	Entry("one", []byte{0xA5}, "1010 0101"),
	Entry("two", []byte{0x5A, 0xC3}, "1100 0011 0101 1010"),
	Entry("three", []byte{0x12, 0x34, 0x56}, "0101 0110 0011 0100 0001 0010"),
	Entry("four", []byte{0xDE, 0xAD, 0xBE, 0xEF},
		"1110 1111 1011 1110 1010 1101 1101 1110"),
)

var _ = DescribeTable("FancyHex",
	func(b []byte, x string) {
		Expect(FancyHex(b)).To(Equal(x))
	},
	Entry("one", []byte{0x00}, ".."),
	Entry("two", []byte{0x01, 0x20}, ".1 2."),
	Entry("three", []byte{0x03, 0x40, 0xA5}, ".3 4. a5"),
	Entry("four", []byte{0x00, 0x00, 0x01, 0x9F}, ".. .. .1 9f"),
)

var _ = DescribeTable("HexWrap",
	func(b []byte, x string) {
		Expect(HexWrap(b, "raw tx: [", "]", "  ", 54, 80)).To(Equal(x))
	},
	Entry("in first",
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		"raw tx: [.. .1 .2 .3 .4 .5 .6 .7 .8 .9 .a .b .c .d .e]"),
	Entry("one next",
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
			14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		`raw tx: [
  .. .1 .2 .3 .4 .5 .6 .7 .8 .9 .a .b .c .d .e .f 1. 11 12 13 14 15 16 17 18 19
]`),
	Entry("two next",
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
		`raw tx: [
  .. .1 .2 .3 .4 .5 .6 .7 .8 .9 .a .b .c .d .e .f 1. 11 12 13 14 15 16 17 18 19
  1a 1b 1c 1d
]`),
)

var _ = DescribeTable("HexWrapper",
	func(b []byte, x string) {
		Expect(HexWrapper(b, "raw tx: [", "]", "  ", 54, 80).String()).
			To(Equal(x))
	},
	Entry("in first",
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		"raw tx: [.. .1 .2 .3 .4 .5 .6 .7 .8 .9 .a .b .c .d .e]"),
	Entry("one next",
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
			14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		`raw tx: [
  .. .1 .2 .3 .4 .5 .6 .7 .8 .9 .a .b .c .d .e .f 1. 11 12 13 14 15 16 17 18 19
]`),
	Entry("two next",
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
			20, 21, 22, 23, 24, 25, 26, 27, 28, 29},
		`raw tx: [
  .. .1 .2 .3 .4 .5 .6 .7 .8 .9 .a .b .c .d .e .f 1. 11 12 13 14 15 16 17 18 19
  1a 1b 1c 1d
]`),
)

var _ = DescribeTable("Indent",
	func(i, s, x string) {
		Expect(Indent(s, i)).To(Equal(x))
	},
	Entry("one", " ", "one", " one"),
	Entry("two", "  ",
		"satu\ndua",
		"  satu\n  dua"),
	Entry("three", "   ",
		"satu\ndua\n\ttiga",
		"   satu\n   dua\n   \ttiga"),
)

var _ = DescribeTable("TrimZero",
	func(s, x string) {
		Expect(TrimZero(s)).To(Equal(x))
	},
	Entry("one",
		"a 12.0 b 34.10 c 567.0 d 890.10 e",
		"a 12 b 34.1 c 567 d 890.1 e"),
	Entry("two",
		"a 12.00 b 34.5900 c 567.00 d 890.1200 e",
		"a 12 b 34.59 c 567 d 890.12 e"),
	Entry("three",
		"a 12.000 b 34.567000 c 567.000 d 890.12000 e",
		"a 12 b 34.567 c 567 d 890.12 e"),
	Entry("dot",
		"a 12.0 b 34.12 c 567.0 d 890.12 e",
		"a 12 b 34.12 c 567 d 890.12 e"),
	Entry("end",
		"a 12.34 b 34.5900 c 567.89 d 890.1200 e",
		"a 12.34 b 34.59 c 567.89 d 890.12 e"),
	Entry("none",
		"a 12.1 b 34.56 c 567.891 d 2345.6789 e",
		"a 12.1 b 34.56 c 567.891 d 2345.6789 e"),
)

var _ = DescribeTable("SpaceZero",
	func(s, x string) {
		Expect(SpaceZero(s)).To(Equal(x))
	},
	Entry("one",
		"a 12.0 b 34.10 c 567.0 d 890.10 e",
		"a 12   b 34.1  c 567   d 890.1  e"),
	Entry("two",
		"a 12.00 b 34.5900 c 567.00 d 890.1200 e",
		"a 12    b 34.59   c 567    d 890.12   e"),
	Entry("three",
		"a 12.000 b 34.567000 c 567.000 d 890.1200 e",
		"a 12     b 34.567    c 567     d 890.12   e"),
	Entry("dot",
		"a 12.0 b 34.12 c 567.0 d 890.12 e",
		"a 12   b 34.12 c 567   d 890.12 e"),
	Entry("end",
		"a 12.34 b 34.5900 c 567.89 d 890.1200 e",
		"a 12.34 b 34.59   c 567.89 d 890.12   e"),
	Entry("none",
		"a 12.1 b 34.56 c 567.891 d 2345.6789 e",
		"a 12.1 b 34.56 c 567.891 d 2345.6789 e"),
)

var _ = Describe("Comma", func() {
	It("should put comma in the right places", func() {
		//      0123456789
		//          ^  ^
		str := "x:   16500"
		Expect(Comma(str, []int{4, 7})).To(Equal("x:    16,500"))
	})
})

var _ = DescribeTable("Round",
	func(n float64, d int, x float64) {
		Expect(Round(n, d)).To(Equal(x))
	},
	Entry("lower 0", 1.45, 0, 1.0),
	Entry("half 0", 1.5, 0, 2.0),
	Entry("lower 1", 1.54, 1, 1.5),
	Entry("half 1", 1.25, 1, 1.3),
	Entry("lower 2", 1.654, 2, 1.65),
	Entry("half 2", 1.225, 2, 1.23),
)

var _ = DescribeTable("TruncDate",
	func(t, x time.Time) {
		Expect(TruncDate(t)).To(Equal(x))
	},
	Entry("local",
		time.Date(2021, time.March, 4, 5, 6, 7, 8, time.Local),
		time.Date(2021, time.March, 4, 0, 0, 0, 0, time.Local)),
	Entry("UTC",
		time.Date(2000, time.January, 2, 13, 14, 15, 16, time.UTC),
		time.Date(2000, time.January, 2, 0, 0, 0, 0, time.UTC)),
)
