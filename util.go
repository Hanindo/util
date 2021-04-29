package util

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	SHORT_TIME = "15:04:05"
)

var (
	dotZeroRE = regexp.MustCompile(`\.0+\b`)
)

type JsonEnc struct {
	V interface{}
}

func (e JsonEnc) String() string {
	b, err := json.Marshal(e.V)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

func IsHex(b []byte) bool {
	for _, c := range b {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func IsAscii(b []byte) bool {
	for _, c := range b {
		if c < 0x20 || c > 0x7f {
			return false
		}
	}
	return true
}

func FancyHex(b []byte) string {
	return strings.ReplaceAll(fmt.Sprintf("% x", string(b)), "0", ".")
}

func Indent(str, indent string) string {
	return strings.ReplaceAll(str, "\n", "\n"+indent)
}

func TrimZero(str string) string {
	return dotZeroRE.ReplaceAllString(str, "")
}

func SpaceZero(str string) string {
	return dotZeroRE.ReplaceAllStringFunc(str, func(s string) string {
		return strings.Repeat(" ", len(s))
	})
}

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

func FormatSec(end int) string {
	var zero time.Time
	return zero.Add(time.Duration(end) * time.Second).Format(SHORT_TIME)
}

func TruncDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func FileNotExist(name string) bool {
	_, err := os.Stat(name)
	return err != nil && os.IsNotExist(err)
}
