// Package util provides commonly used utily functions and types.
package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
)

// JsonEnc is an utility type to easily print the JSON encoded value.
type JsonEnc struct {
	V interface{}
}

// String returns JSON formatted content of V field.
// It will panic on json.Marshal error.
func (e JsonEnc) String() string {
	b, err := json.Marshal(e.V)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

// IsAscii checks whether all b contents are ASCII printable characters.
func IsAscii(b []byte) bool {
	for _, c := range b {
		if c < 0x20 || c > 0x7E {
			return false
		}
	}
	return true
}

// IsHex checks whether all b contents are upper case hexadecimal numbers.
func IsHex(b []byte) bool {
	for _, c := range b {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// BitString returns bits formatted content of b bytes,
// the bits will be grouped by 4 bits for legibility.
// The string will be formatted Little Endian which means the first byte
// will be the last.
func BitString(buff []byte) string {
	var s string
	for i, _ := range buff {
		b := buff[len(buff)-i-1]
		if i > 0 {
			s += " "
		}
		s += fmt.Sprintf("%04b %04b", (b&0xF0)>>4, b&0x0F)
	}
	return s
}

// FancyHex return hex formatted content of b bytes.
// To increase legibility, all ``0'' will be replaced by ``.''
// and the bytes will be space separated.
func FancyHex(b []byte) string {
	return strings.ReplaceAll(fmt.Sprintf("% x", string(b)), "0", ".")
}

// Indent prepend i to all s lines.
func Indent(s, i string) string {
	return i + strings.ReplaceAll(s, "\n", "\n"+i)
}

var dot0RE = regexp.MustCompile(`\.0+\b`)
var end0RE = regexp.MustCompile(`(\.\d*[1-9])0+\b`)

// TrimZero removes all trailing fraction zeroes.
func TrimZero(s string) string {
	return end0RE.ReplaceAllString(dot0RE.ReplaceAllString(s, ""), "$1")
}

// SpaceZero replace all trailing fraction zeroes with spaces.
func SpaceZero(s string) string {
	str := dot0RE.ReplaceAllStringFunc(s, func(m string) string {
		return strings.Repeat(" ", len(m))
	})

	if end0RE.MatchString(str) {
		b := make([]byte, 0, len(str))
		var i int
		for _, sm := range end0RE.FindAllStringSubmatchIndex(str, -1) {
			b = append(b, str[i:sm[3]]...)
			b = append(b, strings.Repeat(" ", sm[1]-sm[3])...)
			i = sm[1]
		}
		b = append(b, str[i:]...)
		return string(b)
	} else {
		return str
	}
}

// Round x number to d decimal points.
func Round(x float64, d int) float64 {
	return math.Round(x*math.Pow10(d)) / math.Pow10(d)
}

// TruncDate truncate t time into whole day.
func TruncDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// FileNotExist checks whether the path is not exist.
func FileNotExist(name string) bool {
	_, err := os.Stat(name)
	return err != nil && errors.Is(err, os.ErrNotExist)
}
