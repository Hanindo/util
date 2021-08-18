package util_test

import (
	"fmt"
	. "github.com/hanindo/util/v2"
)

func ExampleJsonEnc() {
	list := []int{1, 2, 3}
	rec := struct {
		No   int    `json:"no"`
		Name string `json:"name"`
	}{
		No:   1,
		Name: "Jack",
	}
	fmt.Println(JsonEnc{list})
	fmt.Println(JsonEnc{rec})
	// Output:
	// [1,2,3]
	// {"no":1,"name":"Jack"}
}

func ExampleBitString() {
	fmt.Println(BitString([]byte{0xA5}))
	fmt.Println(BitString([]byte{0x5A, 0xC3}))
	fmt.Println(BitString([]byte{0xDE, 0xAD, 0xBE, 0xEF}))
	// Output:
	// 1010 0101
	// 1100 0011 0101 1010
	// 1110 1111 1011 1110 1010 1101 1101 1110
}

func ExampleFancyHex() {
	fmt.Println(FancyHex([]byte{0x0A}))
	fmt.Println(FancyHex([]byte{0, 0, 0x01, 0x23}))
	// Output:
	// .a
	// .. .. .1 23
}

func ExampleIndent() {
	fmt.Println(Indent(`Aku
Kalau sampai waktuku
'Ku mau tak seorang 'kan merayu
Tidak juga kau`,
		"... "))
	// Output:
	// ... Aku
	// ... Kalau sampai waktuku
	// ... 'Ku mau tak seorang 'kan merayu
	// ... Tidak juga kau
}

func ExampleRound() {
	fmt.Println(Round(1.5, 0))
	fmt.Println(Round(1.23, 1))
	fmt.Println(Round(12.3456, 2))
	// Output:
	// 2
	// 1.2
	// 12.35
}

func ExampleSpaceZero() {
	fmt.Println(SpaceZero("| 12.00 34.50 6.789 |"))
	fmt.Println(SpaceZero("|  3.40  1.00 9.100 |"))
	// Output:
	// | 12    34.5  6.789 |
	// |  3.4   1    9.1   |
}

func ExampleTrimZero() {
	fmt.Println(TrimZero("a:12.00 b:34.5000 c:6.780"))
	// Output: a:12 b:34.5 c:6.78
}
