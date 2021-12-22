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
	for i := range buff {
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

// HexWrap returns that FancyHex that wrap.
func HexWrap(b []byte, prefix, postfix, indent string, first, next int) string {
	if len(prefix)+len(postfix)+len(b)*3-1 <= first {
		return prefix + FancyHex(b) + postfix
	}

	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteByte('\n')
	l := (next - len(indent) + 1) / 3
	for i, e := 0, 0; i < len(b); i = e {
		e = i + l
		if e > len(b) {
			e = len(b)
		}
		sb.WriteString(indent)
		sb.WriteString(FancyHex(b[i:e]))
		sb.WriteByte('\n')
	}
	sb.WriteString(postfix)
	return sb.String()
}

// HexWrapper returns a struct that will HexWrap() on String(),
// handy for optional logging.
func HexWrapper(b []byte, pre, post, ind string, first, next int) hexWrapper {
	return hexWrapper{b, pre, post, ind, first, next}
}

type hexWrapper struct {
	b     []byte
	pre   string
	post  string
	ind   string
	first int
	next  int
}

func (h hexWrapper) String() string {
	return HexWrap(h.b, h.pre, h.post, h.ind, h.first, h.next)
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

// Comma put comma or space in list of indexes.
func Comma(s string, idxs []int) string {
	var b strings.Builder
	b.Grow(len(s) + len(idxs))
	var last int
	for _, idx := range idxs {
		b.WriteString(s[last:idx])
		if s[idx-1] != ' ' && s[idx] != ' ' {
			b.WriteByte(',')
		} else {
			b.WriteByte(' ')
		}
		last = idx
	}
	b.WriteString(s[last:])
	return b.String()
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
